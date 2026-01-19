package downloader

import (
	"fmt"
	"os"
	"os/exec"
)

// StartExternalDownload executes external downloader (XDM)
func (d *Downloader) StartExternalDownload(xdmPath, url string) error {
	if _, err := os.Stat(xdmPath); os.IsNotExist(err) {
		return fmt.Errorf("xdm executable not found at: %s", xdmPath)
	}

	// Format: xdm.exe --url "URL"
	// Note: Command line arguments might vary between XDM versions.
	// Standard XDM CLI usually supports just passing the URL as argument or --url
	// We'll try generic execution. Most download managers take URL as last arg or via specific flag.
	// For XDM (Java based wrapper often), passing URL usually works.

	fmt.Printf("\n🚀 Launching XDM...\nUrl: %s\n", url)

	// Command execution is non-blocking (starts separate process)
	cmd := exec.Command(xdmPath, "--url", url)

	// Create a detached process if possible or just start
	err := cmd.Start()
	if err != nil {
		// Try without --url flag if failed (some tools just take the url)
		cmd = exec.Command(xdmPath, url)
		if errRetry := cmd.Start(); errRetry != nil {
			return fmt.Errorf("failed to start XDM: %v (Retry error: %v)", err, errRetry)
		}
	}

	fmt.Println("✅ XDM command sent successfully!")
	return nil
}
