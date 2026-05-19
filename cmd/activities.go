package cmd

import (
	"fmt"
	"os"
	"strings"
)

type catalogActivity struct {
	Name string
	Slug string
}

var activityCatalog = []catalogActivity{
	{Name: "Badminton 40min", Slug: "badminton-40min"},
	{Name: "Badminton 60min", Slug: "badminton-60min"},
	{Name: "Basketball Half Court 40min", Slug: "basketball-half-court-40min"},
	{Name: "Basketball Half Court 60min", Slug: "basketball-half-court-60min"},
	{Name: "Football Indoor 60min", Slug: "football-indoor-60min"},
	{Name: "Netball Indoor 40min", Slug: "netball-indoor-40min"},
	{Name: "Netball Indoor 60min", Slug: "netball-indoor-60min"},
	{Name: "Pickleball 40mins", Slug: "pickleball-40mins"},
	{Name: "Pickleball 60mins", Slug: "pickleball-60mins"},
	{Name: "Table Tennis 60min", Slug: "table-tennis-60min"},
}

func printActivityCatalog() {
	fmt.Fprintln(os.Stderr, "Activities:")
	for i, a := range activityCatalog {
		fmt.Fprintf(os.Stderr, "  %d. %s  [%s]\n", i+1, a.Name, a.Slug)
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 40))
}

func resolveActivity(slugFlag string) (slug, name string, err error) {
	if slugFlag != "" {
		return slugFlag, slugFlag, nil
	}
	printActivityCatalog()
	idx, err := promptChoice("Enter activity number", len(activityCatalog))
	if err != nil {
		return "", "", err
	}
	selected := activityCatalog[idx-1]
	return selected.Slug, selected.Name, nil
}
