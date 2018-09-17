package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// GetPublicIP returns the public IP from the current machine.
func GetPublicIP() (string, error) {
	response, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// GetLastPublicIP returns last known public IP.
func GetLastPublicIP(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SavePublicIP saves the public IP to a file.
func SavePublicIP(filename, ip string) error {
	return ioutil.WriteFile(filename, []byte(ip), 0644)
}
