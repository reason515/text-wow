package battle

// ═══════════════════════════════════════════════════════════
// 技能数据定义 (后续从数据库加载)
// ═══════════════════════════════════════════════════════════

// 职业技能缓存
var skillsCache = map[string][]*Skill{
	// 战士技能
	"warrior": {
		{
			ID: "heroic_strike", Name: "英勇打击", Description: "强力的武器攻击",
			ClassID: "warrior", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 8, ScalingStat: "strength", ScalingRatio: 0.5, ResourceCost: 10, Cooldown: 0,
		},
		{
			ID: "charge", Name: "冲锋", Description: "冲向敌人造成伤害并眩晕",
			ClassID: "warrior", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 5, ScalingStat: "strength", ScalingRatio: 0.3, ResourceCost: 0, Cooldown: 4,
			EffectID: "eff_stun", EffectChance: 1.0,
		},
		{
			ID: "rend", Name: "撕裂", Description: "造成流血效果",
			ClassID: "warrior", Type: SkillDOT, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 2, ScalingStat: "strength", ScalingRatio: 0.15, ResourceCost: 8, Cooldown: 0,
			EffectID: "eff_rend", EffectChance: 1.0,
		},
		{
			ID: "battle_shout", Name: "战斗怒吼", Description: "提升全队攻击力",
			ClassID: "warrior", Type: SkillBuff, TargetType: TargetAllyAll, DamageType: "",
			BaseValue: 0, ResourceCost: 15, Cooldown: 6,
			EffectID: "eff_battle_shout", EffectChance: 1.0,
		},
	},
	// 法师技能
	"mage": {
		{
			ID: "fireball", Name: "火球术", Description: "发射火球，有几率点燃敌人",
			ClassID: "mage", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamageFire,
			BaseValue: 10, ScalingStat: "intellect", ScalingRatio: 0.5, ResourceCost: 5, Cooldown: 0,
			EffectID: "eff_ignite", EffectChance: 0.3,
		},
		{
			ID: "frostbolt", Name: "寒冰箭", Description: "发射寒冰箭，减缓敌人",
			ClassID: "mage", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamageFrost,
			BaseValue: 8, ScalingStat: "intellect", ScalingRatio: 0.4, ResourceCost: 4, Cooldown: 0,
			EffectID: "eff_slow", EffectChance: 0.5,
		},
		{
			ID: "arcane_missiles", Name: "奥术飞弹", Description: "发射多道奥术飞弹",
			ClassID: "mage", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamageMagic,
			BaseValue: 12, ScalingStat: "intellect", ScalingRatio: 0.6, ResourceCost: 7, Cooldown: 2,
		},
		{
			ID: "flamestrike", Name: "烈焰风暴", Description: "对所有敌人造成火焰伤害",
			ClassID: "mage", Type: SkillAttack, TargetType: TargetEnemyAll, DamageType: DamageFire,
			BaseValue: 8, ScalingStat: "intellect", ScalingRatio: 0.4, ResourceCost: 10, Cooldown: 3,
			EffectID: "eff_ignite", EffectChance: 0.2,
		},
	},
	// 盗贼技能
	"rogue": {
		{
			ID: "sinister_strike", Name: "邪恶攻击", Description: "快速攻击",
			ClassID: "rogue", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 6, ScalingStat: "agility", ScalingRatio: 0.4, ResourceCost: 15, Cooldown: 0,
		},
		{
			ID: "ambush", Name: "伏击", Description: "对高血量敌人伤害翻倍",
			ClassID: "rogue", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 10, ScalingStat: "agility", ScalingRatio: 0.5, ResourceCost: 25, Cooldown: 2,
		},
		{
			ID: "deadly_poison", Name: "致命毒药", Description: "使敌人中毒",
			ClassID: "rogue", Type: SkillDOT, TargetType: TargetEnemy, DamageType: DamageNature,
			BaseValue: 2, ScalingStat: "agility", ScalingRatio: 0.1, ResourceCost: 20, Cooldown: 0,
			EffectID: "eff_poison", EffectChance: 1.0,
		},
		{
			ID: "kidney_shot", Name: "肾击", Description: "眩晕敌人",
			ClassID: "rogue", Type: SkillControl, TargetType: TargetEnemy, DamageType: DamagePhysical,
			BaseValue: 4, ScalingStat: "agility", ScalingRatio: 0.2, ResourceCost: 25, Cooldown: 4,
			EffectID: "eff_stun", EffectChance: 1.0,
		},
	},
	// 牧师技能
	"priest": {
		{
			ID: "smite", Name: "惩击", Description: "用神圣能量攻击敌人",
			ClassID: "priest", Type: SkillAttack, TargetType: TargetEnemy, DamageType: DamageHoly,
			BaseValue: 7, ScalingStat: "intellect", ScalingRatio: 0.4, ResourceCost: 4, Cooldown: 0,
		},
		{
			ID: "lesser_heal", Name: "次级治疗术", Description: "恢复生命值",
			ClassID: "priest", Type: SkillHeal, TargetType: TargetAllyLowest, DamageType: "",
			BaseValue: 8, ScalingStat: "spirit", ScalingRatio: 0.4, ResourceCost: 4, Cooldown: 0,
		},
		{
			ID: "renew", Name: "恢复", Description: "持续恢复生命",
			ClassID: "priest", Type: SkillHOT, TargetType: TargetAllyLowest, DamageType: "",
			BaseValue: 3, ScalingStat: "spirit", ScalingRatio: 0.2, ResourceCost: 5, Cooldown: 0,
			EffectID: "eff_renew", EffectChance: 1.0,
		},
		{
			ID: "power_word_shield", Name: "真言术:盾", Description: "创造吸收伤害的护盾",
			ClassID: "priest", Type: SkillShield, TargetType: TargetAllyLowest, DamageType: "",
			BaseValue: 12, ScalingStat: "spirit", ScalingRatio: 0.5, ResourceCost: 8, Cooldown: 4,
			EffectID: "eff_pw_shield", EffectChance: 1.0,
		},
	},
}

// GetSkillsForClass 获取职业技能
func GetSkillsForClass(classID string) []*Skill {
	if skills, ok := skillsCache[classID]; ok {
		return skills
	}
	return []*Skill{}
}

// GetBasicAttack 获取普通攻击技能
func GetBasicAttack() *Skill {
	return &Skill{
		ID:           "basic_attack",
		Name:         "普通攻击",
		Description:  "基础物理攻击",
		Type:         SkillAttack,
		TargetType:   TargetEnemy,
		DamageType:   DamagePhysical,
		BaseValue:    0,
		ScalingStat:  "strength",
		ScalingRatio: 0.5,
		ResourceCost: 0,
		Cooldown:     0,
	}
}

// CreateMonsterSkills 为怪物创建简单技能
func CreateMonsterSkills(level int) []*Skill {
	return []*Skill{
		{
			ID:           "monster_attack",
			Name:         "攻击",
			Type:         SkillAttack,
			TargetType:   TargetEnemy,
			DamageType:   DamagePhysical,
			BaseValue:    3 + level,
			ScalingStat:  "strength",
			ScalingRatio: 0.3,
			ResourceCost: 0,
			Cooldown:     0,
		},
	}
}


