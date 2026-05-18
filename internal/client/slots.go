package client

import (
	"fmt"
	"net/url"
)

type slotLocation struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type bookableSlotPrice struct {
	Formatted string `json:"formatted"`
}

type slotAction struct {
	Status string `json:"status"`
}

// BookableSlot is a court/resource available to book for a chosen time window.
type BookableSlot struct {
	ID                  int               `json:"id"`
	Name                string            `json:"name"`
	Location            slotLocation      `json:"location"`
	SlotDisplayLocation string            `json:"slot_display_location"`
	Price               bookableSlotPrice `json:"price"`
	StartsAt            timeOfDay         `json:"starts_at"`
	EndsAt              timeOfDay         `json:"ends_at"`
	Spaces              int               `json:"spaces"`
	Capacity            int               `json:"capacity"`
	ActionToShow        slotAction        `json:"action_to_show"`
	CompositeKey        string            `json:"composite_key"`
}

type slotsResponse struct {
	Data []BookableSlot `json:"data"`
}

// GetSlots fetches bookable courts for a time window (requires login).
func (c *Client) GetSlots(venueSlug, activitySlug, date, startTime, endTime, compositeKey string) ([]BookableSlot, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("login required: set --auth-token or BETTER_AUTH_TOKEN")
	}

	q := url.Values{}
	q.Set("date", date)
	q.Set("start_time", startTime)
	q.Set("end_time", endTime)
	q.Set("composite_key", compositeKey)

	path := fmt.Sprintf("/api/activities/venue/%s/activity/%s/slots?%s",
		venueSlug, activitySlug, q.Encode())

	var result slotsResponse
	if err := c.adminGET(path, &result); err != nil {
		return nil, fmt.Errorf("get slots: %w", err)
	}
	return result.Data, nil
}
