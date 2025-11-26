import type { Character, Monster, Skill, LogEntry, BattleStrategy } from './types'

let logIdCounter = 0

// 生成日志ID
function nextLogId(): number {
  return ++logIdCounter
}

// 创建日志条目
export function createLog(
  message: string, 
  type: LogEntry['type'] = 'system'
): LogEntry {
  return {
    id: nextLogId(),
    timestamp: Date.now(),
    message,
    type
  }
}

// 计算伤害
export function calculateDamage(
  attacker: { attack: number; level?: number },
  defender: { defense: number },
  skill?: Skill
): { damage: number; isCrit: boolean } {
  const baseDamage = attacker.attack
  const skillMultiplier = skill?.damage ?? 1
  const defense = defender.defense
  
  // 基础伤害公式
  let damage = Math.max(1, Math.floor(baseDamage * skillMultiplier - defense * 0.5))
  
  // 随机浮动 ±20%
  const variance = 0.2
  damage = Math.floor(damage * (1 + (Math.random() * variance * 2 - variance)))
  
  // 暴击判定（15%几率，2倍伤害）
  const isCrit = Math.random() < 0.15
  if (isCrit) {
    damage = Math.floor(damage * 2)
  }
  
  return { damage: Math.max(1, damage), isCrit }
}

// 选择可用技能
export function selectSkill(
  character: Character,
  strategy: BattleStrategy
): Skill | null {
  const availableSkills = character.skills.filter(skill => {
    // 检查冷却
    if (skill.currentCooldown > 0) return false
    // 检查法力
    if (character.combatStats.currentMp < skill.mpCost) return false
    // 治疗技能特殊处理
    if (skill.type === 'heal') {
      const hpPercent = (character.combatStats.currentHp / character.combatStats.maxHp) * 100
      return hpPercent < strategy.useHealAt
    }
    return true
  })
  
  if (availableSkills.length === 0) return null
  
  // 按策略优先级排序
  for (const skillId of strategy.skillPriority) {
    const skill = availableSkills.find(s => s.id === skillId)
    if (skill) return skill
  }
  
  // 默认返回第一个可用技能
  return availableSkills[0]
}

// 玩家回合
export function playerTurn(
  character: Character,
  monster: Monster,
  strategy: BattleStrategy
): LogEntry[] {
  const logs: LogEntry[] = []
  
  // 选择技能
  const skill = selectSkill(character, strategy)
  
  if (skill) {
    // 消耗法力
    character.combatStats.currentMp -= skill.mpCost
    
    if (skill.type === 'heal') {
      // 治疗
      const healAmount = Math.floor(Math.abs(skill.damage) * character.combatStats.attack * 0.5)
      const actualHeal = Math.min(
        healAmount,
        character.combatStats.maxHp - character.combatStats.currentHp
      )
      character.combatStats.currentHp += actualHeal
      logs.push(createLog(
        `你使用了 [${skill.name}]，恢复了 ${actualHeal} 点生命值`,
        'heal'
      ))
    } else {
      // 攻击
      const { damage, isCrit } = calculateDamage(
        { attack: character.combatStats.attack },
        { defense: monster.defense },
        skill
      )
      monster.currentHp -= damage
      
      const critText = isCrit ? ' (暴击!)' : ''
      logs.push(createLog(
        `你使用了 [${skill.name}] 对 ${monster.name} 造成 ${damage} 点伤害${critText}`,
        'damage'
      ))
    }
    
    // 设置冷却
    skill.currentCooldown = skill.cooldown
  } else {
    // 普通攻击
    const { damage, isCrit } = calculateDamage(
      { attack: character.combatStats.attack },
      { defense: monster.defense }
    )
    monster.currentHp -= damage
    
    const critText = isCrit ? ' (暴击!)' : ''
    logs.push(createLog(
      `你攻击了 ${monster.name} 造成 ${damage} 点伤害${critText}`,
      'damage'
    ))
  }
  
  return logs
}

