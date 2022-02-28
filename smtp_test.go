package truemail

import (
	"errors"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/stretchr/testify/assert"
)

func TestValidationSmtpCheck(t *testing.T) {
	blacklistedMailfromEmail, nonExistentEmail := randomEmail(), randomEmail()
	server := startSmtpMock(
		smtpmock.ConfigurationAttr{
			BlacklistedMailfromEmails: []string{blacklistedMailfromEmail},
			NotRegisteredEmails:       []string{nonExistentEmail},
		},
	)
	portNumber := server.PortNumber
	defer func() { _ = server.Stop() }()

	t.Run("SMTP validation: successful after first attempt on first server", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpPort = portNumber
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		new(validationSmtp).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.SmtpDebug)
		assert.Empty(t, validatorResult.usedValidations)
	})

	// // TODO: add for successful case during second attempt, MailServers == 1; MailServers > 1;

	t.Run("SMTP validation: failed after second attempt on second server, safe check scenario is disabled, fail fast scenario is disabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.VerifierEmail, configuration.SmtpPort = blacklistedMailfromEmail, portNumber
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"smtp": "smtp error"}, validatorResult.Errors)
		assert.Equal(t, 2, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("SMTP validation: failed after first attempt on first server, safe check scenario is disabled, fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.VerifierEmail, configuration.SmtpPort, configuration.SmtpFailFast = blacklistedMailfromEmail, portNumber, true
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"smtp": "smtp error"}, validatorResult.Errors)
		assert.Equal(t, 1, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("SMTP validation: successful after first attempt on first server, safe check scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.VerifierEmail, configuration.SmtpPort, configuration.SmtpSafeCheck = blacklistedMailfromEmail, portNumber, true
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Equal(t, 2, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("SMTP validation: successful after first attempt on first server, safe check scenario is enabled, fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.VerifierEmail, configuration.SmtpPort, configuration.SmtpSafeCheck, configuration.SmtpFailFast = blacklistedMailfromEmail, portNumber, true, true
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Equal(t, 1, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("SMTP validation: failure after second attempt on second server, safe check scenario is enabled, fail fast is disabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpPort, configuration.SmtpSafeCheck = portNumber, true
		validatorResult := createSuccessfulValidatorResult(nonExistentEmail, configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"smtp": "smtp error"}, validatorResult.Errors)
		assert.Equal(t, 2, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("SMTP validation: failure after first attempt on first server, safe check scenario is enabled, fail fast is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpPort, configuration.SmtpSafeCheck, configuration.SmtpFailFast = portNumber, true, true
		validatorResult := createSuccessfulValidatorResult(nonExistentEmail, configuration)
		validatorResult.MailServers = append(validatorResult.MailServers, localhostIPv4Address, localhostIPv4Address)
		validationSmtp := new(validationSmtp)
		validationSmtp.check(validatorResult)
		smtpDebug := validatorResult.SmtpDebug

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{"smtp": "smtp error"}, validatorResult.Errors)
		assert.Equal(t, 1, len(smtpDebug))
		assert.Equal(t, smtpDebug, validationSmtp.smtpResults)
		assert.Empty(t, validatorResult.usedValidations)
	})
}

func TestValidationSmtpInitSmtpBuilder(t *testing.T) {
	t.Run("creates SMTP validation SMTP entities builder", func(t *testing.T) {
		validation := new(validationSmtp)
		validation.initSmtpBuilder()

		assert.Equal(t, new(smtpBuilder), validation.builder)
	})
}

func TestValidationSmtpRun(t *testing.T) {
	targetEmail, configuration := randomEmail(), createConfiguration()

	t.Run("when successful session with first server during first attempt", func(t *testing.T) {
		validatorResult := createValidatorResult(targetEmail, configuration)
		targetHostAddress := randomIpAddress()
		validatorResult.MailServers = append(validatorResult.MailServers, targetHostAddress)

		builder, smtpClient := new(smtpBuilderMock), new(smtpClientMock)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts := validation.attempts()

		smtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          targetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress),
			response:      new(smtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Once().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)
		validation.run()

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when successful session with third server during fisrt attempt", func(t *testing.T) {
		validatorResult := createValidatorResult(targetEmail, configuration)
		firstTargetHostAddress, secondTargetHostAddress, thirdTargetHostAddress := randomIpAddress(), randomIpAddress(), randomIpAddress()
		validatorResult.MailServers = append(validatorResult.MailServers, firstTargetHostAddress, secondTargetHostAddress, thirdTargetHostAddress)

		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(smtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts := validation.attempts()

		firstSmtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          firstTargetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, firstTargetHostAddress),
			response:      new(smtpResponse),
		}
		secondSmtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          secondTargetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, secondTargetHostAddress),
			response:      new(smtpResponse),
		}
		thirdSmtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          thirdTargetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, thirdTargetHostAddress),
			response:      new(smtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, firstTargetHostAddress, configuration).Once().Return(firstSmtpReq)
		builder.On("newSmtpClient", firstSmtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, secondTargetHostAddress, configuration).Once().Return(secondSmtpReq)
		builder.On("newSmtpClient", secondSmtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, thirdTargetHostAddress, configuration).Once().Return(thirdSmtpReq)
		builder.On("newSmtpClient", thirdSmtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)
		validation.run()

		assert.Equal(t, []*smtpRequest{firstSmtpReq, secondSmtpReq, thirdSmtpReq}, validation.smtpResults)
		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when failed session during first attempt for each server", func(t *testing.T) {
		validatorResult := createValidatorResult(targetEmail, configuration)
		firstTargetHostAddress, secondTargetHostAddress := randomIpAddress(), randomIpAddress()
		validatorResult.MailServers = append(validatorResult.MailServers, firstTargetHostAddress, secondTargetHostAddress)

		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(smtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts := validation.attempts()

		firstSmtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          firstTargetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, firstTargetHostAddress),
			response:      new(smtpResponse),
		}
		secondSmtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          secondTargetHostAddress,
			configuration: newSmtpRequestConfiguration(configuration, targetEmail, secondTargetHostAddress),
			response:      new(smtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, firstTargetHostAddress, configuration).Once().Return(firstSmtpReq)
		builder.On("newSmtpClient", firstSmtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, secondTargetHostAddress, configuration).Once().Return(secondSmtpReq)
		builder.On("newSmtpClient", secondSmtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)
		validation.run()

		assert.Equal(t, []*smtpRequest{firstSmtpReq, secondSmtpReq}, validation.smtpResults)
		assert.False(t, validation.isIncludesSuccessfulSmtpResponse())
	})
}

