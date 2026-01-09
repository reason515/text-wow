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
// Calculation 相关函数

func (tr *TestRunner) executeCalculatePhysCritDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	critDamage := tr.calculator.CalculatePhysCritDamage(char)
	// 更新角色的属性

	char.PhysCritDamage = critDamage
	tr.context.Characters["character"] = char
	// 存储到断言上下文和Variables

	tr.safeSetContext("phys_crit_damage", critDamage)
	tr.safeSetContext("character.phys_crit_damage", critDamage)
	tr.context.Variables["phys_crit_damage"] = critDamage
	tr.context.Variables["character_phys_crit_damage"] = critDamage
	return nil
}

func (tr *TestRunner) executeCalculateSpellCritDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	critDamage := tr.calculator.CalculateSpellCritDamage(char)
	// 更新角色的属性

	char.SpellCritDamage = critDamage
	tr.context.Characters["character"] = char
	// 存储到断言上下文和Variables

	tr.safeSetContext("spell_crit_damage", critDamage)
	tr.safeSetContext("character.spell_crit_damage", critDamage)
	tr.context.Variables["spell_crit_damage"] = critDamage
	tr.context.Variables["character_spell_crit_damage"] = critDamage
	return nil
}

func (tr *TestRunner) executeCalculateSpeed() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 确保敏捷值正确（从Variables恢复，如果存在）

	if agilityVal, exists := tr.context.Variables["character_agility"]; exists {
		if agility, ok := agilityVal.(int); ok {
			char.Agility = agility
			debugPrint("[DEBUG] executeCalculateSpeed: restored Agility=%d from Variables\n", agility)
		}

	}

	debugPrint("[DEBUG] executeCalculateSpeed: char.Agility=%d\n", char.Agility)
	speed := tr.calculator.CalculateSpeed(char)
	debugPrint("[DEBUG] executeCalculateSpeed: calculated speed=%d\n", speed)
	tr.safeSetContext("speed", speed)
	tr.context.Variables["speed"] = speed
	return nil
}

func (tr *TestRunner) executeCalculateResourceRegen(instruction string) error {
	// 怒气获得不需要角�

	if strings.Contains(instruction, "怒气") || strings.Contains(instruction, "rage") {
		// 解析基础获得值（如"计算怒气获得（基础获得=10）"）

	baseGain := 0
		if strings.Contains(instruction, "基础获得=") {
			parts := strings.Split(instruction, "基础获得=")
			if len(parts) > 1 {
				gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
				gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])
				if gain, err := strconv.Atoi(gainStr); err == nil {
					baseGain = gain
				}

			}

		}

		// 如果没有在指令中指定，尝试从Variables获取

		if baseGain == 0 {
			if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
				if gain, ok := gainVal.(int); ok {
					baseGain = gain
				}

			}

		}

	// 解析加成百分比（从Variables获取�

		bonusPercent := 0.0
		if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
			if percent, ok := percentVal.(float64); ok {
				bonusPercent = percent
			}

		}

		// 默认基础获得值

		if baseGain == 0 {
			baseGain = 10
		}

		regen := tr.calculator.CalculateRageGain(baseGain, bonusPercent)
		tr.safeSetContext("rage_gain", regen)
		tr.context.Variables["rage_gain"] = regen
		return nil
	}

	// 其他资源类型需要角色（但允许nil�

	char, ok := tr.context.Characters["character"]
	if !ok {
		return fmt.Errorf("character not found")
	}

	// 允许char为nil（用于测试nil情况）

	// 解析基础恢复值（如"计算法力恢复（基础恢复=10）"）

	baseRegen := 0
	if strings.Contains(instruction, "基础恢复=") {
		parts := strings.Split(instruction, "基础恢复=")
		if len(parts) > 1 {
			regenStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
			regenStr = strings.TrimSpace(strings.Split(regenStr, "）")[0])
			if regen, err := strconv.Atoi(regenStr); err == nil {
				baseRegen = regen
			}

		}

	}

	// 解析基础获得值（如"计算怒气获得（基础获得=10）"）

	baseGain := 0
	if strings.Contains(instruction, "基础获得=") {
		parts := strings.Split(instruction, "基础获得=")
		if len(parts) > 1 {
			gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
			gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])
			if gain, err := strconv.Atoi(gainStr); err == nil {
				baseGain = gain
			}

		}

	}

	// 如果没有在指令中指定，尝试从Variables获取

	if baseGain == 0 {
		if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
			if gain, ok := gainVal.(int); ok {
				baseGain = gain
			}

		}

	}

	// 解析加成百分比（从Variables获取�

	bonusPercent := 0.0
	if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
		if percent, ok := percentVal.(float64); ok {
			bonusPercent = percent
		}

	}

	// 如果没有在指令中指定基础恢复，尝试从Variables获取

	if baseRegen == 0 {
		if regenVal, exists := tr.context.Variables["mana_base_regen"]; exists {
			if regen, ok := regenVal.(int); ok {
				baseRegen = regen
			}

		}

	}

	// 根据指令确定资源类型

	if strings.Contains(instruction, "法力") || strings.Contains(instruction, "mana") {
		regen := tr.calculator.CalculateManaRegen(char, baseRegen)
		tr.safeSetContext("mana_regen", regen)
		tr.context.Variables["mana_regen"] = regen
	} else if strings.Contains(instruction, "怒气") || strings.Contains(instruction, "rage") {
	// 怒气获得不需要角色，只需要基础获得值和加成百分�

		if baseGain > 0 {
			// 使用基础获得值和加成百分比

			regen := tr.calculator.CalculateRageGain(baseGain, bonusPercent)
			tr.safeSetContext("rage_gain", regen)
			tr.context.Variables["rage_gain"] = regen
		} else {
			// 默认基础获得值

			regen := tr.calculator.CalculateRageGain(10, bonusPercent)
			tr.safeSetContext("rage_gain", regen)
			tr.context.Variables["rage_gain"] = regen
		}

	} else if strings.Contains(instruction, "能量") || strings.Contains(instruction, "energy") {
		regen := tr.calculator.CalculateEnergyRegen(char, baseRegen)
		tr.safeSetContext("energy_regen", regen)
		tr.context.Variables["energy_regen"] = regen
	} else {
	// 默认使用角色的资源类�

		resourceType := char.ResourceType
		if resourceType == "" {
			resourceType = "mana"
		}

		var regen int
		var key string
		switch resourceType {
		case "mana":
			regen = tr.calculator.CalculateManaRegen(char, baseRegen)
			key = "mana_regen"
		case "rage":
			// 从Variables获取基础获得值和加成百分比

			rageBaseGain := 10
			rageBonusPercent := 0.0
			if gainVal, exists := tr.context.Variables["rage_base_gain"]; exists {
				if gain, ok := gainVal.(int); ok {
					rageBaseGain = gain
				}

			}

			if percentVal, exists := tr.context.Variables["rage_bonus_percent"]; exists {
				if percent, ok := percentVal.(float64); ok {
					rageBonusPercent = percent
				}

			}

			regen = tr.calculator.CalculateRageGain(rageBaseGain, rageBonusPercent)
			key = "rage_gain"
		case "energy":
			regen = tr.calculator.CalculateEnergyRegen(char, baseRegen)
			key = "energy_regen"
		default:
			regen = tr.calculator.CalculateManaRegen(char, baseRegen)
			key = "resource_regen"
		}

		tr.safeSetContext(key, regen)
		tr.context.Variables[key] = regen
	}

	return nil
}

