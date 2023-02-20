package truemail

import (
	"errors"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
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
	portNumber := server.PortNumber()
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

		smtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          targetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress),
			Response:      new(SmtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Once().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)
		validation.run()

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when successful session with third server during fisrt attempt", func(t *testing.T) {
		validatorResult := createValidatorResult(targetEmail, configuration)
		firstTargetHostAddress, secondTargetHostAddress, thirdTargetHostAddress := randomIpAddress(), randomIpAddress(), randomIpAddress()
		validatorResult.MailServers = append(validatorResult.MailServers, firstTargetHostAddress, secondTargetHostAddress, thirdTargetHostAddress)

		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(SmtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts := validation.attempts()

		firstSmtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          firstTargetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, firstTargetHostAddress),
			Response:      new(SmtpResponse),
		}
		secondSmtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          secondTargetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, secondTargetHostAddress),
			Response:      new(SmtpResponse),
		}
		thirdSmtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          thirdTargetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, thirdTargetHostAddress),
			Response:      new(SmtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, firstTargetHostAddress, configuration).Once().Return(firstSmtpReq)
		builder.On("newSmtpClient", firstSmtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, secondTargetHostAddress, configuration).Once().Return(secondSmtpReq)
		builder.On("newSmtpClient", secondSmtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, thirdTargetHostAddress, configuration).Once().Return(thirdSmtpReq)
		builder.On("newSmtpClient", thirdSmtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)
		validation.run()

		assert.Equal(t, []*SmtpRequest{firstSmtpReq, secondSmtpReq, thirdSmtpReq}, validation.smtpResults)
		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when failed session during first attempt for each server", func(t *testing.T) {
		validatorResult := createValidatorResult(targetEmail, configuration)
		firstTargetHostAddress, secondTargetHostAddress := randomIpAddress(), randomIpAddress()
		validatorResult.MailServers = append(validatorResult.MailServers, firstTargetHostAddress, secondTargetHostAddress)

		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(SmtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts := validation.attempts()

		firstSmtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          firstTargetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, firstTargetHostAddress),
			Response:      new(SmtpResponse),
		}
		secondSmtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          secondTargetHostAddress,
			Configuration: newSmtpRequestConfiguration(configuration, targetEmail, secondTargetHostAddress),
			Response:      new(SmtpResponse),
		}

		builder.On("newSmtpRequest", attempts, targetEmail, firstTargetHostAddress, configuration).Once().Return(firstSmtpReq)
		builder.On("newSmtpClient", firstSmtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)

		builder.On("newSmtpRequest", attempts, targetEmail, secondTargetHostAddress, configuration).Once().Return(secondSmtpReq)
		builder.On("newSmtpClient", secondSmtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)
		validation.run()

		assert.Equal(t, []*SmtpRequest{firstSmtpReq, secondSmtpReq}, validation.smtpResults)
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
		attempts, smtpResponse := validation.attempts(), new(SmtpResponse)

		smtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          targetHostAddress,
			Configuration: smtpRequestConfiguration,
			Response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Once().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.Configuration).Once().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(true)

		assert.True(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-1, smtpReq.Attempts)
		assert.True(t, smtpResponse.Rcptto)
		assert.Empty(t, smtpResponse.Errors)
		assert.Equal(t, []*SmtpRequest{smtpReq}, validation.smtpResults)
	})

	t.Run("when successful session during second attempt", func(t *testing.T) {
		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(SmtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts, smtpResponse := validation.attempts(), new(SmtpResponse)

		smtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          targetHostAddress,
			Configuration: smtpRequestConfiguration,
			Response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Twice().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.Configuration).Twice().Return(smtpClient)
		smtpClient.On("runSession").Once().Return(false)
		smtpClient.On("sessionError").Once().Return(sessionError)
		smtpClient.On("runSession").Once().Return(true)

		assert.True(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-2, smtpReq.Attempts)
		assert.True(t, smtpResponse.Rcptto)
		assert.Equal(t, []*SmtpClientError{sessionError}, smtpResponse.Errors)
		assert.Equal(t, []*SmtpRequest{smtpReq}, validation.smtpResults)
	})

	t.Run("when failed session during all attempts", func(t *testing.T) {
		builder, smtpClient, sessionError := new(smtpBuilderMock), new(smtpClientMock), new(SmtpClientError)
		validation := &validationSmtp{result: validatorResult, builder: builder}
		attempts, smtpResponse := validation.attempts(), new(SmtpResponse)

		smtpReq := &SmtpRequest{
			Attempts:      attempts,
			Email:         targetEmail,
			Host:          targetHostAddress,
			Configuration: smtpRequestConfiguration,
			Response:      smtpResponse,
		}

		builder.On("newSmtpRequest", attempts, targetEmail, targetHostAddress, configuration).Twice().Return(smtpReq)
		builder.On("newSmtpClient", smtpReq.Configuration).Twice().Return(smtpClient)
		smtpClient.On("runSession").Twice().Return(false)
		smtpClient.On("sessionError").Twice().Return(sessionError)

		assert.False(t, validation.runSmtpSession(targetHostAddress))
		assert.Equal(t, attempts-2, smtpReq.Attempts)
		assert.False(t, smtpResponse.Rcptto)
		assert.Equal(t, []*SmtpClientError{sessionError, sessionError}, smtpResponse.Errors)
		assert.Equal(t, []*SmtpRequest{smtpReq}, validation.smtpResults)
	})
}

