package datastore

import (
	"errors"
	util "pay-with-crypto/app/utility"
	"strings"

	"github.com/gofrs/uuid"
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

func GetManyBy[T All](key string, value interface{}) ([]T, bool) { // used in handler like: datastore.GetBy[datastore.User]("id", id)
	var items []T

	result := Datastore.Where(map[string]interface{}{key: value}).Find(&items)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			util.Error(result.Error, "GetAllCards")
		}

		return items, false
	}

	return items, true
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

func WholeOneUpdate[T All](item T) bool {

	result := Datastore.Model(&item).Updates(item)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "WholeOneUpdate")
		}
		return false

	}

	return true
}

func DeleteBy[T All](key string, value any) bool {
	var item T
	var state bool

	_, found := GetOneBy[T](key, value)
	if found {
		result := Datastore.Model(&item).Delete(map[string]interface{}{key: value})
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
				util.Error(result.Error, "DeleteBy")
			}
			state = false
		}
		state = true
	}
	return state
}

func Auth[T Logineble](login string) (T, bool) {
	var item T

	result := Datastore.Where("login = ?", login).Find(&item)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "Auth")
		}
		return item, false

	}
	return item, true
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

func SearchCardsByTags(rawTags string) ([]Card, bool) {
	var cards []Card

	splitedTags := pq.StringArray(strings.Split(rawTags, "&"))

	result := Datastore.Where("Tags && ?", splitedTags).Find(&cards)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardsByTags")
		}

		return cards, false
	}

	return cards, true
}

// Эту функцию не меняй на дженерик
func GetUserById(userId string) (User, bool) {
	var user User

	result := Datastore.Model(User{}).Where(map[string]interface{}{"id": userId}).Preload("Cards").First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "GetUserById")
		}
		return user, false

	}
	return user, true
}

func IsCardValidToLoginedUser(cardId uuid.UUID, loginedUserId uuid.UUID) bool {
	var card Card

	card, found := GetOneBy[Card]("id", cardId)
	if !found {
		return false
	}

	return card.UserID == loginedUserId
}

func ShowCompanyById(userId uuid.UUID) (User, bool) {
	var state = true
	var user User

	result := Datastore.Select("id", "login", "mail", "link_to_company").Where("id = ?", userId).Find(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "ShowCompanyById")
		}

		state = false
	}

	return user, state
}

func AdminCheck() bool {
	var empty bool
	var admins Admin

	result := Datastore.Find(&admins)
	if r := result.RowsAffected; r == 0 {
		empty = true
	}
	return empty
}
