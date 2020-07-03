package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandler(t *testing.T) {

	t.Run("Successful Request", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"time": "2019-07-01T12:32",
				"from": "Europe/Rome",
				"to":   "Europe/Kiev,Pacific/Auckland"},
		})

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		payload := decodePayload(t, response.Body)

		assertEqual(t, payload.Reference, Timezone{"01 Jul 19 12:32 CEST", "Europe/Rome"}, "Reference mismatch")
		assertEqual(t, len(payload.Others), 2, "Others count mismatch")

		assertEqual(t, payload.Others[0], Timezone{Zone: "Europe/Kiev", LocalTime: "01 Jul 19 13:32 EEST"}, "Timezone mismatch")
		assertEqual(t, payload.Others[1], Timezone{Zone: "Pacific/Auckland", LocalTime: "01 Jul 19 22:32 NZST"}, "Timezone mismatch")
	})

	t.Run("Return default timezones", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"zone": "Europe/Rome",
				"time": "2019-07-01T12:32"},
		})

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		payload := decodePayload(t, response.Body)
		assertEqual(t, len(payload.Others), len(DefaultTimeZones), "Invalid Others returned")
	})
}

func decodePayload(t *testing.T, body string) ComparedTimezones {
	var payload ComparedTimezones
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatal("Unable to unmarshal request")
	}
	return payload
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	message = fmt.Sprintf("%s -- '%v' != '%v'", message, a, b)
	t.Fatal(message)
}
