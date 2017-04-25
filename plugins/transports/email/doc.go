/*
Package email provides a transport to send alerts as emails.

This transport sends alerts one by one.  Each alert is rendered as
plain text using text/template package.  The default template is
built-in as DefaultTemplate.

"Subject" header value will be the alert's Title field value.
"Date" header value will be the alert's Date field value.

The plugin takes these construction parameters:

    Name      Type        Default     Description
    label     string      ""          Arbitrary string label.
    from      string                  Address for From header.  Required.
    host      string      localhost   SMTP server name.
    port      int         25          SMTP server port.
    user      string      ""          Username for SMTP auth.
    password  string      ""          Password for SMTP auth.
    to        []string    nil         Addresses for To header.
    cc        []string    nil         Addresses for Cc header.
    bcc       []string    nil         Addresses of non-disclosed recipients.
    to_file   string      ""          Filename.  See below.
    cc_file   string      ""          Filename.  See below.
    bcc_file  string      ""          Filename.  See below.
    template  string      ""          Filesystem path of the template file.

Recipient addresses are statically provided by "to", "cc", or "bcc".
If "to_file", "cc_file", or "bcc_file" is given, the file contents
will be (re-)read each time when alerts are sent.  The contents shall
list mail addresses line by line.  If a file specified in "to_file",
"cc_file", or "bcc_file" does not exist, the plugin just ignores it.

To customize the email body, set "template" to a template file.
The template must be written for text/template package.

Example snippet for TOML configuration:

    [[route.notify]]
    type        = "email"
    label       = "send alerts to alert@example.com"
    from        = "kkok@example.com"
    to          = ["alert@example.com"]
*/
package email
