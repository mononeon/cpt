package communicate

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

func CommunicateTest(problemName, interactorFile string, testID string, timeLimit time.Duration, verbose bool) {
	fmt.Println("\033[36mStarting communication (run-twice) test...\033[0m")

	solSource := runner.FindSourceFile(problemName)
	if solSource == "" {
		fmt.Printf("\033[1;31mCould not find source file for %s\033[0m\n", problemName)
		return
	}
	solExec := runner.CompileSource(solSource, false)
	interExec := runner.CompileSource(interactorFile, true)

	cptDir := filepath.Join(".cpt", problemName)
	inFile := filepath.Join(cptDir, fmt.Sprintf("in_%s.txt", testID))
	outFile := filepath.Join(cptDir, fmt.Sprintf("out_%s.txt", testID))
	ansFile := filepath.Join(cptDir, fmt.Sprintf("ans_%s.txt", testID))

	if _, err := os.Stat(ansFile); os.IsNotExist(err) {
		os.WriteFile(ansFile, []byte("dummy\n"), 0644)
	}

	interCmd := runner.GetRunCmd(interactorFile, interExec)
	interCmd.Args = append(interCmd.Args, inFile, outFile, ansFile)

	interStdin, _ := interCmd.StdinPipe()
	interRawStdout, _ := interCmd.StdoutPipe()
	interCmd.Stderr = os.Stderr

	if err := interCmd.Start(); err != nil {
		fmt.Printf("\033[1;31mFailed to start interactor: %v\033[0m\n", err)
		return
	}

	// Create a persistent pipe for the solution's stdin.
	// This ensures that unread data (like Phase 2's input) is preserved in the pipe buffer
	// when Phase 1 exits, preventing data loss!
	solReader, solWriter, _ := os.Pipe()

	// Goroutine: Interactor -> solWriter
	go func() {
		var writer io.Writer
		if verbose {
			writer = logger.NewPrefixWriter("\033[36m[JURY]> \033[0m", solWriter)
		} else {
			writer = solWriter
		}

		buf := make([]byte, 1024)
		for {
			n, err := interRawStdout.Read(buf)
			if n > 0 {
				writer.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		solWriter.Close()
	}()

	runPhase := func(phase int) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
		defer cancel()

		solCmd := runner.GetRunCmd(solSource, solExec)
		solCmdWithCtx := exec.CommandContext(ctx, solCmd.Args[0], solCmd.Args[1:]...)

		// Use the persistent reader so the OS buffer is shared between phases
		solCmdWithCtx.Stdin = solReader
		solRawStdout, _ := solCmdWithCtx.StdoutPipe()
		solCmdWithCtx.Stderr = os.Stderr

		if err := solCmdWithCtx.Start(); err != nil {
			return fmt.Errorf("Failed to start solution phase %d: %v", phase, err)
		}

		// Goroutine: Solution -> Interactor
		go func() {
			var writer io.Writer
			if verbose {
				writer = logger.NewPrefixWriter("\033[32m[USER]> \033[0m", interStdin)
			} else {
				writer = interStdin
			}

			buf := make([]byte, 1024)
			for {
				n, err := solRawStdout.Read(buf)
				if n > 0 {
					writer.Write(buf[:n])
				}
				if err != nil {
					break
				}
			}
		}()

		err := solCmdWithCtx.Wait()

		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("Time/Idleness Limit Exceeded")
		}
		return err
	}

	fmt.Println("\033[33m--- Running Phase 1 (Encode) ---\033[0m")
	err1 := runPhase(1)
	if err1 != nil {
		fmt.Printf("\033[1;31mPhase 1 failed: %v\033[0m\n", err1)
	} else {
		fmt.Println("\033[32mPhase 1 completed.\033[0m")

		fmt.Println("\033[33m--- Running Phase 2 (Decode) ---\033[0m")
		err2 := runPhase(2)
		if err2 != nil {
			fmt.Printf("\033[1;31mPhase 2 failed: %v\033[0m\n", err2)
		} else {
			fmt.Println("\033[32mPhase 2 completed.\033[0m")
		}
	}

	// Close interStdin to signal the interactor that no more solutions will connect
	interStdin.Close()

	interErr := interCmd.Wait()
	exitCode := 0
	if interErr != nil {
		if exitError, ok := interErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	fmt.Println("\n\033[36mCommunication test finished.\033[0m")
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
		if outBytes, err := os.ReadFile(outFile); err == nil && len(outBytes) > 0 {
			fmt.Printf("Interactor Output:\n%s\n", string(outBytes))
		}
	}
}
