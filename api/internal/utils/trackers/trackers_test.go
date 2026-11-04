package trackers

import (
	"strings"
	"testing"
)

func TestRemove_RemovesTrackers(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		gone    []string
		want    []Tracker
		keep    []string
		wantOut string
	}{
		{
			name:  "known domain",
			input: `<p>Hi</p><img src="https://mc.us1.list-manage.com/track/open.php?u=1&amp;id=2" alt="">`,
			gone:  []string{"list-manage.com"},
			want:  []Tracker{{"mc.us1.list-manage.com", ElementImg, ReasonKnownDomain}},
			keep:  []string{"<p>Hi</p>"},
		},
		{
			name:  "known domain protocol-relative and uppercase",
			input: `<img src="//U123.CT.SendGrid.NET/wf/open?upn=abc">`,
			gone:  []string{"sendgrid"},
			want:  []Tracker{{"u123.ct.sendgrid.net", ElementImg, ReasonKnownDomain}},
		},
		{
			name:  "url pattern on a branded domain",
			input: `<img src="https://url1234.shop.example/wf/open?upn=abc" width="600">`,
			gone:  []string{"url1234.shop.example"},
			want:  []Tracker{{"url1234.shop.example", ElementImg, ReasonURLPattern}},
		},
		{
			name:  "1x1 attributes",
			input: `<img src="https://news.example.com/i/abc.gif" width="1" height="1" border="0">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonTinyImage}},
		},
		{
			name:  "1px inline css",
			input: `<img src="https://news.example.com/i/abc.gif" style="width:1px;height:1px">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonTinyImage}},
		},
		{
			name:  "display none",
			input: `<img src="https://news.example.com/i/abc.gif" style="DISPLAY: none !important">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonHiddenImage}},
		},
		{
			name:  "visibility hidden",
			input: `<img src="https://news.example.com/i/abc.gif" style="visibility:hidden">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonHiddenImage}},
		},
		{
			name:  "opacity zero",
			input: `<img src="https://news.example.com/i/abc.gif" style="opacity:0.0">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonHiddenImage}},
		},
		{
			name:  "css width zero",
			input: `<img src="https://news.example.com/i/abc.gif" style="width:0;max-height:0">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonHiddenImage}},
		},
		{
			name:  "hidden attribute",
			input: `<img hidden src="https://news.example.com/i/abc.gif">`,
			gone:  []string{"abc.gif"},
			want:  []Tracker{{"news.example.com", ElementImg, ReasonHiddenImage}},
		},
		{
			name:    "srcset entry",
			input:   `<img src="https://cdn.example.com/a.png" srcset="https://cdn.example.com/a.png 1x, https://trk.klaviyomail.com/open/x.png 2x">`,
			want:    []Tracker{{"trk.klaviyomail.com", ElementSrcset, ReasonKnownDomain}},
			wantOut: `<img src="https://cdn.example.com/a.png" srcset="https://cdn.example.com/a.png 1x"/>`,
		},
		{
			name:    "srcset only entry drops the attribute",
			input:   `<picture><source srcset="https://api.mixpanel.com/track?data=x"><img src="https://cdn.example.com/a.png"></picture>`,
			want:    []Tracker{{"api.mixpanel.com", ElementSrcset, ReasonKnownDomain}},
			wantOut: `<picture><source/><img src="https://cdn.example.com/a.png"/></picture>`,
		},
		{
			name:    "background attribute",
			input:   `<table><tbody><tr><td background="https://cl.s7.exct.net/open.aspx?ffcb10">x</td></tr></tbody></table>`,
			want:    []Tracker{{"cl.s7.exct.net", ElementBackground, ReasonKnownDomain}},
			wantOut: `<table><tbody><tr><td>x</td></tr></tbody></table>`,
		},
		{
			name:    "background-image url in inline style",
			input:   `<div style="color:red; background-image: url('https://www.google-analytics.com/collect?v=1')">x</div>`,
			want:    []Tracker{{"www.google-analytics.com", ElementBackground, ReasonKnownDomain}},
			wantOut: `<div style="color:red; background-image: none">x</div>`,
		},
		{
			name:  "background url in style block",
			input: `<html><head><style>.a{background:#fff url("https://r.superhuman.com/abc.png") no-repeat}</style></head><body><div class="a">x</div></body></html>`,
			want:  []Tracker{{"r.superhuman.com", ElementBackground, ReasonKnownDomain}},
			keep:  []string{`.a{background:#fff none no-repeat}`},
			gone:  []string{"superhuman"},
		},
		{
			name:    "prefetch link",
			input:   `<link rel="dns-prefetch" href="//t.yesware.com"><p>x</p>`,
			want:    []Tracker{{"t.yesware.com", ElementLink, ReasonKnownDomain}},
			wantOut: `<p>x</p>`,
		},
		{
			name:  "img in noscript",
			input: `<noscript><img src="https://mailtrack.io/trace/mail/abc.png"></noscript><p>x</p>`,
			want:  []Tracker{{"mailtrack.io", ElementImg, ReasonKnownDomain}},
			gone:  []string{"mailtrack"},
		},
		{
			name:  "backslashes and userinfo read as the browser does",
			input: `<img src="https:\\user@T.Yesware.com:443\t\abc.gif">`,
			want:  []Tracker{{"t.yesware.com", ElementImg, ReasonKnownDomain}},
			gone:  []string{"yesware"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Remove(tt.input)
			if err != nil {
				t.Fatalf("Remove() error = %v", err)
			}
			if len(res.Removed) != len(tt.want) {
				t.Fatalf("Removed = %+v, want %+v", res.Removed, tt.want)
			}
			for i := range tt.want {
				if res.Removed[i] != tt.want[i] {
					t.Errorf("Removed[%d] = %+v, want %+v", i, res.Removed[i], tt.want[i])
				}
			}
			for _, s := range tt.gone {
				if strings.Contains(strings.ToLower(res.HTML), strings.ToLower(s)) {
					t.Errorf("HTML still contains %q: %s", s, res.HTML)
				}
			}
			for _, s := range tt.keep {
				if !strings.Contains(res.HTML, s) {
					t.Errorf("HTML lost %q: %s", s, res.HTML)
				}
			}
			if tt.wantOut != "" && res.HTML != tt.wantOut {
				t.Errorf("HTML = %s, want %s", res.HTML, tt.wantOut)
			}
		})
	}
}

