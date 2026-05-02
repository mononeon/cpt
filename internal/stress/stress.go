package stress

import (
	"cpt/internal/runner"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func StressTest(problemName, generatorFile, bruteFile, checkerFile, validatorFile string, iters int, timeLimit time.Duration, verbose bool) {
	fmt.Println("\033[36mStarting stress testing...\033[0m")

	if bruteFile == "" && checkerFile == "" {
		fmt.Println("\033[1;31mError: You must provide either a brute force solution or a checker!\033[0m")
		return
	}

	solSource := runner.FindSourceFile(problemName)
	if solSource == "" {
		fmt.Printf("\033[1;31mCould not find source file for %s\033[0m\n", problemName)
		return
	}
	solExec := runner.CompileSource(solSource, false)

	var bruteExec, checkerExec, validatorExec string
	if bruteFile != "" {
		bruteExec = runner.CompileSource(bruteFile, false)
	}
	if checkerFile != "" {
		checkerExec = runner.CompileSource(checkerFile, true)
	}
	if validatorFile != "" {
		validatorExec = runner.CompileSource(validatorFile, true)
	}

	var genExec string
	if strings.HasSuffix(generatorFile, ".cpp") {
		genExec = runner.CompileSource(generatorFile, true)
	}

	for i := 1; i <= iters; i++ {
		genCmd := runner.GetRunCmd(generatorFile, genExec)
		inpRaw, stderrRaw, err := runner.RunWithIO(genCmd, "", timeLimit)
		if stderrRaw != "" && verbose {
			fmt.Printf("Generator Stderr:\n%s\n", strings.TrimSpace(stderrRaw))
		}
		if err != nil {
			fmt.Printf("\033[1;31mGenerator error on iter %d: %v\033[0m\n", i, err)
			return
		}

		if validatorFile != "" {
			valCmd := runner.GetRunCmd(validatorFile, validatorExec)
			_, valStderr, err := runner.RunWithIO(valCmd, inpRaw, timeLimit)
			if valStderr != "" && verbose {
				fmt.Printf("Validator Stderr:\n%s\n", strings.TrimSpace(valStderr))
			}
			if err != nil {
				fmt.Printf("\033[1;31mValidator rejected generated input on iter %d: %v\033[0m\n", i, err)
				return
			}
		}

		solCmd := runner.GetRunCmd(solSource, solExec)
		actualRaw, solStderr, err := runner.RunWithIO(solCmd, inpRaw, timeLimit)
		if solStderr != "" && verbose {
			fmt.Printf("Solution Stderr:\n%s\n", strings.TrimSpace(solStderr))
		}
		if err != nil {
			if strings.Contains(err.Error(), "Time Limit Exceeded") {
				fmt.Printf("\033[1;33mSolution TLE on iter %d\033[0m\n", i)
			} else {
				fmt.Printf("\033[1;31mSolution RE on iter %d: %v\033[0m\n", i, err)
			}
			saveStressTest(problemName, inpRaw, "", actualRaw)
			return
		}

		expectedRaw := ""
		if bruteFile != "" {
			bruteCmd := runner.GetRunCmd(bruteFile, bruteExec)
			var bruteStderr string
			expectedRaw, bruteStderr, err = runner.RunWithIO(bruteCmd, inpRaw, timeLimit)
			if bruteStderr != "" && verbose {
				fmt.Printf("Brute Force Stderr:\n%s\n", strings.TrimSpace(bruteStderr))
			}
			if err != nil {
				if strings.Contains(err.Error(), "Time Limit Exceeded") {
					fmt.Printf("\033[1;33mBrute force TLE on iter %d\033[0m\n", i)
				} else {
					fmt.Printf("\033[1;31mBrute force RE on iter %d: %v\033[0m\n", i, err)
				}
				return
			}
		}

		if checkerFile != "" {
			cptDir := filepath.Join(".cpt", problemName)
			os.MkdirAll(cptDir, 0755)

			inFile := filepath.Join(cptDir, "stress_in.txt")
			actualFile := filepath.Join(cptDir, "stress_actual.txt")
			ansFile := filepath.Join(cptDir, "stress_out.txt")

			os.WriteFile(inFile, []byte(inpRaw), 0644)
			os.WriteFile(actualFile, []byte(actualRaw), 0644)
			os.WriteFile(ansFile, []byte(expectedRaw), 0644)

			checkCmd := runner.GetRunCmd(checkerFile, checkerExec)
			checkCmd.Args = append(checkCmd.Args, inFile, actualFile, ansFile)

			checkerOut, checkErr := checkCmd.CombinedOutput()
			exitCode := 0
			if checkErr != nil {
				if exitError, ok := checkErr.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					exitCode = -1
				}
			}

			if exitCode != 0 {
				switch exitCode {
				case 1:
					fmt.Printf("\n\033[1;31mWA (Wrong Answer) on iteration %d!\033[0m\n", i)
				case 2, 7, 8:
					fmt.Printf("\n\033[1;33mPE (Presentation Error) on iteration %d!\033[0m\n", i)
				default:
					fmt.Printf("\n\033[1;31mFAIL (Checker Error) on iteration %d!\033[0m\n", i)
				}
				fmt.Printf("Checker: %s\n", string(checkerOut))
				fmt.Println("\033[33mSaved to .cpt/stress_in.txt, stress_out.txt, stress_actual.txt\033[0m")
				return
			}

		} else {
			// Fallback to exact match
			expected := strings.TrimSpace(expectedRaw)
			actual := strings.TrimSpace(actualRaw)

			if expected != actual {
				fmt.Printf("\n\033[1;31mMismatch (WA) found on iteration %d!\033[0m\n", i)
				saveStressTest(problemName, inpRaw, expectedRaw, actualRaw)
				return
			}
		}

		if verbose {
			fmt.Printf("Iter %d AC (Accepted)\n", i)
			fmt.Printf("Input:\n%s\n", string(inpRaw))
			fmt.Printf("Output:\n%s\n", string(actualRaw))
			if bruteFile != "" {
				fmt.Printf("Expected:\n%s\n", string(expectedRaw))
			}
			fmt.Println()
		} else if i%10 == 0 || i == iters {
			fmt.Printf("Passed %d iterations...\n", i)
		}
	}

	fmt.Println("\033[1;32mStress testing completed successfully! No mismatches found.\033[0m")
}

func saveStressTest(problemName, inp, exp, act string) {
	cptDir := filepath.Join(".cpt", problemName)
	os.MkdirAll(cptDir, 0755)

	os.WriteFile(filepath.Join(cptDir, "stress_in.txt"), []byte(inp), 0644)
	os.WriteFile(filepath.Join(cptDir, "stress_out.txt"), []byte(exp), 0644)
	os.WriteFile(filepath.Join(cptDir, "stress_actual.txt"), []byte(act), 0644)

	fmt.Println("\033[33mSaved failing case to .cpt/stress_in.txt, stress_out.txt, stress_actual.txt\033[0m")
}
