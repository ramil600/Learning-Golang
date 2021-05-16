// This code snippet demonstates the use of decorator(middleware) pattern we use
// to check whether the user is logged in, using session and local storage of users.
// middleware function takes handler and returns same handler if the user is logged in.
package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	uName string
	pwd   string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/new.gohtml"))
}

//imitates local storage of users in the system
var dbUsers = map[string]user{}

//stores sessions that corresponds to valid users
var sessions = map[string]string{}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		myUuid, _ := uuid.NewV4()

		cookie := &http.Cookie{
			Name:     "session",
			Value:    myUuid.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		uName := r.FormValue("username")
		pwd := r.FormValue("password")

		usr := user{
			uName: uName,
			pwd:   pwd,
		}
		dbUsers[uName] = usr

		sessions[myUuid.String()] = uName
		fmt.Println("Cookie:", cookie)
		fmt.Println("User: ", dbUsers[uName])
		http.Redirect(w, r, "/admin", http.StatusSeeOther)

	}

	tpl.Execute(w, nil)
}

// middleware returns handler by using its method ServeHTTP, which will call parent function
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Not logged in", http.StatusForbidden)
			return
		}

		if un, ok := sessions[cookie.Value]; ok {
			if _, val := dbUsers[un]; val {
				next.ServeHTTP(w, r)
			}

		} else {
			http.Error(w, "Not Authorized!", http.StatusForbidden)
		}
	})
}

func admin(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	uid := cookie.Value
	uName := sessions[uid]

	io.WriteString(w, "Hello, ")
	fmt.Fprint(w, uName)
}

func main() {

	http.HandleFunc("/", register)

	// HandlerFunc casts admin function to http.Handler, which is expected by middleware
	http.Handle("/admin", middleware(http.HandlerFunc(admin)))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
