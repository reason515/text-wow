# 测试系统设计文档

> 📌 **核心设计理念**: 使用自然语言描述测试用例，易于理解和维护，支持自动化测试

---

## 📋 目录

1. [系统概览](#系统概览)
2. [测试框架架构](#测试框架架构)
3. [自然语言测试用例格式](#自然语言测试用例格式)
4. [测试用例编写指南](#测试用例编写指南)
5. [测试执行流程](#测试执行流程)
6. [测试报告格式](#测试报告格式)

---

## 系统概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          测试系统架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   测试用例 ──→ 测试解析 ──→ 测试执行 ──→ 断言验证 ──→ 测试报告              │
│      │            │            │            │            │                 │
│   YAML格式    YAML解析器   测试运行器   断言执行器   报告生成器             │
│                                                                             │
│   测试类型:                                                                  │
│   ├─ 单元测试: 测试单个函数或模块                                            │
│   ├─ 集成测试: 测试系统间交互                                                │
│   └─ 端到端测试: 测试完整流程                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 测试框架架构

### 架构层次

```
┌─────────────────────────────────────────────────────────────┐
│                    测试用例层 (Test Cases)                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  自然语言测试用例 (Natural Language Test Cases)        │   │
│  │  - 使用YAML/JSON格式，接近自然语言描述                │   │
│  │  - 易于理解和维护                                      │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    测试执行层 (Test Runner)                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  单元测试     │  │  集成测试     │  │  端到端测试   │     │
│  │  Unit Tests  │  │ Integration  │  │  E2E Tests   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    测试报告层 (Test Reports)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  覆盖率报告   │  │  性能报告     │  │  可视化报告   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构

```
server/test/
├── cases/                    # 自然语言测试用例
│   ├── battle/              # 战斗系统测试用例
│   │   ├── basic_combat.yaml
│   │   ├── damage_calculation.yaml
│   │   └── team_combat.yaml
│   ├── skill/               # 技能系统测试用例
│   │   ├── skill_learning.yaml
│   │   ├── skill_usage.yaml
│   │   └── buff_debuff.yaml
│   ├── equipment/           # 装备系统测试用例
│   │   ├── equip_unequip.yaml
│   │   ├── affix_generation.yaml
│   │   └── equipment_enhancement.yaml
│   ├── calculator/          # 数值计算系统测试用例
│   │   ├── attribute_conversion.yaml
│   │   ├── damage_calculation.yaml
│   │   └── defense_reduction.yaml
│   ├── team/                # 队伍系统测试用例
│   │   ├── team_management.yaml
│   │   └── team_combat.yaml
│   ├── monster/             # 怪物系统测试用例
│   │   ├── monster_generation.yaml
│   │   └── monster_ai.yaml
│   └── economy/             # 经济系统测试用例
│       ├── gold_acquisition.yaml
│       └── gold_consumption.yaml
├── runner/                   # 测试执行器
│   ├── test_runner.go       # 测试运行器
│   ├── yaml_parser.go       # YAML解析器
│   ├── assertion.go         # 断言执行器
│   └── reporter.go          # 报告生成器
└── fixtures/                # 测试数据
    ├── characters.json
    ├── monsters.json
    └── skills.json
```

---

## 自然语言测试用例格式

### YAML格式

```yaml
# 测试套件定义
test_suite: "战斗系统基础测试"
description: "测试战斗系统的核心功能，包括回合流程、伤害计算、胜负判定"
version: "1.0"

# 测试用例列表
tests:
  - name: "基础战斗流程测试"
    description: "测试一个简单的战斗流程，角色攻击怪物直到一方死亡"
    category: "integration"  # unit/integration/e2e
    priority: "high"        # high/medium/low
    
    # 前置条件
    setup:
      - "创建一个1级人类战士角色，HP=25，攻击力=8"
      - "创建一个1级森林狼怪物，HP=20，攻击力=5"
      - "初始化战斗系统"
    
    # 测试步骤
    steps:
      - action: "开始战斗"
        expected: "战斗状态为进行中"
        timeout: 5  # 秒
        
      - action: "执行一个回合"
        expected: "角色攻击怪物，怪物受到伤害"
        assertions:
          - "怪物HP应该减少"
          - "战斗日志应该包含攻击信息"
        
      - action: "怪物反击"
        expected: "怪物攻击角色，角色受到伤害"
        assertions:
          - "角色HP应该减少"
          - "战斗日志应该包含受击信息"
        
      - action: "继续战斗直到怪物死亡"
        expected: "战斗结束，角色获胜，获得经验和金币"
        max_rounds: 20  # 最多20回合
    
    # 断言列表
    assertions:
      - type: "equals"
        target: "character.hp"
        expected: "> 0"
        message: "角色HP应该大于0"
      
      - type: "equals"
        target: "monster.hp"
        expected: "0"
        message: "怪物HP应该等于0"
      
      - type: "greater_than"
        target: "character.exp"
        expected: "0"
        message: "角色应该获得经验值"
      
      - type: "contains"
        target: "battle_logs"
        expected: "攻击"
        message: "战斗日志应该包含攻击信息"
    
    # 清理
    teardown:
      - "清理战斗状态"
      - "重置角色数据"

  - name: "伤害计算测试"
    description: "测试物理伤害计算，包括基础伤害、防御减伤、暴击"
    category: "unit"
    priority: "high"
    
    setup:
      - "创建一个攻击力=20的角色"
      - "创建一个防御力=10的怪物"
    
    steps:
      - action: "角色使用普通攻击"
        expected: "计算基础伤害=20"
        assertions:
          - "基础伤害应该等于20"
      
      - action: "应用防御减伤"
        expected: "最终伤害=20 * (1 - 10/(10+100)) ≈ 18"
        assertions:
          - "最终伤害应该约等于18"
          - "伤害应该考虑防御减伤"
      
      - action: "如果暴击，应用暴击倍率"
        expected: "最终伤害=18 * 1.5 = 27"
        assertions:
          - "如果暴击，伤害应该乘以暴击倍率"
    
    assertions:
      - type: "approximately"
        target: "damage"
        expected: "18"
        tolerance: "1"
        message: "伤害应该约等于18（考虑误差）"
```

### JSON格式（备选）

```json
{
  "test_suite": "战斗系统基础测试",
  "description": "测试战斗系统的核心功能",
  "version": "1.0",
  "tests": [
    {
      "name": "基础战斗流程测试",
      "description": "测试一个简单的战斗流程",
      "category": "integration",
      "priority": "high",
      "setup": [
        "创建一个1级人类战士角色，HP=25，攻击力=8",
        "创建一个1级森林狼怪物，HP=20，攻击力=5"
      ],
      "steps": [
        {
          "action": "开始战斗",
          "expected": "战斗状态为进行中"
        }
      ],
      "assertions": [
        {
          "type": "equals",
          "target": "character.hp",
          "expected": "> 0",
          "message": "角色HP应该大于0"
        }
      ]
    }
  ]
}
```

---

## 测试用例编写指南

### 编写原则

1. **使用自然语言**: 测试用例应该像描述游戏玩法一样自然
2. **清晰明确**: 每个步骤和断言都应该清晰明确
3. **易于理解**: 非技术人员也能理解测试用例
4. **完整覆盖**: 覆盖所有核心功能和边界情况

### 命名规范

#### 测试套件命名

```
{系统名称}_{功能模块}_测试

示例:
- 战斗系统_基础战斗_测试
- 装备系统_词缀生成_测试
- 技能系统_冷却机制_测试
```

#### 测试用例命名

```
{测试场景}_{预期结果}

示例:
- 基础战斗流程_角色获胜
- 伤害计算_防御减伤正确
- 装备强化_属性提升
```

### 测试用例结构

#### 必需字段

- `name`: 测试用例名称
- `description`: 测试用例描述
- `setup`: 前置条件
- `steps`: 测试步骤
- `assertions`: 断言列表

#### 可选字段

- `category`: 测试类型（unit/integration/e2e）
- `priority`: 优先级（high/medium/low）
- `timeout`: 超时时间（秒）
- `teardown`: 清理步骤
- `tags`: 标签列表

### 断言类型

#### 1. equals - 相等断言

```yaml
- type: "equals"
  target: "character.hp"
  expected: "25"
  message: "角色HP应该等于25"
```

#### 2. greater_than - 大于断言

```yaml
- type: "greater_than"
  target: "character.exp"
  expected: "0"
  message: "角色经验值应该大于0"
```

#### 3. less_than - 小于断言

```yaml
- type: "less_than"
  target: "monster.hp"
  expected: "20"
  message: "怪物HP应该小于20"
```

#### 4. contains - 包含断言

```yaml
- type: "contains"
  target: "battle_logs"
  expected: "攻击"
  message: "战斗日志应该包含'攻击'"
```

#### 5. approximately - 近似断言

```yaml
- type: "approximately"
  target: "damage"
  expected: "18"
  tolerance: "1"
  message: "伤害应该约等于18（误差±1）"
```

#### 6. range - 范围断言

```yaml
- type: "range"
  target: "damage"
  expected: "[15, 25]"
  message: "伤害应该在15-25之间"
```

---

## 测试执行流程

### 执行步骤

```
1. 加载测试用例
   ├─ 读取YAML文件
   ├─ 解析测试用例
   └─ 验证测试用例格式

2. 执行前置条件
   ├─ 创建测试环境
   ├─ 初始化测试数据
   └─ 准备测试对象

3. 执行测试步骤
   ├─ 按顺序执行每个步骤
   ├─ 验证每个步骤的预期结果
   └─ 记录执行日志

4. 执行断言
   ├─ 按顺序执行每个断言
   ├─ 验证断言结果
   └─ 记录断言结果

5. 执行清理
   ├─ 清理测试环境
   ├─ 重置测试数据
   └─ 释放资源

6. 生成报告
   ├─ 汇总测试结果
   ├─ 生成测试报告
   └─ 输出测试日志
```

### 执行器接口

```go
type TestRunner interface {
    // 运行测试套件
    RunTestSuite(suitePath string) (*TestSuiteResult, error)
    
    // 运行单个测试用例
    RunTestCase(testCase *TestCase) (*TestCaseResult, error)
    
    // 运行所有测试
    RunAllTests(testDir string) (*TestResult, error)
}
```

---

## 测试报告格式

### 报告结构

```json
{
  "test_suite": "战斗系统基础测试",
  "total_tests": 10,
  "passed_tests": 8,
  "failed_tests": 2,
  "skipped_tests": 0,
  "duration": "2.5s",
  "timestamp": "2025-01-01T12:00:00Z",
  "results": [
    {
      "test_name": "基础战斗流程测试",
      "status": "passed",
      "duration": "0.5s",
      "assertions": [
        {
          "type": "equals",
          "target": "character.hp",
          "expected": "> 0",
          "actual": "15",
          "status": "passed"
        }
      ]
    },
    {
      "test_name": "伤害计算测试",
      "status": "failed",
      "duration": "0.3s",
      "error": "伤害计算错误，期望18，实际20",
      "assertions": [
        {
          "type": "approximately",
          "target": "damage",
          "expected": "18",
          "actual": "20",
          "status": "failed",
          "message": "伤害应该约等于18（误差±1），实际20超出范围"
        }
      ]
    }
  ]
}
```

### 报告类型

#### 1. 控制台报告

```
=== 测试报告 ===
测试套件: 战斗系统基础测试
总测试数: 10
通过: 8
失败: 2
跳过: 0
耗时: 2.5s

[PASS] 基础战斗流程测试 (0.5s)
[FAIL] 伤害计算测试 (0.3s)
  错误: 伤害计算错误，期望18，实际20
[PASS] 技能使用测试 (0.4s)
...
```

#### 2. HTML报告

- 可视化测试结果
- 图表展示通过率
- 详细的错误信息
- 测试覆盖率

#### 3. JSON报告

- 机器可读格式
- 便于集成到CI/CD
- 支持自动化分析

---

## 测试覆盖率

### 覆盖率目标

- **单元测试覆盖率**: > 70%
- **集成测试覆盖率**: > 50%
- **核心功能覆盖率**: 100%

### 覆盖率报告

```
=== 测试覆盖率报告 ===
总覆盖率: 75.3%

按模块:
- 战斗系统: 82.5%
- 技能系统: 78.3%
- 装备系统: 71.2%
- 数值计算: 95.0%
- 队伍系统: 65.8%
```

---

## 总结

### 设计亮点

1. **自然语言测试用例**: 易于理解和维护
2. **灵活的断言系统**: 支持多种断言类型
3. **详细的测试报告**: 清晰的测试结果展示
4. **自动化执行**: 支持CI/CD集成

### 后续扩展

- [ ] 测试用例生成工具
- [ ] 测试数据管理工具
- [ ] 性能测试支持
- [ ] 压力测试支持

---

**文档版本**: v1.0  
**最后更新**: 2025年


