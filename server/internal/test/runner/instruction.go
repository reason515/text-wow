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
	// 将max_rounds存储到上下文中，以便"继续战斗直到"等指令使用
	if step.MaxRounds > 0 {
		tr.context.Variables["step_max_rounds"] = step.MaxRounds
	}
	if err := tr.executeInstruction(step.Action); err != nil {
		return fmt.Errorf("step action failed: %s, error: %w", step.Action, err)
	}
	// 更新断言上下文
	tr.updateAssertionContext()
	return nil
}

// executeInstruction 执行指令（主入口）
func (tr *TestRunner) executeInstruction(instruction string) error {
	// 按类别依次尝试处理指令

	// 1. 装备相关
	if handled, err := tr.tryExecuteEquipmentInstruction(instruction); handled {
		return err
	}

	// 2. 角色相关
	if handled, err := tr.tryExecuteCharacterInstruction(instruction); handled {
		return err
	}

	// 3. 怪物相关
	if handled, err := tr.tryExecuteMonsterInstruction(instruction); handled {
		return err
	}

	// 4. 计算相关
	if handled, err := tr.tryExecuteCalculationInstruction(instruction); handled {
		return err
	}

	// 5. 技能相关
	if handled, err := tr.tryExecuteSkillInstruction(instruction); handled {
		return err
	}

	// 6. 战斗相关
	if handled, err := tr.tryExecuteBattleInstruction(instruction); handled {
		return err
	}

	// 7. 队伍相关
	if handled, err := tr.tryExecuteTeamInstruction(instruction); handled {
		return err
	}

	// 8. 商店相关
	if handled, err := tr.tryExecuteShopInstruction(instruction); handled {
		return err
	}

	// 9. 地图/区域相关
	if handled, err := tr.tryExecuteMapInstruction(instruction); handled {
		return err
	}

	// 10. 策略相关
	if handled, err := tr.tryExecuteStrategyInstruction(instruction); handled {
		return err
	}

	// 11. 其他指令
	if handled, err := tr.tryExecuteOtherInstruction(instruction); handled {
		return err
	}

	return nil
}

// tryExecuteEquipmentInstruction 尝试处理装备相关指令
func (tr *TestRunner) tryExecuteEquipmentInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "掉落") && strings.Contains(instruction, "装备") {
		return true, tr.generateEquipmentFromMonster(instruction)
	}
	if strings.Contains(instruction, "连续") && strings.Contains(instruction, "装备") {
		return true, tr.generateMultipleEquipments(instruction)
	}
	if strings.Contains(instruction, "获得") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		return true, tr.generateEquipmentWithAttributes(instruction)
	}
	if strings.Contains(instruction, "尝试穿戴") || strings.Contains(instruction, "尝试装备") {
		return true, tr.executeTryEquipItem(instruction)
	}
	if strings.Contains(instruction, "穿戴") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		return true, tr.executeEquipItem(instruction)
	}
	if strings.Contains(instruction, "卸下") && (strings.Contains(instruction, "装备") || strings.Contains(instruction, "武器") || strings.Contains(instruction, "护甲") || strings.Contains(instruction, "饰品")) {
		return true, tr.executeUnequipItem(instruction)
	}
	if strings.Contains(instruction, "依次穿戴") && strings.Contains(instruction, "装备") {
		return true, tr.executeEquipAllItems(instruction)
	}
	if strings.Contains(instruction, "检查词缀") || strings.Contains(instruction, "检查词缀数值") || strings.Contains(instruction, "检查词缀类型") || strings.Contains(instruction, "检查词缀Tier") {
		return true, nil
	}
	return false, nil
}

