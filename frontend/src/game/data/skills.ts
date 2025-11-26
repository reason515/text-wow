import type { Skill } from '../types'

export const SKILLS: Record<string, Skill> = {
  // 战士技能
  heroic_strike: {
    id: 'heroic_strike',
    name: '英勇打击',
    description: '一次强力的武器攻击。',
    damage: 1.5,
    mpCost: 15,
    cooldown: 0,
    currentCooldown: 0,
    type: 'physical'
  },
  charge: {
    id: 'charge',
    name: '冲锋',
    description: '向敌人冲锋，造成伤害并眩晕。',
    damage: 1.2,
    mpCost: 10,
    cooldown: 3,
    currentCooldown: 0,
    type: 'physical'
  },
  battle_shout: {
    id: 'battle_shout',
    name: '战斗怒吼',
    description: '提升攻击力。',
    damage: 0,
    mpCost: 20,
    cooldown: 5,
    currentCooldown: 0,
    type: 'physical'
  },
  
  // 圣骑士技能
  crusader_strike: {
    id: 'crusader_strike',
    name: '十字军打击',
    description: '以圣光之力攻击敌人。',
    damage: 1.4,
    mpCost: 20,
    cooldown: 1,
    currentCooldown: 0,
    type: 'physical'
  },
  holy_light: {
    id: 'holy_light',
    name: '圣光术',
    description: '治疗自己。',
    damage: -2.0,  // 负数表示治疗
    mpCost: 30,
    cooldown: 2,
    currentCooldown: 0,
    type: 'heal'
  },
  blessing_of_might: {
    id: 'blessing_of_might',
    name: '力量祝福',
    description: '增加攻击力。',
    damage: 0,
    mpCost: 25,
    cooldown: 10,
    currentCooldown: 0,
    type: 'physical'
  },
  
  // 猎人技能
  arcane_shot: {
    id: 'arcane_shot',
    name: '奥术射击',
    description: '发射一支奥术箭矢。',
    damage: 1.3,
    mpCost: 15,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  serpent_sting: {
    id: 'serpent_sting',
    name: '毒蛇钉刺',
    description: '对敌人造成持续毒素伤害。',
    damage: 0.8,
    mpCost: 20,
    cooldown: 2,
    currentCooldown: 0,
    type: 'physical'
  },
  aspect_of_hawk: {
    id: 'aspect_of_hawk',
    name: '雄鹰守护',
    description: '提升攻击力。',
    damage: 0,
    mpCost: 10,
    cooldown: 10,
    currentCooldown: 0,
    type: 'physical'
  },
  
  // 盗贼技能
  sinister_strike: {
    id: 'sinister_strike',
    name: '邪恶攻击',
    description: '快速的武器攻击。',
    damage: 1.3,
    mpCost: 10,
    cooldown: 0,
    currentCooldown: 0,
    type: 'physical'
  },
  backstab: {
    id: 'backstab',
    name: '背刺',
    description: '从背后发起致命一击。',
    damage: 2.0,
    mpCost: 25,
    cooldown: 2,
    currentCooldown: 0,
    type: 'physical'
  },
  eviscerate: {
    id: 'eviscerate',
    name: '剔骨',
    description: '终结技，造成高额伤害。',
    damage: 2.5,
    mpCost: 35,
    cooldown: 3,
    currentCooldown: 0,
    type: 'physical'
  },
  
  // 牧师技能
  smite: {
    id: 'smite',
    name: '惩击',
    description: '以圣光惩击敌人。',
    damage: 1.2,
    mpCost: 20,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  lesser_heal: {
    id: 'lesser_heal',
    name: '次级治疗术',
    description: '恢复生命值。',
    damage: -1.5,
    mpCost: 25,
    cooldown: 1,
    currentCooldown: 0,
    type: 'heal'
  },
  shadow_word_pain: {
    id: 'shadow_word_pain',
    name: '暗言术：痛',
    description: '对敌人造成暗影伤害。',
    damage: 1.4,
    mpCost: 25,
    cooldown: 1,
    currentCooldown: 0,
    type: 'magical'
  },
  
  // 法师技能
  fireball: {
    id: 'fireball',
    name: '火球术',
    description: '发射一颗炽热的火球。',
    damage: 1.6,
    mpCost: 25,
    cooldown: 1,
    currentCooldown: 0,
    type: 'magical'
  },
  frostbolt: {
    id: 'frostbolt',
    name: '寒冰箭',
    description: '发射一支寒冰箭，减速敌人。',
    damage: 1.3,
    mpCost: 20,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  arcane_missiles: {
    id: 'arcane_missiles',
    name: '奥术飞弹',
    description: '连续发射奥术飞弹。',
    damage: 1.8,
    mpCost: 30,
    cooldown: 2,
    currentCooldown: 0,
    type: 'magical'
  },
  
  // 术士技能
  shadow_bolt: {
    id: 'shadow_bolt',
    name: '暗影箭',
    description: '发射一支暗影箭。',
    damage: 1.5,
    mpCost: 20,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  corruption: {
    id: 'corruption',
    name: '腐蚀术',
    description: '对敌人造成持续暗影伤害。',
    damage: 1.0,
    mpCost: 25,
    cooldown: 2,
    currentCooldown: 0,
    type: 'magical'
  },
  life_tap: {
    id: 'life_tap',
    name: '生命分流',
    description: '消耗生命恢复法力。',
    damage: 0,
    mpCost: 0,
    cooldown: 1,
    currentCooldown: 0,
    type: 'magical'
  },
  
  // 德鲁伊技能
  wrath: {
    id: 'wrath',
    name: '愤怒',
    description: '召唤自然之力攻击敌人。',
    damage: 1.4,
    mpCost: 20,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  rejuvenation: {
    id: 'rejuvenation',
    name: '回春术',
    description: '为目标恢复生命值。',
    damage: -1.2,
    mpCost: 20,
    cooldown: 1,
    currentCooldown: 0,
    type: 'heal'
  },
  moonfire: {
    id: 'moonfire',
    name: '月火术',
    description: '月光造成奥术伤害。',
    damage: 1.2,
    mpCost: 15,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  
  // 萨满技能
  lightning_bolt: {
    id: 'lightning_bolt',
    name: '闪电箭',
    description: '召唤闪电攻击敌人。',
    damage: 1.5,
    mpCost: 22,
    cooldown: 0,
    currentCooldown: 0,
    type: 'magical'
  },
  earth_shock: {
    id: 'earth_shock',
    name: '大地震击',
    description: '用大地之力震击敌人。',
    damage: 1.3,
    mpCost: 18,
    cooldown: 1,
    currentCooldown: 0,
    type: 'magical'
  },
  healing_wave: {
    id: 'healing_wave',
    name: '治疗波',
    description: '用水之力治疗目标。',
    damage: -1.8,
    mpCost: 28,
    cooldown: 2,
    currentCooldown: 0,
    type: 'heal'
  }
}

export function getSkillsByIds(ids: string[]): Skill[] {
  return ids.map(id => ({ ...SKILLS[id] })).filter(Boolean)
}

