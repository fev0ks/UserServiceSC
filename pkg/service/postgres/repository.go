package postgres

import (
	"database/sql"
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/service/errorhandler"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"time"
)

const (
	InsertUserQuery     = "INSERT INTO \"user\"(name, age) VALUES($1, $2) RETURNING id, created_at; "
	InsertItemQuery     = "INSERT INTO \"item\"(name) VALUES %s RETURNING id, name, created_at; "
	InsertUserTypeQuery = "INSERT INTO user_type(user_id, type_id) VALUES($1, $2); "
	InsertUserItemQuery = "INSERT INTO user_item(user_id, item_id) VALUES %s; "
	SelectUserQuery     = "SELECT " +
		"us.id, us.name userName, us.age userAge, type_id userType, us.created_at userCreatedAt, us.updated_at userUpdatedAt, " +
		"item.id itemId, item.name itemName, item.created_at itemCreatedAt, item.updated_at itemUpdatedAt " +
		"FROM \"user\" us " +
		"inner join \"user_type\" on user_type.user_id = us.id " +
		"left join \"user_item\" on user_item.user_id = us.id " +
		"left join \"item\" item on user_item.item_id = item.id " +
		"where us.id = $1; "
	SelectUsersQuery = "SELECT " +
		"us.id, us.name userName, us.age userAge, type_id userType, us.created_at userCreatedAt, us.updated_at userUpdatedAt, " +
		"item.id itemId, item.name itemName, item.created_at itemCreatedAt, item.updated_at itemUpdatedAt " +
		"FROM \"user\" us \ninner join \"user_type\" on user_type.user_id = us.id " +
		"left join \"user_item\" on user_item.user_id = us.id " +
		"left join \"item\" item on user_item.item_id = item.id " +
		"where " +
		"us.id in (select id from \"user\" order by id LIMIT $1 OFFSET $2) " +
		"order by us.id"
	DeleteItemQuery     = "DELETE FROM item where id in (select item_id from user_item where user_id = $1); "
	DeleteUserQuery     = "DELETE FROM \"user\" where id = $1; "
	UpdateUserQuery     = "UPDATE \"user\" set name = $2, age = $3, updated_at = $4 where id = $1; "
	UpdateItemQuery     = "UPDATE \"item\" set name = %s, updated_at = %s where id = %s; "
	UpdateUserTypeQuery = "UPDATE \"user_type\" set type_id = $2 where user_id = $1; "
)

var StorageInstance *Storage

func CreateUser(data *api.CreateUserRequest) (*api.User, error) {

	tx, err := StorageInstance.DB.Begin()
	if err != nil {
		errorhandler.LogMsg("CreateUser: StorageInstance.DB.Begin")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	defer tx.Rollback()

	user, err := createUser(tx, data.GetName(), data.GetAge(), data.GetUserType())
	if err != nil {
		errorhandler.LogMsg("CreateUser: createUser")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	items, err := createItems(tx, user.Id, data.GetItems())
	if err != nil {
		errorhandler.LogMsg("CreateUser: createItems")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	user.Items = items

	if err := tx.Commit(); err != nil {
		errorhandler.LogMsg("CreateUser: tx.Commit")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	return user, nil
}

func createUser(tx *sql.Tx, name string, age int32, userType api.UserType) (*api.User, error) {
	var (
		userId    string
		createdAt time.Time
	)

	stmt, err := tx.Prepare(InsertUserQuery)
	if err != nil {
		errorhandler.LogMsg("CreateUser: tx.Prepare(InsertUserQuery)")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	defer stmt.Close()
	err = stmt.QueryRow(name, age).Scan(&userId, &createdAt)
	if err != nil {
		errorhandler.LogMsg("CreateUser: stmt.QueryRow")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	if err = setUserType(tx, userId, userType); err != nil {
		errorhandler.LogMsg("CreateUser: setUserType")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	return &api.User{
			Id:        userId,
			Name:      name,
			Age:       age,
			UserType:  userType,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: nil},
		nil
}

func setUserType(tx *sql.Tx, userId string, userType api.UserType) error {
	stmt, err := tx.Prepare(InsertUserTypeQuery)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("setUserType: tx.Prepare(%s)", InsertUserTypeQuery))
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(userId, userType); err != nil {
		errorhandler.LogMsg(fmt.Sprintf("setUserType: stmt.Exec(%s, %v)", userId, userType))
		return err
	}
	return nil
}

func createItems(tx *sql.Tx, userId string, data []*api.CreateItemRequest) ([]*api.Item, error) {
	if len(data) > 0 {
		var items = make([]*api.Item, 0, len(data))
		valueStrings := make([]string, 0, len(items))
		valueArgs := make([]interface{}, 0, len(items))
		i := 1
		for _, item := range data {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i))
			valueArgs = append(valueArgs, item.Name)
			i++
		}
		query := fmt.Sprintf(InsertItemQuery, strings.Join(valueStrings, ","))
		stmt, err := tx.Prepare(query)
		if err != nil {
			errorhandler.LogMsg(fmt.Sprintf("createItems: tx.Prepare(%s)", query))
			return nil, errorhandler.NewInternalError(err.Error())
		}

		defer stmt.Close()
		rows, err := stmt.Query(valueArgs...)
		if err != nil {
			errorhandler.LogMsg(fmt.Sprintf("createItems: tmt.Query(%v)", valueArgs))
			return nil, errorhandler.NewInternalError(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var (
				itemId    int
				name      string
				createdAt time.Time
			)
			err := rows.Scan(&itemId, &name, &createdAt)
			if err != nil {
				errorhandler.LogMsg("createItems: rows.Scan")
				return nil, errorhandler.NewInternalError(err.Error())
			}
			items = append(
				items,
				&api.Item{
					Id:        strconv.Itoa(itemId),
					Name:      name,
					UserId:    userId,
					CreatedAt: timestamppb.New(createdAt)})
		}
		err = setUserItem(tx, userId, items)
		if err != nil {
			errorhandler.LogMsg("createItems: setUserItem")
			return nil, errorhandler.NewInternalError(err.Error())
		}
		return items, nil
	} else {
		return nil, nil
	}
}

func setUserItem(tx *sql.Tx, userId string, items []*api.Item) error {
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items))
	i := 0
	for _, item := range items {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, userId)
		valueArgs = append(valueArgs, item.Id)
		i++
	}
	query := fmt.Sprintf(InsertUserItemQuery, strings.Join(valueStrings, ","))
	stmt, err := tx.Prepare(query)
	defer stmt.Close()
	_, err = stmt.Exec(valueArgs...)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("setUserItem: stmt.Exec(%v)", valueArgs))
		return err
	}
	return nil
}

