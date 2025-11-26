export interface Character {
  name: string
  race: string
  class: string
  level: number
  exp: number
  exp_to_next: number
  max_hp: number
  hp: number
  max_mp: number
  mp: number
  attack: number
  defense: number
  gold: number
  total_kills: number
}

export interface Monster {
  id: string
  name: string
  level: number
  max_hp: number
  hp: number
  attack: number
  defense: number
}

export interface BattleLog {
  time: string
  type: 'info' | 'damage' | 'heal' | 'loot' | 'exp' | 'levelup'
  message: string
  color: string
}

export interface BattleStatus {
  is_running: boolean
  current_zone: string
  current_enemy: Monster | null
  battle_count: number
  session_kills: number
  session_gold: number
  session_exp: number
}

export interface BattleResult {
  character: Character
  enemy: Monster | null
  logs: BattleLog[]
  status: BattleStatus
}

export interface Zone {
  id: string
  name: string
  description: string
  min_level: number
}



