package datastore

import (
	"errors"
	util "pay-with-crypto/app/utility"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
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

	result := Datastore.Where(map[string]interface{}{key: value}).First(&i)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Record Not Found" write error to log
			util.Error(result.Error, "GetOneBy")
		}

		return i, false
	}

	return i, true
}

func SearchCardByName(value string) ([]Card, bool) {
	var cards []Card

	result := Datastore.Where("Name LIKE ?", "%"+value+"%").Find(&cards)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardByName")
		}

		return cards, false
	}

	return cards, true
}

func SearchCardsByTags(value string) ([]Card, bool) {
	var cards []Card

	splitedTags := pq.StringArray(strings.Split(value, "&"))

	result := Datastore.Where("Tags && ?", splitedTags).Find(&cards)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardByName")
		}

		return cards, false
	}

	return cards, true
}

func UpdateOneBy[T All](key string, value interface{}, updatedKey string, newValue string) (T, bool) {
	var i T

	result := Datastore.Model(&i).Where(map[string]interface{}{key: value}).Update(updatedKey, newValue)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Record Not Found" write error to log
			util.Error(result.Error, "UpdateOneBy")
		}

		return i, false
	}

	return i, true
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
	return user, true
}
