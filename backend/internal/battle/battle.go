package battle

type Battle struct {
	Left *Team

	Right *Team

	Time float64

	Events []Event
}

func NewDemoBattle() *Battle {

	left := &Team{
		Name: "Left",
		Units: []*Unit{
			{
				ID:     1,
				Name:   "张三",
				Camp:   0,
				MaxHP:  1000,
				HP:     1000,
				Attack: 13,
				Speed:  2.0,
				Alive:  true,
			},
			{
				ID:     2,
				Name:   "李四",
				Camp:   0,
				MaxHP:  1200,
				HP:     1200,
				Attack: 12,
				Speed:  1.5,
				Alive:  true,
			},
			{
				ID:     3,
				Name:   "小火龙",
				Camp:   0,
				MaxHP:  800,
				HP:     800,
				Attack: 20,
				Speed:  1.2,
				Alive:  true,
			},
		},
	}

	right := &Team{
		Name: "Right",
		Units: []*Unit{
			{
				ID:     101,
				Name:   "哥布林A",
				Camp:   1,
				MaxHP:  900,
				HP:     900,
				Attack: 10,
				Speed:  1.8,
				Alive:  true,
			},
			{
				ID:     102,
				Name:   "哥布林B",
				Camp:   1,
				MaxHP:  900,
				HP:     900,
				Attack: 10,
				Speed:  1.8,
				Alive:  true,
			},
			{
				ID:     103,
				Name:   "狼王",
				Camp:   1,
				MaxHP:  1500,
				HP:     1500,
				Attack: 18,
				Speed:  1.0,
				Alive:  true,
			},
		},
	}

	return &Battle{
		Left:  left,
		Right: right,
	}
}
