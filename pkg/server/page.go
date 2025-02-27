package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"www/pkg/model"
)

func pageLk(w http.ResponseWriter, r *http.Request) {

	bytedata, err := io.ReadAll(r.Body)
	reqBodyString := string(bytedata)

	fmt.Println(reqBodyString)
	tmpl, _ := template.ParseFiles("templates/lk.html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Execute) SERVER ERROR: %v", err)))
		return
	}

}

func logiPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Execute) SERVER ERROR: %v", err)))
		return
	}
}

func pageRegistr(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/registr.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Execute) SERVER ERROR: %v", err)))
		return
	}
}

// генерация главной стр
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Execute) SERVER ERROR: %v", err)))
		return
	}
}

// отображение списка героев на стр
func pageHeroSpecificationsHadler(w http.ResponseWriter, r *http.Request) {
	res, err := DB.Query("SELECT heroes.id, heroes.name,attribute.name AS attribute_name,heroes.winrate FROM `heroes` JOIN attribute ON heroes.attribute_id=attribute.id")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
		return
	}

	var heroes = []model.Hero{}
	for res.Next() {
		var hero model.Hero
		err = res.Scan(&hero.ID, &hero.Name, &hero.Attribute, &hero.Winrate)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
		}
		heroes = append(heroes, hero)
	}

	tmpl, err := template.ParseFiles("templates/specifications.html")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
	}
	err = tmpl.Execute(w, heroes)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
	}

}

// отображение стр героя по ид
func pageHeroHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	resHero := DB.QueryRow("SELECT heroes.id, heroes.name,attribute.name AS attribute_name,heroes.winrate,heroes.picha,heroes.history FROM `heroes` JOIN attribute ON heroes.attribute_id=attribute.id WHERE heroes.id = ?", id)
	var hero model.Hero

	err := resHero.Scan(&hero.ID, &hero.Name, &hero.Attribute, &hero.Winrate, &hero.Picha, &hero.History)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Scan) SERVER ERROR: %v", err)))
		fmt.Println(err)
		return
	}
	resItems, err := DB.Query("SELECT  items.* FROM items_hero JOIN heroes ON heroes.id = items_hero.heroesId JOIN items ON items.Id = items_hero.itemsId WHERE heroes.id =? ", id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
		return
	}

	var items []model.Item

	for resItems.Next() {
		var item = model.Item{}
		err = resItems.Scan(&item.ID, &item.Name, &item.Price, &item.Picha)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("SERVER ERROR: %v", err)))
			return
		}
		items = append(items, item)
	}
	hero.Items = items
	tmpl, _ := template.ParseFiles("templates/show_hero.html")
	err = tmpl.Execute(w, hero)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("(Execute) SERVER ERROR: %v", err)))
		return
	}

}
