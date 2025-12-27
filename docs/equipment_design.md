# 📦 装备系统设计文档

> 📌 **核心设计理念**: 装备掉落时随机生成词缀，通过进化系统获得更强形态和传说效果

---

## 📋 目录

1. [系统概览](#系统概览)
2. [数据库设计](#数据库设计)
3. [装备槽位定义](#装备槽位定义)
4. [词缀系统](#词缀系统)
5. [进化系统](#进化系统)
6. [强化系统](#强化系统)
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

---

## 数据库设计

### 1. items - 物品配置表

> 📌 定义所有物品的基础属性

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
| head | 头部 | 头盔、帽子 |
| shoulder | 肩部 | 肩甲 |
| chest | 胸甲 | 胸甲、衣服 |
| hands | 手套 | 手套 |
| legs | 腿部 | 护腿 |
| feet | 脚部 | 靴子 |
| main_hand | 主手武器 | 单手/双手武器 |
| off_hand | 副手 | 盾牌、副手武器 |
| neck | 项链 | 项链 |
| ring | 戒指 | 戒指（可装备2个） |
| trinket | 饰品 | 饰品 |

---

## 词缀系统

### affixes - 词缀配置表

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

---

### 词缀列表

#### 前缀 (攻击/属性向)

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

#### 后缀 (特殊效果向)

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

#### 暴击相关词缀

| 词缀 | 效果 | 适用装备 |
|-----|------|---------|
| of 暴击 | +3-8% 物理暴击率 | 武器 |
| of 致命 | +10-25% 暴击伤害 | 武器 |
| of 法术暴击 | +3-8% 法术暴击率 | 武器/饰品 |
| of 法术致命 | +10-25% 法术暴击伤害 | 武器/饰品 |

#### 闪避相关词缀

| 词缀 | 效果 | 适用装备 |
|-----|------|---------|
| of 敏捷 | +5-15 敏捷（间接增加闪避） | 护甲/饰品 |
| of 闪避 | +2-5% 闪避率 | 护甲/饰品 |
| of 灵巧 | +3-8% 闪避率 | 皮甲 |

#### 仇恨相关词缀

**前缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| 守护者的 | 仇恨生成+20%, 嘲讽CD-1 | 坦克武器 | 稀有 |
| 威压的 | 仇恨生成+15% | 坦克装备 | 精良 |
| 隐秘的 | 仇恨生成-20% | 输出武器 | 精良 |
| 暗影的 | 暴击仇恨-30% | 输出装备 | 稀有 |

**后缀:**

| 词缀 | 效果 | 适用 | 稀有度 |
|-----|------|-----|-------|
| of 威胁 | 仇恨生成+15% | 坦克装备 | 精良 |
| of 守护 | 仇恨生成+10%, 格挡+5% | 盾牌 | 精良 |
| of 隐匿 | 仇恨生成-15% | 输出装备 | 精良 |
| of 消散 | 仇恨衰减+20% | 治疗装备 | 稀有 |

---

## 进化系统

### evolution_paths - 进化路线配置表

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

---

### 武器进化路线

| 路线 | 元素 | 核心加成 | 特殊效果 |
|-----|------|---------|---------|
| 🔥 烈焰 | fire | 火焰伤害+50% | 攻击灼烧敌人 |
| ❄️ 霜寒 | frost | 冰霜伤害+50% | 攻击减速敌人 |
| ⚡ 雷霆 | lightning | 攻速+20%, 雷伤+40% | 伤害连锁跳跃 |
| ✨ 神圣 | holy | 圣伤+40%, 治疗+15% | 攻击回复生命 |
| 🌑 暗影 | shadow | 暗伤+50%, 吸血+5% | 伤害转化HP |
| 🌿 自然 | nature | 自然伤+40%, 再生+20% | 持续恢复HP |
| ⚔️ 物理 | physical | 物伤+30%, 穿透+15% | 无视部分护甲 |

---

### 防具进化路线

| 路线 | 定位 | 核心加成 | 特殊效果 |
|-----|------|---------|---------|
| 🛡️ 守护 | 坦克 | 防御+30%, 生命+20% | 受伤减免 |
| 🌵 荆棘 | 反伤 | 防御+15%, 反伤+25% | 被攻击时反弹伤害 |
| 💨 迅捷 | 闪避 | 闪避+20%, 攻速+15% | 闪避后加速 |
| 💚 再生 | 续航 | 生命+15%, 回复+50% | 每回合恢复HP |

---

### legendary_effects - 传说效果表

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

---

### 传说武器效果

| ID | 名称 | 效果 | 来源路线 |
|---|-----|------|---------|
| legend_inferno | 地狱烈焰 | 攻击使敌人灼烧3回合，每回合3点火伤 | 烈焰 |
| legend_frostmourne | 霜之哀伤 | 击杀敌人后冰冻周围敌人1回合 | 霜寒 |
| legend_thunderfury | 雷霆之怒 | 20%几率触发闪电链，最多跳跃3目标 | 雷霆 |
| legend_ashbringer | 灰烬使者 | 攻击时恢复自身5%最大生命 | 神圣 |
| legend_shadowmourne | 暗影之殇 | 暴击时吸取敌人10%当前生命 | 暗影 |
| legend_earthshatter | 大地粉碎 | 攻击叠加标记，5层后引爆额外伤害 | 自然 |
| legend_gorehowl | 血吼 | 对精英和Boss伤害+50% | 物理 |

---

### 传说防具效果

| ID | 名称 | 效果 | 来源路线 |
|---|-----|------|---------|
| legend_immortal | 不灭意志 | 首次致死伤害免疫 (每场战斗1次) | 守护 |
| legend_retribution | 复仇之刺 | 反弹50%受到的物理伤害 | 荆棘 |
| legend_shadowstep | 暗影步 | 闪避成功后下次攻击必暴击 | 迅捷 |
| legend_lifesource | 生命之泉 | HP<30%时每回合恢复10% HP | 再生 |

---

## 进化阶段与强化

### 进化阶段

| 阶段 | 名称 | 角色等级 | 强化上限 | 进化材料 | 特点 |
|-----|------|---------|---------|---------|------|
| Ⅰ | 基础 | 1-15 | +5 | - | 初始形态 |
| Ⅱ | 精炼 | 16-30 | +10 | 精炼石×5 | 属性提升20% |
| Ⅲ | 进化 | 31-45 | +15 | 进化石×3 + 元素核心 | **选择分支** |
| Ⅳ | 觉醒 | 46-55 | +20 | 觉醒石×1 + 稀有材料 | 解锁特殊效果 |
| Ⅴ | 传说 | 56-60 | +25 | 传说碎片×5 | 获得传说效果 |

---

### 强化效果

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

---

### 词缀继承

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

### 品质分布 (当装备掉落时)

| 品质 | 颜色 | 词缀数 | 掉落率 | 说明 |
|-----|------|-------|-------|------|
| 普通 (Common) | ⬜ 白 | 0 | 30% | 大幅减少垃圾 |
| 优秀 (Uncommon) | 🟩 绿 | 1 | 35% | 单词缀起步装 |
| 精良 (Rare) | 🟦 蓝 | 2 | 25% | 主力装备 |
| 稀有 (Epic) | 🟪 紫 | 3 | 8% | 有培养价值 |
| 史诗 (Legendary) | 🟧 橙 | 4 | 1.8% | 稀有可期待 |
| 传说 (Mythic) | 🟥 红 | 4+传说效果 | 0.2% | 终极追求 |

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

### 在属性计算中的作用

装备加成会直接添加到角色的各项属性中：

#### 生命值 (HP)

```
最大HP = 职业基础HP + (耐力 × 10) + (等级 × HP成长) + 装备加成
```

#### 法力值 (MP)

```
最大MP = 职业基础MP + (智力 × 5) + (等级 × MP成长) + 装备加成
```

#### 物理攻击力

```
物理攻击 = 基础攻击 + (力量 × 2.0) + (敏捷 × 0.5) + 武器伤害 + 装备加成
```

#### 法术攻击力

```
法术攻击 = 基础法伤 + (智力 × 1.5) + 装备加成
```

#### 护甲值

```
护甲值 = 职业基础护甲 + (耐力 × 0.5) + (敏捷 × 0.3) + 装备护甲
```

#### 暴击率

```
物理暴击率 = 5% + (敏捷 ÷ 20)% + 装备加成% + 被动技能加成 + Buff加成
法术暴击率 = 5% + (精神/20)% + 装备加成% + 被动技能加成 + Buff加成
```

#### 暴击伤害

```
物理暴击伤害 = 150% + (力量 ÷ 100 × 10)% + 装备加成%
法术暴击伤害 = 150% + 装备加成%
```

#### 闪避率

```
闪避率 = 5% + (敏捷 ÷ 25)% + 装备加成% + 种族加成%
```

#### 命中率

```
最终命中率 = 95% - 等级差修正 + 装备加成%
```

---

## 装备数值规范

### 装备品质与属性加成范围

| 装备品质 | 属性加成范围 | 示例 |
|---------|------------|------|
| 普通(白) | +1~3 | 破旧之剑: 攻击+2 |
| 优秀(绿) | +3~6 | 民兵之剑: 攻击+4, 力量+2 |
| 精良(蓝) | +5~10 | 迪菲亚军刀: 攻击+7, 敏捷+3 |
| 史诗(紫) | +8~15 | 黑龙之牙: 攻击+12, 暴击+3% |
| 传说(橙) | +12~20 | 雷霆之怒: 攻击+18, 攻速+10% |

---

### 数值平衡检查表

每次添加新装备时，检查以下项目：

- [ ] HP加成不超过+20
- [ ] MP加成不超过+20
- [ ] 攻击加成不超过+20
- [ ] 防御加成不超过+20
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

## 总结

### 设计亮点

1. **词缀系统**: 每件装备都有独特的词缀组合，增加收集乐趣
2. **进化系统**: 装备可以长期培养，不会快速淘汰
3. **奇迹掉落**: 低级区域也有机会获得顶级装备，保持惊喜感
4. **保底机制**: 防止长时间无掉落的挫败感
5. **品质分级**: 清晰的品质体系，从白色到红色传说
6. **强化系统**: 通过强化提升装备属性，增加培养深度

### 后续扩展方向

- [ ] 装备套装系统
- [ ] 装备附魔系统
- [ ] 装备交易系统
- [ ] 装备外观系统
- [ ] 装备图鉴系统

---

**文档版本**: v1.0  
**最后更新**: 2024年

