# 前端集成指南

本文档说明后端更新后，前端需要进行的调整。

## 主要变化

### 1. 速度计算公式变更 ⚠️ **需要调整**

**后端变化：**
- 速度计算公式从 `speed = 10 + (agility / 2)` 简化为 `speed = agility`
- 速度相等时随机确定出手顺序

**前端需要调整：**

#### 1.1 在角色详情面板中显示速度属性

**文件：** `client/src/components/GameScreen.vue`

**位置：** 角色详情面板的属性显示部分（约1960行附近）

**需要添加：**
```vue
<!-- 在敏捷属性后添加速度显示 -->
<div class="character-detail-stat">
  <span class="character-detail-stat-label">速度</span>
  <span class="character-detail-stat-value">{{ selectedCharacter.agility || 0 }}</span>
</div>
```

#### 1.2 更新敏捷属性的Tooltip说明

**文件：** `client/src/components/GameScreen.vue`

**位置：** `getPrimaryStatTooltip` 函数中的 `agility` case（约667行）

**需要修改：**
```typescript
case 'agility': {
  const physAtk = (agi * 0.2).toFixed(1)
  const critRate = ((char.physCritRate ?? (0.05 + agi / 20)) * 100).toFixed(1)
  const dodge = ((char.dodgeRate ?? (0.05 + agi / 20)) * 100).toFixed(1)
  const speed = agi // 速度 = 敏捷
  return [
    '敏捷',
    `- 每点提供 0.2 物理攻击 (当前贡献: +${physAtk})`,
    `- 物理暴击率 = 5% + 敏捷/20 (当前: ${critRate}%)`,
    `- 闪避率 = 5% + 敏捷/20 (当前: ${dodge}%)`,
    `- 速度 = 敏捷 (当前速度: ${speed})`  // 新增
  ].join('\n')
}
```

#### 1.3 添加速度属性的Tooltip处理

**文件：** `client/src/components/GameScreen.vue`

**位置：** 在 `handleAttrTooltip` 函数中添加对 `speed` 的处理

**需要添加：**
```typescript
// 在 handleAttrTooltip 函数中添加
case 'speed': {
  const agi = char.agility || 0
  return `速度 = 敏捷\n当前敏捷: ${agi}\n当前速度: ${agi}`
}
```

### 2. 地图/区域系统 ✅ **已支持**

前端已经实现了完整的地图系统：
- ✅ 地图选择面板 (`showZoneSelector`)
- ✅ 地图列表加载 (`loadZones`)
- ✅ 地图切换功能 (`selectZone`)
- ✅ 探索度显示
- ✅ 解锁条件显示

**无需调整**

### 3. 回合顺序显示（可选增强）

如果需要显示当前回合顺序，可以添加以下功能：

#### 3.1 在战斗状态中显示回合顺序

**文件：** `client/src/components/GameScreen.vue`

**位置：** 战斗状态显示区域

**可选添加：**
```vue
<!-- 回合顺序显示 -->
<div v-if="game.battleStatus?.turnOrder" class="turn-order">
  <div class="turn-order-label">回合顺序：</div>
  <div class="turn-order-list">
    <span 
      v-for="(participant, index) in game.battleStatus.turnOrder" 
      :key="index"
      class="turn-participant"
      :class="{ 'turn-current': index === 0 }"
    >
      {{ participant.type === 'character' ? participant.character?.name : participant.monster?.name }}
    </span>
  </div>
</div>
```

**注意：** 这需要后端API返回 `turnOrder` 数据。

### 4. 其他属性检查

#### 4.1 角色属性类型定义

**文件：** `client/src/types/game.ts`

检查 `Character` 接口是否包含 `speed` 字段（可选，因为速度现在等于敏捷）

#### 4.2 战斗日志格式

**文件：** `client/src/components/GameScreen.vue`

检查 `formatLogMessage` 函数是否需要更新以支持新的日志类型。

## 实施步骤

1. **立即需要：**
   - [ ] 更新敏捷属性的Tooltip，添加速度说明
   - [ ] 在角色详情面板中显示速度（可选，因为速度=敏捷）

2. **可选增强：**
   - [ ] 添加回合顺序显示
   - [ ] 添加速度相关的视觉提示

3. **测试验证：**
   - [ ] 验证角色属性显示正确
   - [ ] 验证地图切换功能正常
   - [ ] 验证战斗日志显示正常
   - [ ] 验证Buff/Debuff显示正常

## 后端API兼容性

### 已确认兼容的API：
- ✅ `/api/battle/status` - 战斗状态
- ✅ `/api/battle/zones` - 地图列表
- ✅ `/api/battle/change-zone` - 切换地图
- ✅ `/api/characters/:id/skills` - 角色技能
- ✅ `/api/characters/:id/passives` - 被动技能

### 可能需要更新的API：
- ⚠️ `/api/battle/status` - 如果需要返回 `turnOrder`，需要后端支持

## 注意事项

1. **速度计算：** 前端不需要单独计算速度，直接使用敏捷值即可
2. **向后兼容：** 如果后端返回了 `speed` 字段，前端应该优先使用，否则使用 `agility`
3. **地图系统：** 前端地图系统已经完整，只需要确保后端API正常返回数据

## 测试清单

- [ ] 角色详情面板显示正确
- [ ] 敏捷Tooltip包含速度说明
- [ ] 地图选择器正常工作
- [ ] 地图切换功能正常
- [ ] 探索度显示正确
- [ ] 战斗日志格式正确
- [ ] Buff/Debuff显示正常
- [ ] 被动技能显示正常











































