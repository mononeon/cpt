# Configuration

`cpt` is highly customizable via a JSON configuration file. 

To open the configuration file in your default editor, run:
```bash
$ cpt config
```

By default, the configuration file is located at `~/.cptconfig.json` (in your user's home directory).

## Example Configuration

```json
{
    "compile_commands": {
        "cpp": "g++ -O2 -Wall -Wextra -std=c++17 {source} -o {output}",
        "py": "",
        "c": "gcc -O2 -Wall {source} -o {output}",
        "rs": "rustc {source} -o {output}",
        "java": "javac {source}",
        "go": "go build -o {output} {source}"
    },
    "run_commands": {
        "cpp": "{executable}",
        "py": "python {source}",
        "c": "{executable}",
        "rs": "{executable}",
        "java": "java {executable}",
        "go": "{executable}"
    }
}
```

## Variables
- **`{source}`**: The source file name (e.g., `A.cpp`).
- **`{output}`**: The desired executable name (e.g., `A.exe` on Windows or `A` on Unix).
- **`{executable}`**: The path to run the executable.

## Custom Languages
You can easily add support for any programming language by adding its file extension to both the `compile_commands` and `run_commands` objects.

If a language doesn't require compilation (like Python), simply leave its `compile_commands` value as an empty string `""`.
