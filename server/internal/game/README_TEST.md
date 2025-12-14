# 战士技能系统测试文档

## 测试文件结构

所有测试文件都使用 `_test.go` 后缀，与正常代码良好区分：

- `buff_manager_test.go` - Buff/Debuff管理器测试
- `passive_skill_manager_test.go` - 被动技能管理器测试
- `skill_manager_test.go` - 技能管理器测试
- `warrior_skills_integration_test.go` - 战士技能集成测试

## 测试覆盖范围

### 1. BuffManager 测试 (`buff_manager_test.go`)

#### 基础功能测试
- `TestBuffManager_ApplyBuff` - 测试应用Buff
- `TestBuffManager_ApplyDebuff` - 测试应用Debuff
- `TestBuffManager_GetBuffValue` - 测试获取Buff值（累加）
- `TestBuffManager_TickBuffs` - 测试Buff持续时间减少
- `TestBuffManager_ClearBuffs` - 测试清除Buff

#### 敌人Debuff系统测试
- `TestBuffManager_ApplyEnemyDebuff` - 测试应用敌人Debuff
- `TestBuffManager_GetEnemyDebuffValue` - 测试获取敌人Debuff值
- `TestBuffManager_TickEnemyDebuffs` - 测试敌人Debuff持续时间减少
- `TestBuffManager_ClearEnemyDebuffs` - 测试清除敌人Debuff

#### Buff效果计算测试
- `TestBuffManager_CalculateDamageTakenWithBuffs` - 测试减伤计算
- `TestBuffManager_CalculateDamageWithBuffs` - 测试伤害加成计算

### 2. PassiveSkillManager 测试 (`passive_skill_manager_test.go`)

#### 基础功能测试
- `TestPassiveSkillManager_GetPassiveModifier` - 测试获取被动技能修正值
- `TestPassiveSkillManager_GetPassiveModifier_Multiple` - 测试多个被动技能累加
- `TestPassiveSkillManager_GetPassiveModifier_MultiAttribute` - 测试多属性被动技能
- `TestPassiveSkillManager_GetPassiveSkills` - 测试获取被动技能列表
- `TestPassiveSkillManager_HasPassiveSkill` - 测试检查被动技能
- `TestPassiveSkillManager_GetPassiveSkillLevel` - 测试获取被动技能等级

### 3. SkillManager 测试 (`skill_manager_test.go`)

#### 技能伤害计算测试
- `TestSkillManager_CalculateSkillDamage_Basic` - 测试基础技能伤害计算
- `TestSkillManager_CalculateSkillDamage_WithPassiveModifiers` - 测试带被动技能加成的伤害计算
- `TestSkillManager_CalculateSkillDamage_WithBuffModifiers` - 测试带Buff加成的伤害计算
- `TestSkillManager_CalculateSkillDamage_ShieldSlam` - 测试盾牌猛击特殊伤害计算
- `TestSkillManager_CalculateSkillDamage_WithEnemyDebuff` - 测试带敌人Debuff的伤害计算

#### 技能效果测试
- `TestSkillManager_ApplySkillEffects_Charge` - 测试冲锋技能效果（怒气获得、眩晕）
- `TestSkillManager_ApplySkillEffects_Bloodthirst` - 测试嗜血技能效果（恢复生命值）

#### 技能冷却测试
- `TestSkillManager_TickCooldowns` - 测试技能冷却时间减少

### 4. 战士技能集成测试 (`warrior_skills_integration_test.go`)

#### 被动技能特殊效果测试
- `TestPassiveSkill_BloodCraze` - 测试血之狂热（攻击回血）
- `TestPassiveSkill_Revenge` - 测试复仇（概率反击）
- `TestPassiveSkill_Unbreakable` - 测试坚韧不拔（保命效果）
- `TestPassiveSkill_ShieldReflection` - 测试盾牌反射被动（反射伤害）

#### 被动技能怒气管理测试
- `TestPassiveSkill_AngerManagement` - 测试愤怒掌握（怒气获得加成）
- `TestPassiveSkill_WarMachine` - 测试战争机器（击杀回怒）

#### 技能效果测试
- `TestSkill_LastStand` - 测试破釜沉舟（立即恢复）
- `TestSkill_UnbreakableBarrier` - 测试不灭壁垒（护盾系统）
- `TestSkill_ShieldReflection` - 测试盾牌反射技能（反射伤害）

#### 敌人Debuff测试
- `TestEnemyDebuff_DemoralizingShout` - 测试挫志怒吼（降低敌人攻击力）
- `TestEnemyDebuff_Whirlwind` - 测试旋风斩（降低敌人防御）
- `TestEnemyDebuff_MortalStrike` - 测试致死打击（降低治疗效果）

## 运行测试

### 运行所有测试
```bash
cd server
go test ./internal/game -v
```

### 运行特定测试
```bash
# 运行BuffManager测试
go test ./internal/game -v -run TestBuffManager

# 运行被动技能测试
go test ./internal/game -v -run TestPassiveSkill

# 运行技能管理器测试
go test ./internal/game -v -run TestSkillManager

# 运行集成测试
go test ./internal/game -v -run TestSkill_
```

### 运行测试并生成覆盖率报告
```bash
go test ./internal/game -cover
go test ./internal/game -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试原则

1. **独立性**：每个测试都是独立的，不依赖其他测试的执行顺序
2. **可重复性**：测试结果应该一致，不依赖外部状态
3. **完整性**：测试覆盖主要功能和边界情况
4. **清晰性**：测试名称和注释清晰说明测试内容
5. **区分性**：测试代码使用 `_test.go` 后缀，与正常代码良好区分

## 注意事项

- 测试使用 `testify/assert` 进行断言
- 测试数据使用模拟对象，不依赖真实数据库
- 概率性测试（如复仇）使用多次尝试确保触发
- 测试中未使用的变量使用 `_` 避免编译警告


























