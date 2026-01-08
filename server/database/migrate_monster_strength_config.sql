-- ═══════════════════════════════════════════════════════════
-- 怪物强度配置表迁移
-- 用于快速调整怪物强度，无需修改seed.sql
-- ═══════════════════════════════════════════════════════════

-- 怪物强度配置表
CREATE TABLE IF NOT EXISTS monster_strength_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_min INTEGER NOT NULL,              -- 等级下限（包含）
    level_max INTEGER NOT NULL,              -- 等级上限（包含）
    hp_multiplier REAL DEFAULT 1.0,          -- 生命值倍数
    attack_multiplier REAL DEFAULT 1.0,      -- 攻击力倍数（物理和魔法）
    defense_multiplier REAL DEFAULT 1.0,     -- 防御力倍数（物理和魔法）
    crit_rate_bonus REAL DEFAULT 0.0,        -- 暴击率加成（绝对值，如0.02表示+2%）
    description TEXT,                        -- 描述
    is_active INTEGER DEFAULT 1,             -- 是否启用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(level_min, level_max)
);

CREATE INDEX IF NOT EXISTS idx_strength_config_level ON monster_strength_config(level_min, level_max);
CREATE INDEX IF NOT EXISTS idx_strength_config_active ON monster_strength_config(is_active);

-- 插入默认配置（当前已应用的强度提升）
INSERT OR IGNORE INTO monster_strength_config (level_min, level_max, hp_multiplier, attack_multiplier, defense_multiplier, crit_rate_bonus, description) VALUES
(1, 10, 1.5, 1.4, 1.4, 0.02, '1-10级：HP +50%, 攻击 +40%, 防御 +40%, 暴击率 +2%'),
(11, 20, 1.45, 1.35, 1.35, 0.02, '10-20级：HP +45%, 攻击 +35%, 防御 +35%, 暴击率 +2%'),
(21, 30, 1.4, 1.35, 1.35, 0.02, '20-30级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +2%'),
(31, 40, 1.4, 1.35, 1.35, 0.03, '30-40级：HP +40%, 攻击 +35%, 防御 +35%, 暴击率 +3%'),
(41, 50, 1.35, 1.3, 1.3, 0.03, '40-50级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%'),
(51, 60, 1.35, 1.3, 1.3, 0.03, '50-60级：HP +35%, 攻击 +30%, 防御 +30%, 暴击率 +3%');












































































