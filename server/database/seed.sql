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
    base_strength, base_agility, base_intellect, base_stamina, base_spirit,
    base_threat_modifier, combat_role, is_ranged) VALUES
-- 战士: 怒气系统 (初始0，通过攻击/受击获得) - 主力坦克
('warrior', '战士', '近战格斗专家，可以承受大量伤害。', 'tank', 'strength',
    'rage', 35, 0, 3, 0, 0, 0,  -- HP35起步，每级+3
    12, 8, 5, 10, 6,
    1.3, 'tank', 0),  -- 高仇恨，坦克，近战
-- 盗贼: 能量系统 (固定100上限，快速恢复) - 近战DPS低仇恨
('rogue', '盗贼', '潜行刺客，擅长连击和爆发伤害。', 'dps', 'agility',
    'energy', 25, 100, 2, 0, 20, 0,  -- HP25起步，每级+2
    8, 12, 5, 7, 6,
    0.7, 'dps', 0),  -- 低仇恨，DPS，近战
-- 法师: 法力系统 (高法力，基于精神恢复) - 远程DPS
('mage', '法师', '强大的奥术施法者，擅长范围伤害。', 'dps', 'intellect',
    'mana', 20, 40, 2, 2, 0, 0.5,  -- HP20起步，MP40起步
    4, 5, 14, 5, 10,
    0.8, 'dps', 1),  -- 中低仇恨，DPS，远程
-- 牧师: 法力系统 (高法力，高精神恢复) - 主力治疗
('priest', '牧师', '治疗者和暗影施法者。', 'healer', 'intellect',
    'mana', 22, 35, 2, 2, 0, 0.8,  -- HP22起步
    4, 5, 12, 6, 14,
    0.6, 'healer', 1),  -- 最低仇恨，治疗，远程
-- 术士: 法力系统 - 远程DPS（宠物分担仇恨）
('warlock', '术士', '黑暗魔法师，召唤恶魔作战。', 'dps', 'intellect',
    'mana', 24, 38, 2, 2, 0, 0.5,
    5, 5, 13, 6, 9,
    0.8, 'dps', 1),  -- 中低仇恨，DPS，远程
-- 德鲁伊: 法力系统 - 混合职业（熊形态为坦克）
('druid', '德鲁伊', '自然的守护者，可变形为多种形态。', 'hybrid', 'intellect',
    'mana', 28, 30, 2, 1, 0, 0.6,
    8, 8, 10, 8, 10,
    1.0, 'hybrid', 0),  -- 默认中仇恨，熊形态×1.2，近战
-- 萨满: 法力系统 - 混合职业
('shaman', '萨满', '元素的操控者，可治疗和增益。', 'hybrid', 'intellect',
    'mana', 28, 32, 2, 1, 0, 0.6,
    9, 7, 11, 8, 10,
    0.9, 'hybrid', 1),  -- 中仇恨，混合，远程
-- 圣骑士: 法力系统 (较低法力) - 副坦/治疗
('paladin', '圣骑士', '神圣战士，可以治疗和保护盟友。', 'tank', 'strength',
    'mana', 32, 20, 3, 1, 0, 0.4,
    10, 6, 8, 10, 10,
    1.2, 'tank', 0),  -- 高仇恨，坦克，近战
-- 猎人: 法力系统 (较低法力) - 远程DPS（宠物分担仇恨）
('hunter', '猎人', '远程物理攻击者，与宠物并肩作战。', 'dps', 'agility',
    'mana', 26, 18, 2, 1, 0, 0.3,
    6, 12, 6, 8, 8,
    0.8, 'dps', 1);  -- 中低仇恨，DPS，远程

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
('eff_flame_shock', '烈焰震击', '每回合2点火焰伤害', 'dot', 0, 0, 1, 4, 'flat', 2, NULL, 'fire', 1),

-- ═══════════════════════════════════════════════════════════
-- 仇恨管理相关效果
-- ═══════════════════════════════════════════════════════════
-- 仇恨降低效果
('eff_threat_reduce', '仇恨降低', '仇恨值降低50%', 'threat_mod', 1, 0, 1, 1, 'percent', -50, 'threat', NULL, 0),
('eff_fade', '渐隐', '仇恨生成降低50%', 'stat_mod', 1, 0, 1, 3, 'percent', -50, 'threat_gen', NULL, 0),
('eff_misdirection', '误导', '将仇恨转移给坦克', 'threat_transfer', 1, 0, 1, 3, NULL, NULL, NULL, NULL, 0),
('eff_soulshatter', '灵魂碎裂', '仇恨值降低50%', 'threat_mod', 1, 0, 1, 1, 'percent', -50, 'threat', NULL, 0),
-- 仇恨清除效果
('eff_vanish', '消失', '清除所有仇恨并进入潜行', 'stealth', 1, 0, 1, 1, NULL, NULL, NULL, NULL, 0),
('eff_invisibility', '隐形', '逐渐隐形并清除仇恨', 'stealth', 1, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
('eff_ice_block', '寒冰屏障', '免疫伤害并暂停仇恨', 'invulnerable', 1, 0, 1, 2, NULL, NULL, NULL, NULL, 0),
-- 仇恨增加效果 (坦克)
('eff_demo_shout', '挫志', '攻击力降低10%', 'stat_mod', 0, 0, 1, 5, 'percent', -10, 'attack', NULL, 1),
('eff_holy_shield', '神圣之盾', '格挡时造成神圣伤害并产生仇恨', 'stat_mod', 1, 0, 1, 4, 'percent', 30, 'block_rate', 'holy', 0);

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
-- 战士技能 (已迁移到 warrior_skills.sql)
-- ═══════════════════════════════════════════════════════════
-- 注意：所有战士技能数据已迁移到 warrior_skills.sql
-- 请确保在运行seed.sql之前或之后运行warrior_skills.sql

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
-- 仇恨管理技能 (嘲讽/仇恨清除)
-- ═══════════════════════════════════════════════════════════
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target_type, damage_type, base_value, scaling_stat, scaling_ratio, resource_cost, cooldown, level_required, effect_id, effect_chance, threat_modifier, threat_type) VALUES
-- 战士仇恨技能 (已迁移到 warrior_skills.sql)

-- 圣骑士仇恨技能
('righteous_defense', '正义防御', '嘲讽攻击队友的敌人。', 'paladin', 'taunt', 'enemy', NULL, 0, NULL, 0, 4, 4, 4, 'eff_taunt', 1.0, 0, 'taunt'),
('holy_shield', '神圣之盾', '格挡时造成神圣伤害并产生仇恨。', 'paladin', 'buff', 'self', NULL, 0, NULL, 0, 5, 4, 8, 'eff_holy_shield', 1.0, 1.5, 'high'),
('avengers_shield', '复仇者之盾', '投掷盾牌造成伤害并产生高仇恨。', 'paladin', 'attack', 'enemy', 'holy', 12, 'strength', 0.5, 6, 3, 10, NULL, 1.0, 2.0, 'high'),

-- 德鲁伊(熊形态)仇恨技能
('growl', '低吼', '嘲讽单个敌人。', 'druid', 'taunt', 'enemy', NULL, 0, NULL, 0, 5, 4, 6, 'eff_taunt', 1.0, 0, 'taunt'),
('swipe', '挥击', '对所有敌人造成伤害并产生仇恨。', 'druid', 'attack', 'enemy_all', 'physical', 5, 'strength', 0.25, 6, 0, 8, NULL, 1.0, 1.5, 'high'),
('maul', '槌击', '强力攻击产生高仇恨。', 'druid', 'attack', 'enemy', 'physical', 10, 'strength', 0.5, 8, 0, 6, NULL, 1.0, 1.8, 'high'),

-- 盗贼仇恨清除技能
('feint', '佯攻', '清除50%仇恨。', 'rogue', 'threat', 'self', NULL, 0, NULL, 0, 20, 4, 8, 'eff_threat_reduce', 1.0, 0, 'reduce'),
('vanish', '消失', '进入潜行并清除所有仇恨。', 'rogue', 'threat', 'self', NULL, 0, NULL, 0, 50, 10, 12, 'eff_vanish', 1.0, 0, 'clear'),

-- 法师仇恨清除技能
('invisibility', '隐形术', '逐渐隐形，清除所有仇恨。', 'mage', 'threat', 'self', NULL, 0, NULL, 0, 15, 10, 14, 'eff_invisibility', 1.0, 0, 'clear'),
('ice_block', '寒冰屏障', '免疫所有伤害并暂停仇恨生成。', 'mage', 'buff', 'self', NULL, 0, NULL, 0, 15, 8, 10, 'eff_ice_block', 1.0, 0, 'reduce'),

-- 牧师仇恨管理技能
('fade', '渐隐术', '3回合内仇恨生成降低50%。', 'priest', 'threat', 'self', NULL, 0, NULL, 0, 5, 6, 6, 'eff_fade', 1.0, 0, 'reduce'),

-- 猎人仇恨管理技能
('misdirection', '误导', '将仇恨转移给坦克，持续3回合。', 'hunter', 'threat', 'ally_tank', NULL, 0, NULL, 0, 5, 6, 10, 'eff_misdirection', 1.0, 0, 'reduce'),

-- 术士仇恨管理技能
('soulshatter', '灵魂碎裂', '清除50%仇恨。', 'warlock', 'threat', 'self', NULL, 0, NULL, 0, 6, 8, 10, 'eff_soulshatter', 1.0, 0, 'reduce');

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

