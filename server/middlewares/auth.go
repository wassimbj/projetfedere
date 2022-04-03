package middlewares

import (
	"fmt"
	"net/http"

	"pfserver/config"

	"pfserver/core"
)

func isUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	data, _ := config.NewSession(req, res).GetUser()

	return data["id"] != nil
}

// block == true ? block not logged in users : block logged in users
func CheckAuth(f http.HandlerFunc, block bool) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loggedInUser, _ := config.NewSession(req, res).Get("user")
		isLoggedIn := loggedInUser.Values["id"] != nil
		fmt.Println(isLoggedIn)
		// login and signup routes blocks logged in users
		// other routes that needs authentication blocks un-auth users
		// if (block && !isLoggedIn) || (!block && isLoggedIn) {
		if block && !isLoggedIn {
			fmt.Println("NOT AUTH 1")
			http.Redirect(res, req, "/login", http.StatusSeeOther)
		} else if !block && isLoggedIn {
			fmt.Println("NOT AUTH 2")
			// fmt.Println("PATH: ", req.URL.Path)
			http.Redirect(res, req, "/", http.StatusSeeOther)
		}

		// run the function
		f(res, req)

	}
}

// check if admin
func CheckAdminAuth(f http.HandlerFunc, block bool) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		s, _ := config.NewSession(req, res).Get("admin")
		isLoggedIn := s.Values["isLoggedIn"] != nil

		// login and signup routes blocks logged in users
		// other routes that needs authentication blocks un-auth users
		if (block && !isLoggedIn) || (!block && isLoggedIn) {
			core.Respond(res, core.ResOpts{
				Status: http.StatusUnauthorized,
				Msg:    "Unauthorized",
			})
			return
		}

		// run the function
		f(res, req)

	}
}
