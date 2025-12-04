package repository

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// SkillRepository 技能数据仓库
type SkillRepository struct{}

// NewSkillRepository 创建技能仓库
func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

// GetSkillByID 根据ID获取技能
func (r *SkillRepository) GetSkillByID(skillID string) (*models.Skill, error) {
	skill := &models.Skill{}
	var tags sql.NullString

	var effectID sql.NullString
	var effectChance sql.NullFloat64
	var icon sql.NullString
	var damageType sql.NullString
	var scalingStat sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, name, description, icon, class_id, type, target_type, damage_type,
		       base_value, scaling_stat, scaling_ratio, resource_cost, cooldown,
		       level_required, effect_id, effect_chance, tags, threat_modifier, threat_type
		FROM skills WHERE id = ?`, skillID,
	).Scan(
		&skill.ID, &skill.Name, &skill.Description, &icon, &skill.ClassID, &skill.Type,
		&skill.TargetType, &damageType, &skill.BaseValue, &scalingStat,
		&skill.ScalingRatio, &skill.ResourceCost, &skill.Cooldown, &skill.LevelRequired,
		&effectID, &effectChance, &tags, &skill.ThreatModifier, &skill.ThreatType,
	)
	if err != nil {
		return nil, err
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

// GetInitialSkills 获取初始技能池（战士）
func (r *SkillRepository) GetInitialSkills(classID string) ([]*models.Skill, error) {
	// 战士的初始技能池ID列表
	initialSkillIDs := []string{
		"warrior_heroic_strike",
		"warrior_taunt",
		"warrior_shield_block",
		"warrior_cleave",
		"warrior_slam",
		"warrior_battle_shout",
		"warrior_demoralizing_shout",
		"warrior_last_stand",
		"warrior_charge",
	}

	skills := make([]*models.Skill, 0)
	missingSkills := make([]string, 0)
	
	for _, skillID := range initialSkillIDs {
		skill, err := r.GetSkillByID(skillID)
		if err != nil {
			missingSkills = append(missingSkills, skillID)
			continue // 跳过不存在的技能
		}
		if skill.ClassID == classID || classID == "" {
			skills = append(skills, skill)
		}
	}
	

	return skills, nil
}

// GetAllActiveSkills 获取所有主动技能（战士）
func (r *SkillRepository) GetAllActiveSkills(classID string) ([]*models.Skill, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, icon, class_id, type, target_type, damage_type,
		       base_value, scaling_stat, scaling_ratio, resource_cost, cooldown,
		       level_required, effect_id, effect_chance, tags, threat_modifier, threat_type
		FROM skills WHERE class_id = ? AND id LIKE 'warrior_%'`, classID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills := make([]*models.Skill, 0)
	for rows.Next() {
		skill := &models.Skill{}
		var tags sql.NullString
		var effectID sql.NullString
		var effectChance sql.NullFloat64
		var icon sql.NullString
		var damageType sql.NullString
		var scalingStat sql.NullString

		err := rows.Scan(
			&skill.ID, &skill.Name, &skill.Description, &icon, &skill.ClassID, &skill.Type,
			&skill.TargetType, &damageType, &skill.BaseValue, &scalingStat,
			&skill.ScalingRatio, &skill.ResourceCost, &skill.Cooldown, &skill.LevelRequired,
			&effectID, &effectChance, &tags, &skill.ThreatModifier, &skill.ThreatType,
		)
		
		if damageType.Valid {
			skill.DamageType = damageType.String
		}
		if scalingStat.Valid {
			skill.ScalingStat = scalingStat.String
		}
		if err != nil {
			continue
		}

		if tags.Valid {
			skill.Tags = tags.String
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// GetRandomActiveSkills 随机获取N个主动技能
func (r *SkillRepository) GetRandomActiveSkills(classID string, count int, excludeIDs []string) ([]*models.Skill, error) {
	allSkills, err := r.GetAllActiveSkills(classID)
	if err != nil {
		return nil, err
	}

	// 排除已拥有的技能
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	availableSkills := make([]*models.Skill, 0)
	for _, skill := range allSkills {
		if !excludeMap[skill.ID] {
			availableSkills = append(availableSkills, skill)
		}
	}

	// 如果可用技能少于请求数量，返回所有可用技能
	if len(availableSkills) <= count {
		return availableSkills, nil
	}

	// 随机选择
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(availableSkills), func(i, j int) {
		availableSkills[i], availableSkills[j] = availableSkills[j], availableSkills[i]
	})

	return availableSkills[:count], nil
}

// GetCharacterSkills 获取角色的所有主动技能
func (r *SkillRepository) GetCharacterSkills(characterID int) ([]*models.CharacterSkill, error) {
	rows, err := database.DB.Query(`
		SELECT cs.id, cs.character_id, cs.skill_id, cs.skill_level, cs.slot, cs.is_auto
		FROM character_skills cs
		WHERE cs.character_id = ?`, characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills := make([]*models.CharacterSkill, 0)
	for rows.Next() {
		cs := &models.CharacterSkill{}
		var slot sql.NullInt64
		var isAuto int

		err := rows.Scan(
			&cs.ID, &cs.CharacterID, &cs.SkillID, &cs.SkillLevel, &slot, &isAuto,
		)
		if err != nil {
			continue
		}

		if slot.Valid {
			s := int(slot.Int64)
			cs.Slot = &s
		}
		cs.IsAuto = isAuto == 1

		// 加载技能详情
		skill, err := r.GetSkillByID(cs.SkillID)
		if err == nil {
			cs.Skill = skill
		}

		skills = append(skills, cs)
	}

	return skills, nil
}

// AddCharacterSkill 添加角色技能
func (r *SkillRepository) AddCharacterSkill(characterID int, skillID string, level int) error {
	_, err := database.DB.Exec(`
		INSERT OR REPLACE INTO character_skills (character_id, skill_id, skill_level)
		VALUES (?, ?, ?)`, characterID, skillID, level,
	)
	return err
}

// UpgradeCharacterSkill 升级角色技能
func (r *SkillRepository) UpgradeCharacterSkill(characterID int, skillID string) error {
	_, err := database.DB.Exec(`
		UPDATE character_skills
		SET skill_level = skill_level + 1
		WHERE character_id = ? AND skill_id = ? AND skill_level < 5`, characterID, skillID,
	)
	return err
}

// GetPassiveSkillByID 根据ID获取被动技能
func (r *SkillRepository) GetPassiveSkillByID(passiveID string) (*models.PassiveSkill, error) {
	passive := &models.PassiveSkill{}

	err := database.DB.QueryRow(`
		SELECT id, name, description, class_id, rarity, tier, effect_type,
		       effect_value, effect_stat, max_level, level_scaling
		FROM passive_skills WHERE id = ?`, passiveID,
	).Scan(
		&passive.ID, &passive.Name, &passive.Description, &passive.ClassID,
		&passive.Rarity, &passive.Tier, &passive.EffectType, &passive.EffectValue,
		&passive.EffectStat, &passive.MaxLevel, &passive.LevelScaling,
	)
	if err != nil {
		return nil, err
	}

	return passive, nil
}

// GetAllPassiveSkills 获取所有被动技能（战士）
func (r *SkillRepository) GetAllPassiveSkills(classID string) ([]*models.PassiveSkill, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, class_id, rarity, tier, effect_type,
		       effect_value, effect_stat, max_level, level_scaling
		FROM passive_skills WHERE class_id = ? AND id LIKE 'warrior_passive_%'`, classID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passives := make([]*models.PassiveSkill, 0)
	for rows.Next() {
		passive := &models.PassiveSkill{}

		err := rows.Scan(
			&passive.ID, &passive.Name, &passive.Description, &passive.ClassID,
			&passive.Rarity, &passive.Tier, &passive.EffectType, &passive.EffectValue,
			&passive.EffectStat, &passive.MaxLevel, &passive.LevelScaling,
		)
		if err != nil {
			continue
		}

		passives = append(passives, passive)
	}

	return passives, nil
}

// GetRandomPassiveSkills 随机获取N个被动技能
func (r *SkillRepository) GetRandomPassiveSkills(classID string, count int, excludeIDs []string) ([]*models.PassiveSkill, error) {
	allPassives, err := r.GetAllPassiveSkills(classID)
	if err != nil {
		return nil, err
	}

	// 排除已拥有的被动技能
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	availablePassives := make([]*models.PassiveSkill, 0)
	for _, passive := range allPassives {
		if !excludeMap[passive.ID] {
			availablePassives = append(availablePassives, passive)
		}
	}

	// 如果可用被动技能少于请求数量，返回所有可用被动技能
	if len(availablePassives) <= count {
		return availablePassives, nil
	}

	// 随机选择
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(availablePassives), func(i, j int) {
		availablePassives[i], availablePassives[j] = availablePassives[j], availablePassives[i]
	})

	return availablePassives[:count], nil
}

// GetCharacterPassiveSkills 获取角色的所有被动技能
func (r *SkillRepository) GetCharacterPassiveSkills(characterID int) ([]*models.CharacterPassiveSkill, error) {
	rows, err := database.DB.Query(`
		SELECT id, character_id, passive_id, level, acquired_at
		FROM character_passive_skills
		WHERE character_id = ?`, characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passives := make([]*models.CharacterPassiveSkill, 0)
	for rows.Next() {
		cps := &models.CharacterPassiveSkill{}

		err := rows.Scan(
			&cps.ID, &cps.CharacterID, &cps.PassiveID, &cps.Level, &cps.AcquiredAt,
		)
		if err != nil {
			continue
		}

		// 加载被动技能详情
		passive, err := r.GetPassiveSkillByID(cps.PassiveID)
		if err == nil {
			cps.Passive = passive
		}

		passives = append(passives, cps)
	}

	return passives, nil
}

// AddCharacterPassiveSkill 添加角色被动技能
func (r *SkillRepository) AddCharacterPassiveSkill(characterID int, passiveID string, level int) error {
	_, err := database.DB.Exec(`
		INSERT OR REPLACE INTO character_passive_skills (character_id, passive_id, level)
		VALUES (?, ?, ?)`, characterID, passiveID, level,
	)
	return err
}

// UpgradeCharacterPassiveSkill 升级角色被动技能
func (r *SkillRepository) UpgradeCharacterPassiveSkill(characterID int, passiveID string) error {
	_, err := database.DB.Exec(`
		UPDATE character_passive_skills
		SET level = level + 1
		WHERE character_id = ? AND passive_id = ? AND level < 5`, characterID, passiveID,
	)
	return err
}

// RecordSkillSelection 记录技能选择历史
func (r *SkillRepository) RecordSkillSelection(characterID, levelMilestone int, selectedSkillID string, isUpgrade bool) error {
	_, err := database.DB.Exec(`
		INSERT INTO skill_selection_history 
		(character_id, level_milestone, offered_skill_1, offered_skill_2, offered_skill_3, selected_skill_id, skill_was_upgrade)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		characterID, levelMilestone, selectedSkillID, "", "", selectedSkillID, boolToInt(isUpgrade),
	)
	return err
}

// Helper function to parse JSON tags
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	var tags []string
	json.Unmarshal([]byte(tagsStr), &tags)
	return tags
}

