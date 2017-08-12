package forms

type Player struct {
  ID       string    `json:"id"`
  Points   float64   `json:"points" binding:"required"`
}
