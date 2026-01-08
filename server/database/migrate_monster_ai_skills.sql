-- ═══════════════════════════════════════════════════════════
-- 怪物AI和技能系统迁移脚本
-- ═══════════════════════════════════════════════════════════

-- 扩展monsters表，添加AI配置和技能配置字段
ALTER TABLE monsters ADD COLUMN IF NOT EXISTS ai_type VARCHAR(32) DEFAULT 'balanced';
ALTER TABLE monsters ADD COLUMN IF NOT EXISTS ai_behavior TEXT;  -- JSON格式的AI行为配置
ALTER TABLE monsters ADD COLUMN IF NOT EXISTS speed INTEGER DEFAULT 10;  -- 行动速度
ALTER TABLE monsters ADD COLUMN IF NOT EXISTS balance_version INTEGER DEFAULT 1;  -- 平衡版本号

-- 创建怪物技能表
CREATE TABLE IF NOT EXISTS monster_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    monster_id VARCHAR(32) NOT NULL,
    skill_id VARCHAR(32) NOT NULL,
    skill_type VARCHAR(32) NOT NULL,  -- attack/defense/control/heal/special
    priority INTEGER DEFAULT 0,  -- 优先级（数字越大优先级越高）
    cooldown INTEGER DEFAULT 0,  -- 冷却时间（回合数）
    use_condition TEXT,  -- JSON格式，使用条件
    FOREIGN KEY (monster_id) REFERENCES monsters(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id)
);

CREATE INDEX IF NOT EXISTS idx_monster_skills_monster_id ON monster_skills(monster_id);
CREATE INDEX IF NOT EXISTS idx_monster_skills_skill_id ON monster_skills(skill_id);

-- 创建怪物AI配置表（可选，用于更复杂的AI配置）
CREATE TABLE IF NOT EXISTS monster_ai_configs (
    id VARCHAR(32) PRIMARY KEY,
    monster_id VARCHAR(32) NOT NULL,
    ai_type VARCHAR(32) NOT NULL,  -- aggressive/defensive/balanced/special
    behavior_config TEXT NOT NULL,  -- JSON格式
    FOREIGN KEY (monster_id) REFERENCES monsters(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_monster_ai_configs_monster_id ON monster_ai_configs(monster_id);






















































