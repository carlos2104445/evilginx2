package stealth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type DomainFronting struct {
	cdnManager    *CDNManager
	frontDomains  map[string]*FrontDomain
	rotationTimer *time.Timer
	config        *FrontingConfig
	mutex         sync.RWMutex
}

type FrontingConfig struct {
	EnableFronting     bool          `yaml:"enable_fronting"`
	RotationInterval   time.Duration `yaml:"rotation_interval"`
	MaxFrontDomains    int           `yaml:"max_front_domains"`
	PreferredProviders []string      `yaml:"preferred_providers"`
	BackupDomains      []string      `yaml:"backup_domains"`
}

type FrontDomain struct {
	Domain       string
	Provider     string
	DistributionID string
	OriginDomain string
	Status       string
	CreatedAt    time.Time
	LastUsed     time.Time
	RequestCount int64
}

type CDNManager struct {
	providers map[string]CDNProvider
	mutex     sync.RWMutex
}

type CDNProvider interface {
	CreateDistribution(domain string, origin string) (*Distribution, error)
	UpdateOrigin(distID, newOrigin string) error
	DeleteDistribution(distID string) error
	GetDistributionStatus(distID string) (string, error)
	ListDistributions() ([]*Distribution, error)
}

type Distribution struct {
	ID           string
	Domain       string
	Origin       string
	Status       string
	Provider     string
	CreatedAt    time.Time
	DomainName   string
	CNAME        string
}

func NewDomainFronting(config *FrontingConfig) *DomainFronting {
	df := &DomainFronting{
		cdnManager:   NewCDNManager(),
		frontDomains: make(map[string]*FrontDomain),
		config:       config,
	}

	if config.EnableFronting && config.RotationInterval > 0 {
		df.startRotationTimer()
	}

	return df
}

func NewCDNManager() *CDNManager {
	cm := &CDNManager{
		providers: make(map[string]CDNProvider),
	}

	cm.RegisterProvider("cloudflare", NewCloudFlareProvider())
	cm.RegisterProvider("aws", NewAWSCloudFrontProvider())
	cm.RegisterProvider("azure", NewAzureCDNProvider())

	return cm
}

func (cm *CDNManager) RegisterProvider(name string, provider CDNProvider) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.providers[name] = provider
}

func (cm *CDNManager) GetProvider(name string) (CDNProvider, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	provider, exists := cm.providers[name]
	return provider, exists
}

func (df *DomainFronting) CreateFrontDomain(originDomain string) (*FrontDomain, error) {
	if !df.config.EnableFronting {
		return nil, fmt.Errorf("domain fronting is disabled")
	}

	df.mutex.Lock()
	defer df.mutex.Unlock()

	if len(df.frontDomains) >= df.config.MaxFrontDomains {
		df.cleanupOldestDomain()
	}

	provider := df.selectProvider()
	frontDomain := df.generateFrontDomain()

	cdnProvider, exists := df.cdnManager.GetProvider(provider)
	if !exists {
		return nil, fmt.Errorf("CDN provider %s not found", provider)
	}

	distribution, err := cdnProvider.CreateDistribution(frontDomain, originDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create distribution: %v", err)
	}

	domain := &FrontDomain{
		Domain:         frontDomain,
		Provider:       provider,
		DistributionID: distribution.ID,
		OriginDomain:   originDomain,
		Status:         distribution.Status,
		CreatedAt:      time.Now(),
		LastUsed:       time.Now(),
		RequestCount:   0,
	}

	df.frontDomains[frontDomain] = domain
	return domain, nil
}

func (df *DomainFronting) GetActiveFrontDomain(originDomain string) (*FrontDomain, error) {
	df.mutex.RLock()
	defer df.mutex.RUnlock()

	for _, domain := range df.frontDomains {
		if domain.OriginDomain == originDomain && domain.Status == "deployed" {
			domain.LastUsed = time.Now()
			domain.RequestCount++
			return domain, nil
		}
	}

	return nil, fmt.Errorf("no active front domain found for origin %s", originDomain)
}

func (df *DomainFronting) RotateDomains() error {
	df.mutex.Lock()
	defer df.mutex.Unlock()

	var domainsToRotate []*FrontDomain
	cutoff := time.Now().Add(-df.config.RotationInterval)

	for _, domain := range df.frontDomains {
		if domain.CreatedAt.Before(cutoff) {
			domainsToRotate = append(domainsToRotate, domain)
		}
	}

	for _, domain := range domainsToRotate {
		newDomain, err := df.createReplacementDomain(domain)
		if err != nil {
			continue
		}

		err = df.deleteFrontDomain(domain.Domain)
		if err != nil {
			continue
		}

		df.frontDomains[newDomain.Domain] = newDomain
	}

	return nil
}

