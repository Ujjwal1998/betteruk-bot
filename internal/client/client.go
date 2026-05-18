package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

const (
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:150.0) Gecko/20100101 Firefox/150.0"
	csrfPageURL  = "https://www.better.org.uk/what-we-offer/activities/badminton"
	baseURL      = "https://www.better.org.uk"
	bookingsURL  = "https://bookings.better.org.uk"
	bookingsRoot = bookingsURL + "/"
	adminURL        = "https://better-admin.org.uk"
	defaultTimeout  = 45 * time.Second
	maxRequestRetry = 3
)

type Client struct {
	http      *http.Client
	jar       *cookiejar.Jar
	csrfToken string
	authToken string
	debug     bool
}

func New(debug bool) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	return &Client{
		http: &http.Client{
			Jar:     jar,
			Timeout: defaultTimeout,
		},
		jar:   jar,
		debug: debug,
	}, nil
}

// FetchCSRF loads the badminton page, extracts the CSRF token from the HTML
// meta tag, and captures the auth token from the cookie jar.
func (c *Client) FetchCSRF() error {
	req, err := http.NewRequest("GET", csrfPageURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := c.readBody(resp)
	if err != nil {
		return err
	}

	csrf, err := extractCSRF(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("extract CSRF token: %w", err)
	}
	c.csrfToken = csrf

	if err := c.warmBookingsSession(); err != nil {
		return fmt.Errorf("warm bookings session: %w", err)
	}

	if c.authToken == "" {
		c.authToken = extractAuthToken(c.jar)
	}
	return nil
}

func (c *Client) warmBookingsSession() error {
	req, err := http.NewRequest("GET", bookingsRoot, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req, bookingsRoot)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = c.readBody(resp)
	return err
}

func extractAuthToken(jar *cookiejar.Jar) string {
	for _, origin := range []string{baseURL, bookingsURL, adminURL} {
		u, err := url.Parse(origin)
		if err != nil {
			continue
		}
		for _, cookie := range jar.Cookies(u) {
			if cookie.Name == "better.org.uk-authToken" || strings.HasSuffix(cookie.Name, "-authToken") {
				return strings.Trim(cookie.Value, `"`)
			}
		}
	}
	return ""
}

// SetAuthToken sets the Bearer token for better-admin.org.uk API calls.
func (c *Client) SetAuthToken(token string) {
	c.authToken = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
}

func (c *Client) CSRFToken() string  { return c.csrfToken }
func (c *Client) AuthToken() string   { return c.authToken }

func (c *Client) do(req *http.Request, referer string) (*http.Response, error) {
	return c.doRetry(req, referer, false)
}

func (c *Client) doRetry(req *http.Request, referer string, retry bool) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
	}

	attempts := 1
	if retry {
		attempts = maxRequestRetry
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}
		if attempt > 1 {
			time.Sleep(time.Duration(attempt-1) * time.Second)
		}

		resp, err := c.http.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if !retry || !isRetryableHTTPError(err) || attempt == attempts {
			return nil, err
		}
		if c.debug {
			fmt.Fprintf(os.Stderr, "[debug] retry %d/%d after: %v\n", attempt, attempts, err)
		}
	}
	return nil, lastErr
}

func isRetryableHTTPError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "EOF") ||
		strings.Contains(msg, "temporary failure")
}

// setBookingsCORSHeaders mirrors the browser request to better-admin.org.uk.
func setBookingsCORSHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", bookingsURL)
	req.Header.Set("Referer", bookingsRoot)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
}

func (c *Client) readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if c.debug {
		fmt.Fprintf(os.Stderr, "\n[debug] %s %s → %d\n%s\n",
			resp.Request.Method, resp.Request.URL, resp.StatusCode, string(body))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := body
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return nil, fmt.Errorf("HTTP %d from %s: %s", resp.StatusCode, resp.Request.URL, snippet)
	}
	return body, nil
}

func extractCSRF(r io.Reader) (string, error) {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.SelfClosingTagToken || tt == html.StartTagToken {
			tok := z.Token()
			if tok.Data != "meta" {
				continue
			}
			var name, content string
			for _, a := range tok.Attr {
				switch a.Key {
				case "name":
					name = a.Val
				case "content":
					content = a.Val
				}
			}
			if name == "csrf-token" && content != "" {
				return content, nil
			}
		}
	}
	return "", fmt.Errorf("csrf-token meta tag not found in page")
}
