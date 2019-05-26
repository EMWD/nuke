package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const authCookieName = "auth"

type AuthCookie struct {
	Login    string
	Password []byte
}

func getHandler(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		login, _, _ := authHandler(w, req)
		templates.ExecuteTemplate(w, name, login)
	}
}

func getAuthHandler(handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		login, isAuth, err := authHandler(w, req)
		if err != nil {
			templates.ExecuteTemplate(w, "auth_required", "внутренняя ошибка сервера")
			return
		}
		if !isAuth {
			templates.ExecuteTemplate(w, "auth_required", "для данной функции требуется авторизация")
			return
		}
		handler(w, req, login)
	}
}

func postHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(w, req)
			return
		}
		handler(w, req)
	}
}

func postAuthHandler(handler func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(w, req)
			return
		}

		login, isAuth, err := authHandler(w, req)
		if err != nil {
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
		if !isAuth {
			http.Error(w, "ошибка авторизации", http.StatusBadRequest)
			return
		}

		handler(w, req, login)
	}
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Expires: time.Unix(0, 0),
	})
}

func setAuthCookie(w http.ResponseWriter, login string, password []byte) {
	auth, _ := json.Marshal(AuthCookie{
		Login:    login,
		Password: password,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   base64.StdEncoding.EncodeToString(auth),
		Expires: time.Now().AddDate(0, 1, 0),
	})
}

func getAuthCookie(w http.ResponseWriter, req *http.Request) *AuthCookie {
	cookies := req.Cookies()
	authStr := ""

	for _, c := range cookies {
		if c.Name == authCookieName {
			authStr = c.Value
			break
		}
	}
	if authStr == "" {
		return nil
	}

	var auth *AuthCookie
	authBytes, err := base64.StdEncoding.DecodeString(authStr)
	if err != nil {
		deleteCookie(w, authCookieName)
		return nil
	}
	err = json.Unmarshal(authBytes, &auth)
	if err != nil {
		deleteCookie(w, authCookieName)
		return nil
	}

	return auth
}

func authHandler(w http.ResponseWriter, req *http.Request) (string, bool, error) {
	auth := getAuthCookie(w, req)
	if auth == nil {
		return "", false, nil
	}

	isOk, err := checkAuth(db, auth.Login, hex.EncodeToString(auth.Password))
	return auth.Login, isOk, err
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	login := strings.TrimSpace(req.FormValue("login"))
	password := req.FormValue("password")

	if len(login) < 5 {
		http.Error(w, "логин должен быть не короче 5 символов", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "пароль должен быть не короче 8 символов", http.StatusBadRequest)
		return
	}

	hash := sha256d([]byte(password))
	err := insertUser(db, login, hex.EncodeToString(hash))
	if err != nil {
		http.Error(w, "данный логин уже занят", http.StatusBadRequest)
		return
	}

	setAuthCookie(w, login, hash)
}

