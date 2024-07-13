package utils

import (
    "fmt"
    "math/rand"
)

func GenerateOTP() (int,string, error) {
    
    otp := rand.Intn(1000000)

    hashedOTP, err := PwdSaltAndHash(fmt.Sprintf("%d", otp))
    // if err != nil {
    //     return "", fmt.Errorf("failed to hash OTP")
    // }

    return otp, hashedOTP, err
}
