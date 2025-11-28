# Text WoW 数据库设计文档

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

#### 属性加成机制

**采用"基础值 + 百分比"混合方案：**

| 组成部分 | 说明 | 示例(兽人力量) |
|---------|------|--------------|
| 基础加成 | 创建角色时一次性加成 | +5点力量 |
| 百分比加成 | 该属性总值的额外加成 | +5%力量 |

**计算公式：**
```
最终属性 = (基础属性 + 职业成长 × 等级 + 装备加成) × (1 + 种族百分比加成)
```

**示例计算（兽人战士 Lv.30）：**
```
基础力量 = 15 (职业基础)
等级成长 = 2 × 30 = 60
种族基础 = +5
装备加成 = +20 (假设)
小计 = 15 + 60 + 5 + 20 = 100

种族百分比 = +5%
最终力量 = 100 × 1.05 = 105

vs 人类战士: 100 × 1.00 = 100 (差5点，约5%)
```

这样设计的好处：
- ✅ 初期差异明显（基础值）
- ✅ 后期仍有意义（百分比）
- ✅ 差距不会过于悬殊（只有5%左右）

---

**联盟种族:**

| 种族 | 属性加成 | 被动特性1 | 被动特性2 |
|-----|---------|---------|---------|
| **人类** | 精神+3, 精神+3% | 💡 适应力：经验获取+10% | 🗡️ 剑术专精：物理伤害+3% |
| **矮人** | 力量+3, 耐力+5% | ❄️ 霜抗：冰霜伤害-15% | 🛡️ 石肤：受到暴击伤害-10% |
| **暗夜精灵** | 敏捷+5, 敏捷+3% | 🌙 暗影之心：伤害+5% | 👁️ 敏锐：闪避率+2% |
| **侏儒** | 智力+5, 智力+5% | ⚡ 灵巧心智：法术暴击+3% | 🔧 工程专精：对机械怪伤害+15% |

**部落种族:**

| 种族 | 属性加成 | 被动特性1 | 被动特性2 |
|-----|---------|---------|---------|
| **兽人** | 力量+5, 力量+5% | 💢 嗜血：HP<30%时攻击+15% | 💪 坚韧：眩晕时间-25% |
| **亡灵** | 智力+3, 暗影伤害+5% | 💀 亡者之触：攻击5%几率恐惧 | 🌑 暗影抗性：暗影伤害-15% |
| **牛头人** | 耐力+5, 最大HP+5% | ❤️ 坚忍：防御+3% | 🌿 自然亲和：受到治疗+10% |
| **巨魔** | 敏捷+3, 攻速+5% | 💚 再生：每回合恢复2%HP | 🐾 野兽杀手：野兽伤害+15% |

**种族特性触发机制:**
- 所有特性都是被动效果，无需玩家操作
- 属性加成在计算面板时自动应用
- 战斗特性由战斗引擎自动检测触发

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

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | ID |
| character_id | INTEGER | NOT NULL FK | 角色ID |
| name | VARCHAR(32) | NOT NULL | 策略名称 |
| priority | INTEGER | NOT NULL | 优先级(越小越优先) |
| condition_type | VARCHAR(32) | NOT NULL | 条件类型 |
| condition_operator | VARCHAR(8) | NOT NULL | 比较运算符 |
| condition_value | REAL | NOT NULL | 条件数值 |
| action_type | VARCHAR(32) | NOT NULL | 动作类型 |
| action_target | VARCHAR(32) | | 动作目标 |
| skill_id | VARCHAR(32) | | 使用技能 |
| item_id | VARCHAR(32) | | 使用物品 |
| is_active | INTEGER | DEFAULT 1 | 是否激活 |

