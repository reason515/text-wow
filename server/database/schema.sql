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
    resource INTEGER NOT NULL,         -- 当前能量值(怒气/能量/法力)
    max_resource INTEGER NOT NULL,     -- 最大能量值
    resource_type VARCHAR(16) NOT NULL, -- 能量类型: mana/rage/energy
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

-- 种族配置表 (基础值+百分比加成，适合放置游戏)
CREATE TABLE IF NOT EXISTS races (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    faction VARCHAR(16) NOT NULL,
    description TEXT,
    -- 基础加成(固定值，创建时一次性加)
    strength_base INTEGER DEFAULT 0,
    agility_base INTEGER DEFAULT 0,
    intellect_base INTEGER DEFAULT 0,
    stamina_base INTEGER DEFAULT 0,
    spirit_base INTEGER DEFAULT 0,
    -- 百分比加成(乘算，随等级保持意义) 5 = 5%
    strength_pct REAL DEFAULT 0,
    agility_pct REAL DEFAULT 0,
    intellect_pct REAL DEFAULT 0,
    stamina_pct REAL DEFAULT 0,
    spirit_pct REAL DEFAULT 0,
    -- 被动特性
    racial_passive_id VARCHAR(32),
    racial_passive2_id VARCHAR(32),
    allowed_classes TEXT,
    FOREIGN KEY (racial_passive_id) REFERENCES effects(id),
    FOREIGN KEY (racial_passive2_id) REFERENCES effects(id)
);

-- 职业配置表 (不同职业使用不同能量类型)
CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    role VARCHAR(16) NOT NULL,
    primary_stat VARCHAR(16) NOT NULL,
    resource_type VARCHAR(16) NOT NULL,   -- mana/rage/energy
    base_hp INTEGER NOT NULL,
    base_resource INTEGER NOT NULL,       -- 基础能量值
    hp_per_level INTEGER NOT NULL,
    resource_per_level INTEGER NOT NULL,  -- 每级能量成长
    resource_regen REAL DEFAULT 0,        -- 每回合固定恢复
    resource_regen_pct REAL DEFAULT 0,    -- 每回合百分比恢复(基于精神)
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
    resource_cost INTEGER DEFAULT 0,  -- 能量消耗(怒气/能量/法力)
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

-- ═══════════════════════════════════════════════════════════
-- 后期玩法系统
-- ═══════════════════════════════════════════════════════════

-- 无尽深渊配置表
CREATE TABLE IF NOT EXISTS abyss_config (
    floor INTEGER PRIMARY KEY,
    monster_level_base INTEGER NOT NULL,
    monster_level_growth REAL DEFAULT 0.5,
    monster_hp_mult REAL DEFAULT 1.0,
    monster_atk_mult REAL DEFAULT 1.0,
    reward_exp_mult REAL DEFAULT 1.0,
    reward_gold_mult REAL DEFAULT 1.0,
    special_reward TEXT,
    boss_id VARCHAR(32),
    FOREIGN KEY (boss_id) REFERENCES monsters(id)
);

-- 无尽深渊进度表
CREATE TABLE IF NOT EXISTS abyss_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    highest_floor INTEGER DEFAULT 0,
    current_floor INTEGER DEFAULT 1,
    weekly_attempts INTEGER DEFAULT 0,
    weekly_reset_at DATETIME,
    total_clears INTEGER DEFAULT 0,
    best_time INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_abyss_highest ON abyss_progress(highest_floor DESC);

-- 装备强化表
CREATE TABLE IF NOT EXISTS equipment_enhance (
    equipment_id INTEGER PRIMARY KEY,
    enhance_level INTEGER DEFAULT 0,       -- 强化等级 0-15
    refine_level INTEGER DEFAULT 0,        -- 精炼等级 0-10
    enchant_id VARCHAR(32),                -- 附魔ID
    gem_slot_1 VARCHAR(32),                -- 宝石槽1
    gem_slot_2 VARCHAR(32),                -- 宝石槽2
    gem_slot_3 VARCHAR(32),                -- 宝石槽3
    FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE
);

-- 附魔配置表
CREATE TABLE IF NOT EXISTS enchants (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    slot_type VARCHAR(16),                 -- 适用槽位
    effect_type VARCHAR(32) NOT NULL,
    effect_value REAL NOT NULL,
    quality VARCHAR(16)
);

