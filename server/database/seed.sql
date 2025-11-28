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
-- 职业数据 (含能量系统)
-- ═══════════════════════════════════════════════════════════
-- resource_type: rage(怒气) / energy(能量) / mana(法力)
-- resource_regen: 每回合固定恢复值
-- resource_regen_pct: 每回合基于精神的百分比恢复(仅法力有效)

-- ═══════════════════════════════════════════════════════════
-- 职业数据 (小数值设计：HP 25-35起步，每级+3)
-- ═══════════════════════════════════════════════════════════
-- resource_type: rage(怒气) / energy(能量) / mana(法力)
-- 数值规范: HP 25~230, MP 15~100, 属性 5~50

INSERT OR REPLACE INTO classes (id, name, description, role, primary_stat, 
    resource_type, base_hp, base_resource, hp_per_level, resource_per_level, 
    resource_regen, resource_regen_pct,
    base_strength, base_agility, base_intellect, base_stamina, base_spirit) VALUES
-- 战士: 怒气系统 (初始0，通过攻击/受击获得)
('warrior', '战士', '近战格斗专家，可以承受大量伤害。', 'tank', 'strength',
    'rage', 35, 0, 3, 0, 0, 0,  -- HP35起步，每级+3
    12, 8, 5, 10, 6),
-- 盗贼: 能量系统 (固定100上限，快速恢复)
('rogue', '盗贼', '潜行刺客，擅长连击和爆发伤害。', 'dps', 'agility',
    'energy', 25, 100, 2, 0, 20, 0,  -- HP25起步，每级+2
    8, 12, 5, 7, 6),
-- 法师: 法力系统 (高法力，基于精神恢复)
('mage', '法师', '强大的奥术施法者，擅长范围伤害。', 'dps', 'intellect',
    'mana', 20, 40, 2, 2, 0, 0.5,  -- HP20起步，MP40起步
    4, 5, 14, 5, 10),
-- 牧师: 法力系统 (高法力，高精神恢复)
('priest', '牧师', '治疗者和暗影施法者。', 'healer', 'intellect',
    'mana', 22, 35, 2, 2, 0, 0.8,  -- HP22起步
    4, 5, 12, 6, 14),
-- 术士: 法力系统
('warlock', '术士', '黑暗魔法师，召唤恶魔作战。', 'dps', 'intellect',
    'mana', 24, 38, 2, 2, 0, 0.5,
    5, 5, 13, 6, 9),
-- 德鲁伊: 法力系统
('druid', '德鲁伊', '自然的守护者，可变形为多种形态。', 'dps', 'intellect',
    'mana', 28, 30, 2, 1, 0, 0.6,
    8, 8, 10, 8, 10),
-- 萨满: 法力系统
('shaman', '萨满', '元素的操控者，可治疗和增益。', 'dps', 'intellect',
    'mana', 28, 32, 2, 1, 0, 0.6,
    9, 7, 11, 8, 10),
-- 圣骑士: 法力系统 (较低法力)
('paladin', '圣骑士', '神圣战士，可以治疗和保护盟友。', 'tank', 'strength',
    'mana', 32, 20, 3, 1, 0, 0.4,
    10, 6, 8, 10, 10),
-- 猎人: 法力系统 (较低法力)
('hunter', '猎人', '远程物理攻击者，与宠物并肩作战。', 'dps', 'agility',
    'mana', 26, 18, 2, 1, 0, 0.3,
    6, 12, 6, 8, 8);