// tryExecuteCharacterInstruction 尝试处理角色相关指令
func (tr *TestRunner) tryExecuteCharacterInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "设置") {
		return true, tr.executeSetVariable(instruction)
	}
	if strings.Contains(instruction, "创建一个nil角色") {
		tr.context.Characters["character"] = nil
		return true, nil
	}
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "队伍") {
		return true, tr.createTeam(instruction)
	}
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "角色") {
		debugPrint("[DEBUG] executeInstruction: matched '创建一个角色' pattern for: %s\n", instruction)
		return true, tr.createCharacter(instruction)
	}
	if (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个角色") && !strings.Contains(instruction, "创建一个")) ||
		(strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") && strings.Contains(instruction, "在")) {
		debugPrint("[DEBUG] executeInstruction: matched '创建N个角色' pattern for: %s\n", instruction)
		return true, tr.createMultipleCharacters(instruction)
	}
	if strings.Contains(instruction, "创建") && strings.Contains(instruction, "角色") {
		debugPrint("[DEBUG] executeInstruction: matched '创建角色' pattern for: %s\n", instruction)
		return true, tr.createCharacter(instruction)
	}
	// 处理Buff相关指令
	if strings.Contains(instruction, "给角色添加") && strings.Contains(instruction, "Buff") {
		return true, tr.executeAddBuff(instruction)
	}
	// 处理护盾相关指令
	if strings.Contains(instruction, "给角色添加") && strings.Contains(instruction, "护盾") {
		return true, tr.executeAddShield(instruction)
	}
	// 处理"创建一个角色和一个怪物"
	if strings.Contains(instruction, "角色") && strings.Contains(instruction, "怪物") && strings.Contains(instruction, "和") {
		// 先创建角色
		if err := tr.createCharacter(instruction); err != nil {
			return true, err
		}
		// 再创建怪物
		return true, tr.createMonster(instruction)
	}
	return false, nil
}

// tryExecuteMonsterInstruction 尝试处理怪物相关指令
func (tr *TestRunner) tryExecuteMonsterInstruction(instruction string) (bool, error) {
	if (strings.Contains(instruction, "创建") && strings.Contains(instruction, "个怪物")) ||
		(strings.Contains(instruction, "创建") && strings.Contains(instruction, "怪物") && strings.Contains(instruction, "在")) {
		return true, tr.createMultipleMonsters(instruction)
	}
	if (strings.Contains(instruction, "创建一个") || strings.Contains(instruction, "创建")) && strings.Contains(instruction, "怪物") {
		return true, tr.createMonster(instruction)
	}
	if strings.Contains(instruction, "击败") && strings.Contains(instruction, "怪物") {
		return true, tr.createMonster(instruction)
	}
	return false, nil
}

// tryExecuteCalculationInstruction 尝试处理计算相关指令
func (tr *TestRunner) tryExecuteCalculationInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "计算物理攻击力") {
		return true, tr.executeCalculatePhysicalAttack()
	}
	if strings.Contains(instruction, "计算法术攻击") {
		return true, tr.executeCalculateMagicAttack()
	}
	if strings.Contains(instruction, "计算最大生命") || strings.Contains(instruction, "计算生命") {
		return true, tr.executeCalculateMaxHP()
	}
	if strings.Contains(instruction, "计算物理暴击") {
		return true, tr.executeCalculatePhysCritRate()
	}
	if strings.Contains(instruction, "计算法术暴击") {
		return true, tr.executeCalculateSpellCritRate()
	}
	if strings.Contains(instruction, "计算物理暴击伤害倍率") {
		return true, tr.executeCalculatePhysCritDamage()
	}
	if strings.Contains(instruction, "计算物理防御") {
		return true, tr.executeCalculatePhysicalDefense()
	}
	if strings.Contains(instruction, "计算魔法防御") {
		return true, tr.executeCalculateMagicDefense()
	}
	if strings.Contains(instruction, "计算法术暴击伤害倍率") {
		return true, tr.executeCalculateSpellCritDamage()
	}
	if strings.Contains(instruction, "计算闪避") {
		return true, tr.executeCalculateDodgeRate()
	}
	if strings.Contains(instruction, "计算速度") {
		return true, tr.executeCalculateSpeed()
	}
	if strings.Contains(instruction, "计算资源回复") || strings.Contains(instruction, "计算法力回复") ||
		strings.Contains(instruction, "计算法力恢复") || strings.Contains(instruction, "计算怒气获得") ||
		strings.Contains(instruction, "计算能量回复") || strings.Contains(instruction, "计算能量恢复") {
		return true, tr.executeCalculateResourceRegen(instruction)
	}
	if strings.Contains(instruction, "计算队伍总攻击力") || strings.Contains(instruction, "计算队伍总生命") {
		tr.syncTeamToContext()
		return true, nil
	}
	if strings.Contains(instruction, "有队伍攻击力") || strings.Contains(instruction, "有队伍生命") {
		return true, tr.executeParseTeamBonus(instruction)
	}
	if strings.Contains(instruction, "计算基础伤害") {
		return true, tr.executeCalculateBaseDamage()
	}
	if strings.Contains(instruction, "应用防御减伤") {
		return true, tr.executeCalculateDefenseReduction()
	}
	if strings.Contains(instruction, "计算防御减伤") || strings.Contains(instruction, "计算减伤后伤害") {
		return true, tr.executeCalculateDefenseReduction()
	}
	if strings.Contains(instruction, "如果触发暴击，应用暴击倍率") || strings.Contains(instruction, "应用暴击倍率") {
		return true, tr.executeApplyCrit()
	}
	if strings.Contains(instruction, "计算伤害") {
		return true, tr.executeCalculateDamage(instruction)
	}
	if strings.Contains(instruction, "角色对怪物进行") && strings.Contains(instruction, "次攻击") {
		return true, tr.executeMultipleAttacks(instruction)
	}
	return false, nil
}