-- 宝石配置表
CREATE TABLE IF NOT EXISTS gems (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    color VARCHAR(16) NOT NULL,            -- red/blue/yellow/green/purple
    stat_type VARCHAR(32) NOT NULL,
    stat_value INTEGER NOT NULL,
    quality VARCHAR(16)
);

-- 成就配置表
CREATE TABLE IF NOT EXISTS achievements (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    category VARCHAR(32) NOT NULL,         -- combat/explore/collect/social/special
    condition_type VARCHAR(32) NOT NULL,
    condition_value INTEGER NOT NULL,
    points INTEGER DEFAULT 10,
    reward_type VARCHAR(32),
    reward_value TEXT,
    icon VARCHAR(64),
    is_hidden INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);

-- 玩家成就表
CREATE TABLE IF NOT EXISTS user_achievements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    achievement_id VARCHAR(32) NOT NULL,
    progress INTEGER DEFAULT 0,
    completed_at DATETIME,
    rewarded INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_user_achievements_user ON user_achievements(user_id);

-- 图鉴配置表
CREATE TABLE IF NOT EXISTS codex (
    id VARCHAR(32) PRIMARY KEY,
    category VARCHAR(32) NOT NULL,         -- monster/item/boss/zone
    target_id VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    unlock_condition TEXT,
    bonus_type VARCHAR(32),
    bonus_value REAL
);

CREATE INDEX IF NOT EXISTS idx_codex_category ON codex(category);

-- 玩家图鉴表
CREATE TABLE IF NOT EXISTS user_codex (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    codex_id VARCHAR(32) NOT NULL,
    unlock_count INTEGER DEFAULT 1,
    first_unlock_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (codex_id) REFERENCES codex(id),
    UNIQUE(user_id, codex_id)
);

CREATE INDEX IF NOT EXISTS idx_user_codex_user ON user_codex(user_id);

-- ═══════════════════════════════════════════════════════════
-- 阵营PVP遭遇战系统
-- ═══════════════════════════════════════════════════════════

