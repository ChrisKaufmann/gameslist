package main
import (
	"../game"
	"github.com/golang/glog"
	"github.com/msbranco/goconfig"
	"flag"
	"database/sql"
	"github.com/ChrisKaufmann/goauth"
	"fmt"
	"github.com/ChrisKaufmann/ebay-go"
	"math/rand"
)

var (
	ebay_application_id string
	eb *ebay.EBay
	user auth.User
	db *sql.DB
)


func init() {
	var err error
	flag.Parse()
	c, err := goconfig.ReadConfigFile("../config")
	db_name, err := c.GetString("DB", "db")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.db): %s", err)
	}
	db_host, err := c.GetString("DB", "host")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.Host): %s", err)
	}
	db_user, err := c.GetString("DB", "user")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.user): %s", err)
	}
	db_pass, err := c.GetString("DB", "pass")
	if err != nil {
		glog.Fatalf("init(): config.GetString(DB.pass): %s", err)
	}
	db, err = sql.Open("mysql", db_user+":"+db_pass+"@tcp("+db_host+")/"+db_name)
	if err != nil {
		glog.Fatalf("Init():sql.Open(mysql, %s:%s@%s/%s: %s", db_user, db_pass, db_host, db_name, err)
	}
	ebay_application_id, err = c.GetString("Ebay","ApplicationID")
	if err != nil {
		glog.Fatalf("init():Config.GetString(Ebay,ApplicationID): %s", err)
	}
	eb = ebay.New(ebay_application_id)
	ebay_url, err := c.GetString("Ebay","URL")
	if err == nil {
		eb.URL = ebay_url
	}
}

func main() {
	defer db.Close()
	auth.DB(db)
	game.GameDB(db)
	game.ConsoleDB(db)
	var err error
	wantedgames, err := game.GetAllWantedGames()
	if err != nil {
		glog.Errorf("game.GetAllWantedGames(): %s", err)
	}
	for i := range wantedgames {
		j := rand.Intn(i + 1)
		wantedgames[i], wantedgames[j] = wantedgames[j], wantedgames[i]
	}
	fmt.Printf("Len(wantedgames): %v\n", len(wantedgames))
	c := make(chan game.Game)
	var ugl []game.Game
	for i := 0; i <10; i++ {
		go ProcessGameSearch(&ugl, c)
	}
	for _,i := range wantedgames {
		c <- i
	}

	for _, g := range ugl {
		fmt.Printf("\t%v\t%s\n", g.EbayPrice, g.Name)
	}
}
func ProcessGameSearch(ugl *[]game.Game, queue chan game.Game) {
	fmt.Printf("ProcessGameSearch\n")
	for r := range queue {
		r.UpdateEbay(eb)
		*ugl = append(*ugl, r)
	}
}