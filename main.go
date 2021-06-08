package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
)

func main() {
	webApi := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode := socketmode.New(
		webApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID
	fmt.Println(selfUserId)

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeSlashCommand:
				socketMode.Debugf("Got slash command: %v", envelope.Type)
				fmt.Println(envelope)
			case socketmode.EventTypeEventsAPI:
				socketMode.Debugf("Got events_api: %v", envelope.Type)
				fmt.Println(envelope.Type)
			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	socketMode.Run()
}
