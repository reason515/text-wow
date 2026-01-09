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
// Team 相关函数

func (tr *TestRunner) createTeam(instruction string) error {
	// 确保用户存在
	user, err := tr.createTestUser()
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 解析队伍成员（通过冒号或逗号分隔�	// 格式：战�HP=100)、牧�HP=100)、法�HP=100)
	var members []string
	if strings.Contains(instruction, "�) {
		parts := strings.Split(instruction, "�)
		if len(parts) > 1 {
			members = strings.Split(parts[1], "�)
		}
	} else if strings.Contains(instruction, ":") {
		parts := strings.Split(instruction, ":")
		if len(parts) > 1 {
			members = strings.Split(parts[1], ",")
		}
	}

	charRepo := repository.NewCharacterRepository()
	slot := 1

	// 先获取用户的所有角色，检查哪些slot已被占用
	existingChars, err := charRepo.GetByUserID(user.ID)
	if err != nil {
		existingChars = []*models.Character{}
	}
	existingSlots := make(map[int]*models.Character)
	for _, c := range existingChars {
		existingSlots[c.TeamSlot] = c
	}

	for _, memberDesc := range members {
		memberDesc = strings.TrimSpace(memberDesc)
		if memberDesc == "" {
			continue
		}

		// 解析职业（战士、牧师、法师等�		classID := "warrior"
		if strings.Contains(memberDesc, "战士") {
			classID = "warrior"
		} else if strings.Contains(memberDesc, "牧师") {
			classID = "priest"
		} else if strings.Contains(memberDesc, "法师") {
			classID = "mage"
		} else if strings.Contains(memberDesc, "盗贼") {
			classID = "rogue"
		}

		// 解析HP（如"HP=100"�		hp := 100
		if strings.Contains(memberDesc, "HP=") {
			parts := strings.Split(memberDesc, "HP=")
			if len(parts) > 1 {
				hpStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])
				if h, err := strconv.Atoi(hpStr); err == nil {
					hp = h
				}
			}
		}

		// 检查该slot是否已存在角�		var createdChar *models.Character
		if existingChar, exists := existingSlots[slot]; exists {
			// 更新已存在的角色
			existingChar.Name = fmt.Sprintf("测试角色%d", slot)
			existingChar.ClassID = classID
			existingChar.HP = hp
			existingChar.MaxHP = hp
			existingChar.Level = 1
			existingChar.Strength = 10
			existingChar.Agility = 10
			existingChar.Intellect = 10
			existingChar.Stamina = 10
			existingChar.Spirit = 10

			// 根据职业设置资源类型
			if classID == "warrior" {
				existingChar.ResourceType = "rage"
				existingChar.MaxResource = 100
				existingChar.Resource = 0
			} else if classID == "rogue" {
				existingChar.ResourceType = "energy"
				existingChar.MaxResource = 100
				existingChar.Resource = 100
			} else {
				existingChar.ResourceType = "mana"
				existingChar.MaxResource = 100
				existingChar.Resource = 100
			}

			// 更新到数据库
			if err := charRepo.Update(existingChar); err != nil {
				return fmt.Errorf("failed to update character in team: %w", err)
			}
			createdChar = existingChar
		} else {
			// 创建新角�			char := &models.Character{
				UserID:    user.ID,
				Name:      fmt.Sprintf("测试角色%d", slot),
				RaceID:    "human",
				ClassID:   classID,
				Faction:   "alliance",
				TeamSlot:  slot,
				Level:     1,
				HP:        hp,
				MaxHP:     hp,
				Strength:  10,
				Agility:   10,
				Intellect: 10,
				Stamina:   10,
				Spirit:    10,
			}

			// 根据职业设置资源类型
			if classID == "warrior" {
				char.ResourceType = "rage"
				char.MaxResource = 100
				char.Resource = 0
			} else if classID == "rogue" {
				char.ResourceType = "energy"
				char.MaxResource = 100
				char.Resource = 100
			} else {
				char.ResourceType = "mana"
				char.MaxResource = 100
				char.Resource = 100
			}

			// 保存到数据库
			var err error
			createdChar, err = charRepo.Create(char)
			if err != nil {
				return fmt.Errorf("failed to create character in team: %w", err)
			}
		}

		// 保存到上下文（使用character_1, character_2等作为key�		key := fmt.Sprintf("character_%d", slot)
		tr.context.Characters[key] = createdChar

		// 第一个角色也保存�character"（向后兼容）
		if slot == 1 {
			tr.context.Characters["character"] = createdChar
		}

		slot++
	}

	return nil
}

func (tr *TestRunner) executeCreateEmptyTeam() error {
	// 清空所有角色（除了character，保留作为默认角色）
	// 实际上，空队伍意味着没有角色在队伍槽位中
	// 我们只需要确保team.character_count�
	tr.context.Variables["team.character_count"] = 0
	tr.safeSetContext("team.character_count", 0)
	return nil
}

