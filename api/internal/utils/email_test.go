package utils

import (
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
			input:    `<div dir="auto">Looks good.</div><div dir="auto"><br></div><div dir="auto">-- <br></div><div dir="auto"> Secured with Tuta Mail: <br></div><div dir="auto"> <a href="https://tuta.com/free-email" rel="noopener noreferrer" target="_blank">https://tuta.com/free-email</a><br></div><div dir="auto"><br></div><div dir="auto"><br></div><div dir="auto">28 Jul 2025 at 12:49 by tepid.closet03+hilje.juraj=proton.me@irelay.work:<br></div><blockquote class="tutanota_quote" style="border-left: 1px solid #93A3B8; padding-left: 10px; margin-left: 5px;"><table style="width: 100%;"><tbody><tr><td><div style="padding: 15px; background: no-repeat rgb(241, 245, 249); color: rgb(75, 85, 99); font-size: 13px; text-align: center; font-family: Arial, Helvetica, sans-serif;">This email was sent to <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:tepid.closet03@irelay.work" rel="noopener noreferrer" target="_blank">tepid.closet03@irelay.work</a> from <a style="color: rgb(59 130 246);text-decoration: none;" href="mailto:hilje.juraj@proton.me" rel="noopener noreferrer" target="_blank">hilje.juraj@proton.me</a><br></div><div dir="auto"><br></div></td></tr></tbody></table><div dir="auto"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;">Hello<br></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">您好<br></p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">这是一次测试。<br></p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'PingFang SC'">此致。<br></p></div><div style="font-family: Arial, sans-serif; font-size: 14px;"><br></div><div style="font-family: Arial, sans-serif; font-size: 14px;" class=""><div class=""><br></div><div class="">Sent with <a href="https://proton.me/mail/home" target="_blank" rel="noopener noreferrer">Proton Mail</a> secure email.<br></div></div></blockquote>`,
			expected: `<div dir="auto">Looks good.</div><div dir="auto">--</div><div dir="auto"> Secured with Tuta Mail:</div><div dir="auto"> <a href="https://tuta.com/free-email" rel="noopener noreferrer" target="_blank">https://tuta.com/free-email</a></div><div dir="auto">28 Jul 2025 at 12:49 by tepid.closet03+hilje.juraj=proton.me@irelay.work:</div><blockquote class="tutanota_quote" style="border-left: 1px solid #93A3B8; padding-left: 10px; margin-left: 5px;"><div style="font-family: Arial, sans-serif; font-size: 14px;">Hello</div><div style="font-family: Arial, sans-serif; font-size: 14px;"><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">您好</p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'Helvetica Neue'">这是一次测试。</p><p style="margin:0.0px 0.0px 0.0px 0.0px;font:13.0px 'PingFang SC'">此致。</p></div><div style="font-family: Arial, sans-serif; font-size: 14px;" class=""><div class="">Sent with <a href="https://proton.me/mail/home" target="_blank" rel="noopener noreferrer">Proton Mail</a> secure email.</div></div></blockquote>`,
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