-- ═══════════════════════════════════════════════════════════
-- 效果配置 (Buff/Debuff) - 每场战斗开始时清空
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 效果配置 (小数值设计：DOT伤害2~5/回合，护盾10~20吸收)
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO effects (id, name, description, type, is_buff, is_stackable, max_stacks, duration, value_type, value, stat_affected, damage_type, can_dispel) VALUES
-- 增益效果 (Buff)
('eff_shield_wall', '盾墙', '受到的伤害降低50%', 'stat_mod', 1, 0, 1, 3, 'percent', -50, 'damage_taken', NULL, 0),
('eff_ice_barrier', '寒冰护体', '吸收15点伤害', 'shield', 1, 0, 1, 5, 'flat', 15, NULL, NULL, 1),
('eff_pw_shield', '真言术:盾', '吸收12点伤害', 'shield', 1, 0, 1, 4, 'flat', 12, NULL, NULL, 1),
('eff_battle_shout', '战斗怒吼', '攻击力提升10%', 'stat_mod', 1, 0, 1, 5, 'percent', 10, 'attack', NULL, 1),
('eff_blade_flurry', '剑刃乱舞', '攻击力提升20%', 'stat_mod', 1, 0, 1, 4, 'percent', 20, 'attack', NULL, 1),
('eff_renew', '恢复', '每回合恢复3点生命', 'hot', 1, 0, 1, 4, 'flat', 3, NULL, 'holy', 1),
('eff_stealth', '潜行', '进入隐身状态', 'stealth', 1, 0, 1, 3, NULL, NULL, NULL, NULL, 1),
('eff_inner_fire', '心灵之火', '防御力提升15%', 'stat_mod', 1, 0, 1, 10, 'percent', 15, 'defense', NULL, 1),
('eff_blessing_might', '力量祝福', '攻击力+3', 'stat_mod', 1, 0, 1, 8, 'flat', 3, 'attack', NULL, 1),
('eff_arcane_intellect', '奥术智慧', '智力提升10%', 'stat_mod', 1, 0, 1, 30, 'percent', 10, 'intellect', NULL, 1),

