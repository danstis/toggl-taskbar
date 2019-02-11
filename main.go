package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/fogleman/gg"
	"github.com/getlantern/systray"
	"gopkg.in/resty.v1"
)

// Version contains the package version
var Version = "0.0.0-default"

const configFile = "config.toml"

type togglTime struct {
	hours   int
	minutes int
}

// Settings contains the application configuration
type Settings struct {
	Token        string `toml:"token"`
	UserID       string `toml:"userId"`
	WorkspaceID  string `toml:"workspaceId"`
	Email        string `toml:"email"`
	SyncInterval int    `toml:"syncInterval"`
}

// Main entry point for the app.
func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	var t togglTime
	var config Settings
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}
	updateIcon("00")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuit.ClickedCh
		log.Println("Applicaiton exiting...")
		systray.Quit()
	}()

	systray.SetTitle("Toggl Weekly Time")

	log.Printf("Refresh interval: %v\n", time.Duration(config.SyncInterval)*time.Minute)

	for {
		t, err = getWeeklyTime(&config)
		if err != nil {
			log.Printf("Failed to get Toggl details: %v\n", err)
		}
		log.Printf("- Got new time %d:%2d\n", t.hours, t.minutes) // TODO: remove this when logging goes away
		updateIcon(fmt.Sprintf("%v", t.hours))
		systray.SetTooltip(fmt.Sprintf("Toggl time tracker: %d:%2d", t.hours, t.minutes))
		time.Sleep(time.Duration(config.SyncInterval) * time.Minute)
	}
}

func onExit() {
	// Cleaning stuff here.
}

func getWeeklyTime(c *Settings) (togglTime, error) {
	closedTime, err := getClosedTimeEntries(c)
	if err != nil {
		return togglTime{}, fmt.Errorf("failed to get closed time entries: %v", err)
	}
	openTime, err := getOpenTimeEntry(c)
	if err != nil {
		return togglTime{}, fmt.Errorf("failed to get open time entry: %v", err)
	}

	return togglTime{
		hours:   getHours(closedTime + openTime),
		minutes: getMinutes(closedTime + openTime),
	}, nil
}

func getClosedTimeEntries(c *Settings) (time.Duration, error) {
	type WeeklyResponse struct {
		TotalGrand int `json:"total_grand"`
	}
	var ct WeeklyResponse

	toggleReports := resty.New().SetHostURL("https://toggl.com/reports/api/v2").SetBasicAuth(c.Token, "api_token")

	_, err := toggleReports.R().
		SetQueryParams(map[string]string{
			"user_agent":   c.Email,
			"workspace_id": c.WorkspaceID,
			"user_ids":     c.UserID,
			"since":        getLastMonday(),
		}).
		SetResult(&ct).
		Get("/weekly")
	if err != nil {
		return time.Duration(0), fmt.Errorf("unable to get summary report from the Toggl API: %v", err)
	}

	return time.Duration(ct.TotalGrand) * time.Millisecond, nil
}

func getOpenTimeEntry(c *Settings) (time.Duration, error) {
	type TimeEntriesResponse struct {
		Data struct {
			Duration int32 `json:"duration"`
		} `json:"data"`
	}
	var ot TimeEntriesResponse

	toggl := resty.New().SetHostURL("https://www.toggl.com/api/v8").SetBasicAuth(c.Token, "api_token")

	_, err := toggl.R().
		SetQueryParams(map[string]string{
			"user_agent": c.Email,
			"wid":        c.WorkspaceID,
		}).
		SetResult(&ot).
		Get("/time_entries/current")
	if err != nil {
		return time.Duration(0), fmt.Errorf("unable to get current time entry from the Toggl API: %v", err)
	}

	// if the returned duration is not negative then there is no open entry.
	if ot.Data.Duration >= 0 {
		return 0, nil
	}

	// Calculate the number of seconds based on the input data.
	// Unix epoc plus returned value of duration = seconds the current entry has been running for.
	od := int32(time.Now().Unix()) + ot.Data.Duration

	return time.Duration(od) * time.Second, nil
}

func getLastMonday() string {
	t := time.Now()
	delta := (int(t.Weekday()) - 1) * -1
	t = t.AddDate(0, 0, delta)

	return t.Format("2006-01-02")
}

func getHours(t time.Duration) int {
	d := t.Round(time.Minute)
	h := d / time.Hour
	return int(h)
}

func getMinutes(t time.Duration) int {
	d := t.Round(time.Minute) % time.Hour
	m := d / time.Minute
	return int(m)
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

	return []byte(img), nil
}