func TestRemove_KeepsLegitimateImages(t *testing.T) {
	inputs := []string{
		`<img src="https://cdn.example.com/hero.jpg" width="600" height="300" alt="Hero">`,
		`<img src="https://gallery.mailchimp.com/logo.png">`,
		`<img src="https://mcusercontent.com/logo.png" style="display:block;width:100%">`,
		`<img src="cid:logo@example.com" width="1" height="1">`,
		`<img src="data:image/gif;base64,R0lGODlhAQABAAAAACw=" width="1" height="1" style="display:none">`,
		`<img src="/relative/pixel.gif" width="1" height="1">`,
		`<img src="https://cdn.example.com/a.png" srcset="data:image/png;base64,iVBORw0KGgo=, https://cdn.example.com/b.png 2x">`,
		`<img src="https://cdn.example.com/a.png" width="1%">`,
		`<img src="https://cdn.example.com/open-sign.png">`,
		`<a href="https://mc.us1.list-manage.com/track/click?u=1">Click</a>`,
		`<div style="background-image:url(https://cdn.example.com/bg.png)">x</div>`,
		`<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto">`,
	}

	for _, input := range inputs {
		res, err := Remove(input)
		if err != nil {
			t.Fatalf("Remove(%s) error = %v", input, err)
		}
		if len(res.Removed) != 0 {
			t.Errorf("Remove(%s) removed %+v, want nothing", input, res.Removed)
		}
		if res.HTML != input {
			t.Errorf("Remove(%s) changed the HTML to %s", input, res.HTML)
		}
	}
}

func TestRemove_Count(t *testing.T) {
	input := `<html><body>
<img src="https://cdn.example.com/hero.jpg" width="600">
<img src="https://mc.us1.list-manage.com/track/open.php?u=1">
<img src="https://news.example.com/p.gif" width="1" height="1">
<img src="https://cdn.example.com/a.png" srcset="https://t.hubspotemail.net/e2t/to/a 1x, https://track.hubspot.com/__ptq.gif 2x">
<table><tr><td background="https://email.mailgun.net/o/abc">x</td></tr></table>
<link rel="preload" href="https://via.intercom.io/o?abc">
<img src="cid:logo">
</body></html>`

	res, err := Remove(input)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if got := len(res.Removed); got != 6 {
		t.Errorf("count = %d, want 6: %+v", got, res.Removed)
	}
	for _, s := range []string{"hero.jpg", "cid:logo", "cdn.example.com/a.png"} {
		if !strings.Contains(res.HTML, s) {
			t.Errorf("HTML lost %q", s)
		}
	}
}

