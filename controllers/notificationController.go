package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
	"github.com/slack-go/slack"
)

// GetNotifications returns all slack notifications of a broadcast channel by authenticated users tenants
func GetNotifications(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	tenants := CheckAuth(c)
	if len(tenants) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// set authorization header Bearer token as util.SLACK_TOKEN
	if util.SLACK_TOKEN == "" {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: util.BroadCastChannelID,
		Limit:     10,
	}

	// make slack api call to get all notifications of the BroadCast channel
	resp, err := util.SlackClient.GetConversationHistory(&historyParams)

	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// parse messages to string slice
	var notifications []string
	for _, message := range resp.Messages {
		if message.Blocks.BlockSet != nil {
			notifications = append(notifications, message.Text)
		}
	}

	if err != nil {
		util.ErrorLogger.Printf("Error getting slack messages: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(notifications)
}
