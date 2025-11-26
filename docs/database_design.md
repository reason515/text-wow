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
| mp | INTEGER | NOT NULL | 当前法力值 |
| max_mp | INTEGER | NOT NULL | 最大法力值 |
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

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 职业ID |
| name | VARCHAR(32) | NOT NULL | 职业名称 |
| description | TEXT | | 描述 |
| role | VARCHAR(16) | NOT NULL | 定位: tank/dps/healer |
| primary_stat | VARCHAR(16) | NOT NULL | 主属性 |
| base_hp | INTEGER | NOT NULL | 基础HP |
| base_mp | INTEGER | NOT NULL | 基础MP |
| hp_per_level | INTEGER | NOT NULL | 每级HP成长 |
| mp_per_level | INTEGER | NOT NULL | 每级MP成长 |
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
    base_hp INTEGER NOT NULL,
    base_mp INTEGER NOT NULL,
    hp_per_level INTEGER NOT NULL,
    mp_per_level INTEGER NOT NULL,
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
| mp_cost | INTEGER | DEFAULT 0 | 法力消耗 |
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
CREATE TABLE skills (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    class_id VARCHAR(32),
    type VARCHAR(16) NOT NULL,
    target VARCHAR(16) NOT NULL,
    damage_type VARCHAR(16),
    base_damage INTEGER DEFAULT 0,
    damage_scaling REAL DEFAULT 1.0,
    mp_cost INTEGER DEFAULT 0,
    cooldown INTEGER DEFAULT 0,
    level_required INTEGER DEFAULT 1,
    effect_type VARCHAR(32),
    effect_value REAL,
    effect_duration INTEGER,
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

CREATE INDEX idx_skills_class_id ON skills(class_id);
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

