import type { RaceInfo, Race, Faction } from '../types'

export const RACES: Record<Race, RaceInfo> = {
  // 联盟
  human: {
    id: 'human',
    name: '人类',
    faction: 'alliance',
    description: '适应力强的人类是艾泽拉斯最年轻也最有活力的种族。',
    bonusStats: { spirit: 2, intellect: 1 }
  },
  dwarf: {
    id: 'dwarf',
    name: '矮人',
    faction: 'alliance',
    description: '矮人以其坚韧和酿酒技术闻名于世。',
    bonusStats: { stamina: 2, strength: 1 }
  },
  nightelf: {
    id: 'nightelf',
    name: '暗夜精灵',
    faction: 'alliance',
    description: '古老的暗夜精灵与自然有着深厚的联系。',
    bonusStats: { agility: 2, intellect: 1 }
  },
  gnome: {
    id: 'gnome',
    name: '侏儒',
    faction: 'alliance',
    description: '聪明绝顶的侏儒是天生的发明家和工程师。',
    bonusStats: { intellect: 3 }
  },
  draenei: {
    id: 'draenei',
    name: '德莱尼',
    faction: 'alliance',
    description: '来自外域的德莱尼人追随圣光的指引。',
    bonusStats: { intellect: 1, spirit: 2 }
  },
  
  // 部落
  orc: {
    id: 'orc',
    name: '兽人',
    faction: 'horde',
    description: '强壮的兽人战士以荣誉为生命的信条。',
    bonusStats: { strength: 3 }
  },
  troll: {
    id: 'troll',
    name: '巨魔',
    faction: 'horde',
    description: '狡猾的巨魔拥有惊人的再生能力。',
    bonusStats: { agility: 2, stamina: 1 }
  },
  undead: {
    id: 'undead',
    name: '亡灵',
    faction: 'horde',
    description: '被遗忘者已经摆脱了巫妖王的控制。',
    bonusStats: { intellect: 2, spirit: 1 }
  },
  tauren: {
    id: 'tauren',
    name: '牛头人',
    faction: 'horde',
    description: '高大的牛头人是大地之母的忠实追随者。',
    bonusStats: { stamina: 3 }
  },
  bloodelf: {
    id: 'bloodelf',
    name: '血精灵',
    faction: 'horde',
    description: '高傲的血精灵对魔法有着天生的亲和力。',
    bonusStats: { intellect: 2, agility: 1 }
  }
}

export function getRacesByFaction(faction: Faction): RaceInfo[] {
  return Object.values(RACES).filter(race => race.faction === faction)
}

