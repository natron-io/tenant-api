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
	responseConversationHistory, err := util.SlackClient.GetConversationHistory(&historyParams)

	if err != nil {
		util.ErrorLogger.Printf("%s", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// get the slack url of the channel

	// create map of notifications with client_msg_id as key and message as value and username

	type Notification struct {
		ClientMsgID   string `json:"client_msg_id"`
		Message       string `json:"message"`
		UserRealName  string `json:"user_real_name"`
		UserAvatarURL string `json:"user_avatar_url"`
		UnixTimestamp string `json:"unix_timestamp"`
		LinkToMessage string `json:"link_to_message"`
	}

	notifications := []Notification{}
	for _, message := range responseConversationHistory.Messages {
		responseUserName, err := util.SlackClient.GetUserInfo(message.User)

		if err != nil {
			util.ErrorLogger.Printf("%s", err)
			continue
		}

		if message.ClientMsgID != "" {
			// json map of notifications
			notifications = append(notifications, Notification{
				ClientMsgID:   message.ClientMsgID,
				Message:       message.Text,
				UserRealName:  responseUserName.Profile.RealName,
				UserAvatarURL: responseUserName.Profile.Image192,
				UnixTimestamp: message.Timestamp,
				LinkToMessage: util.SlackURL + "/archives/" + util.BroadCastChannelID + "/p" + message.Timestamp,
			})
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
