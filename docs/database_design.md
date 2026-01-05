# Text WoW 数据库设计文档

> 📌 **说明**: 本文档描述数据库设计和数值设计规范。详细的系统设计请参考 [系统架构文档](./architecture.md) 和对应的系统设计文档。

---

## 🔢 数值设计规范

> 📌 **核心理念**: 所有数值保持小巧直观，避免"数值膨胀"，玩家可以轻松理解和比较

### 数值设计原则

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          数值设计哲学                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ❌ 避免                              ✅ 追求                               │
│  ├─ 成千上万的伤害数字                ├─ 一眼能看懂的小数字                 │
│  ├─ 后期数值爆炸                      ├─ 全程保持相对稳定                   │
│  ├─ 需要计算器的对比                  ├─ 心算即可判断优劣                   │
│  └─ 意义不明的大数字                  └─ 每个数字都有明确意义               │
│                                                                             │
│  示例对比:                                                                   │
│  ❌ "你造成了 1,234,567 点伤害"    vs    ✅ "你造成了 45 点伤害"           │
│  ❌ "获得 50,000 经验"             vs    ✅ "获得 25 经验"                 │
│  ❌ "HP: 125,000 / 125,000"        vs    ✅ "HP: 85 / 85"                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 数值范围规范

| 类别 | 数值范围 | 说明 |
|-----|---------|------|
| **生命值 (HP)** | 20 ~ 500 | 1级约20-30，满级约200-400，Boss最高500 |
| **法力值 (MP)** | 10 ~ 200 | 法系职业较高，战士无MP |
| **基础属性** | 5 ~ 50 | 力量/敏捷/智力/耐力/精神，初始5-15 |
| **攻击力** | 3 ~ 80 | 包含武器加成，满级约40-80 |
| **防御力** | 1 ~ 50 | 满级坦克约30-50 |
| **技能伤害** | 5 ~ 60 | 普攻5-15，强力技能30-60 |
| **治疗量** | 5 ~ 50 | 与伤害相近 |
| **经验值** | 5 ~ 100 | 每只怪5-30，Boss约100 |
| **金币** | 1 ~ 50 | 普通怪1-5，Boss约20-50 |
| **暴击/闪避率** | 5% ~ 40% | 上限控制在40% |
| **减伤率** | 0% ~ 60% | 护甲减伤上限60% |

### 等级成长曲线

> 💡 采用**平缓线性成长**而非指数成长，保证后期数值不爆炸

```
等级1:   HP≈25   攻击≈8   防御≈3
等级10:  HP≈50   攻击≈15  防御≈8
等级20:  HP≈80   攻击≈25  防御≈12
等级30:  HP≈110  攻击≈35  防御≈18
等级40:  HP≈150  攻击≈45  防御≈24
等级50:  HP≈190  攻击≈55  防御≈30
等级60:  HP≈230  攻击≈65  防御≈36

每级成长: HP +3~4, 攻击 +1, 防御 +0.5
```

### 战斗数值示例

```
╔═══════════════════════════════════════════════════════════════════════════╗
║                     示例战斗 (30级队伍 vs 精英怪)                          ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                           ║
║  【我方队伍】                                                              ║
║  战士 Lv30: HP 120/120, 攻击32, 防御22                                    ║
║  法师 Lv30: HP 65/65,   MP 80/80,  法伤28                                 ║
║  牧师 Lv30: HP 70/70,   MP 65/65,  治疗25                                 ║
║                                                                           ║
║  【敌方】                                                                  ║
║  精英狼人 Lv32: HP 180/180, 攻击38, 防御15                                ║
║                                                                           ║
║  ─────────────────────────────────────────────────────────────────────── ║
║                                                                           ║
║  回合1:                                                                    ║
║  ├─ 战士使用[英勇打击] → 狼人受到 18 点伤害 (HP: 162)                     ║
║  ├─ 法师使用[火球术] → 狼人受到 24 点伤害 (HP: 138) ★暴击                ║
║  ├─ 牧师使用[治疗术] → 战士恢复 12 点生命                                 ║
║  └─ 狼人攻击战士 → 战士受到 14 点伤害 (HP: 106)                           ║
║                                                                           ║
║  战斗持续约 8 回合                                                         ║
║  胜利！获得 35 经验，8 金币                                                ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 经验升级表

| 等级 | 升级所需经验 | 累计经验 | 说明 |
|-----|------------|---------|------|
| 1→2 | 20 | 20 | 约4只怪 |
| 5→6 | 40 | 150 | |
| 10→11 | 60 | 400 | |
| 20→21 | 100 | 1,200 | |
| 30→31 | 140 | 2,600 | |
| 40→41 | 180 | 4,400 | |
| 50→51 | 220 | 6,600 | |
| 59→60 | 280 | 9,200 | 满级 |

> 📝 升级公式: `所需经验 = 20 + (等级 × 4)`，线性增长而非指数

### 装备数值规范

| 装备品质 | 属性加成范围 | 示例 |
|---------|------------|------|
| 普通(白) | +1~3 | 破旧之剑: 攻击+2 |
| 优秀(绿) | +3~6 | 民兵之剑: 攻击+4, 力量+2 |
| 精良(蓝) | +5~10 | 迪菲亚军刀: 攻击+7, 敏捷+3 |
| 史诗(紫) | +8~15 | 黑龙之牙: 攻击+12, 暴击+3% |
| 传说(橙) | +12~20 | 雷霆之怒: 攻击+18, 攻速+10% |

### 金币经济规范

| 阶段 | 怪物掉落 | 商店物品 | 说明 |
|-----|---------|---------|------|
| 初期(1-10级) | 1~3金 | 药水5金，武器20金 | 约10怪买1瓶药 |
| 中期(11-30级) | 3~8金 | 药水15金，武器100金 | 保持购买力 |
| 后期(31-60级) | 5~15金 | 药水30金，武器300金 | 经济稳定 |

### 数值平衡检查表

每次添加新内容时，检查以下项目：

- [ ] HP不超过500
- [ ] 单次伤害不超过80
- [ ] 经验奖励不超过100
- [ ] 金币掉落不超过50
- [ ] 属性加成不超过+20
- [ ] 百分比加成不超过+30%
- [ ] 技能系数不超过2.0

---

## 📊 数据库概览

### 🎮 核心设计：5人小队系统

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         用户小队结构                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│    ┌─────────┐                                                          │
│    │  USER   │ ─────────────────────────────────────────┐               │
│    └─────────┘                                          │               │
│         │                                               │               │
│         │ 1:N (最多5个)                                 │               │
│         ▼                                               ▼               │
│    ┌─────────────────────────────────────────────────────────────┐     │
│    │                        小队 (Team)                           │     │
│    │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐    │     │
│    │  │ Slot 1 │ │ Slot 2 │ │ Slot 3 │ │ Slot 4 │ │ Slot 5 │    │     │
│    │  │ 队长   │ │ 成员   │ │ 成员   │ │ 成员   │ │ 成员   │    │     │
│    │  │ 战士   │ │ 法师   │ │ 牧师   │ │ (空)   │ │ (空)   │    │     │
│    │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘    │     │
│    └─────────────────────────────────────────────────────────────┘     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**小队规则：**
- 初始只有1个角色槽位，需解锁更多
- 最多可解锁5个角色槽位
- 角色通过 `team_slot` (1-5) 确定位置
- `is_active` 控制是否参与战斗
- 金币小队共享，经验平均分配

**槽位解锁条件（建议）：**
| 槽位 | 解锁条件 |
|-----|---------|
| 1 | 初始拥有 |
| 2 | 队伍中任意角色达到 10 级 |
| 3 | 队伍中任意角色达到 20 级 |
| 4 | 队伍中任意角色达到 35 级 |
| 5 | 队伍中任意角色达到 50 级 |

**死亡与复活机制：**
- 角色死亡后需要等待复活
- 复活时间 = 基础时间 × 当前死亡角色数量
- 玩家需通过策略配置尽量避免角色死亡

| 死亡人数 | 复活等待时间(建议) |
|---------|------------------|
| 1人死亡 | 30秒 |
| 2人死亡 | 60秒 (每人) |
| 3人死亡 | 90秒 (每人) |
| 4人死亡 | 120秒 (每人) |
| 5人全灭 | 180秒 (每人) |

> 💡 这个机制鼓励玩家合理配置策略（如低血量时使用治疗技能），而不是无脑挂机

### 📐 数据库架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         TEXT WoW 数据库架构                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐ 1:N ┌─────────────┐     ┌─────────────┐                   │
│  │  users  │────►│ characters  │────►│  inventory  │                   │
│  └─────────┘     │ (最多5个)   │     └─────────────┘                   │
│                  └──────┬──────┘                                        │
│                         │                                               │
│         ┌───────────────┼───────────────┐                               │
│         ▼               ▼               ▼                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                       │
│  │  equipment  │ │char_skills  │ │ strategies  │                       │
│  └─────────────┘ └─────────────┘ └─────────────┘                       │
│                                                                         │
│  ════════════════════ 游戏配置表（只读）════════════════════            │
│                                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │
│  │  races  │  │ classes │  │  items  │  │ skills  │  │  zones  │      │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └────┬────┘      │
│                                                           │            │
│                                                    ┌──────▼──────┐     │
│                                                    │  monsters   │     │
│                                                    └─────────────┘     │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## ⚔️ 属性与战斗系统

> 📌 **透明化设计**: 所有计算公式对玩家完全公开，支持游戏内查询，便于研究战术和培养角色

### 🎯 五大基础属性

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           基础属性总览                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │
│  │ 💪力量  │  │ 🏃敏捷   │  │ 🧠智力  │  │ ❤️耐力  │  │ ✨精神  │      │
│  │Strength │  │ Agility │  │Intellect│  │ Stamina │  │ Spirit  │      │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘      │
│       │            │            │            │            │            │
│       ▼            ▼            ▼            ▼            ▼            │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │
│  │物理攻击 │  │暴击率   │  │法术攻击 │  │最大HP   │  │法力回复 │      │
│  │格挡值   │  │闪避率   │  │最大MP   │  │物理防御 │  │生命回复 │      │
│  │暴击伤害 │  │物理攻击 │  │法术暴击 │  │         │  │治疗效果 │      │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

> 📌 **说明**: 属性转换公式和职业属性优先级已移至 [角色属性系统设计文档](./character_attributes_design.md)

详细的属性转换公式、职业属性优先级等内容请参考角色属性设计文档。

---

### 🔧 游戏规则查询系统

> 📌 **设计理念**: 所有战斗规则对玩家透明，可通过游戏内命令或界面查询

#### game_formulas - 公式配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(64) | PRIMARY KEY | 公式ID |
| category | VARCHAR(32) | NOT NULL | 分类: attribute/combat/skill/resource |
| name | VARCHAR(64) | NOT NULL | 公式名称 |
| formula | TEXT | NOT NULL | 公式表达式 |
| description | TEXT | | 详细说明 |
| variables | TEXT | | 变量说明(JSON) |
| example | TEXT | | 计算示例 |
| display_order | INTEGER | DEFAULT 0 | 显示顺序 |

```sql
CREATE TABLE game_formulas (
    id VARCHAR(64) PRIMARY KEY,
    category VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    formula TEXT NOT NULL,
    description TEXT,
    variables TEXT,
    example TEXT,
    display_order INTEGER DEFAULT 0
);

CREATE INDEX idx_formulas_category ON game_formulas(category);
```

#### 预置公式数据

```sql
INSERT INTO game_formulas (id, category, name, formula, description, variables, example, display_order) VALUES
-- 属性转换
('attr_max_hp', 'attribute', '最大生命值', 
 'class_base_hp + stamina × 10 + level × hp_per_level + equipment_bonus',
 '计算角色的最大生命值',
 '{"class_base_hp":"职业基础HP","stamina":"耐力值","level":"角色等级","hp_per_level":"每级HP成长","equipment_bonus":"装备加成"}',
 '战士(30级,耐力80): 120 + 800 + 360 + 0 = 1280 HP', 1),

('attr_phys_attack', 'attribute', '物理攻击力',
 'base_attack + strength × 2.0 + agility × 0.5 + weapon_damage + equipment_bonus',
 '计算角色的物理攻击力',
 '{"base_attack":"职业基础攻击","strength":"力量","agility":"敏捷","weapon_damage":"武器伤害"}',
 '盗贼(力量40,敏捷100): 12 + 80 + 50 = 142', 2),

('attr_spell_power', 'attribute', '法术攻击力',
 'base_spell + intellect × 1.5 + equipment_bonus',
 '计算角色的法术攻击力',
 '{"base_spell":"职业基础法伤","intellect":"智力"}',
 '法师(智力120): 20 + 180 = 200', 3),

('attr_armor', 'attribute', '护甲值',
 'base_armor + stamina × 0.5 + agility × 0.3 + equipment_armor',
 '计算角色的护甲值',
 '{"base_armor":"职业基础护甲","stamina":"耐力","agility":"敏捷"}',
 '战士(耐力80,敏捷40): 50 + 40 + 12 = 102', 4),

-- 战斗判定
('combat_phys_crit', 'combat', '物理暴击率',
 '5% + agility ÷ 20 + equipment_bonus (上限50%)',
 '计算物理攻击的暴击概率',
 '{"agility":"敏捷值"}',
 '敏捷100: 5% + 5% = 10%', 10),

('combat_spell_crit', 'combat', '法术暴击率',
 '5% + intellect ÷ 30 + equipment_bonus (上限50%)',
 '计算法术攻击的暴击概率',
 '{"intellect":"智力值"}',
 '智力120: 5% + 4% = 9%', 11),

('combat_dodge', 'combat', '闪避率',
 '5% + agility ÷ 25 + equipment_bonus + racial_bonus (上限30%)',
 '计算物理攻击的闪避概率',
 '{"agility":"敏捷值","racial_bonus":"种族加成(暗夜精灵+2%)"}',
 '暗夜精灵(敏捷100): 5% + 4% + 2% = 11%', 12),

('combat_armor_reduction', 'combat', '护甲减伤',
 'armor ÷ (armor + 400 + attacker_level × 10) × 100% (上限75%)',
 '根据护甲值计算物理伤害减免',
 '{"armor":"护甲值","attacker_level":"攻击者等级"}',
 '护甲500,被30级怪攻击: 500÷1200 = 41.7%减伤', 13),

('combat_crit_damage', 'combat', '暴击伤害倍率',
 '物理: 150% + strength ÷ 100 × 10% (上限250%); 法术: 150% + equipment_bonus',
 '暴击时的伤害倍率',
 '{"strength":"力量(物理暴击伤害加成)"}',
 '力量150: 150% + 15% = 165%', 14),

-- 技能伤害
('skill_damage', 'skill', '技能伤害计算',
 '(base_value + scaling_stat × scaling_ratio) × skill_level_mult',
 '计算技能的最终伤害值',
 '{"base_value":"技能基础值","scaling_stat":"成长属性值","scaling_ratio":"成长系数","skill_level_mult":"1+(技能等级-1)×0.1"}',
 '5级火球术(基础50,智力120,系数0.8): (50+96)×1.4 = 204', 20),

('skill_final_damage', 'skill', '最终伤害',
 'skill_damage × crit_mult × (1 - target_reduction) × random(0.9~1.1)',
 '计算造成的实际伤害',
 '{"skill_damage":"技能伤害","crit_mult":"暴击倍率(非暴击=1)","target_reduction":"目标减伤%"}',
 '技能204,暴击1.5倍,目标20%减伤: 204×1.5×0.8 = 244.8', 21),

-- 资源回复
('resource_mana_regen', 'resource', '法力回复(战斗中)',
 'spirit × class_regen_pct × max_mp ÷ 100 (每回合)',
 '战斗中每回合的法力恢复量',
 '{"spirit":"精神","class_regen_pct":"职业恢复系数(牧师0.8%,法师0.5%等)","max_mp":"最大法力"}',
 '牧师(精神80,MP500): 80×0.008×500 = 3.2/回合', 30),

('resource_rage_gain', 'resource', '怒气获得(战士)',
 '攻击命中+5, 暴击额外+10, 受伤+(damage÷max_hp×20)',
 '战士获得怒气的方式',
 '{"damage":"受到伤害","max_hp":"最大生命"}',
 '受到200伤害(最大HP1000): +4怒气', 31),

('resource_energy_regen', 'resource', '能量恢复(盗贼)',
 '每回合固定+20能量',
 '盗贼的能量恢复方式',
 '{}',
 '每回合恢复20,5回合恢复满', 32),

('resource_hp_regen', 'resource', '生命回复',
 '战斗中: spirit × 0.2/回合; 战斗外: spirit × 1.0 + max_hp × 1%/秒',
 '角色生命值的自然恢复',
 '{"spirit":"精神","max_hp":"最大生命"}',
 '精神50,最大HP1000: 战斗外 50+10 = 60/秒', 33);
```

#### 游戏内查询命令

| 命令 | 说明 | 示例 |
|-----|------|-----|
| `/help formulas` | 查看所有公式分类 | 显示属性/战斗/技能/资源分类 |
| `/help formulas [分类]` | 查看某分类的所有公式 | `/help formulas combat` |
| `/help formula [公式ID]` | 查看具体公式详情 | `/help formula combat_phys_crit` |
| `/calc [公式ID] [参数...]` | 计算具体数值 | `/calc attr_max_hp stamina=80 level=30` |
| `/stats` | 查看当前角色详细属性 | 显示所有属性及其来源分解 |
| `/compare [角色1] [角色2]` | 对比两个角色属性 | 对比队伍中不同角色 |

---

## 🎯 角色成长系统

> 📌 **核心机制**: 每升1级获得1个属性点自由分配，每升3级从3个技能中选择1个

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          角色成长系统                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  📈 升级时获得:                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │   每升 1 级  →  +1 属性点 (可自由分配到五大属性)                     │   │
│  │                                                                     │   │
│  │   每升 3 级  →  从 3 个技能中选择 1 个 (含主动/被动技能)             │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  满级 60 级时:                                                               │
│  ├─ 累计可分配属性点: 59 点                                                 │
│  └─ 累计可选技能次数: 20 次 (Lv3, 6, 9, 12...60)                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 🔢 属性点分配系统

#### 分配规则

| 规则 | 说明 |
|-----|------|
| 获得方式 | 每升1级获得1点，满级共59点 |
| 分配目标 | 力量/敏捷/智力/耐力/精神 |
| 分配时机 | 升级后可随时分配，无时间限制 |
| 重置机制 | 消耗金币可重置（费用随等级增加） |

#### 属性点上限（防止极端堆叠）

| 限制类型 | 数值 | 说明 |
|---------|-----|------|
| 单项软上限 | 40点 | 超过后每点效果减半 |
| 单项硬上限 | 50点 | 无法继续分配 |
| 总分配上限 | 59点 | 等于总升级次数 |

#### 分配策略示例

```
╔═══════════════════════════════════════════════════════════════════════════╗
║                    30级战士 - 两种培养流派                                  ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                           ║
║  【输出流】分配29点:                   【坦克流】分配29点:                 ║
║  ├─ 力量: +18 (攻击力优先)             ├─ 力量: +8  (保证基础输出)        ║
║  ├─ 敏捷: +6  (暴击加成)               ├─ 敏捷: +3  (少量闪避)            ║
║  ├─ 智力: +0                          ├─ 智力: +0                        ║
║  ├─ 耐力: +5  (基础生存)               ├─ 耐力: +15 (生存优先)            ║
║  └─ 精神: +0                          └─ 精神: +3  (回复)                ║
║                                                                           ║
║  最终攻击: ~45 (高)                    最终攻击: ~30 (中)                  ║
║  最终HP: ~85 (中)                      最终HP: ~120 (高)                   ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

#### character_stat_allocation - 属性分配表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| character_id | INTEGER | PRIMARY KEY FK | 角色ID |
| unspent_points | INTEGER | DEFAULT 0 | 未分配点数 |
| allocated_strength | INTEGER | DEFAULT 0 | 已分配力量 |
| allocated_agility | INTEGER | DEFAULT 0 | 已分配敏捷 |
| allocated_intellect | INTEGER | DEFAULT 0 | 已分配智力 |
| allocated_stamina | INTEGER | DEFAULT 0 | 已分配耐力 |
| allocated_spirit | INTEGER | DEFAULT 0 | 已分配精神 |
| respec_count | INTEGER | DEFAULT 0 | 重置次数 |
| last_respec_at | DATETIME | | 上次重置时间 |