INSERT OR REPLACE INTO zones (id, name, description, min_level, max_level, faction, exp_modifier, gold_modifier, unlock_zone_id, required_exploration) VALUES
-- ═══════════════════════════════════════════════════════════
-- 联盟区域
-- ═══════════════════════════════════════════════════════════
-- 初始地图（无需解锁）
('elwynn', '艾尔文森林', '人类王国暴风城外的宁静森林，适合新手冒险者。', 1, 10, 'alliance', 1.0, 1.0, NULL, 0),
('dun_morogh', '丹莫罗', '矮人和侏儒的雪域家园，群山环绕的寒冷之地。', 1, 10, 'alliance', 1.0, 1.0, NULL, 0),
('teldrassil', '泰达希尔', '暗夜精灵的世界树，神秘的森林与自然之力。', 1, 10, 'alliance', 1.0, 1.0, NULL, 0),
-- 10-20级地图（需要初始地图探索度50-80）
('westfall', '西部荒野', '曾经肥沃的农田，如今被迪菲亚兄弟会占领。', 10, 20, 'alliance', 1.1, 1.1, 'elwynn', 50),
('loch_modan', '洛克莫丹', '矮人的家园，群山环绕的美丽山谷。', 10, 22, 'alliance', 1.1, 1.1, 'dun_morogh', 50),
('darkshore', '黑海岸', '被诅咒的海岸线，暗夜精灵的领地。', 10, 20, 'alliance', 1.1, 1.1, 'teldrassil', 50),
-- 15-25级地图（需要初始地图探索度80-120）
('redridge', '赤脊山', '被黑石兽人威胁的山脉，联盟的前线。', 15, 25, 'alliance', 1.15, 1.15, 'elwynn', 80),
-- 18-30级地图（需要10-20级地图探索度100-150）
('wetlands', '湿地', '泥泞的沼泽地，连接南北的交通要道。', 18, 30, 'alliance', 1.2, 1.2, 'loch_modan', 100),
('hillsbrad', '希尔斯布莱德丘陵', '联盟的农业区，经常遭受部落袭击。', 20, 30, 'alliance', 1.2, 1.2, 'westfall', 100),
-- 30-40级地图（需要18-30级地图探索度150-200）
('arathi', '阿拉希高地', '联盟与部落争夺的战略要地。', 30, 40, 'alliance', 1.3, 1.3, 'hillsbrad', 150),
-- ═══════════════════════════════════════════════════════════
-- 部落区域
-- ═══════════════════════════════════════════════════════════
-- 初始地图（无需解锁）
('durotar', '杜隆塔尔', '兽人的家园，炎热干燥的红色大地。', 1, 10, 'horde', 1.0, 1.0, NULL, 0),
('mulgore', '莫高雷', '牛头人的家园，广阔的草原。', 1, 10, 'horde', 1.0, 1.0, NULL, 0),
('tirisfal', '提瑞斯法林地', '被遗忘者的家园，阴森的森林。', 1, 10, 'horde', 1.0, 1.0, NULL, 0),
-- 10-25级地图（需要初始地图探索度50-80）
('barrens', '贫瘠之地', '广袤的草原，危险与机遇并存。', 10, 25, 'horde', 1.1, 1.1, 'durotar', 50),
('silverpine', '银松森林', '被诅咒的森林，被遗忘者的领地。', 10, 22, 'horde', 1.1, 1.1, 'tirisfal', 50),
-- 15-25级地图（需要初始地图探索度80-120）
('stonetalon', '石爪山脉', '被联盟威胁的山脉，部落的防御前线。', 15, 25, 'horde', 1.15, 1.15, 'mulgore', 80),
-- 18-30级地图（需要10-25级地图探索度100-150）
('ashenvale', '灰谷', '暗夜精灵与部落的冲突前线。', 18, 30, 'horde', 1.2, 1.2, 'barrens', 100),
('tarren_mill', '塔伦米尔', '部落的据点，与联盟的希尔斯布莱德对峙。', 20, 30, 'horde', 1.2, 1.2, 'silverpine', 100),
-- 25-35级地图（需要15-25级地图探索度150-200）
('thousand_needles', '千针石林', '奇特的石柱群，半人马的家园。', 25, 35, 'horde', 1.25, 1.25, 'stonetalon', 150),
-- 30-40级地图（需要18-30级地图探索度200-250）
('desolace', '凄凉之地', '荒凉的废土，半人马和恶魔的领地。', 30, 40, 'horde', 1.3, 1.3, 'ashenvale', 200),
-- ═══════════════════════════════════════════════════════════
-- PVP/中立区域
-- ═══════════════════════════════════════════════════════════
-- 18-30级PVP地图（需要联盟或部落10-20级地图探索度100）
('duskwood', '暮色森林', '被永恒黑暗笼罩的诡异森林，亡灵与狼人出没。', 18, 30, NULL, 1.2, 1.2, 'westfall', 100),
-- 28-45级PVP地图（需要18-30级地图探索度200）
('stranglethorn', '荆棘谷', '危险的丛林，到处是食人族和野兽。', 28, 45, NULL, 1.3, 1.3, 'duskwood', 200),
-- 32-45级PVP地图（需要28-45级地图探索度250）
('badlands', '荒芜之地', '荒凉的废土，黑铁矮人的家园。', 32, 45, NULL, 1.35, 1.35, 'stranglethorn', 250),
('swamp_of_sorrows', '悲伤沼泽', '被诅咒的沼泽，绿龙军团的领地。', 32, 45, NULL, 1.35, 1.35, 'stranglethorn', 250),
('dustwallow', '尘泥沼泽', '泥泞的沼泽，黑龙的巢穴。', 32, 45, NULL, 1.35, 1.35, 'stranglethorn', 250),
-- 38-50级PVP地图（需要32-45级地图探索度300）
('tanaris', '塔纳利斯', '炎热的沙漠，隐藏着古老的秘密。', 38, 50, NULL, 1.4, 1.4, 'badlands', 300),
('feralas', '菲拉斯', '茂密的丛林，古精灵的遗迹。', 38, 50, NULL, 1.4, 1.4, 'tanaris', 300),
-- 45-55级PVP地图（需要38-50级地图探索度400）
('ungoro', '安戈洛环形山', '史前生物的乐园，充满危险。', 45, 55, NULL, 1.45, 1.45, 'tanaris', 400),
('felwood', '费伍德森林', '被恶魔腐蚀的森林，燃烧军团的痕迹。', 45, 55, NULL, 1.45, 1.45, 'feralas', 400),
-- 48-60级PVP地图（需要45-55级地图探索度500）
('burning_steppes', '燃烧平原', '被黑龙军团占领的焦土。', 48, 60, NULL, 1.5, 1.5, 'ungoro', 500),
-- 52-60级PVP地图（需要48-60级地图探索度600）
('winterspring', '冬泉谷', '永恒的雪域，蓝龙军团的领地。', 52, 60, NULL, 1.5, 1.5, 'burning_steppes', 600),
('silithus', '希利苏斯', '沙漠中的虫巢，其拉虫人的威胁。', 52, 60, NULL, 1.5, 1.5, 'burning_steppes', 600);

-- ═══════════════════════════════════════════════════════════
-- 怪物数据
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 怪物数据 (小数值设计：HP 15~300, 攻击3~50, 经验5~80)
-- ═══════════════════════════════════════════════════════════

