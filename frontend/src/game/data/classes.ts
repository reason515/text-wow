import type { ClassInfo, CharacterClass, Race } from '../types'

export const CLASSES: Record<CharacterClass, ClassInfo> = {
  warrior: {
    id: 'warrior',
    name: '战士',
    description: '近战格斗大师，可以使用各种武器和板甲。',
    primaryStat: 'strength',
    startingSkills: ['heroic_strike', 'charge', 'battle_shout'],
    availableRaces: ['human', 'dwarf', 'nightelf', 'gnome', 'draenei', 'orc', 'troll', 'undead', 'tauren']
  },
  paladin: {
    id: 'paladin',
    name: '圣骑士',
    description: '圣光的勇士，可以治疗、坦克和输出。',
    primaryStat: 'strength',
    startingSkills: ['crusader_strike', 'holy_light', 'blessing_of_might'],
    availableRaces: ['human', 'dwarf', 'draenei', 'bloodelf']
  },
  hunter: {
    id: 'hunter',
    name: '猎人',
    description: '远程物理输出，拥有宠物协助战斗。',
    primaryStat: 'agility',
    startingSkills: ['arcane_shot', 'serpent_sting', 'aspect_of_hawk'],
    availableRaces: ['human', 'dwarf', 'nightelf', 'draenei', 'orc', 'troll', 'tauren', 'bloodelf']
  },
  rogue: {
    id: 'rogue',
    name: '盗贼',
    description: '隐匿的刺客，擅长暴击和控制。',
    primaryStat: 'agility',
    startingSkills: ['sinister_strike', 'backstab', 'eviscerate'],
    availableRaces: ['human', 'dwarf', 'nightelf', 'gnome', 'orc', 'troll', 'undead', 'bloodelf']
  },
  priest: {
    id: 'priest',
    name: '牧师',
    description: '强大的治疗者，也可以使用暗影魔法。',
    primaryStat: 'intellect',
    startingSkills: ['smite', 'lesser_heal', 'shadow_word_pain'],
    availableRaces: ['human', 'dwarf', 'nightelf', 'draenei', 'troll', 'undead', 'bloodelf']
  },
  mage: {
    id: 'mage',
    name: '法师',
    description: '操控火焰、冰霜和奥术的法术大师。',
    primaryStat: 'intellect',
    startingSkills: ['fireball', 'frostbolt', 'arcane_missiles'],
    availableRaces: ['human', 'gnome', 'draenei', 'troll', 'undead', 'bloodelf']
  },
  warlock: {
    id: 'warlock',
    name: '术士',
    description: '召唤恶魔、施放诅咒的黑暗法师。',
    primaryStat: 'intellect',
    startingSkills: ['shadow_bolt', 'corruption', 'life_tap'],
    availableRaces: ['human', 'gnome', 'orc', 'undead', 'bloodelf']
  },
  druid: {
    id: 'druid',
    name: '德鲁伊',
    description: '自然的守护者，可以变形为多种形态。',
    primaryStat: 'intellect',
    startingSkills: ['wrath', 'rejuvenation', 'moonfire'],
    availableRaces: ['nightelf', 'tauren']
  },
  shaman: {
    id: 'shaman',
    name: '萨满',
    description: '元素的使者，可以治疗和造成元素伤害。',
    primaryStat: 'intellect',
    startingSkills: ['lightning_bolt', 'earth_shock', 'healing_wave'],
    availableRaces: ['draenei', 'orc', 'troll', 'tauren']
  }
}

export function getAvailableClasses(race: Race): ClassInfo[] {
  return Object.values(CLASSES).filter(c => c.availableRaces.includes(race))
}

