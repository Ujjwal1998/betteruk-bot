package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	errBack = errors.New("back")
	errDate = errors.New("date")
)

var stdin = bufio.NewScanner(os.Stdin)

func promptChoice(label string, count int) (int, error) {
	for {
		fmt.Fprintf(os.Stderr, "%s (1-%d, d=date, b=back): ", label, count)
		if !stdin.Scan() {
			if err := stdin.Err(); err != nil {
				return 0, fmt.Errorf("read choice: %w", err)
			}
			return 0, fmt.Errorf("no choice entered")
		}

		input := strings.TrimSpace(stdin.Text())
		switch strings.ToLower(input) {
		case "b":
			return 0, errBack
		case "d":
			return 0, errDate
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > count {
			fmt.Fprintf(os.Stderr, "Enter a number from 1 to %d, d, or b.\n", count)
			continue
		}
		return choice, nil
	}
}

func promptDate(current string) (string, error) {
	now := time.Now()
	dates := allowedBookingDates(now)

	defaultIdx := 1
	if current != "" {
		for i, d := range dates {
			if d == current {
				defaultIdx = i + 1
				break
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Dates (today through %d days ahead):\n", maxBookingDaysAhead)
	for i, d := range dates {
		fmt.Fprintf(os.Stderr, "  %d. %s  [%s]\n", i+1, formatDateLabel(d, now), d)
	}
	fmt.Fprintln(os.Stderr, strings.Repeat("-", 40))

	for {
		fmt.Fprintf(os.Stderr, "Enter date number (1-%d) [%d]: ", len(dates), defaultIdx)
		if !stdin.Scan() {
			if err := stdin.Err(); err != nil {
				return "", fmt.Errorf("read date: %w", err)
			}
			return "", fmt.Errorf("no date entered")
		}

		input := strings.TrimSpace(stdin.Text())
		if input == "" {
			return dates[defaultIdx-1], nil
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(dates) {
			fmt.Fprintf(os.Stderr, "Enter a number from 1 to %d, or press Enter for [%d].\n", len(dates), defaultIdx)
			continue
		}
		return dates[choice-1], nil
	}
}

func promptAfterTimes(count int) (action string, choice int, err error) {
	for {
		if count == 0 {
			fmt.Fprint(os.Stderr, "(d=date, b=back, q=quit): ")
		} else {
			fmt.Fprintf(os.Stderr, "(1-%d=bookable courts, d=date, b=back, q=quit): ", count)
		}
		if !stdin.Scan() {
			if err := stdin.Err(); err != nil {
				return "", 0, fmt.Errorf("read choice: %w", err)
			}
			return "", 0, fmt.Errorf("no choice entered")
		}

		input := strings.TrimSpace(stdin.Text())
		switch strings.ToLower(input) {
		case "b":
			return "back", 0, nil
		case "q":
			return "quit", 0, nil
		case "d":
			return "date", 0, nil
		}

		if count == 0 {
			fmt.Fprintln(os.Stderr, "Enter d, b, or q.")
			continue
		}
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > count {
			fmt.Fprintf(os.Stderr, "Enter a number from 1 to %d, d, b, or q.\n", count)
			continue
		}
		return "slot", choice, nil
	}
}
