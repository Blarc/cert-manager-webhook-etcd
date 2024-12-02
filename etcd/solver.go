package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"time"
)

func NewSolver() webhook.Solver {
	return &Solver{}
}

type Solver struct {
	etcdClient *etcdClientV3.Client
}

func (s *Solver) Name() string {
	return "coredns-etcd"
}

func (s *Solver) Present(ch *v1alpha1.ChallengeRequest) error {
	klog.Infof("Presenting txt record: %v %v", ch.ResolvedFQDN, ch.ResolvedZone)
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		klog.Errorf("Load configuration error : %v", err)
		return err
	}

	klog.Infof("Decoded configuration %v", cfg)

	s.NewEtcdClient(cfg)
	defer s.etcdClient.Close()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	value, err := json.Marshal(TxtRecord{
		TTL:  uint32(60),
		Text: ch.Key,
	})
	if err != nil {
		return err
	}

	key := cfg.etcdKeyFor(ch.ResolvedFQDN)
	_, err = s.etcdClient.Put(ctx, key, string(value))
	if err != nil {
		return err
	}

	klog.Infof("Presented txt record %v", ch.ResolvedFQDN)
	return nil
}

func (s *Solver) CleanUp(ch *v1alpha1.ChallengeRequest) (err error) {
	fmt.Printf("Cleaning up a TXT record for dns01\n")
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}
	fmt.Printf("Decoded configuration %v", cfg)

	s.NewEtcdClient(cfg)
	defer s.etcdClient.Close()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	key := cfg.etcdKeyFor(ch.ResolvedFQDN)
	_, err = s.etcdClient.Delete(ctx, key, etcdClientV3.WithPrefix())

	return nil
}

func (s *Solver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}

func (s *Solver) NewEtcdClient(cfg Config) (err error) {
	etcdClientConfig := &etcdClientV3.Config{Endpoints: []string{cfg.EtcdEndpoint}}
	s.etcdClient, err = etcdClientV3.New(*etcdClientConfig)
	if err != nil {
		return err
	}
	return nil
}

type TxtRecord struct {
	TTL  uint32 `json:"ttl,omitempty"`
	Text string `json:"text,omitempty"`
}
