package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ujjwaltalwar/betteruk-bot/internal/client"
)

func printAuthStatus(c *client.Client) {
	if c.AuthToken() == "" {
		fmt.Fprintln(os.Stderr, "Auth: not set (press a to paste token)")
		return
	}
	t := c.AuthToken()
	summary := t
	if len(t) > 16 {
		summary = t[:10] + "…" + t[len(t)-6:]
	}
	fmt.Fprintf(os.Stderr, "Auth: set (%s)\n", summary)
}

func promptSetAuthToken(c *client.Client) error {
	fmt.Fprintln(os.Stderr, "Paste Bearer token (DevTools → Network → Authorization, or BETTER_AUTH_TOKEN value):")
	fmt.Fprint(os.Stderr, "Token: ")
	if !stdin.Scan() {
		if err := stdin.Err(); err != nil {
			return fmt.Errorf("read auth token: %w", err)
		}
		return fmt.Errorf("no token entered")
	}
	token := strings.TrimSpace(stdin.Text())
	if token == "" {
		return fmt.Errorf("no token entered")
	}
	c.SetAuthToken(token)
	fmt.Fprintln(os.Stderr, "Auth token saved for this session.")
	printAuthStatus(c)
	return nil
}

// handleChoicePromptErr handles d=date and a=auth from numbered prompts; returns true if handled.
func handleChoicePromptErr(err error, c *client.Client, date *string) (bool, error) {
	if errors.Is(err, errDate) {
		if err := pickDate(date); err != nil {
			return true, err
		}
		return true, nil
	}
	if errors.Is(err, errAuth) {
		if err := promptSetAuthToken(c); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}
