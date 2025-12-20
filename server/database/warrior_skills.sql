-- ═══════════════════════════════════════════════════════════
-- 战士职业技能数据
-- 根据 docs/warrior_skills_design.md 实现
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 战士技能所需的效果 (Effects)
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO effects (id, name, description, type, is_buff, is_stackable, max_stacks, duration, tick_interval, value_type, value, stat_affected, damage_type, can_dispel, tags) VALUES
-- 嘲讽效果
('eff_taunt', '嘲讽', '强制攻击施法者', 'taunt', 0, 0, 1, 1, 1, NULL, NULL, NULL, NULL, 0, '["control"]'),
-- 盾牌格挡效果 (30%减伤，持续2回合，1级)
('eff_shield_block', '盾牌格挡', '减少受到的物理伤害30%', 'stat_mod', 1, 0, 1, 2, 1, 'percent', -30, 'physical_damage_taken', NULL, 0, '["defensive"]'),
-- 战斗怒吼效果 (10%攻击力，持续5回合，1级)
('eff_battle_shout', '战斗怒吼', '提升物理攻击力10%', 'stat_mod', 1, 0, 1, 5, 1, 'percent', 10, 'attack', NULL, 1, '["buff"]'),
-- 挫志怒吼效果 (15%降低攻击力，持续3回合，1级)
('eff_demoralizing_shout', '挫志怒吼', '降低攻击力15%', 'stat_mod', 0, 0, 1, 3, 1, 'percent', -15, 'attack', NULL, 1, '["debuff"]'),
-- 破釜沉舟效果 (立即恢复30%最大HP，持续3回合，1级)
('eff_last_stand', '破釜沉舟', '立即恢复最大HP的30%', 'heal', 1, 0, 1, 3, 1, 'percent', 30, 'max_hp', NULL, 0, '["survival"]'),
-- 冲锋眩晕效果 (30%概率眩晕1回合，1级)
('eff_charge_stun', '冲锋眩晕', '眩晕1回合', 'stun', 0, 0, 1, 1, 1, NULL, NULL, NULL, NULL, 0, '["control"]'),
-- 旋风斩防御降低效果 (10%降低物理防御，持续2回合，1级)
('eff_whirlwind_debuff', '旋风斩', '降低物理防御10%', 'stat_mod', 0, 0, 1, 2, 1, 'percent', -10, 'physical_defense', NULL, 1, '["debuff"]'),
-- 盾墙效果 (60%减伤，持续2回合，1级)
('eff_shield_wall', '盾墙', '减少受到的伤害60%', 'stat_mod', 1, 0, 1, 2, 1, 'percent', -60, 'damage_taken', NULL, 0, '["defensive"]'),
-- 不灭壁垒护盾效果 (50%最大HP护盾，持续4回合，1级)
('eff_unbreakable_barrier', '不灭壁垒', '获得相当于最大HP50%的护盾', 'shield', 1, 0, 1, 4, 1, 'percent', 50, 'max_hp', NULL, 0, '["defensive", "shield"]'),
-- 盾牌反射效果 (50%反射伤害，持续2回合，1级)
('eff_shield_reflection', '盾牌反射', '反射50%受到的伤害', 'reflect', 1, 0, 1, 2, 1, 'percent', 50, NULL, NULL, 0, '["defensive", "reflect"]'),
-- 挑战怒吼效果 (强制所有敌人攻击自己，持续1回合，1级)
('eff_challenging_shout', '挑战怒吼', '强制所有敌人攻击自己', 'taunt', 0, 0, 1, 1, 1, NULL, NULL, NULL, NULL, 0, '["control", "aoe"]'),
-- 鲁莽效果 (50%暴击率，+20%受到伤害，持续3回合，1级)
('eff_recklessness', '鲁莽', '提升50%暴击率，但受到伤害增加20%', 'stat_mod', 1, 0, 1, 3, 1, 'percent', 50, 'crit_rate', NULL, 0, '["offensive", "tradeoff"]'),
('eff_recklessness_debuff', '鲁莽副作用', '受到伤害增加20%', 'stat_mod', 0, 0, 1, 3, 1, 'percent', 20, 'damage_taken', NULL, 0, '["debuff"]'),
-- 反击风暴效果 (50%反击伤害，持续3回合，1级)
('eff_retaliation', '反击风暴', '每次受到攻击时反击50%物理攻击力', 'counter_attack', 1, 0, 1, 3, 1, 'percent', 50, NULL, NULL, 0, '["defensive", "counter"]'),
-- 狂暴之怒效果 (30%攻击力，每次攻击额外获得5点怒气，持续4回合，1级)
('eff_berserker_rage', '狂暴之怒', '提升30%攻击力，每次攻击额外获得5点怒气', 'stat_mod', 1, 0, 1, 4, 1, 'percent', 30, 'attack', NULL, 0, '["offensive", "rage_generation"]'),
-- 天神下凡效果 (50%攻击力，免疫控制，持续3回合，1级)
('eff_avatar', '天神下凡', '提升50%攻击力，免疫控制效果', 'stat_mod', 1, 0, 1, 3, 1, 'percent', 50, 'attack', NULL, 0, '["offensive", "immune_cc"]'),
-- 致死打击治疗效果降低 (50%降低治疗效果，持续3回合，1级)
('eff_mortal_strike', '致死打击', '降低受到的治疗效果50%', 'stat_mod', 0, 0, 1, 3, 1, 'percent', -50, 'healing_received', NULL, 1, '["debuff", "anti_heal"]');

