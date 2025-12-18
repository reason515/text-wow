-- 迁移脚本：将暴击系统从单一暴击改为物理/法术分离
-- 运行方式: sqlite3 game.db < migrate_crit.sql

-- 添加新的暴击字段
ALTER TABLE characters ADD COLUMN phys_crit_rate REAL DEFAULT 0.05;
ALTER TABLE characters ADD COLUMN phys_crit_damage REAL DEFAULT 1.5;
ALTER TABLE characters ADD COLUMN spell_crit_rate REAL DEFAULT 0.05;
ALTER TABLE characters ADD COLUMN spell_crit_damage REAL DEFAULT 1.5;

-- 从旧字段迁移数据（如果存在）
UPDATE characters SET 
    phys_crit_rate = COALESCE(crit_rate, 0.05),
    phys_crit_damage = COALESCE(crit_damage, 1.5),
    spell_crit_rate = COALESCE(crit_rate, 0.05),
    spell_crit_damage = COALESCE(crit_damage, 1.5);

-- 注意：SQLite 不支持 DROP COLUMN，旧字段会保留但不再使用
























