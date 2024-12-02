package main

import (
	"github.com/blarc/cert-manager-webhook-etcd/etcd"
	"os"
	"testing"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone = getEnv("TEST_ZONE_NAME", "example.com.")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	//

	fixture := acmetest.NewFixture(etcd.NewSolver(),
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/my-custom-solver"),
		acmetest.SetDNSServer("localhost:1053"),
		acmetest.SetUseAuthoritative(false),
	)

	fixture.RunConformance(t)

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