func TestValidationSmtpIsFailFastScenario(t *testing.T) {
	t.Run("when SMTP fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpFailFast = true
		validation := &validationSmtp{result: &ValidatorResult{Configuration: configuration}}

		assert.True(t, validation.isFailFastScenario())
	})

	t.Run("when SMTP fail fast scenario is disabled", func(t *testing.T) {
		validation := &validationSmtp{result: &ValidatorResult{Configuration: createConfiguration()}}

		assert.False(t, validation.isFailFastScenario())
	})
}

func TestValidationSmtpFilteredMailServersByFailFastScenario(t *testing.T) {
	mailServers := []string{randomIpAddress(), randomIpAddress()}

	t.Run("when SMTP fail fast scenario is enabled", func(t *testing.T) {
		configuration := createConfiguration()
		configuration.SmtpFailFast = true
		validation := &validationSmtp{result: &ValidatorResult{MailServers: mailServers, Configuration: configuration}}

		assert.Equal(t, mailServers[:1], validation.filteredMailServersByFailFastScenario())
	})

	t.Run("when SMTP fail fast scenario is disabled", func(t *testing.T) {
		validation := &validationSmtp{result: &ValidatorResult{MailServers: mailServers, Configuration: createConfiguration()}}

		assert.Equal(t, mailServers, validation.filteredMailServersByFailFastScenario())
	})
}

func TestValidationSmtpIsMoreThanOneMailServer(t *testing.T) {
	t.Run("when more than one mail server", func(t *testing.T) {
		validation := &validationSmtp{
			result: &ValidatorResult{
				MailServers: []string{randomIpAddress(), randomIpAddress()},
			},
		}

		assert.True(t, validation.isMoreThanOneMailServer())
	})

	t.Run("when less than two mail servers", func(t *testing.T) {
		validation := &validationSmtp{result: &ValidatorResult{}}

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
			result: &ValidatorResult{
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
			result: &ValidatorResult{
				MailServers:   []string{randomIpAddress()},
				Configuration: configuration,
			},
		}

		assert.Equal(t, connectionAttempts, validation.attempts())
	})
}

func TestValidationSmtpIsIncludesSuccessfulSmtpResponse(t *testing.T) {
	failedSmtpRequest, successfulSmtpRequest := &SmtpRequest{Response: new(SmtpResponse)}, &SmtpRequest{Response: &SmtpResponse{Rcptto: true}}

	t.Run("when smtpResults is empty", func(t *testing.T) {
		validation := &validationSmtp{result: &ValidatorResult{}}

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when first successful SMTP response found", func(t *testing.T) {
		validation := &validationSmtp{
			result:      &ValidatorResult{},
			smtpResults: []*SmtpRequest{failedSmtpRequest, failedSmtpRequest, successfulSmtpRequest},
		}

		assert.True(t, validation.isIncludesSuccessfulSmtpResponse())
	})

	t.Run("when successful SMTP response not found", func(t *testing.T) {
		validation := &validationSmtp{
			result:      &ValidatorResult{},
			smtpResults: []*SmtpRequest{failedSmtpRequest},
		}

		assert.False(t, validation.isIncludesSuccessfulSmtpResponse())
	})
}

func TestValidationSmtpIsSmtpSafeCheckEnabled(t *testing.T) {
	t.Run("when SMTP safe check is disabled", func(t *testing.T) {
		validation := &validationSmtp{
			result: &ValidatorResult{Configuration: createConfiguration()},
		}

		assert.False(t, validation.isSmtpSafeCheckEnabled())
	})

	t.Run("when SMTP safe check is enabled", func(t *testing.T) {
		validation := &validationSmtp{
			result: &ValidatorResult{
				Configuration: &Configuration{SmtpSafeCheck: true},
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
			smtpResults: []*SmtpRequest{
				{
					Response: &SmtpResponse{
						Rcptto: true,
						Errors: []*SmtpClientError{
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
			smtpResults: []*SmtpRequest{
				{
					Response: &SmtpResponse{
						Rcptto: true,
						Errors: []*SmtpClientError{
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
