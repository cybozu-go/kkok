/*
Package exec provides a transport to send alerts via external commands.

Alerts are serialized as JSON and can be read from stdin.

If "all" construction parameter is true, the input will be a list of
kkok.Alert objects exported to JSON.  Otherwise, the input will be
a kkok.Alert object exported to JSON.

The plugin takes these construction parameters:

    Name      Type       Default   Description
    label     string     ""        Arbitrary string label.
    command   []string             Command and arguments.  Required.
    all       bool       false     See above description.
    timeout   int        5         Seconds before killing the command.
                                   If 0, the command will not be killed.

Example snippet for TOML configuration:

    [[route.notify]]
    type        = "exec"
    label       = "send alerts via curl"
    command     = ["curl", "--data-binary", "@-", "-f", "-s",
                   "-H", "Content-Type: application/json",
                   "http://some.service.com/"]
    all         = true

This example sends all alerts to an external service in one request.
If curl cannot connect to the external service or if the service
returns status code other than 200 (by -f option), an error will
be logged.
*/
package exec
