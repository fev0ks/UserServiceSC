package validation

import (
	"errors"
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
)

type AgeData interface {
	GetAge() int32
}

type NameData interface {
	GetName() string
}

type IdData interface {
	GetId() string
}

type CreateUserData interface {
	AgeData
	NameData
	GetItems() []*api.CreateItemRequest
}

type UserData interface {
	AgeData
	NameData
	IdData
	GetItems() []*api.UpdateItemRequest
}

type NewItemData interface {
	NameData
}

type ItemData interface {
	NameData
	IdData
}

type PageFilterData interface {
	GetPageFilter() *api.PageFilter
}

func ValidateCreateUserRequestData(userData CreateUserData) error {
	if err := ValidateAge(userData); err != nil {
		return createUserValidationErrorFmt(userData, err)
	}
	if err := ValidateName(userData); err != nil {
		return createUserValidationErrorFmt(userData, err)
	}
	for _, item := range userData.GetItems() {
		if err := validateCreateItemRequestData(item); err != nil {
			return errors.New(fmt.Sprintf("Item validation failed: item - '%v', err - %v", item.String(), err.Error()))
		}
	}
	return nil
}

func ValidateUserRequestData(userData UserData) error {
	if err := ValidateId(userData); err != nil {
		return userValidationErrorFmt(userData, err)
	}
	if err := ValidateAge(userData); err != nil {
		return userValidationErrorFmt(userData, err)
	}
	if err := ValidateName(userData); err != nil {
		return userValidationErrorFmt(userData, err)
	}
	for _, item := range userData.GetItems() {
		if err := validateItemRequestData(item); err != nil {
			return errors.New(fmt.Sprintf("Item validation failed: item - '%v', err - '%v'", item.String(), err.Error()))
		}
	}
	return nil
}

func userValidationErrorFmt(userData UserData, err error) error {
	return errors.New(fmt.Sprintf("User validation failed: user - '%v', err - %v", userData, err.Error()))
}

func createUserValidationErrorFmt(userData CreateUserData, err error) error {
	return errors.New(fmt.Sprintf("User validation failed: user - '%v', err - %v", userData, err.Error()))
}

func validateCreateItemRequestData(itemDate NewItemData) error {
	if err := ValidateName(itemDate); err != nil {
		return err
	}
	return nil
}

func validateItemRequestData(itemData ItemData) error {
	if err := ValidateId(itemData); err != nil {
		return err
	}
	if err := ValidateName(itemData); err != nil {
		return err
	}
	return nil
}

func ValidateId(idData IdData) error {
	if idData.GetId() == "" {
		return errors.New("id is missed")
	}
	return nil
}

//ValidateAge TODO age may be bigger than MAX int32 -> as result age < 0
func ValidateAge(userData AgeData) error {
	if userData.GetAge() <= 0 {
		return errors.New(fmt.Sprintf("age of user must be positive, age = %d", userData.GetAge()))
	}
	return nil
}

func ValidateName(userData NameData) error {
	if userData.GetName() == "" {
		return errors.New("name is missed")
	}
	return nil
}

//ValidatePageFilter TODO page and limit are uint type if input value = -n then result value = MAX.INT-n ...
func ValidatePageFilter(pageFilterData PageFilterData) error {
	if pageFilterData.GetPageFilter() == nil {
		msg := fmt.Sprintf("pageFilter is missed")
		return errors.New(msg)
	}
	if pageFilterData.GetPageFilter().Page <= 0 {
		msg := fmt.Sprintf("page must be > 0, page = %d", pageFilterData.GetPageFilter().Page)
		return errors.New(msg)
	}
	if pageFilterData.GetPageFilter().Limit <= 0 {
		msg := fmt.Sprintf("limit must be > 0, limit = %d", pageFilterData.GetPageFilter().Limit)
		return errors.New(msg)
	}
	return nil
}
