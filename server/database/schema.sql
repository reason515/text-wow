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
    physical_attack INTEGER DEFAULT 10,
    magic_attack INTEGER DEFAULT 10,
    physical_defense INTEGER DEFAULT 5,
    magic_defense INTEGER DEFAULT 5,
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
-- 游戏规则公式表 (玩家可查询)
-- ═══════════════════════════════════════════════════════════

-- 公式配置表 - 存储所有战斗计算规则，供玩家查询
CREATE TABLE IF NOT EXISTS game_formulas (
    id VARCHAR(64) PRIMARY KEY,
    category VARCHAR(32) NOT NULL,          -- attribute/combat/skill/resource
    name VARCHAR(64) NOT NULL,              -- 公式名称
    formula TEXT NOT NULL,                  -- 公式表达式
    description TEXT,                       -- 详细说明
    variables TEXT,                         -- 变量说明(JSON)
    example TEXT,                           -- 计算示例
    display_order INTEGER DEFAULT 0         -- 显示顺序
);

CREATE INDEX IF NOT EXISTS idx_formulas_category ON game_formulas(category);

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
    base_spirit INTEGER DEFAULT 10,
    -- 仇恨系统字段
    base_threat_modifier REAL DEFAULT 1.0, -- 基础仇恨系数
    combat_role VARCHAR(16) DEFAULT 'dps', -- 战斗定位: tank/healer/dps/hybrid
    is_ranged INTEGER DEFAULT 0            -- 是否远程: 1是 0否 (影响OT阈值)
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
    -- 仇恨系统字段
    threat_modifier REAL DEFAULT 1.0, -- 仇恨系数 (2.0=双倍仇恨, 0.5=半仇恨)
    threat_type VARCHAR(16) DEFAULT 'normal', -- 仇恨类型: normal/high/taunt/reduce/clear
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

-- 装备词缀表
CREATE TABLE IF NOT EXISTS affixes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    type VARCHAR(16) NOT NULL,            -- prefix/suffix
    slot_type VARCHAR(16) DEFAULT 'all',  -- weapon/armor/accessory/all
    rarity VARCHAR(16) NOT NULL DEFAULT 'common',
    effect_type VARCHAR(32) NOT NULL,
    effect_stat VARCHAR(32),
    min_value REAL NOT NULL,
    max_value REAL NOT NULL,
    value_type VARCHAR(16) NOT NULL,      -- flat/percent
    description TEXT,
    level_required INTEGER DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_affixes_type ON affixes(type);
CREATE INDEX IF NOT EXISTS idx_affixes_rarity ON affixes(rarity);
CREATE INDEX IF NOT EXISTS idx_affixes_slot ON affixes(slot_type);

-- 装备实例表 (玩家获得的装备)
CREATE TABLE IF NOT EXISTS equipment_instance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id VARCHAR(32) NOT NULL,
    owner_id INTEGER NOT NULL,
    character_id INTEGER,                  -- NULL表示在背包中
    slot VARCHAR(16),
    quality VARCHAR(16) NOT NULL DEFAULT 'common',
    enhance_level INTEGER DEFAULT 0,       -- 强化等级 0-25
    evolution_stage INTEGER DEFAULT 1,     -- 进化阶段 1-5
    evolution_path VARCHAR(32),            -- 进化路线
    prefix_id VARCHAR(32),                 -- 前缀词缀ID
    prefix_value REAL,                     -- 前缀数值
    suffix_id VARCHAR(32),                 -- 后缀词缀ID
    suffix_value REAL,                     -- 后缀数值
    bonus_affix_1 VARCHAR(32),             -- 额外词缀1 (紫色+)
    bonus_affix_1_value REAL,
    bonus_affix_2 VARCHAR(32),             -- 额外词缀2 (橙色+)
    bonus_affix_2_value REAL,
    legendary_effect_id VARCHAR(32),       -- 传说效果ID
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_locked INTEGER DEFAULT 0,           -- 防误分解
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE SET NULL,
    FOREIGN KEY (prefix_id) REFERENCES affixes(id),
    FOREIGN KEY (suffix_id) REFERENCES affixes(id)
);

