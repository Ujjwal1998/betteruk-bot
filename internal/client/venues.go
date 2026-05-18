package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	venueSearchURL  = baseURL + "/api/venue_searches"
	dynamicPanelID  = "13629" // badminton activity widget
	activityTypeIDs = "7"     // 7 = badminton
)

var (
	venueSlugRe = regexp.MustCompile(`bookings\.better\.org\.uk/location/([a-z0-9-]+)`)
	venueNameRe = regexp.MustCompile(`facility-finder__result-link\\?"[^>]*>([^<]+)<`)
)

type Venue struct {
	Slug     string  `json:"slug"`
	Name     string  `json:"name"`
	Distance float64 `json:"distance"`
}

type venueSearchResponse struct {
	Venues []Venue `json:"venues"`
}

// SearchVenues posts a venue search and returns venues in result order.
func (c *Client) SearchVenues(postcode string) ([]Venue, error) {
	postcode = strings.ToUpper(strings.TrimSpace(postcode))

	form := url.Values{}
	form.Set("utf8", "✓")
	form.Set("venue_search[dynamic_panel_id]", dynamicPanelID)
	form.Set("venue_search[activity_type_ids]", activityTypeIDs)
	form.Set("venue_search[searchterm]", postcode)
	// Rails checkbox pattern: send 0 (hidden) then 1 (checked)
	form.Add("venue_search[only_open]", "0")
	form.Add("venue_search[only_open]", "1")
	form.Set("commit", "Search")

	req, err := http.NewRequest("POST", venueSearchURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-CSRF-Token", c.csrfToken)
	req.Header.Set("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")

	resp, err := c.doRetry(req, csrfPageURL, true)
	if err != nil {
		return nil, fmt.Errorf("venue search request: %w", err)
	}
	defer resp.Body.Close()

	body, err := c.readBody(resp)
	if err != nil {
		return nil, err
	}

	venues, err := parseVenueSearchBody(body)
	if err != nil {
		return nil, fmt.Errorf("parse venue search response: %w", err)
	}
	return venues, nil
}

func parseVenueSearchBody(body []byte) ([]Venue, error) {
	var jsonResult venueSearchResponse
	if err := json.Unmarshal(body, &jsonResult); err == nil && len(jsonResult.Venues) > 0 {
		return jsonResult.Venues, nil
	}

	venues := parseVenueSearchHTML(string(body))
	if len(venues) == 0 {
		snippet := body
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return nil, fmt.Errorf("no venues in response: %s", snippet)
	}
	return venues, nil
}

// parseVenueSearchHTML extracts venues from the JavaScript innerHTML payload.
func parseVenueSearchHTML(body string) []Venue {
	chunks := strings.Split(body, `facility-finder__result `)
	if len(chunks) < 2 {
		chunks = strings.Split(body, `facility-finder__result\"`)
	}
	if len(chunks) < 2 {
		return nil
	}

	seen := make(map[string]struct{})
	var venues []Venue
	for _, chunk := range chunks[1:] {
		if !strings.Contains(chunk, "bookings.better.org.uk/location/") {
			continue
		}

		slugMatch := venueSlugRe.FindStringSubmatch(chunk)
		if len(slugMatch) < 2 {
			continue
		}
		slug := slugMatch[1]
		if _, ok := seen[slug]; ok {
			continue
		}

		name := slug
		if nameMatch := venueNameRe.FindStringSubmatch(chunk); len(nameMatch) >= 2 {
			name = html.UnescapeString(nameMatch[1])
		}

		seen[slug] = struct{}{}
		venues = append(venues, Venue{Slug: slug, Name: name})
	}
	return venues
}
