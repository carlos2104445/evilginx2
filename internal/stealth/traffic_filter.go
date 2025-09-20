package stealth

import (
	"net/http"
	"regexp"
	"strings"
)

type TrafficFilter struct {
	botDetector    *BotDetector
	geoService     *GeoLocationService
	fingerprintDB  *FingerprintDatabase
	config         *FilterConfig
}

type FilterConfig struct {
	EnableBotFiltering   bool     `yaml:"enable_bot_filtering"`
	AllowedCountries     []string `yaml:"allowed_countries"`
	BlockVPN             bool     `yaml:"block_vpn"`
	BlockTor             bool     `yaml:"block_tor"`
	BlockCloudProviders  bool     `yaml:"block_cloud_providers"`
	CustomBlockedUAs     []string `yaml:"custom_blocked_uas"`
	MaxRequestsPerMinute int      `yaml:"max_requests_per_minute"`
}

type FilterResult struct {
	ShouldBlock bool
	Reason      string
	Score       int
}

func NewTrafficFilter(config *FilterConfig) *TrafficFilter {
	return &TrafficFilter{
		botDetector:   NewBotDetector(),
		geoService:    NewGeoLocationService(),
		fingerprintDB: NewFingerprintDatabase(),
		config:        config,
	}
}

func (tf *TrafficFilter) ShouldBlock(req *http.Request) (*FilterResult, error) {
	result := &FilterResult{
		ShouldBlock: false,
		Reason:      "",
		Score:       0,
	}

	if !tf.config.EnableBotFiltering {
		return result, nil
	}

	clientIP := tf.getClientIP(req)
	userAgent := req.Header.Get("User-Agent")

	if tf.checkBotUserAgent(userAgent) {
		result.ShouldBlock = true
		result.Reason = "Bot user agent detected"
		result.Score += 100
		return result, nil
	}

	if tf.checkSecurityScannerSignatures(req) {
		result.ShouldBlock = true
		result.Reason = "Security scanner detected"
		result.Score += 100
		return result, nil
	}

	if tf.checkAutomatedToolPatterns(req) {
		result.ShouldBlock = true
		result.Reason = "Automated tool detected"
		result.Score += 90
		return result, nil
	}

	if tf.checkHeadlessBrowser(req) {
		result.ShouldBlock = true
		result.Reason = "Headless browser detected"
		result.Score += 80
		return result, nil
	}

	geoResult, err := tf.geoService.CheckIP(clientIP)
	if err == nil {
		if tf.shouldBlockByGeo(geoResult) {
			result.ShouldBlock = true
			result.Reason = "Blocked by geolocation policy"
			result.Score += 70
			return result, nil
		}
	}

	if tf.checkRateLimiting(clientIP) {
		result.ShouldBlock = true
		result.Reason = "Rate limit exceeded"
		result.Score += 60
		return result, nil
	}

	behaviorScore := tf.analyzeBehavior(req, clientIP)
	result.Score += behaviorScore

	if result.Score >= 50 {
		result.ShouldBlock = true
		result.Reason = "Suspicious behavior pattern"
	}

	return result, nil
}

func (tf *TrafficFilter) getClientIP(req *http.Request) string {
	forwarded := req.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := req.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return strings.Split(req.RemoteAddr, ":")[0]
}

func (tf *TrafficFilter) checkBotUserAgent(userAgent string) bool {
	botPatterns := []string{
		"(?i)bot",
		"(?i)crawler",
		"(?i)spider",
		"(?i)scraper",
		"(?i)curl",
		"(?i)wget",
		"(?i)python-requests",
		"(?i)go-http-client",
		"(?i)nessus",
		"(?i)openvas",
		"(?i)nikto",
		"(?i)sqlmap",
		"(?i)burp",
		"(?i)zap",
		"(?i)nuclei",
	}

	for _, pattern := range botPatterns {
		matched, _ := regexp.MatchString(pattern, userAgent)
		if matched {
			return true
		}
	}

	for _, customUA := range tf.config.CustomBlockedUAs {
		if strings.Contains(strings.ToLower(userAgent), strings.ToLower(customUA)) {
			return true
		}
	}

	return false
}

