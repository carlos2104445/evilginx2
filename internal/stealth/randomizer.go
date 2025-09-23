package stealth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type PatternRandomizer struct {
	requestPatterns []RequestPattern
	trafficShaper   *TrafficShaper
	config          *RandomizerConfig
}

type RandomizerConfig struct {
	EnablePatternRandomization bool          `yaml:"enable_pattern_randomization"`
	MinDelayBetweenRequests   time.Duration `yaml:"min_delay_between_requests"`
	MaxDelayBetweenRequests   time.Duration `yaml:"max_delay_between_requests"`
	FakeResourceCount         int           `yaml:"fake_resource_count"`
	DecoyRequestProbability   float64       `yaml:"decoy_request_probability"`
}

type RequestPattern struct {
	Method      string
	Path        string
	Headers     map[string]string
	Probability float64
	Delay       time.Duration
}

type TrafficShaper struct {
	patterns       []TrafficPattern
	activeShaping  bool
	lastRequest    time.Time
}

type TrafficPattern struct {
	Name        string
	Intervals   []time.Duration
	Variance    time.Duration
	Description string
}

func NewPatternRandomizer(config *RandomizerConfig) *PatternRandomizer {
	pr := &PatternRandomizer{
		trafficShaper: NewTrafficShaper(),
		config:        config,
	}

	pr.initializeRequestPatterns()
	return pr
}

func NewTrafficShaper() *TrafficShaper {
	ts := &TrafficShaper{
		activeShaping: true,
		lastRequest:   time.Now(),
	}

	ts.initializeTrafficPatterns()
	return ts
}

func (pr *PatternRandomizer) initializeRequestPatterns() {
	pr.requestPatterns = []RequestPattern{
		{
			Method: "GET",
			Path:   "/favicon.ico",
			Headers: map[string]string{
				"Accept": "image/webp,image/apng,image/*,*/*;q=0.8",
			},
			Probability: 0.8,
			Delay:       time.Millisecond * 100,
		},
		{
			Method: "GET",
			Path:   "/robots.txt",
			Headers: map[string]string{
				"Accept": "text/plain,*/*;q=0.8",
			},
			Probability: 0.3,
			Delay:       time.Millisecond * 200,
		},
		{
			Method: "GET",
			Path:   "/sitemap.xml",
			Headers: map[string]string{
				"Accept": "application/xml,text/xml,*/*;q=0.8",
			},
			Probability: 0.2,
			Delay:       time.Millisecond * 150,
		},
		{
			Method: "GET",
			Path:   "/css/bootstrap.min.css",
			Headers: map[string]string{
				"Accept": "text/css,*/*;q=0.1",
			},
			Probability: 0.9,
			Delay:       time.Millisecond * 50,
		},
		{
			Method: "GET",
			Path:   "/js/jquery.min.js",
			Headers: map[string]string{
				"Accept": "*/*",
			},
			Probability: 0.9,
			Delay:       time.Millisecond * 75,
		},
		{
			Method: "GET",
			Path:   "/api/health",
			Headers: map[string]string{
				"Accept": "application/json,*/*;q=0.8",
			},
			Probability: 0.1,
			Delay:       time.Millisecond * 300,
		},
	}
}

func (ts *TrafficShaper) initializeTrafficPatterns() {
	ts.patterns = []TrafficPattern{
		{
			Name: "human_browsing",
			Intervals: []time.Duration{
				time.Second * 2,
				time.Second * 5,
				time.Second * 1,
				time.Second * 8,
				time.Second * 3,
			},
			Variance:    time.Second * 2,
			Description: "Simulates human browsing behavior",
		},
		{
			Name: "mobile_browsing",
			Intervals: []time.Duration{
				time.Second * 3,
				time.Second * 7,
				time.Second * 2,
				time.Second * 10,
				time.Second * 4,
			},
			Variance:    time.Second * 3,
			Description: "Simulates mobile device browsing",
		},
		{
			Name: "research_pattern",
			Intervals: []time.Duration{
				time.Second * 15,
				time.Second * 30,
				time.Second * 45,
				time.Second * 20,
				time.Second * 60,
			},
			Variance:    time.Second * 10,
			Description: "Simulates research/reading behavior",
		},
	}
}

