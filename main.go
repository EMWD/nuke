package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

var (
	templates *template.Template
	db        *sql.DB
)

func sha256d(data []byte) []byte {
	result := sha256.Sum256(data)
	result = sha256.Sum256(result[:])
	return result[:]
}

func main() {
	var err error

	connStr := "host=127.0.0.1 port=5432 user=undefined password=killer777 dbname=representationdb"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"add":    func(a, b int) int { return a + b },
		"decode": func(enc string) string { data, _ := base64.StdEncoding.DecodeString(enc); return string(data) },
	}).ParseGlob("html/*.html"))

	// HTTP
	http.HandleFunc("/", getHandler("index"))
	http.HandleFunc("/signup", getHandler("signup"))
	http.HandleFunc("/signin", getHandler("signin"))
	http.HandleFunc("/donate", getHandler("donate"))
	http.HandleFunc("/showpresentation", getHandler("showPresentation"))

	http.HandleFunc("/editor", getAuthHandler(editorHandler))
	http.HandleFunc("/profile", getAuthHandler(profileHandler))
	http.HandleFunc("/review", getAuthHandler(reviewHandler))
	http.HandleFunc("/exit", getAuthHandler(exitHandler))

	// REST
	http.HandleFunc("/signup_rest", postHandler(signupHandler))
	http.HandleFunc("/signin_rest", postHandler(signinHandler))
	http.HandleFunc("/add_presentation_rest", postAuthHandler(addPresentationHandler))
	http.HandleFunc("/get_presentation_rest", postHandler(getPresentationHandler))
	http.HandleFunc("/del_presentation_rest", postAuthHandler(delPresentationHandler))
	http.HandleFunc("/update_presentation_rest", postAuthHandler(updatePresentationHandler))
	http.HandleFunc("/change_password_rest", postAuthHandler(changePasswordHandler))
	http.HandleFunc("/delete_account_rest", postAuthHandler(deleteAccountHandler))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.ListenAndServe(":9000", nil)
}
