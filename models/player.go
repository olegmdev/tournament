package models

import (
  "math"
  "time"
  "errors"

  "tournament/db"
  "tournament/forms"

  "github.com/satori/go.uuid"
  "github.com/jinzhu/gorm"
)

type Player struct {
  Base
  ID        string    `gorm:"primary_key; column:id; type:varchar(255)"`
  Balance   float64   `gorm:"not null; column:balance"`
  CreatedAt time.Time `gorm:"not null; column:created_at"`
  UpdatedAt time.Time `gorm:"not null; column:updated_at"`
}

func (player *Player) BeforeCreate(scope *gorm.Scope) error {
  if _, err := uuid.FromString(player.ID); err != nil {
    return errors.New("Invalid ID")
  }

  return nil
}

func createPlayer(data forms.Player, tx *gorm.DB) (Player, error) {
  var id string
  if id = data.ID; id == "" {
    id = uuid.NewV4().String()
  }

  var player Player = Player{
    ID: id,
    Balance: math.Abs(data.Points),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  err := tx.Create(&player).Error
  return player, err
}

func (model Player) verify(data forms.Player) error {
  if model.ID == "" {
    return errors.New("Player not found")
  }

  if model.Balance < data.Points {
    return errors.New("Points not enough")
  }

  return nil
}

// Funds a player with specified amount of points.
// If player not exists it will be created with proper balance
func (model Player) Fund(data forms.Player) (Base, error) {
  if data.Points < 0 {
    return nil, errors.New("Points amount can not be negative")
  }

  return Sync(func (tx *gorm.DB) (Base, error) {
    var player Player

    if data.ID == "" {
      return createPlayer(data, tx)
    }

    db.Lock(tx).Where(&Player{ID: data.ID}).First(&player)

    if player.ID != "" {
      player.Balance += math.Abs(data.Points)
      tx.Save(&player)
      return player, nil
    }

    return createPlayer(data, tx)
  })
}

// Takes specified amount of points (if exist) from player
func (model Player) Take(data forms.Player) (Base, error) {
  if data.ID == "" {
    return nil, errors.New("ID not specified")
  }

  if data.Points < 0 {
    return nil, errors.New("Points amount can not be negative")
  }

  return Sync(func (tx *gorm.DB) (Base, error) {
    var player Player

    db.Lock(tx).Where(&Player{ID: data.ID}).First(&player)

    if err := player.verify(data); err != nil {
      return player, err
    }

    player.Balance -= math.Abs(data.Points)
    tx.Save(&player)
    return player, nil
  })
}

// Gets player balance
func (model Player) Get(playerId string) (Base, error) {
  var player Player
  db := db.GetDB()

  db.Where(&Player{ID: playerId}).First(&player)
  if player.ID == "" {
    return player, errors.New("Player not found")
  }

  return player, nil
}
