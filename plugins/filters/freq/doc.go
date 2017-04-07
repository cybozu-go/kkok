/*
Package freq provides a filter to calculate frequency of the given alerts.

Frequency is calculated as the number of alerts received for a given
period.  For instance, if the period is set to 60 seconds and there were
two of such alerts in the last 60 second, then the next alert's frequency
value will become 3.

For convenience, the value can be divided by a constant in order to
normalize the value as alerts per second, minute, etc.

The calculated frequency value is stored into Stats field which is
a map[string]float64.  The map key is configurable but the default
is the filter ID.

This filter can classify alerts and calculate frequencies for
each class.  To do so, the filter has an optional construction parameter
"foreach".  The value of "foreach" is a JavaScript expression whose
evaluated value will be used to classify alerts.

In addition to the standard filter construction parameters, this
plugin takes these parameters:

    Name            Type           Default       Description
    duration        int            600           Seconds for collection.
    divisor         float64        10.0          Constant divisor.
    foreach         string         ""            JavaScript expression.
    key             string         nil           Key of Stats.

Example snippet for TOML configuration:

    [[filter]]
    type        = "freq"
    id          = "failed_process"
    if          = "alert.From == 'process monitor'"
    foreach     = "alert.Host"

This filter calculates the frequencies of alerts from "process monitor"
for each Host as the average number of alerts per minute
(600 / 10 = 60 seconds) for the latest 10 minutes (600 seconds).
The calculated value will be stored in alert's Stats["failed_process"].
*/
package freq