-- 地图阵营控制表
-- 记录每个地图的阵营控制状态和效率加成
CREATE TABLE IF NOT EXISTS zone_faction_control (
    zone_id VARCHAR(32) PRIMARY KEY,
    controlling_faction VARCHAR(16) DEFAULT 'neutral', -- alliance/horde/neutral
    alliance_wins INTEGER DEFAULT 0,          -- 联盟近期胜场
    horde_wins INTEGER DEFAULT 0,             -- 部落近期胜场
    alliance_win_rate REAL DEFAULT 0.5,       -- 联盟胜率
    control_score INTEGER DEFAULT 0,          -- 控制积分 (正=联盟, 负=部落, 范围-100~+100)
    efficiency_bonus REAL DEFAULT 0,          -- 控制方效率加成 (0.0-0.2)
    efficiency_penalty REAL DEFAULT 0,        -- 被控方效率惩罚 (0.0-0.1)
    last_battle_at DATETIME,                  -- 最后一次战斗时间
    stats_reset_at DATETIME,                  -- 统计重置时间 (每周)
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

-- PVP遭遇战记录表
-- 记录每一场PVP遭遇战的详细信息
CREATE TABLE IF NOT EXISTS pvp_encounters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    zone_id VARCHAR(32) NOT NULL,             -- 发生区域
    attacker_user_id INTEGER NOT NULL,        -- 进攻方用户ID
    defender_user_id INTEGER NOT NULL,        -- 防守方用户ID
    attacker_faction VARCHAR(16) NOT NULL,    -- 进攻方阵营
    defender_faction VARCHAR(16) NOT NULL,    -- 防守方阵营
    winner_user_id INTEGER,                   -- 胜利方用户ID (NULL=平局)
    winner_faction VARCHAR(16),               -- 胜利阵营
    attacker_team_info TEXT,                  -- 进攻方队伍信息 (JSON)
    defender_team_info TEXT,                  -- 防守方队伍信息 (JSON)
    battle_rounds INTEGER DEFAULT 0,          -- 战斗回合数
    battle_duration INTEGER DEFAULT 0,        -- 战斗时长(秒)
    attacker_damage_dealt INTEGER DEFAULT 0,  -- 进攻方造成伤害
    defender_damage_dealt INTEGER DEFAULT 0,  -- 防守方造成伤害
    honor_reward INTEGER DEFAULT 0,           -- 荣誉奖励
    battle_log TEXT,                          -- 战斗日志 (JSON)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (attacker_user_id) REFERENCES users(id),
    FOREIGN KEY (defender_user_id) REFERENCES users(id),
    FOREIGN KEY (winner_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_pvp_encounters_zone ON pvp_encounters(zone_id);
CREATE INDEX IF NOT EXISTS idx_pvp_encounters_attacker ON pvp_encounters(attacker_user_id);
CREATE INDEX IF NOT EXISTS idx_pvp_encounters_defender ON pvp_encounters(defender_user_id);
CREATE INDEX IF NOT EXISTS idx_pvp_encounters_time ON pvp_encounters(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_pvp_encounters_winner ON pvp_encounters(winner_faction);

-- 玩家荣誉表
-- 记录玩家PVP战绩和荣誉积累
CREATE TABLE IF NOT EXISTS user_honor (
    user_id INTEGER PRIMARY KEY,
    faction VARCHAR(16) NOT NULL,             -- 所属阵营
    total_honor INTEGER DEFAULT 0,            -- 累计荣誉值
    current_honor INTEGER DEFAULT 0,          -- 当前可用荣誉
    honor_rank INTEGER DEFAULT 0,             -- 荣誉军衔等级 (0-14)
    pvp_wins INTEGER DEFAULT 0,               -- PVP胜场
    pvp_losses INTEGER DEFAULT 0,             -- PVP败场
    pvp_draws INTEGER DEFAULT 0,              -- PVP平局
    win_streak INTEGER DEFAULT 0,             -- 当前连胜
    best_win_streak INTEGER DEFAULT 0,        -- 最高连胜
    total_kills INTEGER DEFAULT 0,            -- 总击杀角色数
    total_deaths INTEGER DEFAULT 0,           -- 总死亡角色数
    total_damage_dealt INTEGER DEFAULT 0,     -- 总造成伤害
    weekly_honor INTEGER DEFAULT 0,           -- 本周荣誉
    weekly_reset_at DATETIME,                 -- 周重置时间
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_honor_rank ON user_honor(honor_rank DESC);
CREATE INDEX IF NOT EXISTS idx_user_honor_wins ON user_honor(pvp_wins DESC);
CREATE INDEX IF NOT EXISTS idx_user_honor_faction ON user_honor(faction);

-- 全服公告表
-- 记录并推送PVP战斗结果到全服
CREATE TABLE IF NOT EXISTS server_announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(32) NOT NULL,                -- pvp_victory/zone_captured/kill_streak/zone_contested/faction_dominant
    content TEXT NOT NULL,                    -- 公告内容
    zone_id VARCHAR(32),                      -- 相关区域
    winner_user_id INTEGER,                   -- 胜利者ID
    loser_user_id INTEGER,                    -- 失败者ID
    pvp_encounter_id INTEGER,                 -- 关联PVP战斗ID
    importance INTEGER DEFAULT 1,             -- 重要程度 (1-5)
    expires_at DATETIME,                      -- 过期时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (winner_user_id) REFERENCES users(id),
    FOREIGN KEY (loser_user_id) REFERENCES users(id),
    FOREIGN KEY (pvp_encounter_id) REFERENCES pvp_encounters(id)
);

CREATE INDEX IF NOT EXISTS idx_announcements_time ON server_announcements(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_announcements_type ON server_announcements(type);

-- 荣誉商店配置表
-- 使用荣誉值兑换独特奖励
CREATE TABLE IF NOT EXISTS honor_shop (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,                -- 商品名称
    description TEXT,                         -- 描述
    item_type VARCHAR(32) NOT NULL,           -- equipment/consumable/cosmetic/title
    item_id VARCHAR(32),                      -- 关联物品ID
    honor_cost INTEGER NOT NULL,              -- 荣誉花费
    rank_required INTEGER DEFAULT 0,          -- 需求军衔等级
    faction VARCHAR(16),                      -- 阵营限制 (NULL=双阵营)
    weekly_limit INTEGER,                     -- 每周购买限制
    stock INTEGER,                            -- 库存 (NULL=无限)
    is_active INTEGER DEFAULT 1,              -- 是否在售
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE INDEX IF NOT EXISTS idx_honor_shop_rank ON honor_shop(rank_required);

