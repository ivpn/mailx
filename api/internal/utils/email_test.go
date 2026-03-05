package utils

import (
	"strings"
	"testing"
)

func TestRemoveHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes header with newline",
			input:    "This email was sent to user@example.com from sender@example.com\nActual email content here",
			expected: "Actual email content here",
		},
		{
			name:     "removes header without newline",
			input:    "This email was sent to user@example.com from sender@example.com",
			expected: "",
		},
		{
			name:     "removes header at beginning of multiline text",
			input:    "This email was sent to test@domain.com from noreply@service.com\nLine 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "no header to remove",
			input:    "Regular email content without header",
			expected: "Regular email content without header",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "header in middle",
			input:    "Some content\nThis email was sent to {{.alias}} from {{.from}}\nMore content",
			expected: "Some content\nMore content",
		},
		{
			name:     "multiple headers at start",
			input:    "This email was sent to user1@example.com from sender@example.com\nThis email was sent to user2@example.com from sender@example.com\nContent",
			expected: "Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveHeader(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveHeader() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRemoveHtmlHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "removes HTML header",
			input: `<table style="width: 100%;"><tr><td><div>This email was sent to <a href="mailto:user@example.com">user@example.com</a> from <a href="mailto:sender@example.com">sender@example.com</a></div><br></td></tr></table><br>
<p>Actual email content here</p>`,
			expected: `<p>Actual email content here</p>`,
		},
		{
			name: "removes HTML header with extra whitespace",
			input: `<table class="header">  <tr>  <td>  <div style="color:gray;">  This email was sent to  <a href="mailto:test@domain.com">test@domain.com</a>  from  <a href="mailto:noreply@service.com">noreply@service.com</a>  </div>  <br>  </td>  </tr>  </table>  <br>
		<div>Content</div>`,
			expected: `<div>Content</div>`,
		},
		{
			name:     "no HTML header to remove",
			input:    `<p>Regular email content without header</p>`,
			expected: `<p>Regular email content without header</p>`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "HTML header without content after",
			input:    `<table style="width: 100%;"><tr><td><div>This email was sent to <a href="mailto:user@example.com">user@example.com</a> from <a href="mailto:sender@example.com">sender@example.com</a></div><br></td></tr></table><br>`,
			expected: ``,
		},
		{
			name: "multiple HTML headers",
			input: `<table style="width: 100%;"><tr><td><div>This email was sent to <a href="mailto:user1@example.com">user1@example.com</a> from <a href="mailto:sender@example.com">sender@example.com</a></div><br></td></tr></table><br>
<table style="width: 100%;"><tr><td><div>This email was sent to <a href="mailto:user2@example.com">user2@example.com</a> from <a href="mailto:sender@example.com">sender@example.com</a></div><br></td></tr></table><br>
<p>Content</p>`,
			expected: `<p>Content</p>`,
		},
		{
			name: "HTML header in middle of content",
			input: `<p>Some content</p>
<table style="width: 100%;">
    <tr>
        <td>
            <div style="padding: 15px;background: rgb(241 245 249);color: rgb(75 85 99);font-size: 13px;text-align: center;font-family: Arial, Helvetica, sans-serif;">
                This email was sent to
                <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:{{.alias}}">{{.alias}}</a>
                from
                <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:{{.from}}">{{.from}}</a>
            </div>
            <br>
        </td>
    </tr>
</table>`,
			expected: `<p>Some content</p>
`,
		},
		{
			name:     "HTML header in middle of content",
			input:    `<div dir="auto">Looks good.</div><div dir="auto"><br></div><div dir="auto">-- <br></div><div dir="auto"> Secured with Mail Provider: <br></div><div dir="auto"> <a href="" rel="noopener noreferrer" target="_blank"></a><br></div><div dir="auto"><br></div><div dir="auto"><br></div><div dir="auto">28 Jul 2025 at 12:49 by sender@example.com:<br></div><blockquote class="quote" style="border-left: 1px solid #93A3B8; padding-left: 10px; margin-left: 5px;"><table style="width: 100%;"><tbody><tr><td><div style="padding: 15px; background: no-repeat rgb(241, 245, 249); color: rgb(75, 85, 99); font-size: 13px; text-align: center; font-family: Arial, Helvetica, sans-serif;">This email was sent to <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:recipient@example.com" rel="noopener noreferrer" target="_blank">recipient@example.com</a> from <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:sender@example.com" rel="noopener noreferrer" target="_blank">sender@example.com</a><br></div><div dir="auto"><br></div></td></tr></tbody></table><div dir="auto"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;">Hello<br></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">您好<br></p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">这是一次测试。<br></p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'PingFang SC'">此致。<br></p></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;" class=""><div class=""><br></div><div class="">Sent with <a href="" target="_blank" rel="noopener noreferrer">Mail</a> secure email.<br></div></div></blockquote>`,
			expected: `<div dir="auto">Looks good.</div><div dir="auto">--</div><div dir="auto"> Secured with Mail Provider:</div><div dir="auto"> <a href="" rel="noopener noreferrer" target="_blank"></a></div><div dir="auto">28 Jul 2025 at 12:49 by sender@example.com:</div><blockquote class="quote" style="border-left: 1px solid #93A3B8; padding-left: 10px; margin-left: 5px;"><div style="font-family: Arial, sans-serif; font-size: 14px;">Hello</div><div style="font-family: Arial, sans-serif; font-size: 14px;"><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">您好</p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">这是一次测试。</p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'PingFang SC'">此致。</p></div><div style="font-family: Arial, sans-serif; font-size: 14px;" class=""><div class="">Sent with <a href="" target="_blank" rel="noopener noreferrer">Mail</a> secure email.</div></div></blockquote>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveHtmlHeader(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveHtmlHeader() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSafeDecodeAddressName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain name without encoding",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "UTF-8 base64 encoded word",
			input:    "=?UTF-8?B?SmFuZSBEb2U=?=",
			expected: "Jane Doe",
		},
		{
			// =?UTF-8?Q?John_D=C5=8De?= encodes "John Dōe" (ō = U+014D)
			name:     "UTF-8 QP encoded word with multibyte char",
			input:    "=?UTF-8?Q?John_D=C5=8De?=",
			expected: "John D\u014de",
		},
		{
			// When the charset is iso-8859-2, Go's mime.WordDecoder returns raw bytes
			// that are not valid UTF-8.  SafeDecodeAddressName must fall back to
			// stripping the encoded word and returning only the plain-text prefix.
			name:     "iso-8859-2 encoded word falls back to plain-text prefix",
			input:    "John =?iso-8859-2?q?D=F6e?=",
			expected: "John",
		},
		{
			name:     "no encoded words passed through unchanged",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			name:     "us-ascii QP encoded word is decoded normally",
			input:    "=?us-ascii?Q?John_Doe?=",
			expected: "John Doe",
		},
		{
			// =??q?...?= has an empty charset.  SafeDecodeAddressName decoding
			// fails (empty charset lookup), so the fallback regex strips it.
			// In practice, PreprocessEmailData fixes this before it reaches here,
			// but SafeDecodeAddressName must be safe against it regardless.
			name:     "empty charset encoded word stripped by fallback",
			input:    "Service =??q?Support_Team?=",
			expected: "Service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeDecodeAddressName(tt.input)
			if result != tt.expected {
				t.Errorf("SafeDecodeAddressName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFixEmptyCharsetEncodedWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty charset QP encoded word is fixed",
			input:    "=??q?Hello_World?=",
			expected: "=?UTF-8?q?Hello_World?=",
		},
		{
			name:     "empty charset base64 encoded word is fixed",
			input:    "=??B?SGVsbG8gV29ybGQ=?=",
			expected: "=?UTF-8?B?SGVsbG8gV29ybGQ=?=",
		},
		{
			name:     "well-formed UTF-8 encoded word is unchanged",
			input:    "=?UTF-8?q?Hello_World?=",
			expected: "=?UTF-8?q?Hello_World?=",
		},
		{
			name:     "well-formed iso-8859-1 encoded word is unchanged",
			input:    "=?iso-8859-1?q?Hello?=",
			expected: "=?iso-8859-1?q?Hello?=",
		},
		{
			name:     "plain text without encoded words is unchanged",
			input:    "John Doe <john.doe@example.com>",
			expected: "John Doe <john.doe@example.com>",
		},
		{
			// Regression: display name with empty charset followed by address with '=' in local part.
			name:     "empty charset encoded word in display name alongside address with equals",
			input:    "=??q?Service_Name?= <jane.doe+tag=example.net@example.com>",
			expected: "=?UTF-8?q?Service_Name?= <jane.doe+tag=example.net@example.com>",
		},
		{
			name:     "multiple encoded words, one with empty charset",
			input:    "=??q?First?= =?UTF-8?q?Second?=",
			expected: "=?UTF-8?q?First?= =?UTF-8?q?Second?=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixEmptyCharsetEncodedWords(tt.input)
			if result != tt.expected {
				t.Errorf("fixEmptyCharsetEncodedWords(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPreprocessEmailData_EmptyCharsetEncodedWord(t *testing.T) {
	// Regression: letters.ParseEmail fails with "cannot lookup encoding" when
	// the From or To header contains an RFC 2047 encoded word with an empty
	// charset field (=??q?…?=).  PreprocessEmailData must rewrite these to
	// =?UTF-8?q?…?= before returning.
	input := strings.Join([]string{
		"From: =??q?Service_Support?= <support@example.com>",
		"To: John Doe <john.doe+tag=example.net@example.com>",
		"Subject: Hello",
		"",
		"Body text",
	}, "\r\n")

	processed, err := PreprocessEmailData([]byte(input))
	if err != nil {
		t.Fatalf("PreprocessEmailData() error = %v", err)
	}

	processedStr := string(processed)
	if strings.Contains(processedStr, "=??q?") {
		t.Errorf("processed data still contains empty-charset encoded word: %q", processedStr)
	}
	if !strings.Contains(processedStr, "=?UTF-8?q?Service_Support?=") {
		t.Errorf("processed data missing fixed encoded word; got: %q", processedStr)
	}
	// The To address with '=' in the local part must not be corrupted.
	if !strings.Contains(processedStr, "john.doe+tag=example.net@example.com") {
		t.Errorf("processed data corrupted To address; got: %q", processedStr)
	}
}

func TestNormalizeAddressSeparators(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain address unchanged",
			in:   "alice@example.com",
			want: "alice@example.com",
		},
		{
			name: "trailing semicolon removed",
			in:   "alice@example.com;",
			want: "alice@example.com",
		},
		{
			name: "semicolon separator converted to comma",
			in:   "alice@example.com; bob@example.com",
			want: "alice@example.com, bob@example.com",
		},
		{
			name: "multiple semicolons converted",
			in:   "alice@example.com; bob@example.com; carol@example.com",
			want: "alice@example.com, bob@example.com, carol@example.com",
		},
		{
			name: "comma-separated list unchanged",
			in:   "alice@example.com, bob@example.com",
			want: "alice@example.com, bob@example.com",
		},
		{
			name: "display name with angle brackets",
			in:   "Alice <alice@example.com>; Bob <bob@example.com>",
			want: "Alice <alice@example.com>, Bob <bob@example.com>",
		},
		{
			name: "trailing comma trimmed",
			in:   "alice@example.com,",
			want: "alice@example.com",
		},
		{
			name: "leading and trailing whitespace trimmed",
			in:   "  alice@example.com  ",
			want: "alice@example.com",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "only semicolon",
			in:   ";",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeAddressSeparators(tt.in)
			if got != tt.want {
				t.Errorf("NormalizeAddressSeparators(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