// 怪物回合
export function monsterTurn(
  monster: Monster,
  character: Character
): LogEntry[] {
  const logs: LogEntry[] = []
  
  // 闪避判定（基于敏捷）
  const dodgeChance = Math.min(30, character.stats.agility * 0.5)
  if (Math.random() * 100 < dodgeChance) {
    logs.push(createLog(
      `${monster.name} 的攻击被你闪避了！`,
      'system'
    ))
    return logs
  }
  
  const { damage, isCrit } = calculateDamage(
    { attack: monster.attack },
    { defense: character.combatStats.defense }
  )
  
  character.combatStats.currentHp -= damage
  
  const critText = isCrit ? ' (暴击!)' : ''
  logs.push(createLog(
    `${monster.name} 攻击你造成 ${damage} 点伤害${critText}`,
    'damage'
  ))
  
  return logs
}

// 回合结束处理
export function endTurn(character: Character): void {
  // 技能冷却减少
  character.skills.forEach(skill => {
    if (skill.currentCooldown > 0) {
      skill.currentCooldown--
    }
  })
  
  // 法力回复（基于精神）
  const mpRegen = Math.floor(2 + character.stats.spirit * 0.3)
  character.combatStats.currentMp = Math.min(
    character.combatStats.maxMp,
    character.combatStats.currentMp + mpRegen
  )
}

// 战斗结束 - 获取奖励
export function getVictoryRewards(
  character: Character,
  monster: Monster
): { logs: LogEntry[]; leveledUp: boolean } {
  const logs: LogEntry[] = []
  
  // 经验奖励
  const expGain = monster.expReward
  character.exp += expGain
  logs.push(createLog(
    `获得 ${expGain} 点经验值`,
    'exp'
  ))
  
  // 金币奖励
  const [minGold, maxGold] = monster.goldReward
  const goldGain = Math.floor(Math.random() * (maxGold - minGold + 1)) + minGold
  character.gold += goldGain
  logs.push(createLog(
    `获得 ${goldGain} 金币`,
    'loot'
  ))
  
  // 掉落判定
  monster.lootTable.forEach(loot => {
    if (Math.random() * 100 < loot.dropRate) {
      const [minQty, maxQty] = loot.quantity
      const qty = Math.floor(Math.random() * (maxQty - minQty + 1)) + minQty
      logs.push(createLog(
        `掉落: [${loot.name}] x${qty}`,
        'loot'
      ))
    }
  })
  
  // 升级判定
  let leveledUp = false
  while (character.exp >= character.expToNextLevel) {
    character.exp -= character.expToNextLevel
    character.level++
    character.expToNextLevel = calculateExpToNextLevel(character.level)
    leveledUp = true
    
    // 升级属性增长
    levelUpStats(character)
    
    logs.push(createLog(
      `★ 恭喜升级! 你现在是 ${character.level} 级了! ★`,
      'levelup'
    ))
  }
  
  return { logs, leveledUp }
}

// 计算升级所需经验
export function calculateExpToNextLevel(level: number): number {
  return Math.floor(100 * Math.pow(1.5, level - 1))
}

// 升级时增加属性
function levelUpStats(character: Character): void {
  // 每级增加基础属性
  character.stats.strength += 2
  character.stats.agility += 2
  character.stats.intellect += 2
  character.stats.stamina += 2
  character.stats.spirit += 1
  
  // 重新计算战斗属性
  recalculateCombatStats(character)
  
  // 回满HP/MP
  character.combatStats.currentHp = character.combatStats.maxHp
  character.combatStats.currentMp = character.combatStats.maxMp
}

// 重新计算战斗属性
export function recalculateCombatStats(character: Character): void {
  const stats = character.stats
  
  character.combatStats.maxHp = 100 + stats.stamina * 10
  character.combatStats.maxMp = 50 + stats.intellect * 5
  character.combatStats.attack = 10 + stats.strength * 2 + stats.agility * 0.5
  character.combatStats.defense = 5 + stats.stamina * 0.5 + stats.agility * 0.3
  character.combatStats.critRate = Math.min(50, 5 + stats.agility * 0.3)
  character.combatStats.dodgeRate = Math.min(30, stats.agility * 0.5)
}

// 战斗后恢复（挂机间隙）
export function restBetweenBattles(character: Character): void {
  // 恢复20%生命和法力
  const hpRegen = Math.floor(character.combatStats.maxHp * 0.2)
  const mpRegen = Math.floor(character.combatStats.maxMp * 0.3)
  
  character.combatStats.currentHp = Math.min(
    character.combatStats.maxHp,
    character.combatStats.currentHp + hpRegen
  )
  character.combatStats.currentMp = Math.min(
    character.combatStats.maxMp,
    character.combatStats.currentMp + mpRegen
  )
}

