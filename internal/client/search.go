package client

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

const defaultSearchWorkers = 5

// DefaultSearchWorkers returns the concurrent worker count for venue scans.
func DefaultSearchWorkers() int { return defaultSearchWorkers }

// VenueTime is one available time window at a specific venue.
type VenueTime struct {
	Venue Venue
	Slot  TimeSlot
}

// SearchTimesAcrossVenues fetches times for an activity at each venue concurrently.
// Per-venue errors are skipped; venues with no times are omitted.
func (c *Client) SearchTimesAcrossVenues(venues []Venue, activity, date string, workers int, availableOnly bool) ([]VenueTime, error) {
	if workers < 1 {
		workers = defaultSearchWorkers
	}

	jobs := make(chan Venue)
	var (
		mu      sync.Mutex
		results []VenueTime
		wg      sync.WaitGroup
	)

	worker := func() {
		defer wg.Done()
		for venue := range jobs {
			times, err := c.GetTimes(venue.Slug, activity, date)
			if err != nil {
				continue
			}
			for _, slot := range times {
				if availableOnly && slot.Spaces <= 0 {
					continue
				}
				mu.Lock()
				results = append(results, VenueTime{Venue: venue, Slot: slot})
				mu.Unlock()
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	for _, v := range venues {
		jobs <- v
	}
	close(jobs)
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		ti := results[i].Slot.StartsAt.Format24Hour
		tj := results[j].Slot.StartsAt.Format24Hour
		if ti != tj {
			return ti < tj
		}
		return strings.Compare(results[i].Venue.Name, results[j].Venue.Name) < 0
	})

	return results, nil
}

// FormatSearchSummary returns a short summary line for stderr.
func FormatSearchSummary(found, scanned int, activity, date string) string {
	return fmt.Sprintf("Found %d available slot(s) across %d venue(s) for %s on %s", found, scanned, activity, date)
}
