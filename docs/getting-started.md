# Getting Started with CPT

This guide provides detailed explanations and examples for all `cpt` commands to help you streamline your competitive programming workflow.

---

## 1. Automated Problem Parsing
Integrate with the [Competitive Companion](https://github.com/jmerle/competitive-companion) browser extension to automatically parse test cases.

```bash
# Start the local listener
$ cpt listen
```
While running, click the Competitive Companion icon on any supported platform (Codeforces, AtCoder, etc.). `cpt` will automatically create a `.cpt/<problem-id>` directory containing all inputs and expected outputs.

---

## 2. Basic Testing
Test your solution against the downloaded sample test cases. `cpt` will automatically compile your code if needed.

```bash
# Test a solution (e.g., A.cpp) against all test cases
$ cpt test A

# Test only specific cases
$ cpt test A --tests 1,3

# Use a custom testlib.h validator or checker
$ cpt test A --validator val.cpp --checker check.cpp
```

---

## 3. Stress Testing
Hunt down edge cases by running your solution against a brute-force approach using randomized inputs.

```bash
# Run 100 random iterations
$ cpt stress A --gen gen.cpp --brute brute.cpp --iters 100

# See detailed outputs for every single iteration
$ cpt stress A --gen gen.cpp --brute brute.cpp -v
```
If a mismatch or runtime error occurs, the failing input, expected output, and actual output will be saved in the `.cpt/A/` directory.

---

## 4. Interactive & Communication Problems
`cpt` supports real-time, terminal-based visual logging for interactive queries. It automatically pipes data between your solution and a local interactor.

### Standard Interactive Problems
For standard ping-pong style interactive problems:
```bash
$ cpt interact A --interactor interact.cpp -v
```

### Two-Stage Communication Problems
For Codeforces-style run-twice communication problems (Encode/Decode phases):
```bash
$ cpt communicate A --interactor interact.cpp -v
```

---

## 5. Test Case Management
Easily manage custom test cases from the CLI without navigating directories.

```bash
# Manually add a test case via a terminal prompt
$ cpt add A

# Open a specific test case in your default text editor
$ cpt edit A 2

# Remove a test case
$ cpt rm A 2
```

---

## 6. Boilerplate Generation
Instantly generate `testlib.h` boilerplates for writing custom testing tools:

```bash
# Generate a generator (gen.cpp)
$ cpt template gen

# Generate a checker (checker.cpp)
$ cpt template check

# Generate a validator (validator.cpp)
$ cpt template val

# Generate an interactor (interactor.cpp)
$ cpt template interact
```
