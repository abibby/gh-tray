package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/getlantern/systray"
)

//go:generate go run build/main.go icon.png
//go:generate go run build/main.go icon-active.png
var files map[string][]byte = map[string][]byte{}

func try(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	systray.Run(onReady, onExit)
}

func setIcon(count int) {
	if count > 0 {
		systray.SetIcon(files["icon-active.png"])
	} else {
		systray.SetIcon(files["icon.png"])
	}
}

func onReady() {
	setIcon(0)
	mOpen := systray.AddMenuItem("Open Web", "Open the github notifications page")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	notificationsChan := ghNotificationsSub()

	for {

		select {
		case <-mQuit.ClickedCh:
			fmt.Println("Requesting quit")
			systray.Quit()
			fmt.Println("Finished quitting")

		case <-mOpen.ClickedCh:
			xdgOpen("https://github.com/notifications?query=is%3Aunread")

		case notifications := <-notificationsChan:
			mOpen.SetTitle(fmt.Sprintf("%d notifications", len(notifications)))
			// spew.Dump(notifications)
			setIcon(len(notifications))
		}
	}
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