func (pr *PatternRandomizer) ShouldGenerateDecoyRequest() bool {
	if !pr.config.EnablePatternRandomization {
		return false
	}

	randomFloat, _ := rand.Int(rand.Reader, big.NewInt(1000))
	probability := float64(randomFloat.Int64()) / 1000.0

	return probability < pr.config.DecoyRequestProbability
}

func (pr *PatternRandomizer) GenerateDecoyRequest() *RequestPattern {
	if len(pr.requestPatterns) == 0 {
		return nil
	}

	var eligiblePatterns []RequestPattern
	for _, pattern := range pr.requestPatterns {
		randomFloat, _ := rand.Int(rand.Reader, big.NewInt(1000))
		probability := float64(randomFloat.Int64()) / 1000.0
		
		if probability < pattern.Probability {
			eligiblePatterns = append(eligiblePatterns, pattern)
		}
	}

	if len(eligiblePatterns) == 0 {
		return &pr.requestPatterns[0]
	}

	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(eligiblePatterns))))
	return &eligiblePatterns[index.Int64()]
}

func (pr *PatternRandomizer) CalculateRequestDelay() time.Duration {
	if !pr.config.EnablePatternRandomization {
		return 0
	}

	minDelay := pr.config.MinDelayBetweenRequests
	maxDelay := pr.config.MaxDelayBetweenRequests

	if maxDelay <= minDelay {
		return minDelay
	}

	delayRange := maxDelay - minDelay
	randomDelay, _ := rand.Int(rand.Reader, big.NewInt(int64(delayRange)))

	return minDelay + time.Duration(randomDelay.Int64())
}

func (ts *TrafficShaper) GetNextRequestDelay() time.Duration {
	if !ts.activeShaping {
		return 0
	}

	pattern := ts.selectTrafficPattern()
	baseDelay := ts.getRandomInterval(pattern)
	variance := ts.applyVariance(baseDelay, pattern.Variance)

	ts.lastRequest = time.Now()
	return variance
}

func (ts *TrafficShaper) selectTrafficPattern() TrafficPattern {
	if len(ts.patterns) == 0 {
		return TrafficPattern{
			Intervals: []time.Duration{time.Second * 2},
			Variance:  time.Second * 1,
		}
	}

	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ts.patterns))))
	return ts.patterns[index.Int64()]
}

func (ts *TrafficShaper) getRandomInterval(pattern TrafficPattern) time.Duration {
	if len(pattern.Intervals) == 0 {
		return time.Second * 2
	}

	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pattern.Intervals))))
	return pattern.Intervals[index.Int64()]
}

func (ts *TrafficShaper) applyVariance(baseDelay, variance time.Duration) time.Duration {
	if variance == 0 {
		return baseDelay
	}

	maxVariance := int64(variance)
	randomVariance, _ := rand.Int(rand.Reader, big.NewInt(maxVariance*2))
	actualVariance := time.Duration(randomVariance.Int64() - maxVariance)

	result := baseDelay + actualVariance
	if result < 0 {
		result = baseDelay / 2
	}

	return result
}

func (pr *PatternRandomizer) GenerateFakeResources() []string {
	if !pr.config.EnablePatternRandomization {
		return []string{}
	}

	resources := []string{}
	count := pr.config.FakeResourceCount

	resourceTypes := []string{
		"/images/bg-%d.jpg",
		"/css/theme-%d.css",
		"/js/module-%d.js",
		"/fonts/font-%d.woff2",
		"/api/data-%d.json",
		"/assets/icon-%d.svg",
	}

	for i := 0; i < count; i++ {
		typeIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(resourceTypes))))
		resourceType := resourceTypes[typeIndex.Int64()]
		
		resourceNum, _ := rand.Int(rand.Reader, big.NewInt(1000))
		resource := fmt.Sprintf(resourceType, resourceNum.Int64())
		
		resources = append(resources, resource)
	}

	return resources
}

