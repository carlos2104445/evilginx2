package stealth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Obfuscator struct {
	randomizer *Randomizer
	config     *ObfuscationConfig
}

type ObfuscationConfig struct {
	EnableURLObfuscation     bool     `yaml:"enable_url_obfuscation"`
	EnableContentObfuscation bool     `yaml:"enable_content_obfuscation"`
	EnableTrafficRandomization bool   `yaml:"enable_traffic_randomization"`
	RandomizeJSVariables     bool     `yaml:"randomize_js_variables"`
	RandomizeCSSClasses      bool     `yaml:"randomize_css_classes"`
	InjectHTMLComments       bool     `yaml:"inject_html_comments"`
	RandomizeResourcePaths   bool     `yaml:"randomize_resource_paths"`
	DecoyParameters          []string `yaml:"decoy_parameters"`
	DecoyHeaders             []string `yaml:"decoy_headers"`
}

type Randomizer struct {
	urlPatterns     map[string]string
	jsVariables     map[string]string
	cssClasses      map[string]string
	resourcePaths   map[string]string
	htmlComments    []string
	decoyResources  []string
}

func NewObfuscator(config *ObfuscationConfig) *Obfuscator {
	return &Obfuscator{
		randomizer: NewRandomizer(),
		config:     config,
	}
}

func NewRandomizer() *Randomizer {
	return &Randomizer{
		urlPatterns:   make(map[string]string),
		jsVariables:   make(map[string]string),
		cssClasses:    make(map[string]string),
		resourcePaths: make(map[string]string),
		htmlComments: []string{
			"<!-- Analytics tracking -->",
			"<!-- Performance optimization -->",
			"<!-- Security headers -->",
			"<!-- Cache control -->",
			"<!-- SEO metadata -->",
			"<!-- Social media integration -->",
			"<!-- Third-party scripts -->",
			"<!-- Accessibility features -->",
		},
		decoyResources: []string{
			"/js/analytics.js",
			"/css/bootstrap.min.css",
			"/images/logo.png",
			"/fonts/roboto.woff2",
			"/api/health",
			"/favicon.ico",
			"/robots.txt",
			"/sitemap.xml",
		},
	}
}

func (o *Obfuscator) RandomizeURL(originalURL string) (string, error) {
	if !o.config.EnableURLObfuscation {
		return originalURL, nil
	}

	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return originalURL, err
	}

	if o.config.RandomizeResourcePaths {
		parsedURL.Path = o.randomizer.RandomizePath(parsedURL.Path)
	}

	query := parsedURL.Query()
	
	for _, param := range o.config.DecoyParameters {
		query.Add(param, o.generateRandomValue())
	}

	query.Add("_t", fmt.Sprintf("%d", time.Now().Unix()))
	query.Add("_r", o.generateRandomString(8))

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func (o *Obfuscator) ObfuscateContent(content string, contentType string) string {
	if !o.config.EnableContentObfuscation {
		return content
	}

	switch {
	case strings.Contains(contentType, "text/html"):
		return o.obfuscateHTML(content)
	case strings.Contains(contentType, "application/javascript"):
		return o.obfuscateJavaScript(content)
	case strings.Contains(contentType, "text/css"):
		return o.obfuscateCSS(content)
	default:
		return content
	}
}

func (o *Obfuscator) obfuscateHTML(content string) string {
	if o.config.InjectHTMLComments {
		content = o.injectRandomComments(content)
	}

	if o.config.RandomizeCSSClasses {
		content = o.randomizer.RandomizeCSSClasses(content)
	}

	content = o.injectDecoyElements(content)

	return content
}

func (o *Obfuscator) obfuscateJavaScript(content string) string {
	if o.config.RandomizeJSVariables {
		content = o.randomizer.RandomizeJSVariables(content)
	}

	content = o.injectDecoyJSCode(content)

	return content
}

func (o *Obfuscator) obfuscateCSS(content string) string {
	if o.config.RandomizeCSSClasses {
		content = o.randomizer.RandomizeCSSClasses(content)
	}

	content = o.injectDecoyCSS(content)

	return content
}

func (o *Obfuscator) injectRandomComments(content string) string {
	comments := o.randomizer.htmlComments
	
	for i := 0; i < 3; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(comments))))
		comment := comments[index.Int64()]
		
		insertPos, _ := rand.Int(rand.Reader, big.NewInt(int64(len(content))))
		pos := insertPos.Int64()
		
		content = content[:pos] + comment + "\n" + content[pos:]
	}

	return content
}

func (o *Obfuscator) injectDecoyElements(content string) string {
	decoyElements := []string{
		`<div style="display:none" id="analytics-tracker"></div>`,
		`<img src="/images/pixel.gif" style="display:none" alt="">`,
		`<script>/* Performance monitoring */</script>`,
		`<link rel="preload" href="/fonts/system.woff2" as="font">`,
		`<meta name="generator" content="WordPress 6.2">`,
	}

	for _, element := range decoyElements {
		if strings.Contains(content, "<head>") {
			content = strings.Replace(content, "<head>", "<head>\n"+element, 1)
		}
	}

	return content
}

