package validation

import (
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

type createUserRequestTC struct {
	caseName          string
	createUserRequest *api.CreateUserRequest
	expectedErrorMsg  string
}

func TestValidateCreateUserRequestData_shouldNotReturnError_whenRequestIsValid(t *testing.T) {
	validTestCases := []createUserRequestTC{
		{
			caseName: "User with 0 items",
			createUserRequest: &api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    initCreateItemRequest(createItems(0)...),
			},
		},
		{
			caseName: "User with 2 items",
			createUserRequest: &api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    initCreateItemRequest(createItems(2)...),
			},
		},
	}

	for _, tc := range validTestCases {
		t.Run(tc.caseName, func(t *testing.T) {
			assert.Equal(t, nil, ValidateCreateUserRequestData(tc.createUserRequest))
		})
	}
}

func TestValidateCreateUserRequestData_shouldReturnError_whenRequestIsNotValid(t *testing.T) {
	validTestCases := []createUserRequestTC{
		{
			caseName: "User with 0 items",
			createUserRequest: &api.CreateUserRequest{
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    initCreateItemRequest(createItems(0)...),
			},
			expectedErrorMsg: "User validation failed: user - 'age:123 user_type:EMPLOYEE_USER_TYPE', err - name is missed",
		},
		{
			caseName: "User with 2 items",
			createUserRequest: &api.CreateUserRequest{
				Name:     "testName",
				Age:      -123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    initCreateItemRequest(createItems(2)...),
			},
			expectedErrorMsg: "User validation failed: user - 'name:\"testName\" age:-123 user_type:EMPLOYEE_USER_TYPE items:{name:\"Im item #1\"} items:{name:\"Im item #2\"}', err - age of user must be positive, age = -123",
		},
		{
			caseName: "User with 2 items",
			createUserRequest: &api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    initCreateItemRequest(createInvalidItems(1)...),
			},
			expectedErrorMsg: "Item validation failed: item - '', err - name is missed",
		},
	}

	for _, tc := range validTestCases {
		t.Run(tc.caseName, func(t *testing.T) {
			assert.Equal(t, tc.expectedErrorMsg, ValidateCreateUserRequestData(tc.createUserRequest).Error())
		})
	}
}

func TestSomethingElse(t *testing.T) {
	//etc
}

func createInvalidItems(count int) []*api.Item {
	items := make([]*api.Item, 0, count)
	for i := 1; i <= count; i++ {
		items = append(items, &api.Item{Name: ""})
	}
	return items
}

func createItems(count int) []*api.Item {
	items := make([]*api.Item, 0, count)
	for i := 1; i <= count; i++ {
		itemName := fmt.Sprintf("Im item #%d", i)
		items = append(items, &api.Item{Name: itemName})
	}
	return items
}

func initCreateItemRequest(items ...*api.Item) []*api.CreateItemRequest {
	createItemRequests := make([]*api.CreateItemRequest, 0, len(items))
	for _, item := range items {
		createItemRequests = append(createItemRequests, &api.CreateItemRequest{Name: item.GetName()})
	}
	return createItemRequests
}