```sql
CREATE TABLE character_stat_allocation (
    character_id INTEGER PRIMARY KEY,
    unspent_points INTEGER DEFAULT 0,
    allocated_strength INTEGER DEFAULT 0,
    allocated_agility INTEGER DEFAULT 0,
    allocated_intellect INTEGER DEFAULT 0,
    allocated_stamina INTEGER DEFAULT 0,
    allocated_spirit INTEGER DEFAULT 0,
    respec_count INTEGER DEFAULT 0,
    last_respec_at DATETIME,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);
```

---

### ⚔️ 技能选择系统

#### 选择规则

| 规则 | 说明 |
|-----|------|
| 触发时机 | 每升3级触发一次 (Lv3, 6, 9, 12...60) |
| 选项数量 | 每次提供3个技能选项 |
| 技能类型 | 主动技能 + 被动技能 混合 |
| 选择限制 | 必须选择1个，选后不可更改 |

#### 技能池分层

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           技能池分层设计                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  【基础层】Lv3-15 (5次选择)                                                  │
│  ├─ 核心攻击技能 (必出1个)                                                  │
│  ├─ 基础增益/减益                                                           │
│  └─ 入门被动技能                                                            │
│                                                                             │
│  【进阶层】Lv18-36 (7次选择)                                                 │
│  ├─ 特化技能 (AOE/单体/控制)                                                │
│  ├─ 职业专属强化                                                            │
│  └─ 进阶被动技能                                                            │
│                                                                             │
│  【大师层】Lv39-60 (8次选择)                                                 │
│  ├─ 终极技能 (大招)                                                         │
│  ├─ 稀有被动                                                                │
│  └─ 跨职业通用技能                                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 技能选择机制

| 机制 | 说明 |
|-----|------|
| **重复技能** | 选到已有技能 → 技能升级(+20%效果) |
| **保底机制** | 每次至少1个未拥有技能 |
| **稀有度权重** | 普通60% / 精良30% / 史诗10% |
| **职业限制** | 70%职业专属 + 30%通用技能 |

#### 被动技能示例

| 被动技能 | 效果 | 适合职业 | 稀有度 |
|---------|------|---------|-------|
| 利刃专精 | 物理伤害+8% | 战士/盗贼 | 普通 |
| 护甲掌握 | 护甲值+15% | 战士/圣骑 | 普通 |
| 法力涌流 | 法力回复+20% | 法系职业 | 普通 |
| 致命一击 | 暴击伤害+12% | 全职业 | 精良 |
| 生命汲取 | 伤害的3%转化为HP | 术士 | 精良 |
| 格挡专精 | 格挡几率+10% | 坦克 | 精良 |
| 嗜血本能 | HP<30%时攻击+25% | 战士/盗贼 | 史诗 |
| 魔法屏障 | 法术伤害-15% | 法系职业 | 史诗 |
| 不灭意志 | 首次致死伤害免疫 | 全职业 | 史诗 |

#### skill_selection_history - 技能选择记录表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| level_milestone | INTEGER | NOT NULL | 选择时等级(3,6,9...) |
| offered_skill_1 | VARCHAR(32) | NOT NULL | 选项1技能ID |
| offered_skill_2 | VARCHAR(32) | NOT NULL | 选项2技能ID |
| offered_skill_3 | VARCHAR(32) | NOT NULL | 选项3技能ID |
| selected_skill_id | VARCHAR(32) | NOT NULL | 选中的技能ID |
| skill_was_upgrade | INTEGER | DEFAULT 0 | 是否为技能升级 |
| selected_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 选择时间 |

```sql
CREATE TABLE skill_selection_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    level_milestone INTEGER NOT NULL,
    offered_skill_1 VARCHAR(32) NOT NULL,
    offered_skill_2 VARCHAR(32) NOT NULL,
    offered_skill_3 VARCHAR(32) NOT NULL,
    selected_skill_id VARCHAR(32) NOT NULL,
    skill_was_upgrade INTEGER DEFAULT 0,
    selected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (offered_skill_1) REFERENCES skills(id),
    FOREIGN KEY (offered_skill_2) REFERENCES skills(id),
    FOREIGN KEY (offered_skill_3) REFERENCES skills(id),
    FOREIGN KEY (selected_skill_id) REFERENCES skills(id),
    UNIQUE(character_id, level_milestone)
);

CREATE INDEX idx_skill_selection_char ON skill_selection_history(character_id);
```

#### passive_skills - 被动技能配置表

> 📌 被动技能独立于主动技能，无冷却无消耗，永久生效

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 被动ID |
| name | VARCHAR(32) | NOT NULL | 名称 |
| description | TEXT | | 描述 |
| icon | VARCHAR(64) | | 图标 |
| class_id | VARCHAR(32) | | 限定职业(NULL=通用) |
| rarity | VARCHAR(16) | DEFAULT 'common' | 稀有度: common/rare/epic |
| tier | INTEGER | DEFAULT 1 | 层级: 1基础/2进阶/3大师 |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_value | REAL | NOT NULL | 效果数值 |
| effect_stat | VARCHAR(32) | | 影响的属性 |
| max_level | INTEGER | DEFAULT 5 | 最大升级次数 |
| level_scaling | REAL | DEFAULT 0.2 | 每级提升比例(20%) |

```sql
CREATE TABLE passive_skills (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    icon VARCHAR(64),
    class_id VARCHAR(32),
    rarity VARCHAR(16) DEFAULT 'common',
    tier INTEGER DEFAULT 1,
    effect_type VARCHAR(32) NOT NULL,
    effect_value REAL NOT NULL,
    effect_stat VARCHAR(32),
    max_level INTEGER DEFAULT 5,
    level_scaling REAL DEFAULT 0.2,
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

CREATE INDEX idx_passive_skills_class ON passive_skills(class_id);
CREATE INDEX idx_passive_skills_tier ON passive_skills(tier);
```

#### character_passive_skills - 角色被动技能表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| passive_id | VARCHAR(32) | NOT NULL FK | 被动技能ID |
| level | INTEGER | DEFAULT 1 | 当前等级 |
| acquired_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 获得时间 |

```sql
CREATE TABLE character_passive_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    passive_id VARCHAR(32) NOT NULL,
    level INTEGER DEFAULT 1,
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (passive_id) REFERENCES passive_skills(id),
    UNIQUE(character_id, passive_id)
);

CREATE INDEX idx_char_passive_char ON character_passive_skills(character_id);
```

---

### 升级里程碑一览

| 等级 | 累计属性点 | 技能选择次数 | 里程碑事件 |
|-----|----------|------------|----------|
| 1 | 0 | 0 | 角色创建 |
| 3 | 2 | 1 | 首次技能选择 |
| 6 | 5 | 2 | |
| 10 | 9 | 3 | 解锁第2个队伍槽位 |
| 15 | 14 | 5 | 进入进阶技能池 |
| 20 | 19 | 6 | 解锁第3个队伍槽位 |
| 30 | 29 | 10 | 中期里程碑 |
| 35 | 34 | 11 | 解锁第4个队伍槽位 |
| 39 | 38 | 13 | 进入大师技能池 |
| 50 | 49 | 16 | 解锁第5个队伍槽位 |
| 60 | 59 | 20 | 满级 |

---

### 相关游戏命令

| 命令 | 说明 | 示例 |
|-----|------|-----|
| `/points` | 查看可分配点数 | 显示当前未分配点数 |
| `/allocate [属性] [点数]` | 分配属性点 | `/allocate strength 5` |
| `/respec` | 重置属性点 | 消耗金币重置所有分配 |
| `/passives` | 查看已有被动技能 | 列出所有被动及等级 |
| `/build` | 查看角色培养方案 | 显示属性分配和技能选择 |

---

## 📋 表结构设计

### 1. users - 用户表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 用户ID |
| username | VARCHAR(32) | UNIQUE NOT NULL | 用户名 |
| password_hash | VARCHAR(256) | NOT NULL | 密码哈希 |
| email | VARCHAR(128) | UNIQUE | 邮箱（可选） |
| max_team_size | INTEGER | DEFAULT 5 | 最大队伍人数上限 |
| unlocked_slots | INTEGER | DEFAULT 1 | 已解锁槽位数(初始1个) |
| gold | INTEGER | DEFAULT 0 | 金币(小队共享) |
| current_zone_id | VARCHAR(32) | DEFAULT 'elwynn' | 当前区域(小队共享) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 注册时间 |
| last_login_at | DATETIME | | 最后登录时间 |
| status | INTEGER | DEFAULT 1 | 状态: 1正常 0禁用 |

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(32) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    email VARCHAR(128) UNIQUE,
    max_team_size INTEGER DEFAULT 5,
    unlocked_slots INTEGER DEFAULT 1,  -- 已解锁槽位数(初始1个)
    gold INTEGER DEFAULT 0,
    current_zone_id VARCHAR(32) DEFAULT 'elwynn',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME,
    status INTEGER DEFAULT 1
);

CREATE INDEX idx_users_username ON users(username);
```

---

### 2. characters - 角色表

> 📌 **小队系统**: 每个用户可以拥有最多5个角色组成小队，共同参与战斗。

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 角色ID |
| user_id | INTEGER | NOT NULL FK | 所属用户 |
| name | VARCHAR(32) | NOT NULL | 角色名 |
| race_id | VARCHAR(32) | NOT NULL | 种族ID |
| class_id | VARCHAR(32) | NOT NULL | 职业ID |
| faction | VARCHAR(16) | NOT NULL | 阵营: alliance/horde |
| team_slot | INTEGER | NOT NULL | 队伍位置: 1-5 (1=队长) |
| is_active | INTEGER | DEFAULT 1 | 是否出战: 1是 0否 |
| is_dead | INTEGER | DEFAULT 0 | 是否死亡: 1是 0否 |
| revive_at | DATETIME | NULL | 复活时间(NULL表示存活) |
| level | INTEGER | DEFAULT 1 | 等级 |
| exp | INTEGER | DEFAULT 0 | 当前经验 |
| exp_to_next | INTEGER | DEFAULT 100 | 升级所需经验 |
| hp | INTEGER | NOT NULL | 当前生命值 |
| max_hp | INTEGER | NOT NULL | 最大生命值 |
| resource | INTEGER | NOT NULL | 当前能量值(怒气/能量/法力) |
| max_resource | INTEGER | NOT NULL | 最大能量值 |
| resource_type | VARCHAR(16) | NOT NULL | 能量类型: mana/rage/energy (继承自职业) |
| strength | INTEGER | DEFAULT 10 | 力量 |
| agility | INTEGER | DEFAULT 10 | 敏捷 |
| intellect | INTEGER | DEFAULT 10 | 智力 |
| stamina | INTEGER | DEFAULT 10 | 耐力 |
| spirit | INTEGER | DEFAULT 10 | 精神 |
| attack | INTEGER | DEFAULT 10 | 攻击力 |
| defense | INTEGER | DEFAULT 5 | 防御力 |
| crit_rate | REAL | DEFAULT 0.05 | 暴击率 |
| crit_damage | REAL | DEFAULT 1.5 | 暴击伤害倍率 |
| current_zone_id | VARCHAR(32) | DEFAULT 'elwynn' | 当前区域(跟随队长) |
| total_kills | INTEGER | DEFAULT 0 | 总击杀数 |
| total_deaths | INTEGER | DEFAULT 0 | 总死亡数 |
| total_exp_gained | INTEGER | DEFAULT 0 | 总获得经验 |
| total_gold_gained | INTEGER | DEFAULT 0 | 总获得金币 |
| play_time | INTEGER | DEFAULT 0 | 游戏时长(秒) |
| offline_time | DATETIME | | 离线时间 |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

```sql
CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    race_id VARCHAR(32) NOT NULL,
    class_id VARCHAR(32) NOT NULL,
    faction VARCHAR(16) NOT NULL,
    team_slot INTEGER NOT NULL,
    is_active INTEGER DEFAULT 1,
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    exp_to_next INTEGER DEFAULT 100,
    hp INTEGER NOT NULL,
    max_hp INTEGER NOT NULL,
    mp INTEGER NOT NULL,
    max_mp INTEGER NOT NULL,
    strength INTEGER DEFAULT 10,
    agility INTEGER DEFAULT 10,
    intellect INTEGER DEFAULT 10,
    stamina INTEGER DEFAULT 10,
    spirit INTEGER DEFAULT 10,
    attack INTEGER DEFAULT 10,
    defense INTEGER DEFAULT 5,
    crit_rate REAL DEFAULT 0.05,
    crit_damage REAL DEFAULT 1.5,
    gold INTEGER DEFAULT 0,
    current_zone_id VARCHAR(32) DEFAULT 'elwynn',
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    total_exp_gained INTEGER DEFAULT 0,
    total_gold_gained INTEGER DEFAULT 0,
    play_time INTEGER DEFAULT 0,
    offline_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, team_slot)
);

CREATE INDEX idx_characters_user_id ON characters(user_id);
CREATE INDEX idx_characters_level ON characters(level);
CREATE INDEX idx_characters_team ON characters(user_id, team_slot);
```

---

### 3. races - 种族配置表

> 📌 **种族差异化**: 每个种族有独特的属性加成、主动技能和被动特性

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 种族ID |
| name | VARCHAR(32) | NOT NULL | 种族名称 |
| faction | VARCHAR(16) | NOT NULL | 阵营 |
| description | TEXT | | 描述 |
| strength_base | INTEGER | DEFAULT 0 | 力量基础加成(固定值) |
| strength_pct | REAL | DEFAULT 0 | 力量百分比加成 |
| agility_base | INTEGER | DEFAULT 0 | 敏捷基础加成 |
| agility_pct | REAL | DEFAULT 0 | 敏捷百分比加成 |
| intellect_base | INTEGER | DEFAULT 0 | 智力基础加成 |
| intellect_pct | REAL | DEFAULT 0 | 智力百分比加成 |
| stamina_base | INTEGER | DEFAULT 0 | 耐力基础加成 |
| stamina_pct | REAL | DEFAULT 0 | 耐力百分比加成 |
| spirit_base | INTEGER | DEFAULT 0 | 精神基础加成 |
| spirit_pct | REAL | DEFAULT 0 | 精神百分比加成 |
| racial_passive_id | VARCHAR(32) | | 种族被动特性1 |
| racial_passive2_id | VARCHAR(32) | | 种族被动特性2 |
| allowed_classes | TEXT | | 可选职业(JSON数组,null为全部) |

```sql
CREATE TABLE races (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    faction VARCHAR(16) NOT NULL,
    description TEXT,
    -- 基础加成(固定值，创建时一次性加)
    strength_base INTEGER DEFAULT 0,
    agility_base INTEGER DEFAULT 0,
    intellect_base INTEGER DEFAULT 0,
    stamina_base INTEGER DEFAULT 0,
    spirit_base INTEGER DEFAULT 0,
    -- 百分比加成(乘算，随等级保持意义)
    strength_pct REAL DEFAULT 0,
    agility_pct REAL DEFAULT 0,
    intellect_pct REAL DEFAULT 0,
    stamina_pct REAL DEFAULT 0,
    spirit_pct REAL DEFAULT 0,
    -- 被动特性
    racial_passive_id VARCHAR(32),
    racial_passive2_id VARCHAR(32),
    allowed_classes TEXT,
    FOREIGN KEY (racial_passive_id) REFERENCES effects(id),
    FOREIGN KEY (racial_passive2_id) REFERENCES effects(id)
);
```

### 种族差异化设计

> 📌 **放置游戏适配**: 所有种族特性都设计为自动触发或被动效果，适合挂机场景

> 📌 **说明**: 种族属性加成的详细设计已移至 [角色属性系统设计文档](./character_attributes_design.md)

种族属性加成采用"基础值 + 百分比"混合方案，详细计算公式和种族列表请参考角色属性设计文档。

---

### 4. classes - 职业配置表

> 📌 **能量系统差异**: 不同职业使用不同的能量类型，体现职业特色

#### 能量类型设计

> 📌 **职业差异化核心**: 不同职业的能量系统完全不同，这是职业特色的重要体现

| 能量类型 | 使用职业 | 初始值 | 最大值 | 恢复机制 | 特色 |
|---------|---------|-------|-------|---------|------|
| **怒气 (rage)** | 战士 | 0 | 100 | 攻击/受击获得，脱战衰减 | 越打越强 |
| **能量 (energy)** | 盗贼 | 100 | 100 | 每回合+20，上限固定 | 快速循环 |
| **法力 (mana)** | 法师、牧师、术士、德鲁伊、萨满 | 满 | 基于智力成长 | 每回合恢复(精神×系数) | 持续施法 |
| **法力 (mana)** | 圣骑士、猎人 | 满 | 较低 | 每回合恢复(较慢) | 混合职业 |

**怒气机制详解 (战士):**
```
获得怒气:
  - 普通攻击命中: +5 怒气
  - 造成暴击: +10 怒气  
  - 受到伤害: +怒气 (受伤比例 × 2)
  - 使用某些技能: +怒气
  
消耗怒气:
  - 使用技能消耗怒气 (10-30点)
  - 脱战后每回合 -10 怒气 (衰减到0)
  
策略要点:
  - 开局无怒气，需要先攻击积累
  - 受伤也能获得怒气，坦克优势
  - 需要持续战斗保持怒气值
```

**能量机制详解 (盗贼):**
```
恢复能量:
  - 每回合自动 +20 能量
  - 上限固定100，不会超过
  - 某些技能/天赋可加速恢复
  
消耗能量:
  - 使用技能消耗能量 (20-40点)
  - 不会自然衰减
  - 能量不足时无法使用技能
  
策略要点:
  - 能量恢复快，适合频繁使用技能
  - 需要合理分配能量，避免浪费
  - 连击技能组合需要能量管理
```

**法力机制详解 (法系职业):**
```
恢复法力:
  - 每回合恢复 = 精神 × 恢复系数
  - 法师/术士: 精神 × 0.5%
  - 牧师: 精神 × 0.8% (治疗职业)
  - 德鲁伊/萨满: 精神 × 0.6%
  - 圣骑士/猎人: 精神 × 0.3-0.4% (混合职业)
  
消耗法力:
  - 使用技能消耗法力 (15-50点)
  - 法力不足时无法施法
  
策略要点:
  - 精神属性影响恢复速度
  - 需要平衡输出和治疗
  - 法力耗尽会失去战斗力