-- ═══════════════════════════════════════════════════════════
-- 战士主动技能 (22个)
-- ═══════════════════════════════════════════════════════════
-- 注意：伤害计算为 base_value + scaling_stat * scaling_ratio
-- 设计文档中的百分比伤害，base_value=0，scaling_ratio=1.0表示100%
-- 升级效果在代码中处理，数据库只存储1级基础值

-- 初始技能池 (9个)
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance, threat_modifier, threat_type, tags) VALUES
-- A1. 英勇打击 (Heroic Strike)
('warrior_heroic_strike', '英勇打击', '基础攻击技能，造成物理伤害并产生较高仇恨。1级：100%物理攻击力，每级+15%伤害', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.0, 15, 0, 1, NULL, 1.0, 1.2, 'normal', '["basic", "rage_builder"]'),
-- A2. 嘲讽 (Taunt)
('warrior_taunt', '嘲讽', '核心坦克技能，强制敌人攻击自己。1级：冷却2回合，每级-0.5回合冷却', 'warrior', 'control', 'enemy', NULL, 0, NULL, 0, 0, 2, 1, 'eff_taunt', 1.0, 3.0, 'taunt', '["tank", "threat", "control"]'),
-- A3. 盾牌格挡 (Shield Block)
('warrior_shield_block', '盾牌格挡', '提升防御能力，减少受到的伤害。1级：减少30%物理伤害，持续2回合，每级+5%减伤', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 10, 3, 1, 'eff_shield_block', 1.0, 0.5, 'normal', '["defensive", "survival"]'),
-- A4. 顺劈斩 (Cleave)
('warrior_cleave', '顺劈斩', '对主目标造成伤害，同时攻击主目标相邻的敌人（最多2个）。1级：主目标100%，相邻80%，每级+10%', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.0, 20, 0, 1, NULL, 1.0, 1.3, 'normal', '["aoe", "positional", "rage_builder"]'),
-- A5. 重击 (Slam)
('warrior_slam', '重击', '强力单体攻击，造成大量伤害和仇恨。1级：150%物理攻击力，每级+20%伤害', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.5, 25, 0, 1, NULL, 1.0, 1.5, 'normal', '["high_damage", "rage_spender"]'),
-- A6. 战斗怒吼 (Battle Shout)
('warrior_battle_shout', '战斗怒吼', '团队增益技能，提升整体输出。1级：提升所有友方10%物理攻击力，持续5回合，每级+2%攻击力，+1回合持续时间', 'warrior', 'buff', 'ally_all', NULL, 0, NULL, 0, 10, 0, 1, 'eff_battle_shout', 1.0, 0.3, 'normal', '["buff", "support", "team"]'),
-- A7. 挫志怒吼 (Demoralizing Shout)
('warrior_demoralizing_shout', '挫志怒吼', '降低敌人伤害，保护团队。1级：降低所有敌人15%攻击力，持续3回合，每级+3%降低比例，+0.5回合持续时间', 'warrior', 'debuff', 'enemy_all', NULL, 0, NULL, 0, 10, 0, 1, 'eff_demoralizing_shout', 1.0, 1.1, 'normal', '["debuff", "aoe", "defensive"]'),
-- A8. 破釜沉舟 (Last Stand)
('warrior_last_stand', '破釜沉舟', '紧急恢复技能，关键时刻保命。1级：立即恢复最大HP的30%，持续3回合，每级+5%恢复量，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 10, 1, 'eff_last_stand', 1.0, 0.0, 'normal', '["survival", "emergency"]'),
-- A9. 冲锋 (Charge)
('warrior_charge', '冲锋', '战士的招牌技能！快速冲向敌人并造成伤害，获得怒气，有概率眩晕目标。1级：80%伤害，+15怒气，30%眩晕，冷却5回合，每级+10%伤害，+3怒气，+5%眩晕概率，-0.5回合冷却', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 0.8, 0, 5, 1, 'eff_charge_stun', 0.3, 1.2, 'normal', '["signature", "mobility", "rage_generation", "control"]'),

-- 所有主动技能池 (剩余13个)
-- A10. 旋风斩 (Whirlwind)
('warrior_whirlwind', '旋风斩', '强力AOE技能，对所有敌人造成伤害并降低其防御。1级：100%伤害，降低10%物理防御持续2回合，每级+15%伤害，+2%降低比例，+0.5回合持续时间', 'warrior', 'attack', 'enemy_all', 'physical', 0, 'strength', 1.0, 25, 2, 1, 'eff_whirlwind_debuff', 1.0, 1.4, 'normal', '["aoe", "high_damage", "debuff", "rage_spender"]'),
-- A11. 盾牌猛击 (Shield Slam)
('warrior_shield_slam', '盾牌猛击', '造成高额伤害，伤害基于攻击力和防御力。1级：120%攻击力+50%防御力，每级+15%攻击力加成，+10%防御力加成', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.2, 20, 2, 1, NULL, 1.0, 1.8, 'normal', '["defensive_offensive", "high_damage"]'),
-- A12. 斩杀 (Execute)
('warrior_execute', '斩杀', '终结技能，对低血量敌人造成巨额伤害。1级：200%物理攻击力（仅对HP<20%），每级+30%伤害', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 2.0, 30, 0, 1, NULL, 1.0, 1.6, 'normal', '["execute", "finisher", "high_damage"]'),
-- A13. 盾墙 (Shield Wall)
('warrior_shield_wall', '盾墙', '强大的防御技能，大幅减少伤害。1级：减少60%受到的伤害，持续2回合，每级+5%减伤，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 8, 1, 'eff_shield_wall', 1.0, 0.8, 'normal', '["defensive", "cooldown"]'),
-- A14. 不灭壁垒 (Unbreakable Barrier)
('warrior_unbreakable_barrier', '不灭壁垒', '终极防御技能。1级：获得相当于最大HP50%的护盾，持续4回合，每级+10%护盾量，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 12, 1, 'eff_unbreakable_barrier', 1.0, 0.0, 'normal', '["ultimate", "defensive", "shield"]'),
-- A15. 盾牌反射 (Shield Reflection)
('warrior_shield_reflection', '盾牌反射', '反射伤害，适合高防御战士。1级：反射50%受到的伤害给攻击者，持续2回合，每级+10%反射比例，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 8, 1, 'eff_shield_reflection', 1.0, 0.8, 'normal', '["defensive", "reflect"]'),
-- A16. 挑战怒吼 (Challenging Shout)
('warrior_challenging_shout', '挑战怒吼', '群体嘲讽，保护所有队友。1级：强制所有敌人攻击自己，持续1回合，冷却5回合，每级+0.5回合持续时间，-0.5回合冷却', 'warrior', 'control', 'enemy_all', NULL, 0, NULL, 0, 5, 5, 1, 'eff_challenging_shout', 1.0, 2.5, 'taunt', '["tank", "aoe", "threat", "control"]'),
-- A17. 鲁莽 (Recklessness)
('warrior_recklessness', '鲁莽', '高风险高回报技能，提升输出但降低生存。1级：提升50%暴击率，但受到伤害增加20%，持续3回合，每级+10%暴击率，-2%受到伤害增加，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 8, 1, 'eff_recklessness', 1.0, 1.0, 'normal', '["offensive", "tradeoff", "cooldown"]'),
-- A18. 反击风暴 (Retaliation)
('warrior_retaliation', '反击风暴', '反击技能，受到攻击时自动反击。1级：每次受到攻击时，对攻击者造成50%物理攻击力的反击伤害，持续3回合，每级+10%反击伤害，+0.5回合持续时间，-1回合冷却', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 10, 1, 'eff_retaliation', 1.0, 1.0, 'normal', '["defensive_offensive", "cooldown", "counter"]'),
-- A19. 狂暴之怒 (Berserker Rage)
('warrior_berserker_rage', '狂暴之怒', '强化输出和怒气获取。1级：提升30%攻击力，每次攻击额外获得5点怒气，持续4回合，每级+5%攻击力，+1点额外怒气，+0.5回合持续时间', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 6, 1, 'eff_berserker_rage', 1.0, 1.2, 'normal', '["offensive", "rage_generation"]'),
-- A20. 天神下凡 (Avatar)
('warrior_avatar', '天神下凡', '终极输出技能。1级：提升50%攻击力，免疫控制效果，持续3回合，每级+10%攻击力，+0.5回合持续时间，-1回合冷却', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 0, 12, 1, 'eff_avatar', 1.0, 1.5, 'normal', '["ultimate", "offensive"]'),
-- A21. 致死打击 (Mortal Strike)
('warrior_mortal_strike', '致死打击', '造成高额伤害并降低目标治疗效果。1级：180%物理攻击力，降低目标受到的治疗效果50%持续3回合，每级+25%伤害，+5%降低比例，+0.5回合持续时间', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.8, 30, 0, 1, 'eff_mortal_strike', 1.0, 1.6, 'normal', '["high_damage", "debuff", "anti_heal"]'),
-- A22. 嗜血 (Bloodthirst)
('warrior_bloodthirst', '嗜血', '造成伤害并恢复生命值。1级：120%物理攻击力，恢复造成伤害的30%生命值，每级+15%伤害，+5%恢复比例', 'warrior', 'attack', 'enemy', 'physical', 0, 'strength', 1.2, 25, 0, 1, NULL, 1.0, 1.3, 'normal', '["high_damage", "sustain", "heal"]');

-- ═══════════════════════════════════════════════════════════
-- 战士被动技能 (21个)
-- ═══════════════════════════════════════════════════════════
-- 注意：effect_value为基础值（1级），升级效果在代码中处理

INSERT OR REPLACE INTO passive_skills (id, name, description, class_id, rarity, tier, effect_type, effect_value, effect_stat, max_level, level_scaling) VALUES
-- 属性增强类
-- P1. 坚韧 (Toughness)
('warrior_passive_toughness', '坚韧', '提升最大生命值，增强生存能力。1级：+5%最大HP，每级+2%', 'warrior', 'common', 1, 'stat_mod', 5.0, 'max_hp', 5, 0.4),
-- P2. 护甲专精 (Armor Specialization)
('warrior_passive_armor_spec', '护甲专精', '提升物理防御能力。1级：+5%物理防御，每级+3%', 'warrior', 'common', 1, 'stat_mod', 5.0, 'physical_defense', 5, 0.6),
-- P3. 战斗专注 (Battle Focus)
('warrior_passive_battle_focus', '战斗专注', '提升输出能力。1级：+10%物理攻击力，每级+5%', 'warrior', 'common', 1, 'stat_mod', 10.0, 'attack', 5, 0.5),
-- P4. 武器专精 (Weapon Expertise)
('warrior_passive_weapon_expertise', '武器专精', '提升暴击率。1级：+8%暴击率，每级+4%', 'warrior', 'rare', 2, 'stat_mod', 8.0, 'crit_rate', 5, 0.5),
-- P5. 防御姿态 (Defensive Stance) - 多属性：仇恨+防御
('warrior_passive_defensive_stance', '防御姿态', '被动提升仇恨和防御（类似姿态）。1级：+15%仇恨生成，+10%物理防御，每级+5%仇恨，+3%防御', 'warrior', 'rare', 2, 'stat_mod', 15.0, 'threat_and_defense', 5, 0.33),
-- P6. 武器大师 (Weapon Mastery) - 多属性：攻击+暴击
('warrior_passive_weapon_mastery', '武器大师', '大幅提升输出能力。1级：+10%物理攻击力，+5%暴击率，每级+5%攻击力，+2%暴击率', 'warrior', 'epic', 3, 'stat_mod', 10.0, 'attack_and_crit', 5, 0.5),
-- P7. 战斗大师 (Battle Master) - 多属性：伤害+仇恨
('warrior_passive_battle_master', '战斗大师', '全面提升战斗能力。1级：+15%所有伤害，+10%仇恨生成，每级+5%伤害，+5%仇恨', 'warrior', 'epic', 3, 'stat_mod', 15.0, 'damage_and_threat', 5, 0.33),
-- P8. 守护者 (Guardian) - 多属性：HP+防御
('warrior_passive_guardian', '守护者', '全面提升防御能力。1级：+20%最大HP，+15%物理防御，每级+10%HP，+5%防御', 'warrior', 'epic', 3, 'stat_mod', 20.0, 'hp_and_defense', 5, 0.5),

-- 怒气管理类
-- P9. 愤怒掌握 (Anger Management)
('warrior_passive_anger_management', '愤怒掌握', '提升怒气获取速度。1级：+10%怒气获得量，每级+5%', 'warrior', 'common', 1, 'rage_generation', 10.0, NULL, 5, 0.5),
-- P10. 战争机器 (War Machine)
('warrior_passive_war_machine', '战争机器', '击杀回怒，适合多目标战斗。1级：每次击杀敌人时，立即获得30点怒气，每级+10点', 'warrior', 'epic', 3, 'rage_generation', 30.0, NULL, 5, 0.33),

-- 特殊效果类
-- P11. 血之狂热 (Blood Craze)
('warrior_passive_blood_craze', '血之狂热', '攻击时恢复生命值。1级：每次攻击恢复1%最大HP，每级+0.5%', 'warrior', 'rare', 2, 'on_hit_heal', 1.0, NULL, 5, 0.5),
-- P12. 坚韧不拔 (Unbreakable)
('warrior_passive_unbreakable', '坚韧不拔', '防止被秒杀。1级：受到致命伤害时，保留1点HP（每场战斗1次），每级增加1次触发次数', 'warrior', 'rare', 2, 'survival', 1.0, NULL, 5, 1.0),
-- P13. 盾牌反射 (Shield Reflection Passive)
('warrior_passive_shield_reflection', '盾牌反射', '被动反射伤害。1级：受到物理攻击时，反射10%伤害给攻击者，每级+5%', 'warrior', 'rare', 2, 'reflect', 10.0, NULL, 5, 0.5),
-- P14. 钢铁意志 (Iron Will)
('warrior_passive_iron_will', '钢铁意志', '提升控制抗性。1级：+20%控制效果抗性（眩晕、沉默等），每级+15%', 'warrior', 'epic', 3, 'resistance', 20.0, NULL, 5, 0.75),
-- P15. 复仇 (Revenge)
('warrior_passive_revenge', '复仇', '被动反击技能。1级：受到攻击时，有15%概率立即反击，造成100%物理攻击力伤害，每级+10%触发概率，+20%反击伤害', 'warrior', 'epic', 3, 'counter_attack', 15.0, NULL, 5, 0.67),
-- P16. 不灭意志 (Unbreakable Will)
('warrior_passive_unbreakable_will', '不灭意志', '低血量时提升生存能力。1级：HP低于30%时，受到的伤害减少25%，每级+10%减伤，-5%触发阈值', 'warrior', 'epic', 3, 'survival', 25.0, NULL, 5, 0.4),
-- P17. 狂暴之心 (Berserker Heart)
('warrior_passive_berserker_heart', '狂暴之心', '低血量时提升输出。1级：HP低于50%时，攻击力+20%，每级+10%攻击力，-5%触发阈值', 'warrior', 'epic', 3, 'stat_mod', 20.0, 'attack', 5, 0.5),
-- P18. 战争领主 (War Lord) - 多属性：伤害+暴击
('warrior_passive_war_lord', '战争领主', '终极输出被动。1级：+20%所有伤害，+15%暴击率，每级+10%伤害，+5%暴击率', 'warrior', 'epic', 3, 'stat_mod', 20.0, 'damage_and_crit', 5, 0.5),
-- P19. 钢铁堡垒 (Iron Fortress) - 多属性：HP+防御+抗性
('warrior_passive_iron_fortress', '钢铁堡垒', '终极防御被动。1级：+25%最大HP，+20%物理防御，+15%控制抗性，每级+5%HP，+5%防御，+5%抗性', 'warrior', 'epic', 3, 'stat_mod', 25.0, 'hp_defense_resistance', 5, 0.2),
-- P20. 战神 (War God) - 多属性：伤害+暴击+仇恨
('warrior_passive_war_god', '战神', '传说级输出被动。1级：+25%所有伤害，+20%暴击率，+15%仇恨生成，每级+5%伤害，+5%暴击率，+5%仇恨', 'warrior', 'legendary', 4, 'stat_mod', 25.0, 'damage_crit_threat', 5, 0.2),
-- P21. 不灭守护 (Immortal Guardian) - 多属性：HP+防御+抗性+免疫
('warrior_passive_immortal_guardian', '不灭守护', '传说级防御被动。1级：+30%最大HP，+25%物理防御，+20%控制抗性，每场战斗免疫一次致命伤害，每级+5%HP，+5%防御，+5%抗性，+1次免疫次数', 'warrior', 'legendary', 4, 'stat_mod', 30.0, 'hp_defense_resistance_immune', 5, 0.17);

