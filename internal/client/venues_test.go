package client

import "testing"

const sampleVenueSearchJS = `var elmId = 'venue-search-results-13629';
var elm = document.getElementById(elmId);
elm.innerHTML = "<section class=\"facility-finder__results\">\n      <article class=\"facility-finder__result \">\n        <a class=\"facility-finder__result-link\" href=\"/leisure-centre/london/camden/kings-cross-fitness\">King&#39;s Cross Fitness</a>\n      <a class=\"call-to-action call-to-action--primary call-to-action--join\" href=\"https://bookings.better.org.uk/location/kings-cross-fitness\">Book activity</a>\n  </article>\n\n      <article class=\"facility-finder__result \">\n        <a class=\"facility-finder__result-link\" href=\"/leisure-centre/london/islington/sobell\">Sobell Leisure Centre</a>\n      <a class=\"call-to-action call-to-action--primary call-to-action--join\" href=\"https://bookings.better.org.uk/location/sobell-leisure-centre\">Book activity</a>\n  </article>\n</section>";`

func TestParseVenueSearchHTML(t *testing.T) {
	venues := parseVenueSearchHTML(sampleVenueSearchJS)
	if len(venues) != 2 {
		t.Fatalf("got %d venues, want 2", len(venues))
	}
	if venues[0].Slug != "kings-cross-fitness" {
		t.Errorf("venue[0].Slug = %q, want kings-cross-fitness", venues[0].Slug)
	}
	if venues[0].Name != "King's Cross Fitness" {
		t.Errorf("venue[0].Name = %q, want King's Cross Fitness", venues[0].Name)
	}
	if venues[1].Slug != "sobell-leisure-centre" {
		t.Errorf("venue[1].Slug = %q, want sobell-leisure-centre", venues[1].Slug)
	}
	if venues[1].Name != "Sobell Leisure Centre" {
		t.Errorf("venue[1].Name = %q, want Sobell Leisure Centre", venues[1].Name)
	}
}

func TestParseVenueSearchBody_JSON(t *testing.T) {
	body := []byte(`{"venues":[{"slug":"sobell-leisure-centre","name":"Sobell","distance":1.2}]}`)
	venues, err := parseVenueSearchBody(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(venues) != 1 || venues[0].Slug != "sobell-leisure-centre" {
		t.Fatalf("unexpected venues: %+v", venues)
	}
}
