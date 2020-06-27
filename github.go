package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var ghClient *github.Client

func personalAccessToken() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tokenFile := path.Join(d, "gh-tray/token")

	b, err := ioutil.ReadFile(tokenFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if b == nil {
		err = xdgOpen("https://github.com/settings/tokens/new?description=gh-tray&scopes=notifications")
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter token: ")
		b, err = reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		b = bytes.Trim(b, " \t\n")

		err = os.MkdirAll(path.Dir(tokenFile), 0755)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(tokenFile, b, 0644)
		if err != nil {
			return "", err
		}

	}

	return string(b), nil
}

func getGHClient(ctx context.Context) *github.Client {
	if ghClient == nil {
		token, err := personalAccessToken()
		if err != nil {
			log.Printf("error getting github-token: %v\n", err)
			os.Exit(1)
		}
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)

		ghClient = github.NewClient(tc)
	}
	return ghClient
}

func pushNotifications(notificationsChan chan []*github.Notification) {
	ctx := context.Background()
	client := getGHClient(ctx)

	notifs, _, err := client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	notificationsChan <- notifs
}

func ghNotificationsSub() chan []*github.Notification {
	notifications := make(chan []*github.Notification)
	go func() {
		pushNotifications(notifications)
		for range time.Tick(time.Second * 10) {
			pushNotifications(notifications)
		}
	}()
	return notifications
}
