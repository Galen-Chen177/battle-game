# CLAUDE.md

本文件为 Claude Code（claude.ai/code）在此仓库中工作时提供指导。

## 约定

当项目发生设计决策、架构调整、新增/删除包、或重要的实现方式变更时，**需同步更新本文件**。如果你提醒我修改，我会改；如果我在做决策时发现应该记录，我会主动改。

命名约定：
- JSON tag 使用**驼峰命名**（如 `maxHp`、`nextAttack`、`targetHp`），不使用 snake_case
- 不需要执行编译命令（`go build`、`go run` 等），用户自己编译测试

## 常用命令

```bash
# 运行战斗引擎（打印 3v3 演示战斗到 stdout）
cd backend && go run ./cmd/

# 编译
cd backend && go build ./cmd/

# 运行所有测试（目前还没有 _test.go 文件）
cd backend && go test ./...

# 运行单个包的测试
cd backend && go test ./internal/battle/
cd backend && go test ./internal/engine/

# 运行单个测试函数
cd backend && go test ./internal/battle/ -run TestName

# 整理依赖
cd backend && go mod tidy
```

## 架构

```
backend/
  cmd/main.go           → 入口：构建演示战斗，运行引擎，打印事件
  internal/
    battle/             → 领域类型（Unit, Team, Event, Battle）—— 纯数据，无逻辑
    engine/             → 模拟循环（Engine.Run）—— 读取状态，产出 Event
```

**各包职责：**

- **`battle`** — 领域模型。
  - `Unit`：ID、Name、HP、MaxHP、Attack、Speed（攻击频率，次/秒）、NextAttack（下次行动时间）、Camp（0=左方 1=右方）、Alive
  - `Team`：队伍名 + 单位指针切片，方法 `AliveCount() int`
  - `Event`：Time、From、To、Damage、TargetHP、Dead，方法 `String() string` 用于格式化输出
  - `Battle`：Left/Right 两个队伍、当前 Time、累积的 Events 切片
  - `NewDemoBattle()`：硬编码 3v3 场景（张三/李四/小火龙 vs 哥布林A/哥布林B/狼王）

- **`engine`** — 核心模拟循环，`Engine.Run()`：
  1. 初始化所有单位的 `NextAttack = 1/Speed`（速度越快，首次行动越早）
  2. 进入循环：找出 `NextAttack` 最小的存活单位 → 推进时间 → 从敌方阵营随机选一个存活目标 → 造成伤害 → 记录 Event → 为该攻击者安排下一次行动时间（`NextAttack += 1/Speed`）
  3. 任一方全部死亡 或 `Time >= 999`（硬编码超时保护，后续会改为可配置）时结束

**数据流：** `Battle 状态 → Engine.Run() → []Event → 消费方（CLI，未来 React/Unity/Godot）`

## 设计原则

项目采用**事件驱动架构**，战斗逻辑与表现层严格分离：

1. **引擎只负责计算** — 它不是游戏，是确定性模拟器。整场战斗在几毫秒内完成，不使用 Sleep/Tick/Timer。
2. **Event 是唯一输出** — 所有动作（攻击、死亡、未来的治疗/Buff/召唤）都产出 Event。消费方按自己的节奏回放（1 倍速、2 倍速、跳过、回退），无需重新计算。
3. **引擎与前端无关** — 不引入 React、Gin、WebSocket、Unity 或任何 UI 库。同一引擎应能服务于 CLI、Web 和游戏引擎。
4. **Event 后续会升级为接口** — 当前是单一结构体，后续改为 `type Event interface { Time() float64; Type() string }`，再派生出 `AttackEvent`、`HealEvent`、`DeadEvent`、`BuffAddEvent` 等。
5. **所有规则逐步配置化** — 当前硬编码的演示单位和 `999` 超时都是临时的，后续版本通过 JSON 配置文件来定义单位、技能、怪物和 Buff（见 v0.8）。
6. **超时是正常结果，不是异常** — 无法结束的战斗（如治疗超过伤害）会产出 `BattleTimeoutEvent` 并返回有效的 `BattleResult`，不应报错。

## 版本路线图

当前处于 **v0.1**（基础 3v3 自动战斗）。后续规划：

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
