package emailver

import "github.com/L2SH-Dev/admissions/internal/mailing"

type EmailVerificationService interface {
	SendVerificationEmail(email string, registrationID uint) error
	VerifyEmail(token string) (uint, error)
}

type EmailVerificationServiceImpl struct {
	repo EmailVerificationRepo
}

func NewEmailVerificationService(repo EmailVerificationRepo) EmailVerificationService {
	return &EmailVerificationServiceImpl{repo: repo}
}

func (s *EmailVerificationServiceImpl) SendVerificationEmail(email string, registrationID uint) error {
	token, err := s.repo.CreateVerificationToken(registrationID)
	if err != nil {
		return err
	}

	err = mailing.SendEmailConfirmation(email, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailVerificationServiceImpl) VerifyEmail(token string) (uint, error) {
	registrationID, err := s.repo.GetRegistrationIDByToken(token)
	if err != nil {
		return 0, err
	}
	err = s.repo.DeleteToken(token)
	if err != nil {
		return 0, err
	}

	return registrationID, nil
}