func UpdateUser(data *api.UpdateUserRequest) (*api.User, error) {
	tx, err := StorageInstance.DB.Begin()
	if err != nil {
		errorhandler.LogMsg("UpdateUser: StorageInstance.DB.Begin")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	defer tx.Rollback()

	err = updateUser(tx, data)
	if err != nil {
		errorhandler.LogMsg("UpdateUser: updateUser")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	err = updateUserType(tx, data.GetId(), data.GetUserType())
	if err != nil {
		errorhandler.LogMsg("UpdateUser: updateUserType")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	err = updateItems(tx, data.GetItems())
	if err != nil {
		errorhandler.LogMsg("UpdateUser: createItems")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	if err := tx.Commit(); err != nil {
		errorhandler.LogMsg("UpdateUser: tx.Commit")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	user, err := getUserById(data.GetId())
	if err != nil {
		errorhandler.LogMsg("UpdateUser: getUserById")
		return nil, errorhandler.NewInternalError(err.Error())
	}

	return user, nil
}

func updateUser(tx *sql.Tx, data *api.UpdateUserRequest) error {
	stmt, err := tx.Prepare(UpdateUserQuery)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("updateUser: tx.Prepare(%s)", UpdateUserTypeQuery))
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.GetId(), data.GetName(), data.GetAge(), time.Now())
	if err != nil {
		errorhandler.LogMsg("updateUser: row.Scan")
		return err
	}
	return nil
}

func updateUserType(tx *sql.Tx, userId string, newTypeId api.UserType) error {
	stmt, err := tx.Prepare(UpdateUserTypeQuery)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("updateItems: tx.Prepare(%s)", UpdateUserTypeQuery))
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, newTypeId)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("setUserItem: stmt.Exec(%s, %s)", userId, newTypeId))
		return err
	}
	return nil
}

func updateItems(tx *sql.Tx, data []*api.UpdateItemRequest) error {
	if len(data) > 0 {
		var items = make([]*api.Item, 0, len(data))
		valueArgs := make([]interface{}, 0, len(items))
		query := ""
		i := 0
		for _, item := range data {
			query += fmt.Sprintf(UpdateItemQuery, "$"+strconv.Itoa(i*2+2), "$"+strconv.Itoa(i*2+3), "$"+strconv.Itoa(i*2+1))
			valueArgs = append(valueArgs, item.Id, item.Name, time.Now())
			i++
		}
		stmt, err := tx.Prepare(query)
		if err != nil {
			errorhandler.LogMsg(fmt.Sprintf("updateItems: tx.Prepare(%s)", query))
			return err
		}

		defer stmt.Close()
		_, err = stmt.Exec(valueArgs...)
		if err != nil {
			errorhandler.LogMsg(fmt.Sprintf("updateItems: stmt.Query(%v)", valueArgs))
			return err
		}
	}
	return nil
}

