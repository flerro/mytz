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
				"zone":  "Europe/Rome",
				"time":  "2019-07-01T12:32",
				"other": "Europe/Kiev,Pacific/Auckland"},
		})

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		payload := decodePayload(t, response.Body)

		assertEqual(t, payload.RefTime, "01 Jul 19 12:32 CEST", "RefTime mismatch")
		assertEqual(t, len(payload.Timezones), 2, "Too many timezones returned")

		assertEqual(t, payload.Timezones["Europe/Kiev"], "01 Jul 19 13:32 EEST", "Europe/Kiev mismatch")
		assertEqual(t, payload.Timezones["Pacific/Auckland"], "01 Jul 19 22:32 NZST", "Europe/Kiev mismatch")
	})

	t.Run("All timezones", func(t *testing.T) {
		response, err := handler(events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"zone": "Europe/Rome",
				"time": "2019-07-01T12:32"},
		})

		if err != nil {
			t.Fatal("Everything should be ok")
		}

		payload := decodePayload(t, response.Body)
		assertEqual(t, len(payload.Timezones), len(DefaultTimeZones), "Invalid Timezones returned")
	})
}

func decodePayload(t *testing.T, body string) ResponsePayload {
	var payload ResponsePayload
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
