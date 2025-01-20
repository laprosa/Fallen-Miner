package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var getLastInputInfo = user32.NewProc("GetLastInputInfo")

func lastInputTicks() (uint32, error) {
	// set up struct to contain returned information
	var lastInputInfo struct {
		cbSize uint32
		dwTime uint32
	}
	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))

	// call the WinAPI function with a pointer to the struct
	errBool, _, err := getLastInputInfo.Call(
		uintptr(unsafe.Pointer(&lastInputInfo)),
	)

	if errBool == 0 {
		return 0, err
	}

	return lastInputInfo.dwTime, nil
}

var kernel32 = syscall.MustLoadDLL("kernel32.dll")
var getTickCount = kernel32.MustFindProc("GetTickCount")

func currentTicks() uint32 {
	ticks, _, _ := getTickCount.Call()
	return uint32(ticks)
}

// Duration returns the time since the last user input
func Get() (time.Duration, error) {
	lastInput, err := lastInputTicks()
	if err != nil {
		return 0, err
	}
	idleTicks := currentTicks() - lastInput
	idleDuration := time.Duration(idleTicks) * time.Millisecond
	return idleDuration, nil
}

func checkIdle() bool {
	var idleTime time.Duration
	idleTime, _ = Get()
	return idleTime.Seconds() > float64(idle_time)

}

func isProcessRunning(pid int) bool {
	// Run the "tasklist" command and capture the output.
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running tasklist command:", err)
		return false
	}

	// Check if the output contains the PID.
	outputStr := string(output)
	return strings.Contains(outputStr, fmt.Sprintf("%d", pid))
}

func checkProcessesForMatch(processNames []string) bool {
	// Run the "tasklist" command to get all processes.
	cmd := exec.Command("tasklist")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running tasklist command:", err)
		return false
	}

	// Convert output to a string and split into lines.
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// Check each line for the presence of process names.
	for _, processName := range processNames {
		for _, line := range lines {
			if strings.Contains(line, processName) {
				return true
			}
		}
	}
	return false
}
