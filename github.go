package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var ghClient *github.Client

func pass(key string) (string, error) {
	result := shell("pass " + key)
	if strings.HasPrefix("Error: ", result) {
		return "", fmt.Errorf(result)
	}
	return result, nil
}

func getGHClient(ctx context.Context) *github.Client {
	if ghClient == nil {
		token, err := pass("github-token")
		if err != nil {
			log.Printf("error getting github-token: %v\n", err)
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
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
