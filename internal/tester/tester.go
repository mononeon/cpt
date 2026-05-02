package tester

import (
	"cpt/internal/diff"
	"cpt/internal/runner"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func RunTests(problemName string, specificTests []string, useStdin bool, checkerFile, validatorFile string) {
	sourceFile := runner.FindSourceFile(problemName)
	if sourceFile == "" {
		fmt.Printf("\033[1;31mCould not find source file for %s\033[0m\n", problemName)
		return
	}

	executable := runner.CompileSource(sourceFile, false)
	cmd := runner.GetRunCmd(sourceFile, executable)

	var checkerExec, validatorExec string
	if checkerFile != "" {
		checkerExec = runner.CompileSource(checkerFile, true)
	}
	if validatorFile != "" {
		validatorExec = runner.CompileSource(validatorFile, true)
	}

	if useStdin {
		fmt.Printf("\033[1;36mRunning %s with provided input...\033[0m\n", sourceFile)
		fmt.Println("Enter input (Ctrl+Z and Enter on Windows, Ctrl+D on Unix to EOF):")
		inputBytes, _ := io.ReadAll(os.Stdin)
		out, stderrRaw, err := runner.RunWithIO(cmd, string(inputBytes), 10*time.Second)
		if stderrRaw != "" {
			fmt.Printf("Standard Error:\n%s\n", strings.TrimSpace(stderrRaw))
		}
		if err != nil {
			fmt.Printf("\033[31m%v\033[0m\n", err)
		} else {
			fmt.Println("\033[32mOutput:\033[0m")
			fmt.Println(out)
		}
		return
	}

	cptDir := filepath.Join(".cpt", problemName)
	if _, err := os.Stat(cptDir); os.IsNotExist(err) {
		fmt.Printf("\033[33mNo test cases found in %s\033[0m\n", cptDir)
		return
	}

	files, err := os.ReadDir(cptDir)
	if err != nil {
		fmt.Printf("Error reading test directory: %v\n", err)
		return
	}

	var testIDs []string
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "in_") && strings.HasSuffix(name, ".txt") {
			id := strings.TrimSuffix(strings.TrimPrefix(name, "in_"), ".txt")
			if len(specificTests) > 0 {
				found := false
				for _, t := range specificTests {
					if t == id {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			testIDs = append(testIDs, id)
		}
	}

	if len(testIDs) == 0 {
		fmt.Println("\033[33mNo specific test cases found to run.\033[0m")
		return
	}

	sort.Strings(testIDs)

	passed := 0
	total := len(testIDs)

	for _, id := range testIDs {
		inFile := filepath.Join(cptDir, fmt.Sprintf("in_%s.txt", id))
		outFile := filepath.Join(cptDir, fmt.Sprintf("out_%s.txt", id))

		inData, _ := os.ReadFile(inFile)
		expectedData, _ := os.ReadFile(outFile)

		// Validation
		if validatorFile != "" {
			valCmd := runner.GetRunCmd(validatorFile, validatorExec)
			_, _, err := runner.RunWithIO(valCmd, string(inData), 5*time.Second)
			if err != nil {
				fmt.Printf("\n\033[36m^ TC %s\033[0m \033[1;31mValidator FAIL\033[0m\n", id)
				fmt.Println(err)
				continue
			}
		}

		expected := strings.TrimSpace(string(expectedData))

		startTime := time.Now()
		runCmd := runner.GetRunCmd(sourceFile, executable)
		actualRaw, stderrRaw, err := runner.RunWithIO(runCmd, string(inData), 5*time.Second)
		duration := time.Since(startTime).Milliseconds()

		actual := strings.TrimSpace(actualRaw)

		fmt.Printf("\n\033[36m^ TC %s\033[0m ", id)

		if stderrRaw != "" {
			fmt.Printf("Standard Error:\n%s\n", strings.TrimSpace(stderrRaw))
		}

		if err != nil {
			if strings.Contains(err.Error(), "Time Limit Exceeded") {
				fmt.Printf("\033[1;33mTLE\033[0m \033[2m%dms\033[0m\n", duration)
			} else {
				// Treat crashes/memory issues as RE
				fmt.Printf("\033[1;31mRE (Runtime Error)\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Printf("\033[31m%v\033[0m\n", err)
			}
			continue
		}

		if checkerFile != "" {
			// Write actual to a temp file for testlib
			actualFile := filepath.Join(cptDir, fmt.Sprintf("actual_%s.txt", id))
			os.WriteFile(actualFile, []byte(actualRaw), 0644)
			
			// checker <input> <output> <answer>
			checkCmd := runner.GetRunCmd(checkerFile, checkerExec)
			checkCmd.Args = append(checkCmd.Args, inFile, actualFile, outFile)
			
			checkerOut, checkErr := checkCmd.CombinedOutput()
			
			exitCode := 0
			if checkErr != nil {
				if exitError, ok := checkErr.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					exitCode = -1
				}
			}

			// testlib exit codes: 0=OK, 1=WA, 2=PE, 3=FAIL
			switch exitCode {
			case 0:
				passed++
				fmt.Printf("\033[1;32mAC (Accepted)\033[0m \033[2m%dms\033[0m\n", duration)
			case 1:
				fmt.Printf("\033[1;31mWA (Wrong Answer)\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Printf("Checker: %s\n", string(checkerOut))
			case 2, 7, 8:
				fmt.Printf("\033[1;33mPE (Presentation Error)\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Printf("Checker: %s\n", string(checkerOut))
			default:
				fmt.Printf("\033[1;31mFAIL (Checker Error)\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Printf("Checker: %s\n", string(checkerOut))
			}

		} else {
			// Default Exact/LCS Diff
			if expected == "" {
				fmt.Printf("\033[33mNo Expected Output\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Println("Received Output:")
				fmt.Println(actualRaw)
			} else if diff.DiffWords(expected, actual) == expected || actual == expected {
				passed++
				fmt.Printf("\033[1;32mAC (Accepted)\033[0m \033[2m%dms\033[0m\n", duration)
			} else {
				fmt.Printf("\033[1;31mWA (Wrong Answer)\033[0m \033[2m%dms\033[0m\n", duration)
				fmt.Println("Input:")
				fmt.Println(string(inData))
				fmt.Println("Expected:")
				fmt.Println(expected)
				fmt.Println("Received:")
				fmt.Println(actual)
				fmt.Println("Diff:")
				fmt.Println(diff.DiffWords(expected, actual))
			}
		}
	}

	fmt.Printf("\n\033[1mSummary:\033[0m %d/%d passed.\n", passed, total)
}

func AddTest(problemName string) {
	cptDir := filepath.Join(".cpt", problemName)
	os.MkdirAll(cptDir, 0755)

	files, _ := os.ReadDir(cptDir)
	maxID := 0
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "in_") {
			var id int
			fmt.Sscanf(name, "in_%d.txt", &id)
			if id > maxID {
				maxID = id
			}
		}
	}
	nextID := maxID + 1

	fmt.Printf("\033[36mAdding Test %d for %s\033[0m\n", nextID, problemName)
	fmt.Println("Paste input (Ctrl+Z and Enter on Windows to EOF):")
	inputBytes, _ := io.ReadAll(os.Stdin)

	fmt.Println("Paste expected output (Ctrl+Z and Enter on Windows to EOF):")
	outputBytes, _ := io.ReadAll(os.Stdin)

	os.WriteFile(filepath.Join(cptDir, fmt.Sprintf("in_%d.txt", nextID)), inputBytes, 0644)
	os.WriteFile(filepath.Join(cptDir, fmt.Sprintf("out_%d.txt", nextID)), outputBytes, 0644)

	fmt.Printf("\033[32mTest %d added successfully.\033[0m\n", nextID)
}

func EditTest(problemName string, testID string) {
	cptDir := filepath.Join(".cpt", problemName)
	inFile := filepath.Join(cptDir, fmt.Sprintf("in_%s.txt", testID))
	outFile := filepath.Join(cptDir, fmt.Sprintf("out_%s.txt", testID))

	if _, err := os.Stat(inFile); os.IsNotExist(err) {
		fmt.Printf("\033[31mTest %s does not exist.\033[0m\n", testID)
		return
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad" // Fallback on Windows
	}

	cmd1 := exec.Command(editor, inFile)
	cmd1.Run()
	cmd2 := exec.Command(editor, outFile)
	cmd2.Run()

	fmt.Printf("\033[32mTest %s edited.\033[0m\n", testID)
}

func RmTest(problemName string, testID string) {
	cptDir := filepath.Join(".cpt", problemName)
	inFile := filepath.Join(cptDir, fmt.Sprintf("in_%s.txt", testID))
	outFile := filepath.Join(cptDir, fmt.Sprintf("out_%s.txt", testID))

	deleted := false
	if err := os.Remove(inFile); err == nil {
		deleted = true
	}
	if err := os.Remove(outFile); err == nil {
		deleted = true
	}

	if deleted {
		fmt.Printf("\033[32mTest %s removed.\033[0m\n", testID)
	} else {
		fmt.Printf("\033[31mTest %s not found.\033[0m\n", testID)
	}
}
