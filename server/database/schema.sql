-- ═══════════════════════════════════════════════════════════
-- Text WoW 数据库初始化脚本
-- SQLite 版本
-- ═══════════════════════════════════════════════════════════

-- 开启外键约束
PRAGMA foreign_keys = ON;

-- 开启 WAL 模式
PRAGMA journal_mode = WAL;

-- ═══════════════════════════════════════════════════════════
-- 用户相关表
-- ═══════════════════════════════════════════════════════════

-- 用户表 (小队级别数据)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(32) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    email VARCHAR(128) UNIQUE,
    max_team_size INTEGER DEFAULT 5,          -- 最大队伍人数上限
    unlocked_slots INTEGER DEFAULT 1,         -- 已解锁槽位数(初始1个)
    gold INTEGER DEFAULT 0,                   -- 金币(小队共享)
    current_zone_id VARCHAR(32) DEFAULT 'elwynn', -- 当前区域(小队共享)
    total_kills INTEGER DEFAULT 0,            -- 总击杀(小队统计)
    total_gold_gained INTEGER DEFAULT 0,      -- 总获得金币
    play_time INTEGER DEFAULT 0,              -- 游戏时长(秒)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME,
    status INTEGER DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 角色表 (每个用户最多5个角色组成小队)
CREATE TABLE IF NOT EXISTS characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    race_id VARCHAR(32) NOT NULL,
    class_id VARCHAR(32) NOT NULL,
    faction VARCHAR(16) NOT NULL,
    team_slot INTEGER NOT NULL,       -- 队伍位置: 1-5 (1=队长)
    is_active INTEGER DEFAULT 1,      -- 是否出战: 1是 0否
    is_dead INTEGER DEFAULT 0,        -- 是否死亡: 1是 0否
    revive_at DATETIME,               -- 复活时间(NULL表示存活)
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    exp_to_next INTEGER DEFAULT 100,
    hp INTEGER NOT NULL,
    max_hp INTEGER NOT NULL,
    mp INTEGER NOT NULL,
    max_mp INTEGER NOT NULL,
    strength INTEGER DEFAULT 10,
    agility INTEGER DEFAULT 10,
    intellect INTEGER DEFAULT 10,
    stamina INTEGER DEFAULT 10,
    spirit INTEGER DEFAULT 10,
    attack INTEGER DEFAULT 10,
    defense INTEGER DEFAULT 5,
    crit_rate REAL DEFAULT 0.05,
    crit_damage REAL DEFAULT 1.5,
    total_kills INTEGER DEFAULT 0,            -- 该角色击杀数
    total_deaths INTEGER DEFAULT 0,           -- 该角色死亡数
    total_damage_dealt INTEGER DEFAULT 0,     -- 该角色总伤害
    total_healing_done INTEGER DEFAULT 0,     -- 该角色总治疗
    offline_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, team_slot)        -- 每个位置只能有一个角色
);

CREATE INDEX IF NOT EXISTS idx_characters_user_id ON characters(user_id);
CREATE INDEX IF NOT EXISTS idx_characters_level ON characters(level);
CREATE INDEX IF NOT EXISTS idx_characters_team ON characters(user_id, team_slot);

-- ═══════════════════════════════════════════════════════════
-- 游戏配置表
-- ═══════════════════════════════════════════════════════════

-- 种族配置表
CREATE TABLE IF NOT EXISTS races (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    faction VARCHAR(16) NOT NULL,
    description TEXT,
    strength_mod INTEGER DEFAULT 0,
    agility_mod INTEGER DEFAULT 0,
    intellect_mod INTEGER DEFAULT 0,
    stamina_mod INTEGER DEFAULT 0,
    spirit_mod INTEGER DEFAULT 0,
    racial_skill_id VARCHAR(32)
);

