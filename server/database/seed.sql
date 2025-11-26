-- ═══════════════════════════════════════════════════════════
-- Text WoW 初始游戏数据
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 槽位解锁条件配置
-- ═══════════════════════════════════════════════════════════

-- 注：槽位解锁逻辑在代码中实现，这里记录配置说明
-- 槽位1: 初始拥有
-- 槽位2: 队伍中任意角色达到 10 级
-- 槽位3: 队伍中任意角色达到 20 级  
-- 槽位4: 队伍中任意角色达到 35 级
-- 槽位5: 队伍中任意角色达到 50 级

-- ═══════════════════════════════════════════════════════════
-- 种族被动效果 (全部为被动，适合放置游戏自动战斗)
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO effects (id, name, description, type, is_buff, is_stackable, max_stacks, duration, value_type, value, stat_affected, can_dispel) VALUES
-- 人类被动
('racial_human_exp', '适应力', '经验获取+10%', 'stat_mod', 1, 0, 1, 999, 'percent', 10, 'exp_gain', 0),
('racial_human_sword', '剑术专精', '物理伤害+3%', 'stat_mod', 1, 0, 1, 999, 'percent', 3, 'physical_damage', 0),
-- 矮人被动
('racial_dwarf_frost', '霜抗', '冰霜伤害-15%', 'stat_mod', 1, 0, 1, 999, 'percent', -15, 'frost_damage_taken', 0),
('racial_dwarf_stone', '石肤', '受到暴击伤害-10%', 'stat_mod', 1, 0, 1, 999, 'percent', -10, 'crit_damage_taken', 0),
-- 暗夜精灵被动
('racial_nelf_shadow', '暗影之心', '夜间伤害+8%', 'stat_mod', 1, 0, 1, 999, 'percent', 8, 'damage_dealt', 0),
('racial_nelf_dodge', '敏锐', '闪避率+2%', 'stat_mod', 1, 0, 1, 999, 'percent', 2, 'dodge_rate', 0),
-- 侏儒被动
('racial_gnome_crit', '灵巧心智', '法术暴击+3%', 'stat_mod', 1, 0, 1, 999, 'percent', 3, 'spell_crit', 0),
('racial_gnome_mech', '工程专精', '对机械怪伤害+15%', 'stat_mod', 1, 0, 1, 999, 'percent', 15, 'damage_vs_mechanical', 0),
-- 兽人被动
('racial_orc_fury', '嗜血', 'HP<30%时攻击+15%', 'stat_mod', 1, 0, 1, 999, 'percent', 15, 'attack_low_hp', 0),
('racial_orc_stun', '坚韧', '眩晕时间-25%', 'stat_mod', 1, 0, 1, 999, 'percent', -25, 'stun_duration', 0),
-- 亡灵被动
('racial_undead_touch', '亡者之触', '攻击5%几率恐惧', 'proc', 1, 0, 1, 999, 'percent', 5, 'fear_on_hit', 0),
('racial_undead_shadow', '暗影抗性', '暗影伤害-15%', 'stat_mod', 1, 0, 1, 999, 'percent', -15, 'shadow_damage_taken', 0),
-- 牛头人被动
('racial_tauren_hp', '坚忍', '最大HP+5%', 'stat_mod', 1, 0, 1, 999, 'percent', 5, 'max_hp', 0),
('racial_tauren_heal', '自然亲和', '受到治疗+10%', 'stat_mod', 1, 0, 1, 999, 'percent', 10, 'healing_taken', 0),
-- 巨魔被动
('racial_troll_regen', '再生', '每回合恢复2%HP', 'hot', 1, 0, 1, 999, 'percent', 2, NULL, 0),
('racial_troll_beast', '野兽杀手', '对野兽伤害+15%', 'stat_mod', 1, 0, 1, 999, 'percent', 15, 'damage_vs_beast', 0);

-- ═══════════════════════════════════════════════════════════
-- 种族数据 (纯被动设计，适合自动战斗)
-- ═══════════════════════════════════════════════════════════

-- 种族数据说明:
-- strength_base等: 基础固定加成
-- strength_pct等: 百分比加成(5 = 5%)
-- 公式: 最终属性 = (基础 + 成长 + 装备 + 种族基础) × (1 + 种族百分比/100)

