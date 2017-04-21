/*
Package slack provides a transport to send alerts to Slack.

This transport uses Slack's incoming webhooks.  For details about
incoming webhooks, read https://api.slack.com/incoming-webhooks .

Alerts are sent as attachments.  For details, see:
https://api.slack.com/docs/message-attachments

The plugin takes these construction parameters:

    Name       Type        Default     Description
    label      string      ""          Arbitrary string label.
    url        string                  Incoming webhook URL.  Required.
    max_retry  int         3           Max retry count when server returns 500.
    name       string      ""          Customize the user name.
    icon       string      ""          Customize the user icon.  Emoji only.
    channel    string      ""          Override the default channel.
    color      string      ""          JavaScript expression to choose color.
    template   string      ""          Filesystem path of the template file.

An incoming webhook has the default user name, icon, and a channel
to post messages.  "name", "icon", and "channel" construction parameters
can override these defaults.

"color" is a JavaScript expression that should evaluates to a string.
The string should be either an RGB hex code such as "#D0B011" or one of
"good", "warning", "danger".  If it returns an empty string, the color
will not be specified.

To customize the message body, set "template" to a template file.
The template must be written for text/template package.  To escape
special characters, the template provides a non-standard function "slack"
to escape strings for Slack.

Example snippet for TOML configuration:

    [[route.notify]]
    type        = "slack"
    url         = "https://hooks.slack.com/services/xxxx/yyyy/zzzz"
    color       = "alert.Info.severity"

This example assumes that an alert may have "good", "warning", or
"danger" in its Info["severity"] field.
*/
package slack
