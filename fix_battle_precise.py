#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
精确修复 battle.go 文件中的编码问题
"""

def fix_precise():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    for i, line in enumerate(lines, 1):
        # 修复第315行
        if i == 315 and '解析回合数（' in line and 'roundNum := 1' in line:
            fixed_lines.append('	// 解析回合数（如"执行X回合"或"执行一个回合"）\n')
            fixed_lines.append('	roundNum := 1\n')
            continue
        
        # 修复第335行
        if i == 335 and '减少技能冷却时' in line and 'skillManager := game.NewSkillManager()' in line:
            fixed_lines.append('	// 减少技能冷却时间\n')
            fixed_lines.append('	skillManager := game.NewSkillManager()\n')
            continue
        
        # 修复第340行
        if i == 340 and '先减少冷却时' in line and 'skillManager.TickCooldowns' in line:
            fixed_lines.append('			// 先减少冷却时间\n')
            fixed_lines.append('			skillManager.TickCooldowns(char.ID)\n')
            continue
        
        # 修复第341行
        if i == 341 and '减少Buff持续时间（每回合' in line and 'if buffDuration' in line:
            fixed_lines.append('			// 减少Buff持续时间（每回合减少1）\n')
            fixed_lines.append('			if buffDuration, exists := tr.context.Variables["character_buff_duration"]; exists {\n')
            continue
        
        # 修复第353行
        if i == 353 and '减少护盾持续时间（每回合' in line and 'if shieldDuration' in line:
            fixed_lines.append('			// 减少护盾持续时间（每回合减少1）\n')
            fixed_lines.append('			if shieldDuration, exists := tr.context.Variables["character.shield_duration"]; exists {\n')
            continue
        
        # 修复第365行
        if i == 365 and '获取技能状态' in line and 'skillID, exists' in line:
            fixed_lines.append('			// 获取技能状态，检查是否可用（不再从Variables读取Skill对象，避免序列化错误）\n')
            fixed_lines.append('			skillID, exists := tr.context.Variables["skill_id"]\n')
            continue
        
        # 修复第375行
        if i == 375 and '如果技能状态不存在' in line and 'cooldown := 0' in line:
            fixed_lines.append('					// 如果技能状态不存在，从Variables获取冷却时间并计算\n')
            fixed_lines.append('					cooldown := 0\n')
            continue
        
        # 修复第381-384行
        if i == 381 and '假设' in line and '回合使用了技能' in line:
            fixed_lines.append('					// 假设第1回合使用了技能，冷却时间3，那么：\n')
            fixed_lines.append('					// 第2回合：冷却剩余2，不可用\n')
            fixed_lines.append('					// 第3回合：冷却剩余1，不可用\n')
            fixed_lines.append('					// 第4回合：冷却剩余0，可用\n')
            # 跳过第382-384行
            i += 3
            continue
        
        # 修复第395行
        if i == 395 and '如果角色没有技能' in line and 'if _, exists' in line:
            fixed_lines.append('			// 如果角色没有技能，从上下文获取技能信息（不再从Variables读取Skill对象）\n')
            fixed_lines.append('			if _, exists := tr.context.Variables["skill_id"]; exists {\n')
            continue
        
        # 修复第396行
        if i == 396 and '从Variables获取冷却时间并计' in line and 'cooldown := 0' in line:
            fixed_lines.append('				// 从Variables获取冷却时间并计算\n')
            fixed_lines.append('				cooldown := 0\n')
            continue
        
        # 修复第465行
        if i == 465 and '检查是否处于休息状' in line and 'isResting, exists' in line:
            fixed_lines.append('	// 检查是否处于休息状态\n')
            fixed_lines.append('	isResting, exists := tr.context.Variables["is_resting"]\n')
            continue
        
        # 修复第468行
        if i == 468 and '如果不在休息状态' in line and 'if err := tr.checkAndEnterRest()' in line:
            fixed_lines.append('		// 如果不在休息状态，先进入休息状态\n')
            fixed_lines.append('		if err := tr.checkAndEnterRest(); err != nil {\n')
            continue
        
        # 修复第491行
        if i == 491 and '设置战斗状' in line and 'if isVictory' in line:
            fixed_lines.append('	// 设置战斗状态\n')
            fixed_lines.append('	if isVictory {\n')
            continue
        
        # 修复第502行
        if i == 502 and '检查是否应该进入休息状' in line and 'if err := tr.checkAndEnterRest()' in line:
            fixed_lines.append('		// 检查是否应该进入休息状态\n')
            fixed_lines.append('		if err := tr.checkAndEnterRest(); err != nil {\n')
            continue
        
        # 修复第438行
        if i == 438 and '解析剩余怪物数量' in line and 'expectedCount := 0' in line:
            fixed_lines.append('	// 解析剩余怪物数量（如"剩余2个怪物攻击角色"）\n')
            fixed_lines.append('	expectedCount := 0\n')
            continue
        
        # 修复第300行
        if i == 300 and '假设获' in line and '点）' in line:
            fixed_lines.append('	// 战士受到伤害时获得怒气（假设获得5点）\n')
            continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_precise()
