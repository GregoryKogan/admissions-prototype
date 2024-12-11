package exams

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/spf13/viper"
)

var (
	ErrAlreadyRegistered      = errors.New("user is already registered for this exam")
	ErrGradeMismatch          = errors.New("user grade does not match exam grade")
	ErrRegistrationNotAllowed = errors.New("registration is not allowed")
)

type ExamsService interface {
	Create(exam *Exam) error
	Delete(examID uint) error
	CreateDefaultExamTypes() error
	List() ([]*Exam, error)
	ListAvailable(*users.User) ([]*Exam, error)
	Register(user *users.User, examID uint) error
}

type ExamsServiceImpl struct {
	repo           ExamsRepo
	regDataService regdata.RegistrationDataService
}

func NewExamsService(repo ExamsRepo) ExamsService {
	return &ExamsServiceImpl{repo: repo}
}

func (s *ExamsServiceImpl) Create(exam *Exam) error {
	return s.repo.Create(exam)
}

func (s *ExamsServiceImpl) Delete(examID uint) error {
	return s.repo.Delete(examID)
}

func (s *ExamsServiceImpl) CreateDefaultExamTypes() error {
	examTypesConfig := viper.Get("exams.types").([]interface{})

	for _, examTypeData := range examTypesConfig {
		data := examTypeData.(map[string]interface{})
		examType := ExamType{
			Title:      data["title"].(string),
			Order:      data["order"].(int),
			Dismissing: data["dismissing"].(bool),
			HasPoints:  data["has_points"].(bool),
		}

		err := s.repo.CreateExamType(&examType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ExamsServiceImpl) List() ([]*Exam, error) {
	return s.repo.List()
}

func (s *ExamsServiceImpl) ListAvailable(user *users.User) ([]*Exam, error) {
	allExams, err := s.List()
	if err != nil {
		return nil, err
	}

	var availableExams []*Exam
	for _, exam := range allExams {
		allowed, err := s.isAllowedToRegister(user, exam)
		if err != nil {
			return nil, err
		}

		if allowed {
			availableExams = append(availableExams, exam)
		}
	}

	return availableExams, nil
}

func (s *ExamsServiceImpl) Register(user *users.User, examID uint) error {
	exam, err := s.repo.GetByID(examID)
	if err != nil {
		return err
	}

	allowed, err := s.isAllowedToRegister(user, exam)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrRegistrationNotAllowed
	}

	return s.repo.CreateRegistration(user.ID, examID)
}

func (s *ExamsServiceImpl) isAllowedToRegister(user *users.User, exam *Exam) (bool, error) {
	registered, err := s.repo.IsRegistered(user.ID, exam.ID)
	if err != nil {
		return false, err
	}
	if registered {
		return false, ErrAlreadyRegistered
	}

	regData, err := s.regDataService.GetByID(user.RegistrationDataID)
	if err != nil {
		return false, err
	}

	// TODO: Implement full allow logic
	// Should take order of the last passed exam into account
	if regData.Grade != exam.Grade {
		return false, ErrGradeMismatch
	}

	currentNumOfRegistrations, err := s.repo.CountRegistrations(exam.ID)
	if err != nil {
		return false, err
	}

	if currentNumOfRegistrations >= exam.Capacity {
		return false, nil
	}

	return true, nil
}
