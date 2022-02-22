package truemail

import (
	"fmt"
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
	"github.com/stretchr/testify/assert"
)

func TestValidationMxCheck(t *testing.T) {
	targetUserName, targetHostName := "niña@", "mañana.com"
	targetEmail := targetUserName + targetHostName
	configuration := createConfiguration()

	t.Run("MX validation: successful, servers extracted by MX records resolver", func(t *testing.T) {
		mxHostnameFirst, mxHostnameSecond := randomDnsHostName(), randomDnsHostName()
		resolvedIpAddressFirst, resolvedIpAddressSecond := randomIpAddress(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(targetHostName)): {
				MX: []net.MX{
					{Host: mxHostnameFirst, Pref: uint16(5)},
					{Host: mxHostnameSecond, Pref: uint16(10)},
				},
			},
			mxHostnameFirst: {
				A: []string{
					resolvedIpAddressFirst,
				},
			},
			mxHostnameSecond: {
				A: []string{
					resolvedIpAddressFirst,
					resolvedIpAddressSecond,
				},
			},
		}
		configuration.Dns = runMockDnsServer(dnsRecords)
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Equal(t, []string{resolvedIpAddressFirst, resolvedIpAddressSecond}, validatorResult.MailServers)
	})

	t.Run("MX validation: successful, servers extracted by CNAME record resolver", func(t *testing.T) {
		resolvedCnameHostName := randomDnsHostName()
		resolvedAHostAddress := "1.2.3.4"
		resolvedARdnsHostAddress := "4.3.2.1.in-addr.arpa."
		resolvedPtrHostNameFirst, resolvedPtrHostNameSecond := randomDnsHostName(), randomDnsHostName()
		resolvedMxHostnameFirst, resolvedMxHostnameSecond := randomDnsHostName(), randomDnsHostName()
		resolvedIpAddressFirst, resolvedIpAddressSecond := randomIpAddress(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(targetHostName)): {
				CNAME: resolvedCnameHostName,
			},
			resolvedCnameHostName: {
				A: []string{resolvedAHostAddress, randomIpAddress(), randomIp6Address()},
			},
			resolvedARdnsHostAddress: {
				PTR: []string{resolvedPtrHostNameFirst, resolvedPtrHostNameSecond},
			},
			resolvedPtrHostNameFirst: {
				MX: []net.MX{
					{Host: resolvedMxHostnameFirst, Pref: uint16(10)},
					{Host: resolvedMxHostnameSecond, Pref: uint16(5)},
				},
			},
			resolvedPtrHostNameSecond: {
				MX: []net.MX{
					{Host: resolvedMxHostnameFirst, Pref: uint16(5)},
				},
			},
			resolvedMxHostnameFirst: {
				A: []string{resolvedIpAddressFirst, randomIp6Address()},
			},
			resolvedMxHostnameSecond: {
				A: []string{resolvedIpAddressSecond},
			},
		}
		configuration.Dns = runMockDnsServer(dnsRecords)
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Equal(t, []string{resolvedIpAddressSecond, resolvedIpAddressFirst}, validatorResult.MailServers)
	})

	t.Run("MX validation: successful, servers extracted by A record resolver", func(t *testing.T) {
		resolvedIpAddressFirst := randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(targetHostName)): {
				A: []string{resolvedIpAddressFirst, randomIpAddress()},
			},
		}
		configuration.Dns = runMockDnsServer(dnsRecords)
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Equal(t, []string{resolvedIpAddressFirst}, validatorResult.MailServers)
	})

	t.Run("MX validation: failure, servers not found using MX, CNAME and A record resolvers", func(t *testing.T) {
		configuration.Dns = runMockDnsServer(map[string]mockdns.Zone{})
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"mx": mxErrorContext}, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Empty(t, validatorResult.MailServers)
	})

	t.Run("MX validation: failure, null MX record found using MX resolver", func(t *testing.T) {
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(targetHostName)): {
				MX: []net.MX{
					{Host: ".", Pref: uint16(0)},
				},
				A: []string{randomIpAddress()},
			},
		}
		configuration.Dns = runMockDnsServer(dnsRecords)
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"mx": mxErrorContext}, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Empty(t, validatorResult.MailServers)
	})

	t.Run("MX validation: failure, servers not found using MX resolver, not RFC MX lookup enabled", func(t *testing.T) {
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(targetHostName)): {
				A: []string{randomIpAddress()},
			},
		}
		configuration.Dns, configuration.NotRfcMxLookupFlow = runMockDnsServer(dnsRecords), true
		validatorResult := createSuccessfulValidatorResult(targetEmail, configuration)
		new(validationMx).check(validatorResult)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"mx": mxErrorContext}, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
		assert.Equal(t, punycodeDomain(targetHostName), validatorResult.punycodeDomain)
		assert.Equal(t, targetUserName+punycodeDomain(targetHostName), validatorResult.punycodeEmail)
		assert.Empty(t, validatorResult.MailServers)
	})
}

