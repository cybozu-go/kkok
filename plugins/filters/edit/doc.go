/*
Package edit provides a filter to edit alerts by JavaScript.

Unlike other areas of kkok where alerts are Go objects natively
exported to JavaScript, this filter converts Go alert object to
a pure JavaScript object to grant users to edit it freely.

Specifically, an alert in this filter is a JavaScript object with
these properties:

    Name     Type           Reference
    From     string
    Date     Date           https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date
    Host     string
    Title    string
    Message  string
    Routes   Array          https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array
    Info     Object         https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object
    Stats    Object         ditto
    Sub      []*kkok.Alert  (Go slice, only for reference)

This filter does not (yet) support "all" contruction parameter.
Use "exec" filter in case this filter is too limited.

In addition to the standard filter construction parameters, this
plugin takes one extra parameter:

    Name     Type       Default   Description
    code     string               JavaScript code.  Required.

Example snippet for TOML configuration:

    [[filter]]
    type        = "edit"
    id          = "addprefix"
    label       = "add prefix to alert Title"
    if          = "alert.From=='foo monitor'"
    code        = "alert.Title = '[foo] ' + alert.Title;"

Another example to remove a route:

    [[filter]]
    type        = "edit"
    id          = "removeroute"
    label       = "remove r1 route."
    code        = """
    routes = new Array();
    for( i = 0; i < alert.Routes.length; i++ ) {
        if( alert.Routes[i] != "r1" ) {
            routes.push(alert.Routes[i]);
        }
    }
    alert.Routes = routes;
    """

See test cases in convert_test.go for more examples.
*/
package edit