```

**能量系统对比表:**

| 特性 | 怒气 | 能量 | 法力 |
|-----|------|------|------|
| 开局状态 | 0 (需积累) | 100 (满) | 满 |
| 恢复方式 | 战斗获得 | 每回合+20 | 每回合恢复(精神) |
| 上限 | 100 | 100 | 成长型 |
| 衰减 | 脱战衰减 | 无 | 无 |
| 策略重点 | 保持战斗节奏 | 快速循环 | 资源管理 |

---

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 职业ID |
| name | VARCHAR(32) | NOT NULL | 职业名称 |
| description | TEXT | | 描述 |
| role | VARCHAR(16) | NOT NULL | 定位: tank/dps/healer |
| primary_stat | VARCHAR(16) | NOT NULL | 主属性 |
| resource_type | VARCHAR(16) | NOT NULL | 能量类型: mana/rage/energy |
| base_hp | INTEGER | NOT NULL | 基础HP |
| base_resource | INTEGER | NOT NULL | 基础能量值 |
| hp_per_level | INTEGER | NOT NULL | 每级HP成长 |
| resource_per_level | INTEGER | NOT NULL | 每级能量成长 |
| resource_regen | REAL | DEFAULT 0 | 每回合能量恢复(固定值) |
| resource_regen_pct | REAL | DEFAULT 0 | 每回合能量恢复(百分比) |
| base_strength | INTEGER | DEFAULT 10 | 基础力量 |
| base_agility | INTEGER | DEFAULT 10 | 基础敏捷 |
| base_intellect | INTEGER | DEFAULT 10 | 基础智力 |
| base_stamina | INTEGER | DEFAULT 10 | 基础耐力 |
| base_spirit | INTEGER | DEFAULT 10 | 基础精神 |

```sql
CREATE TABLE classes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    role VARCHAR(16) NOT NULL,
    primary_stat VARCHAR(16) NOT NULL,
    resource_type VARCHAR(16) NOT NULL,  -- mana/rage/energy
    base_hp INTEGER NOT NULL,
    base_resource INTEGER NOT NULL,       -- 基础能量值
    hp_per_level INTEGER NOT NULL,
    resource_per_level INTEGER NOT NULL,  -- 每级能量成长
    resource_regen REAL DEFAULT 0,        -- 每回合固定恢复
    resource_regen_pct REAL DEFAULT 0,    -- 每回合百分比恢复
    base_strength INTEGER DEFAULT 10,
    base_agility INTEGER DEFAULT 10,
    base_intellect INTEGER DEFAULT 10,
    base_stamina INTEGER DEFAULT 10,
    base_spirit INTEGER DEFAULT 10
);
```

---

### 5. skills - 技能配置表

> 📌 **Buff/Debuff机制**: 所有增益/减益效果仅在战斗中生效，每场战斗开始前自动清空。

#### 回合制战斗技能设计原则

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       技能设计规范 (回合制自动战斗)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  【设计原则】                                                                 │
│  ├─ 所有技能适配自动战斗，无需玩家手动操作                                    │
│  ├─ 冷却时间以"回合"计算，而非实时秒数                                        │
│  ├─ 资源消耗平衡，保证技能循环可持续                                          │
│  └─ 强力技能必须有冷却，避免无脑堆叠                                          │
│                                                                             │
│  【技能类型平衡】                                                             │
│  ├─ 基础技能: 无/低消耗，无CD，提供稳定输出                                   │
│  ├─ 核心技能: 中等消耗，短CD(2-3回合)，主要伤害来源                           │
│  ├─ 强力技能: 高消耗，中CD(4-6回合)，爆发伤害                                 │
│  └─ 终极技能: 高消耗，长CD(8-15回合)，战斗转折点                              │
│                                                                             │
│  【目标类型适配】                                                             │
│  ├─ enemy: AI自动选择优先目标                                                 │
│  ├─ enemy_lowest_hp: 自动斩杀低血量敌人                                       │
│  ├─ ally_lowest_hp: 自动优先治疗低血队友                                      │
│  └─ enemy_all/ally_all: AOE效果，自动全体作用                                 │
│                                                                             │
│  【不适合自动战斗的概念】(已移除/调整)                                         │
│  ├─ ❌ 位置相关: "背刺从背后攻击" → ✅ "伏击对满血敌人加倍"                   │
│  ├─ ❌ 隐身概念: "消失潜行" → ✅ "闪避提升闪避率"                              │
│  └─ ❌ 无限堆叠: "斩杀无CD" → ✅ "斩杀3回合CD"                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 各职业技能数量统计

| 职业 | 技能数量 | 能量类型 | 特点 |
|-----|---------|---------|------|
| 战士 | 8 | 怒气 (100上限) | 攻击/受伤获取，强力爆发 |
| 法师 | 8 | 法力 (40基础) | 高伤害，法系主力 |
| 盗贼 | 8 | 能量 (100上限) | 快速恢复，连击体系 |
| 牧师 | 9 | 法力 (35基础) | 治疗为主，辅助输出 |
| 圣骑士 | 8 | 法力 (20基础) | 混合职业，攻防兼备 |
| 猎人 | 8 | 法力 (18基础) | 远程物理+宠物 |
| 术士 | 8 | 法力 (38基础) | DOT+吸血+召唤 |
| 德鲁伊 | 8 | 法力 (30基础) | 变形+治疗+DOT |
| 萨满 | 8 | 法力 (32基础) | 元素+图腾+治疗 |
| 通用 | 1 | 无 | 普通攻击 |

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 技能ID |
| name | VARCHAR(32) | NOT NULL | 技能名称 |
| description | TEXT | | 描述 |
| icon | VARCHAR(64) | | 图标标识(预留) |
| class_id | VARCHAR(32) | | 所属职业(null为通用) |
| type | VARCHAR(16) | NOT NULL | 类型(见下表) |
| target_type | VARCHAR(16) | NOT NULL | 目标类型(见下表) |
| damage_type | VARCHAR(16) | | 伤害类型: physical/magic/true/nature/fire/frost/shadow/holy |
| base_value | INTEGER | DEFAULT 0 | 基础数值(伤害/治疗/效果强度) |
| scaling_stat | VARCHAR(16) | | 成长属性: strength/agility/intellect/spirit |
| scaling_ratio | REAL | DEFAULT 1.0 | 属性加成系数 |
| resource_cost | INTEGER | DEFAULT 0 | 能量消耗(怒气/能量/法力) |
| cooldown | INTEGER | DEFAULT 0 | 冷却时间(回合) |
| level_required | INTEGER | DEFAULT 1 | 需求等级 |
| effect_id | VARCHAR(32) | | 附加效果ID(关联effects表) |
| effect_chance | REAL | DEFAULT 1.0 | 效果触发概率(0-1) |
| tags | TEXT | | 标签(JSON数组，用于分类筛选) |

**技能类型 (type):**
| 类型 | 说明 |
|-----|------|
| `attack` | 造成伤害 |
| `heal` | 恢复生命 |
| `buff` | 增益效果 |
| `debuff` | 减益效果 |
| `dot` | 持续伤害(Damage over Time) |
| `hot` | 持续治疗(Heal over Time) |
| `shield` | 伤害吸收护盾 |
| `summon` | 召唤(预留) |
| `dispel` | 驱散效果 |
| `interrupt` | 打断施法 |
| `control` | 控制(眩晕/沉默等) |

**目标类型 (target_type):**
| 类型 | 说明 |
|-----|------|
| `self` | 自身 |
| `ally` | 友方单体 |
| `ally_all` | 友方全体 |
| `ally_lowest_hp` | 血量最低的友方 |
| `enemy` | 敌方单体 |
| `enemy_all` | 敌方全体 |
| `enemy_random` | 随机敌人 |
| `enemy_lowest_hp` | 血量最低的敌人 |

---

### 5.1 effects - 效果配置表(Buff/Debuff)

> 📌 **扩展性设计**: 通过独立的效果表，支持技能附带各种复杂效果

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 效果ID |
| name | VARCHAR(32) | NOT NULL | 效果名称 |
| description | TEXT | | 描述 |
| icon | VARCHAR(64) | | 图标标识(预留) |
| type | VARCHAR(16) | NOT NULL | 效果类型(见下表) |
| is_buff | INTEGER | NOT NULL | 1=增益 0=减益 |
| is_stackable | INTEGER | DEFAULT 0 | 是否可叠加 |
| max_stacks | INTEGER | DEFAULT 1 | 最大叠加层数 |
| duration | INTEGER | NOT NULL | 持续回合数 |
| tick_interval | INTEGER | DEFAULT 1 | 触发间隔(回合) |
| value_type | VARCHAR(16) | | 数值类型: flat/percent |
| value | REAL | | 效果数值 |
| stat_affected | VARCHAR(32) | | 影响的属性 |
| damage_type | VARCHAR(16) | | DOT伤害类型 |
| can_dispel | INTEGER | DEFAULT 1 | 是否可驱散 |
| tags | TEXT | | 标签(JSON数组) |

**效果类型 (type):**
| 类型 | 说明 | 示例 |
|-----|------|-----|
| `stat_mod` | 属性修改 | 攻击力+10% |
| `dot` | 持续伤害 | 中毒、燃烧 |
| `hot` | 持续治疗 | 回春术 |
| `shield` | 伤害吸收 | 真言术:盾 |
| `stun` | 眩晕 | 无法行动 |
| `silence` | 沉默 | 无法施法 |
| `slow` | 减速 | 攻击速度降低(预留) |
| `root` | 定身 | 无法移动(预留) |
| `taunt` | 嘲讽 | 强制攻击自己 |
| `immunity` | 免疫 | 免疫某类伤害 |
| `reflect` | 反射 | 反弹伤害 |
| `lifesteal` | 吸血 | 造成伤害时回血 |
| `thorns` | 荆棘 | 被攻击时反伤 |
| `stealth` | 潜行 | 隐身状态 |
| `invulnerable` | 无敌 | 免疫所有伤害 |

**可影响的属性 (stat_affected):**
- `attack` - 攻击力
- `defense` - 防御力
- `max_hp` - 最大生命值
- `max_mp` - 最大法力值
- `crit_rate` - 暴击率
- `crit_damage` - 暴击伤害
- `hit_rate` - 命中率(预留)
- `dodge_rate` - 闪避率(预留)
- `damage_taken` - 受到的伤害
- `damage_dealt` - 造成的伤害
- `healing_taken` - 受到的治疗
- `healing_done` - 造成的治疗

```sql
-- 技能表定义见上方"5. skills - 技能配置表"章节的完整定义
-- 主要字段包括: resource_cost (能量消耗，支持怒气/能量/法力)
```

---

### 6. character_skills - 角色技能表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| skill_id | VARCHAR(32) | NOT NULL FK | 技能ID |
| skill_level | INTEGER | DEFAULT 1 | 技能等级 |
| slot | INTEGER | | 技能槽位(null为未装备) |
| is_auto | INTEGER | DEFAULT 1 | 自动释放: 1是 0否 |

```sql
CREATE TABLE character_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    skill_level INTEGER DEFAULT 1,
    slot INTEGER,
    is_auto INTEGER DEFAULT 1,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id),
    UNIQUE(character_id, skill_id)
);

CREATE INDEX idx_char_skills_char_id ON character_skills(character_id);
```

---

### 7. items - 物品配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 物品ID |
| name | VARCHAR(64) | NOT NULL | 物品名称 |
| description | TEXT | | 描述 |
| type | VARCHAR(16) | NOT NULL | 类型: equipment/consumable/material/quest |
| subtype | VARCHAR(16) | | 子类型 |
| quality | VARCHAR(16) | DEFAULT 'common' | 品质: common/uncommon/rare/epic/legendary |
| level_required | INTEGER | DEFAULT 1 | 需求等级 |
| class_required | VARCHAR(32) | | 需求职业 |
| slot | VARCHAR(16) | | 装备槽位 |
| stackable | INTEGER | DEFAULT 0 | 可堆叠: 1是 0否 |
| max_stack | INTEGER | DEFAULT 1 | 最大堆叠数 |
| sell_price | INTEGER | DEFAULT 0 | 售价 |
| buy_price | INTEGER | DEFAULT 0 | 购买价 |
| strength | INTEGER | DEFAULT 0 | 力量加成 |
| agility | INTEGER | DEFAULT 0 | 敏捷加成 |
| intellect | INTEGER | DEFAULT 0 | 智力加成 |
| stamina | INTEGER | DEFAULT 0 | 耐力加成 |
| spirit | INTEGER | DEFAULT 0 | 精神加成 |
| attack | INTEGER | DEFAULT 0 | 攻击加成 |
| defense | INTEGER | DEFAULT 0 | 防御加成 |
| hp_bonus | INTEGER | DEFAULT 0 | HP加成 |
| mp_bonus | INTEGER | DEFAULT 0 | MP加成 |
| crit_rate | REAL | DEFAULT 0 | 暴击率加成 |
| effect_type | VARCHAR(32) | | 使用效果类型 |
| effect_value | INTEGER | | 使用效果数值 |

```sql
CREATE TABLE items (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    type VARCHAR(16) NOT NULL,
    subtype VARCHAR(16),
    quality VARCHAR(16) DEFAULT 'common',
    level_required INTEGER DEFAULT 1,
    class_required VARCHAR(32),
    slot VARCHAR(16),
    stackable INTEGER DEFAULT 0,
    max_stack INTEGER DEFAULT 1,
    sell_price INTEGER DEFAULT 0,
    buy_price INTEGER DEFAULT 0,
    strength INTEGER DEFAULT 0,
    agility INTEGER DEFAULT 0,
    intellect INTEGER DEFAULT 0,
    stamina INTEGER DEFAULT 0,
    spirit INTEGER DEFAULT 0,
    attack INTEGER DEFAULT 0,
    defense INTEGER DEFAULT 0,
    hp_bonus INTEGER DEFAULT 0,
    mp_bonus INTEGER DEFAULT 0,
    crit_rate REAL DEFAULT 0,
    effect_type VARCHAR(32),
    effect_value INTEGER
);

CREATE INDEX idx_items_type ON items(type);
CREATE INDEX idx_items_quality ON items(quality);
```

**装备槽位 (slot) 定义:**
- `head` - 头部
- `shoulder` - 肩部
- `chest` - 胸甲
- `hands` - 手套
- `legs` - 腿部
- `feet` - 脚部
- `main_hand` - 主手武器
- `off_hand` - 副手
- `neck` - 项链
- `ring` - 戒指
- `trinket` - 饰品

---

### 8. inventory - 背包表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| item_id | VARCHAR(32) | NOT NULL FK | 物品ID |
| quantity | INTEGER | DEFAULT 1 | 数量 |
| slot | INTEGER | | 背包槽位 |

```sql
CREATE TABLE inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    quantity INTEGER DEFAULT 1,
    slot INTEGER,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX idx_inventory_char_id ON inventory(character_id);
```

---

### 9. equipment - 装备表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| slot | VARCHAR(16) | NOT NULL | 装备槽位 |
| item_id | VARCHAR(32) | NOT NULL FK | 物品ID |

```sql
CREATE TABLE equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    slot VARCHAR(16) NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id),
    UNIQUE(character_id, slot)
);

