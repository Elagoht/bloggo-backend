package useragent

import (
	"strings"
)

type DeviceInfo struct {
	DeviceType string // "desktop", "mobile", "tablet", "unknown"
	OS         string // "Windows", "macOS", "Linux", "iOS", "Android", "Unknown"
	Browser    string // "Chrome", "Firefox", "Safari", "Edge", "Unknown"
}

// ParseUserAgent extracts device type, OS, and browser from user agent string
func ParseUserAgent(userAgent string) DeviceInfo {
	if userAgent == "" {
		return DeviceInfo{
			DeviceType: "unknown",
			OS:         "Unknown",
			Browser:    "Unknown",
		}
	}

	ua := strings.ToLower(userAgent)

	return DeviceInfo{
		DeviceType: detectDeviceType(ua),
		OS:         detectOS(ua),
		Browser:    detectBrowser(ua),
	}
}

func detectDeviceType(ua string) string {
	// Mobile indicators
	mobileIndicators := []string{
		"mobile", "android", "iphone", "ipod", "blackberry",
		"windows phone", "palm", "symbian", "iemobile", "opera mini",
	}

	for _, indicator := range mobileIndicators {
		if strings.Contains(ua, indicator) {
			return "mobile"
		}
	}

	// Tablet indicators
	tabletIndicators := []string{
		"ipad", "tablet", "kindle", "nexus 7", "nexus 9", "nexus 10",
	}

	for _, indicator := range tabletIndicators {
		if strings.Contains(ua, indicator) {
			return "tablet"
		}
	}

	// If contains typical desktop indicators or none of the above
	if strings.Contains(ua, "windows") || strings.Contains(ua, "macintosh") ||
		strings.Contains(ua, "linux") || strings.Contains(ua, "x11") {
		return "desktop"
	}

	return "unknown"
}

func detectOS(ua string) string {
	osPatterns := map[string][]string{
		"Windows": {
			"windows nt 10", "windows nt 6.3", "windows nt 6.2",
			"windows nt 6.1", "windows nt 6.0", "windows nt 5",
			"windows", "win32", "win64",
		},
		"macOS": {
			"mac os x", "macos", "macintosh", "darwin",
		},
		"iOS": {
			"iphone os", "ios", "iphone", "ipad", "ipod",
		},
		"Android": {
			"android",
		},
		"Linux": {
			"linux", "ubuntu", "debian", "fedora", "centos", "x11",
		},
		"Chrome OS": {
			"cros", "chromium os",
		},
	}

	for os, patterns := range osPatterns {
		for _, pattern := range patterns {
			if strings.Contains(ua, pattern) {
				return os
			}
		}
	}

	return "Unknown"
}

func detectBrowser(ua string) string {
	// Order matters - check more specific patterns first
	browserPatterns := map[string][]string{
		"Edge": {
			"edg/", "edge/", "edgios/", "edga/",
		},
		"Chrome": {
			"chrome/", "chromium/", "crios/",
		},
		"Firefox": {
			"firefox/", "fxios/",
		},
		"Safari": {
			"safari/", "version/", // Safari usually has both
		},
		"Opera": {
			"opera/", "opr/", "opios/",
		},
		"Internet Explorer": {
			"msie", "trident/",
		},
	}

	// Special case for Safari - needs both Safari and Version in UA
	if strings.Contains(ua, "safari/") && strings.Contains(ua, "version/") &&
		!strings.Contains(ua, "chrome/") && !strings.Contains(ua, "chromium/") {
		return "Safari"
	}

	for browser, patterns := range browserPatterns {
		if browser == "Safari" {
			continue // Already handled above
		}
		for _, pattern := range patterns {
			if strings.Contains(ua, pattern) {
				return browser
			}
		}
	}

	return "Unknown"
}

// GetDeviceTypeDistribution analyzes a slice of user agents and returns device type distribution
func GetDeviceTypeDistribution(userAgents []string) map[string]int {
	distribution := map[string]int{
		"desktop": 0,
		"mobile":  0,
		"tablet":  0,
		"unknown": 0,
	}

	for _, ua := range userAgents {
		deviceInfo := ParseUserAgent(ua)
		distribution[deviceInfo.DeviceType]++
	}

	return distribution
}

// GetOSDistribution analyzes a slice of user agents and returns OS distribution
func GetOSDistribution(userAgents []string) map[string]int {
	distribution := make(map[string]int)

	for _, ua := range userAgents {
		deviceInfo := ParseUserAgent(ua)
		distribution[deviceInfo.OS]++
	}

	return distribution
}

// GetBrowserDistribution analyzes a slice of user agents and returns browser distribution
func GetBrowserDistribution(userAgents []string) map[string]int {
	distribution := make(map[string]int)

	for _, ua := range userAgents {
		deviceInfo := ParseUserAgent(ua)
		distribution[deviceInfo.Browser]++
	}

	return distribution
}
