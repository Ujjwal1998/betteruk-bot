package cmd

import (
	"bufio"
	"errors"
	"os"
	"testing"
	"time"
)

func TestPromptChoiceBack(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("b\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	_, err = promptChoice("test", 3)
	if !errors.Is(err, errBack) {
		t.Fatalf("got err %v, want errBack", err)
	}
}

func TestPromptChoiceNumber(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("2\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	n, err := promptChoice("test", 3)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("got %d, want 2", n)
	}
}

func TestPromptDateDefault(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	want := time.Now().Format("2006-01-02")
	got, err := promptDate("")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPromptDateChoice(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("3\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	want := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	got, err := promptDate("")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPromptAfterSearchResultsAuthAction(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("a\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	action, _, err := promptAfterSearchResults(3)
	if err != nil {
		t.Fatal(err)
	}
	if action != "auth" {
		t.Fatalf("got action %q, want auth", action)
	}
}

func TestPromptAfterTimesAuthAction(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("a\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	action, _, err := promptAfterTimes(3)
	if err != nil {
		t.Fatal(err)
	}
	if action != "auth" {
		t.Fatalf("got action %q, want auth", action)
	}
}

func TestPromptAfterTimesDateAction(t *testing.T) {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r
	stdin = bufio.NewScanner(os.Stdin)
	t.Cleanup(func() {
		os.Stdin = oldStdin
		stdin = bufio.NewScanner(os.Stdin)
	})

	if _, err := w.WriteString("d\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	action, _, err := promptAfterTimes(3)
	if err != nil {
		t.Fatal(err)
	}
	if action != "date" {
		t.Fatalf("got action %q, want date", action)
	}
}
