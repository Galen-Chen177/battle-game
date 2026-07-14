package battle_test

import (
	"testing"

	"battle-game/internal/battle"
	"battle-game/internal/engine"
)

func TestNewDemoBattle(t *testing.T) {
	b := NewDemoBattle()

	if b.Left == nil {
		t.Fatal("Left team is nil")
	}
	if b.Right == nil {
		t.Fatal("Right team is nil")
	}
	if len(b.Left.Units) == 0 {
		t.Fatal("Left team has no units")
	}
	if len(b.Right.Units) == 0 {
		t.Fatal("Right team has no units")
	}

	// 验证 camp
	for _, u := range b.Left.Units {
		if u.Camp != 0 {
			t.Errorf("Left unit %s has camp=%d, want 0", u.Name, u.Camp)
		}
	}
	for _, u := range b.Right.Units {
		if u.Camp != 1 {
			t.Errorf("Right unit %s has camp=%d, want 1", u.Name, u.Camp)
		}
	}

	// 跑战斗
	e := engine.New(b, engine.PriorityAttackRandom)
	e.Run()

	// 至少有一方全灭 或 超时
	leftAlive := b.Left.AliveCount()
	rightAlive := b.Right.AliveCount()
	if leftAlive > 0 && rightAlive > 0 && b.Time < 999 {
		t.Errorf("battle did not end: left=%d alive, right=%d alive, time=%.2f", leftAlive, rightAlive, b.Time)
	}

	// 有事件产生
	if len(b.Events) == 0 {
		t.Error("no events produced")
	}

	t.Logf("Battle finished: time=%.2f, events=%d, left alive=%d, right alive=%d",
		b.Time, len(b.Events), leftAlive, rightAlive)
}

func NewDemoBattle() *battle.Battle {

	left := &battle.Team{
		Name: "Left",
		Units: []*battle.Unit{
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

	right := &battle.Team{
		Name: "Right",
		Units: []*battle.Unit{
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

	return &battle.Battle{
		Left:  left,
		Right: right,
	}
}
