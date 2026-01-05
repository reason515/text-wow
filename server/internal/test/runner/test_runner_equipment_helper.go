package runner

import (
	"fmt"
	"os"

	"text-wow/internal/database"
)

// createOrUpdateTestItem 创建或更新测试item配置
func (tr *TestRunner) createOrUpdateTestItem(item map[string]interface{}) error {
	itemID, _ := item["id"].(string)
	if itemID == "" {
		return fmt.Errorf("item ID is required")
	}
	
	// 检查item是否已存在
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = ?)", itemID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check item existence: %w", err)
	}
	
	// 获取字段值
	name, _ := item["name"].(string)
	if name == "" {
		name = "测试装备"
	}
	itemType, _ := item["type"].(string)
	if itemType == "" {
		itemType = "equipment"
	}
	slot, _ := item["slot"].(string)
	quality, _ := item["quality"].(string)
	if quality == "" {
		quality = "rare"
	}
	levelRequired, _ := item["level_required"].(int)
	classRequired, _ := item["class_required"].(string)
	// 注意：items表没有strength_required等字段，这些要求需要通过其他方式实现
	// 暂时只支持level_required和class_required
	
	// 获取基础属性（如果有）
	strength, _ := item["strength"].(int)
	agility, _ := item["agility"].(int)
	intellect, _ := item["intellect"].(int)
	attack, _ := item["attack"].(int)
	defense, _ := item["defense"].(int)
	
	if exists {
		// 更新现有item
		_, err = database.DB.Exec(`
			UPDATE items 
			SET name = ?, type = ?, slot = ?, quality = ?, 
			    level_required = ?, class_required = ?,
			    strength = ?, agility = ?, intellect = ?,
			    attack = ?, defense = ?
			WHERE id = ?`,
			name, itemType, slot, quality,
			levelRequired, classRequired,
			strength, agility, intellect,
			attack, defense,
			itemID,
		)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] createOrUpdateTestItem: updated item %s with level_required=%d, class_required=%s\n", 
			itemID, levelRequired, classRequired)
	} else {
		// 创建新item
		_, err = database.DB.Exec(`
			INSERT INTO items (
				id, name, type, slot, quality,
				level_required, class_required,
				strength, agility, intellect,
				attack, defense, stackable, max_stack
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 1)`,
			itemID, name, itemType, slot, quality,
			levelRequired, classRequired,
			strength, agility, intellect,
			attack, defense,
		)
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] createOrUpdateTestItem: created item %s with level_required=%d, class_required=%s\n", 
			itemID, levelRequired, classRequired)
	}
	
	return nil
}

