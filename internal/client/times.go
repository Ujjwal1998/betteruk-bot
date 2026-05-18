package client

import (
	"encoding/json"
	"fmt"
	"sort"
)

type timeOfDay struct {
	Format12Hour string `json:"format_12_hour"`
	Format24Hour string `json:"format_24_hour"`
}

type slotPrice struct {
	FormattedAmount string `json:"formatted_amount"`
}

type TimeSlot struct {
	StartsAt     timeOfDay `json:"starts_at"`
	EndsAt       timeOfDay `json:"ends_at"`
	Price        slotPrice `json:"price"`
	Spaces       int       `json:"spaces"`
	Location     string    `json:"location"`
	Name         string    `json:"name"`
	CategorySlug string    `json:"category_slug"`
	CompositeKey string    `json:"composite_key"`
}

type timesResponse struct {
	Data timesData `json:"data"`
}

// timesData accepts API "data" as either a JSON array or an object keyed by court/id.
type timesData []TimeSlot

func (d *timesData) UnmarshalJSON(b []byte) error {
	slots, err := parseTimesData(b)
	if err != nil {
		return err
	}
	*d = slots
	return nil
}

func parseTimesData(b []byte) ([]TimeSlot, error) {
	var slots []TimeSlot
	if err := json.Unmarshal(b, &slots); err == nil {
		sortTimeSlots(slots)
		return slots, nil
	}

	var byKey map[string]TimeSlot
	if err := json.Unmarshal(b, &byKey); err == nil {
		slots = make([]TimeSlot, 0, len(byKey))
		for _, slot := range byKey {
			slots = append(slots, slot)
		}
		sortTimeSlots(slots)
		return slots, nil
	}

	return nil, fmt.Errorf("times data is neither array nor object")
}

func sortTimeSlots(slots []TimeSlot) {
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartsAt.Format24Hour < slots[j].StartsAt.Format24Hour
	})
}

// GetTimes fetches available time slots for an activity at a venue on the given date.
func (c *Client) GetTimes(venueSlug, activitySlug, date string) ([]TimeSlot, error) {
	var result timesResponse
	path := fmt.Sprintf("/api/activities/venue/%s/activity/%s/times?date=%s", venueSlug, activitySlug, date)
	if err := c.adminGET(path, &result); err != nil {
		return nil, fmt.Errorf("get times: %w", err)
	}
	return result.Data, nil
}
