package plugins

import (
	"fmt"
	"image"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/orsinium-labs/tellowerk/icons"
)

type UI struct {
	app fyne.App
	win fyne.Window

	battery    *canvas.Text
	warns      *canvas.Text
	speed      *canvas.Text
	video      *canvas.Image
	warnsState map[string]bool

	verticalSpeed int16
	northSpeed    int16
	eastSpeed     int16
}

var _ StateHandler = &UI{}

func NewUI() *UI {
	return &UI{
		app:        app.New(),
		warnsState: make(map[string]bool),
	}
}

func (ui *UI) Connect(pl *Plugins) {
}

func (ui *UI) Wait() {
	ui.app.Run()
}

func (ui *UI) Stop() error {
	ui.app.Quit()
	return nil
}

func (UI) icon(res fyne.Resource) *canvas.Image {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillOriginal
	img.SetMinSize(fyne.NewSize(24, 24))
	return img
}

func (ui *UI) Start() error {
	// ui.app.Settings().SetTheme(theme.LightTheme())
	ui.win = ui.app.NewWindow("tellowerk")

	ui.battery = canvas.NewText("?%", theme.ForegroundColor())
	ui.warns = canvas.NewText("", theme.ForegroundColor())
	ui.speed = canvas.NewText("? cm/s", theme.ForegroundColor())
	ui.video = canvas.NewImageFromImage(
		image.NewRGBA(image.Rect(0, 0, frameX, frameY)),
	)
	ui.video.SetMinSize(fyne.NewSize(frameX, frameY))
	content := container.NewHBox(
		container.New(
			layout.NewGridLayout(1),
			container.NewHBox(ui.icon(icons.BatteryStdOutlinedIconThemed), ui.battery),
			container.NewHBox(ui.icon(icons.SpeedOutlinedIconThemed), ui.speed),
			container.NewHBox(ui.icon(icons.WarningOutlinedIconThemed), ui.warns),
		),
		ui.video,
	)
	ui.win.SetContent(content)
	ui.win.Show()
	return nil
}

func (ui *UI) SetBattery(val int8) {
	if ui.battery == nil {
		return
	}
	ui.battery.Text = fmt.Sprintf("battery %d%%", val)
	if val <= 20 {
		ui.warnsState["low battery"] = true
	}
	ui.battery.Refresh()
}

func (ui *UI) SetNorthSpeed(val int16) {
	ui.northSpeed = val
	ui.updateSpeed()
}
func (ui *UI) SetEastSpeed(val int16) {
	ui.eastSpeed = val
	ui.updateSpeed()
}
func (ui *UI) SetVerticalSpeed(val int16) {
	ui.verticalSpeed = val
	ui.updateSpeed()
}

func (ui *UI) updateSpeed() {
	s := math.Sqrt(float64(
		ui.northSpeed*ui.northSpeed +
			ui.eastSpeed*ui.eastSpeed +
			ui.verticalSpeed*ui.verticalSpeed))
	ui.speed.Text = fmt.Sprintf("%d cm/s", int(s))
}

func (ui *UI) SetWarning(msg string, state bool) {
	ui.warnsState[msg] = state

	text := ""
	for msg, state := range ui.warnsState {
		if state {
			text += msg + "\n"
		}
	}
	ui.warns.Text = text
	ui.warns.Refresh()
}

func (ui *UI) SetFrame(img *RGB) {
	ui.video.File = ""
	ui.video.Image = img
	ui.video.Refresh()
}
