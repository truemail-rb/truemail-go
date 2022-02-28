package truemail

import (
	"fmt"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/stretchr/testify/assert"
)

func TestNewSmtpRequestConfiguration(t *testing.T) {
	t.Run("creates new smtp request configuration with settings specified in configuration", func(t *testing.T) {
		configuration, email, server := createConfiguration(), randomEmail(), randomIpAddress()
		configuration.SmtpPort, configuration.ConnectionTimeout, configuration.ResponseTimeout = randomPortNumber(), randomPositiveNumber(), randomPositiveNumber()
		smtpRequestConfiguration := newSmtpRequestConfiguration(configuration, email, server)

		assert.Equal(t, configuration.VerifierDomain, smtpRequestConfiguration.verifierDomain)
		assert.Equal(t, configuration.VerifierEmail, smtpRequestConfiguration.verifierEmail)
		assert.Equal(t, email, smtpRequestConfiguration.targetEmail)
		assert.Equal(t, server, smtpRequestConfiguration.targetServerAddress)
		assert.Equal(t, configuration.SmtpPort, smtpRequestConfiguration.targetServerPortNumber)
		assert.Equal(t, configuration.ConnectionTimeout, smtpRequestConfiguration.connectionTimeout)
		assert.Equal(t, configuration.ResponseTimeout, smtpRequestConfiguration.responseTimeout)
	})
}

func TestNewSmtpClient(t *testing.T) {
	t.Run("creates new smtp client with settings specified in smtpRequestConfiguration", func(t *testing.T) {
		smtpRequestConfig := &smtpRequestConfiguration{
			verifierDomain:         randomDomain(),
			verifierEmail:          randomEmail(),
			targetEmail:            randomEmail(),
			targetServerAddress:    randomIpAddress(),
			targetServerPortNumber: randomPortNumber(),
			connectionTimeout:      randomPositiveNumber(),
			responseTimeout:        randomPositiveNumber(),
		}
		smtpClient := newSmtpClient(smtpRequestConfig)

		assert.Equal(t, smtpRequestConfig.verifierDomain, smtpClient.verifierDomain)
		assert.Equal(t, smtpRequestConfig.verifierEmail, smtpClient.verifierEmail)
		assert.Equal(t, smtpRequestConfig.targetEmail, smtpClient.targetEmail)
		assert.Equal(t, smtpRequestConfig.targetServerAddress, smtpClient.targetServerAddress)
		assert.Equal(t, smtpRequestConfig.targetServerPortNumber, smtpClient.targetServerPortNumber)
		assert.Equal(t, tcpTransportLayer, smtpClient.networkProtocol)
		assert.Equal(t, time.Duration(smtpRequestConfig.connectionTimeout)*time.Second, smtpClient.connectionTimeout)
		assert.Equal(t, time.Duration(smtpRequestConfig.responseTimeout)*time.Second, smtpClient.responseTimeout)
	})
}

func TestSmtpInitConnection(t *testing.T) {
	t.Run("when connection succesful", func(t *testing.T) {
		server := startSmtpMock(smtpmock.ConfigurationAttr{})
		portNumber := server.PortNumber
		defer func() { _ = server.Stop() }()

		smtpClient := &smtpClient{
			networkProtocol:        tcpTransportLayer,
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumber,
		}
		targetServerWithPortNumber := serverWithPortNumber(smtpClient.targetServerAddress, smtpClient.targetServerPortNumber)
		connection, err := smtpClient.initConnection()

		assert.Equal(t, targetServerWithPortNumber, connection.RemoteAddr().String())
		assert.NoError(t, err)
	})

	t.Run("when connection failed", func(t *testing.T) {
		smtpClient := &smtpClient{
			networkProtocol:        tcpTransportLayer,
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: 1,
			connectionTimeout:      0,
		}
		targetServerWithPortNumber := serverWithPortNumber(smtpClient.targetServerAddress, smtpClient.targetServerPortNumber)
		errorMessage := fmt.Sprintf("dial tcp %s: connect: connection refused", targetServerWithPortNumber)
		connection, err := smtpClient.initConnection()

		assert.Nil(t, connection)
		assert.EqualError(t, err, errorMessage)
	})
}

func TestSmtpClientSessionError(t *testing.T) {
	t.Run("when error does not exist", func(t *testing.T) {
		assert.Nil(t, new(smtpClient).sessionError())
	})

	t.Run("when error exists", func(t *testing.T) {
		err := new(smtpClientError)
		smtpClient := &smtpClient{err: err}

		assert.Equal(t, err, smtpClient.sessionError())
	})
}