func (tf *TrafficFilter) checkSecurityScannerSignatures(req *http.Request) bool {
	scannerHeaders := []string{
		"X-Scanner",
		"X-Forwarded-Proto",
		"X-Originating-IP",
		"X-Remote-IP",
		"X-Remote-Addr",
	}

	for _, header := range scannerHeaders {
		if req.Header.Get(header) != "" {
			return true
		}
	}

	suspiciousParams := []string{
		"<script>",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"../",
		"..\\",
		"union select",
		"' or 1=1",
		"' or '1'='1",
	}

	queryString := req.URL.RawQuery
	for _, param := range suspiciousParams {
		if strings.Contains(strings.ToLower(queryString), strings.ToLower(param)) {
			return true
		}
	}

	return false
}

func (tf *TrafficFilter) checkAutomatedToolPatterns(req *http.Request) bool {
	userAgent := req.Header.Get("User-Agent")
	
	automatedPatterns := []string{
		"(?i)python",
		"(?i)java",
		"(?i)perl",
		"(?i)ruby",
		"(?i)php",
		"(?i)node",
		"(?i)axios",
		"(?i)okhttp",
		"(?i)apache-httpclient",
	}

	for _, pattern := range automatedPatterns {
		matched, _ := regexp.MatchString(pattern, userAgent)
		if matched {
			return true
		}
	}

	acceptHeader := req.Header.Get("Accept")
	if acceptHeader == "*/*" || acceptHeader == "" {
		return true
	}

	if req.Header.Get("Accept-Language") == "" {
		return true
	}

	return false
}

func (tf *TrafficFilter) checkHeadlessBrowser(req *http.Request) bool {
	userAgent := req.Header.Get("User-Agent")
	
	headlessPatterns := []string{
		"(?i)headlesschrome",
		"(?i)phantomjs",
		"(?i)selenium",
		"(?i)webdriver",
		"(?i)puppeteer",
		"(?i)playwright",
	}

	for _, pattern := range headlessPatterns {
		matched, _ := regexp.MatchString(pattern, userAgent)
		if matched {
			return true
		}
	}

	if strings.Contains(userAgent, "Chrome") && !strings.Contains(userAgent, "Safari") {
		return true
	}

	webglHeader := req.Header.Get("X-WebGL-Vendor")
	if webglHeader == "Brian Paul" || webglHeader == "Mesa" {
		return true
	}

	return false
}

func (tf *TrafficFilter) shouldBlockByGeo(geoResult *GeoResult) bool {
	if len(tf.config.AllowedCountries) > 0 {
		allowed := false
		for _, country := range tf.config.AllowedCountries {
			if country == geoResult.Country {
				allowed = true
				break
			}
		}
		if !allowed {
			return true
		}
	}

	if tf.config.BlockVPN && geoResult.IsVPN {
		return true
	}

	if tf.config.BlockTor && geoResult.IsTor {
		return true
	}

	if tf.config.BlockCloudProviders && geoResult.IsCloudProvider {
		return true
	}

	return false
}

func (tf *TrafficFilter) checkRateLimiting(clientIP string) bool {
	if tf.config.MaxRequestsPerMinute <= 0 {
		return false
	}

	return tf.botDetector.CheckRateLimit(clientIP, tf.config.MaxRequestsPerMinute)
}

func (tf *TrafficFilter) analyzeBehavior(req *http.Request, clientIP string) int {
	score := 0

	if req.Header.Get("Referer") == "" {
		score += 10
	}

	if len(req.Header) < 5 {
		score += 15
	}

	if req.Method != "GET" && req.Method != "POST" {
		score += 20
	}

	timing := tf.botDetector.AnalyzeRequestTiming(clientIP)
	if timing.IsRobotic {
		score += 25
	}

	return score
}