-- 职业配置表
CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    role VARCHAR(16) NOT NULL,
    primary_stat VARCHAR(16) NOT NULL,
    base_hp INTEGER NOT NULL,
    base_mp INTEGER NOT NULL,
    hp_per_level INTEGER NOT NULL,
    mp_per_level INTEGER NOT NULL,
    base_strength INTEGER DEFAULT 10,
    base_agility INTEGER DEFAULT 10,
    base_intellect INTEGER DEFAULT 10,
    base_stamina INTEGER DEFAULT 10,
    base_spirit INTEGER DEFAULT 10
);

-- 技能配置表 (扩展版)
CREATE TABLE IF NOT EXISTS skills (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    icon VARCHAR(64),                 -- 图标标识(预留)
    class_id VARCHAR(32),
    type VARCHAR(16) NOT NULL,        -- attack/heal/buff/debuff/dot/hot/shield/control等
    target_type VARCHAR(16) NOT NULL, -- self/ally/enemy/ally_all/enemy_all等
    damage_type VARCHAR(16),          -- physical/magic/fire/frost/shadow/holy/nature
    base_value INTEGER DEFAULT 0,     -- 基础数值
    scaling_stat VARCHAR(16),         -- 成长属性: strength/agility/intellect/spirit
    scaling_ratio REAL DEFAULT 1.0,   -- 属性加成系数
    mp_cost INTEGER DEFAULT 0,
    cooldown INTEGER DEFAULT 0,
    level_required INTEGER DEFAULT 1,
    effect_id VARCHAR(32),            -- 附加效果ID
    effect_chance REAL DEFAULT 1.0,   -- 效果触发概率
    tags TEXT,                        -- 标签(JSON数组)
    FOREIGN KEY (class_id) REFERENCES classes(id),
    FOREIGN KEY (effect_id) REFERENCES effects(id)
);

CREATE INDEX IF NOT EXISTS idx_skills_class_id ON skills(class_id);
CREATE INDEX IF NOT EXISTS idx_skills_type ON skills(type);

-- 效果配置表 (Buff/Debuff)
-- 每场战斗开始时清空所有效果
CREATE TABLE IF NOT EXISTS effects (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    icon VARCHAR(64),
    type VARCHAR(16) NOT NULL,        -- stat_mod/dot/hot/shield/stun/silence等
    is_buff INTEGER NOT NULL,         -- 1=增益 0=减益
    is_stackable INTEGER DEFAULT 0,   -- 是否可叠加
    max_stacks INTEGER DEFAULT 1,     -- 最大叠加层数
    duration INTEGER NOT NULL,        -- 持续回合数
    tick_interval INTEGER DEFAULT 1,  -- 触发间隔
    value_type VARCHAR(16),           -- flat=固定值 percent=百分比
    value REAL,                       -- 效果数值
    stat_affected VARCHAR(32),        -- 影响的属性
    damage_type VARCHAR(16),          -- DOT伤害类型
    can_dispel INTEGER DEFAULT 1,     -- 是否可驱散
    tags TEXT                         -- 标签(JSON数组)
);

CREATE INDEX IF NOT EXISTS idx_effects_type ON effects(type);

-- 物品配置表
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    type VARCHAR(16) NOT NULL,
    subtype VARCHAR(16),
    quality VARCHAR(16) DEFAULT 'common',
    level_required INTEGER DEFAULT 1,
    class_required VARCHAR(32),
    slot VARCHAR(16),
    stackable INTEGER DEFAULT 0,
    max_stack INTEGER DEFAULT 1,
    sell_price INTEGER DEFAULT 0,
    buy_price INTEGER DEFAULT 0,
    strength INTEGER DEFAULT 0,
    agility INTEGER DEFAULT 0,
    intellect INTEGER DEFAULT 0,
    stamina INTEGER DEFAULT 0,
    spirit INTEGER DEFAULT 0,
    attack INTEGER DEFAULT 0,
    defense INTEGER DEFAULT 0,
    hp_bonus INTEGER DEFAULT 0,
    mp_bonus INTEGER DEFAULT 0,
    crit_rate REAL DEFAULT 0,
    effect_type VARCHAR(32),
    effect_value INTEGER
);

