package usecase

import (
	"io"
	"time"

	"github.com/JIeeiroSst/core-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type ConnectFondyInput struct {
	SchoolID         primitive.ObjectID
	MerchantID       string
	MerchantPassword string
}

type ConnectSendPulseInput struct {
	SchoolID primitive.ObjectID
	ID       string
	Secret   string
	ListID   string
}

type StudentSignUpInput struct {
	Name         string
	Email        string
	Password     string
	SchoolID     primitive.ObjectID
	SchoolDomain string
	Verified     bool
}

type SchoolSignInInput struct {
	Email    string
	Password string
	SchoolID primitive.ObjectID
}

type UploadInput struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
	SchoolID    primitive.ObjectID
	Type        domain.FileType
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
	Domain           string
}

type StudentPurchaseSuccessfulEmailInput struct {
	Email      string
	Name       string
	CourseName string
}

type UpdateCourseInput struct {
	CourseID    string
	SchoolID    string
	Name        *string
	ImageURL    *string
	Description *string
	Color       *string
	Published   *bool
}

type CreateModuleInput struct {
	SchoolID string
	CourseID string
	Name     string
	Position uint
}

type UpdateModuleInput struct {
	ID        string
	SchoolID  string
	Name      string
	Position  *uint
	Published *bool
}

type AddLessonInput struct {
	ModuleID string
	SchoolID string
	Name     string
	Position uint
}

type UpdateLessonInput struct {
	LessonID  string
	SchoolID  string
	Name      string
	Content   string
	Position  *uint
	Published *bool
}

type CreatePackageInput struct {
	CourseID string
	SchoolID string
	Name     string
	Modules  []string
}

type UpdatePackageInput struct {
	ID       string
	SchoolID string
	Name     string
	Modules  []string
}

type CreateSurveyInput struct {
	ModuleID primitive.ObjectID
	SchoolID primitive.ObjectID
	Survey   domain.Survey
}

type SaveStudentAnswersInput struct {
	ModuleID  primitive.ObjectID
	StudentID primitive.ObjectID
	SchoolID  primitive.ObjectID
	Answers   []domain.SurveyAnswer
}

type CreatePromoCodeInput struct {
	SchoolID           primitive.ObjectID
	Code               string
	DiscountPercentage int
	ExpiresAt          time.Time
	OfferIDs           []primitive.ObjectID
}

type CreateOfferInput struct {
	Name          string
	Description   string
	Benefits      []string
	SchoolID      primitive.ObjectID
	Price         domain.Price
	Packages      []string
	PaymentMethod domain.PaymentMethod
}

type UpdateOfferInput struct {
	ID            string
	SchoolID      string
	Name          string
	Description   string
	Benefits      []string
	Price         *domain.Price
	Packages      []string
	PaymentMethod *domain.PaymentMethod
}

func (i UpdateOfferInput) ValidatePayment() error {
	if i.PaymentMethod == nil {
		return nil
	}

	if !i.PaymentMethod.UsesProvider {
		return nil
	}

	return i.PaymentMethod.Validate()
}
