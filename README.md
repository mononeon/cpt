# Competitive Programming Tool (cpt)

A fast, lightweight command-line utility written in Go designed to perfectly streamline your competitive programming workflow. Automate downloading test cases, executing solutions, stress testing, and debugging interactive problems entirely from your terminal.

## Features
- **Parse Problems:** Automatically download sample test cases via the Competitive Companion extension.
- **Run Tests:** Compile and run your code against downloaded or custom test cases.
- **Stress Test:** Find edge cases by stress testing against brute-force solutions with randomized inputs.
- **Interactive Support:** Test interactive and run-twice communication problems with real-time visual chat logs.
- **Testlib Integrations:** Built-in support for `testlib.h` validators, checkers, and interactors.
- **Manage Cases:** Quickly add, edit, or remove custom test cases via CLI commands.

For detailed documentation, read [docs/getting-started.md](docs/getting-started.md).

## Supported Platforms
`cpt` supports parsing from any online judge supported by [Competitive Companion](https://github.com/jmerle/competitive-companion), including **Codeforces**, **AtCoder**, **LeetCode**, **CodeChef**, and many more.

---

## How to Install

`cpt` is distributed as a lightning-fast, standalone binary. You do **not** need Go installed to use it.

1. Download the latest binary for your operating system (Windows, macOS, or Linux) from the [Releases](../../releases/latest) page.
2. If you are on macOS or Linux, make the file executable: `chmod +x cpt`
3. Rename the binary to `cpt` (or `cpt.exe` on Windows).
4. Add the binary to your system's `PATH`.

For detailed installation instructions and troubleshooting, read [docs/INSTALL.md](docs/INSTALL.md).

---

## How to Use

```bash
$ cpt listen                           # Listen for Competitive Companion payloads
$ cpt test A                           # Test your code against sample cases
$ cpt stress A --gen G --brute B       # Stress test your code to find failing inputs
$ cpt interact A --interactor I        # Test standard interactive problems
$ cpt communicate A --interactor I     # Test run-twice communication problems
$ cpt template [gen|check|val|interact]# Generate testlib boilerplate code
$ cpt add A                            # Manually add a custom test case
$ cpt config                           # Edit compiler/run flags and settings
```

For full details, flags, and advanced examples, see:
- [Getting Started & Commands](docs/getting-started.md)
- [Configuration Guide](docs/config.md)
