package main

import (
	"encoding/json"
	"reflect"
	"strings"

	//"flag"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	DefaultTimeZones = []string{
		"America/Argentina/Buenos_Aires",
		"America/Los_Angeles",
		"America/New_York",
		"America/Sao_Paulo",
		"Asia/Hong_Kong",
		"Asia/Shanghai",
		"Australia/Sydney",
		"Europe/Berlin",
		"Europe/Istanbul",
		"Europe/London",
		"Europe/Moscow",
		"Europe/Kiev",
		"Pacific/Auckland",
	}
)

type ResponsePayload struct {
	RefTime   string            `json:"referenceTime"`
	Timezones map[string]string `json:"timezones"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var ok bool
	var localTimeZone, localTime, other string
	var otherZones []string

	if localTimeZone, ok = request.QueryStringParameters["zone"]; !ok {
		localTimeZone = "UTC"
	}

	if localTime, ok = request.QueryStringParameters["time"]; !ok {
		localTime = time.Now().Format(time.RFC822)
	}

	fmt.Println(">>>", reflect.ValueOf(request.QueryStringParameters).MapKeys())

	if other, ok = request.QueryStringParameters["other"]; !ok {
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

func compareToTimeZones(localTimeZone string, localTime string, timeZones []string) ResponsePayload {

	if localTimeZone == "" {
		return ResponsePayload{}
	}

	refLocation, err := time.LoadLocation(localTimeZone)
	if err != nil {
		return ResponsePayload{}
	}

	time.Local = refLocation

	refTime, err := dateparse.ParseLocal(localTime)
	if err != nil {
		return ResponsePayload{}
	}

	results := ResponsePayload{RefTime: refTime.Format(time.RFC822), Timezones: map[string]string{}}

	for _, tz := range timeZones {
		targetLocation, err := time.LoadLocation(tz)
		if err != nil {
			results.Timezones[tz] = err.Error()
		}

		results.Timezones[tz] = refTime.In(targetLocation).Format(time.RFC822)
	}

	return results
}

func main() {
	lambda.Start(handler)
}
