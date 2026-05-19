package display

import (
	"encoding/json"
	"testing"

	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
)

func TestBookableSlotsUniform(t *testing.T) {
	const raw = `[{"starts_at":{"format_24_hour":"08:20"},"ends_at":{"format_24_hour":"09:00"},"price":{"formatted":"£10.20"},"spaces":1,"action_to_show":{"status":"BOOK"},"slot_display_location":"Multiple"},{"starts_at":{"format_24_hour":"08:20"},"ends_at":{"format_24_hour":"09:00"},"price":{"formatted":"£10.20"},"spaces":1,"action_to_show":{"status":"BOOK"},"slot_display_location":"Multiple"}]`
	var slots []client.BookableSlot
	if err := json.Unmarshal([]byte(raw), &slots); err != nil {
		t.Fatal(err)
	}
	if !bookableSlotsUniform(slots) {
		t.Fatal("expected uniform slots")
	}
}
