package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
	"github.com/ujjwaltalwar/betteruk-bot/internal/display"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Find activity availability across venues near a postcode",
	Example: `  betteruk-bot search -p "N7 8AN" -a badminton-60min -d 2026-05-23
  betteruk-bot search -p N1C4EP`,
	RunE: runSearch,
}

var (
	searchPostcode      string
	searchActivity      string
	searchDate          string
	searchScanVenues    int
	searchAvailableOnly bool
	searchAuthToken     string
	searchDebug         bool
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchPostcode, "postcode", "p", "", "UK postcode to search near (required)")
	searchCmd.Flags().StringVarP(&searchActivity, "activity", "a", "", "Activity slug (interactive catalog if omitted)")
	searchCmd.Flags().StringVarP(&searchDate, "date", "d", "", "Date YYYY-MM-DD, today through 5 days ahead")
	searchCmd.Flags().IntVar(&searchScanVenues, "scan-venues", 10, "Number of nearby venues to scan for availability")
	searchCmd.Flags().BoolVar(&searchAvailableOnly, "available-only", true, "Only include times with spaces > 0")
	searchCmd.Flags().StringVar(&searchAuthToken, "auth-token", "", "Bearer token for bookable courts API (or BETTER_AUTH_TOKEN)")
	searchCmd.Flags().BoolVar(&searchDebug, "debug", false, "Print raw HTTP responses to stderr")
	_ = searchCmd.MarkFlagRequired("postcode")
}

func runSearch(cmd *cobra.Command, args []string) error {
	postcode, err := normalizePostcode(searchPostcode)
	if err != nil {
		return err
	}

	date := searchDate
	if err := validateDateFlag(date); err != nil {
		return err
	}

	authToken := resolveAuthToken(searchAuthToken)
	c, err := initClient(searchDebug, authToken)
	if err != nil {
		return err
	}

	venues, err := fetchVenuesNear(c, postcode, searchScanVenues)
	if err != nil {
		return err
	}

	activitySlug, activityName, err := resolveActivity(searchActivity)
	if err != nil {
		return err
	}

	if date == "" {
		if err := pickDate(&date); err != nil {
			return err
		}
	}

	return runSearchSession(c, venues, activitySlug, activityName, date, searchAvailableOnly)
}

func runSearchSession(c *client.Client, venues []client.Venue, activitySlug, activityName, date string, availableOnly bool) error {
	for {
		printCurrentDate(date)
		printAuthStatus(c)
		fmt.Fprintf(os.Stderr, "Scanning %d venues for %s on %s...\n", len(venues), activityName, date)

		results, err := c.SearchTimesAcrossVenues(venues, activitySlug, date, client.DefaultSearchWorkers(), availableOnly)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, client.FormatSearchSummary(len(results), len(venues), activityName, date))

		rows := display.BuildSearchRows(results, activitySlug, date)
		fmt.Println()
		display.PrintSearchResults(activityName, date, rows)

		if err := promptSearchResults(c, rows, &date); err != nil {
			if errors.Is(err, errRestartScan) {
				continue
			}
			return err
		}
		return nil
	}
}

func promptSearchResults(c *client.Client, rows []display.SearchRow, date *string) error {
	for {
		action, choice, err := promptAfterSearchResults(len(rows))
		if err != nil {
			return err
		}
		switch action {
		case "quit":
			return nil
		case "date":
			if err := pickDate(date); err != nil {
				return err
			}
			return errRestartScan
		case "auth":
			if err := promptSetAuthToken(c); err != nil {
				return err
			}
			continue
		case "slot":
			if len(rows) == 0 {
				continue
			}
			row := rows[choice-1]
			if c.AuthToken() == "" {
				fmt.Fprintln(os.Stderr, "Login required for bookable courts. Press a to paste token, or set BETTER_AUTH_TOKEN.")
				continue
			}
			fmt.Fprintf(os.Stderr, "Fetching bookable courts at %s for %s–%s...\n",
				row.Venue.Name, row.Slot.StartsAt.Format24Hour, row.Slot.EndsAt.Format24Hour)
			bookable, err := c.GetSlots(
				row.Venue.Slug, row.ActivitySlug, row.Date,
				row.Slot.StartsAt.Format24Hour, row.Slot.EndsAt.Format24Hour, row.Slot.CompositeKey,
			)
			fmt.Println()
			display.PrintBookableSlotsForVenue(row.Venue.Name, bookable, err)
		}
	}
}
