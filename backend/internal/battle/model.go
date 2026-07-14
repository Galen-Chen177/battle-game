package battle

import "fmt"

type Unit struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"maxHp"`
	Attack     int     `json:"attack"`
	Speed      float64 `json:"speed"`
	NextAttack float64 `json:"nextAttack"`
	Camp       int     `json:"camp"` // 0: Left 1: Right
	Alive      bool    `json:"alive"`
}

type Team struct {
	Name  string
	Units []*Unit
}

type Event struct {
	Time     float64 `json:"time"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Damage   int     `json:"damage"`
	TargetHP int     `json:"targetHp"`
	Dead     bool    `json:"dead"`
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
