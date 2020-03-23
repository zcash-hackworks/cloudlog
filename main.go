package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber"
	flog "github.com/gofiber/logger"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

var slackURL string

type slackMessage struct {
	Blocks []slack.MessageBlockType
}

func main() {
	versionFlag := flag.Bool("version", false, "print version information")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("(version=%s, branch=%s, gitcommit=%s)\n", Version, Branch, GitCommit)
		fmt.Printf("(go=%s, user=%s, date=%s)\n", GoVersion, BuildUser, BuildDate)
		os.Exit(0)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()
	app := fiber.New()
	app.Use(flog.New())

	slackURL = os.Getenv("slackURL")
	if slackURL == "" {
		log.Fatalw("slackURL ENV value missing")
	}

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("cloudlog!")
	})

	app.Post("/inbox", func(c *fiber.Ctx) {
		var message map[string]interface{}

		if err := c.BodyParser(&message); err != nil {
			log.Error(err)
			c.Status(500).Send(err)
			return
		}
		log.Infow("Got an inbox message!", "message", message)
		cloudEvent, err := coerceCloudEvent(message)
		if err != nil {
			log.Error(err)
			c.Status(500).Send(err)
			return
		}

		if err := sendSlackMessage(cloudEvent); err != nil {
			log.Error(err)
			c.Status(500).Send(err)
			return
		}
		c.JSON(message)
	})

	app.Listen(3000)
}

func coerceCloudEvent(message map[string]interface{}) (te *TektonCloudEvent, err error) {
	var taskRun, status, condition map[string]interface{}
	var conditions []interface{}
	var ok bool
	var podName, reason, cMessage string
	if taskRun, ok = message["taskRun"].(map[string]interface{}); !ok {
		return nil, errors.New("taskRun not found in message")
	}

	if status, ok = taskRun["status"].(map[string]interface{}); !ok {
		return nil, errors.New("status not found in taskRun")
	}

	if podName, ok = status["podName"].(string); !ok {
		return nil, errors.New("podName not found in status")
	}

	if conditions, ok = status["conditions"].([]interface{}); !ok {
		return nil, errors.New("conditions not found in status")
	}
	if condition, ok = conditions[0].(map[string]interface{}); !ok {
		fmt.Println("condition not found in conditions")
	}
	if reason, ok = condition["reason"].(string); !ok {
		fmt.Println("reason not found in condition")
	}
	if cMessage, ok = condition["message"].(string); !ok {
		fmt.Println("message not found in condition")
	}

	te = &TektonCloudEvent{
		Title:   podName + " " + reason,
		Reason:  reason,
		Message: cMessage,
		PodName: podName,
	}

	return te, nil
}

type TektonCloudEvent struct {
	Title   string `json:"string"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	PodName string `json:"pod_name"`
}

func sendSlackMessage(cloudEvent *TektonCloudEvent) error {
	text := fmt.Sprintf("*%s*\n\n```%s```\nPodName: %s", cloudEvent.Title, cloudEvent.Message, cloudEvent.PodName)
	msg := slack.WebhookMessage{
		Text: text,
	}

	err := slack.PostWebhook(slackURL, &msg)
	if err != nil {
		return err
	}

	return nil
}
