package emailver

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type EmailVerificationRepo interface {
	CreateVerificationToken(registrationID uint) (string, error)
	GetRegistrationIDByToken(token string) (uint, error)
	DeleteToken(token string) error
}

type EmailVerificationRepoImpl struct {
	storage datastore.Storage
}

func NewEmailVerificationRepo(storage datastore.Storage) EmailVerificationRepo {
	return &EmailVerificationRepoImpl{storage: storage}
}

func (r *EmailVerificationRepoImpl) CreateVerificationToken(registrationID uint) (string, error) {
	token := uuid.New().String()
	if token == "" {
		return "", errors.New("failed to generate verification token")
	}
	err := r.storage.Cache().Set(
		context.Background(),
		fmt.Sprintf("email-token:%s", token),
		registrationID,
		viper.GetDuration("email_verification.token_lifetime"),
	).Err()

	if err != nil {
		return "", errors.New("failed to cache verification token")
	}

	return token, nil
}

func (r *EmailVerificationRepoImpl) GetRegistrationIDByToken(token string) (uint, error) {
	registrationIDString, err := r.storage.Cache().Get(
		context.Background(),
		fmt.Sprintf("email-token:%s", token),
	).Result()
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	registrationID, err := strconv.ParseUint(registrationIDString, 10, 32)
	if err != nil {
		return 0, errors.New("failed to parse registration ID")
	}

	return uint(registrationID), nil
}

func (r *EmailVerificationRepoImpl) DeleteToken(token string) error {
	err := r.storage.Cache().Del(
		context.Background(),
		fmt.Sprintf("email-token:%s", token),
	).Err()
	if err != nil {
		return errors.Join(errors.New("failed to delete token"), err)
	}

	return nil
}
