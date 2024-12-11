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
	ListAvailable(*users.User) ([]*Exam, error)
	ListUserExams(user *users.User) ([]*Exam, error)
	Register(user *users.User, examID uint) error
	ListTypes() ([]*ExamType, error)
	Allocation(examID uint) (*allocation, error)
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

func (s *ExamsServiceImpl) ListUserExams(user *users.User) ([]*Exam, error) {
	return s.repo.ListUserExams(user.ID)
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

func (s *ExamsServiceImpl) isAllowedToRegister(user *users.User, exam *Exam) (bool, error) {
	registered, err := s.repo.IsRegistered(user.ID, exam.ID)
	if err != nil {
		return false, err
	}
	if registered {
		return false, nil
	}

	regData, err := s.regDataService.GetByID(user.RegistrationDataID)
	if err != nil {
		return false, err
	}

	// TODO: Implement full allow logic
	// Should take order of the last passed exam into account
	if regData.Grade != exam.Grade {
		return false, nil
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
