package app

import (
	"context"

	"github.com/JIeeiroSst/core-backend/pkg/dns"
	"github.com/JIeeiroSst/core-backend/pkg/email/smtp"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"

	"github.com/JIeeiroSst/core-backend/pkg/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/JIeeiroSst/core-backend/internal/config"
	delivery "github.com/JIeeiroSst/core-backend/internal/delivery/http/v1"
	"github.com/JIeeiroSst/core-backend/internal/repository"
	"github.com/JIeeiroSst/core-backend/internal/usecase"
	"github.com/JIeeiroSst/core-backend/pkg/auth"
	"github.com/JIeeiroSst/core-backend/pkg/cache"
	"github.com/JIeeiroSst/core-backend/pkg/database/mongodb"
	"github.com/JIeeiroSst/core-backend/pkg/hash"
	"github.com/JIeeiroSst/core-backend/pkg/logger"
	"github.com/JIeeiroSst/core-backend/pkg/otp"
)

// @title Creatly API
// @version 1.0
// @description REST API for Creatly App

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run(configPath string, gin *gin.Engine) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)

		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	otpGenerator := otp.NewGOTPGenerator()

	storageProvider, err := newStorageProvider(cfg)
	if err != nil {
		logger.Error(err)

		return
	}

	cloudflareClient, err := cloudflare.New(cfg.Cloudflare.ApiKey, cfg.Cloudflare.Email)
	if err != nil {
		logger.Error(err)

		return
	}

	dnsService := dns.NewService(cloudflareClient, cfg.Cloudflare.ZoneEmail, cfg.Cloudflare.CnameTarget)

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := usecase.NewUsecase(usecase.Deps{
		Repos:                  repos,
		Cache:                  memCache,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		EmailSender:            emailSender,
		EmailConfig:            cfg.Email,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		FondyCallbackURL:       cfg.Payment.FondyCallbackURL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		OtpGenerator:           otpGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		StorageProvider:        storageProvider,
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
		DNS:                    dnsService,
	})
	handlers := delivery.NewHandler(services, tokenManager)

	services.Files.InitStorageUploaderWorkers(context.Background())

	handlers.Init(gin)
}

func newStorageProvider(cfg *config.Config) (storage.Provider, error) {
	client, err := minio.New(cfg.FileStorage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.FileStorage.AccessKey, cfg.FileStorage.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	provider := storage.NewFileStorage(client, cfg.FileStorage.Bucket, cfg.FileStorage.Endpoint)

	return provider, nil
}
