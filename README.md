# Quote Catalog

I needed to learn [Go](https://golang.org/) for an upcoming much more complicated project, so  I thought I'd do something simple to learn how to build a basic Go web service.

The service has some crud-like functions for an in-memory data store, including a moderation queue. It dumps the contents of the DB on approval, and loads it on start.

Also included is a Makefile, mostly for `make install`, which installs the app as systemd service.

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

#### GET `/api/submit/?body=[BODY]&uri=[URI]`
Create a new quote with `body` as the quote text and `uri` as the URI where the quote is from. This creates a new quote where `approved` will be set to `false`.

#### GET `/api/approve/?id=[ID]`
If a quote with `id` exists, set `approved` to `true`

#### GET `/api/disapprove/?id=[ID]`
If a quote with `id` exists, set `approved` to `true`

There is *no authentication* for any endpoint. When I set it up on the test server, I put nginx in front of it, a restricted `/api/approve` and `/api/disapprove` behind a HTTP Basic Auth password.