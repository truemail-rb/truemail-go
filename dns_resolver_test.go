package truemail

import (
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
	"github.com/stretchr/testify/assert"
)

func TestNewDnsResolver(t *testing.T) {
	t.Run("creates dnsResolver with custom gateway", func(t *testing.T) {
		connectionTimeout, dns, configuration := 42, randomDnsServer(), createConfiguration()
		configuration.ConnectionTimeout, configuration.DNS = connectionTimeout, dns
		dnsResolver := newDnsResolver(configuration)

		assert.Equal(t, connectionTimeout, dnsResolver.connectionTimeout)
		assert.Equal(t, dns, dnsResolver.dnsServer)
	})

	// Integration test with external DNS request
	t.Run("when DNS gateway not specified dnsResolver uses default system DNS gateway", func(t *testing.T) {
		dnsResolver := newDnsResolver(createConfiguration())
		resolvedHostAddresses, _ := dnsResolver.aRecords("dns.google")

		assert.True(t, isIncluded(resolvedHostAddresses, "8.8.8.8"))
	})

	// Integration test with internal DNS request
	t.Run("when DNS gateway not specified dnsResolver uses default system DNS gateway", func(t *testing.T) {
		hostName, hostAddress := randomDomain(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{hostName + ".": {A: []string{hostAddress}}}
		connectionTimeout, dns, configuration := 1, runMockDnsServer(dnsRecords), createConfiguration()
		configuration.ConnectionTimeout, configuration.DNS = connectionTimeout, dns
		dnsResolver := newDnsResolver(configuration)
		resolvedHostAddresses, _ := dnsResolver.aRecords(hostName)

		assert.True(t, isIncluded(resolvedHostAddresses, hostAddress))
	})
}

func TestDnsResolverDnsNameToHostName(t *testing.T) {
	domain, dnsResolver := randomDomain(), new(dnsResolver)

	t.Run("when domain consists dot at the end", func(t *testing.T) {
		assert.Equal(t, domain, dnsResolver.dnsNameToHostName(domain+"."))
	})

	t.Run("when domain not consists dot at the end", func(t *testing.T) {
		assert.Equal(t, domain, dnsResolver.dnsNameToHostName(domain))
	})
}

func TestDnsResolverRejectIp6Addresses(t *testing.T) {
	ip4First, ip4Second, dnsResolver := randomIpAddress(), randomIpAddress(), new(dnsResolver)
	ip4Addresses := []string{ip4First, ip4Second}

	t.Run("when slice includes ip4 and ip6 adresses", func(t *testing.T) {
		assert.Equal(t, ip4Addresses, dnsResolver.rejectIp6Addresses([]string{ip4First, randomIp6Address(), ip4Second}))
	})

	t.Run("when slice includes ip4 adresses only", func(t *testing.T) {
		assert.Equal(t, ip4Addresses, dnsResolver.rejectIp6Addresses(ip4Addresses))
	})

	t.Run("when slice includes ip6 adresses only", func(t *testing.T) {
		assert.Empty(t, dnsResolver.rejectIp6Addresses([]string{randomIp6Address()}))
	})

	t.Run("when slice is empty", func(t *testing.T) {
		assert.Empty(t, dnsResolver.rejectIp6Addresses([]string(nil)))
	})
}

func TestDnsResolverARecords(t *testing.T) {
	domain := randomDomain()

	t.Run("when target A record found", func(t *testing.T) {
		ip4First, ip4Second := randomIpAddress(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{domain + ".": {A: []string{ip4First, randomIp6Address(), ip4Second}}}
		dnsResolver := createDnsResolver(dnsRecords)
		resolvedIp4Addresses, err := dnsResolver.aRecords(domain)

		assert.Equal(t, []string{ip4First, ip4Second}, resolvedIp4Addresses)
		assert.Nil(t, err)
	})

	t.Run("when target A record not found", func(t *testing.T) {
		dnsResolver := createDnsResolverWithEpmtyRecords()
		resolvedIp4Addresses, err := dnsResolver.aRecords(domain)

		assert.Empty(t, resolvedIp4Addresses)
		assert.EqualError(t, err, dnsErrorMessage(domain))
	})
}

func TestDnsResolverARecord(t *testing.T) {
	domain := randomDomain()

	t.Run("when target A record found", func(t *testing.T) {
		ip4First, ip4Second := randomIpAddress(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{domain + ".": {A: []string{ip4First, randomIp6Address(), ip4Second}}}
		dnsResolver := createDnsResolver(dnsRecords)
		resolvedIp4Address, err := dnsResolver.aRecord(domain)

		assert.Equal(t, ip4First, resolvedIp4Address)
		assert.Nil(t, err)
	})

	t.Run("when target A record not found", func(t *testing.T) {
		dnsResolver := createDnsResolverWithEpmtyRecords()
		resolvedIp4Address, err := dnsResolver.aRecord(domain)

		assert.Empty(t, resolvedIp4Address)
		assert.EqualError(t, err, dnsErrorMessage(domain))
	})
}

func TestDnsResolverCnameRecord(t *testing.T) {
	domain := randomDomain()

	t.Run("when target CNAME record found", func(t *testing.T) {
		otherDomain := randomDomain()
		dnsRecords := map[string]mockdns.Zone{domain: {CNAME: otherDomain + "."}} // TODO: different go-mockdns behaviour for domain key
		dnsResolver := createDnsResolver(dnsRecords)
		resolvedHostName, err := dnsResolver.cnameRecord(domain)

		assert.Equal(t, otherDomain, resolvedHostName)
		assert.Nil(t, err)
	})

	t.Run("when target CNAME record not found", func(t *testing.T) {
		dnsResolver := createDnsResolverWithEpmtyRecords()
		resolvedHostName, err := dnsResolver.cnameRecord(domain)

		assert.Empty(t, resolvedHostName)
		assert.EqualError(t, err, dnsErrorMessage(domain))
	})
}

func TestDnsResolverMxRecords(t *testing.T) {
	domain := randomDomain()

	t.Run("when target MX record found", func(t *testing.T) {
		mxFirst, mxSecond := randomDomain(), randomDomain()
		dnsRecords := map[string]mockdns.Zone{
			domain + ".": {
				MX: []net.MX{
					{Host: mxFirst, Pref: 20},
					{Host: mxSecond, Pref: 10},
				},
			},
		}
		dnsResolver := createDnsResolver(dnsRecords)
		resolvedMxHostNames, err := dnsResolver.mxRecords(domain)

		assert.Equal(t, []string{mxSecond, mxFirst}, resolvedMxHostNames)
		assert.Nil(t, err)
	})

	t.Run("when target MX record not found", func(t *testing.T) {
		dnsResolver := createDnsResolverWithEpmtyRecords()
		resolvedMxHostNames, err := dnsResolver.mxRecords(domain)

		assert.Empty(t, resolvedMxHostNames)
		assert.EqualError(t, err, dnsErrorMessage(domain))
	})
}

func TestDnsResolverPtrRecords(t *testing.T) {
	hostAddress, rdnsHostAddress := "1.2.3.4", "4.3.2.1.in-addr.arpa."

	t.Run("when target RDNS record found", func(t *testing.T) {
		hostNameFirst, hostNameSecond := randomDomain(), randomDomain()
		dnsRecords := map[string]mockdns.Zone{
			rdnsHostAddress: {
				PTR: []string{
					hostNameFirst,
					hostNameSecond,
				},
			},
		}
		dnsResolver := createDnsResolver(dnsRecords)
		resolvedPtrHostNames, err := dnsResolver.ptrRecords(hostAddress)

		assert.Equal(t, []string{hostNameFirst, hostNameSecond}, resolvedPtrHostNames)
		assert.Nil(t, err)
	})

	t.Run("when target RDNS record not found", func(t *testing.T) {
		dnsResolver := createDnsResolverWithEpmtyRecords()
		resolvedPtrHostNames, err := dnsResolver.ptrRecords(hostAddress)

		assert.Empty(t, resolvedPtrHostNames)
		assert.EqualError(t, err, dnsErrorMessage(rdnsHostAddress))
	})
}