func signinHandler(w http.ResponseWriter, req *http.Request) {
	login := strings.TrimSpace(req.FormValue("login"))
	password := req.FormValue("password")

	hash := sha256d([]byte(password))
	isAuth, err := checkAuth(db, login, hex.EncodeToString(hash))
	if err != nil {
		http.Error(w, "внутрення ошибка сервера", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		http.Error(w, "неверный логин или пароль", http.StatusBadRequest)
		return
	}

	setAuthCookie(w, login, hash)
}

func editorHandler(w http.ResponseWriter, req *http.Request, login string) {
	styles, err := getStyles(db)
	if err != nil {
		return
	}

	templates.ExecuteTemplate(w, "editor", struct {
		Login  string
		Styles []string
	}{
		Login:  login,
		Styles: styles,
	})
}

func addPresentationHandler(w http.ResponseWriter, req *http.Request, login string) {
	title := strings.TrimSpace(req.FormValue("title"))
	style := strings.TrimSpace(req.FormValue("style"))
	code := strings.TrimSpace(req.FormValue("code"))

	// validation
	if len(title) == 0 {
		http.Error(w, "название - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(style) == 0 {
		http.Error(w, "стиль - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(code) == 0 {
		http.Error(w, "код презентации - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(code) >= 65537 {
		http.Error(w, "код презентации должен быть меньше 65537 символов", http.StatusBadRequest)
		return
	}

	// create file
	err := os.MkdirAll(login, 0777)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	fileName := filepath.Join(login, title)
	err = ioutil.WriteFile(fileName, []byte(code), 0644)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = insertPrsentation(db, login, style, title, fileName)
	if err != nil {
		http.Error(w, "презентация с таким именем уже существует", http.StatusBadRequest)
	}
}

func profileHandler(w http.ResponseWriter, req *http.Request, login string) {
	presentations, err := getPresentations(db, login)
	if err != nil {
		return
	}

	templates.ExecuteTemplate(w, "profile", struct {
		Login         string
		Presentations []string
	}{
		Login:         login,
		Presentations: presentations,
	})
}

func reviewHandler(w http.ResponseWriter, req *http.Request, login string) {
	templates.ExecuteTemplate(w, "review", login)
}

func getPresentationHandler(w http.ResponseWriter, req *http.Request) {
	user := strings.TrimSpace(req.FormValue("user"))
	title := strings.TrimSpace(req.FormValue("title"))

	if len(user) == 0 {
		http.NotFound(w, req)
		return
	}

	if len(title) == 0 {
		http.NotFound(w, req)
		return
	}

	filepath, style, err := getPresentation(db, user, title)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	buf, _ = json.Marshal(struct {
		Style        string
		Presentation string
	}{
		Style:        style,
		Presentation: string(buf),
	})

	w.Write(buf)
}

func delPresentationHandler(w http.ResponseWriter, req *http.Request, login string) {
	title := strings.TrimSpace(req.FormValue("title"))

	if len(title) == 0 {
		http.Error(w, "название - обязательный аргумент", http.StatusBadRequest)
		return
	}

	filepath, _, err := getPresentation(db, login, title)
	if err != nil {
		http.Error(w, "презентации с таким название не существует", http.StatusBadRequest)
		return
	}

	err = deletePresentation(db, login, title)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Remove(filepath)
}

func updatePresentationHandler(w http.ResponseWriter, req *http.Request, login string) {
	title := strings.TrimSpace(req.FormValue("title"))
	new_title := strings.TrimSpace(req.FormValue("new_title"))
	new_style := strings.TrimSpace(req.FormValue("new_style"))
	new_code := strings.TrimSpace(req.FormValue("new_code"))

	// validation
	if len(title) == 0 {
		http.Error(w, "текущее название - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(new_title) == 0 {
		http.Error(w, "новое название - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(new_style) == 0 {
		http.Error(w, "новый стиль - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(new_code) == 0 {
		http.Error(w, "новый код презентации - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(new_code) >= 65537 {
		http.Error(w, "новый код презентации должен быть меньше 65537 символов", http.StatusBadRequest)
		return
	}

	// update file
	new_filepath := filepath.Join(login, new_title)

	filepath, _, err := getPresentation(db, login, title)
	if err != nil {
		http.Error(w, "презентации с таким названием не существует", http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile(filepath, []byte(new_code), 0644)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = os.Rename(filepath, new_filepath)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = updatePresentation(db, login, title, new_title, new_style, new_filepath)
	if err != nil {
		http.Error(w, "презентация с таким именем уже существует", http.StatusBadRequest)
	}
}

func changePasswordHandler(w http.ResponseWriter, req *http.Request, login string) {
	old_password := req.FormValue("old_password")
	new_password := req.FormValue("new_password")

	if len(old_password) == 0 {
		http.Error(w, "текущий пароль - обязательный аргумент", http.StatusBadRequest)
		return
	}

	if len(new_password) == 0 {
		http.Error(w, "новый пароль - обязательный аргумент", http.StatusBadRequest)
		return
	}

	isAuth, err := checkAuth(db, login, hex.EncodeToString(sha256d([]byte(old_password))))
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		http.Error(w, "неверный пароль", http.StatusBadRequest)
		return
	}

	err = changePassword(db, login, hex.EncodeToString(sha256d([]byte(new_password))))
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}

	deleteCookie(w, "auth")
}

func exitHandler(w http.ResponseWriter, req *http.Request, login string) {
	deleteCookie(w, "auth")
	http.Redirect(w, req, "/index", http.StatusFound)
}

func deleteAccountHandler(w http.ResponseWriter, req *http.Request, login string) {
	password := req.FormValue("password")

	if len(password) == 0 {
		http.Error(w, "текущий пароль - обязательный аргумент", http.StatusBadRequest)
		return
	}

	isAuth, err := checkAuth(db, login, hex.EncodeToString(sha256d([]byte(password))))
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		http.Error(w, "неверный пароль", http.StatusBadRequest)
		return
	}

	err = deleteUser(db, login)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	os.RemoveAll(login)
	deleteCookie(w, "auth")
}
