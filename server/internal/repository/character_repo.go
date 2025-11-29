package repository

import (
	"database/sql"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// CharacterRepository 角色数据仓库
type CharacterRepository struct{}

// NewCharacterRepository 创建角色仓库
func NewCharacterRepository() *CharacterRepository {
	return &CharacterRepository{}
}

// Create 创建角色
func (r *CharacterRepository) Create(char *models.Character) (*models.Character, error) {
	result, err := database.DB.Exec(`
		INSERT INTO characters (
			user_id, name, race_id, class_id, faction, team_slot,
			is_active, is_dead, level, exp, exp_to_next,
			hp, max_hp, resource, max_resource, resource_type,
			strength, agility, intellect, stamina, spirit,
			attack, defense, crit_rate, crit_damage,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		char.UserID, char.Name, char.RaceID, char.ClassID, char.Faction, char.TeamSlot,
		boolToInt(char.IsActive), boolToInt(char.IsDead), char.Level, char.Exp, char.ExpToNext,
		char.HP, char.MaxHP, char.Resource, char.MaxResource, char.ResourceType,
		char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit,
		char.Attack, char.Defense, char.CritRate, char.CritDamage,
		time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	char.ID = int(id)
	return char, nil
}

// GetByID 根据ID获取角色
func (r *CharacterRepository) GetByID(id int) (*models.Character, error) {
	char := &models.Character{}
	var isActive, isDead int
	var reviveAt sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, user_id, name, race_id, class_id, faction, team_slot,
		       is_active, is_dead, revive_at, level, exp, exp_to_next,
		       hp, max_hp, resource, max_resource, resource_type,
		       strength, agility, intellect, stamina, spirit,
		       attack, defense, crit_rate, crit_damage,
		       total_kills, total_deaths, created_at
		FROM characters WHERE id = ?`, id,
	).Scan(
		&char.ID, &char.UserID, &char.Name, &char.RaceID, &char.ClassID, &char.Faction, &char.TeamSlot,
		&isActive, &isDead, &reviveAt, &char.Level, &char.Exp, &char.ExpToNext,
		&char.HP, &char.MaxHP, &char.Resource, &char.MaxResource, &char.ResourceType,
		&char.Strength, &char.Agility, &char.Intellect, &char.Stamina, &char.Spirit,
		&char.Attack, &char.Defense, &char.CritRate, &char.CritDamage,
		&char.TotalKills, &char.TotalDeaths, &char.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	char.IsActive = isActive == 1
	char.IsDead = isDead == 1
	if reviveAt.Valid {
		char.ReviveAt = &reviveAt.Time
	}

	return char, nil
}

// GetByUserID 获取用户的所有角色
func (r *CharacterRepository) GetByUserID(userID int) ([]*models.Character, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, name, race_id, class_id, faction, team_slot,
		       is_active, is_dead, revive_at, level, exp, exp_to_next,
		       hp, max_hp, resource, max_resource, resource_type,
		       strength, agility, intellect, stamina, spirit,
		       attack, defense, crit_rate, crit_damage,
		       total_kills, total_deaths, created_at
		FROM characters WHERE user_id = ? ORDER BY team_slot`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		char := &models.Character{}
		var isActive, isDead int
		var reviveAt sql.NullTime

		err := rows.Scan(
			&char.ID, &char.UserID, &char.Name, &char.RaceID, &char.ClassID, &char.Faction, &char.TeamSlot,
			&isActive, &isDead, &reviveAt, &char.Level, &char.Exp, &char.ExpToNext,
			&char.HP, &char.MaxHP, &char.Resource, &char.MaxResource, &char.ResourceType,
			&char.Strength, &char.Agility, &char.Intellect, &char.Stamina, &char.Spirit,
			&char.Attack, &char.Defense, &char.CritRate, &char.CritDamage,
			&char.TotalKills, &char.TotalDeaths, &char.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		char.IsActive = isActive == 1
		char.IsDead = isDead == 1
		if reviveAt.Valid {
			char.ReviveAt = &reviveAt.Time
		}

		characters = append(characters, char)
	}

	return characters, nil
}

// GetActiveByUserID 获取用户的活跃角色（小队成员）
func (r *CharacterRepository) GetActiveByUserID(userID int) ([]*models.Character, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, name, race_id, class_id, faction, team_slot,
		       is_active, is_dead, revive_at, level, exp, exp_to_next,
		       hp, max_hp, resource, max_resource, resource_type,
		       strength, agility, intellect, stamina, spirit,
		       attack, defense, crit_rate, crit_damage,
		       total_kills, total_deaths, created_at
		FROM characters WHERE user_id = ? AND is_active = 1 ORDER BY team_slot`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character
	for rows.Next() {
		char := &models.Character{}
		var isActive, isDead int
		var reviveAt sql.NullTime

		err := rows.Scan(
			&char.ID, &char.UserID, &char.Name, &char.RaceID, &char.ClassID, &char.Faction, &char.TeamSlot,
			&isActive, &isDead, &reviveAt, &char.Level, &char.Exp, &char.ExpToNext,
			&char.HP, &char.MaxHP, &char.Resource, &char.MaxResource, &char.ResourceType,
			&char.Strength, &char.Agility, &char.Intellect, &char.Stamina, &char.Spirit,
			&char.Attack, &char.Defense, &char.CritRate, &char.CritDamage,
			&char.TotalKills, &char.TotalDeaths, &char.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		char.IsActive = isActive == 1
		char.IsDead = isDead == 1
		if reviveAt.Valid {
			char.ReviveAt = &reviveAt.Time
		}

		characters = append(characters, char)
	}

	return characters, nil
}

// CountByUserID 获取用户角色数量
func (r *CharacterRepository) CountByUserID(userID int) (int, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM characters WHERE user_id = ?`, userID,
	).Scan(&count)
	return count, err
}

// NameExists 检查角色名是否存在
func (r *CharacterRepository) NameExists(name string) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM characters WHERE name = ?`, name,
	).Scan(&count)
	return count > 0, err
}

// GetNextSlot 获取下一个可用的队伍槽位
func (r *CharacterRepository) GetNextSlot(userID int) (int, error) {
	var maxSlot sql.NullInt64
	err := database.DB.QueryRow(`
		SELECT MAX(team_slot) FROM characters WHERE user_id = ?`, userID,
	).Scan(&maxSlot)
	if err != nil {
		return 1, err
	}
	if !maxSlot.Valid {
		return 1, nil
	}
	return int(maxSlot.Int64) + 1, nil
}

// Update 更新角色
func (r *CharacterRepository) Update(char *models.Character) error {
	_, err := database.DB.Exec(`
		UPDATE characters SET
			is_active = ?, is_dead = ?, revive_at = ?,
			level = ?, exp = ?, exp_to_next = ?,
			hp = ?, max_hp = ?, resource = ?, max_resource = ?,
			strength = ?, agility = ?, intellect = ?, stamina = ?, spirit = ?,
			attack = ?, defense = ?, crit_rate = ?, crit_damage = ?,
			total_kills = ?, total_deaths = ?, updated_at = ?
		WHERE id = ?`,
		boolToInt(char.IsActive), boolToInt(char.IsDead), char.ReviveAt,
		char.Level, char.Exp, char.ExpToNext,
		char.HP, char.MaxHP, char.Resource, char.MaxResource,
		char.Strength, char.Agility, char.Intellect, char.Stamina, char.Spirit,
		char.Attack, char.Defense, char.CritRate, char.CritDamage,
		char.TotalKills, char.TotalDeaths, time.Now(),
		char.ID,
	)
	return err
}

// SetActive 设置角色激活状态
func (r *CharacterRepository) SetActive(id int, active bool) error {
	_, err := database.DB.Exec(`
		UPDATE characters SET is_active = ?, updated_at = ? WHERE id = ?`,
		boolToInt(active), time.Now(), id,
	)
	return err
}

// SetDead 设置角色死亡状态
func (r *CharacterRepository) SetDead(id int, dead bool, reviveAt *time.Time) error {
	_, err := database.DB.Exec(`
		UPDATE characters SET is_dead = ?, revive_at = ?, updated_at = ? WHERE id = ?`,
		boolToInt(dead), reviveAt, time.Now(), id,
	)
	return err
}

// Delete 删除角色
func (r *CharacterRepository) Delete(id int) error {
	_, err := database.DB.Exec(`DELETE FROM characters WHERE id = ?`, id)
	return err
}

// helper function
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