func TestSmtpClientRunSession(t *testing.T) {
	verifierDomain, verifierEmail, targetEmail := randomDomain(), randomEmail(), randomEmail()
	msgHeloBlacklistedDomain := "421 msgHeloBlacklistedDomain"
	msgMailfromBlacklistedEmail := "421 msgMailfromBlacklistedEmail"
	msgRcpttoNotRegisteredEmail := "550 MsgRcpttoNotRegisteredEmail"

	server := startSmtpMock(
		smtpmock.ConfigurationAttr{
			BlacklistedHeloDomains:      []string{verifierDomain},
			MsgHeloBlacklistedDomain:    msgHeloBlacklistedDomain,
			BlacklistedMailfromEmails:   []string{verifierEmail},
			MsgMailfromBlacklistedEmail: msgMailfromBlacklistedEmail,
			NotRegisteredEmails:         []string{targetEmail},
			MsgRcpttoNotRegisteredEmail: msgRcpttoNotRegisteredEmail,
		},
	)
	defer func() { _ = server.Stop() }()
	portNumber := server.PortNumber

	t.Run("iteracting with external SMTP server, no errors", func(t *testing.T) {
		client := &smtpClient{
			verifierDomain:         randomDomain(),
			verifierEmail:          randomEmail(),
			targetEmail:            randomEmail(),
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumber,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.True(t, client.runSession())
		assert.Nil(t, client.err)
	})

	t.Run("iteracting with external SMTP server, connection timeout", func(t *testing.T) {
		invalidPortNumber := 1
		errorMessage := fmt.Sprintf("dial tcp %s:%d: connect: connection refused", localhostIPv4Address, invalidPortNumber)
		client := &smtpClient{
			verifierDomain:         randomDomain(),
			verifierEmail:          randomEmail(),
			targetEmail:            randomEmail(),
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: invalidPortNumber,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.False(t, client.runSession())
		err := client.err
		assert.EqualError(t, err, errorMessage)
		assert.True(t, err.isConnection)
		assert.False(t, err.isResponseTimeout)
		assert.False(t, err.isHello)
		assert.False(t, err.isMailFrom)
		assert.False(t, err.isRecptTo)
	})

	t.Run("iteracting with external SMTP server, HELO error", func(t *testing.T) {
		client := &smtpClient{
			verifierDomain:         verifierDomain,
			verifierEmail:          randomEmail(),
			targetEmail:            randomEmail(),
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumber,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.False(t, client.runSession())
		err := client.err
		assert.EqualError(t, err, msgHeloBlacklistedDomain)
		assert.False(t, err.isConnection)
		assert.False(t, err.isResponseTimeout)
		assert.True(t, err.isHello)
		assert.False(t, err.isMailFrom)
		assert.False(t, err.isRecptTo)
	})

	t.Run("iteracting with external SMTP server, response timeout during HELO command", func(t *testing.T) {
		serverWithDelay := startSmtpMock(smtpmock.ConfigurationAttr{ResponseDelayHelo: 2})
		defer func() { _ = serverWithDelay.Stop() }()

		portNumberServerWithDelay := serverWithDelay.PortNumber
		errorMessage := fmt.Sprintf("->%s:%d: use of closed network connection", localhostIPv4Address, portNumberServerWithDelay)
		client := &smtpClient{
			verifierDomain:         randomDomain(),
			verifierEmail:          randomEmail(),
			targetEmail:            randomEmail(),
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumberServerWithDelay,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.False(t, client.runSession())
		err := client.err
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errorMessage)
		assert.False(t, err.isConnection)
		assert.False(t, err.isResponseTimeout)
		assert.True(t, err.isHello)
		assert.False(t, err.isMailFrom)
		assert.False(t, err.isRecptTo)
	})

	t.Run("iteracting with external SMTP server, MAIL FROM error", func(t *testing.T) {
		client := &smtpClient{
			verifierDomain:         randomDomain(),
			verifierEmail:          verifierEmail,
			targetEmail:            randomEmail(),
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumber,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.False(t, client.runSession())
		err := client.err
		assert.EqualError(t, err, msgMailfromBlacklistedEmail)
		assert.False(t, err.isConnection)
		assert.False(t, err.isResponseTimeout)
		assert.False(t, err.isHello)
		assert.True(t, err.isMailFrom)
		assert.False(t, err.isRecptTo)
	})

	t.Run("iteracting with external SMTP server, RCPT TO error", func(t *testing.T) {
		client := &smtpClient{
			verifierDomain:         randomDomain(),
			verifierEmail:          randomEmail(),
			targetEmail:            targetEmail,
			targetServerAddress:    localhostIPv4Address,
			targetServerPortNumber: portNumber,
			networkProtocol:        tcpTransportLayer,
			connectionTimeout:      time.Duration(1) * time.Second,
			responseTimeout:        time.Duration(1) * time.Second,
		}

		assert.False(t, client.runSession())
		err := client.err
		assert.EqualError(t, err, msgRcpttoNotRegisteredEmail)
		assert.False(t, err.isConnection)
		assert.False(t, err.isResponseTimeout)
		assert.False(t, err.isHello)
		assert.False(t, err.isMailFrom)
		assert.True(t, err.isRecptTo)
	})
}
