package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/natron-io/tenant-api/util"
)

func GithubLogin(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=read:org&client_id=%s&redirect_uri=%s",
		util.CLIENT_ID, util.CALLBACK_URL+"/login/github/callback")

	return c.Redirect(redirectURL)
}

func FrontendGithubLogin(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// get access_token from data
	if githubCode := data["github_code"]; githubCode == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	} else {
		// util.InfoLogger.Printf("Received code: %s", accessToken)

		githubAccessToken := util.GetGithubAccessToken(githubCode)

		// util.InfoLogger.Printf("Received github data: %s", githubData)

		githubData := util.GetGithubTeams(githubAccessToken)

		return LoggedIn(c, githubData)
	}

}

func GithubCallback(c *fiber.Ctx) error {

	util.InfoLogger.Printf("%s %s %s", c.IP(), c.Method(), c.Path())

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

	if githubTeamSlugs == nil {
		// return unauthorized
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims := jwt.MapClaims{
		"github_team_slugs": githubTeamSlugs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(util.SECRET_KEY))

	if !util.FRONTENDAUTH_ENABLED {
		cookie := &fiber.Cookie{
			Name:    "tenant-api-token",
			Value:   tokenString,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		}

		c.Cookie(cookie)
	}

	// return token
	return c.JSON(fiber.Map{
		"token": tokenString,
	})

}

func CheckAuth(c *fiber.Ctx) []string {
	var token *jwt.Token
	var tokenString string
	cookie := c.Cookies("tenant-api-token")

	if util.FRONTENDAUTH_ENABLED {

		var body map[string]string
		if err := c.BodyParser(&body); err != nil {
			return nil
		}

		reqToken := c.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		util.InfoLogger.Printf("reqToken %s", reqToken)

		if len(splitToken) == 2 && splitToken[1] != "" {
			tokenString = splitToken[1]
		} else if body["token"] != "" {
			tokenString = body["token"]
		} else {
			return nil
		}

		token, _ = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(util.SECRET_KEY), nil
		})

	} else {
		if cookie == "" {
			util.WarningLogger.Printf("IP %s is not authorized", c.IP())
			return nil
		}

		token, _ = jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(util.SECRET_KEY), nil
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["github_team_slugs"] == nil {
		util.WarningLogger.Printf("IP %s is not authorized", c.IP())
		return nil
	}

	var githubTeamSlugs []string
	for _, githubTeam := range claims["github_team_slugs"].([]interface{}) {
		githubTeamSlugs = append(githubTeamSlugs, githubTeam.(string))
	}

	return githubTeamSlugs
}

func Logout(c *fiber.Ctx) error {
	cookie := &fiber.Cookie{
		Name:     "tenant-api-token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Logged out",
	})
}
