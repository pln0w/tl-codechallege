package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// CallSlave - calls external server at given URL and returns response as string
func CallSlave(url string) string {

	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	// Dial slave server
	response, err := netClient.Get(url)
	if err != nil {
		log.Error(err.Error())
	}

	bodyString := ""
	bodyBytes, rErr := ioutil.ReadAll(response.Body)
	if rErr != nil {
		log.Error(rErr.Error())
	} else {
		bodyString = string(bodyBytes)
	}

	return bodyString
}
