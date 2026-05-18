package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
	"github.com/ujjwaltalwar/betteruk-bot/internal/display"
)

var postcodeRe = regexp.MustCompile(`(?i)^[A-Z]{1,2}[0-9][0-9A-Z]?\s*[0-9][A-Z]{2}$`)

type flowStep int

const (
	stepVenue flowStep = iota
	stepCategory
	stepActivity
	stepDate
)

var rootCmd = &cobra.Command{
	Use:   "betteruk-bot",
	Short: "Check Better UK activity slot availability near a postcode",
	Example: `  betteruk-bot -p "N1 0SB"
  betteruk-bot -p SW1A1AA -c sports-hall-activities -a badminton-40min -d 2026-05-20
  betteruk-bot -p EC1A1BB --debug`,
	RunE: run,
}

var (
	flagPostcode      string
	flagCategory      string
	flagActivity      string
	flagDate          string
	flagMaxVenues     int
	flagAvailableOnly bool
	flagAuthToken     string
	flagDebug         bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagPostcode, "postcode", "p", "", "UK postcode to search near (required)")
	rootCmd.Flags().StringVarP(&flagCategory, "category", "c", "", "Category slug (skip category prompt)")
	rootCmd.Flags().StringVarP(&flagActivity, "activity", "a", "", "Activity slug (skip activity prompt)")
	rootCmd.Flags().StringVarP(&flagDate, "date", "d", "", "Date YYYY-MM-DD, today through 5 days ahead (skip prompt if set)")
	rootCmd.Flags().IntVarP(&flagMaxVenues, "max-venues", "n", 50, "Max venues to list for selection")
	rootCmd.Flags().BoolVar(&flagAvailableOnly, "available-only", true, "Show only slots with spaces > 0")
	rootCmd.Flags().StringVar(&flagAuthToken, "auth-token", "", "Bearer token for bookable courts API (or BETTER_AUTH_TOKEN)")
	rootCmd.Flags().BoolVar(&flagDebug, "debug", false, "Print raw HTTP responses to stderr")
	_ = rootCmd.MarkFlagRequired("postcode")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	postcode := strings.ToUpper(strings.TrimSpace(flagPostcode))
	if !postcodeRe.MatchString(postcode) {
		return fmt.Errorf("invalid UK postcode: %q", postcode)
	}

	c, err := client.New(flagDebug)
	if err != nil {
		return fmt.Errorf("init client: %w", err)
	}

	authToken := strings.TrimSpace(flagAuthToken)
	if authToken == "" {
		authToken = strings.TrimSpace(os.Getenv("BETTER_AUTH_TOKEN"))
	}
	if authToken != "" {
		c.SetAuthToken(authToken)
	}

	fmt.Fprintln(os.Stderr, "Fetching session and CSRF token...")
	if err := c.FetchCSRF(); err != nil {
		return fmt.Errorf("fetch CSRF: %w", err)
	}
	if authToken != "" {
		c.SetAuthToken(authToken)
	}

	fmt.Fprintf(os.Stderr, "Searching venues near %s...\n", postcode)
	venues, err := c.SearchVenues(postcode)
	if err != nil {
		return fmt.Errorf("venue search: %w", err)
	}
	if len(venues) == 0 {
		return fmt.Errorf("no venues found near %s", postcode)
	}
	if len(venues) > flagMaxVenues {
		venues = venues[:flagMaxVenues]
	}

	date := strings.TrimSpace(flagDate)
	if date != "" {
		if err := validateBookingDate(date, time.Now()); err != nil {
			return err
		}
	}

	return runInteractiveSession(c, venues, date)
}

func pickDate(date *string) error {
	newDate, err := promptDate(*date)
	if err != nil {
		return err
	}
	*date = newDate
	return nil
}