func TestValidationMxPunycodeDomain(t *testing.T) {
	t.Run("returns domain punycode representation", func(t *testing.T) {
		internationalizedDomain := "mañana.cøm"
		asciiDomain := "xn--maana-pta.xn--cm-lka"

		assert.Equal(t, asciiDomain, new(validationMx).punycodeDomain(internationalizedDomain))
	})
}

func TestValidationMxSetValidatorResultPunycodeRepresentation(t *testing.T) {
	t.Run("returns domain punycode representation", func(t *testing.T) {
		internationalizedUser, internationalizedDomain := "niña", "mañana.cøm"
		internationalizedEmail := internationalizedUser + "@" + internationalizedDomain
		validatorResult := &validatorResult{Email: internationalizedEmail}
		validation := &validationMx{result: validatorResult}
		validation.setValidatorResultPunycodeRepresentation()
		asciiDomain := "xn--maana-pta.xn--cm-lka"

		assert.Equal(t, internationalizedUser+"@"+asciiDomain, validatorResult.punycodeEmail)
		assert.Equal(t, asciiDomain, validatorResult.punycodeDomain)
	})
}

func TestValidationMxInitDnsResolver(t *testing.T) {
	// Integration test of validation resolver with internal DNS request

	t.Run("initializes DNS resolver from configuration settings", func(t *testing.T) {
		hostName, hostAddress := randomDomain(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{hostName + ".": {A: []string{hostAddress}}}
		dns, configuration := runMockDnsServer(dnsRecords), createConfiguration()
		configuration.Dns = dns
		validation := &validationMx{result: &validatorResult{Configuration: configuration}}
		validation.initDnsResolver()
		resolvedHostAddress, _ := validation.resolver.aRecord(hostName)

		assert.Equal(t, hostAddress, resolvedHostAddress)
	})
}

func TestValidationMxIsMailServerNotFound(t *testing.T) {
	t.Run("when mail servers none", func(t *testing.T) {
		validation := &validationMx{result: &validatorResult{}}

		assert.True(t, validation.isMailServerNotFound())
	})

	t.Run("when mail servers exist", func(t *testing.T) {
		validation := &validationMx{result: &validatorResult{MailServers: []string{randomIp6Address()}}}

		assert.False(t, validation.isMailServerNotFound())
	})
}

func TestValidationMxIsMailServerFound(t *testing.T) {
	t.Run("when mail servers exist", func(t *testing.T) {
		validation := &validationMx{result: &validatorResult{MailServers: []string{randomIp6Address()}}}

		assert.True(t, validation.isMailServerFound())
	})

	t.Run("when mail servers none", func(t *testing.T) {
		validation := &validationMx{result: &validatorResult{}}

		assert.False(t, validation.isMailServerFound())
	})
}

func TestValidationMxFetchTargetHosts(t *testing.T) {
	t.Run("addes only uniques hosts with keeping sequence of grabbed hosts", func(t *testing.T) {
		hostFirst, hostSecond, hostThird := randomIpAddress(), randomIpAddress(), randomIpAddress()
		validatorResult := &validatorResult{MailServers: []string{hostFirst}}
		hosts := []string{hostSecond, hostFirst, hostThird, hostSecond}
		validation := &validationMx{result: validatorResult}
		validation.fetchTargetHosts(hosts...)

		assert.Equal(t, []string{hostFirst, hostSecond, hostThird}, validatorResult.MailServers)
	})
}

func TestIsConnectionAttemptsAvailable(t *testing.T) {
	t.Run("when connection attempts are available", func(t *testing.T) {
		assert.True(t, new(validationMx).isConnectionAttemptsAvailable(1))
	})

	t.Run("when connection attempts aren't available", func(t *testing.T) {
		assert.False(t, new(validationMx).isConnectionAttemptsAvailable(0))
	})
}

func TestValidationMxIsDnsNotFoundError(t *testing.T) {
	t.Run("when DNS not found error", func(t *testing.T) {
		assert.True(t, new(validationMx).isDnsNotFoundError(createDnsNotFoundError()))
	})

	t.Run("when another error", func(t *testing.T) {
		assert.False(t, new(validationMx).isDnsNotFoundError(new(validationError)))
	})
}

func TestValidationMxIsNullMxError(t *testing.T) {
	t.Run("when null MX found error", func(t *testing.T) {
		assert.True(t, new(validationMx).isNullMxError(&validationError{isNullMxFound: true}))
	})

	t.Run("when another error", func(t *testing.T) {
		assert.False(t, new(validationMx).isNullMxError(new(validationError)))
	})
}

func TestValidationMxARecords(t *testing.T) {
	method, email, connectionAttempts, ipAddresses := "aRecords", randomEmail(), 2, []string{randomIpAddress()}
	hostName := emailDomain(email)
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: email, ConnectionAttempts: connectionAttempts})

	t.Run("when A records was found during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Once().Return(ipAddresses, nil)
		resolvedIpAddresses, err := validation.aRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, ipAddresses, resolvedIpAddresses)
		assert.NoError(t, err)
	})

	t.Run("when A records was not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Once().Return([]string{}, createDnsNotFoundError())
		resolvedIpAddresses, err := validation.aRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when connection issues", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Times(connectionAttempts).Return([]string{}, new(validationError))
		resolvedIpAddresses, err := validation.aRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedIpAddresses)
		assert.Error(t, err)
	})
}

