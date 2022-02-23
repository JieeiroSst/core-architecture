package usecase

import (
	"context"

	"github.com/JIeeiroSst/core-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentLessons interface {
	AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error
}

type StudentLessonsService struct {
	repo repository.StudentLessons
}

func NewStudentLessonsService(repo repository.StudentLessons) *StudentLessonsService {
	return &StudentLessonsService{
		repo: repo,
	}
}

func (s *StudentLessonsService) AddFinished(ctx context.Context, studentID, lessonID primitive.ObjectID) error {
	return s.repo.AddFinished(ctx, studentID, lessonID)
}

func (s *StudentLessonsService) SetLastOpened(ctx context.Context, studentID, lessonID primitive.ObjectID) error {
	return s.repo.SetLastOpened(ctx, studentID, lessonID)
}

