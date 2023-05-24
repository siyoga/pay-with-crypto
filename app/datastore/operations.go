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

func GetOneUnscopedBy[T All](key string, value interface{}) (T, bool) { // used in handler like: datastore.GetBy[datastore.User]("id", id)
	var i T

	result := Datastore.Unscoped().Where(map[string]interface{}{key: value}).First(&i)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Record Not Found" write error to log
			util.Error(result.Error, "GetUnscopedOneBy")
		}

		return i, false

	}

	return i, true
}

func GetUnscopedBy[T All](key string, value interface{}) ([]T, bool) {
	var i []T

	result := Datastore.Unscoped().Where(map[string]interface{}{key: value}).Find(&i)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Record Not Found" write error to log
			util.Error(result.Error, "GetOneUnscopedBy")
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

func GetAllOrdered[T All](key string, value interface{}, order string) ([]T, bool) { // used in handler like: datastore.GetBy[datastore.User]("id", id)
	var items []T

	result := Datastore.Where(map[string]interface{}{key: value}).Order(order).Find(&items)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			util.Error(result.Error, "GetAllCards")
		}

		return items, false
	}

	return items, true
}

func UpdateOneBy[T All](key string, value interface{}, updatedKey string, newValue interface{}) (T, bool) {
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

func UpdateOneUnscopedBy[T All](key string, value interface{}, updatedKey string, newValue interface{}) (T, bool) {
	var i T

	result := Datastore.Model(&i).Unscoped().Where(map[string]interface{}{key: value}).Update(updatedKey, newValue)

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
		result := Datastore.Model(&item).Where(map[string]interface{}{key: value}).Delete(&item)
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

func UnscopeCompanyByIdWithCards(companyID uuid.UUID) bool {
	company, _ := GetUserById(companyID.String())

	if len(company.Cards) > 0 {
		for _, card := range company.Cards {
			if ok := DeleteBy[Card]("id", &card.ID); !ok {
				return false
			}
		}
	}

	result := Datastore.Model(&company).Where(map[string]interface{}{"id": company.ID}).Delete(&company)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "UnscopeCompanyByIdWithCards")
		}
		return false
	}

	return true
}

func ScopeCompanyByIdWithCards(companyID uuid.UUID) bool {
	company, _ := GetUnscopedCompanyById(companyID.String())

	if _, ok := UpdateOneUnscopedBy[Company]("id", company.ID, "is_del", 0); !ok {
		return false
	}

	cards, _ := GetUnscopedBy[Card]("company_id", company.ID)
	if len(cards) > 0 {
		for _, card := range cards {

			if _, ok := UpdateOneUnscopedBy[Card]("id", card.ID, "is_del", 0); !ok {
				return false
			}
		}
	}

	return true
}

func Auth[T Authable](name string) (T, bool) {
	var item T

	result := Datastore.Where("name = ?", name).Find(&item)
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

	if result := Datastore.
		Where("name LIKE ?", "%"+value+"%").
		Find(&cards); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardByName")
		}
		return cards, false
	}

	return cards, true
}

func SearchCard(name string, tags string) ([]Card, bool) {
	var cards []Card

	splitedTags := pq.StringArray(strings.Split(tags, "&"))

	if result := Datastore.Where("name LIKE ? AND Tags && ?", "%"+name+"%", splitedTags).Find(&cards); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCard")
		}
		return cards, false
	}

	return cards, true
}

func SearchCardsByTags(rawTags string) ([]Card, bool) {
	var cards []Card

	splitedTags := pq.StringArray(strings.Split(rawTags, "&"))

	if result := Datastore.Where("Tags && ?", splitedTags).Find(&cards); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "SearchCardsByTags")
		}
		return cards, false
	}

	return cards, true
}

// Эту функцию не меняй на дженерик
func GetUserById(companyId string) (Company, bool) {
	var company Company

	result := Datastore.Model(Company{}).Where(map[string]interface{}{"id": companyId}).Preload("Cards").First(&company)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "GetUserById")
		}
		return company, false

	}
	return company, true
}

func GetUnscopedCompanyById(companyId string) (Company, bool) {
	var company Company

	result := Datastore.Model(Company{}).Unscoped().Where(map[string]interface{}{"id": companyId}).First(&company)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) { // if error NOT "Records Not Found" write error to log
			util.Error(result.Error, "GetUnscopedCompanyById")
		}
		return company, false

	}
	return company, true
}

func IsValid[T comparable](firstItem T, secondItem T) bool {

	return firstItem == secondItem
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
