package services

import (
	"JWT-Authentication-go/config"
	"JWT-Authentication-go/constants"
	"JWT-Authentication-go/data/cache"
	"JWT-Authentication-go/pkg/logging"
	"JWT-Authentication-go/pkg/service_errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)


type OtpService struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type OtpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpService{
		logger:      logger,
		cfg:         cfg,
		redisClient: redis,
	}
}

func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
		key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	val := &OtpDto{
		Value: otp,
		Used:  false,
	}
	res, err := cache.Get[OtpDto](s.redisClient, key)

		if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OptExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	err = cache.Set(s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (s *OtpService) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[OtpDto](s.redisClient, key)
	if err != nil {
		return err
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if err == nil && !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	} else if err == nil && !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(s.redisClient, key, res, s.cfg.Otp.ExpireTime*time.Second)

	}
	return nil
}