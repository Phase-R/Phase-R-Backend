package utils

import (
    "log"
    "strconv"
    "math/rand"
    "github.com/alexedwards/argon2id"
)

func GenerateOTP() (int,string, error) {
    
    
    otp := rand.Intn(900000)+100000

    otpString := strconv.Itoa(otp)

	// Hash the OTP string using argon2id
	hashedOTP, err := argon2id.CreateHash(otpString, argon2id.DefaultParams)
	if err != nil {
		log.Fatal("could not hash OTP", err)
		return 0, "", err
	}

    return otp, hashedOTP, err
}
