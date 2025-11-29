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
-- 效果配置 (回合制战斗优化版)
-- ═══════════════════════════════════════════════════════════
-- 设计原则:
-- 1. 持续时间以"回合"计算
-- 2. DOT/HOT伤害保持小数值 (2-5/回合)
-- 3. 护盾吸收量与HP匹配 (10-20点)

INSERT OR REPLACE INTO effects (id, name, description, type, is_buff, is_stackable, max_stacks, duration, value_type, value, stat_affected, damage_type, can_dispel) VALUES
-- ═══════════════════════════════════════════════════════════
-- 增益效果 (Buff)
-- ═══════════════════════════════════════════════════════════
-- 战士
('eff_shield_wall', '盾墙', '受到的伤害降低50%', 'stat_mod', 1, 0, 1, 3, 'percent', -50, 'damage_taken', NULL, 0),
('eff_battle_shout', '战斗怒吼', '攻击力提升10%', 'stat_mod', 1, 0, 1, 5, 'percent', 10, 'attack', NULL, 1),
-- 法师
('eff_ice_barrier', '寒冰护体', '吸收15点伤害', 'shield', 1, 0, 1, 5, 'flat', 15, NULL, NULL, 1),
('eff_arcane_intellect', '奥术智慧', '智力提升10%', 'stat_mod', 1, 0, 1, 99, 'percent', 10, 'intellect', NULL, 1),
-- 盗贼
('eff_blade_flurry', '剑刃乱舞', '攻击力提升20%', 'stat_mod', 1, 0, 1, 4, 'percent', 20, 'attack', NULL, 1),
('eff_evasion', '闪避', '闪避率提升50%', 'stat_mod', 1, 0, 1, 3, 'percent', 50, 'dodge_rate', NULL, 0),
-- 牧师
('eff_pw_shield', '真言术:盾', '吸收12点伤害', 'shield', 1, 0, 1, 4, 'flat', 12, NULL, NULL, 1),
('eff_renew', '恢复', '每回合恢复3点生命', 'hot', 1, 0, 1, 4, 'flat', 3, NULL, 'holy', 1),
('eff_inner_fire', '心灵之火', '防御力提升15%', 'stat_mod', 1, 0, 1, 99, 'percent', 15, 'defense', NULL, 1),
-- 圣骑士
('eff_blessing_might', '力量祝福', '攻击力+3', 'stat_mod', 1, 0, 1, 99, 'flat', 3, 'attack', NULL, 1),
('eff_divine_shield', '圣盾术', '免疫所有伤害', 'invulnerable', 1, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
('eff_consecration', '奉献', '每回合对敌人造成神圣伤害', 'dot', 1, 0, 1, 4, 'flat', 3, NULL, 'holy', 1),
-- 猎人
('eff_rapid_fire', '急速射击', '攻击速度提升30%', 'stat_mod', 1, 0, 1, 4, 'percent', 30, 'attack_speed', NULL, 1),
('eff_feign_death', '假死', '无法被攻击', 'untargetable', 1, 0, 1, 1, NULL, NULL, NULL, NULL, 0),
('eff_bestial_wrath', '狂野怒火', '伤害提升50%', 'stat_mod', 1, 0, 1, 3, 'percent', 50, 'damage_dealt', NULL, 0),
-- 术士
('eff_drain_life', '吸取生命', '造成伤害的50%转化为生命', 'lifesteal', 1, 0, 1, 1, 'percent', 50, NULL, NULL, 0),
('eff_soul_link', '灵魂链接', '受到伤害降低30%', 'stat_mod', 1, 0, 1, 99, 'percent', -30, 'damage_taken', NULL, 0),
-- 德鲁伊
('eff_rejuvenation', '回春术', '每回合恢复4点生命', 'hot', 1, 0, 1, 4, 'flat', 4, NULL, 'nature', 1),
('eff_regrowth', '愈合', '每回合恢复2点生命', 'hot', 1, 0, 1, 3, 'flat', 2, NULL, 'nature', 1),
('eff_barkskin', '树皮术', '受到伤害降低25%', 'stat_mod', 1, 0, 1, 3, 'percent', -25, 'damage_taken', NULL, 0),
('eff_tranquility', '宁静', '每回合恢复5点生命', 'hot', 1, 0, 1, 3, 'flat', 5, NULL, 'nature', 0),
-- 萨满
('eff_windfury', '风怒武器', '20%几率额外攻击', 'proc', 1, 0, 1, 99, 'percent', 20, 'extra_attack', NULL, 1),
('eff_bloodlust', '嗜血', '攻击速度提升30%', 'stat_mod', 1, 0, 1, 4, 'percent', 30, 'attack_speed', NULL, 0),

-- ═══════════════════════════════════════════════════════════
-- 减益效果 (Debuff)
-- ═══════════════════════════════════════════════════════════
-- 控制效果
('eff_stun', '眩晕', '无法行动', 'stun', 0, 0, 1, 1, NULL, NULL, NULL, NULL, 1),
('eff_silence', '沉默', '无法施放法术', 'silence', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 1),
('eff_fear', '恐惧', '无法控制行动', 'stun', 0, 0, 1, 2, NULL, NULL, 'shadow', NULL, 1),
('eff_root', '缠绕', '无法移动和行动', 'stun', 0, 0, 1, 2, NULL, NULL, 'nature', NULL, 1),
('eff_interrupt', '打断', '施法被打断', 'interrupt', 0, 0, 1, 1, NULL, NULL, NULL, NULL, 0),
-- 减益效果
('eff_slow', '减速', '攻击速度降低30%', 'slow', 0, 0, 1, 3, 'percent', -30, 'attack_speed', NULL, 1),
('eff_frostbite', '冻伤', '攻击速度降低20%', 'stat_mod', 0, 0, 1, 3, 'percent', -20, 'attack_speed', 'frost', 1),
('eff_sunder_armor', '破甲', '防御力降低20%', 'stat_mod', 0, 1, 5, 5, 'percent', -20, 'defense', NULL, 0),
('eff_curse_weakness', '虚弱诅咒', '造成的伤害降低20%', 'stat_mod', 0, 0, 1, 4, 'percent', -20, 'damage_dealt', 'shadow', 1),
('eff_taunt', '嘲讽', '强制攻击施法者', 'taunt', 0, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
-- DOT效果 (持续伤害)
('eff_rend', '撕裂', '每回合2点流血伤害', 'dot', 0, 1, 3, 3, 'flat', 2, NULL, 'physical', 1),
('eff_rupture', '割裂', '每回合3点流血伤害', 'dot', 0, 0, 1, 6, 'flat', 3, NULL, 'physical', 1),
('eff_ignite', '点燃', '每回合2点火焰伤害', 'dot', 0, 1, 5, 3, 'flat', 2, NULL, 'fire', 1),
('eff_sw_pain', '暗言术:痛', '每回合3点暗影伤害', 'dot', 0, 0, 1, 4, 'flat', 3, NULL, 'shadow', 1),
('eff_poison', '中毒', '每回合2点自然伤害', 'dot', 0, 1, 5, 5, 'flat', 2, NULL, 'nature', 1),
('eff_serpent_sting', '毒蛇钉刺', '每回合2点自然伤害', 'dot', 0, 0, 1, 5, 'flat', 2, NULL, 'nature', 1),
('eff_corruption', '腐蚀术', '每回合2点暗影伤害', 'dot', 0, 0, 1, 6, 'flat', 2, NULL, 'shadow', 1),
('eff_agony', '痛苦诅咒', '每回合1-4点递增伤害', 'dot', 0, 0, 1, 5, 'flat', 2, NULL, 'shadow', 1),
('eff_immolate', '献祭', '每回合3点火焰伤害', 'dot', 0, 0, 1, 5, 'flat', 3, NULL, 'fire', 1),
('eff_moonfire', '月火术', '每回合2点自然伤害', 'dot', 0, 0, 1, 4, 'flat', 2, NULL, 'nature', 1),
('eff_flame_shock', '烈焰震击', '每回合2点火焰伤害', 'dot', 0, 0, 1, 4, 'flat', 2, NULL, 'fire', 1);

-- ═══════════════════════════════════════════════════════════
-- 技能数据 (扩展版)
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 技能数据 (回合制战斗优化版)
-- ═══════════════════════════════════════════════════════════
-- 设计原则:
-- 1. 所有技能适配自动战斗，无需玩家操作
-- 2. 冷却时间以"回合"计算
-- 3. 资源消耗平衡，保证技能循环
-- 4. 强力技能必须有冷却

-- ═══════════════════════════════════════════════════════════
-- 战士技能 (怒气系统: 上限100, 攻击+5, 受伤+怒气)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('heroic_strike', '英勇打击', '强力的武器攻击。', 'warrior', 'attack', 'enemy', 'physical', 8, 'strength', 0.5, 10, 0, 1, NULL, 1.0),
('charge', '冲锋', '冲向敌人造成伤害并眩晕1回合。', 'warrior', 'attack', 'enemy', 'physical', 5, 'strength', 0.3, 0, 4, 1, 'eff_stun', 1.0),
('rend', '撕裂', '造成流血，持续3回合。', 'warrior', 'dot', 'enemy', 'physical', 2, 'strength', 0.15, 8, 0, 2, 'eff_rend', 1.0),
('thunder_clap', '雷霆一击', '对所有敌人造成伤害并减速。', 'warrior', 'attack', 'enemy_all', 'physical', 6, 'strength', 0.3, 15, 3, 4, 'eff_slow', 0.8),
('sunder_armor', '破甲攻击', '降低敌人防御20%，可叠加。', 'warrior', 'debuff', 'enemy', 'physical', 4, 'strength', 0.2, 12, 0, 6, 'eff_sunder_armor', 1.0),
('execute', '斩杀', '对HP<30%的敌人造成巨额伤害。', 'warrior', 'attack', 'enemy_lowest_hp', 'physical', 25, 'strength', 1.0, 20, 3, 8, NULL, 1.0),
('shield_wall', '盾墙', '3回合内受到伤害降低50%。', 'warrior', 'buff', 'self', NULL, 0, NULL, 0, 30, 10, 10, 'eff_shield_wall', 1.0),
('battle_shout', '战斗怒吼', '全队攻击力提升10%，持续5回合。', 'warrior', 'buff', 'ally_all', NULL, 0, NULL, 0, 15, 6, 3, 'eff_battle_shout', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 法师技能 (法力系统: 基础40, 每回合+精神×0.5%)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('fireball', '火球术', '发射火球，30%几率点燃敌人。', 'mage', 'attack', 'enemy', 'fire', 10, 'intellect', 0.5, 5, 0, 1, 'eff_ignite', 0.3),
('frostbolt', '寒冰箭', '发射寒冰箭，50%几率减速敌人。', 'mage', 'attack', 'enemy', 'frost', 8, 'intellect', 0.4, 4, 0, 1, 'eff_frostbite', 0.5),
('arcane_missiles', '奥术飞弹', '发射3道奥术飞弹。', 'mage', 'attack', 'enemy', 'magic', 12, 'intellect', 0.6, 7, 2, 4, NULL, 1.0),
('flamestrike', '烈焰风暴', '对所有敌人造成火焰伤害。', 'mage', 'attack', 'enemy_all', 'fire', 8, 'intellect', 0.4, 10, 3, 6, 'eff_ignite', 0.2),
('pyroblast', '炎爆术', '蓄力后释放巨大火球。', 'mage', 'attack', 'enemy', 'fire', 22, 'intellect', 0.9, 15, 4, 8, 'eff_ignite', 0.8),
('ice_barrier', '寒冰护体', '创造可吸收15点伤害的护盾。', 'mage', 'shield', 'self', NULL, 15, 'intellect', 0.5, 12, 6, 10, 'eff_ice_barrier', 1.0),
('arcane_intellect', '奥术智慧', '提升目标智力10%，持续整场战斗。', 'mage', 'buff', 'ally', NULL, 0, NULL, 0, 6, 0, 2, 'eff_arcane_intellect', 1.0),
('blizzard', '暴风雪', '对所有敌人造成冰霜伤害并减速。', 'mage', 'attack', 'enemy_all', 'frost', 6, 'intellect', 0.3, 12, 4, 12, 'eff_frostbite', 0.6);

-- ═══════════════════════════════════════════════════════════
-- 盗贼技能 (能量系统: 上限100, 每回合+20)
-- 调整: 降低能量消耗，添加冷却平衡
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('sinister_strike', '邪恶攻击', '快速攻击，积累连击点。', 'rogue', 'attack', 'enemy', 'physical', 6, 'agility', 0.4, 15, 0, 1, NULL, 1.0),
('ambush', '伏击', '对HP>80%的敌人伤害翻倍。', 'rogue', 'attack', 'enemy', 'physical', 10, 'agility', 0.5, 25, 2, 1, NULL, 1.0),
('deadly_poison', '致命毒药', '使敌人中毒，持续5回合。', 'rogue', 'dot', 'enemy', 'nature', 2, 'agility', 0.1, 20, 0, 3, 'eff_poison', 1.0),
('eviscerate', '剔骨', '消耗连击点造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 15, 'agility', 0.6, 30, 2, 4, NULL, 1.0),
('kidney_shot', '肾击', '眩晕敌人1回合。', 'rogue', 'control', 'enemy', 'physical', 4, 'agility', 0.2, 25, 4, 6, 'eff_stun', 1.0),
('blade_flurry', '剑刃乱舞', '4回合内攻击力提升20%。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 30, 6, 8, 'eff_blade_flurry', 1.0),
('evasion', '闪避', '3回合内闪避率提升50%。', 'rogue', 'buff', 'self', NULL, 0, NULL, 0, 35, 8, 10, 'eff_evasion', 1.0),
('rupture', '割裂', '造成强力流血，持续6回合。', 'rogue', 'dot', 'enemy', 'physical', 3, 'agility', 0.2, 35, 0, 12, 'eff_rupture', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 牧师技能 (法力系统: 基础35, 每回合+精神×0.8%)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('smite', '惩击', '用神圣能量攻击敌人。', 'priest', 'attack', 'enemy', 'holy', 7, 'intellect', 0.4, 4, 0, 1, NULL, 1.0),
('shadow_word_pain', '暗言术:痛', '对敌人施加暗影DOT，持续4回合。', 'priest', 'dot', 'enemy', 'shadow', 3, 'intellect', 0.15, 4, 0, 1, 'eff_sw_pain', 1.0),
('lesser_heal', '次级治疗术', '恢复少量生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 8, 'spirit', 0.4, 4, 0, 1, NULL, 1.0),
('renew', '恢复', '持续恢复生命，持续4回合。', 'priest', 'hot', 'ally_lowest_hp', 'holy', 3, 'spirit', 0.2, 5, 0, 2, 'eff_renew', 1.0),
('heal', '治疗术', '恢复大量生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 15, 'spirit', 0.6, 8, 0, 4, NULL, 1.0),
('inner_fire', '心灵之火', '提升自身防御力15%，持续整场战斗。', 'priest', 'buff', 'self', NULL, 0, NULL, 0, 6, 0, 4, 'eff_inner_fire', 1.0),
('power_word_shield', '真言术:盾', '为血量最低队友创造吸收12点伤害的护盾。', 'priest', 'shield', 'ally_lowest_hp', 'holy', 12, 'spirit', 0.5, 8, 4, 6, 'eff_pw_shield', 1.0),
('flash_heal', '快速治疗', '快速恢复中等生命值。', 'priest', 'heal', 'ally_lowest_hp', 'holy', 12, 'spirit', 0.5, 6, 0, 8, NULL, 1.0),
('prayer_of_healing', '治疗祷言', '恢复全队生命值。', 'priest', 'heal', 'ally_all', 'holy', 6, 'spirit', 0.3, 12, 4, 12, NULL, 1.0);

-- ═══════════════════════════════════════════════════════════
-- 圣骑士技能 (法力系统: 基础20, 混合输出/治疗/坦克)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('crusader_strike', '十字军打击', '神圣武器攻击。', 'paladin', 'attack', 'enemy', 'physical', 8, 'strength', 0.5, 3, 0, 1, NULL, 1.0),
('judgement', '审判', '释放神圣审判造成伤害。', 'paladin', 'attack', 'enemy', 'holy', 10, 'strength', 0.4, 4, 2, 1, NULL, 1.0),
('holy_light', '圣光术', '恢复生命值。', 'paladin', 'heal', 'ally_lowest_hp', 'holy', 12, 'intellect', 0.5, 5, 0, 2, NULL, 1.0),
('blessing_of_might', '力量祝福', '提升全队攻击力，持续整场战斗。', 'paladin', 'buff', 'ally_all', NULL, 0, NULL, 0, 4, 0, 3, 'eff_blessing_might', 1.0),
('consecration', '奉献', '在脚下创造神圣区域，每回合伤害敌人。', 'paladin', 'dot', 'enemy_all', 'holy', 4, 'strength', 0.2, 6, 3, 6, 'eff_consecration', 1.0),
('divine_shield', '圣盾术', '2回合内免疫所有伤害。', 'paladin', 'buff', 'self', NULL, 0, NULL, 0, 8, 10, 10, 'eff_divine_shield', 1.0),
('lay_on_hands', '圣疗术', '完全恢复目标生命值。', 'paladin', 'heal', 'ally_lowest_hp', 'holy', 100, NULL, 0, 10, 15, 14, NULL, 1.0),
('hammer_of_justice', '制裁之锤', '眩晕敌人2回合。', 'paladin', 'control', 'enemy', 'holy', 5, 'strength', 0.3, 5, 5, 8, 'eff_stun', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 猎人技能 (法力系统: 基础18, 远程物理+宠物)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('arcane_shot', '奥术射击', '发射附魔箭矢。', 'hunter', 'attack', 'enemy', 'magic', 8, 'agility', 0.5, 3, 0, 1, NULL, 1.0),
('serpent_sting', '毒蛇钉刺', '使敌人中毒，持续5回合。', 'hunter', 'dot', 'enemy', 'nature', 2, 'agility', 0.15, 4, 0, 1, 'eff_serpent_sting', 1.0),
('multi_shot', '多重射击', '对所有敌人射击。', 'hunter', 'attack', 'enemy_all', 'physical', 6, 'agility', 0.3, 5, 2, 4, NULL, 1.0),
('aimed_shot', '瞄准射击', '精准射击造成高伤害。', 'hunter', 'attack', 'enemy', 'physical', 15, 'agility', 0.7, 6, 3, 6, NULL, 1.0),
('concussive_shot', '震荡射击', '减速敌人3回合。', 'hunter', 'debuff', 'enemy', 'physical', 4, 'agility', 0.2, 3, 2, 3, 'eff_slow', 1.0),
('rapid_fire', '急速射击', '4回合内攻击速度提升30%。', 'hunter', 'buff', 'self', NULL, 0, NULL, 0, 5, 6, 10, 'eff_rapid_fire', 1.0),
('feign_death', '假死', '脱离战斗1回合，期间不会被攻击。', 'hunter', 'buff', 'self', NULL, 0, NULL, 0, 4, 8, 8, 'eff_feign_death', 1.0),
('bestial_wrath', '狂野怒火', '宠物伤害提升50%，持续3回合。', 'hunter', 'buff', 'self', NULL, 0, NULL, 0, 6, 8, 12, 'eff_bestial_wrath', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 术士技能 (法力系统: 基础38, DOT+吸血+召唤)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('shadow_bolt', '暗影箭', '发射暗影能量。', 'warlock', 'attack', 'enemy', 'shadow', 10, 'intellect', 0.5, 4, 0, 1, NULL, 1.0),
('corruption', '腐蚀术', '使敌人腐蚀，持续6回合。', 'warlock', 'dot', 'enemy', 'shadow', 2, 'intellect', 0.15, 4, 0, 1, 'eff_corruption', 1.0),
('curse_of_agony', '痛苦诅咒', '造成逐渐增强的痛苦，持续5回合。', 'warlock', 'dot', 'enemy', 'shadow', 3, 'intellect', 0.2, 5, 0, 2, 'eff_agony', 1.0),
('drain_life', '吸取生命', '吸取敌人生命，治疗自己。', 'warlock', 'attack', 'enemy', 'shadow', 8, 'intellect', 0.4, 6, 2, 4, 'eff_drain_life', 1.0),
('fear', '恐惧', '使敌人恐惧2回合，无法行动。', 'warlock', 'control', 'enemy', 'shadow', 0, NULL, 0, 6, 5, 6, 'eff_fear', 1.0),
('immolate', '献祭', '燃烧敌人，造成火焰DOT。', 'warlock', 'dot', 'enemy', 'fire', 5, 'intellect', 0.3, 5, 0, 3, 'eff_immolate', 1.0),
('hellfire', '地狱烈焰', '对所有敌人造成火焰伤害，自身也受伤。', 'warlock', 'attack', 'enemy_all', 'fire', 10, 'intellect', 0.4, 8, 4, 10, NULL, 1.0),
('soul_link', '灵魂链接', '与宠物分担30%伤害，持续整场战斗。', 'warlock', 'buff', 'self', NULL, 0, NULL, 0, 6, 0, 8, 'eff_soul_link', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 德鲁伊技能 (法力系统: 基础30, 变形+治疗+DOT)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('wrath', '愤怒', '用自然之力攻击敌人。', 'druid', 'attack', 'enemy', 'nature', 9, 'intellect', 0.5, 4, 0, 1, NULL, 1.0),
('moonfire', '月火术', '造成即时伤害和DOT。', 'druid', 'dot', 'enemy', 'nature', 5, 'intellect', 0.25, 4, 0, 1, 'eff_moonfire', 1.0),
('rejuvenation', '回春术', '持续恢复生命，持续4回合。', 'druid', 'hot', 'ally_lowest_hp', 'nature', 4, 'spirit', 0.25, 5, 0, 2, 'eff_rejuvenation', 1.0),
('regrowth', '愈合', '即时治疗并附带HOT效果。', 'druid', 'heal', 'ally_lowest_hp', 'nature', 10, 'spirit', 0.5, 7, 0, 4, 'eff_regrowth', 1.0),
('entangling_roots', '纠缠根须', '使敌人无法行动2回合。', 'druid', 'control', 'enemy', 'nature', 0, NULL, 0, 5, 4, 6, 'eff_root', 1.0),
('swipe', '横扫', '对所有敌人造成物理伤害。', 'druid', 'attack', 'enemy_all', 'physical', 6, 'strength', 0.3, 5, 2, 5, NULL, 1.0),
('barkskin', '树皮术', '3回合内受到伤害降低25%。', 'druid', 'buff', 'self', NULL, 0, NULL, 0, 4, 5, 8, 'eff_barkskin', 1.0),
('tranquility', '宁静', '持续3回合治疗全队。', 'druid', 'heal', 'ally_all', 'nature', 5, 'spirit', 0.3, 15, 10, 12, 'eff_tranquility', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 萨满技能 (法力系统: 基础32, 元素+图腾+治疗)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('lightning_bolt', '闪电箭', '召唤闪电攻击敌人。', 'shaman', 'attack', 'enemy', 'nature', 10, 'intellect', 0.5, 4, 0, 1, NULL, 1.0),
('earth_shock', '地震术', '造成自然伤害并打断施法。', 'shaman', 'attack', 'enemy', 'nature', 7, 'intellect', 0.35, 4, 2, 1, 'eff_interrupt', 0.8),
('flame_shock', '烈焰震击', '造成火焰伤害和DOT。', 'shaman', 'dot', 'enemy', 'fire', 5, 'intellect', 0.25, 4, 0, 2, 'eff_flame_shock', 1.0),
('healing_wave', '治疗波', '恢复生命值。', 'shaman', 'heal', 'ally_lowest_hp', 'nature', 12, 'spirit', 0.5, 6, 0, 2, NULL, 1.0),
('chain_lightning', '闪电链', '闪电跳跃攻击多个敌人。', 'shaman', 'attack', 'enemy_all', 'nature', 6, 'intellect', 0.3, 7, 3, 6, NULL, 1.0),
('windfury_weapon', '风怒武器', '攻击时有几率额外攻击，持续整场战斗。', 'shaman', 'buff', 'self', NULL, 0, NULL, 0, 5, 0, 4, 'eff_windfury', 1.0),
('purge', '净化', '驱散敌人的增益效果。', 'shaman', 'dispel', 'enemy', NULL, 0, NULL, 0, 4, 2, 8, NULL, 1.0),
('bloodlust', '嗜血', '全队攻击速度提升30%，持续4回合。', 'shaman', 'buff', 'ally_all', NULL, 0, NULL, 0, 10, 12, 12, 'eff_bloodlust', 1.0);

-- ═══════════════════════════════════════════════════════════
-- 通用技能
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance) VALUES
('basic_attack', '普通攻击', '基础物理攻击。', NULL, 'attack', 'enemy', 'physical', 0, 'strength', 0.5, 0, 0, 1, NULL, 1.0);

-- ═══════════════════════════════════════════════════════════
-- 被动技能数据 (每3级技能选择时可能出现)
-- ═══════════════════════════════════════════════════════════
-- tier: 1=基础(Lv3-15), 2=进阶(Lv18-36), 3=大师(Lv39-60)
-- rarity: common=60%出现率, rare=30%, epic=10%

INSERT OR REPLACE INTO passive_skills (id, name, description, class_id, rarity, tier, effect_type, effect_value, effect_stat, max_level, level_scaling) VALUES
-- ═══════════════════════════════════════════════════════════
-- 基础层被动 (Tier 1) - 通用
-- ═══════════════════════════════════════════════════════════
('passive_blade_mastery', '利刃专精', '物理伤害+8%', NULL, 'common', 1, 'stat_mod_pct', 8, 'physical_damage', 5, 0.2),
('passive_armor_mastery', '护甲掌握', '护甲值+15%', NULL, 'common', 1, 'stat_mod_pct', 15, 'armor', 5, 0.2),
('passive_vitality', '生命力', '最大HP+10%', NULL, 'common', 1, 'stat_mod_pct', 10, 'max_hp', 5, 0.2),
('passive_quick_learner', '快速学习', '经验获取+5%', NULL, 'common', 1, 'stat_mod_pct', 5, 'exp_gain', 5, 0.2),
('passive_fortune', '幸运', '金币掉落+10%', NULL, 'common', 1, 'stat_mod_pct', 10, 'gold_gain', 5, 0.2),
('passive_toughness', '坚韧', '受到伤害-5%', NULL, 'common', 1, 'stat_mod_pct', -5, 'damage_taken', 5, 0.2),

-- 基础层被动 - 职业专属
('passive_rage_mastery', '怒气掌控', '怒气获取+15%', 'warrior', 'common', 1, 'stat_mod_pct', 15, 'rage_gain', 5, 0.2),
('passive_mana_flow', '法力涌流', '法力回复+20%', 'mage', 'common', 1, 'stat_mod_pct', 20, 'mana_regen', 5, 0.2),
('passive_energy_flow', '能量循环', '能量恢复+10%', 'rogue', 'common', 1, 'stat_mod_pct', 10, 'energy_regen', 5, 0.2),
('passive_holy_light', '圣光祝福', '治疗效果+10%', 'priest', 'common', 1, 'stat_mod_pct', 10, 'healing_done', 5, 0.2),
('passive_righteousness', '正义之心', '神圣伤害+10%', 'paladin', 'common', 1, 'stat_mod_pct', 10, 'holy_damage', 5, 0.2),
('passive_steady_aim', '稳固射击', '远程攻击+12%', 'hunter', 'common', 1, 'stat_mod_pct', 12, 'ranged_damage', 5, 0.2),
('passive_fel_power', '邪能掌握', '暗影伤害+10%', 'warlock', 'common', 1, 'stat_mod_pct', 10, 'shadow_damage', 5, 0.2),
('passive_nature_bond', '自然契约', '自然伤害+10%', 'druid', 'common', 1, 'stat_mod_pct', 10, 'nature_damage', 5, 0.2),
('passive_elemental_focus', '元素专注', '法术暴击+5%', 'shaman', 'common', 1, 'stat_mod_pct', 5, 'spell_crit', 5, 0.2),

-- ═══════════════════════════════════════════════════════════
-- 进阶层被动 (Tier 2)
-- ═══════════════════════════════════════════════════════════
('passive_critical_edge', '致命一击', '暴击伤害+12%', NULL, 'rare', 2, 'stat_mod_pct', 12, 'crit_damage', 5, 0.2),
('passive_evasion', '闪避本能', '闪避率+5%', NULL, 'rare', 2, 'stat_mod_pct', 5, 'dodge_rate', 5, 0.2),
('passive_precision', '精准打击', '命中率+5%', NULL, 'rare', 2, 'stat_mod_pct', 5, 'hit_rate', 5, 0.2),
('passive_battle_focus', '战斗专注', '暴击率+4%', NULL, 'rare', 2, 'stat_mod_pct', 4, 'crit_rate', 5, 0.2),
('passive_regeneration', '再生', '每回合恢复1%HP', NULL, 'rare', 2, 'regen_pct', 1, 'hp', 5, 0.2),
('passive_mana_shield', '法力护盾', '受到伤害时消耗MP抵消10%', 'mage', 'rare', 2, 'damage_absorb_mp', 10, NULL, 5, 0.2),
('passive_life_steal', '生命汲取', '伤害的3%转化为HP', 'warlock', 'rare', 2, 'lifesteal', 3, NULL, 5, 0.2),
('passive_block_mastery', '格挡专精', '格挡几率+8%', 'warrior', 'rare', 2, 'stat_mod_pct', 8, 'block_rate', 5, 0.2),
('passive_assassin', '刺客本能', '对HP>80%敌人伤害+15%', 'rogue', 'rare', 2, 'conditional_damage', 15, 'enemy_hp_high', 5, 0.2),
('passive_healing_aura', '治愈光环', '队友每回合恢复1HP', 'priest', 'rare', 2, 'team_regen', 1, 'hp', 5, 0.2),
('passive_holy_shield', '神圣护盾', '格挡时反弹3点伤害', 'paladin', 'rare', 2, 'block_reflect', 3, NULL, 5, 0.2),
('passive_pet_bond', '宠物羁绊', '宠物伤害+20%', 'hunter', 'rare', 2, 'stat_mod_pct', 20, 'pet_damage', 5, 0.2),
('passive_soul_siphon', '灵魂虹吸', 'DOT伤害+15%', 'warlock', 'rare', 2, 'stat_mod_pct', 15, 'dot_damage', 5, 0.2),
('passive_natural_regen', '自然再生', '战斗中每回合恢复2%HP', 'druid', 'rare', 2, 'regen_pct', 2, 'hp', 5, 0.2),
('passive_totemic_focus', '图腾专注', '图腾效果+15%', 'shaman', 'rare', 2, 'stat_mod_pct', 15, 'totem_effect', 5, 0.2),

-- ═══════════════════════════════════════════════════════════
-- 大师层被动 (Tier 3)
-- ═══════════════════════════════════════════════════════════
('passive_berserker', '嗜血本能', 'HP<30%时攻击+25%', NULL, 'epic', 3, 'conditional_stat', 25, 'attack_low_hp', 5, 0.2),
('passive_magic_barrier', '魔法屏障', '法术伤害-15%', NULL, 'epic', 3, 'stat_mod_pct', -15, 'magic_damage_taken', 5, 0.2),
('passive_undying', '不灭意志', '首次致死伤害免疫(每战1次)', NULL, 'epic', 3, 'death_prevention', 1, NULL, 3, 0.5),
('passive_executioner', '处刑者', '对HP<20%敌人伤害+50%', NULL, 'epic', 3, 'conditional_damage', 50, 'enemy_hp_low', 5, 0.2),
('passive_iron_will', '钢铁意志', '控制效果持续时间-30%', NULL, 'epic', 3, 'stat_mod_pct', -30, 'cc_duration', 5, 0.2),
('passive_double_strike', '二连击', '普通攻击15%几率攻击两次', 'warrior', 'epic', 3, 'proc_chance', 15, 'double_attack', 5, 0.2),
('passive_spell_echo', '法术回响', '技能15%几率不消耗法力', 'mage', 'epic', 3, 'proc_chance', 15, 'free_cast', 5, 0.2),
('passive_shadow_dance', '暗影之舞', '暴击后下次攻击必暴击', 'rogue', 'epic', 3, 'proc_chain', 100, 'guaranteed_crit', 3, 0.3),
('passive_divine_favor', '神恩术', '治疗技能25%几率暴击', 'paladin', 'epic', 3, 'proc_chance', 25, 'heal_crit', 5, 0.2),
('passive_spirit_bond', '灵魂纽带', '治疗宠物时自己也恢复50%', 'hunter', 'epic', 3, 'heal_split', 50, NULL, 5, 0.2),
('passive_shadow_mastery', '暗影掌控', '暗影技能15%几率重置冷却', 'warlock', 'epic', 3, 'proc_chance', 15, 'reset_cooldown', 5, 0.2),
('passive_omen_of_clarity', '清晰预兆', '攻击15%几率下次技能免费', 'druid', 'epic', 3, 'proc_chance', 15, 'free_cast', 5, 0.2),
('passive_elemental_mastery', '元素掌握', '下次法术必暴击(每10回合)', 'shaman', 'epic', 3, 'cooldown_proc', 100, 'guaranteed_crit', 3, 0.3),
('passive_holy_guardian', '神圣守护', 'HP<20%时自动触发圣盾(每战1次)', 'priest', 'epic', 3, 'death_prevention', 1, NULL, 3, 0.5),
('passive_guardian_angel', '守护天使', '治疗暴击率+15%', 'priest', 'epic', 3, 'stat_mod_pct', 15, 'heal_crit_rate', 5, 0.2);

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


-- ═══════════════════════════════════════════════════════════
-- 装备词缀数据
-- ═══════════════════════════════════════════════════════════

-- 前缀词缀 (攻击/属性向)
INSERT OR REPLACE INTO affixes (id, name, type, slot_type, rarity, effect_type, effect_stat, min_value, max_value, value_type, description, level_required) VALUES
-- 武器前缀 - 普通
('affix_sharp', '锋利的', 'prefix', 'weapon', 'common', 'stat_bonus', 'attack', 2, 5, 'flat', '攻击力 +{value}', 1),
('affix_fiery', '炽热的', 'prefix', 'weapon', 'common', 'elemental_damage', 'fire', 1, 4, 'flat', '火焰伤害 +{value}', 1),
('affix_frozen', '冰霜的', 'prefix', 'weapon', 'common', 'elemental_damage', 'frost', 1, 4, 'flat', '冰霜伤害 +{value}', 1),
('affix_charged', '雷击的', 'prefix', 'weapon', 'common', 'elemental_damage', 'lightning', 1, 4, 'flat', '雷电伤害 +{value}', 1),
-- 武器前缀 - 精良
('affix_holy', '神圣的', 'prefix', 'weapon', 'uncommon', 'elemental_damage', 'holy', 2, 5, 'flat', '神圣伤害 +{value}', 10),
('affix_shadow', '暗影的', 'prefix', 'weapon', 'uncommon', 'elemental_damage', 'shadow', 2, 5, 'flat', '暗影伤害 +{value}', 10),
('affix_brutal', '残暴的', 'prefix', 'weapon', 'uncommon', 'stat_bonus', 'attack', 4, 8, 'flat', '攻击力 +{value}', 15),
-- 武器前缀 - 稀有
('affix_vampiric', '吸血鬼的', 'prefix', 'weapon', 'rare', 'lifesteal', NULL, 2, 5, 'percent', '生命偷取 {value}%', 20),
('affix_berserker', '狂战士的', 'prefix', 'weapon', 'rare', 'stat_bonus_pct', 'attack', 8, 15, 'percent', '攻击力 +{value}%', 25),
-- 武器前缀 - 史诗
('affix_devastating', '毁灭的', 'prefix', 'weapon', 'epic', 'stat_bonus_pct', 'attack', 15, 25, 'percent', '攻击力 +{value}%', 35),
('affix_annihilating', '湮灭的', 'prefix', 'weapon', 'epic', 'stat_bonus_pct', 'damage', 12, 20, 'percent', '伤害 +{value}%', 40),

-- 防具前缀 - 普通
('affix_sturdy', '坚固的', 'prefix', 'armor', 'common', 'stat_bonus', 'defense', 2, 5, 'flat', '防御力 +{value}', 1),
('affix_vital', '活力的', 'prefix', 'armor', 'common', 'stat_bonus', 'max_hp', 5, 15, 'flat', '生命值 +{value}', 1),
('affix_fortified', '强化的', 'prefix', 'armor', 'common', 'stat_bonus', 'armor', 3, 8, 'flat', '护甲 +{value}', 1),
-- 防具前缀 - 精良
('affix_scholarly', '智者的', 'prefix', 'armor', 'uncommon', 'stat_bonus', 'intellect', 2, 4, 'flat', '智力 +{value}', 10),
('affix_resilient', '坚韧的', 'prefix', 'armor', 'uncommon', 'stat_bonus_pct', 'max_hp', 5, 10, 'percent', '生命值 +{value}%', 15),
-- 防具前缀 - 稀有
('affix_unyielding', '不屈的', 'prefix', 'armor', 'rare', 'damage_reduction', NULL, 3, 8, 'percent', '受伤减免 {value}%', 25),
('affix_guardian', '守护的', 'prefix', 'armor', 'rare', 'stat_bonus_pct', 'defense', 10, 18, 'percent', '防御力 +{value}%', 30),
-- 防具前缀 - 史诗
('affix_indomitable', '不可阻挡的', 'prefix', 'armor', 'epic', 'damage_reduction', NULL, 8, 15, 'percent', '受伤减免 {value}%', 40),

-- 后缀词缀 (特殊效果向)
-- 通用后缀 - 普通
('affix_of_strength', 'of 力量', 'suffix', 'all', 'common', 'stat_bonus', 'strength', 1, 3, 'flat', '力量 +{value}', 1),
('affix_of_agility', 'of 敏捷', 'suffix', 'all', 'common', 'stat_bonus', 'agility', 1, 3, 'flat', '敏捷 +{value}', 1),
('affix_of_intellect', 'of 智力', 'suffix', 'all', 'common', 'stat_bonus', 'intellect', 1, 3, 'flat', '智力 +{value}', 1),
('affix_of_stamina', 'of 耐力', 'suffix', 'all', 'common', 'stat_bonus', 'stamina', 1, 3, 'flat', '耐力 +{value}', 1),
('affix_of_spirit', 'of 精神', 'suffix', 'all', 'common', 'stat_bonus', 'spirit', 1, 3, 'flat', '精神 +{value}', 1),

-- 武器后缀 - 精良
('affix_of_haste', 'of 迅捷', 'suffix', 'weapon', 'uncommon', 'stat_bonus_pct', 'attack_speed', 5, 15, 'percent', '攻击速度 +{value}%', 10),
('affix_of_piercing', 'of 穿刺', 'suffix', 'weapon', 'uncommon', 'armor_penetration', NULL, 5, 15, 'percent', '无视 {value}% 护甲', 15),
-- 武器后缀 - 稀有
('affix_of_crit', 'of 暴击', 'suffix', 'weapon', 'rare', 'stat_bonus_pct', 'crit_rate', 3, 8, 'percent', '暴击率 +{value}%', 20),
('affix_of_lethality', 'of 致命', 'suffix', 'weapon', 'rare', 'stat_bonus_pct', 'crit_damage', 10, 25, 'percent', '暴击伤害 +{value}%', 25),
('affix_of_leech', 'of 吸血', 'suffix', 'weapon', 'rare', 'lifesteal', NULL, 2, 4, 'percent', '伤害的 {value}% 转化为生命', 30),
-- 武器后缀 - 史诗
('affix_of_slaying', 'of 斩杀', 'suffix', 'weapon', 'epic', 'execute_damage', NULL, 15, 30, 'percent', '对低血量敌人伤害 +{value}%', 40),
('affix_of_fury', 'of 狂怒', 'suffix', 'weapon', 'epic', 'berserk', NULL, 20, 35, 'percent', 'HP<30%时伤害 +{value}%', 45),

-- 防具后缀 - 精良
('affix_of_blocking', 'of 守护', 'suffix', 'armor', 'uncommon', 'stat_bonus_pct', 'block_rate', 5, 10, 'percent', '格挡率 +{value}%', 10),
('affix_of_evasion', 'of 闪避', 'suffix', 'armor', 'uncommon', 'stat_bonus_pct', 'dodge_rate', 3, 8, 'percent', '闪避率 +{value}%', 15),
-- 防具后缀 - 稀有
('affix_of_thorns', 'of 反射', 'suffix', 'armor', 'rare', 'reflect_damage', NULL, 5, 15, 'percent', '反弹 {value}% 受到的伤害', 25),
('affix_of_regen', 'of 再生', 'suffix', 'armor', 'rare', 'hp_regen', NULL, 1, 3, 'flat', '每回合恢复 {value} HP', 20),
('affix_of_wisdom', 'of 智慧', 'suffix', 'armor', 'rare', 'stat_bonus_pct', 'mana_regen', 10, 20, 'percent', '法力恢复 +{value}%', 20),
-- 防具后缀 - 史诗
('affix_of_immortality', 'of 不朽', 'suffix', 'armor', 'epic', 'hp_regen_pct', NULL, 2, 5, 'percent', '每回合恢复 {value}% HP', 40),
('affix_of_retribution', 'of 惩戒', 'suffix', 'armor', 'epic', 'reflect_damage', NULL, 15, 30, 'percent', '反弹 {value}% 受到的伤害', 45);

-- ═══════════════════════════════════════════════════════════
-- 进化路线数据
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO evolution_paths (id, name, element, description, slot_type, stat_bonus_type, stat_bonus_value, special_effect, material_required) VALUES
-- 武器进化路线
('evo_fire', '烈焰', 'fire', '以火焰之力锻造', 'weapon', 'fire_damage_pct', 50, '攻击有几率灼烧敌人', '{"fire_essence": 5, "evolution_stone": 3}'),
('evo_frost', '霜寒', 'frost', '以冰霜之力锻造', 'weapon', 'frost_damage_pct', 50, '攻击有几率减速敌人', '{"frost_essence": 5, "evolution_stone": 3}'),
('evo_lightning', '雷霆', 'lightning', '以雷电之力锻造', 'weapon', 'lightning_damage_pct', 40, '伤害可连锁跳跃', '{"lightning_essence": 5, "evolution_stone": 3}'),
('evo_holy', '神圣', 'holy', '以圣光之力锻造', 'weapon', 'holy_damage_pct', 40, '攻击时回复生命', '{"holy_essence": 5, "evolution_stone": 3}'),
('evo_shadow', '暗影', 'shadow', '以暗影之力锻造', 'weapon', 'shadow_damage_pct', 50, '伤害转化为生命', '{"shadow_essence": 5, "evolution_stone": 3}'),
('evo_nature', '自然', 'nature', '以自然之力锻造', 'weapon', 'nature_damage_pct', 40, '持续恢复生命', '{"nature_essence": 5, "evolution_stone": 3}'),
('evo_physical', '物理', 'physical', '以纯粹力量锻造', 'weapon', 'physical_damage_pct', 30, '无视部分护甲', '{"steel_ingot": 10, "evolution_stone": 3}'),

-- 防具进化路线
('evo_guardian', '守护', 'guardian', '守护之道', 'armor', 'defense_pct', 30, '受伤减免提升', '{"iron_core": 5, "evolution_stone": 3}'),
('evo_thorns', '荆棘', 'thorns', '荆棘之道', 'armor', 'reflect_pct', 25, '被攻击时反弹伤害', '{"thorn_vine": 5, "evolution_stone": 3}'),
('evo_agile', '迅捷', 'agile', '迅捷之道', 'armor', 'dodge_pct', 20, '闪避后获得加速', '{"swift_feather": 5, "evolution_stone": 3}'),
('evo_vitality', '再生', 'vitality', '生命之道', 'armor', 'hp_regen_pct', 50, '每回合恢复生命', '{"life_crystal": 5, "evolution_stone": 3}');

-- ═══════════════════════════════════════════════════════════
-- 传说效果数据
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO legendary_effects (id, name, description, slot_type, evolution_path, trigger_type, trigger_chance, effect_type, effect_value, cooldown) VALUES
-- 武器传说效果
('legend_inferno', '地狱烈焰', '攻击使敌人灼烧3回合，每回合3点火焰伤害', 'weapon', 'evo_fire', 'on_hit', 1.0, 'apply_dot', 3, 0),
('legend_frostmourne', '霜之哀伤', '击杀敌人后冰冻周围敌人1回合', 'weapon', 'evo_frost', 'on_kill', 1.0, 'aoe_freeze', 1, 0),
('legend_thunderfury', '雷霆之怒', '攻击时20%几率触发闪电链，最多跳跃3个目标', 'weapon', 'evo_lightning', 'on_hit', 0.2, 'chain_lightning', 3, 0),
('legend_ashbringer', '灰烬使者', '攻击时恢复自身5%最大生命值', 'weapon', 'evo_holy', 'on_hit', 1.0, 'heal_pct', 5, 0),
('legend_shadowmourne', '暗影之殇', '暴击时吸取敌人10%当前生命值', 'weapon', 'evo_shadow', 'on_crit', 1.0, 'drain_hp_pct', 10, 0),
('legend_earthshatter', '大地粉碎', '攻击叠加自然标记，5层后引爆造成额外伤害', 'weapon', 'evo_nature', 'on_hit', 1.0, 'stack_explode', 5, 0),
('legend_gorehowl', '血吼', '对精英和Boss伤害+50%', 'weapon', 'evo_physical', 'passive', 1.0, 'elite_damage_pct', 50, 0),

-- 防具传说效果
('legend_immortal', '不灭意志', '首次致死伤害免疫 (每场战斗1次)', 'armor', 'evo_guardian', 'on_fatal', 1.0, 'prevent_death', 1, 0),
('legend_retribution', '复仇之刺', '反弹50%受到的物理伤害', 'armor', 'evo_thorns', 'on_damaged', 1.0, 'reflect_pct', 50, 0),
('legend_shadowstep', '暗影步', '闪避成功后下次攻击必定暴击', 'armor', 'evo_agile', 'on_dodge', 1.0, 'guaranteed_crit', 1, 0),
('legend_lifesource', '生命之泉', 'HP低于30%时每回合恢复10%最大生命', 'armor', 'evo_vitality', 'passive', 1.0, 'low_hp_regen', 10, 0);

-- ═══════════════════════════════════════════════════════════
-- 装备掉落配置
-- ═══════════════════════════════════════════════════════════
-- quality_weights格式: {"common":30,"uncommon":35,"rare":25,"epic":8,"legendary":1.8,"mythic":0.2}

INSERT OR REPLACE INTO drop_config (id, monster_type, base_drop_rate, quality_weights, miracle_rate, pity_threshold, pity_min_quality) VALUES
-- 普通怪物: 5%掉率，正常品质分布
('drop_normal', 'normal', 0.05, 
 '{"common":30,"uncommon":35,"rare":25,"epic":8,"legendary":1.8,"mythic":0.2}',
 0, 40, 'rare'),

-- 精英怪物: 15%掉率，更好的品质
('drop_elite', 'elite', 0.15,
 '{"common":20,"uncommon":30,"rare":30,"epic":15,"legendary":4.5,"mythic":0.5}',
 0, 30, 'rare'),

-- Boss: 50%掉率，优质品质
('drop_boss', 'boss', 0.50,
 '{"common":10,"uncommon":25,"rare":35,"epic":22,"legendary":7,"mythic":1}',
 0, 10, 'epic'),

-- 深渊Boss: 100%掉率，高品质保证
('drop_abyss_boss', 'abyss_boss', 1.00,
 '{"common":0,"uncommon":15,"rare":40,"epic":30,"legendary":12,"mythic":3}',
 0, 1, 'epic');

-- ═══════════════════════════════════════════════════════════
-- 区域奇迹掉落率配置 (按区域等级范围)
-- ═══════════════════════════════════════════════════════════
-- 注：奇迹掉落率在代码中根据区域等级动态计算
-- 1-10级区: 0.5%
-- 11-30级区: 0.3%
-- 31-50级区: 0.1%
-- 51-60级区: 0%

