package validation

import (
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/service/errorhandler"
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

type NewUserData interface {
	AgeData
	NameData
}

type UserData interface {
	AgeData
	NameData
	IdData
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

func ValidateBaseUserRequestData(userData NewUserData) error {
	if err := ValidateAge(userData); err != nil {
		return err
	}
	if err := ValidateName(userData); err != nil {
		return err
	}
	return nil
}

func ValidateUserRequestData(userData UserData) error {
	if err := ValidateId(userData); err != nil {
		return err
	}
	if err := ValidateBaseUserRequestData(userData); err != nil {
		return err
	}
	return nil
}

// ValidateBaseItemRequestData TODO check items data as well
func ValidateBaseItemRequestData(itemData NewItemData) error {
	if err := ValidateName(itemData); err != nil {
		return err
	}
	return nil
}

// ValidateItemRequestData TODO check items data as well
func ValidateItemRequestData(itemData ItemData) error {
	if err := ValidateId(itemData); err != nil {
		return err
	}
	if err := ValidateBaseItemRequestData(itemData); err != nil {
		return err
	}
	return nil
}

func ValidateId(idData IdData) error {
	if idData.GetId() == "" {
		return errorhandler.NewInvalidArgumentError("User Id is missed")
	}
	return nil
}

//ValidateAge TODO age may be bigger than MAX int32 -> as result age < 0
func ValidateAge(userData AgeData) error {
	if userData.GetAge() <= 0 {
		return errorhandler.NewInvalidArgumentError(fmt.Sprintf("Age of User must be positive, age = %d", userData.GetAge()))
	}
	return nil
}

func ValidateName(userData NameData) error {
	if userData.GetName() == "" {
		return errorhandler.NewInvalidArgumentError("User Name is missed")
	}
	return nil
}

//ValidatePageFilter TODO page and limit are uint type if input value = -n then result value = MAX.INT-n ...
func ValidatePageFilter(pageFilterData PageFilterData) error {
	if pageFilterData.GetPageFilter() == nil {
		msg := fmt.Sprintf("PageFilter is missed")
		return errorhandler.NewInvalidArgumentError(msg)
	}
	if pageFilterData.GetPageFilter().Page <= 0 {
		msg := fmt.Sprintf("Page must be > 0, page = %d", pageFilterData.GetPageFilter().Page)
		return errorhandler.NewInvalidArgumentError(msg)
	}
	if pageFilterData.GetPageFilter().Limit <= 0 {
		msg := fmt.Sprintf("Limit must be > 0, limit = %d", pageFilterData.GetPageFilter().Limit)
		return errorhandler.NewInvalidArgumentError(msg)
	}
	return nil
}