CREATE INDEX idx_equipment_char_id ON equipment(character_id);
```

---

### 10. zones - 区域配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 区域ID |
| name | VARCHAR(64) | NOT NULL | 区域名称 |
| description | TEXT | | 描述 |
| min_level | INTEGER | DEFAULT 1 | 最低等级 |
| max_level | INTEGER | DEFAULT 60 | 最高等级 |
| faction | VARCHAR(16) | | 阵营限制 |
| parent_zone_id | VARCHAR(32) | | 父区域 |
| exp_modifier | REAL | DEFAULT 1.0 | 经验倍率 |
| gold_modifier | REAL | DEFAULT 1.0 | 金币倍率 |
| drop_modifier | REAL | DEFAULT 1.0 | 掉落倍率 |
| is_dungeon | INTEGER | DEFAULT 0 | 是否副本 |
| unlock_condition | TEXT | | 解锁条件(JSON) |

```sql
CREATE TABLE zones (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    min_level INTEGER DEFAULT 1,
    max_level INTEGER DEFAULT 60,
    faction VARCHAR(16),
    parent_zone_id VARCHAR(32),
    exp_modifier REAL DEFAULT 1.0,
    gold_modifier REAL DEFAULT 1.0,
    drop_modifier REAL DEFAULT 1.0,
    is_dungeon INTEGER DEFAULT 0,
    unlock_condition TEXT,
    FOREIGN KEY (parent_zone_id) REFERENCES zones(id)
);
```

---

### 11. monsters - 怪物配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 怪物ID |
| zone_id | VARCHAR(32) | NOT NULL FK | 所属区域 |
| name | VARCHAR(64) | NOT NULL | 怪物名称 |
| description | TEXT | | 描述 |
| level | INTEGER | NOT NULL | 等级 |
| type | VARCHAR(16) | DEFAULT 'normal' | 类型: normal/elite/boss |
| hp | INTEGER | NOT NULL | 生命值 |
| mp | INTEGER | DEFAULT 0 | 法力值 |
| attack | INTEGER | NOT NULL | 攻击力 |
| defense | INTEGER | NOT NULL | 防御力 |
| exp_reward | INTEGER | NOT NULL | 经验奖励 |
| gold_min | INTEGER | DEFAULT 0 | 最小金币 |
| gold_max | INTEGER | DEFAULT 0 | 最大金币 |
| spawn_weight | INTEGER | DEFAULT 100 | 生成权重 |
| skills | TEXT | | 技能列表(JSON) |

```sql
CREATE TABLE monsters (
    id VARCHAR(32) PRIMARY KEY,
    zone_id VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL,
    type VARCHAR(16) DEFAULT 'normal',
    hp INTEGER NOT NULL,
    mp INTEGER DEFAULT 0,
    attack INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    exp_reward INTEGER NOT NULL,
    gold_min INTEGER DEFAULT 0,
    gold_max INTEGER DEFAULT 0,
    spawn_weight INTEGER DEFAULT 100,
    skills TEXT,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

CREATE INDEX idx_monsters_zone_id ON monsters(zone_id);
CREATE INDEX idx_monsters_level ON monsters(level);
```

---

### 12. monster_drops - 怪物掉落表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| monster_id | VARCHAR(32) | NOT NULL FK | 怪物ID |
| item_id | VARCHAR(32) | NOT NULL FK | 物品ID |
| drop_rate | REAL | NOT NULL | 掉落率(0-1) |
| min_quantity | INTEGER | DEFAULT 1 | 最小数量 |
| max_quantity | INTEGER | DEFAULT 1 | 最大数量 |

```sql
CREATE TABLE monster_drops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    monster_id VARCHAR(32) NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    drop_rate REAL NOT NULL,
    min_quantity INTEGER DEFAULT 1,
    max_quantity INTEGER DEFAULT 1,
    FOREIGN KEY (monster_id) REFERENCES monsters(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX idx_monster_drops_monster_id ON monster_drops(monster_id);
```

---

### 13. battle_strategies - 战斗策略表

> 📎 **详细设计请参阅**: [作战策略系统设计](strategy_design.md)

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 策略ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| name | VARCHAR(32) | NOT NULL | 策略名称 |
| is_active | INTEGER | DEFAULT 0 | 是否当前使用 |
| skill_priority | TEXT | | 技能优先级 (JSON数组) |
| conditional_rules | TEXT | | 条件规则 (JSON数组) |
| target_priority | VARCHAR(32) | DEFAULT 'lowest_hp' | 默认目标选择策略 |
| skill_target_overrides | TEXT | | 技能目标覆盖 (JSON对象) |
| resource_threshold | INTEGER | DEFAULT 0 | 资源阈值 |
| reserved_skills | TEXT | | 保留技能 (JSON数组) |
| auto_target_settings | TEXT | | 智能目标设置 (JSON对象) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | | 更新时间 |

```sql
CREATE TABLE IF NOT EXISTS battle_strategies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    is_active INTEGER DEFAULT 0,
    skill_priority TEXT,
    conditional_rules TEXT,
    target_priority VARCHAR(32) DEFAULT 'lowest_hp',
    skill_target_overrides TEXT,
    resource_threshold INTEGER DEFAULT 0,
    reserved_skills TEXT,
    auto_target_settings TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX idx_battle_strategies_character ON battle_strategies(character_id);
CREATE INDEX idx_battle_strategies_active ON battle_strategies(character_id, is_active);
```

---

### 14. game_sessions - 游戏会话表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| session_start | DATETIME | NOT NULL | 会话开始 |
| session_end | DATETIME | | 会话结束 |
| kills | INTEGER | DEFAULT 0 | 击杀数 |
| exp_gained | INTEGER | DEFAULT 0 | 获得经验 |
| gold_gained | INTEGER | DEFAULT 0 | 获得金币 |
| deaths | INTEGER | DEFAULT 0 | 死亡次数 |

```sql
CREATE TABLE game_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    session_start DATETIME NOT NULL,
    session_end DATETIME,
    kills INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_char_id ON game_sessions(character_id);
```

---

## 📊 战斗数据分析系统

> 📌 **数据驱动决策**: 完整记录每场战斗和每个角色的详细数据，帮助玩家分析团队表现、优化战术配置

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          战斗数据分析系统                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        实时战斗面板                                   │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │   │
│  │  │ 战士     │ │ 法师     │ │ 牧师     │ │ 盗贼     │ │ 猎人     │  │   │
│  │  │ DPS:156  │ │ DPS:203  │ │ HPS:89   │ │ DPS:178  │ │ DPS:145  │  │   │
│  │  │ 承伤:45% │ │ 承伤:8%  │ │ 承伤:12% │ │ 承伤:15% │ │ 承伤:20% │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌──────────────────┐    │
│  │ 📈 伤害统计          │  │ 🛡️ 承伤统计          │  │ 💚 治疗统计     │    │
│  │ 总伤害: 125,432     │  │ 总承伤: 45,678      │  │ 总治疗: 32,100  │    │
│  │ 暴击率: 23.5%       │  │ 闪避次数: 156       │  │ 过量治疗: 12%   │    │
│  │ 最高单次: 1,234     │  │ 格挡次数: 89        │  │ 治疗目标分布    │    │
│  └─────────────────────┘  └─────────────────────┘  └──────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 核心统计指标

| 类别 | 指标 | 说明 | 用途 |
|-----|------|-----|------|
| **输出** | 总伤害 (Damage Dealt) | 造成的总伤害值 | 评估输出能力 |
| | DPS (每回合伤害) | 伤害 ÷ 战斗回合数 | 对比输出效率 |
| | 暴击次数/率 | 暴击触发统计 | 验证暴击收益 |
| | 技能伤害占比 | 各技能贡献比例 | 优化技能选择 |
| **承伤** | 总承伤 (Damage Taken) | 受到的总伤害值 | 评估坦克压力 |
| | DTPS (每回合承伤) | 承伤 ÷ 战斗回合数 | 评估生存压力 |
| | 闪避/格挡/减伤 | 防御触发统计 | 验证防御收益 |
| **治疗** | 总治疗 (Healing Done) | 产生的总治疗量 | 评估治疗能力 |
| | HPS (每回合治疗) | 治疗 ÷ 战斗回合数 | 对比治疗效率 |
| | 过量治疗 (Overheal) | 溢出的治疗量 | 优化治疗时机 |
| | 受到治疗 | 收到的治疗量 | 评估被关注度 |
| **效率** | 击杀数 | 击杀的敌人数量 | 综合贡献 |
| | 死亡数 | 被击倒次数 | 生存能力 |
| | 技能命中率 | 技能实际命中比例 | 技能效率 |

---

### 15. battle_records - 战斗记录表

> 📌 **单场战斗记录**: 记录每一场战斗的基本信息

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 战斗ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| zone_id | VARCHAR(32) | NOT NULL FK | 区域ID |
| battle_type | VARCHAR(16) | NOT NULL | 类型: pve/pvp/boss/abyss |
| monster_id | VARCHAR(32) | | 怪物ID(PVE) |
| opponent_user_id | INTEGER | | 对手ID(PVP) |
| total_rounds | INTEGER | DEFAULT 0 | 总回合数 |
| duration_seconds | INTEGER | DEFAULT 0 | 战斗时长(秒) |
| result | VARCHAR(16) | NOT NULL | 结果: victory/defeat/draw/flee |
| team_damage_dealt | INTEGER | DEFAULT 0 | 队伍总输出伤害 |
| team_damage_taken | INTEGER | DEFAULT 0 | 队伍总承受伤害 |
| team_healing_done | INTEGER | DEFAULT 0 | 队伍总治疗量 |
| exp_gained | INTEGER | DEFAULT 0 | 获得经验 |
| gold_gained | INTEGER | DEFAULT 0 | 获得金币 |
| battle_log | TEXT | | 详细战斗日志(JSON) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 战斗时间 |

```sql
CREATE TABLE battle_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    zone_id VARCHAR(32) NOT NULL,
    battle_type VARCHAR(16) NOT NULL,
    monster_id VARCHAR(32),
    opponent_user_id INTEGER,
    total_rounds INTEGER DEFAULT 0,
    duration_seconds INTEGER DEFAULT 0,
    result VARCHAR(16) NOT NULL,
    team_damage_dealt INTEGER DEFAULT 0,
    team_damage_taken INTEGER DEFAULT 0,
    team_healing_done INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    battle_log TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (monster_id) REFERENCES monsters(id),
    FOREIGN KEY (opponent_user_id) REFERENCES users(id)
);

CREATE INDEX idx_battle_records_user ON battle_records(user_id);
CREATE INDEX idx_battle_records_time ON battle_records(created_at DESC);
CREATE INDEX idx_battle_records_zone ON battle_records(zone_id);
CREATE INDEX idx_battle_records_type ON battle_records(battle_type);
```

---

### 16. battle_character_stats - 战斗角色统计表

> 📌 **单场角色数据**: 记录每场战斗中每个角色的详细表现

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| battle_id | INTEGER | NOT NULL FK | 战斗记录ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| team_slot | INTEGER | NOT NULL | 队伍位置 |
| **─── 伤害统计 ───** |
| damage_dealt | INTEGER | DEFAULT 0 | 造成总伤害 |
| physical_damage | INTEGER | DEFAULT 0 | 物理伤害 |
| magic_damage | INTEGER | DEFAULT 0 | 魔法伤害 |
| fire_damage | INTEGER | DEFAULT 0 | 火焰伤害 |
| frost_damage | INTEGER | DEFAULT 0 | 冰霜伤害 |
| shadow_damage | INTEGER | DEFAULT 0 | 暗影伤害 |
| holy_damage | INTEGER | DEFAULT 0 | 神圣伤害 |
| nature_damage | INTEGER | DEFAULT 0 | 自然伤害 |
| dot_damage | INTEGER | DEFAULT 0 | DOT伤害 |
| **─── 暴击统计 ───** |
| crit_count | INTEGER | DEFAULT 0 | 暴击次数 |
| crit_damage | INTEGER | DEFAULT 0 | 暴击造成的总伤害 |
| max_crit | INTEGER | DEFAULT 0 | 最高单次暴击 |
| **─── 承伤统计 ───** |
| damage_taken | INTEGER | DEFAULT 0 | 受到总伤害 |
| physical_taken | INTEGER | DEFAULT 0 | 物理承伤 |
| magic_taken | INTEGER | DEFAULT 0 | 魔法承伤 |
| damage_blocked | INTEGER | DEFAULT 0 | 格挡伤害 |
| damage_absorbed | INTEGER | DEFAULT 0 | 护盾吸收 |
| **─── 闪避统计 ───** |
| dodge_count | INTEGER | DEFAULT 0 | 闪避次数 |
| block_count | INTEGER | DEFAULT 0 | 格挡次数 |
| hit_count | INTEGER | DEFAULT 0 | 被命中次数 |
| **─── 治疗统计 ───** |
| healing_done | INTEGER | DEFAULT 0 | 造成治疗量 |
| healing_received | INTEGER | DEFAULT 0 | 受到治疗量 |
| overhealing | INTEGER | DEFAULT 0 | 过量治疗 |
| self_healing | INTEGER | DEFAULT 0 | 自我治疗 |
| hot_healing | INTEGER | DEFAULT 0 | HOT治疗 |
| **─── 技能统计 ───** |
| skill_uses | INTEGER | DEFAULT 0 | 技能使用次数 |
| skill_hits | INTEGER | DEFAULT 0 | 技能命中次数 |
| skill_misses | INTEGER | DEFAULT 0 | 技能未命中 |
| **─── 控制统计 ───** |
| cc_applied | INTEGER | DEFAULT 0 | 施加控制次数 |
| cc_received | INTEGER | DEFAULT 0 | 受到控制次数 |
| dispels | INTEGER | DEFAULT 0 | 驱散次数 |
| interrupts | INTEGER | DEFAULT 0 | 打断次数 |
| **─── 其他统计 ───** |
| kills | INTEGER | DEFAULT 0 | 击杀数(最后一击) |
| deaths | INTEGER | DEFAULT 0 | 死亡次数 |
| resurrects | INTEGER | DEFAULT 0 | 复活次数 |
| resource_used | INTEGER | DEFAULT 0 | 消耗能量总量 |
| resource_generated | INTEGER | DEFAULT 0 | 获得能量总量 |

```sql
CREATE TABLE battle_character_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    team_slot INTEGER NOT NULL,
    -- 伤害统计
    damage_dealt INTEGER DEFAULT 0,
    physical_damage INTEGER DEFAULT 0,
    magic_damage INTEGER DEFAULT 0,
    fire_damage INTEGER DEFAULT 0,
    frost_damage INTEGER DEFAULT 0,
    shadow_damage INTEGER DEFAULT 0,
    holy_damage INTEGER DEFAULT 0,
    nature_damage INTEGER DEFAULT 0,
    dot_damage INTEGER DEFAULT 0,
    -- 暴击统计
    crit_count INTEGER DEFAULT 0,
    crit_damage INTEGER DEFAULT 0,
    max_crit INTEGER DEFAULT 0,
    -- 承伤统计
    damage_taken INTEGER DEFAULT 0,
    physical_taken INTEGER DEFAULT 0,
    magic_taken INTEGER DEFAULT 0,
    damage_blocked INTEGER DEFAULT 0,
    damage_absorbed INTEGER DEFAULT 0,
    -- 闪避统计
    dodge_count INTEGER DEFAULT 0,
    block_count INTEGER DEFAULT 0,
    hit_count INTEGER DEFAULT 0,
    -- 治疗统计
    healing_done INTEGER DEFAULT 0,
    healing_received INTEGER DEFAULT 0,
    overhealing INTEGER DEFAULT 0,
    self_healing INTEGER DEFAULT 0,
    hot_healing INTEGER DEFAULT 0,
    -- 技能统计
    skill_uses INTEGER DEFAULT 0,
    skill_hits INTEGER DEFAULT 0,
    skill_misses INTEGER DEFAULT 0,
    -- 控制统计
    cc_applied INTEGER DEFAULT 0,
    cc_received INTEGER DEFAULT 0,
    dispels INTEGER DEFAULT 0,
    interrupts INTEGER DEFAULT 0,
    -- 其他统计
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    resurrects INTEGER DEFAULT 0,
    resource_used INTEGER DEFAULT 0,
    resource_generated INTEGER DEFAULT 0,
    FOREIGN KEY (battle_id) REFERENCES battle_records(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX idx_battle_char_stats_battle ON battle_character_stats(battle_id);
CREATE INDEX idx_battle_char_stats_char ON battle_character_stats(character_id);
```

---

### 17. character_lifetime_stats - 角色生涯统计表

> 📌 **累计统计数据**: 角色从创建至今的所有战斗数据汇总

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| character_id | INTEGER | PRIMARY KEY FK | 角色ID |
| **─── 战斗场次 ───** |
| total_battles | INTEGER | DEFAULT 0 | 总战斗场数 |
| victories | INTEGER | DEFAULT 0 | 胜利场数 |
| defeats | INTEGER | DEFAULT 0 | 失败场数 |
| pve_battles | INTEGER | DEFAULT 0 | PVE战斗数 |
| pvp_battles | INTEGER | DEFAULT 0 | PVP战斗数 |
| boss_kills | INTEGER | DEFAULT 0 | Boss击杀数 |
| **─── 累计伤害 ───** |
| total_damage_dealt | INTEGER | DEFAULT 0 | 总造成伤害 |
| total_physical_damage | INTEGER | DEFAULT 0 | 物理总伤害 |
| total_magic_damage | INTEGER | DEFAULT 0 | 魔法总伤害 |
| total_crit_damage | INTEGER | DEFAULT 0 | 暴击总伤害 |
| total_crit_count | INTEGER | DEFAULT 0 | 总暴击次数 |
| highest_damage_single | INTEGER | DEFAULT 0 | 单次最高伤害 |
| highest_damage_battle | INTEGER | DEFAULT 0 | 单场最高伤害 |
| **─── 累计承伤 ───** |
| total_damage_taken | INTEGER | DEFAULT 0 | 总承受伤害 |
| total_damage_blocked | INTEGER | DEFAULT 0 | 总格挡伤害 |
| total_damage_absorbed | INTEGER | DEFAULT 0 | 总吸收伤害 |
| total_dodge_count | INTEGER | DEFAULT 0 | 总闪避次数 |
| **─── 累计治疗 ───** |
| total_healing_done | INTEGER | DEFAULT 0 | 总治疗量 |
| total_healing_received | INTEGER | DEFAULT 0 | 总受到治疗 |
| total_overhealing | INTEGER | DEFAULT 0 | 总过量治疗 |
| highest_healing_single | INTEGER | DEFAULT 0 | 单次最高治疗 |
| highest_healing_battle | INTEGER | DEFAULT 0 | 单场最高治疗 |
| **─── 击杀与死亡 ───** |
| total_kills | INTEGER | DEFAULT 0 | 总击杀数 |
| total_deaths | INTEGER | DEFAULT 0 | 总死亡数 |
| kill_streak_best | INTEGER | DEFAULT 0 | 最长连杀 |
| current_kill_streak | INTEGER | DEFAULT 0 | 当前连杀 |
| **─── 技能使用 ───** |
| total_skill_uses | INTEGER | DEFAULT 0 | 技能总使用次数 |
| total_skill_hits | INTEGER | DEFAULT 0 | 技能总命中数 |
| **─── 资源统计 ───** |
| total_resource_used | INTEGER | DEFAULT 0 | 总消耗能量 |
| total_rounds | INTEGER | DEFAULT 0 | 总战斗回合数 |
| total_battle_time | INTEGER | DEFAULT 0 | 总战斗时间(秒) |
| **─── 最后更新 ───** |
| last_battle_at | DATETIME | | 最后战斗时间 |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

```sql
CREATE TABLE character_lifetime_stats (
    character_id INTEGER PRIMARY KEY,
    -- 战斗场次
    total_battles INTEGER DEFAULT 0,
    victories INTEGER DEFAULT 0,
    defeats INTEGER DEFAULT 0,
    pve_battles INTEGER DEFAULT 0,
    pvp_battles INTEGER DEFAULT 0,
    boss_kills INTEGER DEFAULT 0,
    -- 累计伤害
    total_damage_dealt INTEGER DEFAULT 0,
    total_physical_damage INTEGER DEFAULT 0,
    total_magic_damage INTEGER DEFAULT 0,
    total_crit_damage INTEGER DEFAULT 0,
    total_crit_count INTEGER DEFAULT 0,
    highest_damage_single INTEGER DEFAULT 0,
    highest_damage_battle INTEGER DEFAULT 0,
    -- 累计承伤
    total_damage_taken INTEGER DEFAULT 0,
    total_damage_blocked INTEGER DEFAULT 0,
    total_damage_absorbed INTEGER DEFAULT 0,
    total_dodge_count INTEGER DEFAULT 0,
    -- 累计治疗
    total_healing_done INTEGER DEFAULT 0,
    total_healing_received INTEGER DEFAULT 0,
    total_overhealing INTEGER DEFAULT 0,
    highest_healing_single INTEGER DEFAULT 0,
    highest_healing_battle INTEGER DEFAULT 0,
    -- 击杀与死亡
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    kill_streak_best INTEGER DEFAULT 0,
    current_kill_streak INTEGER DEFAULT 0,
    -- 技能使用
    total_skill_uses INTEGER DEFAULT 0,
    total_skill_hits INTEGER DEFAULT 0,
    -- 资源统计
    total_resource_used INTEGER DEFAULT 0,
    total_rounds INTEGER DEFAULT 0,
    total_battle_time INTEGER DEFAULT 0,
    -- 最后更新
    last_battle_at DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);
```

---

### 18. battle_skill_breakdown - 战斗技能明细表

> 📌 **技能使用明细**: 记录每场战斗中各技能的使用和效果

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| battle_id | INTEGER | NOT NULL FK | 战斗记录ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| skill_id | VARCHAR(32) | NOT NULL FK | 技能ID |
| use_count | INTEGER | DEFAULT 0 | 使用次数 |
| hit_count | INTEGER | DEFAULT 0 | 命中次数 |
| crit_count | INTEGER | DEFAULT 0 | 暴击次数 |
| total_damage | INTEGER | DEFAULT 0 | 造成总伤害 |
| total_healing | INTEGER | DEFAULT 0 | 造成总治疗 |
| resource_cost | INTEGER | DEFAULT 0 | 总消耗能量 |

```sql
CREATE TABLE battle_skill_breakdown (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    use_count INTEGER DEFAULT 0,
    hit_count INTEGER DEFAULT 0,
    crit_count INTEGER DEFAULT 0,
    total_damage INTEGER DEFAULT 0,
    total_healing INTEGER DEFAULT 0,
    resource_cost INTEGER DEFAULT 0,
    FOREIGN KEY (battle_id) REFERENCES battle_records(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

CREATE INDEX idx_skill_breakdown_battle ON battle_skill_breakdown(battle_id);
CREATE INDEX idx_skill_breakdown_char ON battle_skill_breakdown(character_id);
```

---

### 19. daily_statistics - 每日统计汇总表

> 📌 **每日快照**: 记录每日战斗数据，便于趋势分析

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| stat_date | DATE | NOT NULL | 统计日期 |
| battles_count | INTEGER | DEFAULT 0 | 战斗次数 |
| victories | INTEGER | DEFAULT 0 | 胜利次数 |
| defeats | INTEGER | DEFAULT 0 | 失败次数 |
| total_damage | INTEGER | DEFAULT 0 | 总伤害 |
| total_healing | INTEGER | DEFAULT 0 | 总治疗 |
| total_damage_taken | INTEGER | DEFAULT 0 | 总承伤 |
| exp_gained | INTEGER | DEFAULT 0 | 获得经验 |
| gold_gained | INTEGER | DEFAULT 0 | 获得金币 |
| play_time | INTEGER | DEFAULT 0 | 游戏时长(秒) |
| kills | INTEGER | DEFAULT 0 | 击杀数 |
| deaths | INTEGER | DEFAULT 0 | 死亡数 |

```sql
CREATE TABLE daily_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    stat_date DATE NOT NULL,
    battles_count INTEGER DEFAULT 0,
    victories INTEGER DEFAULT 0,
    defeats INTEGER DEFAULT 0,
    total_damage INTEGER DEFAULT 0,
    total_healing INTEGER DEFAULT 0,
    total_damage_taken INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    play_time INTEGER DEFAULT 0,
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, stat_date)
);

CREATE INDEX idx_daily_stats_user ON daily_statistics(user_id);
CREATE INDEX idx_daily_stats_date ON daily_statistics(stat_date DESC);
```

---

### 数据分析查询命令

> 💡 玩家可通过游戏内命令查询战斗数据

| 命令 | 说明 | 示例 |
|-----|------|-----|
| `/stats` | 查看当前角色生涯统计 | 显示总伤害、总治疗、胜率等 |
| `/stats [角色名]` | 查看指定角色统计 | `/stats 小明` |
| `/team stats` | 查看全队统计对比 | 各角色DPS/HPS对比 |
| `/battle last` | 查看上一场战斗详情 | 详细伤害/治疗分解 |
| `/battle history [N]` | 查看最近N场战斗 | `/battle history 10` |
| `/dps` | 查看实时DPS排行 | 当前战斗中的输出排名 |
| `/hps` | 查看实时HPS排行 | 当前战斗中的治疗排名 |
| `/skill stats` | 查看技能使用统计 | 各技能伤害/命中率 |
| `/trend [天数]` | 查看数据趋势 | `/trend 7` 最近7天趋势 |

### 数据可视化面板

```
╔═══════════════════════════════════════════════════════════════════════════╗
║                        战斗数据面板 - 最近一战                              ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                           ║
║  📊 伤害输出                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────┐ ║
║  │ 法师·火焰   ████████████████████████████████████████  45,230 (38%)│ ║
║  │ 盗贼·暗影   ██████████████████████████████           32,100 (27%)│ ║
║  │ 战士·破坏   ████████████████████████                 25,890 (22%)│ ║
║  │ 猎人·精准   ████████████████                         15,432 (13%)│ ║
║  └─────────────────────────────────────────────────────────────────────┘ ║
║                                                                           ║
║  🛡️ 伤害承受                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────┐ ║
║  │ 战士·守护   ████████████████████████████████████████  28,500 (55%)│ ║
║  │ 盗贼·暗影   █████████████                             8,230 (16%)│ ║
║  │ 猎人·精准   ████████████                              7,850 (15%)│ ║
║  │ 法师·火焰   ███████                                   4,520  (9%)│ ║
║  │ 牧师·神圣   ███                                       2,100  (4%)│ ║
║  └─────────────────────────────────────────────────────────────────────┘ ║
║                                                                           ║
║  💚 治疗输出                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────┐ ║
║  │ 牧师·神圣   ████████████████████████████████████████  32,100 (92%)│ ║
║  │ 战士·守护   ███                                       2,800  (8%)│ ║
║  │ (自我治疗)                                                         │ ║
║  └─────────────────────────────────────────────────────────────────────┘ ║
║                                                                           ║
║  📈 关键数据                                                               ║
║  ├─ 战斗时长: 2分35秒 (38回合)                                            ║
║  ├─ 团队DPS: 3,068/回合                                                   ║
║  ├─ 团队HPS: 922/回合                                                     ║
║  ├─ 最高单次: 法师·炎爆术 暴击 2,456                                      ║
║  └─ 总过量治疗: 4,230 (13%)                                               ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 战术优化建议

基于战斗数据，系统可以提供智能建议：

| 数据表现 | 可能问题 | 优化建议 |
|---------|---------|---------|
| 坦克承伤<40% | 仇恨不稳定 | 调整战斗策略，优先使用嘲讽 |
| DPS差距大 | 装备/等级差异 | 优先培养落后角色 |
| 过量治疗>30% | 治疗时机不佳 | 调整治疗触发条件(HP<70%) |
| 暴击率<10% | 敏捷/智力不足 | 提升对应属性或装备 |
| 治疗者频繁死亡 | 仇恨问题或位置 | 调整治疗策略，减少过量治疗 |
| 技能命中率低 | 命中属性不足 | 攻击比自己低级的怪物，或提升命中 |

---

## 🔗 ER 关系图

```
                                    ┌──────────────┐
                                    │    users     │
                                    └──────┬───────┘
                                           │ 1:N
                                           ▼
┌─────────────┐    N:1    ┌──────────────────────────┐    1:N    ┌─────────────┐
│   races     │◄──────────│       characters         │──────────►│  inventory  │
└─────────────┘           └─────────────┬────────────┘           └─────────────┘
                                        │                               │
┌─────────────┐    N:1                  │                               │ N:1
│  classes    │◄────────────────────────┤                               ▼
└─────────────┘                         │                        ┌─────────────┐
                                        │                        │    items    │
                          ┌─────────────┼─────────────┐          └──────▲──────┘
                          │             │             │                 │
                          ▼             ▼             ▼                 │
                   ┌────────────┐ ┌───────────┐ ┌───────────┐          │
                   │ equipment  │ │char_skills│ │strategies │          │
                   └─────┬──────┘ └─────┬─────┘ └───────────┘          │
                         │              │                              │
                         │ N:1          │ N:1                          │
                         ▼              ▼                              │
                  ┌─────────────┐ ┌─────────────┐                      │
                  │    items    │ │   skills    │                      │
                  └─────────────┘ └─────────────┘                      │
                                                                       │
                  ┌─────────────┐ 1:N ┌─────────────┐ N:1 ┌────────────┘
                  │    zones    │────►│  monsters   │────►│monster_drops│
                  └─────────────┘     └─────────────┘     └─────────────┘
```

---

## 📊 初始数据

### 种族数据

```sql
INSERT INTO races (id, name, faction, description, strength_mod, agility_mod, intellect_mod, stamina_mod, spirit_mod) VALUES
('human', '人类', 'alliance', '适应力强的种族，各项属性平衡。', 1, 0, 0, 0, 1),
('dwarf', '矮人', 'alliance', '坚韧的山地种族，擅长近战和工艺。', 2, 0, 0, 2, 0),
('nightelf', '暗夜精灵', 'alliance', '古老的精灵种族，与自然和谐共存。', 0, 2, 0, 0, 1),
('gnome', '侏儒', 'alliance', '聪明的小型种族，擅长魔法和机械。', 0, 0, 3, 0, 0),
('orc', '兽人', 'horde', '强壮的战士种族，崇尚力量和荣耀。', 3, 0, 0, 1, 0),
('undead', '亡灵', 'horde', '不死的存在，对暗影魔法有天赋。', 0, 0, 2, 0, 2),
('tauren', '牛头人', 'horde', '高大温和的种族，与大地之母相连。', 2, 0, 0, 2, 1),
('troll', '巨魔', 'horde', '敏捷的丛林种族，拥有快速再生能力。', 0, 2, 0, 1, 0);
```

### 职业数据

```sql
INSERT INTO classes (id, name, description, role, primary_stat, base_hp, base_mp, hp_per_level, mp_per_level, base_strength, base_agility, base_intellect, base_stamina, base_spirit) VALUES
('warrior', '战士', '近战格斗专家，可以承受大量伤害。', 'tank', 'strength', 120, 20, 12, 2, 15, 10, 5, 14, 8),
('paladin', '圣骑士', '神圣战士，可以治疗和保护盟友。', 'tank', 'strength', 110, 60, 10, 6, 13, 8, 10, 13, 12),
('hunter', '猎人', '远程物理攻击者，与宠物并肩作战。', 'dps', 'agility', 90, 40, 8, 4, 8, 15, 8, 10, 10),
('rogue', '盗贼', '潜行刺客，擅长连击和爆发伤害。', 'dps', 'agility', 85, 50, 7, 5, 10, 16, 6, 9, 8),
('priest', '牧师', '治疗者和暗影施法者。', 'healer', 'intellect', 70, 100, 5, 12, 5, 6, 15, 8, 16),
('mage', '法师', '强大的奥术施法者，擅长范围伤害。', 'dps', 'intellect', 65, 120, 4, 15, 4, 6, 18, 6, 12),
('warlock', '术士', '黑暗魔法师，召唤恶魔作战。', 'dps', 'intellect', 75, 110, 5, 13, 5, 6, 17, 8, 10),
('druid', '德鲁伊', '自然的守护者，可变形为多种形态。', 'dps', 'intellect', 85, 80, 7, 10, 10, 10, 13, 10, 12),
('shaman', '萨满', '元素的操控者，可治疗和增益。', 'dps', 'intellect', 90, 90, 8, 10, 12, 8, 14, 11, 12);
```

---

## 📈 索引策略

| 表 | 索引 | 用途 |
|---|------|-----|
| users | username | 登录查询 |
| characters | user_id | 用户角色查询 |
| characters | level | 排行榜 |
| inventory | character_id | 背包查询 |
| equipment | character_id | 装备查询 |
| character_skills | character_id | 技能查询 |
| monsters | zone_id | 区域怪物查询 |
| monsters | level | 等级匹配 |
| battle_strategies | character_id | 策略查询 |

---

## 🔒 数据完整性

1. **外键约束** - 所有关联使用外键，CASCADE删除
2. **唯一约束** - 用户名、装备槽位等
3. **默认值** - 所有数值字段设置合理默认值
4. **触发器** - 自动更新 updated_at 时间戳

```sql
-- 自动更新 updated_at 触发器
CREATE TRIGGER update_character_timestamp 
AFTER UPDATE ON characters
BEGIN
    UPDATE characters SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

---

## 📝 注意事项

1. **SQLite 特性**
   - 使用 WAL 模式提高并发性能
   - VARCHAR 实际等同于 TEXT
   - 外键需要手动开启: `PRAGMA foreign_keys = ON`

2. **扩展性考虑**
   - ID使用VARCHAR便于配置数据管理
   - 预留了JSON字段用于灵活扩展
   - 统计字段分离，避免频繁更新主表

3. **迁移到其他数据库**
   - 表结构兼容 MySQL/PostgreSQL
   - 需调整自增语法和部分数据类型

---

## 🎮 后期玩法系统

### 一、无尽深渊系统

> 📌 **核心后期挑战**: 无限层数的挑战塔，层数越高难度越大，奖励越丰富

#### abyss_config - 深渊配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| floor | INTEGER | PRIMARY KEY | 层数 |
| monster_level_base | INTEGER | NOT NULL | 怪物基础等级 |
| monster_level_growth | REAL | DEFAULT 0.5 | 每层等级增长 |
| monster_hp_mult | REAL | DEFAULT 1.0 | 怪物HP倍率 |
| monster_atk_mult | REAL | DEFAULT 1.0 | 怪物攻击倍率 |
| reward_exp_mult | REAL | DEFAULT 1.0 | 经验奖励倍率 |
| reward_gold_mult | REAL | DEFAULT 1.0 | 金币奖励倍率 |
| special_reward | TEXT | | 特殊奖励(JSON) |
| boss_id | VARCHAR(32) | | Boss怪物ID(每10层) |

```sql
CREATE TABLE abyss_config (
    floor INTEGER PRIMARY KEY,
    monster_level_base INTEGER NOT NULL,
    monster_level_growth REAL DEFAULT 0.5,
    monster_hp_mult REAL DEFAULT 1.0,
    monster_atk_mult REAL DEFAULT 1.0,
    reward_exp_mult REAL DEFAULT 1.0,
    reward_gold_mult REAL DEFAULT 1.0,
    special_reward TEXT,
    boss_id VARCHAR(32),
    FOREIGN KEY (boss_id) REFERENCES monsters(id)
);
```

#### abyss_progress - 深渊进度表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| user_id | INTEGER | NOT NULL FK UNIQUE | 用户ID |
| highest_floor | INTEGER | DEFAULT 0 | 最高通关层数 |
| current_floor | INTEGER | DEFAULT 1 | 当前挑战层数 |
| weekly_attempts | INTEGER | DEFAULT 0 | 本周已挑战次数 |
| weekly_reset_at | DATETIME | | 周重置时间 |
| total_clears | INTEGER | DEFAULT 0 | 总通关次数 |
| best_time | INTEGER | | 最快通关时间(秒) |

```sql
CREATE TABLE abyss_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    highest_floor INTEGER DEFAULT 0,
    current_floor INTEGER DEFAULT 1,
    weekly_attempts INTEGER DEFAULT 0,
    weekly_reset_at DATETIME,
    total_clears INTEGER DEFAULT 0,
    best_time INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_abyss_highest ON abyss_progress(highest_floor DESC);
```

---

### 二、装备系统：词缀 + 进化链

> 📌 **核心设计**: 装备掉落时随机生成词缀，通过进化系统获得更强形态和传说效果

#### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          装备系统架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   装备掉落 ──→ 随机词缀 ──→ 强化升级 ──→ 进化分支 ──→ 传说形态              │
│      │           │           │           │           │                     │
│   基础属性    前缀+后缀    +1~+25      选择路线    独特效果                   │
│                                                                             │
│   品质决定词缀数量:                                                          │
│   ⬜白(0) → 🟩绿(1) → 🟦蓝(2) → 🟪紫(3) → 🟧橙(4) → 🟥红(4+传说)            │
│                                                                             │
│   进化阶段:                                                                  │
│   Ⅰ基础 → Ⅱ精炼 → Ⅲ进化(选分支) → Ⅳ觉醒 → Ⅴ传说                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 2.1 equipment_instance - 装备实例表

> 📌 记录玩家获得的每一件装备及其词缀

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 装备实例ID |
| item_id | VARCHAR(32) | NOT NULL FK | 基础物品ID |
| owner_id | INTEGER | NOT NULL FK | 拥有者用户ID |
| character_id | INTEGER | FK | 装备者角色ID (NULL=背包中) |
| slot | VARCHAR(16) | | 装备槽位 |
| quality | VARCHAR(16) | NOT NULL | 品质: common/uncommon/rare/epic/legendary/mythic |
| enhance_level | INTEGER | DEFAULT 0 | 强化等级 (0-25) |
| evolution_stage | INTEGER | DEFAULT 1 | 进化阶段 (1-5) |
| evolution_path | VARCHAR(32) | | 进化路线: fire/frost/lightning/holy/shadow/nature/physical |
| prefix_id | VARCHAR(32) | FK | 前缀词缀ID |
| prefix_value | REAL | | 前缀数值 (词缀效果的具体数值) |
| suffix_id | VARCHAR(32) | FK | 后缀词缀ID |
| suffix_value | REAL | | 后缀数值 |
| bonus_affix_1 | VARCHAR(32) | FK | 额外词缀1 (紫色+) |
| bonus_affix_1_value | REAL | | 额外词缀1数值 |
| bonus_affix_2 | VARCHAR(32) | FK | 额外词缀2 (橙色+) |
| bonus_affix_2_value | REAL | | 额外词缀2数值 |
| legendary_effect_id | VARCHAR(32) | FK | 传说效果ID (红色品质) |
| acquired_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 获得时间 |
| is_locked | INTEGER | DEFAULT 0 | 是否锁定 (防误分解) |

```sql
CREATE TABLE equipment_instance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id VARCHAR(32) NOT NULL,
    owner_id INTEGER NOT NULL,
    character_id INTEGER,
    slot VARCHAR(16),
    quality VARCHAR(16) NOT NULL DEFAULT 'common',
    enhance_level INTEGER DEFAULT 0,
    evolution_stage INTEGER DEFAULT 1,
    evolution_path VARCHAR(32),
    prefix_id VARCHAR(32),
    prefix_value REAL,
    suffix_id VARCHAR(32),
    suffix_value REAL,
    bonus_affix_1 VARCHAR(32),
    bonus_affix_1_value REAL,
    bonus_affix_2 VARCHAR(32),
    bonus_affix_2_value REAL,
    legendary_effect_id VARCHAR(32),
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_locked INTEGER DEFAULT 0,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE SET NULL,
    FOREIGN KEY (prefix_id) REFERENCES affixes(id),
    FOREIGN KEY (suffix_id) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_1) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_2) REFERENCES affixes(id),
    FOREIGN KEY (legendary_effect_id) REFERENCES legendary_effects(id)
);

