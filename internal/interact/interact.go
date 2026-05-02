package interact

import (
	"context"
	"cpt/internal/logger"
	"cpt/internal/runner"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func InteractTest(problemName, interactorFile string, testID string, timeLimit time.Duration, verbose bool) {
	fmt.Println("\033[36mStarting interactive test...\033[0m")

	solSource := runner.FindSourceFile(problemName)
	if solSource == "" {
		fmt.Printf("\033[1;31mCould not find source file for %s\033[0m\n", problemName)
		return
	}
	solExec := runner.CompileSource(solSource, false)
	solCmd := runner.GetRunCmd(solSource, solExec)

	interExec := runner.CompileSource(interactorFile, true)
	interCmd := runner.GetRunCmd(interactorFile, interExec)

	cptDir := filepath.Join(".cpt", problemName)
	inFile := filepath.Join(cptDir, fmt.Sprintf("in_%s.txt", testID))
	ansFile := filepath.Join(cptDir, fmt.Sprintf("out_%s.txt", testID))
	outFile := filepath.Join(cptDir, fmt.Sprintf("interactor_out_%s.txt", testID))

	if _, err := os.Stat(inFile); os.IsNotExist(err) {
		fmt.Printf("\033[33mWarning: Input file %s not found. Using empty dummy file.\033[0m\n", inFile)
		os.MkdirAll(cptDir, 0755)
		os.WriteFile(inFile, []byte(""), 0644)
	}
	if _, err := os.Stat(ansFile); os.IsNotExist(err) {
		os.WriteFile(ansFile, []byte(""), 0644)
	}

	interCmd.Args = append(interCmd.Args, inFile, outFile, ansFile)

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	solCmdWithCtx := exec.CommandContext(ctx, solCmd.Args[0], solCmd.Args[1:]...)
	interCmdWithCtx := exec.CommandContext(ctx, interCmd.Args[0], interCmd.Args[1:]...)

	solStdout, _ := solCmdWithCtx.StdoutPipe()
	interStdin, _ := interCmdWithCtx.StdinPipe()
	interStdout, _ := interCmdWithCtx.StdoutPipe()
	solStdin, _ := solCmdWithCtx.StdinPipe()

	// Goroutine: Solution -> Interactor
	go func() {
		defer interStdin.Close()
		var writer io.Writer
		if verbose {
			writer = logger.NewPrefixWriter("\033[32m[USER]> \033[0m", interStdin)
		} else {
			writer = interStdin
		}

		buf := make([]byte, 1024)
		for {
			n, err := solStdout.Read(buf)
			if n > 0 {
				writer.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	// Goroutine: Interactor -> Solution
	go func() {
		defer solStdin.Close()
		var writer io.Writer
		if verbose {
			writer = logger.NewPrefixWriter("\033[36m[JURY]> \033[0m", solStdin)
		} else {
			writer = solStdin
		}

		buf := make([]byte, 1024)
		for {
			n, err := interStdout.Read(buf)
			if n > 0 {
				writer.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	err := solCmdWithCtx.Start()
	if err != nil {
		fmt.Println("Error starting solution:", err)
		return
	}

	err = interCmdWithCtx.Start()
	if err != nil {
		fmt.Println("Error starting interactor:", err)
		return
	}

	solCmdWithCtx.Wait()
	interErr := interCmdWithCtx.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("\n\033[1;33mTime Limit / Idleness Limit Exceeded\033[0m")
		return
	}

	exitCode := 0
	if interErr != nil {
		if exitError, ok := interErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	fmt.Println("\n\033[36mInteractive test finished.\033[0m")
	switch exitCode {
	case 0:
		fmt.Println("\033[1;32mAC (Accepted)\033[0m")
	case 1:
		fmt.Println("\033[1;31mWA (Wrong Answer)\033[0m")
	case 2, 7, 8:
		fmt.Println("\033[1;33mPE (Presentation Error)\033[0m")
	default:
		fmt.Println("\033[1;31mFAIL (Interactor Error)\033[0m")
	}

	if !verbose {
		if outData, err := os.ReadFile(outFile); err == nil && len(outData) > 0 {
			fmt.Printf("Interactor Output:\n%s\n", string(outData))
		}
	}
}
