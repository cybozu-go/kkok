[![GitHub release](https://img.shields.io/github/release/cybozu-go/kkok.svg?maxAge=60)][releases]
[![GoDoc](https://godoc.org/github.com/cybozu-go/kkok?status.svg)][godoc]
[![CircleCI](https://circleci.com/gh/cybozu-go/kkok.svg?style=svg)](https://circleci.com/gh/cybozu-go/kkok)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/kkok)](https://goreportcard.com/report/github.com/cybozu-go/kkok)

**kkok** (taken from Japanese word 警告 - *keikoku* -, in English *alert*) is a service to process alerts nicely.  It gathers alerts from miscellaneous sources, applies filters to edit or route them, then sends the processed alerts via email, SMS (Twilio), Slack, etc.

Architecture
------------

kkok sends alerts through these steps:

1. Generate alerts from sources.
2. Collect and pool alerts for some duration.
3. Edit/route collected alerts by filters.
4. Send alerts along with the given routes.

Please read [Architecture.md](docs/Architecture.md) for more details.

Features
--------

* Generators:

    * HTTP REST API.
    * `maildir`: generate alerts from mails in a [Maildir][] directory.

* Filters:

    * `freq`: calculate and add frequency information to alerts.
    * `discard`: discard alerts based on the given conditions.
    * `group`: merge alerts into groups by field values.
    * `route`: add or replace routes to alert receivers.
    * `edit`: edit alerts by JavaScript.
    * `exec`: invoke an external command to edit alerts.

* Transports:

    * `email`: format and send alerts via email.
    * `slack`: format and send alerts to a [Slack][] channel.
    * `twilio`: format and send SMS via [Twilio][].
    * `exec`: invoke an external command to send alerts.

Build
-----

Use Go 1.7 or better.

Run the command below exactly as shown, including the ellipsis.
They are significant - see `go help packages`.

```
go get -u github.com/cybozu-go/kkok/...
```

Usage
-----

Read [Usage.md](docs/Usage.md) and [API.md](docs/API.md).

License
-------

[MIT][]

Authors & Contributors
----------------------

* Yamamoto, Hirotaka [@ymmt2005](https://github.com/ymmt2005)

[releases]: https://github.com/cybozu-go/kkok/releases
[godoc]: https://godoc.org/github.com/cybozu-go/kkok
[Maildir]: https://en.wikipedia.org/wiki/Maildir
[Twilio]: https://www.twilio.com/
[Slack]: https://slack.com/
[MIT]: https://opensource.org/licenses/MIT
