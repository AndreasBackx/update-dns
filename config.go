package main

import (
	"encoding/json"
	"github.com/gobuffalo/packr"
	"golang.org/x/oauth2"
)

// TokenSource used for OAuth2.
type TokenSource struct {
	AccessToken string `json:"access_token"`
}

// Config required for app.
type Config struct {
	Domain     string `json:"domain"`
	Hostname   string `json:"hostname"`
	IPFilePath string `json:"ip_file_path"`

	TokenSource TokenSource `json:"token_source"`
}

// Token used for OAuth2.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// GetConfig returns the saved config.
func GetConfig() (*Config, error) {
	box := packr.NewBox("./secret")
	data := box.Bytes("config.json")
	config := &Config{}
	err := json.Unmarshal(data, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
