Architecture
============

How it works
------------

kkok sends alerts with these steps:

1. Generate alerts from sources.
2. Collect and pool alerts for some duration.
3. Edit/route collected alerts by filters.
4. Send alerts along with the given routes.

The rest of this document describes terminologies and concepts
appearing in kkok.

Alert
-----

An alert is represented as an object with these fields:

| Name | Optional | Type | Description |
| ---- | -------- | ---- | ----------- |
| `From` | No | string | Who sent this alert. |
| `Date` | No | string | RFC3339 format date string. |
| `Host` | No | string | Where this alert was generated. |
| `Title` | No | string | One-line description of the alert. |
| `Message` | Yes | string | Optional multi-line description of the alert. |
| `Routes` | No | array of strings | List of routes for alert receivers. |
| `Info` | No | object | Additional fields. |
| `Sub` | Yes | array of objects | A list of sub-alerts for grouped alert. |

Additionally, an alert has `Stats` that is a `map[string]float64`
to bring dynamically calculated values between filters.  `Stats`
is not exported to nor imported from JSON.

An alert can be represented as a simple JSON object like this:

```javascript
{
    "From": "NTP monitor",
    "Date": "2016-11-02T11:23:45.789Z",
    "Host": "host1",
    "Title": "[ntp] lost sync over 300 seconds",
    "Message": "output of ntpq -p:\n...",
    "Routes": ["emergency", "mailbox"],
    "Info": {"Domain": "www.cybozu.com", "User": "ymmt2005"}
}
```

Generator
---------

Generators generates alert objects.  Generated alerts are pooled
by kkok for some duration.  The most basic generator is REST API
to post an alert directly to kkok through HTTP.

kkok configures generators statically at process start.

`Routes` of new alerts are empty as routing should be done by filters.

Routes
------

A route is a set of transportation means to send alerts to receivers.

For example, an "emergency" route can consist of SMS notifications
and emails destined for on-call SREs while a "moderate" route posts
only to a Slack channel.

This example can be expressed in [TOML][] as:

```
[[route.emergency]]
type = "email"
from = "kkok@example.com"
to = ["ymmt2005@example.com"]
tofile = "/etc/kkok/mailto"

[[route.emergency]]
type = "twilio"
sid = "*******************"
token = "*******************"
from = "999888777"
to = ["0123456789", "111222333"]

[[route.moderate]]
type = "slack"
url = "https://hooks.slack.com/services/**********"
```

Filter
------

Filters edit alerts and/or do whatever for alert handling.

Filters can be defined statically in configuration files or dynamically
via REST API.  All filters have the following properties:

| Name       | Type | Description |
| ---------- | ---- | ----------- |
| `type`     | string | Filter type such as `discard`, `group`, `route`. |
| `id`       | string | The unique ID of the filter. |
| `label`    | string | Arbitrary string label. |
| `disabled` | bool | Disabled filters will not be used. |
| `all`      | bool | If `true`, the filter works for all alerts (not one-by-one). |
| `if`       | string/array of strings | Filter condition. See below. |
| `scripts`  | []string | JavaScript files to be loaded. |

Filters with `if` will only work for alerts matching the given condition.

`if` may be either:
1) a string of JavaScript boolean expression to test an alert, or
2) an array of alerts should be filtered, or
3) an array of strings to invoke an external command.

For JavaScript expressions, when `all` is `false`, the filter will
assign each alert as `alert` variable and evaluate the JavaScript
expression.  When `all` is `true`, the filter will assign an array of
all alerts as `alerts`.

In addition, `scripts` may be specified to load JavaScript files
before evaluating `if` expressions.  Scripts are reloaded from
files every time a filter starts evaluation.

For external commands, the filter executes the command by passing
the array of strings to `os/exec.Command`.  If the command exits
successfully, the filter work for the alerts.  When `all` is `false`,
the filter feeds a JSON object representing an alert via stdin.

For example, the following filter groups all alerts if the number of
pooled alerts is larger than 10.

```
[[filters]]
type = "group"
id = "toomanyalerts"
label = "Too many alerts"
all = true
if = "alerts.length > 10"
```

### Filter ordering

Filters defined first will be applied first.

This means dynamic filters will always be applied after static filters.

### Disabling filters

Filters can be disabled completely or temporarily until a given time.
This helps operators suppress alerts.

Disabling filters temporarily is implemented as _temporary inactivation_ of
filters.

### Temporary filters

Dynamic filters have an extra property called `expire` whose value
is an RFC3339 format date string.  Filters will be removed automatically
when they *expire* as of the `expire` property values.

Dynamic filters with expiration dates can be used to, for example,
change routes temporarily.

Temporary filters can be considered _temporarily enabled filters_.
For this reason, kkok does not provide temporary activation of filters.

Mute
----

In order to avoid sending emergency alerts too frequently, _route filter_
has an option to inactivate itself for a given period.

Instead of route filter, we could add such an option to _route_.
However, doing so would also suppress unrelated emergency alerts hence
avoided.

[TOML]: https://github.com/toml-lang/toml
