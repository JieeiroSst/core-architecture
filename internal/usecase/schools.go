package usecase

import (
	"context"

	"github.com/JIeeiroSst/core-backend/pkg/payment"
	"github.com/JIeeiroSst/core-backend/pkg/payment/fondy"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/JIeeiroSst/core-backend/internal/domain"
	"github.com/JIeeiroSst/core-backend/internal/repository"
	"github.com/JIeeiroSst/core-backend/pkg/cache"
)

type Schools interface {
	Create(ctx context.Context, name string) (primitive.ObjectID, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
	UpdateSettings(ctx context.Context, schoolId primitive.ObjectID, input domain.UpdateSchoolSettingsInput) error
	ConnectFondy(ctx context.Context, input ConnectFondyInput) error
	ConnectSendPulse(ctx context.Context, input ConnectSendPulseInput) error
}

type SchoolsService struct {
	repo  repository.Schools
	cache cache.Cache
	ttl   int64
}

func NewSchoolsService(repo repository.Schools, cache cache.Cache, ttl int64) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache, ttl: ttl}
}

func (s *SchoolsService) Create(ctx context.Context, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, name)
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	err = s.cache.Set(domainName, school, s.ttl)

	return school, err
}

func (s *SchoolsService) GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error) {
	return s.repo.GetById(ctx, id)
}

func (s *SchoolsService) UpdateSettings(ctx context.Context, schoolId primitive.ObjectID, inp domain.UpdateSchoolSettingsInput) error {
	return s.repo.UpdateSettings(ctx, schoolId, inp)
}

func (s *SchoolsService) ConnectFondy(ctx context.Context, input ConnectFondyInput) error {
	client := fondy.NewFondyClient(input.MerchantID, input.MerchantPassword)

	id := primitive.NewObjectID()

	_, err := client.GeneratePaymentLink(payment.GeneratePaymentLinkInput{
		OrderId:   id.Hex(),
		Amount:    1000,
		Currency:  "USD",
		OrderDesc: "CREATLY - TESTING FONDY CREDENTIALS",
	})
	if err != nil {
		return err
	}

	creds := domain.Fondy{
		MerchantPassword: input.MerchantPassword,
		MerchantID:       input.MerchantID,
		Connected:        true,
	}

	return s.repo.SetFondyCredentials(ctx, input.SchoolID, creds)
}

func (s *SchoolsService) ConnectSendPulse(ctx context.Context, input ConnectSendPulseInput) error {
	// todo
	return nil
}
