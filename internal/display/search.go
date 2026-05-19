package display

import (
	"fmt"

	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
)

// SearchRow is one selectable row in aggregated search results.
type SearchRow struct {
	Venue        client.Venue
	Slot         client.TimeSlot
	ActivitySlug string
	Date         string
}

// BuildSearchRows converts client results into display rows.
func BuildSearchRows(results []client.VenueTime, activitySlug, date string) []SearchRow {
	rows := make([]SearchRow, len(results))
	for i, r := range results {
		rows[i] = SearchRow{
			Venue:        r.Venue,
			Slot:         r.Slot,
			ActivitySlug: activitySlug,
			Date:         date,
		}
	}
	return rows
}

// PrintSearchResults prints a numbered table of availability across venues.
func PrintSearchResults(activityName, date string, rows []SearchRow) {
	fmt.Printf("--- %s on %s ---\n", activityName, date)
	if len(rows) == 0 {
		fmt.Println("  No available times at scanned venues")
		fmt.Println()
		return
	}
	fmt.Printf("  %-3s %-28s %-11s %-8s  %s\n", "#", "Venue", "Time", "Price", "Spaces")
	for i, r := range rows {
		price := r.Slot.Price.FormattedAmount
		if price == "" {
			price = "—"
		}
		name := r.Venue.Name
		if len(name) > 28 {
			name = name[:25] + "..."
		}
		fmt.Printf("  %-3d %-28s %s–%s  %-8s  %s\n",
			i+1,
			name,
			r.Slot.StartsAt.Format24Hour,
			r.Slot.EndsAt.Format24Hour,
			price,
			spacesLabel(r.Slot.Spaces),
		)
	}
	fmt.Println()
}
