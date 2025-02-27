package server

import (
	"database/sql"
	"net/http"
)

var DB *sql.DB

func HandleRequests(db *sql.DB) {
	DB = db
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/site_static/", http.StripPrefix("/site_static/", fs))

	http.HandleFunc("/", homeHandler)

	// отрисовка страниц
	http.HandleFunc("/page/hero/specifications/", pageHeroSpecificationsHadler)
	http.HandleFunc("/page/hero/", pageHeroHandler)
	http.HandleFunc("/page/registr/", pageRegistr)
	http.HandleFunc("/login/page/", logiPage)
	http.HandleFunc("/page/lk/", pageLk)
	// ауетнификация
	http.HandleFunc("/login/", login)
	http.HandleFunc("/registr/", registUser)
	http.HandleFunc("/user/", lkUser)

	http.ListenAndServe(":8090", nil)
}
