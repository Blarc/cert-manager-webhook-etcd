package etcd

import (
	"encoding/json"
	"fmt"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"strings"
)

type Config struct {
	EtcdEndpoint  string `json:"etcdEndpoint"`
	CoreDNSPrefix string `json:"coreDNSPrefix"`
}

// Returns the etcd key path for a record
// from https://github.com/kubernetes-sigs/external-dns/blob/b03da005e217b2a0a5da098cb6cd669372a89799/provider/coredns/coredns.go#L454
func (c Config) etcdKeyFor(dnsName string) string {
	domains := strings.Split(dnsName, ".")
	reverse(domains)
	return c.CoreDNSPrefix + strings.Join(domains, "/")
}

func loadConfig(cfgJSON *extapi.JSON) (Config, error) {
	cfg := Config{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}
	if cfg.EtcdEndpoint == "" {
		cfg.EtcdEndpoint = "http://localhost:2379"
	}
	if cfg.CoreDNSPrefix == "" {
		cfg.CoreDNSPrefix = "/skydns"
	}

	return cfg, nil
}

func reverse(slice []string) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
}
