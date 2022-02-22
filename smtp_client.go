package truemail

import (
	"net"
	"net/smtp"
	"time"
)

// SMTP request configuration. Provides connection/request settings for SMTP client
type smtpRequestConfiguration struct {
	verifierDomain, verifierEmail, targetEmail, targetServerAddress string
	targetServerPortNumber, connectionTimeout, responseTimeout      int
}

// smtpRequestConfiguration builder. Creates SMTP request configuration with settings from configuration
func newSmtpRequestConfiguration(config *configuration, targetEmail, targetServerAddress string) *smtpRequestConfiguration {
	return &smtpRequestConfiguration{
		verifierDomain:         config.VerifierDomain,
		verifierEmail:          config.VerifierEmail,
		targetEmail:            targetEmail,
		targetServerAddress:    targetServerAddress,
		targetServerPortNumber: config.SmtpPort,
		connectionTimeout:      config.ConnectionTimeout,
		responseTimeout:        config.ResponseTimeout,
	}
}

// SMTP response structure. Includes RCPTTO successful request marker
// and SMTP client error pointers slice
type smtpResponse struct {
	rcptto bool
	errors []*smtpClientError
}

// SMTP request structure. Includes attempts count, target email & host address,
// pointers to SMTP request configuration and SMTP response
type smtpRequest struct {
	attempts      int
	email, host   string
	configuration *smtpRequestConfiguration
	response      *smtpResponse
}

// SMTP validation client interface
type client interface {
	runSession() bool
	sessionError() *smtpClientError
}

// SMTP client structure. Provides possibility to interact with target SMTP server
type smtpClient struct {
	verifierDomain, verifierEmail, targetEmail, targetServerAddress, networkProtocol string
	targetServerPortNumber                                                           int
	connectionTimeout, responseTimeout                                               time.Duration
	client                                                                           *smtp.Client
	err                                                                              *smtpClientError
}

// smtpClient builder. Creates SMTP client with settings from smtpRequestConfiguration
func newSmtpClient(config *smtpRequestConfiguration) *smtpClient {
	return &smtpClient{
		verifierDomain:         config.verifierDomain,
		verifierEmail:          config.verifierEmail,
		targetEmail:            config.targetEmail,
		targetServerAddress:    config.targetServerAddress,
		targetServerPortNumber: config.targetServerPortNumber,
		networkProtocol:        tcpTransportLayer,
		connectionTimeout:      time.Duration(config.connectionTimeout) * time.Second,
		responseTimeout:        time.Duration(config.responseTimeout) * time.Second,
	}
}

// smtpClient methods

// Initializes SMTP client connection with connection timeout
func (smtpClient *smtpClient) initConnection() (net.Conn, error) {
	targetAddress := serverWithPortNumber(smtpClient.targetServerAddress, smtpClient.targetServerPortNumber)
	connection, error := net.DialTimeout(smtpClient.networkProtocol, targetAddress, smtpClient.connectionTimeout)
	return connection, error
}

// interface implementation

// Returns pointer to current SMTP client custom error
func (smtpClient *smtpClient) sessionError() *smtpClientError {
	return smtpClient.err
}

// Runs SMTP session with target mail server. Assigns smtpClient.error
// for failure case and return false. Otherwise returns true
func (smtpClient *smtpClient) runSession() bool {
	var err error
	connection, err := smtpClient.initConnection()

	if err != nil {
		smtpClient.err = &smtpClientError{isConnection: true, err: err}
		return false
	}

	closeConnection := func() {
		connection.Close()
		smtpClient.err = &smtpClientError{isResponseTimeout: true, err: err}
	}

	client, _ := smtp.NewClient(connection, smtpClient.targetServerAddress)
	smtpClient.client = client
	defer client.Close()

	timerHello := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Hello(smtpClient.verifierDomain)
	if err != nil {
		smtpClient.err = &smtpClientError{isHello: true, err: err}
		return false
	}
	defer timerHello.Stop()

	timerMailFrom := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Mail(smtpClient.verifierEmail)
	if err != nil {
		smtpClient.err = &smtpClientError{isMailFrom: true, err: err}
		return false
	}
	defer timerMailFrom.Stop()

	timerRcptTo := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Rcpt(smtpClient.targetEmail)
	if err != nil {
		smtpClient.err = &smtpClientError{isRecptTo: true, err: err}
		return false
	}
	defer timerRcptTo.Stop()

	// TODO: What about client.Quit() ?
	return true
}
