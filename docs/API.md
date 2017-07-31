REST API Specifications
=======================

This document describes kkok REST API.

Common
------

### Content-Type

Requests that send JSON by POST or PUT methods must supply

    Content-Type: application/json

header.

Response body is formatted in JSON *if and only if* the response header's
Content-Type is "application/json".  Response body may be formatted
in "text/plain" or "text/html".

### Authentication

kkok can be configured to require authentication token for API usage.
If configured, clients need to send the token with

    Authorization: Bearer TOKEN

header.  Bearer scheme is defined in [RFC6750][]

Unless otherwise noted, all APIs require the token if configured so.

### Status code

When a request goes successful, the HTTP response status will be 200.
Other statuses indicate errors.

For APIs taking an ID, the HTTP status will be 404 if ID is not found.

Normally, an error response's content-type is "text/plain" and the
response body contains error messages in UTF-8.

### Method override

For clients that can send only GET and POST requests, kkok provides
HTTP method override by a special HTTP header `X-HTTP-Method-Override`.

Kkok regards POST requests with the header as if requests of the
method specified by the header value.  For instance, the following
request is handled as `DELETE` instead of `POST`.

```
POST /filters/ID HTTP/1.1
Host: hostname.domain:80
Content-Length: 0
X-HTTP-Method-Override: DELETE

```

API List
--------

Each subsection describes the method and the end point (URL) of an API.

* [GET /version](#get-version)
* [GET /alerts](#get-alerts)
* [POST /alerts](#post-alerts)
* [GET /filters](#get-filters)
* [PUT /filters/ID](#put-filtersid)
* [GET /filters/ID](#get-filtersid)
* [DELETE /filters/ID](#delete-filtersid)
* [PUT /filters/ID/enable](#put-filtersidenable)
* [PUT /filters/ID/disable](#put-filtersiddisable)
* [PUT /filters/ID/inactivate](#put-filtersidinactivate)
* [GET /routes](#get-routes)
* [PUT /routes/ID](#put-routesid)
* [GET /routes/ID](#get-routesid)

### GET /version

Return version string like `0.1.1` as "text/plain".

This API does *not* require the authentication token.

### GET /alerts

Return a JSON array of pending alert objects.

### POST /alerts

* Content-Type: application/json
* Body: see below.

Post a new alert.  The JSON must be an object with these fields:

| Name | Required | Type | Description |
| ---- | -------- | ---- | ----------- |
| `From` | Yes | string | Who sent this alert. |
| `Title` | Yes | string | One-line description of the alert. |
| `Date` | No | string | RFC3339 format date string. |
| `Host` | No | string | Where this alert was generated. |
| `Message` | No | string | Multi-line description of the alert. |
| `Info` | No | object | Additional fields. |

If `Date` is omitted, the current date will be used for the alert.

If `Host` is omitted, the request client's IP address is used.

### GET /filters

Return all filter IDs as a JSON array.
The IDs are ordered the same as the filters are ordered.

### PUT /filters/ID

* Content-Type: application/json
* Body: see below.

`ID` must match this regexp: `^[a-zA-Z0-9_-]+$`

Create a new filter with `ID`, or edit the existing filter matching `ID`.
The JSON object must be an object with these fields:

| Name | Required | Type | Description |
| ---- | -------- | ---- | ----------- |
| `type` | Yes | string | Filter type such as `discard`, `group`, `route`. |
| `disabled` | No | bool | If `true`, the filter will not be used. |
| `all` | No | bool | If `true`, the filter works for all alerts (not one-by-one). |
| `if` | No | string/array of strings | Filter condition. |
| `expire` | No | string | RFC3339 date string. |

Other fields may be used depending on the filter type.

The default values of `disabled` and `all` are `false`.

`if` may be either a string of JavaScript boolean expression to
test an alert or an array of alerts should be filtered, or an array
of strings to invoke an external command.

For JavaScript expressions, when `all` is `false`, the filter will
assign each alert as `alert` variable and evaluate the JavaScript
expression.  When `all` is `true`, the filter will assign an array of
all alerts as `alerts`.

For external commands, the filter executes the command by passing
the array of strings to `os/exec.Command`.  If the command exits
successfully, the filter work for the alerts.  When `all` is `false`,
the filter feeds a JSON object representing an alert via stdin.

Not all filters can be configured by `all`.  For example, `group` filter
always works as if `all` is `true`.

If `expire` is given, the filter will automatically be removed
at the given date.

### GET /filters/ID

Return a JSON representation of the filter specified by `ID`.

In addition to fields for PUT method, GET may return these read-only fields:

| Name       | Type   | Description                                     |
| ---------- | ------ | ----------------------------------------------- |
| `inactive` | string | RFC3339 date string. See /filters/ID/inactivate |


### DELETE /filters/ID

Delete a filter specified by `ID`.

### PUT /filters/ID/enable

Enable the filter specified by `ID`.
The body should be empty (Content-Length is 0).

If the filter is currently inactive, it is activated.

### PUT /filters/ID/disable

Disable the filter specified by `ID`.
The body should be empty (Content-Length is 0).

### PUT /filters/ID/inactivate

Inactivate the specified filter until the given time.
The body must be a JSON object with the following single field:

| Name    | Required | Type   | Description          |
| ------- | -------- | ------ | -------------------- |
| `until` | Yes      | string | RFC3339 date string. |

### GET /routes

Return all route IDs as a JSON array.

### PUT /routes/ID

* Content-Type: application/json
* Body: JSON array of objects.

`ID` must match this regexp: `^[a-zA-Z0-9_-]+$`

Create a new route or replace existing one.

The body must be a JSON array of objects.

Each object must have these fields and may have optional fields:

| Name | Required | Type | Description |
| ---- | -------- | ---- | ----------- |
| `type` | Yes | string | Transport type such as "email" or "slack". |

An example JSON may look like:

```javascript
[
    {
        "type": "email",
        "from": "kkok@example.org",
        "to": ["ymmt2005@example.org"]
    },
    {
        "type": "slack",
        "url": "https://hooks.slack.com/xxxxxxx"
    }
]
```

### GET /routes/ID

Return a JSON representation of the route specified by `ID`.

[JSON]: http://json.org/
[RFC6750]: https://tools.ietf.org/html/rfc6750
