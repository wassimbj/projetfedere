package handlers

import (
	"net/http"

	"pfserver/config"
	"pfserver/core"

	"github.com/jackc/pgx/v4"

	"pfserver/services"

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
}

func (a *UserAuth) Signup(res http.ResponseWriter, req *http.Request) {

	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	password := req.FormValue("password")

	// check if the phone is already registered
	_, existingUserErr := services.User().GetUserData(req.Context(), services.GetUserBy{
		Email: email,
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
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	// var newUserId int64

	// auto created = created while creating a new transaction
	// we will just update the neccessary details
	_, createErr := services.User().CreateAccount(
		req.Context(),
		services.CreateAccountData{
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
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
	// newUserId = userId

	// log the user in (phone is not confirmed yet and identity too)
	// config.NewSession(req, res).Save("user", config.SessData{
	// 	"id":    newUserId,
	// 	"email": email,
	// })

	http.Redirect(res, req, "/login", http.StatusSeeOther)

}

type LoginData struct {
	Password string `json:"password" validate:"required,max=20,min=5"`
	Email    string `json:"email" validate:"required,max=250,min=12"`
}

func (a *UserAuth) Login(res http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	password := req.FormValue("password")

	// validate the data
	// validationErr := utils.Validator().Struct(loginInfo)
	// if validationErr != nil {
	// 	core.Respond(res, core.ResOpts{
	// 		Status: http.StatusBadRequest,
	// 		Msg:    validationErr.Error(),
	// 	})
	// 	return
	// }

	user, err := services.User().GetUserData(req.Context(), services.GetUserBy{
		Email: email,
	})

	// err may happen if no rows are found or something else
	if err != nil {
		core.Respond(res, core.ResOpts{
			Status: http.StatusUnauthorized,
			Msg:    "invalid account details",
		})
		return
	}

	matchFailed := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func (a *UserAuth) Logout(res http.ResponseWriter, req *http.Request) {

	config.NewSession(req, res).Del("user")

	http.Redirect(res, req, "/", http.StatusSeeOther)
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