func (pr *PatternRandomizer) RandomizeHeaders(baseHeaders map[string]string) map[string]string {
	headers := make(map[string]string)
	
	for k, v := range baseHeaders {
		headers[k] = v
	}

	decoyHeaders := map[string]string{
		"X-Requested-With":   "XMLHttpRequest",
		"X-CSRF-Token":       pr.generateRandomToken(),
		"X-Client-Version":   pr.generateRandomVersion(),
		"X-Request-ID":       pr.generateRandomID(),
		"X-Correlation-ID":   pr.generateRandomID(),
		"X-Session-Token":    pr.generateRandomToken(),
		"X-API-Key":          pr.generateRandomAPIKey(),
		"X-Client-Platform":  pr.getRandomPlatform(),
		"X-App-Version":      pr.generateRandomVersion(),
		"X-Device-ID":        pr.generateRandomDeviceID(),
	}

	addCount, _ := rand.Int(rand.Reader, big.NewInt(4))
	headersToAdd := int(addCount.Int64()) + 1

	headerKeys := make([]string, 0, len(decoyHeaders))
	for k := range decoyHeaders {
		headerKeys = append(headerKeys, k)
	}

	for i := 0; i < headersToAdd && i < len(headerKeys); i++ {
		keyIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(headerKeys))))
		key := headerKeys[keyIndex.Int64()]
		headers[key] = decoyHeaders[key]
		
		headerKeys = append(headerKeys[:keyIndex.Int64()], headerKeys[keyIndex.Int64()+1:]...)
	}

	return headers
}

func (pr *PatternRandomizer) generateRandomToken() string {
	return pr.generateRandomString(32)
}

func (pr *PatternRandomizer) generateRandomVersion() string {
	major, _ := rand.Int(rand.Reader, big.NewInt(10))
	minor, _ := rand.Int(rand.Reader, big.NewInt(20))
	patch, _ := rand.Int(rand.Reader, big.NewInt(50))
	
	return fmt.Sprintf("%d.%d.%d", major.Int64(), minor.Int64(), patch.Int64())
}

func (pr *PatternRandomizer) generateRandomID() string {
	return pr.generateRandomString(16)
}

func (pr *PatternRandomizer) generateRandomAPIKey() string {
	return "ak_" + pr.generateRandomString(28)
}

func (pr *PatternRandomizer) getRandomPlatform() string {
	platforms := []string{"web", "mobile", "desktop", "tablet", "api"}
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(platforms))))
	return platforms[index.Int64()]
}

func (pr *PatternRandomizer) generateRandomDeviceID() string {
	return pr.generateRandomString(20)
}

func (pr *PatternRandomizer) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[index.Int64()]
	}
	return string(b)
}

func (ts *TrafficShaper) SetActiveShaping(active bool) {
	ts.activeShaping = active
}

func (ts *TrafficShaper) AddCustomPattern(pattern TrafficPattern) {
	ts.patterns = append(ts.patterns, pattern)
}

func (ts *TrafficShaper) GetPatterns() []TrafficPattern {
	return ts.patterns
}

func (pr *PatternRandomizer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"pattern_randomization_enabled": pr.config.EnablePatternRandomization,
		"min_delay":                     pr.config.MinDelayBetweenRequests.String(),
		"max_delay":                     pr.config.MaxDelayBetweenRequests.String(),
		"fake_resource_count":           pr.config.FakeResourceCount,
		"decoy_request_probability":     pr.config.DecoyRequestProbability,
		"available_patterns":            len(pr.requestPatterns),
		"traffic_patterns":              len(pr.trafficShaper.patterns),
	}
}
