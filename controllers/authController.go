package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/natron-io/tenant-api/util"
)

var SECRET_KEY string

func GithubLogin(c *fiber.Ctx) error {
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=read:org&client_id=%s&redirect_uri=%s",
		util.CLIENT_ID, "http://127.0.0.1:56668/login/github/callback") //api.natron.io

	return c.Redirect(redirectURL)
}

func GithubCallback(c *fiber.Ctx) error {
	// get code from "code" query param
	code := c.Query("code")

	// util.InfoLogger.Printf("Received code: %s", code)

	githubAccessToken := util.GetGithubAccessToken(code)

	// util.InfoLogger.Printf("Received access token: %s", githubAccessToken)

	githubData := util.GetGithubTeams(githubAccessToken)

	// util.InfoLogger.Printf("Received github data: %s", githubData)

	return LoggedIn(c, githubData)
}

func LoggedIn(c *fiber.Ctx, githubData string) error {
	if githubData == "" {
		// return unauthorized
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// parse responsebody to map array
	var githubDataMap []map[string]interface{}
	json.Unmarshal([]byte(githubData), &githubDataMap)

	// get each github team slug
	var githubTeamSlugs []string
	for _, githubTeam := range githubDataMap {
		githubTeamSlugs = append(githubTeamSlugs, githubTeam["slug"].(string))
	}

	claims := jwt.MapClaims{
		"github_team_slugs": githubTeamSlugs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(SECRET_KEY))

	cookie := &fiber.Cookie{
		Name:    "tenant-api-token",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24),
		Path:    "/",
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Logged in",
		"data":    githubTeamSlugs,
	})

}

func CheckAuth(c *fiber.Ctx) []string {
	cookie := c.Cookies("tenant-api-token")

	if cookie == "" {
		util.WarningLogger.Printf("IP %s is not authorized", c.IP())
		return nil
	}

	token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	var githubTeamSlugs []string
	for _, githubTeam := range claims["github_team_slugs"].([]interface{}) {
		githubTeamSlugs = append(githubTeamSlugs, githubTeam.(string))
	}

	return githubTeamSlugs
}
