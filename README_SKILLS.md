# 技能数据加载指南

## 问题症状

创建角色时出现提示：
```
⚠ 无法获取初始技能列表，请检查数据库是否已加载技能数据
```

## 原因

数据库中没有加载战士技能数据（`warrior_skills.sql`）。

## 解决方案

### 方法1：使用加载工具（推荐）

```bash
cd server/cmd/load_skills
go run main.go
```

这会自动：
1. 检查当前技能数据
2. 加载 `warrior_skills.sql` 文件
3. 验证所有初始技能是否已加载

### 方法2：检查并手动加载

```bash
# 检查技能数据是否存在
cd server/cmd/load_skills
go run main.go -check

# 如果不存在，加载技能数据
go run main.go
```

### 方法3：重启服务器（自动加载）

服务器启动时会自动检查并加载技能数据。如果技能数据缺失，会在日志中显示警告，并尝试自动加载。

## 验证

加载完成后，工具会显示：
```
✅ 技能数据加载完成！
   加载前: 0 个技能
   加载后: 22 个技能

检查初始技能池:
  ✅ warrior_heroic_strike
  ✅ warrior_taunt
  ✅ warrior_shield_block
  ✅ warrior_cleave
  ✅ warrior_slam
  ✅ warrior_battle_shout
  ✅ warrior_demoralizing_shout
  ✅ warrior_last_stand
  ✅ warrior_charge

✅ 所有初始技能都已加载
```

## 故障排除

### 问题1：找不到 warrior_skills.sql 文件

**错误**：`warrior_skills.sql not found`

**解决**：确保 `server/database/warrior_skills.sql` 文件存在

### 问题2：SQL执行失败

**错误**：`failed to execute warrior_skills.sql`

**解决**：
1. 检查数据库文件权限
2. 检查SQL文件格式是否正确
3. 查看详细错误信息

### 问题3：技能数据已存在但仍有错误

**解决**：使用强制重新加载
```bash
cd server/cmd/load_skills
go run main.go -force
```

## 自动加载机制

服务器启动时（`database.Init()`）会自动：
1. 检查是否有基础数据（races表）
2. 检查是否有战士技能数据
3. 如果缺失，自动加载 `warrior_skills.sql`

如果自动加载失败，会在日志中显示警告，但不影响服务器启动。此时需要手动运行加载工具。

## 相关文件

- `server/database/warrior_skills.sql` - 战士技能数据文件
- `server/cmd/load_skills/main.go` - 技能加载工具
- `server/internal/database/db.go` - 数据库初始化代码





































