package model

import (
	"encoding/json"
	"fmt"
	"ioutil"
	"net/http"
	"sort"

	"github.com/labstack/echo"
)

// メンションに必要な情報
type UserInfo struct {
	ID       int    `db:"id"`
	TraqID   string `db:"traq_id`
	UserId   string `db:"user_id"`   // traQ内部のUUID
	UserType string `db:"user_type"` // user or group
}

// traQが返却するユーザーオブジェクト
type TraqUserObject struct {
	UserId        string `json:"userId db:"user_id`
	Name          string `json:"name db:"traq_id`
	DisplayName   string `json:"displayName`
	IconField     string `json:"iconField`
	Bot           bool   `json:"bot"`
	TwitterId     string `json:"twitterId`
	LastOnline    string `json:"lastOnline`
	IsOnline      string `json:"isOnline`
	Suspended     bool   `json:"suspended`
	AccountStatus int    `json:"accountStatus`
}

// traQに問い合わせてユーザーの情報をDBに取り込む
func FetchUserInfo() error {
	url := "https://q.trap.jp/api/1.0/users/"
	req, err := http.NewRequest("GET", url)
	if err != nil {
		return err
	}

	// set authorization header
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", "hoge"))

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var users []TraqUserObject
	err := json.Unmarshal(body, &users)
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err

	}

	// この辺修正しないと
	stmt, err := tx.PrepareNamed("INSERT INTO user_info (traq_id, user_id, user_type) VALUES (:name, :user_id, \"user\")")
	if err != nil {
		return err
	}

	for _, user := range users {
		if _, err := stmt.Queryx(user); err != nil {
			return err
		}
	}

	tx.Commit()
}

// traQに問い合わせてグループの情報をDBに取り込む
func FetchGroupsInfo() error {
	url := "https://q.trap.jp/api/1.0/groups/"
	req, err := http.NewRequest("GET", url)
	if err != nil {
		return nil, err
	}

	// set authorization header
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", "hoge"))

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var groups []TraqUserObject
	err := json.Unmarshal(body, &groups)
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	stmt, err := tx.Preparex("INSERT INTO user_info (user_name, uuid, user_type) VALUES (?, ?, group)")
	if err != nil {
		return err
	}

	for _, group := range groups {
		if _, err := stmt.Queryx(group); err != nil {
			return err
		}
	}

	tx.Commit()
}

func GetAllUserInfo() ([]UserInfo, error) {
	allUserInfo := []UserInfo{}

	if err := db.Select(&allUserInfo, "SELECT * FROM user_info ORDER BY user_name"); err != nil {
		return nil, err
	}

	return allUserInfo, nil
}

func MakeMentionText(user UserInfo) string {
	return fmt.Sprintf(`!{"type":"%s","raw":"@%s","id":"%s"}`, user.UserType, user.TraqID, user.Uuid)
}

func MakeMentionTexts(traqIds []string) ([]string, error) {
	allUserInfo, err := GetAllUserInfo()
	if err != nil {
		return nil, err
	}

	sort.Sort(traqIds)
	var ret []string

	j := 0
	for _, user := range allUserInfo {
		if user.TraqID == traqIds[j] {
			append(ret, MakeMentionText(user))
			j++
		} else if user.TraqID > traqIds[j] {
			// データにない人が出てきたのでDBを更新して再度呼び出す(ほんまか)
			// もうちょっといい書き方ありそう
			if err := FetchUserInfo(); err != nil {
				return nil, err
			}
			if err := FetchGroupsInfo(); err != nil {
				return nil, err
			}
			return MakeMentionTexts(traqIds)
		}
	}

	return ret, nil
}
