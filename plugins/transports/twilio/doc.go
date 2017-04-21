/*
Package twilio provides a transport to send alerts via SMS using Twilio.

This transport uses Twilio SMS messaging API.  For details about
SMS API, read https://www.twilio.com/docs/api/rest/sending-messages

Message segment per second (MPS) is limited to 1 for US/Canada or 10
for other countries.  The plugin automatically adjust sending rates to
comply these rate limits.

The plugin takes these construction parameters:

    Name        Type        Default     Description
    label       string      ""          Arbitrary string label.
    account     string                  AccountSID.  Required.
    key_sid     string      ""          API key SID.  Optional.  See below.
    token       string                  AuthToken or API key Secret.  Required.
    from        string                  Caller phone number.  Required.
    to          []string    nil         Destination phone numbers.
    to_file     string      ""          Filename.  See below.
    max_length  int         160         The maximum body length in characters.
    max_retry   int         3           Max retry count when server returns 500.
    template    string      ""          Filesystem path of the template file.
    count_only  bool        false       See below.

If "key_sid" is empty, "token" must be AuthToken for the account.
To use REST API key, "key_sid" and "token" need to be the key SID
and secret, respectively.
For details, read https://www.twilio.com/docs/api/rest/keys

Recipients are statically provided by "to".  If "to_file" is given,
the file contents will be (re-)read each time when alerts are sent.
The contents shall list recipient phone numbers line by line.
If a file specified by "to_file" does not exist, the plugin just
ignores it.

SMS message body is rendered by a text/template template.
The default template is built-in as DefaultTemplate.
To customize the message body, set "template" to a template file.

Messages body that exceeds "max_length" characters will be truncated
to "max_length".

Alternatively, if count_only is true, SMS body contains only the
number of alerts instead of renderng each alert.

Example snippet for TOML configuration:

    [[route.notify]]
    type        = "twilio"
    account     = "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    token       = "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
    from        = "+1484xxxxxxx"
    to          = ["+8180zzzzzzzz", "+8170yyyyyyyy"]
    count_only  = true

This example sends an SMS notifying just the number of alerts received
instead of each alert details.
*/
package twilio
