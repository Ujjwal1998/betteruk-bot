package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
)

func normalizePostcode(postcode string) (string, error) {
	postcode = strings.ToUpper(strings.TrimSpace(postcode))
	if !postcodeRe.MatchString(postcode) {
		return "", fmt.Errorf("invalid UK postcode: %q", postcode)
	}
	return postcode, nil
}

func resolveAuthToken(flagToken string) string {
	token := strings.TrimSpace(flagToken)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("BETTER_AUTH_TOKEN"))
	}
	return token
}

func initClient(debug bool, authToken string) (*client.Client, error) {
	c, err := client.New(debug)
	if err != nil {
		return nil, fmt.Errorf("init client: %w", err)
	}
	if authToken != "" {
		c.SetAuthToken(authToken)
	}
	fmt.Fprintln(os.Stderr, "Fetching session and CSRF token...")
	if err := c.FetchCSRF(); err != nil {
		return nil, fmt.Errorf("fetch CSRF: %w", err)
	}
	if authToken != "" {
		c.SetAuthToken(authToken)
	}
	return c, nil
}

func fetchVenuesNear(c *client.Client, postcode string, limit int) ([]client.Venue, error) {
	fmt.Fprintf(os.Stderr, "Searching venues near %s...\n", postcode)
	venues, err := c.SearchVenues(postcode)
	if err != nil {
		return nil, fmt.Errorf("venue search: %w", err)
	}
	if len(venues) == 0 {
		return nil, fmt.Errorf("no venues found near %s", postcode)
	}
	if limit > 0 && len(venues) > limit {
		venues = venues[:limit]
	}
	return venues, nil
}

func validateDateFlag(date string) error {
	if date == "" {
		return nil
	}
	return validateBookingDate(date, time.Now())
}
