# 📦 装备系统设计文档

> 📌 **核心设计理念**: 装备掉落时随机生成词缀，每件装备都有独特的属性组合。借鉴暗黑2的设计理念，装备系统是驱动玩家持续挑战的核心。

---

## 🎮 暗黑2风格设计理念

### 核心设计原则

1. **随机词缀系统**: 每件装备掉落时随机生成词缀，让每次掉落都有期待
2. **品质分级系统**: 从白色到橙色，清晰的品质体系
3. **稀有装备掉落**: 高品质装备稀有，驱动玩家持续刷怪
4. **装备强化系统**: 通过材料强化装备，提供成长路径
5. **装备交易系统**: 稀有装备可以交易，促进玩家间互动
6. **装备收集展示**: 装备图鉴和收藏系统，增加收集乐趣

### 暗黑2核心机制

- **随机词缀**: 前缀+后缀组合，每件装备都独特
- **品质决定词缀数**: 白色0个，绿色1个，蓝色2个，紫色3个，橙色4个
- **独特装备**: 固定词缀+特殊效果，类似暗金装备
- **稀有掉落率**: 高品质装备稀有，增加获取难度和成就感

---

## 📋 目录

1. [系统概览](#系统概览)
2. [数据库设计](#数据库设计)
3. [装备槽位定义](#装备槽位定义)
4. [词缀系统](#词缀系统)
5. [传说装备系统](#传说装备系统)
6. [材料强化系统](#材料强化系统)
7. [装备掉落系统](#装备掉落系统)
8. [装备属性加成](#装备属性加成)
9. [装备数值规范](#装备数值规范)
10. [背包管理](#背包管理)

---

## 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          装备系统架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   装备掉落 ──→ 随机词缀 ──→ 材料强化                                        │
│      │           │           │                                             │
│   基础属性    前缀+后缀    重铸/添加/强化                                     │
│                                                                             │
│   品质决定词缀数量:                                                          │
│   ⬜白(0) → 🟩绿(1) → 🟦蓝(2) → 🟪紫(3) → 🟧橙(4) → 🟨传说(固定词缀+特效) │
│                                                                             │
│   传说装备特点:                                                              │
│   ✨ 固定词缀组合（类似暗金装备）                                             │
│   ⭐ 独特的特殊效果                                                          │
│   💎 无法通过材料强化改变词缀                                                │
│   🎯 只能通过催化剂强化词缀数值                                               │
│                                                                             │
│   材料类型:                                                                  │
│   🔧基础材料 → 重铸/添加词缀                                                 │
│   💎精华材料 → 保证特定词缀                                                  │
│   ⚗️催化剂 → 强化已有词缀                                                    │
│   🛡️保护材料 → 锁定词缀不被改变                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 数据库设计

### 1. items - 物品配置表

> 📌 **用途说明**: 
> - 定义所有物品的**底材基础属性**（装备在没有词缀时的基础属性）
> - 这些基础属性会根据底材等级（normal/exceptional/elite）和品质（common/uncommon/rare/epic等）进行倍率调整
> - 词缀系统是独立的，通过`equipment_instance`表记录，词缀属性会叠加到底材基础属性上
> - 对于防具，基础防御值由护甲类型和槽位决定（见"防具基础防御值"章节），此表中的`physical_defense`和`magic_defense`为额外加成

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 物品ID |
| name | VARCHAR(64) | NOT NULL | 物品名称 |
| description | TEXT | | 描述 |
| type | VARCHAR(16) | NOT NULL | 类型: equipment/consumable/material/quest |
| subtype | VARCHAR(16) | | 子类型（武器类型：sword/axe/mace/dagger/fist/staff/wand/bow/crossbow/polearm等） |
| armor_type | VARCHAR(16) | | 护甲类型: cloth/leather/mail/plate (仅防具) |
| base_tier | VARCHAR(16) | DEFAULT 'normal' | 底材等级: normal/exceptional/elite |
| quality | VARCHAR(16) | DEFAULT 'common' | 品质: common/uncommon/rare/epic/legendary/unique |
| level_required | INTEGER | DEFAULT 1 | 需求等级 |
| class_required | VARCHAR(32) | | 需求职业 |
| strength_required | INTEGER | DEFAULT 0 | 需求力量（护甲类型要求） |
| agility_required | INTEGER | DEFAULT 0 | 需求敏捷（护甲类型要求） |
| intellect_required | INTEGER | DEFAULT 0 | 需求智力（护甲类型要求） |
| slot | VARCHAR(16) | | 装备槽位 |
| stackable | INTEGER | DEFAULT 0 | 可堆叠: 1是 0否 |
| max_stack | INTEGER | DEFAULT 1 | 最大堆叠数 |
| sell_price | INTEGER | DEFAULT 0 | 售价 |
| buy_price | INTEGER | DEFAULT 0 | 购买价 |
| **基础属性加成** | | | |
| strength | INTEGER | DEFAULT 0 | 力量加成 |
| agility | INTEGER | DEFAULT 0 | 敏捷加成 |
| intellect | INTEGER | DEFAULT 0 | 智力加成 |
| stamina | INTEGER | DEFAULT 0 | 耐力加成 |
| spirit | INTEGER | DEFAULT 0 | 精神加成 |
| **攻击属性** | | | |
| physical_attack | INTEGER | DEFAULT 0 | 物理攻击加成（武器） |
| magic_attack | INTEGER | DEFAULT 0 | 魔法攻击加成（法杖、魔杖） |
| **防御属性** | | | |
| physical_defense | INTEGER | DEFAULT 0 | 物理防御额外加成（防具，基础防御值由护甲类型和槽位决定） |
| magic_defense | INTEGER | DEFAULT 0 | 魔法防御额外加成（防具，基础防御值由护甲类型和槽位决定） |
| **生命与资源** | | | |
| hp_bonus | INTEGER | DEFAULT 0 | HP加成 |
| mp_bonus | INTEGER | DEFAULT 0 | MP加成 |
| **伤害加成（百分比）** | | | |
| physical_damage_pct | REAL | DEFAULT 0 | 物理伤害+% |
| fire_damage_pct | REAL | DEFAULT 0 | 火焰伤害+% |
| frost_damage_pct | REAL | DEFAULT 0 | 冰霜伤害+% |
| lightning_damage_pct | REAL | DEFAULT 0 | 雷电伤害+% |
| shadow_damage_pct | REAL | DEFAULT 0 | 暗影伤害+% |
| holy_damage_pct | REAL | DEFAULT 0 | 神圣伤害+% |
| nature_damage_pct | REAL | DEFAULT 0 | 自然伤害+% |
| elemental_damage_pct | REAL | DEFAULT 0 | 全元素伤害+% |
| **伤害加成（固定值）** | | | |
| physical_damage_flat | INTEGER | DEFAULT 0 | 物理伤害+（固定数值） |
| fire_damage_flat | INTEGER | DEFAULT 0 | 火焰伤害+（固定数值） |
| frost_damage_flat | INTEGER | DEFAULT 0 | 冰霜伤害+（固定数值） |
| lightning_damage_flat | INTEGER | DEFAULT 0 | 雷电伤害+（固定数值） |
| shadow_damage_flat | INTEGER | DEFAULT 0 | 暗影伤害+（固定数值） |
| holy_damage_flat | INTEGER | DEFAULT 0 | 神圣伤害+（固定数值） |
| nature_damage_flat | INTEGER | DEFAULT 0 | 自然伤害+（固定数值） |
| **抗性加成** | | | |
| physical_resistance_pct | REAL | DEFAULT 0 | 物理抗性+% |
| fire_resistance_pct | REAL | DEFAULT 0 | 火焰抗性+% |
| frost_resistance_pct | REAL | DEFAULT 0 | 冰霜抗性+% |
| lightning_resistance_pct | REAL | DEFAULT 0 | 雷电抗性+% |
| shadow_resistance_pct | REAL | DEFAULT 0 | 暗影抗性+% |
| holy_resistance_pct | REAL | DEFAULT 0 | 神圣抗性+% |
| nature_resistance_pct | REAL | DEFAULT 0 | 自然抗性+% |
| elemental_resistance_pct | REAL | DEFAULT 0 | 全元素抗性+% |
| **其他属性** | | | |
| crit_rate | REAL | DEFAULT 0 | 暴击率加成 |
| crit_damage_pct | REAL | DEFAULT 0 | 暴击伤害+% |
| dodge_rate | REAL | DEFAULT 0 | 闪避率+% |
| damage_reduction_pct | REAL | DEFAULT 0 | 受伤减免+%（盾牌基础属性） |
| resource_gain_pct | REAL | DEFAULT 0 | 资源获取+% |
| initiative | INTEGER | DEFAULT 0 | 先手值+ |
| **使用效果** | | | |
| effect_type | VARCHAR(32) | | 使用效果类型（消耗品） |
| effect_value | INTEGER | | 使用效果数值（消耗品） |

```sql
CREATE TABLE items (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    type VARCHAR(16) NOT NULL,
    subtype VARCHAR(16),
    armor_type VARCHAR(16),
    base_tier VARCHAR(16) DEFAULT 'normal',
    quality VARCHAR(16) DEFAULT 'common',
    level_required INTEGER DEFAULT 1,
    class_required VARCHAR(32),
    strength_required INTEGER DEFAULT 0,
    agility_required INTEGER DEFAULT 0,
    intellect_required INTEGER DEFAULT 0,
    slot VARCHAR(16),
    stackable INTEGER DEFAULT 0,
    max_stack INTEGER DEFAULT 1,
    sell_price INTEGER DEFAULT 0,
    buy_price INTEGER DEFAULT 0,
    -- 基础属性加成
    strength INTEGER DEFAULT 0,
    agility INTEGER DEFAULT 0,
    intellect INTEGER DEFAULT 0,
    stamina INTEGER DEFAULT 0,
    spirit INTEGER DEFAULT 0,
    -- 攻击属性
    physical_attack INTEGER DEFAULT 0,
    magic_attack INTEGER DEFAULT 0,
    -- 防御属性（额外加成，基础防御值由护甲类型和槽位决定）
    physical_defense INTEGER DEFAULT 0,
    magic_defense INTEGER DEFAULT 0,
    -- 生命与资源
    hp_bonus INTEGER DEFAULT 0,
    mp_bonus INTEGER DEFAULT 0,
    -- 伤害加成（百分比）
    physical_damage_pct REAL DEFAULT 0,
    fire_damage_pct REAL DEFAULT 0,
    frost_damage_pct REAL DEFAULT 0,
    lightning_damage_pct REAL DEFAULT 0,
    shadow_damage_pct REAL DEFAULT 0,
    holy_damage_pct REAL DEFAULT 0,
    nature_damage_pct REAL DEFAULT 0,
    elemental_damage_pct REAL DEFAULT 0,
    -- 伤害加成（固定值）
    physical_damage_flat INTEGER DEFAULT 0,
    fire_damage_flat INTEGER DEFAULT 0,
    frost_damage_flat INTEGER DEFAULT 0,
    lightning_damage_flat INTEGER DEFAULT 0,
    shadow_damage_flat INTEGER DEFAULT 0,
    holy_damage_flat INTEGER DEFAULT 0,
    nature_damage_flat INTEGER DEFAULT 0,
    -- 抗性加成
    physical_resistance_pct REAL DEFAULT 0,
    fire_resistance_pct REAL DEFAULT 0,
    frost_resistance_pct REAL DEFAULT 0,
    lightning_resistance_pct REAL DEFAULT 0,
    shadow_resistance_pct REAL DEFAULT 0,
    holy_resistance_pct REAL DEFAULT 0,
    nature_resistance_pct REAL DEFAULT 0,
    elemental_resistance_pct REAL DEFAULT 0,
    -- 其他属性
    crit_rate REAL DEFAULT 0,
    crit_damage_pct REAL DEFAULT 0,
    dodge_rate REAL DEFAULT 0,
    damage_reduction_pct REAL DEFAULT 0,
    resource_gain_pct REAL DEFAULT 0,
    initiative INTEGER DEFAULT 0,
    -- 使用效果（消耗品）
    effect_type VARCHAR(32),
    effect_value INTEGER
);

CREATE INDEX idx_items_type ON items(type);
CREATE INDEX idx_items_quality ON items(quality);
CREATE INDEX idx_items_armor_type ON items(armor_type);
CREATE INDEX idx_items_base_tier ON items(base_tier);
CREATE INDEX idx_items_slot ON items(slot);
```

**items表字段说明：**

1. **基础属性加成**：力量、敏捷、智力、耐力、精神
   - 这些属性会直接加到角色基础属性上
   - 根据底材等级和品质进行倍率调整

2. **攻击属性**：
   - `physical_attack`：物理攻击加成（武器）
   - `magic_attack`：魔法攻击加成（法杖、魔杖）

3. **防御属性**：
   - `physical_defense`：物理防御额外加成（防具）
   - `magic_defense`：魔法防御额外加成（防具）
   - **注意**：防具的基础防御值由护甲类型和槽位决定（见"防具基础防御值"章节），此字段为额外加成

4. **伤害加成**：
   - 百分比加成：`*_damage_pct`（如`fire_damage_pct`）
   - 固定值加成：`*_damage_flat`（如`fire_damage_flat`）

5. **抗性加成**：
   - 各种元素抗性和物理抗性的百分比加成

6. **其他属性**：
   - 暴击率、暴击伤害、闪避率、受伤减免、资源获取、先手值等

---

### 2. equipment_instance - 装备实例表

> 📌 记录玩家获得的每一件装备及其词缀

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 装备实例ID |
| item_id | VARCHAR(32) | NOT NULL FK | 基础物品ID |
| owner_id | INTEGER | NOT NULL FK | 拥有者用户ID |
| character_id | INTEGER | FK | 装备者角色ID (NULL=背包中) |
| slot | VARCHAR(16) | | 装备槽位 |
| quality | VARCHAR(16) | NOT NULL | 品质: common/uncommon/rare/epic/legendary/unique |
| prefix_id | VARCHAR(32) | FK | 前缀词缀ID |
| prefix_value | REAL | | 前缀数值 (词缀效果的具体数值) |
| suffix_id | VARCHAR(32) | FK | 后缀词缀ID |
| suffix_value | REAL | | 后缀数值 |
| bonus_affix_1 | VARCHAR(32) | FK | 额外词缀1 (紫色+) |
| bonus_affix_1_value | REAL | | 额外词缀1数值 |
| bonus_affix_2 | VARCHAR(32) | FK | 额外词缀2 (橙色+) |
| bonus_affix_2_value | REAL | | 额外词缀2数值 |
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
    prefix_id VARCHAR(32),
    prefix_value REAL,
    suffix_id VARCHAR(32),
    suffix_value REAL,
    bonus_affix_1 VARCHAR(32),
    bonus_affix_1_value REAL,
    bonus_affix_2 VARCHAR(32),
    bonus_affix_2_value REAL,
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_locked INTEGER DEFAULT 0,
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE SET NULL,
    FOREIGN KEY (prefix_id) REFERENCES affixes(id),
    FOREIGN KEY (suffix_id) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_1) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_2) REFERENCES affixes(id)
);

CREATE INDEX idx_equipment_owner ON equipment_instance(owner_id);
CREATE INDEX idx_equipment_character ON equipment_instance(character_id);
CREATE INDEX idx_equipment_quality ON equipment_instance(quality);
```

---

### 3. equipment - 装备表（当前装备）

> 📌 记录角色当前装备的物品

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

### 4. inventory - 背包表

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

## 装备槽位定义

| 槽位 | 英文 | 说明 |
|-----|------|------|
| head | 头部 | 头盔、帽子（可装备布甲/皮甲/锁甲/板甲） |
| armor | 盔甲 | 胸甲和裤子（合并为一个槽位，可装备布甲/皮甲/锁甲/板甲） |
| hands | 手套 | 手套（可装备布甲/皮甲/锁甲/板甲） |
| feet | 靴子 | 靴子（可装备布甲/皮甲/锁甲/板甲） |
| main_hand | 主手 | 主手武器（单手/双手武器） |
| off_hand | 副手 | 副手武器或盾牌（盾牌不区分护甲类型） |
| accessory | 首饰 | 首饰（可装备2个，不再区分戒指和项链） |

**护甲类型说明：**
- 头部、盔甲、手套、靴子可以装备布甲、皮甲、锁甲、板甲中的任意一种
- 不同护甲类型提供不同比例的物理防御和魔法防御
- 每种护甲类型对角色基础属性有要求（不满足时防御值降低50%）
- 盾牌不区分护甲类型，统一提供物理防御和魔法防御

---

## 词缀系统

### affixes - 词缀配置表

> 📌 定义所有可能的装备词缀

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 词缀ID |
| name | VARCHAR(32) | NOT NULL | 词缀名称 |
| type | VARCHAR(16) | NOT NULL | 类型: prefix/suffix |
| slot_type | VARCHAR(16) | | 适用槽位: weapon/armor/accessory/all (weapon=主手/副手, armor=头部/盔甲/手套/靴子, accessory=首饰) |
| rarity | VARCHAR(16) | NOT NULL | 稀有度: common/uncommon/rare/epic |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_stat | VARCHAR(32) | | 影响的属性 |
| min_value | REAL | NOT NULL | 最小数值 |
| max_value | REAL | NOT NULL | 最大数值 |
| value_type | VARCHAR(16) | NOT NULL | 数值类型: flat/percent |
| description | TEXT | | 描述模板 (用{value}占位) |
| level_required | INTEGER | DEFAULT 1 | 最低出现等级 |
| tier | INTEGER | DEFAULT 1 | 词缀等级 (1-5) |

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
    level_required INTEGER DEFAULT 1,
    tier INTEGER DEFAULT 1
);

CREATE INDEX idx_affixes_type ON affixes(type);
CREATE INDEX idx_affixes_rarity ON affixes(rarity);
CREATE INDEX idx_affixes_slot ON affixes(slot_type);
```

---

### 词缀分级系统

> 📌 **设计理念**: 根据装备等级，词缀有不同的数值范围，高等级装备可以获得更强的词缀

**词缀等级 (Tier) 与装备等级对应关系：**

| 词缀等级 | 装备等级范围 | 数值倍率 | 说明 |
|---------|------------|---------|------|
| Tier 1 | 1-20级 | 1.0x | 基础词缀 |
| Tier 2 | 21-40级 | 1.8x | 数值提升80% |
| Tier 3 | 41-60级 | 2.5x | 数值提升150% |

**词缀生成规则：**
- 装备掉落时，根据装备等级随机生成对应Tier的词缀
- 低等级装备只能获得低Tier词缀
- 高等级装备可能获得低Tier或高Tier词缀（高Tier概率较低）
- 词缀数值 = 基础数值 × Tier倍率
- **不同Tier使用不同的词缀名称，便于直观识别**

**词缀命名规则（创意命名）：**
- Tier 1: 基础名称（如"of 力量"、"锋利的"）
- Tier 2: 进阶名称（使用更有创意的词汇，如"of 蛮力"、"寒芒的"）
- Tier 3: 传说名称（使用史诗感强的词汇，如"of 泰坦之力"、"龙牙的"）

**命名风格：**
- 力量系：力量 → 蛮力 → 泰坦之力
- 敏捷系：敏捷 → 疾风 → 影舞
- 智力系：智力 → 睿智 → 奥术之心
- 攻击系：锋利的 → 寒芒的 → 龙牙的
- 元素系：炽热的 → 熔岩的 → 凤凰之怒
- 抗性系：火焰抗性 → 烈焰护体 → 不灭之火

**示例：**
```
力量词缀分级：
Tier 1 (1-20级):   of 力量 +1~3
Tier 2 (21-40级):  of 蛮力 +2~5   (1.8倍)
Tier 3 (41-60级):  of 泰坦之力 +3~8 (2.5倍)

物理伤害词缀分级：
Tier 1 (1-20级):   of 物理伤害 +3~8%
Tier 2 (21-40级):  of 撕裂 +5~14%
Tier 3 (41-60级):  of 粉碎 +8~20%

攻击力词缀分级（前缀）：
Tier 1 (1-20级):   锋利的 +2~5
Tier 2 (21-40级):  寒芒的 +4~9
Tier 3 (41-60级):  龙牙的 +5~13

火焰抗性词缀分级：
Tier 1 (1-20级):   of 火焰抗性 +3~10%
Tier 2 (21-40级):  of 烈焰护体 +5~18%
Tier 3 (41-60级):  of 不灭之火 +8~25%

敏捷词缀分级：
Tier 1 (1-20级):   of 敏捷 +1~3
Tier 2 (21-40级):  of 疾风 +2~5
Tier 3 (41-60级):  of 影舞 +3~8

智力词缀分级：
Tier 1 (1-20级):   of 智力 +1~3
Tier 2 (21-40级):  of 睿智 +2~5
Tier 3 (41-60级):  of 奥术之心 +3~8

火焰伤害词缀分级（前缀）：
Tier 1 (1-20级):   炽热的 +1~4
Tier 2 (21-40级):  熔岩的 +2~7
Tier 3 (41-60级):  凤凰之怒 +3~10
```

**词缀Tier生成概率：**
- 装备等级在Tier范围内：100%生成对应Tier
- 装备等级高于Tier范围：有15%概率生成更高Tier（最多高1级）
- 装备等级低于Tier范围：无法生成该Tier词缀

---

### 词缀列表

#### 前缀 (攻击/属性向)

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_sharp | 锋利的 | 攻击力+ | 2~5 (T1) | 普通 | 武器 |
| affix_sharp_t2 | 寒芒的 | 攻击力+ | 4~9 (T2) | 普通 | 武器 |
| affix_sharp_t3 | 龙牙的 | 攻击力+ | 5~13 (T3) | 普通 | 武器 |
| affix_physical | 物理的 | 物理伤害+ | 2~5 (T1) | 普通 | 武器 |
| affix_physical_t2 | 撕裂的 | 物理伤害+ | 4~9 (T2) | 普通 | 武器 |
| affix_physical_t3 | 粉碎的 | 物理伤害+ | 5~13 (T3) | 普通 | 武器 |
| affix_fiery | 炽热的 | 火焰伤害+ | 1~4 (T1) | 普通 | 武器 |
| affix_fiery_t2 | 熔岩的 | 火焰伤害+ | 2~7 (T2) | 普通 | 武器 |
| affix_fiery_t3 | 凤凰之怒 | 火焰伤害+ | 3~10 (T3) | 普通 | 武器 |
| affix_frozen | 冰霜的 | 冰霜伤害+ | 1~4 (T1) | 普通 | 武器 |
| affix_frozen_t2 | 寒冰的 | 冰霜伤害+ | 2~7 (T2) | 普通 | 武器 |
| affix_frozen_t3 | 冰龙之息 | 冰霜伤害+ | 3~10 (T3) | 普通 | 武器 |
| affix_charged | 雷击的 | 雷电伤害+ | 1~4 (T1) | 普通 | 武器 |
| affix_charged_t2 | 闪电的 | 雷电伤害+ | 2~7 (T2) | 普通 | 武器 |
| affix_charged_t3 | 雷神之怒 | 雷电伤害+ | 3~10 (T3) | 普通 | 武器 |
| affix_holy | 神圣的 | 神圣伤害+ | 2~5 (T1) | 精良 | 武器 |
| affix_holy_t2 | 圣光的 | 神圣伤害+ | 4~9 (T2) | 精良 | 武器 |
| affix_holy_t3 | 天使之翼 | 神圣伤害+ | 5~13 (T3) | 精良 | 武器 |
| affix_vampiric | 吸血鬼的 | 生命偷取% | 2~5 (T1) | 稀有 | 武器 |
| affix_vampiric_t2 | 嗜血的 | 生命偷取% | 4~9 (T2) | 稀有 | 武器 |
| affix_vampiric_t3 | 血魔之触 | 生命偷取% | 5~13 (T3) | 稀有 | 武器 |
| affix_devastating | 毁灭的 | 攻击力+% | 15~25 (T1) | 史诗 | 武器 |
| affix_devastating_t2 | 天灾的 | 攻击力+% | 27~45 (T2) | 史诗 | 武器 |
| affix_devastating_t3 | 灭世之威 | 攻击力+% | 38~63 (T3) | 史诗 | 武器 |
| affix_vital | 活力的 | 生命值+ | 5~15 (T1) | 普通 | 防具 |
| affix_vital_t2 | 强健的 | 生命值+ | 9~27 (T2) | 普通 | 防具 |
| affix_vital_t3 | 生命之源 | 生命值+ | 13~38 (T3) | 普通 | 防具 |
| affix_scholarly | 智者的 | 智力+ | 2~4 (T1) | 精良 | 防具 |
| affix_scholarly_t2 | 博学的 | 智力+ | 4~7 (T2) | 精良 | 防具 |
| affix_scholarly_t3 | 奥术大师 | 智力+ | 5~10 (T3) | 精良 | 防具 |
| affix_unyielding | 不屈的 | 受伤减免% | 3~8 (T1) | 稀有 | 防具 |
| affix_unyielding_t2 | 坚韧的 | 受伤减免% | 5~14 (T2) | 稀有 | 防具 |
| affix_unyielding_t3 | 不灭之盾 | 受伤减免% | 8~20 (T3) | 稀有 | 防具 |

#### 后缀 (特殊效果向)

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_strength | of 力量 | 力量+ | 1~3 (T1) | 普通 | 全部 |
| affix_of_strength_t2 | of 蛮力 | 力量+ | 2~5 (T2) | 普通 | 全部 |
| affix_of_strength_t3 | of 泰坦之力 | 力量+ | 3~8 (T3) | 普通 | 全部 |
| affix_of_physical_damage | of 物理伤害 | 物理伤害+% | 3~8% (T1) | 普通 | 武器 |
| affix_of_physical_damage_t2 | of 撕裂 | 物理伤害+% | 5~14% (T2) | 普通 | 武器 |
| affix_of_physical_damage_t3 | of 粉碎 | 物理伤害+% | 8~20% (T3) | 普通 | 武器 |
| affix_of_haste | of 迅捷 | 资源获取+% | 10~20 (T1) | 精良 | 武器 |
| affix_of_haste_t2 | of 疾风 | 资源获取+% | 18~36 (T2) | 精良 | 武器 |
| affix_of_haste_t3 | of 风暴之速 | 资源获取+% | 25~50 (T3) | 精良 | 武器 |
| affix_of_piercing | of 穿刺 | 无视防御% | 5~15 (T1) | 精良 | 武器 |
| affix_of_piercing_t2 | of 穿透 | 无视防御% | 9~27 (T2) | 精良 | 武器 |
| affix_of_piercing_t3 | of 破甲之锋 | 无视防御% | 13~38 (T3) | 精良 | 武器 |
| affix_of_crit | of 暴击 | 暴击率+% | 3~8 (T1) | 稀有 | 武器 |
| affix_of_crit_t2 | of 致命 | 暴击率+% | 5~14 (T2) | 稀有 | 武器 |
| affix_of_crit_t3 | of 必杀 | 暴击率+% | 8~20 (T3) | 稀有 | 武器 |
| affix_of_lethality | of 致命 | 暴击伤害+% | 10~25 (T1) | 稀有 | 武器 |
| affix_of_lethality_t2 | of 杀戮 | 暴击伤害+% | 18~45 (T2) | 稀有 | 武器 |
| affix_of_lethality_t3 | of 绝杀 | 暴击伤害+% | 25~63 (T3) | 稀有 | 武器 |
| affix_of_leech | of 吸血 | 伤害转HP% | 2~4 (T1) | 稀有 | 武器 |
| affix_of_leech_t2 | of 嗜血 | 伤害转HP% | 4~7 (T2) | 稀有 | 武器 |
| affix_of_leech_t3 | of 血魔 | 伤害转HP% | 5~10 (T3) | 稀有 | 武器 |
| affix_of_blocking | of 守护 | 受伤减免+% | 3~6 (T1) | 精良 | 副手 |
| affix_of_blocking_t2 | of 防护 | 受伤减免+% | 5~11 (T2) | 精良 | 副手 |
| affix_of_blocking_t3 | of 坚盾 | 受伤减免+% | 8~15 (T3) | 精良 | 副手 |
| affix_of_thorns | of 反射 | 反弹伤害% | 5~15 (T1) | 稀有 | 防具 |
| affix_of_thorns_t2 | of 荆棘 | 反弹伤害% | 9~27 (T2) | 稀有 | 防具 |
| affix_of_thorns_t3 | of 荆棘之甲 | 反弹伤害% | 13~38 (T3) | 稀有 | 防具 |
| affix_of_regen | of 再生 | 每回合恢复HP | 1~3 (T1) | 稀有 | 防具 |
| affix_of_regen_t2 | of 恢复 | 每回合恢复HP | 2~5 (T2) | 稀有 | 防具 |
| affix_of_regen_t3 | of 生命之泉 | 每回合恢复HP | 3~8 (T3) | 稀有 | 防具 |
| affix_of_wisdom | of 智慧 | 法力恢复+% | 10~20 (T1) | 精良 | 防具 |
| affix_of_wisdom_t2 | of 睿智 | 法力恢复+% | 18~36 (T2) | 精良 | 防具 |
| affix_of_wisdom_t3 | of 奥术之源 | 法力恢复+% | 25~50 (T3) | 精良 | 防具 |

#### 暴击相关词缀

| 词缀 | 效果 | 适用装备 |
|-----|------|---------|
| of 暴击 | +3-8% 物理暴击率 | 武器 |
| of 致命 | +10-25% 暴击伤害 | 武器 |
| of 法术暴击 | +3-8% 法术暴击率 | 武器/首饰 |
| of 法术致命 | +10-25% 法术暴击伤害 | 武器/首饰 |

#### 闪避相关词缀

| 词缀 | 效果 | 适用装备 |
|-----|------|---------|
| of 敏捷 | +5-15 敏捷（间接增加闪避） | 盔甲/首饰 |
| of 闪避 | +2-5% 闪避率 | 盔甲/首饰 |
| of 灵巧 | +3-8% 闪避率 | 皮甲 |

#### 仇恨相关词缀

**前缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| 守护者的 | 仇恨生成+20%, 嘲讽资源消耗-20% | 坦克武器 | 稀有 |
| 威压的 | 仇恨生成+15% | 坦克装备 | 精良 |
| 隐秘的 | 仇恨生成-20% | 输出武器 | 精良 |
| 暗影的 | 暴击仇恨-30% | 输出装备 | 稀有 |

**后缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| of 威胁 | 仇恨生成+15% | 坦克装备 | 精良 |
| of 守护 | 仇恨生成+10%, 受伤减免+3% | 副手 | 精良 |
| of 隐匿 | 仇恨生成-15% | 输出装备 | 精良 |
| of 消散 | 仇恨衰减+20% | 治疗装备 | 稀有 |

---

#### 掉落与收益相关词缀

> 📌 这些词缀主要出现在首饰和特定装备上，提升玩家的收益和效率

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_lucky | 幸运的 | 物品掉落率+% | 5~15 (T1) | 精良 | 首饰 |
| affix_lucky_t2 | 祝福的 | 物品掉落率+% | 9~27 (T2) | 精良 | 首饰 |
| affix_lucky_t3 | 天运之佑 | 物品掉落率+% | 13~38 (T3) | 精良 | 首饰 |
| affix_blessed_drop | 祝福的 | 装备掉落率+% | 3~8 (T1) | 稀有 | 首饰 |
| affix_blessed_drop_t2 | 恩赐的 | 装备掉落率+% | 5~14 (T2) | 稀有 | 首饰 |
| affix_blessed_drop_t3 | 神恩之赐 | 装备掉落率+% | 8~20 (T3) | 稀有 | 首饰 |
| affix_wealthy | 富有的 | 金币获取+% | 10~25 (T1) | 精良 | 首饰 |
| affix_wealthy_t2 | 富足的 | 金币获取+% | 18~45 (T2) | 精良 | 首饰 |
| affix_wealthy_t3 | 黄金之触 | 金币获取+% | 25~63 (T3) | 精良 | 首饰 |
| affix_wise | 智慧的 | 经验获取+% | 10~20 (T1) | 精良 | 首饰 |
| affix_wise_t2 | 博学的 | 经验获取+% | 18~36 (T2) | 精良 | 首饰 |
| affix_wise_t3 | 智慧之光 | 经验获取+% | 25~50 (T3) | 精良 | 首饰 |
| affix_prosperous | 繁荣的 | 金币+经验获取+% | 8~15 (T1) | 稀有 | 首饰 |
| affix_prosperous_t2 | 丰饶的 | 金币+经验获取+% | 14~27 (T2) | 稀有 | 首饰 |
| affix_prosperous_t3 | 繁荣之冠 | 金币+经验获取+% | 20~38 (T3) | 稀有 | 首饰 |
| affix_magical | 魔法的 | 装备稀有度提升 | +1品质等级 | 史诗 | 首饰 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_fortune | of 财富 | 金币获取+% | 5~12 (T1) | 普通 | 首饰 |
| affix_of_fortune_t2 | of 富足 | 金币获取+% | 9~22 (T2) | 普通 | 首饰 |
| affix_of_fortune_t3 | of 黄金之触 | 金币获取+% | 13~30 (T3) | 普通 | 首饰 |
| affix_of_knowledge | of 知识 | 经验获取+% | 5~12 (T1) | 普通 | 首饰 |
| affix_of_knowledge_t2 | of 智慧 | 经验获取+% | 9~22 (T2) | 普通 | 首饰 |
| affix_of_knowledge_t3 | of 智慧之光 | 经验获取+% | 13~30 (T3) | 普通 | 首饰 |
| affix_of_plenty | of 丰饶 | 物品掉落率+% | 3~8 (T1) | 精良 | 首饰 |
| affix_of_plenty_t2 | of 丰收 | 物品掉落率+% | 5~14 (T2) | 精良 | 首饰 |
| affix_of_plenty_t3 | of 天运之佑 | 物品掉落率+% | 8~20 (T3) | 精良 | 首饰 |
| affix_of_treasure | of 宝藏 | 装备掉落率+% | 2~5 (T1) | 稀有 | 首饰 |
| affix_of_treasure_t2 | of 秘宝 | 装备掉落率+% | 4~9 (T2) | 稀有 | 首饰 |
| affix_of_treasure_t3 | of 神恩之赐 | 装备掉落率+% | 5~13 (T3) | 稀有 | 首饰 |
| affix_of_riches | of 富足 | 金币+经验+% | 6~12 (T1) | 精良 | 首饰 |
| affix_of_riches_t2 | of 繁荣 | 金币+经验+% | 11~22 (T2) | 精良 | 首饰 |
| affix_of_riches_t3 | of 繁荣之冠 | 金币+经验+% | 15~30 (T3) | 精良 | 首饰 |
| affix_of_legend | of 传说 | 装备稀有度提升 | +1品质等级 | 史诗 | 首饰 |

---

#### 效率与速度相关词缀

> 📌 提升战斗和探索效率的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_swift | 迅捷的 | 资源获取+% | 8~15 (T1) | 普通 | 武器 |
| affix_swift_t2 | 疾风的 | 资源获取+% | 14~27 (T2) | 普通 | 武器 |
| affix_swift_t3 | 风暴之速 | 资源获取+% | 20~38 (T3) | 普通 | 武器 |
| affix_energetic | 精力充沛的 | 战斗后休息时间-% | 10~25 (T1) | 精良 | 防具/首饰 |
| affix_energetic_t2 | 充满活力的 | 战斗后休息时间-% | 18~45 (T2) | 精良 | 防具/首饰 |
| affix_energetic_t3 | 永动之体 | 战斗后休息时间-% | 25~63 (T3) | 精良 | 防具/首饰 |
| affix_vigorous | 充满活力的 | 战斗后休息时间-% | 15~30 (T1) | 稀有 | 防具/首饰 |
| affix_vigorous_t2 | 强健的 | 战斗后休息时间-% | 27~54 (T2) | 稀有 | 防具/首饰 |
| affix_vigorous_t3 | 不竭之力 | 战斗后休息时间-% | 38~75 (T3) | 稀有 | 防具/首饰 |
| affix_rapid | 快速的 | 先手值+ | 3~8 (T1) | 精良 | 靴子 |
| affix_rapid_t2 | 疾行的 | 先手值+ | 5~14 (T2) | 精良 | 靴子 |
| affix_rapid_t3 | 影舞之步 | 先手值+ | 8~20 (T3) | 精良 | 靴子 |
| affix_lightning | 闪电般的 | 资源获取+% + 先手值+ | 12~20% + 3~8 (T1) | 稀有 | 武器/靴子 |
| affix_lightning_t2 | 雷霆的 | 资源获取+% + 先手值+ | 22~36% + 5~14 (T2) | 稀有 | 武器/靴子 |
| affix_lightning_t3 | 雷神之速 | 资源获取+% + 先手值+ | 30~50% + 8~20 (T3) | 稀有 | 武器/靴子 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_swiftness | of 敏捷 | 资源获取+% | 5~12 (T1) | 普通 | 武器 |
| affix_of_swiftness_t2 | of 疾风 | 资源获取+% | 9~22 (T2) | 普通 | 武器 |
| affix_of_swiftness_t3 | of 风暴之速 | 资源获取+% | 13~30 (T3) | 普通 | 武器 |
| affix_of_rest | of 休息 | 战斗后休息时间-% | 8~15 (T1) | 精良 | 防具/首饰 |
| affix_of_rest_t2 | of 恢复 | 战斗后休息时间-% | 14~27 (T2) | 精良 | 防具/首饰 |
| affix_of_rest_t3 | of 永动之体 | 战斗后休息时间-% | 20~38 (T3) | 精良 | 防具/首饰 |
| affix_of_recovery | of 恢复 | 战斗后休息时间-% | 12~20 (T1) | 稀有 | 防具/首饰 |
| affix_of_recovery_t2 | of 强健 | 战斗后休息时间-% | 22~36 (T2) | 稀有 | 防具/首饰 |
| affix_of_recovery_t3 | of 不竭之力 | 战斗后休息时间-% | 30~50 (T3) | 稀有 | 防具/首饰 |
| affix_of_speed | of 速度 | 先手值+ | 2~6 (T1) | 精良 | 靴子 |
| affix_of_speed_t2 | of 疾行 | 先手值+ | 4~11 (T2) | 精良 | 靴子 |
| affix_of_speed_t3 | of 影舞之步 | 先手值+ | 5~15 (T3) | 精良 | 靴子 |
| affix_of_quickness | of 迅速 | 资源获取+% + 先手值+ | 8~15% + 2~5 (T1) | 稀有 | 武器/靴子 |
| affix_of_quickness_t2 | of 雷霆 | 资源获取+% + 先手值+ | 14~27% + 4~9 (T2) | 稀有 | 武器/靴子 |
| affix_of_quickness_t3 | of 雷神之速 | 资源获取+% + 先手值+ | 20~38% + 5~13 (T3) | 稀有 | 武器/靴子 |

---

#### 属性增强词缀

> 📌 直接提升角色基础属性的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_mighty | 强大的 | 力量+ | 2~5 (T1) | 普通 | 武器/防具 |
| affix_mighty_t2 | 蛮力的 | 力量+ | 4~9 (T2) | 普通 | 武器/防具 |
| affix_mighty_t3 | 泰坦之力 | 力量+ | 5~13 (T3) | 普通 | 武器/防具 |
| affix_nimble | 灵巧的 | 敏捷+ | 2~5 (T1) | 普通 | 武器/防具 |
| affix_nimble_t2 | 疾风的 | 敏捷+ | 4~9 (T2) | 普通 | 武器/防具 |
| affix_nimble_t3 | 影舞 | 敏捷+ | 5~13 (T3) | 普通 | 武器/防具 |
| affix_brilliant | 聪慧的 | 智力+ | 2~5 (T1) | 普通 | 武器/防具 |
| affix_brilliant_t2 | 睿智的 | 智力+ | 4~9 (T2) | 普通 | 武器/防具 |
| affix_brilliant_t3 | 奥术之心 | 智力+ | 5~13 (T3) | 普通 | 武器/防具 |
| affix_robust | 强壮的 | 耐力+ | 2~5 (T1) | 普通 | 防具 |
| affix_robust_t2 | 坚韧的 | 耐力+ | 4~9 (T2) | 普通 | 防具 |
| affix_robust_t3 | 山岳之体 | 耐力+ | 5~13 (T3) | 普通 | 防具 |
| affix_spiritual | 精神的 | 精神+ | 2~5 (T1) | 普通 | 防具/首饰 |
| affix_spiritual_t2 | 灵性的 | 精神+ | 4~9 (T2) | 普通 | 防具/首饰 |
| affix_spiritual_t3 | 灵魂之火 | 精神+ | 5~13 (T3) | 普通 | 防具/首饰 |
| affix_balanced | 平衡的 | 全属性+ | 1~3 (T1) | 稀有 | 首饰 |
| affix_balanced_t2 | 全能的 | 全属性+ | 2~5 (T2) | 稀有 | 首饰 |
| affix_balanced_t3 | 万象归一 | 全属性+ | 3~8 (T3) | 稀有 | 首饰 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_strength_enhanced | of 力量 | 力量+ | 1~4 (T1) | 普通 | 全部 |
| affix_of_strength_enhanced_t2 | of 蛮力 | 力量+ | 2~7 (T2) | 普通 | 全部 |
| affix_of_strength_enhanced_t3 | of 泰坦之力 | 力量+ | 3~10 (T3) | 普通 | 全部 |
| affix_of_agility_enhanced | of 敏捷 | 敏捷+ | 1~4 (T1) | 普通 | 全部 |
| affix_of_agility_enhanced_t2 | of 疾风 | 敏捷+ | 2~7 (T2) | 普通 | 全部 |
| affix_of_agility_enhanced_t3 | of 影舞 | 敏捷+ | 3~10 (T3) | 普通 | 全部 |
| affix_of_intellect | of 智力 | 智力+ | 1~4 (T1) | 普通 | 全部 |
| affix_of_intellect_t2 | of 睿智 | 智力+ | 2~7 (T2) | 普通 | 全部 |
| affix_of_intellect_t3 | of 奥术之心 | 智力+ | 3~10 (T3) | 普通 | 全部 |
| affix_of_stamina | of 耐力 | 耐力+ | 1~4 (T1) | 普通 | 防具 |
| affix_of_stamina_t2 | of 坚韧 | 耐力+ | 2~7 (T2) | 普通 | 防具 |
| affix_of_stamina_t3 | of 山岳之体 | 耐力+ | 3~10 (T3) | 普通 | 防具 |
| affix_of_spirit | of 精神 | 精神+ | 1~4 (T1) | 普通 | 防具/首饰 |
| affix_of_spirit_t2 | of 灵性 | 精神+ | 2~7 (T2) | 普通 | 防具/首饰 |
| affix_of_spirit_t3 | of 灵魂之火 | 精神+ | 3~10 (T3) | 普通 | 防具/首饰 |
| affix_of_might | of 威力 | 全属性+ | 1~2 (T1) | 稀有 | 首饰 |
| affix_of_might_t2 | of 全能 | 全属性+ | 2~4 (T2) | 稀有 | 首饰 |
| affix_of_might_t3 | of 万象归一 | 全属性+ | 3~5 (T3) | 稀有 | 首饰 |

---

#### 防御与生存相关词缀

> 📌 提升角色生存能力的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_tough | 坚韧的 | 物理防御+ | 3~8 (T1) | 普通 | 防具 |
| affix_tough_t2 | 坚硬的 | 物理防御+ | 5~14 (T2) | 普通 | 防具 |
| affix_tough_t3 | 钢铁之躯 | 物理防御+ | 8~20 (T3) | 普通 | 防具 |
| affix_resistant | 抗性的 | 魔法防御+ | 3~8 (T1) | 普通 | 防具 |
| affix_resistant_t2 | 防护的 | 魔法防御+ | 5~14 (T2) | 普通 | 防具 |
| affix_resistant_t3 | 魔法壁垒 | 魔法防御+ | 8~20 (T3) | 普通 | 防具 |
| affix_warded | 防护的 | 全元素抗性+% | 2~5% (T1) | 精良 | 防具 |
| affix_warded_t2 | 元素防护的 | 全元素抗性+% | 4~9% (T2) | 精良 | 防具 |
| affix_warded_t3 | 元素壁垒 | 全元素抗性+% | 5~13% (T3) | 精良 | 防具 |
| affix_fireproof | 防火的 | 火焰抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_fireproof_t2 | 烈焰的 | 火焰抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_fireproof_t3 | 不灭之火 | 火焰抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_frostproof | 防冻的 | 冰霜抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_frostproof_t2 | 寒冰的 | 冰霜抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_frostproof_t3 | 永恒之冰 | 冰霜抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_insulated | 绝缘的 | 雷电抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_insulated_t2 | 雷电的 | 雷电抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_insulated_t3 | 避雷之盾 | 雷电抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_shadowward | 暗影防护的 | 暗影抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_shadowward_t2 | 暗影的 | 暗影抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_shadowward_t3 | 暗影之壁 | 暗影抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_holyward | 神圣防护的 | 神圣抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_holyward_t2 | 神圣的 | 神圣抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_holyward_t3 | 圣光之盾 | 神圣抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_natureward | 自然防护的 | 自然抗性+% | 5~12% (T1) | 精良 | 防具 |
| affix_natureward_t2 | 自然的 | 自然抗性+% | 9~22% (T2) | 精良 | 防具 |
| affix_natureward_t3 | 生命之盾 | 自然抗性+% | 13~30% (T3) | 精良 | 防具 |
| affix_elementalward | 元素防护的 | 全元素抗性+% | 3~8% (T1) | 稀有 | 防具 |
| affix_elementalward_t2 | 元素护体的 | 全元素抗性+% | 5~14% (T2) | 稀有 | 防具 |
| affix_elementalward_t3 | 元素壁垒 | 全元素抗性+% | 8~20% (T3) | 稀有 | 防具 |
| affix_immortal | 不朽的 | 最大生命+% | 5~12 (T1) | 精良 | 防具 |
| affix_immortal_t2 | 永生的 | 最大生命+% | 9~22 (T2) | 精良 | 防具 |
| affix_immortal_t3 | 生命之源 | 最大生命+% | 13~30 (T3) | 精良 | 防具 |
| affix_regenerating | 再生的 | 每回合生命恢复+ | 2~5 (T1) | 稀有 | 防具 |
| affix_regenerating_t2 | 恢复的 | 每回合生命恢复+ | 4~9 (T2) | 稀有 | 防具 |
| affix_regenerating_t3 | 生命之泉 | 每回合生命恢复+ | 5~13 (T3) | 稀有 | 防具 |
| affix_undying | 不死的 | 生命值+% + 生命恢复+ | 8~15% + 3~6 (T1) | 史诗 | 防具 |
| affix_undying_t2 | 不灭的 | 生命值+% + 生命恢复+ | 14~27% + 5~11 (T2) | 史诗 | 防具 |
| affix_undying_t3 | 永恒之体 | 生命值+% + 生命恢复+ | 20~38% + 8~15 (T3) | 史诗 | 防具 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_protection | of 保护 | 物理防御+ | 2~6 (T1) | 普通 | 防具 |
| affix_of_protection_t2 | of 防护 | 物理防御+ | 4~11 (T2) | 普通 | 防具 |
| affix_of_protection_t3 | of 钢铁之躯 | 物理防御+ | 5~15 (T3) | 普通 | 防具 |
| affix_of_warding | of 防护 | 魔法防御+ | 2~6 (T1) | 普通 | 防具 |
| affix_of_warding_t2 | of 抗性 | 魔法防御+ | 4~11 (T2) | 普通 | 防具 |
| affix_of_warding_t3 | of 魔法壁垒 | 魔法防御+ | 5~15 (T3) | 普通 | 防具 |
| affix_of_resistance | of 抗性 | 全元素抗性+% | 1~4% (T1) | 精良 | 防具 |
| affix_of_resistance_t2 | of 元素防护 | 全元素抗性+% | 2~7% (T2) | 精良 | 防具 |
| affix_of_resistance_t3 | of 元素壁垒 | 全元素抗性+% | 3~10% (T3) | 精良 | 防具 |
| affix_of_fire_resistance | of 火焰抗性 | 火焰抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_fire_resistance_t2 | of 烈焰护体 | 火焰抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_fire_resistance_t3 | of 不灭之火 | 火焰抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_frost_resistance | of 冰霜抗性 | 冰霜抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_frost_resistance_t2 | of 寒冰护体 | 冰霜抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_frost_resistance_t3 | of 永恒之冰 | 冰霜抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_lightning_resistance | of 雷电抗性 | 雷电抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_lightning_resistance_t2 | of 雷电护体 | 雷电抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_lightning_resistance_t3 | of 避雷之盾 | 雷电抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_shadow_resistance | of 暗影抗性 | 暗影抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_shadow_resistance_t2 | of 暗影护体 | 暗影抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_shadow_resistance_t3 | of 暗影之壁 | 暗影抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_holy_resistance | of 神圣抗性 | 神圣抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_holy_resistance_t2 | of 神圣护体 | 神圣抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_holy_resistance_t3 | of 圣光之盾 | 神圣抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_nature_resistance | of 自然抗性 | 自然抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_nature_resistance_t2 | of 自然护体 | 自然抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_nature_resistance_t3 | of 生命之盾 | 自然抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_elemental_resistance | of 元素抗性 | 全元素抗性+% | 2~6% (T1) | 精良 | 防具 |
| affix_of_elemental_resistance_t2 | of 元素护体 | 全元素抗性+% | 4~11% (T2) | 精良 | 防具 |
| affix_of_elemental_resistance_t3 | of 元素壁垒 | 全元素抗性+% | 5~15% (T3) | 精良 | 防具 |
| affix_of_physical_resistance | of 物理抗性 | 物理抗性+% | 3~10% (T1) | 普通 | 防具 |
| affix_of_physical_resistance_t2 | of 钢铁护体 | 物理抗性+% | 5~18% (T2) | 普通 | 防具 |
| affix_of_physical_resistance_t3 | of 金刚之盾 | 物理抗性+% | 8~25% (T3) | 普通 | 防具 |
| affix_of_vitality | of 活力 | 最大生命+% | 3~8 (T1) | 精良 | 防具 |
| affix_of_vitality_t2 | of 强健 | 最大生命+% | 5~14 (T2) | 精良 | 防具 |
| affix_of_vitality_t3 | of 生命之源 | 最大生命+% | 8~20 (T3) | 精良 | 防具 |
| affix_of_health | of 健康 | 每回合生命恢复+ | 1~4 (T1) | 稀有 | 防具 |
| affix_of_health_t2 | of 恢复 | 每回合生命恢复+ | 2~7 (T2) | 稀有 | 防具 |
| affix_of_health_t3 | of 生命之泉 | 每回合生命恢复+ | 3~10 (T3) | 稀有 | 防具 |
| affix_of_immortality | of 不朽 | 生命值+% + 生命恢复+ | 5~10% + 2~4 (T1) | 史诗 | 防具 |
| affix_of_immortality_t2 | of 不灭 | 生命值+% + 生命恢复+ | 9~18% + 4~7 (T2) | 史诗 | 防具 |
| affix_of_immortality_t3 | of 永恒之体 | 生命值+% + 生命恢复+ | 13~25% + 5~10 (T3) | 史诗 | 防具 |

---

#### 物理伤害相关词缀

> 📌 提升物理伤害的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_physical | 物理的 | 物理伤害+ | 2~5 (T1) | 普通 | 武器 |
| affix_mighty_phys | 强大的 | 物理伤害+% | 8~15% (T1) | 精良 | 武器 |
| affix_mighty_phys_t2 | 狂暴的 | 物理伤害+% | 14~27% (T2) | 精良 | 武器 |
| affix_mighty_phys_t3 | 毁灭之怒 | 物理伤害+% | 20~38% (T3) | 精良 | 武器 |
| affix_brutal | 残暴的 | 物理伤害+% | 12~20% (T1) | 稀有 | 武器 |
| affix_brutal_t2 | 嗜血的 | 物理伤害+% | 22~36% (T2) | 稀有 | 武器 |
| affix_brutal_t3 | 屠戮之刃 | 物理伤害+% | 30~50% (T3) | 稀有 | 武器 |
| affix_devastating_phys | 毁灭的 | 物理伤害+% | 15~25% (T1) | 史诗 | 武器 |
| affix_devastating_phys_t2 | 天灾的 | 物理伤害+% | 27~45% (T2) | 史诗 | 武器 |
| affix_devastating_phys_t3 | 灭世之威 | 物理伤害+% | 38~63% (T3) | 史诗 | 武器 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_physical_damage | of 物理伤害 | 物理伤害+% | 3~8% (T1) | 普通 | 武器 |
| affix_of_physical_damage_t2 | of 撕裂 | 物理伤害+% | 5~14% (T2) | 普通 | 武器 |
| affix_of_physical_damage_t3 | of 粉碎 | 物理伤害+% | 8~20% (T3) | 普通 | 武器 |
| affix_of_might | of 威力 | 物理伤害+% | 5~12% (T1) | 精良 | 武器 |
| affix_of_might_t2 | of 破甲 | 物理伤害+% | 9~22% (T2) | 精良 | 武器 |
| affix_of_might_t3 | of 崩山 | 物理伤害+% | 13~30% (T3) | 精良 | 武器 |
| affix_of_brutality | of 残暴 | 物理伤害+% | 8~15% (T1) | 稀有 | 武器 |
| affix_of_brutality_t2 | of 嗜血 | 物理伤害+% | 14~27% (T2) | 稀有 | 武器 |
| affix_of_brutality_t3 | of 屠戮 | 物理伤害+% | 20~38% (T3) | 稀有 | 武器 |

**物理伤害词缀分级示例：**
```
of 物理伤害 (百分比加成):
Tier 1 (1-20级):   +3~8%
Tier 2 (21-40级):  +5~14%
Tier 3 (41-60级):  +8~20%

物理的 (固定数值加成):
Tier 1 (1-20级):   +2~5
Tier 2 (21-40级):  +4~9
Tier 3 (41-60级):  +5~13
```

---

#### 元素伤害相关词缀

> 📌 提升特定元素伤害的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_burning | 燃烧的 | 火焰伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_burning_t2 | 熔岩的 | 火焰伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_burning_t3 | 凤凰之怒 | 火焰伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_freezing | 冰冻的 | 冰霜伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_freezing_t2 | 寒冰的 | 冰霜伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_freezing_t3 | 冰龙之息 | 冰霜伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_shocking | 电击的 | 雷电伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_shocking_t2 | 闪电的 | 雷电伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_shocking_t3 | 雷神之怒 | 雷电伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_corrupting | 腐蚀的 | 暗影伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_corrupting_t2 | 暗影的 | 暗影伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_corrupting_t3 | 暗影之魂 | 暗影伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_holy_damage | 神圣的 | 神圣伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_holy_damage_t2 | 圣光的 | 神圣伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_holy_damage_t3 | 天使之翼 | 神圣伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_natural | 自然的 | 自然伤害+% | 10~20 (T1) | 精良 | 武器 |
| affix_natural_t2 | 自然之力的 | 自然伤害+% | 18~36 (T2) | 精良 | 武器 |
| affix_natural_t3 | 自然之灵 | 自然伤害+% | 25~50 (T3) | 精良 | 武器 |
| affix_elemental | 元素的 | 全元素伤害+% | 8~15 (T1) | 稀有 | 武器 |
| affix_elemental_t2 | 元素之力的 | 全元素伤害+% | 14~27 (T2) | 稀有 | 武器 |
| affix_elemental_t3 | 元素之魂 | 全元素伤害+% | 20~38 (T3) | 稀有 | 武器 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_flame | of 火焰 | 火焰伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_flame_t2 | of 熔岩 | 火焰伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_flame_t3 | of 凤凰之怒 | 火焰伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_frost | of 冰霜 | 冰霜伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_frost_t2 | of 寒冰 | 冰霜伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_frost_t3 | of 冰龙之息 | 冰霜伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_lightning | of 雷电 | 雷电伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_lightning_t2 | of 闪电 | 雷电伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_lightning_t3 | of 雷神之怒 | 雷电伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_shadow | of 暗影 | 暗影伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_shadow_t2 | of 暗影之力 | 暗影伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_shadow_t3 | of 暗影之魂 | 暗影伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_holy | of 神圣 | 神圣伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_holy_t2 | of 圣光 | 神圣伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_holy_t3 | of 天使之翼 | 神圣伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_nature | of 自然 | 自然伤害+% | 5~12% (T1) | 普通 | 武器 |
| affix_of_nature_t2 | of 自然之力 | 自然伤害+% | 9~22% (T2) | 普通 | 武器 |
| affix_of_nature_t3 | of 自然之灵 | 自然伤害+% | 13~30% (T3) | 普通 | 武器 |
| affix_of_elements | of 元素 | 全元素伤害+% | 4~8% (T1) | 稀有 | 武器 |
| affix_of_elements_t2 | of 元素之力 | 全元素伤害+% | 7~14% (T2) | 稀有 | 武器 |
| affix_of_elements_t3 | of 元素之魂 | 全元素伤害+% | 10~20% (T3) | 稀有 | 武器 |

---

#### 特殊效果词缀

> 📌 提供独特战斗效果的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_poisonous | 剧毒的 | 攻击附加中毒效果 | 每回合2~5点伤害，持续3回合 (T1) | 稀有 | 武器 |
| affix_poisonous_t2 | 剧毒的 | 攻击附加中毒效果 | 每回合4~9点伤害，持续3回合 (T2) | 稀有 | 武器 |
| affix_poisonous_t3 | 毒龙之息 | 攻击附加中毒效果 | 每回合5~13点伤害，持续3回合 (T3) | 稀有 | 武器 |
| affix_cursed | 诅咒的 | 攻击附加虚弱效果 | 降低目标5~10%攻击力，持续3回合 (T1) | 稀有 | 武器 |
| affix_cursed_t2 | 虚弱的 | 攻击附加虚弱效果 | 降低目标9~18%攻击力，持续3回合 (T2) | 稀有 | 武器 |
| affix_cursed_t3 | 诅咒之触 | 攻击附加虚弱效果 | 降低目标13~25%攻击力，持续3回合 (T3) | 稀有 | 武器 |
| affix_blessed_heal | 祝福的 | 攻击附加治疗 | 造成伤害的5~10%转化为生命 (T1) | 稀有 | 武器 |
| affix_blessed_heal_t2 | 治愈的 | 攻击附加治疗 | 造成伤害的9~18%转化为生命 (T2) | 稀有 | 武器 |
| affix_blessed_heal_t3 | 生命之触 | 攻击附加治疗 | 造成伤害的13~25%转化为生命 (T3) | 稀有 | 武器 |
| affix_echoing | 回响的 | 资源获取+% | 10~18 (T1) | 稀有 | 武器/首饰 |
| affix_echoing_t2 | 共鸣的 | 资源获取+% | 18~32 (T2) | 稀有 | 武器/首饰 |
| affix_echoing_t3 | 回响之音 | 资源获取+% | 25~45 (T3) | 稀有 | 武器/首饰 |
| affix_empowering | 强化的 | 技能伤害+% | 8~15 (T1) | 稀有 | 武器 |
| affix_empowering_t2 | 增强的 | 技能伤害+% | 14~27 (T2) | 稀有 | 武器 |
| affix_empowering_t3 | 威力之触 | 技能伤害+% | 20~38 (T3) | 稀有 | 武器 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_poison | of 剧毒 | 攻击附加中毒 | 每回合1~3点伤害，持续2回合 (T1) | 精良 | 武器 |
| affix_of_poison_t2 | of 剧毒 | 攻击附加中毒 | 每回合2~5点伤害，持续2回合 (T2) | 精良 | 武器 |
| affix_of_poison_t3 | of 毒龙之息 | 攻击附加中毒 | 每回合3~8点伤害，持续2回合 (T3) | 精良 | 武器 |
| affix_of_weakness | of 虚弱 | 攻击附加虚弱 | 降低目标3~6%攻击力，持续2回合 (T1) | 精良 | 武器 |
| affix_of_weakness_t2 | of 虚弱 | 攻击附加虚弱 | 降低目标5~11%攻击力，持续2回合 (T2) | 精良 | 武器 |
| affix_of_weakness_t3 | of 诅咒之触 | 攻击附加虚弱 | 降低目标8~15%攻击力，持续2回合 (T3) | 精良 | 武器 |
| affix_of_healing | of 治疗 | 攻击附加治疗 | 造成伤害的3~6%转化为生命 (T1) | 精良 | 武器 |
| affix_of_healing_t2 | of 治愈 | 攻击附加治疗 | 造成伤害的5~11%转化为生命 (T2) | 精良 | 武器 |
| affix_of_healing_t3 | of 生命之触 | 攻击附加治疗 | 造成伤害的8~15%转化为生命 (T3) | 精良 | 武器 |
| affix_of_energy | of 能量 | 资源获取+% | 8~15 (T1) | 稀有 | 武器/首饰 |
| affix_of_energy_t2 | of 充能 | 资源获取+% | 14~27 (T2) | 稀有 | 武器/首饰 |
| affix_of_energy_t3 | of 回响之音 | 资源获取+% | 20~38 (T3) | 稀有 | 武器/首饰 |
| affix_of_empowerment | of 威力 | 技能伤害+% | 5~10 (T1) | 稀有 | 武器 |
| affix_of_empowerment_t2 | of 增强 | 技能伤害+% | 9~18 (T2) | 稀有 | 武器 |
| affix_of_empowerment_t3 | of 威力之触 | 技能伤害+% | 13~25 (T3) | 稀有 | 武器 |

---

#### 资源管理词缀

> 📌 提升资源获取和管理的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_energizing | 充能的 | 法力恢复+% | 10~20 (T1) | 精良 | 防具/首饰 |
| affix_energizing_t2 | 充盈的 | 法力恢复+% | 18~36 (T2) | 精良 | 防具/首饰 |
| affix_energizing_t3 | 奥术之源 | 法力恢复+% | 25~50 (T3) | 精良 | 防具/首饰 |
| affix_raging | 狂暴的 | 怒气获取+% | 15~25 (T1) | 精良 | 武器 |
| affix_raging_t2 | 狂怒的 | 怒气获取+% | 27~45 (T2) | 精良 | 武器 |
| affix_raging_t3 | 狂暴之怒 | 怒气获取+% | 38~63 (T3) | 精良 | 武器 |
| affix_energetic | 充满能量的 | 能量恢复+% | 15~25 (T1) | 精良 | 武器 |
| affix_energetic_t2 | 充盈的 | 能量恢复+% | 27~45 (T2) | 精良 | 武器 |
| affix_energetic_t3 | 能量之源 | 能量恢复+% | 38~63 (T3) | 精良 | 武器 |
| affix_manabound | 法力束缚的 | 最大法力+% | 8~15 (T1) | 精良 | 防具/首饰 |
| affix_manabound_t2 | 法力充盈的 | 最大法力+% | 14~27 (T2) | 精良 | 防具/首饰 |
| affix_manabound_t3 | 奥术之海 | 最大法力+% | 20~38 (T3) | 精良 | 防具/首饰 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_mana | of 法力 | 法力恢复+% | 5~12 (T1) | 普通 | 防具/首饰 |
| affix_of_mana_t2 | of 充盈 | 法力恢复+% | 9~22 (T2) | 普通 | 防具/首饰 |
| affix_of_mana_t3 | of 奥术之源 | 法力恢复+% | 13~30 (T3) | 普通 | 防具/首饰 |
| affix_of_rage | of 怒气 | 怒气获取+% | 8~15 (T1) | 精良 | 武器 |
| affix_of_rage_t2 | of 狂怒 | 怒气获取+% | 14~27 (T2) | 精良 | 武器 |
| affix_of_rage_t3 | of 狂暴之怒 | 怒气获取+% | 20~38 (T3) | 精良 | 武器 |
| affix_of_energy | of 能量 | 能量恢复+% | 8~15 (T1) | 精良 | 武器 |
| affix_of_energy_t2 | of 充盈 | 能量恢复+% | 14~27 (T2) | 精良 | 武器 |
| affix_of_energy_t3 | of 能量之源 | 能量恢复+% | 20~38 (T3) | 精良 | 武器 |
| affix_of_capacity | of 容量 | 最大法力+% | 5~10 (T1) | 精良 | 防具/首饰 |
| affix_of_capacity_t2 | of 充盈 | 最大法力+% | 9~18 (T2) | 精良 | 防具/首饰 |
| affix_of_capacity_t3 | of 奥术之海 | 最大法力+% | 13~25 (T3) | 精良 | 防具/首饰 |

---

#### 战斗效率词缀

> 📌 提升战斗整体效率的词缀

**前缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_efficient | 高效的 | 战斗时间-% | 5~10 (T1) | 稀有 | 首饰 |
| affix_efficient_t2 | 极速的 | 战斗时间-% | 9~18 (T2) | 稀有 | 首饰 |
| affix_efficient_t3 | 时间之流 | 战斗时间-% | 13~25 (T3) | 稀有 | 首饰 |
| affix_swift | 迅速的 | 回合间隔时间-% | 3~8 (T1) | 稀有 | 首饰 |
| affix_swift_t2 | 疾速的 | 回合间隔时间-% | 5~14 (T2) | 稀有 | 首饰 |
| affix_swift_t3 | 时间加速 | 回合间隔时间-% | 8~20 (T3) | 稀有 | 首饰 |
| affix_energetic | 精力充沛的 | 战斗后恢复速度+% | 20~40 (T1) | 精良 | 首饰 |
| affix_energetic_t2 | 充满活力的 | 战斗后恢复速度+% | 36~72 (T2) | 精良 | 首饰 |
| affix_energetic_t3 | 永动之体 | 战斗后恢复速度+% | 50~100 (T3) | 精良 | 首饰 |

**后缀:**

| ID | 名称 | 效果 | 数值范围 | 稀有度 | 适用 |
|---|-----|------|---------|-------|------|
| affix_of_efficiency | of 效率 | 战斗时间-% | 3~6 (T1) | 精良 | 首饰 |
| affix_of_efficiency_t2 | of 极速 | 战斗时间-% | 5~11 (T2) | 精良 | 首饰 |
| affix_of_efficiency_t3 | of 时间之流 | 战斗时间-% | 8~15 (T3) | 精良 | 首饰 |
| affix_of_turn_speed | of 迅速 | 回合间隔时间-% | 2~5 (T1) | 精良 | 首饰 |
| affix_of_turn_speed_t2 | of 疾速 | 回合间隔时间-% | 4~9 (T2) | 精良 | 首饰 |
| affix_of_turn_speed_t3 | of 时间加速 | 回合间隔时间-% | 5~13 (T3) | 精良 | 首饰 |
| affix_of_recovery | of 恢复 | 战斗后恢复速度+% | 15~30 (T1) | 精良 | 首饰 |
| affix_of_recovery_t2 | of 活力 | 战斗后恢复速度+% | 27~54 (T2) | 精良 | 首饰 |
| affix_of_recovery_t3 | of 永动之体 | 战斗后恢复速度+% | 38~75 (T3) | 精良 | 首饰 |

---

#### 词缀分类总结

| 词缀类别 | 前缀数量 | 后缀数量 | 主要适用装备 |
|---------|---------|---------|------------|
| 攻击/属性向 | 8 | 10 | 武器/防具 |
| 掉落与收益 | 6 | 6 | 首饰 |
| 效率与速度 | 5 | 5 | 武器/防具/靴子 |
| 属性增强 | 6 | 6 | 全部 |
| 防御与生存 | 14 | 15 | 防具 |
| 物理伤害 | 4 | 3 | 武器 |
| 元素伤害 | 7 | 7 | 武器 |
| 特殊效果 | 5 | 5 | 武器/首饰 |
| 资源管理 | 4 | 4 | 武器/防具/首饰 |
| 战斗效率 | 3 | 3 | 首饰 |
| **总计** | **57** | **64** | - |

---

## 材料强化系统

> 📌 **设计理念**: 类似POE的材料系统，通过消耗各种材料来强化和改造装备，让每件装备都有培养价值

### 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          材料强化系统架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   装备强化操作:                                                              │
│                                                                             │
│   1. 重铸 (Reforge)                                                         │
│      └─ 消耗: 重铸石 × 1                                                    │
│      └─ 效果: 随机重新生成所有词缀（保留品质）                               │
│                                                                             │
│   2. 添加词缀 (Add Affix)                                                   │
│      └─ 消耗: 词缀石 × 1                                                    │
│      └─ 效果: 随机添加一个新词缀（如果槽位未满）                              │
│                                                                             │
│   3. 精华重铸 (Essence Reforge)                                             │
│      └─ 消耗: 精华 × 1                                                      │
│      └─ 效果: 重铸装备，但保证出现精华指定的词缀                              │
│                                                                             │
│   4. 词缀强化 (Affix Enhancement)                                           │
│      └─ 消耗: 催化剂 × 1                                                    │
│      └─ 效果: 提升指定词缀的数值（+10%~30%）                                  │
│                                                                             │
│   5. 锁定重铸 (Locked Reforge)                                              │
│      └─ 消耗: 重铸石 × 1 + 锁定石 × N（N=锁定词缀数）                        │
│      └─ 效果: 重铸未锁定的词缀，保留锁定的词缀                                │
│                                                                             │
│   注意: 传说装备只能使用操作4（词缀强化），其他操作均不可用                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 材料分类

#### 1. 基础材料

| 材料ID | 名称 | 图标 | 用途 | 获取方式 |
|-------|------|------|------|---------|
| material_reforge | 重铸石 | 🔧 | 重铸装备词缀 | 分解装备、商店购买 |
| material_affix | 词缀石 | 📝 | 添加新词缀 | 分解装备、商店购买 |
| material_lock | 锁定石 | 🔒 | 锁定词缀 | 分解高品质装备、商店购买 |

**重铸石 (Reforge Stone)**
- 效果: 随机重新生成装备的所有词缀
- 保留: 装备品质、基础属性
- 消耗: 1个/次
- 风险: 可能生成更差的词缀组合

**词缀石 (Affix Stone)**
- 效果: 如果装备词缀槽位未满，随机添加一个新词缀
- 限制: 白色装备(0词缀) → 最多1个词缀
- 限制: 绿色装备(1词缀) → 最多2个词缀
- 限制: 蓝色装备(2词缀) → 最多3个词缀
- 限制: 紫色装备(3词缀) → 最多4个词缀
- 限制: 橙色装备(4词缀) → 无法添加（已满）
- 限制: 传说装备 → 无法添加（固定词缀，不可改变）

**锁定石 (Lock Stone)**
- 效果: 锁定一个词缀，使其在重铸时不会被改变
- 消耗: 锁定1个词缀需要1个锁定石
- 限制: 最多可锁定2个词缀
- 持续时间: 锁定后持续到手动解锁或装备被分解

---

#### 2. 精华材料

> 📌 精华可以保证装备出现特定类型的词缀，是定向强化的核心材料

| 精华ID | 名称 | 颜色 | 保证词缀类型 | 稀有度 |
|-------|------|------|------------|-------|
| essence_fire | 火焰精华 | 🔥 | 火焰伤害相关词缀 | 普通 |
| essence_frost | 冰霜精华 | ❄️ | 冰霜伤害相关词缀 | 普通 |
| essence_lightning | 雷电精华 | ⚡ | 雷电伤害相关词缀 | 普通 |
| essence_physical | 物理精华 | ⚔️ | 物理攻击相关词缀 | 普通 |
| essence_crit | 暴击精华 | 💥 | 暴击率/暴击伤害词缀 | 精良 |
| essence_life | 生命精华 | ❤️ | 生命值/生命恢复词缀 | 普通 |
| essence_defense | 防御精华 | 🛡️ | 防御/护甲相关词缀 | 普通 |
| essence_speed | 速度精华 | 💨 | 资源获取/先手值词缀 | 精良 |
| essence_leech | 吸血精华 | 🩸 | 生命偷取/伤害转化词缀 | 稀有 |

**精华使用规则:**
- 消耗1个精华 + 1个重铸石
- 重铸装备，但保证出现精华指定的词缀类型
- 如果装备已有该类型词缀，会替换为新的（数值重新随机）
- 如果装备没有该类型词缀，会添加一个（如果槽位未满）

**精华等级:**
- 普通精华: 保证普通/精良词缀
- 精良精华: 保证精良/稀有词缀
- 稀有精华: 保证稀有/史诗词缀

---

#### 3. 催化剂

> 📌 催化剂用于强化已有词缀的数值，是后期优化的关键材料

| 催化剂ID | 名称 | 图标 | 效果 | 适用词缀类型 |
|---------|------|------|------|------------|
| catalyst_attack | 攻击催化剂 | ⚔️ | 词缀数值+15% | 攻击力、伤害相关 |
| catalyst_defense | 防御催化剂 | 🛡️ | 词缀数值+15% | 防御、护甲相关 |
| catalyst_attr | 属性催化剂 | 💪 | 词缀数值+20% | 力量/敏捷/智力等 |
| catalyst_crit | 暴击催化剂 | 💥 | 词缀数值+20% | 暴击率、暴击伤害 |
| catalyst_life | 生命催化剂 | ❤️ | 词缀数值+15% | 生命值、生命恢复 |
| catalyst_speed | 速度催化剂 | 💨 | 词缀数值+20% | 资源获取、先手值 |

**催化剂使用规则:**
- 选择一个已有词缀进行强化
- 消耗1个对应类型的催化剂
- 词缀数值提升10%~30%（随机，平均15%）
- 可以多次强化同一词缀（每次消耗催化剂）
- 强化上限: 词缀数值最多提升到原值的200%

**强化示例:**
```
原始词缀: of 暴击 (+5% 暴击率)
使用暴击催化剂 × 1
强化后: of 暴击 (+6% 暴击率) [提升20%]

再次使用暴击催化剂 × 1
强化后: of 暴击 (+7% 暴击率) [累计提升40%]
```

---

#### 4. 保护材料

| 材料ID | 名称 | 图标 | 用途 |
|-------|------|------|------|
| material_protect | 保护石 | 🛡️ | 保护所有词缀不被重铸改变 |

**保护石 (Protect Stone)**
- 效果: 使用重铸石时，保护所有词缀不被改变
- 消耗: 1个重铸石 + 1个保护石
- 用途: 只改变装备的基础属性（如果可改变）或品质
- 限制: 不能与锁定石同时使用

---

### 材料强化操作流程

#### 操作1: 基础重铸

```
步骤:
1. 选择装备（传说装备除外）
2. 消耗: 重铸石 × 1
3. 随机重新生成所有词缀
4. 保留装备品质和基础属性

风险: ⚠️ 可能生成更差的词缀组合
限制: ❌ 传说装备无法使用此操作（词缀固定）
```

#### 操作2: 添加词缀

```
步骤:
1. 选择装备（词缀槽位未满，传说装备除外）
2. 消耗: 词缀石 × 1
3. 随机添加一个新词缀
4. 如果添加后达到品质上限，可能提升品质

示例:
白色装备(0词缀) + 词缀石 → 绿色装备(1词缀)
绿色装备(1词缀) + 词缀石 → 蓝色装备(2词缀)

限制: ❌ 传说装备无法使用此操作（词缀固定）
```

#### 操作3: 精华重铸

```
步骤:
1. 选择装备（传说装备除外）
2. 消耗: 精华 × 1 + 重铸石 × 1
3. 重铸所有词缀，但保证出现精华指定的词缀类型
4. 如果装备已有该类型词缀，替换为新词缀

示例:
装备: 炽热的钢剑 of 力量 (+2火伤, +1力量)
使用: 暴击精华 + 重铸石
结果: 炽热的钢剑 of 暴击 (+2火伤, +4%暴击率)
      ↑保留火焰词缀    ↑保证暴击词缀

限制: ❌ 传说装备无法使用此操作（词缀固定）
```

#### 操作4: 词缀强化

```
步骤:
1. 选择装备和要强化的词缀
2. 消耗: 对应类型的催化剂 × 1
3. 词缀数值提升10%~30%
4. 可以多次强化（每次消耗催化剂）

限制:
- 词缀数值最多提升到原值的200%
- 每次强化消耗1个催化剂
```

#### 操作5: 锁定重铸

```
步骤:
1. 选择装备（传说装备除外）
2. 锁定要保留的词缀（消耗锁定石 × N，N=锁定词缀数）
3. 消耗: 重铸石 × 1
4. 只重铸未锁定的词缀

示例:
装备: 炽热的钢剑 of 暴击 of 力量
锁定: of 暴击 (消耗锁定石 × 1)
重铸: 只改变"炽热的"前缀
结果: 锋利的钢剑 of 暴击 of 力量
      ↑已改变    ↑保留    ↑保留

限制: ❌ 传说装备无法使用此操作（词缀固定）
```

---

### 材料获取方式

| 材料类型 | 主要获取方式 | 次要获取方式 |
|---------|------------|------------|
| 重铸石 | 分解白色/绿色装备 | 商店购买(10金/个) |
| 词缀石 | 分解蓝色装备 | 商店购买(50金/个) |
| 锁定石 | 分解紫色/橙色装备 | 商店购买(200金/个) |
| 精华 | 分解对应属性装备、怪物掉落 | 商店购买(100-500金/个) |
| 催化剂 | 分解高品质装备、Boss掉落 | 商店购买(300-1000金/个) |
| 保护石 | 分解橙色装备、Boss掉落 | 商店购买(500金/个) |

**分解规则:**
- 白色装备 → 重铸石 × 1
- 绿色装备 → 重铸石 × 1-2
- 蓝色装备 → 词缀石 × 1 + 重铸石 × 1
- 紫色装备 → 锁定石 × 1 + 词缀石 × 1-2
- 橙色装备 → 锁定石 × 1-2 + 保护石 × 1 + 随机催化剂 × 1
- 传说装备 → 锁定石 × 2-3 + 保护石 × 2 + 随机催化剂 × 2-3 + 传说碎片 × 1

---

### 数据库设计

#### crafting_materials - 材料配置表

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | VARCHAR(32) | PRIMARY KEY | 材料ID |
| name | VARCHAR(64) | NOT NULL | 材料名称 |
| type | VARCHAR(16) | NOT NULL | 类型: reforge/affix/essence/catalyst/protect/lock |
| category | VARCHAR(32) | | 分类（精华/催化剂的子类型） |
| rarity | VARCHAR(16) | DEFAULT 'common' | 稀有度 |
| effect_type | VARCHAR(32) | NOT NULL | 效果类型 |
| effect_value | REAL | | 效果数值 |
| description | TEXT | | 描述 |
| sell_price | INTEGER | DEFAULT 0 | 售价 |

```sql
CREATE TABLE crafting_materials (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,
    category VARCHAR(32),
    rarity VARCHAR(16) DEFAULT 'common',
    effect_type VARCHAR(32) NOT NULL,
    effect_value REAL,
    description TEXT,
    sell_price INTEGER DEFAULT 0
);

CREATE INDEX idx_materials_type ON crafting_materials(type);
CREATE INDEX idx_materials_rarity ON crafting_materials(rarity);
```

---

#### equipment_crafting_log - 装备强化记录表

> 📌 记录装备的强化历史，用于分析和追踪

| 字段 | 类型 | 约束 | 说明 |
|-----|------|-----|------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| equipment_instance_id | INTEGER | NOT NULL FK | 装备实例ID |
| operation_type | VARCHAR(32) | NOT NULL | 操作类型 |
| material_id | VARCHAR(32) | FK | 使用的材料ID |
| material_count | INTEGER | DEFAULT 1 | 材料数量 |
| before_state | TEXT | | 操作前状态(JSON) |
| after_state | TEXT | | 操作后状态(JSON) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | 操作时间 |

```sql
CREATE TABLE equipment_crafting_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    equipment_instance_id INTEGER NOT NULL,
    operation_type VARCHAR(32) NOT NULL,
    material_id VARCHAR(32),
    material_count INTEGER DEFAULT 1,
    before_state TEXT,
    after_state TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (equipment_instance_id) REFERENCES equipment_instance(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES crafting_materials(id)
);

CREATE INDEX idx_crafting_log_equipment ON equipment_crafting_log(equipment_instance_id);
CREATE INDEX idx_crafting_log_time ON equipment_crafting_log(created_at DESC);
```

---

### 强化策略示例

#### 策略1: 从白色装备开始培养

```
1. 获得白色装备（0词缀）
2. 使用词缀石 × 1 → 绿色装备（1词缀）
3. 使用词缀石 × 1 → 蓝色装备（2词缀）
4. 使用精华重铸 → 保证出现想要的词缀类型
5. 使用催化剂强化 → 提升词缀数值

注意: 传说装备无法使用此策略（词缀固定，只能使用催化剂强化）
```

#### 策略2: 优化已有装备

```
1. 获得蓝色装备（已有2个词缀，其中1个很好）
2. 使用锁定石 × 1 → 锁定好词缀
3. 使用重铸石 × 1 → 只改变未锁定的词缀
4. 重复步骤2-3，直到满意
5. 使用词缀石 → 添加第3个词缀
6. 使用催化剂 → 强化所有词缀
```

#### 策略3: 定向打造

```
1. 获得任意品质装备
2. 使用精华重铸 → 保证出现核心词缀（如暴击）
3. 使用锁定石锁定核心词缀
4. 使用词缀石添加其他词缀
5. 使用催化剂强化所有词缀
```

---

### 设计亮点

1. **材料多样性**: 不同类型的材料提供不同的强化路径
2. **风险与收益**: 重铸有风险，但锁定和保护机制降低风险
3. **定向强化**: 精华系统允许玩家定向打造装备
4. **渐进优化**: 催化剂系统允许逐步提升装备
5. **策略深度**: 锁定+重铸的组合提供丰富的策略选择
6. **材料循环**: 分解装备获得材料，形成经济循环

---

## 装备掉落系统

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

---

### 基础掉落率 (每场战斗)

| 怪物类型 | 掉落几率 | 说明 |
|---------|---------|------|
| 普通怪 | 5% | 基础掉率，大部分战斗不掉装备 |
| 精英怪 | 15% | 精英更值得挑战 |
| Boss | 50% | Boss有一半几率掉装备 |
| 深渊Boss | 100% | 每10层Boss必掉 |

---

### 品质分布 (当装备掉落时) - 暗黑2风格

| 品质 | 颜色 | 词缀数 | 掉落率 | 说明 |
|-----|------|-------|-------|------|
| 普通 (Common) | ⬜ 白 | 0 | 30% | 大幅减少垃圾 |
| 优秀 (Uncommon) | 🟩 绿 | 1 | 35% | 单词缀起步装 |
| 精良 (Rare) | 🟦 蓝 | 2 | 25% | 主力装备 |
| 稀有 (Epic) | 🟪 紫 | 3 | 8% | 有培养价值 |
| 史诗 (Legendary) | 🟧 橙 | 4 | 1.8% | 稀有可期待 |
| 传说 (Unique) | 🟨 金 | 固定词缀+特效 | 0.2% | 独特装备（类似暗金装备） |

### 稀有装备掉落机制（暗黑2风格）

#### 掉落率设计原则

1. **高品质装备稀有**: 橙色和独特装备掉落率极低，增加获取难度和成就感
2. **任何区域都有惊喜**: 低级区域也有极小概率掉落高级装备（奇迹掉落）
3. **Boss必掉机制**: Boss有更高概率掉落高品质装备
4. **保底机制**: 长时间无掉落时，提高掉落率

#### 怪物掉落配置

| 怪物类型 | 白色 | 绿色 | 蓝色 | 紫色 | 橙色 | 独特 |
|---------|------|------|------|------|------|------|
| **普通怪物** | 80% | 18% | 2% | 0% | 0% | 0% |
| **精英怪物** | 40% | 35% | 20% | 4% | 0.9% | 0.1% |
| **Boss怪物** | 10% | 25% | 35% | 20% | 8% | 2% |
| **特殊Boss** | 0% | 10% | 30% | 40% | 15% | 5% |

#### 区域掉落倍率

| 区域等级 | 掉落倍率 | 说明 |
|---------|---------|------|
| 1-10级 | 1.0 | 基础倍率 |
| 11-20级 | 1.2 | 略微提升 |
| 21-30级 | 1.5 | 明显提升 |
| 31-40级 | 2.0 | 大幅提升 |
| 41-50级 | 2.5 | 顶级区域 |
| 51-60级 | 3.0 | 最高倍率 |

#### 奇迹掉落系统

> 📌 **核心**: 任何区域都有极小概率掉落超出等级的顶级装备（类似暗黑2的"惊喜掉落"）

**机制**:
- 所有区域都有0.01%概率掉落橙色装备
- 所有区域都有0.001%概率掉落独特装备
- 即使1级区域，也有机会获得60级装备

**设计目的**:
- 保持低级区域的期待感
- 增加游戏的随机性和惊喜
- 让每次战斗都有期待

---

### 奇迹掉落系统

> 📌 **核心**: 任何区域都有极小概率掉落超出等级的顶级装备

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          奇迹掉落                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   场景: 玩家在1级新手村刷怪                                                   │
│                                                                             │
│   正常掉落: 只能掉落1-5级的装备（普通级底材）                                 │
│                                                                             │
│   奇迹触发 (0.5%): 无视等级限制，从全装备池随机                                │
│                    → 可能掉落60级史诗装备或传说装备！                          │
│                    → 可能掉落精英级底材（即使等级很低）！                      │
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

---

### 保底机制

> 📌 防止长时间无掉落的挫败感

| 连续无掉落次数 | 效果 |
|--------------|------|
| 1-19 | 正常掉率 |
| 20-29 | 掉率×2 |
| 30-39 | 掉率×4 |
| 40+ | 保底掉落 (🟦精良或以上) |

---

### drop_config - 掉落配置表

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

---

### user_drop_pity - 玩家保底计数表

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

## 装备属性加成

> 📌 **说明**: 装备属性加成的详细计算公式和伤害计算规则已移至 [角色属性系统设计文档](./character_attributes_design.md)

装备加成会直接添加到角色的各项属性中，具体计算公式请参考角色属性设计文档。

---

## 装备数值规范

### 装备底材等级系统

> 📌 **设计理念**: 参考暗黑2的底材系统，装备分为普通级、扩展级、精英级三个层次，不同层次的底材提供不同级别的基础属性加成

#### 底材等级定义

| 底材等级 | 英文 | 出现等级范围 | 基础属性倍率 | 说明 |
|---------|------|------------|------------|------|
| **普通级** | Normal | 1-20级 | 1.0x | 基础底材，属性加成最低 |
| **扩展级** | Exceptional | 21-40级 | 1.5x | 进阶底材，属性加成提升50% |
| **精英级** | Elite | 41-60级 | 2.0x | 顶级底材，属性加成提升100% |

**底材等级说明：**
- 底材等级决定装备的**基础属性加成范围**（攻击、防御、属性点等）
- 底材等级与装备等级相关，但不完全绑定（可通过奇迹掉落获得高等级底材）
- 同一品质的装备，精英级底材的基础属性是普通级的2倍
- 底材等级不影响词缀系统（词缀系统独立计算）

#### 装备品质与属性加成范围（按底材等级）

**普通级底材（Normal）：**

| 装备品质 | 属性加成范围 | 示例 |
|---------|------------|------|
| 普通(白) | +1~3 | 破旧之剑: 攻击+2 |
| 优秀(绿) | +3~6 | 民兵之剑: 攻击+4, 力量+2 |
| 精良(蓝) | +5~10 | 迪菲亚军刀: 攻击+7, 敏捷+3 |
| 史诗(紫) | +8~15 | 黑龙之牙: 攻击+12, 暴击+3% |
| 史诗(橙) | +12~20 | 泰坦之刃: 攻击+18, 资源获取+15% |
| 传说(金) | 固定词缀 | 霜之哀伤: 固定词缀+特殊效果 |

**扩展级底材（Exceptional，基础属性×1.5）：**

| 装备品质 | 属性加成范围 | 示例 |
|---------|------------|------|
| 普通(白) | +1.5~4.5 | 精制之剑: 攻击+3 |
| 优秀(绿) | +4.5~9 | 精锐之剑: 攻击+6, 力量+3 |
| 精良(蓝) | +7.5~15 | 精工军刀: 攻击+11, 敏捷+4.5 |
| 史诗(紫) | +12~22.5 | 精炼之牙: 攻击+18, 暴击+4.5% |
| 史诗(橙) | +18~30 | 精铸之刃: 攻击+27, 资源获取+22.5% |
| 传说(金) | 固定词缀 | 固定词缀+特殊效果 |

**精英级底材（Elite，基础属性×2.0）：**

| 装备品质 | 属性加成范围 | 示例 |
|---------|------------|------|
| 普通(白) | +2~6 | 大师之剑: 攻击+4 |
| 优秀(绿) | +6~12 | 大师之剑: 攻击+8, 力量+4 |
| 精良(蓝) | +10~20 | 大师军刀: 攻击+14, 敏捷+6 |
| 史诗(紫) | +16~30 | 大师之牙: 攻击+24, 暴击+6% |
| 史诗(橙) | +24~40 | 大师之刃: 攻击+36, 资源获取+30% |
| 传说(金) | 固定词缀 | 固定词缀+特殊效果 |

**底材等级对基础属性的影响：**

```
最终基础属性 = 品质基础属性范围 × 底材等级倍率

例如：
- 普通级精良武器：攻击+5~10
- 扩展级精良武器：攻击+7.5~15（×1.5）
- 精英级精良武器：攻击+10~20（×2.0）
```

**底材等级掉落规则：**

| 区域等级 | 普通级概率 | 扩展级概率 | 精英级概率 | 说明 |
|---------|----------|----------|----------|------|
| 1-20级 | 100% | 0% | 0% | 只能掉落普通级 |
| 21-30级 | 70% | 30% | 0% | 开始出现扩展级 |
| 31-40级 | 40% | 60% | 0% | 扩展级为主 |
| 41-50级 | 20% | 60% | 20% | 开始出现精英级 |
| 51-60级 | 10% | 40% | 50% | 精英级为主 |

**奇迹掉落：**
- 任何区域都有0.5%概率触发奇迹掉落
- 奇迹掉落可以无视等级限制，从全底材池随机选择
- 1级区域也可能掉落精英级底材（极低概率）

---

### 装备底材详细设计

> 📌 **设计理念**: 参考魔兽世界和暗黑2的设计，为每个装备槽位设计具体的底材名称，不同底材等级使用不同的名称，增加收集和辨识的乐趣

#### 武器底材设计

**武器分类：**
- **单手武器**：可主手装备，部分职业可双持（主手+副手各一把）
- **双手武器**：占用主手和副手两个槽位
- **远程武器**：弓、弩、枪械（仅主手）

**职业武器限制：**

| 职业 | 可装备武器类型 | 双持 | 盾牌 | 说明 |
|-----|-------------|------|------|------|
| **战士** | 单手剑、双手剑、单手斧、双手斧、单手锤、双手锤、长柄武器、拳套、弓、弩 | ✅ | ✅ | 武器大师，几乎可装备所有武器 |
| **盗贼** | 匕首、单手剑、拳套 | ✅ | ❌ | 双持专家，无法装备盾牌 |
| **法师** | 法杖、魔杖、单手剑、匕首 | ❌ | ❌ | 法系武器为主 |
| **牧师** | 法杖、魔杖、单手锤、匕首 | ❌ | ❌ | 法系武器为主 |
| **术士** | 法杖、魔杖、单手剑、匕首 | ❌ | ❌ | 法系武器为主 |
| **德鲁伊** | 法杖、单手锤、匕首、拳套 | ❌ | ❌ | 自然系武器 |
| **圣骑士** | 单手剑、双手剑、单手锤、双手锤 | ❌ | ✅ | 可装备盾牌 |
| **猎人** | 弓、弩、长柄武器、单手剑、双手剑 | ❌ | ❌ | 远程武器为主 |
| **萨满** | 单手锤、双手锤、法杖、匕首 | ❌ | ✅ | 可装备盾牌 |

**单手武器底材：**

| 武器类型 | 普通级 | 扩展级 | 精英级 | 基础攻击范围（精良品质） |
|---------|--------|--------|--------|----------------------|
| **单手剑** | 短剑 | 长剑 | 符文剑 | 5~10 / 7.5~15 / 10~20 |
| **单手斧** | 手斧 | 战斧 | 战刃 | 5~10 / 7.5~15 / 10~20 |
| **单手锤** | 小锤 | 战锤 | 重锤 | 5~10 / 7.5~15 / 10~20 |
| **匕首** | 小刀 | 短刃 | 影刃 | 4~8 / 6~12 / 8~16 |
| **拳套** | 拳套 | 利爪 | 钢爪 | 4~8 / 6~12 / 8~16 |
| **魔杖** | 法杖 | 奥术杖 | 秘法杖 | 3~6（魔法攻击）/ 4.5~9 / 6~12 |

**双手武器底材：**

| 武器类型 | 普通级 | 扩展级 | 精英级 | 基础攻击范围（精良品质） |
|---------|--------|--------|--------|----------------------|
| **双手剑** | 大剑 | 巨剑 | 双手巨剑 | 8~16 / 12~24 / 16~32 |
| **双手斧** | 大斧 | 战斧 | 巨斧 | 8~16 / 12~24 / 16~32 |
| **双手锤** | 大锤 | 重锤 | 巨锤 | 8~16 / 12~24 / 16~32 |
| **长柄武器** | 长矛 | 战戟 | 巨戟 | 7~14 / 10.5~21 / 14~28 |
| **法杖** | 木杖 | 法杖 | 奥术法杖 | 6~12（魔法攻击）/ 9~18 / 12~24 |

**远程武器底材：**

| 武器类型 | 普通级 | 扩展级 | 精英级 | 基础攻击范围（精良品质） |
|---------|--------|--------|--------|----------------------|
| **弓** | 短弓 | 长弓 | 强弓 | 6~12 / 9~18 / 12~24 |
| **弩** | 轻弩 | 重弩 | 强弩 | 7~14 / 10.5~21 / 14~28 |

**武器底材示例：**

```
单手剑系列：
- 普通级：破旧短剑、民兵短剑、精制短剑
- 扩展级：战士长剑、符文长剑、精工长剑
- 精英级：符文剑、符文之刃、符文大师之剑

双手剑系列：
- 普通级：破旧大剑、民兵大剑、精制大剑
- 扩展级：战士巨剑、符文巨剑、精工巨剑
- 精英级：双手巨剑、符文巨剑、符文大师之剑

法杖系列：
- 普通级：学徒木杖、法师木杖、精制木杖
- 扩展级：奥术法杖、秘法法杖、精工法杖
- 精英级：奥术法杖、秘法法杖、大法师法杖
```

#### 防具底材设计

**头部底材：**

| 护甲类型 | 普通级 | 扩展级 | 精英级 |
|---------|--------|--------|--------|
| **布甲** | 布帽 | 法师帽 | 大法师帽 |
| **皮甲** | 皮帽 | 猎手帽 | 影舞者帽 |
| **锁甲** | 链甲帽 | 战盔 | 战将头盔 |
| **板甲** | 铁盔 | 板甲头盔 | 板甲战盔 |

**盔甲底材：**

| 护甲类型 | 普通级 | 扩展级 | 精英级 |
|---------|--------|--------|--------|
| **布甲** | 布袍 | 法师袍 | 大法师袍 |
| **皮甲** | 皮甲 | 猎手皮甲 | 影舞者皮甲 |
| **锁甲** | 链甲 | 战甲 | 战将锁甲 |
| **板甲** | 板甲 | 重板甲 | 板甲战甲 |

**手套底材：**

| 护甲类型 | 普通级 | 扩展级 | 精英级 |
|---------|--------|--------|--------|
| **布甲** | 布手套 | 法师手套 | 大法师手套 |
| **皮甲** | 皮手套 | 猎手手套 | 影舞者手套 |
| **锁甲** | 链甲手套 | 战手套 | 战将手套 |
| **板甲** | 板甲手套 | 重板甲手套 | 板甲战手套 |

**靴子底材：**

| 护甲类型 | 普通级 | 扩展级 | 精英级 |
|---------|--------|--------|--------|
| **布甲** | 布靴 | 法师靴 | 大法师靴 |
| **皮甲** | 皮靴 | 猎手靴 | 影舞者靴 |
| **锁甲** | 链甲靴 | 战靴 | 战将靴 |
| **板甲** | 板甲靴 | 重板甲靴 | 板甲战靴 |

**盾牌底材：**

| 底材等级 | 名称 | 基础受伤减免（精良品质） | 说明 |
|---------|------|---------------------|------|
| **普通级** | 木盾 | 2~4% | 基础减伤 |
| **扩展级** | 铁盾 | 3~6% | 进阶减伤 |
| **精英级** | 符文盾 | 4~8% | 顶级减伤 |

**盾牌基础属性：**
- 盾牌提供物理防御和魔法防御（与其他防具相同）
- 盾牌额外提供基础受伤减免（按底材等级）
- 受伤减免作用于所有类型的伤害（物理和魔法）
- 受伤减免与词缀的受伤减免叠加

**防具底材示例：**

```
布甲系列（法师、术士、牧师）：
- 普通级：学徒布帽、学徒布袍、学徒布手套、学徒布靴
- 扩展级：法师帽、法师袍、法师手套、法师靴
- 精英级：大法师帽、大法师袍、大法师手套、大法师靴

板甲系列（战士、圣骑士）：
- 普通级：新兵铁盔、新兵板甲、新兵板甲手套、新兵板甲靴
- 扩展级：战士板甲头盔、战士板甲、战士板甲手套、战士板甲靴
- 精英级：板甲战盔、板甲战甲、板甲战手套、板甲战靴
```

#### 首饰底材设计

**首饰不区分底材等级，统一使用以下名称：**

| 首饰类型 | 名称 |
|---------|------|
| **戒指** | 戒指 |
| **项链** | 项链 |

首饰的基础属性加成由品质决定，不受底材等级影响。

---

#### 底材命名规则

**命名格式：**
```
[前缀] + [底材名称] + [后缀]

前缀：品质相关（破旧、精制、精工等）
底材名称：根据底材等级和类型
后缀：品质相关（可选）
```

**示例：**
- 普通级精良单手剑：`精制短剑`
- 扩展级史诗单手剑：`符文长剑`
- 精英级传说单手剑：`符文大师之剑`

---

### 防具护甲类型系统

> 📌 **设计理念**: 参考魔兽世界设计，防具分为布甲、皮甲、锁甲、板甲四类，不同护甲类型提供不同比例的物理防御和魔法防御，并对角色基础属性有要求

#### 护甲类型定义

| 护甲类型 | 英文 | 物理防御基础值 | 魔法防御基础值 | 属性要求 | 说明 |
|---------|------|--------------|--------------|---------|------|
| **布甲** | Cloth | 低 | 高 | 智力 ≥ 装备等级×2 | 轻便灵活，魔法防御高，物理防御低 |
| **皮甲** | Leather | 中 | 中（略高） | 敏捷 ≥ 装备等级×2 | 平衡型，物理和魔法防御均衡，略偏魔法 |
| **锁甲** | Mail | 中高 | 中低 | 力量 ≥ 装备等级×1.5 或 敏捷 ≥ 装备等级×1.5 | 中等防护，物理防御较高，略偏物理 |
| **板甲** | Plate | 高 | 低 | 力量 ≥ 装备等级×2 | 重装防护，物理防御最高，魔法防御最低 |

**属性要求说明：**
- 不满足属性要求时，装备仍可穿戴，但防御值会降低（降低50%）
- 满足属性要求时，获得完整的防御值
- 属性要求按装备等级动态计算，确保高等级装备需要更高的属性

#### 防具基础防御值（按装备等级1-60级范围）

> 📌 **设计理念**: 不同类型的防具在基础防御值上就有明显区别，布甲偏向魔法防御，板甲偏向物理防御，皮甲和锁甲在两者之间平衡，确保总体防御值平衡

**防具基础防御值（按护甲类型和槽位）：**

| 护甲类型 | 槽位 | 物理防御基础值 | 魔法防御基础值 | 总防御值 | 说明 |
|---------|------|--------------|--------------|---------|------|
| **布甲** | 头部 | 1~4 | 3~8 | 4~12 | 魔法防御为主 |
| | 盔甲 | 2~8 | 8~20 | 10~28 | 魔法防御为主 |
| | 手套 | 0.5~2 | 1.5~5 | 2~7 | 魔法防御为主 |
| | 靴子 | 0.5~2 | 1.5~5 | 2~7 | 魔法防御为主 |
| **皮甲** | 头部 | 2~6 | 2~6 | 4~12 | 平衡型，略偏魔法 |
| | 盔甲 | 4~12 | 5~15 | 9~27 | 平衡型，略偏魔法 |
| | 手套 | 1~3 | 1~3 | 2~6 | 平衡型 |
| | 靴子 | 1~3 | 1~3 | 2~6 | 平衡型 |
| **锁甲** | 头部 | 3~7 | 1.5~4 | 4.5~11 | 平衡型，略偏物理 |
| | 盔甲 | 6~18 | 3~10 | 9~28 | 平衡型，略偏物理 |
| | 手套 | 1.5~4 | 0.5~2 | 2~6 | 平衡型，略偏物理 |
| | 靴子 | 1.5~4 | 0.5~2 | 2~6 | 平衡型，略偏物理 |
| **板甲** | 头部 | 4~10 | 1~3 | 5~13 | 物理防御为主 |
| | 盔甲 | 10~25 | 2~6 | 12~31 | 物理防御为主 |
| | 手套 | 2~5 | 0.5~1.5 | 2.5~6.5 | 物理防御为主 |
| | 靴子 | 2~5 | 0.5~1.5 | 2.5~6.5 | 物理防御为主 |
| **盾牌** | 副手 | 3~15 | 2~8 | 5~23 | 不区分护甲类型，统一计算 |

**防御值平衡说明：**
- 所有护甲类型的**总防御值**（物理+魔法）在同一槽位上基本平衡
- 布甲：总防御值中，魔法防御占70-80%，物理防御占20-30%
- 皮甲：总防御值中，魔法防御占50-60%，物理防御占40-50%
- 锁甲：总防御值中，物理防御占60-70%，魔法防御占30-40%
- 板甲：总防御值中，物理防御占75-85%，魔法防御占15-25%

**防御值计算公式（按装备等级）：**

```
1. 根据护甲类型和槽位获取基础值范围:
   物理防御基础值 = 护甲类型对应槽位的物理防御基础值范围
   魔法防御基础值 = 护甲类型对应槽位的魔法防御基础值范围

2. 按等级计算实际基础值:
   物理防御基础值 = 基础值最小值 + (基础值最大值 - 基础值最小值) × (等级 / 60)
   魔法防御基础值 = 基础值最小值 + (基础值最大值 - 基础值最小值) × (等级 / 60)

3. 检查属性要求:
   IF (角色属性 < 装备属性要求):
       物理防御基础值 = 物理防御基础值 × 0.5
       魔法防御基础值 = 魔法防御基础值 × 0.5
```

**示例（30级装备，头部）：**

| 护甲类型 | 物理防御基础值 | 魔法防御基础值 | 总防御值 | 防御偏向 |
|---------|--------------|--------------|---------|---------|
| **布甲** | 1 + (4-1)×(30/60) = **2.5** | 3 + (8-3)×(30/60) = **5.5** | **8** | 魔法防御为主 |
| **皮甲** | 2 + (6-2)×(30/60) = **4** | 2 + (6-2)×(30/60) = **4** | **8** | 平衡型 |
| **锁甲** | 3 + (7-3)×(30/60) = **5** | 1.5 + (4-1.5)×(30/60) = **2.75** | **7.75** | 物理防御为主 |
| **板甲** | 4 + (10-4)×(30/60) = **7** | 1 + (3-1)×(30/60) = **2** | **9** | 物理防御为主 |

**示例（30级装备，盔甲）：**

| 护甲类型 | 物理防御基础值 | 魔法防御基础值 | 总防御值 | 防御偏向 |
|---------|--------------|--------------|---------|---------|
| **布甲** | 2 + (8-2)×(30/60) = **5** | 8 + (20-8)×(30/60) = **14** | **19** | 魔法防御为主 |
| **皮甲** | 4 + (12-4)×(30/60) = **8** | 5 + (15-5)×(30/60) = **10** | **18** | 平衡型，略偏魔法 |
| **锁甲** | 6 + (18-6)×(30/60) = **12** | 3 + (10-3)×(30/60) = **6.5** | **18.5** | 物理防御为主 |
| **板甲** | 10 + (25-10)×(30/60) = **17.5** | 2 + (6-2)×(30/60) = **4** | **21.5** | 物理防御为主 |

**属性要求示例（30级装备）：**

| 护甲类型 | 属性要求 | 说明 |
|---------|---------|------|
| 布甲 | 智力 ≥ 60 | 适合法师、术士、牧师 |
| 皮甲 | 敏捷 ≥ 60 | 适合盗贼、猎人、德鲁伊 |
| 锁甲 | 力量 ≥ 45 或 敏捷 ≥ 45 | 适合萨满、猎人（部分） |
| 板甲 | 力量 ≥ 60 | 适合战士、圣骑士、死亡骑士 |

**品质对基础防御的影响：**

| 品质 | 防御倍率 | 说明 |
|-----|---------|------|
| 普通(白) | 1.0x | 基础值 |
| 优秀(绿) | 1.1x | +10% |
| 精良(蓝) | 1.2x | +20% |
| 史诗(紫) | 1.3x | +30% |
| 史诗(橙) | 1.5x | +50% |
| 传说(金) | 固定值 | 根据装备设计固定 |

**最终防御值计算：**

```
最终物理防御 = (物理防御基础值 × 属性要求修正 × 底材等级倍率) × 品质倍率 + 词缀物理防御加成
最终魔法防御 = (魔法防御基础值 × 属性要求修正 × 底材等级倍率) × 品质倍率 + 词缀魔法防御加成

其中：
- 属性要求修正 = 1.0（满足要求）或 0.5（不满足要求）
- 底材等级倍率：普通级(1.0x) → 扩展级(1.5x) → 精英级(2.0x)
- 品质倍率：普通(1.0x) → 优秀(1.1x) → 精良(1.2x) → 史诗(1.3x) → 史诗(1.5x)
```

**底材等级对防御值的影响示例（30级盔甲，布甲，满足属性要求，精良品质）：**

| 底材等级 | 物理防御基础值 | 魔法防御基础值 | 底材倍率 | 品质倍率 | 最终物理防御 | 最终魔法防御 |
|---------|--------------|--------------|---------|---------|------------|------------|
| **普通级** | 5 | 14 | 1.0x | 1.2x | 6 | 16.8 |
| **扩展级** | 5 | 14 | 1.5x | 1.2x | 9 | 25.2 |
| **精英级** | 5 | 14 | 2.0x | 1.2x | 12 | 33.6 |

**防御值平衡验证（60级装备，盔甲槽位，满足属性要求，普通品质）：**

| 护甲类型 | 物理防御 | 魔法防御 | 总防御值 | 物理占比 | 魔法占比 |
|---------|---------|---------|---------|---------|---------|
| **布甲** | 8 | 20 | 28 | 28.6% | 71.4% |
| **皮甲** | 12 | 15 | 27 | 44.4% | 55.6% |
| **锁甲** | 18 | 10 | 28 | 64.3% | 35.7% |
| **板甲** | 25 | 6 | 31 | 80.6% | 19.4% |

**平衡说明：**
- 所有护甲类型的总防御值基本平衡（27-31之间）
- 布甲和板甲在总防御值上略有优势，但偏向不同防御类型
- 皮甲和锁甲提供平衡的防御，适合混合职业

**护甲类型选择建议：**

| 职业类型 | 推荐护甲类型 | 原因 |
|---------|------------|------|
| 法师、术士、牧师 | 布甲 | 高魔法防御，适合法系职业 |
| 盗贼、猎人、德鲁伊 | 皮甲 | 平衡防御，适合敏捷职业 |
| 萨满、部分猎人 | 锁甲 | 中等物理防御，适合混合职业 |
| 战士、圣骑士 | 板甲 | 最高物理防御，适合近战坦克 |

---

### 数值平衡检查表

每次添加新装备时，检查以下项目：

- [ ] HP加成不超过+20
- [ ] MP加成不超过+20
- [ ] 攻击加成不超过+20
- [ ] 物理防御加成不超过+20（词缀加成，基础防御值按槽位和等级计算）
- [ ] 魔法防御加成不超过+20（词缀加成，基础防御值按槽位和等级计算）
- [ ] 属性点加成不超过+5
- [ ] 百分比加成不超过+30%
- [ ] 暴击率加成不超过+8%
- [ ] 暴击伤害加成不超过+25%

---

## 背包管理

### 功能列表

| 功能 | 说明 |
|-----|------|
| **自动分解** | 可设置自动分解⬜白色装备 |
| **快速分解** | 一键分解所有未锁定的低品质装备 |
| **背包上限** | 100格，满时自动分解最低品质 |
| **掉落预览** | 战斗后显示掉落，可选拾取或直接分解 |
| **装备锁定** | 锁定重要装备防止误分解 |

---

## 装备交易系统

### 系统概览

装备交易系统是装备系统的核心组成部分，借鉴暗黑2的交易机制，促进玩家间装备流通，增加游戏社交性。

### 交易方式

#### 1. 玩家间直接交易

**功能**:
- 两个玩家面对面交易
- 实时确认交易
- 支持装备和金币交换

**流程**:
```
1. 玩家A发起交易请求
2. 玩家B接受请求
3. 双方放入交易物品
4. 双方确认交易
5. 交易完成，物品交换
```

#### 2. 拍卖行系统

**功能**:
- 玩家上架装备到拍卖行
- 其他玩家浏览和购买
- 支持价格排序和筛选

**流程**:
```
1. 玩家上架装备，设置价格
2. 支付上架费用（装备价格的5%）
3. 装备在拍卖行展示24-72小时
4. 其他玩家购买
5. 交易完成，卖家获得金币（扣除10%手续费）
```

### 交易规则

#### 1. 交易限制

- **绑定装备**: 已装备的装备不能交易（需先卸下）
- **锁定装备**: 锁定的装备不能交易
- **等级限制**: 交易装备有等级要求
- **交易冷却**: 装备获得后24小时内不能交易（防止刷装备）

#### 2. 交易费用

- **上架费用**: 装备价格的5%
- **交易手续费**: 成交价格的10%
- **最低费用**: 至少1金币

#### 3. 交易安全

- **确认机制**: 交易前需要双方确认
- **交易日志**: 记录所有交易，防止纠纷
- **反作弊**: 检测异常交易，防止刷装备

### 数据库设计

```sql
-- 拍卖行表
CREATE TABLE auction_house (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    seller_id INTEGER NOT NULL,
    equipment_instance_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    listed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    status VARCHAR(16) DEFAULT 'active',  -- active/sold/expired/cancelled
    buyer_id INTEGER,
    sold_at DATETIME,
    FOREIGN KEY (seller_id) REFERENCES users(id),
    FOREIGN KEY (equipment_instance_id) REFERENCES equipment_instance(id),
    FOREIGN KEY (buyer_id) REFERENCES users(id)
);

CREATE INDEX idx_auction_status ON auction_house(status);
CREATE INDEX idx_auction_expires ON auction_house(expires_at);

-- 交易记录表
CREATE TABLE trade_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trade_type VARCHAR(16) NOT NULL,  -- direct/auction
    seller_id INTEGER NOT NULL,
    buyer_id INTEGER,
    equipment_instance_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    fee INTEGER NOT NULL,  -- 交易手续费
    status VARCHAR(16) DEFAULT 'pending',  -- pending/completed/cancelled
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (seller_id) REFERENCES users(id),
    FOREIGN KEY (buyer_id) REFERENCES users(id),
    FOREIGN KEY (equipment_instance_id) REFERENCES equipment_instance(id)
);
```

### 交易价值

#### 1. 装备定价

- **基础价格**: 装备基础价格 × 品质倍率
- **词缀加成**: 稀有词缀增加装备价值
- **强化加成**: 强化等级增加装备价值
- **市场供需**: 根据市场供需调整价格

#### 2. 价格建议

系统可以根据装备属性自动建议价格：

```
建议价格 = 基础价格 × 品质倍率 × 词缀倍率 × 强化倍率

品质倍率:
├─ 蓝色: 1.5
├─ 紫色: 2.5
├─ 橙色: 5.0
└─ 独特: 10.0

词缀倍率:
├─ 普通词缀: 1.0
├─ 稀有词缀: 1.2
└─ 传说词缀: 1.5

强化倍率:
├─ 未强化: 1.0
├─ +1: 1.1
├─ +2: 1.2
├─ +3: 1.3
├─ +4: 1.4
└─ +5: 1.5
```

---

## 总结

### 设计亮点

1. **词缀系统**: 每件装备都有独特的词缀组合，增加收集乐趣
2. **材料强化系统**: 类似POE的材料系统，提供丰富的装备培养路径
3. **奇迹掉落**: 低级区域也有机会获得顶级装备，保持惊喜感
4. **保底机制**: 防止长时间无掉落的挫败感
5. **品质分级**: 清晰的品质体系，从白色到橙色史诗，以及独特的传说装备

### 后续扩展方向

- [ ] 装备套装系统
- [ ] 装备附魔系统
- [ ] 装备交易系统
- [ ] 装备外观系统
- [ ] 装备图鉴系统

---

**文档版本**: v1.0  
**最后更新**: 2024年

