# Battle Engine

一个使用 **Go** 编写的轻量级战斗引擎（Battle Engine）。

项目目标不是开发一个完整的游戏，而是实现一个**可扩展、可复用、可配置**的战斗模拟核心。

未来可以作为：

* 网页游戏战斗核心
* 放置类游戏
* RPG 战斗系统
* 宠物养成游戏
* 自动战斗模拟器
* AI 对战模拟器

甚至可以作为 Unity、Godot、React 等前端项目的后端战斗引擎。

---

# 项目目标

整个项目遵循几个原则：

* **战斗逻辑与表现层完全分离**
* **所有战斗结果都通过 Event 输出**
* **前端永远不参与战斗计算**
* **所有规则均可逐步配置化**
* **先保证简单，再逐步扩展**

项目最终希望做到：

```
React
Unity
Godot
CLI
        │
        ▼
   Battle Engine
        │
        ▼
 Battle Events
```

Battle Engine 不关心：

* React
* WebSocket
* 数据库
* UI

它只负责：

* 谁行动
* 打谁
* 造成多少伤害
* 是否死亡
* 战斗是否结束

---

# 当前功能（v0.1）

目前已经支持：

* 3v3 自动战斗
* 每个角色拥有基础属性
  * HP
  * Attack
  * Speed
* 根据 Speed 决定攻击频率
* 优先攻击策略（随机 / maxHP 最多 / 当前 HP 最多 / 当前 HP 最少）
* 敌方默认随机攻击，我方可选优先策略
* 角色死亡后退出战斗
* 战斗结束判定
* 每一次攻击都会生成 Event
* HTTP API：`POST /battle/start`（Gin 框架）
* 前端演示页面：战斗动画、伤害数字、计时器、变速播放

目前还没有：

* 技能
* Buff
* 暴击
* 闪避
* 护盾
* 治疗
* 宠物
* 装备
* 属性克制

这些都会在后续版本逐步加入。

---

# 当前目录

```
battle-game/
├── README.md
├── CLAUDE.md
├── docs/
│   └── 决策设计01.md            # 设计决策：事件驱动架构
├── backend/
│   ├── cmd/
│   │   └── main.go              # 入口：Gin HTTP 服务 (:8080)
│   ├── internal/
│   │   ├── battle/
│   │   │   ├── battle.go        # Battle 结构体、NewBattle()、NewDemoBattle()
│   │   │   ├── model.go         # Unit、Team、Event 领域类型
│   │   │   └── battle_test.go   # Demo 战斗测试
│   │   ├── engine/
│   │   │   └── engine.go        # Engine 模拟循环、优先攻击目标选择
│   │   └── handler/
│   │       └── battle.go        # POST /battle/start handler
│   ├── go.mod
│   └── go.sum
└── frontend/
    └── index.html               # 战斗演示页面（浏览器直接打开）
```

---

# 项目设计思想

Battle Engine 只有一个职责：

> 根据当前战场状态，计算下一步发生什么。

所有结果统一输出 Event，例如：

```
0.50 张三 -> 哥布林A  Damage:15

0.82 狼王 -> 李四 Damage:20

1.30 张三 -> 哥布林A Damage:15 Dead
```

以后：

普通攻击 → 治疗 → Buff → 召唤 → 反击 → 中毒

全部都会转换成 Event。

因此：任何前端都可以直接播放 Event。

---

# 运行方式

**启动后端：**

```bash
cd backend
go run ./cmd/
# 服务运行在 http://localhost:8080
```

**打开前端：**

浏览器直接打开 `frontend/index.html`，点击"开始战斗"即可。

**API 调用示例：**

```bash
curl -X POST http://localhost:8080/battle/start \
  -H 'Content-Type: application/json' \
  -d '{
    "units": [
      {"name":"Hero","hp":100,"maxHp":100,"attack":15,"speed":2.0,"camp":0},
      {"name":"Monster","hp":80,"maxHp":80,"attack":10,"speed":1.5,"camp":1}
    ],
    "priorityAttackOptions": 0
  }'
```

`priorityAttackOptions`：`0`=随机，`1`=优先 maxHP 最多，`2`=优先当前 HP 最多，`3`=优先当前 HP 最少。

---

# Roadmap

## v0.1

基础战斗（当前版本）

* [x] 3v3
* [x] 自动攻击
* [x] Speed 调度
* [x] Event
* [x] 战斗结束
* [x] HTTP API
* [x] 优先攻击策略
* [x] 前端演示页面

---

## v0.2

角色状态优化

计划加入：

* [ ] HP 百分比
* [ ] Alive() 方法
* [ ] BattleResult
* [ ] 战斗统计
* [ ] MVP 输出

目标：让 Battle Engine 不只是能战斗，还能输出完整战斗结果。

---

## v0.3

攻击系统

计划加入：

* [ ] 暴击
* [ ] 随机伤害
* [ ] 命中
* [ ] 闪避
* [ ] 防御
* [ ] Damage Calculator

此版本开始，将伤害计算独立出来。

---

## v0.4

技能系统（第一版）

计划加入：

* [ ] 主动技能
* [ ] 冷却时间
* [ ] 单体技能
* [ ] 群体技能
* [ ] 技能 Event

目标：第一次拥有真正意义上的技能。

---

## v0.5

Buff 系统

计划加入：

* [ ] 中毒
* [ ] 灼烧
* [ ] 护盾
* [ ] 加速
* [ ] 减速
* [ ] 持续治疗

Buff 将成为一个独立系统。

---

## v0.6

宠物 / 召唤物

计划加入：

* [ ] 宠物
* [ ] 召唤物
* [ ] 临时单位
* [ ] 死亡消失

战场单位不再固定为 3v3。

---

## v0.7

AI

计划加入：

* [ ] 优先攻击后排
* [ ] 仇恨系统
* [ ] 技能释放策略

不同角色可以拥有不同 AI。

---

## v0.8

配置化

计划加入：

* [ ] JSON 配置角色
* [ ] JSON 配置技能
* [ ] JSON 配置怪物
* [ ] JSON 配置 Buff

之后新增角色无需修改 Go 代码。

---

## v0.9

Replay

计划加入：

* [ ] 战斗录像
* [ ] Event 回放
* [ ] 战斗导出
* [ ] 战斗重播

Battle Engine 将支持完整回放。

---

## v1.0

前端

计划加入：

* React 页面
* WebSocket
* 实时播放
* 血条
* 战斗日志
* 战斗速度调节
* 暂停 / 继续

完成后即可拥有一个真正可以玩的自动战斗模拟器。

---

# 长期目标

Battle Engine 希望最终成为一个可以独立维护的 Go 战斗框架。

未来支持：CLI / Web / Unity / Godot / 桌面客户端。

任何前端都可以直接复用 Battle Engine。

整个项目始终坚持：

> Battle Engine 只负责计算，不负责显示。