INSERT OR REPLACE INTO races (id, name, faction, description, 
    strength_base, agility_base, intellect_base, stamina_base, spirit_base,
    strength_pct, agility_pct, intellect_pct, stamina_pct, spirit_pct,
    racial_passive_id, racial_passive2_id) VALUES
-- 联盟
('human', '人类', 'alliance', '适应力强的种族，学习能力出众，剑术精湛。',
    0, 0, 0, 0, 3,    -- 基础: 精神+3
    0, 0, 0, 0, 3,    -- 百分比: 精神+3%
    'racial_human_exp', 'racial_human_sword'),
('dwarf', '矮人', 'alliance', '坚韧的山地种族，皮糙肉厚，抵抗寒冷。',
    3, 0, 0, 0, 0,    -- 基础: 力量+3
    0, 0, 0, 5, 0,    -- 百分比: 耐力+5%
    'racial_dwarf_frost', 'racial_dwarf_stone'),
('nightelf', '暗夜精灵', 'alliance', '古老的精灵种族，夜间行动敏捷。',
    0, 5, 0, 0, 0,    -- 基础: 敏捷+5
    0, 3, 0, 0, 0,    -- 百分比: 敏捷+3%
    'racial_nelf_shadow', 'racial_nelf_dodge'),
('gnome', '侏儒', 'alliance', '聪明的小型种族，精通魔法和机械。',
    0, 0, 5, 0, 0,    -- 基础: 智力+5
    0, 0, 5, 0, 0,    -- 百分比: 智力+5%
    'racial_gnome_crit', 'racial_gnome_mech'),
-- 部落
('orc', '兽人', 'horde', '强壮的战士种族，濒死时爆发战斗本能。',
    5, 0, 0, 0, 0,    -- 基础: 力量+5
    5, 0, 0, 0, 0,    -- 百分比: 力量+5%
    'racial_orc_fury', 'racial_orc_stun'),
('undead', '亡灵', 'horde', '不死的存在，攻击带有恐惧之力。',
    0, 0, 3, 0, 0,    -- 基础: 智力+3
    0, 0, 0, 0, 5,    -- 百分比: 精神+5% (暗影亲和)
    'racial_undead_touch', 'racial_undead_shadow'),
('tauren', '牛头人', 'horde', '高大温和的种族，生命力顽强。',
    0, 0, 0, 5, 0,    -- 基础: 耐力+5
    0, 0, 0, 5, 0,    -- 百分比: 耐力+5%
    'racial_tauren_hp', 'racial_tauren_heal'),
('troll', '巨魔', 'horde', '敏捷的丛林种族，拥有惊人的再生能力。',
    0, 3, 0, 0, 0,    -- 基础: 敏捷+3
    0, 5, 0, 0, 0,    -- 百分比: 敏捷+5%
    'racial_troll_regen', 'racial_troll_beast');

-- ═══════════════════════════════════════════════════════════
-- 职业数据
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO classes (id, name, description, role, primary_stat, base_hp, base_mp, hp_per_level, mp_per_level, base_strength, base_agility, base_intellect, base_stamina, base_spirit) VALUES
('warrior', '战士', '近战格斗专家，可以承受大量伤害。', 'tank', 'strength', 120, 20, 12, 2, 15, 10, 5, 14, 8),
('paladin', '圣骑士', '神圣战士，可以治疗和保护盟友。', 'tank', 'strength', 110, 60, 10, 6, 13, 8, 10, 13, 12),
('hunter', '猎人', '远程物理攻击者，与宠物并肩作战。', 'dps', 'agility', 90, 40, 8, 4, 8, 15, 8, 10, 10),
('rogue', '盗贼', '潜行刺客，擅长连击和爆发伤害。', 'dps', 'agility', 85, 50, 7, 5, 10, 16, 6, 9, 8),
('priest', '牧师', '治疗者和暗影施法者。', 'healer', 'intellect', 70, 100, 5, 12, 5, 6, 15, 8, 16),
('mage', '法师', '强大的奥术施法者，擅长范围伤害。', 'dps', 'intellect', 65, 120, 4, 15, 4, 6, 18, 6, 12),
('warlock', '术士', '黑暗魔法师，召唤恶魔作战。', 'dps', 'intellect', 75, 110, 5, 13, 5, 6, 17, 8, 10),
('druid', '德鲁伊', '自然的守护者，可变形为多种形态。', 'dps', 'intellect', 85, 80, 7, 10, 10, 10, 13, 10, 12),
('shaman', '萨满', '元素的操控者，可治疗和增益。', 'dps', 'intellect', 90, 90, 8, 10, 12, 8, 14, 11, 12);

