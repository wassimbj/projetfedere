package handlers

import (
	"encoding/json"
	"net/http"

	"pfserver/config"
	"pfserver/core"

	"github.com/jackc/pgx/v4"

	"pfserver/services"
	"pfserver/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserAuth struct{}

func Auth() *UserAuth {
	return &UserAuth{}
}

type SignupData struct {
	Firstname string `json:"firstname" validate:"required,max=100,min=3"`
	Lastname  string `json:"lastname" validate:"required,max=100,min=3"`
	Email     string `json:"email" validate:"required,max=250,min=12"`
	Password  string `json:"password" validate:"required,max=20,min=8"`
	PhoneCode string `json:"phoneCode" validate:"required,max=5,min=1"` // e.g: +216
}

func (a *UserAuth) Signup(res http.ResponseWriter, req *http.Request) {

	body, err := core.ReadBody(req.Body)
	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusBadRequest,
			Msg:    "Coudln't read the body",
		})
		return
	}
	var userInfo SignupData
	json.Unmarshal([]byte(body), &userInfo)

	// validate the data
	validationErr := utils.Validator().Struct(userInfo)
	if validationErr != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusBadRequest,
			Msg:    validationErr.Error(),
		})
		return
	}

	// check if the phone is already registered
	_, existingUserErr := services.User().GetUserData(req.Context(), services.GetUserBy{
		Email: userInfo.Email,
	})
	// exist, _ := services.User().IsUserPhoneExist(req.Context(), userInfo.Phone); exist
	if existingUserErr != pgx.ErrNoRows && existingUserErr != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusUnauthorized,
			Msg:    "Invalid details",
		})
		return
	}

	// create user account
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 10)
	var newUserId int64

	// auto created = created while creating a new transaction
	// we will just update the neccessary details
	userId, createErr := services.User().CreateAccount(
		req.Context(),
		services.CreateAccountData{
			FirstName: userInfo.Firstname,
			LastName:  userInfo.Lastname,
			Email:     userInfo.Email,
			Password:  string(hashedPass),
		},
	)
	if createErr != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusInternalServerError,
			Msg:    "something went wrong...",
		})
		return
	}
	newUserId = userId

	// log the user in (phone is not confirmed yet and identity too)
	config.NewSession(req, res).Save("user", config.SessData{
		"id":    newUserId,
		"email": userInfo.Email,
	})

	core.Respond(res, core.ResOpts{
		Status: http.StatusCreated,
		Msg:    "success",
	})

}

type LoginData struct {
	Password string `json:"password" validate:"required,max=20,min=5"`
	Email    string `json:"email" validate:"required,max=250,min=12"`
}

func (a *UserAuth) Login(res http.ResponseWriter, req *http.Request) {

	body, err := core.ReadBody(req.Body)
	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusBadRequest,
			Msg:    "Bad Request",
		})
		return
	}
	var loginInfo LoginData
	json.Unmarshal([]byte(body), &loginInfo)

	// validate the data
	validationErr := utils.Validator().Struct(loginInfo)
	if validationErr != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusBadRequest,
			Msg:    validationErr.Error(),
		})
		return
	}

	user, err := services.User().GetUserData(req.Context(), services.GetUserBy{
		Email: loginInfo.Email,
	})

	// err may happen if no rows are found or something else
	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusUnauthorized,
			Msg:    "invalid account details",
		})
		return
	}

	matchFailed := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if matchFailed != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusUnauthorized,
			Msg:    "invalid account details",
		})
		return
	}

	config.NewSession(req, res).Save("user", config.SessData{
		"id":    user.Id,
		"email": user.Email,
	})

	type LoggedInUserData struct {
		Id    int64  `json:"id"`
		Email string `json:"email"`
	}
	core.Respond(res, core.ResOpts{
		Status: http.StatusOK,
		Msg: LoggedInUserData{
			Id:    user.Id,
			Email: user.Email,
		},
	})
}

func (a *UserAuth) Logout(res http.ResponseWriter, req *http.Request) {

	err := config.NewSession(req, res).Del("user")

	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusInternalServerError,
			Msg:    "something went wrong",
		})
		return
	}

	core.Respond(res, core.ResOpts{
		Status: http.StatusOK,
		Msg:    "Logged out !!",
	})
}

type ConfirmPhoneData struct {
	ConfirmToken string `json:"confirmCode"`
}

type LoggedInUser struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

// return logged in user details
func (a *UserAuth) Me(res http.ResponseWriter, req *http.Request) {

	sess := config.NewSession(req, res)
	userSess, err := sess.Get("user")

	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusInternalServerError,
			Msg:    "Sorry, pleas come back later",
		})
		return
	}

	userId := userSess.Values["id"].(int64)
	Email := userSess.Values["email"].(string)

	user := LoggedInUser{
		Id:    userId,
		Email: Email,
	}

	core.Respond(res, core.ResOpts{
		Status: http.StatusOK,
		Msg:    user,
	})
}
