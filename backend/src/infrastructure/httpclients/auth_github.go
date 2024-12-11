package httpclients

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"recycle-waste-management-backend/src/domain/entities"
	"strings"

	"github.com/goccy/go-json"
)

type authGithub struct {
	url          string
	urlClient    string
	http         *http.Client
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type IAuthGithub interface {
	GitHubAccessToken(code string) (entities.TokenResponseGithub, error)
	GetUserGithub(token string) (entities.UserGithub, error)
}

func NewAuthGithub() IAuthGithub {
	return &authGithub{
		url:          "https://github.com/login/oauth/access_token",
		urlClient:    "https://api.github.com/user",
		http:         &http.Client{},
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("GITHUB_REDIRECT_URI"),
	}
}

func (a *authGithub) GitHubAccessToken(code string) (entities.TokenResponseGithub, error) {
	data := url.Values{}
	data.Set("client_id", a.ClientID)
	data.Set("client_secret", a.ClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", a.url, strings.NewReader(data.Encode()))
	if err != nil {
		return entities.TokenResponseGithub{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := a.http
	resp, err := client.Do(req)
	if err != nil {
		return entities.TokenResponseGithub{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.TokenResponseGithub{}, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return entities.TokenResponseGithub{}, fmt.Errorf("failed to get access token: %s", body)
	}

	var tokenResponse entities.TokenResponseGithub

	params, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing query string:", err)
	}

	for key, values := range params {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		switch key {
		case "error":
			tokenResponse.Error = value
		case "error_description":
			tokenResponse.ErrorDesc = value
		case "access_token":
			tokenResponse.AccessToken = value
		case "scope":
			tokenResponse.Scope = value
		case "token_type":
			tokenResponse.TokenType = value
		}
	}
	if tokenResponse.Error != "" && tokenResponse.ErrorDesc != "" {
		return tokenResponse, fmt.Errorf("failed to get access token: %s", tokenResponse.ErrorDesc)
	}
	return tokenResponse, nil
}

func (a *authGithub) GetUserGithub(token string) (entities.UserGithub, error) {
	url := a.urlClient
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return entities.UserGithub{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return entities.UserGithub{}, err
	}

	var user entities.UserGithub
	if err := json.Unmarshal(body, &user); err != nil {
		return user, err
	}
	return user, nil
}
