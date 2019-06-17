package main

import (
	"net/http"

	"github.com/Webprojekt/app/controller"
)

func main() {
	server := http.Server{Addr: ":8080"}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", controller.AuthIndex(controller.Index))
	http.HandleFunc("/karteikasten", controller.AuthKartei(controller.Karteikasten))
	http.HandleFunc("/meinekarteien", controller.AuthIndex(controller.Meinekarteien))
	http.HandleFunc("/karteikasten/view/", controller.AuthView(controller.View))
	http.HandleFunc("/karteikasten/edit/", controller.AuthKartei(controller.Edit))
	http.HandleFunc("/karteikasten/edit/add-Kasten/", controller.AddKarteikasten)
	http.HandleFunc("/karteikasten/edit-karten/", controller.AuthIndex(controller.EditKarten))
	http.HandleFunc("/karteikasten/edit-karten/update-karte/", controller.UpdateKarte)
	http.HandleFunc("/karteikasten/lern/", controller.AuthView(controller.Lern))
	http.HandleFunc("/karteikasten/lern/update/", controller.AuthView(controller.UpdateKarteiLern))
	http.HandleFunc("/profil", controller.AuthIndex(controller.Profil))

	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/add-user", controller.AddUser)
	http.HandleFunc("/profil/update-User", controller.UpdateUser)
	http.HandleFunc("/profil/del-User", controller.DelUser)
	//http.HandleFunc("/meinekarteien/del-kartei", controller.DelKartei)
	http.HandleFunc("/authenticate-user", controller.AuthenticateUser)
	http.HandleFunc("/karteikasten/authenticate-user", controller.AuthenticateUserKartei)

	http.HandleFunc("/logout", controller.Logout)

	server.ListenAndServe()
}
