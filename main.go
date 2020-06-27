package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"strings"

	"github.com/getlantern/systray"
)

func try(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	systray.Run(onReady, onExit)
}

func setIcon(count int) {
	f, err := os.Open("./icon.png")
	try(err)
	defer f.Close()
	img, err := png.Decode(f)
	try(err)
	img = addLabel(img, 1, 1, "6")
	b := bytes.NewBuffer([]byte{})
	try(png.Encode(b, img))
	systray.SetIcon(b.Bytes())
}

func onReady() {
	setIcon(0)
	mOpen := systray.AddMenuItem("Open Web", "Open the github notifications page")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	notificationsChan := ghNotificationsSub()

	go func() {
		for {

			select {
			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")

			case <-mOpen.ClickedCh:
				shell("xdg-open 'https://github.com/notifications?query=is%3Aunread'")

			case notifications := <-notificationsChan:
				mOpen.SetTitle(fmt.Sprintf("%d notifications", len(notifications)))
				// spew.Dump(notifications)
				setIcon(len(notifications))
			}
		}
	}()
}
func insideCircle(cx, cy, px, py, radius int) bool {
	dx := float64(cx - px)
	dy := float64(cy - py)
	distance_squared := dx*dx + dy*dy
	return distance_squared <= float64(radius*radius)
}
func addLabel(src image.Image, x, y int, label string) image.Image {
	b := src.Bounds()
	img := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)

	col := color.RGBA{0xE9, 0x54, 0x20, 255}

	size := 12

	for x := 32 - size; x < 32; x++ {
		for y := 0; y < size; y++ {
			if insideCircle(32-(size/2), size/2, x, y, size/2) {
				img.Set(x, y, col)
			}
		}
	}

	return img
}

func onExit() {
	// clean up here
}

func shell(cmd string) string {
	b, _ := exec.Command("bash", "-c", cmd).Output()

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	return lines[len(lines)-1]
}
