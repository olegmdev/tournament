package models

import (
  "math"
  "time"
  "errors"
  "fmt"
  "database/sql/driver"

  "tournament/db"
  "tournament/forms"

  "github.com/satori/go.uuid"
  "github.com/jinzhu/gorm"
  "github.com/ahl5esoft/golang-underscore"
)

type TournamentStatus string
const (
  Active TournamentStatus = "active"
  Closed TournamentStatus = "closed"
)

func (u *TournamentStatus) Scan(value interface{}) error {
  *u = TournamentStatus(value.([]byte))
  return nil
}

func (u TournamentStatus) Value() (driver.Value, error) {
  return string(u), nil
}

type Tournament struct {
  Base
  ID          string            `gorm:"primary_key; column:id; type:varchar(255)"`
  Deposit     float64           `gorm:"not null; column:deposit;"`
  Status      TournamentStatus  `gorm:"not null; column:status; type:ENUM('active', 'closed')"`
  CreatedAt   time.Time         `gorm:"not null; column:created_at"`
  UpdatedAt   time.Time         `gorm:"not null; column:updated_at"`
}

func (tournament *Tournament) BeforeCreate(scope *gorm.Scope) error {
  if _, err := uuid.FromString(tournament.ID); err != nil {
    return errors.New("Invalid ID")
  }

  return nil
}

func createTournament(data forms.Tournament) (Tournament, error) {
  var id string
  if id = data.ID; id == "" {
    id = uuid.NewV4().String()
  }

  var tournament Tournament = Tournament{
    ID: id,
    Deposit: math.Abs(data.Deposit),
    Status: Active,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }
  db := db.GetDB()
  err := db.Create(&tournament).Error
  return tournament, err
}

func checkTournament(tournament Tournament) error {
  if tournament.ID != "" && tournament.Status == Closed {
    return errors.New("Tournament had been closed")
  }

  if tournament.ID == "" {
    return errors.New("Tournament not found")
  }

  return nil
}

func reward(bet *Bet, prizes map[string]float64, tx *gorm.DB) error {
  var player Player
  db.Lock(tx).Where(&Player{ID: bet.PlayerID}).First(&player)
  if player.ID == "" {
    return errors.New("Winner not found")
  }

  var backers []Player
  db.Lock(tx).Model(bet).Related(&backers, "Backers")

  if amount := len(backers); amount > 0 {
    prize := prizes[bet.PlayerID] / float64(amount + 1)

    for index, _ := range backers {
      backer := &backers[index]
      backer.Balance += prize
      if err := tx.Save(&backer).Error; err != nil {
        return err
      }
    }

    player.Balance += prize
  } else {
    player.Balance += prizes[bet.PlayerID]
  }

  err := tx.Save(&player).Error

  return err
}

// Announces a tournament with specified deposit value
// Tournament ID may be specified optionally. Only valid uuid IDs are allowed.
func (model Tournament) Announce(data forms.Tournament) (Base, error) {
  var tournament Tournament
  db := db.GetDB()

  if data.ID == "" {
    return createTournament(data)
  }

  db.Where(&Tournament{ID: data.ID}).First(&tournament)

  if tournament.ID != "" {
    if tournament.Status == Closed {
      return tournament, errors.New("Tournament had been closed")
    }

    return tournament, errors.New("Tournament is already announced")
  }

  return createTournament(data)
}

// Join players to tournaments. Backing is not mandatory but possible.
// Player can play on his own money (points). Funds will be deducted based on participants amount.
func (model Tournament) Join(data forms.TournamentJoin) (Base, error) {
  return Sync(func (tx *gorm.DB) (Base, error) {
    var tournament Tournament
    db.Lock(tx).Where(&Tournament{ID: data.TournamentID}).First(&tournament)
    if err := checkTournament(tournament); err != nil {
      return nil, err
    }

    // calculate a value to subtract from player and backers if exist
    var price float64 = tournament.Deposit / float64(len(data.Backers) + 1)

    var player Player
    db.Lock(tx).Where(&Player{ID: data.PlayerID}).First(&player)
    if player.Balance < price {
      return nil, errors.New("Player doesn't have proper amount of points")
    }
    player.Balance -= price
    if err := tx.Save(&player).Error; err != nil {
      return nil, err
    }

    var backers []Player
    db.Lock(tx).Where(data.Backers).Find(&backers)
    for index, _ := range backers {
      backer := &backers[index]
      if backer.Balance < price {
        return nil, errors.New(fmt.Sprintf("Backer %s doesn't have proper amount of points", backer.ID))
      }

      backer.Balance -= price
      if err := tx.Save(&backer).Error; err != nil {
        return nil, err
      }
    }

    bet := Bet{}.Create(tournament.ID, player.ID, backers)
    err := tx.Save(&bet).Error

    return bet, err
  })
}

// Finish the tournament and reward winners with proper prize value
func (model Tournament) Result(data forms.TournamentResult) (Base, error) {
  return Sync(func (tx *gorm.DB) (Base, error) {
    var tournament Tournament
    db.Lock(tx).Where(&Tournament{ID: data.TournamentID}).First(&tournament)
    if err := checkTournament(tournament); err != nil {
      return nil, err
    }

    var bets []Bet
    prizes := make(map[string]float64)
    for _, winner := range data.Winners {
      prizes[winner.PlayerID] = winner.Prize
    }
    db.Lock(tx).Where(&Bet{TournamentID: tournament.ID}).Where("player_id in (?)", underscore.Pluck(data.Winners, "PlayerID")).Find(&bets)
    for index, _ := range bets {
      bet := &bets[index]
      if err := reward(bet, prizes, tx); err != nil {
        return tournament, err
      }
    }

    tournament.Status = Closed
    err := tx.Save(&tournament).Error

    return tournament, err
  })
}
