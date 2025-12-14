<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { get, post, put, del } from '@/api/client'
import type { 
  BattleStrategy, 
  ConditionalRule, 
  RuleCondition,
  ConditionTypeInfo,
  TargetPriorityInfo,
  StrategyTemplate,
  AutoTargetSettings
} from '@/types/game'

const props = defineProps<{
  characterId: number
  characterSkills?: Array<{ skillId: string; skill?: { name: string } }>
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

// 状态
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const strategies = ref<BattleStrategy[]>([])
const currentStrategy = ref<BattleStrategy | null>(null)
const activeTab = ref('rules')  // rules, skills, target, advanced
const conditionTypes = ref<ConditionTypeInfo[]>([])
const targetPriorities = ref<TargetPriorityInfo[]>([])
const templates = ref<StrategyTemplate[]>([])

// 确保策略字段有默认值（防御性处理后端返回 null 的情况）
function ensureStrategyDefaults(strategy: BattleStrategy): BattleStrategy {
  return {
    ...strategy,
    skillPriority: strategy.skillPriority || [],
    conditionalRules: strategy.conditionalRules || [],
    skillTargetOverrides: strategy.skillTargetOverrides || {},
    reservedSkills: strategy.reservedSkills || [],
    autoTargetSettings: strategy.autoTargetSettings || {
      positionalAutoOptimize: true,
      executeAutoTarget: true,
      healAutoTarget: true
    }
  }
}

const addSkillSelectEl = ref<HTMLSelectElement | null>(null)

function addSkillToPriority() {
  const strategy = currentStrategy.value
  const select = addSkillSelectEl.value
  if (!strategy || !select) {
    return
  }

  const skillId = select.value
  if (!skillId) {
    return
  }

  if (!Array.isArray(strategy.skillPriority)) {
    return
  }

  if (strategy.skillPriority.includes(skillId)) {
    return
  }

  strategy.skillPriority.push(skillId)
  select.value = ''
}

// 新建策略弹窗
const showNewDialog = ref(false)
const newStrategyName = ref('')
const newStrategyTemplate = ref('')

// 加载策略列表
async function loadStrategies() {
  loading.value = true
  error.value = ''
  try {
    const res = await get<BattleStrategy[]>(`/characters/${props.characterId}/strategies`)
    if (res.success && res.data) {
      // 确保所有策略字段有默认值
      strategies.value = res.data.map(ensureStrategyDefaults)
      // 选中激活的策略
      const active = strategies.value.find(s => s.isActive)
      if (active) {
        currentStrategy.value = active
      } else if (strategies.value.length > 0) {
        currentStrategy.value = strategies.value[0]
      }
    } else {
      error.value = res.error || '加载失败'
    }
  } catch (e) {
    error.value = '加载策略失败'
  } finally {
    loading.value = false
  }
}

// 加载条件类型
async function loadConditionTypes() {
  try {
    const res = await get<{ conditionTypes: ConditionTypeInfo[]; targetPriorities: TargetPriorityInfo[] }>('/strategy-condition-types')
    if (res.success && res.data) {
      conditionTypes.value = res.data.conditionTypes
      targetPriorities.value = res.data.targetPriorities
    }
  } catch (e) {
    console.error('Failed to load condition types', e)
  }
}

// 加载模板
async function loadTemplates() {
  try {
    const res = await get<{ templateList: StrategyTemplate[] }>('/strategy-templates')
    if (res.success && res.data) {
      templates.value = res.data.templateList
    }
  } catch (e) {
    console.error('Failed to load templates', e)
  }
}

// 创建策略
async function createStrategy() {
  if (!newStrategyName.value.trim()) {
    error.value = '请输入策略名称'
    return
  }
  saving.value = true
  try {
    const res = await post<BattleStrategy>(`/characters/${props.characterId}/strategies`, {
      characterId: props.characterId,
      name: newStrategyName.value.trim(),
      fromTemplate: newStrategyTemplate.value || undefined
    })
    if (res.success && res.data) {
      const normalizedStrategy = ensureStrategyDefaults(res.data)
      strategies.value.push(normalizedStrategy)
      currentStrategy.value = normalizedStrategy
      showNewDialog.value = false
      newStrategyName.value = ''
      newStrategyTemplate.value = ''
    } else {
      error.value = res.error || '创建失败'
    }
  } catch (e) {
    error.value = '创建策略失败'
  } finally {
    saving.value = false
  }
}

// 保存策略
async function saveStrategy() {
  if (!currentStrategy.value) return
  saving.value = true
  error.value = ''
  try {
    const res = await put<BattleStrategy>(`/strategies/${currentStrategy.value.id}`, {
      name: currentStrategy.value.name,
      skillPriority: currentStrategy.value.skillPriority,
      conditionalRules: currentStrategy.value.conditionalRules,
      targetPriority: currentStrategy.value.targetPriority,
      skillTargetOverrides: currentStrategy.value.skillTargetOverrides,
      resourceThreshold: currentStrategy.value.resourceThreshold,
      reservedSkills: currentStrategy.value.reservedSkills,
      autoTargetSettings: currentStrategy.value.autoTargetSettings
    })
    if (res.success) {
      error.value = ''
      // 显示成功提示
      const idx = strategies.value.findIndex(s => s.id === currentStrategy.value!.id)
      if (idx >= 0 && res.data) {
        strategies.value[idx] = res.data
        currentStrategy.value = res.data
      }
    } else {
      error.value = res.error || '保存失败'
    }
  } catch (e) {
    error.value = '保存策略失败'
  } finally {
    saving.value = false
  }
}

// 删除策略
async function deleteStrategy() {
  if (!currentStrategy.value) return
  if (!confirm('确定要删除这个策略吗？')) return
  
  try {
    const res = await del(`/strategies/${currentStrategy.value.id}`)
    if (res.success) {
      strategies.value = strategies.value.filter(s => s.id !== currentStrategy.value!.id)
      currentStrategy.value = strategies.value[0] || null
    } else {
      error.value = res.error || '删除失败'
    }
  } catch (e) {
    error.value = '删除策略失败'
  }
}

// 激活策略
async function activateStrategy() {
  if (!currentStrategy.value) return
  try {
    const res = await post(`/strategies/${currentStrategy.value.id}/activate`)
    if (res.success) {
      // 更新本地状态
      strategies.value.forEach(s => s.isActive = false)
      currentStrategy.value.isActive = true
    } else {
      error.value = res.error || '激活失败'
    }
  } catch (e) {
    error.value = '激活策略失败'
  }
}

// 添加条件规则
function addRule() {
  if (!currentStrategy.value) return
  const newRule: ConditionalRule = {
    id: `rule_${Date.now()}`,
    priority: currentStrategy.value.conditionalRules.length + 1,
    enabled: true,
    condition: {
      type: 'self_hp_percent',
      operator: '<',
      value: 30
    },
    action: {
      type: 'use_skill',
      skillId: ''
    }
  }
  currentStrategy.value.conditionalRules.push(newRule)
}

// 删除条件规则
function removeRule(index: number) {
  if (!currentStrategy.value) return
  currentStrategy.value.conditionalRules.splice(index, 1)
  // 重新计算优先级
  currentStrategy.value.conditionalRules.forEach((rule, idx) => {
    rule.priority = idx + 1
  })
}

// 移动规则
function moveRule(index: number, direction: 'up' | 'down') {
  if (!currentStrategy.value) return
  const rules = currentStrategy.value.conditionalRules
  const newIndex = direction === 'up' ? index - 1 : index + 1
  if (newIndex < 0 || newIndex >= rules.length) return
  
  [rules[index], rules[newIndex]] = [rules[newIndex], rules[index]]
  rules.forEach((rule, idx) => {
    rule.priority = idx + 1
  })
}

// 获取条件类型名称
function getConditionTypeName(type: string): string {
  const ct = conditionTypes.value.find(c => c.type === type)
  return ct?.name || type
}

// 获取目标策略名称
function getTargetPriorityName(value: string): string {
  const tp = targetPriorities.value.find(t => t.value === value)
  return tp?.label || value
}

// 获取技能名称
function getSkillName(skillId: string): string {
  const skill = props.characterSkills?.find(s => s.skillId === skillId)
  return skill?.skill?.name || skillId || '选择技能'
}

// 条件类型分类
const conditionCategories = computed(() => {
  const categories: Record<string, ConditionTypeInfo[]> = {}
  conditionTypes.value.forEach(ct => {
    if (!categories[ct.category]) {
      categories[ct.category] = []
    }
    categories[ct.category].push(ct)
  })
  return categories
})

const categoryNames: Record<string, string> = {
  self: '自身状态',
  enemy: '敌人状态',
  ally: '队友状态',
  battle: '战斗状态'
}

onMounted(() => {
  loadStrategies()
  loadConditionTypes()
  loadTemplates()
})
</script>

<template>
  <div class="strategy-editor">
    <!-- 头部 -->
    <div class="editor-header">
      <h3>⚔️ 作战策略</h3>
      <button class="close-btn" @click="emit('close')">×</button>
    </div>

    <!-- 策略选择栏 -->
    <div class="strategy-selector">
      <select v-model="currentStrategy" class="strategy-select">
        <option v-for="s in strategies" :key="s.id" :value="s">
          {{ s.name }} {{ s.isActive ? '✓' : '' }}
        </option>
      </select>
      <button class="btn-new" @click="showNewDialog = true" title="新建策略">+</button>
      <button 
        class="btn-activate" 
        @click="activateStrategy" 
        :disabled="!currentStrategy || currentStrategy.isActive"
        title="激活此策略"
      >
        激活
      </button>
      <button 
        class="btn-delete" 
        @click="deleteStrategy" 
        :disabled="!currentStrategy"
        title="删除策略"
      >
        删除
      </button>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-message">{{ error }}</div>

    <!-- 加载中 -->
    <div v-if="loading" class="loading">加载中...</div>

    <!-- 策略编辑区 -->
    <div v-else-if="currentStrategy" class="editor-content">
      <!-- 标签页 -->
      <div class="tabs">
        <button 
          :class="['tab', { active: activeTab === 'rules' }]" 
          @click="activeTab = 'rules'"
        >
          条件规则
        </button>
        <button 
          :class="['tab', { active: activeTab === 'skills' }]" 
          @click="activeTab = 'skills'"
        >
          技能顺序
        </button>
        <button 
          :class="['tab', { active: activeTab === 'target' }]" 
          @click="activeTab = 'target'"
        >
          目标选择
        </button>
        <button 
          :class="['tab', { active: activeTab === 'advanced' }]" 
          @click="activeTab = 'advanced'"
        >
          高级设置
        </button>
      </div>

      <!-- 条件规则标签页 -->
      <div v-show="activeTab === 'rules'" class="tab-content">
        <div class="section-hint">
          💡 条件规则按优先级从上到下执行，满足条件时使用对应技能
        </div>
        
        <div class="rules-list">
          <div 
            v-for="(rule, index) in currentStrategy.conditionalRules" 
            :key="rule.id"
            class="rule-card"
          >
            <div class="rule-header">
              <label class="rule-checkbox">
                <input type="checkbox" v-model="rule.enabled">
                <span>#{{ index + 1 }}</span>
              </label>
              <div class="rule-actions">
                <button @click="moveRule(index, 'up')" :disabled="index === 0">↑</button>
                <button @click="moveRule(index, 'down')" :disabled="index === currentStrategy.conditionalRules.length - 1">↓</button>
                <button @click="removeRule(index)" class="btn-remove">×</button>
              </div>
            </div>
            
            <div class="rule-body">
              <div class="rule-condition">
                <span>当</span>
                <select v-model="rule.condition.type" class="condition-type">
                  <optgroup v-for="(types, category) in conditionCategories" :key="category" :label="categoryNames[category] || category">
                    <option v-for="ct in types" :key="ct.type" :value="ct.type">
                      {{ ct.name }}
                    </option>
                  </optgroup>
                </select>
                <select v-model="rule.condition.operator" class="condition-operator">
                  <option value="<">&lt;</option>
                  <option value=">">&gt;</option>
                  <option value="<=">&le;</option>
                  <option value=">=">&ge;</option>
                  <option value="=">=</option>
                </select>
                <input 
                  type="number" 
                  v-model.number="rule.condition.value" 
                  class="condition-value"
                  min="0"
                  max="100"
                >
                <span v-if="rule.condition.type.includes('percent')">%</span>
              </div>
              
              <div class="rule-action">
                <span>使用</span>
                <select v-model="rule.action.skillId" class="skill-select">
                  <option value="">选择技能</option>
                  <option v-for="skill in characterSkills" :key="skill.skillId" :value="skill.skillId">
                    {{ skill.skill?.name || skill.skillId }}
                  </option>
                  <option value="__normal_attack__">普通攻击</option>
                </select>
              </div>
            </div>
          </div>
        </div>

        <button class="btn-add-rule" @click="addRule">+ 添加条件规则</button>
      </div>

      <!-- 技能顺序标签页 -->
      <div v-show="activeTab === 'skills'" class="tab-content">
        <div class="section-hint">
          💡 当没有条件规则触发时，按以下顺序选择可用技能
        </div>
        
        <div class="skills-priority">
          <div 
            v-for="(skillId, index) in currentStrategy.skillPriority" 
            :key="skillId"
            class="skill-item"
          >
            <span class="skill-order">{{ index + 1 }}</span>
            <span class="skill-name">{{ getSkillName(skillId) }}</span>
            <button @click="currentStrategy.skillPriority.splice(index, 1)" class="btn-remove">×</button>
          </div>
        </div>

        <div class="add-skill">
          <select ref="addSkillSelectEl" class="skill-select">
            <option value="">添加技能到优先级列表</option>
            <option 
              v-for="skill in characterSkills" 
              :key="skill.skillId" 
              :value="skill.skillId"
              :disabled="currentStrategy.skillPriority.includes(skill.skillId)"
            >
              {{ skill.skill?.name || skill.skillId }}
            </option>
          </select>
          <button @click="addSkillToPriority">添加</button>
        </div>
      </div>

      <!-- 目标选择标签页 -->
      <div v-show="activeTab === 'target'" class="tab-content">
        <div class="section-hint">
          💡 配置攻击/治疗目标的选择策略
        </div>

        <div class="target-section">
          <h4>默认目标策略</h4>
          <div class="target-options">
            <label v-for="tp in targetPriorities" :key="tp.value" class="target-option">
              <input 
                type="radio" 
                :value="tp.value" 
                v-model="currentStrategy.targetPriority"
              >
              {{ tp.label }}
            </label>
          </div>
        </div>

        <div class="target-section">
          <h4>智能目标</h4>
          <div class="auto-target-options">
            <label class="checkbox-option">
              <input type="checkbox" v-model="currentStrategy.autoTargetSettings.positionalAutoOptimize">
              位置技能自动优化（顺劈斩等技能自动选择能波及最多敌人的位置）
            </label>
            <label class="checkbox-option">
              <input type="checkbox" v-model="currentStrategy.autoTargetSettings.executeAutoTarget">
              斩杀技能自动找残血（斩杀等终结技能自动选择HP低于20%的敌人）
            </label>
            <label class="checkbox-option">
              <input type="checkbox" v-model="currentStrategy.autoTargetSettings.healAutoTarget">
              治疗技能自动选低血队友（治疗技能自动选择HP最低的队友）
            </label>
          </div>
        </div>
      </div>

      <!-- 高级设置标签页 -->
      <div v-show="activeTab === 'advanced'" class="tab-content">
        <div class="section-hint">
          💡 资源管理和保留技能设置
        </div>

        <div class="advanced-section">
          <h4>资源阈值</h4>
          <div class="threshold-setting">
            <input 
              type="range" 
              v-model.number="currentStrategy.resourceThreshold" 
              min="0" 
              max="100" 
              class="threshold-slider"
            >
            <span class="threshold-value">{{ currentStrategy.resourceThreshold }}</span>
          </div>
          <p class="setting-hint">低于此值时优先使用普通攻击积攒资源</p>
        </div>

        <div class="advanced-section">
          <h4>策略名称</h4>
          <input 
            type="text" 
            v-model="currentStrategy.name" 
            class="strategy-name-input"
            maxlength="32"
          >
        </div>
      </div>
    </div>

    <!-- 无策略提示 -->
    <div v-else-if="!loading" class="no-strategy">
      <p>暂无策略，点击 + 创建第一个策略</p>
    </div>

    <!-- 底部按钮 -->
    <div class="editor-footer">
      <button class="btn-save" @click="saveStrategy" :disabled="saving || !currentStrategy">
        {{ saving ? '保存中...' : '保存策略' }}
      </button>
      <button class="btn-cancel" @click="emit('close')">关闭</button>
    </div>

    <!-- 新建策略弹窗 -->
    <div v-if="showNewDialog" class="dialog-overlay" @click.self="showNewDialog = false">
      <div class="dialog">
        <h4>新建策略</h4>
        <div class="dialog-body">
          <div class="form-group">
            <label>策略名称</label>
            <input 
              type="text" 
              v-model="newStrategyName" 
              placeholder="输入策略名称"
              maxlength="32"
            >
          </div>
          <div class="form-group">
            <label>从模板创建 (可选)</label>
            <select v-model="newStrategyTemplate">
              <option value="">空白策略</option>
              <option v-for="t in templates" :key="t.id" :value="t.id">
                {{ t.name }} - {{ t.description }}
              </option>
            </select>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-save" @click="createStrategy" :disabled="saving">
            {{ saving ? '创建中...' : '创建' }}
          </button>
          <button class="btn-cancel" @click="showNewDialog = false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.strategy-editor {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  max-width: 95vw;
  max-height: 85vh;
  background: #1a1a2e;
  border: 1px solid #33ff33;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  z-index: 1000;
  font-family: 'Consolas', 'Monaco', monospace;
  color: #33ff33;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #33ff33;
}

