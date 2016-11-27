package srv

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Davmuz/gqt"
	// postgres dialect
	_ "github.com/lib/pq"
)

// PG gorm Connection to Postgres
var PG *sql.DB

// InitPostgres opens a database connection to postgres
func InitPostgres() {

	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		"dock",
		"dock",
		"dock",
		"179.16.238.11",
		"5432",
	)

	var err error
	PG, err = sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)
		if err = PG.Ping(); err == nil {
			log.Println("\nPostgres server is running")
			break
		}
		log.Println(err)
	}

	gqt.Add("templates/sql", "*.sql")
	// TODO: users / admins schema
	// PG.AutoMigrate(db.User{})
}

// MigrateData is to be used only by dev environment to bootstrap the database schema
func MigrateData() {
	gqt.Add("templates/pgsql", "*.pgsql")
	_, err := PG.Query(gqt.Get("createSchema"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = PG.Query(gqt.Get("migrateData"))
	if err != nil {
		log.Fatal(err)
	}
}
