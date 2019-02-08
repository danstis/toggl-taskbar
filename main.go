package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/fogleman/gg"
	"github.com/getlantern/systray"
)

// Version contains the package version
var Version = "0.0.0-default"

type togglTime struct {
	hours   int
	minutes int
}

// Main entry point for the app.
func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	var t togglTime
	updateIcon("00")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuit.ClickedCh
		log.Println("Applicaiton exiting...")
		systray.Quit()
	}()

	systray.SetTitle("Toggl Weekly Time")

	for {
		t = getWeeklyTime()
		log.Printf("- Got new time %v:%v\n", t.hours, t.minutes)
		updateIcon(fmt.Sprintf("%v", t.hours))
		systray.SetTooltip(fmt.Sprintf("Toggl time tracker: %v:%v", t.hours, t.minutes))
		time.Sleep(60 * time.Second)
	}
}

func onExit() {
	// Cleaning stuff here.
}

func getWeeklyTime() togglTime {
	return togglTime{
		hours:   0,
		minutes: 0,
	}
}

func updateIcon(hours string) {
	icon, err := createIcon(16, 16, hours)
	if err != nil {
		log.Fatalf("Error generating icon: %v", err)
	}
	systray.SetIcon(icon)
}

func createIcon(x, y int, label string) ([]byte, error) {
	dc := gg.NewContext(x, y)
	// Create a trasparant background
	dc.SetColor(color.Transparent)
	dc.Clear()
	// Add Text
	dc.SetHexColor("#FFFFFF")
	if err := dc.LoadFontFace("assets/fonts/Go-Bold.ttf", 14); err != nil {
		return []byte{}, err
	}
	dc.DrawStringAnchored(label, float64(x/2), float64(y/2), 0.5, 0.5)

	buf := new(bytes.Buffer)
	err := ico.Encode(buf, dc.Image())
	if err != nil {
		return []byte{}, err
	}
	img := buf.Bytes()

	icoimg, _ := os.Create("new.ico")
	defer icoimg.Close()
	_ = ico.Encode(icoimg, dc.Image())

	return []byte(img), nil
}
