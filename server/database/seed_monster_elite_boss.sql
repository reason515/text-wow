-- ═══════════════════════════════════════════════════════════
-- 精英和Boss怪物数据
-- ═══════════════════════════════════════════════════════════

-- 首先需要确保怪物技能表存在
-- 这个脚本应该在 migrate_monster_ai_skills.sql 之后执行

-- ═══════════════════════════════════════════════════════════
-- 精英怪物 (Elite) - 艾尔文森林
-- ═══════════════════════════════════════════════════════════

-- 精英狼人 - 高攻击+撕裂技能
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'elite_werewolf', 'elwynn', '精英狼人', '狂暴的狼人，攻击力极强', 8, 'elite',
    120, 50, 25, 5, 8, 6,
    'physical', 0.10, 1.8, 0.05, 1.5,
    0.08, 12, 40, 8, 15, 10,
    'aggressive', '{"target_priority": ["lowest_hp", "lowest_defense"], "skill_priority": ["high_damage", "execute"], "defense_threshold": 0.2, "random_factor": 0.1}'
);

-- 精英法师 - 法术攻击+护盾技能
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'elite_mage', 'elwynn', '精英法师', '强大的奥术法师，擅长法术攻击', 10, 'elite',
    80, 150, 5, 30, 4, 12,
    'magic', 0.05, 1.5, 0.15, 2.0,
    0.05, 8, 50, 10, 20, 8,
    'balanced', '{"target_priority": ["lowest_hp", "random"], "skill_priority": ["defense", "high_damage"], "defense_threshold": 0.4, "random_factor": 0.2}'
);

-- 精英治疗者 - 低攻击+治疗技能
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'elite_healer', 'elwynn', '精英治疗者', '拥有强大治疗能力的怪物', 9, 'elite',
    100, 120, 8, 15, 10, 10,
    'magic', 0.05, 1.5, 0.10, 1.8,
    0.06, 10, 45, 9, 18, 5,
    'defensive', '{"target_priority": ["highest_threat"], "skill_priority": ["heal", "defense", "attack"], "defense_threshold": 0.5, "random_factor": 0.1}'
);

-- ═══════════════════════════════════════════════════════════
-- Boss怪物 (Boss) - 艾尔文森林
-- ═══════════════════════════════════════════════════════════

-- 森林之王 - 高HP+召唤小怪+范围攻击
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'boss_forest_king', 'elwynn', '森林之王', '统治森林的强大存在，拥有召唤和范围攻击能力', 15, 'boss',
    500, 200, 40, 25, 20, 15,
    'physical', 0.12, 2.0, 0.10, 2.0,
    0.10, 15, 150, 50, 100, 1,
    'special', '{"target_priority": ["highest_threat", "lowest_hp"], "skill_priority": ["special", "high_damage"], "defense_threshold": 0.3, "random_factor": 0.05, "phases": [{"hp_threshold": 1.0, "behavior": "aggressive", "skills": ["boss_summon", "boss_cleave"]}, {"hp_threshold": 0.5, "behavior": "defensive", "skills": ["boss_heal", "boss_rage"]}]}'
);

-- 暗影法师Boss - 法术攻击+控制技能+护盾
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'boss_shadow_mage', 'elwynn', '暗影法师', '掌握暗影魔法的强大法师，擅长控制和护盾', 18, 'boss',
    350, 300, 10, 50, 12, 25,
    'magic', 0.08, 1.8, 0.20, 2.5,
    0.12, 12, 200, 80, 150, 1,
    'special', '{"target_priority": ["lowest_hp", "random"], "skill_priority": ["control", "defense", "high_damage"], "defense_threshold": 0.4, "random_factor": 0.05, "phases": [{"hp_threshold": 1.0, "behavior": "balanced", "skills": ["boss_shadow_bolt", "boss_shield"]}, {"hp_threshold": 0.3, "behavior": "aggressive", "skills": ["boss_mind_control", "boss_shadow_nova"]}]}'
);

-- ═══════════════════════════════════════════════════════════
-- 特殊怪物 (Special) - 艾尔文森林
-- ═══════════════════════════════════════════════════════════

-- 暗影幽灵 - 物理攻击无效，只能用法术攻击
INSERT OR REPLACE INTO monsters (
    id, zone_id, name, description, level, type, hp, mp, 
    physical_attack, magic_attack, physical_defense, magic_defense,
    attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
    dodge_rate, speed, exp_reward, gold_min, gold_max, spawn_weight,
    ai_type, ai_behavior
) VALUES (
    'special_shadow_wraith', 'elwynn', '暗影幽灵', '物理攻击无效的幽灵，只能用法术攻击', 12, 'special',
    150, 100, 0, 35, 999, 8,  -- 物理防御极高，魔法防御正常
    'magic', 0.0, 1.0, 0.15, 2.0,
    0.15, 14, 80, 30, 60, 3,
    'aggressive', '{"target_priority": ["lowest_hp"], "skill_priority": ["high_damage"], "defense_threshold": 0.3, "random_factor": 0.1}'
);

-- ═══════════════════════════════════════════════════════════
-- 怪物技能配置
-- ═══════════════════════════════════════════════════════════

-- 注意：这些技能ID需要在skills表中存在，如果不存在需要先创建

-- 精英狼人技能：撕裂（持续伤害）
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('elite_werewolf', 'monster_rend', 'attack', 3, 3, '{"target_hp_min": 0.3}');

-- 精英法师技能：护盾
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('elite_mage', 'monster_shield', 'defense', 5, 5, '{"hp_max": 0.6}'),
('elite_mage', 'monster_fireball', 'attack', 2, 2, NULL);

-- 精英治疗者技能：治疗
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('elite_healer', 'monster_heal', 'heal', 5, 4, '{"hp_max": 0.7}');

-- Boss森林之王技能：召唤、范围攻击
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('boss_forest_king', 'boss_summon', 'special', 4, 6, '{"hp_max": 0.8}'),
('boss_forest_king', 'boss_cleave', 'attack', 3, 4, NULL),
('boss_forest_king', 'boss_heal', 'heal', 5, 5, '{"hp_max": 0.5}'),
('boss_forest_king', 'boss_rage', 'buff', 4, 8, '{"hp_max": 0.5}');

-- Boss暗影法师技能：暗影箭、护盾、控制、范围攻击
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('boss_shadow_mage', 'boss_shadow_bolt', 'attack', 3, 2, NULL),
('boss_shadow_mage', 'boss_shield', 'defense', 5, 6, '{"hp_max": 0.7}'),
('boss_shadow_mage', 'boss_mind_control', 'control', 4, 8, '{"hp_max": 0.4}'),
('boss_shadow_mage', 'boss_shadow_nova', 'attack', 4, 5, '{"hp_max": 0.3}');

-- 特殊怪物暗影幽灵技能：暗影箭
INSERT OR REPLACE INTO monster_skills (monster_id, skill_id, skill_type, priority, cooldown, use_condition) VALUES
('special_shadow_wraith', 'monster_shadow_bolt', 'attack', 2, 2, NULL);







