CREATE INDEX idx_equipment_owner ON equipment_instance(owner_id);
CREATE INDEX idx_equipment_character ON equipment_instance(character_id);
CREATE INDEX idx_equipment_quality ON equipment_instance(quality);
```

#### 2.2 affixes - 词缀配置表

> 📌 定义所有可能的装备词缀

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 词缀ID |
| name | VARCHAR(32) | NOT NULL | 词缀名称 |
| type | VARCHAR(16) | NOT NULL | 类型: prefix/suffix |
| slot_type | VARCHAR(16) | | 适用槽位: weapon/armor/accessory/all |
| rarity | VARCHAR(16) | NOT NULL | 稀有度: common/uncommon/rare/epic |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_stat | VARCHAR(32) | | 影响的属性 |
| min_value | REAL | NOT NULL | 最小数值 |
| max_value | REAL | NOT NULL | 最大数值 |
| value_type | VARCHAR(16) | NOT NULL | 数值类型: flat/percent |
| description | TEXT | | 描述模板 (用{value}占位) |
| level_required | INTEGER | DEFAULT 1 | 最低出现等级 |

```sql
CREATE TABLE affixes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    type VARCHAR(16) NOT NULL,
    slot_type VARCHAR(16) DEFAULT 'all',
    rarity VARCHAR(16) NOT NULL DEFAULT 'common',
    effect_type VARCHAR(32) NOT NULL,
    effect_stat VARCHAR(32),
    min_value REAL NOT NULL,
    max_value REAL NOT NULL,
    value_type VARCHAR(16) NOT NULL DEFAULT 'flat',
    description TEXT,
    level_required INTEGER DEFAULT 1
);

CREATE INDEX idx_affixes_type ON affixes(type);
CREATE INDEX idx_affixes_rarity ON affixes(rarity);
CREATE INDEX idx_affixes_slot ON affixes(slot_type);
```

#### 词缀列表

**前缀 (攻击/属性向):**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_sharp | 锋利的 | 攻击力+ | 2~5 | 普通 | 武器 |
| affix_fiery | 炽热的 | 火焰伤害+ | 1~4 | 普通 | 武器 |
| affix_frozen | 冰霜的 | 冰霜伤害+ | 1~4 | 普通 | 武器 |
| affix_charged | 雷击的 | 雷电伤害+ | 1~4 | 普通 | 武器 |
| affix_holy | 神圣的 | 神圣伤害+ | 2~5 | 精良 | 武器 |
| affix_vampiric | 吸血鬼的 | 生命偷取% | 2~5 | 稀有 | 武器 |
| affix_devastating | 毁灭的 | 攻击力+% | 15~25 | 史诗 | 武器 |
| affix_sturdy | 坚固的 | 防御力+ | 2~5 | 普通 | 防具 |
| affix_vital | 活力的 | 生命值+ | 5~15 | 普通 | 防具 |
| affix_scholarly | 智者的 | 智力+ | 2~4 | 精良 | 防具 |
| affix_unyielding | 不屈的 | 受伤减免% | 3~8 | 稀有 | 防具 |

**后缀 (特殊效果向):**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_strength | of 力量 | 力量+ | 1~3 | 普通 | 全部 |
| affix_of_agility | of 敏捷 | 敏捷+ | 1~3 | 普通 | 全部 |
| affix_of_haste | of 迅捷 | 攻击速度+% | 5~15 | 精良 | 武器 |
| affix_of_piercing | of 穿刺 | 无视防御% | 5~15 | 精良 | 武器 |
| affix_of_crit | of 暴击 | 暴击率+% | 3~8 | 稀有 | 武器 |
| affix_of_lethality | of 致命 | 暴击伤害+% | 10~25 | 稀有 | 武器 |
| affix_of_leech | of 吸血 | 伤害转HP% | 2~4 | 稀有 | 武器 |
| affix_of_blocking | of 守护 | 格挡率+% | 5~10 | 精良 | 盾牌 |
| affix_of_thorns | of 反射 | 反弹伤害% | 5~15 | 稀有 | 防具 |
| affix_of_regen | of 再生 | 每回合恢复HP | 1~3 | 稀有 | 防具 |
| affix_of_wisdom | of 智慧 | 法力恢复+% | 10~20 | 精良 | 防具 |

---

#### 2.3 evolution_paths - 进化路线配置表

> 📌 定义装备可选择的进化分支

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 路线ID |
| name | VARCHAR(32) | NOT NULL | 路线名称 |
| element | VARCHAR(16) | NOT NULL | 元素类型 |
| description | TEXT | | 描述 |
| slot_type | VARCHAR(16) | NOT NULL | 适用槽位: weapon/armor |
| stat_bonus_type | VARCHAR(32) | | 属性加成类型 |
| stat_bonus_value | REAL | | 属性加成数值 |
| special_effect | TEXT | | 特殊效果描述 |
| material_required | TEXT | | 所需材料 (JSON) |

```sql
CREATE TABLE evolution_paths (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    element VARCHAR(16) NOT NULL,
    description TEXT,
    slot_type VARCHAR(16) NOT NULL,
    stat_bonus_type VARCHAR(32),
    stat_bonus_value REAL,
    special_effect TEXT,
    material_required TEXT
);
```

**武器进化路线:**

| 路线 | 元素 | 核心加成 | 特殊效果 |
|-----|------|---------|---------|
| 🔥 烈焰 | fire | 火焰伤害+50% | 攻击灼烧敌人 |
| ❄️ 霜寒 | frost | 冰霜伤害+50% | 攻击减速敌人 |
| ⚡ 雷霆 | lightning | 攻速+20%, 雷伤+40% | 伤害连锁跳跃 |
| ✨ 神圣 | holy | 圣伤+40%, 治疗+15% | 攻击回复生命 |
| 🌑 暗影 | shadow | 暗伤+50%, 吸血+5% | 伤害转化HP |
| 🌿 自然 | nature | 自然伤+40%, 再生+20% | 持续恢复HP |
| ⚔️ 物理 | physical | 物伤+30%, 穿透+15% | 无视部分护甲 |

**防具进化路线:**

| 路线 | 定位 | 核心加成 | 特殊效果 |
|-----|------|---------|---------|
| 🛡️ 守护 | 坦克 | 防御+30%, 生命+20% | 受伤减免 |
| 🌵 荆棘 | 反伤 | 防御+15%, 反伤+25% | 被攻击时反弹伤害 |
| 💨 迅捷 | 闪避 | 闪避+20%, 攻速+15% | 闪避后加速 |
| 💚 再生 | 续航 | 生命+15%, 回复+50% | 每回合恢复HP |

---

#### 2.4 legendary_effects - 传说效果表

> 📌 进化到最终阶段解锁的独特效果

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 传说效果ID |
| name | VARCHAR(32) | NOT NULL | 效果名称 |
| description | TEXT | NOT NULL | 效果描述 |
| slot_type | VARCHAR(16) | NOT NULL | 适用槽位 |
| evolution_path | VARCHAR(32) | | 关联进化路线 |
| trigger_type | VARCHAR(32) | | 触发类型: on_hit/on_kill/on_damaged/passive |
| trigger_chance | REAL | DEFAULT 1.0 | 触发概率 |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_value | REAL | | 效果数值 |
| cooldown | INTEGER | DEFAULT 0 | 冷却回合 |

```sql
CREATE TABLE legendary_effects (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT NOT NULL,
    slot_type VARCHAR(16) NOT NULL,
    evolution_path VARCHAR(32),
    trigger_type VARCHAR(32),
    trigger_chance REAL DEFAULT 1.0,
    effect_type VARCHAR(32) NOT NULL,
    effect_value REAL,
    cooldown INTEGER DEFAULT 0,
    FOREIGN KEY (evolution_path) REFERENCES evolution_paths(id)
);
```

**传说武器效果:**

| ID | 名称 | 效果 | 来源路线 |
|---|-----|------|---------|
| legend_inferno | 地狱烈焰 | 攻击使敌人灼烧3回合，每回合3点火伤 | 烈焰 |
| legend_frostmourne | 霜之哀伤 | 击杀敌人后冰冻周围敌人1回合 | 霜寒 |
| legend_thunderfury | 雷霆之怒 | 20%几率触发闪电链，最多跳跃3目标 | 雷霆 |
| legend_ashbringer | 灰烬使者 | 攻击时恢复自身5%最大生命 | 神圣 |
| legend_shadowmourne | 暗影之殇 | 暴击时吸取敌人10%当前生命 | 暗影 |
| legend_earthshatter | 大地粉碎 | 攻击叠加标记，5层后引爆额外伤害 | 自然 |
| legend_gorehowl | 血吼 | 对精英和Boss伤害+50% | 物理 |

**传说防具效果:**

| ID | 名称 | 效果 | 来源路线 |
|---|-----|------|---------|
| legend_immortal | 不灭意志 | 首次致死伤害免疫 (每场战斗1次) | 守护 |
| legend_retribution | 复仇之刺 | 反弹50%受到的物理伤害 | 荆棘 |
| legend_shadowstep | 暗影步 | 闪避成功后下次攻击必暴击 | 迅捷 |
| legend_lifesource | 生命之泉 | HP<30%时每回合恢复10% HP | 再生 |

---

#### 2.5 进化阶段与强化

**进化阶段:**

| 阶段 | 名称 | 角色等级 | 强化上限 | 进化材料 | 特点 |
|-----|------|---------|---------|---------|------|
| Ⅰ | 基础 | 1-15 | +5 | - | 初始形态 |
| Ⅱ | 精炼 | 16-30 | +10 | 精炼石×5 | 属性提升20% |
| Ⅲ | 进化 | 31-45 | +15 | 进化石×3 + 元素核心 | **选择分支** |
| Ⅳ | 觉醒 | 46-55 | +20 | 觉醒石×1 + 稀有材料 | 解锁特殊效果 |
| Ⅴ | 传说 | 56-60 | +25 | 传说碎片×5 | 获得传说效果 |

**强化效果:**

```
强化属性加成 = 基础属性 × (1 + 强化等级 × 0.02)

