package main

import (
	"database/sql"
)

// User
func checkAuth(db *sql.DB, login, password string) (bool, error) {
	var rec_pass string
	row := db.QueryRow("SELECT password FROM users WHERE login=$1", login)
	err := row.Scan(&rec_pass)
	if err != nil {
		return false, err
	}
	return (rec_pass == password), nil
}

// Styles
func getStyleId(db *sql.DB, title string) (int, error) {
	var style_id int
	row := db.QueryRow("SELECT id FROM styles WHERE title=$1", title)
	err := row.Scan(&style_id)
	return style_id, err
}

func getStyles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT title FROM styles")
	if err != nil {
		return nil, err
	}

	var style string
	var styles []string
	for rows.Next() {
		err = rows.Scan(&style)
		if err != nil {
			continue
		}
		styles = append(styles, style)
	}

	return styles, nil
}

// User
func getUserId(db *sql.DB, login string) (int, error) {
	var user_id int
	row := db.QueryRow("SELECT id FROM users WHERE login=$1", login)
	err := row.Scan(&user_id)
	return user_id, err
}

func insertUser(db *sql.DB, login, password string) error {
	_, err := db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	return err
}

func deleteUser(db *sql.DB, login string) error {
	user_id, err := getUserId(db, login)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM presentations WHERE user_id=$1", user_id)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE login=$1", login)
	return err
}

func changePassword(db *sql.DB, login, password string) error {
	_, err := db.Exec("UPDATE users SET password=$1 WHERE login=$2", password, login)
	return err
}

// Prsentation
func getPresentations(db *sql.DB, login string) ([]string, error) {
	rows, err := db.Query("SELECT presentations.title FROM presentations INNER JOIN users ON presentations.user_id = users.id AND users.login=$1", login)
	if err != nil {
		return nil, err
	}

	var title string
	var titles []string
	for rows.Next() {
		err = rows.Scan(&title)
		if err != nil {
			continue
		}
		titles = append(titles, title)
	}

	return titles, nil
}

func getPresentation(db *sql.DB, login, title string) (string, string, error) {
	var filepath, style string
	row := db.QueryRow("SELECT presentations.filepath, styles.title FROM presentations INNER JOIN users ON presentations.user_id = users.id AND users.login = $1 AND presentations.title = $2 INNER JOIN styles ON presentations.style_id = styles.id", login, title)
	err := row.Scan(&filepath, &style)
	if err != nil {
		return "", "", err
	}
	return filepath, style, nil
}

func insertPrsentation(db *sql.DB, login, style, title, filepath string) error {
	user_id, err := getUserId(db, login)
	if err != nil {
		return err
	}

	style_id, err := getStyleId(db, style)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO presentations (user_id, style_id, title, filepath) VALUES ($1, $2, $3, $4)", user_id, style_id, title, filepath)
	return nil
}

func deletePresentation(db *sql.DB, login, title string) error {
	_, err := db.Exec("DELETE FROM presentations WHERE exists (SELECT id FROM users WHERE users.id=presentations.user_id and users.login=$1 and presentations.title=$2)", login, title)
	return err
}

func updatePresentation(db *sql.DB, login, title, new_title, new_style, new_filepath string) error {
	user_id, err := getUserId(db, login)
	if err != nil {
		return err
	}

	style_id, err := getStyleId(db, new_style)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE presentations SET style_id=$1, title=$2, filepath=$3 WHERE user_id=$4 AND title=$5", style_id, new_title, new_filepath, user_id, title)
	return err
}