CREATE INDEX IF NOT EXISTS idx_items_type ON items(type);
CREATE INDEX IF NOT EXISTS idx_items_quality ON items(quality);

-- 区域配置表
CREATE TABLE IF NOT EXISTS zones (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    min_level INTEGER DEFAULT 1,
    max_level INTEGER DEFAULT 60,
    faction VARCHAR(16),
    parent_zone_id VARCHAR(32),
    exp_modifier REAL DEFAULT 1.0,
    gold_modifier REAL DEFAULT 1.0,
    drop_modifier REAL DEFAULT 1.0,
    is_dungeon INTEGER DEFAULT 0,
    unlock_condition TEXT,
    FOREIGN KEY (parent_zone_id) REFERENCES zones(id)
);

-- 怪物配置表
CREATE TABLE IF NOT EXISTS monsters (
    id VARCHAR(32) PRIMARY KEY,
    zone_id VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    level INTEGER NOT NULL,
    type VARCHAR(16) DEFAULT 'normal',
    hp INTEGER NOT NULL,
    mp INTEGER DEFAULT 0,
    attack INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    exp_reward INTEGER NOT NULL,
    gold_min INTEGER DEFAULT 0,
    gold_max INTEGER DEFAULT 0,
    spawn_weight INTEGER DEFAULT 100,
    skills TEXT,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

CREATE INDEX IF NOT EXISTS idx_monsters_zone_id ON monsters(zone_id);
CREATE INDEX IF NOT EXISTS idx_monsters_level ON monsters(level);

-- 怪物掉落表
CREATE TABLE IF NOT EXISTS monster_drops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    monster_id VARCHAR(32) NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    drop_rate REAL NOT NULL,
    min_quantity INTEGER DEFAULT 1,
    max_quantity INTEGER DEFAULT 1,
    FOREIGN KEY (monster_id) REFERENCES monsters(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX IF NOT EXISTS idx_monster_drops_monster_id ON monster_drops(monster_id);

-- ═══════════════════════════════════════════════════════════
-- 玩家数据表
-- ═══════════════════════════════════════════════════════════

-- 角色技能表
CREATE TABLE IF NOT EXISTS character_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    skill_level INTEGER DEFAULT 1,
    slot INTEGER,
    is_auto INTEGER DEFAULT 1,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id),
    UNIQUE(character_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_char_skills_char_id ON character_skills(character_id);

-- 背包表
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    quantity INTEGER DEFAULT 1,
    slot INTEGER,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_char_id ON inventory(character_id);

-- 装备表
CREATE TABLE IF NOT EXISTS equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    slot VARCHAR(16) NOT NULL,
    item_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES items(id),
    UNIQUE(character_id, slot)
);

CREATE INDEX IF NOT EXISTS idx_equipment_char_id ON equipment(character_id);

-- 战斗策略表
CREATE TABLE IF NOT EXISTS battle_strategies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,
    priority INTEGER NOT NULL,
    condition_type VARCHAR(32) NOT NULL,
    condition_operator VARCHAR(8) NOT NULL,
    condition_value REAL NOT NULL,
    action_type VARCHAR(32) NOT NULL,
    action_target VARCHAR(32),
    skill_id VARCHAR(32),
    item_id VARCHAR(32),
    is_active INTEGER DEFAULT 1,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX IF NOT EXISTS idx_strategies_char_id ON battle_strategies(character_id);

-- 游戏会话表
CREATE TABLE IF NOT EXISTS game_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    session_start DATETIME NOT NULL,
    session_end DATETIME,
    kills INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_char_id ON game_sessions(character_id);

-- ═══════════════════════════════════════════════════════════
-- 触发器
-- ═══════════════════════════════════════════════════════════

-- 自动更新 updated_at
CREATE TRIGGER IF NOT EXISTS update_character_timestamp 
AFTER UPDATE ON characters
BEGIN
    UPDATE characters SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

