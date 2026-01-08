package main

import (
	"fmt"
	"text-wow/internal/game"
)

func main() {
	fmt.Println("=== 后端系统验证 ===")
	
	// 检查BattleManager是否使用新功能
	manager := game.GetBattleManager()
	if manager == nil {
		fmt.Println("❌ BattleManager 未初始化")
		return
	}
	
	fmt.Println("✅ BattleManager 已初始化")
	
	// 检查是否有BuffManager
	buffMgr := manager.GetBuffManager()
	if buffMgr == nil {
		fmt.Println("❌ BuffManager 未初始化")
	} else {
		fmt.Println("✅ BuffManager 已初始化")
	}
	
	// 检查是否有SkillManager
	skillMgr := manager.GetSkillManager()
	if skillMgr == nil {
		fmt.Println("❌ SkillManager 未初始化")
	} else {
		fmt.Println("✅ SkillManager 已初始化")
	}
	
	// 检查是否有PassiveSkillManager
	passiveMgr := manager.GetPassiveSkillManager()
	if passiveMgr == nil {
		fmt.Println("❌ PassiveSkillManager 未初始化")
	} else {
		fmt.Println("✅ PassiveSkillManager 已初始化")
	}
	
	fmt.Println("\n=== 验证完成 ===")
	fmt.Println("如果所有组件都显示 ✅，说明重构后的系统正在运行")
}





















































