/*
Package maildir reads mails in a Maildir directory to generate alerts.

Maildir is a mail spool format first implemented in qmail:
https://cr.yp.to/proto/maildir.html

This plugin scans a Maildir directory repeatedly at regular intervals.
Specifically, mails in "new" directory will be scanned, processed, then
removed by the plugin.

Construction parameters:

    Name        Type               Default       Description
    dir         string                           Absolute path to a Maildir directory.
    interval    int                10            Scanning interval (seconds).

Example snippet for TOML configuration:

    [[source]]
    type     = "maildir"
    dir      = "/var/mail/kkok"
    interval = 60

Alerts are generated from mail headers and headers-in-mail-body.
Specifically, "From" is taken from the mail's From header value,
"Date" is taken from the mail's Date header, "Title" is taken from
the mail's Subject header value, "Message" is taken from the mail
body text.

Headers-in-mail-body are pseudo headers written in mail body.
"Host" and other fields can be given by these pseudo headers, like this:
Pseudo headers that do not match kkok.Alert field names become members of
"Info" field.

    From: foobar <foobar@example.com>
    Subject: alert from zoo monitor

    Host: east-asian-zoo                               <- pseudo header
    From: east asian zoo operator                      <- ditto
    Option1: 123
    Option2: abc def

    Message body.

Fields specified by pseudo headers take precedence over normal headers.
In the above example, "From" field value will be "east asian zoo operator",
and "Info" field value will be a map {"Option1":123,"Option2":"abc def"}.
*/
package maildir
