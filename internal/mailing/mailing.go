package mailing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type emailRequest struct {
	To      string      `json:"to"`
	Payment string      `json:"payment"`
	Params  interface{} `json:"params"`
}

type verificationParams struct {
	Email            string `json:"email"`
	VerificationLink string `json:"verification_link"`
}

type loginCredentialsParams struct {
	Email     string `json:"email"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	LoginLink string `json:"login_link"`
}

type rejectionParams struct {
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

func SendVerificationEmail(email string, token string) error {
	protocol := viper.GetString("server.protocol")
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	params := &verificationParams{
		Email:            email,
		VerificationLink: fmt.Sprintf("%s://%s:%s/verification?token=%s", protocol, host, port, token),
	}

	request := &emailRequest{
		To:      email,
		Payment: "credit",
		Params:  params,
	}

	return sendEmail("1350749", request)
}

func SendLoginAndPassword(email, login, password string) error {
	protocol := viper.GetString("server.protocol")
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	params := &loginCredentialsParams{
		Email:     email,
		Login:     login,
		Password:  password,
		LoginLink: fmt.Sprintf("%s://%s:%s/login", protocol, host, port),
	}

	request := &emailRequest{
		To:      email,
		Payment: "credit",
		Params:  params,
	}

	return sendEmail("1354145", request)
}

func SendRegistrationRejection(email, reason string) error {
	params := &rejectionParams{
		Email:  email,
		Reason: reason,
	}

	request := &emailRequest{
		To:      email,
		Payment: "credit",
		Params:  params,
	}

	return sendEmail("1365773", request)
}

func sendEmail(templateID string, request *emailRequest) error {
	if !viper.GetBool("mailing.enabled") {
		return nil
	}

	apiBase := viper.GetString("mailing.api_base")
	apiKey := viper.GetString("secrets.mail_api_key")

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/email/templates/%s/messages", apiBase, templateID), bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
