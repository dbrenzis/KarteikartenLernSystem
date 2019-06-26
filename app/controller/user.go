package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Webprojekt/app/model"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore

func init() {
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewCookieStore(key)
}

//Register Für die Register Seite
func Register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/base.tmpl", "template/register.tmpl")
	log.Println(err)
	register := model.RegisterDat{}
	register.BaseDaten = model.GetBaseData("", "")

	t.Execute(w, register)
}

//DelUser controller
func DelUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)

	user, _ := model.GetUserByUsername(username)

	err := model.DelKasten(username)
	err = user.Del()
	if err == nil {

		http.Redirect(w, r, "/logout", http.StatusFound)
	}
}

//UpdateImage controller
func UpdateImage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)

	user, _ := model.GetUserByUsername(username)

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	filebyte, _ := ioutil.ReadAll(file)

	encodedString := base64.StdEncoding.EncodeToString(filebyte)

	user.Bild = encodedString

	user.Update()

	http.Redirect(w, r, "/", http.StatusFound)

}

//UpdateUser controller
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	oldPassword := r.FormValue("oldPassword")
	password := r.FormValue("Password")

	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)

	user, _ := model.GetUserByUsername(username)
	// decode base64 String to []byte
	passwordDB, _ := base64.StdEncoding.DecodeString(user.Passwort)
	err := bcrypt.CompareHashAndPassword(passwordDB, []byte(oldPassword))
	if err == nil {

		if model.Checkmail(email) {
			if email != "" {
				user.Email = email
			}
			user.Passwort = password
			err = user.Update()
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			t, _ := template.ParseFiles("template/base.tmpl", "template/profil.tmpl")
			profildata, _ := model.GetProfilData(username)
			profildata.BaseDaten.LoggedIn = loggedIn
			profildata.ErrorEmail = "Mail existiert bereits"
			t.Execute(w, profildata)
		}
	} else {
		t, _ := template.ParseFiles("template/base.tmpl", "template/profil.tmpl")
		profildata, _ := model.GetProfilData(username)
		profildata.BaseDaten.LoggedIn = loggedIn
		profildata.ErrorPW = "Falsches Passwort"
		t.Execute(w, profildata)
	}

}

// AddUser controller
func AddUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	passwort := r.FormValue("password")
	email := r.FormValue("email")

	user := model.UserData{}
	user.Name = name
	user.Passwort = passwort
	user.Email = email

	err := user.Add()
	fmt.Println(err)

	if err == nil {
		session, _ := store.Get(r, "session")

		// Set user as authenticated
		session.Values["loggedIn"] = true
		session.Values["name"] = user.Name
		session.Save(r, w)
		http.Redirect(w, r, "/meinekarteien", http.StatusFound)
	} else {

		var register = model.RegisterDat{}
		register.BaseDaten = model.GetBaseData("", "")
		t, _ := template.ParseFiles("template/base.tmpl", "template/register.tmpl")

		if strings.Contains(err.Error(), "Name existiert") {
			register.ErrorUsername = "Dieser Nutzername ist bereits vergeben!"
		}

		if strings.Contains(err.Error(), "Email exisitert") {
			register.ErrorEmail = "Diese Mail ist bereist vergeben!"
		}

		if strings.Contains(err.Error(), "Kein Benutzername") {
			register.ErrorUsername = "Bitte einen gültigen Benutzernamen eingeben"
		}

		if strings.Contains(err.Error(), "Keine Email") {
			register.ErrorEmail = "Bitte eine gültige E-Mail eingeben"
		}

		t.Execute(w, register)
	}
}

// AuthenticateUser controller
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var index = model.IndexDat{}
	index, _ = model.GetIndexData("")
	index.ErrorMsg = "Benutzername und/oder Passwort ist falsch!"

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authentication
	user, err := model.GetUserByUsername(username)
	if err == nil {
		// decode base64 String to []byte
		passwordDB, _ := base64.StdEncoding.DecodeString(user.Passwort)
		err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
		if err == nil {
			session, _ := store.Get(r, "session")

			// Set user as authenticated
			session.Values["loggedIn"] = true
			session.Values["name"] = username
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			t, _ := template.ParseFiles("template/base.tmpl", "template/index.tmpl")
			t.Execute(w, index)
		}
	} else {
		t, _ := template.ParseFiles("template/base.tmpl", "template/index.tmpl")
		t.Execute(w, index)
	}
}

// AuthenticateUserKartei controller Karteikasten
func AuthenticateUserKartei(w http.ResponseWriter, r *http.Request) {
	var err error
	var user = model.UserData{}
	var kartei = model.KarteikastenDat{}
	kartei, _ = model.GetKarteikastenData("")
	kartei.ErrorMsg = "Benutzername und/oder Passwort ist falsch!"

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authentication
	user, err = model.GetUserByUsername(username)
	if err == nil {
		// decode base64 String to []byte
		passwordDB, _ := base64.StdEncoding.DecodeString(user.Passwort)
		err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))
		if err == nil {
			session, _ := store.Get(r, "session")

			// Set user as authenticated
			session.Values["loggedIn"] = true
			session.Values["name"] = username
			session.Save(r, w)
			http.Redirect(w, r, "/karteikasten", http.StatusFound)
		} else {
			t, _ := template.ParseFiles("template/base.tmpl", "template/karteikasten.tmpl")
			t.Execute(w, kartei)
		}
	} else {
		t, _ := template.ParseFiles("template/base.tmpl", "template/karteikasten.tmpl")
		t.Execute(w, kartei)
	}
}

// Logout controller
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["loggedIn"] = false
	session.Values["name"] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

//AuthIndex is an authentication handler
func AuthIndex(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated
		if auth, ok := session.Values["loggedIn"].(bool); !ok || !auth {
			t, _ := template.ParseFiles("template/base.tmpl", "template/index.tmpl")
			index, _ := model.GetIndexData("")
			t.Execute(w, index)
		} else {
			h(w, r)
		}
	}
}

//AuthKartei is an authentication handler
func AuthKartei(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated
		if auth, ok := session.Values["loggedIn"].(bool); !ok || !auth {
			t, _ := template.ParseFiles("template/base.tmpl", "template/karteikasten.tmpl")
			karteien, _ := model.GetKarteikastenData("")
			t.Execute(w, karteien)
		} else {
			h(w, r)
		}
	}
}

//AuthView is an authentication handler
func AuthView(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		id := r.FormValue("id")
		nr := r.FormValue("nr")

		// Check if user is authenticated
		if auth, ok := session.Values["loggedIn"].(bool); !ok || !auth {
			t, _ := template.ParseFiles("template/base.tmpl", "template/view.tmpl")
			karteien, _ := model.GetViewData("", id, nr)
			t.Execute(w, karteien)
		} else {
			h(w, r)
		}
	}
}