// tryExecuteSkillInstruction 尝试处理技能相关指令
func (tr *TestRunner) tryExecuteSkillInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "学习技能") || strings.Contains(instruction, "角色学习技能") {
		return true, tr.executeLearnSkill(instruction)
	}
	if strings.Contains(instruction, "怪物使用") && strings.Contains(instruction, "技能") {
		return true, tr.executeMonsterUseSkill(instruction)
	}
	if strings.Contains(instruction, "使用技能") || strings.Contains(instruction, "角色使用技能") ||
		(strings.Contains(instruction, "使用") && strings.Contains(instruction, "技能")) {
		return true, tr.executeUseSkill(instruction)
	}
	if strings.Contains(instruction, "创建一个") && strings.Contains(instruction, "技能") {
		return true, tr.createSkill(instruction)
	}
	if strings.Contains(instruction, "给怪物添加") && strings.Contains(instruction, "技能") {
		// TODO: 实现给怪物添加技能的功能
		return true, nil
	}
	return false, nil
}

// tryExecuteBattleInstruction 尝试处理战斗相关指令
func (tr *TestRunner) tryExecuteBattleInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "执行") && strings.Contains(instruction, "回合") {
		return true, tr.executeBattleRound(instruction)
	}
	if strings.Contains(instruction, "构建回合顺序") {
		return true, tr.executeBuildTurnOrder()
	}
	if strings.Contains(instruction, "开始战斗") {
		return true, tr.executeStartBattle()
	}
	if strings.Contains(instruction, "检查战斗初始状态") || strings.Contains(instruction, "检查战斗状态") {
		return true, tr.executeCheckBattleState(instruction)
	}
	if strings.Contains(instruction, "检查战斗结束状态") {
		return true, tr.executeCheckBattleEndState()
	}
	if strings.Contains(instruction, "角色攻击怪物") || strings.Contains(instruction, "攻击怪物") {
		return true, tr.executeAttackMonster()
	}
	if strings.Contains(instruction, "怪物攻击角色") {
		return true, tr.executeMonsterAttack()
	}
	if strings.Contains(instruction, "获取角色数据") || strings.Contains(instruction, "获取战斗状态") {
		return true, tr.executeGetCharacterData()
	}
	if strings.Contains(instruction, "检查角色属性") || strings.Contains(instruction, "检查角色") {
		return true, tr.executeCheckCharacterAttributes()
	}
	if strings.Contains(instruction, "初始化战斗系统") {
		return true, nil
	}
	if strings.Contains(instruction, "继续战斗直到") {
		return true, tr.executeContinueBattleUntil(instruction)
	}
	if strings.Contains(instruction, "所有怪物攻击") || strings.Contains(instruction, "所有敌人攻击") {
		return true, tr.executeAllMonstersAttack(instruction)
	}
	if strings.Contains(instruction, "剩余") && strings.Contains(instruction, "个怪物攻击") {
		return true, tr.executeRemainingMonstersAttack(instruction)
	}
	if strings.Contains(instruction, "角色攻击") && strings.Contains(instruction, "个怪物") {
		return true, tr.executeAttackSpecificMonster(instruction)
	}
	if strings.Contains(instruction, "怪物反击") {
		return true, tr.executeMonsterAttack()
	}
	if strings.Contains(instruction, "等待休息恢复") {
		return true, tr.executeWaitRestRecovery()
	}
	if strings.Contains(instruction, "进入休息状态") {
		return true, tr.executeEnterRestState()
	}
	if strings.Contains(instruction, "记录战斗") {
		return true, nil
	}
	if strings.Contains(instruction, "角色击败怪物") {
		// 角色击败怪物，给予经验和金币奖励
		return true, tr.executeAttackMonster()
	}
	return false, nil
}

