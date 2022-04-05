package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"pfserver/config"
	"pfserver/services"
	"pfserver/utils"
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
