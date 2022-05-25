package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

//var conn *gorm.DB // database
var username, password, dbName, dbHost, dbPort string

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	v4uuid := uuid.NewV4()
	return scope.SetColumn("ID", v4uuid)
}

// always runs
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	username = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")

	migrateDatabase()
}

func migrateDatabase() {
	db := ConnectDatabase()
	// migrate Article struct
	db.Debug().AutoMigrate(
		&Article{},
	)
}

func ConnectDatabase() *gorm.DB {
	pgDbUri := fmt.Sprintf("postgres://%v@%v:%v/%v?sslmode=disable&password=%v", username, dbHost, dbPort, dbName, password)

	// connect with |postgres| dialect.
	conn, err := gorm.Open("postgres", pgDbUri)
	if err != nil {
		log.Println(err)
	}

	return conn
}
