package server

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
	"www/pkg/model"
)

var (
	SecretKeyToken = "qwevsd"
	ExpiresAtToken = 300 * time.Hour
)

type userJwtResponse struct {
	Token string `json:"token"`
}

func registUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	accountID := r.FormValue("account_id")
	if name == "" || password == "" || accountID == "" {

		fmt.Fprintf(w, "не все данные введены")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	res := DB.QueryRow("SELECT name FROM users WHERE name = ?", name)
	var userName string
	err = res.Scan(&userName)
	if err == nil {
		fmt.Fprintf(w, "такой пользователь уже есть ")
		return
	}

	type getPlayersByAccountIDSelectWLResponse struct {
		Win  int `json:"win"`
		Loss int `json:"lose"`
	}
	winrate := &getPlayersByAccountIDSelectWLResponse{}
	resp, err := http.Get(fmt.Sprintf("https://api.opendota.com/api/players/%v/wl", accountID))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "неверно введен аккаунт апи", err)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(winrate)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	//TODO: расскомментировать после теста
	if winrate.Win < 0 && winrate.Loss < 0 {
		fmt.Fprintf(w, "не верный id akk")
		return
	}
	winr := (float32(winrate.Win) / (float32(winrate.Win) + float32(winrate.Loss))) * float32(100)
	fmt.Println(winr)

	result, err := DB.Exec("INSERT INTO users (`name`,`password`,account_id,winrate) VALUES (?,?,?,?)", name, hash, accountID, winr)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	userID := strconv.Itoa(int(lastID))
	token, err := newJWT(userID, ExpiresAtToken, SecretKeyToken)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}
	fmt.Println(token)
	//TODO: 143693439 dota id
	response := userJwtResponse{
		Token: token,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	var user model.Users
	row := DB.QueryRow("SELECT password,id FROM users WHERE name = ?", name)
	err := row.Scan(&user.Password, &user.ID)

	if err != nil {
		fmt.Fprintf(w, "пользователь с таким логином не найден, %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Fprintf(w, "пароли не совпадают, %v", err)
		return
	}

	token, err := newJWT(strconv.FormatInt(user.ID, 10), ExpiresAtToken, SecretKeyToken)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}

	response := userJwtResponse{
		Token: token,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Fprintf(w, "err: %v", err)
		return
	}
}

func lkUser(w http.ResponseWriter, r *http.Request) {
	userId, err := authHeader(r)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	var userInfo model.Users
	row := DB.QueryRow("SELECT name,winrate,account_id FROM users WHERE id  = ?", userId)
	err = row.Scan(&userInfo.Name, &userInfo.Winrate, &userInfo.Account_id)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	info, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	_, err = w.Write(info)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
}
