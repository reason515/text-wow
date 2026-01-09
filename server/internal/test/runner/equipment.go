package runner

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"text-wow/internal/database"
	"text-wow/internal/game"
	"text-wow/internal/models"
	"text-wow/internal/repository"

	"gopkg.in/yaml.v3"
)
// Equipment 相关函数

func (tr *TestRunner) generateMultipleEquipments(action string) error {
	// 解析数量：如"连续获得10件蓝色装�?
	count := 10
	numStr := ""
	for _, r := range action {
		if r >= '0' && r <= '9' {
			numStr += string(r)
		} else if numStr != "" {
			break
		}
	}
	if numStr != "" {
		if n, err := strconv.Atoi(numStr); err == nil {
			count = n
		}
	}

	// 解析品质
	quality := "rare"
	if strings.Contains(action, "白色") || strings.Contains(action, "white") || strings.Contains(action, "common") {
		quality = "common"
	} else if strings.Contains(action, "绿色") || strings.Contains(action, "green") || strings.Contains(action, "uncommon") {
		quality = "uncommon"
	} else if strings.Contains(action, "蓝色") || strings.Contains(action, "blue") || strings.Contains(action, "rare") {
		quality = "rare"
	} else if strings.Contains(action, "紫色") || strings.Contains(action, "purple") || strings.Contains(action, "epic") {
		quality = "epic"
	}

	// 获取角色等级
	level := 1
	if char, ok := tr.context.Characters["character"]; ok {
		level = char.Level
	}

	// 确保用户和角色存�?	ownerID := 1
	if char, ok := tr.context.Characters["character"]; ok {
		ownerID = char.UserID
	} else {
		userRepo := repository.NewUserRepository()
		user, err := userRepo.GetByUsername("test_user")
		if err != nil {
			passwordHash := "test_hash"
			user, err = userRepo.Create("test_user", passwordHash, "test@test.com")
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		}
		ownerID = user.ID

		charRepo := repository.NewCharacterRepository()
		char, err := charRepo.Create(&models.Character{
			UserID:   user.ID,
			Name:     "测试角色",
			RaceID:   "human",
			ClassID:  "warrior",
			Faction:  "alliance",
			TeamSlot: 1,
			Level:    level,
		})
		if err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}
		tr.context.Characters["character"] = char
	}

	// 生成多件装备并统计唯一词缀组合
	uniqueCombinations := make(map[string]bool)
	itemID := "worn_sword"

	for i := 0; i < count; i++ {
		equipment, err := tr.equipmentManager.GenerateEquipment(itemID, quality, level, ownerID)
		if err != nil {
			continue
		}

		// 构建词缀组合字符�?		prefixID := "none"
		suffixID := "none"
		if equipment.PrefixID != nil {
			prefixID = *equipment.PrefixID
		}
		if equipment.SuffixID != nil {
			suffixID = *equipment.SuffixID
		}
		combination := fmt.Sprintf("%s_%s", prefixID, suffixID)
		uniqueCombinations[combination] = true

			// 存储最后一件装备到上下文（只存储基本字段，不存储整个对象）
		if i == count-1 {
			tr.context.Variables["equipment_id"] = equipment.ID
			tr.context.Variables["equipment_item_id"] = equipment.ItemID
			tr.context.Variables["equipment_quality"] = equipment.Quality
			tr.context.Variables["equipment_slot"] = equipment.Slot
		}
	}

	// 设置唯一词缀组合数量
	tr.context.Variables["unique_affix_combinations"] = len(uniqueCombinations)

	return nil
}

func (tr *TestRunner) generateEquipmentFromMonster(action string) error {
	// 解析品质：如"怪物掉落一件白色装�?
	quality := "common"
	if strings.Contains(action, "白色") || strings.Contains(action, "white") || strings.Contains(action, "common") {
		quality = "common"
	} else if strings.Contains(action, "绿色") || strings.Contains(action, "green") || strings.Contains(action, "uncommon") {
		quality = "uncommon"
	} else if strings.Contains(action, "蓝色") || strings.Contains(action, "blue") || strings.Contains(action, "rare") {
		quality = "rare"
	} else if strings.Contains(action, "紫色") || strings.Contains(action, "purple") || strings.Contains(action, "epic") {
		quality = "epic"
	} else if strings.Contains(action, "橙色") || strings.Contains(action, "orange") || strings.Contains(action, "legendary") {
		quality = "legendary"
	}

	// 处理"Boss掉落"的情�?	if strings.Contains(action, "Boss") || strings.Contains(action, "boss") {
		// 如果没有怪物，创建一个Boss怪物
		if len(tr.context.Monsters) == 0 {
			monster := &models.Monster{
				ID:              "boss_monster",
				Name:            "Boss怪物",
				Type:            "boss",
				Level:           30,
				HP:              0, // 被击�?				MaxHP:           1000,
				PhysicalAttack:  50,
				MagicAttack:     50,
				PhysicalDefense: 20,
				MagicDefense:    20,
				DodgeRate:       0.1,
			}
			tr.context.Monsters["monster"] = monster
		}
	}

	// 获取怪物等级
	level := 1
	for _, monster := range tr.context.Monsters {
		level = monster.Level
		break
	}

	// 确保用户和角色存�?	ownerID := 1
	if char, ok := tr.context.Characters["character"]; ok {
		ownerID = char.UserID
	} else {
		user, err := tr.createTestUser()
		if err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		ownerID = user.ID

		char, err := tr.createTestCharacter(user.ID, level)
		if err != nil {
			return fmt.Errorf("failed to create test character: %w", err)
		}
		tr.context.Characters["character"] = char
	}

	// 生成装备（使用数据库中存在的itemID�?	itemID := "worn_sword" // 使用seed.sql中存在的itemID
	equipment, err := tr.equipmentManager.GenerateEquipment(itemID, quality, level, ownerID)
	if err != nil {
		return fmt.Errorf("failed to generate equipment: %w", err)
	}

	// 存储到上下文（只存储基本字段，不存储整个对象�?	tr.context.Variables["equipment_id"] = equipment.ID
	tr.context.Variables["equipment_item_id"] = equipment.ItemID
	tr.context.Variables["equipment_quality"] = equipment.Quality
	tr.context.Variables["equipment_slot"] = equipment.Slot
	tr.context.Equipments[fmt.Sprintf("%d", equipment.ID)] = equipment

	return nil
}