func DeleteUser(data *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	tx, err := StorageInstance.DB.Begin()
	if err != nil {
		errorhandler.LogMsg("UpdateUser: StorageInstance.DB.Begin")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	defer tx.Rollback()

	if err := deleteItem(tx, data.Id); err != nil {
		errorhandler.LogMsg(fmt.Sprintf("DeleteUser: deleteItem(tx, %s)", data.Id))
		return nil, err
	}
	if err := deleteUser(tx, data.Id); err != nil {
		errorhandler.LogMsg(fmt.Sprintf("DeleteUser: deleteItem(tx, %s)", data.Id))
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		errorhandler.LogMsg("DeleteUser: tx.Commit")
		return nil, errorhandler.NewInternalError(err.Error())
	}
	return &api.DeleteUserResponse{}, nil
}

func deleteUser(tx *sql.Tx, userId string) error {
	stmt, err := tx.Prepare(DeleteUserQuery)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("deleteUser: tx.Prepare(%s)", DeleteItemQuery))
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("deleteUser: stmtUser.Exec(%s)", userId))
		return err
	}
	return nil
}

func deleteItem(tx *sql.Tx, userId string) error {
	stmt, err := tx.Prepare(DeleteItemQuery)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("deleteItem: tx.Prepare(%s)", DeleteItemQuery))
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("deleteItem: stmtItem.Exec(%s)", userId))
		return err
	}
	return nil
}

func ListUser(data *api.ListUserRequest) (*api.ListUserResponse, error) {
	rows, err := StorageInstance.DB.Query(SelectUsersQuery,
		data.GetPageFilter().GetLimit(),
		data.GetPageFilter().GetLimit()*(data.GetPageFilter().GetPage()-1))
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("ListUser: StorageInstance.DB.Query(%v, %s)", SelectUserQuery, data.GetPageFilter()))
		return nil, errorhandler.NewInternalError(err.Error())
	}
	defer rows.Close()
	users, err := retrieveUsers(rows)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("ListUser: retrieveUsers(rows), error = %v", err))
		return nil, errorhandler.NewInternalError(err.Error())
	}
	return &api.ListUserResponse{Users: users}, nil
}

func GetUser(data *api.GetUserRequest) (*api.User, error) {
	return getUserById(data.GetId())
}

func getUserById(userId string) (*api.User, error) {
	var user *api.User = nil
	rows, err := StorageInstance.DB.Query(SelectUserQuery, userId)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("GetUser: StorageInstance.DB.Query(%v, %s)", SelectUserQuery, userId))
		return nil, errorhandler.NewInternalError(err.Error())
	}
	users, err := retrieveUsers(rows)
	if err != nil {
		errorhandler.LogMsg(fmt.Sprintf("GetUser: retrieveUsers(rows), error = %v", err))
		return nil, errorhandler.NewInternalError(err.Error())
	}
	if len(users) == 1 {
		user = users[0]
		err = nil
	} else if len(users) == 0 {
		msg := fmt.Sprintf("GetUser: User not found by id = %s", userId)
		errorhandler.LogMsg(msg)
		err = errorhandler.NewNotFoundError(msg)
	} else if len(users) > 1 {
		msg := fmt.Sprintf("GetUser: There are more than 1 user by id %s", userId)
		errorhandler.LogMsg(msg)
		err = errorhandler.NewInternalError(msg)
	}
	return user, err
}

func retrieveUsers(rows *sql.Rows) ([]*api.User, error) {
	var userIdToUser = make(map[string]*api.User, 0)
	defer rows.Close()
	for rows.Next() {
		var (
			userId        string
			userName      string
			userAge       int32
			userType      api.UserType
			userCreatedAt time.Time
			userUpdatedAt pq.NullTime
			itemId        sql.NullString
			itemName      sql.NullString
			itemCreatedAt pq.NullTime
			itemUpdatedAt pq.NullTime
		)
		if err := rows.Scan(&userId, &userName, &userAge, &userType, &userCreatedAt, &userUpdatedAt, &itemId, &itemName, &itemCreatedAt, &itemUpdatedAt); err != nil {
			errorhandler.LogMsg(fmt.Sprintf("GetUser: StorageInstance.DB.Query(%v, %s)", SelectUserQuery, userId))
			return nil, errorhandler.NewInternalError(err.Error())
		}

		if userIdToUser[userId] == nil {
			user := &api.User{
				Id:        userId,
				Name:      userName,
				Age:       userAge,
				UserType:  userType,
				CreatedAt: timestamppb.New(userCreatedAt),
				UpdatedAt: getTimestamp(userUpdatedAt)}
			userIdToUser[userId] = user
		}
		if userIdToUser[userId] != nil && itemId.Valid && itemName.Valid {
			item := &api.Item{
				Id:        itemId.String,
				Name:      itemName.String,
				UserId:    userId,
				CreatedAt: getTimestamp(itemCreatedAt),
				UpdatedAt: getTimestamp(itemUpdatedAt)}
			userIdToUser[userId].Items = append(userIdToUser[userId].Items, item)
		}
	}

	var users = make([]*api.User, 0, len(userIdToUser))
	for _, value := range userIdToUser {
		users = append(users, value)
	}

	return users, nil
}

func getTimestamp(value pq.NullTime) *timestamppb.Timestamp {
	if value.Valid {
		return timestamppb.New(value.Time)
	}
	return nil
}