.editor-header h3 {
  margin: 0;
  font-size: 16px;
}

.close-btn {
  background: none;
  border: none;
  color: #ff4444;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.strategy-selector {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid #333;
}

.strategy-select {
  flex: 1;
  background: #0a0a14;
  border: 1px solid #33ff33;
  color: #33ff33;
  padding: 6px;
  border-radius: 4px;
}

.btn-new, .btn-activate, .btn-delete {
  background: #0a0a14;
  border: 1px solid #33ff33;
  color: #33ff33;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-new:hover, .btn-activate:hover {
  background: #1a3a1a;
}

.btn-delete {
  border-color: #ff4444;
  color: #ff4444;
}

.btn-delete:hover {
  background: #3a1a1a;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-message {
  background: #3a1a1a;
  color: #ff4444;
  padding: 8px 16px;
  margin: 8px 16px;
  border-radius: 4px;
}

.loading, .no-strategy {
  padding: 32px;
  text-align: center;
  color: #888;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  min-height: 300px;
}

.tabs {
  display: flex;
  border-bottom: 1px solid #333;
}

.tab {
  flex: 1;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #888;
  padding: 12px;
  cursor: pointer;
  font-family: inherit;
}

.tab:hover {
  color: #33ff33;
}

.tab.active {
  color: #33ff33;
  border-bottom-color: #33ff33;
}

