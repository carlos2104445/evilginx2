package stealth

import (
	"net"
	"strings"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

type GeoLocationService struct {
	cityDB    *geoip2.Reader
	asnDB     *geoip2.Reader
	vpnRanges map[string]bool
	torNodes  map[string]bool
	cloudIPs  map[string]bool
	mutex     sync.RWMutex
}

type GeoResult struct {
	Country         string
	CountryCode     string
	City            string
	Region          string
	Latitude        float64
	Longitude       float64
	ISP             string
	ASN             uint
	Organization    string
	IsVPN           bool
	IsTor           bool
	IsCloudProvider bool
	IsProxy         bool
	ThreatLevel     int
}

type IPRange struct {
	Network     *net.IPNet
	Description string
	Type        string
}

func NewGeoLocationService() *GeoLocationService {
	service := &GeoLocationService{
		vpnRanges: make(map[string]bool),
		torNodes:  make(map[string]bool),
		cloudIPs:  make(map[string]bool),
	}

	service.loadVPNRanges()
	service.loadTorNodes()
	service.loadCloudProviderRanges()

	return service
}

func (gls *GeoLocationService) LoadDatabases(cityDBPath, asnDBPath string) error {
	if cityDBPath != "" {
		cityDB, err := geoip2.Open(cityDBPath)
		if err != nil {
			return err
		}
		gls.cityDB = cityDB
	}

	if asnDBPath != "" {
		asnDB, err := geoip2.Open(asnDBPath)
		if err != nil {
			return err
		}
		gls.asnDB = asnDB
	}

	return nil
}

func (gls *GeoLocationService) CheckIP(ipStr string) (*GeoResult, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &GeoError{Message: "Invalid IP address"}
	}

	result := &GeoResult{
		ThreatLevel: 0,
	}

	if gls.cityDB != nil {
		cityRecord, err := gls.cityDB.City(ip)
		if err == nil {
			result.Country = cityRecord.Country.Names["en"]
			result.CountryCode = cityRecord.Country.IsoCode
			result.City = cityRecord.City.Names["en"]
			if len(cityRecord.Subdivisions) > 0 {
				result.Region = cityRecord.Subdivisions[0].Names["en"]
			}
			result.Latitude = float64(cityRecord.Location.Latitude)
			result.Longitude = float64(cityRecord.Location.Longitude)
		}
	}

	if gls.asnDB != nil {
		asnRecord, err := gls.asnDB.ASN(ip)
		if err == nil {
			result.ASN = asnRecord.AutonomousSystemNumber
			result.Organization = asnRecord.AutonomousSystemOrganization
			result.ISP = asnRecord.AutonomousSystemOrganization
		}
	}

	result.IsVPN = gls.checkVPN(ipStr, result.Organization)
	result.IsTor = gls.checkTor(ipStr)
	result.IsCloudProvider = gls.checkCloudProvider(ipStr, result.Organization)
	result.IsProxy = gls.checkProxy(ipStr, result.Organization)

	result.ThreatLevel = gls.calculateThreatLevel(result)

	return result, nil
}

func (gls *GeoLocationService) checkVPN(ip, organization string) bool {
	gls.mutex.RLock()
	defer gls.mutex.RUnlock()

	if gls.vpnRanges[ip] {
		return true
	}

	vpnKeywords := []string{
		"vpn", "proxy", "tunnel", "anonymizer", "privacy",
		"nordvpn", "expressvpn", "surfshark", "cyberghost",
		"purevpn", "hotspot shield", "windscribe", "protonvpn",
		"mullvad", "private internet access", "pia",
	}

	orgLower := strings.ToLower(organization)
	for _, keyword := range vpnKeywords {
		if strings.Contains(orgLower, keyword) {
			return true
		}
	}

	return false
}

func (gls *GeoLocationService) checkTor(ip string) bool {
	gls.mutex.RLock()
	defer gls.mutex.RUnlock()

	return gls.torNodes[ip]
}

func (gls *GeoLocationService) checkCloudProvider(ip, organization string) bool {
	gls.mutex.RLock()
	defer gls.mutex.RUnlock()

	if gls.cloudIPs[ip] {
		return true
	}

	cloudKeywords := []string{
		"amazon", "aws", "microsoft", "azure", "google", "gcp",
		"digitalocean", "linode", "vultr", "hetzner", "ovh",
		"cloudflare", "fastly", "akamai", "maxcdn", "keycdn",
		"rackspace", "ibm cloud", "oracle cloud", "alibaba cloud",
	}

	orgLower := strings.ToLower(organization)
	for _, keyword := range cloudKeywords {
		if strings.Contains(orgLower, keyword) {
			return true
		}
	}

	return false
}