-- ═══════════════════════════════════════════════════════════
-- 效果配置 (Buff/Debuff) - 每场战斗开始时清空
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO effects (id, name, description, type, is_buff, is_stackable, max_stacks, duration, value_type, value, stat_affected, damage_type, can_dispel) VALUES
-- 增益效果 (Buff)
('eff_shield_wall', '盾墙', '受到的伤害降低50%', 'stat_mod', 1, 0, 1, 3, 'percent', -50, 'damage_taken', NULL, 0),
('eff_ice_barrier', '寒冰护体', '吸收伤害的护盾', 'shield', 1, 0, 1, 5, 'flat', 100, NULL, NULL, 1),
('eff_pw_shield', '真言术:盾', '吸收伤害', 'shield', 1, 0, 1, 4, 'flat', 80, NULL, NULL, 1),
('eff_battle_shout', '战斗怒吼', '攻击力提升10%', 'stat_mod', 1, 0, 1, 5, 'percent', 10, 'attack', NULL, 1),
('eff_blade_flurry', '剑刃乱舞', '攻击力提升20%', 'stat_mod', 1, 0, 1, 4, 'percent', 20, 'attack', NULL, 1),
('eff_renew', '恢复', '每回合恢复生命', 'hot', 1, 0, 1, 4, 'flat', 15, NULL, 'holy', 1),
('eff_stealth', '潜行', '进入隐身状态', 'stealth', 1, 0, 1, 3, NULL, NULL, NULL, NULL, 1),
('eff_inner_fire', '心灵之火', '防御力提升15%', 'stat_mod', 1, 0, 1, 10, 'percent', 15, 'defense', NULL, 1),
('eff_blessing_might', '力量祝福', '攻击力提升', 'stat_mod', 1, 0, 1, 8, 'flat', 20, 'attack', NULL, 1),
('eff_arcane_intellect', '奥术智慧', '智力提升10%', 'stat_mod', 1, 0, 1, 30, 'percent', 10, 'intellect', NULL, 1),

