# 怪物系统设计文档

> 📌 **核心设计理念**: 丰富的怪物类型带来不同挑战，避免一招走遍天下。所有怪物数据配置化，方便平衡性调整。

---

## 📋 目录

1. [系统概览](#系统概览)
2. [怪物分类](#怪物分类)
3. [怪物属性设计](#怪物属性设计)
4. [怪物AI系统](#怪物ai系统)
5. [怪物技能系统](#怪物技能系统)
6. [怪物掉落系统](#怪物掉落系统)
7. [配置化设计](#配置化设计)
8. [平衡性调整](#平衡性调整)

---

## 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          怪物系统架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   怪物生成 ──→ 怪物AI ──→ 战斗行为 ──→ 掉落奖励                              │
│      │           │           │           │                                 │
│   配置表读取   行为模式    技能释放    掉落表查询                             │
│                                                                             │
│   怪物分类:                                                                  │
│   ⬜ 普通怪物 (Normal)    → 基础属性，简单AI                                 │
│   🟦 精英怪物 (Elite)    → 增强属性，特殊技能                               │
│   🟧 Boss怪物 (Boss)     → 高属性，复杂技能组合                             │
│   🟨 特殊怪物 (Special)  → 独特机制，特殊挑战                               │
│                                                                             │
│   配置化设计:                                                                │
│   ├─ 所有怪物数据存储在配置表中                                             │
│   ├─ 支持热更新，无需重启服务                                                │
│   ├─ 版本化管理，支持回滚                                                    │
│   └─ 减少数据迁移复杂度                                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 怪物分类

### 1. 普通怪物 (Normal)

**特点**:
- 基础属性，无特殊加成
- 简单的AI行为（攻击、防御）
- 基础技能或普通攻击
- 掉落：白色/绿色装备

**设计目的**: 提供基础挑战，让玩家熟悉战斗机制

**示例**:
- 森林狼：高攻击，低防御，优先攻击
- 野猪：高防御，低攻击，优先防御
- 哥布林：平衡属性，随机行为

### 2. 精英怪物 (Elite)

**特点**:
- 属性增强（1.5-2倍普通怪物）
- 特殊技能（AOE、控制、治疗等）
- 更智能的AI（优先攻击低血量、使用技能）
- 掉落：蓝色/紫色装备

**设计目的**: 提供中等挑战，需要策略应对

**示例**:
- 精英狼人：高攻击+撕裂技能（持续伤害）
- 精英法师：法术攻击+护盾技能
- 精英治疗者：低攻击+治疗技能

### 3. Boss怪物 (Boss)

**特点**:
- 高属性（3-5倍普通怪物）
- 多个技能组合
- 阶段机制（血量低于阈值时改变行为）
- 掉落：紫色/橙色/独特装备

**设计目的**: 提供高难度挑战，需要完整队伍配合

**示例**:
- 森林之王：高HP+召唤小怪+范围攻击
- 暗影法师：法术攻击+控制技能+护盾
- 龙族Boss：高攻击+火焰吐息+飞行（闪避高）

### 4. 特殊怪物 (Special)

**特点**:
- 独特的机制和挑战
- 特殊的行为模式
- 特殊的奖励
- 掉落：特殊装备/材料

**设计目的**: 提供独特挑战，增加游戏趣味性

**示例**:
- 幽灵：物理攻击无效，只能用法术攻击
- 元素：免疫对应元素伤害
- 分裂怪：死亡时分裂成多个小怪

---

## 怪物属性设计

### 基础属性

| 属性 | 说明 | 计算公式 |
|-----|------|---------|
| **HP** | 生命值 | 基础HP × 等级系数 × 类型系数 |
| **物理攻击** | 物理攻击力 | 基础攻击 × 等级系数 × 类型系数 |
| **法术攻击** | 法术攻击力 | 基础攻击 × 等级系数 × 类型系数 |
| **物理防御** | 物理防御力 | 基础防御 × 等级系数 × 类型系数 |
| **法术防御** | 法术防御力 | 基础防御 × 等级系数 × 类型系数 |
| **暴击率** | 暴击概率 | 基础暴击率 + 类型加成 |
| **闪避率** | 闪避概率 | 基础闪避率 + 类型加成 |
| **速度** | 行动速度 | 决定行动顺序 |

### 类型系数

| 怪物类型 | HP系数 | 攻击系数 | 防御系数 | 说明 |
|---------|--------|---------|---------|------|
| 普通 | 1.0 | 1.0 | 1.0 | 基准值 |
| 精英 | 1.8 | 1.5 | 1.3 | 全面增强 |
| Boss | 4.0 | 2.5 | 2.0 | 大幅增强 |
| 特殊 | 变化 | 变化 | 变化 | 根据机制调整 |

### 属性配置表

```sql
-- 怪物基础配置表
CREATE TABLE monsters (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,  -- normal/elite/boss/special
    level INTEGER NOT NULL,
    
    -- 基础属性（配置化）
    base_hp INTEGER NOT NULL,
    base_physical_attack INTEGER NOT NULL,
    base_magic_attack INTEGER NOT NULL,
    base_physical_defense INTEGER NOT NULL,
    base_magic_defense INTEGER NOT NULL,
    
    -- 属性系数（可调整）
    hp_multiplier REAL DEFAULT 1.0,
    attack_multiplier REAL DEFAULT 1.0,
    defense_multiplier REAL DEFAULT 1.0,
    
    -- 特殊属性
    crit_rate REAL DEFAULT 0.05,
    crit_damage REAL DEFAULT 1.5,
    dodge_rate REAL DEFAULT 0.05,
    speed INTEGER DEFAULT 10,
    
    -- AI配置
    ai_type VARCHAR(32),  -- aggressive/defensive/balanced/special
    ai_behavior TEXT,     -- JSON格式的AI行为配置
    
    -- 技能配置
    skill_ids TEXT,       -- JSON数组，技能ID列表
    
    -- 掉落配置
    drop_table_id VARCHAR(32),  -- 掉落表ID
    
    -- 平衡性调整字段
    balance_version INTEGER DEFAULT 1,  -- 平衡版本号
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 怪物AI系统

### AI类型

#### 1. 攻击型 (Aggressive)

**行为模式**:
- 优先攻击低血量目标
- 优先使用高伤害技能
- 忽略自身防御

**适用怪物**: 高攻击、低防御的怪物

**配置示例**:
```json
{
  "ai_type": "aggressive",
  "target_priority": ["lowest_hp", "lowest_defense"],
  "skill_priority": ["high_damage", "execute"],
  "defense_threshold": 0.2  // HP低于20%时考虑防御
}
```

#### 2. 防御型 (Defensive)

**行为模式**:
- 优先使用防御技能
- 优先攻击高威胁目标
- 保持安全血量

**适用怪物**: 高防御、低攻击的怪物

**配置示例**:
```json
{
  "ai_type": "defensive",
  "target_priority": ["highest_threat"],
  "skill_priority": ["defense", "heal", "attack"],
  "defense_threshold": 0.5  // HP低于50%时优先防御
}
```

#### 3. 平衡型 (Balanced)

**行为模式**:
- 根据情况选择攻击或防御
- 平衡使用技能
- 随机性行为

**适用怪物**: 属性平衡的怪物

**配置示例**:
```json
{
  "ai_type": "balanced",
  "target_priority": ["random", "lowest_hp"],
  "skill_priority": ["balanced"],
  "random_factor": 0.3  // 30%随机性
}
```

#### 4. 特殊型 (Special)

**行为模式**:
- 自定义行为逻辑
- 特殊机制触发
- 阶段转换

**适用怪物**: Boss和特殊怪物

**配置示例**:
```json
{
  "ai_type": "special",
  "phases": [
    {
      "hp_threshold": 1.0,
      "behavior": "aggressive",
      "skills": ["skill1", "skill2"]
    },
    {
      "hp_threshold": 0.5,
      "behavior": "defensive",
      "skills": ["skill3", "skill4"]
    }
  ]
}
```

### AI行为配置

```sql
-- 怪物AI配置表
CREATE TABLE monster_ai_configs (
    id VARCHAR(32) PRIMARY KEY,
    monster_id VARCHAR(32) NOT NULL,
    ai_type VARCHAR(32) NOT NULL,
    behavior_config TEXT NOT NULL,  -- JSON格式
    FOREIGN KEY (monster_id) REFERENCES monsters(id)
);
```

---

## 怪物技能系统

### 技能类型

#### 1. 攻击技能

- **普通攻击**: 基础物理/法术攻击
- **重击**: 高伤害，有冷却
- **范围攻击**: AOE伤害
- **持续伤害**: DOT效果

#### 2. 防御技能

- **护盾**: 吸收伤害
- **格挡**: 减少伤害
- **闪避提升**: 临时提高闪避率

#### 3. 控制技能

- **眩晕**: 无法行动
- **沉默**: 无法使用技能
- **减速**: 降低行动速度

#### 4. 治疗技能

- **自我治疗**: 恢复HP
- **持续恢复**: HOT效果

#### 5. 特殊技能

- **召唤**: 召唤小怪
- **变身**: 改变形态
- **阶段转换**: 改变行为模式

### 技能配置

```sql
-- 怪物技能配置表
CREATE TABLE monster_skills (
    id VARCHAR(32) PRIMARY KEY,
    monster_id VARCHAR(32) NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    skill_type VARCHAR(32) NOT NULL,
    cooldown INTEGER DEFAULT 0,
    use_condition TEXT,  -- JSON格式，使用条件
    priority INTEGER DEFAULT 0,  -- 优先级
    FOREIGN KEY (monster_id) REFERENCES monsters(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);
```

---

## 怪物掉落系统

### 掉落表设计

```sql
-- 掉落表
CREATE TABLE drop_tables (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT
);

-- 掉落项配置
CREATE TABLE drop_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drop_table_id VARCHAR(32) NOT NULL,
    item_id VARCHAR(32),
    item_type VARCHAR(32),  -- equipment/material/gold/exp
    min_quantity INTEGER DEFAULT 1,
    max_quantity INTEGER DEFAULT 1,
    drop_rate REAL NOT NULL,  -- 掉落概率 (0.0-1.0)
    min_level INTEGER DEFAULT 1,
    max_level INTEGER DEFAULT 100,
    FOREIGN KEY (drop_table_id) REFERENCES drop_tables(id)
);
```

### 掉落规则

| 怪物类型 | 装备品质 | 掉落率 | 说明 |
|---------|---------|--------|------|
| 普通 | 白色 | 80% | 基础装备 |
| 普通 | 绿色 | 20% | 魔法装备 |
| 精英 | 蓝色 | 60% | 稀有装备 |
| 精英 | 紫色 | 30% | 史诗装备 |
| 精英 | 橙色 | 10% | 传说装备 |
| Boss | 紫色 | 40% | 史诗装备 |
| Boss | 橙色 | 50% | 传说装备 |
| Boss | 独特 | 10% | 独特装备 |

### 掉落配置示例

```json
{
  "drop_table_id": "elite_wolf",
  "items": [
    {
      "item_type": "equipment",
      "quality": "blue",
      "drop_rate": 0.6,
      "min_quantity": 1,
      "max_quantity": 1
    },
    {
      "item_type": "equipment",
      "quality": "purple",
      "drop_rate": 0.3,
      "min_quantity": 1,
      "max_quantity": 1
    },
    {
      "item_type": "gold",
      "drop_rate": 1.0,
      "min_quantity": 10,
      "max_quantity": 50
    }
  ]
}
```

---

## 配置化设计

### 配置管理

所有怪物数据存储在配置表中，支持：

1. **热更新**: 修改配置表后，通过ConfigManager重新加载，无需重启服务
2. **版本管理**: 配置变更记录版本，支持回滚
3. **批量调整**: 通过SQL批量更新，快速调整平衡性

### 配置加载流程

```
配置表更新
    ↓
ConfigManager 检测变更
    ↓
重新加载配置到内存
    ↓
通知 MonsterManager 更新缓存
    ↓
新战斗使用新配置
```

### 配置接口

```go
// MonsterConfigManager 怪物配置管理器
type MonsterConfigManager interface {
    // 加载怪物配置
    LoadMonsterConfig(monsterID string) (*MonsterConfig, error)
    
    // 批量加载配置
    LoadAllMonsterConfigs() error
    
    // 热更新配置
    ReloadMonsterConfig(monsterID string) error
    
    // 获取配置版本
    GetConfigVersion(monsterID string) (int, error)
}
```

---

## 平衡性调整

### 调整原则

1. **数据驱动**: 所有平衡数据存储在配置表中
2. **快速调整**: 修改配置表即可，无需代码修改
3. **版本管理**: 记录每次调整，支持回滚
4. **测试验证**: 调整后通过测试验证效果

### 调整字段

所有可调整的字段都在配置表中：

- **属性系数**: `hp_multiplier`, `attack_multiplier`, `defense_multiplier`
- **特殊属性**: `crit_rate`, `crit_damage`, `dodge_rate`
- **掉落率**: `drop_rate` 在 `drop_items` 表中
- **技能强度**: 在 `skills` 配置表中

### 调整示例

#### 示例1: 调整怪物强度

```sql
-- 降低精英狼人的攻击力（平衡性调整）
UPDATE monsters 
SET attack_multiplier = 1.3,  -- 从1.5降低到1.3
    balance_version = balance_version + 1,
    last_updated = CURRENT_TIMESTAMP
WHERE id = 'elite_wolf';
```

#### 示例2: 调整掉落率

```sql
-- 提高Boss的橙色装备掉落率
UPDATE drop_items 
SET drop_rate = 0.6,  -- 从0.5提高到0.6
    last_updated = CURRENT_TIMESTAMP
WHERE drop_table_id = 'boss_drop_table' 
  AND item_type = 'equipment' 
  AND quality = 'orange';
```

#### 示例3: 批量调整

```sql
-- 批量调整所有精英怪物的HP（增加难度）
UPDATE monsters 
SET hp_multiplier = hp_multiplier * 1.1,  -- 增加10%
    balance_version = balance_version + 1,
    last_updated = CURRENT_TIMESTAMP
WHERE type = 'elite';
```

### 平衡性测试

调整后需要通过测试验证：

1. **战斗测试**: 测试不同等级队伍的战斗难度
2. **掉落测试**: 验证掉落率是否符合预期
3. **数据统计**: 收集战斗数据，分析平衡性

---

## 避免"一招走遍天下"

### 设计策略

#### 1. 多样化的怪物类型

- **物理免疫怪**: 只能用法术攻击
- **法术免疫怪**: 只能用物理攻击
- **高闪避怪**: 需要提高命中率
- **高防御怪**: 需要穿透或高攻击力
- **治疗怪**: 需要优先击杀或控制

#### 2. 不同的技能组合

- **AOE怪**: 需要分散站位或快速击杀
- **控制怪**: 需要解控技能或免疫
- **召唤怪**: 需要AOE技能清理小怪
- **阶段怪**: 不同阶段需要不同策略

#### 3. 属性克制

- **火系怪**: 弱冰系，强火系
- **暗影怪**: 弱神圣，强暗影
- **物理怪**: 弱法术，强物理

#### 4. 队伍搭配要求

- **单一DPS**: 无法应对所有情况
- **需要治疗**: 某些Boss需要持续治疗
- **需要控制**: 某些精英需要控制技能
- **需要坦克**: 高伤害Boss需要坦克吸收伤害

### 配置示例

```sql
-- 物理免疫怪配置
INSERT INTO monsters (id, name, type, ...) VALUES (
    'shadow_wraith',
    '暗影幽灵',
    'special',
    ...
);

-- 添加物理免疫技能
INSERT INTO monster_skills (monster_id, skill_id, skill_type) VALUES (
    'shadow_wraith',
    'physical_immunity',
    'passive'
);
```

---

## 总结

### 设计亮点

1. **丰富的怪物类型**: 普通、精英、Boss、特殊，带来不同挑战
2. **配置化设计**: 所有数据可配置，方便平衡性调整
3. **避免一招走遍天下**: 多样化机制，需要不同策略
4. **版本化管理**: 支持配置回滚，减少数据迁移复杂度

### 后续扩展

- [ ] 怪物图鉴系统
- [ ] 怪物刷新机制
- [ ] 怪物组合系统（多个怪物协同）
- [ ] 动态难度调整

---

**文档版本**: v1.0  
**最后更新**: 2025年


