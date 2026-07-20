package battle

import "fmt"

// Unit 战斗单位
type Unit struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	HP         int     `json:"hp"`
	MaxHP      int     `json:"maxHp"`
	Attack     int     `json:"attack"`
	Speed      float64 `json:"speed"`
	NextAttack float64 `json:"nextAttack"`
	Camp       int     `json:"camp"` // 0: 我方 1: 敌方
	Alive      bool    `json:"alive"`

	CriticalHitRate   int `json:"criticalHitRate"`   // 暴击率
	AntiViolenceRate  int `json:"antiViolenceRate"`  // 免暴率
	CriticalHitDamage int `json:"criticalHitDamage"` // 暴击伤害，默认是100

	HitRate     int `json:"hitRate"`     // 命中率，默认是100
	EvasionRate int `json:"evasionRate"` // 闪避率
}

type Team struct {
	Name  string
	Units []*Unit
}

type Event struct {
	Time       float64    `json:"time"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	Damage     int        `json:"damage"`
	TargetHP   int        `json:"targetHp"`
	Dead       bool       `json:"dead"`
	DamageType DamageType `json:"damageType"` // 伤害类型，暴击，普通，毒之类的
}

type DamageType int

const (
	NormalDamage   DamageType = 0
	CriticalDamage DamageType = 1 // 暴击伤害
	EvasionDamage  DamageType = 2 // 闪避成功
)

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
