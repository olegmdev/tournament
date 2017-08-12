package models

import (
  "fmt"
  "tournament/db"

  "github.com/jinzhu/gorm"
)

type Base interface{}

type modelAction func(*gorm.DB) (Base, error)

// Decorates model actions with a database transaction
func Sync(fn modelAction) (Base, error) {
  db := db.GetDB()
  tx := db.Begin()

  model, err := fn(tx)
  if err != nil {
    tx.Rollback()
    return model, err
  }

  tx.Commit()
  return model, nil
}

// Truncates database tables
func Reset() error {
  db := db.GetDB()

  tables := []string{"tournaments", "players"}
  for _, table := range tables {
    if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error; err != nil {
      return err
    }
  }

  return nil
}
