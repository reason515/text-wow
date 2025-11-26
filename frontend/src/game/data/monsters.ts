import type { Monster, Zone } from '../types'

// 怪物模板
export const MONSTER_TEMPLATES: Record<string, Omit<Monster, 'currentHp'>> = {
  // 艾尔文森林 (1-10级)
  forest_wolf: {
    id: 'forest_wolf',
    name: '森林狼',
    level: 2,
    maxHp: 45,
    attack: 8,
    defense: 2,
    expReward: 25,
    goldReward: [2, 5],
    lootTable: [
      { itemId: 'wolf_pelt', name: '狼皮', dropRate: 60, quantity: [1, 1] },
      { itemId: 'wolf_meat', name: '狼肉', dropRate: 40, quantity: [1, 2] }
    ]
  },
  young_boar: {
    id: 'young_boar',
    name: '小野猪',
    level: 1,
    maxHp: 35,
    attack: 6,
    defense: 1,
    expReward: 18,
    goldReward: [1, 3],
    lootTable: [
      { itemId: 'boar_meat', name: '野猪肉', dropRate: 70, quantity: [1, 2] }
    ]
  },
  forest_spider: {
    id: 'forest_spider',
    name: '森林蜘蛛',
    level: 3,
    maxHp: 55,
    attack: 10,
    defense: 3,
    expReward: 32,
    goldReward: [3, 7],
    lootTable: [
      { itemId: 'spider_silk', name: '蛛丝', dropRate: 50, quantity: [1, 3] },
      { itemId: 'spider_venom', name: '蜘蛛毒液', dropRate: 25, quantity: [1, 1] }
    ]
  },
  defias_thug: {
    id: 'defias_thug',
    name: '迪菲亚暴徒',
    level: 5,
    maxHp: 85,
    attack: 15,
    defense: 5,
    expReward: 55,
    goldReward: [8, 15],
    lootTable: [
      { itemId: 'linen_cloth', name: '亚麻布', dropRate: 45, quantity: [1, 2] },
      { itemId: 'red_bandana', name: '红色头巾', dropRate: 20, quantity: [1, 1] }
    ]
  },
  murloc_scout: {
    id: 'murloc_scout',
    name: '鱼人斥候',
    level: 4,
    maxHp: 65,
    attack: 12,
    defense: 4,
    expReward: 42,
    goldReward: [5, 10],
    lootTable: [
      { itemId: 'murloc_fin', name: '鱼人鳍', dropRate: 55, quantity: [1, 2] },
      { itemId: 'clam_meat', name: '蛤肉', dropRate: 30, quantity: [1, 1] }
    ]
  },

  // 西部荒野 (10-20级)
  harvest_golem: {
    id: 'harvest_golem',
    name: '收割傀儡',
    level: 12,
    maxHp: 180,
    attack: 28,
    defense: 12,
    expReward: 120,
    goldReward: [15, 25],
    lootTable: [
      { itemId: 'golem_core', name: '傀儡核心', dropRate: 30, quantity: [1, 1] },
      { itemId: 'iron_scrap', name: '废铁', dropRate: 50, quantity: [1, 3] }
    ]
  },
  defias_pillager: {
    id: 'defias_pillager',
    name: '迪菲亚掠夺者',
    level: 14,
    maxHp: 220,
    attack: 35,
    defense: 14,
    expReward: 150,
    goldReward: [20, 35],
    lootTable: [
      { itemId: 'wool_cloth', name: '毛料', dropRate: 50, quantity: [1, 2] },
      { itemId: 'defias_dagger', name: '迪菲亚匕首', dropRate: 15, quantity: [1, 1] }
    ]
  },
  coyote: {
    id: 'coyote',
    name: '郊狼',
    level: 11,
    maxHp: 160,
    attack: 25,
    defense: 10,
    expReward: 100,
    goldReward: [12, 20],
    lootTable: [
      { itemId: 'coyote_fang', name: '郊狼牙', dropRate: 45, quantity: [1, 2] }
    ]
  },

  // 暮色森林 (20-30级)
  skeletal_warrior: {
    id: 'skeletal_warrior',
    name: '骷髅战士',
    level: 22,
    maxHp: 380,
    attack: 55,
    defense: 25,
    expReward: 280,
    goldReward: [30, 50],
    lootTable: [
      { itemId: 'bone_fragment', name: '骨头碎片', dropRate: 60, quantity: [1, 3] },
      { itemId: 'tarnished_sword', name: '锈蚀的剑', dropRate: 20, quantity: [1, 1] }
    ]
  },
  dire_wolf: {
    id: 'dire_wolf',
    name: '恐狼',
    level: 24,
    maxHp: 420,
    attack: 62,
    defense: 28,
    expReward: 320,
    goldReward: [35, 55],
    lootTable: [
      { itemId: 'dire_wolf_fang', name: '恐狼牙', dropRate: 40, quantity: [1, 1] },
      { itemId: 'thick_leather', name: '厚皮革', dropRate: 35, quantity: [1, 2] }
    ]
  },
  worgen: {
    id: 'worgen',
    name: '狼人',
    level: 26,
    maxHp: 500,
    attack: 72,
    defense: 32,
    expReward: 380,
    goldReward: [45, 70],
    lootTable: [
      { itemId: 'worgen_claw', name: '狼人之爪', dropRate: 35, quantity: [1, 1] },
      { itemId: 'shadow_gem', name: '暗影宝石', dropRate: 15, quantity: [1, 1] }
    ]
  }
}

// 战斗区域
export const ZONES: Record<string, Zone> = {
  elwynn_forest: {
    id: 'elwynn_forest',
    name: '艾尔文森林',
    description: '暴风城附近的宁静森林，适合初出茅庐的冒险者。',
    levelRange: [1, 10],
    monsters: ['young_boar', 'forest_wolf', 'forest_spider', 'murloc_scout', 'defias_thug'],
    unlockLevel: 1
  },
  westfall: {
    id: 'westfall',
    name: '西部荒野',
    description: '曾经繁荣的农业区，如今被迪菲亚兄弟会占据。',
    levelRange: [10, 20],
    monsters: ['coyote', 'harvest_golem', 'defias_pillager'],
    unlockLevel: 10
  },
  duskwood: {
    id: 'duskwood',
    name: '暮色森林',
    description: '被永恒黑暗笼罩的诅咒之地，充满亡灵和狼人。',
    levelRange: [20, 30],
    monsters: ['skeletal_warrior', 'dire_wolf', 'worgen'],
    unlockLevel: 20
  }
}

// 创建怪物实例
export function createMonster(templateId: string): Monster {
  const template = MONSTER_TEMPLATES[templateId]
  if (!template) {
    throw new Error(`Monster template not found: ${templateId}`)
  }
  return {
    ...template,
    currentHp: template.maxHp
  }
}

// 根据等级获取可用区域
export function getUnlockedZones(level: number): Zone[] {
  return Object.values(ZONES).filter(zone => zone.unlockLevel <= level)
}

// 从区域随机选择怪物
export function getRandomMonsterFromZone(zoneId: string): Monster {
  const zone = ZONES[zoneId]
  if (!zone) {
    throw new Error(`Zone not found: ${zoneId}`)
  }
  const randomIndex = Math.floor(Math.random() * zone.monsters.length)
  return createMonster(zone.monsters[randomIndex])
}

