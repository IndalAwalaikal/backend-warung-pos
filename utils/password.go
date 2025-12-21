package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func CheckPassword(hash, pw string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
    return err == nil
}
