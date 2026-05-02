package runner

import (
	"bytes"
	"context"
	"cpt/internal/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"cpt/internal/testlib"
)

func FindSourceFile(problemName string) string {
	exts := []string{"cpp", "py", "c", "rs", "java", "go"}
	for _, ext := range exts {
		file := fmt.Sprintf("%s.%s", problemName, ext)
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

func getExtension(sourceFile string) string {
	parts := strings.Split(sourceFile, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

func CompileSource(sourceFile string, isTestingCode bool) string {
	ext := getExtension(sourceFile)
	cmdTemplate, ok := config.GlobalConfig.CompileCommands[ext]
	if !ok || cmdTemplate == "" {
		return "" // No compilation needed
	}

	if isTestingCode && ext == "cpp" {
		includeDir, err := testlib.EnsureTestlib()
		if err == nil && includeDir != "" {
			// Append the include directory flag to the compile command
			cmdTemplate = strings.Replace(cmdTemplate, "{source}", fmt.Sprintf("-I%s {source}", includeDir), 1)
		}
	}

	baseName := strings.TrimSuffix(sourceFile, "."+ext)
	outputName := baseName
	if runtime.GOOS == "windows" {
		outputName += ".exe"
	}

	cmdStr := strings.ReplaceAll(cmdTemplate, "{source}", sourceFile)
	cmdStr = strings.ReplaceAll(cmdStr, "{output}", outputName)

	fmt.Printf("\033[2mCompiling: %s\033[0m\n", cmdStr)

	// Since we might have complex shell commands, let's use the shell
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\033[1;31mCompilation Failed\033[0m\n")
		fmt.Println(string(output))
		os.Exit(1)
	}

	return outputName
}

func GetRunCmd(sourceFile, executable string) *exec.Cmd {
	ext := getExtension(sourceFile)
	cmdTemplate, ok := config.GlobalConfig.RunCommands[ext]
	if !ok {
		fmt.Printf("\033[1;31mNo run command configured for extension '%s'\033[0m\n", ext)
		os.Exit(1)
	}

	if executable == "" {
		executable = sourceFile
	}

	if runtime.GOOS == "windows" && strings.HasSuffix(executable, ".exe") && filepath.Dir(executable) == "." {
		executable = ".\\" + executable
	}

	cmdStr := strings.ReplaceAll(cmdTemplate, "{source}", sourceFile)
	cmdStr = strings.ReplaceAll(cmdStr, "{executable}", executable)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	return cmd
}

func RunWithIO(cmd *exec.Cmd, inputData string, timeout time.Duration) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdWithCtx := exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)
	cmdWithCtx.Stdin = strings.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmdWithCtx.Stdout = &stdout
	cmdWithCtx.Stderr = &stderr

	err := cmdWithCtx.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", "", fmt.Errorf("Time Limit Exceeded (%v)", timeout)
	}

	if err != nil {
		return "", stderr.String(), fmt.Errorf("Runtime Error (Exit Code %v):\n%s", err, stderr.String())
	}

	return stdout.String(), stderr.String(), nil
}
