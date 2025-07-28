package utils

import "regexp"

func RemoveHeader(text string) string {
	re := regexp.MustCompile(`(?m)^This email was sent to .+? from .+\n?`)
	return re.ReplaceAllString(text, "")
}

func RemoveHtmlHeader(html string) string {
	// Relaxed regex: match any <table> containing "This email was sent to" and ending at </table>
	re := regexp.MustCompile(`(?is)<table[^>]*>.*?This email was sent to.*?</table>`)
	cleaned := re.ReplaceAllString(html, "")

	// Optionally clean up one or more immediate trailing <br> tags or empty <div><br></div>
	cleaned = regexp.MustCompile(`(?i)(\s*<br\s*/?>\s*|<div[^>]*>\s*(<br\s*/?>)?\s*</div>)+`).ReplaceAllString(cleaned, "")

	return cleaned
}
