package forms

type Tournament struct {
  ID        string    `json:"id"`
  Deposit   float64   `json:"deposit" binding:"required"`
}

type TournamentJoin struct {
  PlayerID      string    `json:"player_id" binding:"required"`
  TournamentID  string    `json:"tournament_id" binding:"required"`
  Backers       []string  `json:"backers"`
}

type Winner struct {
  PlayerID  string    `json:"player_id" binding:"required"`
  Prize     float64   `json:"prize" binding:"required"`
}

type TournamentResult struct {
  TournamentID  string    `json:"tournament_id" binding:"required"`
  Winners       []Winner  `json:"winners" binding:"required"`
}
