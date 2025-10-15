// package services

// import (
// 	"JWT-Authentication-go/config"
// 	"JWT-Authentication-go/pkg/logging"

// 	"github.com/go-redis/redis"
// 	"github.com/go-redis/redis/v7"
// 	"github.com/gofiber/fiber/v2/middleware/cache"
// )


// type OtpService struct {
// 	logger      logging.Logger
// 	cfg         *config.Config
// 	redisClient *redis.Client
// }

// type OtpDto struct {
// 	Value string
// 	Used  bool
// }

// func NewOtpService(cfg *config.Config) *OtpService {
// 	logger := logging.NewLogger(cfg)
// 	redis := cache.GetRed
// }
