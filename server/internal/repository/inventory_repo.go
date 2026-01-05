package repository

import (
	"database/sql"
	"fmt"
	"text-wow/internal/database"
)

// InventoryRepository 背包数据仓库
type InventoryRepository struct{}

// NewInventoryRepository 创建背包仓库
func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{}
}

// AddItem 添加物品到背包
// 如果物品可堆叠且已存在，则增加数量；否则创建新记录
// 使用事务确保操作的原子性
func (r *InventoryRepository) AddItem(characterID int, itemID string, quantity int) error {
	return WithTransaction(func(tx *sql.Tx) error {
		// 检查物品是否可堆叠
		var stackable, maxStack int
		err := tx.QueryRow(`
			SELECT COALESCE(stackable, 0), COALESCE(max_stack, 1)
			FROM items WHERE id = ?`, itemID,
		).Scan(&stackable, &maxStack)
		if err != nil {
			// 如果物品不存在，仍然尝试添加（可能是装备类型）
			stackable = 0
			maxStack = 1
		}

		// 如果可堆叠，尝试更新现有记录
		if stackable > 0 {
			var existingID, existingQuantity int
			err := tx.QueryRow(`
				SELECT id, quantity FROM inventory
				WHERE character_id = ? AND item_id = ?
				LIMIT 1`, characterID, itemID,
			).Scan(&existingID, &existingQuantity)
			
			if err == nil {
				// 物品已存在，更新数量
				newQuantity := existingQuantity + quantity
				if newQuantity > maxStack {
					newQuantity = maxStack
				}
				_, err = tx.Exec(`
					UPDATE inventory SET quantity = ?
					WHERE id = ?`, newQuantity, existingID,
				)
				return err
			}
			// 如果查询失败（物品不存在），继续创建新记录
		}

		// 创建新记录
		// 找到下一个可用的槽位
		var nextSlot sql.NullInt64
		err = tx.QueryRow(`
			SELECT MAX(slot) FROM inventory WHERE character_id = ?`, characterID,
		).Scan(&nextSlot)
		
		slot := 1
		if err == nil && nextSlot.Valid {
			slot = int(nextSlot.Int64) + 1
		}

		_, err = tx.Exec(`
			INSERT INTO inventory (character_id, item_id, quantity, slot)
			VALUES (?, ?, ?, ?)`, characterID, itemID, quantity, slot,
		)
		return err
	})
}

// GetByCharacterID 获取角色的所有背包物品
func (r *InventoryRepository) GetByCharacterID(characterID int) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT i.id, i.item_id, i.quantity, i.slot,
		       it.name, it.type, it.quality, it.description
		FROM inventory i
		JOIN items it ON i.item_id = it.id
		WHERE i.character_id = ?
		ORDER BY i.slot`, characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, quantity, slot int
		var itemID, name, itemType, quality string
		var description sql.NullString

		err := rows.Scan(&id, &itemID, &quantity, &slot, &name, &itemType, &quality, &description)
		if err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":       id,
			"item_id":  itemID,
			"quantity": quantity,
			"slot":     slot,
			"name":     name,
			"type":     itemType,
			"quality":  quality,
		}
		if description.Valid {
			item["description"] = description.String
		}
		items = append(items, item)
	}

	return items, nil
}

// GetInventoryCapacity 获取背包容量（当前物品数量）
func (r *InventoryRepository) GetInventoryCapacity(characterID int) (int, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM inventory WHERE character_id = ?`, characterID,
	).Scan(&count)
	return count, err
}

// GetInventoryMaxCapacity 获取背包最大容量（默认100）
func (r *InventoryRepository) GetInventoryMaxCapacity(characterID int) int {
	// 默认容量100，未来可以从用户配置或角色配置中获取
	return 100
}

// IsInventoryFull 检查背包是否已满
func (r *InventoryRepository) IsInventoryFull(characterID int) (bool, error) {
	current, err := r.GetInventoryCapacity(characterID)
	if err != nil {
		return false, err
	}
	max := r.GetInventoryMaxCapacity(characterID)
	return current >= max, nil
}

// SortInventory 排序背包物品（按品质、类型、名称等）
func (r *InventoryRepository) SortInventory(characterID int, sortBy string) ([]map[string]interface{}, error) {
	var orderBy string
	switch sortBy {
	case "quality":
		orderBy = "it.quality DESC, i.slot"
	case "type":
		orderBy = "it.type, it.quality DESC, i.slot"
	case "name":
		orderBy = "it.name, i.slot"
	default:
		orderBy = "i.slot"
	}

	rows, err := database.DB.Query(fmt.Sprintf(`
		SELECT i.id, i.item_id, i.quantity, i.slot,
		       it.name, it.type, it.quality, it.description
		FROM inventory i
		JOIN items it ON i.item_id = it.id
		WHERE i.character_id = ?
		ORDER BY %s`, orderBy), characterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, quantity, slot int
		var itemID, name, itemType, quality string
		var description sql.NullString

		err := rows.Scan(&id, &itemID, &quantity, &slot, &name, &itemType, &quality, &description)
		if err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":       id,
			"item_id":  itemID,
			"quantity": quantity,
			"slot":     slot,
			"name":     name,
			"type":     itemType,
			"quality":  quality,
		}
		if description.Valid {
			item["description"] = description.String
		}
		items = append(items, item)
	}

	return items, nil
}

// FilterInventory 筛选背包物品（按类型、品质等）
func (r *InventoryRepository) FilterInventory(characterID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	query := `
		SELECT i.id, i.item_id, i.quantity, i.slot,
		       it.name, it.type, it.quality, it.description
		FROM inventory i
		JOIN items it ON i.item_id = it.id
		WHERE i.character_id = ?`
	
	args := []interface{}{characterID}

	// 添加筛选条件
	if itemType, ok := filters["type"].(string); ok && itemType != "" {
		query += " AND it.type = ?"
		args = append(args, itemType)
	}
	if quality, ok := filters["quality"].(string); ok && quality != "" {
		query += " AND it.quality = ?"
		args = append(args, quality)
	}

	query += " ORDER BY i.slot"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var id, quantity, slot int
		var itemID, name, itemType, quality string
		var description sql.NullString

		err := rows.Scan(&id, &itemID, &quantity, &slot, &name, &itemType, &quality, &description)
		if err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":       id,
			"item_id":  itemID,
			"quantity": quantity,
			"slot":     slot,
			"name":     name,
			"type":     itemType,
			"quality":  quality,
		}
		if description.Valid {
			item["description"] = description.String
		}
		items = append(items, item)
	}

	return items, nil
}