func TestValidationMxHostsFromMxRecords(t *testing.T) {
	email, connectionAttempts := randomEmail(), 2
	ipAddresses, hostNames, priorities := []string{randomIpAddress()}, []string{randomDomain()}, []uint16{5}
	hostName := emailDomain(email)
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: email, ConnectionAttempts: connectionAttempts})

	t.Run("when MX records was found during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return(priorities, hostNames, nil)
		resolver.On("aRecords", hostNames[0]).Once().Return(ipAddresses, nil)
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, ipAddresses, resolvedIpAddresses)
		assert.NoError(t, err)
	})

	t.Run("when null MX records was found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}
		errorMessage := fmt.Sprintf("%v includes null MX record", hostName)

		resolver.On("mxRecords", hostName).Once().Return([]uint16{0}, []string{""}, nil)
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.EqualError(t, err, errorMessage)
		assert.True(t, isNullMxError(err))
	})

	t.Run("when MX records was not found, mxRecords error", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{}, []string(nil), createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when MX records was not found, aRecords error", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return(priorities, hostNames, nil)
		resolver.On("aRecords", hostNames[0]).Once().Return([]string{}, createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when connection issues, mxRecords error", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Times(connectionAttempts).Return([]uint16{}, []string(nil), new(validationError))
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when connection issues, aRecords error", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return(priorities, hostNames, nil)
		resolver.On("aRecords", hostNames[0]).Times(connectionAttempts).Return([]string{}, new(validationError))
		resolvedIpAddresses, err := validation.hostsFromMxRecords(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})
}

