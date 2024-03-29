package truemail

import (
	"context"
	"net"
	"sort"
	"strings"
	"time"
)

type gateway interface {
	LookupHost(context.Context, string) ([]string, error)
	LookupCNAME(context.Context, string) (string, error)
	LookupMX(context.Context, string) ([]*net.MX, error)
	LookupAddr(context.Context, string) ([]string, error)
}

// dnsResolver structure. Provides possibility to send DNS requests
// via system or custom DNS gateway
type dnsResolver struct {
	connectionTimeout int
	dnsServer         string
	gateway
}

// dnsResolver builder. Creates custom resolver with
// connection timeout and DNS gateway from configuration
func newDnsResolver(configuration *Configuration) *dnsResolver {
	connectionTimeout, dnsServer := configuration.ConnectionTimeout, configuration.Dns

	return &dnsResolver{
		connectionTimeout: connectionTimeout,
		dnsServer:         dnsServer,
		gateway: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, networkProtocol, customDnsIpAddress string) (net.Conn, error) {
				dialer := net.Dialer{Timeout: time.Duration(connectionTimeout) * time.Second}
				if dnsServer != emptyString {
					customDnsIpAddress = dnsServer
				}
				return dialer.DialContext(ctx, networkProtocol, customDnsIpAddress)
			},
		},
	}
}

// dnsResolver methods

// Helper method. Removes last dot from dns hostname representation, example.com. => example.com
func (dnsResolver *dnsResolver) dnsNameToHostName(dnsName string) string {
	return strings.TrimSuffix(dnsName, ".")
}

// Helper method. Filter out ipv6 ip addresses from mixed collection
func (dnsResolver *dnsResolver) rejectIp6Addresses(ipAddresses []string) (ip4Addresses []string) {
	for _, ipAddress := range ipAddresses {
		if matchRegex(ipAddress, regexIpAddressPattern) && ipAddress != "0.0.0.0" { // ipv6 can be converted to 0.0.0.0
			ip4Addresses = append(ip4Addresses, ipAddress)
		}
		continue
	}

	return ip4Addresses
}

// Returns all A records by hostname
func (dnsResolver *dnsResolver) aRecords(hostName string) ([]string, error) {
	ipAddresses, err := dnsResolver.gateway.LookupHost(context.Background(), hostName)
	if err != nil {
		return []string{}, wrapDnsError(err)
	}

	return dnsResolver.rejectIp6Addresses(ipAddresses), nil
}

// Returns first A record by hostname
func (dnsResolver *dnsResolver) aRecord(hostName string) (string, error) {
	ipAddresses, err := dnsResolver.aRecords(hostName)
	if err != nil {
		return emptyString, err
	}

	return ipAddresses[0], nil
}

// Returns CNAME record by hostname for case when CNAME is different as hostname only
func (dnsResolver *dnsResolver) cnameRecord(hostName string) (resolvedHostName string, err error) {
	cName, err := dnsResolver.gateway.LookupCNAME(context.Background(), hostName)
	if err != nil {
		return resolvedHostName, wrapDnsError(err)
	}

	cName = dnsResolver.dnsNameToHostName(cName)
	if cName != hostName {
		resolvedHostName = cName
	}

	return resolvedHostName, nil
}

// Returns MX records priorities and hostnames sorted by record priority
func (dnsResolver *dnsResolver) mxRecords(hostName string) (priorities []uint16, hostNames []string, err error) {
	mxRecords, err := dnsResolver.gateway.LookupMX(context.Background(), hostName)
	if err != nil {
		return priorities, hostNames, wrapDnsError(err)
	}

	// sorting MX records by priority
	sort.SliceStable(mxRecords, func(i, j int) bool {
		return mxRecords[i].Pref < mxRecords[j].Pref
	})

	for _, mxRecord := range mxRecords {
		priorities = append(priorities, mxRecord.Pref)
		hostNames = append(hostNames, dnsResolver.dnsNameToHostName(mxRecord.Host))
	}

	return priorities, hostNames, nil
}

// Returns PTR records by host address
func (dnsResolver *dnsResolver) ptrRecords(hostAddress string) (hostNames []string, err error) {
	hostNames, err = dnsResolver.gateway.LookupAddr(context.Background(), hostAddress)
	if err != nil {
		return hostNames, wrapDnsError(err)
	}

	for index, hostName := range hostNames {
		hostNames[index] = dnsResolver.dnsNameToHostName(hostName)
	}

	return hostNames, nil
}
