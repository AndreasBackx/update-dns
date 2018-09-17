package main

import "github.com/sirupsen/logrus"

func main() {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	ip, err := GetPublicIP()
	if err != nil {
		logrus.Errorln("Could not get the public IP address.")
		panic(err)
	}

	logrus.Infof("Current public IP address: %s.", ip)
	lastIP, _ := GetLastPublicIP(config.IPFilePath)
	// Don't update the DNS if the ip remains the same.
	if lastIP == ip {
		logrus.Infoln("IP Address did not change, exiting...")
		return
	}

	err = UpdateDNS(ip, config.Domain, config.Hostname, &config.TokenSource)
	if err != nil {
		logrus.Errorln("Could not update the DNS.")
		panic(err)
	}

	err = SavePublicIP(config.IPFilePath, ip)
	if err != nil {
		logrus.Errorf("Could not save the public IP address to %s.", config.IPFilePath)
		panic(err)
	}
	logrus.Infoln("Successfully updated the DNS record.")
}
