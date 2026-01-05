package repository

import (
	"database/sql"
	"text-wow/internal/database"
	"text-wow/internal/models"
)

// GameRepository 游戏配置数据仓库
type GameRepository struct{}

// NewGameRepository 创建游戏配置仓库
func NewGameRepository() *GameRepository {
	return &GameRepository{}
}

// GetRaces 获取所有种族
func (r *GameRepository) GetRaces() ([]models.Race, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, faction, description,
		       strength_base, agility_base, intellect_base, stamina_base, spirit_base,
		       strength_pct, agility_pct, intellect_pct, stamina_pct, spirit_pct
		FROM races ORDER BY faction, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []models.Race
	for rows.Next() {
		race := models.Race{}
		err := rows.Scan(
			&race.ID, &race.Name, &race.Faction, &race.Description,
			&race.StrengthBase, &race.AgilityBase, &race.IntellectBase, &race.StaminaBase, &race.SpiritBase,
			&race.StrengthPct, &race.AgilityPct, &race.IntellectPct, &race.StaminaPct, &race.SpiritPct,
		)
		if err != nil {
			return nil, err
		}
		races = append(races, race)
	}

	return races, nil
}

// GetRaceByID 根据ID获取种族
func (r *GameRepository) GetRaceByID(id string) (*models.Race, error) {
	race := &models.Race{}
	err := database.DB.QueryRow(`
		SELECT id, name, faction, description,
		       strength_base, agility_base, intellect_base, stamina_base, spirit_base,
		       strength_pct, agility_pct, intellect_pct, stamina_pct, spirit_pct
		FROM races WHERE id = ?`, id,
	).Scan(
		&race.ID, &race.Name, &race.Faction, &race.Description,
		&race.StrengthBase, &race.AgilityBase, &race.IntellectBase, &race.StaminaBase, &race.SpiritBase,
		&race.StrengthPct, &race.AgilityPct, &race.IntellectPct, &race.StaminaPct, &race.SpiritPct,
	)
	if err != nil {
		return nil, err
	}
	return race, nil
}

// GetClasses 获取所有职业
func (r *GameRepository) GetClasses() ([]models.Class, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, role, primary_stat, resource_type,
		       base_hp, base_resource, hp_per_level, resource_per_level,
		       resource_regen, resource_regen_pct,
		       base_strength, base_agility, base_intellect, base_stamina, base_spirit,
		       base_threat_modifier, combat_role, is_ranged
		FROM classes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		class := models.Class{}
		var isRanged int
		err := rows.Scan(
			&class.ID, &class.Name, &class.Description, &class.Role, &class.PrimaryStat, &class.ResourceType,
			&class.BaseHP, &class.BaseResource, &class.HPPerLevel, &class.ResourcePerLevel,
			&class.ResourceRegen, &class.ResourceRegenPct,
			&class.BaseStrength, &class.BaseAgility, &class.BaseIntellect, &class.BaseStamina, &class.BaseSpirit,
			&class.BaseThreatModifier, &class.CombatRole, &isRanged,
		)
		if err != nil {
			return nil, err
		}
		class.IsRanged = isRanged == 1
		classes = append(classes, class)
	}

	return classes, nil
}

// GetClassByID 根据ID获取职业
func (r *GameRepository) GetClassByID(id string) (*models.Class, error) {
	class := &models.Class{}
	var isRanged int
	err := database.DB.QueryRow(`
		SELECT id, name, description, role, primary_stat, resource_type,
		       base_hp, base_resource, hp_per_level, resource_per_level,
		       resource_regen, resource_regen_pct,
		       base_strength, base_agility, base_intellect, base_stamina, base_spirit,
		       base_threat_modifier, combat_role, is_ranged
		FROM classes WHERE id = ?`, id,
	).Scan(
		&class.ID, &class.Name, &class.Description, &class.Role, &class.PrimaryStat, &class.ResourceType,
		&class.BaseHP, &class.BaseResource, &class.HPPerLevel, &class.ResourcePerLevel,
		&class.ResourceRegen, &class.ResourceRegenPct,
		&class.BaseStrength, &class.BaseAgility, &class.BaseIntellect, &class.BaseStamina, &class.BaseSpirit,
		&class.BaseThreatModifier, &class.CombatRole, &isRanged,
	)
	if err != nil {
		return nil, err
	}
	class.IsRanged = isRanged == 1
	return class, nil
}