-- 减益效果 (Debuff)
('eff_stun', '眩晕', '无法行动', 'stun', 0, 0, 1, 1, NULL, NULL, NULL, NULL, 1),
('eff_silence', '沉默', '无法施放法术', 'silence', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 1),
('eff_slow', '减速', '攻击速度降低30%', 'slow', 0, 0, 1, 3, 'percent', -30, 'attack_speed', NULL, 1),
('eff_sw_pain', '暗言术:痛', '每回合造成暗影伤害', 'dot', 0, 0, 1, 4, 'flat', 12, NULL, 'shadow', 1),
('eff_rend', '撕裂', '每回合造成物理流血伤害', 'dot', 0, 1, 3, 3, 'flat', 8, NULL, 'physical', 1),
('eff_ignite', '点燃', '每回合造成火焰伤害', 'dot', 0, 1, 5, 3, 'flat', 10, NULL, 'fire', 1),
('eff_frostbite', '冻伤', '攻击力降低15%', 'stat_mod', 0, 0, 1, 3, 'percent', -15, 'attack', 'frost', 1),
('eff_curse_weakness', '虚弱诅咒', '造成的伤害降低20%', 'stat_mod', 0, 0, 1, 4, 'percent', -20, 'damage_dealt', 'shadow', 1),
('eff_sunder_armor', '破甲', '防御力降低20%', 'stat_mod', 0, 1, 5, 5, 'percent', -20, 'defense', NULL, 0),
('eff_poison', '中毒', '每回合造成自然伤害', 'dot', 0, 1, 5, 5, 'flat', 6, NULL, 'nature', 1),
('eff_taunt', '嘲讽', '强制攻击施法者', 'taunt', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
('eff_fear', '恐惧', '无法控制行动', 'stun', 0, 0, 1, 2, NULL, NULL, 'shadow', NULL, 1);

-- ═══════════════════════════════════════════════════════════
-- 技能数据 (扩展版)
-- ═══════════════════════════════════════════════════════════

-- 战士技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, mp_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('heroic_strike', '英勇打击', '一次强力的武器攻击。', 'warrior', 'attack', 'enemy', 'physical', 15, 'strength', 1.2, 5, 0, 1, NULL, 1.0),
('charge', '冲锋', '冲向敌人，造成伤害并眩晕。', 'warrior', 'attack', 'enemy', 'physical', 10, 'strength', 0.8, 10, 3, 1, 'eff_stun', 1.0),
('rend', '撕裂', '造成流血效果。', 'warrior', 'dot', 'enemy', 'physical', 5, 'strength', 0.3, 8, 0, 2, 'eff_rend', 1.0),
('thunder_clap', '雷霆一击', '对所有敌人造成伤害并减速。', 'warrior', 'attack', 'enemy_all', 'physical', 20, 'strength', 0.6, 15, 2, 4, 'eff_slow', 0.8),
('sunder_armor', '破甲攻击', '降低敌人防御力。', 'warrior', 'debuff', 'enemy', 'physical', 8, 'strength', 0.5, 10, 0, 6, 'eff_sunder_armor', 1.0),
('execute', '斩杀', '对低血量敌人造成大量伤害。', 'warrior', 'attack', 'enemy_lowest_hp', 'physical', 60, 'strength', 1.8, 20, 0, 8, NULL, 1.0),
('shield_wall', '盾墙', '大幅减少受到的伤害。', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 30, 10, 10, 'eff_shield_wall', 1.0),
('battle_shout', '战斗怒吼', '提升全队攻击力。', 'warrior', 'buff', 'ally_all', NULL, 0, NULL, 0, 15, 5, 3, 'eff_battle_shout', 1.0);

-- 法师技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, mp_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('fireball', '火球术', '发射火球，有几率点燃敌人。', 'mage', 'attack', 'enemy', 'fire', 25, 'intellect', 1.3, 15, 0, 1, 'eff_ignite', 0.3),
('frostbolt', '寒冰箭', '发射寒冰箭，减缓敌人。', 'mage', 'attack', 'enemy', 'frost', 18, 'intellect', 1.1, 12, 0, 1, 'eff_frostbite', 0.5),
('arcane_missiles', '奥术飞弹', '发射多道奥术飞弹。', 'mage', 'attack', 'enemy', 'magic', 30, 'intellect', 1.4, 20, 2, 4, NULL, 1.0),
('flamestrike', '烈焰风暴', '对所有敌人造成火焰伤害。', 'mage', 'attack', 'enemy_all', 'fire', 35, 'intellect', 1.0, 30, 3, 6, 'eff_ignite', 0.2),
('pyroblast', '炎爆术', '施放巨大的火球。', 'mage', 'attack', 'enemy', 'fire', 80, 'intellect', 2.0, 45, 5, 8, 'eff_ignite', 0.8),
('ice_barrier', '寒冰护体', '创造吸收伤害的护盾。', 'mage', 'shield', 'self', NULL, 0, 'intellect', 1.0, 35, 8, 10, 'eff_ice_barrier', 1.0),
('arcane_intellect', '奥术智慧', '提升智力。', 'mage', 'buff', 'ally', NULL, 0, NULL, 0, 20, 0, 2, 'eff_arcane_intellect', 1.0);

-- 盗贼技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, mp_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('sinister_strike', '邪恶攻击', '快速的攻击。', 'rogue', 'attack', 'enemy', 'physical', 15, 'agility', 1.1, 8, 0, 1, NULL, 1.0),
('backstab', '背刺', '从背后攻击造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 35, 'agility', 1.6, 18, 0, 1, NULL, 1.0),
('deadly_poison', '致命毒药', '使敌人中毒。', 'rogue', 'dot', 'enemy', 'nature', 0, 'agility', 0.2, 12, 0, 3, 'eff_poison', 1.0),
('eviscerate', '剔骨', '造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 50, 'agility', 1.8, 25, 0, 4, NULL, 1.0),
('kidney_shot', '肾击', '眩晕敌人。', 'rogue', 'control', 'enemy', 'physical', 10, 'agility', 0.5, 20, 4, 6, 'eff_stun', 1.0),
('blade_flurry', '剑刃乱舞', '攻击力大幅提升。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 30, 8, 8, 'eff_blade_flurry', 1.0),
('vanish', '消失', '进入潜行状态。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 40, 10, 10, 'eff_stealth', 1.0);

-- 牧师技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, mp_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('smite', '惩击', '用神圣能量攻击敌人。', 'priest', 'attack', 'enemy', 'holy', 18, 'intellect', 1.0, 10, 0, 1, NULL, 1.0),
('shadow_word_pain', '暗言术:痛', '对敌人施加持续伤害。', 'priest', 'dot', 'enemy', 'shadow', 0, 'intellect', 0.3, 12, 0, 1, 'eff_sw_pain', 1.0),
('lesser_heal', '次级治疗术', '恢复生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 25, 'spirit', 1.0, 15, 0, 1, NULL, 1.0),
('renew', '恢复', '持续恢复生命。', 'priest', 'hot', 'ally_lowest_hp', 'holy', 0, 'spirit', 0.5, 18, 0, 2, 'eff_renew', 1.0),
('heal', '治疗术', '恢复大量生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 50, 'spirit', 1.5, 25, 0, 4, NULL, 1.0),
('inner_fire', '心灵之火', '提升防御力。', 'priest', 'buff', 'self', NULL, 0, NULL, 0, 20, 0, 4, 'eff_inner_fire', 1.0),
('power_word_shield', '真言术:盾', '创造吸收伤害的护盾。', 'priest', 'shield', 'ally_lowest_hp', 'holy', 0, 'spirit', 1.2, 25, 4, 6, 'eff_pw_shield', 1.0),
('flash_heal', '快速治疗', '快速恢复生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 40, 'spirit', 1.2, 20, 0, 8, NULL, 1.0),
('silence', '沉默', '使敌人无法施法。', 'priest', 'control', 'enemy', 'shadow', 0, NULL, 0, 25, 6, 10, 'eff_silence', 1.0);

-- 通用技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, mp_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('basic_attack', '普通攻击', '基础的物理攻击。', NULL, 'attack', 'enemy', 'physical', 0, 'strength', 1.0, 0, 0, 1, NULL, 1.0);

-- ═══════════════════════════════════════════════════════════
-- 区域数据
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier) VALUES
('elwynn', '艾尔文森林', '人类王国暴风城外的宁静森林，适合新手冒险者。', 1, 10, 'alliance', 1.0, 1.0),
('durotar', '杜隆塔尔', '兽人的家园，炎热干燥的红色大地。', 1, 10, 'horde', 1.0, 1.0),
('westfall', '西部荒野', '曾经肥沃的农田，如今被迪菲亚兄弟会占领。', 10, 20, 'alliance', 1.1, 1.1),
('barrens', '贫瘠之地', '广袤的草原，危险与机遇并存。', 10, 25, 'horde', 1.1, 1.1),
('duskwood', '暮色森林', '被永恒黑暗笼罩的诡异森林，亡灵与狼人出没。', 20, 30, NULL, 1.2, 1.2),
('stranglethorn', '荆棘谷', '危险的丛林，到处是食人族和野兽。', 30, 45, NULL, 1.3, 1.3),
('tanaris', '塔纳利斯', '炎热的沙漠，隐藏着古老的秘密。', 40, 50, NULL, 1.4, 1.4),
('burning_steppes', '燃烧平原', '被黑龙军团占领的焦土。', 50, 60, NULL, 1.5, 1.5);

-- ═══════════════════════════════════════════════════════════
-- 怪物数据
-- ═══════════════════════════════════════════════════════════

-- 艾尔文森林 (1-10级)
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('wolf', 'elwynn', '森林狼', 1, 'normal', 30, 5, 1, 15, 1, 3, 100),
('young_boar', 'elwynn', '小野猪', 1, 'normal', 25, 4, 2, 12, 1, 2, 100),
('kobold_worker', 'elwynn', '狗头人矿工', 2, 'normal', 40, 6, 2, 20, 2, 5, 80),
('kobold_tunneler', 'elwynn', '狗头人掘地工', 3, 'normal', 50, 7, 3, 25, 2, 6, 60),
('defias_thug', 'elwynn', '迪菲亚暴徒', 4, 'normal', 60, 9, 3, 35, 3, 8, 50),
('defias_bandit', 'elwynn', '迪菲亚劫匪', 5, 'normal', 75, 11, 4, 45, 4, 10, 40),
('murloc', 'elwynn', '鱼人', 3, 'normal', 45, 8, 2, 28, 2, 7, 70),
('murloc_warrior', 'elwynn', '鱼人战士', 5, 'normal', 70, 12, 5, 50, 5, 12, 30),
('prowler', 'elwynn', '潜伏者', 6, 'normal', 85, 14, 5, 60, 5, 15, 25),
('hogger', 'elwynn', '霍格', 8, 'elite', 300, 25, 10, 200, 20, 50, 5);

-- 西部荒野 (10-20级)
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('harvest_golem', 'westfall', '收割傀儡', 10, 'normal', 150, 20, 8, 80, 8, 20, 100),
('defias_rogue', 'westfall', '迪菲亚盗贼', 11, 'normal', 170, 22, 9, 95, 10, 25, 80),
('defias_highwayman', 'westfall', '迪菲亚拦路贼', 12, 'normal', 190, 25, 10, 110, 12, 30, 60),
('gnoll_brute', 'westfall', '豺狼人蛮兵', 13, 'normal', 210, 28, 11, 125, 14, 35, 50),
('gnoll_mystic', 'westfall', '豺狼人秘法师', 14, 'normal', 180, 32, 9, 140, 16, 40, 40),
('defias_overlord', 'westfall', '迪菲亚霸主', 16, 'elite', 500, 45, 18, 350, 40, 100, 5),
('dust_devil', 'westfall', '尘土恶魔', 15, 'normal', 250, 35, 12, 160, 18, 45, 30);

-- 暮色森林 (20-30级)
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('skeleton_warrior', 'duskwood', '骷髅战士', 20, 'normal', 300, 40, 15, 200, 20, 50, 100),
('skeleton_mage', 'duskwood', '骷髅法师', 21, 'normal', 250, 50, 12, 220, 22, 55, 80),
('ghoul', 'duskwood', '食尸鬼', 22, 'normal', 350, 45, 16, 240, 25, 60, 70),
('dire_wolf', 'duskwood', '恐狼', 23, 'normal', 380, 48, 18, 260, 28, 65, 60),
('worgen', 'duskwood', '狼人', 24, 'normal', 420, 55, 20, 300, 30, 75, 50),
('worgen_alpha', 'duskwood', '狼人首领', 26, 'elite', 800, 80, 30, 600, 60, 150, 5),
('abomination', 'duskwood', '憎恶', 28, 'elite', 1000, 90, 35, 800, 80, 200, 3),
('stitches', 'duskwood', '缝合怪', 30, 'boss', 2000, 120, 45, 1500, 150, 400, 1);

-- ═══════════════════════════════════════════════════════════
-- 物品数据
-- ═══════════════════════════════════════════════════════════

-- 消耗品
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, stackable, max_stack, sell_price, buy_price, effect_type, effect_value) VALUES
('minor_healing_potion', '初级治疗药水', '恢复50点生命值。', 'consumable', 'potion', 'common', 1, 20, 1, 5, 'heal_hp', 50),
('healing_potion', '治疗药水', '恢复150点生命值。', 'consumable', 'potion', 'uncommon', 1, 20, 5, 25, 'heal_hp', 150),
('greater_healing_potion', '强效治疗药水', '恢复400点生命值。', 'consumable', 'potion', 'rare', 1, 20, 20, 100, 'heal_hp', 400),
('minor_mana_potion', '初级法力药水', '恢复50点法力值。', 'consumable', 'potion', 'common', 1, 20, 1, 5, 'heal_mp', 50),
('mana_potion', '法力药水', '恢复150点法力值。', 'consumable', 'potion', 'uncommon', 1, 20, 5, 25, 'heal_mp', 150);

-- 装备 - 武器
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, level_required, slot, sell_price, attack, strength) VALUES
('worn_sword', '破旧的剑', '一把破旧的铁剑。', 'equipment', 'weapon', 'common', 1, 'main_hand', 5, 3, 1),
('militia_sword', '民兵之剑', '联盟民兵的标准武器。', 'equipment', 'weapon', 'common', 5, 'main_hand', 20, 8, 2),
('outlaw_sabre', '逃犯军刀', '从迪菲亚成员身上缴获。', 'equipment', 'weapon', 'uncommon', 10, 'main_hand', 80, 15, 4),
('corpsemaker', '尸体收割者', '沉重的双手斧。', 'equipment', 'weapon', 'rare', 20, 'main_hand', 300, 30, 8);

-- 装备 - 护甲
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, level_required, slot, sell_price, defense, stamina) VALUES
('worn_leather_vest', '破旧皮甲', '一件破旧的皮甲。', 'equipment', 'armor', 'common', 1, 'chest', 3, 2, 1),
('militia_chain_vest', '民兵锁甲', '联盟民兵的标准护甲。', 'equipment', 'armor', 'common', 5, 'chest', 15, 5, 2),
('defias_leather_vest', '迪菲亚皮甲', '迪菲亚兄弟会的制服。', 'equipment', 'armor', 'uncommon', 10, 'chest', 60, 10, 4),
('blackened_defias_armor', '黑化迪菲亚护甲', '高级成员的护甲。', 'equipment', 'armor', 'rare', 15, 'chest', 150, 18, 7);

-- 材料
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, stackable, max_stack, sell_price) VALUES
('wolf_pelt', '狼皮', '可以出售给商人。', 'material', 'leather', 'common', 1, 99, 2),
('linen_cloth', '亚麻布', '基础的布料。', 'material', 'cloth', 'common', 1, 99, 1),
('copper_ore', '铜矿石', '基础的矿石。', 'material', 'ore', 'common', 1, 99, 2),
('kobold_candle', '狗头人蜡烛', '你不许拿走蜡烛！', 'material', 'junk', 'common', 1, 99, 1);

-- ═══════════════════════════════════════════════════════════
-- 怪物掉落表
-- ═══════════════════════════════════════════════════════════

-- 森林狼掉落
INSERT OR REPLACE INTO monster_drops (monster_id, item_id, drop_rate, min_quantity, max_quantity) VALUES
('wolf', 'wolf_pelt', 0.5, 1, 1),
('wolf', 'minor_healing_potion', 0.1, 1, 1);

-- 狗头人掉落
INSERT OR REPLACE INTO monster_drops (monster_id, item_id, drop_rate, min_quantity, max_quantity) VALUES
('kobold_worker', 'kobold_candle', 0.3, 1, 2),
('kobold_worker', 'copper_ore', 0.2, 1, 1),
('kobold_worker', 'linen_cloth', 0.15, 1, 2);

-- 迪菲亚掉落
INSERT OR REPLACE INTO monster_drops (monster_id, item_id, drop_rate, min_quantity, max_quantity) VALUES
('defias_thug', 'linen_cloth', 0.3, 1, 2),
('defias_thug', 'worn_sword', 0.05, 1, 1),
('defias_bandit', 'linen_cloth', 0.35, 1, 3),
('defias_bandit', 'militia_sword', 0.03, 1, 1),
('defias_bandit', 'worn_leather_vest', 0.03, 1, 1);

-- 霍格掉落
INSERT OR REPLACE INTO monster_drops (monster_id, item_id, drop_rate, min_quantity, max_quantity) VALUES
('hogger', 'outlaw_sabre', 0.15, 1, 1),
('hogger', 'defias_leather_vest', 0.15, 1, 1),
('hogger', 'healing_potion', 0.5, 1, 2);


