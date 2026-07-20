package engine

import (
	"battle-game/internal/battle"
	"math"
	"math/rand"
)

// PriorityAttackOptions 优先攻击选项
type PriorityAttackOptions uint

const (
	// 随机攻击
	PriorityAttackRandom PriorityAttackOptions = 0
	// 优先攻击 maxHP 最多的
	PriorityAttackMaxHP PriorityAttackOptions = 1
	// 优先攻击当前 HP 最多的
	PriorityAttackHP PriorityAttackOptions = 2
	// 优先攻击当前 HP 最少的
	PriorityAttackMinHP PriorityAttackOptions = 3
)

type Engine struct {
	battle                *battle.Battle
	priorityAttackOptions PriorityAttackOptions
}

func New(b *battle.Battle, opts PriorityAttackOptions) *Engine {
	return &Engine{
		battle:                b,
		priorityAttackOptions: opts,
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
		target := e.pickTarget(attacker.Camp)

		if target == nil {
			return
		}

		// 计算伤害

		damageType := battle.NormalDamage
		damage := attacker.Attack

		// 是否命中
		if attacker.HitRate-target.EvasionRate > 0 &&
			attacker.HitRate-target.EvasionRate-rand.Intn(100) > 0 {
			// 是否暴击
			if attacker.CriticalHitRate-target.AntiViolenceRate > 0 &&
				(attacker.CriticalHitRate-target.AntiViolenceRate)-rand.Intn(100) > 0 {
				// 暴击伤害
				damage = attacker.Attack * attacker.CriticalHitDamage / 100
				damageType = battle.CriticalDamage
			}
		} else {
			damage = 0
			damageType = battle.EvasionDamage
		}

		target.HP -= damage

		dead := false

		if target.HP <= 0 {
			target.HP = 0
			target.Alive = false
			dead = true
			target.NextAttack = math.MaxFloat64
		}

		// 记录事件
		e.battle.Events = append(e.battle.Events, battle.Event{
			Time:       e.battle.Time,
			From:       attacker.Name,
			To:         target.Name,
			Damage:     damage,
			TargetHP:   target.HP,
			Dead:       dead,
			DamageType: damageType,
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

// pickTarget 根据优先级选出攻击的目标
func (e *Engine) pickTarget(camp int) *battle.Unit {

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

	// 敌军默认随机
	if camp == 1 {
		return alive[rand.Intn(len(alive))]
	}

	// 我方根据选择来
	switch e.priorityAttackOptions {
	case PriorityAttackMaxHP:
		target := alive[0]
		best := alive[0].MaxHP
		for _, u := range alive[1:] {
			if u.MaxHP > best {
				best = u.MaxHP
				target = u
			}
		}
		return target

	case PriorityAttackHP:
		target := alive[0]
		best := alive[0].HP
		for _, u := range alive[1:] {
			if u.HP > best {
				best = u.HP
				target = u
			}
		}
		return target

	case PriorityAttackMinHP:
		target := alive[0]
		best := alive[0].HP
		for _, u := range alive[1:] {
			if u.HP < best {
				best = u.HP
				target = u
			}
		}
		return target

	default:
		return alive[rand.Intn(len(alive))]
	}
}
