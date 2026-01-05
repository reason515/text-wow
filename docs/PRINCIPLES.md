# Text WoW 开发基本原则

> ⚠️ **重要**: 在处理任何开发任务前，请先阅读本文档，确保遵循项目的基本原则。

---

## 📋 必须遵循的原则

### 1. 执行任务前先阅读相关设计文档

**开始任何开发任务前：**
- ✅ **必须先阅读相关的设计文档**，了解系统设计和实现要求
- ✅ 确认理解设计意图和实现方式
- ✅ 如有疑问，先查阅文档或询问，不要盲目实现

**相关文档位置：**
- 所有设计文档位于 `docs/` 目录
- 系统架构：`docs/architecture.md`
- 游戏机制索引：`docs/game_mechanics.md`
- 具体系统设计文档请参考 `docs/game_mechanics.md` 中的索引

### 2. 代码、测试用例和文档同步更新

**每次改动必须确保：**
- ✅ 代码修改后，同步更新相关的测试用例
- ✅ 同步更新相关的文档（设计文档、API文档等）
- ✅ 确保所有测试用例通过

**提交前检查：**
- [ ] 代码已修改
- [ ] 测试用例已更新并通过
- [ ] 相关文档已更新

### 3. 提交信息使用英文

**每次提交代码时：**
- ✅ **必须使用英文 commit 信息**
- ❌ **禁止使用中文 commit 信息**（GitLab 会显示乱码）

**提交信息格式：**
```
<type>(<scope>): <subject>

<body>

<footer>
```

**示例：**
```
feat(battle): add multi-character battle support

- Support up to 5 characters in battle
- Add threat system
- Add team synergy skills

Closes #123
```

```
fix(equipment): fix equipment enhancement cost calculation

- Fix enhancement level calculation error
- Fix material consumption calculation error

Fixes #456
```

### 4. 测试用例使用 YAML 自然语言描述

**编写测试用例时：**
- ✅ **必须使用 YAML 格式的自然语言测试用例**
- ✅ 测试用例应使用中文自然语言描述，易于理解
- ❌ **禁止使用代码形式的测试用例**（除非特殊情况）

**测试用例格式参考：**
- 详细格式请参考 `docs/testing_system_design.md`
- 测试用例文件存放在 `server/internal/test/cases/` 目录

**示例：**
```yaml
- name: "基础战斗流程测试"
  description: "测试一个简单的战斗流程，角色攻击怪物直到一方死亡"
  setup:
    - "创建一个1级人类战士角色，HP=25，攻击力=8"
    - "创建一个1级森林狼怪物，HP=20，攻击力=5"
  steps:
    - action: "开始战斗"
      expected: "战斗状态为进行中"
```

---

**文档版本**: v1.0  
**最后更新**: 2025年
