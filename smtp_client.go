package truemail

import (
	"net"
	"net/smtp"
	"time"
)

// SMTP request configuration. Provides connection/request settings for SMTP client
type SmtpRequestConfiguration struct {
	VerifierDomain, VerifierEmail, TargetEmail, TargetServerAddress string
	TargetServerPortNumber, ConnectionTimeout, ResponseTimeout      int
}

// smtpRequestConfiguration builder. Creates SMTP request configuration with settings from configuration
func newSmtpRequestConfiguration(config *Configuration, targetEmail, targetServerAddress string) *SmtpRequestConfiguration {
	return &SmtpRequestConfiguration{
		VerifierDomain:         config.VerifierDomain,
		VerifierEmail:          config.VerifierEmail,
		TargetEmail:            targetEmail,
		TargetServerAddress:    targetServerAddress,
		TargetServerPortNumber: config.SmtpPort,
		ConnectionTimeout:      config.ConnectionTimeout,
		ResponseTimeout:        config.ResponseTimeout,
	}
}

// SMTP response structure. Includes RCPTTO successful request marker
// and SMTP client error pointers slice
type SmtpResponse struct {
	Rcptto bool
	Errors []*SmtpClientError
}

// SMTP request structure. Includes attempts count, target email & host address,
// pointers to SMTP request configuration and SMTP response
type SmtpRequest struct {
	Attempts      int
	Email, Host   string
	Configuration *SmtpRequestConfiguration
	Response      *SmtpResponse
}

// SMTP validation client interface
type client interface {
	runSession() bool
	sessionError() *SmtpClientError
}

// SMTP client structure. Provides possibility to interact with target SMTP server
type smtpClient struct {
	verifierDomain, verifierEmail, targetEmail, targetServerAddress, networkProtocol string
	targetServerPortNumber                                                           int
	connectionTimeout, responseTimeout                                               time.Duration
	client                                                                           *smtp.Client
	err                                                                              *SmtpClientError
}

// smtpClient builder. Creates SMTP client with settings from smtpRequestConfiguration
func newSmtpClient(config *SmtpRequestConfiguration) *smtpClient {
	return &smtpClient{
		verifierDomain:         config.VerifierDomain,
		verifierEmail:          config.VerifierEmail,
		targetEmail:            config.TargetEmail,
		targetServerAddress:    config.TargetServerAddress,
		targetServerPortNumber: config.TargetServerPortNumber,
		networkProtocol:        tcpTransportLayer,
		connectionTimeout:      time.Duration(config.ConnectionTimeout) * time.Second,
		responseTimeout:        time.Duration(config.ResponseTimeout) * time.Second,
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
func (smtpClient *smtpClient) sessionError() *SmtpClientError {
	return smtpClient.err
}

// Runs SMTP session with target mail server. Assigns smtpClient.error
// for failure case and return false. Otherwise returns true
func (smtpClient *smtpClient) runSession() bool {
	var err error
	connection, err := smtpClient.initConnection()

	if err != nil {
		smtpClient.err = &SmtpClientError{isConnection: true, err: err}
		return false
	}

	closeConnection := func() {
		connection.Close()
		smtpClient.err = &SmtpClientError{isResponseTimeout: true, err: err}
	}

	client, _ := smtp.NewClient(connection, smtpClient.targetServerAddress)
	smtpClient.client = client
	defer client.Close()

	timerHello := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Hello(smtpClient.verifierDomain)
	if err != nil {
		smtpClient.err = &SmtpClientError{isHello: true, err: err}
		return false
	}
	defer timerHello.Stop()

	timerMailFrom := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Mail(smtpClient.verifierEmail)
	if err != nil {
		smtpClient.err = &SmtpClientError{isMailFrom: true, err: err}
		return false
	}
	defer timerMailFrom.Stop()

	timerRcptTo := time.AfterFunc(smtpClient.responseTimeout, closeConnection)
	err = client.Rcpt(smtpClient.targetEmail)
	if err != nil {
		smtpClient.err = &SmtpClientError{isRecptTo: true, err: err}
		return false
	}
	defer timerRcptTo.Stop()

	// TODO: What about client.Quit() ?
	return true
}