func (tr *TestRunner) executeCalculateBaseDamage() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 基础伤害 = 攻击�× 技能系数（默认1.0�

	baseDamage := char.PhysicalAttack
	tr.safeSetContext("base_damage", baseDamage)
	tr.context.Variables["base_damage"] = baseDamage
	return nil
}

func (tr *TestRunner) executeCalculateDefenseReduction() error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	monster, ok := tr.context.Monsters["monster"]
	if !ok || monster == nil {
		return fmt.Errorf("monster not found")
	}

	// 获取基础伤害（如果已计算）

	baseDamage := char.PhysicalAttack
	if val, exists := tr.context.Variables["base_damage"]; exists {
		if bd, ok := val.(int); ok {
			baseDamage = bd
		}

	}

	// 应用防御减伤（减法公式）

	damageAfterDefense := baseDamage - monster.PhysicalDefense
	if damageAfterDefense < 1 {
		damageAfterDefense = 1 // 至少1点伤害
	tr.safeSetContext("damage_after_defense", damageAfterDefense)
	tr.context.Variables["damage_after_defense"] = damageAfterDefense
	// 如果没有最终伤害，使用减伤后伤害作为最终伤害

	if _, exists := tr.context.Variables["final_damage"]; !exists {
		tr.safeSetContext("final_damage", damageAfterDefense)
		tr.context.Variables["final_damage"] = damageAfterDefense
	}

	return nil
}

func (tr *TestRunner) executeApplyCrit() error {
	// 从上下文中获取伤害值

	var baseDamage int
	if val, exists := tr.context.Variables["damage_after_defense"]; exists {
		if bd, ok := val.(int); ok {
			baseDamage = bd
		}

	}

	if baseDamage == 0 {
		// 如果没有伤害值，尝试从角色和怪物计算

		char, ok := tr.context.Characters["character"]
		if !ok || char == nil {
			return fmt.Errorf("character not found")
		}

		monster, ok := tr.context.Monsters["monster"]
		if !ok || monster == nil {
			return fmt.Errorf("monster not found")
		}

		baseDamage = char.PhysicalAttack - monster.PhysicalDefense
		if baseDamage < 1 {
			baseDamage = 1
		}

		// 更新上下文

		tr.safeSetContext("damage_after_defense", baseDamage)
		tr.context.Variables["damage_after_defense"] = baseDamage
	}

	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 假设暴击（实际应该随机判断）

	// 注意：PhysCritDamage是倍率，如1.5表示150%

	finalDamage := int(float64(baseDamage) * char.PhysCritDamage)
	tr.safeSetContext("final_damage", finalDamage)
	tr.context.Variables["final_damage"] = finalDamage
	return nil
}