-- 减益效果 (Debuff)
('eff_stun', '眩晕', '无法行动', 'stun', 0, 0, 1, 1, NULL, NULL, NULL, NULL, 1),
('eff_silence', '沉默', '无法施放法术', 'silence', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 1),
('eff_slow', '减速', '攻击速度降低30%', 'slow', 0, 0, 1, 3, 'percent', -30, 'attack_speed', NULL, 1),
('eff_sw_pain', '暗言术:痛', '每回合3点暗影伤害', 'dot', 0, 0, 1, 4, 'flat', 3, NULL, 'shadow', 1),
('eff_rend', '撕裂', '每回合2点流血伤害', 'dot', 0, 1, 3, 3, 'flat', 2, NULL, 'physical', 1),
('eff_ignite', '点燃', '每回合2点火焰伤害', 'dot', 0, 1, 5, 3, 'flat', 2, NULL, 'fire', 1),
('eff_frostbite', '冻伤', '攻击力降低15%', 'stat_mod', 0, 0, 1, 3, 'percent', -15, 'attack', 'frost', 1),
('eff_curse_weakness', '虚弱诅咒', '造成的伤害降低20%', 'stat_mod', 0, 0, 1, 4, 'percent', -20, 'damage_dealt', 'shadow', 1),
('eff_sunder_armor', '破甲', '防御力降低20%', 'stat_mod', 0, 1, 5, 5, 'percent', -20, 'defense', NULL, 0),
('eff_poison', '中毒', '每回合2点自然伤害', 'dot', 0, 1, 5, 5, 'flat', 2, NULL, 'nature', 1),
('eff_taunt', '嘲讽', '强制攻击施法者', 'taunt', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
('eff_fear', '恐惧', '无法控制行动', 'stun', 0, 0, 1, 2, NULL, NULL, 'shadow', NULL, 1);

-- ═══════════════════════════════════════════════════════════
-- 技能数据 (扩展版)
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 技能数据 (小数值设计：伤害5~40, 消耗3~20)
-- ═══════════════════════════════════════════════════════════

-- 战士技能 (怒气消耗5~25)
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('heroic_strike', '英勇打击', '一次强力的武器攻击。', 'warrior', 'attack', 'enemy', 'physical', 8, 'strength', 0.5, 5, 0, 1, NULL, 1.0),
('charge', '冲锋', '冲向敌人，造成伤害并眩晕。', 'warrior', 'attack', 'enemy', 'physical', 5, 'strength', 0.3, 8, 3, 1, 'eff_stun', 1.0),
('rend', '撕裂', '造成流血效果。', 'warrior', 'dot', 'enemy', 'physical', 2, 'strength', 0.15, 6, 0, 2, 'eff_rend', 1.0),
('thunder_clap', '雷霆一击', '对所有敌人造成伤害并减速。', 'warrior', 'attack', 'enemy_all', 'physical', 6, 'strength', 0.3, 12, 2, 4, 'eff_slow', 0.8),
('sunder_armor', '破甲攻击', '降低敌人防御力。', 'warrior', 'debuff', 'enemy', 'physical', 4, 'strength', 0.2, 8, 0, 6, 'eff_sunder_armor', 1.0),
('execute', '斩杀', '对低血量敌人造成大量伤害。', 'warrior', 'attack', 'enemy_lowest_hp', 'physical', 20, 'strength', 0.8, 18, 0, 8, NULL, 1.0),
('shield_wall', '盾墙', '大幅减少受到的伤害。', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 25, 10, 10, 'eff_shield_wall', 1.0),
('battle_shout', '战斗怒吼', '提升全队攻击力。', 'warrior', 'buff', 'ally_all', NULL, 0, NULL, 0, 10, 5, 3, 'eff_battle_shout', 1.0);

-- 法师技能 (法力消耗4~18)
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('fireball', '火球术', '发射火球，有几率点燃敌人。', 'mage', 'attack', 'enemy', 'fire', 10, 'intellect', 0.5, 6, 0, 1, 'eff_ignite', 0.3),
('frostbolt', '寒冰箭', '发射寒冰箭，减缓敌人。', 'mage', 'attack', 'enemy', 'frost', 8, 'intellect', 0.4, 5, 0, 1, 'eff_frostbite', 0.5),
('arcane_missiles', '奥术飞弹', '发射多道奥术飞弹。', 'mage', 'attack', 'enemy', 'magic', 12, 'intellect', 0.6, 8, 2, 4, NULL, 1.0),
('flamestrike', '烈焰风暴', '对所有敌人造成火焰伤害。', 'mage', 'attack', 'enemy_all', 'fire', 10, 'intellect', 0.4, 12, 3, 6, 'eff_ignite', 0.2),
('pyroblast', '炎爆术', '施放巨大的火球。', 'mage', 'attack', 'enemy', 'fire', 25, 'intellect', 1.0, 18, 5, 8, 'eff_ignite', 0.8),
('ice_barrier', '寒冰护体', '创造吸收伤害的护盾。', 'mage', 'shield', 'self', NULL, 15, 'intellect', 0.5, 14, 8, 10, 'eff_ice_barrier', 1.0),
('arcane_intellect', '奥术智慧', '提升智力。', 'mage', 'buff', 'ally', NULL, 0, NULL, 0, 8, 0, 2, 'eff_arcane_intellect', 1.0);

-- 盗贼技能 (能量消耗15~40)
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('sinister_strike', '邪恶攻击', '快速的攻击。', 'rogue', 'attack', 'enemy', 'physical', 6, 'agility', 0.4, 20, 0, 1, NULL, 1.0),
('backstab', '背刺', '从背后攻击造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 12, 'agility', 0.6, 35, 0, 1, NULL, 1.0),
('deadly_poison', '致命毒药', '使敌人中毒。', 'rogue', 'dot', 'enemy', 'nature', 2, 'agility', 0.1, 25, 0, 3, 'eff_poison', 1.0),
('eviscerate', '剔骨', '造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 18, 'agility', 0.7, 40, 0, 4, NULL, 1.0),
('kidney_shot', '肾击', '眩晕敌人。', 'rogue', 'control', 'enemy', 'physical', 4, 'agility', 0.2, 30, 4, 6, 'eff_stun', 1.0),
('blade_flurry', '剑刃乱舞', '攻击力大幅提升。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 35, 8, 8, 'eff_blade_flurry', 1.0),
('vanish', '消失', '进入潜行状态。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 50, 10, 10, 'eff_stealth', 1.0);

-- 牧师技能 (法力消耗4~12)
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('smite', '惩击', '用神圣能量攻击敌人。', 'priest', 'attack', 'enemy', 'holy', 7, 'intellect', 0.4, 4, 0, 1, NULL, 1.0),
('shadow_word_pain', '暗言术:痛', '对敌人施加持续伤害。', 'priest', 'dot', 'enemy', 'shadow', 2, 'intellect', 0.15, 5, 0, 1, 'eff_sw_pain', 1.0),
('lesser_heal', '次级治疗术', '恢复生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 8, 'spirit', 0.4, 5, 0, 1, NULL, 1.0),
('renew', '恢复', '持续恢复生命。', 'priest', 'hot', 'ally_lowest_hp', 'holy', 3, 'spirit', 0.2, 6, 0, 2, 'eff_renew', 1.0),
('heal', '治疗术', '恢复大量生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 15, 'spirit', 0.6, 10, 0, 4, NULL, 1.0),
('inner_fire', '心灵之火', '提升防御力。', 'priest', 'buff', 'self', NULL, 0, NULL, 0, 8, 0, 4, 'eff_inner_fire', 1.0),
('power_word_shield', '真言术:盾', '创造吸收伤害的护盾。', 'priest', 'shield', 'ally_lowest_hp', 'holy', 12, 'spirit', 0.5, 10, 4, 6, 'eff_pw_shield', 1.0),
('flash_heal', '快速治疗', '快速恢复生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 12, 'spirit', 0.5, 8, 0, 8, NULL, 1.0),
('silence', '沉默', '使敌人无法施法。', 'priest', 'control', 'enemy', 'shadow', 0, NULL, 0, 10, 6, 10, 'eff_silence', 1.0);

-- 通用技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('basic_attack', '普通攻击', '基础的物理攻击。', NULL, 'attack', 'enemy', 'physical', 0, 'strength', 0.5, 0, 0, 1, NULL, 1.0);

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

-- ═══════════════════════════════════════════════════════════
-- 怪物数据 (小数值设计：HP 15~300, 攻击3~50, 经验5~80)
-- ═══════════════════════════════════════════════════════════

-- 艾尔文森林 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('wolf', 'elwynn', '森林狼', 1, 'normal', 15, 3, 1, 5, 1, 2, 100),
('young_boar', 'elwynn', '小野猪', 1, 'normal', 12, 3, 1, 4, 1, 1, 100),
('kobold_worker', 'elwynn', '狗头人矿工', 2, 'normal', 18, 4, 1, 6, 1, 2, 80),
('kobold_tunneler', 'elwynn', '狗头人掘地工', 3, 'normal', 22, 5, 2, 7, 1, 3, 60),
('defias_thug', 'elwynn', '迪菲亚暴徒', 4, 'normal', 26, 6, 2, 8, 2, 3, 50),
('defias_bandit', 'elwynn', '迪菲亚劫匪', 5, 'normal', 30, 7, 3, 10, 2, 4, 40),
('murloc', 'elwynn', '鱼人', 3, 'normal', 20, 5, 2, 7, 1, 3, 70),
('murloc_warrior', 'elwynn', '鱼人战士', 5, 'normal', 28, 7, 3, 9, 2, 4, 30),
('prowler', 'elwynn', '潜伏者', 6, 'normal', 32, 8, 3, 11, 2, 5, 25),
('hogger', 'elwynn', '霍格', 8, 'elite', 80, 12, 5, 35, 5, 12, 5);

-- 西部荒野 (10-20级) - HP: 35~80, 攻击: 9~18
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('harvest_golem', 'westfall', '收割傀儡', 10, 'normal', 40, 10, 5, 14, 2, 5, 100),
('defias_rogue', 'westfall', '迪菲亚盗贼', 11, 'normal', 44, 11, 5, 16, 3, 6, 80),
('defias_highwayman', 'westfall', '迪菲亚拦路贼', 12, 'normal', 48, 12, 6, 18, 3, 7, 60),
('gnoll_brute', 'westfall', '豺狼人蛮兵', 13, 'normal', 52, 13, 6, 20, 4, 8, 50),
('gnoll_mystic', 'westfall', '豺狼人秘法师', 14, 'normal', 45, 15, 5, 22, 4, 9, 40),
('defias_overlord', 'westfall', '迪菲亚霸主', 16, 'elite', 120, 20, 10, 50, 10, 20, 5),
('dust_devil', 'westfall', '尘土恶魔', 15, 'normal', 55, 14, 6, 24, 5, 10, 30);

-- 暮色森林 (20-30级) - HP: 55~150, 攻击: 16~35
INSERT OR REPLACE INTO monsters (id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight) VALUES
('skeleton_warrior', 'duskwood', '骷髅战士', 20, 'normal', 65, 18, 9, 28, 5, 10, 100),
('skeleton_mage', 'duskwood', '骷髅法师', 21, 'normal', 55, 22, 7, 30, 5, 11, 80),
('ghoul', 'duskwood', '食尸鬼', 22, 'normal', 70, 20, 9, 32, 6, 12, 70),
('dire_wolf', 'duskwood', '恐狼', 23, 'normal', 75, 22, 10, 34, 6, 13, 60),
('worgen', 'duskwood', '狼人', 24, 'normal', 85, 25, 11, 36, 7, 14, 50),
('worgen_alpha', 'duskwood', '狼人首领', 26, 'elite', 180, 35, 15, 65, 12, 25, 5),
('abomination', 'duskwood', '憎恶', 28, 'elite', 220, 40, 18, 75, 15, 30, 3),
('stitches', 'duskwood', '缝合怪', 30, 'boss', 350, 50, 22, 100, 25, 50, 1);

-- ═══════════════════════════════════════════════════════════
-- 物品数据
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 物品数据 (小数值设计：药水恢复10~40，装备加成1~12)
-- ═══════════════════════════════════════════════════════════

-- 消耗品 (恢复量：初级10，中级20，高级40)
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, stackable, max_stack, sell_price, buy_price, effect_type, effect_value) VALUES
('minor_healing_potion', '初级治疗药水', '恢复10点生命值。', 'consumable', 'potion', 'common', 1, 20, 1, 3, 'heal_hp', 10),
('healing_potion', '治疗药水', '恢复20点生命值。', 'consumable', 'potion', 'uncommon', 1, 20, 2, 8, 'heal_hp', 20),
('greater_healing_potion', '强效治疗药水', '恢复40点生命值。', 'consumable', 'potion', 'rare', 1, 20, 5, 20, 'heal_hp', 40),
('minor_mana_potion', '初级法力药水', '恢复8点法力值。', 'consumable', 'potion', 'common', 1, 20, 1, 3, 'heal_mp', 8),
('mana_potion', '法力药水', '恢复16点法力值。', 'consumable', 'potion', 'uncommon', 1, 20, 2, 8, 'heal_mp', 16);

-- 装备 - 武器 (攻击加成1~10)
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, level_required, slot, sell_price, attack, strength) VALUES
('worn_sword', '破旧的剑', '一把破旧的铁剑。', 'equipment', 'weapon', 'common', 1, 'main_hand', 2, 1, 0),
('militia_sword', '民兵之剑', '联盟民兵的标准武器。', 'equipment', 'weapon', 'common', 5, 'main_hand', 5, 2, 1),
('outlaw_sabre', '逃犯军刀', '从迪菲亚成员身上缴获。', 'equipment', 'weapon', 'uncommon', 10, 'main_hand', 15, 4, 2),
('corpsemaker', '尸体收割者', '沉重的双手斧。', 'equipment', 'weapon', 'rare', 20, 'main_hand', 40, 8, 4);

-- 装备 - 护甲 (防御加成1~8)
INSERT OR REPLACE INTO items (id, name, description, type, subtype, quality, level_required, slot, sell_price, defense, stamina) VALUES
('worn_leather_vest', '破旧皮甲', '一件破旧的皮甲。', 'equipment', 'armor', 'common', 1, 'chest', 1, 1, 0),
('militia_chain_vest', '民兵锁甲', '联盟民兵的标准护甲。', 'equipment', 'armor', 'common', 5, 'chest', 4, 2, 1),
('defias_leather_vest', '迪菲亚皮甲', '迪菲亚兄弟会的制服。', 'equipment', 'armor', 'uncommon', 10, 'chest', 12, 4, 2),
('blackened_defias_armor', '黑化迪菲亚护甲', '高级成员的护甲。', 'equipment', 'armor', 'rare', 15, 'chest', 30, 6, 3);

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

-- ═══════════════════════════════════════════════════════════
-- 游戏公式配置 (玩家可查询)
-- ═══════════════════════════════════════════════════════════
-- 所有战斗计算规则对玩家透明，可通过游戏内命令查询

INSERT OR REPLACE INTO game_formulas (id, category, name, formula, description, variables, example, display_order) VALUES
-- ═══════════════════════════════════════════════════════════
-- 属性转换公式
-- ═══════════════════════════════════════════════════════════
('attr_max_hp', 'attribute', '最大生命值', 
 'class_base_hp + stamina × 10 + level × hp_per_level + equipment_bonus',
 '计算角色的最大生命值。耐力是最主要的生命值来源。',
 '{"class_base_hp":"职业基础HP(战士120,法师65等)","stamina":"耐力值","level":"角色等级","hp_per_level":"每级HP成长(战士12,法师4等)","equipment_bonus":"装备HP加成"}',
 '战士(30级,耐力80): 120 + 800 + 360 + 0 = 1280 HP', 1),

('attr_max_mp', 'attribute', '最大法力值',
 'class_base_mp + intellect × 5 + level × mp_per_level + equipment_bonus',
 '计算角色的最大法力值。智力是最主要的法力来源。仅法力职业有效。',
 '{"class_base_mp":"职业基础MP","intellect":"智力值","level":"角色等级","mp_per_level":"每级MP成长","equipment_bonus":"装备MP加成"}',
 '法师(30级,智力100): 120 + 500 + 450 + 0 = 1070 MP', 2),

('attr_phys_attack', 'attribute', '物理攻击力',
 'base_attack + strength × 2.0 + agility × 0.5 + weapon_damage + equipment_bonus',
 '计算角色的物理攻击力。力量是主要来源，敏捷提供少量加成。',
 '{"base_attack":"职业基础攻击(战士15,盗贼12,法师5等)","strength":"力量值","agility":"敏捷值","weapon_damage":"武器伤害","equipment_bonus":"其他装备加成"}',
 '盗贼(力量40,敏捷100): 12 + 80 + 50 = 142', 3),

('attr_spell_power', 'attribute', '法术攻击力',
 'base_spell + intellect × 1.5 + equipment_bonus',
 '计算角色的法术攻击力。智力是唯一的属性来源。',
 '{"base_spell":"职业基础法伤(法师20,术士18,牧师15等)","intellect":"智力值","equipment_bonus":"装备法伤加成"}',
 '法师(智力120): 20 + 180 = 200', 4),

('attr_armor', 'attribute', '护甲值',
 'base_armor + stamina × 0.5 + agility × 0.3 + equipment_armor',
 '计算角色的护甲值。装备提供主要护甲，耐力和敏捷提供少量加成。',
 '{"base_armor":"职业基础护甲(战士50,法师10等)","stamina":"耐力值","agility":"敏捷值","equipment_armor":"装备护甲"}',
 '战士(耐力80,敏捷40): 50 + 40 + 12 = 102', 5),

-- ═══════════════════════════════════════════════════════════
-- 战斗判定公式
-- ═══════════════════════════════════════════════════════════
('combat_phys_crit', 'combat', '物理暴击率',
 '5% + agility ÷ 20 + equipment_bonus% (上限50%)',
 '计算物理攻击的暴击概率。每20点敏捷增加1%暴击率。',
 '{"agility":"敏捷值","equipment_bonus":"装备暴击加成%"}',
 '敏捷100: 5% + 5% = 10%', 10),

('combat_spell_crit', 'combat', '法术暴击率',
 '5% + intellect ÷ 30 + equipment_bonus% (上限50%)',
 '计算法术攻击的暴击概率。每30点智力增加1%暴击率。',
 '{"intellect":"智力值","equipment_bonus":"装备法术暴击加成%"}',
 '智力120: 5% + 4% = 9%', 11),

('combat_crit_damage', 'combat', '暴击伤害倍率',
 '物理: 150% + strength ÷ 100 × 10%; 法术: 150% + equipment_bonus% (上限250%)',
 '暴击时的伤害倍率。物理暴击伤害受力量影响，法术暴击伤害仅受装备影响。',
 '{"strength":"力量值(仅影响物理暴击)","equipment_bonus":"装备暴击伤害加成%"}',
 '力量150: 150% + 15% = 165% 暴击伤害', 12),

('combat_dodge', 'combat', '闪避率',
 '5% + agility ÷ 25 + equipment_bonus% + racial_bonus% (上限30%)',
 '计算物理攻击的闪避概率。每25点敏捷增加1%闪避率。',
 '{"agility":"敏捷值","equipment_bonus":"装备闪避加成%","racial_bonus":"种族加成(暗夜精灵+2%)"}',
 '暗夜精灵(敏捷100): 5% + 4% + 2% = 11%', 13),

('combat_hit', 'combat', '命中率',
 '95% - (target_level - self_level) × 1% + equipment_bonus% (范围75%~100%)',
 '计算攻击命中的概率。攻击比自己高级的目标会降低命中率。',
 '{"target_level":"目标等级","self_level":"自身等级","equipment_bonus":"装备命中加成%"}',
 '30级攻击35级目标: 95% - 5% = 90%', 14),

('combat_armor_reduction', 'combat', '护甲减伤%',
 'armor ÷ (armor + 400 + attacker_level × 10) × 100% (上限75%)',
 '根据护甲值计算物理伤害减免百分比。减伤存在上限。',
 '{"armor":"护甲值","attacker_level":"攻击者等级"}',
 '护甲500,被30级怪攻击: 500÷1200 = 41.7%减伤', 15),

('combat_resist_reduction', 'combat', '法术抗性减伤%',
 'resistance ÷ (resistance + attacker_level × 5) × 75% (上限75%)',
 '根据抗性值计算对应法术类型的伤害减免。',
 '{"resistance":"对应类型抗性值(火焰/冰霜/暗影等)","attacker_level":"攻击者等级"}',
 '冰霜抗性100,被30级法师攻击: 100÷250×75% = 30%减伤', 16),

-- ═══════════════════════════════════════════════════════════
-- 技能伤害公式
-- ═══════════════════════════════════════════════════════════
('skill_base_damage', 'skill', '技能基础伤害',
 '(base_value + scaling_stat × scaling_ratio) × skill_level_mult',
 '计算技能的基础伤害值，不含暴击和减伤。',
 '{"base_value":"技能基础值(每个技能不同)","scaling_stat":"成长属性值(力量/敏捷/智力/精神)","scaling_ratio":"成长系数(每个技能不同)","skill_level_mult":"1+(技能等级-1)×0.1"}',
 '5级火球术(基础50,智力120,系数0.8): (50+96)×1.4 = 204', 20),

('skill_final_damage', 'skill', '最终伤害',
 'base_damage × crit_mult × (1 - target_reduction%) × random(0.9~1.1)',
 '计算造成的实际伤害，包含暴击、减伤和随机波动。',
 '{"base_damage":"技能基础伤害","crit_mult":"暴击倍率(非暴击=1)","target_reduction":"目标减伤%(护甲或抗性)","random":"±10%随机波动"}',
 '技能204,暴击1.5倍,目标20%减伤: 204×1.5×0.8×1.0 = 244.8', 21),

('skill_healing', 'skill', '治疗量',
 '(base_value + spirit × scaling_ratio) × skill_level_mult × (1 + healing_bonus%)',
 '计算治疗技能的恢复量。精神是主要的治疗成长属性。',
 '{"base_value":"技能基础治疗量","spirit":"精神值","scaling_ratio":"成长系数","skill_level_mult":"技能等级系数","healing_bonus":"治疗效果加成%(装备/种族)"}',
 '治疗术(基础50,精神80,系数1.5): (50+120)×1.0 = 170 HP', 22),

-- ═══════════════════════════════════════════════════════════
-- 能量回复公式
-- ═══════════════════════════════════════════════════════════
('resource_mana_regen_combat', 'resource', '法力回复(战斗中)',
 'spirit × class_regen_pct% × max_mp ÷ 100 (每回合)',
 '战斗中每回合的法力恢复量。不同职业有不同的恢复系数。',
 '{"spirit":"精神值","class_regen_pct":"职业恢复系数(牧师0.8,德鲁伊0.6,法师0.5,圣骑士0.4,猎人0.3)","max_mp":"最大法力值"}',
 '牧师(精神80,MP500,系数0.8%): 80×0.008×500 = 3.2/回合', 30),

('resource_mana_regen_rest', 'resource', '法力回复(战斗外)',
 '战斗中回复量 × 3 (每秒)',
 '战斗外每秒的法力恢复量，是战斗中的3倍。',
 '{}',
 '战斗中3.2/回合 → 战斗外9.6/秒', 31),

('resource_rage_gain_attack', 'resource', '怒气获得(攻击)',
 '普通攻击命中: +5; 造成暴击: 额外+10',
 '战士通过攻击获得怒气的方式。',
 '{}',
 '普通攻击+5怒气，暴击共+15怒气', 32),

('resource_rage_gain_damage', 'resource', '怒气获得(受伤)',
 'damage_taken ÷ max_hp × 20',
 '战士通过受到伤害获得怒气。受伤越重获得越多。',
 '{"damage_taken":"受到的伤害值","max_hp":"最大生命值"}',
 '受到200伤害(最大HP1000): 200÷1000×20 = +4怒气', 33),

('resource_rage_decay', 'resource', '怒气衰减',
 '脱战后每回合: -5怒气',
 '战士脱离战斗后怒气会逐渐消退。',
 '{}',
 '50怒气脱战，10回合后归零', 34),

('resource_energy_regen', 'resource', '能量恢复(盗贼)',
 '每回合固定: +20能量',
 '盗贼的能量恢复方式。能量恢复快但技能消耗也高。',
 '{}',
 '每回合+20，5回合恢复满100能量', 35),

('resource_hp_regen_combat', 'resource', '生命回复(战斗中)',
 'spirit × 0.2 (每回合)',
 '战斗中每回合的自然生命恢复量。精神提供少量恢复。',
 '{"spirit":"精神值"}',
 '精神50: 每回合恢复10 HP', 36),

('resource_hp_regen_rest', 'resource', '生命回复(战斗外)',
 'spirit × 1.0 + max_hp × 1% (每秒)',
 '战斗外每秒的生命恢复量。',
 '{"spirit":"精神值","max_hp":"最大生命值"}',
 '精神50,最大HP1000: 50+10 = 60 HP/秒', 37),

('resource_hp_regen_troll', 'resource', '巨魔再生(种族)',
 '每回合额外恢复: max_hp × 2%',
 '巨魔种族的特殊被动，提供额外的生命恢复。',
 '{"max_hp":"最大生命值"}',
 '最大HP500: 每回合额外恢复10 HP', 38),

-- ═══════════════════════════════════════════════════════════
-- 经验与等级公式
-- ═══════════════════════════════════════════════════════════
('exp_required', 'progression', '升级所需经验',
 'base_exp × level × 1.1^(level-1)',
 '计算升到下一级所需的经验值。等级越高需要的经验越多。',
 '{"base_exp":"基础经验值(100)","level":"当前等级"}',
 '10级升11级: 100×10×1.1^9 = 2358 经验', 40),

('exp_gain', 'progression', '击杀经验获取',
 'monster_exp × (1 + level_diff × 0.1) × zone_modifier × racial_bonus',
 '计算击杀怪物获得的经验值。击杀高级怪物获得更多经验。',
 '{"monster_exp":"怪物基础经验","level_diff":"怪物等级-角色等级(正数有加成，负数有惩罚)","zone_modifier":"区域经验倍率","racial_bonus":"人类种族+10%"}',
 '基础100经验,怪物高3级,人类: 100×1.3×1.1 = 143 经验', 41),

('exp_penalty', 'progression', '低级怪物经验惩罚',
 '角色等级 - 怪物等级 > 5: 经验 × (1 - (等级差-5) × 0.1), 最低10%',
 '击杀比自己低太多级的怪物会减少经验获取。',
 '{"level_diff":"角色等级-怪物等级"}',
 '30级打15级怪: 经验×(1-10×0.1) = 0 (最低10%)', 42),

('exp_distribution', 'progression', '小队经验分配',
 '总经验 ÷ 参战角色数 (平均分配)',
 '小队中所有参战角色平均分配经验值。',
 '{}',
 '获得1000经验,3人参战: 每人333经验', 43);


