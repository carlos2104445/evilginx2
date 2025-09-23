package stealth

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type FingerprintDatabase struct {
	fingerprints map[string]*Fingerprint
	mutex        sync.RWMutex
}

type Fingerprint struct {
	ID          string
	Headers     map[string]string
	UserAgent   string
	AcceptLang  string
	Encoding    string
	Connection  string
	DNT         string
	Timestamp   time.Time
	RequestCount int
	IsBot       bool
	TrustScore  int
}

type BotDetector struct {
	requestHistory map[string][]time.Time
	timingPatterns map[string]*TimingPattern
	mutex          sync.RWMutex
}

type TimingPattern struct {
	Intervals   []time.Duration
	IsRobotic   bool
	Confidence  float64
	LastRequest time.Time
}

func NewFingerprintDatabase() *FingerprintDatabase {
	return &FingerprintDatabase{
		fingerprints: make(map[string]*Fingerprint),
	}
}

func NewBotDetector() *BotDetector {
	return &BotDetector{
		requestHistory: make(map[string][]time.Time),
		timingPatterns: make(map[string]*TimingPattern),
	}
}

func (fdb *FingerprintDatabase) GenerateFingerprint(req *http.Request) *Fingerprint {
	fdb.mutex.Lock()
	defer fdb.mutex.Unlock()

	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	fingerprint := &Fingerprint{
		Headers:     headers,
		UserAgent:   req.Header.Get("User-Agent"),
		AcceptLang:  req.Header.Get("Accept-Language"),
		Encoding:    req.Header.Get("Accept-Encoding"),
		Connection:  req.Header.Get("Connection"),
		DNT:         req.Header.Get("DNT"),
		Timestamp:   time.Now(),
		RequestCount: 1,
		TrustScore:  50,
	}

	fingerprint.ID = fdb.calculateFingerprintID(fingerprint)
	
	if existing, exists := fdb.fingerprints[fingerprint.ID]; exists {
		existing.RequestCount++
		existing.Timestamp = time.Now()
		return existing
	}

	fingerprint.IsBot = fdb.analyzeFingerprint(fingerprint)
	fdb.fingerprints[fingerprint.ID] = fingerprint

	return fingerprint
}

func (fdb *FingerprintDatabase) calculateFingerprintID(fp *Fingerprint) string {
	var keys []string
	for k := range fp.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var headerString strings.Builder
	for _, k := range keys {
		headerString.WriteString(fmt.Sprintf("%s:%s;", k, fp.Headers[k]))
	}

	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		fp.UserAgent,
		fp.AcceptLang,
		fp.Encoding,
		fp.Connection,
		fp.DNT,
		headerString.String(),
	)

	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (fdb *FingerprintDatabase) analyzeFingerprint(fp *Fingerprint) bool {
	score := 0

	if fp.UserAgent == "" {
		score += 50
	}

	if fp.AcceptLang == "" {
		score += 30
	}

	if fp.Encoding == "" {
		score += 20
	}

	if len(fp.Headers) < 5 {
		score += 25
	}

	if fdb.checkCommonBotHeaders(fp) {
		score += 40
	}

	if fdb.checkUserAgentConsistency(fp) {
		score += 35
	}

	if fdb.checkHeaderOrderAnomalies(fp) {
		score += 30
	}

	return score >= 70
}

func (fdb *FingerprintDatabase) checkCommonBotHeaders(fp *Fingerprint) bool {
	botHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Originating-IP",
		"X-Remote-IP",
		"X-Remote-Addr",
		"X-ProxyUser-Ip",
	}

	count := 0
	for _, header := range botHeaders {
		if _, exists := fp.Headers[header]; exists {
			count++
		}
	}

	return count >= 2
}

func (fdb *FingerprintDatabase) checkUserAgentConsistency(fp *Fingerprint) bool {
	ua := strings.ToLower(fp.UserAgent)
	
	if strings.Contains(ua, "chrome") {
		if !strings.Contains(ua, "webkit") || !strings.Contains(ua, "safari") {
			return true
		}
	}

	if strings.Contains(ua, "firefox") {
		if strings.Contains(ua, "chrome") || strings.Contains(ua, "safari") {
			return true
		}
	}

	if strings.Contains(ua, "safari") && !strings.Contains(ua, "version") {
		return true
	}

	return false
}

func (fdb *FingerprintDatabase) checkHeaderOrderAnomalies(fp *Fingerprint) bool {
	expectedOrder := []string{
		"Host",
		"User-Agent",
		"Accept",
		"Accept-Language",
		"Accept-Encoding",
		"Connection",
	}

	headerOrder := make([]string, 0, len(fp.Headers))
	for header := range fp.Headers {
		headerOrder = append(headerOrder, header)
	}

	orderScore := 0
	for i, expected := range expectedOrder {
		if i < len(headerOrder) && headerOrder[i] == expected {
			orderScore++
		}
	}

	return float64(orderScore)/float64(len(expectedOrder)) < 0.6
}

