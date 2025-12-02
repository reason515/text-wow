<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCharacterStore } from '@/stores/character'
import { useAuthStore } from '@/stores/auth'
import type { Race, Class } from '@/types/game'
import { CLASS_COLORS, getClassColorClass } from '@/types/game'

const emit = defineEmits<{
  created: []
  back: []
}>()

const charStore = useCharacterStore()
const authStore = useAuthStore()

// 状态
const step = ref(1) // 1=选阵营, 2=选种族, 3=选职业, 4=命名
const selectedFaction = ref<'alliance' | 'horde' | null>(null)
const selectedRace = ref<Race | null>(null)
const selectedClass = ref<Class | null>(null)
const characterName = ref('')

// 计算属性
const availableRaces = computed(() => {
  if (!selectedFaction.value) return []
  return selectedFaction.value === 'alliance' 
    ? charStore.allianceRaces 
    : charStore.hordeRaces
})

const canProceed = computed(() => {
  switch (step.value) {
    case 1: return !!selectedFaction.value
    case 2: return !!selectedRace.value
    case 3: return !!selectedClass.value
    case 4: return characterName.value.length >= 2 && characterName.value.length <= 32
    default: return false
  }
})

const factionColors = {
  alliance: { primary: '#4a90d9', secondary: '#1a5aa0' },
  horde: { primary: '#c41e3a', secondary: '#8c1a2c' }
}

const roleIcons: Record<string, string> = {
  tank: '🛡️',
  healer: '💚',
  dps: '⚔️',
  hybrid: '🔄'
}

const resourceIcons: Record<string, string> = {
  mana: '💙',
  rage: '❤️',
  energy: '💛'
}

// 方法
function selectFaction(faction: 'alliance' | 'horde') {
  selectedFaction.value = faction
  selectedRace.value = null
  selectedClass.value = null
  step.value = 2
}

function selectRace(race: Race) {
  selectedRace.value = race
  step.value = 3
}

function selectClass(cls: Class) {
  selectedClass.value = cls
  step.value = 4
}

function goBack() {
  if (step.value > 1) {
    step.value--
    if (step.value === 1) {
      selectedFaction.value = null
      selectedRace.value = null
      selectedClass.value = null
    } else if (step.value === 2) {
      selectedRace.value = null
      selectedClass.value = null
    } else if (step.value === 3) {
      selectedClass.value = null
    }
  } else {
    emit('back')
  }
}

async function createCharacter() {
  if (!selectedRace.value || !selectedClass.value || !characterName.value) return

  const char = await charStore.createCharacter({
    name: characterName.value,
    raceId: selectedRace.value.id,
    classId: selectedClass.value.id,
  })

  if (char) {
    emit('created')
  }
}

// 初始化
onMounted(async () => {
  if (charStore.races.length === 0) {
    await charStore.fetchRaces()
  }
  if (charStore.classes.length === 0) {
    await charStore.fetchClasses()
  }
})
</script>

