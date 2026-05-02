# Installation Instructions

`cpt` is a standalone executable written in Go. You do not need any runtime environments (like Python or Node.js) installed to use it.

## Windows

1. Navigate to the [Latest Release](../../releases/latest) page on GitHub.
2. Download the `cpt-windows-amd64.exe` file.
3. Rename the downloaded file to `cpt.exe`.
4. Create a permanent directory for the tool (e.g., `C:\Tools\cpt\`) and move `cpt.exe` into it.
5. **Add to System PATH**:
   - Open the Start Menu, type "Environment Variables", and press Enter.
   - Click the "Environment Variables..." button.
   - Under "System variables" or "User variables", find the `Path` variable and click "Edit".
   - Click "New" and add the path to the directory you created (e.g., `C:\Tools\cpt\`).
   - Click "OK" on all windows to save the changes.
6. Restart your terminal (Command Prompt, PowerShell, or VS Code terminal).
7. Verify the installation by running:
   ```bash
   $ cpt --help
   ```

## macOS & Linux

1. Navigate to the [Latest Release](../../releases/latest) page on GitHub.
2. Download the appropriate binary for your system:
   - macOS (Apple Silicon): `cpt-darwin-arm64`
   - macOS (Intel): `cpt-darwin-amd64`
   - Linux: `cpt-linux-amd64`
3. Rename the downloaded file to `cpt`.
4. Open a terminal and navigate to your download directory.
5. Make the binary executable:
   ```bash
   $ chmod +x cpt
   ```
6. Move the binary to a directory that is in your `$PATH` (e.g., `/usr/local/bin`):
   ```bash
   $ sudo mv cpt /usr/local/bin/
   ```
7. Verify the installation by running:
   ```bash
   $ cpt --help
   ```
