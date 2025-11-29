package repository

import (
	"database/sql"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// UserRepository 用户数据仓库
type UserRepository struct{}

// NewUserRepository 创建用户仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(username, passwordHash, email string) (*models.User, error) {
	result, err := database.DB.Exec(`
		INSERT INTO users (username, password_hash, email, created_at) 
		VALUES (?, ?, ?, ?)`,
		username, passwordHash, email, time.Now(),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	var lastLogin sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, username, COALESCE(email, ''), max_team_size, unlocked_slots, 
		       gold, current_zone_id, total_kills, total_gold_gained, play_time,
		       created_at, last_login_at
		FROM users WHERE id = ?`, id,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.MaxTeamSize, &user.UnlockedSlots,
		&user.Gold, &user.CurrentZoneID, &user.TotalKills, &user.TotalGoldGained, &user.PlayTime,
		&user.CreatedAt, &lastLogin,
	)
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	var lastLogin sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, username, COALESCE(email, ''), max_team_size, unlocked_slots,
		       gold, current_zone_id, total_kills, total_gold_gained, play_time,
		       created_at, last_login_at
		FROM users WHERE username = ?`, username,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.MaxTeamSize, &user.UnlockedSlots,
		&user.Gold, &user.CurrentZoneID, &user.TotalKills, &user.TotalGoldGained, &user.PlayTime,
		&user.CreatedAt, &lastLogin,
	)
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	return user, nil
}

// GetPasswordHash 获取用户密码哈希
func (r *UserRepository) GetPasswordHash(username string) (int, string, error) {
	var id int
	var hash string
	err := database.DB.QueryRow(`
		SELECT id, password_hash FROM users WHERE username = ?`, username,
	).Scan(&id, &hash)
	return id, hash, err
}

// UsernameExists 检查用户名是否存在
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE username = ?`, username,
	).Scan(&count)
	return count > 0, err
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(id int) error {
	_, err := database.DB.Exec(`
		UPDATE users SET last_login_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

// UpdateGold 更新金币
func (r *UserRepository) UpdateGold(id int, gold int) error {
	_, err := database.DB.Exec(`
		UPDATE users SET gold = gold + ?, total_gold_gained = total_gold_gained + ? WHERE id = ?`,
		gold, gold, id,
	)
	return err
}

// UpdateZone 更新当前区域
func (r *UserRepository) UpdateZone(id int, zoneID string) error {
	_, err := database.DB.Exec(`
		UPDATE users SET current_zone_id = ? WHERE id = ?`,
		zoneID, id,
	)
	return err
}

// IncrementKills 增加击杀数
func (r *UserRepository) IncrementKills(id int) error {
	_, err := database.DB.Exec(`
		UPDATE users SET total_kills = total_kills + 1 WHERE id = ?`, id,
	)
	return err
}

