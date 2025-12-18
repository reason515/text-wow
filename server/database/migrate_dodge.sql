-- ═══════════════════════════════════════════════════════════
-- 闪避系统迁移脚本
-- 添加 dodge_rate 字段到 characters 和 monsters 表
-- ═══════════════════════════════════════════════════════════

-- 为角色表添加闪避率字段
ALTER TABLE characters ADD COLUMN dodge_rate REAL DEFAULT 0.05;

-- 为怪物表添加闪避率字段
ALTER TABLE monsters ADD COLUMN dodge_rate REAL DEFAULT 0.05;

-- 更新现有角色的闪避率（基于敏捷计算：5% + 敏捷/20）
-- 注意：这是初始迁移，实际闪避率会在游戏运行时动态计算
UPDATE characters SET dodge_rate = 0.05 + (agility / 2000.0) WHERE dodge_rate IS NULL OR dodge_rate = 0.05;

-- 更新现有怪物的闪避率（默认5%）
UPDATE monsters SET dodge_rate = 0.05 WHERE dodge_rate IS NULL;

-- ═══════════════════════════════════════════════════════════
-- 使用说明：
-- 在 SQLite 中执行此脚本：
-- sqlite3 game.db < migrate_dodge.sql
-- ═══════════════════════════════════════════════════════════
























