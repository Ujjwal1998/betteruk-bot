package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
)

type VenueResult struct {
	Venue client.Venue
	Slots []client.TimeSlot
	Err   error
}

// PrintTimes prints numbered available time windows.
func PrintTimes(venueName string, times []client.TimeSlot, availableOnly bool) {
	fmt.Printf("--- %s ---\n", venueName)

	slots := times
	if availableOnly {
		filtered := slots[:0]
		for _, s := range slots {
			if s.Spaces > 0 {
				filtered = append(filtered, s)
			}
		}
		slots = filtered
	}

	if len(slots) == 0 {
		fmt.Println("  No available slots")
		fmt.Println()
		return
	}

	for i, s := range slots {
		start := s.StartsAt.Format24Hour
		end := s.EndsAt.Format24Hour
		spaces := spacesLabel(s.Spaces)
		loc := ""
		if s.Location != "" {
			loc = fmt.Sprintf("  [%s]", s.Location)
		}
		price := s.Price.FormattedAmount
		if price == "" {
			price = "—"
		}
		fmt.Printf("  %2d. %s–%s   %-8s  %s%s\n", i+1, start, end, price, spaces, loc)
	}
	fmt.Println()
}

// PrintBookableSlots prints courts/resources returned by the authenticated slots API.
func PrintBookableSlots(slots []client.BookableSlot, err error) {
	PrintBookableSlotsForVenue("", slots, err)
}

// PrintBookableSlotsForVenue prints bookable courts, with an optional venue heading.
func PrintBookableSlotsForVenue(venueName string, slots []client.BookableSlot, err error) {
	if err != nil {
		fmt.Printf("  [error: %v]\n\n", err)
		return
	}
	if len(slots) == 0 {
		fmt.Println("  No bookable courts")
		fmt.Println()
		return
	}

	s := slots[0]
	if venueName != "" {
		fmt.Printf("--- %s · %s–%s ---\n", venueName, s.StartsAt.Format24Hour, s.EndsAt.Format24Hour)
	}

	if len(slots) > 1 && bookableSlotsUniform(slots) {
		price := s.Price.Formatted
		if price == "" {
			price = "—"
		}
		status := s.ActionToShow.Status
		loc := courtLabel(s)
		fmt.Printf("  %d courts · %s–%s   %-8s  %s each  %s\n",
			len(slots),
			s.StartsAt.Format24Hour,
			s.EndsAt.Format24Hour,
			price,
			spacesLabel(s.Spaces),
			status,
		)
		if loc != "" && loc != "Multiple" {
			fmt.Printf("  (%s)\n", loc)
		}
		fmt.Println()
		return
	}

	for i, slot := range slots {
		loc := courtLabel(slot)
		price := slot.Price.Formatted
		if price == "" {
			price = "—"
		}
		status := slot.ActionToShow.Status
		if status == "" {
			status = "—"
		}
		fmt.Printf("  %d. %s–%s   %-8s  %s  [%s]  %s\n",
			i+1,
			slot.StartsAt.Format24Hour,
			slot.EndsAt.Format24Hour,
			price,
			spacesLabel(slot.Spaces),
			loc,
			status,
		)
	}
	fmt.Println()
}

func courtLabel(s client.BookableSlot) string {
	if s.SlotDisplayLocation != "" && s.SlotDisplayLocation != "Multiple" {
		return s.SlotDisplayLocation
	}
	if s.Location.Name != "" {
		return s.Location.Name
	}
	return s.SlotDisplayLocation
}

func bookableSlotsUniform(slots []client.BookableSlot) bool {
	if len(slots) < 2 {
		return false
	}
	s0 := slots[0]
	for _, s := range slots[1:] {
		if s.StartsAt.Format24Hour != s0.StartsAt.Format24Hour ||
			s.EndsAt.Format24Hour != s0.EndsAt.Format24Hour ||
			s.Price.Formatted != s0.Price.Formatted ||
			s.Spaces != s0.Spaces ||
			s.ActionToShow.Status != s0.ActionToShow.Status {
			return false
		}
	}
	return true
}

func Print(results []VenueResult, availableOnly bool) {
	for _, r := range results {
		dist := ""
		if r.Venue.Distance > 0 {
			dist = fmt.Sprintf(" (%.1f mi)", r.Venue.Distance)
		}
		fmt.Printf("--- %s%s ---\n", r.Venue.Name, dist)

		if r.Err != nil {
			fmt.Printf("  [error: %v]\n\n", r.Err)
			continue
		}

		slots := r.Slots
		if availableOnly {
			filtered := slots[:0]
			for _, s := range slots {
				if s.Spaces > 0 {
					filtered = append(filtered, s)
				}
			}
			slots = filtered
		}

		if len(slots) == 0 {
			fmt.Println("  No available slots")
			fmt.Println()
			continue
		}

		for _, s := range slots {
			start := s.StartsAt.Format24Hour
			end := s.EndsAt.Format24Hour
			spaces := spacesLabel(s.Spaces)
			cat := ""
			if s.Location != "" {
				cat = fmt.Sprintf("  [%s]", s.Location)
			}
			price := s.Price.FormattedAmount
			if price == "" {
				price = "—"
			}
			fmt.Printf("  %s–%s   %-8s  %s%s\n", start, end, price, spaces, cat)
		}
		fmt.Println()
	}
}

func spacesLabel(n int) string {
	switch {
	case n == 0:
		return "FULL"
	case n == 1:
		return "1 space"
	default:
		return fmt.Sprintf("%d spaces", n)
	}
}

// PrintCategories prints a numbered list of categories on stderr.
func PrintCategories(categories []client.Category) {
	fmt.Fprintln(os.Stderr, "Categories:")
	for i, cat := range categories {
		fmt.Fprintf(os.Stderr, "  %d. %s  [%s]\n", i+1, cat.Name, cat.Slug)
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 40))
}

// PrintActivities prints a numbered list of activities on stderr.
func PrintActivities(activities []client.Activity) {
	fmt.Fprintln(os.Stderr, "Activities:")
	for i, a := range activities {
		fmt.Fprintf(os.Stderr, "  %d. %s  [%s]\n", i+1, a.Name, a.Slug)
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 40))
}

// PrintVenues prints a numbered list of venues on stderr for interactive selection.
func PrintVenues(venues []client.Venue) {
	fmt.Fprintln(os.Stderr, "Venues found:")
	for i, v := range venues {
		dist := ""
		if v.Distance > 0 {
			dist = fmt.Sprintf(" (%.1f mi)", v.Distance)
		}
		fmt.Fprintf(os.Stderr, "  %d. %s%s  [%s]\n", i+1, v.Name, dist, v.Slug)
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 40))
}
