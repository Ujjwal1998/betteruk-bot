package client

import (
	"sort"
	"testing"
)

func TestSearchResultsSortedByTime(t *testing.T) {
	results := []VenueTime{
		{Venue: Venue{Name: "B"}, Slot: TimeSlot{StartsAt: timeOfDay{Format24Hour: "11:00"}}},
		{Venue: Venue{Name: "A"}, Slot: TimeSlot{StartsAt: timeOfDay{Format24Hour: "09:00"}}},
		{Venue: Venue{Name: "C"}, Slot: TimeSlot{StartsAt: timeOfDay{Format24Hour: "11:00"}}},
	}

	sort.Slice(results, func(i, j int) bool {
		ti := results[i].Slot.StartsAt.Format24Hour
		tj := results[j].Slot.StartsAt.Format24Hour
		if ti != tj {
			return ti < tj
		}
		return results[i].Venue.Name < results[j].Venue.Name
	})

	if results[0].Slot.StartsAt.Format24Hour != "09:00" {
		t.Fatalf("first slot time = %q", results[0].Slot.StartsAt.Format24Hour)
	}
	if results[1].Venue.Name != "B" || results[2].Venue.Name != "C" {
		t.Fatalf("tie-break sort failed: %+v", results)
	}
}