示例:
+0: 100%
+5: 110% (+10%)
+10: 120% (+20%)
+15: 130% (+30%)
+20: 140% (+40%)
+25: 150% (+50%)
```

**词缀继承:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          进化时词缀处理                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ✅ 所有词缀完整保留                                                         │
│  ✅ 词缀数值随进化阶段提升 (+15%/阶段)                                        │
│  ✅ 品质可能提升 (蓝→紫 30%几率, 紫→橙 10%几率)                               │
│                                                                             │
│  示例:                                                                       │
│  Lv15 炽热的(+2火伤) 钢剑 of迅捷(+10%攻速) [蓝]                              │
│         ↓ 进化到阶段Ⅲ                                                       │
│  Lv30 炽热的(+3火伤) 火焰秘银剑 of迅捷(+12%攻速) [紫]                         │
│         ↓ 进化到阶段Ⅴ                                                       │
│  Lv50 炽热的(+4火伤) 烈焰之刃 of迅捷(+14%攻速) [橙]                           │
│        + 传说效果: 攻击使敌人灼烧                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

#### 2.6 装备掉落系统

> 📌 **设计理念**: 少而精的掉落，每次掉落都值得关注，任何区域都有惊喜可能

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          掉落设计理念                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ❌ 避免                                 ✅ 追求                            │
│   ├─ 满背包白色/绿色垃圾                   ├─ 每次掉落都值得关注               │
│   ├─ 频繁整理背包的烦恼                    ├─ "叮！"的惊喜感                   │
│   ├─ 掉落=分解材料                        ├─ 低级区也有小概率出神装             │
│   └─ 数量堆砌，毫无期待                    └─ 稀有但有意义的掉落                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### 基础掉落率 (每场战斗)

| 怪物类型 | 掉落几率 | 说明 |
|---------|---------|------|
| 普通怪 | 5% | 基础掉率，大部分战斗不掉装备 |
| 精英怪 | 15% | 精英更值得挑战 |
| Boss | 50% | Boss有一半几率掉装备 |
| 深渊Boss | 100% | 每10层Boss必掉 |

##### 品质分布 (当装备掉落时)

| 品质 | 颜色 | 词缀数 | 掉落率 | 说明 |
|-----|------|-------|-------|------|
| 普通 (Common) | ⬜ 白 | 0 | 30% | 大幅减少垃圾 |
| 优秀 (Uncommon) | 🟩 绿 | 1 | 35% | 单词缀起步装 |
| 精良 (Rare) | 🟦 蓝 | 2 | 25% | 主力装备 |
| 稀有 (Epic) | 🟪 紫 | 3 | 8% | 有培养价值 |
| 史诗 (Legendary) | 🟧 橙 | 4 | 1.8% | 稀有可期待 |
| 传说 (Mythic) | 🟥 红 | 4+传说效果 | 0.2% | 终极追求 |

##### 奇迹掉落系统

> 📌 **核心**: 任何区域都有极小概率掉落超出等级的顶级装备

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          奇迹掉落                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   场景: 玩家在1级新手村刷怪                                                   │
│                                                                             │
│   正常掉落: 只能掉落1-5级的装备                                               │
│                                                                             │
│   奇迹触发 (0.5%): 无视等级限制，从全装备池随机                                │
│                    → 可能掉落60级传说装备！                                    │
│                                                                             │
│   📢 全服公告: "玩家【小明】在艾尔文森林获得了【霜之哀伤】！"                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

| 区域等级 | 奇迹触发率 | 说明 |
|---------|-----------|------|
| 1-10级区 | 0.5% | 新手也有希望 |
| 11-30级区 | 0.3% | 逐渐降低 |
| 31-50级区 | 0.1% | 本身装备等级已高 |
| 51-60级区 | 0% | 已是顶级区域 |

##### 保底机制

> 📌 防止长时间无掉落的挫败感

| 连续无掉落次数 | 效果 |
|--------------|------|
| 1-19 | 正常掉率 |
| 20-29 | 掉率×2 |
| 30-39 | 掉率×4 |
| 40+ | 保底掉落 (🟦精良或以上) |

##### 背包管理

| 功能 | 说明 |
|-----|------|
| **自动分解** | 可设置自动分解⬜白色装备 |
| **快速分解** | 一键分解所有未锁定的低品质装备 |
| **背包上限** | 100格，满时自动分解最低品质 |
| **掉落预览** | 战斗后显示掉落，可选拾取或直接分解 |

##### drop_config - 掉落配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 配置ID |
| monster_type | VARCHAR(16) | NOT NULL | 怪物类型: normal/elite/boss/abyss_boss |
| base_drop_rate | REAL | NOT NULL | 基础掉落率 |
| quality_weights | TEXT | NOT NULL | 品质权重 (JSON) |
| miracle_rate | REAL | DEFAULT 0 | 奇迹掉落率 |
| pity_threshold | INTEGER | DEFAULT 40 | 保底触发次数 |
| pity_min_quality | VARCHAR(16) | DEFAULT 'rare' | 保底最低品质 |

```sql
CREATE TABLE IF NOT EXISTS drop_config (
    id VARCHAR(32) PRIMARY KEY,
    monster_type VARCHAR(16) NOT NULL,
    base_drop_rate REAL NOT NULL,
    quality_weights TEXT NOT NULL,
    miracle_rate REAL DEFAULT 0,
    pity_threshold INTEGER DEFAULT 40,
    pity_min_quality VARCHAR(16) DEFAULT 'rare'
);
```

##### user_drop_pity - 玩家保底计数表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| user_id | INTEGER | PRIMARY KEY FK | 用户ID |
| no_drop_count | INTEGER | DEFAULT 0 | 连续无掉落次数 |
| last_drop_at | DATETIME | | 上次掉落时间 |
| total_drops | INTEGER | DEFAULT 0 | 总掉落次数 |
| miracle_drops | INTEGER | DEFAULT 0 | 奇迹掉落次数 |

```sql
CREATE TABLE IF NOT EXISTS user_drop_pity (
    user_id INTEGER PRIMARY KEY,
    no_drop_count INTEGER DEFAULT 0,
    last_drop_at DATETIME,
    total_drops INTEGER DEFAULT 0,
    miracle_drops INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

### 三、收集与成就系统

#### 3.1 achievements - 成就配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 成就ID |
| name | VARCHAR(64) | NOT NULL | 成就名称 |
| description | TEXT | | 描述 |
| category | VARCHAR(32) | NOT NULL | 分类: combat/explore/collect/social |
| condition_type | VARCHAR(32) | NOT NULL | 条件类型 |
| condition_value | INTEGER | NOT NULL | 条件数值 |
| points | INTEGER | DEFAULT 10 | 成就点数 |
| reward_type | VARCHAR(32) | | 奖励类型 |
| reward_value | TEXT | | 奖励内容(JSON) |
| icon | VARCHAR(64) | | 图标 |
| is_hidden | INTEGER | DEFAULT 0 | 是否隐藏成就 |

```sql
CREATE TABLE achievements (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    category VARCHAR(32) NOT NULL,
    condition_type VARCHAR(32) NOT NULL,
    condition_value INTEGER NOT NULL,
    points INTEGER DEFAULT 10,
    reward_type VARCHAR(32),
    reward_value TEXT,
    icon VARCHAR(64),
    is_hidden INTEGER DEFAULT 0
);

CREATE INDEX idx_achievements_category ON achievements(category);
```

**成就分类:**
- `combat` - 战斗成就 (击杀数、伤害、连胜...)
- `explore` - 探索成就 (区域、副本、深渊层数...)
- `collect` - 收集成就 (装备、图鉴、宠物...)
- `social` - 社交成就 (组队、公会、PvP...)
- `special` - 特殊成就 (限时、首杀、极限挑战...)

#### 3.2 user_achievements - 玩家成就表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| achievement_id | VARCHAR(32) | NOT NULL FK | 成就ID |
| progress | INTEGER | DEFAULT 0 | 当前进度 |
| completed_at | DATETIME | | 完成时间 |
| rewarded | INTEGER | DEFAULT 0 | 是否已领奖 |

```sql
CREATE TABLE user_achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    achievement_id VARCHAR(32) NOT NULL,
    progress INTEGER DEFAULT 0,
    completed_at DATETIME,
    rewarded INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user ON user_achievements(user_id);
```

#### 3.3 codex - 图鉴配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 图鉴ID |
| category | VARCHAR(32) | NOT NULL | 分类: monster/item/boss/zone |
| target_id | VARCHAR(32) | NOT NULL | 关联目标ID |
| name | VARCHAR(64) | NOT NULL | 名称 |
| description | TEXT | | 描述 |
| unlock_condition | TEXT | | 解锁条件 |
| bonus_type | VARCHAR(32) | | 收集奖励类型 |
| bonus_value | REAL | | 收集奖励数值 |

```sql
CREATE TABLE codex (
    id VARCHAR(32) PRIMARY KEY,
    category VARCHAR(32) NOT NULL,
    target_id VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    unlock_condition TEXT,
    bonus_type VARCHAR(32),
    bonus_value REAL
);

CREATE INDEX idx_codex_category ON codex(category);
```

#### 3.4 user_codex - 玩家图鉴表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| codex_id | VARCHAR(32) | NOT NULL FK | 图鉴ID |
| unlock_count | INTEGER | DEFAULT 1 | 解锁次数/击杀次数 |
| first_unlock_at | DATETIME | NOT NULL | 首次解锁时间 |

```sql
CREATE TABLE user_codex (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    codex_id VARCHAR(32) NOT NULL,
    unlock_count INTEGER DEFAULT 1,
    first_unlock_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (codex_id) REFERENCES codex(id),
    UNIQUE(user_id, codex_id)
);

CREATE INDEX idx_user_codex_user ON user_codex(user_id);
```

#### 图鉴收集奖励示例

| 图鉴类型 | 收集进度 | 奖励 |
|---------|---------|------|
| 怪物图鉴 | 10% | 伤害+1% |
| 怪物图鉴 | 50% | 伤害+3% |
| 怪物图鉴 | 100% | 伤害+5% + 称号"百科全书" |
| Boss图鉴 | 首杀任意Boss | 专属称号 |
| Boss图鉴 | 击杀全部Boss | 传说外观 |

---

## ⚔️ 阵营PVP遭遇战系统

> 📌 **核心竞争机制**: 在同一地图挂机的联盟与部落玩家会随机发生PVP遭遇战，胜负影响该地图的阵营效率加成

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          阵营PVP遭遇战系统                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────────────────────────────────────────────────────────────┐     │
│   │                        地图: 艾尔文森林                            │     │
│   │  ┌─────────────────┐                  ┌─────────────────┐        │     │
│   │  │  联盟玩家群体    │    ⚔️ 遭遇战     │   部落玩家群体   │        │     │
│   │  │  ├─ Player_A    │ ◄─────────────► │   ├─ Player_X   │        │     │
│   │  │  ├─ Player_B    │                  │   └─ Player_Y   │        │     │
│   │  │  └─ Player_C    │                  │                 │        │     │
│   │  └─────────────────┘                  └─────────────────┘        │     │
│   │                                                                   │     │
│   │  当前控制方: 联盟 (胜率65%)                                        │     │
│   │  联盟效率: +10% 经验/金币/掉落                                     │     │
│   │  部落效率: -5% 经验/金币/掉落                                      │     │
│   └──────────────────────────────────────────────────────────────────┘     │
│                                                                             │
│   📢 全服公告: [艾尔文森林] 联盟<勇士小明>击败了部落<暗影猎手>！            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 核心机制

| 机制 | 说明 |
|-----|------|
| **遭遇触发** | 同一地图的敌对阵营玩家随机匹配，触发自动PVP战斗 |
| **战斗方式** | 使用玩家设置的战斗策略，自动化对决 |
| **胜负判定** | 击杀对方或对方投降/逃跑 |
| **荣誉奖励** | 胜者获得荣誉值，可用于兑换奖励 |
| **地图控制** | 根据近期胜率计算阵营控制权 |
| **效率加成** | 控制方获得挂机效率提升，被控方效率降低 |
| **全服公告** | 所有PVP战斗结果向全服玩家广播 |

---

### 四、地图阵营控制表

#### zone_faction_control - 地图阵营控制表

> 📌 **地图归属**: 记录每个地图的阵营控制状态和效率加成

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| zone_id | VARCHAR(32) | PRIMARY KEY FK | 区域ID |
| controlling_faction | VARCHAR(16) | | 当前控制阵营: alliance/horde/neutral |
| alliance_wins | INTEGER | DEFAULT 0 | 联盟近期胜场 |
| horde_wins | INTEGER | DEFAULT 0 | 部落近期胜场 |
| alliance_win_rate | REAL | DEFAULT 0.5 | 联盟胜率 |
| control_score | INTEGER | DEFAULT 0 | 控制积分 (正=联盟, 负=部落) |
| efficiency_bonus | REAL | DEFAULT 0 | 控制方效率加成 (0.0-0.2) |
| efficiency_penalty | REAL | DEFAULT 0 | 被控方效率惩罚 (0.0-0.1) |
| last_battle_at | DATETIME | | 最后一次战斗时间 |
| stats_reset_at | DATETIME | | 统计重置时间 (每周) |

```sql
CREATE TABLE zone_faction_control (
    zone_id VARCHAR(32) PRIMARY KEY,
    controlling_faction VARCHAR(16) DEFAULT 'neutral',
    alliance_wins INTEGER DEFAULT 0,
    horde_wins INTEGER DEFAULT 0,
    alliance_win_rate REAL DEFAULT 0.5,
    control_score INTEGER DEFAULT 0,
    efficiency_bonus REAL DEFAULT 0,
    efficiency_penalty REAL DEFAULT 0,
    last_battle_at DATETIME,
    stats_reset_at DATETIME,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);
```

**控制权计算规则:**
```
控制积分变化:
  - 联盟胜利: +1 分
  - 部落胜利: -1 分
  - 积分范围: -100 ~ +100

控制权判定:
  - 积分 >= +20: 联盟控制
  - 积分 <= -20: 部落控制
  - -20 < 积分 < +20: 中立/争夺中

效率加成计算:
  - 控制方: +效率加成 (最高+20%)
  - 被控方: -效率惩罚 (最高-10%)
  - 加成比例 = |积分| / 100 × 最大加成
```

---

### 五、PVP遭遇战记录表

#### pvp_encounters - PVP遭遇战记录表

> 📌 **战斗记录**: 记录每一场PVP遭遇战的详细信息

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 战斗ID |
| zone_id | VARCHAR(32) | NOT NULL FK | 发生区域 |
| attacker_user_id | INTEGER | NOT NULL FK | 进攻方用户ID |
| defender_user_id | INTEGER | NOT NULL FK | 防守方用户ID |
| attacker_faction | VARCHAR(16) | NOT NULL | 进攻方阵营 |
| defender_faction | VARCHAR(16) | NOT NULL | 防守方阵营 |
| winner_user_id | INTEGER | | 胜利方用户ID (NULL=平局) |
| winner_faction | VARCHAR(16) | | 胜利阵营 |
| attacker_team_info | TEXT | | 进攻方队伍信息 (JSON) |
| defender_team_info | TEXT | | 防守方队伍信息 (JSON) |
| battle_rounds | INTEGER | DEFAULT 0 | 战斗回合数 |
| battle_duration | INTEGER | DEFAULT 0 | 战斗时长(秒) |
| attacker_damage_dealt | INTEGER | DEFAULT 0 | 进攻方造成伤害 |
| defender_damage_dealt | INTEGER | DEFAULT 0 | 防守方造成伤害 |
| honor_reward | INTEGER | DEFAULT 0 | 荣誉奖励 |
| battle_log | TEXT | | 战斗日志 (JSON) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 战斗时间 |

```sql
CREATE TABLE pvp_encounters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    zone_id VARCHAR(32) NOT NULL,
    attacker_user_id INTEGER NOT NULL,
    defender_user_id INTEGER NOT NULL,
    attacker_faction VARCHAR(16) NOT NULL,
    defender_faction VARCHAR(16) NOT NULL,
    winner_user_id INTEGER,
    winner_faction VARCHAR(16),
    attacker_team_info TEXT,
    defender_team_info TEXT,
    battle_rounds INTEGER DEFAULT 0,
    battle_duration INTEGER DEFAULT 0,
    attacker_damage_dealt INTEGER DEFAULT 0,
    defender_damage_dealt INTEGER DEFAULT 0,
    honor_reward INTEGER DEFAULT 0,
    battle_log TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (attacker_user_id) REFERENCES users(id),
    FOREIGN KEY (defender_user_id) REFERENCES users(id),
    FOREIGN KEY (winner_user_id) REFERENCES users(id)
);

CREATE INDEX idx_pvp_encounters_zone ON pvp_encounters(zone_id);
CREATE INDEX idx_pvp_encounters_attacker ON pvp_encounters(attacker_user_id);
CREATE INDEX idx_pvp_encounters_defender ON pvp_encounters(defender_user_id);
CREATE INDEX idx_pvp_encounters_time ON pvp_encounters(created_at DESC);
CREATE INDEX idx_pvp_encounters_winner ON pvp_encounters(winner_faction);
```

---

### 六、玩家荣誉表

#### user_honor - 玩家荣誉表

> 📌 **荣誉系统**: 记录玩家PVP战绩和荣誉积累

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| user_id | INTEGER | PRIMARY KEY FK | 用户ID |
| faction | VARCHAR(16) | NOT NULL | 所属阵营 |
| total_honor | INTEGER | DEFAULT 0 | 累计荣誉值 |
| current_honor | INTEGER | DEFAULT 0 | 当前可用荣誉 |
| honor_rank | INTEGER | DEFAULT 0 | 荣誉军衔等级 (0-14) |
| pvp_wins | INTEGER | DEFAULT 0 | PVP胜场 |
| pvp_losses | INTEGER | DEFAULT 0 | PVP败场 |
| pvp_draws | INTEGER | DEFAULT 0 | PVP平局 |
| win_streak | INTEGER | DEFAULT 0 | 当前连胜 |
| best_win_streak | INTEGER | DEFAULT 0 | 最高连胜 |
| total_kills | INTEGER | DEFAULT 0 | 总击杀角色数 |
| total_deaths | INTEGER | DEFAULT 0 | 总死亡角色数 |
| total_damage_dealt | INTEGER | DEFAULT 0 | 总造成伤害 |
| weekly_honor | INTEGER | DEFAULT 0 | 本周荣誉 |
| weekly_reset_at | DATETIME | | 周重置时间 |

```sql
CREATE TABLE user_honor (
    user_id INTEGER PRIMARY KEY,
    faction VARCHAR(16) NOT NULL,
    total_honor INTEGER DEFAULT 0,
    current_honor INTEGER DEFAULT 0,
    honor_rank INTEGER DEFAULT 0,
    pvp_wins INTEGER DEFAULT 0,
    pvp_losses INTEGER DEFAULT 0,
    pvp_draws INTEGER DEFAULT 0,
    win_streak INTEGER DEFAULT 0,
    best_win_streak INTEGER DEFAULT 0,
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    total_damage_dealt INTEGER DEFAULT 0,
    weekly_honor INTEGER DEFAULT 0,
    weekly_reset_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_honor_rank ON user_honor(honor_rank DESC);
CREATE INDEX idx_user_honor_wins ON user_honor(pvp_wins DESC);
CREATE INDEX idx_user_honor_faction ON user_honor(faction);
```

**荣誉军衔体系:**

| 等级 | 联盟军衔 | 部落军衔 | 所需累计荣誉 |
|-----|---------|---------|------------|
| 0 | 无军衔 | 无军衔 | 0 |
| 1 | 列兵 | 斥候 | 100 |
| 2 | 下士 | 步兵 | 500 |
| 3 | 中士 | 中士 | 1,000 |
| 4 | 军士长 | 高级中士 | 2,000 |
| 5 | 准尉 | 一等军士长 | 5,000 |
| 6 | 少尉 | 石卫士 | 10,000 |
| 7 | 中尉 | 血卫士 | 20,000 |
| 8 | 上尉 | 军团士兵 | 35,000 |
| 9 | 少校 | 百夫长 | 50,000 |
| 10 | 中校 | 勇士 | 75,000 |
| 11 | 上校 | 将军 | 100,000 |
| 12 | 准将 | 军阀 | 150,000 |
| 13 | 元帅 | 高阶督军 | 250,000 |
| 14 | 大元帅 | 大督军 | 500,000 |

---

### 七、全服公告表

#### server_announcements - 全服公告表

> 📌 **实时广播**: 记录并推送PVP战斗结果到全服

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 公告ID |
| type | VARCHAR(32) | NOT NULL | 公告类型 |
| content | TEXT | NOT NULL | 公告内容 |
| zone_id | VARCHAR(32) | | 相关区域 |
| winner_user_id | INTEGER | | 胜利者ID |
| loser_user_id | INTEGER | | 失败者ID |
| pvp_encounter_id | INTEGER | | 关联PVP战斗ID |
| importance | INTEGER | DEFAULT 1 | 重要程度 (1-5) |
| expires_at | DATETIME | | 过期时间 |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

```sql
CREATE TABLE server_announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(32) NOT NULL,
    content TEXT NOT NULL,
    zone_id VARCHAR(32),
    winner_user_id INTEGER,
    loser_user_id INTEGER,
    pvp_encounter_id INTEGER,
    importance INTEGER DEFAULT 1,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (winner_user_id) REFERENCES users(id),
    FOREIGN KEY (loser_user_id) REFERENCES users(id),
    FOREIGN KEY (pvp_encounter_id) REFERENCES pvp_encounters(id)
);

CREATE INDEX idx_announcements_time ON server_announcements(created_at DESC);
CREATE INDEX idx_announcements_type ON server_announcements(type);
```

**公告类型:**

| 类型 | 说明 | 示例 |
|-----|------|-----|
| `pvp_victory` | PVP胜利 | [西部荒野] 联盟「勇士小明」击败了部落「暗影猎手」！ |
| `zone_captured` | 区域易主 | ⚔️ 联盟已占领「石爪山脉」！该区域联盟效率+15% |
| `kill_streak` | 连杀公告 | 🔥 联盟「死亡骑士」达成5连杀！ |
| `zone_contested` | 区域争夺激烈 | ⚠️ 「灰谷」正在激烈争夺中！双方势均力敌 |
| `faction_dominant` | 阵营优势 | 👑 部落本周在全服7个区域保持控制！ |

---

### 八、荣誉商店表

#### honor_shop - 荣誉商店配置表

> 📌 **荣誉兑换**: 使用荣誉值兑换独特奖励

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 商品ID |
| name | VARCHAR(64) | NOT NULL | 商品名称 |
| description | TEXT | | 描述 |
| item_type | VARCHAR(32) | NOT NULL | 物品类型: equipment/consumable/cosmetic/title |
| item_id | VARCHAR(32) | | 关联物品ID |
| honor_cost | INTEGER | NOT NULL | 荣誉花费 |
| rank_required | INTEGER | DEFAULT 0 | 需求军衔等级 |
| faction | VARCHAR(16) | | 阵营限制 (NULL=双阵营) |
| weekly_limit | INTEGER | | 每周购买限制 |
| stock | INTEGER | | 库存 (NULL=无限) |
| is_active | INTEGER | DEFAULT 1 | 是否在售 |

```sql
CREATE TABLE honor_shop (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    item_type VARCHAR(32) NOT NULL,
    item_id VARCHAR(32),
    honor_cost INTEGER NOT NULL,
    rank_required INTEGER DEFAULT 0,
    faction VARCHAR(16),
    weekly_limit INTEGER,
    stock INTEGER,
    is_active INTEGER DEFAULT 1,
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX idx_honor_shop_rank ON honor_shop(rank_required);
```

**荣誉商店示例商品:**

| 商品 | 类型 | 荣誉花费 | 需求军衔 |
|-----|------|---------|---------|
| PVP套装·头盔 | 装备 | 2,000 | 5 (准尉) |
| PVP套装·胸甲 | 装备 | 3,500 | 7 (中尉) |
| 战斗徽章 | 饰品 | 1,500 | 3 (中士) |
| 荣誉药剂 | 消耗品 | 100 | 0 |
| 征服者称号 | 称号 | 10,000 | 10 (中校) |
| 战马坐骑皮肤 | 外观 | 50,000 | 12 (准将) |

---

### PVP遭遇战流程图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PVP遭遇战触发流程                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 匹配检测                                                                 │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │ 每X秒检测各地图的联盟/部落在线玩家                                 │    │
│     │ IF 双方都有玩家 AND 随机触发(概率Y%)                              │    │
│     │ THEN 随机匹配一对敌对玩家进行PVP                                  │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                       ▼                                     │
│  2. 战斗执行                                                                 │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │ 双方使用各自的战斗策略自动对战                                     │    │
│     │ 战斗引擎模拟回合制PVP                                             │    │
│     │ 记录战斗日志和伤害统计                                            │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                       ▼                                     │
│  3. 结算奖励                                                                 │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │ 胜者: +荣誉值 (基于对手等级/军衔)                                  │    │
│     │ 败者: 无惩罚 (避免挫败感)                                         │    │
│     │ 更新双方PVP统计                                                   │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                       ▼                                     │
│  4. 更新地图控制                                                             │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │ 更新 zone_faction_control 的胜败统计                              │    │
│     │ 重新计算控制积分和效率加成                                         │    │
│     │ IF 控制权变化 THEN 触发区域易主公告                               │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                       ▼                                     │
│  5. 发布公告                                                                 │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │ 生成战斗结果公告                                                   │    │
│     │ 推送到全服在线玩家                                                 │    │
│     │ 记录到 server_announcements 表供查询                              │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 效率加成计算示例

```
场景: 艾尔文森林 - 联盟控制 (积分+45)

联盟玩家挂机效率:
  基础经验: 100/分钟
  控制加成: +9% (积分45 ÷ 100 × 20%)
  实际经验: 100 × 1.09 = 109/分钟

部落玩家挂机效率:
  基础经验: 100/分钟
  被控惩罚: -4.5% (积分45 ÷ 100 × 10%)
  实际经验: 100 × 0.955 = 95.5/分钟

差距: 联盟比部落多获得 14% 经验
```

> 💡 这个机制激励玩家提升战力和优化策略来争夺地图控制权，同时保持竞争的持续性

---

## 🐌 游戏节奏与体力系统

> 📌 **设计理念**: 慢节奏沉浸式体验，让玩家有时间探索和享受游戏深度

### 节奏设计原则

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          慢节奏游戏设计                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ❌ 快餐式体验                          ✅ 沉浸式体验                        │
│   ├─ 3天满级，1周毕业                    ├─ 数周探索，数月精通                 │
│   ├─ 秒杀怪物，无脑挂机                  ├─ 每场战斗都有意义                   │
│   ├─ 数值膨胀，装备快速淘汰              ├─ 装备长期有价值，培养值得            │
│   └─ 内容消耗殆尽                        └─ 持续发现新东西                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 节奏控制点

| 方面 | 设计 | 效果 |
|-----|------|------|
| **战斗时长** | 每场战斗平均8-15回合 | 有足够时间展现策略 |
| **升级曲线** | 1-30级约1周，31-60级约3周 | 中后期有深度探索空间 |
| **体力系统** | 每场战斗消耗体力 | 防止无限刷，鼓励策略优化 |
| **探索节奏** | 新区域需要解锁条件 | 有目标感和成就感 |
| **装备成长** | 好装备值得培养数周 | 减少装备焦虑 |

### 体力系统

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           体力系统                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   体力上限: 100 点                                                           │
│   恢复速度: 1 点 / 10 分钟 (每天144点)                                        │
│                                                                             │
│   消耗:                                                                      │
│   ├─ 普通战斗: 1 点                                                          │
│   ├─ 精英战斗: 2 点                                                          │
│   ├─ Boss战斗: 5 点                                                          │
│   └─ 深渊挑战: 3 点/层                                                       │
│                                                                             │
│   玩家可以:                                                                  │
│   ├─ 每天打约100-150场普通战斗                                               │
│   ├─ 或者打20-30场Boss                                                       │
│   └─ 需要规划每天的挂机目标                                                   │
│                                                                             │
│   离线收益: 体力溢出时转化为离线经验/金币 (效率50%)                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### user_stamina - 玩家体力表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| user_id | INTEGER | PRIMARY KEY FK | 用户ID |
| current_stamina | INTEGER | DEFAULT 100 | 当前体力 |
| max_stamina | INTEGER | DEFAULT 100 | 最大体力 |
| last_regen_at | DATETIME | | 上次恢复时间 |
| overflow_exp | INTEGER | DEFAULT 0 | 溢出转化的经验 |
| overflow_gold | INTEGER | DEFAULT 0 | 溢出转化的金币 |

```sql
CREATE TABLE IF NOT EXISTS user_stamina (
    user_id INTEGER PRIMARY KEY,
    current_stamina INTEGER DEFAULT 100,
    max_stamina INTEGER DEFAULT 100,
    last_regen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    overflow_exp INTEGER DEFAULT 0,
    overflow_gold INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

## 🎯 作战策略系统

> 📎 **详细设计已迁移至独立文档**: [作战策略系统设计](strategy_design.md)
>
> 策略系统包含：技能优先级、条件规则、目标选择、资源管理等功能。
> 数据表结构请参阅上方 [13. battle_strategies](#13-battle_strategies---战斗策略表)。

---

## 📊 战斗数据分析系统

> 📌 **核心理念**: 让玩家通过数据分析验证和优化作战策略

### 统计周期

| 周期 | 数据内容 | 用途 |
|-----|---------|------|
| **上一场** | 完整回合记录、每个技能使用、详细伤害流水 | 即时验证策略调整效果 |
| **今日** | 汇总统计、技能效率、胜率 | 当天表现评估 |
| **本周** | 趋势对比、进步曲线 | 中期策略优化 |
| **本月** | 长期数据、里程碑 | 整体成长回顾 |
| **全部** | 历史总计 | 终极成就统计 |

### 多角色战斗统计示例

```
╔═══════════════════════════════════════════════════════════════════════════╗
║  📋 上一场战斗详情                                      2024-11-29 14:32:15 ║
╠═══════════════════════════════════════════════════════════════════════════╣
║                                                                           ║
║  地点: 暮色森林            敌人: 亡灵小队 (3体)          结果: ✓ 胜利      ║
║  回合数: 12                耗时: 28秒                   获得: 85经验 24金币║
║                                                                           ║
║  ═══════════════════════════════════════════════════════════════════════  ║
║  👥 我方队伍 (3/5)                                                         ║
║  ───────────────────────────────────────────────────────────────────────  ║
║  │ 角色         │ 职业   │ 伤害输出 │ 承受伤害 │ 治疗量 │ DPR  │ 状态     │ ║
║  │ 破晓         │ 战士   │ 89  45%  │ 52  62%  │ -      │ 7.4  │ HP 43/95│ ║
║  │ 寒霜         │ 法师   │ 76  38%  │ 18  21%  │ -      │ 6.3  │ HP 52/60│ ║
║  │ 圣光         │ 牧师   │ 32  17%  │ 14  17%  │ 45     │ 2.7  │ HP 55/70│ ║
║  ├──────────────┴────────┴──────────┴──────────┴────────┴──────┴─────────┤ ║
║  │ 队伍合计              │ 197      │ 84       │ 45     │ 16.4 │          │ ║
║  ───────────────────────────────────────────────────────────────────────  ║
║                                                                           ║
║  ═══════════════════════════════════════════════════════════════════════  ║
║  💀 敌方队伍 (3体)                                                         ║
║  ───────────────────────────────────────────────────────────────────────  ║
║  │ 敌人         │ 等级   │ 伤害输出 │ 承受伤害 │ 控制时间 │ 击杀者 │ 状态  │ ║
║  │ 骷髅战士     │ Lv.10  │ 35  42%  │ 78       │ 2回合   │ 破晓   │ R8死亡│ ║
║  │ 骷髅法师     │ Lv.10  │ 32  38%  │ 65       │ 1回合   │ 寒霜   │ R10死亡║
║  │ 骷髅弓手     │ Lv.9   │ 17  20%  │ 54       │ 0回合   │ 破晓   │ R12死亡║
║  ├──────────────┴────────┴──────────┴──────────┴─────────┴────────┴───────┤ ║
║  │ 敌方合计              │ 84       │ 197      │ 3回合   │        │        │ ║
║  ───────────────────────────────────────────────────────────────────────  ║
║                                                                           ║
║  ═══════════════════════════════════════════════════════════════════════  ║
║  📈 队伍配合分析                                                           ║
║  ───────────────────────────────────────────────────────────────────────  ║
║  │ 指标               │ 数值    │ 评价                                    │ ║
║  │ 伤害分布           │ 45/38/17│ 战士主力，法师辅助，牧师补刀 ✓           │ ║
║  │ 承伤分布           │ 62/21/17│ 战士扛伤合理 ✓                          │ ║
║  │ 控制覆盖           │ 25%     │ 3/12回合有敌人被控 ✓                    │ ║
║  │ 治疗效率           │ 89%     │ 过量治疗11% (良好)                       │ ║
║  │ 集火效率           │ 78%     │ 优先击杀高威胁目标 ✓                     │ ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
```

### 关键分析指标

| 指标 | 计算方式 | 用途 |
|-----|---------|------|
| **DPR** | 总伤害 / 回合数 | 评估每回合输出效率 |
| **技能效率** | 技能伤害 / 资源消耗 | 判断技能性价比 |
| **生存率** | 存活战斗 / 总战斗 | 策略稳定性 |
| **资源利用率** | 实际消耗 / 可用资源 | 资源管理评估 |
| **控制贡献** | 控制回合数 / 总回合数 | 控制策略效果 |
| **治疗效率** | 有效治疗 / 总治疗 | 过量治疗检测 |
| **集火效率** | 目标切换次数 | 是否集中火力 |
| **承伤分布** | 各角色承伤占比 | 坦克是否有效 |

#### detailed_battle_logs - 详细战斗日志表

> 📌 记录每一场战斗的完整回合数据，用于"上一场"分析

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 日志ID |
| battle_id | VARCHAR(64) | NOT NULL | 战斗ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| zone_id | VARCHAR(32) | | 区域ID |
| battle_type | VARCHAR(16) | NOT NULL | 战斗类型: pve/pvp/abyss |
| result | VARCHAR(16) | NOT NULL | 结果: victory/defeat/draw |
| total_turns | INTEGER | NOT NULL | 总回合数 |
| duration_seconds | INTEGER | | 战斗时长 |
| player_team_data | TEXT | NOT NULL | 我方队伍数据 (JSON) |
| enemy_team_data | TEXT | NOT NULL | 敌方队伍数据 (JSON) |
| turn_logs | TEXT | NOT NULL | 回合日志 (JSON) |
| exp_gained | INTEGER | DEFAULT 0 | 获得经验 |
| gold_gained | INTEGER | DEFAULT 0 | 获得金币 |
| items_dropped | TEXT | | 掉落物品 (JSON) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 战斗时间 |

```sql
CREATE TABLE IF NOT EXISTS detailed_battle_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id VARCHAR(64) NOT NULL,
    user_id INTEGER NOT NULL,
    zone_id VARCHAR(32),
    battle_type VARCHAR(16) NOT NULL,
    result VARCHAR(16) NOT NULL,
    total_turns INTEGER NOT NULL,
    duration_seconds INTEGER,
    player_team_data TEXT NOT NULL,
    enemy_team_data TEXT NOT NULL,
    turn_logs TEXT NOT NULL,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    items_dropped TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

CREATE INDEX idx_detailed_logs_user ON detailed_battle_logs(user_id);
CREATE INDEX idx_detailed_logs_time ON detailed_battle_logs(created_at DESC);
CREATE INDEX idx_detailed_logs_battle ON detailed_battle_logs(battle_id);

-- 只保留最近100场详细日志，定期清理
```

**player_team_data / enemy_team_data JSON格式:**
```json
{
  "units": [
    {
      "id": "char_1",
      "name": "破晓",
      "class": "warrior",
      "level": 25,
      "initial_hp": 95,
      "final_hp": 43,
      "max_hp": 95,
      "damage_dealt": 89,
      "damage_taken": 52,
      "healing_done": 0,
      "healing_taken": 32,
      "is_dead": false,
      "skills_used": [
        {"skill_id": "charge", "count": 2, "damage": 24, "crits": 0, "effects": ["stun"]},
        {"skill_id": "rend", "count": 3, "damage": 33, "crits": 0, "effects": ["bleed"]},
        {"skill_id": "heroic_strike", "count": 5, "damage": 58, "crits": 1, "effects": []}
      ]
    }
  ],
  "total_damage": 197,
  "total_healing": 45
}
```

**turn_logs JSON格式:**
```json
[
  {
    "turn": 1,
    "actions": [
      {"source": "破晓", "skill": "冲锋", "target": "骷髅法师", "damage": 12, "crit": false, "effect": "眩晕"},
      {"source": "寒霜", "skill": "寒冰箭", "target": "骷髅战士", "damage": 15, "crit": false, "effect": "减速"},
      {"source": "圣光", "skill": "真言术:盾", "target": "破晓", "shield": 12},
      {"source": "骷髅战士", "skill": "攻击", "target": "破晓", "damage": 8, "absorbed": 8},
      {"source": "骷髅法师", "status": "眩晕", "skipped": true},
      {"source": "骷髅弓手", "skill": "射击", "target": "寒霜", "damage": 6}
    ]
  }
]
```

#### character_battle_stats - 角色战斗统计表

> 📌 汇总统计，用于今日/本周/本月/全部分析

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| stat_date | DATE | NOT NULL | 统计日期 |
| battles_total | INTEGER | DEFAULT 0 | 总战斗场次 |
| battles_won | INTEGER | DEFAULT 0 | 胜利场次 |
| battles_lost | INTEGER | DEFAULT 0 | 失败场次 |
| total_damage_dealt | INTEGER | DEFAULT 0 | 总造成伤害 |
| total_damage_taken | INTEGER | DEFAULT 0 | 总承受伤害 |
| total_healing_done | INTEGER | DEFAULT 0 | 总治疗量 |
| total_turns | INTEGER | DEFAULT 0 | 总回合数 |
| deaths | INTEGER | DEFAULT 0 | 死亡次数 |
| kills | INTEGER | DEFAULT 0 | 击杀数 |
| crits | INTEGER | DEFAULT 0 | 暴击次数 |
| dodges | INTEGER | DEFAULT 0 | 闪避次数 |
| skills_used | TEXT | | 技能使用统计 (JSON) |

```sql
CREATE TABLE IF NOT EXISTS character_battle_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    stat_date DATE NOT NULL,
    battles_total INTEGER DEFAULT 0,
    battles_won INTEGER DEFAULT 0,
    battles_lost INTEGER DEFAULT 0,
    total_damage_dealt INTEGER DEFAULT 0,
    total_damage_taken INTEGER DEFAULT 0,
    total_healing_done INTEGER DEFAULT 0,
    total_turns INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    kills INTEGER DEFAULT 0,
    crits INTEGER DEFAULT 0,
    dodges INTEGER DEFAULT 0,
    skills_used TEXT,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    UNIQUE(character_id, stat_date)
);

