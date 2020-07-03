package main

import (
	"strings"
	"time"

	"encoding/json"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	DefaultTimeZones = []string{
		"America/Los_Angeles",
		"America/New_York",
		"Asia/Shanghai",
		"Europe/Istanbul",
	}
)

type ComparedTimezones struct {
	Reference Timezone   `json:"ref"`
	Others    []Timezone `json:"others"`
}

type Timezone struct {
	LocalTime string `json:"localTime"`
	Zone      string `json:"zone"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var ok bool
	var localTimeZone, localTime, other string
	var otherZones []string

	if localTimeZone, ok = request.QueryStringParameters["from"]; !ok {
		localTimeZone = "UTC"
	}

	if localTime, ok = request.QueryStringParameters["time"]; !ok {
		localTime = time.Now().Format(time.RFC822)
	}

	if other, ok = request.QueryStringParameters["to"]; !ok {
		otherZones = DefaultTimeZones
	} else {
		otherZones = strings.Split(other, ",")
	}

	payload := compareToTimeZones(localTimeZone, localTime, otherZones)
	jsonPayload, _ := json.Marshal(payload)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonPayload),
		StatusCode: 200,
	}, nil
}

func compareToTimeZones(localTimeZone string, localTime string, timeZones []string) ComparedTimezones {
	if localTimeZone == "" {
		return ComparedTimezones{}
	}

	refLocation, err := time.LoadLocation(localTimeZone)
	if err != nil {
		return ComparedTimezones{}
	}

	time.Local = refLocation

	refTime, err := dateparse.ParseLocal(localTime)
	if err != nil {
		return ComparedTimezones{}
	}

	localTz := Timezone{Zone: localTimeZone, LocalTime: refTime.Format(time.RFC822)}
	results := ComparedTimezones{Reference: localTz, Others: []Timezone{}}

	for _, tz := range timeZones {
		targetLocation, err := time.LoadLocation(tz)
		if err != nil {
			continue
		}

		ltime := refTime.In(targetLocation).Format(time.RFC822)
		ltimezone := Timezone{LocalTime: ltime, Zone: tz}
		results.Others = append(results.Others, ltimezone)
	}

	return results
}

func main() {
	lambda.Start(handler)
}
