package consul

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

//Lookup - Given a consul DNS name it returns the host and port of an random instance
//of that service. Consul DNS names must end in .consul
func Lookup(service string) (*string, error) {
	return lookupWihLookupSRV(service, net.LookupSRV)
}

type lookupSRV func(service, proto, name string) (cname string, addrs []*net.SRV, err error)

func lookupWihLookupSRV(service string, lookup lookupSRV) (*string, error) {
	if !strings.HasSuffix(service, ".consul") {
		return nil, nil
	}

	_, srvs, err := lookup("", "", service)

	if err != nil {
		return nil, fmt.Errorf("Error performing SRV DNS Lookup for %s - %s", service, err.Error())
	}

	if srvs == nil || len(srvs) == 0 {
		return nil, nil
	}

	srv := srvs[rand.Intn(len(srvs))]
	target := fmt.Sprintf("%s:%d", strings.TrimSuffix(srv.Target, "."), srv.Port)

	return &target, nil

}