.tab-content {
  padding: 16px;
}

.section-hint {
  color: #888;
  font-size: 12px;
  margin-bottom: 16px;
  padding: 8px;
  background: #0a0a14;
  border-radius: 4px;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.rule-card {
  background: #0a0a14;
  border: 1px solid #333;
  border-radius: 4px;
  padding: 12px;
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.rule-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rule-checkbox input {
  accent-color: #33ff33;
}

.rule-actions {
  display: flex;
  gap: 4px;
}

.rule-actions button {
  background: #1a1a2e;
  border: 1px solid #333;
  color: #888;
  padding: 2px 8px;
  cursor: pointer;
  border-radius: 2px;
}

.rule-actions button:hover:not(:disabled) {
  color: #33ff33;
  border-color: #33ff33;
}

.btn-remove {
  color: #ff4444 !important;
}

.btn-remove:hover {
  border-color: #ff4444 !important;
}

.rule-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rule-condition, .rule-action {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.condition-type {
  min-width: 120px;
}

.condition-operator {
  width: 60px;
}

.condition-value {
  width: 60px;
}

.skill-select {
  min-width: 150px;
}

select, input[type="text"], input[type="number"] {
  background: #1a1a2e;
  border: 1px solid #333;
  color: #33ff33;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: inherit;
}

select:focus, input:focus {
  outline: none;
  border-color: #33ff33;
}

.btn-add-rule {
  width: 100%;
  background: #0a0a14;
  border: 1px dashed #33ff33;
  color: #33ff33;
  padding: 12px;
  cursor: pointer;
  border-radius: 4px;
}

.btn-add-rule:hover {
  background: #1a3a1a;
}

.skills-priority {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.skill-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #0a0a14;
  padding: 8px 12px;
  border-radius: 4px;
}

.skill-order {
  color: #888;
  min-width: 24px;
}

.skill-name {
  flex: 1;
}

.add-skill {
  display: flex;
  gap: 8px;
}

.add-skill select {
  flex: 1;
}

.add-skill button {
  background: #0a0a14;
  border: 1px solid #33ff33;
  color: #33ff33;
  padding: 4px 12px;
  cursor: pointer;
  border-radius: 4px;
}

.target-section, .advanced-section {
  margin-bottom: 24px;
}

.target-section h4, .advanced-section h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #33ff33;
}

.target-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.target-option {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.target-option input {
  accent-color: #33ff33;
}

.auto-target-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.checkbox-option {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  cursor: pointer;
  font-size: 13px;
  line-height: 1.4;
}

.checkbox-option input {
  accent-color: #33ff33;
  margin-top: 3px;
}

.threshold-setting {
  display: flex;
  align-items: center;
  gap: 12px;
}

.threshold-slider {
  flex: 1;
  accent-color: #33ff33;
}

.threshold-value {
  min-width: 40px;
  text-align: right;
}

.setting-hint {
  color: #888;
  font-size: 12px;
  margin-top: 8px;
}

.strategy-name-input {
  width: 100%;
  padding: 8px;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 12px 16px;
  border-top: 1px solid #333;
}

.btn-save {
  background: #1a3a1a;
  border: 1px solid #33ff33;
  color: #33ff33;
  padding: 8px 24px;
  cursor: pointer;
  border-radius: 4px;
}

.btn-save:hover:not(:disabled) {
  background: #2a4a2a;
}

.btn-cancel {
  background: #1a1a2e;
  border: 1px solid #888;
  color: #888;
  padding: 8px 24px;
  cursor: pointer;
  border-radius: 4px;
}

.btn-cancel:hover {
  color: #33ff33;
  border-color: #33ff33;
}

/* 弹窗样式 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
}

.dialog {
  background: #1a1a2e;
  border: 1px solid #33ff33;
  border-radius: 8px;
  padding: 20px;
  min-width: 350px;
}

.dialog h4 {
  margin: 0 0 16px 0;
  color: #33ff33;
}

.dialog-body {
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 12px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  color: #888;
}

.form-group input, .form-group select {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 滚动条样式 */
.editor-content::-webkit-scrollbar {
  width: 8px;
}

.editor-content::-webkit-scrollbar-track {
  background: #0a0a14;
}

.editor-content::-webkit-scrollbar-thumb {
  background: #333;
  border-radius: 4px;
}

.editor-content::-webkit-scrollbar-thumb:hover {
  background: #33ff33;
}
</style>



