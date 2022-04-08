package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"pfserver/config"
	"pfserver/services"
	"pfserver/utils"

	"golang.org/x/crypto/bcrypt"
)

type U struct{}

func Users() U {
	return U{}
}

func (U) MembersList(res http.ResponseWriter, req *http.Request) {
	homeTmpl, _ := template.ParseFiles(utils.TemplatePath("home.html"))
	loggedInUser, _ := config.NewSession(req, res).GetUser()
	members, err := services.User().GetAllMembers(
		req.Context(),
		"",
		int(loggedInUser["id"].(int64)),
	)

	if err != nil {
		homeTmpl.Execute(res, nil)
		fmt.Println("Members List Error: ", err)
		return
	}

	type Result struct {
		Members []*services.UserDetails
	}
	homeTmpl.Execute(res, Result{
		Members: members,
	})

}

func (U) Settings(res http.ResponseWriter, req *http.Request) {
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	password := req.FormValue("password")

	loggedInUser, _ := config.NewSession(req, res).GetUser()
	loggedInUserId := loggedInUser["id"].(int64)

	newPassword := ""

	if password != "" {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
		newPassword = string(hashedPass)
	}

	err := services.User().UpdateUser(req.Context(), loggedInUserId, services.UserDataToUpdate{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  newPassword,
	})

	isSuccess := "true"
	if err != nil {
		isSuccess = "false"
	}
	http.Redirect(res, req, fmt.Sprintf("/settings?success=%s", isSuccess), http.StatusSeeOther)
}
