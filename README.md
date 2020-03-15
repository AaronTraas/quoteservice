# Quote Catalog

This application is a basic web service that allows the creation and
moderation of a list of quotes. It stores the quotes in an in-memory DB. At
time of any moderation, it dumps the DB to a JSON file. It automatically loads
the JSON file on startup.

Also included is a Makefile, mostly for `make install`, which installs the app
as systemd service.

## Installation

Requirements:

- Go 1.14+
- Make
- Linux system with systemd (for installation)

### Instructions

1. Clone the repo
2. type `make run` to build and run in your terminal
3. type `make install` to build, install the systemd task, and start it

## Using the API

#### GET `/api/quotes`

Return the list of quotes as JSON.

#### POST `/api/quote/[ID]`

Parameters:
- `ID`: integer; unique id of quote from list of quotes.

If a quote with `id` exists, returns that quote.

#### POST `/api/submit/`

Parameters:
- `body`: main body of the quote. Assumed URL encoded
- `url`: url that contains the qutoe.

Create a new quote with `body` as the quote text and `uri` as the URI where
the quote is from. This creates a new quote where `approved` will be set to
`false`.

#### POST `/api/approve/?id=[ID]`

Parameters:
- `id`: integer; unique id of quote from list of quotes.

If a quote with `id` exists, set `approved` to `true`

#### POST `/api/disapprove/?id=[ID]`

Parameters:
- `id`: integer; unique id of quote from list of quotes.

If a quote with `id` exists, set `approved` to `true`

There is *no authentication* for any endpoint. When I set it up on the test
server, I put nginx in front of it, a restricted `/api/approve` and
`/api/disapprove` behind a HTTP Basic Auth password.