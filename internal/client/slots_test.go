package client

import (
	"encoding/json"
	"testing"
)

const sampleSlotsJSON = `{"data":[{"id":91112446,"name":"Pickleball 60mins","slot_display_location":"Court 1, Clissold Leisure Centre","price":{"formatted":"£12.70"},"starts_at":{"format_24_hour":"07:00"},"ends_at":{"format_24_hour":"08:00"},"spaces":1,"capacity":1,"action_to_show":{"status":"BOOK"},"composite_key":"63f4393f","location":{"name":"Court 1"}}]}`

func TestParseSlotsResponse(t *testing.T) {
	var result slotsResponse
	if err := json.Unmarshal([]byte(sampleSlotsJSON), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("got %d slots, want 1", len(result.Data))
	}
	s := result.Data[0]
	if s.ID != 91112446 || s.ActionToShow.Status != "BOOK" {
		t.Fatalf("unexpected slot: %+v", s)
	}
}
