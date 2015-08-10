package auth

import (
	"code.google.com/p/goauth2/oauth"
	"database/sql"
	"encoding/json"
	"fmt"
	u "github.com/ChrisKaufmann/goutils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/msbranco/goconfig"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var oauthCfg = &oauth.Config{
	AuthURL:    "https://accounts.google.com/o/oauth2/auth",
	TokenURL:   "https://accounts.google.com/o/oauth2/token",
	Scope:      "https://www.googleapis.com/auth/userinfo.email",
	TokenCache: oauth.CacheFile(cachefile),
}

const profileInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo"
const cachefile = "/dev/null"

var (
	MyURL       string
	db          *sql.DB
	cookieName  string
	environment string
)

func CookieName(c string) {
	cookieName = c
}
func Environment(e string) {
	environment = e
}
func DB(d *sql.DB) {
	db = d
}
func init() {
	c, err := goconfig.ReadConfigFile("config")
	if err != nil {
		panic(err)
	}
	oauthCfg.ClientId, err = c.GetString("Google", "ClientId")
	if err != nil {
		panic(err)
	}
	oauthCfg.ClientSecret, err = c.GetString("Google", "ClientSecret")
	if err != nil {
		panic(err)
	}
	url, err := c.GetString("Web", "url")
	MyURL = url
	if err != nil {
		panic(err)
	}
	oauthCfg.RedirectURL = url + "oauth2callback"
}

// Start the authorization process
func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	print("In handleauth\n")
	//Get the Google URL which shows the Authentication page to the user
	url := oauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	//Get the code from the response
	code := r.FormValue("code")

	print("code=" + code)
	t := &oauth.Transport{Config: oauthCfg}

	// Exchange the received code for a token
	tok, err := oauthCfg.TokenCache.Token()
	if err != nil {
		print(err)
		tok, err = t.Exchange(code)
		if err != nil {
			print(err)
			panic(err.Error())
		}
		fmt.Printf("token cached in %v\n", oauthCfg.TokenCache)
	}
	print(tok)

	// Make the request.
	req, err := t.Client().Get(profileInfoURL)
	if err != nil {
		print(err)
		print("\n")
		panic(err.Error())
		return
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	log.Println(string(body))
	//print(string(body))
	//body.id is the google id to use
	//set a cookie with the id, and random hash. then save the id/hash pair to db for lookup
	var f interface{}
	err = json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	print(m["email"].(string))
	if err != nil {
		panic(err.Error())
	}
	var authString = u.RandomString(64)
	makeSureUserExists(m["email"].(string))
	stmtCookieIns := u.Sth(db, "INSERT INTO sessions (name,userid) VALUES( ? ,?  )")
	_, err = stmtCookieIns.Exec(m["email"], u.GetHash(authString))

	if err != nil {
		panic(err.Error())
	}
	//set the cookie
	expire := time.Now().AddDate(1, 0, 0) // year expirey seems reasonable
	cookie := http.Cookie{Name: cookieName, Value: authString, Expires: expire}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusFound)
}
func makeSureUserExists(e string) {
	var userId string
	stmtGetUsername := u.Sth(db, "select username from users where email = ?")
	err := stmtGetUsername.QueryRow(e).Scan(&userId)
	if err != nil {
		// in this case, emans there's no username - create one
		stmtInsertUser := u.Sth(db, "insert into users (username,email) values ( ?, ?) ")
		_, err = stmtInsertUser.Exec(e, e)
	}
}
func LoggedIn(w http.ResponseWriter, r *http.Request) (loggedIn bool, userName string) {
	if environment == "test" {
		userName = "chris"
		loggedIn = true
		return loggedIn, userName
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		//just means that the cookie doesn't exist or we couldn't read it
		loggedIn = false
		userName = ""
		return loggedIn, userName
	}
	tok := cookie.Value
	tokHash := u.GetHash(tok)
	var userId string
	stmtGetUser := u.Sth(db, "select u.username from users as u, sessions as s where s.userid = ? and s.name=u.email")
	err = stmtGetUser.QueryRow(tokHash).Scan(&userId) // WHERE number = 13
	if err != nil {                                   //probably no results in query
		loggedIn = false
		userName = ""
		return loggedIn, userName
	}

	if userId != "" {
		userName = userId
		loggedIn = true
		return loggedIn, userName
	}
	loggedIn = false
	userName = ""
	return loggedIn, userName
}
