package engine

import (
	"battle-game/internal/battle"
	"math"
	"math/rand"
)

type Engine struct {
	battle *battle.Battle
}

func New(b *battle.Battle) *Engine {
	return &Engine{
		battle: b,
	}
}

func (e *Engine) Run() {

	// 初始化每个人第一次攻击时间
	for _, u := range e.battle.Left.Units {
		u.NextAttack = 1 / u.Speed
	}

	for _, u := range e.battle.Right.Units {
		u.NextAttack = 1 / u.Speed
	}

	for {
		// 战斗结束
		if e.battle.Left.AliveCount() == 0 ||
			e.battle.Right.AliveCount() == 0 {
			return
		}

		// 时间到了，结束
		if e.battle.Time >= 999 {
			return
		}

		// 找下一位行动的人
		attacker := e.nextActor()

		// 已经死了
		if !attacker.Alive {
			attacker.NextAttack = math.MaxFloat64
			continue
		}

		// 更新时间
		e.battle.Time = attacker.NextAttack

		// 找敌人
		target := e.randomTarget(attacker.Camp)

		if target == nil {
			return
		}

		// 攻击
		target.HP -= attacker.Attack

		dead := false

		if target.HP <= 0 {
			target.HP = 0
			target.Alive = false
			dead = true
			target.NextAttack = math.MaxFloat64
		}

		// 记录事件
		e.battle.Events = append(e.battle.Events, battle.Event{
			Time:     e.battle.Time,
			From:     attacker.Name,
			To:       target.Name,
			Damage:   attacker.Attack,
			TargetHP: target.HP,
			Dead:     dead,
		})

		// 下一次攻击
		attacker.NextAttack += 1 / attacker.Speed
	}
}

func (e *Engine) nextActor() *battle.Unit {

	var actor *battle.Unit

	min := math.MaxFloat64

	all := []*battle.Unit{}

	all = append(all, e.battle.Left.Units...)
	all = append(all, e.battle.Right.Units...)

	for _, u := range all {

		if !u.Alive {
			continue
		}

		if u.NextAttack < min {
			min = u.NextAttack
			actor = u
		}
	}

	return actor
}

func (e *Engine) randomTarget(camp int) *battle.Unit {

	var enemies []*battle.Unit

	if camp == 0 {
		enemies = e.battle.Right.Units
	} else {
		enemies = e.battle.Left.Units
	}

	var alive []*battle.Unit

	for _, u := range enemies {

		if u.Alive {
			alive = append(alive, u)
		}
	}

	if len(alive) == 0 {
		return nil
	}

	return alive[rand.Intn(len(alive))]
}