func runInteractiveSession(c *client.Client, venues []client.Venue, date string) error {
	step := stepVenue
	var venue client.Venue
	var categories []client.Category
	var categorySlug string
	var activities []client.Activity
	var activitySlug, activityName string

	for {
		switch step {
		case stepVenue:
			printCurrentDate(date)
			display.PrintVenues(venues)
			idx, err := promptChoice("Enter venue number", len(venues))
			if errors.Is(err, errDate) {
				if err := pickDate(&date); err != nil {
					return err
				}
				continue
			}
			if errors.Is(err, errBack) {
				fmt.Fprintln(os.Stderr, "Already at venue selection.")
				continue
			}
			if err != nil {
				return err
			}
			venue = venues[idx-1]

			if flagActivity != "" {
				activitySlug = flagActivity
				activityName = flagActivity
				step = stepDate
				continue
			}
			if flagCategory != "" {
				categorySlug = flagCategory
				step = stepActivity
				continue
			}
			step = stepCategory

		case stepCategory:
			printCurrentDate(date)
			fmt.Fprintf(os.Stderr, "Fetching categories at %s...\n", venue.Slug)
			var err error
			categories, err = c.GetCategories(venue.Slug)
			if err != nil {
				return err
			}
			if len(categories) == 0 {
				return fmt.Errorf("no categories at venue %s", venue.Slug)
			}

			display.PrintCategories(categories)
			idx, err := promptChoice("Enter category number", len(categories))
			if errors.Is(err, errDate) {
				if err := pickDate(&date); err != nil {
					return err
				}
				continue
			}
			if errors.Is(err, errBack) {
				step = stepVenue
				continue
			}
			if err != nil {
				return err
			}
			categorySlug = categories[idx-1].Slug
			step = stepActivity

		case stepActivity:
			printCurrentDate(date)
			fmt.Fprintf(os.Stderr, "Fetching activities in %s...\n", categorySlug)
			var err error
			activities, err = c.GetCategoryActivities(venue.Slug, categorySlug)
			if err != nil {
				return err
			}
			if len(activities) == 0 {
				return fmt.Errorf("no activities in category %q", categorySlug)
			}

			display.PrintActivities(activities)
			idx, err := promptChoice("Enter activity number", len(activities))
			if errors.Is(err, errDate) {
				if err := pickDate(&date); err != nil {
					return err
				}
				continue
			}
			if errors.Is(err, errBack) {
				if flagCategory != "" {
					step = stepVenue
				} else {
					step = stepCategory
				}
				continue
			}
			if err != nil {
				return err
			}
			selected := activities[idx-1]
			activitySlug = selected.Slug
			activityName = selected.Name
			step = stepDate

		case stepDate:
			if date == "" {
				if err := pickDate(&date); err != nil {
					return err
				}
			}

			again, err := showTimesAndPrompt(c, venue, activitySlug, activityName, &date)
			if err != nil {
				return err
			}
			if again {
				if flagActivity != "" {
					step = stepVenue
				} else {
					step = stepActivity
				}
				continue
			}
			return nil
		}
	}
}

func showTimesAndPrompt(c *client.Client, venue client.Venue, activitySlug, activityName string, date *string) (again bool, err error) {
	for {
		fmt.Fprintf(os.Stderr, "Fetching %s times at %s for %s...\n", activityName, venue.Name, *date)
		times, err := c.GetTimes(venue.Slug, activitySlug, *date)
		if err != nil {
			return false, err
		}

		filtered := times
		if flagAvailableOnly {
			filtered = filtered[:0]
			for _, t := range times {
				if t.Spaces > 0 {
					filtered = append(filtered, t)
				}
			}
		}

		fmt.Println()
		if len(filtered) == 0 {
			fmt.Printf("--- %s (%s) ---\n  No available times\n\n", venue.Name, *date)
		} else {
			display.PrintTimes(fmt.Sprintf("%s (%s)", venue.Name, *date), times, flagAvailableOnly)
		}

		for {
			action, choice, err := promptAfterTimes(len(filtered))
			if err != nil {
				return false, err
			}
			switch action {
			case "back":
				return true, nil
			case "quit":
				return false, nil
			case "date":
				if err := pickDate(date); err != nil {
					return false, err
				}
				break // refetch times for new date
			case "slot":
				t := filtered[choice-1]
				if c.AuthToken() == "" {
					fmt.Fprintln(os.Stderr, "Login required for bookable courts. Set BETTER_AUTH_TOKEN or --auth-token.")
					continue
				}
				fmt.Fprintf(os.Stderr, "Fetching bookable courts for %s–%s...\n",
					t.StartsAt.Format24Hour, t.EndsAt.Format24Hour)
				bookable, err := c.GetSlots(
					venue.Slug, activitySlug, *date,
					t.StartsAt.Format24Hour, t.EndsAt.Format24Hour, t.CompositeKey,
				)
				fmt.Println()
				display.PrintBookableSlots(bookable, err)
			}
		}
	}
}