// tryExecuteTeamInstruction 尝试处理队伍相关指令
func (tr *TestRunner) tryExecuteTeamInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "创建一个空队伍") {
		return true, tr.executeCreateEmptyTeam()
	}
	if strings.Contains(instruction, "创建一个队伍") && (strings.Contains(instruction, "槽位") || strings.Contains(instruction, "包含")) {
		return true, tr.executeCreateTeamWithMembers(instruction)
	}
	if strings.Contains(instruction, "将角色") && strings.Contains(instruction, "添加到槽") {
		return true, tr.executeAddCharacterToTeamSlot(instruction)
	}
	if strings.Contains(instruction, "尝试将角色") && strings.Contains(instruction, "添加到槽") {
		return true, tr.executeTryAddCharacterToTeamSlot(instruction)
	}
	if strings.Contains(instruction, "从槽位") && strings.Contains(instruction, "移除角色") {
		return true, tr.executeRemoveCharacterFromTeamSlot(instruction)
	}
	if strings.Contains(instruction, "解锁槽位") {
		return true, tr.executeUnlockTeamSlot(instruction)
	}
	if strings.Contains(instruction, "尝试将角色添加到槽位") {
		return true, tr.executeTryAddCharacterToUnlockedSlot(instruction)
	}
	return false, nil
}

// tryExecuteShopInstruction 尝试处理商店相关指令
func (tr *TestRunner) tryExecuteShopInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "创建一个物品") {
		return true, tr.executeCreateItem(instruction)
	}
	if strings.Contains(instruction, "角色购买物品") || strings.Contains(instruction, "购买物品") {
		return true, tr.executePurchaseItem(instruction)
	}
	if strings.Contains(instruction, "角色尝试购买物品") {
		return true, tr.executeTryPurchaseItem(instruction)
	}
	if strings.Contains(instruction, "初始化商店") || strings.Contains(instruction, "初始化商店系统") {
		return true, tr.executeInitializeShop(instruction)
	}
	if strings.Contains(instruction, "查看商店物品列表") {
		return true, tr.executeViewShopItems()
	}
	if strings.Contains(instruction, "角色获得") && strings.Contains(instruction, "金币") {
		return true, tr.executeGainGold(instruction)
	}
	return false, nil
}

// tryExecuteMapInstruction 尝试处理地图/区域相关指令
func (tr *TestRunner) tryExecuteMapInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "初始化地图管理器") {
		return true, tr.executeInitializeMapManager()
	}
	if strings.Contains(instruction, "加载区域") {
		return true, tr.executeLoadZone(instruction)
	}
	if strings.Contains(instruction, "切换到区域") || strings.Contains(instruction, "尝试切换") {
		return true, tr.executeSwitchZone(instruction)
	}
	if strings.Contains(instruction, "创建一个区域") {
		return true, tr.executeCreateZone(instruction)
	}
	if strings.Contains(instruction, "计算该区域") && strings.Contains(instruction, "倍率") {
		return true, tr.executeCalculateZoneMultiplier(instruction)
	}
	if strings.Contains(instruction, "检查区域") && strings.Contains(instruction, "解锁状态") {
		return true, tr.executeCheckZoneUnlockStatus(instruction)
	}
	if strings.Contains(instruction, "查询") && strings.Contains(instruction, "可用区域") {
		return true, tr.executeQueryAvailableZones(instruction)
	}
	if strings.Contains(instruction, "角色") && strings.Contains(instruction, "区域击杀") && strings.Contains(instruction, "个怪物") {
		return true, tr.executeKillMonsterInZoneForExploration(instruction)
	}
	if strings.Contains(instruction, "角色") && strings.Contains(instruction, "区域击杀") {
		return true, tr.executeKillMonsterInZone(instruction)
	}
	if strings.Contains(instruction, "用户获得") && strings.Contains(instruction, "点探索度") {
		return true, tr.executeGainExploration(instruction)
	}
	if strings.Contains(instruction, "设置区域解锁要求") {
		return true, tr.executeSetZoneUnlockRequirement(instruction)
	}
	return false, nil
}

