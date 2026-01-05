package repository

import (
	"database/sql"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// EquipmentRepository 装备数据仓库
type EquipmentRepository struct{}

// NewEquipmentRepository 创建装备仓库
func NewEquipmentRepository() *EquipmentRepository {
	return &EquipmentRepository{}
}

// Create 创建装备实例
func (r *EquipmentRepository) Create(equipment *models.EquipmentInstance) (*models.EquipmentInstance, error) {
	result, err := database.DB.Exec(`
		INSERT INTO equipment_instance (
			item_id, owner_id, character_id, slot, quality,
			evolution_stage, evolution_path,
			prefix_id, prefix_value, suffix_id, suffix_value,
			bonus_affix_1, bonus_affix_1_value,
			bonus_affix_2, bonus_affix_2_value,
			legendary_effect_id, acquired_at, is_locked
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		equipment.ItemID, equipment.OwnerID, equipment.CharacterID, equipment.Slot, equipment.Quality,
		equipment.EvolutionStage, equipment.EvolutionPath,
		equipment.PrefixID, equipment.PrefixValue, equipment.SuffixID, equipment.SuffixValue,
		equipment.BonusAffix1, equipment.BonusAffix1Value,
		equipment.BonusAffix2, equipment.BonusAffix2Value,
		equipment.LegendaryEffectID, time.Now(), boolToInt(equipment.IsLocked),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	equipment.ID = int(id)
	return equipment, nil
}

// GetByID 根据ID获取装备
func (r *EquipmentRepository) GetByID(id int) (*models.EquipmentInstance, error) {
	equipment := &models.EquipmentInstance{}
	var characterID sql.NullInt64
	var evolutionPath, prefixID, suffixID, bonusAffix1, bonusAffix2, legendaryEffectID sql.NullString
	var prefixValue, suffixValue, bonusAffix1Value, bonusAffix2Value sql.NullFloat64
	var isLocked int

	err := database.DB.QueryRow(`
		SELECT id, item_id, owner_id, character_id, slot, quality,
		       evolution_stage, evolution_path,
		       prefix_id, prefix_value, suffix_id, suffix_value,
		       bonus_affix_1, bonus_affix_1_value,
		       bonus_affix_2, bonus_affix_2_value,
		       legendary_effect_id, acquired_at, is_locked
		FROM equipment_instance WHERE id = ?`, id,
	).Scan(
		&equipment.ID, &equipment.ItemID, &equipment.OwnerID, &characterID, &equipment.Slot, &equipment.Quality,
		&equipment.EvolutionStage, &evolutionPath,
		&prefixID, &prefixValue, &suffixID, &suffixValue,
		&bonusAffix1, &bonusAffix1Value,
		&bonusAffix2, &bonusAffix2Value,
		&legendaryEffectID, &equipment.AcquiredAt, &isLocked,
	)
	if err != nil {
		return nil, err
	}

	if characterID.Valid {
		id := int(characterID.Int64)
		equipment.CharacterID = &id
	}
	if evolutionPath.Valid {
		equipment.EvolutionPath = &evolutionPath.String
	}
	if prefixID.Valid {
		equipment.PrefixID = &prefixID.String
	}
	if prefixValue.Valid {
		equipment.PrefixValue = &prefixValue.Float64
	}
	if suffixID.Valid {
		equipment.SuffixID = &suffixID.String
	}
	if suffixValue.Valid {
		equipment.SuffixValue = &suffixValue.Float64
	}
	if bonusAffix1.Valid {
		equipment.BonusAffix1 = &bonusAffix1.String
	}
	if bonusAffix1Value.Valid {
		equipment.BonusAffix1Value = &bonusAffix1Value.Float64
	}
	if bonusAffix2.Valid {
		equipment.BonusAffix2 = &bonusAffix2.String
	}
	if bonusAffix2Value.Valid {
		equipment.BonusAffix2Value = &bonusAffix2Value.Float64
	}
	if legendaryEffectID.Valid {
		equipment.LegendaryEffectID = &legendaryEffectID.String
	}
	equipment.IsLocked = intToBool(isLocked)

	return equipment, nil
}

// GetByOwnerID 获取用户的所有装备
func (r *EquipmentRepository) GetByOwnerID(ownerID int) ([]*models.EquipmentInstance, error) {
	rows, err := database.DB.Query(`
		SELECT id, item_id, owner_id, character_id, slot, quality,
		       evolution_stage, evolution_path,
		       prefix_id, prefix_value, suffix_id, suffix_value,
		       bonus_affix_1, bonus_affix_1_value,
		       bonus_affix_2, bonus_affix_2_value,
		       legendary_effect_id, acquired_at, is_locked
		FROM equipment_instance WHERE owner_id = ?`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []*models.EquipmentInstance
	for rows.Next() {
		equipment := &models.EquipmentInstance{}
		var characterID sql.NullInt64
		var evolutionPath, prefixID, suffixID, bonusAffix1, bonusAffix2, legendaryEffectID sql.NullString
		var prefixValue, suffixValue, bonusAffix1Value, bonusAffix2Value sql.NullFloat64
		var isLocked int

		err := rows.Scan(
			&equipment.ID, &equipment.ItemID, &equipment.OwnerID, &characterID, &equipment.Slot, &equipment.Quality,
			&equipment.EvolutionStage, &evolutionPath,
			&prefixID, &prefixValue, &suffixID, &suffixValue,
			&bonusAffix1, &bonusAffix1Value,
			&bonusAffix2, &bonusAffix2Value,
			&legendaryEffectID, &equipment.AcquiredAt, &isLocked,
		)
		if err != nil {
			return nil, err
		}

		if characterID.Valid {
			id := int(characterID.Int64)
			equipment.CharacterID = &id
		}
		if evolutionPath.Valid {
			equipment.EvolutionPath = &evolutionPath.String
		}
		if prefixID.Valid {
			equipment.PrefixID = &prefixID.String
		}
		if prefixValue.Valid {
			equipment.PrefixValue = &prefixValue.Float64
		}
		if suffixID.Valid {
			equipment.SuffixID = &suffixID.String
		}
		if suffixValue.Valid {
			equipment.SuffixValue = &suffixValue.Float64
		}
		if bonusAffix1.Valid {
			equipment.BonusAffix1 = &bonusAffix1.String
		}
		if bonusAffix1Value.Valid {
			equipment.BonusAffix1Value = &bonusAffix1Value.Float64
		}
		if bonusAffix2.Valid {
			equipment.BonusAffix2 = &bonusAffix2.String
		}
		if bonusAffix2Value.Valid {
			equipment.BonusAffix2Value = &bonusAffix2Value.Float64
		}
		if legendaryEffectID.Valid {
			equipment.LegendaryEffectID = &legendaryEffectID.String
		}
		equipment.IsLocked = intToBool(isLocked)

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// GetByCharacterID 获取角色装备的所有装备
func (r *EquipmentRepository) GetByCharacterID(characterID int) ([]*models.EquipmentInstance, error) {
	rows, err := database.DB.Query(`
		SELECT id, item_id, owner_id, character_id, slot, quality,
		       evolution_stage, evolution_path,
		       prefix_id, prefix_value, suffix_id, suffix_value,
		       bonus_affix_1, bonus_affix_1_value,
		       bonus_affix_2, bonus_affix_2_value,
		       legendary_effect_id, acquired_at, is_locked
		FROM equipment_instance WHERE character_id = ?`, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []*models.EquipmentInstance
	for rows.Next() {
		equipment := &models.EquipmentInstance{}
		var characterID sql.NullInt64
		var evolutionPath, prefixID, suffixID, bonusAffix1, bonusAffix2, legendaryEffectID sql.NullString
		var prefixValue, suffixValue, bonusAffix1Value, bonusAffix2Value sql.NullFloat64
		var isLocked int

		err := rows.Scan(
			&equipment.ID, &equipment.ItemID, &equipment.OwnerID, &characterID, &equipment.Slot, &equipment.Quality,
			&equipment.EvolutionStage, &evolutionPath,
			&prefixID, &prefixValue, &suffixID, &suffixValue,
			&bonusAffix1, &bonusAffix1Value,
			&bonusAffix2, &bonusAffix2Value,
			&legendaryEffectID, &equipment.AcquiredAt, &isLocked,
		)
		if err != nil {
			return nil, err
		}

		if characterID.Valid {
			id := int(characterID.Int64)
			equipment.CharacterID = &id
		}
		if evolutionPath.Valid {
			equipment.EvolutionPath = &evolutionPath.String
		}
		if prefixID.Valid {
			equipment.PrefixID = &prefixID.String
		}
		if prefixValue.Valid {
			equipment.PrefixValue = &prefixValue.Float64
		}
		if suffixID.Valid {
			equipment.SuffixID = &suffixID.String
		}
		if suffixValue.Valid {
			equipment.SuffixValue = &suffixValue.Float64
		}
		if bonusAffix1.Valid {
			equipment.BonusAffix1 = &bonusAffix1.String
		}
		if bonusAffix1Value.Valid {
			equipment.BonusAffix1Value = &bonusAffix1Value.Float64
		}
		if bonusAffix2.Valid {
			equipment.BonusAffix2 = &bonusAffix2.String
		}
		if bonusAffix2Value.Valid {
			equipment.BonusAffix2Value = &bonusAffix2Value.Float64
		}
		if legendaryEffectID.Valid {
			equipment.LegendaryEffectID = &legendaryEffectID.String
		}
		equipment.IsLocked = intToBool(isLocked)

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// GetByCharacterAndSlot 获取角色指定槽位的装备
func (r *EquipmentRepository) GetByCharacterAndSlot(characterID int, slot string) (*models.EquipmentInstance, error) {
	equipment := &models.EquipmentInstance{}
	var characterIDVal sql.NullInt64
	var evolutionPath, prefixID, suffixID, bonusAffix1, bonusAffix2, legendaryEffectID sql.NullString
	var prefixValue, suffixValue, bonusAffix1Value, bonusAffix2Value sql.NullFloat64
	var isLocked int

	err := database.DB.QueryRow(`
		SELECT id, item_id, owner_id, character_id, slot, quality,
		       evolution_stage, evolution_path,
		       prefix_id, prefix_value, suffix_id, suffix_value,
		       bonus_affix_1, bonus_affix_1_value,
		       bonus_affix_2, bonus_affix_2_value,
		       legendary_effect_id, acquired_at, is_locked
		FROM equipment_instance WHERE character_id = ? AND slot = ?`, characterID, slot,
	).Scan(
		&equipment.ID, &equipment.ItemID, &equipment.OwnerID, &characterIDVal, &equipment.Slot, &equipment.Quality,
		&equipment.EvolutionStage, &evolutionPath,
		&prefixID, &prefixValue, &suffixID, &suffixValue,
		&bonusAffix1, &bonusAffix1Value,
		&bonusAffix2, &bonusAffix2Value,
		&legendaryEffectID, &equipment.AcquiredAt, &isLocked,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if characterIDVal.Valid {
		id := int(characterIDVal.Int64)
		equipment.CharacterID = &id
	}
	if evolutionPath.Valid {
		equipment.EvolutionPath = &evolutionPath.String
	}
	if prefixID.Valid {
		equipment.PrefixID = &prefixID.String
	}
	if prefixValue.Valid {
		equipment.PrefixValue = &prefixValue.Float64
	}
	if suffixID.Valid {
		equipment.SuffixID = &suffixID.String
	}
	if suffixValue.Valid {
		equipment.SuffixValue = &suffixValue.Float64
	}
	if bonusAffix1.Valid {
		equipment.BonusAffix1 = &bonusAffix1.String
	}
	if bonusAffix1Value.Valid {
		equipment.BonusAffix1Value = &bonusAffix1Value.Float64
	}
	if bonusAffix2.Valid {
		equipment.BonusAffix2 = &bonusAffix2.String
	}
	if bonusAffix2Value.Valid {
		equipment.BonusAffix2Value = &bonusAffix2Value.Float64
	}
	if legendaryEffectID.Valid {
		equipment.LegendaryEffectID = &legendaryEffectID.String
	}
	equipment.IsLocked = intToBool(isLocked)

	return equipment, nil
}

// Update 更新装备
func (r *EquipmentRepository) Update(equipment *models.EquipmentInstance) error {
	_, err := database.DB.Exec(`
		UPDATE equipment_instance SET
			item_id = ?, owner_id = ?, character_id = ?, slot = ?, quality = ?,
			evolution_stage = ?, evolution_path = ?,
			prefix_id = ?, prefix_value = ?, suffix_id = ?, suffix_value = ?,
			bonus_affix_1 = ?, bonus_affix_1_value = ?,
			bonus_affix_2 = ?, bonus_affix_2_value = ?,
			legendary_effect_id = ?, is_locked = ?
		WHERE id = ?`,
		equipment.ItemID, equipment.OwnerID, equipment.CharacterID, equipment.Slot, equipment.Quality,
		equipment.EvolutionStage, equipment.EvolutionPath,
		equipment.PrefixID, equipment.PrefixValue, equipment.SuffixID, equipment.SuffixValue,
		equipment.BonusAffix1, equipment.BonusAffix1Value,
		equipment.BonusAffix2, equipment.BonusAffix2Value,
		equipment.LegendaryEffectID, boolToInt(equipment.IsLocked),
		equipment.ID,
	)
	return err
}

// Delete 删除装备
func (r *EquipmentRepository) Delete(id int) error {
	_, err := database.DB.Exec(`DELETE FROM equipment_instance WHERE id = ?`, id)
	return err
}

