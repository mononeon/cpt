package testlib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const testlibURL = "https://raw.githubusercontent.com/MikeMirzayanov/testlib/master/testlib.h"

// EnsureTestlib checks if testlib.h exists in ~/.cpt/include/ and downloads it if not.
// Returns the absolute path to the include directory so the compiler can use it.
func EnsureTestlib() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	includeDir := filepath.Join(home, ".cpt", "include")
	if err := os.MkdirAll(includeDir, 0755); err != nil {
		return "", err
	}

	testlibPath := filepath.Join(includeDir, "testlib.h")
	if _, err := os.Stat(testlibPath); os.IsNotExist(err) {
		fmt.Printf("\033[36mDownloading testlib.h...\033[0m\n")
		resp, err := http.Get(testlibURL)
		if err != nil {
			return "", fmt.Errorf("failed to download testlib.h: %v", err)
		}
		defer resp.Body.Close()

		out, err := os.Create(testlibPath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}
		fmt.Printf("\033[32mtestlib.h successfully downloaded to %s\033[0m\n", testlibPath)
	}

	return includeDir, nil
}
