# CLAUDE.md

本文件为 Claude Code（claude.ai/code）在此仓库中工作时提供指导。

## 约定

当项目发生设计决策、架构调整、新增/删除包、或重要的实现方式变更时，**需同步更新本文件**。如果你提醒我修改，我会改；如果我在做决策时发现应该记录，我会主动改。

命名约定：
- JSON tag 使用**驼峰命名**（如 `maxHp`、`nextAttack`、`targetHp`），不使用 snake_case
- 不需要执行编译命令（`go build`、`go run` 等），用户自己编译测试

## 常用命令

```bash
# 启动后端 HTTP 服务
cd backend && go run ./cmd/

# 编译
cd backend && go build ./cmd/

# 运行所有测试
cd backend && go test ./...

# 运行单个包的测试
cd backend && go test ./internal/battle/
cd backend && go test ./internal/engine/

# 运行单个测试函数
cd backend && go test ./internal/battle/ -run TestNewDemoBattle

# 整理依赖
cd backend && go mod tidy
```

## 架构

```
backend/
  cmd/main.go             → 入口：Gin HTTP 服务，监听 :8080
  internal/
    battle/               → 领域类型（Unit, Team, Event, Battle）—— 纯数据，无逻辑
    engine/               → 模拟循环（Engine.Run）—— 读取状态，产出 Event
    handler/              → HTTP handler（POST /battle/start）
frontend/
  index.html              → 战斗演示页面（纯 HTML，浏览器直接打开）
```

**各包职责：**

- **`battle`** — 领域模型。
  - `Unit`：ID、Name、HP、MaxHP、Attack、Speed（攻击频率，次/秒）、NextAttack（下次行动时间）、Camp（0=我方 1=敌军）、Alive。所有字段有 JSON tag（驼峰）。
  - `Team`：队伍名 + 单位指针切片，方法 `AliveCount() int`
  - `Event`：Time、From、To、Damage、TargetHP、Dead，方法 `String() string`
  - `Battle`：Left/Right 两个队伍、当前 Time、累积的 Events 切片
  - `NewBattle(left, right)`：通用构造函数
  - `NewDemoBattle()`：硬编码 3v3 场景，已移至 `battle_test.go`

- **`engine`** — 核心模拟循环。
  - `Engine.Run()`：初始化 NextAttack → 循环找最小 NextAttack 的存活单位 → 调 `pickTarget()` 选目标 → 造成伤害 → 记录 Event → 安排下次行动。任一方全灭或 `Time >= 999` 时结束。
  - `pickTarget(camp)`：**Camp 1（敌军）始终随机攻击**；Camp 0（我方）根据 `PriorityAttackOptions` 选择目标。
  - `PriorityAttackOptions`：`0`=随机，`1`=优先 maxHP 最多，`2`=优先当前 HP 最多，`3`=优先当前 HP 最少。

- **`handler`** — HTTP 层。
  - `POST /battle/start`：入参 `{units: [...], priorityAttackOptions: 0}`，按 Camp 拆队 → 构建 Battle → 跑引擎 → 返回 `[]Event`。
  - 自动校验 camp 值、双方非空、归一化 Alive/MaxHP。

**数据流：** `HTTP JSON → Handler 拆队 → Battle 状态 → Engine.Run() → []Event → JSON 响应 → 前端播放`

## 设计原则

项目采用**事件驱动架构**，战斗逻辑与表现层严格分离：

1. **引擎只负责计算** — 它不是游戏，是确定性模拟器。整场战斗在几毫秒内完成，不使用 Sleep/Tick/Timer。
2. **Event 是唯一输出** — 所有动作都产出 Event。消费方按自己的节奏回放（1 倍速、2 倍速、跳过、回退），无需重新计算。
3. **引擎与前端无关** — 引擎本身不引入 Gin、React 或任何 UI 库。Gin 只在 handler 层和 main.go 中使用，引擎包不依赖 HTTP。
4. **Event 后续会升级为接口** — 当前是单一结构体，后续改为 `type Event interface { Time() float64; Type() string }`，派生出 `AttackEvent`、`HealEvent` 等。
5. **所有规则逐步配置化** — 当前硬编码的演示单位和 `999` 超时都是临时的，后续版本通过 JSON 配置文件来定义。
6. **超时是正常结果，不是异常** — 无法结束的战斗会产出 `BattleTimeoutEvent` 并返回有效的 `BattleResult`。

## API

### POST /battle/start

```json
// 请求
{
  "units": [
    {"id":1,"name":"Hero","hp":100,"maxHp":100,"attack":15,"speed":2.0,"camp":0},
    {"id":2,"name":"Monster","hp":80,"maxHp":80,"attack":10,"speed":1.5,"camp":1}
  ],
  "priorityAttackOptions": 0
}

// 响应 — []Event
[
  {"time":0.50,"from":"Hero","to":"Monster","damage":15,"targetHp":65,"dead":false},
  ...
]
```

`priorityAttackOptions`：`0`=随机，`1`=优先 maxHP 最多，`2`=优先当前 HP 最多，`3`=优先当前 HP 最少。

## 版本路线图

当前处于 **v0.1**（基础 3v3 自动战斗 + HTTP API + 前端演示）。后续规划：

| 版本 | 主题 | 核心内容 |
|------|------|----------|
| v0.2 | 单位状态 | HP 百分比、BattleResult、战斗统计、MVP |
| v0.3 | 攻击系统 | 暴击、随机伤害、命中/闪避、防御、Damage Calculator |
| v0.4 | 技能系统 | 主动技能、冷却、单体/AoE、技能 Event |
| v0.5 | Buff 系统 | 中毒、灼烧、护盾、加速/减速、持续治疗 |
| v0.6 | 宠物/召唤 | 宠物、召唤物、临时单位、死亡清理 |
| v0.7 | AI | 目标优先级（后排/低血量）、仇恨系统、技能策略 |
| v0.8 | 配置化 | JSON 配置角色、技能、怪物、Buff |
| v0.9 | 回放 | 战斗录像、Event 回放、导出/导入 |
| v1.0 | 前端 | React 页面、WebSocket、实时播放、血条、速度控制 |

设计决策文档：`docs/决策设计01.md`
