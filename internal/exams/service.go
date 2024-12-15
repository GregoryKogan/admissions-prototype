package exams

import (
	"errors"
	"log/slog"

	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/spf13/viper"
)

var (
	ErrRegistrationNotAllowed = errors.New("registration is not allowed")
)

type allocation struct {
	Capacity uint `json:"capacity"`
	Occupied uint `json:"occupied"`
}

type ExamsService interface {
	Create(exam *Exam) error
	Delete(examID uint) error
	CreateDefaultExamTypes() error
	List() ([]*Exam, error)
	Register(user *users.User, examID uint) error
	ListTypes() ([]*ExamType, error)
	Allocation(examID uint) (*allocation, error)
	History(user *users.User) ([]*Exam, error)
	Available(user *users.User) ([]*Exam, error)
	RegistrationStatus(user *users.User, examID uint) (bool, bool, error)
}

type ExamsServiceImpl struct {
	repo           ExamsRepo
	regDataService regdata.RegistrationDataService
}

func NewExamsService(repo ExamsRepo, regDataService regdata.RegistrationDataService) ExamsService {
	return &ExamsServiceImpl{repo: repo, regDataService: regDataService}
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

		title := data["title"].(string)
		exists, err := s.repo.TypeExistsByTitle(title)
		if err != nil {
			slog.Warn("Failed to check if exam type exists by title", slog.Any("title", title), slog.Any("err", err))
			continue
		}

		if exists {
			continue
		}

		examType := ExamType{
			Title:      data["title"].(string),
			Order:      data["order"].(int),
			Dismissing: data["dismissing"].(bool),
			HasPoints:  data["has_points"].(bool),
		}

		err = s.repo.CreateExamType(&examType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ExamsServiceImpl) List() ([]*Exam, error) {
	return s.repo.List()
}

func (s *ExamsServiceImpl) Register(user *users.User, examID uint) error {
	exam, err := s.repo.GetByID(examID)
	if err != nil {
		return err
	}

	// Check if the user is allowed to register
	if err := s.canRegister(user, exam); err != nil {
		return err
	}

	return s.repo.CreateRegistration(user.ID, exam.ID)
}

func (s *ExamsServiceImpl) canRegister(user *users.User, exam *Exam) error {
	// Check if the user is already registered to the exam
	registered, err := s.repo.IsRegistered(user.ID, exam.ID)
	if err != nil {
		return err
	}
	if registered {
		return ErrAlreadyRegistered
	}

	// Retrieve the user's registration data
	regData, err := s.regDataService.GetByID(user.RegistrationDataID)
	if err != nil {
		return err
	}

	// Check if the user's grade matches the exam's grade
	if regData.Grade != exam.Grade {
		return ErrInvalidGrade
	}

	// Check if the exam capacity has not been exceeded
	currentRegistrations, err := s.repo.CountRegistrations(exam.ID)
	if err != nil {
		return err
	}
	if currentRegistrations >= exam.Capacity {
		return ErrExamFull
	}

	// Get the next required exam type order for the user
	nextOrder, err := s.repo.GetNextExamTypeOrder(user.ID)
	if err != nil {
		return err
	}

	// Check if the exam's type order matches the next required order
	if exam.ExamType.Order != nextOrder {
		return ErrInvalidExamOrder
	}

	return nil
}

func (s *ExamsServiceImpl) ListTypes() ([]*ExamType, error) {
	return s.repo.ListTypes()
}

func (s *ExamsServiceImpl) Allocation(examID uint) (*allocation, error) {
	exam, err := s.repo.GetByID(examID)
	if err != nil {
		return nil, err
	}

	occupied, err := s.repo.CountRegistrations(examID)
	if err != nil {
		return nil, err
	}

	return &allocation{Capacity: exam.Capacity, Occupied: occupied}, nil
}

func (s *ExamsServiceImpl) History(user *users.User) ([]*Exam, error) {
	return s.repo.History(user.ID)
}

func (s *ExamsServiceImpl) Available(user *users.User) ([]*Exam, error) {
	regData, err := s.regDataService.GetByID(user.RegistrationDataID)
	if err != nil {
		return nil, err
	}

	return s.repo.Available(user.ID, regData.Grade)
}

func (s *ExamsServiceImpl) RegistrationStatus(user *users.User, examID uint) (bool, bool, error) {
	registeredToExam, registeredToSameType, err := s.repo.RegistrationStatus(user.ID, examID)
	if err != nil {
		return false, false, err
	}
	return registeredToExam, registeredToSameType, nil
}
