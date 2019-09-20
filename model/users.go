package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"os"

	"github.com/labstack/echo"
)

type (
	// メンションに必要な情報(anke-to DBが持つ情報)
	UserInfo struct {
		Name     string `db:"name"`      // traQ ID
		UserId   string `db:"user_id"`   // traQ内部のUUID
		UserType string `db:"user_type"` // user or group
	}

	// traQが返却するユーザーオブジェクト
	TraqUser struct {
		UserId        string `json:"userId" db:"user_id"`
		Name          string `json:"name" db:"name"`
		DisplayName   string `json:"displayName"`
		IconField     string `json:"iconField"`
		Bot           bool   `json:"bot"`
		TwitterId     string `json:"twitterId"`
		LastOnline    string `json:"lastOnline"`
		IsOnline      string `json:"isOnline"`
		Suspended     bool   `json:"suspended"`
		AccountStatus int    `json:"accountStatus"`
	}
	// traQが返却するグループオブジェクト
	TraqUserGroup struct {
		GroupId     string   `json:"groupId" db:"user_id"`
		Name        string   `json:"name" db:"name"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		AdminUserId string   `json:"adminUserId"`
		Members     []string `json:"members"`
		CreatedAt   string   `json:"createdAt"`
		UpdatedAt   string   `json:"updatedAt"`
	}
)

// traQに問い合わせてユーザーの情報をDBに取り込む
func FetchUserInfo(names []string) ([]UserInfo, error) {
	// ユーザー一覧の取得
	userUrl := "https://q.trap.jp/api/1.0/users"
	reqUser, err := http.NewRequest("GET", userUrl, nil)
	if err != nil {
		return nil, err
	}
	reqUser.Header.Set(echo.HeaderAuthorization, "Bearer " + os.Getenv("TRAQ_BOT_ACCESSTOKEN"))
	client := &http.Client{}

	respUser, err := client.Do(reqUser)
	if err != nil {
		return nil, err
	}
	defer respUser.Body.Close()

	userBody, err := ioutil.ReadAll(respUser.Body)
	if err != nil {
		return nil, err
	}
	var users []TraqUser
	if err := json.Unmarshal(userBody, users); err != nil {
		return nil, err
	}

	// グループ一覧の取得
	groupUrl := "https://q.trap.jp/api/1.0/groups"
	reqGroup, err := http.NewRequest("GET", groupUrl, nil)
	if err != nil {
		return nil, err
	}
	reqGroup.Header.Set(echo.HeaderAuthorization, "Bearer " + os.Getenv("TRAQ_BOT_ACCESSTOKEN"))

	respGroup, err := client.Do(reqGroup)
	if err != nil {
		return nil, err
	}
	defer respGroup.Body.Close()

	groupBody, err := ioutil.ReadAll(respGroup.Body)
	if err != nil {
		return nil, err
	}
	var groups []TraqUserGroup
	if err := json.Unmarshal(groupBody, groups); err != nil {
		return nil, err
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	newUsers := make([]UserInfo, 0, len(names))
	// TODO: なんかもうちょいいいのを考える
	u, g := 0, 0
	for _, name := range names {
		if users[u].Name == name {
			_, err := tx.NamedExec("INSERT INTO userinfo (name, user_id, user_type) VALUES (:name, :user_id, \"user\")", &users[u])
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			newUsers = append(newUsers, UserInfo{users[u].Name, users[u].UserId, "user"})
			u++
		} else if groups[g].Name == name {
			_, err := tx.NamedExec("INSERT INTO userinfo (name, user_id, user_type) VALUES (:name, :user_id, \"group\")", &groups[g])
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			newUsers = append(newUsers, UserInfo{users[u].Name, users[u].UserId, "group"})
			g++
		} else {
			u++
			g++
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return newUsers, nil
}

func GetAllUserInfo() ([]UserInfo, error) {
	allUserInfo := []UserInfo{}

	if err := db.Select(&allUserInfo, "SELECT * FROM user_info ORDER BY user_name"); err != nil {
		return nil, err
	}

	return allUserInfo, nil
}

func MakeMentionText(user UserInfo) string {
	return fmt.Sprintf(`!{"type":"%s","raw":"@%s","id":"%s"}`, user.UserType, user.Name, user.UserId)
}

func MakeMentionTexts(traqIds []string) ([]string, error) {
	allUserInfo, err := GetAllUserInfo()
	if err != nil {
		return nil, err
	}

	sort.Strings(traqIds)

	mentions := make([]string, 0, len(traqIds))
	unknowns := make([]string, 0, len(traqIds)) // DBになかった人の一覧

	for i, j := 0, 0;i < len(allUserInfo) && j < len(traqIds); {
		if allUserInfo[i].Name == traqIds[j] {
			mentions = append(mentions, MakeMentionText(allUserInfo[i]))
			i++
			j++
		} else if allUserInfo[i].Name < traqIds[j] {
			i++
		} else {
			// その人の情報がなかった
			unknowns = append(unknowns, traqIds[j])
			j++
		}
	}

	// いなかった人の分を追加し、メンションのテキスト群に追加
	newUsers, err := FetchUserInfo(unknowns)
	if err != nil {
		return nil, err
	}

	for _, newUser := range newUsers {
		mentions = append(mentions, MakeMentionText(newUser))
	}

	return mentions, nil
}