func (df *DomainFronting) createReplacementDomain(oldDomain *FrontDomain) (*FrontDomain, error) {
	provider := df.selectProvider()
	frontDomain := df.generateFrontDomain()

	cdnProvider, exists := df.cdnManager.GetProvider(provider)
	if !exists {
		return nil, fmt.Errorf("CDN provider %s not found", provider)
	}

	distribution, err := cdnProvider.CreateDistribution(frontDomain, oldDomain.OriginDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create replacement distribution: %v", err)
	}

	return &FrontDomain{
		Domain:         frontDomain,
		Provider:       provider,
		DistributionID: distribution.ID,
		OriginDomain:   oldDomain.OriginDomain,
		Status:         distribution.Status,
		CreatedAt:      time.Now(),
		LastUsed:       time.Now(),
		RequestCount:   0,
	}, nil
}

func (df *DomainFronting) deleteFrontDomain(domain string) error {
	frontDomain, exists := df.frontDomains[domain]
	if !exists {
		return fmt.Errorf("front domain %s not found", domain)
	}

	cdnProvider, exists := df.cdnManager.GetProvider(frontDomain.Provider)
	if !exists {
		return fmt.Errorf("CDN provider %s not found", frontDomain.Provider)
	}

	err := cdnProvider.DeleteDistribution(frontDomain.DistributionID)
	if err != nil {
		return fmt.Errorf("failed to delete distribution: %v", err)
	}

	delete(df.frontDomains, domain)
	return nil
}

func (df *DomainFronting) selectProvider() string {
	if len(df.config.PreferredProviders) > 0 {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(df.config.PreferredProviders))))
		return df.config.PreferredProviders[index.Int64()]
	}

	providers := []string{"cloudflare", "aws", "azure"}
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(providers))))
	return providers[index.Int64()]
}

func (df *DomainFronting) generateFrontDomain() string {
	prefixes := []string{
		"api", "cdn", "static", "assets", "media", "content",
		"files", "images", "js", "css", "fonts", "data",
		"cache", "edge", "global", "secure", "fast",
	}

	suffixes := []string{
		"cloudfront.net", "azureedge.net", "fastly.com",
		"maxcdn.com", "keycdn.com", "stackpathcdn.com",
		"jsdelivr.net", "unpkg.com", "cdnjs.cloudflare.com",
	}

	prefixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
	suffixIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(suffixes))))

	randomNum, _ := rand.Int(rand.Reader, big.NewInt(10000))

	return fmt.Sprintf("%s%d.%s",
		prefixes[prefixIndex.Int64()],
		randomNum.Int64(),
		suffixes[suffixIndex.Int64()],
	)
}

func (df *DomainFronting) cleanupOldestDomain() {
	var oldest *FrontDomain
	var oldestDomain string

	for domain, frontDomain := range df.frontDomains {
		if oldest == nil || frontDomain.CreatedAt.Before(oldest.CreatedAt) {
			oldest = frontDomain
			oldestDomain = domain
		}
	}

	if oldest != nil {
		df.deleteFrontDomain(oldestDomain)
	}
}

func (df *DomainFronting) startRotationTimer() {
	df.rotationTimer = time.NewTimer(df.config.RotationInterval)
	go func() {
		for {
			select {
			case <-df.rotationTimer.C:
				df.RotateDomains()
				df.rotationTimer.Reset(df.config.RotationInterval)
			}
		}
	}()
}

func (df *DomainFronting) Stop() {
	if df.rotationTimer != nil {
		df.rotationTimer.Stop()
	}
}

func (df *DomainFronting) GetStats() map[string]interface{} {
	df.mutex.RLock()
	defer df.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_domains":    len(df.frontDomains),
		"active_domains":   0,
		"pending_domains":  0,
		"failed_domains":   0,
		"providers":        make(map[string]int),
		"total_requests":   int64(0),
	}

	for _, domain := range df.frontDomains {
		switch domain.Status {
		case "deployed":
			stats["active_domains"] = stats["active_domains"].(int) + 1
		case "pending":
			stats["pending_domains"] = stats["pending_domains"].(int) + 1
		case "failed":
			stats["failed_domains"] = stats["failed_domains"].(int) + 1
		}

		providers := stats["providers"].(map[string]int)
		providers[domain.Provider]++
		stats["providers"] = providers

		stats["total_requests"] = stats["total_requests"].(int64) + domain.RequestCount
	}

	return stats
}

func (df *DomainFronting) UpdateOrigin(frontDomain, newOrigin string) error {
	df.mutex.Lock()
	defer df.mutex.Unlock()

	domain, exists := df.frontDomains[frontDomain]
	if !exists {
		return fmt.Errorf("front domain %s not found", frontDomain)
	}

	cdnProvider, exists := df.cdnManager.GetProvider(domain.Provider)
	if !exists {
		return fmt.Errorf("CDN provider %s not found", domain.Provider)
	}

	err := cdnProvider.UpdateOrigin(domain.DistributionID, newOrigin)
	if err != nil {
		return fmt.Errorf("failed to update origin: %v", err)
	}

	domain.OriginDomain = newOrigin
	return nil
}

func (df *DomainFronting) ListFrontDomains() []*FrontDomain {
	df.mutex.RLock()
	defer df.mutex.RUnlock()

	domains := make([]*FrontDomain, 0, len(df.frontDomains))
	for _, domain := range df.frontDomains {
		domains = append(domains, domain)
	}

	return domains
}
