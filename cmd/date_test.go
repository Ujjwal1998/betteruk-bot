package cmd

import (
	"testing"
	"time"
)

func TestAllowedBookingDates(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	dates := allowedBookingDates(now)
	if len(dates) != 6 {
		t.Fatalf("got %d dates, want 6", len(dates))
	}
	if dates[0] != "2026-05-18" || dates[5] != "2026-05-23" {
		t.Fatalf("unexpected range: %v", dates)
	}
}

func TestValidateBookingDate(t *testing.T) {
	now := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	if err := validateBookingDate("2026-05-23", now); err != nil {
		t.Fatalf("expected valid: %v", err)
	}
	if err := validateBookingDate("2026-05-24", now); err == nil {
		t.Fatal("expected error for date too far ahead")
	}
	if err := validateBookingDate("2026-05-17", now); err == nil {
		t.Fatal("expected error for past date")
	}
}