func TestValidationSmtpRunSmtpSession(t *testing.T) {
	targetEmail, targetHostAddress, configuration := randomEmail(), randomIpAddress(), createConfiguration()
	validatorResult := createValidatorResult(targetEmail, configuration)
	smtpRequestConfiguration := newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress)

	t.Run("when successful session during first attempt", func(t *testing.T) {
		builder, smtpClient := new(smtpBuilderMock), new(smtpClientMock)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts, smtpResponse := validation.attempts(), new(smtpResponse)

		smtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          targetHostAddress,
			configuration: smtpRequestConfiguration,
			response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Once().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)

		assert.True(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-1, smtpReq.attempts)
		assert.True(t, smtpResponse.rcptto)
		assert.Empty(t, smtpResponse.errors)
		assert.Equal(t, []*smtpRequest{smtpReq}, validation.smtpResults)
	})

	t.Run("when successful session during second attempt", func(t *testing.T) {
		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(smtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts, smtpResponse := validation.attempts(), new(smtpResponse)

		smtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          targetHostAddress,
			configuration: smtpRequestConfiguration,
			response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Twice().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.configuration).Twice().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)
		smtpClient.On("runSession").Once().Return(true)

		assert.True(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-2, smtpReq.attempts)
		assert.True(t, smtpResponse.rcptto)
		assert.Equal(t, []*smtpClientError{sessionError}, smtpResponse.errors)
		assert.Equal(t, []*smtpRequest{smtpReq}, validation.smtpResults)
	})

	t.Run("when failed session during all attempts", func(t *testing.T) {
		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(smtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts, smtpResponse := validation.attempts(), new(smtpResponse)

		smtpReq := &smtpRequest{
			attempts:      attempts,
			email:         targetEmail,
			host:          targetHostAddress,
			configuration: smtpRequestConfiguration,
			response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Twice().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.configuration).Twice().Return(smtpClient)
		smtpClient.On("runSession").Twice().Return(false)
		smtpClient.On("sessionError").Twice().Return(sessionError)

		assert.False(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-2, smtpReq.attempts)
		assert.False(t, smtpResponse.rcptto)
		assert.Equal(t, []*smtpClientError{sessionError, sessionError}, smtpResponse.errors)
		assert.Equal(t, []*smtpRequest{smtpReq}, validation.smtpResults)
	})
}

func TestValidationSmtpIsFailFastScenario(t *testing.T) {
	t.Run("when SMTP fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpFailFast = true
		validation := &validationSmtp{result: &validatorResult{Configuration: configuration}}

		assert.True(t, validation.isFailFastScenario())
	})

	t.Run("when SMTP fail fast scenario is disabled", func(t *testing.T) {
		validation := &validationSmtp{result: &validatorResult{Configuration: createConfiguration()}}

		assert.False(t, validation.isFailFastScenario())
	})
}

func TestValidationSmtpFilteredMailServersByFailFastScenario(t *testing.T) {
	mailServers := []string{randomIpAddress(), randomIpAddress()}

	t.Run("when SMTP fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpFailFast = true
		validation := &validationSmtp{result: &validatorResult{MailServers: mailServers, Configuration: configuration}}

		assert.Equal(t, mailServers[:1], validation.filteredMailServersByFailFastScenario())
	})

	t.Run("when SMTP fail fast scenario is disabled", func(t *testing.T) {
		validation := &validationSmtp{result: &validatorResult{MailServers: mailServers, Configuration: createConfiguration()}}

		assert.Equal(t, mailServers, validation.filteredMailServersByFailFastScenario())
	})
}