// GetZones 获取所有区域
func (r *GameRepository) GetZones() ([]models.Zone, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, min_level, max_level, COALESCE(faction, ''), 
		       COALESCE(exp_modifier, 1.0), COALESCE(gold_modifier, 1.0),
		       unlock_zone_id, COALESCE(required_exploration, 0)
		FROM zones ORDER BY min_level`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []models.Zone
	for rows.Next() {
		zone := models.Zone{}
		var unlockZoneID sql.NullString
		err := rows.Scan(
			&zone.ID, &zone.Name, &zone.Description, &zone.MinLevel, &zone.MaxLevel,
			&zone.Faction, &zone.ExpMulti, &zone.GoldMulti,
			&unlockZoneID, &zone.RequiredExploration,
		)
		if err != nil {
			return nil, err
		}
		if unlockZoneID.Valid {
			zone.UnlockZoneID = &unlockZoneID.String
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

// GetZoneByID 根据ID获取区域
func (r *GameRepository) GetZoneByID(id string) (*models.Zone, error) {
	zone := &models.Zone{}
	var unlockZoneID sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, name, description, min_level, max_level, COALESCE(faction, ''), 
		       COALESCE(exp_modifier, 1.0), COALESCE(gold_modifier, 1.0),
		       unlock_zone_id, COALESCE(required_exploration, 0)
		FROM zones WHERE id = ?`, id,
	).Scan(
		&zone.ID, &zone.Name, &zone.Description, &zone.MinLevel, &zone.MaxLevel,
		&zone.Faction, &zone.ExpMulti, &zone.GoldMulti,
		&unlockZoneID, &zone.RequiredExploration,
	)
	if err != nil {
		return nil, err
	}
	if unlockZoneID.Valid {
		zone.UnlockZoneID = &unlockZoneID.String
	}
	return zone, nil
}

