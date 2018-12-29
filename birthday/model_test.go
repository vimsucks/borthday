package birthday

import (
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init () {
	var err error
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	db, err = sqlx.Connect("sqlite3", cwd + "../borthday.db")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func TestCRUD(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int63()
	fmt.Println(uid)
	b := Birthday{
		UID: uid,
		Name: "test",
		LunarBirthday: "2018-11-26",
		SolarBirthday: "2018-11-22",
	}
	err := CreateBirthday(&b)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
	query := sqrl.Eq{}
	query["uid"] = b.UID
	birthday, err := GetBirthday(query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(birthday)
	err = DeleteBirthday(query)
	if err != nil {
		t.Error(err)
	}
}

