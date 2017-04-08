/*
Package group provides a filter to merge alerts into groups.

The filter can evaluate a JavaScript expression for each alert to
find a grouping key.  For example, if the filter is configured with
"alert.Host" JavaScript expression, the filter will merge alerts
having the same Host field value into one.

The field values of a merged alert are determined as follows:

    Sub             Original alerts.
    From            If alerts in Sub share the same From, then that.
                    Otherwise, "from" construction parameter value.
    Date            The current date.
    Host            If alerts in Sub share the same Host, then that.
                    Otherwise, "localhost".
    Title           If alerts in Sub share the same Title, then that.
                    Otherwise, "title" construction parameter value.
    Message         If alerts in Sub share the same Message, then that.
                    Otherwise, "message" construction parameter value.
    Routes          "routes" construction parameter value.
    Info            nil.

In addition to the standard filter construction parameters, this
plugin takes these parameters:

    Name            Type        Default         Description
    by              string      ""              JavaScript expression.
    from            string      "filter:<ID>"   <ID> is the filter ID string.
    title           string      "merged alert"  see above.
    message         string      ""              see above.
    routes          []string    nil             see above.

Example snippet for TOML configuration:

    [[filter]]
    type        = "group"
    id          = "group_by_host"
    if          = "alert.From == 'process monitor'"
    by          = "alert.Host"
    title       = "some processes died"

This filter merges alerts from "process monitor" grouped by Host field value.
*/
package group