CREATE INDEX idx_char_stats_date ON character_battle_stats(character_id, stat_date DESC);
```

### 策略建议生成

> 📌 基于数据分析自动生成优化建议

| 检测项 | 条件 | 建议 |
|-------|------|------|
| 普攻过多 | 普攻占比 > 30% | 降低资源阈值，增加技能使用 |
| 治疗不足 | 死亡率 > 20% | 调整治疗触发条件或增加治疗优先级 |
| 过量治疗 | 治疗效率 < 70% | 提高治疗触发的HP阈值 |
| 控制浪费 | 控制重叠 | 调整控制技能间隔 |
| 爆发时机 | 大招未命中低血敌人 | 添加斩杀类条件规则 |
| 资源浪费 | 资源溢出 > 20% | 降低资源阈值 |

---

## 🛡️ 仇恨值系统

> 📌 **核心机制**: 敌人攻击仇恨值最高的目标，形成坦克-治疗-输出的"铁三角"战术体系

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          仇恨值系统概览                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   敌人始终攻击仇恨值最高的目标                                                │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────┐      │
│   │                     仇恨列表 (某只怪物)                          │      │
│   │  ─────────────────────────────────────────────────────────────  │      │
│   │  1. 破晓 (战士)  ████████████████████████████░░  2850 仇恨 ← 攻击│      │
│   │  2. 寒霜 (法师)  ██████████████████░░░░░░░░░░░░  1820 仇恨      │      │
│   │  3. 暗影 (盗贼)  █████████████░░░░░░░░░░░░░░░░░  1340 仇恨      │      │
│   │  4. 圣光 (牧师)  ████░░░░░░░░░░░░░░░░░░░░░░░░░░   520 仇恨      │      │
│   └─────────────────────────────────────────────────────────────────┘      │
│                                                                             │
│   铁三角职责:                                                                │
│   ├─ 🛡️ 坦克: 高仇恨技能，保持仇恨第一，承受伤害                              │
│   ├─ 💚 治疗: 低仇恨治疗，保持队伍存活，注意仇恨控制                           │
│   └─ ⚔️ 输出: 高伤害但需控制仇恨，避免OT(仇恨超过坦克)                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 仇恨值生成规则

| 行为 | 基础仇恨 | 说明 |
|-----|---------|------|
| **造成伤害** | 伤害值 × 1.0 | 每1点伤害产生1点仇恨 |
| **治疗队友** | 治疗量 × 0.5 | 治疗产生一半仇恨 (分摊给所有敌人) |
| **施加Buff** | 固定值 | 部分增益技能产生少量仇恨 |
| **嘲讽** | 特殊 | 立即成为第一仇恨 + 10%额外仇恨 |
| **进入战斗** | 100 | 初始仇恨值 |

### 职业仇恨系数

| 职业 | 仇恨系数 | 定位 | 说明 |
|-----|---------|------|------|
| **战士** | ×1.3 | 坦克 | 高仇恨生成，主力坦克 |
| **圣骑士** | ×1.2 | 坦克/治疗 | 半坦半奶，灵活定位 |
| **德鲁伊** | ×1.2 (熊) / ×0.8 (其他) | 坦克/治疗/输出 | 变形影响仇恨 |
| **盗贼** | ×0.7 | 输出 | 低仇恨，有仇恨清除 |
| **猎人** | ×0.8 | 输出 | 宠物分担仇恨 |
| **法师** | ×0.8 | 输出 | 高爆发需注意OT |
| **术士** | ×0.8 | 输出 | 宠物分担仇恨 |
| **萨满** | ×0.9 | 输出/治疗 | 均衡职业 |
| **牧师** | ×0.6 | 治疗 | 最低仇恨系数 |

### OT (Over Threat) 判定

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          OT判定规则                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   近战夺取仇恨: 需要超过当前目标仇恨的 110%                                    │
│   远程夺取仇恨: 需要超过当前目标仇恨的 130%                                    │
│                                                                             │
│   示例:                                                                      │
│   战士当前仇恨: 1000                                                         │
│   法师(远程)要OT: 需要 1000 × 130% = 1300 仇恨                               │
│   盗贼(近战)要OT: 需要 1000 × 110% = 1100 仇恨                               │
│                                                                             │
│   这个缓冲机制让坦克有反应时间使用嘲讽                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 仇恨相关技能

#### 坦克技能 (高仇恨)

| 技能 | 职业 | 仇恨效果 | CD | 说明 |
|-----|------|---------|-----|------|
| **嘲讽** | 战士 | 强制第一 +10% | 4回合 | 单体嘲讽 |
| **挑战怒吼** | 战士 | AOE嘲讽 | 8回合 | 群体嘲讽 |
| **盾牌猛击** | 战士 | 伤害×2仇恨 | 0 | 高仇恨攻击 |
| **破甲攻击** | 战士 | 伤害×1.5仇恨 | 0 | 削弱+仇恨 |
| **复仇** | 战士 | 格挡后+仇恨 | 被动 | 被动仇恨 |
| **正义之盾** | 圣骑 | 伤害×1.5仇恨 | 2回合 | 神圣+仇恨 |
| **奉献** | 圣骑 | AOE持续仇恨 | 3回合 | 范围仇恨 |
| **正义防御** | 圣骑 | 嘲讽 | 4回合 | 单体嘲讽 |

#### 仇恨管理技能 (输出/治疗用)

| 技能 | 职业 | 仇恨效果 | CD | 说明 |
|-----|------|---------|-----|------|
| **佯攻** | 盗贼 | 清除50%仇恨 | 4回合 | 紧急脱仇 |
| **消失** | 盗贼 | 清除100%仇恨 | 10回合 | 完全脱仇 |
| **隐形术** | 法师 | 清除100%仇恨 | 10回合 | 完全脱仇 |
| **寒冰屏障** | 法师 | 暂停仇恨生成 | 8回合 | 冰块期间不产生仇恨 |
| **渐隐术** | 牧师 | 仇恨生成-50% | 6回合 | 持续3回合 |
| **假死** | 猎人 | 清除100%仇恨 | 8回合 | 紧急脱仇 |
| **误导** | 猎人 | 仇恨转移给坦克 | 6回合 | 仇恨辅助 |
| **灵魂碎裂** | 术士 | 宠物嘲讽 | 4回合 | 转移仇恨 |

### 仇恨相关装备词缀

**新增前缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| 守护者的 | 仇恨生成+20%, 嘲讽CD-1 | 坦克武器 | 稀有 |
| 威压的 | 仇恨生成+15% | 坦克装备 | 精良 |
| 隐秘的 | 仇恨生成-20% | 输出武器 | 精良 |
| 暗影的 | 暴击仇恨-30% | 输出装备 | 稀有 |

**新增后缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| of 威胁 | 仇恨生成+15% | 坦克装备 | 精良 |
| of 守护 | 仇恨生成+10%, 格挡+5% | 盾牌 | 精良 |
| of 隐匿 | 仇恨生成-15% | 输出装备 | 精良 |
| of 消散 | 仇恨衰减+20% | 治疗装备 | 稀有 |

### 敌人AI目标选择

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       敌人AI目标选择逻辑                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   每回合开始:                                                                │
│   1. 检查仇恨列表                                                            │
│   2. 如果当前目标仍是最高仇恨 → 继续攻击                                       │
│   3. 如果有新目标仇恨超过阈值 → 切换目标 (OT发生)                              │
│   4. 特殊情况:                                                               │
│      - 被嘲讽 → 强制攻击嘲讽者 (持续到嘲讽结束)                               │
│      - 当前目标死亡 → 攻击次高仇恨                                            │
│      - 仇恨清零 → 重新选择目标                                               │
│                                                                             │
│   智能敌人 (精英/Boss):                                                       │
│   - 可能无视嘲讽攻击治疗 (需要特殊处理)                                        │
│   - 可能有仇恨重置技能                                                        │
│   - 阶段转换时可能重置仇恨                                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 策略系统扩展

#### 新增条件规则

| 条件类型 | 可选值 | 示例 |
|---------|-------|------|
| **自身仇恨排名** | =, >, < | IF 仇恨排名 > 1 THEN 使用 [佯攻] |
| **仇恨百分比** | 相对坦克 | IF 仇恨 > 坦克的90% THEN 停止攻击 |
| **队友仇恨** | 检测OT风险 | IF 治疗仇恨排名=1 THEN 使用 [嘲讽] |
| **被攻击状态** | 是/否 | IF 被攻击 AND 不是坦克 THEN 使用 [消失] |

#### 角色定位设置

| 定位 | AI行为 |
|-----|-------|
| **坦克** | 优先高仇恨技能，监控队友仇恨，及时嘲讽 |
| **输出** | 正常输出，仇恨接近坦克时减缓/使用仇恨清除 |
| **治疗** | 专注治疗，使用仇恨降低技能 |

### 数据库表更新

#### skills表新增字段

| 字段 | 类型 | 说明 |
|-----|------|------|
| threat_modifier | REAL | 仇恨系数 (默认1.0) |
| threat_type | VARCHAR(16) | 仇恨类型: normal/high/taunt/reduce/clear |

#### classes表新增字段

| 字段 | 类型 | 说明 |
|-----|------|------|
| base_threat_modifier | REAL | 基础仇恨系数 |
| combat_role | VARCHAR(16) | 战斗定位: tank/healer/dps/hybrid |
| is_ranged | INTEGER | 是否远程 (影响OT阈值) |

#### 新增 battle_threat_log 表

> 📌 记录战斗中的仇恨变化，用于分析

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY | ID |
| battle_id | VARCHAR(64) | NOT NULL | 战斗ID |
| turn | INTEGER | NOT NULL | 回合数 |
| enemy_id | VARCHAR(32) | NOT NULL | 敌人ID |
| threat_snapshot | TEXT | NOT NULL | 仇恨快照 (JSON) |
| target_changed | INTEGER | DEFAULT 0 | 是否切换目标 |
| ot_occurred | INTEGER | DEFAULT 0 | 是否发生OT |
| taunt_used | INTEGER | DEFAULT 0 | 是否使用嘲讽 |

```sql
CREATE TABLE IF NOT EXISTS battle_threat_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id VARCHAR(64) NOT NULL,
    turn INTEGER NOT NULL,
    enemy_id VARCHAR(32) NOT NULL,
    threat_snapshot TEXT NOT NULL,
    target_changed INTEGER DEFAULT 0,
    ot_occurred INTEGER DEFAULT 0,
    taunt_used INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_threat_log_battle ON battle_threat_log(battle_id);
