<script setup lang="ts">
import { onMounted } from 'vue'
import { useGameStore } from './stores/game'
import CreateCharacter from './components/CreateCharacter.vue'
import GameScreen from './components/GameScreen.vue'

const game = useGameStore()

onMounted(async () => {
  await game.fetchCharacter()
  if (game.hasCharacter) {
    await game.fetchBattleLogs()
    await game.fetchBattleStatus()
    await game.fetchZones()
    
    // 如果战斗正在进行，启动循环
    if (game.isRunning) {
      game.startBattleLoop()
    }
  }
})
</script>

<template>
  <div class="terminal-container">
    <div class="terminal-window">
      <!-- 标题栏 -->
      <div class="terminal-header">
        <span class="terminal-title">TEXT WoW v0.1 - 放置类文字RPG</span>
        <div class="terminal-controls">
          <button class="terminal-btn">─</button>
          <button class="terminal-btn">□</button>
          <button class="terminal-btn">×</button>
        </div>
      </div>

      <!-- 主内容 -->
      <div class="terminal-body">
        <CreateCharacter v-if="!game.hasCharacter" />
        <GameScreen v-else />
      </div>
    </div>

    <!-- CRT扫描线效果 -->
    <div class="crt-overlay"></div>
  </div>
</template>



