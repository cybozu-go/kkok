/*
Package discard provides a filter to eliminate alerts matching
given conditions.

There are no additional construction parameters for this filter.

Example snippet for TOML configuration:

    [[filter]]
    type        = "discard"
    id          = "ignore_host1"
    if          = "alert.Host == 'host1'"
*/
package discard
