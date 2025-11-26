-- ═══════════════════════════════════════════════════════════
-- Text WoW 初始游戏数据
-- ═══════════════════════════════════════════════════════════

-- ═══════════════════════════════════════════════════════════
-- 种族数据
-- ═══════════════════════════════════════════════════════════

INSERT OR REPLACE INTO races (id, name, faction, description, strength_mod, agility_mod, intellect_mod, stamina_mod, spirit_mod) VALUES
('human', '人类', 'alliance', '适应力强的种族，各项属性平衡。', 1, 0, 0, 0, 1),
('dwarf', '矮人', 'alliance', '坚韧的山地种族，擅长近战和工艺。', 2, 0, 0, 2, 0),
('nightelf', '暗夜精灵', 'alliance', '古老的精灵种族，与自然和谐共存。', 0, 2, 0, 0, 1),
('gnome', '侏儒', 'alliance', '聪明的小型种族，擅长魔法和机械。', 0, 0, 3, 0, 0),
('orc', '兽人', 'horde', '强壮的战士种族，崇尚力量和荣耀。', 3, 0, 0, 1, 0),
('undead', '亡灵', 'horde', '不死的存在，对暗影魔法有天赋。', 0, 0, 2, 0, 2),
('tauren', '牛头人', 'horde', '高大温和的种族，与大地之母相连。', 2, 0, 0, 2, 1),
('troll', '巨魔', 'horde', '敏捷的丛林种族，拥有快速再生能力。', 0, 2, 0, 1, 0);

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
-- 技能数据
-- ═══════════════════════════════════════════════════════════

-- 战士技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target, damage_type, base_damage, damage_scaling, mp_cost, cooldown, level_required) VALUES
('heroic_strike', '英勇打击', '一次强力的武器攻击。', 'warrior', 'attack', 'enemy', 'physical', 10, 1.2, 5, 0, 1),
('charge', '冲锋', '冲向敌人，造成伤害并眩晕。', 'warrior', 'attack', 'enemy', 'physical', 5, 0.8, 10, 3, 1),
('thunder_clap', '雷霆一击', '对周围敌人造成伤害。', 'warrior', 'attack', 'enemy', 'physical', 15, 0.6, 15, 2, 4),
('execute', '斩杀', '对低血量敌人造成大量伤害。', 'warrior', 'attack', 'enemy', 'physical', 50, 1.5, 20, 0, 8),
('shield_wall', '盾墙', '大幅减少受到的伤害。', 'warrior', 'buff', 'self', NULL, 0, 0, 30, 10, 10);

-- 法师技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target, damage_type, base_damage, damage_scaling, mp_cost, cooldown, level_required) VALUES
('fireball', '火球术', '发射一个火球攻击敌人。', 'mage', 'attack', 'enemy', 'magic', 20, 1.3, 15, 0, 1),
('frostbolt', '寒冰箭', '发射寒冰箭，减缓敌人。', 'mage', 'attack', 'enemy', 'magic', 15, 1.1, 12, 0, 1),
('arcane_missiles', '奥术飞弹', '发射多道奥术飞弹。', 'mage', 'attack', 'enemy', 'magic', 25, 1.4, 20, 2, 4),
('pyroblast', '炎爆术', '施放巨大的火球。', 'mage', 'attack', 'enemy', 'magic', 60, 2.0, 40, 5, 8),
('ice_barrier', '寒冰护体', '创造吸收伤害的护盾。', 'mage', 'buff', 'self', NULL, 0, 0, 35, 8, 10);

-- 盗贼技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target, damage_type, base_damage, damage_scaling, mp_cost, cooldown, level_required) VALUES
('sinister_strike', '邪恶攻击', '快速的攻击，积累连击点。', 'rogue', 'attack', 'enemy', 'physical', 12, 1.1, 8, 0, 1),
('backstab', '背刺', '从背后攻击造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 30, 1.5, 15, 0, 1),
('eviscerate', '剔骨', '消耗连击点造成大量伤害。', 'rogue', 'attack', 'enemy', 'physical', 40, 1.8, 25, 0, 4),
('blade_flurry', '剑刃乱舞', '攻击速度大幅提升。', 'rogue', 'buff', 'self', NULL, 0, 0, 30, 8, 8),
('vanish', '消失', '进入潜行状态。', 'rogue', 'buff', 'self', NULL, 0, 0, 40, 10, 10);

-- 牧师技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target, damage_type, base_damage, damage_scaling, mp_cost, cooldown, level_required) VALUES
('smite', '惩击', '用神圣能量攻击敌人。', 'priest', 'attack', 'enemy', 'magic', 15, 1.0, 10, 0, 1),
('shadow_word_pain', '暗言术：痛', '对敌人施加持续伤害。', 'priest', 'debuff', 'enemy', 'magic', 8, 0.5, 12, 0, 1),
('heal', '治疗术', '恢复自身生命值。', 'priest', 'heal', 'self', NULL, 30, 1.2, 20, 0, 1),
('flash_heal', '快速治疗', '快速恢复生命值。', 'priest', 'heal', 'self', NULL, 20, 0.8, 15, 0, 4),
('power_word_shield', '真言术：盾', '创造吸收伤害的护盾。', 'priest', 'buff', 'self', NULL, 0, 0, 25, 4, 6);

-- 通用技能
INSERT OR REPLACE INTO skills (id, name, description, class_id, type, target, damage_type, base_damage, damage_scaling, mp_cost, cooldown, level_required) VALUES
('basic_attack', '普通攻击', '基础的物理攻击。', NULL, 'attack', 'enemy', 'physical', 0, 1.0, 0, 0, 1);

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