// tryExecuteStrategyInstruction 尝试处理策略相关指令
func (tr *TestRunner) tryExecuteStrategyInstruction(instruction string) (bool, error) {
	if strings.Contains(instruction, "配置策略") {
		return true, tr.executeConfigureStrategy(instruction)
	}
	if strings.Contains(instruction, "执行策略判断") || strings.Contains(instruction, "执行策略选择") {
		return true, tr.executeStrategyDecision(instruction)
	}
	if strings.Contains(instruction, "配置技能优先级") {
		return true, tr.executeConfigureSkillPriority(instruction)
	}
	return false, nil
}

// tryExecuteOtherInstruction 尝试处理其他指令
func (tr *TestRunner) tryExecuteOtherInstruction(instruction string) (bool, error) {
	// 暂时没有其他指令
	return false, nil
}

// executeParseTeamBonus 解析队伍加成
func (tr *TestRunner) executeParseTeamBonus(instruction string) error {
	if strings.Contains(instruction, "队伍攻击") && strings.Contains(instruction, "+") && strings.Contains(instruction, "%") {
		parts := strings.Split(instruction, "队伍攻击")
		if len(parts) > 1 {
			bonusPart := parts[1]
			if plusIdx := strings.Index(bonusPart, "+"); plusIdx >= 0 {
				bonusStr := bonusPart[plusIdx+1:]
				bonusStr = strings.TrimSpace(strings.Split(bonusStr, "%")[0])
				if bonus, err := strconv.ParseFloat(bonusStr, 64); err == nil {
					tr.context.Variables["team_attack_bonus"] = bonus / 100.0
				}
			}
		}
	}
	if strings.Contains(instruction, "队伍生命") && strings.Contains(instruction, "+") && strings.Contains(instruction, "%") {
		parts := strings.Split(instruction, "队伍生命")
		if len(parts) > 1 {
			bonusPart := parts[1]
			if plusIdx := strings.Index(bonusPart, "+"); plusIdx >= 0 {
				bonusStr := bonusPart[plusIdx+1:]
				bonusStr = strings.TrimSpace(strings.Split(bonusStr, "%")[0])
				if bonus, err := strconv.ParseFloat(bonusStr, 64); err == nil {
					tr.context.Variables["team_hp_bonus"] = bonus / 100.0
				}
			}
		}
	}
	return nil
}

// executeTeardown 执行清理
func (tr *TestRunner) executeTeardown(teardown []string) error {
	// TODO: 实现清理逻辑
	// 例如：清理战斗状态、重置角色数据等
	return nil
}

// executeCreateItem 创建物品
func (tr *TestRunner) executeCreateItem(instruction string) error {
	// 解析物品价格，如"创建一个物品，价格=30"
	price := 30 // 默认价格
	if strings.Contains(instruction, "价格=") {
		parts := strings.Split(instruction, "价格=")
		if len(parts) > 1 {
			priceStr := strings.TrimSpace(strings.Split(parts[1], "，")[0])
			priceStr = strings.TrimSpace(strings.Split(priceStr, ")")[0])
			if p, err := strconv.Atoi(priceStr); err == nil {
				price = p
			}
		}
	}
	tr.context.Variables["item_price"] = price
	tr.safeSetContext("item_price", price)
	return nil
}

