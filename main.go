package main

import (
	"bytes"
	"fmt"
	"image"
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

func addLabel(src image.Image, x, y int, label string) image.Image {
	b := src.Bounds()
	img := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)

	// pixfont.DrawString(img, 10, 10, "Hello, World!", color.RGBA{0, 255, 255, 255})

	// col := color.RGBA{200, 100, 0, 255}
	// point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	// d := &font.Drawer{
	// 	Dst:  img,
	// 	Src:  image.NewUniform(col),
	// 	Face: basicfont.Face7x13,
	// 	Dot:  point,
	// }
	// d.DrawString(label)
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