// GetMonsterSkills 获取怪物的技能列表
func (r *GameRepository) GetMonsterSkills(monsterID string) ([]*models.MonsterSkill, error) {
	rows, err := database.DB.Query(`
		SELECT id, monster_id, skill_id, skill_type, priority, cooldown, use_condition
		FROM monster_skills
		WHERE monster_id = ?
		ORDER BY priority DESC`, monsterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*models.MonsterSkill
	for rows.Next() {
		skill := &models.MonsterSkill{}
		var useCondition sql.NullString
		err := rows.Scan(
			&skill.ID, &skill.MonsterID, &skill.SkillID, &skill.SkillType,
			&skill.Priority, &skill.Cooldown, &useCondition,
		)
		if err != nil {
			return nil, err
		}
		if useCondition.Valid {
			skill.UseCondition = useCondition.String
		}
		skill.CooldownLeft = 0
		
		// 加载技能详情
		skillDetail, err := r.GetSkillByID(skill.SkillID)
		if err == nil && skillDetail != nil {
			skill.Skill = skillDetail
		}
		
		skills = append(skills, skill)
	}

	return skills, nil
}

// GetMonstersByZone 获取区域内的怪物
func (r *GameRepository) GetMonstersByZone(zoneID string) ([]models.Monster, error) {
	rows, err := database.DB.Query(`
		SELECT id, zone_id, name, level, type, hp, COALESCE(mp, 0), physical_attack, magic_attack, physical_defense, magic_defense,
		       COALESCE(attack_type, 'physical'),
		       COALESCE(phys_crit_rate, 0.05), COALESCE(phys_crit_damage, 1.5),
		       COALESCE(spell_crit_rate, 0.05), COALESCE(spell_crit_damage, 1.5),
		       COALESCE(dodge_rate, 0.05), COALESCE(speed, 10), exp_reward, gold_min, gold_max, spawn_weight,
		       COALESCE(ai_type, 'balanced'), COALESCE(ai_behavior, '')
		FROM monsters WHERE zone_id = ? ORDER BY level`, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monsters []models.Monster
	for rows.Next() {
		m := models.Monster{}
		var aiType sql.NullString
		var aiBehavior sql.NullString
		err := rows.Scan(
			&m.ID, &m.ZoneID, &m.Name, &m.Level, &m.Type, &m.HP, &m.MP, &m.PhysicalAttack, &m.MagicAttack, &m.PhysicalDefense, &m.MagicDefense,
			&m.AttackType, &m.PhysCritRate, &m.PhysCritDamage, &m.SpellCritRate, &m.SpellCritDamage, &m.DodgeRate,
			&m.Speed, &m.ExpReward, &m.GoldMin, &m.GoldMax, &m.SpawnWeight,
			&aiType, &aiBehavior,
		)
		if err != nil {
			return nil, err
		}
		if m.AttackType == "" {
			m.AttackType = "physical"
		}
		m.MaxHP = m.HP
		m.MaxMP = m.MP
		if aiType.Valid {
			m.AIType = aiType.String
		} else {
			m.AIType = "balanced"
		}
		if aiBehavior.Valid {
			m.AIBehavior = aiBehavior.String
		}
		
		// 加载怪物技能
		skills, err := r.GetMonsterSkills(m.ID)
		if err == nil {
			m.MonsterSkills = skills
		}
		
		monsters = append(monsters, m)
	}

	return monsters, nil
}

// GetMonsterByID 根据ID获取怪物
func (r *GameRepository) GetMonsterByID(monsterID string) (*models.Monster, error) {
	m := &models.Monster{}
	var aiType sql.NullString
	var aiBehavior sql.NullString
	
	err := database.DB.QueryRow(`
		SELECT id, zone_id, name, level, type, hp, COALESCE(mp, 0), physical_attack, magic_attack, physical_defense, magic_defense,
		       COALESCE(attack_type, 'physical'),
		       COALESCE(phys_crit_rate, 0.05), COALESCE(phys_crit_damage, 1.5),
		       COALESCE(spell_crit_rate, 0.05), COALESCE(spell_crit_damage, 1.5),
		       COALESCE(dodge_rate, 0.05), COALESCE(speed, 10), exp_reward, gold_min, gold_max, spawn_weight,
		       COALESCE(ai_type, 'balanced'), COALESCE(ai_behavior, '')
		FROM monsters WHERE id = ?`, monsterID,
	).Scan(
		&m.ID, &m.ZoneID, &m.Name, &m.Level, &m.Type, &m.HP, &m.MP, &m.PhysicalAttack, &m.MagicAttack, &m.PhysicalDefense, &m.MagicDefense,
		&m.AttackType, &m.PhysCritRate, &m.PhysCritDamage, &m.SpellCritRate, &m.SpellCritDamage, &m.DodgeRate,
		&m.Speed, &m.ExpReward, &m.GoldMin, &m.GoldMax, &m.SpawnWeight,
		&aiType, &aiBehavior,
	)
	if err != nil {
		return nil, err
	}
	
	if m.AttackType == "" {
		m.AttackType = "physical"
	}
	m.MaxHP = m.HP
	m.MaxMP = m.MP
	if aiType.Valid {
		m.AIType = aiType.String
	} else {
		m.AIType = "balanced"
	}
	if aiBehavior.Valid {
		m.AIBehavior = aiBehavior.String
	}
	
	// 加载怪物技能
	skills, err := r.GetMonsterSkills(m.ID)
	if err == nil {
		m.MonsterSkills = skills
	}
	
	return m, nil
}

// GetItemByID 根据ID获取物品
func (r *GameRepository) GetItemByID(itemID string) (map[string]interface{}, error) {
	var id, name, itemType, quality, slot string
	var description, subtype, classRequired sql.NullString
	var levelRequired, stackable, maxStack, sellPrice, buyPrice int
	var strength, agility, intellect, stamina, spirit int
	var attack, defense, hpBonus, mpBonus int
	var critRate float64
	var effectType sql.NullString
	var effectValue sql.NullInt64
	
	err := database.DB.QueryRow(`
		SELECT id, name, COALESCE(description, ''), type, COALESCE(subtype, ''), quality,
		       level_required, class_required, slot, stackable, max_stack,
		       sell_price, buy_price, strength, agility, intellect, stamina, spirit,
		       attack, defense, hp_bonus, mp_bonus, crit_rate, effect_type, effect_value
		FROM items WHERE id = ?`, itemID,
	).Scan(
		&id, &name, &description, &itemType, &subtype, &quality,
		&levelRequired, &classRequired, &slot, &stackable, &maxStack,
		&sellPrice, &buyPrice, &strength, &agility, &intellect, &stamina, &spirit,
		&attack, &defense, &hpBonus, &mpBonus, &critRate, &effectType, &effectValue,
	)
	if err != nil {
		return nil, err
	}
	
	item := map[string]interface{}{
		"id":            id,
		"name":          name,
		"type":          itemType,
		"quality":      quality,
		"level_required": levelRequired,
		"slot":          slot,
		"stackable":    stackable,
		"max_stack":    maxStack,
		"sell_price":   sellPrice,
		"buy_price":    buyPrice,
		"strength":     strength,
		"agility":      agility,
		"intellect":    intellect,
		"stamina":      stamina,
		"spirit":       spirit,
		"attack":       attack,
		"defense":      defense,
		"hp_bonus":     hpBonus,
		"mp_bonus":     mpBonus,
		"crit_rate":    critRate,
	}
	
	if description.Valid {
		item["description"] = description.String
	}
	if subtype.Valid {
		item["subtype"] = subtype.String
	}
	if classRequired.Valid {
		item["class_required"] = classRequired.String
	}
	if effectType.Valid {
		item["effect_type"] = effectType.String
	}
	if effectValue.Valid {
		item["effect_value"] = effectValue.Int64
	}
	
	return item, nil
}

// GetSkillByID 根据ID获取技能
func (r *GameRepository) GetSkillByID(skillID string) (*models.Skill, error) {
	skill := &models.Skill{}
	var description sql.NullString
	var damageType sql.NullString
	var scalingStat sql.NullString
	var tags sql.NullString
	
	err := database.DB.QueryRow(`
		SELECT id, name, COALESCE(description, ''), class_id, type, target_type,
		       COALESCE(damage_type, ''), base_value, COALESCE(scaling_stat, ''),
		       scaling_ratio, resource_cost, cooldown, level_required,
		       threat_modifier, threat_type, tags
		FROM skills WHERE id = ?`, skillID,
	).Scan(
		&skill.ID, &skill.Name, &description, &skill.ClassID, &skill.Type, &skill.TargetType,
		&damageType, &skill.BaseValue, &scalingStat,
		&skill.ScalingRatio, &skill.ResourceCost, &skill.Cooldown, &skill.LevelRequired,
		&skill.ThreatModifier, &skill.ThreatType, &tags,
	)
	if err != nil {
		return nil, err
	}
	
	if description.Valid {
		skill.Description = description.String
	}
	if damageType.Valid {
		skill.DamageType = damageType.String
	}
	if scalingStat.Valid {
		skill.ScalingStat = scalingStat.String
	}
	if tags.Valid {
		skill.Tags = tags.String
	}
	
	return skill, nil
}

// MonsterDrop 怪物掉落项
type MonsterDrop struct {
	ItemID     string
	DropRate   float64
	MinQuantity int
	MaxQuantity int
}

// GetMonsterDrops 获取怪物的掉落表
func (r *GameRepository) GetMonsterDrops(monsterID string) ([]MonsterDrop, error) {
	rows, err := database.DB.Query(`
		SELECT item_id, drop_rate, min_quantity, max_quantity
		FROM monster_drops
		WHERE monster_id = ?
		ORDER BY drop_rate DESC
	`, monsterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drops []MonsterDrop
	for rows.Next() {
		var drop MonsterDrop
		err := rows.Scan(&drop.ItemID, &drop.DropRate, &drop.MinQuantity, &drop.MaxQuantity)
		if err != nil {
			return nil, err
		}
		drops = append(drops, drop)
	}

	return drops, nil
}
