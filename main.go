//go:generate powershell -NoLogo -NoProfile -ExecutionPolicy Unrestricted -File ./.version.ps1
package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/fogleman/gg"
	"github.com/getlantern/systray"
	"gopkg.in/resty.v1"
)

// Version contains the package version
var Version = "0.0.0-dev"

const (
	configFile     = "config.toml"
	templateFormat = "%s: %d:%02d"
)

type togglTime struct {
	hours   int
	minutes int
}

// Settings contains the application configuration
type Settings struct {
	Token              string `toml:"token"`
	Email              string `toml:"email"`
	SyncInterval       int    `toml:"syncInterval"`
	HighlightThreshold int    `toml:"highlightThreshold"`
	UserID             string
	Workspaces         []Workspaces
}

// Workspaces stores the available toggl workspaces for the user
type Workspaces struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Main entry point for the app.
func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	// Get the settings
	var config Settings
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}
	err = config.getUserDetail()
	if err != nil {
		log.Printf("Failed to get Toggl user details: %v\n", err)
	}
	log.Printf("- Workspaces:%v\n", config.Workspaces)

	// Configure the systray item
	updateIcon(0, config.HighlightThreshold)
	mTitle := systray.AddMenuItem("Toggl Weekly Tracker", "Title")
	mTitle.Disable()
	mVersion := systray.AddMenuItem(fmt.Sprintf("v%v", Version), "Version")
	mVersion.Disable()
	systray.AddSeparator()
	systray.SetTitle("Toggl Weekly Time")
	menuItems := make(map[int32]*systray.MenuItem)
	for _, item := range config.Workspaces {
		menuItems[item.ID] = systray.AddMenuItem(fmt.Sprintf(templateFormat, item.Name, 0, 0), item.Name)
	}
	menuItems[0] = systray.AddMenuItem(fmt.Sprintf(templateFormat, "Total", 0, 0), "Total")
	systray.AddSeparator()
	mRefresh := systray.AddMenuItem("Refresh", "Refresh the data")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				log.Println("Application exiting...")
				systray.Quit()
				return
			case <-mRefresh.ClickedCh:
				log.Println("Manual refresh triggered...")
				refreshData(&config, menuItems)
			}
		}
	}()

	go func() {
		for {
			refreshData(&config, menuItems)
			time.Sleep(time.Duration(config.SyncInterval) * time.Minute)
		}
	}()
}

func refreshData(config *Settings, menuItems map[int32]*systray.MenuItem) {
	totalTime := togglTime{hours: 0, minutes: 0}
	for _, item := range config.Workspaces {
		t, err := getWeeklyTime(config, fmt.Sprint(item.ID))
		if err != nil {
			log.Printf("Failed to get Toggl details: %v\n", err)
		}
		log.Printf("- %s [%s] time: %d:%02d\n", item.Name, fmt.Sprint(item.ID), t.hours, t.minutes)
		// Set the title of the menuItem to contain the time for the individual workspace
		menuItems[item.ID].SetTitle(fmt.Sprintf(templateFormat, item.Name, t.hours, t.minutes))
		totalTime.add(t)
	}

	log.Printf("- Got new total time %d:%02d\n", totalTime.hours, totalTime.minutes)
	updateIcon(int(totalTime.hours), config.HighlightThreshold)
	systray.SetTooltip(fmt.Sprintf("Toggl time tracker: %d:%02d", totalTime.hours, totalTime.minutes))
	menuItems[0].SetTitle(fmt.Sprintf(templateFormat, "Total", totalTime.hours, totalTime.minutes))
}

func onExit() {
	// Cleaning stuff here.
}

func (c *Settings) getUserDetail() error {
	type UserResponse struct {
		Data struct {
			ID         int32        `json:"id"`
			Workspaces []Workspaces `json:"workspaces"`
		} `json:"data"`
	}
	var ur UserResponse
	toggl := resty.New().SetHostURL("https://api.track.toggl.com/api/v8").SetBasicAuth(c.Token, "api_token")

	_, err := toggl.R().
		SetQueryParams(map[string]string{
			"user_agent": c.Email,
		}).
		SetResult(&ur).
		Get("/me?with_related_data=true")
	if err != nil {
		return fmt.Errorf("unable to get user details from the Toggl API: %v", err)
	}

	c.UserID = fmt.Sprint(ur.Data.ID)
	c.Workspaces = ur.Data.Workspaces
	return nil
}

func getWeeklyTime(c *Settings, w string) (togglTime, error) {
	closedTime, err := getClosedTimeEntries(c, w)
	if err != nil {
		return togglTime{}, fmt.Errorf("failed to get closed time entries: %v", err)
	}
	openTime, err := getOpenTimeEntry(c, w)
	if err != nil {
		return togglTime{}, fmt.Errorf("failed to get open time entry: %v", err)
	}

	return togglTime{
		hours:   getHours(closedTime + openTime),
		minutes: getMinutes(closedTime + openTime),
	}, nil
}

func getClosedTimeEntries(c *Settings, w string) (time.Duration, error) {
	type WeeklyResponse struct {
		TotalGrand int `json:"total_grand"`
	}
	var ct WeeklyResponse

	toggleReports := resty.New().SetHostURL("https://api.track.toggl.com/reports/api/v2").SetBasicAuth(c.Token, "api_token")

	_, err := toggleReports.R().
		SetQueryParams(map[string]string{
			"user_agent":   c.Email,
			"workspace_id": w,
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

func getOpenTimeEntry(c *Settings, w string) (time.Duration, error) {
	type TimeEntriesResponse struct {
		Data struct {
			WID      int32 `json:"wid"`
			Duration int32 `json:"duration"`
		} `json:"data"`
	}
	var ot TimeEntriesResponse

	toggl := resty.New().SetHostURL("https://api.track.toggl.com/api/v8").SetBasicAuth(c.Token, "api_token")

	_, err := toggl.R().
		SetQueryParams(map[string]string{
			"user_agent": c.Email,
			"wid":        w,
		}).
		SetResult(&ot).
		Get("/time_entries/current")
	if err != nil {
		return time.Duration(0), fmt.Errorf("unable to get current time entry from the Toggl API: %v", err)
	}

	// if the returned duration is not negative then there is no open entry.
	// we also filter entries that do not match the workspace here.
	if ot.Data.Duration >= 0 || fmt.Sprint(ot.Data.WID) != w {
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

func (t *togglTime) add(n togglTime) {
	t.minutes += n.minutes
	t.hours += n.hours
	if t.minutes >= 60 {
		t.hours++
		t.minutes -= 60
	}
}

func updateIcon(hours, threshold int) {
	icon, err := createIcon(16, 16, hours, threshold)
	if err != nil {
		log.Fatalf("Error generating icon: %v", err)
	}
	systray.SetIcon(icon)
}

func createIcon(x, y, hours, threshold int) ([]byte, error) {
	dc := gg.NewContext(x, y)
	if hours >= threshold {
		// Create a red background
		dc.SetHexColor("#9E0000")
		dc.Clear()
	}
	// Add Text
	dc.SetHexColor("#FFFFFF")
	if err := dc.LoadFontFace("assets/fonts/Go-Bold.ttf", 14); err != nil {
		return []byte{}, err
	}
	dc.DrawStringAnchored(fmt.Sprintf("%v", hours), float64(x/2), float64(y/2), 0.5, 0.5)

	buf := new(bytes.Buffer)
	err := ico.Encode(buf, dc.Image())
	if err != nil {
		return []byte{}, err
	}
	img := buf.Bytes()

	return []byte(img), nil
}
