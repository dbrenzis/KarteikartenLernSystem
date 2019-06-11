package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//User Data
type UserData struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Passwort string `json:"passwort"`
	Bild     string `json:"bild"`
	Erstellt string `json:"erstellt"`
}

//Checkmail Prüft, ob die Mail schon vorhanden ist
func Checkmail(email string) bool {
	query := `
	{
		"selector": {
			 "type": "Benutzer",
			 "email": "%s"
		}
	}`
	u, err := btDBS.QueryJSON(fmt.Sprintf(query, email))
	if err == nil && len(u) == 0 {
		return true
	} else {
		return false
	}
}

//Del Löscht den User
func (user UserData) Del() (err error) {
	err = btDBS.Delete(user.ID)
	return err
}

//Update Erneuert die Benutzer Daten
func (user UserData) Update() (err error) {

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Passwort), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

	user.Passwort = b64HashedPwd

	// Convert Todo struct to map[string]interface as required by Save() method
	u, err := user2Map(user)

	//delete(u, "_rev")

	err = btDBS.Set(user.ID, u)

	return err

}

// Add User
func (user UserData) Add() (err error) {
	// Check wether username already exists
	userInDB, err := GetUserByUsernameAndMail(user.Name, user.Email)
	if err == nil {
		if userInDB.Name == user.Name {
			errors.New("Name existiert")

		}
		if userInDB.Email == user.Email {
			err = errors.New("Email exisitert")
		}
	}

	if err != nil {
		return err
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Passwort), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

	user.Passwort = b64HashedPwd
	user.Type = "Benutzer"

	t := time.Now()
	user.Erstellt = t.Format("2006-01-02 15:04:05")

	// Convert Todo struct to map[string]interface as required by Save() method
	u, err := user2Map(user)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(u, "_id")
	delete(u, "_rev")

	// Add todo to DB
	_, _, err = btDBS.Save(u, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}

	return err
}

// GetUserByUsernameAndMail check
func GetUserByUsernameAndMail(username string, email string) (user UserData, err error) {
	if username == "" {
		errors.New("Kein Benutzername")
	}
	if email == "" {
		errors.New("Keine Email")
	}

	if err != nil {
		return UserData{}, err
	}

	query := `
	{
		"selector": {
			 "type": "Benutzer",
			 "name": "%s"
		}
	}`
	u, err := btDBS.QueryJSON(fmt.Sprintf(query, username))
	if err != nil || len(u) != 1 {
		return UserData{}, err
	}

	user, err = map2User(u[0])
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}

//Get User By Username
func GetUserByUsername(username string) (user UserData, err error) {
	if username == "" {
		return UserData{}, errors.New("Kein Benutzername")
	}

	query := `
	{
		"selector": {
			 "type": "Benutzer",
			 "name": "%s"
		}
	}`
	u, err := btDBS.QueryJSON(fmt.Sprintf(query, username))
	if err != nil || len(u) != 1 {
		return UserData{}, err
	}

	user, err = map2User(u[0])
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}

// ---------------------------------------------------------------------------
// Internal helper functions
// ---------------------------------------------------------------------------

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func user2Map(u UserData) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2User(user map[string]interface{}) (u UserData, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
