# betteruk-bot

A CLI to browse **Better UK** leisure centre availability: search venues by postcode, pick a category and activity, view time slots, and (when logged in) inspect bookable courts.

## Requirements

- Go 1.25+
- Network access to `better.org.uk` and `better-admin.org.uk`

## Install

```bash
git clone <repo-url>
cd betteruk-bot
go build -o betteruk-bot .
```

Or run without installing:

```bash
go run main.go -p "N7 8AN"
```

## Quick start

```bash
# Interactive session (postcode required)
betteruk-bot -p "N7 8AN"

# Skip prompts with flags
betteruk-bot -p SW1A1AA -c sports-hall-activities -a badminton-40min -d 2026-05-20

# Bookable courts (copy Bearer token from browser DevTools)
export BETTER_AUTH_TOKEN='v4.local....'
betteruk-bot -p "N7 8AN"
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--postcode` | `-p` | *(required)* | UK postcode to search near |
| `--category` | `-c` | | Category slug (skip category prompt) |
| `--activity` | `-a` | | Activity slug (skip activity prompt) |
| `--date` | `-d` | | Date `YYYY-MM-DD` (today through 5 days ahead; skip date prompt) |
| `--max-venues` | `-n` | `50` | Max venues shown in the selection list |
| `--available-only` | | `true` | Only show times with spaces &gt; 0 |
| `--auth-token` | | | Bearer token for the bookable courts API |
| `--debug` | | `false` | Print raw HTTP responses to stderr |

Environment variable: `BETTER_AUTH_TOKEN` (same as `--auth-token`).

## Interactive flow

After startup the bot:

1. Fetches a CSRF session from Better UK
2. Searches venues near your postcode
3. Guides you through numbered menus

```
Postcode search
    → Venue (1–N)
    → Category (1–N)     e.g. Sports Hall Activities
    → Activity (1–N)     e.g. Badminton 40min
    → Date (1–6)         today … +5 days (if not set via -d)
    → Times              available windows
    → [optional] Bookable courts for a chosen time
```

### Keyboard shortcuts

Available at **every** numbered prompt and after times are shown:

| Key | Action |
|-----|--------|
| `1`–`N` | Select list item |
| `d` | Set or change date |
| `b` | Go back one step |
| `q` | Quit (after times are displayed) |
| Enter | Accept default date when prompted |

When a date is selected it appears at the top of each step:

```
Date: Tue 20 May 2026  [2026-05-20]
```

### Date range

The Better API only accepts dates from **today through 5 days ahead**. The CLI enforces the same range. Dates further out return HTTP 422 from the API.

## API endpoints used

The bot mirrors the Better bookings website flow:

| Step | Method | Endpoint |
|------|--------|----------|
| Session | `GET` | `https://www.better.org.uk/what-we-offer/activities/badminton` |
| Venue search | `POST` | `https://www.better.org.uk/api/venue_searches` |
| Categories | `GET` | `https://better-admin.org.uk/api/activities/venue/{venue}/categories` |
| Activities | `GET` | `https://better-admin.org.uk/api/activities/venue/{venue}/categories/{category}` |
| Times | `GET` | `https://better-admin.org.uk/api/activities/venue/{venue}/activity/{activity}/times?date=…` |
| Bookable courts | `GET` | `https://better-admin.org.uk/api/activities/venue/{venue}/activity/{activity}/slots?date=…&start_time=…&end_time=…&composite_key=…` |

- **Times** — public; no login required
- **Slots** (bookable courts) — requires `Authorization: Bearer …` (see below)

## Authentication (bookable courts)

To drill into a specific time and see courts/resources (`BOOK` status, exact price, location):

1. Log in at [bookings.better.org.uk](https://bookings.better.org.uk)
2. Open browser DevTools → Network
3. Copy the `Bearer` token from any `better-admin.org.uk` request
4. Set it before running:

```bash
export BETTER_AUTH_TOKEN='v4.local....'
```

Without a token you can still browse venues, categories, activities, and times.

## Project layout

```
betteruk-bot/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Cobra CLI, interactive session loop
│   ├── prompt.go           # stdin prompts (choice, date, after-times menu)
│   └── date.go             # Allowed date range and validation
└── internal/
    ├── client/             # HTTP client and Better API
    └── display/            # Terminal output formatting
```

## Package reference

### `cmd`

| Function | Description |
|----------|-------------|
| `Execute()` | Run the Cobra root command |
| `run()` | Validate postcode, init client, search venues, start session |
| `runInteractiveSession()` | State machine: venue → category → activity → date → times |
| `showTimesAndPrompt()` | Fetch times, display menu, optional bookable-court lookup |
| `pickDate()` | Interactive date picker wrapper |
| `promptChoice()` | Numbered list input (`b` = back, `d` = date) |
| `promptDate()` | Date list input (today … +5 days) |
| `promptAfterTimes()` | Post-times menu (`1–N`, `d`, `b`, `q`) |
| `allowedBookingDates()` | Build slice of valid `YYYY-MM-DD` strings |
| `validateBookingDate()` | Validate `-d` flag value |
| `printCurrentDate()` | Show selected date on stderr |
| `formatDateLabel()` | Human-readable date label (`today`, `tomorrow`, etc.) |

### `internal/client`

| Type / function | Description |
|-----------------|-------------|
| `Client` | HTTP client with cookie jar, CSRF, and optional auth token |
| `New(debug)` | Create client (45s timeout, retry on venue search) |
| `FetchCSRF()` | Load badminton page, extract CSRF, warm bookings session |
| `SetAuthToken(token)` | Set Bearer token for `better-admin.org.uk` |
| `CSRFToken()` / `AuthToken()` | Access stored tokens |
| `SearchVenues(postcode)` | Find leisure centres near a postcode |
| `GetCategories(venueSlug)` | List activity categories at a venue |
| `GetCategoryActivities(venue, category)` | List activities in a category |
| `GetTimes(venue, activity, date)` | Available time windows (array or object JSON) |
| `GetSlots(venue, activity, date, start, end, compositeKey)` | Bookable courts for one time (login required) |
| `Venue` | `Slug`, `Name`, `Distance` |
| `Category` | `Slug`, `Name`, `HasChildren` |
| `Activity` | `Slug`, `Name` |
| `TimeSlot` | Start/end, price, spaces, location, `CompositeKey` |
| `BookableSlot` | Court details, price, `BOOK` status, capacity |

### `internal/display`

| Function | Description |
|----------|-------------|
| `PrintVenues()` | Numbered venue list (stderr) |
| `PrintCategories()` | Numbered category list (stderr) |
| `PrintActivities()` | Numbered activity list (stderr) |
| `PrintTimes()` | Numbered time windows with price and spaces (stdout) |
| `PrintBookableSlots()` | Bookable court details after slots API call |
| `Print()` | Legacy combined venue + slots output |

## Examples

```bash
# Full interactive flow
betteruk-bot -p "N1 0SB"

# Pre-select activity and date; only prompt for venue
betteruk-bot -p "N7 8AN" -c sports-hall-activities -a badminton-40min -d 2026-05-23

# Show all times including full ones
betteruk-bot -p "N7 8AN" --available-only=false

# Debug raw HTTP
betteruk-bot -p "N7 8AN" --debug
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o betteruk-bot .
```

## Notes

- Venue search responses are JavaScript/HTML payloads, not JSON; the client parses embedded HTML.
- Times API `data` may be a JSON **array** or an **object** keyed by court ID; both are handled.
- HTTP requests retry up to 3 times on timeout for venue search.
- This tool checks availability only; it does not complete bookings.
