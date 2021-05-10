package service

import (
	"context"
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/service/postgres"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

//init
//TODO:
// mock db,
// sometimes some tests have additional spaces in message :(
func init() {
	lis = bufconn.Listen(bufSize)
	log.Println("server is started")
	server := grpc.NewServer()
	api.RegisterUserServiceServer(server, &GRPCServer{})

	dbConnection := postgres.OpenDataBaseConnection()
	postgres.StorageInstance = postgres.NewStorage(dbConnection)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		caseName          string
		createUserRequest api.CreateUserRequest
		result            *api.User
		isPositive        bool
		errCode           codes.Code
		errMsg            string
	}{
		{
			caseName: "Valid CreateUserRequest without items",
			createUserRequest: api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    createItemRequest(createItemsData(0)...),
			},
			isPositive: true,
		},
		{
			caseName: "Valid CreateUserRequest with 1 item",
			createUserRequest: api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    createItemRequest(createItemsData(1)...),
			},
			isPositive: true,
		},
		{
			caseName: "Invalid CreateUserRequest, missedName",
			createUserRequest: api.CreateUserRequest{
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    createItemRequest(createItemsData(0)...),
			},
			isPositive: false,
			errCode:    codes.InvalidArgument,
			errMsg:     "User validation failed: user - 'age:123  user_type:EMPLOYEE_USER_TYPE', err - name is missed",
		},
		{
			caseName: "Invalid CreateUserRequest, age is negative",
			createUserRequest: api.CreateUserRequest{
				Age:      -1,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    createItemRequest(createItemsData(0)...),
			},
			isPositive: false,
			errCode:    codes.InvalidArgument,
			errMsg:     "User validation failed: user - 'age:-1  user_type:EMPLOYEE_USER_TYPE', err - age of user must be positive, age = -1",
		},
		{
			caseName: "Invalid CreateUserRequest, missedName in item",
			createUserRequest: api.CreateUserRequest{
				Name:     "testName",
				Age:      123,
				UserType: api.UserType_EMPLOYEE_USER_TYPE,
				Items:    createItemRequest(createInvalidItemsData(1)...),
			},
			isPositive: false,
			errCode:    codes.InvalidArgument,
			errMsg:     "Item validation failed: item - '', err - name is missed",
		},
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(bufDialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := api.NewUserServiceClient(conn)

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			user, err := client.CreateUser(ctx, &tc.createUserRequest)
			if tc.isPositive {
				assert.Empty(t, err)
				assert.NotEmpty(t, user.Id)
				assert.Equal(t, tc.createUserRequest.Name, user.Name)
				assert.Equal(t, tc.createUserRequest.Age, user.Age)
				assert.Equal(t, tc.createUserRequest.UserType, user.UserType)
				assert.NotEmpty(t, user.CreatedAt)
				assert.Empty(t, user.UpdatedAt)
				deleteUser(t, ctx, client, user.Id)
			} else {
				assert.NotEmpty(t, err)
				fromError, _ := status.FromError(err)
				assert.Equal(t, tc.errCode, fromError.Code())
				assert.Equal(t, tc.errMsg, fromError.Message())
			}
		})
	}
}

