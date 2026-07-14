package battle

import "fmt"

type Unit struct {
	ID         uint
	Name       string
	HP         int
	MaxHP      int
	Attack     int
	Speed      float64
	NextAttack float64
	Camp       int // 0: Left 1: Right
	Alive      bool
}

type Team struct {
	Name  string
	Units []*Unit
}

type Event struct {
	Time float64

	From string

	To string

	Damage int

	TargetHP int

	Dead bool
}

func (t *Team) AliveCount() int {

	count := 0

	for _, unit := range t.Units {
		if unit.Alive {
			count++
		}
	}

	return count
}

func (e Event) String() string {
	dead := ""
	if e.Dead {
		dead = "   DEAD"
	}

	return fmt.Sprintf(
		"[%7.2f] %-12s -> %-12s Damage:%3d TargetHP:%3d%s",
		e.Time,
		e.From,
		e.To,
		e.Damage,
		e.TargetHP,
		dead,
	)
}
