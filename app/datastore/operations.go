package datastore

import (
	"errors"
	util "pay-with-crypto/app/utility"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

func Add[T All](i T) bool {
	result := Datastore.Create(&i)

	if result.Error != nil {
		util.Error(result.Error, "Add")
		return false
	}

	return true
}

func GetOneBy[T All](key string, value interface{}) (T, bool) { // used in handler like: datastore.GetBy[datastore.User]("id", id)
	var i T

	result := Datastore.Where(&value).First(&i)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Record Not Found" write error to log
			util.Error(result.Error, "GetOneBy")
		}

		return i, false
	}

	return i, true
}

func SearchCardByName(nameOfCard string) ([]Card, bool) {
	var cards []Card

	result := Datastore.Where("Name LIKE %?%", nameOfCard).Find(&cards)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardByName")
		}

		return cards, false
	}

	return cards, true
}

func UserAuth(name string, password string) (User, bool) {
	var user User

	result := Datastore.Where("Company_Name = ?", name).Find(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "UserAuth")
		}
		return user, false

	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, false
	}

	return user, true
}

func GeneratToken(key string, payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(key)

	return t, err

}
