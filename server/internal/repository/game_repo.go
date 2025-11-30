package repository

import (
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
		       COALESCE(exp_modifier, 1.0), COALESCE(gold_modifier, 1.0)
		FROM zones ORDER BY min_level`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []models.Zone
	for rows.Next() {
		zone := models.Zone{}
		err := rows.Scan(
			&zone.ID, &zone.Name, &zone.Description, &zone.MinLevel, &zone.MaxLevel,
			&zone.Faction, &zone.ExpMulti, &zone.GoldMulti,
		)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

// GetZoneByID 根据ID获取区域
func (r *GameRepository) GetZoneByID(id string) (*models.Zone, error) {
	zone := &models.Zone{}
	err := database.DB.QueryRow(`
		SELECT id, name, description, min_level, max_level, COALESCE(faction, ''), 
		       COALESCE(exp_modifier, 1.0), COALESCE(gold_modifier, 1.0)
		FROM zones WHERE id = ?`, id,
	).Scan(
		&zone.ID, &zone.Name, &zone.Description, &zone.MinLevel, &zone.MaxLevel,
		&zone.Faction, &zone.ExpMulti, &zone.GoldMulti,
	)
	if err != nil {
		return nil, err
	}
	return zone, nil
}

// GetMonstersByZone 获取区域内的怪物
func (r *GameRepository) GetMonstersByZone(zoneID string) ([]models.Monster, error) {
	rows, err := database.DB.Query(`
		SELECT id, zone_id, name, level, type, hp, attack, defense, exp_reward, gold_min, gold_max, spawn_weight
		FROM monsters WHERE zone_id = ? ORDER BY level`, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monsters []models.Monster
	for rows.Next() {
		m := models.Monster{}
		err := rows.Scan(
			&m.ID, &m.ZoneID, &m.Name, &m.Level, &m.Type, &m.HP, &m.Attack, &m.Defense,
			&m.ExpReward, &m.GoldMin, &m.GoldMax, &m.SpawnWeight,
		)
		if err != nil {
			return nil, err
		}
		m.MaxHP = m.HP
		monsters = append(monsters, m)
	}

	return monsters, nil
}

