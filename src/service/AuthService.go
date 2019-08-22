package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../entity"

	"github.com/gentwolf-shen/gohelper/convert"
	"github.com/gentwolf-shen/gohelper/cryptohelper/aes"
	"github.com/gentwolf-shen/gohelper/timehelper"
)

var (
	Auth = &AuthService{}
)

type AuthService struct {
	authItems map[string]entity.AuthConfig
}

func (this *AuthService) Init() {
	b, err := ioutil.ReadFile(filepath.Dir(os.Args[0]) + "/config/auth.json")
	if err != nil {
		panic(err)
	}

	this.authItems = make(map[string]entity.AuthConfig)
	err = json.Unmarshal(b, &this.authItems)
	if err != nil {
		panic(err)
	}
}

func (this *AuthService) CheckToken(token string) string {
	if token == "" {
		return ""
	}

	tmp := strings.Split(token, ":")
	if len(tmp) != 2 {
		return ""
	}

	appKey := tmp[0]
	config, ok := this.authItems[appKey]
	if !ok {
		return ""
	}

	rawStr, err := aes.New(aes.CBC, []byte(config.Secret[0:16]), []byte(config.Secret[16:])).DecryptToString(tmp[1])
	if err != nil {
		return ""
	}

	tmp = strings.Split(string(rawStr), "|")
	if len(tmp) != 3 {
		return ""
	}

	span := timehelper.Timestamp() - convert.ToInt64(tmp[2])
	if tmp[1] != appKey || span > 60 || span < -60 {
		return ""
	}

	return appKey
}

func (this *AuthService) checkAction(actions []string, action string) error {
	bl := false

	for _, v := range actions {
		if action == v {
			bl = true
			break
		}
	}

	if !bl {
		return errors.New("\"" + action + "\" is not allowed to execute")
	}

	return nil
}

func (this *AuthService) CheckSql(appKey, sql string) (string, error) {
	config := this.authItems[appKey]

	err := this.checkAction(config.Actions, strings.ToUpper(sql[0:6]))
	if err != nil {
		return "", err
	}

	return config.Database, nil
}