// executePurchaseItem 购买物品
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
		parts := strings.Split(instruction, "价格=")
		if len(parts) > 1 {
			priceStr := strings.TrimSpace(strings.Split(parts[1], ")")[0])
			if p, err := strconv.Atoi(priceStr); err == nil {
				price = p
			}
		}
	}

	// 解析物品名称
	itemName := "物品A"
	if strings.Contains(instruction, "购买物品") {
		parts := strings.Split(instruction, "购买物品")
		if len(parts) > 1 {
			namePart := strings.TrimSpace(strings.Split(parts[1], ")")[0])
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

	// 检查金币是否足够
	if user.Gold < price {
		tr.context.Variables["purchase_success"] = false
		tr.safeSetContext("purchase_success", false)
		return fmt.Errorf("insufficient gold: need %d, have %d", price, user.Gold)
	}

	// 扣除金币
	user.Gold -= price
	_, err = database.DB.Exec(`UPDATE users SET gold = ? WHERE id = ?`, user.Gold, char.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user gold: %w", err)
	}

	// 标记角色拥有该物品
	itemKey := fmt.Sprintf("character.has_%s", strings.ToLower(strings.ReplaceAll(itemName, " ", "_")))
	tr.context.Variables[itemKey] = true
	tr.safeSetContext(itemKey, true)

	// 更新上下文
	tr.context.Variables["character.gold"] = user.Gold
	tr.safeSetContext("character.gold", user.Gold)
	tr.context.Variables["purchase_success"] = true
	tr.safeSetContext("purchase_success", true)

	return nil
}

// executeTryPurchaseItem 尝试购买物品（用于测试失败情况）
func (tr *TestRunner) executeTryPurchaseItem(instruction string) error {
	err := tr.executePurchaseItem(instruction)
	if err != nil {
		tr.context.Variables["purchase_success"] = false
		tr.safeSetContext("purchase_success", false)
		return nil // 不返回错误，因为这是预期的失败
	}
	return nil
}

// executeViewShopItems 查看商店物品
func (tr *TestRunner) executeViewShopItems() error {
	// 这个操作主要是为了测试，实际不需要做什么
	// 物品列表已经在initializeShop中设置了
	return nil
}

// executeGainGold 角色获得金币
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
			// 更新用户金币
			userRepo := repository.NewUserRepository()
			user, err := userRepo.GetByID(char.UserID)
			if err == nil && user != nil {
				user.Gold += gold
				_, err = database.DB.Exec(`UPDATE users SET gold = ?, total_gold_gained = total_gold_gained + ? WHERE id = ?`,
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

// executeInitializeShop 初始化商店
func (tr *TestRunner) executeInitializeShop(instruction string) error {
	// 解析商店物品，如"初始化商店，包含物品A（价格=50）"
	itemsCount := 0
	if strings.Contains(instruction, "包含") {
		if strings.Contains(instruction, "多个物品") {
			itemsCount = 3 // 默认3个物品
		} else if strings.Contains(instruction, "物品A") {
			itemsCount = 1
			// 解析价格
			if strings.Contains(instruction, "价格=") {
				parts := strings.Split(instruction, "价格=")
				if len(parts) > 1 {
					priceStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])
					if price, err := strconv.Atoi(priceStr); err == nil {
						tr.context.Variables["shop_item_a_price"] = price
						tr.safeSetContext("shop_item_a_price", price)
					}
				}
			}
		}
	}

	tr.context.Variables["shop.items_count"] = itemsCount
	tr.safeSetContext("shop.items_count", itemsCount)

	return nil
}

// 确保所有导入的包都被使用（避免编译错误）
var _ = fmt.Sprintf
var _ = math.Abs
var _ = rand.Intn
var _ = os.Getenv
var _ = filepath.Join
var _ = reflect.TypeOf
var _ = regexp.MustCompile
var _ = sort.Ints
var _ = strconv.Itoa
var _ = strings.Contains
var _ = time.Now
var _ = database.DB
var _ = game.GetBattleManager
var _ = models.Character{}
var _ = repository.NewCharacterRepository
var _ = yaml.Marshal
