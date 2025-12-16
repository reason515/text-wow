-- ═══════════════════════════════════════════════════════════
-- 怪物暴击与攻击类型迁移
-- 用途：为 monsters 表添加暴击字段和攻击类型，用于敌人暴击与魔法攻击
-- 使用方式：sqlite3 game.db < migrate_monster_crit.sql
-- ═══════════════════════════════════════════════════════════

ALTER TABLE monsters ADD COLUMN attack_type VARCHAR(16) DEFAULT 'physical';
ALTER TABLE monsters ADD COLUMN phys_crit_rate REAL DEFAULT 0.05;
ALTER TABLE monsters ADD COLUMN phys_crit_damage REAL DEFAULT 1.5;
ALTER TABLE monsters ADD COLUMN spell_crit_rate REAL DEFAULT 0.05;
ALTER TABLE monsters ADD COLUMN spell_crit_damage REAL DEFAULT 1.5;

-- 为已存在的数据填充默认值
UPDATE monsters
SET attack_type = COALESCE(attack_type, 'physical'),
    phys_crit_rate = COALESCE(phys_crit_rate, 0.05),
    phys_crit_damage = COALESCE(phys_crit_damage, 1.5),
    spell_crit_rate = COALESCE(spell_crit_rate, 0.05),
    spell_crit_damage = COALESCE(spell_crit_damage, 1.5);




















