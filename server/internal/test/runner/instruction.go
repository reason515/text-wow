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
// Instruction 相关函数

func (tr *TestRunner) executeSetup(setup []string) error {
	for _, instruction := range setup {
		debugPrint("[DEBUG] executeSetup: processing instruction: %s\n", instruction)
		if err := tr.executeInstruction(instruction); err != nil {
			debugPrint("[DEBUG] executeSetup: instruction failed: %s, error: %v\n", instruction, err)
			return fmt.Errorf("setup instruction failed: %w", err)
		}
		debugPrint("[DEBUG] executeSetup: instruction completed: %s, characters=%d\n", instruction, len(tr.context.Characters))
	}
	return nil
}

func (tr *TestRunner) executeStep(step TestStep) error {
	// 将max_rounds存储到上下文中，�继续战斗直到"等指令使�	if step.MaxRounds > 0 {
		tr.context.Variables["step_max_rounds"] = step.MaxRounds
	}
	if err := tr.executeInstruction(step.Action); err != nil {
		return fmt.Errorf("step action failed: %s, error: %w", step.Action, err)
	}
	// 更新断言上下�	tr.updateAssertionContext()
	return nil
}

func (tr *TestRunner) executeInstruction(instruction string) error {
	// 处理装备相关操作
	if strings.Contains(instruction, "掉落") && strings.Contains(instruction, "装备") {
		return tr.generateEquipmentFromMonster(instruction)
	} else if strings.Contains(instruction, "连续") && strings.Contains(instruction, "装备") {
		return tr.generateMultipleEquipments(instruction)
	} else if strings.Contains(instruction, "获得") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"获得一件X级武器，攻击�X"这样的setup指令
		return tr.generateEquipmentWithAttributes(instruction)
	} else if strings.Contains(instruction, "尝试穿戴") || strings.Contains(instruction, "尝试装备") {
		// 处理"角色尝试穿戴武器"等action（用于测试失败情况）
		// 必须�穿戴"之前检查，因为"尝试穿戴"包含"穿戴"
		return tr.executeTryEquipItem(instruction)
	} else if strings.Contains(instruction, "穿戴") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"角色穿戴武器"�角色穿戴装备"等action
		return tr.executeEquipItem(instruction)
	} else if strings.Contains(instruction, "卸下") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		// 处理"角色卸下武器"�角色卸下装备"等action
		return tr.executeUnequipItem(instruction)
	} else if strings.Contains(instruction, "依次穿戴") && strings.Contains(instruction, "装备") {
		// 处理"角色依次穿戴所有装�
		return tr.executeEquipAllItems(instruction)
	} else if strings.Contains(instruction, "检查词缀") || strings.Contains(instruction, "检查词缀数值") || strings.Contains(instruction, "检查词缀类型") || strings.Contains(instruction, "检查词缀Tier") {
		// 这些操作已经在updateAssertionContext中处�		return nil
	} else if strings.Contains(instruction, "设置") {
		return tr.executeSetVariable(instruction)
	} else if strings.Contains(instruction, "创建一个nil角色") {
		// 创建一个nil角色（用于测试nil情况�		tr.context.Characters["character"] = nil
		return nil
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "队伍") {
		// 创建多人队伍（如"创建一�人队伍：战士(HP=100)、牧�HP=100)、法�HP=100)"�		return tr.createTeam(instruction)
	} else if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "角色") {
		// 必须�创建N个角�之前检查，因为"创建一个角�也包�创建"�个角�
		debugPrint("[DEBUG] executeInstruction: matched '创建一个角� pattern for: %s\n", instruction)
		return tr.createCharacter(instruction)
	} else if (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") && strings.Contains(instruction, "在")) {
		// 处理"创建3个角色：角色1（敏�30），角色2（敏�50�这样的指�		// 注意：必须排�创建一个角�，因为上面已经处理了
		debugPrint("[DEBUG] executeInstruction: matched '创建N个角� pattern for: %s\n", instruction)
		return tr.createMultipleCharacters(instruction)
	} else if strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") {
		// 处理"创建角色"（没�一��N�）的情况
		debugPrint("[DEBUG] executeInstruction: matched '创建角色' pattern for: %s\n", instruction)
		return tr.createCharacter(instruction)
	} else if (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个怪物")) || (strings.Contains(instruction, "创建") && strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在")) {
		// 处理"创建3个怪物：怪物1（速度=40），怪物2（速度=80�这样的指�		return tr.createMultipleMonsters(instruction)
	} else if (strings.Contains(instruction, "创建一个") || strings.Contains(instruction, "创建")) && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "击败") && strings.Contains(instruction, "怪物") {
		return tr.createMonster(instruction)
	} else if strings.Contains(instruction, "计算物理攻击力") {
		return tr.executeCalculatePhysicalAttack()
	} else if strings.Contains(instruction, "计算法术攻击�) {
		return tr.executeCalculateMagicAttack()
	} else if strings.Contains(instruction, "计算最大生命�) || strings.Contains(instruction, "计算生命�) {
		return tr.executeCalculateMaxHP()
	} else if strings.Contains(instruction, "计算物理暴击�) {
		return tr.executeCalculatePhysCritRate()
	} else if strings.Contains(instruction, "计算法术暴击�) {
		return tr.executeCalculateSpellCritRate()
	} else if strings.Contains(instruction, "计算物理暴击伤害倍率") {
		return tr.executeCalculatePhysCritDamage()
	} else if strings.Contains(instruction, "计算物理防御�) {
		return tr.executeCalculatePhysicalDefense()
	} else if strings.Contains(instruction, "计算魔法防御�) {
		return tr.executeCalculateMagicDefense()
	} else if strings.Contains(instruction, "计算法术暴击伤害倍率") {
		return tr.executeCalculateSpellCritDamage()
	} else if strings.Contains(instruction, "计算闪避�) {
		return tr.executeCalculateDodgeRate()
	} else if strings.Contains(instruction, "角色对怪物进行") && strings.Contains(instruction, "次攻�) {
		return tr.executeMultipleAttacks(instruction)
	} else if strings.Contains(instruction, "计算速度") {
		return tr.executeCalculateSpeed()
	} else if strings.Contains(instruction, "计算资源回复") || strings.Contains(instruction, "计算法力回复") || strings.Contains(instruction, "计算法力恢复") || strings.Contains(instruction, "计算怒气获得") || strings.Contains(instruction, "计算能量回复") || strings.Contains(instruction, "计算能量恢复") {
		return tr.executeCalculateResourceRegen(instruction)
	} else if strings.Contains(instruction, "计算队伍总攻击力") || strings.Contains(instruction, "计算队伍总生命�) {
		// 计算队伍属性（会调用syncTeamToContext�		tr.syncTeamToContext()
		return nil
	} else if strings.Contains(instruction, "有队伍攻击力") || strings.Contains(instruction, "有队伍生命�) {
		// 解析"角色1有队伍攻击力+10%的被动技��角色2有队伍生命�15%的被动技�
		if strings.Contains(instruction, "队伍攻击�) && strings.Contains(instruction, "+") && strings.Contains(instruction, "%") {
			// 解析攻击力加成百分比
			parts := strings.Split(instruction, "队伍攻击�)
			if len(parts) > 1 {
				bonusPart := parts[1]
				if plusIdx := strings.Index(bonusPart, "+"); plusIdx >= 0 {
					bonusStr := bonusPart[plusIdx+1:]
					bonusStr = strings.TrimSpace(strings.Split(bonusStr, "%")[0])
					if bonus, err := strconv.ParseFloat(bonusStr, 64); err == nil {
						tr.context.Variables["team_attack_bonus"] = bonus / 100.0 // 转换为小数（10% -> 0.1�					}
				}
			}
		}
		if strings.Contains(instruction, "队伍生命�) && strings.Contains(instruction, "+") && strings.Contains(instruction, "%") {
			// 解析生命值加成百分比
			parts := strings.Split(instruction, "队伍生命�)
			if len(parts) > 1 {
				bonusPart := parts[1]
				if plusIdx := strings.Index(bonusPart, "+"); plusIdx >= 0 {
					bonusStr := bonusPart[plusIdx+1:]
					bonusStr = strings.TrimSpace(strings.Split(bonusStr, "%")[0])
					if bonus, err := strconv.ParseFloat(bonusStr, 64); err == nil {
						tr.context.Variables["team_hp_bonus"] = bonus / 100.0 // 转换为小数（15% -> 0.15�					}
				}
			}
		}
		return nil
	} else if strings.Contains(instruction, "计算基础伤害") {
		return tr.executeCalculateBaseDamage()
	} else if strings.Contains(instruction, "应用防御减伤") {
		return tr.executeCalculateDefenseReduction()
	} else if strings.Contains(instruction, "计算防御减伤") || strings.Contains(instruction, "计算减伤后伤�) {
		return tr.executeCalculateDefenseReduction()
	} else if strings.Contains(instruction, "如果触发暴击，应用暴击倍率") || strings.Contains(instruction, "应用暴击倍率") {
		return tr.executeApplyCrit()
	} else if strings.Contains(instruction, "计算伤害") {
		return tr.executeCalculateDamage(instruction)
	} else if strings.Contains(instruction, "学习技�) || strings.Contains(instruction, "角色学习技�) {
		return tr.executeLearnSkill(instruction)
	} else if strings.Contains(instruction, "怪物使用") && strings.Contains(instruction, "技�) {
		// 怪物使用技能（包括Buff、Debuff、AOE、治疗等，必须在角色使用技能之前检查）
		return tr.executeMonsterUseSkill(instruction)
	} else if strings.Contains(instruction, "使用技�) || strings.Contains(instruction, "角色使用技�) || (strings.Contains(instruction, "使用") && strings.Contains(instruction, "技�)) {
		return tr.executeUseSkill(instruction)
	} else if strings.Contains(instruction, "创建一�) && strings.Contains(instruction, "技�) {
		return tr.createSkill(instruction)
	} else if strings.Contains(instruction, "执行�) && strings.Contains(instruction, "回合") {
		return tr.executeBattleRound(instruction)
	} else if strings.Contains(instruction, "构建回合顺序") {
		return tr.executeBuildTurnOrder()
	} else if strings.Contains(instruction, "开始战�) {
		return tr.executeStartBattle()
	} else if strings.Contains(instruction, "检查战斗初始状�) || strings.Contains(instruction, "检查战斗状�) {
		// 检查战斗状态，确保战士怒气�
		return tr.executeCheckBattleState(instruction)
	} else if strings.Contains(instruction, "检查战斗结束状�) {
		// 检查战斗结束状态，确保战士怒气�
		return tr.executeCheckBattleEndState()
	} else if strings.Contains(instruction, "角色攻击怪物") || strings.Contains(instruction, "攻击怪物") {
		return tr.executeAttackMonster()
	} else if strings.Contains(instruction, "怪物攻击角色") {
		return tr.executeMonsterAttack()
	} else if strings.Contains(instruction, "获取角色数据") || strings.Contains(instruction, "获取战斗状�) {
		// 获取角色数据或战斗状态，确保战士怒气正确
		return tr.executeGetCharacterData()
	} else if strings.Contains(instruction, "检查角色属�) || strings.Contains(instruction, "检查角�) {
		// 检查角色属性，确保所有属性都基于角色属性正确计�		return tr.executeCheckCharacterAttributes()
	} else if strings.Contains(instruction, "给怪物添加") && strings.Contains(instruction, "技�) {
		// 给怪物添加技�		return tr.executeAddMonsterSkill(instruction)
	} else if strings.Contains(instruction, "初始化战斗系�) {
		// 初始化战斗系统（空操作，战斗系统在开始战斗时自动初始化）
		return nil
	} else if strings.Contains(instruction, "继续战斗直到") {
		// 处理"继续战斗直到怪物死亡"�继续战斗直到所有怪物死亡"
		return tr.executeContinueBattleUntil(instruction)
	} else if strings.Contains(instruction, "所有怪物攻击") || strings.Contains(instruction, "所有敌人攻�) {
		// 处理"所有怪物攻击角色"�所有怪物攻击队伍"
		return tr.executeAllMonstersAttack(instruction)
	} else if strings.Contains(instruction, "剩余") && strings.Contains(instruction, "个怪物攻击") {
		// 处理"剩余2个怪物攻击角色"
		return tr.executeRemainingMonstersAttack(instruction)
	} else if strings.Contains(instruction, "角色攻击�) && strings.Contains(instruction, "个怪物") {
		// 处理"角色攻击第一个怪物"�角色攻击第二个怪物"
		return tr.executeAttackSpecificMonster(instruction)
	} else if strings.Contains(instruction, "怪物反击") {
		// 处理"怪物反击"（等同于"怪物攻击角色"�		return tr.executeMonsterAttack()
	} else if strings.Contains(instruction, "等待休息恢复") {
		// 处理"等待休息恢复"
		return tr.executeWaitRestRecovery()
	} else if strings.Contains(instruction, "进入休息状�) {
		// 处理"进入休息状态，休息速度倍率=X"
		return tr.executeEnterRestState(instruction)
	} else if strings.Contains(instruction, "记录战斗�) {
		// 处理"记录战斗后HP和Resource"（空操作，用于测试文档说明）
		return nil
	} else if strings.Contains(instruction, "创建一个空队伍") {
		// 处理"创建一个空队伍"
		return tr.executeCreateEmptyTeam()
	} else if strings.Contains(instruction, "创建一个队�) && (strings.Contains(instruction, "槽位") || strings.Contains(instruction, "包含")) {
		// 处理"创建一个队伍，槽位1已有角色1"�创建一个队伍，包含3个角�
		return tr.executeCreateTeamWithMembers(instruction)
	} else if strings.Contains(instruction, "将角�) && strings.Contains(instruction, "添加到槽�) {
		// 处理"将角�添加到槽�"
		return tr.executeAddCharacterToTeamSlot(instruction)
	} else if strings.Contains(instruction, "尝试将角�) && strings.Contains(instruction, "添加到槽�) {
		// 处理"尝试将角�添加到槽�"（用于测试失败情况）
		return tr.executeTryAddCharacterToTeamSlot(instruction)
	} else if strings.Contains(instruction, "从槽�) && strings.Contains(instruction, "移除角色") {
		// 处理"从槽�移除角色"
		return tr.executeRemoveCharacterFromTeamSlot(instruction)
	} else if strings.Contains(instruction, "解锁槽位") {
		// 处理"解锁槽位2"
		return tr.executeUnlockTeamSlot(instruction)
	} else if strings.Contains(instruction, "尝试将角色添加到槽位") {
		// 处理"尝试将角色添加到槽位2"（槽位未解锁的情况）
		return tr.executeTryAddCharacterToUnlockedSlot(instruction)
	} else if strings.Contains(instruction, "角色击败怪物") {
		// 处理"角色击败怪物"（给予经验和金币奖励�		return tr.executeDefeatMonster()
	} else if strings.Contains(instruction, "创建一个物�) {
		// 处理"创建一个物品，价格=30"
		return tr.executeCreateItem(instruction)
	} else if strings.Contains(instruction, "角色购买物品") || strings.Contains(instruction, "购买物品") {
		// 处理"角色购买物品"�购买物品A"
		return tr.executePurchaseItem(instruction)
	} else if strings.Contains(instruction, "角色尝试购买物品") {
		// 处理"角色尝试购买物品"（用于测试失败情况）
		return tr.executeTryPurchaseItem(instruction)
	} else if strings.Contains(instruction, "初始化商�) || strings.Contains(instruction, "初始化商店系�) {
		// 处理"初始化商店系��初始化商店，包含物品A（价�50�
		return tr.executeInitializeShop(instruction)
	} else if strings.Contains(instruction, "查看商店物品列表") {
		// 处理"查看商店物品列表"
		return tr.executeViewShopItems()
	} else if strings.Contains(instruction, "角色获得") && strings.Contains(instruction, "金币") {
		// 处理"角色获得1000金币"
		return tr.executeGainGold(instruction)
	} else if strings.Contains(instruction, "初始化地图管理器") {
		// 处理"初始化地图管理器"
		return tr.executeInitializeMapManager()
	} else if strings.Contains(instruction, "加载区域") {
		// 处理"加载区域 elwynn"
		return tr.executeLoadZone(instruction)
	} else if strings.Contains(instruction, "切换到区�) || strings.Contains(instruction, "尝试切换�) {
		// 处理"切换到区�elwynn"�尝试切换到需要等�0的区�
		return tr.executeSwitchZone(instruction)
	} else if strings.Contains(instruction, "创建一个区�) {
		// 处理"创建一个区域，经验倍率=1.5"�创建一个区域，经验倍率=1.5，金币倍率=1.2"
		return tr.executeCreateZone(instruction)
	} else if strings.Contains(instruction, "计算该区�) && strings.Contains(instruction, "倍率") {
		// 处理"计算该区域的经验倍率"�计算该区域的金币倍率"
		return tr.executeCalculateZoneMultiplier(instruction)
	} else if strings.Contains(instruction, "检查区�) && strings.Contains(instruction, "解锁状�) {
		// 处理"检查区�elwynn 的解锁状�
		return tr.executeCheckZoneUnlockStatus(instruction)
	} else if strings.Contains(instruction, "查询") && strings.Contains(instruction, "可用区域") {
		// 处理"查询等级10、阵营alliance的可用区�
		return tr.executeQueryAvailableZones(instruction)
	} else if strings.Contains(instruction, "角色�) && strings.Contains(instruction, "区域击杀") {
		// 处理"角色在该区域击杀怪物（基础经验=10，基础金币=5�
		return tr.executeKillMonsterInZone(instruction)
	} else if strings.Contains(instruction, "配置策略") {
		// 处理"配置策略：如果HP<60%，使用治疗技�
		return tr.executeConfigureStrategy(instruction)
	} else if strings.Contains(instruction, "执行策略判断") || strings.Contains(instruction, "执行策略选择") {
		// 处理"执行策略判断"�执行策略选择"
		return tr.executeStrategyDecision(instruction)
	} else if strings.Contains(instruction, "配置技能优先级") {
		// 处理"配置技能优先级：治疗（优先�0� 攻击（优先级5� 防御（优先级1�
		return tr.executeConfigureSkillPriority(instruction)
	} else if strings.Contains(instruction, "角色�) && strings.Contains(instruction, "区域击杀") && strings.Contains(instruction, "个怪物") {
		// 处理"角色�elwynn 区域击杀1个怪物"
		return tr.executeKillMonsterInZoneForExploration(instruction)
	} else if strings.Contains(instruction, "用户获得") && strings.Contains(instruction, "点探索度") {
		// 处理"用户获得10点探索度"
		return tr.executeGainExploration(instruction)
	} else if strings.Contains(instruction, "设置区域解锁要求") {
		// 处理"设置区域解锁要求：需�0点探索度"
		return tr.executeSetZoneUnlockRequirement(instruction)
	}
	return nil
}

func (tr *TestRunner) executeTeardown(teardown []string) error {
	// TODO: 实现清理逻辑
	// 例如：清理战斗状态、重置角色数据等
	return nil
}

func (tr *TestRunner) executeCreateItem(instruction string) error {
	// 解析物品价格，如"创建一个物品，价格=30"
	price := 0
	if strings.Contains(instruction, "价格=") {
		parts := strings.Split(instruction, "价格=")
		if len(parts) > 1 {
			priceStr := strings.TrimSpace(strings.Split(parts[1], "�)[0])
			if p, err := strconv.Atoi(priceStr); err == nil {
				price = p
			}
		}
	}

	// 存储物品信息到上下文
	tr.context.Variables["item_price"] = price
	tr.safeSetContext("item_price", price)

	return nil
}

func (tr *TestRunner) executePurchaseItem(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 获取物品价格
	price := 0
	if priceVal, exists := tr.context.Variables["item_price"]; exists {
		if p, ok := priceVal.(int); ok {
			price = p
		}
	} else if strings.Contains(instruction, "价格=") {
		// 从指令中解析价格，如"购买物品A（价�50�
		parts := strings.Split(instruction, "价格=")
		if len(parts) > 1 {
			priceStr := strings.TrimSpace(strings.Split(parts[1], "�)[0])
			if p, err := strconv.Atoi(priceStr); err == nil {
				price = p
			}
		}
	}

	// 解析物品名称（如"购买物品A"�	itemName := "物品A"
	if strings.Contains(instruction, "购买物品") {
		parts := strings.Split(instruction, "购买物品")
		if len(parts) > 1 {
			namePart := strings.TrimSpace(strings.Split(parts[1], "�)[0])
			if namePart != "" {
				itemName = namePart
			}
		}
	}

	// 获取用户金币
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(char.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 检查金币是否足�	if user.Gold < price {
		tr.context.Variables["purchase_success"] = false
		tr.safeSetContext("purchase_success", false)
		return fmt.Errorf("insufficient gold: need %d, have %d", price, user.Gold)
	}

	// 扣除金币
	user.Gold -= price
	_, err = database.DB.Exec(`UPDATE users SET gold =  WHERE id = `, user.Gold, char.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user gold: %w", err)
	}

	// 标记角色拥有该物�	itemKey := fmt.Sprintf("character.has_%s", strings.ToLower(strings.ReplaceAll(itemName, " ", "_")))
	tr.context.Variables[itemKey] = true
	tr.safeSetContext(itemKey, true)

	// 更新上下�	tr.context.Variables["character.gold"] = user.Gold
	tr.safeSetContext("character.gold", user.Gold)
	tr.context.Variables["purchase_success"] = true
	tr.safeSetContext("purchase_success", true)

	return nil
}

func (tr *TestRunner) executeTryPurchaseItem(instruction string) error {
	err := tr.executePurchaseItem(instruction)
	if err != nil {
		// 购买失败，设置purchase_success为false
		tr.context.Variables["purchase_success"] = false
		tr.safeSetContext("purchase_success", false)
		return nil // 不返回错误，因为这是预期的失�	}
	return nil
}

func (tr *TestRunner) executeViewShopItems() error {
	// 这个操作主要是为了测试，实际不需要做什�	// 物品列表已经在initializeShop中设置了
	return nil
}

func (tr *TestRunner) executeGainGold(instruction string) error {
	char, ok := tr.context.Characters["character"]
	if !ok || char == nil {
		return fmt.Errorf("character not found")
	}

	// 解析金币数量，如"角色获得1000金币"
	parts := strings.Split(instruction, "获得")
	if len(parts) > 1 {
		goldStr := strings.TrimSpace(strings.Split(parts[1], "金币")[0])
		if gold, err := strconv.Atoi(goldStr); err == nil {
			// 更新用户金币（Gold在User模型中）
			userRepo := repository.NewUserRepository()
			user, err := userRepo.GetByID(char.UserID)
			if err == nil && user != nil {
				user.Gold += gold
				_, err = database.DB.Exec(`UPDATE users SET gold = , total_gold_gained = total_gold_gained +  WHERE id = `, 
					user.Gold, gold, char.UserID)
				if err != nil {
					debugPrint("[DEBUG] executeGainGold: failed to update user gold: %v\n", err)
				}
				tr.context.Variables["character.gold"] = user.Gold
				tr.safeSetContext("character.gold", user.Gold)
			}
		}
	}

	return nil
}