func (gls *GeoLocationService) checkProxy(ip, organization string) bool {
	proxyKeywords := []string{
		"proxy", "datacenter", "hosting", "server", "colocation",
		"colo", "dedicated", "vps", "virtual private server",
	}

	orgLower := strings.ToLower(organization)
	for _, keyword := range proxyKeywords {
		if strings.Contains(orgLower, keyword) {
			return true
		}
	}

	return false
}

func (gls *GeoLocationService) calculateThreatLevel(result *GeoResult) int {
	level := 0

	if result.IsTor {
		level += 50
	}

	if result.IsVPN {
		level += 30
	}

	if result.IsCloudProvider {
		level += 20
	}

	if result.IsProxy {
		level += 15
	}

	suspiciousCountries := []string{"CN", "RU", "KP", "IR"}
	for _, country := range suspiciousCountries {
		if result.CountryCode == country {
			level += 25
			break
		}
	}

	if level > 100 {
		level = 100
	}

	return level
}

func (gls *GeoLocationService) loadVPNRanges() {
	vpnRanges := []string{
		"185.220.100.0/22",
		"185.220.101.0/24",
		"199.87.154.0/24",
		"162.247.74.0/24",
		"107.189.0.0/16",
		"192.42.116.0/22",
	}

	gls.mutex.Lock()
	defer gls.mutex.Unlock()

	for _, rangeStr := range vpnRanges {
		_, network, err := net.ParseCIDR(rangeStr)
		if err != nil {
			continue
		}

		for ip := network.IP.Mask(network.Mask); network.Contains(ip); gls.incrementIP(ip) {
			gls.vpnRanges[ip.String()] = true
		}
	}
}

func (gls *GeoLocationService) loadTorNodes() {
	torNodes := []string{
		"199.87.154.255",
		"162.247.74.201",
		"107.189.29.107",
		"192.42.116.16",
		"185.220.100.240",
		"185.220.101.32",
	}

	gls.mutex.Lock()
	defer gls.mutex.Unlock()

	for _, node := range torNodes {
		gls.torNodes[node] = true
	}
}

func (gls *GeoLocationService) loadCloudProviderRanges() {
	cloudRanges := []string{
		"52.0.0.0/8",
		"54.0.0.0/8",
		"13.0.0.0/8",
		"3.0.0.0/8",
		"18.0.0.0/8",
		"35.0.0.0/8",
		"34.0.0.0/8",
		"104.154.0.0/15",
		"130.211.0.0/16",
		"146.148.0.0/17",
		"162.216.148.0/22",
		"162.222.176.0/21",
		"173.252.64.0/18",
		"173.252.70.0/24",
		"104.16.0.0/12",
		"172.64.0.0/13",
		"131.0.72.0/22",
	}

	gls.mutex.Lock()
	defer gls.mutex.Unlock()

	for _, rangeStr := range cloudRanges {
		_, network, err := net.ParseCIDR(rangeStr)
		if err != nil {
			continue
		}

		start := network.IP.Mask(network.Mask)
		end := make(net.IP, len(start))
		copy(end, start)

		for i := range end {
			end[i] |= ^network.Mask[i]
		}

		for ip := start; !ip.Equal(end); gls.incrementIP(ip) {
			gls.cloudIPs[ip.String()] = true
		}
		gls.cloudIPs[end.String()] = true
	}
}

func (gls *GeoLocationService) incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (gls *GeoLocationService) UpdateTorNodes(nodes []string) {
	gls.mutex.Lock()
	defer gls.mutex.Unlock()

	gls.torNodes = make(map[string]bool)
	for _, node := range nodes {
		gls.torNodes[node] = true
	}
}

func (gls *GeoLocationService) UpdateVPNRanges(ranges []string) {
	gls.mutex.Lock()
	defer gls.mutex.Unlock()

	gls.vpnRanges = make(map[string]bool)
	for _, rangeStr := range ranges {
		_, network, err := net.ParseCIDR(rangeStr)
		if err != nil {
			continue
		}

		for ip := network.IP.Mask(network.Mask); network.Contains(ip); gls.incrementIP(ip) {
			gls.vpnRanges[ip.String()] = true
		}
	}
}

func (gls *GeoLocationService) IsCountryAllowed(countryCode string, allowedCountries []string) bool {
	if len(allowedCountries) == 0 {
		return true
	}

	for _, allowed := range allowedCountries {
		if strings.EqualFold(countryCode, allowed) {
			return true
		}
	}

	return false
}

func (gls *GeoLocationService) Close() error {
	if gls.cityDB != nil {
		gls.cityDB.Close()
	}
	if gls.asnDB != nil {
		gls.asnDB.Close()
	}
	return nil
}

type GeoError struct {
	Message string
}

func (e *GeoError) Error() string {
	return e.Message
}
