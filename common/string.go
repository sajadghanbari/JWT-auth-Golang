package common

import (
	"JWT-Authentication-go/config"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GenerateOtp() string {
	cfg := config.GetConfig()
	rand.Seed(time.Now().UnixNano())
	min := int(math.Pow(10, float64(cfg.Otp.Digits-1))) 
	max := int(math.Pow(10, float64(cfg.Otp.Digits)) - 1) 

	var num = rand.Intn(max-min) + min
	return strconv.Itoa(num)
}