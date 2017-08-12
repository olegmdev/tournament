package models

type Backer struct {
  Base
  PlayerID   string   `gorm:"primary_key; not null; column:player_id" `
  BetID      string   `gorm:"primary_key; not null; column:bet_id"`
}
