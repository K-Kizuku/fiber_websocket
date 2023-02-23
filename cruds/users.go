package cruds

import (
	"errors"
	"fiber_websocket/db"
	"fiber_websocket/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(name string, email string, password string) (res_user db.UserInfo, err error) {
	if err = db.Db.Where("email = ?", email).First(&db.UserInfo{}).Error; err == nil {
		err = errors.New("email is already exist")
		return
	}
	hash_pass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	res_user = db.UserInfo{Email: email, Name: name, PasswordHash: string(hash_pass)}
	err = db.Db.Create(&res_user).Error
	return
}

func GenerateJWT(email string, password string) (jwtInfo db.JWTInfo, err error) {
	var (
		u     db.UserInfo
		token string
	)

	if err = db.Db.Where("email = ?", email).First(&u).Error; err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return
	}

	token, err = generateToken(u.ID)
	if err != nil {
		return
	}

	jwtInfo = db.JWTInfo{JWT: token}
	return
}

func generateToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(utils.SigningKey)
}
