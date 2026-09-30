package trackers

import "regexp"

// knownDomains are hosts that only serve open-tracking pixels. A host matches when it or any
// of its parent domains is listed. Hosts that also serve content images (e.g. sendgrid.net,
// hubspotemail.net, customeriomail.com) are listed by their tracking subdomain only.
var knownDomains = map[string]struct{}{
	// Mailchimp / Mandrill
	"list-manage.com": {},
	"mandrillapp.com": {},
	// SendGrid (content images are served from other sendgrid.net subdomains)
	"ct.sendgrid.net": {},
	// HubSpot
	"t.hubspotemail.net": {},
	"t.hubspotfree.net":  {},
	"track.hubspot.com":  {},
	// Mailgun
	"mailgun.net": {},
	"mailgun.org": {},
	// Klaviyo
	"trk.klaviyo.com":     {},
	"ctrk.klaviyo.com":    {},
	"trk.klaviyomail.com": {},
	// Braze (content images are served from braze-images.com / appboy-images.com)
	"braze.com": {},
	"braze.eu":  {},
	// Customer.io
	"track.customer.io": {},
	// Mixpanel
	"api.mixpanel.com": {},
	// Google Analytics
	"google-analytics.com": {},
	// Salesforce Marketing Cloud (ExactTarget)
	"exct.net": {},
	// Intercom
	"via.intercom.io": {},
	// Superhuman
	"r.superhuman.com": {},
	// Mailtrack / Mailsuite
	"mailtrack.io":  {},
	"mailsuite.com": {},
	// Yesware
	"yesware.com": {},
	// Postmark
	"pstmrk.it": {},
	// Amazon SES
	"awstrack.me": {},
	// Constant Contact
	"rs6.net": {},
	// Litmus
	"emltrk.com": {},
	// Streak
	"mailfoogae.appspot.com": {},
	// Bananatag
	"bl-1.com": {},
}

// urlPatterns match the lowercased URL path (without query) of common open-tracking endpoints,
// including those served from a sender's own branded tracking domain.
var urlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`/track/open`),
	regexp.MustCompile(`/wf/open`),
	regexp.MustCompile(`/open(\.(php|aspx|jsp|gif|png))?$`),
	regexp.MustCompile(`/o\.gif$`),
	regexp.MustCompile(`/pixel(\.(gif|png|php))?$`),
}
