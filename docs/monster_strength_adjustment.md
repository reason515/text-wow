# 怪物强度调整系统

## 概述

本系统提供了一个快速调整怪物强度的机制，无需手动修改 `seed.sql` 文件。通过配置表和工具脚本，可以快速调整所有怪物的属性。

## 系统组成

### 1. 配置表 (`monster_strength_config`)

存储不同等级段的强度调整系数：

- `level_min`, `level_max`: 等级范围
- `hp_multiplier`: 生命值倍数
- `attack_multiplier`: 攻击力倍数（物理和魔法）
- `defense_multiplier`: 防御力倍数（物理和魔法）
- `crit_rate_bonus`: 暴击率加成（绝对值，如0.02表示+2%）
- `is_active`: 是否启用

### 2. Python工具 (`adjust_monster_strength.py`)

快速调整工具，支持命令行操作。

### 3. Go工具 (`server/cmd/adjust_monster_strength/main.go`)

Go语言实现的调整工具，功能与Python版本相同。

## 使用方法

### Python工具

```bash
# 列出所有配置
python adjust_monster_strength.py --list

# 设置配置（等级1-10，HP×1.5，攻击×1.4，防御×1.4，暴击率+2%）
python adjust_monster_strength.py --set 1 10 1.5 1.4 1.4 0.02

# 应用配置到怪物数据
python adjust_monster_strength.py --apply

# 重置为默认配置
python adjust_monster_strength.py --reset
```

### Go工具

```bash
cd server/cmd/adjust_monster_strength
go run main.go -list
go run main.go -min 1 -max 10 -hp 1.5 -attack 1.4 -defense 1.4 -crit 0.02
go run main.go -apply
```

## 工作流程

1. **设置配置**：使用工具设置不同等级段的强度系数
2. **应用配置**：将配置应用到数据库中的怪物数据
3. **验证结果**：使用 `--list` 查看当前配置

## 默认配置

当前默认配置（已应用的强度提升）：

- 1-10级：HP +50%, 攻击 +40%, 防御 +40%, 暴击率 +2%
- 10-20级：HP +45%, 攻击 +35%, 防御 +35%, 暴击率 +2%
- 20-30级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +2%
- 30-40级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +3%
- 40-50级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%
- 50-60级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%

## 注意事项

1. **应用配置会直接修改数据库**：建议在应用前备份数据库
2. **配置是累积的**：多次应用会重复计算，建议先重置或从原始数据开始
3. **暴击率上限**：系统会自动限制暴击率不超过40%
4. **数值四舍五入**：所有计算结果会四舍五入到整数

## 数据库迁移

首次使用时，需要运行迁移脚本创建配置表：

```sql
-- 运行 server/database/migrate_monster_strength_config.sql
```

或者使用工具时会自动创建表。

## 示例场景

### 场景1：整体提升10%强度

```bash
# 为所有等级段提升10%
python adjust_monster_strength.py --set 1 10 1.1 1.1 1.1 0.0
python adjust_monster_strength.py --set 11 20 1.1 1.1 1.1 0.0
# ... 其他等级段
python adjust_monster_strength.py --apply
```

### 场景2：只调整特定等级段

```bash
# 只提升1-10级的强度
python adjust_monster_strength.py --set 1 10 1.2 1.2 1.2 0.01
python adjust_monster_strength.py --apply
```

### 场景3：恢复原始强度

```bash
# 将所有倍数设为1.0，加成设为0
python adjust_monster_strength.py --set 1 60 1.0 1.0 1.0 0.0
python adjust_monster_strength.py --apply
```







