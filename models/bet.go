package models

import (
  "time"

  "github.com/satori/go.uuid"
)

type Bet struct {
  Base
  ID              string      `gorm:"primary_key; column:id; type:varchar(255)"`
  PlayerID        string      `gorm:"not null; column:player_id" `
  TournamentID    string      `gorm:"not null; column:tournament_id"`
  CreatedAt       time.Time   `gorm:"not null; column:created_at"`
  UpdatedAt       time.Time   `gorm:"not null; column:updated_at"`
  Backers         []Player    `gorm:"many2many:backers"`
}

func (model Bet) Create(tournamentId string, playerId string, backers []Player) (Bet) {
  return Bet{
    ID: uuid.NewV4().String(),
    PlayerID: playerId,
    TournamentID: tournamentId,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
    Backers: backers,
  }
}
