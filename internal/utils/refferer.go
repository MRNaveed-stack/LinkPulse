package utils

import (
	"net/url"
	"strings"
)

func NormalizeReferrer(referrer string) string {
	if referrer == "" {
		return "Direct"
	}
	u, err := url.Parse(referrer)
	if err != nil {
		return "Unknown"
	}

	host := strings.ToLower(u.Host)
	switch {
	case strings.Contains(host, "instagram.com"):
		return "Instagram"

	case strings.Contains(host, "facebook.com"):
		return "Facebook"

	case strings.Contains(host, "twitter.com"):
		return "Twitter"

	case strings.Contains(host, "linkedin.com"):
		return "LinkedIn"

	case strings.Contains(host, "tiktok.com"):
		return "TikTok"
	case strings.Contains(host, "reddit.com"):
		return "Reddit"
	case strings.Contains(host, "pinterest.com"):
		return "Pinterest"
	case strings.Contains(host, "youtube.com"):
		return "YouTube"
	case strings.Contains(host, "whatsapp.com"):
		return "WhatsApp"
	case strings.Contains(host, "telegram.org"):
		return "Telegram"
	case strings.Contains(host, "discord.com"):
		return "Discord"
	case strings.Contains(host, "snapchat.com"):
		return "Snapchat"
	case strings.Contains(host, "github.com"):
		return "GitHub"
	case strings.Contains(host, "medium.com"):
		return "Medium"
	case strings.Contains(host, "quora.com"):
		return "Quora"
	case strings.Contains(host, "tumblr.com"):
		return "Tumblr"
	default:
		return "Other"
	}
}