func (tr *TestRunner) executeCreateTeamWithMembers(instruction string) error {
	// 解析指令，如"创建一个队伍，槽位1已有角色1"�创建一个队伍，包含3个角�
	if strings.Contains(instruction, "槽位") && strings.Contains(instruction, "已有") {
		// 解析槽位和角色ID
		// �槽位1已有角色1"
		parts := strings.Split(instruction, "槽位")
		if len(parts) > 1 {
			slotPart := strings.TrimSpace(strings.Split(parts[1], "已有")[0])
			if slot, err := strconv.Atoi(slotPart); err == nil {
				// 解析角色ID
				charIDPart := strings.TrimSpace(strings.Split(parts[1], "角色")[1])
				if charID, err := strconv.Atoi(charIDPart); err == nil {
					// 创建或获取角�					char, err := tr.getOrCreateCharacterByID(charID, slot)
					if err != nil {
						return err
					}
					key := fmt.Sprintf("character_%d", slot)
					tr.context.Characters[key] = char
					tr.context.Variables["team.character_count"] = 1
					tr.safeSetContext("team.character_count", 1)
					tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_id", slot), char.ID)
				}
			}
		}
	} else if strings.Contains(instruction, "包含") && strings.Contains(instruction, "个角�) {
		// 解析角色数量，如"包含3个角�
		parts := strings.Split(instruction, "包含")
		if len(parts) > 1 {
			countStr := strings.TrimSpace(strings.Split(parts[1], "�)[0])
			if count, err := strconv.Atoi(countStr); err == nil {
				// 创建指定数量的角�				for i := 1; i <= count; i++ {
					char, err := tr.getOrCreateCharacterByID(i, i)
					if err != nil {
						return err
					}
					key := fmt.Sprintf("character_%d", i)
					tr.context.Characters[key] = char
				}
				tr.context.Variables["team.character_count"] = count
				tr.safeSetContext("team.character_count", count)
				// 创建队伍后，同步队伍信息到上下文
				tr.syncTeamToContext()
			}
		}
	}
	return nil
}

func (tr *TestRunner) executeAddCharacterToTeamSlot(instruction string) error {
	// 解析指令，如"将角�添加到槽�"
	parts := strings.Split(instruction, "将角�)
	if len(parts) < 2 {
		return fmt.Errorf("invalid instruction: %s", instruction)
	}
	
	charIDPart := strings.TrimSpace(strings.Split(parts[1], "添加到槽�)[0])
	charID, err := strconv.Atoi(charIDPart)
	if err != nil {
		return fmt.Errorf("failed to parse character ID: %w", err)
	}
	
	slotPart := strings.TrimSpace(strings.Split(parts[1], "槽位")[1])
	slot, err := strconv.Atoi(slotPart)
	if err != nil {
		return fmt.Errorf("failed to parse slot: %w", err)
	}
	
	// 检查槽位是否已被占�	slotKey := fmt.Sprintf("character_%d", slot)
	if existingChar, exists := tr.context.Characters[slotKey]; exists && existingChar != nil {
		return fmt.Errorf("slot %d is already occupied", slot)
	}
	
	// 检查槽位是否解锁（简化：假设�个槽位默认解锁）
	if slot > 5 {
		// 检查unlocked_slots
		unlockedSlots := 1
		if unlockedVal, exists := tr.context.Variables["team.unlocked_slots"]; exists {
			if u, ok := unlockedVal.(int); ok {
				unlockedSlots = u
			}
		}
		if slot > unlockedSlots {
			tr.context.Variables["operation_success"] = false
			tr.safeSetContext("operation_success", false)
			return fmt.Errorf("slot %d is not unlocked", slot)
		}
	}
	
	// 获取或创建角�	char, err := tr.getOrCreateCharacterByID(charID, slot)
	if err != nil {
		return err
	}
	
	// 添加到槽�	tr.context.Characters[slotKey] = char
	
	// 更新队伍角色�	teamCount := 0
	for key, c := range tr.context.Characters {
		if c != nil && (strings.HasPrefix(key, "character_") || key == "character") {
			teamCount++
		}
	}
	tr.context.Variables["team.character_count"] = teamCount
	tr.safeSetContext("team.character_count", teamCount)
	tr.safeSetContext(fmt.Sprintf("team.slot_%d.character_id", slot), char.ID)
	
	tr.context.Variables["operation_success"] = true
	tr.safeSetContext("operation_success", true)
	
	return nil
}

func (tr *TestRunner) executeTryAddCharacterToTeamSlot(instruction string) error {
	err := tr.executeAddCharacterToTeamSlot(instruction)
	if err != nil {
		// 操作失败，设置operation_success为false
		tr.context.Variables["operation_success"] = false
		tr.safeSetContext("operation_success", false)
		return nil // 不返回错误，因为这是预期的失�	}
	tr.context.Variables["operation_success"] = true
	tr.safeSetContext("operation_success", true)
	return nil
}

func (tr *TestRunner) executeUnlockTeamSlot(instruction string) error {
	// 解析指令，如"解锁槽位2"
	parts := strings.Split(instruction, "槽位")
	if len(parts) < 2 {
		return fmt.Errorf("invalid instruction: %s", instruction)
	}
	
	slotPart := strings.TrimSpace(parts[1])
	slot, err := strconv.Atoi(slotPart)
	if err != nil {
		return fmt.Errorf("failed to parse slot: %w", err)
	}
	
	// 更新解锁槽位�	tr.context.Variables["team.unlocked_slots"] = slot
	tr.safeSetContext("team.unlocked_slots", slot)
	
	return nil
}

func (tr *TestRunner) executeTryAddCharacterToUnlockedSlot(instruction string) error {
	// 这个函数会尝试添加，但应该失�	return tr.executeTryAddCharacterToTeamSlot(instruction)
}