func (o *Obfuscator) injectDecoyJSCode(content string) string {
	decoyCode := []string{
		`var _analytics = window._analytics || {};`,
		`function trackEvent(e,t){return!0}`,
		`var performance=window.performance||{};`,
		`(function(){var e=document.createElement("script");e.async=!0})();`,
		`window.dataLayer=window.dataLayer||[];`,
	}

	for _, code := range decoyCode {
		content = code + "\n" + content
	}

	return content
}

func (o *Obfuscator) injectDecoyCSS(content string) string {
	decoyCSS := []string{
		`.analytics-hidden{display:none!important}`,
		`.performance-marker{opacity:0;position:absolute}`,
		`.tracking-pixel{width:1px;height:1px;overflow:hidden}`,
		`.seo-helper{text-indent:-9999px}`,
	}

	for _, css := range decoyCSS {
		content = css + "\n" + content
	}

	return content
}

func (r *Randomizer) RandomizePath(path string) string {
	if cached, exists := r.resourcePaths[path]; exists {
		return cached
	}

	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if segment != "" && !strings.Contains(segment, ".") {
			randomized := r.generateRandomSegment()
			r.resourcePaths[segment] = randomized
			segments[i] = randomized
		}
	}

	newPath := strings.Join(segments, "/")
	r.resourcePaths[path] = newPath
	return newPath
}

func (r *Randomizer) RandomizeJSVariables(content string) string {
	varPattern := regexp.MustCompile(`\b(var|let|const)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=`)
	
	return varPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := varPattern.FindStringSubmatch(match)
		if len(parts) >= 3 {
			keyword := parts[1]
			varName := parts[2]
			
			if cached, exists := r.jsVariables[varName]; exists {
				return fmt.Sprintf("%s %s =", keyword, cached)
			}
			
			newName := r.generateRandomVariableName()
			r.jsVariables[varName] = newName
			return fmt.Sprintf("%s %s =", keyword, newName)
		}
		return match
	})
}

func (r *Randomizer) RandomizeCSSClasses(content string) string {
	classPattern := regexp.MustCompile(`class="([^"]*)"`)
	
	return classPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := classPattern.FindStringSubmatch(match)
		if len(parts) >= 2 {
			classes := strings.Fields(parts[1])
			var newClasses []string
			
			for _, class := range classes {
				if cached, exists := r.cssClasses[class]; exists {
					newClasses = append(newClasses, cached)
				} else {
					newClass := r.generateRandomClassName()
					r.cssClasses[class] = newClass
					newClasses = append(newClasses, newClass)
				}
			}
			
			return fmt.Sprintf(`class="%s"`, strings.Join(newClasses, " "))
		}
		return match
	})
}

func (r *Randomizer) generateRandomSegment() string {
	prefixes := []string{"api", "v1", "v2", "data", "content", "assets", "static"}
	suffixes := []string{"service", "handler", "controller", "manager", "provider"}
	
	prefixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
	suffixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(suffixes))))
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(1000))
	
	return fmt.Sprintf("%s%d%s", 
		prefixes[prefixIndex.Int64()], 
		randomNum.Int64(), 
		suffixes[suffixIndex.Int64()])
}

func (r *Randomizer) generateRandomVariableName() string {
	prefixes := []string{"_", "__", "app", "ctx", "obj", "data", "cfg", "opt"}
	suffixes := []string{"", "Obj", "Data", "Config", "Handler", "Manager"}
	
	prefixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
	suffixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(suffixes))))
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	
	return fmt.Sprintf("%s%d%s", 
		prefixes[prefixIndex.Int64()], 
		randomNum.Int64(), 
		suffixes[suffixIndex.Int64()])
}

func (r *Randomizer) generateRandomClassName() string {
	prefixes := []string{"ui", "app", "page", "content", "layout", "widget"}
	suffixes := []string{"container", "wrapper", "element", "component", "section"}
	
	prefixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
	suffixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(suffixes))))
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(1000))
	
	return fmt.Sprintf("%s-%d-%s", 
		prefixes[prefixIndex.Int64()], 
		randomNum.Int64(), 
		suffixes[suffixIndex.Int64()])
}

func (o *Obfuscator) generateRandomValue() string {
	return o.generateRandomString(12)
}

func (o *Obfuscator) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[index.Int64()]
	}
	return string(b)
}

func (o *Obfuscator) GenerateDecoyTraffic() []string {
	if !o.config.EnableTrafficRandomization {
		return []string{}
	}

	var requests []string
	
	for _, resource := range o.randomizer.decoyResources {
		randomizedURL, _ := o.RandomizeURL(resource)
		requests = append(requests, randomizedURL)
	}

	return requests
}

func (o *Obfuscator) AddRandomDelay() time.Duration {
	if !o.config.EnableTrafficRandomization {
		return 0
	}

	maxDelay, _ := rand.Int(rand.Reader, big.NewInt(2000))
	return time.Duration(maxDelay.Int64()) * time.Millisecond
}

func (o *Obfuscator) GenerateDecoyHeaders() map[string]string {
	headers := make(map[string]string)
	
	for _, header := range o.config.DecoyHeaders {
		headers[header] = o.generateRandomValue()
	}

	headers["X-Request-ID"] = o.generateRandomString(16)
	headers["X-Trace-ID"] = o.generateRandomString(32)
	headers["X-Session-ID"] = o.generateRandomString(24)

	return headers
}
