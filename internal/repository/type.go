package repository

import (
	"github.com/JIeeiroSst/core-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateCourseInput struct {
	ID          primitive.ObjectID
	SchoolID    primitive.ObjectID
	Name        *string
	ImageURL    *string
	Description *string
	Color       *string
	Published   *bool
}

type UpdateModuleInput struct {
	ID        primitive.ObjectID
	SchoolID  primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type UpdateLessonInput struct {
	ID        primitive.ObjectID
	SchoolID  primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type UpdatePackageInput struct {
	ID       primitive.ObjectID
	SchoolID primitive.ObjectID
	Name     string
}

type UpdateOfferInput struct {
	ID            primitive.ObjectID
	SchoolID      primitive.ObjectID
	Name          string
	Description   string
	Benefits      []string
	Price         *domain.Price
	Packages      []primitive.ObjectID
	PaymentMethod *domain.PaymentMethod
}
