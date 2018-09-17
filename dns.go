package main

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// UpdateDNS updates DigitalOcean DNS with the new IP address.
func UpdateDNS(ip, domain, hostname string, tokenSource *TokenSource) error {
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	options := &godo.ListOptions{}
	editRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: hostname,
		Data: ip,
	}

	records, _, err := client.Domains.Records(context.Background(), domain, options)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type == "A" && record.Name == hostname {
			_, _, err = client.Domains.EditRecord(context.Background(), domain, record.ID, editRequest)
			if err == nil {
				return nil
			}

			logrus.Error("An error occurred when editing the domain record:")
			logrus.Error(err)
		}
	}

	// None found, create one
	_, _, err = client.Domains.CreateRecord(context.Background(), domain, editRequest)

	if err != nil {
		logrus.Error("An error occurred when creating the domain record.")
	}
	return err
}
