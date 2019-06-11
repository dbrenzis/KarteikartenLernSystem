package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Webprojekt/app/model"
)

//Index Generiert die Index Seite
func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/index.tmpl")
	log.Println(err)

	data, _ := model.GetIndexData(username)

	data.BaseDaten.LoggedIn = loggedIn
	t.Execute(w, data)
}

//Karteikasten Für die Karteikasten Seite
func Karteikasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/karteikasten.tmpl")
	log.Println(err)

	karteikastenDaten, _ := model.GetKarteikastenData(username)
	karteikastenDaten.BaseDaten.LoggedIn = loggedIn
	t.Execute(w, karteikastenDaten)
}

// Meinekarteien Für die MeineKarteikasten Seite
func Meinekarteien(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/meinekarteien.tmpl")
	log.Println(err)

	meinkarteikastenData, _ := model.GetMeineKarteikastenData(username)
	meinkarteikastenData.BaseDaten.LoggedIn = loggedIn

	err = t.Execute(w, meinkarteikastenData)
	fmt.Println(err)
}

//View Für die View Seite
func View(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/view.tmpl")
	log.Println(err)

	id := r.FormValue("id")
	nr := r.FormValue("nr")

	viewData, _ := model.GetViewData(username, id, nr)
	viewData.BaseDaten.LoggedIn = loggedIn

	t.Execute(w, viewData)
}

//Lern Für die Lern Seite
func Lern(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/lern.tmpl")
	log.Println(err)

	id := r.FormValue("id")
	aufgedeckt := r.FormValue("aufgedeckt")
	kartenNr := r.FormValue("kartennr")

	viewData, _ := model.GetLernData(username, id, aufgedeckt, kartenNr)
	viewData.BaseDaten.LoggedIn = loggedIn

	t.Execute(w, viewData)
}

//Edit Für die Edit Seite
func Edit(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/edit.tmpl")
	log.Println(err)

	id := r.FormValue("id")

	editDaten, _ := model.GetEditData(username, id)
	editDaten.BaseDaten.LoggedIn = loggedIn

	t.Execute(w, editDaten)
}

//Profil Für die Profil Seite
func Profil(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	t, err := template.ParseFiles("template/base.tmpl", "template/profil.tmpl")

	log.Println(err)
	profildata, _ := model.GetProfilData(username)
	profildata.BaseDaten.LoggedIn = loggedIn
	t.Execute(w, profildata)
}

//UpdateKarteiLern Für das Lernen
func UpdateKarteiLern(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	choise := r.FormValue("choise")
	kartenNr := r.FormValue("kartennr")

	model.UpdateKarteiLernen(id, choise, kartenNr)

	http.Redirect(w, r, "/karteikasten/lern/?id="+id, http.StatusFound)
}

//AddKarteikasten Updatet, oder erstellt einen neuen Karteikasten
func AddKarteikasten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	id := r.FormValue("id")
	name := r.FormValue("Titel")
	beschreibung := r.FormValue("Beschreibung")
	kategorie := r.PostFormValue("Ueberkategorie")
	oeffentlich := r.FormValue("Öffentlich")

	var sichtbarkeit string
	if oeffentlich == "on" {
		sichtbarkeit = "Öffentlich"
	} else {
		sichtbarkeit = "Privat"
	}

	if id == "" {
		aktlid, _ := model.CreateKarteikasten(username, name, beschreibung, kategorie, sichtbarkeit)
		http.Redirect(w, r, "/karteikasten/edit-karten/?id="+aktlid, http.StatusFound)
	} else {
		model.UpdateKarteikasten(id, name, beschreibung, kategorie, sichtbarkeit)
		http.Redirect(w, r, "/karteikasten/edit-karten/?id="+id, http.StatusFound)
	}
}

//EditKarten Holt das EditNData Struct für das Tmpl
func EditKarten(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["name"].(string)
	loggedIn := session.Values["loggedIn"].(bool)
	id := r.FormValue("id")
	KartenNr := r.FormValue("nr")
	t, _ := template.ParseFiles("template/base.tmpl", "template/editnext.tmpl")

	editnData, _ := model.GetEditNextData(username, id, KartenNr)
	editnData.BaseDaten.LoggedIn = loggedIn

	t.Execute(w, editnData)
}

//UpdateKarte Update der Karte im Karteikasten
func UpdateKarte(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	kartenNr := r.FormValue("nr")
	kartenNamen := r.FormValue("KartenName")
	antwort := r.PostFormValue("Antwort")
	frage := r.FormValue("Frage")

	var Karteikarte struct {
		KartenName string
		Fach       string
		Frage      string
		Antwort    string
		KartenNr   string
	}

	Karteikarte.KartenNr = kartenNr
	Karteikarte.Fach = "0"
	Karteikarte.KartenName = kartenNamen
	Karteikarte.Antwort = antwort
	Karteikarte.Frage = frage

	

	http.Redirect(w, r, "/karteikasten/edit-karten/?id="+id+"&nr="+kartenNr, http.StatusFound)
}
