package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Announcement struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Link  string `json:"link"`
}

// @Summary Get announcements
// @Description Get list of announcements from the public announcements file
// @Produce json
// @Success 200 {array} Announcement
// @Failure 500 {object} ErrorRes
// @Router /announcements [get]
func (h *Handler) GetAnnouncements(c *fiber.Ctx) error {
	resp, err := http.Get(h.Cfg.AnnouncementsURL) //nolint:noctx
	if err != nil {
		return c.Status(500).JSON(ErrorRes{Error: "failed to fetch announcements"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(500).JSON(ErrorRes{Error: fmt.Sprintf("unexpected status fetching announcements: %d", resp.StatusCode)})
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(ErrorRes{Error: "failed to read announcements"})
	}

	announcements := parseAnnouncements(string(data))
	return c.JSON(announcements)
}

var linkPattern = regexp.MustCompile(`^\[.*?\]\((https?://[^)]+)\)$`)
var headingPattern = regexp.MustCompile(`^(#{3,6}) (.+)`)
var inlineLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]+)\)`)

func parseAnnouncements(content string) []Announcement {
	lines := strings.Split(content, "\n")
	var announcements []Announcement
	var current *Announcement
	var bodyLines []string

	flush := func() {
		if current == nil {
			return
		}
		current.Body = bodyLinesToHTML(bodyLines)
		announcements = append(announcements, *current)
		current = nil
		bodyLines = nil
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.HasPrefix(line, "## ") {
			flush()
			current = &Announcement{Title: strings.TrimPrefix(line, "## ")}
			continue
		}

		if current == nil {
			continue
		}

		if m := linkPattern.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			current.Link = m[1]
			continue
		}

		bodyLines = append(bodyLines, line)
	}

	flush()
	return announcements
}

func bodyLinesToHTML(lines []string) string {
	var sb strings.Builder
	var listItems []string
	var paraLines []string

	flushList := func() {
		if len(listItems) == 0 {
			return
		}
		sb.WriteString("<ul>")
		for _, item := range listItems {
			sb.WriteString("<li>")
			sb.WriteString(item)
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul>")
		listItems = nil
	}

	flushPara := func() {
		text := strings.TrimSpace(strings.Join(paraLines, " "))
		if text != "" {
			sb.WriteString("<p>")
			sb.WriteString(inlineLinkPattern.ReplaceAllString(text, `<a href="$2">$1</a>`))
			sb.WriteString("</p>")
		}
		paraLines = nil
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			flushPara()
			flushList()
			continue
		}
		if m := headingPattern.FindStringSubmatch(line); m != nil {
			flushPara()
			flushList()
			tag := fmt.Sprintf("h%d", len(m[1]))
			sb.WriteString("<" + tag + ">")
			sb.WriteString(inlineLinkPattern.ReplaceAllString(m[2], `<a href="$2">$1</a>`))
			sb.WriteString("</" + tag + ">")
			continue
		}
		if strings.HasPrefix(line, "- ") {
			flushPara()
			item := inlineLinkPattern.ReplaceAllString(strings.TrimPrefix(line, "- "), `<a href="$2">$1</a>`)
			listItems = append(listItems, item)
			continue
		}
		flushList()
		paraLines = append(paraLines, line)
	}

	flushPara()
	flushList()

	return sb.String()
}
