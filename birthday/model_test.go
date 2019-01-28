package birthday

import (
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"github.com/vimsucks/borthday/util"
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
	dbFilePath := fmt.Sprintf("%s%s..%sborthday.db", cwd, string(os.PathSeparator), string(os.PathSeparator))
	db, err = sqlx.Connect("sqlite3", (dbFilePath))
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
		LunarBirthday: util.ParseDate("2018-11-26"),
		SolarBirthday: util.ParseDate("2018-11-22"),
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