func TestRemove_DocumentKeepsStructure(t *testing.T) {
	input := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Hi</title></head>` +
		`<body style="margin:0" bgcolor="#eeeeee"><!--[if mso]><table><tr><td><![endif]-->` +
		`<p>Grüße &amp; ✓</p><img src="https://ea.pstmrk.it/open?m=abc" width="1" height="1"></body></html>`

	res, err := Remove(input)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if len(res.Removed) != 1 {
		t.Fatalf("Removed = %+v, want 1", res.Removed)
	}
	for _, s := range []string{
		"<!DOCTYPE html>", `<html lang="en">`, `<meta charset="utf-8"/>`, `<body style="margin:0" bgcolor="#eeeeee">`,
		"<!--[if mso]><table><tr><td><![endif]-->", "<p>Grüße &amp; ✓</p>",
	} {
		if !strings.Contains(res.HTML, s) {
			t.Errorf("HTML lost %q: %s", s, res.HTML)
		}
	}
}

func TestRemove_FragmentGetsNoWrappers(t *testing.T) {
	res, err := Remove(`<div>Hello</div><img src="https://mailtrack.io/trace/x.png">`)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if res.HTML != "<div>Hello</div>" {
		t.Errorf("HTML = %q, want %q", res.HTML, "<div>Hello</div>")
	}
}

func TestRemove_SurvivesMalformedHTML(t *testing.T) {
	inputs := []string{
		`<div><p>unclosed <b>bold <img src="https://mailtrack.io/x.png" <table><td>cell`,
		`<img src="https://mailtrack.io/x.png"`,
		`<<<>>><img/src="https://mailtrack.io/x.png"/></div></div></span>`,
		`<table><tr><td><img src="https://mailtrack.io/x.png"></td></table></tr>`,
		"\x00\xff\xfe<img src=\"https://mailtrack.io/x.png\">",
		`<style>.a{background:url(https://mailtrack.io/x.png</style>`,
		`<img srcset=",,, ,">`,
		`<img src="https:" width="1">`,
		`<img src="https://[::1" width="1">`,
	}

	for _, input := range inputs {
		input = "<p>keep</p>" + input
		res, err := Remove(input)
		if err != nil {
			t.Errorf("Remove(%q) error = %v", input, err)
			continue
		}
		if !strings.Contains(res.HTML, "<p>keep</p>") {
			t.Errorf("Remove(%q) = %q, lost the content", input, res.HTML)
		}
	}
}

func TestRemove_NoTrackersReturnsInputUnchanged(t *testing.T) {
	// Would be normalized by a parse/render round trip
	input := `<TABLE><TR><TD><P>unclosed<BR>&nbsp;<img src=https://cdn.example.com/a.png></TD></TR></TABLE>`

	res, err := Remove(input)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if res.HTML != input {
		t.Errorf("HTML = %q, want the input unchanged", res.HTML)
	}
}

func TestRemove_TooLarge(t *testing.T) {
	input := `<img src="https://mailtrack.io/x.png">` + strings.Repeat("a", MaxHTMLSize)

	res, err := Remove(input)
	if err != ErrTooLarge {
		t.Fatalf("Remove() error = %v, want ErrTooLarge", err)
	}
	if res.HTML != input || len(res.Removed) != 0 {
		t.Error("Remove() changed a body over the size limit")
	}
}

func TestReport_RoundTrip(t *testing.T) {
	removed := []Tracker{
		{"mc.us1.list-manage.com", ElementImg, ReasonKnownDomain},
		{"news.example.com", ElementImg, ReasonTinyImage},
	}

	encoded := EncodeReport(removed)
	if strings.ContainsAny(encoded, "+/=") {
		t.Errorf("EncodeReport() = %q, want base64url without padding", encoded)
	}

	total, got, err := DecodeReport(encoded)
	if err != nil {
		t.Fatalf("DecodeReport() error = %v", err)
	}
	if total != 2 || len(got) != 2 || got[0] != removed[0] || got[1] != removed[1] {
		t.Errorf("DecodeReport() = %d %+v, want 2 %+v", total, got, removed)
	}
}

func TestReport_CapsEntries(t *testing.T) {
	removed := make([]Tracker, MaxReportEntries+25)
	for i := range removed {
		removed[i] = Tracker{"mailtrack.io", ElementImg, ReasonKnownDomain}
	}

	total, got, err := DecodeReport(EncodeReport(removed))
	if err != nil {
		t.Fatalf("DecodeReport() error = %v", err)
	}
	if total != len(removed) {
		t.Errorf("total = %d, want %d", total, len(removed))
	}
	if len(got) != MaxReportEntries {
		t.Errorf("entries = %d, want %d", len(got), MaxReportEntries)
	}
}

func FuzzRemove(f *testing.F) {
	f.Add(`<img src="https://mailtrack.io/x.png" width="1">`)
	f.Add(`<html><body><td background="//exct.net/a" style="background:url(https://exct.net/b)"></body></html>`)
	f.Add(`<noscript><img src="https://mailtrack.io/x.png"></noscript>`)

	f.Fuzz(func(t *testing.T, input string) {
		res, err := Remove(input)
		if err != nil && res.HTML != input {
			t.Errorf("Remove() errored but changed the HTML")
		}
		if len(res.Removed) == 0 && res.HTML != input {
			t.Errorf("Remove() removed nothing but changed the HTML")
		}
	})
}