CREATE INDEX IF NOT EXISTS idx_equipment_owner ON equipment_instance(owner_id);
CREATE INDEX IF NOT EXISTS idx_equipment_character ON equipment_instance(character_id);
CREATE INDEX IF NOT EXISTS idx_equipment_quality ON equipment_instance(quality);

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
    physical_attack INTEGER NOT NULL,
    magic_attack INTEGER NOT NULL,
    physical_defense INTEGER NOT NULL,
    magic_defense INTEGER NOT NULL,
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

-- ═══════════════════════════════════════════════════════════
-- 角色成长系统
-- ═══════════════════════════════════════════════════════════

-- 属性分配表 - 记录玩家自由分配的属性点
CREATE TABLE IF NOT EXISTS character_stat_allocation (
    character_id INTEGER PRIMARY KEY,
    unspent_points INTEGER DEFAULT 0,           -- 未分配点数
    allocated_strength INTEGER DEFAULT 0,       -- 已分配力量
    allocated_agility INTEGER DEFAULT 0,        -- 已分配敏捷
    allocated_intellect INTEGER DEFAULT 0,      -- 已分配智力
    allocated_stamina INTEGER DEFAULT 0,        -- 已分配耐力
    allocated_spirit INTEGER DEFAULT 0,         -- 已分配精神
    respec_count INTEGER DEFAULT 0,             -- 重置次数
    last_respec_at DATETIME,                    -- 上次重置时间
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

-- 被动技能配置表
CREATE TABLE IF NOT EXISTS passive_skills (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    description TEXT,
    icon VARCHAR(64),
    class_id VARCHAR(32),                       -- NULL=通用
    rarity VARCHAR(16) DEFAULT 'common',        -- common/rare/epic
    tier INTEGER DEFAULT 1,                     -- 1基础/2进阶/3大师
    effect_type VARCHAR(32) NOT NULL,           -- 效果类型
    effect_value REAL NOT NULL,                 -- 效果数值
    effect_stat VARCHAR(32),                    -- 影响的属性
    max_level INTEGER DEFAULT 5,                -- 最大升级次数
    level_scaling REAL DEFAULT 0.2,             -- 每级提升比例(20%)
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

CREATE INDEX IF NOT EXISTS idx_passive_skills_class ON passive_skills(class_id);
CREATE INDEX IF NOT EXISTS idx_passive_skills_tier ON passive_skills(tier);

-- 角色被动技能表
CREATE TABLE IF NOT EXISTS character_passive_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    passive_id VARCHAR(32) NOT NULL,
    level INTEGER DEFAULT 1,                    -- 当前等级
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (passive_id) REFERENCES passive_skills(id),
    UNIQUE(character_id, passive_id)
);

CREATE INDEX IF NOT EXISTS idx_char_passive_char ON character_passive_skills(character_id);

-- 技能选择记录表 - 记录每3级的技能选择历史
CREATE TABLE IF NOT EXISTS skill_selection_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    level_milestone INTEGER NOT NULL,           -- 选择时等级(3,6,9...)
    offered_skill_1 VARCHAR(32) NOT NULL,       -- 选项1
    offered_skill_2 VARCHAR(32) NOT NULL,       -- 选项2
    offered_skill_3 VARCHAR(32) NOT NULL,       -- 选项3
    selected_skill_id VARCHAR(32) NOT NULL,     -- 选中的技能
    skill_was_upgrade INTEGER DEFAULT 0,        -- 是否为技能升级
    selected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    UNIQUE(character_id, level_milestone)
);

CREATE INDEX IF NOT EXISTS idx_skill_selection_char ON skill_selection_history(character_id);

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
-- 战斗数据分析系统
-- ═══════════════════════════════════════════════════════════

