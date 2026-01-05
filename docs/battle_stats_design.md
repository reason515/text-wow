# 战斗统计系统设计文档

> 📌 **核心设计理念**: 详细收集战斗数据，提供数据分析和策略优化建议，帮助玩家优化角色和策略

---

## 📋 目录

1. [系统概览](#系统概览)
2. [数据收集机制](#数据收集机制)
3. [统计分析功能](#统计分析功能)
4. [数据可视化](#数据可视化)
5. [策略优化建议](#策略优化建议)

---

## 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          战斗统计系统架构                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   战斗进行 ──→ 数据收集 ──→ 统计分析 ──→ 数据展示 ──→ 策略建议              │
│      │            │            │            │            │                 │
│   实时收集    存储数据    计算指标    图表展示    优化建议                   │
│                                                                             │
│   统计维度:                                                                  │
│   ├─ 伤害统计: 总伤害、DPS、技能伤害分布                                    │
│   ├─ 承伤统计: 总承伤、减伤效果                                             │
│   ├─ 治疗统计: 总治疗、HPS                                                  │
│   ├─ 暴击/闪避统计: 暴击率、闪避率                                          │
│   └─ 技能使用统计: 技能使用频率、效果                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 数据收集机制

### 收集时机

#### 1. 战斗进行中

- **实时收集**: 每次攻击、治疗、技能使用都记录
- **内存存储**: 战斗数据临时存储在内存中
- **性能优化**: 批量写入，减少数据库压力

#### 2. 战斗结束时

- **批量保存**: 战斗结束后批量保存统计数据
- **数据汇总**: 汇总单场战斗的所有数据
- **持久化**: 保存到数据库

### 收集数据

#### 伤害数据

```go
type DamageStats struct {
    TotalDamage    int  // 总伤害
    PhysicalDamage int  // 物理伤害
    MagicDamage    int  // 魔法伤害
    FireDamage     int  // 火焰伤害
    FrostDamage    int  // 冰霜伤害
    ShadowDamage   int  // 暗影伤害
    HolyDamage     int  // 神圣伤害
    NatureDamage   int  // 自然伤害
    DotDamage      int  // DOT伤害
    CritCount      int  // 暴击次数
    CritDamage     int  // 暴击总伤害
    MaxCrit        int  // 最高单次暴击
}
```

#### 承伤数据

```go
type DamageTakenStats struct {
    TotalDamageTaken int  // 总承伤
    PhysicalTaken    int  // 物理承伤
    MagicTaken       int  // 魔法承伤
    DamageBlocked    int  // 格挡伤害
    DamageAbsorbed   int  // 护盾吸收
    DodgeCount       int  // 闪避次数
    BlockCount       int  // 格挡次数
    HitCount         int  // 被命中次数
}
```

#### 治疗数据

```go
type HealingStats struct {
    TotalHealing     int  // 总治疗
    HealingReceived  int  // 受到治疗
    Overhealing      int  // 过量治疗
    SelfHealing      int  // 自我治疗
    HotHealing       int  // HOT治疗
}
```

#### 技能使用数据

```go
type SkillUsageStats struct {
    SkillID      string
    UseCount     int  // 使用次数
    HitCount     int  // 命中次数
    CritCount    int  // 暴击次数
    TotalDamage  int  // 总伤害
    TotalHealing int  // 总治疗
    ResourceCost int  // 总消耗资源
}
```

---

## 统计分析功能

### 伤害分析

#### DPS计算

```
DPS = 总伤害 / 战斗时长(秒)

平均DPS = 总伤害 / 战斗回合数
```

#### 伤害分布

```
物理伤害占比 = 物理伤害 / 总伤害 × 100%
魔法伤害占比 = 魔法伤害 / 总伤害 × 100%
元素伤害占比 = 元素伤害 / 总伤害 × 100%
```

#### 技能伤害分析

```
技能伤害占比 = 技能伤害 / 总伤害 × 100%
技能DPS = 技能总伤害 / 战斗时长
技能效率 = 技能总伤害 / 技能消耗资源
```

### 承伤分析

#### 减伤效果

```
实际减伤率 = (理论伤害 - 实际伤害) / 理论伤害 × 100%
平均每次承伤 = 总承伤 / 被命中次数
```

#### 生存能力

```
生存时间 = 战斗时长
平均每秒承伤 = 总承伤 / 战斗时长
闪避率 = 闪避次数 / (闪避次数 + 被命中次数) × 100%
```

### 治疗分析

#### HPS计算

```
HPS = 总治疗 / 战斗时长(秒)

平均HPS = 总治疗 / 战斗回合数
```

#### 治疗效率

```
治疗效率 = 实际治疗 / (实际治疗 + 过量治疗) × 100%
平均每次治疗 = 总治疗 / 治疗次数
```

### 暴击/闪避分析

#### 暴击统计

```
实际暴击率 = 暴击次数 / 总攻击次数 × 100%
平均暴击伤害 = 暴击总伤害 / 暴击次数
暴击伤害占比 = 暴击总伤害 / 总伤害 × 100%
```

#### 闪避统计

```
实际闪避率 = 闪避次数 / (闪避次数 + 被命中次数) × 100%
闪避减少伤害 = 理论伤害 × 闪避次数
```

---

## 数据可视化

### 图表类型

#### 1. 伤害分布饼图

- **物理伤害**: 蓝色
- **魔法伤害**: 紫色
- **火焰伤害**: 红色
- **冰霜伤害**: 浅蓝色
- **暗影伤害**: 黑色
- **神圣伤害**: 金色
- **自然伤害**: 绿色

#### 2. 技能伤害柱状图

- **X轴**: 技能名称
- **Y轴**: 伤害数值
- **颜色**: 根据伤害类型

#### 3. 时间序列折线图

- **X轴**: 战斗时间
- **Y轴**: DPS/HPS
- **多条线**: 不同角色或技能

#### 4. 对比雷达图

- **维度**: 伤害、承伤、治疗、暴击、闪避
- **对比**: 不同角色或不同战斗

### 数据展示接口

```go
type StatsVisualization interface {
    // 获取伤害分布数据
    GetDamageDistribution(battleID int) (*DamageDistribution, error)
    
    // 获取技能使用数据
    GetSkillUsageData(battleID int) ([]*SkillUsageData, error)
    
    // 获取时间序列数据
    GetTimeSeriesData(battleID int) (*TimeSeriesData, error)
    
    // 获取对比数据
    GetComparisonData(battleIDs []int) (*ComparisonData, error)
}
```

---

## 策略优化建议

### 建议类型

#### 1. 伤害优化建议

```
IF 物理伤害占比 < 50% AND 角色是物理职业 THEN
    建议: "考虑提升物理攻击力，物理伤害占比偏低"
END IF

IF 技能伤害占比 < 30% THEN
    建议: "考虑多使用技能，技能伤害占比偏低"
END IF

IF 暴击率 < 10% THEN
    建议: "考虑提升暴击率，当前暴击率较低"
END IF
```

#### 2. 生存优化建议

```
IF 平均每秒承伤 > 最大HP / 10 THEN
    建议: "承伤过高，考虑提升防御或生命值"
END IF

IF 闪避率 < 5% THEN
    建议: "考虑提升闪避率，当前闪避率较低"
END IF

IF 治疗效率 < 70% THEN
    建议: "过量治疗较多，考虑优化治疗时机"
END IF
```

#### 3. 技能优化建议

```
IF 技能A使用频率 > 50% AND 技能A伤害占比 < 20% THEN
    建议: "技能A使用频率高但伤害低，考虑替换或升级"
END IF

IF 技能B伤害效率 > 技能A伤害效率 × 1.5 THEN
    建议: "技能B效率更高，考虑优先使用技能B"
END IF
```

#### 4. 属性优化建议

```
IF 力量 < 敏捷 × 2 AND 角色是物理职业 THEN
    建议: "考虑提升力量，力量对物理职业更重要"
END IF

IF 耐力 < 力量 AND 角色是坦克 THEN
    建议: "考虑提升耐力，坦克需要更高的生命值"
END IF
```

### 建议生成接口

```go
type StrategyOptimizer interface {
    // 生成优化建议
    GenerateOptimizationSuggestions(battleStats *BattleStats) ([]*OptimizationSuggestion, error)
    
    // 分析技能使用
    AnalyzeSkillUsage(battleStats *BattleStats) (*SkillUsageAnalysis, error)
    
    // 分析属性分配
    AnalyzeAttributeAllocation(character *Character, battleStats *BattleStats) (*AttributeAnalysis, error)
}
```

---

## 数据库设计

### 战斗记录表

```sql
-- 战斗记录表 (已存在)
CREATE TABLE battle_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    zone_id VARCHAR(32),
    battle_type VARCHAR(16),  -- pve/pvp/boss
    total_rounds INTEGER,
    duration_seconds INTEGER,
    result VARCHAR(16),  -- victory/defeat
    team_damage_dealt INTEGER,
    team_damage_taken INTEGER,
    team_healing_done INTEGER,
    exp_gained INTEGER,
    gold_gained INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 角色战斗统计表

```sql
-- 角色战斗统计表 (已存在)
CREATE TABLE battle_character_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    damage_dealt INTEGER DEFAULT 0,
    damage_taken INTEGER DEFAULT 0,
    healing_done INTEGER DEFAULT 0,
    crit_count INTEGER DEFAULT 0,
    dodge_count INTEGER DEFAULT 0,
    -- ... 更多统计字段
    FOREIGN KEY (battle_id) REFERENCES battle_records(id),
    FOREIGN KEY (character_id) REFERENCES characters(id)
);
```

### 技能使用统计表

```sql
-- 技能使用统计表 (已存在)
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
    FOREIGN KEY (battle_id) REFERENCES battle_records(id),
    FOREIGN KEY (character_id) REFERENCES characters(id),
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);
```

---

## 总结

### 设计亮点

1. **详细的数据收集**: 收集所有战斗相关数据
2. **丰富的统计分析**: 多维度分析战斗数据
3. **直观的数据可视化**: 图表展示数据
4. **智能的策略建议**: 基于数据分析提供优化建议

### 后续扩展

- [ ] 实时数据展示
- [ ] 数据导出功能
- [ ] 数据对比工具
- [ ] 自动优化系统

---

**文档版本**: v1.0  
**最后更新**: 2025年