```

**threat_snapshot JSON格式:**
```json
{
  "current_target": "char_1",
  "threat_list": [
    {"unit_id": "char_1", "name": "破晓", "threat": 2850, "percent": 43.7},
    {"unit_id": "char_2", "name": "寒霜", "threat": 1820, "percent": 27.9},
    {"unit_id": "char_3", "name": "暗影", "threat": 1340, "percent": 20.5},
    {"unit_id": "char_4", "name": "圣光", "threat": 520, "percent": 8.0}
  ]
}
```

### 仇恨分析统计

| 指标 | 计算方式 | 用途 |
|-----|---------|------|
| **仇恨占比** | 角色仇恨 / 总仇恨 | 评估仇恨分布 |
| **OT次数** | 仇恨超过坦克次数 | 输出控制评估 |
| **坦克保持率** | 坦克为第一仇恨的回合占比 | 坦克表现 |
| **嘲讽效率** | 成功嘲讽 / 嘲讽使用 | 嘲讽使用评估 |
| **被攻击分布** | 各角色被攻击回合占比 | 伤害分担评估 |
| **仇恨清除使用** | 仇恨清除技能使用次数 | 输出策略评估 |

### 铁三角平衡设计

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          铁三角平衡                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                           🛡️ 坦克                                           │
│                          /        \                                         │
│                         /          \                                        │
│                    仇恨管理      承受伤害                                     │
│                       /              \                                      │
│                      /                \                                     │
│                 ⚔️ 输出 ───治疗依赖─── 💚 治疗                               │
│                                                                             │
│   平衡点:                                                                    │
│   ├─ 坦克仇恨不足 → 治疗/输出被打 → 团灭                                      │
│   ├─ 输出仇恨过高 → OT → 被打死                                              │
│   ├─ 治疗仇恨过高 → 被攻击 → 无法治疗 → 团灭                                  │
│   └─ 平衡良好 → 坦克抗伤，输出输出，治疗治疗                                   │
│                                                                             │
│   策略深度:                                                                  │
│   ├─ 坦克: 开场建立仇恨，监控OT，及时嘲讽                                     │
│   ├─ 输出: 控制爆发时机，必要时使用仇恨清除                                    │
│   └─ 治疗: 使用仇恨降低技能，避免过量治疗                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 💬 阵营聊天系统

> 📌 **核心设计**: 联盟和部落各自独立的聊天频道，跨阵营无法直接交流

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          聊天系统概览                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   联盟频道                              部落频道                             │
│   ┌─────────────────────┐              ┌─────────────────────┐             │
│   │ [世界] 人类勇士说:   │              │ [世界] 兽人战士说:   │             │
│   │   为了联盟！         │              │   为了部落！         │             │
│   │                     │              │                     │             │
│   │ [区域] 矮人猎人:     │              │ [区域] 牛头人德鲁伊: │             │
│   │   有人组队吗？       │              │   这里有精英怪       │             │
│   │                     │              │                     │             │
│   │ [私聊] 来自暗夜精灵  │              │ [私聊] 来自巨魔       │             │
│   └─────────────────────┘              └─────────────────────┘             │
│                    ╲                    ╱                                   │
│                     ╲                  ╱                                    │
│                      ╲   ❌ 不互通   ╱                                     │
│                       ╲            ╱                                        │
│                        ╲          ╱                                         │
│                         遭遇战时可见                                         │
│                         (乱码显示)                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 频道类型

| 频道 | 代码 | 范围 | 颜色 | 说明 |
|-----|------|------|-----|------|
| **世界频道** | world | 全阵营 | 🟡 金色 | 同阵营所有玩家可见 |
| **区域频道** | zone | 同区域 | 🟠 橙色 | 当前区域内同阵营玩家 |
| **交易频道** | trade | 全阵营 | 🟢 绿色 | 买卖装备/材料 |
| **组队频道** | lfg | 全阵营 | 🔵 蓝色 | 找队友 |
| **私聊** | whisper | 1对1 | 🟣 紫色 | 只能和同阵营玩家私聊 |
| **系统消息** | system | 全服 | ⚪ 白色 | 系统公告、上下线通知 |
| **战场频道** | battlefield | PvP遭遇 | 🔴 红色 | PvP时敌方消息(乱码化) |

### 跨阵营遭遇战特殊处理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PvP遭遇时的聊天显示                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   联盟玩家视角:                                                              │
│   ┌─────────────────────────────────────────────────────────┐              │
│   │ [战斗] 兽人战士(部落) 对你使用了 [英勇打击]               │              │
│   │ [战场] <部落玩家> 说: "Kek zug zug!"  ← 乱码化            │              │
│   │ [战场] 你说: "为了联盟！"                                 │              │
│   └─────────────────────────────────────────────────────────┘              │
│                                                                             │
│   部落玩家视角:                                                              │
│   ┌─────────────────────────────────────────────────────────┐              │
│   │ [战斗] 你对 人类战士(联盟) 使用了 [英勇打击]              │              │
│   │ [战场] 你说: "Lok'tar ogar!"                             │              │
│   │ [战场] <联盟玩家> 说: "Bur gull..."  ← 乱码化             │              │
│   └─────────────────────────────────────────────────────────┘              │
│                                                                             │
│   乱码规则: 将消息转换为对方阵营的"语言"                                      │
│   联盟→部落: 使用兽人语音节替换                                              │
│   部落→联盟: 使用通用语音节替换                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 数据库表设计

#### chat_messages - 聊天消息表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY | 消息ID |
| channel | VARCHAR(16) | NOT NULL | 频道类型 |
| faction | VARCHAR(16) | | 阵营 (world/zone/trade/lfg用) |
| zone_id | VARCHAR(32) | | 区域ID (zone频道用) |
| sender_id | INTEGER | NOT NULL FK | 发送者用户ID |
| sender_name | VARCHAR(32) | NOT NULL | 发送者角色名 |
| sender_class | VARCHAR(32) | | 发送者职业 |
| receiver_id | INTEGER | FK | 接收者ID (私聊用) |
| content | TEXT | NOT NULL | 消息内容 |
| created_at | DATETIME | DEFAULT NOW | 发送时间 |

```sql
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel VARCHAR(16) NOT NULL,
    faction VARCHAR(16),
    zone_id VARCHAR(32),
    sender_id INTEGER NOT NULL,
    sender_name VARCHAR(32) NOT NULL,
    sender_class VARCHAR(32),
    receiver_id INTEGER,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_chat_channel ON chat_messages(channel, faction);
CREATE INDEX IF NOT EXISTS idx_chat_zone ON chat_messages(zone_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_whisper ON chat_messages(sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_chat_time ON chat_messages(created_at DESC);
```

#### chat_blocks - 屏蔽列表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY | ID |
| user_id | INTEGER | NOT NULL FK | 用户ID |
| blocked_id | INTEGER | NOT NULL FK | 被屏蔽用户ID |
| created_at | DATETIME | DEFAULT NOW | 屏蔽时间 |

```sql
CREATE TABLE IF NOT EXISTS chat_blocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    blocked_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, blocked_id)
);
```

#### user_online_status - 在线状态表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| user_id | INTEGER | PRIMARY KEY | 用户ID |
| character_name | VARCHAR(32) | | 当前角色名 |
| faction | VARCHAR(16) | | 阵营 |
| zone_id | VARCHAR(32) | | 当前区域 |
| last_active | DATETIME | | 最后活跃时间 |
| is_online | INTEGER | DEFAULT 0 | 是否在线 |

```sql
CREATE TABLE IF NOT EXISTS user_online_status (
    user_id INTEGER PRIMARY KEY,
    character_name VARCHAR(32),
    faction VARCHAR(16),
    zone_id VARCHAR(32),
    last_active DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_online INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_online_faction ON user_online_status(faction, is_online);
CREATE INDEX IF NOT EXISTS idx_online_zone ON user_online_status(zone_id, is_online);
```

### 命令系统

| 命令 | 说明 | 示例 |
|-----|------|------|
| `/s 消息` | 发送到世界频道 | `/s 有人组队吗？` |
| `/z 消息` | 发送到区域频道 | `/z 这里有精英怪` |
| `/t 消息` | 发送到交易频道 | `/t 收购铜矿石` |
| `/lfg 消息` | 发送到组队频道 | `/lfg 找T和奶` |
| `/w 玩家 消息` | 私聊 | `/w 破晓 你好` |
| `/r 消息` | 回复上一个私聊 | `/r 好的` |
| `/block 玩家` | 屏蔽玩家 | `/block 讨厌鬼` |
| `/unblock 玩家` | 取消屏蔽 | `/unblock 讨厌鬼` |

### 内容安全

| 规则 | 设置 | 说明 |
|-----|------|------|
| 敏感词过滤 | 替换为 *** | 自动过滤不当内容 |
| 刷屏检测 | 间隔 > 3秒 | 同一消息重复发送限制 |
| 消息长度 | 最多200字符 | 防止超长消息 |
| 频率限制 | 每分钟20条 | 防止刷屏 |

