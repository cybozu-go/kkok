/*
Package route provides a filter to add or replace routes.

A route consists of a set of kkok.Transport objects and is
identified by a unique ID.  This filter can add or replace
such route IDs to alerts in order to control who receive them.

In addition to the standard filter construction parameters, this
plugin takes these parameters:

    Name            Type           Default       Description
    routes          []string       nil           Array of route IDs.
    replace         bool           false         If true, replace routes.
    auto_mute       bool           false         If true, mute itself for a while.
    mute_seconds    int            60            Seconds for auto mute.
    mute_routes     []string       nil           Alternative route while muting.

With default parameters, route filter adds route IDs given by "routes" to
the current set of route IDs.  If "replace" is true, route filter replaces
the current set with "routes".

If "auto_mute" is true, route filter becomes inactive for "mute_seconds".
This can be used to suppress too frequent emergency alerts.  If
"mute_routes" is not empty, the filter will use it instead of "routes"
while muting.

Example snippet for TOML configuration:

    [[filter]]
    type        = "route"
    id          = "emergency_filter"
    if          = "alert.From == 'alive monitor'"
    routes      = ["slack", "mail"]
    auto_mute   = true
    mute_routes = ["mail"]

This filter adds "slack" and "mail" routes to alerts whose From
field is "alive monitor".  It will be automatically muted for 60
seconds after processing one such alert.  While muting, it adds
only "mail" route to such alerts.
*/
package route
