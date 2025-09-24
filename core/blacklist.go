package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/kgretzky/evilginx2/log"
)

type BlockIP struct {
	ipv4 net.IP
	mask *net.IPNet
}

type Blacklist struct {
	ips        map[string]*BlockIP
	masks      []*BlockIP
	configPath string
	verbose    bool
	mu         sync.RWMutex
}

func NewBlacklist(path string) (*Blacklist, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bl := &Blacklist{
		ips:        make(map[string]*BlockIP),
		configPath: path,
		verbose:    true,
	}

	fs := bufio.NewScanner(f)
	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		l := fs.Text()
		// remove comments
		if n := strings.Index(l, ";"); n > -1 {
			l = l[:n]
		}
		l = strings.Trim(l, " ")

		if len(l) > 0 {
			if strings.Contains(l, "/") {
				ipv4, mask, err := net.ParseCIDR(l)
				if err == nil {
					bl.masks = append(bl.masks, &BlockIP{ipv4: ipv4, mask: mask})
				} else {
					log.Error("blacklist: invalid ip/mask address: %s", l)
				}
			} else {
				ipv4 := net.ParseIP(l)
				if ipv4 != nil {
					bl.ips[ipv4.String()] = &BlockIP{ipv4: ipv4, mask: nil}
				} else {
					log.Error("blacklist: invalid ip address: %s", l)
				}
			}
		}
	}

	log.Info("blacklist: loaded %d ip addresses and %d ip masks", len(bl.ips), len(bl.masks))
	return bl, nil
}

func (bl *Blacklist) GetStats() (int, int) {
	return len(bl.ips), len(bl.masks)
}

func (bl *Blacklist) AddIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}
	
	if bl.IsBlacklisted(ip) {
		return nil
	}

	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return fmt.Errorf("invalid ip address: %s", ip)
	}

	bl.mu.Lock()
	bl.ips[ipv4.String()] = &BlockIP{ipv4: ipv4, mask: nil}
	bl.mu.Unlock()

	f, err := os.OpenFile(bl.configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open blacklist file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(ipv4.String() + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to blacklist file: %w", err)
	}

	return nil
}

func (bl *Blacklist) IsBlacklisted(ip string) bool {
	if ip == "" {
		return false
	}
	
	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return false
	}

	bl.mu.RLock()
	defer bl.mu.RUnlock()
	
	if _, ok := bl.ips[ip]; ok {
		return true
	}
	for _, m := range bl.masks {
		if m.mask != nil && m.mask.Contains(ipv4) {
			return true
		}
	}
	return false
}

func (bl *Blacklist) SetVerbose(verbose bool) {
	bl.verbose = verbose
}

func (bl *Blacklist) IsVerbose() bool {
	return bl.verbose
}

func (bl *Blacklist) IsWhitelisted(ip string) bool {
	if ip == "" {
		return false
	}
	
	ipv4 := net.ParseIP(ip)
	if ipv4 == nil {
		return false
	}
	
	if ip == "127.0.0.1" || ip == "::1" {
		return true
	}
	
	if ipv4.IsLoopback() || ipv4.IsPrivate() {
		return true
	}
	
	return false
}