func TestGetUser(t *testing.T) {

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(bufDialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := api.NewUserServiceClient(conn)

	userER := createUser(t, ctx, client, 1)

	testCases := []struct {
		caseName       string
		getUserRequest api.GetUserRequest
		result         *api.User
		isPositive     bool
		errCode        codes.Code
		errMsg         string
	}{
		{
			caseName: "Get userER by id",
			getUserRequest: api.GetUserRequest{
				Id: userER.Id,
			},
			isPositive: true,
		},
		{
			caseName: "Get userER by id",
			getUserRequest: api.GetUserRequest{
				Id: "12324789",
			},
			isPositive: false,
			errMsg:     "GetUser: User not found by id = 12324789",
			errCode:    codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			userAR, err := client.GetUser(ctx, &tc.getUserRequest)
			if tc.isPositive {
				assert.Empty(t, err)
				assert.Equal(t, userER.Id, userAR.Id)
				assert.Equal(t, userER.Name, userAR.Name)
				assert.Equal(t, userER.Age, userAR.Age)
				assert.Equal(t, userER.UserType, userAR.UserType)
				assert.NotEmpty(t, userAR.CreatedAt)
				assert.Empty(t, userAR.UpdatedAt)
				assert.NotEmpty(t, userAR.Items)
				assert.Equal(t, 1, len(userAR.Items))
				assert.Equal(t, userAR.Id, userAR.Items[0].UserId)
				assert.NotEmpty(t, userAR.Items[0].Id)
				assert.NotEmpty(t, userAR.Items[0].Name)
				assert.NotEmpty(t, userAR.Items[0].CreatedAt)
				assert.Empty(t, userAR.Items[0].UpdatedAt)
			} else {
				assert.NotEmpty(t, err)
				fromError, _ := status.FromError(err)
				assert.Equal(t, tc.errCode, fromError.Code())
				assert.Equal(t, tc.errMsg, fromError.Message())
			}
		})
	}
	deleteUser(t, ctx, client, userER.Id)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(bufDialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := api.NewUserServiceClient(conn)
	userER := createUser(t, ctx, client, 1)

	testCases := []struct {
		caseName          string
		updateUserRequest api.UpdateUserRequest
		result            *api.User
		isPositive        bool
		errCode           codes.Code
		errMsg            string
	}{
		{
			caseName: "Update User",
			updateUserRequest: api.UpdateUserRequest{
				Id:       userER.Id,
				Name:     userER.Name,
				Age:      999,
				UserType: 0,
				Items: []*api.UpdateItemRequest{
					{
						Id:   userER.Items[0].Id,
						Name: "updatedItem",
					}},
			},
			isPositive: true,
		},
		{
			caseName: "Update User, missed user id",
			updateUserRequest: api.UpdateUserRequest{
				Name:     userER.Name,
				Age:      999,
				UserType: 0,
				Items: []*api.UpdateItemRequest{
					{
						Id:   userER.Items[0].Id,
						Name: "updatedItem",
					}},
			},
			isPositive: false,
			errMsg:     fmt.Sprintf("User validation failed: user - 'name:\"testName\"  age:999  items:{id:\"%s\"  name:\"updatedItem\"}', err - id is missed", userER.Items[0].Id),
			errCode:    codes.InvalidArgument,
		},
		{
			caseName: "Update User, missed item id",
			updateUserRequest: api.UpdateUserRequest{
				Id:       userER.Id,
				Name:     userER.Name,
				Age:      999,
				UserType: 0,
				Items: []*api.UpdateItemRequest{
					{
						Name: "updatedItem",
					}},
			},
			isPositive: false,
			errMsg:     "Item validation failed: item - 'name:\"updatedItem\"', err - 'id is missed'",
			errCode:    codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			userAR, err := client.UpdateUser(ctx, &tc.updateUserRequest)
			if tc.isPositive {
				assert.Empty(t, err)
				assert.NotEmpty(t, userAR.Id)
				assert.Equal(t, userER.Name, userAR.Name)
				assert.NotEqual(t, userER.Age, userAR.Age)
				assert.Equal(t, int32(999), userAR.Age)
				assert.Equal(t, api.UserType_INVALID_USER_TYPE, userAR.UserType)
				assert.NotEmpty(t, userAR.CreatedAt)
				assert.NotEmpty(t, userAR.UpdatedAt)
				assert.NotEmpty(t, userAR.Items)
				assert.Equal(t, 1, len(userAR.Items))
				assert.Equal(t, userAR.Id, userAR.Items[0].UserId)
				assert.NotEmpty(t, userAR.Items[0].Id)
				assert.NotEmpty(t, userAR.Items[0].CreatedAt)
				assert.NotEmpty(t, userAR.Items[0].UpdatedAt)
				assert.Equal(t, "updatedItem", userAR.Items[0].Name)
			} else {
				assert.NotEmpty(t, err)
				fromError, _ := status.FromError(err)
				assert.Equal(t, tc.errCode, fromError.Code())
				assert.Equal(t, tc.errMsg, fromError.Message())
			}
		})
	}
	deleteUser(t, ctx, client, userER.Id)
}

