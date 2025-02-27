package processor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"www/pkg/model"
)

type OpenDotaHero struct {
	ID             int    `json:"id"`
	PrimaryAttr    string `json:"primary_attr"`
	Localized_name string `json:"localized_Name"`
}

func MegaUpdateHero(db *sql.DB) {
	for {
		updateHero(db)
		updaeteItems(db)
		time.Sleep(30 * time.Minute)
	}
}

// ----------------

func updateHero(db *sql.DB) {

	res, err := httpGet("https://api.opendota.com/api/heroes")
	if err != nil {
		fmt.Println(err)
		return
	}
	var heroes []OpenDotaHero
	err = json.NewDecoder(res.Body).Decode(&heroes)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, h := range heroes {
		fmt.Println(fmt.Sprintf("start process id:%v, name:%v", h.ID, h.Localized_name))

		idAtribute, err := updateAttribyte(db, h.PrimaryAttr)
		if err != nil {
			fmt.Println(err)
			return
		}

		winrate, err := calcWinrate(h)
		if err != nil {
			fmt.Println(err)
			return
		}

		var hero = model.Hero{
			ID:          h.ID,
			Name:        h.Localized_name,
			AttributeID: idAtribute,
			Winrate:     winrate,
		}

		result := db.QueryRow("SELECT id, winrate FROM heroes where id = ? ", h.ID)
		var idDB int
		var winrateDB int
		err = result.Scan(&idDB, &winrateDB)
		if idDB != 0 && winrate != winrateDB {
			_, err = db.Exec("UPDATE heroes SET winrate = ? WHERE id = ?", winrate, hero.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(fmt.Sprintf("update hero wirate %v on %v", winrateDB, winrate))
		}
		if idDB == 0 {
			_, err = db.Exec("INSERT INTO heroes (`id`,`name`,`winrate`,`attribute_id`,`picha`,`history`) VALUES (?,?,?,?,?,?)", hero.ID, hero.Name, hero.Winrate, hero.AttributeID, "", "")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(fmt.Sprintf("New hero %v ", hero.Name))
		}

		updateItemsHero(db, hero)

	}
}
func updaeteItems(db *sql.DB) {
	fmt.Println("items parse start")

	type OpenDotaitems struct {
		ID    int    `json:"id"`
		Name  string `json:"dname"`
		Price int    `json:"cost"`
	}

	var items = map[string]OpenDotaitems{}
	resp, err := httpGet("https://api.opendota.com/api/constants/items")
	if err != nil {
		fmt.Println(err)
		return
	}
	//var item model.Item
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		fmt.Println(err)
		return
	}
	var countAddItems = 0
	for _, value := range items {
		var idDBItem int
		resul := db.QueryRow("SELECT id FROM items WHERE id = ?", value.ID)
		_ = resul.Scan(&idDBItem)
		if idDBItem == 0 {
			_, err := db.Exec("INSERT INTO items (id,name,price,picha) VALUES (?,?,?,?)", value.ID, value.Name, value.Price, "")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(fmt.Sprintf("new item name: %v, id: %v", value.Name, value.ID))
			countAddItems++
		}

	}

	fmt.Println(fmt.Sprintf("items parse end count :%v", countAddItems))

}

// ----------------

func updateItemsHero(db *sql.DB, hero model.Hero) {
	type ItemsPopulariti struct {
		StartGameItems map[string]int `json:"start_game_items"`
		EarlyGameItems map[string]int `json:"early_game_items"`
		LateGameItems  map[string]int `json:"late_game_items"`
	}

	res, err := httpGet(fmt.Sprintf("https://api.opendota.com/api/heroes/%v/itemPopularity", hero.ID))
	if err != nil {
		fmt.Println(err)
		return
	}
	var items ItemsPopulariti
	err = json.NewDecoder(res.Body).Decode(&items)
	if err != nil {
		fmt.Println(err)
		return
	}
	i := 0
	itemIds := []string{}
	for kay := range items.EarlyGameItems {
		i++
		if i > 2 {
			i = 0
			break
		}
		itemIds = append(itemIds, kay)
	}
	for kay := range items.EarlyGameItems {
		i++
		if i > 2 {
			i = 0
			break
		}
		itemIds = append(itemIds, kay)
	}
	for kay := range items.LateGameItems {
		i++
		if i > 2 {
			i = 0
			break
		}
		itemIds = append(itemIds, kay)
	}
	itemIds = removeDuplicates(itemIds)

	_, err = db.Query("DELETE FROM items_hero WHERE heroesId = ?", hero.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range itemIds {
		_, err = db.Exec("INSERT INTO items_hero (heroesId,itemsId) VALUES (?,?)", hero.ID, v)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func calcWinrate(h OpenDotaHero) (int, error) {
	res, err := httpGet(fmt.Sprintf("https://api.opendota.com/api/heroes/%v/matches", h.ID))
	if err != nil {
		return 0, fmt.Errorf("err get matches from hero %v, err: %v", h.Localized_name, err)
	}

	type WinLoss struct {
		RadiantWin bool `json:"radiant_win"`
		Radiant    bool `json:"radian"`
	}
	var winLoss []WinLoss
	err = json.NewDecoder(res.Body).Decode(&winLoss)
	if err != nil {
		return 0, fmt.Errorf("err decode matches from hero %v, err: %v", h.Localized_name, err)
	}

	winrate := 0
	for _, w := range winLoss {
		if (w.RadiantWin || w.Radiant == true) && (w.RadiantWin || w.Radiant == false) {
			winrate++
		}
	}

	return winrate, nil
}
func updateAttribyte(db *sql.DB, primaryAttr string) (int64, error) {
	var idAtribute int64
	row := db.QueryRow("SELECT id FROM attribute WHERE name = ?", primaryAttr)
	err := row.Scan(&idAtribute)
	if err != nil {
		result, err := db.Exec("INSERT INTO attribute(name ) VALUES (?)", primaryAttr)
		if err != nil {
			return 0, err
		}
		idAtribute, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}
	return idAtribute, nil
}

// ------- helpers ---------

func httpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MyApp/1.0")
	for {
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if res.StatusCode == http.StatusOK {
			return res, nil
		}
		fmt.Println("ждем", res.StatusCode)
		time.Sleep(60 * time.Second)
	}

}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if contains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}