-- 战斗记录表 - 每场战斗的基本信息
CREATE TABLE IF NOT EXISTS battle_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    zone_id VARCHAR(32) NOT NULL,
    battle_type VARCHAR(16) NOT NULL,           -- pve/pvp/boss/abyss
    monster_id VARCHAR(32),                     -- PVE怪物ID
    opponent_user_id INTEGER,                   -- PVP对手ID
    total_rounds INTEGER DEFAULT 0,             -- 总回合数
    duration_seconds INTEGER DEFAULT 0,         -- 战斗时长
    result VARCHAR(16) NOT NULL,                -- victory/defeat/draw/flee
    team_damage_dealt INTEGER DEFAULT 0,        -- 队伍总输出
    team_damage_taken INTEGER DEFAULT 0,        -- 队伍总承伤
    team_healing_done INTEGER DEFAULT 0,        -- 队伍总治疗
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    battle_log TEXT,                            -- 详细战斗日志(JSON)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (zone_id) REFERENCES zones(id),
    FOREIGN KEY (monster_id) REFERENCES monsters(id),
    FOREIGN KEY (opponent_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_battle_records_user ON battle_records(user_id);
CREATE INDEX IF NOT EXISTS idx_battle_records_time ON battle_records(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_battle_records_zone ON battle_records(zone_id);
CREATE INDEX IF NOT EXISTS idx_battle_records_type ON battle_records(battle_type);

-- 战斗角色统计表 - 单场战斗中每个角色的详细数据
CREATE TABLE IF NOT EXISTS battle_character_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    team_slot INTEGER NOT NULL,
    -- 伤害统计
    damage_dealt INTEGER DEFAULT 0,             -- 造成总伤害
    physical_damage INTEGER DEFAULT 0,          -- 物理伤害
    magic_damage INTEGER DEFAULT 0,             -- 魔法伤害
    fire_damage INTEGER DEFAULT 0,              -- 火焰伤害
    frost_damage INTEGER DEFAULT 0,             -- 冰霜伤害
    shadow_damage INTEGER DEFAULT 0,            -- 暗影伤害
    holy_damage INTEGER DEFAULT 0,              -- 神圣伤害
    nature_damage INTEGER DEFAULT 0,            -- 自然伤害
    dot_damage INTEGER DEFAULT 0,               -- DOT伤害
    -- 暴击统计
    crit_count INTEGER DEFAULT 0,               -- 暴击次数
    crit_damage INTEGER DEFAULT 0,              -- 暴击总伤害
    max_crit INTEGER DEFAULT 0,                 -- 最高单次暴击
    -- 承伤统计
    damage_taken INTEGER DEFAULT 0,             -- 受到总伤害
    physical_taken INTEGER DEFAULT 0,           -- 物理承伤
    magic_taken INTEGER DEFAULT 0,              -- 魔法承伤
    damage_blocked INTEGER DEFAULT 0,           -- 格挡伤害
    damage_absorbed INTEGER DEFAULT 0,          -- 护盾吸收
    -- 闪避统计
    dodge_count INTEGER DEFAULT 0,              -- 闪避次数
    block_count INTEGER DEFAULT 0,              -- 格挡次数
    hit_count INTEGER DEFAULT 0,                -- 被命中次数
    -- 治疗统计
    healing_done INTEGER DEFAULT 0,             -- 造成治疗
    healing_received INTEGER DEFAULT 0,         -- 受到治疗
    overhealing INTEGER DEFAULT 0,              -- 过量治疗
    self_healing INTEGER DEFAULT 0,             -- 自我治疗
    hot_healing INTEGER DEFAULT 0,              -- HOT治疗
    -- 技能统计
    skill_uses INTEGER DEFAULT 0,               -- 技能使用次数
    skill_hits INTEGER DEFAULT 0,               -- 技能命中次数
    skill_misses INTEGER DEFAULT 0,             -- 技能未命中
    -- 控制统计
    cc_applied INTEGER DEFAULT 0,               -- 施加控制次数
    cc_received INTEGER DEFAULT 0,              -- 受到控制次数
    dispels INTEGER DEFAULT 0,                  -- 驱散次数
    interrupts INTEGER DEFAULT 0,               -- 打断次数
    -- 其他统计
    kills INTEGER DEFAULT 0,                    -- 击杀数
    deaths INTEGER DEFAULT 0,                   -- 死亡次数
    resurrects INTEGER DEFAULT 0,               -- 复活次数
    resource_used INTEGER DEFAULT 0,            -- 消耗能量
    resource_generated INTEGER DEFAULT 0,       -- 获得能量
    FOREIGN KEY (battle_id) REFERENCES battle_records(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_battle_char_stats_battle ON battle_character_stats(battle_id);
CREATE INDEX IF NOT EXISTS idx_battle_char_stats_char ON battle_character_stats(character_id);

-- 角色生涯统计表 - 累计统计数据
CREATE TABLE IF NOT EXISTS character_lifetime_stats (
    character_id INTEGER PRIMARY KEY,
    -- 战斗场次
    total_battles INTEGER DEFAULT 0,
    victories INTEGER DEFAULT 0,
    defeats INTEGER DEFAULT 0,
    pve_battles INTEGER DEFAULT 0,
    pvp_battles INTEGER DEFAULT 0,
    boss_kills INTEGER DEFAULT 0,
    -- 累计伤害
    total_damage_dealt INTEGER DEFAULT 0,
    total_physical_damage INTEGER DEFAULT 0,
    total_magic_damage INTEGER DEFAULT 0,
    total_crit_damage INTEGER DEFAULT 0,
    total_crit_count INTEGER DEFAULT 0,
    highest_damage_single INTEGER DEFAULT 0,
    highest_damage_battle INTEGER DEFAULT 0,
    -- 累计承伤
    total_damage_taken INTEGER DEFAULT 0,
    total_damage_blocked INTEGER DEFAULT 0,
    total_damage_absorbed INTEGER DEFAULT 0,
    total_dodge_count INTEGER DEFAULT 0,
    -- 累计治疗
    total_healing_done INTEGER DEFAULT 0,
    total_healing_received INTEGER DEFAULT 0,
    total_overhealing INTEGER DEFAULT 0,
    highest_healing_single INTEGER DEFAULT 0,
    highest_healing_battle INTEGER DEFAULT 0,
    -- 击杀与死亡
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    kill_streak_best INTEGER DEFAULT 0,
    current_kill_streak INTEGER DEFAULT 0,
    -- 技能使用
    total_skill_uses INTEGER DEFAULT 0,
    total_skill_hits INTEGER DEFAULT 0,
    -- 资源统计
    total_resource_used INTEGER DEFAULT 0,
    total_rounds INTEGER DEFAULT 0,
    total_battle_time INTEGER DEFAULT 0,
    -- 最后更新
    last_battle_at DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

-- 战斗技能明细表 - 每场战斗中各技能的使用和效果
CREATE TABLE IF NOT EXISTS battle_skill_breakdown (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id INTEGER NOT NULL,
    character_id INTEGER NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    use_count INTEGER DEFAULT 0,                -- 使用次数
    hit_count INTEGER DEFAULT 0,                -- 命中次数
    crit_count INTEGER DEFAULT 0,               -- 暴击次数
    total_damage INTEGER DEFAULT 0,             -- 造成总伤害
    total_healing INTEGER DEFAULT 0,            -- 造成总治疗
    resource_cost INTEGER DEFAULT 0,            -- 总消耗能量
    FOREIGN KEY (battle_id) REFERENCES battle_records(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

CREATE INDEX IF NOT EXISTS idx_skill_breakdown_battle ON battle_skill_breakdown(battle_id);
CREATE INDEX IF NOT EXISTS idx_skill_breakdown_char ON battle_skill_breakdown(character_id);

-- 每日统计汇总表 - 每日战斗数据快照
CREATE TABLE IF NOT EXISTS daily_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    stat_date DATE NOT NULL,
    battles_count INTEGER DEFAULT 0,
    victories INTEGER DEFAULT 0,
    defeats INTEGER DEFAULT 0,
    total_damage INTEGER DEFAULT 0,
    total_healing INTEGER DEFAULT 0,
    total_damage_taken INTEGER DEFAULT 0,
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    play_time INTEGER DEFAULT 0,
    kills INTEGER DEFAULT 0,
    deaths INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, stat_date)
);

CREATE INDEX IF NOT EXISTS idx_daily_stats_user ON daily_statistics(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_stats_date ON daily_statistics(stat_date DESC);

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

-- ═══════════════════════════════════════════════════════════
-- 装备系统：词缀 + 进化链
-- ═══════════════════════════════════════════════════════════

-- 装备实例表 (玩家获得的每件装备)
CREATE TABLE IF NOT EXISTS equipment_instance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id VARCHAR(32) NOT NULL,          -- 基础物品ID
    owner_id INTEGER NOT NULL,             -- 拥有者用户ID
    character_id INTEGER,                  -- 装备者角色ID (NULL=背包中)
    slot VARCHAR(16),                      -- 装备槽位
    quality VARCHAR(16) NOT NULL DEFAULT 'common', -- common/uncommon/rare/epic/legendary/mythic
    enhance_level INTEGER DEFAULT 0,       -- 强化等级 0-25
    evolution_stage INTEGER DEFAULT 1,     -- 进化阶段 1-5
    evolution_path VARCHAR(32),            -- 进化路线: fire/frost/lightning/holy/shadow/nature/physical
    prefix_id VARCHAR(32),                 -- 前缀词缀ID
    prefix_value REAL,                     -- 前缀数值
    suffix_id VARCHAR(32),                 -- 后缀词缀ID
    suffix_value REAL,                     -- 后缀数值
    bonus_affix_1 VARCHAR(32),             -- 额外词缀1 (紫色+)
    bonus_affix_1_value REAL,
    bonus_affix_2 VARCHAR(32),             -- 额外词缀2 (橙色+)
    bonus_affix_2_value REAL,
    legendary_effect_id VARCHAR(32),       -- 传说效果ID (红色品质)
    acquired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_locked INTEGER DEFAULT 0,           -- 是否锁定
    FOREIGN KEY (item_id) REFERENCES items(id),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE SET NULL,
    FOREIGN KEY (prefix_id) REFERENCES affixes(id),
    FOREIGN KEY (suffix_id) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_1) REFERENCES affixes(id),
    FOREIGN KEY (bonus_affix_2) REFERENCES affixes(id),
    FOREIGN KEY (legendary_effect_id) REFERENCES legendary_effects(id)
);

CREATE INDEX IF NOT EXISTS idx_equipment_owner ON equipment_instance(owner_id);
CREATE INDEX IF NOT EXISTS idx_equipment_character ON equipment_instance(character_id);
CREATE INDEX IF NOT EXISTS idx_equipment_quality ON equipment_instance(quality);

-- 词缀配置表
CREATE TABLE IF NOT EXISTS affixes (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,             -- 词缀名称
    type VARCHAR(16) NOT NULL,             -- prefix/suffix
    slot_type VARCHAR(16) DEFAULT 'all',   -- weapon/armor/accessory/all
    rarity VARCHAR(16) NOT NULL DEFAULT 'common', -- common/uncommon/rare/epic
    effect_type VARCHAR(32) NOT NULL,      -- 效果类型
    effect_stat VARCHAR(32),               -- 影响的属性
    min_value REAL NOT NULL,               -- 最小数值
    max_value REAL NOT NULL,               -- 最大数值
    value_type VARCHAR(16) NOT NULL DEFAULT 'flat', -- flat/percent
    description TEXT,                      -- 描述模板
    level_required INTEGER DEFAULT 1       -- 最低出现等级
);

CREATE INDEX IF NOT EXISTS idx_affixes_type ON affixes(type);
CREATE INDEX IF NOT EXISTS idx_affixes_rarity ON affixes(rarity);
CREATE INDEX IF NOT EXISTS idx_affixes_slot ON affixes(slot_type);

-- 进化路线配置表
CREATE TABLE IF NOT EXISTS evolution_paths (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,             -- 路线名称
    element VARCHAR(16) NOT NULL,          -- 元素类型
    description TEXT,
    slot_type VARCHAR(16) NOT NULL,        -- weapon/armor
    stat_bonus_type VARCHAR(32),           -- 属性加成类型
    stat_bonus_value REAL,                 -- 属性加成数值
    special_effect TEXT,                   -- 特殊效果描述
    material_required TEXT                 -- 所需材料 (JSON)
);

-- 传说效果表
CREATE TABLE IF NOT EXISTS legendary_effects (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,             -- 效果名称
    description TEXT NOT NULL,             -- 效果描述
    slot_type VARCHAR(16) NOT NULL,        -- weapon/armor/accessory
    evolution_path VARCHAR(32),            -- 关联进化路线
    trigger_type VARCHAR(32),              -- on_hit/on_kill/on_damaged/passive
    trigger_chance REAL DEFAULT 1.0,       -- 触发概率
    effect_type VARCHAR(32) NOT NULL,      -- 效果类型
    effect_value REAL,                     -- 效果数值
    cooldown INTEGER DEFAULT 0,            -- 冷却回合
    FOREIGN KEY (evolution_path) REFERENCES evolution_paths(id)
);

-- 掉落配置表
CREATE TABLE IF NOT EXISTS drop_config (
    id VARCHAR(32) PRIMARY KEY,
    monster_type VARCHAR(16) NOT NULL,     -- normal/elite/boss/abyss_boss
    base_drop_rate REAL NOT NULL,          -- 基础掉落率
    quality_weights TEXT NOT NULL,         -- 品质权重 (JSON)
    miracle_rate REAL DEFAULT 0,           -- 奇迹掉落率
    pity_threshold INTEGER DEFAULT 40,     -- 保底触发次数
    pity_min_quality VARCHAR(16) DEFAULT 'rare' -- 保底最低品质
);

-- 玩家保底计数表
CREATE TABLE IF NOT EXISTS user_drop_pity (
    user_id INTEGER PRIMARY KEY,
    no_drop_count INTEGER DEFAULT 0,       -- 连续无掉落次数
    last_drop_at DATETIME,                 -- 上次掉落时间
    total_drops INTEGER DEFAULT 0,         -- 总掉落次数
    miracle_drops INTEGER DEFAULT 0,       -- 奇迹掉落次数
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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

-- ═══════════════════════════════════════════════════════════
-- 体力系统
-- ═══════════════════════════════════════════════════════════

-- 玩家体力表
CREATE TABLE IF NOT EXISTS user_stamina (
    user_id INTEGER PRIMARY KEY,
    current_stamina INTEGER DEFAULT 100,   -- 当前体力
    max_stamina INTEGER DEFAULT 100,       -- 最大体力
    last_regen_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 上次恢复时间
    overflow_exp INTEGER DEFAULT 0,        -- 溢出转化的经验
    overflow_gold INTEGER DEFAULT 0,       -- 溢出转化的金币
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ═══════════════════════════════════════════════════════════
-- 作战策略系统
-- ═══════════════════════════════════════════════════════════

-- 战斗策略表
CREATE TABLE IF NOT EXISTS battle_strategies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    name VARCHAR(32) NOT NULL,             -- 策略名称
    is_active INTEGER DEFAULT 0,           -- 是否当前使用
    skill_priority TEXT,                   -- 技能优先级 (JSON数组)
    conditional_rules TEXT,                -- 条件规则 (JSON数组)
    target_priority VARCHAR(32) DEFAULT 'lowest_hp', -- 目标选择策略
    resource_threshold INTEGER DEFAULT 0,  -- 资源阈值
    reserved_skills TEXT,                  -- 保留技能 (JSON数组)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_battle_strategies_character ON battle_strategies(character_id);
CREATE INDEX IF NOT EXISTS idx_battle_strategies_active ON battle_strategies(character_id, is_active);

-- ═══════════════════════════════════════════════════════════
-- 战斗数据分析系统
-- ═══════════════════════════════════════════════════════════

-- 详细战斗日志表 (用于"上一场"分析)
CREATE TABLE IF NOT EXISTS detailed_battle_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id VARCHAR(64) NOT NULL,        -- 战斗ID
    user_id INTEGER NOT NULL,
    zone_id VARCHAR(32),
    battle_type VARCHAR(16) NOT NULL,      -- pve/pvp/abyss
    result VARCHAR(16) NOT NULL,           -- victory/defeat/draw
    total_turns INTEGER NOT NULL,          -- 总回合数
    duration_seconds INTEGER,              -- 战斗时长
    player_team_data TEXT NOT NULL,        -- 我方队伍数据 (JSON)
    enemy_team_data TEXT NOT NULL,         -- 敌方队伍数据 (JSON)
    turn_logs TEXT NOT NULL,               -- 回合日志 (JSON)
    exp_gained INTEGER DEFAULT 0,
    gold_gained INTEGER DEFAULT 0,
    items_dropped TEXT,                    -- 掉落物品 (JSON)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (zone_id) REFERENCES zones(id)
);

CREATE INDEX IF NOT EXISTS idx_detailed_logs_user ON detailed_battle_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_detailed_logs_time ON detailed_battle_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_detailed_logs_battle ON detailed_battle_logs(battle_id);

-- 角色战斗统计表 (汇总统计)
CREATE TABLE IF NOT EXISTS character_battle_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    stat_date DATE NOT NULL,               -- 统计日期
    battles_total INTEGER DEFAULT 0,       -- 总战斗场次
    battles_won INTEGER DEFAULT 0,         -- 胜利场次
    battles_lost INTEGER DEFAULT 0,        -- 失败场次
    total_damage_dealt INTEGER DEFAULT 0,  -- 总造成伤害
    total_damage_taken INTEGER DEFAULT 0,  -- 总承受伤害
    total_healing_done INTEGER DEFAULT 0,  -- 总治疗量
    total_turns INTEGER DEFAULT 0,         -- 总回合数
    deaths INTEGER DEFAULT 0,              -- 死亡次数
    kills INTEGER DEFAULT 0,               -- 击杀数
    crits INTEGER DEFAULT 0,               -- 暴击次数
    dodges INTEGER DEFAULT 0,              -- 闪避次数
    skills_used TEXT,                      -- 技能使用统计 (JSON)
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    UNIQUE(character_id, stat_date)
);

CREATE INDEX IF NOT EXISTS idx_char_stats_date ON character_battle_stats(character_id, stat_date DESC);

-- 战斗仇恨日志 (仇恨系统分析)
CREATE TABLE IF NOT EXISTS battle_threat_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    battle_id VARCHAR(64) NOT NULL,
    turn INTEGER NOT NULL,             -- 回合数
    enemy_id VARCHAR(32) NOT NULL,     -- 敌人单位ID
    threat_snapshot TEXT NOT NULL,     -- 仇恨快照 (JSON: {current_target, threat_list[]})
    target_changed INTEGER DEFAULT 0,  -- 是否切换目标: 1是 0否
    ot_occurred INTEGER DEFAULT 0,     -- 是否发生OT: 1是 0否
    taunt_used INTEGER DEFAULT 0,      -- 是否使用嘲讽: 1是 0否
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_threat_log_battle ON battle_threat_log(battle_id);
CREATE INDEX IF NOT EXISTS idx_threat_log_enemy ON battle_threat_log(battle_id, enemy_id);

-- ═══════════════════════════════════════════════════════════
-- 聊天系统表
-- ═══════════════════════════════════════════════════════════

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel VARCHAR(16) NOT NULL,          -- world/zone/trade/lfg/whisper/system/battlefield
    faction VARCHAR(16),                   -- alliance/horde (公共频道用)
    zone_id VARCHAR(32),                   -- 区域ID (zone频道用)
    sender_id INTEGER NOT NULL,
    sender_name VARCHAR(32) NOT NULL,
    sender_class VARCHAR(32),
    receiver_id INTEGER,                   -- 私聊目标用户ID
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_chat_channel ON chat_messages(channel, faction);
CREATE INDEX IF NOT EXISTS idx_chat_zone ON chat_messages(zone_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_whisper ON chat_messages(sender_id, receiver_id);
CREATE INDEX IF NOT EXISTS idx_chat_time ON chat_messages(created_at DESC);

-- 屏蔽列表
CREATE TABLE IF NOT EXISTS chat_blocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    blocked_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_blocks_user ON chat_blocks(user_id);

-- 在线状态表
CREATE TABLE IF NOT EXISTS user_online_status (
    user_id INTEGER PRIMARY KEY,
    character_id INTEGER,
    character_name VARCHAR(32),
    faction VARCHAR(16),
    zone_id VARCHAR(32),
    last_active DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_online INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES characters(id)
);

CREATE INDEX IF NOT EXISTS idx_online_faction ON user_online_status(faction, is_online);
CREATE INDEX IF NOT EXISTS idx_online_zone ON user_online_status(zone_id, is_online);