//TestListUser TODO update to check response users
func TestListUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(bufDialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := api.NewUserServiceClient(conn)
	user1 := createUser(t, ctx, client, 0)
	user2 := createUser(t, ctx, client, 0)
	user3 := createUser(t, ctx, client, 0)
	user4 := createUser(t, ctx, client, 0)

	testCases := []struct {
		caseName        string
		listUserRequest api.ListUserRequest
		resultLen       int
		isPositive      bool
		errCode         codes.Code
		errMsg          string
	}{
		{
			caseName: "ListUser: limit = 1, page = 1",
			listUserRequest: api.ListUserRequest{
				PageFilter: &api.PageFilter{
					Limit: 1,
					Page:  1,
				},
			},
			resultLen:  1,
			isPositive: true,
		},
		{
			caseName: "ListUser: limit = 2, page = 1",
			listUserRequest: api.ListUserRequest{
				PageFilter: &api.PageFilter{
					Limit: 2,
					Page:  1,
				},
			},
			resultLen:  2,
			isPositive: true,
		},
		{
			caseName: "ListUser: limit = 2, page = 2",
			listUserRequest: api.ListUserRequest{
				PageFilter: &api.PageFilter{
					Limit: 2,
					Page:  2,
				},
			},
			resultLen:  2,
			isPositive: true,
		},
		{
			caseName: "ListUser: limit = 2, page = 0",
			listUserRequest: api.ListUserRequest{
				PageFilter: &api.PageFilter{
					Limit: 2,
					Page:  0,
				},
			},
			isPositive: false,
			errMsg:     "page must be > 0, page = 0",
			errCode:    codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			users, err := client.ListUser(ctx, &tc.listUserRequest)
			if tc.isPositive {
				assert.Empty(t, err)
				assert.Equal(t, tc.resultLen, len(users.Users))
			} else {
				assert.NotEmpty(t, err)
				fromError, _ := status.FromError(err)
				assert.Equal(t, tc.errCode, fromError.Code())
				assert.Equal(t, tc.errMsg, fromError.Message())
			}
		})
	}
	deleteUser(t, ctx, client, user1.Id)
	deleteUser(t, ctx, client, user2.Id)
	deleteUser(t, ctx, client, user3.Id)
	deleteUser(t, ctx, client, user4.Id)
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(bufDialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := api.NewUserServiceClient(conn)
	user := createUser(t, ctx, client, 0)
	deleteUser(t, ctx, client, user.Id)
	userAfterDelete, err := getUser(ctx, client, user.Id)

	assert.NotNil(t, err)
	fromError, _ := status.FromError(err)
	assert.Nil(t, userAfterDelete)
	assert.Equal(t, codes.NotFound, fromError.Code())
	assert.Equal(t, fmt.Sprintf("GetUser: User not found by id = %s", user.GetId()), fromError.Message())
}

func deleteUser(t *testing.T, ctx context.Context, client api.UserServiceClient, userId string) {
	deleteUserRequest := &api.DeleteUserRequest{Id: userId}
	_, err := client.DeleteUser(ctx, deleteUserRequest)
	assert.Empty(t, err)
}

func getUser(ctx context.Context, client api.UserServiceClient, userId string) (*api.User, error) {
	getUserRequest := &api.GetUserRequest{
		Id: userId,
	}
	return client.GetUser(ctx, getUserRequest)
}

func createUser(t *testing.T, ctx context.Context, client api.UserServiceClient, numberOFItems int) *api.User {
	createUserRequest := &api.CreateUserRequest{
		Name:     "testName",
		Age:      123,
		UserType: api.UserType_EMPLOYEE_USER_TYPE,
		Items:    createItemRequest(createItemsData(numberOFItems)...)}
	userER, err := client.CreateUser(ctx, createUserRequest)
	assert.Empty(t, err)
	assert.NotNil(t, userER)
	assert.NotEmpty(t, userER.Id)
	return userER
}

func createItemsData(count int) []*api.Item {
	items := make([]*api.Item, 0, count)
	for i := 1; i <= count; i++ {
		itemName := fmt.Sprintf("Im item #%d", i)
		items = append(items, &api.Item{Name: itemName})
	}
	return items
}

func createItemRequest(items ...*api.Item) []*api.CreateItemRequest {
	createItemRequests := make([]*api.CreateItemRequest, 0, len(items))
	for _, item := range items {
		createItemRequests = append(createItemRequests, &api.CreateItemRequest{Name: item.GetName()})
	}
	return createItemRequests
}

func createInvalidItemsData(count int) []*api.Item {
	items := make([]*api.Item, 0, count)
	for i := 1; i <= count; i++ {
		items = append(items, &api.Item{Name: ""})
	}
	return items
}
