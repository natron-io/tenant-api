package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var (
	CLIENT_ID     string
	CLIENT_SECRET string
	letters       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	SECRET_KEY    string
	CALLBACK_URL  string
	DEBUG         bool
)

// GetGithubAccessToken returns a github access token
func GetGithubAccessToken(code string) string {
	requestBodyMap := map[string]string{"client_id": CLIENT_ID, "client_secret": CLIENT_SECRET, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		ErrorLogger.Printf("Request creation failed: %s", reqerr)
	}

	InfoLogger.Printf("Request body: %s", requestJSON)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		ErrorLogger.Printf("Request failed: %s", respErr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type GithubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var githubAccessTokenResponse GithubAccessTokenResponse
	json.Unmarshal(respbody, &githubAccessTokenResponse)

	return githubAccessTokenResponse.AccessToken
}

// GetGithubData returns the github data
func GetGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		ErrorLogger.Printf("Request creation failed: %s", reqerr)
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		ErrorLogger.Printf("Request failed: %s", respErr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}

// GetGithubTeams returns the github teams
func GetGithubTeams(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/orgs/natron-io/teams", nil)
	if reqerr != nil {
		ErrorLogger.Printf("Request creation failed: %s", reqerr)
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		ErrorLogger.Printf("Request failed: %s", respErr)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}

// RandomStringBytes returns a random string of length n
func RandomStringBytes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
