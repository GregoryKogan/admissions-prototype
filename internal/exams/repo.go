package exams

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"gorm.io/gorm"
)

type ExamsRepo interface {
	Create(exam *Exam) error
	Delete(examID uint) error
	CreateExamType(examType *ExamType) error
	List() ([]*Exam, error)
	GetByID(examID uint) (*Exam, error)
	CreateRegistration(userID, examID uint) error
	IsRegistered(userID, examID uint) (bool, error)
	CountRegistrations(examID uint) (uint, error)
	ListTypes() ([]*ExamType, error)
	TypeExistsByTitle(title string) (bool, error)
	History(userID uint) ([]*Exam, error)
	Available(userID uint, grade uint) ([]*Exam, error)
}

type ExamsRepoImpl struct {
	storage datastore.Storage
}

func NewExamsRepo(storage datastore.Storage) ExamsRepo {
	if err := storage.DB().AutoMigrate(&Exam{}, &ExamType{}, &ExamRegistration{}, &ExamResult{}); err != nil {
		panic(err)
	}
	return &ExamsRepoImpl{storage: storage}
}

func (r *ExamsRepoImpl) Create(exam *Exam) error {
	err := r.storage.DB().Create(exam).Error
	if err != nil {
		return err
	}

	if exam.ID == 0 {
		return errors.New("exam creation failed: exam ID is not set")
	}

	return nil
}

func (r *ExamsRepoImpl) Delete(examID uint) error {
	err := r.storage.DB().Delete(&Exam{}, examID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ExamsRepoImpl) CreateExamType(examType *ExamType) error {
	err := r.storage.DB().Create(examType).Error
	if err != nil {
		return err
	}

	if examType.ID == 0 {
		return errors.New("exam type creation failed: exam type ID is not set")
	}

	return nil
}

func (r *ExamsRepoImpl) List() ([]*Exam, error) {
	var exams []*Exam
	err := r.storage.DB().Preload("ExamType").Order("start DESC").Find(&exams).Error
	if err != nil {
		return nil, err
	}

	return exams, nil
}

func (r *ExamsRepoImpl) GetByID(examID uint) (*Exam, error) {
	var exam Exam
	err := r.storage.DB().Preload("ExamType").First(&exam, examID).Error
	if err != nil {
		return nil, err
	}

	return &exam, nil
}

func (r *ExamsRepoImpl) CreateRegistration(userID, examID uint) error {
	reg := &ExamRegistration{UserID: userID, ExamID: examID}
	err := r.storage.DB().Create(reg).Error
	if err != nil {
		return err
	}

	if reg.ID == 0 {
		return errors.New("exam registration creation failed: registration ID is not set")
	}

	return nil
}

func (r *ExamsRepoImpl) IsRegistered(userID, examID uint) (bool, error) {
	var count int64
	err := r.storage.DB().Model(&ExamRegistration{}).Where("user_id = ? AND exam_id = ?", userID, examID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ExamsRepoImpl) CountRegistrations(examID uint) (uint, error) {
	var count int64
	err := r.storage.DB().Model(&ExamRegistration{}).Where("exam_id = ?", examID).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (r *ExamsRepoImpl) ListTypes() ([]*ExamType, error) {
	var types []*ExamType
	err := r.storage.DB().Find(&types).Error
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (r *ExamsRepoImpl) TypeExistsByTitle(title string) (bool, error) {
	var count int64
	err := r.storage.DB().Model(&ExamType{}).Where("title = ?", title).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ExamsRepoImpl) History(userID uint) ([]*Exam, error) {
	var exams []*Exam
	err := r.storage.DB().
		Preload("ExamType").
		Joins("JOIN exam_registrations ON exams.id = exam_registrations.exam_id").
		Where("exam_registrations.user_id = ? AND exams.end < NOW()", userID).
		Find(&exams).Error
	if err != nil {
		return nil, err
	}

	return exams, nil
}

func (r *ExamsRepoImpl) Available(userID uint, grade uint) ([]*Exam, error) {
	// Get the last passed exam's exam type order for the user
	var lastPassedExamResult ExamResult
	err := r.storage.DB().
		Select("exam_results.*").
		Joins("JOIN exams ON exam_results.exam_id = exams.id").
		Joins("JOIN exam_types ON exams.exam_type_id = exam_types.id").
		Where("exam_results.user_id = ? AND exam_results.result = ?", userID, "PASSED").
		Order("exam_types.order DESC").
		First(&lastPassedExamResult).Error

	var lastOrder int
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No passed exams, set lastOrder to zero
			lastOrder = 0
		} else {
			return nil, err
		}
	} else {
		lastOrder = lastPassedExamResult.Exam.ExamType.Order
	}

	// Get the next exam type order after the last passed exam
	var nextOrder int
	err = r.storage.DB().
		Model(&ExamType{}).
		Select("MIN(`order`)").
		Where("`order` > ?", lastOrder).
		Scan(&nextOrder).Error
	if err != nil {
		return nil, err
	}

	// If no next exam type is found, return an empty list
	if nextOrder == 0 {
		return []*Exam{}, nil
	}

	// Retrieve available exams matching the criteria
	var exams []*Exam
	err = r.storage.DB().
		Preload("ExamType").
		Joins("JOIN exam_types ON exams.exam_type_id = exam_types.id").
		Where("exams.grade = ? AND exams.start > NOW() AND exam_types.order = ?", grade, nextOrder).
		Find(&exams).Error
	if err != nil {
		return nil, err
	}

	return exams, nil
}
