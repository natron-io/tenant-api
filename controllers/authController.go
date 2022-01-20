package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/natron-io/tenant-api/util"
)

func GithubLogin(c *fiber.Ctx) error {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		util.CLIENT_ID, "https://api.natron.io/login/github/callback")

	return c.Redirect(redirectURL)
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	githubAccessToken := util.GetGithubAccessToken(code)

	githubData := util.GetGithubData(githubAccessToken)

	return LoggedIn(c, githubData)
}

func LoggedIn(c *fiber.Ctx, githubData string) error {
	if githubData != "" {
		// return unauthorized
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// set response header to application/json
	c.Set("Content-Type", "application/json")

	var prettyJSON bytes.Buffer

	// Pretty-print the JSON
	err := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	// Return the JSON
	return c.JSON(prettyJSON.String())
}
