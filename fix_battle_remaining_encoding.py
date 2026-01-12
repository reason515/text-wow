#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 battle.go 文件中剩余的编码问题
"""

def fix_remaining():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    skip_count = 0
    
    for i, line in enumerate(lines):
        # 跳过已处理的几行
        if skip_count > 0:
            skip_count -= 1
            continue
        
        # 修复第392-394行：删除旧的编码问题注释
        if i == 391 and '// 第4回合：冷却剩余0，可用' in line:
            fixed_lines.append(line)
            # 检查接下来的3行是否是旧的编码问题
            if i + 1 < len(lines) and '回合：冷却剩' in lines[i + 1]:
                # 跳过这3行，直接添加 cooldownLeft
                for j in range(i + 1, min(i + 5, len(lines))):
                    if 'cooldownLeft := cooldown' in lines[j]:
                        fixed_lines.append('					cooldownLeft := cooldown - (roundNum - 1)\n')
                        skip_count = j - i
                        break
                continue
        
        # 修复第424行
        if i == 423 and '处理怪物技能冷却时间' in line and 'if monsterSkillID' in line:
            fixed_lines.append('	// 处理怪物技能冷却时间（不再从Variables读取Skill对象，避免序列化错误）\n')
            fixed_lines.append('	if monsterSkillID, exists := tr.context.Variables["monster_skill_id"]; exists && monsterSkillID != nil {\n')
            continue
        
        # 修复第425行
        if i == 424 and '从Variables获取怪物技能冷却时' in line and 'monsterCooldown := 0' in line:
            fixed_lines.append('		// 从Variables获取怪物技能冷却时间\n')
            fixed_lines.append('		monsterCooldown := 0\n')
            continue
        
        # 修复第431行
        if i == 430 and '获取上次使用技能的回合' in line and 'lastUsedRound := 1' in line:
            fixed_lines.append('		// 获取上次使用技能的回合\n')
            fixed_lines.append('		lastUsedRound := 1\n')
            continue
        
        # 修复第495行
        if i == 494 and '更新上下文文' in line:
            fixed_lines.append('	// 更新上下文\n')
            continue
        
        # 修复第546行
        if i == 545 and '设置角色死亡状' in line and 'if char != nil' in line:
            fixed_lines.append('	// 设置角色死亡状态\n')
            fixed_lines.append('	if char != nil {\n')
            continue
        
        # 修复第554行
        if i == 553 and '计算经验奖励（基于怪物数量' in line:
            fixed_lines.append('			// 计算经验奖励（基于怪物数量）\n')
            continue
        
        # 修复第561行
        if i == 560 and '计算金币奖励（简化：每个怪物10-30金币' in line:
            fixed_lines.append('			// 计算金币奖励（简化：每个怪物10-30金币）\n')
            continue
        
        # 修复第571行
        if i == 570 and '设置team_total_exp（单角色时等于character.exp' in line:
            fixed_lines.append('			// 设置team_total_exp（单角色时等于character.exp）\n')
            continue
        
        # 修复第575行
        if i == 574 and '失败时，exp_gained和gold_gained' in line:
            fixed_lines.append('			// 失败时，exp_gained和gold_gained为0\n')
            continue
        
        # 修复第581行
        if i == 580 and '设置team_alive_count（单角色时，如果角色死亡则为0，否则为1' in line:
            fixed_lines.append('		// 设置team_alive_count（单角色时，如果角色死亡则为0，否则为1）\n')
            continue
        
        # 修复第598行
        if i == 597 and '如果角色是战士，确保怒气' in line:
            fixed_lines.append('		// 如果角色是战士，确保怒气为0\n')
            continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_remaining()
