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
	"errors"
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
	stmtCookieIns *sql.Stmt
	stmtGetUserID *sql.Stmt
	stmtInsertUser *sql.Stmt
	stmtGetUser	*sql.Stmt
)

func CookieName(c string) {
	cookieName = c
}
func Environment(e string) {
	environment = e
}
func DB(d *sql.DB) {
	db = d
	var err error
	stmtCookieIns,err = u.Sth(db, "INSERT INTO sessions (user_id,session_hash) VALUES( ? ,?  )")
	if err != nil {	err.Error();fmt.Println(err)}
	stmtGetUserID,err = u.Sth(db, "select id from users where email = ?")
	if err != nil {	err.Error();fmt.Println(err)}
	stmtInsertUser,err = u.Sth(db, "insert into users (email) values (?) ")
	if err != nil {	err.Error();fmt.Println(err)}
	stmtGetUser,err = u.Sth(db, "select user_id from sessions as s where s.session_hash = ?")
	if err != nil {	err.Error();fmt.Println(err)}
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
	oauthCfg.RedirectURL = url + "/oauth2callback"
}

// Start the authorization process
func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	print("In handleauth\n")
	//Get the Google URL which shows the Authentication page to the user
	url := oauthCfg.AuthCodeURL("")

	//redirect user to that page
	http.Redirect(w, r, url, http.StatusFound)
}

// Function that handles the callback from the Google server
func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
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
	//body.id is the google id to use
	//set a cookie with the id, and random hash. then save the id/hash pair to db for lookup
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {fmt.Println(err);err.Error();return}
	m := f.(map[string]interface{})
	print(m["email"].(string))
	var authString = u.RandomString(64)
	email := m["email"].(string)
	var uid int
	if ! UserExists(email) {
		uid, err = AddUser(email)
		if err != nil {fmt.Println(err);err.Error();return }
	} else {
		uid, err = GetUser(email)
		if err != nil {fmt.Println(err);err.Error();return }
	}

	print("Inserting authstring "+authString+" into DB\n")
	_, err = stmtCookieIns.Exec(uid, authString)

	if err != nil {	err.Error();fmt.Println(err);return	}
	//set the cookie
	expire := time.Now().AddDate(1, 0, 0) // year expirey seems reasonable
	cookie := http.Cookie{Name: cookieName, Value: authString, Expires: expire}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusFound)
}
func UserExists(email string)(exists bool) {
	var uid int
	err := stmtGetUserID.QueryRow(email).Scan(&uid)
	switch {
		case err == sql.ErrNoRows:
			exists = false
		case err != nil:
			fmt.Println(err);err.Error();exists=false
		default:
			exists = true
	}
	return exists
}
func AddUser(e string)(uid int, err error) {
	if UserExists(e) {
		err = stmtGetUserID.QueryRow(e).Scan(&uid)
		return uid,err
	}
	result, err := stmtInsertUser.Exec(e)
	if err != nil { fmt.Println(err);err.Error();return uid,err}
    lid, err := result.LastInsertId()
    uid=int(lid)
	return uid, err
}
func GetUser(e string)(uid int, err error) {
	if !UserExists(e) {
		err=errors.New("User Doesn't exist")
		return 0, err
	}
	err = stmtGetUserID.QueryRow(e).Scan(&uid)
	if err != nil {	fmt.Println(err);err.Error()}
	return uid, err
}
func LoggedIn(w http.ResponseWriter, r *http.Request) (bool, int) {
	var userId int
	if environment == "test" {
		return true, 1
	}
	print("cookieName="+cookieName+"\n")
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		//just means that the cookie doesn't exist or we couldn't read it
		print("No cookie found")
		return false, 0
	}
	tokHash := cookie.Value
	print("tokhash="+tokHash+"\n")
	err = stmtGetUser.QueryRow(tokHash).Scan(&userId)
	print("userid="+u.Tostr(userId)+"\n")
	switch {
		case err == sql.ErrNoRows:
			return false,0  //no rows returned
		case err != nil:                     //probably no results in query
			fmt.Println(err)
			err.Error()
			return false,0
		default:
			if userId > 0 {
				return true, userId
			}
	}
	return false, 0
}
