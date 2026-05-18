package client

import (
	"encoding/json"
	"testing"
)

const sampleTimesJSONArray = `{"data":[{"starts_at":{"format_24_hour":"09:00"},"ends_at":{"format_24_hour":"09:40"},"price":{"formatted_amount":"£8.45"},"spaces":1,"location":"Court 1","category_slug":"badminton-40min"}]}`

const sampleTimesJSONObject = `{"data":{"4":{"starts_at":{"format_24_hour":"15:00"},"ends_at":{"format_24_hour":"16:00"},"price":{"formatted_amount":"£22.00"},"spaces":3,"location":"Multiple","category_slug":"pickleball-60mins"}}}`

func TestParseTimesResponseArray(t *testing.T) {
	testParseTimesResponse(t, sampleTimesJSONArray, 1, "09:00")
}

func TestParseTimesResponseObject(t *testing.T) {
	testParseTimesResponse(t, sampleTimesJSONObject, 1, "15:00")
}

func testParseTimesResponse(t *testing.T, raw string, wantLen int, wantStart string) {
	t.Helper()
	var result timesResponse
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != wantLen {
		t.Fatalf("got %d slots, want %d", len(result.Data), wantLen)
	}
	if result.Data[0].StartsAt.Format24Hour != wantStart {
		t.Errorf("start = %q, want %q", result.Data[0].StartsAt.Format24Hour, wantStart)
	}
}

func TestParseTimesResponseLegacy(t *testing.T) {
	testParseTimesResponse(t, sampleTimesJSONArray, 1, "09:00")
	s := mustUnmarshalTimes(t, sampleTimesJSONArray).Data[0]
	if s.Price.FormattedAmount != "£8.45" {
		t.Errorf("price = %q", s.Price.FormattedAmount)
	}
	if s.Spaces != 1 {
		t.Errorf("spaces = %d", s.Spaces)
	}
}

func mustUnmarshalTimes(t *testing.T, raw string) timesResponse {
	t.Helper()
	var result timesResponse
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatal(err)
	}
	return result
}
