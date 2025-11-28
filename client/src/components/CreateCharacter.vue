<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGameStore } from '../stores/game'

const game = useGameStore()

const name = ref('')
const selectedRace = ref('')
const selectedClass = ref('')

const races = [
  { id: 'human', name: '人类', faction: 'alliance' },
  { id: 'dwarf', name: '矮人', faction: 'alliance' },
  { id: 'nightelf', name: '暗夜精灵', faction: 'alliance' },
  { id: 'gnome', name: '侏儒', faction: 'alliance' },
  { id: 'orc', name: '兽人', faction: 'horde' },
  { id: 'undead', name: '亡灵', faction: 'horde' },
  { id: 'tauren', name: '牛头人', faction: 'horde' },
  { id: 'troll', name: '巨魔', faction: 'horde' },
]

const classes = [
  { id: 'warrior', name: '战士', desc: '近战 | 力量型' },
  { id: 'mage', name: '法师', desc: '远程 | 智力型' },
  { id: 'rogue', name: '盗贼', desc: '近战 | 敏捷型' },
  { id: 'priest', name: '牧师', desc: '治疗 | 精神型' },
]

const canCreate = computed(() => {
  return name.value.trim() && selectedRace.value && selectedClass.value
})

async function handleCreate() {
  if (!canCreate.value) return
  await game.createCharacter(name.value.trim(), selectedRace.value, selectedClass.value)
}
</script>

<template>
  <div class="create-character">
    <h1 class="create-title">═══ 创建角色 ═══</h1>

    <!-- 角色名 -->
    <div class="create-section">
      <label class="create-label">> 输入你的名字:</label>
      <input 
        v-model="name"
        type="text" 
        class="create-input" 
        placeholder="请输入角色名..."
        maxlength="12"
      />
    </div>

    <!-- 种族选择 -->
    <div class="create-section">
      <label class="create-label">> 选择种族:</label>
      <div class="option-grid">
        <button
          v-for="race in races"
          :key="race.id"
          :class="['option-btn', { selected: selectedRace === race.id }]"
          @click="selectedRace = race.id"
        >
          {{ race.name }}
          <span style="font-size: 12px; display: block; opacity: 0.6">
            {{ race.faction === 'alliance' ? '[联盟]' : '[部落]' }}
          </span>
        </button>
      </div>
    </div>

    <!-- 职业选择 -->
    <div class="create-section">
      <label class="create-label">> 选择职业:</label>
      <div class="option-grid">
        <button
          v-for="cls in classes"
          :key="cls.id"
          :class="['option-btn', { selected: selectedClass === cls.id }]"
          @click="selectedClass = cls.id"
        >
          {{ cls.name }}
          <span style="font-size: 12px; display: block; opacity: 0.6">
            {{ cls.desc }}
          </span>
        </button>
      </div>
    </div>

    <!-- 开始按钮 -->
    <button 
      class="start-btn" 
      :disabled="!canCreate || game.isLoading"
      @click="handleCreate"
    >
      {{ game.isLoading ? '创建中...' : '[ 开始冒险 ]' }}
    </button>
  </div>
</template>






