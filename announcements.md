# Announcements

## Catch-all aliases, custom alias names, and a forwarding privacy fix - 2026-09-01

### Create aliases from catch-all mail (custom domains)
Turn on "Create alias" for a verified domain and Mailx creates a real alias the first time an address receives mail. Those addresses now appear in your Aliases list, so you can disable or delete one without turning off catch-all. Capped at 10 new aliases per hour.

### Wider character set for custom aliases
Custom alias names used with custom domains are no longer limited to letters and numbers. They can now include dots, hyphens, plus signs and the other characters an email address allows.

### Forwarded mail no longer carries your real address
The To header of a forwarded message now shows the alias the sender wrote to, not your real address. Previously it carried your real address, which some mail clients exposed in the quoted text when you replied.

## Web app outage on 23 August - 2026-08-24

On Sunday 23 August the Mailx web app was unavailable between 08:00 and 21:30 UTC. Mailx has been running normally since, and no action is required from you.

### What was affected
- Web interface, including sign in, account creation, new alias creation from the web app, adding a custom domain.

### What was unaffected
- Email forwarding was not affected and no mail was lost - all messages sent to Mailx aliases during that window were forwarded and delivered as normal.
- Browser extension, creating aliases from the extension was available.

### What caused it
The cause was one of our servers failing. Failover should have moved traffic over automatically, but a configuration error meant it did not, and the web app stayed down until a manual fix. Delays in detection were caused by our monitoring setup, which checks servers only and not the availability of the web app.

### What we are changing
We apologise for the downtime, it's not acceptable for a service like Mailx. We are changing two things to reduce the chance of it happening again:
1. Fixing the failover configuration that caused this and reviewing the rest of our failover setup for the same class of error.
2. Adding independent availability checks against both the servers and the web app, so we are alerted quickly rather than hours in.

## Catch-all domains and alias restore added - 2026-07-29

### Catch-all for custom domains
Turn on catch-all for a verified domain to forward mail sent to any address on it, not just addresses you set up as aliases. Set the catch-all recipient and from name per domain.

### Restore deleted aliases
A deleted alias can be restored 90 days after deletion. It stays inactive and does not forward while deleted. Filter the Aliases list with "Show: Active / Deleted / All" to see deleted aliases and restore them.

## Recent Mailx updates - 2026-06-14

### Custom domains and custom aliases
Add and verify your own domain, then use it for new aliases instead of the built-in options. Verification is DNS-based and starts on the "Domains" page. Aliases now have a "Custom" option: set the alias text yourself instead of a randomly generated name/ID when using your own domain. Both work in the web app and the browser extension.

### Browser extension
The Mailx extension adds a small button next to email fields on any site. Click it to generate a new alias and autofill it into the form, no copy/pasting. The extension popup also lets you create and check aliases directly. Available for [Chrome](https://chromewebstore.google.com/detail/mailx/ogdmgaidomgkpefdgciobijgbobjmkpk), [Firefox](https://addons.mozilla.org/en-US/firefox/addon/mailx/) and [Edge](https://microsoftedge.microsoft.com/addons/detail/mailx/bpddflmpdckdchojiejdngkgphjgpjai). 

### Diagnostics updates
Diagnostic logs now catch more delivery issues, including failed and deferred forwarding, not just discarded messages. Turn on logging in settings if you want to see what's happening and debug it yourself. Optional.
