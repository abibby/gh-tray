package main

import (
	"os/exec"
)

func xdgOpen(address string) error {
	return exec.Command("xdg-open", address).Run()
}
