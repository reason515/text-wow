# 怪物强度快速调整指南

## 快速开始

### 1. 查看当前配置

```bash
python adjust_monster_strength.py --list
```

### 2. 调整强度（例如：整体提升10%）

```bash
# 为1-10级提升10%
python adjust_monster_strength.py --set 1 10 1.1 1.1 1.1 0.0

# 为11-20级提升10%
python adjust_monster_strength.py --set 11 20 1.1 1.1 1.1 0.0

# ... 其他等级段类似

# 应用所有配置
python adjust_monster_strength.py --apply
```

### 3. 重置为默认值

```bash
python adjust_monster_strength.py --reset
python adjust_monster_strength.py --apply
```

## 参数说明

`--set` 参数格式：`等级下限 等级上限 HP倍数 攻击倍数 防御倍数 暴击率加成`

- **等级下限/上限**：要调整的等级范围（包含）
- **HP倍数**：生命值倍数，1.5表示+50%
- **攻击倍数**：攻击力倍数（物理和魔法），1.4表示+40%
- **防御倍数**：防御力倍数（物理和魔法），1.4表示+40%
- **暴击率加成**：暴击率绝对值加成，0.02表示+2%

## 常用场景

### 场景1：整体提升20%强度

```bash
python adjust_monster_strength.py --set 1 60 1.2 1.2 1.2 0.0
python adjust_monster_strength.py --apply
```

### 场景2：只调整低等级（1-20级）提升30%

```bash
python adjust_monster_strength.py --set 1 10 1.3 1.3 1.3 0.0
python adjust_monster_strength.py --set 11 20 1.3 1.3 1.3 0.0
python adjust_monster_strength.py --apply
```

### 场景3：只提升攻击力，不提升防御

```bash
python adjust_monster_strength.py --set 1 60 1.0 1.2 1.0 0.0
python adjust_monster_strength.py --apply
```

### 场景4：恢复原始强度

```bash
python adjust_monster_strength.py --set 1 60 1.0 1.0 1.0 0.0
python adjust_monster_strength.py --apply
```

## 注意事项

⚠️ **重要**：
- `--apply` 会直接修改数据库中的怪物数据
- 建议在应用前备份数据库：`cp server/game.db server/game.db.backup`
- 配置是累积的，多次应用会重复计算
- 如果要从原始数据开始，需要先重新导入 `seed.sql`

## 系统架构

1. **配置表** (`monster_strength_config`)：存储强度调整系数
2. **Python工具** (`adjust_monster_strength.py`)：快速调整工具
3. **Go工具** (`server/cmd/adjust_monster_strength/`)：Go语言实现（功能相同）
4. **自动迁移**：服务器启动时自动创建配置表

详细文档请参考：`docs/monster_strength_adjustment.md`































