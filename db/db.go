package db

import (
  "fmt"
  "tournament/config"
  _ "github.com/lib/pq"
  "github.com/mattes/migrate"
  "github.com/mattes/migrate/database/postgres"
  _ "github.com/mattes/migrate/source/file"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func Init() {
  var err error
  config := config.GetConfig()

  if DB, err = gorm.Open("postgres", config.GetString("database")); err != nil {
    panic(fmt.Sprintf("No error should happen when connecting to database, but got err=%+v", err))
  }

  // Run database migrations
  driver, err := postgres.WithInstance(DB.DB(), &postgres.Config{})
  migrator, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
  migrator.Up()
}

func GetDB() *gorm.DB {
  return DB
}

func Lock(instance *gorm.DB) *gorm.DB {
  return instance.Set("gorm:query_option", "FOR UPDATE")
}
