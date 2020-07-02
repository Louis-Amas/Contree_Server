package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // pg
)

var userSchema = `
CREATE TABLE users (
	id serial primary key,
	name text unique,
	password text,
    secret text
);

`
var salt = "fdsakfmsdlakfjmlasdofj"

// HashString hash string
func HashString(str string) string {
	hash := sha256.Sum256([]byte(salt + ":" + str))
	return string(base64.StdEncoding.EncodeToString(hash[:32]))
}

var db *sqlx.DB

// "user=contree password=contree dbname=contree sslmode=disable"
// "host=db user=contree password=contree dbname=contree sslmode=disable"
// InitDB initialize db
func InitDB(initDB bool) {
	var err error
	db, err = sqlx.Connect("postgres", "host=db user=contree password=contree dbname=contree sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if initDB {
		initSchema()
	}
}

func initSchema() {
	fmt.Println("[DB] Init Schemas")
	tx := db.MustBegin()
	tx.MustExec(userSchema)
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	AddUsers()
}

func getDB() *sqlx.DB {
	return db
}

// AddUsers add users to db
func AddUsers() {
	fmt.Println("[DB] Add Users")
	u := newUser("louis", "123")
	u.InsertUser()
	u = newUser("bob", "123")
	u.InsertUser()
	u = newUser("john", "123")
	u.InsertUser()
	u = newUser("doe", "123")
	u.InsertUser()
}