func (bd *BotDetector) CheckRateLimit(clientIP string, maxRequestsPerMinute int) bool {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Minute)

	if _, exists := bd.requestHistory[clientIP]; !exists {
		bd.requestHistory[clientIP] = []time.Time{}
	}

	history := bd.requestHistory[clientIP]
	
	var recentRequests []time.Time
	for _, timestamp := range history {
		if timestamp.After(cutoff) {
			recentRequests = append(recentRequests, timestamp)
		}
	}

	recentRequests = append(recentRequests, now)
	bd.requestHistory[clientIP] = recentRequests

	return len(recentRequests) > maxRequestsPerMinute
}

func (bd *BotDetector) AnalyzeRequestTiming(clientIP string) *TimingPattern {
	bd.mutex.RLock()
	defer bd.mutex.RUnlock()

	if pattern, exists := bd.timingPatterns[clientIP]; exists {
		return pattern
	}

	history, exists := bd.requestHistory[clientIP]
	if !exists || len(history) < 3 {
		return &TimingPattern{
			IsRobotic:  false,
			Confidence: 0.0,
		}
	}

	intervals := make([]time.Duration, 0, len(history)-1)
	for i := 1; i < len(history); i++ {
		interval := history[i].Sub(history[i-1])
		intervals = append(intervals, interval)
	}

	pattern := &TimingPattern{
		Intervals:   intervals,
		IsRobotic:   bd.isRoboticTiming(intervals),
		Confidence:  bd.calculateTimingConfidence(intervals),
		LastRequest: history[len(history)-1],
	}

	bd.timingPatterns[clientIP] = pattern
	return pattern
}

func (bd *BotDetector) isRoboticTiming(intervals []time.Duration) bool {
	if len(intervals) < 3 {
		return false
	}

	variance := bd.calculateVariance(intervals)
	mean := bd.calculateMean(intervals)

	coefficientOfVariation := float64(variance) / float64(mean)

	if coefficientOfVariation < 0.1 {
		return true
	}

	regularPatterns := 0
	for i := 1; i < len(intervals); i++ {
		diff := intervals[i] - intervals[i-1]
		if diff < time.Millisecond*100 && diff > -time.Millisecond*100 {
			regularPatterns++
		}
	}

	return float64(regularPatterns)/float64(len(intervals)-1) > 0.8
}

func (bd *BotDetector) calculateVariance(intervals []time.Duration) time.Duration {
	if len(intervals) == 0 {
		return 0
	}

	mean := bd.calculateMean(intervals)
	var sum time.Duration

	for _, interval := range intervals {
		diff := interval - mean
		sum += time.Duration(int64(diff) * int64(diff) / int64(time.Nanosecond))
	}

	return sum / time.Duration(len(intervals))
}

func (bd *BotDetector) calculateMean(intervals []time.Duration) time.Duration {
	if len(intervals) == 0 {
		return 0
	}

	var sum time.Duration
	for _, interval := range intervals {
		sum += interval
	}

	return sum / time.Duration(len(intervals))
}

func (bd *BotDetector) calculateTimingConfidence(intervals []time.Duration) float64 {
	if len(intervals) < 3 {
		return 0.0
	}

	variance := bd.calculateVariance(intervals)
	mean := bd.calculateMean(intervals)

	if mean == 0 {
		return 0.0
	}

	coefficientOfVariation := float64(variance) / float64(mean)
	
	confidence := 1.0 - coefficientOfVariation
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

func (fdb *FingerprintDatabase) GetFingerprint(id string) (*Fingerprint, bool) {
	fdb.mutex.RLock()
	defer fdb.mutex.RUnlock()

	fp, exists := fdb.fingerprints[id]
	return fp, exists
}

func (fdb *FingerprintDatabase) UpdateTrustScore(id string, delta int) {
	fdb.mutex.Lock()
	defer fdb.mutex.Unlock()

	if fp, exists := fdb.fingerprints[id]; exists {
		fp.TrustScore += delta
		if fp.TrustScore < 0 {
			fp.TrustScore = 0
		}
		if fp.TrustScore > 100 {
			fp.TrustScore = 100
		}
	}
}

func (fdb *FingerprintDatabase) CleanupOldFingerprints(maxAge time.Duration) {
	fdb.mutex.Lock()
	defer fdb.mutex.Unlock()

	cutoff := time.Now().Add(-maxAge)
	for id, fp := range fdb.fingerprints {
		if fp.Timestamp.Before(cutoff) {
			delete(fdb.fingerprints, id)
		}
	}
}
