package services

import (
	"context"
	"fmt"
	"strings"

	"pfserver/db"
	"pfserver/utils"

	"github.com/jackc/pgx/v4"
)

type U struct{}

func User() *U {
	return &U{}
}

type UserData struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// get user data by id or phone
type GetUserBy struct {
	Id    int64
	Email string
}

func (*U) GetUserData(ctx context.Context, getBy GetUserBy) (UserData, error) {

	var col interface{} // default
	var colData interface{}
	if getBy.Id != 0 {
		col = "id"
		colData = getBy.Id
	} else if getBy.Email != "" {
		col = "email"
		colData = getBy.Email
	}

	var data UserData

	// else get the data from the DB and save it in redis :)
	err := db.Conn().QueryRow(
		ctx,
		fmt.Sprintf("SELECT id, firstname, lastname, email, password FROM users WHERE %s = $1", col),
		colData,
	).Scan(&data.Id, &data.Firstname, &data.Lastname, &data.Email, &data.Password)

	return data, err
}

// check if user exists by his phone
func (*U) IsUserEmailExist(ctx context.Context, email string) (bool, error) {
	cmd, err := db.Conn().Exec(ctx, "SELECT id FROM users WHERE email = $1", email)
	if err == pgx.ErrNoRows {
		return false, err
	}

	return cmd.RowsAffected() > 0, err
}

// insert a new user
type CreateAccountData struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (*U) CreateAccount(ctx context.Context, data CreateAccountData) (int64, error) {
	var newUserId int64
	err := db.Conn().QueryRow(
		ctx,
		`INSERT INTO users(firstname, lastname, email, password)
		 VALUES($1, $2, $3, $4) RETURNING(id)`,
		data.FirstName, data.LastName, data.Email, data.Password,
	).Scan(&newUserId)

	return newUserId, err
}

// the user data that can be updated in the DB
//? pointers are used to make the default value "nil",
//? to check later if the value is set or not, and update only the set fields
type UserDataToUpdate struct {
	Firstname string `db_col:"firstname"`
	Lastname  string `db_col:"lastname"`
	Email     string `db_col:"email"`
	Password  string `db_col:"password"`
}

func (*U) UpdateUser(ctx context.Context, userId int64, data UserDataToUpdate) error {

	reflektUserDataToUpdate := utils.Reflekt(&data)
	var fieldsValue []interface{}
	var fieldsToUpdate string
	var counter int = 1
	for i := 0; i < reflektUserDataToUpdate.Elem.NumField(); i++ {
		field := reflektUserDataToUpdate.GetFieldDetails(i, "db_col")
		isFieldEmpty := reflektUserDataToUpdate.IsFieldEmpty(field.Value)

		if !isFieldEmpty {
			fieldsToUpdate += fmt.Sprintf("%s = $%d,", field.Tags["db_col"], counter)
			fieldsValue = append(fieldsValue, field.Value.Interface())
			counter++
		}
	}

	// remove that last ",", coz it causes an error
	fieldsToUpdate = strings.TrimSuffix(fieldsToUpdate, ",")

	_, err := db.Conn().Exec(
		ctx,
		fmt.Sprintf("UPDATE users SET %s WHERE id = %d", fieldsToUpdate, userId),
		fieldsValue...,
	)

	return err
}