func TestValidationSmtpIsMoreThanOneMailServer(t *testing.T) {
	t.Run("when more than one mail server", func(t *testing.T) {
		validation := &validationSmtp{
			result: &validatorResult{
				MailServers: []string{randomIpAddress(), randomIpAddress()},
			},
		}

		assert.True(t, validation.isMoreThanOneMailServer())
	})

	t.Run("when less than two mail servers", func(t *testing.T) {
		validation := &validationSmtp{result: &validatorResult{}}

		assert.False(t, validation.isMoreThanOneMailServer())
	})
}

func TestValidationSmtpAttempts(t *testing.T) {
	t.Run("when SMTP fail fast scenario enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpFailFast = true
		validation := &validationSmtp{result: createValidatorResult(randomEmail(), configuration)}

		assert.Equal(t, 1, validation.attempts())
	})

	t.Run("when more than one mail server", func(t *testing.T) {
		validation := &validationSmtp{
			result: &validatorResult{
				MailServers:   []string{randomIpAddress(), randomIpAddress()},
				Configuration: createConfiguration(),
			},
		}

		assert.Equal(t, 1, validation.attempts())
	})

	t.Run("when less than two mail servers", func(t *testing.T) {
		connectionAttempts := 42
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				ConnectionAttempts: connectionAttempts,
			},
		)
		validation := &validationSmtp{
			result: &validatorResult{
				MailServers:   []string{randomIpAddress()},
				Configuration: configuration,
			},
		}

		assert.Equal(t, connectionAttempts, validation.attempts())
	})
}

func TestValidationSmtpIsIncludesSuccessfulSmtpResponse(t *testing.T) {
	failedSmtpRequest, successfulSmtpRequest := &smtpRequest{response: new(smtpResponse)}, &smtpRequest{response: &smtpResponse{rcptto: true}}

	t.Run("when smtpResults is empty", func(t *testing.T) {
		validation := &validationSmtp{result: &validatorResult{}}

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when first successful SMTP response found", func(t *testing.T) {
		validation := &validationSmtp{
			result:      &validatorResult{},
			smtpResults: []*smtpRequest{failedSmtpRequest, failedSmtpRequest, successfulSmtpRequest},
		}

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when successful SMTP response not found", func(t *testing.T) {
		validation := &validationSmtp{
			result:      &validatorResult{},
			smtpResults: []*smtpRequest{failedSmtpRequest},
		}

		assert.False(t, validation.isIncludesSuccessfulSmtpResponse())
	})
}

func TestValidationSmtpIsSmtpSafeCheckEnabled(t *testing.T) {
	t.Run("when SMTP safe check is disabled", func(t *testing.T) {
		validation := &validationSmtp{
			result: &validatorResult{Configuration: createConfiguration()},
		}

		assert.False(t, validation.isSmtpSafeCheckEnabled())
	})

	t.Run("when SMTP safe check is enabled", func(t *testing.T) {
		validation := &validationSmtp{
			result: &validatorResult{
				Configuration: &configuration{SmtpSafeCheck: true},
			},
		}

		assert.True(t, validation.isSmtpSafeCheckEnabled())
	})
}

func TestValidationSmtpIsNotIncludeUserNotFoundErrors(t *testing.T) {
	t.Run("when smtpResults is empty", func(t *testing.T) {
		assert.True(t, new(validationSmtp).isNotIncludeUserNotFoundErrors())
	})

	t.Run("when does not contain recognized UserNotFound errors ", func(t *testing.T) {
		validation := &validationSmtp{
			result: createValidatorResult(randomEmail(), createConfiguration()),
			smtpResults: []*smtpRequest{
				{
					response: &smtpResponse{
						rcptto: true,
						errors: []*smtpClientError{
							{
								isConnection: true,
								err:          errors.New("Some connection error"),
							},
							{
								isHello: true,
								err:     errors.New("Some HELO error"),
							},
							{
								isRecptTo: true,
								err:       errors.New("Some RCPT TO error"),
							},
						},
					},
				},
			},
		}

		assert.True(t, validation.isNotIncludeUserNotFoundErrors())
	})

	t.Run("when contains recognized UserNotFound errors ", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:        randomEmail(),
				SmtpErrorBodyPattern: `RCPTTO ERROR`,
			},
		)
		validation := &validationSmtp{
			result: createValidatorResult(randomEmail(), configuration),
			smtpResults: []*smtpRequest{
				{
					response: &smtpResponse{
						rcptto: true,
						errors: []*smtpClientError{
							{
								isConnection: true,
								err:          errors.New("Some connection error"),
							},
							{
								isHello: true,
								err:     errors.New("Some HELO error"),
							},
							{
								isRecptTo: true,
								err:       errors.New("Some RCPTTO ERROR"),
							},
						},
					},
				},
			},
		}

		assert.False(t, validation.isNotIncludeUserNotFoundErrors())
	})
}