-- 艾尔文森林 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('wolf', 'elwynn', '森林狼', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('young_boar', 'elwynn', '小野猪', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('kobold_worker', 'elwynn', '狗头人矿工', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('kobold_tunneler', 'elwynn', '狗头人掘地工', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('defias_thug', 'elwynn', '迪菲亚暴徒', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('defias_bandit', 'elwynn', '迪菲亚劫匪', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('murloc', 'elwynn', '鱼人', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('murloc_warrior', 'elwynn', '鱼人战士', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('kobold_geomancer', 'elwynn', '狗头人地卜师', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('prowler', 'elwynn', '潜伏者', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('hogger', 'elwynn', '霍格', 8, 'elite', 120, 17, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 35, 5, 12, 5);

-- 丹莫罗 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('frost_wolf', 'dun_morogh', '霜狼', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('snow_boar', 'dun_morogh', '雪野猪', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('trogg_worker', 'dun_morogh', '穴居人苦工', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('ice_troll', 'dun_morogh', '冰霜巨魔', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('frostmane_scout', 'dun_morogh', '霜鬃斥候', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('frostmane_warrior', 'dun_morogh', '霜鬃战士', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('crag_boar', 'dun_morogh', '峭壁野猪', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('young_wendigo', 'dun_morogh', '幼年温迪戈', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('frostmane_shaman', 'dun_morogh', '霜鬃萨满', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('wendigo', 'dun_morogh', '温迪戈', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('great_father_arctica', 'dun_morogh', '大熊阿卡提卡', 8, 'elite', 80, 12, 0, 5, 5, 'physical', 0.08, 1.6, 0.05, 1.5, 35, 5, 12, 5);

-- 泰达希尔 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('young_nightsaber', 'teldrassil', '幼年夜刃豹', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('webwood_lurker', 'teldrassil', '蛛网潜伏者', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('greymane_cub', 'teldrassil', '灰鬃幼崽', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('moonkin', 'teldrassil', '枭兽', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('furbolg_scout', 'teldrassil', '熊怪斥候', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('furbolg_warrior', 'teldrassil', '熊怪战士', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('nightsaber', 'teldrassil', '夜刃豹', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('greymane', 'teldrassil', '灰鬃', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('furbolg_shaman', 'teldrassil', '熊怪萨满', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('webwood_spider', 'teldrassil', '蛛网蜘蛛', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('ursal_the_mauler', 'teldrassil', '猛击者乌萨尔', 8, 'elite', 80, 12, 0, 5, 5, 'physical', 0.08, 1.6, 0.05, 1.5, 35, 5, 12, 5);

-- 西部荒野 (10-20级) - HP: 35~80, 攻击: 9~18
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('harvest_golem', 'westfall', '收割傀儡', 10, 'normal', 60, 14, 0, 7, 7, 'physical', 0.08, 1.5, 0.07, 1.5, 14, 2, 5, 100),
('defias_rogue', 'westfall', '迪菲亚盗贼', 11, 'normal', 64, 15, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 16, 3, 6, 80),
('defias_highwayman', 'westfall', '迪菲亚拦路贼', 12, 'normal', 70, 16, 0, 8, 8, 'physical', 0.10, 1.6, 0.07, 1.5, 18, 3, 7, 60),
('gnoll_brute', 'westfall', '豺狼人蛮兵', 13, 'normal', 75, 18, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 20, 4, 8, 50),
('gnoll_mystic', 'westfall', '豺狼人秘法师', 14, 'normal', 65, 11, 20, 7, 11, 'magic', 0.07, 1.5, 0.12, 1.6, 22, 4, 9, 40),
('defias_pyromancer', 'westfall', '迪菲亚纵火者', 15, 'normal', 72, 12, 24, 8, 11, 'magic', 0.07, 1.5, 0.14, 1.65, 26, 6, 12, 25),
('defias_overlord', 'westfall', '迪菲亚霸主', 16, 'elite', 174, 27, 0, 14, 14, 'physical', 0.14, 1.7, 0.07, 1.5, 50, 10, 20, 5),
('dust_devil', 'westfall', '尘土恶魔', 15, 'normal', 55, 14, 0, 6, 6, 'physical', 0.07, 1.5, 0.05, 1.5, 24, 5, 10, 30);

-- 洛克莫丹 (10-20级) - HP: 35~80, 攻击: 9~18
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('mountain_boar', 'loch_modan', '山猪', 10, 'normal', 60, 14, 0, 7, 7, 'physical', 0.08, 1.5, 0.07, 1.5, 14, 2, 5, 100),
('trogg_brute', 'loch_modan', '穴居人蛮兵', 11, 'normal', 64, 15, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 16, 3, 6, 80),
('dark_iron_dwarf', 'loch_modan', '黑铁矮人', 12, 'normal', 70, 16, 0, 8, 8, 'physical', 0.10, 1.6, 0.07, 1.5, 18, 3, 7, 60),
('mountain_cougar', 'loch_modan', '山猫', 13, 'normal', 75, 18, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 20, 4, 8, 50),
('dark_iron_sorcerer', 'loch_modan', '黑铁术士', 14, 'normal', 65, 11, 20, 7, 11, 'magic', 0.07, 1.5, 0.12, 1.6, 22, 4, 9, 40),
('stone_elemental', 'loch_modan', '石元素', 15, 'normal', 72, 12, 24, 8, 11, 'magic', 0.07, 1.5, 0.14, 1.65, 26, 6, 12, 25),
('dark_iron_commander', 'loch_modan', '黑铁指挥官', 16, 'elite', 174, 27, 0, 14, 14, 'physical', 0.14, 1.7, 0.07, 1.5, 50, 10, 20, 5),
('elder_mountain_boar', 'loch_modan', '老山猪', 15, 'normal', 55, 14, 0, 6, 6, 'physical', 0.07, 1.5, 0.05, 1.5, 24, 5, 10, 30);

-- 黑海岸 (10-20级) - HP: 35~80, 攻击: 9~18
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('darkshore_thresher', 'darkshore', '黑海岸鞭尾鱼', 10, 'normal', 60, 14, 0, 7, 7, 'physical', 0.08, 1.5, 0.07, 1.5, 14, 2, 5, 100),
('greymist_warrior', 'darkshore', '灰雾战士', 11, 'normal', 64, 15, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 16, 3, 6, 80),
('moonstalker', 'darkshore', '月夜豹', 12, 'normal', 70, 16, 0, 8, 8, 'physical', 0.10, 1.6, 0.07, 1.5, 18, 3, 7, 60),
('vile_sprite', 'darkshore', '邪恶精灵', 13, 'normal', 75, 18, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 20, 4, 8, 50),
('darkstrider', 'darkshore', '暗行者', 14, 'normal', 65, 11, 20, 7, 11, 'magic', 0.07, 1.5, 0.12, 1.6, 22, 4, 9, 40),
('naga_myrmidon', 'darkshore', '纳迦战士', 15, 'normal', 72, 12, 24, 8, 11, 'magic', 0.07, 1.5, 0.14, 1.65, 26, 6, 12, 25),
('naga_siren', 'darkshore', '纳迦海妖', 16, 'elite', 174, 27, 0, 14, 14, 'physical', 0.14, 1.7, 0.07, 1.5, 50, 10, 20, 5),
('ancient_of_war', 'darkshore', '战争古树', 15, 'normal', 55, 14, 0, 6, 6, 'physical', 0.07, 1.5, 0.05, 1.5, 24, 5, 10, 30);

-- 暮色森林 (20-30级) - HP: 55~150, 攻击: 16~35
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('skeleton_warrior', 'duskwood', '骷髅战士', 20, 'normal', 94, 24, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 5, 10, 100),
('skeleton_mage', 'duskwood', '骷髅法师', 21, 'normal', 77, 11, 30, 9, 16, 'magic', 0.07, 1.5, 0.14, 1.65, 30, 5, 11, 80),
('ghoul', 'duskwood', '食尸鬼', 22, 'normal', 98, 27, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 6, 12, 70),
('dire_wolf', 'duskwood', '恐狼', 23, 'normal', 105, 30, 0, 14, 14, 'physical', 0.11, 1.6, 0.07, 1.5, 34, 6, 13, 60),
('worgen', 'duskwood', '狼人', 24, 'normal', 119, 34, 0, 15, 15, 'physical', 0.12, 1.6, 0.07, 1.5, 36, 7, 14, 50),
('shadow_cultist', 'duskwood', '暗影邪教徒', 25, 'normal', 126, 16, 38, 15, 19, 'magic', 0.08, 1.5, 0.16, 1.7, 40, 8, 18, 35),
('worgen_alpha', 'duskwood', '狼人首领', 26, 'elite', 252, 47, 0, 20, 20, 'physical', 0.16, 1.7, 0.07, 1.5, 65, 12, 25, 5),
('abomination', 'duskwood', '憎恶', 28, 'elite', 308, 54, 0, 24, 24, 'physical', 0.17, 1.7, 0.07, 1.5, 75, 15, 30, 3),
('stitches', 'duskwood', '缝合怪', 30, 'boss', 490, 68, 0, 30, 30, 'physical', 0.20, 1.8, 0.08, 1.5, 100, 25, 50, 1);

-- 莫高雷 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('plainstrider', 'mulgore', '平原陆行鸟', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('prairie_wolf', 'mulgore', '草原狼', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('windfury_harpy', 'mulgore', '风怒鹰身人', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('bristleback_quillboar', 'mulgore', '刺背野猪人', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('bristleback_thornweaver', 'mulgore', '刺背织棘者', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('bristleback_battleboar', 'mulgore', '刺背战猪', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('adult_plainstrider', 'mulgore', '成年平原陆行鸟', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('elder_plainstrider', 'mulgore', '老平原陆行鸟', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('bristleback_shaman', 'mulgore', '刺背萨满', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('windfury_matron', 'mulgore', '风怒主母', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('chief_bloodhoof', 'mulgore', '血蹄酋长', 8, 'elite', 80, 12, 0, 5, 5, 'physical', 0.08, 1.6, 0.05, 1.5, 35, 5, 12, 5);

-- 提瑞斯法林地 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('rotting_dead', 'tirisfal', '腐烂的尸体', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('vile_fang', 'tirisfal', '邪恶之牙', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('scavenger', 'tirisfal', '食腐者', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('darkhound', 'tirisfal', '暗影猎犬', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('skeletal_warrior', 'tirisfal', '骷髅战士', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('skeletal_mage', 'tirisfal', '骷髅法师', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('duskbat', 'tirisfal', '暮色蝙蝠', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('grave_robber', 'tirisfal', '盗墓者', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('necrotic_shade', 'tirisfal', '死灵阴影', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('cursed_undead', 'tirisfal', '被诅咒的不死生物', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('captain_perolde', 'tirisfal', '佩罗尔德船长', 8, 'elite', 80, 12, 0, 5, 5, 'physical', 0.08, 1.6, 0.05, 1.5, 35, 5, 12, 5);

-- 杜隆塔尔 (1-10级) - HP: 15~45, 攻击: 3~10
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('scorpid_durotar', 'durotar', '蝎子', 1, 'normal', 22, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 5, 1, 2, 100),
('dire_wolf_durotar', 'durotar', '恐狼', 1, 'normal', 18, 4, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 4, 1, 1, 100),
('razormane_scout_durotar', 'durotar', '钢鬃斥候', 2, 'normal', 27, 6, 0, 1, 1, 'physical', 0.07, 1.5, 0.07, 1.5, 6, 1, 2, 80),
('razormane_warrior_durotar', 'durotar', '钢鬃战士', 3, 'normal', 33, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 60),
('razormane_thornweaver_durotar', 'durotar', '钢鬃织棘者', 4, 'normal', 39, 8, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 8, 2, 3, 50),
('razormane_battleboar_durotar', 'durotar', '钢鬃战猪', 5, 'normal', 45, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 10, 2, 4, 40),
('venomtail_scorpid_durotar', 'durotar', '毒尾蝎', 3, 'normal', 30, 7, 0, 3, 3, 'physical', 0.07, 1.5, 0.07, 1.5, 7, 1, 3, 70),
('elder_scorpid_durotar', 'durotar', '老蝎子', 5, 'normal', 42, 10, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 9, 2, 4, 30),
('razormane_geomancer_durotar', 'durotar', '钢鬃地卜师', 6, 'normal', 45, 6, 14, 3, 4, 'magic', 0.07, 1.5, 0.10, 1.6, 12, 2, 5, 35),
('razormane_champion_durotar', 'durotar', '钢鬃勇士', 6, 'normal', 48, 11, 0, 4, 4, 'physical', 0.07, 1.5, 0.07, 1.5, 11, 2, 5, 25),
('captain_flat_tusk', 'durotar', '平牙队长', 8, 'elite', 120, 17, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 35, 5, 12, 5);

-- 银松森林 (10-20级) - HP: 35~80, 攻击: 9~18
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('rotting_dead', 'silverpine', '腐烂的尸体', 10, 'normal', 60, 14, 0, 7, 7, 'physical', 0.08, 1.5, 0.07, 1.5, 14, 2, 5, 100),
('vile_fang', 'silverpine', '邪恶之牙', 11, 'normal', 64, 15, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 16, 3, 6, 80),
('shadowfang_worgen', 'silverpine', '影牙狼人', 12, 'normal', 70, 16, 0, 8, 8, 'physical', 0.10, 1.6, 0.07, 1.5, 18, 3, 7, 60),
('skeletal_warrior', 'silverpine', '骷髅战士', 13, 'normal', 75, 18, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 20, 4, 8, 50),
('skeletal_mage', 'silverpine', '骷髅法师', 14, 'normal', 65, 11, 20, 7, 11, 'magic', 0.07, 1.5, 0.12, 1.6, 22, 4, 9, 40),
('shadowfang_cultist', 'silverpine', '影牙邪教徒', 15, 'normal', 72, 12, 24, 8, 11, 'magic', 0.07, 1.5, 0.14, 1.65, 26, 6, 12, 25),
('archmage_aratus', 'silverpine', '大法师阿拉图斯', 16, 'elite', 174, 27, 0, 14, 14, 'physical', 0.14, 1.7, 0.07, 1.5, 50, 10, 20, 5),
('cursed_undead', 'silverpine', '被诅咒的不死生物', 15, 'normal', 55, 14, 0, 6, 6, 'physical', 0.07, 1.5, 0.05, 1.5, 24, 5, 10, 30);

-- 贫瘠之地 (10-25级) - HP: 35~90, 攻击: 9~20
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('plainstrider', 'barrens', '平原陆行鸟', 10, 'normal', 60, 14, 0, 7, 7, 'physical', 0.08, 1.5, 0.07, 1.5, 14, 2, 5, 100),
('zhevra', 'barrens', '斑马', 11, 'normal', 64, 15, 0, 7, 7, 'physical', 0.10, 1.6, 0.07, 1.5, 16, 3, 6, 80),
('bristleback_quillboar', 'barrens', '刺背野猪人', 12, 'normal', 70, 16, 0, 8, 8, 'physical', 0.10, 1.6, 0.07, 1.5, 18, 3, 7, 60),
('bristleback_thornweaver', 'barrens', '刺背织棘者', 13, 'normal', 75, 18, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 20, 4, 8, 50),
('bristleback_shaman', 'barrens', '刺背萨满', 14, 'normal', 65, 11, 20, 7, 11, 'magic', 0.07, 1.5, 0.12, 1.6, 22, 4, 9, 40),
('kolkar_brute', 'barrens', '科卡尔蛮兵', 15, 'normal', 72, 12, 24, 8, 11, 'magic', 0.07, 1.5, 0.14, 1.65, 26, 6, 12, 25),
('kolkar_chieftain', 'barrens', '科卡尔酋长', 18, 'elite', 203, 32, 0, 16, 16, 'physical', 0.15, 1.7, 0.07, 1.5, 55, 12, 22, 5),
('razormane_warrior', 'barrens', '钢鬃战士', 16, 'normal', 80, 19, 0, 8, 8, 'physical', 0.09, 1.5, 0.07, 1.5, 24, 5, 10, 30),
('razormane_geomancer', 'barrens', '钢鬃地卜师', 17, 'normal', 87, 14, 27, 9, 12, 'magic', 0.07, 1.5, 0.15, 1.65, 28, 6, 13, 25),
('razormane_champion', 'barrens', '钢鬃勇士', 20, 'elite', 160, 28, 0, 13, 13, 'physical', 0.14, 1.7, 0.05, 1.5, 60, 14, 25, 3);

-- 赤脊山 (15-25级) - HP: 45~95, 攻击: 11~22
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('blackrock_orc', 'redridge', '黑石兽人', 15, 'normal', 72, 15, 0, 9, 9, 'physical', 0.09, 1.5, 0.07, 1.5, 26, 6, 12, 100),
('blackrock_grunt', 'redridge', '黑石步兵', 16, 'normal', 78, 16, 0, 9, 9, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 6, 13, 80),
('blackrock_warrior', 'redridge', '黑石战士', 17, 'normal', 84, 18, 0, 11, 11, 'physical', 0.10, 1.6, 0.07, 1.5, 30, 7, 14, 60),
('blackrock_shaman', 'redridge', '黑石萨满', 18, 'normal', 75, 12, 22, 9, 12, 'magic', 0.07, 1.5, 0.13, 1.65, 32, 7, 15, 50),
('blackrock_warlock', 'redridge', '黑石术士', 19, 'normal', 81, 14, 26, 11, 14, 'magic', 0.07, 1.5, 0.15, 1.7, 34, 8, 16, 40),
('blackrock_ogre', 'redridge', '黑石食人魔', 20, 'normal', 102, 24, 0, 12, 12, 'physical', 0.11, 1.6, 0.07, 1.5, 36, 8, 17, 30),
('blackrock_commander', 'redridge', '黑石指挥官', 22, 'elite', 210, 40, 0, 19, 19, 'physical', 0.15, 1.7, 0.07, 1.5, 60, 12, 24, 5),
('redridge_basilisk', 'redridge', '赤脊山蜥蜴', 21, 'normal', 65, 16, 0, 8, 8, 'physical', 0.08, 1.6, 0.05, 1.5, 38, 8, 18, 25);

-- 湿地 (20-30级) - HP: 55~120, 攻击: 16~28
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('mire_creeper', 'wetlands', '泥沼爬行者', 20, 'normal', 94, 24, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 5, 10, 100),
('dragonmaw_orc', 'wetlands', '龙喉兽人', 21, 'normal', 98, 26, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 30, 5, 11, 80),
('dragonmaw_warrior', 'wetlands', '龙喉战士', 22, 'normal', 105, 27, 0, 14, 14, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 6, 12, 70),
('dragonmaw_shaman', 'wetlands', '龙喉萨满', 23, 'normal', 95, 15, 32, 12, 16, 'magic', 0.07, 1.5, 0.15, 1.7, 34, 6, 13, 60),
('dragonmaw_warlock', 'wetlands', '龙喉术士', 24, 'normal', 101, 16, 35, 14, 18, 'magic', 0.07, 1.5, 0.16, 1.7, 36, 7, 14, 50),
('green_dragon_whelp', 'wetlands', '绿龙幼崽', 25, 'normal', 112, 30, 0, 15, 15, 'physical', 0.12, 1.6, 0.07, 1.5, 38, 7, 15, 40),
('dragonmaw_commander', 'wetlands', '龙喉指挥官', 26, 'elite', 252, 47, 0, 20, 20, 'physical', 0.16, 1.7, 0.07, 1.5, 65, 12, 25, 5),
('mire_lurker', 'wetlands', '泥沼潜伏者', 24, 'normal', 78, 21, 0, 10, 10, 'physical', 0.09, 1.6, 0.05, 1.5, 36, 7, 14, 35);

-- 希尔斯布莱德丘陵 (20-30级) - HP: 55~120, 攻击: 16~28
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('hillsbrad_farmer', 'hillsbrad', '希尔斯布莱德农夫', 20, 'normal', 94, 24, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 5, 10, 100),
('syndicate_thug', 'hillsbrad', '辛迪加暴徒', 21, 'normal', 98, 26, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 30, 5, 11, 80),
('syndicate_rogue', 'hillsbrad', '辛迪加盗贼', 22, 'normal', 105, 27, 0, 14, 14, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 6, 12, 70),
('syndicate_mage', 'hillsbrad', '辛迪加法师', 23, 'normal', 95, 15, 32, 12, 16, 'magic', 0.07, 1.5, 0.15, 1.7, 34, 6, 13, 60),
('syndicate_assassin', 'hillsbrad', '辛迪加刺客', 24, 'normal', 101, 16, 35, 14, 18, 'magic', 0.07, 1.5, 0.16, 1.7, 36, 7, 14, 50),
('hillsbrad_peasant', 'hillsbrad', '希尔斯布莱德农民', 25, 'normal', 112, 30, 0, 15, 15, 'physical', 0.12, 1.6, 0.07, 1.5, 38, 7, 15, 40),
('syndicate_master', 'hillsbrad', '辛迪加首领', 26, 'elite', 252, 47, 0, 20, 20, 'physical', 0.16, 1.7, 0.07, 1.5, 65, 12, 25, 5),
('hillsbrad_guard', 'hillsbrad', '希尔斯布莱德守卫', 24, 'normal', 78, 21, 0, 10, 10, 'physical', 0.09, 1.6, 0.05, 1.5, 36, 7, 14, 35);

-- 阿拉希高地 (30-40级) - HP: 75~140, 攻击: 22~35
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('ogre_brute', 'arathi', '食人魔蛮兵', 30, 'normal', 133, 38, 0, 16, 16, 'physical', 0.12, 1.6, 0.07, 1.5, 42, 8, 16, 100),
('ogre_mage', 'arathi', '食人魔法师', 31, 'normal', 123, 14, 40, 15, 19, 'magic', 0.08, 1.5, 0.18, 1.7, 44, 8, 17, 80),
('ogre_warrior', 'arathi', '食人魔战士', 32, 'normal', 140, 40, 0, 18, 18, 'physical', 0.14, 1.6, 0.08, 1.5, 46, 9, 18, 70),
('ogre_lord', 'arathi', '食人魔领主', 33, 'elite', 280, 51, 0, 22, 22, 'physical', 0.18, 1.7, 0.08, 1.5, 70, 14, 28, 5),
('witherbark_troll', 'arathi', '枯木巨魔', 34, 'normal', 147, 43, 0, 19, 19, 'physical', 0.15, 1.6, 0.08, 1.5, 48, 9, 19, 60),
('witherbark_shaman', 'arathi', '枯木萨满', 35, 'normal', 137, 16, 43, 18, 20, 'magic', 0.08, 1.5, 0.19, 1.7, 50, 10, 20, 50),
('witherbark_chieftain', 'arathi', '枯木酋长', 36, 'elite', 308, 54, 0, 23, 23, 'physical', 0.19, 1.7, 0.08, 1.5, 75, 15, 30, 3),
('hammerfall_guard', 'arathi', '落锤守卫', 34, 'normal', 103, 31, 0, 13, 13, 'physical', 0.11, 1.6, 0.05, 1.5, 48, 9, 19, 40);

-- 石爪山脉 (15-25级) - HP: 45~95, 攻击: 11~22
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('windfury_harpy', 'stonetalon', '风怒鹰身人', 15, 'normal', 72, 15, 0, 9, 9, 'physical', 0.09, 1.5, 0.07, 1.5, 26, 6, 12, 100),
('windfury_rogue', 'stonetalon', '风怒盗贼', 16, 'normal', 78, 16, 0, 9, 9, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 6, 13, 80),
('windfury_witch', 'stonetalon', '风怒女巫', 17, 'normal', 70, 11, 22, 8, 12, 'magic', 0.07, 1.5, 0.13, 1.65, 30, 7, 14, 60),
('windfury_matron', 'stonetalon', '风怒主母', 18, 'normal', 84, 18, 0, 11, 11, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 7, 15, 50),
('grimtotem_tauren', 'stonetalon', '恐怖图腾牛头人', 19, 'normal', 90, 19, 0, 11, 11, 'physical', 0.10, 1.6, 0.07, 1.5, 34, 8, 16, 40),
('grimtotem_shaman', 'stonetalon', '恐怖图腾萨满', 20, 'normal', 81, 12, 24, 9, 14, 'magic', 0.07, 1.5, 0.14, 1.65, 36, 8, 17, 30),
('windfury_queen', 'stonetalon', '风怒女王', 22, 'elite', 210, 40, 0, 19, 19, 'physical', 0.15, 1.7, 0.07, 1.5, 60, 12, 24, 5),
('stonetalon_bear', 'stonetalon', '石爪山熊', 21, 'normal', 65, 16, 0, 8, 8, 'physical', 0.08, 1.6, 0.05, 1.5, 38, 8, 18, 25);

-- 灰谷 (18-30级) - HP: 50~120, 攻击: 14~28
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('ashenvale_bear', 'ashenvale', '灰谷熊', 18, 'normal', 87, 22, 0, 11, 11, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 7, 15, 100),
('satyr_rogue', 'ashenvale', '萨特盗贼', 19, 'normal', 93, 23, 0, 11, 11, 'physical', 0.10, 1.6, 0.07, 1.5, 34, 7, 16, 80),
('satyr_shadowdancer', 'ashenvale', '萨特暗影舞者', 20, 'normal', 99, 24, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 36, 8, 17, 70),
('satyr_felsworn', 'ashenvale', '萨特恶魔信徒', 21, 'normal', 87, 14, 30, 11, 15, 'magic', 0.07, 1.5, 0.15, 1.7, 38, 8, 18, 60),
('satyr_hellcaller', 'ashenvale', '萨特地狱召唤者', 22, 'normal', 92, 15, 32, 12, 16, 'magic', 0.07, 1.5, 0.16, 1.7, 40, 9, 19, 50),
('furbolg_warrior', 'ashenvale', '熊怪战士', 23, 'normal', 101, 27, 0, 14, 14, 'physical', 0.11, 1.6, 0.07, 1.5, 42, 9, 20, 40),
('satyr_lord', 'ashenvale', '萨特领主', 25, 'elite', 266, 49, 0, 20, 20, 'physical', 0.16, 1.7, 0.07, 1.5, 68, 13, 26, 5),
('ashenvale_stag', 'ashenvale', '灰谷雄鹿', 24, 'normal', 70, 19, 0, 9, 9, 'physical', 0.08, 1.6, 0.05, 1.5, 40, 9, 19, 35);

-- 塔伦米尔 (20-30级) - HP: 55~120, 攻击: 16~28
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('syndicate_thug', 'tarren_mill', '辛迪加暴徒', 20, 'normal', 94, 24, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 28, 5, 10, 100),
('syndicate_rogue', 'tarren_mill', '辛迪加盗贼', 21, 'normal', 98, 26, 0, 12, 12, 'physical', 0.10, 1.6, 0.07, 1.5, 30, 5, 11, 80),
('syndicate_mage', 'tarren_mill', '辛迪加法师', 22, 'normal', 105, 27, 0, 14, 14, 'physical', 0.10, 1.6, 0.07, 1.5, 32, 6, 12, 70),
('syndicate_assassin', 'tarren_mill', '辛迪加刺客', 23, 'normal', 95, 15, 32, 12, 16, 'magic', 0.07, 1.5, 0.15, 1.7, 34, 6, 13, 60),
('syndicate_warlock', 'tarren_mill', '辛迪加术士', 24, 'normal', 101, 16, 35, 14, 18, 'magic', 0.07, 1.5, 0.16, 1.7, 36, 7, 14, 50),
('hillsbrad_peasant', 'tarren_mill', '希尔斯布莱德农民', 25, 'normal', 112, 30, 0, 15, 15, 'physical', 0.12, 1.6, 0.07, 1.5, 38, 7, 15, 40),
('syndicate_master', 'tarren_mill', '辛迪加首领', 26, 'elite', 252, 47, 0, 20, 20, 'physical', 0.16, 1.7, 0.07, 1.5, 65, 12, 25, 5),
('tarren_mill_guard', 'tarren_mill', '塔伦米尔守卫', 24, 'normal', 78, 21, 0, 10, 10, 'physical', 0.09, 1.6, 0.05, 1.5, 36, 7, 14, 35);

-- 千针石林 (25-35级) - HP: 65~130, 攻击: 18~32
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('centaur_brave', 'thousand_needles', '半人马勇士', 25, 'normal', 119, 34, 0, 15, 15, 'physical', 0.12, 1.6, 0.07, 1.5, 38, 7, 15, 100),
('centaur_warrior', 'thousand_needles', '半人马战士', 26, 'normal', 126, 35, 0, 16, 16, 'physical', 0.12, 1.6, 0.07, 1.5, 40, 8, 16, 80),
('centaur_shaman', 'thousand_needles', '半人马萨满', 27, 'normal', 116, 16, 38, 15, 19, 'magic', 0.07, 1.5, 0.17, 1.7, 42, 8, 17, 70),
('centaur_outrunner', 'thousand_needles', '半人马斥候', 28, 'normal', 133, 38, 0, 18, 18, 'physical', 0.13, 1.6, 0.07, 1.5, 44, 9, 18, 60),
('centaur_chieftain', 'thousand_needles', '半人马酋长', 30, 'elite', 280, 51, 0, 22, 22, 'physical', 0.17, 1.7, 0.07, 1.5, 70, 14, 28, 5),
('windfury_harpy', 'thousand_needles', '风怒鹰身人', 29, 'normal', 123, 35, 0, 16, 16, 'physical', 0.12, 1.6, 0.07, 1.5, 46, 9, 19, 50),
('windfury_matron', 'thousand_needles', '风怒主母', 31, 'normal', 129, 36, 0, 18, 18, 'physical', 0.14, 1.6, 0.08, 1.5, 48, 10, 20, 40),
('thousand_needles_basilisk', 'thousand_needles', '千针石林蜥蜴', 30, 'normal', 90, 28, 0, 12, 12, 'physical', 0.1, 1.6, 0.05, 1.5, 46, 9, 19, 35);

-- 凄凉之地 (30-40级) - HP: 75~140, 攻击: 22~35
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('centaur_brave', 'desolace', '半人马勇士', 30, 'normal', 133, 38, 0, 16, 16, 'physical', 0.12, 1.6, 0.07, 1.5, 42, 8, 16, 100),
('centaur_warrior', 'desolace', '半人马战士', 31, 'normal', 140, 39, 0, 18, 18, 'physical', 0.14, 1.6, 0.08, 1.5, 44, 8, 17, 80),
('centaur_shaman', 'desolace', '半人马萨满', 32, 'normal', 130, 18, 40, 16, 20, 'magic', 0.08, 1.5, 0.19, 1.7, 46, 9, 18, 70),
('centaur_outrunner', 'desolace', '半人马斥候', 33, 'normal', 147, 42, 0, 19, 19, 'physical', 0.15, 1.6, 0.08, 1.5, 48, 9, 19, 60),
('centaur_chieftain', 'desolace', '半人马酋长', 35, 'elite', 308, 54, 0, 23, 23, 'physical', 0.19, 1.7, 0.08, 1.5, 75, 15, 30, 5),
('felhound', 'desolace', '地狱犬', 34, 'normal', 137, 19, 43, 18, 22, 'magic', 0.08, 1.5, 0.20, 1.7, 50, 10, 20, 50),
('doomguard', 'desolace', '末日守卫', 36, 'elite', 336, 57, 0, 24, 24, 'physical', 0.20, 1.7, 0.08, 1.5, 80, 16, 32, 3),
('desolace_basilisk', 'desolace', '凄凉之地蜥蜴', 33, 'normal', 103, 30, 0, 13, 13, 'physical', 0.11, 1.6, 0.05, 1.5, 48, 9, 19, 40);

-- 荆棘谷 (30-45级) - HP: 75~160, 攻击: 22~40
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('jungle_stalker', 'stranglethorn', '丛林潜伏者', 30, 'normal', 133, 38, 0, 16, 16, 'physical', 0.12, 1.6, 0.07, 1.5, 42, 8, 16, 100),
('gurubashi_troll', 'stranglethorn', '古拉巴什巨魔', 32, 'normal', 140, 40, 0, 18, 18, 'physical', 0.14, 1.6, 0.08, 1.5, 46, 9, 18, 80),
('gurubashi_warrior', 'stranglethorn', '古拉巴什战士', 34, 'normal', 147, 43, 0, 19, 19, 'physical', 0.15, 1.6, 0.08, 1.5, 50, 10, 20, 70),
('gurubashi_shaman', 'stranglethorn', '古拉巴什萨满', 35, 'normal', 137, 19, 43, 18, 22, 'magic', 0.08, 1.5, 0.19, 1.7, 52, 10, 21, 60),
('gurubashi_berserker', 'stranglethorn', '古拉巴什狂战士', 36, 'normal', 154, 46, 0, 20, 20, 'physical', 0.16, 1.6, 0.08, 1.5, 54, 11, 22, 50),
('panther', 'stranglethorn', '黑豹', 37, 'normal', 151, 45, 0, 19, 19, 'physical', 0.15, 1.6, 0.08, 1.5, 56, 11, 23, 40),
('gurubashi_chieftain', 'stranglethorn', '古拉巴什酋长', 38, 'elite', 336, 57, 0, 24, 24, 'physical', 0.20, 1.7, 0.08, 1.5, 85, 17, 34, 5),
('tiger', 'stranglethorn', '猛虎', 39, 'normal', 161, 49, 0, 22, 22, 'physical', 0.17, 1.6, 0.08, 1.5, 58, 12, 24, 30),
('bloodscalp_troll', 'stranglethorn', '血顶巨魔', 40, 'normal', 168, 51, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 25),
('bloodscalp_chieftain', 'stranglethorn', '血顶酋长', 42, 'elite', 260, 45, 0, 19, 19, 'physical', 0.18, 1.7, 0.05, 1.5, 90, 18, 36, 3);

-- 荒芜之地 (35-45级) - HP: 85~160, 攻击: 25~40
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('rock_elemental', 'badlands', '石元素', 35, 'normal', 137, 19, 43, 18, 22, 'magic', 0.08, 1.5, 0.19, 1.7, 52, 10, 21, 100),
('dark_iron_dwarf', 'badlands', '黑铁矮人', 36, 'normal', 154, 46, 0, 20, 20, 'physical', 0.16, 1.6, 0.08, 1.5, 54, 11, 22, 80),
('dark_iron_warrior', 'badlands', '黑石战士', 37, 'normal', 161, 47, 0, 22, 22, 'physical', 0.17, 1.6, 0.08, 1.5, 56, 11, 23, 70),
('dark_iron_sorcerer', 'badlands', '黑铁术士', 38, 'normal', 151, 20, 46, 20, 23, 'magic', 0.08, 1.5, 0.20, 1.7, 58, 12, 24, 60),
('dark_iron_commander', 'badlands', '黑铁指挥官', 40, 'elite', 364, 61, 0, 26, 26, 'physical', 0.21, 1.7, 0.08, 1.5, 90, 18, 36, 5),
('basilisk', 'badlands', '蜥蜴', 39, 'normal', 168, 50, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 50),
('rock_giant', 'badlands', '岩石巨人', 42, 'elite', 378, 62, 0, 26, 26, 'physical', 0.22, 1.7, 0.08, 1.5, 95, 19, 38, 3),
('badlands_scorpion', 'badlands', '荒芜之地蝎子', 41, 'normal', 125, 39, 0, 18, 18, 'physical', 0.16, 1.6, 0.05, 1.5, 62, 13, 26, 40);

-- 悲伤沼泽 (35-45级) - HP: 85~160, 攻击: 25~40
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('mire_creeper', 'swamp_of_sorrows', '泥沼爬行者', 35, 'normal', 137, 19, 43, 18, 22, 'magic', 0.08, 1.5, 0.19, 1.7, 52, 10, 21, 100),
('green_dragon_whelp', 'swamp_of_sorrows', '绿龙幼崽', 36, 'normal', 154, 46, 0, 20, 20, 'physical', 0.16, 1.6, 0.08, 1.5, 54, 11, 22, 80),
('green_dragon_scalebane', 'swamp_of_sorrows', '绿龙鳞片守卫', 37, 'normal', 161, 47, 0, 22, 22, 'physical', 0.17, 1.6, 0.08, 1.5, 56, 11, 23, 70),
('green_dragon_mage', 'swamp_of_sorrows', '绿龙法师', 38, 'normal', 151, 20, 46, 20, 23, 'magic', 0.08, 1.5, 0.20, 1.7, 58, 12, 24, 60),
('green_dragon_guardian', 'swamp_of_sorrows', '绿龙守护者', 40, 'elite', 364, 61, 0, 26, 26, 'physical', 0.21, 1.7, 0.08, 1.5, 90, 18, 36, 5),
('marsh_lurker', 'swamp_of_sorrows', '沼泽潜伏者', 39, 'normal', 168, 50, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 50),
('green_dragon_ancient', 'swamp_of_sorrows', '绿龙古龙', 42, 'elite', 378, 62, 0, 26, 26, 'physical', 0.22, 1.7, 0.08, 1.5, 95, 19, 38, 3),
('swamp_elemental', 'swamp_of_sorrows', '沼泽元素', 41, 'normal', 125, 16, 36, 18, 19, 'magic', 0.05, 1.5, 0.18, 1.7, 62, 13, 26, 40);

-- 尘泥沼泽 (35-45级) - HP: 85~160, 攻击: 25~40
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('mire_creeper', 'dustwallow', '泥沼爬行者', 35, 'normal', 137, 19, 43, 18, 22, 'magic', 0.08, 1.5, 0.19, 1.7, 52, 10, 21, 100),
('black_dragon_whelp', 'dustwallow', '黑龙幼崽', 36, 'normal', 154, 46, 0, 20, 20, 'physical', 0.16, 1.6, 0.08, 1.5, 54, 11, 22, 80),
('black_dragon_scalebane', 'dustwallow', '黑龙鳞片守卫', 37, 'normal', 161, 47, 0, 22, 22, 'physical', 0.17, 1.6, 0.08, 1.5, 56, 11, 23, 70),
('black_dragon_mage', 'dustwallow', '黑龙法师', 38, 'normal', 151, 20, 46, 20, 23, 'magic', 0.08, 1.5, 0.20, 1.7, 58, 12, 24, 60),
('black_dragon_guardian', 'dustwallow', '黑龙守护者', 40, 'elite', 364, 61, 0, 26, 26, 'physical', 0.21, 1.7, 0.08, 1.5, 90, 18, 36, 5),
('marsh_lurker', 'dustwallow', '沼泽潜伏者', 39, 'normal', 168, 50, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 50),
('black_dragon_ancient', 'dustwallow', '黑龙古龙', 42, 'elite', 378, 62, 0, 26, 26, 'physical', 0.22, 1.7, 0.08, 1.5, 95, 19, 38, 3),
('dustwallow_elemental', 'dustwallow', '尘泥元素', 41, 'normal', 125, 16, 36, 18, 19, 'magic', 0.05, 1.5, 0.18, 1.7, 62, 13, 26, 40);

-- 塔纳利斯 (40-50级) - HP: 95~180, 攻击: 28~45
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('sand_elemental', 'tanaris', '沙元素', 40, 'normal', 168, 51, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 100),
('wastewander_bandit', 'tanaris', '废土强盗', 41, 'normal', 169, 51, 0, 23, 23, 'physical', 0.19, 1.6, 0.08, 1.5, 62, 13, 26, 80),
('wastewander_rogue', 'tanaris', '废土盗贼', 42, 'normal', 176, 52, 0, 25, 25, 'physical', 0.20, 1.6, 0.08, 1.5, 64, 13, 27, 70),
('wastewander_mage', 'tanaris', '废土法师', 43, 'normal', 166, 22, 49, 23, 26, 'magic', 0.08, 1.5, 0.22, 1.7, 66, 14, 28, 60),
('wastewander_chieftain', 'tanaris', '废土酋长', 45, 'elite', 405, 65, 0, 29, 29, 'physical', 0.23, 1.8, 0.08, 1.5, 100, 20, 40, 5),
('sand_giant', 'tanaris', '沙巨人', 44, 'normal', 182, 55, 0, 26, 26, 'physical', 0.21, 1.6, 0.08, 1.5, 68, 14, 29, 50),
('ancient_sand_elemental', 'tanaris', '古沙元素', 46, 'elite', 432, 68, 0, 30, 30, 'physical', 0.24, 1.8, 0.08, 1.5, 105, 21, 42, 3),
('tanaris_scorpion', 'tanaris', '塔纳利斯蝎子', 43, 'normal', 128, 41, 0, 19, 19, 'physical', 0.17, 1.6, 0.05, 1.5, 66, 14, 28, 40);

-- 菲拉斯 (40-50级) - HP: 95~180, 攻击: 28~45
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('feralas_bear', 'feralas', '菲拉斯熊', 40, 'normal', 168, 51, 0, 23, 23, 'physical', 0.18, 1.6, 0.08, 1.5, 60, 12, 25, 100),
('feralas_panther', 'feralas', '菲拉斯黑豹', 41, 'normal', 169, 51, 0, 23, 23, 'physical', 0.19, 1.6, 0.08, 1.5, 62, 13, 26, 80),
('feralas_tiger', 'feralas', '菲拉斯猛虎', 42, 'normal', 176, 52, 0, 25, 25, 'physical', 0.20, 1.6, 0.08, 1.5, 64, 13, 27, 70),
('highborne_ruins_guardian', 'feralas', '上层精灵废墟守卫', 43, 'normal', 166, 22, 49, 23, 26, 'magic', 0.08, 1.5, 0.22, 1.7, 66, 14, 28, 60),
('highborne_ruins_mage', 'feralas', '上层精灵废墟法师', 44, 'normal', 173, 23, 52, 25, 27, 'magic', 0.08, 1.5, 0.23, 1.7, 68, 14, 29, 50),
('highborne_ruins_ancient', 'feralas', '上层精灵废墟古灵', 46, 'elite', 432, 68, 0, 30, 30, 'physical', 0.24, 1.8, 0.08, 1.5, 105, 21, 42, 5),
('feralas_basilisk', 'feralas', '菲拉斯蜥蜴', 45, 'normal', 182, 55, 0, 26, 26, 'physical', 0.21, 1.6, 0.08, 1.5, 70, 15, 30, 40),
('feralas_giant', 'feralas', '菲拉斯巨人', 47, 'elite', 340, 54, 0, 24, 24, 'physical', 0.22, 1.8, 0.05, 1.5, 110, 22, 44, 3);

-- 安戈洛环形山 (48-55级) - HP: 110~200, 攻击: 32~50
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('diemetradon', 'ungoro', '双帆龙', 48, 'normal', 189, 57, 0, 26, 26, 'physical', 0.21, 1.6, 0.08, 1.5, 72, 15, 30, 100),
('pterrordax', 'ungoro', '翼手龙', 49, 'normal', 196, 58, 0, 27, 27, 'physical', 0.22, 1.6, 0.08, 1.5, 74, 15, 31, 80),
('devilsaur', 'ungoro', '魔暴龙', 50, 'normal', 202, 60, 0, 29, 29, 'physical', 0.23, 1.6, 0.08, 1.5, 76, 16, 32, 70),
('ungoro_thunderer', 'ungoro', '安戈洛雷霆蜥蜴', 51, 'normal', 209, 61, 0, 30, 30, 'physical', 0.24, 1.6, 0.08, 1.5, 78, 16, 33, 60),
('ungoro_gorger', 'ungoro', '安戈洛吞噬者', 52, 'normal', 216, 62, 0, 31, 31, 'physical', 0.25, 1.6, 0.08, 1.5, 80, 17, 34, 50),
('devilsaur_alpha', 'ungoro', '魔暴龙王', 53, 'elite', 486, 75, 0, 34, 34, 'physical', 0.26, 1.8, 0.08, 1.5, 115, 23, 46, 5),
('ungoro_elemental', 'ungoro', '安戈洛元素', 52, 'normal', 213, 25, 55, 31, 32, 'magic', 0.08, 1.5, 0.24, 1.7, 80, 17, 34, 40),
('ancient_devilsaur', 'ungoro', '古魔暴龙', 54, 'elite', 380, 60, 0, 27, 27, 'physical', 0.24, 1.8, 0.05, 1.5, 120, 24, 48, 3);

-- 费伍德森林 (48-55级) - HP: 110~200, 攻击: 32~50
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('felwood_satyr', 'felwood', '费伍德萨特', 48, 'normal', 189, 57, 0, 26, 26, 'physical', 0.21, 1.6, 0.08, 1.5, 72, 15, 30, 100),
('felwood_felguard', 'felwood', '费伍德恶魔守卫', 49, 'normal', 196, 58, 0, 27, 27, 'physical', 0.22, 1.6, 0.08, 1.5, 74, 15, 31, 80),
('felwood_doomguard', 'felwood', '费伍德末日守卫', 50, 'normal', 202, 60, 0, 29, 29, 'physical', 0.23, 1.6, 0.08, 1.5, 76, 16, 32, 70),
('felwood_infernal', 'felwood', '费伍德地狱火', 51, 'normal', 209, 26, 57, 30, 31, 'magic', 0.08, 1.5, 0.25, 1.7, 78, 16, 33, 60),
('felwood_demon_lord', 'felwood', '费伍德恶魔领主', 52, 'elite', 486, 75, 0, 34, 34, 'physical', 0.26, 1.8, 0.08, 1.5, 115, 23, 46, 5),
('corrupted_treant', 'felwood', '被腐蚀的古树', 51, 'normal', 207, 61, 0, 29, 29, 'physical', 0.24, 1.6, 0.08, 1.5, 78, 16, 33, 50),
('felwood_ancient', 'felwood', '费伍德古灵', 53, 'elite', 513, 78, 0, 35, 35, 'physical', 0.27, 1.8, 0.08, 1.5, 120, 24, 48, 3),
('felwood_corrupted_bear', 'felwood', '被腐蚀的熊', 50, 'normal', 148, 45, 0, 21, 21, 'physical', 0.19, 1.6, 0.05, 1.5, 76, 16, 32, 40);

-- 冬泉谷 (55-60级) - HP: 130~220, 攻击: 38~55
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('winterspring_frostsaber', 'winterspring', '冬泉谷霜刃豹', 55, 'normal', 230, 68, 0, 32, 32, 'physical', 0.25, 1.6, 0.08, 1.5, 84, 17, 35, 100),
('winterspring_owl', 'winterspring', '冬泉谷猫头鹰', 56, 'normal', 236, 69, 0, 34, 34, 'physical', 0.26, 1.6, 0.08, 1.5, 86, 18, 36, 80),
('winterspring_furbolg', 'winterspring', '冬泉谷熊怪', 57, 'normal', 243, 70, 0, 35, 35, 'physical', 0.27, 1.6, 0.08, 1.5, 88, 18, 37, 70),
('blue_dragon_whelp', 'winterspring', '蓝龙幼崽', 58, 'normal', 250, 72, 0, 36, 36, 'physical', 0.28, 1.6, 0.08, 1.5, 90, 19, 38, 60),
('blue_dragon_scalebane', 'winterspring', '蓝龙鳞片守卫', 59, 'normal', 256, 73, 0, 38, 38, 'physical', 0.29, 1.6, 0.08, 1.5, 92, 19, 39, 50),
('blue_dragon_guardian', 'winterspring', '蓝龙守护者', 60, 'elite', 540, 84, 0, 39, 39, 'physical', 0.31, 1.8, 0.08, 1.5, 130, 26, 52, 5),
('winterspring_elemental', 'winterspring', '冬泉谷元素', 58, 'normal', 247, 29, 62, 36, 38, 'magic', 0.08, 1.5, 0.27, 1.7, 90, 19, 38, 40),
('blue_dragon_ancient', 'winterspring', '蓝龙古龙', 60, 'boss', 450, 70, 0, 32, 32, 'physical', 0.3, 1.9, 0.05, 1.5, 150, 30, 60, 1);

-- 希利苏斯 (55-60级) - HP: 130~220, 攻击: 38~55
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('silithid_worker', 'silithus', '其拉工虫', 55, 'normal', 230, 68, 0, 32, 32, 'physical', 0.25, 1.6, 0.08, 1.5, 84, 17, 35, 100),
('silithid_warrior', 'silithus', '其拉战士', 56, 'normal', 236, 69, 0, 34, 34, 'physical', 0.26, 1.6, 0.08, 1.5, 86, 18, 36, 80),
('silithid_royal', 'silithus', '其拉皇家守卫', 57, 'normal', 243, 70, 0, 35, 35, 'physical', 0.27, 1.6, 0.08, 1.5, 88, 18, 37, 70),
('silithid_mage', 'silithus', '其拉法师', 58, 'normal', 234, 30, 62, 35, 36, 'magic', 0.08, 1.5, 0.27, 1.7, 90, 19, 38, 60),
('silithid_queen', 'silithus', '其拉女王', 59, 'elite', 567, 88, 0, 40, 40, 'physical', 0.30, 1.8, 0.08, 1.5, 135, 27, 54, 5),
('anubisath_guardian', 'silithus', '阿努比萨斯守卫', 58, 'normal', 250, 72, 0, 36, 36, 'physical', 0.28, 1.6, 0.08, 1.5, 90, 19, 38, 50),
('qiraji_warrior', 'silithus', '其拉战士', 59, 'normal', 256, 73, 0, 38, 38, 'physical', 0.29, 1.6, 0.08, 1.5, 92, 19, 39, 40),
('c_thun_minion', 'silithus', '克苏恩仆从', 60, 'boss', 450, 70, 0, 32, 32, 'physical', 0.3, 1.9, 0.05, 1.5, 150, 30, 60, 1);

-- 燃烧平原 (50-60级) - HP: 120~220, 攻击: 35~55
INSERT OR REPLACE INTO monsters (
id, zone_id, name, level, type, hp, physical_attack, magic_attack, physical_defense, magic_defense,
attack_type, phys_crit_rate, phys_crit_damage, spell_crit_rate, spell_crit_damage,
exp_reward, gold_min, gold_max, spawn_weight
) VALUES
('blackrock_orc', 'burning_steppes', '黑石兽人', 50, 'normal', 202, 60, 0, 29, 29, 'physical', 0.23, 1.6, 0.08, 1.5, 76, 16, 32, 100),
('blackrock_warrior', 'burning_steppes', '黑石战士', 52, 'normal', 216, 62, 0, 31, 31, 'physical', 0.25, 1.6, 0.08, 1.5, 80, 17, 34, 80),
('blackrock_warlock', 'burning_steppes', '黑石术士', 53, 'normal', 209, 26, 57, 30, 32, 'magic', 0.08, 1.5, 0.25, 1.7, 82, 17, 35, 70),
('black_dragon_whelp', 'burning_steppes', '黑龙幼崽', 54, 'normal', 223, 65, 0, 32, 32, 'physical', 0.26, 1.6, 0.08, 1.5, 84, 18, 36, 60),
('black_dragon_scalebane', 'burning_steppes', '黑龙鳞片守卫', 55, 'normal', 230, 68, 0, 34, 34, 'physical', 0.27, 1.6, 0.08, 1.5, 86, 18, 37, 50),
('black_dragon_guardian', 'burning_steppes', '黑龙守护者', 56, 'elite', 540, 84, 0, 39, 39, 'physical', 0.31, 1.8, 0.08, 1.5, 130, 26, 52, 5),
('blackrock_commander', 'burning_steppes', '黑石指挥官', 57, 'elite', 554, 86, 0, 40, 40, 'physical', 0.32, 1.8, 0.08, 1.5, 135, 27, 54, 3),
('black_dragon_ancient', 'burning_steppes', '黑龙古龙', 58, 'elite', 567, 88, 0, 42, 42, 'physical', 0.33, 1.8, 0.08, 1.5, 140, 28, 56, 3),
('nefarian', 'burning_steppes', '奈法利安', 60, 'boss', 500, 75, 0, 35, 35, 'physical', 0.35, 2.0, 0.05, 1.5, 200, 40, 80, 1);

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
 'class_base_hp + stamina × 2 + level × hp_per_level + equipment_bonus',
 '计算角色的最大生命值。耐力是最主要的生命值来源。',
 '{"class_base_hp":"职业基础HP(战士35,法师20等)","stamina":"耐力值","level":"角色等级","hp_per_level":"每级HP成长","equipment_bonus":"装备HP加成"}',
 '战士(1级,耐力10): 35 + 10×2 = 55 HP', 1),

('attr_max_mp', 'attribute', '最大法力值',
 'class_base_mp + spirit × 2 + level × mp_per_level + equipment_bonus',
 '计算角色的最大法力值。精神是最主要的法力来源。仅法力职业有效。',
 '{"class_base_mp":"职业基础MP","spirit":"精神值","level":"角色等级","mp_per_level":"每级MP成长","equipment_bonus":"装备MP加成"}',
 '法师(1级,精神10): 40 + 10×2 = 60 MP', 2),

('attr_phys_attack', 'attribute', '物理攻击力',
 'strength × 0.4 + agility × 0.2 + weapon_damage + equipment_bonus',
 '计算角色的物理攻击力。力量是主要来源，敏捷提供辅助加成。',
 '{"strength":"力量值","agility":"敏捷值","weapon_damage":"武器伤害","equipment_bonus":"其他装备加成"}',
 '战士(力量12,敏捷8): 12×0.4 + 8×0.2 = 6.4', 3),

('attr_spell_power', 'attribute', '法术攻击力',
 'intellect × 1.0 + spirit × 0.2 + equipment_bonus',
 '计算角色的法术攻击力。智力是主要来源，精神提供辅助加成。',
 '{"intellect":"智力值","spirit":"精神值","equipment_bonus":"装备法伤加成"}',
 '法师(智力14,精神10): 14×1.0 + 10×0.2 = 16', 4),

('attr_phys_defense', 'attribute', '物理防御',
 'strength × 0.1 + stamina × 0.3 + equipment_armor',
 '计算角色的物理防御。力量和耐力共同提供加成。',
 '{"strength":"力量值","stamina":"耐力值","equipment_armor":"装备护甲"}',
 '战士(力量12,耐力10): 12×0.1 + 10×0.3 = 4.2', 5),

('attr_magic_defense', 'attribute', '魔法防御',
 'intellect × 0.2 + spirit × 0.3 + equipment_bonus',
 '计算角色的魔法防御。智力和精神共同提供加成。',
 '{"intellect":"智力值","spirit":"精神值","equipment_bonus":"装备魔抗加成"}',
 '法师(智力14,精神10): 14×0.2 + 10×0.3 = 6', 6),

-- ═══════════════════════════════════════════════════════════
-- 战斗判定公式
-- ═══════════════════════════════════════════════════════════
('combat_phys_crit', 'combat', '物理暴击率',
 '5% + agility ÷ 20 + equipment_bonus% (上限100%)',
 '计算物理攻击的暴击概率。每20点敏捷增加1%暴击率。',
 '{"agility":"敏捷值","equipment_bonus":"装备暴击加成%"}',
 '敏捷100: 5% + 5% = 10%', 10),

('combat_spell_crit', 'combat', '法术暴击率',
 '5% + spirit ÷ 20 + equipment_bonus% (上限100%)',
 '计算法术攻击的暴击概率。每20点精神增加1%法术暴击率。',
 '{"spirit":"精神值","equipment_bonus":"装备法术暴击加成%"}',
 '精神100: 5% + 5% = 10%', 11),

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

-- ═══════════════════════════════════════════════════════════
-- 装备词缀数据
-- ═══════════════════════════════════════════════════════════

-- 前缀 (攻击/属性向)
INSERT OR REPLACE INTO affixes (id, name, type, slot_type, rarity, effect_type, effect_stat, min_value, max_value, value_type, description, level_required) VALUES
-- 基础攻击前缀
('affix_sharp', '锋利的', 'prefix', 'weapon', 'common', 'stat_mod', 'attack', 2, 5, 'flat', '攻击力 +{value}', 1),
('affix_fiery', '炽热的', 'prefix', 'weapon', 'common', 'damage_mod', 'fire_damage', 1, 4, 'flat', '火焰伤害 +{value}', 1),
('affix_frozen', '冰霜的', 'prefix', 'weapon', 'common', 'damage_mod', 'frost_damage', 1, 4, 'flat', '冰霜伤害 +{value}', 1),
('affix_charged', '雷击的', 'prefix', 'weapon', 'common', 'damage_mod', 'lightning_damage', 1, 4, 'flat', '雷电伤害 +{value}', 1),
('affix_holy', '神圣的', 'prefix', 'weapon', 'uncommon', 'damage_mod', 'holy_damage', 2, 5, 'flat', '神圣伤害 +{value}', 10),
('affix_vampiric', '吸血鬼的', 'prefix', 'weapon', 'rare', 'stat_mod', 'lifesteal', 2, 5, 'percent', '生命偷取 +{value}%', 20),
('affix_devastating', '毁灭的', 'prefix', 'weapon', 'epic', 'stat_mod', 'attack', 15, 25, 'percent', '攻击力 +{value}%', 30),
-- 防御前缀
('affix_sturdy', '坚固的', 'prefix', 'armor', 'common', 'stat_mod', 'defense', 2, 5, 'flat', '防御力 +{value}', 1),
('affix_vital', '活力的', 'prefix', 'armor', 'common', 'stat_mod', 'max_hp', 5, 15, 'flat', '生命值 +{value}', 1),
('affix_scholarly', '智者的', 'prefix', 'armor', 'uncommon', 'stat_mod', 'intellect', 2, 4, 'flat', '智力 +{value}', 10),
('affix_unyielding', '不屈的', 'prefix', 'armor', 'rare', 'stat_mod', 'damage_reduction', 3, 8, 'percent', '受伤减免 +{value}%', 20),
-- 仇恨相关前缀 (坦克)
('affix_guardian', '守护者的', 'prefix', 'weapon', 'rare', 'stat_mod', 'threat_gen', 20, 25, 'percent', '仇恨生成 +{value}%, 嘲讽CD -1', 15),
('affix_imposing', '威压的', 'prefix', 'armor', 'uncommon', 'stat_mod', 'threat_gen', 10, 15, 'percent', '仇恨生成 +{value}%', 10),
-- 仇恨相关前缀 (DPS)
('affix_stealthy', '隐秘的', 'prefix', 'weapon', 'uncommon', 'stat_mod', 'threat_gen', -15, -20, 'percent', '仇恨生成 {value}%', 10),
('affix_shadow', '暗影的', 'prefix', 'armor', 'rare', 'stat_mod', 'crit_threat', -25, -30, 'percent', '暴击仇恨 {value}%', 20);

-- 后缀 (特殊效果向)
INSERT OR REPLACE INTO affixes (id, name, type, slot_type, rarity, effect_type, effect_stat, min_value, max_value, value_type, description, level_required) VALUES
('affix_of_power', 'of 力量', 'suffix', 'all', 'common', 'stat_mod', 'strength', 1, 3, 'flat', '力量 +{value}', 1),
('affix_of_agility', 'of 敏捷', 'suffix', 'all', 'common', 'stat_mod', 'agility', 1, 3, 'flat', '敏捷 +{value}', 1),
('affix_of_intellect', 'of 智力', 'suffix', 'all', 'common', 'stat_mod', 'intellect', 1, 3, 'flat', '智力 +{value}', 1),
('affix_of_stamina', 'of 耐力', 'suffix', 'all', 'common', 'stat_mod', 'stamina', 1, 3, 'flat', '耐力 +{value}', 1),
('affix_of_spirit', 'of 精神', 'suffix', 'all', 'common', 'stat_mod', 'spirit', 1, 3, 'flat', '精神 +{value}', 1),
('affix_of_tiger', 'of 猛虎', 'suffix', 'all', 'uncommon', 'stat_mod', 'attack,agility', 3, 5, 'flat', '攻击力和敏捷 +{value}', 15),
('affix_of_bear', 'of 巨熊', 'suffix', 'all', 'uncommon', 'stat_mod', 'stamina,defense', 3, 6, 'flat', '耐力和防御 +{value}', 15),
('affix_of_crit', 'of 暴击', 'suffix', 'weapon', 'rare', 'stat_mod', 'crit_rate', 3, 8, 'percent', '暴击率 +{value}%', 20),
('affix_of_haste', 'of 迅捷', 'suffix', 'weapon', 'rare', 'stat_mod', 'attack_speed', 8, 15, 'percent', '攻击速度 +{value}%', 20),
('affix_of_slay', 'of 屠戮', 'suffix', 'weapon', 'epic', 'stat_mod', 'crit_damage', 20, 35, 'percent', '暴击伤害 +{value}%', 30),
-- 仇恨相关后缀 (坦克)
('affix_of_threat', 'of 威胁', 'suffix', 'armor', 'uncommon', 'stat_mod', 'threat_gen', 10, 15, 'percent', '仇恨生成 +{value}%', 10),
('affix_of_guardian', 'of 守护', 'suffix', 'shield', 'uncommon', 'stat_mod', 'threat_gen,block', 8, 12, 'percent', '仇恨 +{value}%, 格挡 +5%', 15),
-- 仇恨相关后缀 (DPS/治疗)
('affix_of_stealth', 'of 隐匿', 'suffix', 'armor', 'uncommon', 'stat_mod', 'threat_gen', -10, -15, 'percent', '仇恨生成 {value}%', 10),
('affix_of_fade', 'of 消散', 'suffix', 'armor', 'rare', 'stat_mod', 'threat_decay', 15, 20, 'percent', '仇恨衰减 +{value}%', 20);

