package utils

import (
	"fmt"

	"github.com/mileusna/useragent"
)

func ParseUserAgent(userAgent string) string {
	ua := useragent.Parse(userAgent)
	if ua.IsUnknown() {
		return "Unknown"
	}

	normalized := fmt.Sprintf("%s, %s-%s", ua.Name, ua.OS, ua.OSVersion)

	if ua.Bot {
		normalized += " (Bot)"
	} else if ua.Mobile {
		normalized += " (Mobile)"
	} else if ua.Tablet {
		normalized += " (Tablet)"
	} else if ua.Desktop {
		normalized += " (Desktop)"
	}

	return normalized
}