func TestValidationMxHostFromARecord(t *testing.T) {
	method, email, connectionAttempts, ipAddress := "aRecord", randomEmail(), 2, randomIpAddress()
	hostName := emailDomain(email)
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: email, ConnectionAttempts: connectionAttempts})

	t.Run("when A record was found during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Once().Return(ipAddress, nil)
		resolvedIpAddress, err := validation.hostFromARecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, ipAddress, resolvedIpAddress)
		assert.NoError(t, err)
	})

	t.Run("when A record was not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Once().Return("", createDnsNotFoundError())
		resolvedIpAddress, err := validation.hostFromARecord(hostName)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedIpAddress)
		assert.Error(t, err)
	})

	t.Run("when connection issues", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostName).Times(connectionAttempts).Return("", new(validationError))
		resolvedIpAddress, err := validation.hostFromARecord(hostName)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedIpAddress)
		assert.Error(t, err)
	})
}

func TestValidationMxPtrRecords(t *testing.T) {
	method, email, connectionAttempts, hostAddress, hostNames := "ptrRecords", randomEmail(), 2, randomDomain(), []string{randomIpAddress()}
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: email, ConnectionAttempts: connectionAttempts})

	t.Run("when PTR records was found during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostAddress).Once().Return(hostNames, nil)
		resolvedHostNames, err := validation.ptrRecords(hostAddress)
		resolver.AssertExpectations(t)
		assert.Equal(t, hostNames, resolvedHostNames)
		assert.NoError(t, err)
	})

	t.Run("when PTR records was not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostAddress).Once().Return([]string{}, createDnsNotFoundError())
		resolvedHostNames, err := validation.ptrRecords(hostAddress)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedHostNames)
		assert.Error(t, err)
	})

	t.Run("when connection issues", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On(method, hostAddress).Times(connectionAttempts).Return([]string{}, new(validationError))
		resolvedHostNames, err := validation.ptrRecords(hostAddress)
		resolver.AssertExpectations(t)
		assert.Empty(t, resolvedHostNames)
		assert.Error(t, err)
	})
}

func TestValidationMxHostsFromCnameRecord(t *testing.T) {
	email, connectionAttempts := randomEmail(), 2
	hostName := emailDomain(email)
	resolvedHostNameByCname := randomDomain()
	resolvedIpAddressByARecord := randomIpAddress()
	resolvedHostNameByPtrRecord := randomDomain()
	resolvedHostNamesByPtrRecords := []string{resolvedHostNameByPtrRecord}
	resolvedHostNameByMxRecord := randomDomain()
	priorities, resolvedHostNamesByMxRecords := []uint16{5}, []string{resolvedHostNameByMxRecord}
	resolvedIpAddressesByARecords := []string{randomIpAddress()}
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: email, ConnectionAttempts: connectionAttempts})

	t.Run("when host addresses was found during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return(resolvedHostNamesByPtrRecords, nil)
		resolver.On("mxRecords", resolvedHostNameByPtrRecord).Once().Return(priorities, resolvedHostNamesByMxRecords, nil)
		resolver.On("aRecords", resolvedHostNameByMxRecord).Once().Return(resolvedIpAddressesByARecords, nil)
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, resolvedIpAddressesByARecords, resolvedIpAddresses)
		assert.NoError(t, err)
	})

	t.Run("when CNAME record not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return("", createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when A record not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return("", createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when PTR records not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return([]string(nil), createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when MX records not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return(resolvedHostNamesByPtrRecords, nil)
		resolver.On("mxRecords", resolvedHostNameByPtrRecord).Once().Return([]uint16(nil), []string(nil), createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when null MX record found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return(resolvedHostNamesByPtrRecords, nil)
		resolver.On("mxRecords", resolvedHostNameByPtrRecord).Once().Return([]uint16{0}, []string{""}, nil)
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when A records by MX records not found", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return(resolvedHostNamesByPtrRecords, nil)
		resolver.On("mxRecords", resolvedHostNameByPtrRecord).Once().Return(priorities, resolvedHostNamesByMxRecords, nil)
		resolver.On("aRecords", resolvedHostNameByMxRecord).Once().Return([]string(nil), createDnsNotFoundError())
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})

	t.Run("when connection issues, cnameRecord error", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("cnameRecord", hostName).Times(connectionAttempts).Return("", new(validationError))
		resolvedIpAddresses, err := validation.hostsFromCnameRecord(hostName)
		resolver.AssertExpectations(t)
		assert.Equal(t, []string(nil), resolvedIpAddresses)
		assert.Error(t, err)
	})
}

