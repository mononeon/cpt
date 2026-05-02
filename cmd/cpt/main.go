package main

import (
	"cpt/internal/communicate"
	"cpt/internal/companion"
	"cpt/internal/config"
	"cpt/internal/interact"
	"cpt/internal/stress"
	"cpt/internal/templates"
	"cpt/internal/tester"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "listen":
		listenCmd := flag.NewFlagSet("listen", flag.ExitOnError)
		port := listenCmd.Int("port", 10043, "Port to listen on")
		format := listenCmd.String("format", "", "Problem naming format (veryshortform, shortform, longform)")
		listenCmd.Parse(os.Args[2:])
		companion.Listen(*port, *format)

	case "test":
		testCmd := flag.NewFlagSet("test", flag.ExitOnError)
		tests := testCmd.String("tests", "", "Comma-separated list of test cases")
		stdin := testCmd.Bool("stdin", false, "Read input from stdin")
		checker := testCmd.String("checker", "", "Checker file (e.g. checker.cpp)")
		validator := testCmd.String("validator", "", "Validator file (e.g. validator.cpp)")
		testCmd.Parse(os.Args[2:])

		args := testCmd.Args()
		if len(args) < 1 {
			fmt.Println("Usage: cpt test <problem> [--checker <file>] [--validator <file>]")
			os.Exit(1)
		}

		problem := args[0]
		var specificTests []string
		if *tests != "" {
			specificTests = append(specificTests, *tests)
		}

		tester.RunTests(problem, specificTests, *stdin, *checker, *validator)

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cpt add <problem>")
			os.Exit(1)
		}
		tester.AddTest(os.Args[2])

	case "edit":
		if len(os.Args) < 4 {
			fmt.Println("Usage: cpt edit <problem> <test_id>")
			os.Exit(1)
		}
		tester.EditTest(os.Args[2], os.Args[3])

	case "rm":
		if len(os.Args) < 4 {
			fmt.Println("Usage: cpt rm <problem> <test_id>")
			os.Exit(1)
		}
		tester.RmTest(os.Args[2], os.Args[3])

	case "stress":
		stressCmd := flag.NewFlagSet("stress", flag.ExitOnError)
		generator := stressCmd.String("gen", "", "Generator file")
		brute := stressCmd.String("brute", "", "Brute force solution")
		checker := stressCmd.String("checker", "", "Custom checker file")
		validator := stressCmd.String("validator", "", "Input validator file")
		iters := stressCmd.Int("iters", 100, "Number of stress test iterations")
		timeLimit := stressCmd.Duration("tl", 5*time.Second, "Time limit for each run")
		verbose := stressCmd.Bool("v", false, "Verbose mode to print inputs and outputs")
		stressCmd.Parse(os.Args[2:])

		args := stressCmd.Args()
		if len(args) < 1 || *generator == "" {
			fmt.Println("Usage: cpt stress [--gen <file>] [--brute <file>] [--checker <file>] [--validator <file>] [--iters N] [--tl duration] [-v] <problem>")
			os.Exit(1)
		}
		stress.StressTest(args[0], *generator, *brute, *checker, *validator, *iters, *timeLimit, *verbose)

	case "interact":
		interactCmd := flag.NewFlagSet("interact", flag.ExitOnError)
		interactor := interactCmd.String("interactor", "", "Interactor file")
		testID := interactCmd.String("test", "1", "Test ID to run against")
		timeLimit := interactCmd.Duration("tl", 10*time.Second, "Time limit for each run")
		verbose := interactCmd.Bool("v", false, "Verbose mode to show interaction")
		interactCmd.Parse(os.Args[2:])

		args := interactCmd.Args()
		if len(args) < 1 || *interactor == "" {
			fmt.Println("Usage: cpt interact [--interactor <file>] [--test id] [--tl duration] [-v] <problem>")
			os.Exit(1)
		}
		interact.InteractTest(args[0], *interactor, *testID, *timeLimit, *verbose)

	case "communicate":
		commCmd := flag.NewFlagSet("communicate", flag.ExitOnError)
		interactor := commCmd.String("interactor", "", "Interactor file")
		testID := commCmd.String("test", "1", "Test ID to run against")
		timeLimit := commCmd.Duration("tl", 10*time.Second, "Time limit for each run")
		verbose := commCmd.Bool("v", false, "Verbose mode to show interaction")
		commCmd.Parse(os.Args[2:])

		args := commCmd.Args()
		if len(args) < 1 || *interactor == "" {
			fmt.Println("Usage: cpt communicate [--interactor <file>] [--test id] [--tl duration] [-v] <problem>")
			os.Exit(1)
		}
		communicate.CommunicateTest(args[0], *interactor, *testID, *timeLimit, *verbose)

	case "template":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cpt template <gen|check|val|interact>")
			os.Exit(1)
		}
		templates.GenerateTemplate(os.Args[2])

	case "config":
		path, _ := config.GetConfigPath()
		fmt.Printf("Config file path: %s\n", path)
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "notepad" // Fallback on Windows
		}
		exec.Command(editor, path).Run()

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`Competitive Programming Tool (cpt)

Commands:
  listen     Listen for Competitive Companion
  test       Test a solution against cases
  add        Add a custom test case
  edit       Edit an existing test case
  rm         Remove a test case
  stress     Stress test a solution
  interact   Test an interactive problem
  communicate Test a two-stage communication problem
  template   Generate a testing boilerplate (gen|check|val|interact)
  config     Open the configuration file`)
}
