/*
Package exec provides a filter to edit alerts by an external command.

External commands should read a JSON from stdin, and write JSON to stdout.

If "all" is specified at the filter construction, the input and output
shall be a list of kkok.Alert objects exported to JSON.  Otherwise,
the input and output shall be a kkok.Alert object exported to JSON.

In addition to the standard filter construction parameters, this
plugin takes these parameters:

    Name      Type       Default   Description
    command   []string             Command and arguments.  Required.
    timeout   int        5         Seconds before killing the command.
                                   If max_time is 0, the command will not be killed.

Example snippet for TOML configuration:

    [[filter]]
    type        = "exec"
    id          = "headonly"
    label       = "pick the first alert"
    all         = true
    command     = ["jq", "[.[0]]"]

Another example to process alerts one by one:

    [[filter]]
    type        = "exec"
    id          = "edit_subject"
    label       = "add prefix to subjects for emergency alerts"
    command     = ["jq", """
if .Routes|contains(["emergency"]) then
    . + {"Title": ("[warn] " + .Title)}
else
    .
end"""]
*/
package exec
