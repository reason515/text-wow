package repository

import (
	"database/sql"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/models"
)

// ExplorationRepository 探索度数据仓库
type ExplorationRepository struct{}

// NewExplorationRepository 创建探索度仓库
func NewExplorationRepository() *ExplorationRepository {
	return &ExplorationRepository{}
}

// GetExploration 获取玩家在指定地图的探索度
func (r *ExplorationRepository) GetExploration(userID int, zoneID string) (*models.ZoneExploration, error) {
	exploration := &models.ZoneExploration{
		UserID: userID,
		ZoneID: zoneID,
	}

	err := database.DB.QueryRow(`
		SELECT exploration, kills
		FROM user_zone_exploration
		WHERE user_id = ? AND zone_id = ?
	`, userID, zoneID).Scan(&exploration.Exploration, &exploration.Kills)

	if err == sql.ErrNoRows {
		// 如果没有记录，返回默认值
		return exploration, nil
	}
	if err != nil {
		return nil, err
	}

	return exploration, nil
}

// AddExploration 增加玩家在指定地图的探索度
func (r *ExplorationRepository) AddExploration(userID int, zoneID string, explorationGain int) error {
	// 使用 INSERT OR REPLACE 语法（SQLite 兼容）
	_, err := database.DB.Exec(`
		INSERT OR REPLACE INTO user_zone_exploration (user_id, zone_id, exploration, kills, last_updated)
		VALUES (
			?,
			?,
			COALESCE((SELECT exploration FROM user_zone_exploration WHERE user_id = ? AND zone_id = ?), 0) + ?,
			COALESCE((SELECT kills FROM user_zone_exploration WHERE user_id = ? AND zone_id = ?), 0) + 1,
			?
		)
	`, userID, zoneID, userID, zoneID, explorationGain, userID, zoneID, time.Now())

	return err
}

// GetAllExplorations 获取玩家所有地图的探索度
func (r *ExplorationRepository) GetAllExplorations(userID int) (map[string]*models.ZoneExploration, error) {
	rows, err := database.DB.Query(`
		SELECT zone_id, exploration, kills
		FROM user_zone_exploration
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	explorations := make(map[string]*models.ZoneExploration)
	for rows.Next() {
		exp := &models.ZoneExploration{
			UserID: userID,
		}
		err := rows.Scan(&exp.ZoneID, &exp.Exploration, &exp.Kills)
		if err != nil {
			return nil, err
		}
		explorations[exp.ZoneID] = exp
	}

	return explorations, nil
}

// IsZoneUnlocked 检查地图是否已解锁（基于探索度）
func (r *ExplorationRepository) IsZoneUnlocked(userID int, zone *models.Zone) (bool, error) {
	// 如果没有解锁条件，直接返回true
	if zone.UnlockZoneID == nil || zone.RequiredExploration == 0 {
		return true, nil
	}

	// 检查前置地图的探索度
	exploration, err := r.GetExploration(userID, *zone.UnlockZoneID)
	if err != nil {
		return false, err
	}

	return exploration.Exploration >= zone.RequiredExploration, nil
}

