package utils

import "regexp"

func RemoveHeader(text string) string {
	re := regexp.MustCompile(`(?m)^This email was sent to .+? from .+\n?`)
	return re.ReplaceAllString(text, "")
}

func RemoveHtmlHeader(html string) string {
	re := regexp.MustCompile(`(?s)<table[^>]*>\s*<tr>\s*<td>\s*<div[^>]*>\s*This email was sent to\s*<a[^>]*mailto:[^"]+">[^<]+</a>\s*from\s*<a[^>]*mailto:[^"]+">[^<]+</a>\s*</div>\s*<br>\s*</td>\s*</tr>\s*</table>\s*<br>\s*`)
	return re.ReplaceAllString(html, "")
}