```sql
CREATE TABLE battle_strategies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    priority INTEGER NOT NULL,
    condition_type VARCHAR(32) NOT NULL,
    condition_operator VARCHAR(8) NOT NULL,
    condition_value REAL NOT NULL,
    action_type VARCHAR(32) NOT NULL,
    action_target VARCHAR(32),
    skill_id VARCHAR(32),
    item_id VARCHAR(32),
    is_active INTEGER DEFAULT 1,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX idx_strategies_char_id ON battle_strategies(character_id);
```

**条件类型 (condition_type):**
- `self_hp_percent` - 自身HP百分比
- `self_mp_percent` - 自身MP百分比
- `enemy_hp_percent` - 敌人HP百分比
- `battle_round` - 战斗回合数
- `always` - 始终触发

**动作类型 (action_type):**
- `use_skill` - 使用技能
- `use_item` - 使用物品
- `normal_attack` - 普通攻击
- `flee` - 逃跑

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

### 二、装备强化系统

> 📌 **装备深度**: 通过强化、精炼、镶嵌、附魔让装备持续成长

#### 2.1 equipment_enhance - 装备强化表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| equipment_id | INTEGER | PRIMARY KEY FK | 装备记录ID |
| enhance_level | INTEGER | DEFAULT 0 | 强化等级 (0-15) |
| refine_level | INTEGER | DEFAULT 0 | 精炼等级 (0-10) |
| enchant_id | VARCHAR(32) | | 附魔ID |
| gem_slot_1 | VARCHAR(32) | | 宝石槽1 |
| gem_slot_2 | VARCHAR(32) | | 宝石槽2 |
| gem_slot_3 | VARCHAR(32) | | 宝石槽3 |

```sql
CREATE TABLE equipment_enhance (
    equipment_id INTEGER PRIMARY KEY,
    enhance_level INTEGER DEFAULT 0,
    refine_level INTEGER DEFAULT 0,
    enchant_id VARCHAR(32),
    gem_slot_1 VARCHAR(32),
    gem_slot_2 VARCHAR(32),
    gem_slot_3 VARCHAR(32),
    FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE,
    FOREIGN KEY (enchant_id) REFERENCES enchants(id),
    FOREIGN KEY (gem_slot_1) REFERENCES gems(id),
    FOREIGN KEY (gem_slot_2) REFERENCES gems(id),
    FOREIGN KEY (gem_slot_3) REFERENCES gems(id)
);
```

#### 2.2 enchants - 附魔配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 附魔ID |
| name | VARCHAR(32) | NOT NULL | 附魔名称 |
| description | TEXT | | 描述 |
| slot_type | VARCHAR(16) | | 适用槽位 |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_value | REAL | NOT NULL | 效果数值 |
| quality | VARCHAR(16) | | 品质 |

```sql
CREATE TABLE enchants (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    slot_type VARCHAR(16),
    effect_type VARCHAR(32) NOT NULL,
    effect_value REAL NOT NULL,
    quality VARCHAR(16)
);
```

#### 2.3 gems - 宝石配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 宝石ID |
| name | VARCHAR(32) | NOT NULL | 宝石名称 |
| color | VARCHAR(16) | NOT NULL | 颜色: red/blue/yellow/green/purple |
| stat_type | VARCHAR(32) | NOT NULL | 属性类型 |
| stat_value | INTEGER | NOT NULL | 属性数值 |
| quality | VARCHAR(16) | | 品质 |

```sql
CREATE TABLE gems (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    color VARCHAR(16) NOT NULL,
    stat_type VARCHAR(32) NOT NULL,
    stat_value INTEGER NOT NULL,
    quality VARCHAR(16)
);
```

#### 强化效果说明

| 系统 | 效果 | 上限 | 消耗 |
|-----|------|-----|------|
| **强化** | 基础属性 +3%/级 | +15 (+45%) | 金币 |
| **精炼** | 基础属性 +5%/级 | +10 (+50%) | 精炼石 |
| **附魔** | 特殊效果 | 1个 | 附魔材料 |
| **宝石** | 额外属性 | 3颗 | 宝石 |

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

