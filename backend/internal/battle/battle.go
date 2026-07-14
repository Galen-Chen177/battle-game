package battle

type Battle struct {
	Left *Team

	Right *Team

	Time float64

	Events []Event
}

func NewBattle(left, right *Team) *Battle {
	return &Battle{
		Left:  left,
		Right: right,
	}
}