func TestValidationMxRunMxLookup(t *testing.T) {
	email, notFoundDnsError := randomEmail(), createDnsNotFoundError()
	hostName, configuration := emailDomain(email), createConfiguration()

	t.Run("when mail servers found by MX records during first attempt", func(t *testing.T) {
		ipAddress := randomIpAddress()
		resolvedIpAddressesByARecords, hostNames, priorities := []string{ipAddress, ipAddress}, []string{randomDomain()}, []uint16{5}
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return(priorities, hostNames, nil)
		resolver.On("aRecords", hostNames[0]).Once().Return(resolvedIpAddressesByARecords, nil)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Equal(t, uniqStrings(resolvedIpAddressesByARecords), validatorResult.MailServers)
	})

	t.Run("when mail servers found by CNAME record during first attempt ", func(t *testing.T) {
		resolvedHostNameByCname := randomDomain()
		resolvedIpAddressByARecord := randomIpAddress()
		resolvedHostNameByPtrRecord := randomDomain()
		resolvedHostNamesByPtrRecords := []string{resolvedHostNameByPtrRecord}
		resolvedHostNameByMxRecord := randomDomain()
		priorities, resolvedHostNamesByMxRecords := []uint16{5}, []string{resolvedHostNameByMxRecord}
		ipAddress := randomIpAddress()
		resolvedIpAddressesByARecords := []string{ipAddress, ipAddress}
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{}, []string(nil), notFoundDnsError)
		resolver.On("cnameRecord", hostName).Once().Return(resolvedHostNameByCname, nil)
		resolver.On("aRecord", resolvedHostNameByCname).Once().Return(resolvedIpAddressByARecord, nil)
		resolver.On("ptrRecords", resolvedIpAddressByARecord).Once().Return(resolvedHostNamesByPtrRecords, nil)
		resolver.On("mxRecords", resolvedHostNameByPtrRecord).Once().Return(priorities, resolvedHostNamesByMxRecords, nil)
		resolver.On("aRecords", resolvedHostNameByMxRecord).Once().Return(resolvedIpAddressesByARecords, nil)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Equal(t, uniqStrings(resolvedIpAddressesByARecords), validatorResult.MailServers)
	})

	t.Run("when mail servers found by A record during first attempt", func(t *testing.T) {
		resolvedIpAddressByARecord := randomIpAddress()
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{}, []string(nil), notFoundDnsError)
		resolver.On("cnameRecord", hostName).Once().Return("", notFoundDnsError)
		resolver.On("aRecord", hostName).Once().Return(resolvedIpAddressByARecord, nil)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Equal(t, uniqStrings([]string{resolvedIpAddressByARecord}), validatorResult.MailServers)
	})

	t.Run("when mail servers not found by MX records, not RFC MX lookup flow enabled", func(t *testing.T) {
		otherConfiguration := copyConfigurationByPointer(configuration)
		otherConfiguration.NotRfcMxLookupFlow = true
		validatorResult, resolver := createSuccessfulValidatorResult(email, otherConfiguration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{}, []string(nil), notFoundDnsError)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Empty(t, validatorResult.MailServers)
	})

	t.Run("when mail servers not found by MX records, found null MX record", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{0}, []string{""}, nil)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Empty(t, validatorResult.MailServers)
	})

	t.Run("when mail servers not found by A record during first attempt", func(t *testing.T) {
		validatorResult, resolver := createSuccessfulValidatorResult(email, configuration), new(dnsResolverMock)
		validatorResult.punycodeDomain = hostName
		validation := &validationMx{result: validatorResult, resolver: resolver}

		resolver.On("mxRecords", hostName).Once().Return([]uint16{}, []string(nil), notFoundDnsError)
		resolver.On("cnameRecord", hostName).Once().Return("", notFoundDnsError)
		resolver.On("aRecord", hostName).Once().Return("", notFoundDnsError)
		validation.runMxLookup()
		resolver.AssertExpectations(t)
		assert.Empty(t, validatorResult.MailServers)
	})
}
