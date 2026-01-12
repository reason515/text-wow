#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 battle.go 文件中所有编码问题
"""

def fix_encoding():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    # 修复所有编码问题
    replacements = [
        # 注释和代码混在一起的情况
        (r'// 解析回合数（[^）]*?）\troundNum := 1', '// 解析回合数（如"执行X回合"或"执行一个回合"）\n\troundNum := 1'),
        (r'// 减少技能冷却时[^\t]*?\tskillManager := game.NewSkillManager\(\)', '// 减少技能冷却时间\n\tskillManager := game.NewSkillManager()'),
        (r'// 先减少冷却时[^\t]*?\tskillManager.TickCooldowns\(char.ID\)', '// 先减少冷却时间\n\t\t\tskillManager.TickCooldowns(char.ID)'),
        (r'// 减少Buff持续时间（每回合[^\t]*?\tif buffDuration', '// 减少Buff持续时间（每回合减少1）\n\t\t\tif buffDuration'),
        (r'// 减少护盾持续时间（每回合[^\t]*?\tif shieldDuration', '// 减少护盾持续时间（每回合减少1）\n\t\t\tif shieldDuration'),
        (r'// 获取技能状态，检查是否可用（不再从Variables读取Skill对象，避免序列化错误[^\t]*?\tskillID, exists', '// 获取技能状态，检查是否可用（不再从Variables读取Skill对象，避免序列化错误）\n\t\t\tskillID, exists'),
        (r'// 如果技能状态不存在，从Variables获取冷却时间并计[^\t]*?\tcooldown := 0', '// 如果技能状态不存在，从Variables获取冷却时间并计算\n\t\t\t\t\tcooldown := 0'),
        (r'// 假设[^回]*?回合使用了技能，冷却时间[^，]*?，那么：\s*// [^回]*?回合：冷却剩[^，]*?，不可用\s*// [^回]*?回合：冷却剩[^，]*?，不可用\s*// [^回]*?回合：冷却剩[^，]*?，可[^\t]*?\tcooldownLeft', 
         '// 假设第1回合使用了技能，冷却时间3，那么：\n\t\t\t\t\t// 第2回合：冷却剩余2，不可用\n\t\t\t\t\t// 第3回合：冷却剩余1，不可用\n\t\t\t\t\t// 第4回合：冷却剩余0，可用\n\t\t\t\t\tcooldownLeft'),
        (r'// 如果角色没有技能，从上下文获取技能信息（不再从Variables读取Skill对象[^\t]*?\tif _, exists', '// 如果角色没有技能，从上下文获取技能信息（不再从Variables读取Skill对象）\n\t\t\tif _, exists'),
        (r'// 从Variables获取冷却时间并计[^\t]*?\tcooldown := 0', '// 从Variables获取冷却时间并计算\n\t\t\t\tcooldown := 0'),
        (r'// 检查是否处于休息状[^\t]*?\tisResting, exists', '// 检查是否处于休息状态\n\tisResting, exists'),
        (r'// 如果不在休息状态，先进入休息状[^\t]*?\tif err := tr.checkAndEnterRest\(\)', '// 如果不在休息状态，先进入休息状态\n\t\tif err := tr.checkAndEnterRest()'),
        (r'// 设置战斗状[^\t]*?\tif isVictory', '// 设置战斗状态\n\tif isVictory'),
        (r'// 检查是否应该进入休息状[^\t]*?\tif err := tr.checkAndEnterRest\(\)', '// 检查是否应该进入休息状态\n\t\tif err := tr.checkAndEnterRest()'),
        (r'// 解析剩余怪物数量（如"剩余2个怪物攻击角色"[^\t]*?\texpectedCount := 0', '// 解析剩余怪物数量（如"剩余2个怪物攻击角色"）\n\texpectedCount := 0'),
        (r'// 战士受到伤害时获得怒气（假设获[^）]*?）', '// 战士受到伤害时获得怒气（假设获得5点）'),
    ]
    
    import re
    for pattern, replacement in replacements:
        content = re.sub(pattern, replacement, content, flags=re.MULTILINE)
    
    # 修复缩进问题
    content = re.sub(r'^\t\t\t\t\t\t\t\t\t\tcooldownLeft := cooldown', '\t\t\t\t\t\tcooldownLeft := cooldown', content, flags=re.MULTILINE)
    content = re.sub(r'^\t\t\t\t\t\t\t\t\t\tif cooldownLeft < 0', '\t\t\t\t\t\tif cooldownLeft < 0', content, flags=re.MULTILINE)
    
    # 修复缩进问题：第402行的 cooldownLeft
    content = re.sub(r'^\t\t\t\t\t\t// 根据冷却时间计算\s*\n\t\t\t\t\t\tcooldownLeft', '\t\t\t\t// 根据冷却时间计算\n\t\t\t\tcooldownLeft', content, flags=re.MULTILINE)
    content = re.sub(r'^\t\t\t\t\t\tif cooldownLeft < 0', '\t\t\t\tif cooldownLeft < 0', content, flags=re.MULTILINE)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_encoding()