<template>
  <div class="create-character">
    <!-- 顶部信息 -->
    <div class="header">
      <div class="back-btn" @click="goBack">
        ← {{ step === 1 ? '返回' : '上一步' }}
      </div>
      <div class="step-indicator">
        步骤 {{ step }}/4
      </div>
    </div>

    <!-- 步骤1: 选择阵营 -->
    <div v-if="step === 1" class="faction-select">
      <h2>选择你的阵营</h2>
      <p class="hint">为了艾泽拉斯，你将为谁而战？</p>
      
      <div class="faction-options">
        <div 
          class="faction-card alliance"
          :class="{ selected: selectedFaction === 'alliance' }"
          @click="selectFaction('alliance')"
        >
          <div class="faction-icon">🦁</div>
          <div class="faction-name">联盟</div>
          <div class="faction-desc">荣耀、正义、秩序</div>
          <div class="faction-races">人类 · 矮人 · 暗夜精灵 · 侏儒</div>
        </div>

        <div 
          class="faction-card horde"
          :class="{ selected: selectedFaction === 'horde' }"
          @click="selectFaction('horde')"
        >
          <div class="faction-icon">🐺</div>
          <div class="faction-name">部落</div>
          <div class="faction-desc">力量、荣誉、自由</div>
          <div class="faction-races">兽人 · 亡灵 · 牛头人 · 巨魔</div>
        </div>
      </div>
    </div>

    <!-- 步骤2: 选择种族 -->
    <div v-if="step === 2" class="race-select">
      <h2>选择你的种族</h2>
      <p class="hint">{{ selectedFaction === 'alliance' ? '联盟' : '部落' }}的勇士们</p>
      
      <div class="race-grid">
        <div 
          v-for="race in availableRaces" 
          :key="race.id"
          class="race-card"
          :class="{ selected: selectedRace?.id === race.id }"
          @click="selectRace(race)"
        >
          <div class="race-name">{{ race.name }}</div>
          <div class="race-desc">{{ race.description }}</div>
          <div class="race-bonuses">
            <span v-if="race.strengthBase" class="bonus str">+{{ race.strengthBase }} 力量</span>
            <span v-if="race.agilityBase" class="bonus agi">+{{ race.agilityBase }} 敏捷</span>
            <span v-if="race.intellectBase" class="bonus int">+{{ race.intellectBase }} 智力</span>
            <span v-if="race.staminaBase" class="bonus sta">+{{ race.staminaBase }} 耐力</span>
            <span v-if="race.spiritBase" class="bonus spi">+{{ race.spiritBase }} 精神</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤3: 选择职业 -->
    <div v-if="step === 3" class="class-select">
      <h2>选择你的职业</h2>
      <p class="hint">{{ selectedRace?.name }}可以成为...</p>
      
      <div class="class-grid">
        <div 
          v-for="cls in charStore.classes" 
          :key="cls.id"
          class="class-card"
          :class="[{ selected: selectedClass?.id === cls.id }, getClassColorClass(cls.id)]"
          :style="{ '--class-color': CLASS_COLORS[cls.id] || '#33ff33' }"
          @click="selectClass(cls)"
        >
          <div class="class-header">
            <span class="class-role">{{ roleIcons[cls.combatRole] || '⚔️' }}</span>
            <span class="class-name" :style="{ color: CLASS_COLORS[cls.id] }">{{ cls.name }}</span>
            <span class="class-resource">{{ resourceIcons[cls.resourceType] || '💙' }}</span>
          </div>
          <div class="class-desc">{{ cls.description }}</div>
          <div class="class-info">
            <span class="class-tag role">{{ cls.combatRole === 'tank' ? '坦克' : cls.combatRole === 'healer' ? '治疗' : cls.combatRole === 'hybrid' ? '混合' : '输出' }}</span>
            <span class="class-tag resource">{{ cls.resourceType === 'mana' ? '法力' : cls.resourceType === 'rage' ? '怒气' : '能量' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 步骤4: 命名 -->
    <div v-if="step === 4" class="name-input">
      <h2>为你的角色命名</h2>
      
      <div class="character-preview">
        <div class="preview-faction" :class="selectedFaction">
          {{ selectedFaction === 'alliance' ? '联盟' : '部落' }}
        </div>
        <div class="preview-info">
          <span>{{ selectedRace?.name }}</span>
          <span class="preview-sep">·</span>
          <span :style="{ color: CLASS_COLORS[selectedClass?.id || ''], textShadow: `0 0 8px ${CLASS_COLORS[selectedClass?.id || '']}` }">
            {{ selectedClass?.name }}
          </span>
        </div>
      </div>

      <div class="name-form">
        <input 
          v-model="characterName" 
          type="text" 
          placeholder="输入角色名称 (2-32字符)"
          maxlength="32"
          autofocus
        />
        <div class="name-hint">
          角色名称将在整个艾泽拉斯可见
        </div>
      </div>

      <div class="error-message" v-if="charStore.error">
        <span class="error-icon">⚠</span> {{ charStore.error }}
      </div>

      <button 
        class="create-btn"
        :disabled="!canProceed || charStore.loading"
        @click="createCharacter"
      >
        <span v-if="charStore.loading">创建中...</span>
        <span v-else>创建角色</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.create-character {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.back-btn {
  color: var(--terminal-cyan);
  cursor: pointer;
  transition: color 0.3s;
}

.back-btn:hover {
  color: var(--terminal-green);
}

.step-indicator {
  color: var(--terminal-gray);
  font-size: 12px;
}

h2 {
  color: var(--terminal-gold);
  text-align: center;
  margin-bottom: 10px;
  font-size: 18px;
}

.hint {
  color: var(--terminal-gray);
  text-align: center;
  margin-bottom: 30px;
  font-size: 12px;
}

/* 阵营选择 */
.faction-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.faction-card {
  border: 2px solid var(--terminal-gray);
  padding: 30px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.faction-card.alliance {
  border-color: #4a90d9;
}

.faction-card.alliance:hover,
.faction-card.alliance.selected {
  background: rgba(74, 144, 217, 0.1);
  box-shadow: 0 0 20px rgba(74, 144, 217, 0.3);
}

.faction-card.horde {
  border-color: #c41e3a;
}

.faction-card.horde:hover,
.faction-card.horde.selected {
  background: rgba(196, 30, 58, 0.1);
  box-shadow: 0 0 20px rgba(196, 30, 58, 0.3);
}

.faction-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.faction-name {
  font-size: 18px;
  margin-bottom: 10px;
}

.faction-card.alliance .faction-name { color: #4a90d9; }
.faction-card.horde .faction-name { color: #c41e3a; }

.faction-desc {
  color: var(--terminal-gray);
  font-size: 12px;
  margin-bottom: 15px;
}

.faction-races {
  color: var(--terminal-cyan);
  font-size: 12px;
}

/* 种族和职业网格 */
.race-grid,
.class-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.race-card,
.class-card {
  border: 1px solid var(--terminal-gray);
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.race-card:hover,
.race-card.selected {
  border-color: var(--terminal-green);
  background: rgba(0, 255, 0, 0.05);
}

.class-card:hover,
.class-card.selected {
  border-color: var(--class-color, var(--terminal-green));
  background: rgba(255, 255, 255, 0.05);
  box-shadow: 0 0 15px color-mix(in srgb, var(--class-color, #33ff33) 30%, transparent);
}

.race-name {
  color: var(--terminal-green);
  font-size: 14px;
  margin-bottom: 8px;
}

.class-name {
  font-size: 14px;
  margin-bottom: 8px;
  text-shadow: 0 0 8px currentColor;
}

.race-desc,
.class-desc {
  color: var(--terminal-gray);
  font-size: 12px;
  margin-bottom: 10px;
  line-height: 1.4;
}

.race-bonuses {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.bonus {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.1);
}

.bonus.str { color: #ff6b6b; }
.bonus.agi { color: #69db7c; }
.bonus.int { color: #74c0fc; }
.bonus.sta { color: #ffd43b; }
.bonus.spi { color: #da77f2; }

.class-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.class-role {
  font-size: 16px;
}

.class-resource {
  margin-left: auto;
  font-size: 14px;
}

.class-info {
  display: flex;
  gap: 8px;
}

.class-tag {
  font-size: 10px;
  padding: 2px 8px;
  border: 1px solid var(--terminal-gray);
}

.class-tag.role { color: var(--terminal-cyan); }
.class-tag.resource { color: var(--terminal-purple); }

/* 命名 */
.character-preview {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid var(--terminal-gray);
}

.preview-faction {
  font-size: 14px;
  margin-bottom: 10px;
}

.preview-faction.alliance { color: #4a90d9; }
.preview-faction.horde { color: #c41e3a; }

.preview-info {
  color: var(--terminal-green);
  font-size: 14px;
}

.preview-sep {
  color: var(--terminal-gray);
  margin: 0 8px;
}

.name-form {
  margin-bottom: 20px;
}

.name-form input {
  width: 100%;
  padding: 15px;
  background: rgba(0, 0, 0, 0.5);
  border: 2px solid var(--terminal-gray);
  color: var(--terminal-green);
  font-family: inherit;
  font-size: 14px;
  text-align: center;
}

.name-form input:focus {
  outline: none;
  border-color: var(--terminal-green);
}

.name-hint {
  color: var(--terminal-gray);
  font-size: 12px;
  text-align: center;
  margin-top: 10px;
}

.error-message {
  background: rgba(255, 0, 0, 0.1);
  border: 1px solid var(--terminal-red);
  color: var(--terminal-red);
  padding: 10px 15px;
  margin-bottom: 20px;
  font-size: 12px;
  text-align: center;
}

.create-btn {
  width: 100%;
  padding: 15px;
  background: transparent;
  border: 2px solid var(--terminal-gold);
  color: var(--terminal-gold);
  font-family: inherit;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.create-btn:hover:not(:disabled) {
  background: var(--terminal-gold);
  color: var(--terminal-bg);
}

.create-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 600px) {
  .faction-options,
  .race-grid,
  .class-grid {
    grid-template-columns: 1fr;
  }
}
</style>
