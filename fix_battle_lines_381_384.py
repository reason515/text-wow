#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复第381-384行的特殊情况
"""

def fix_lines():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    skip_next = False
    
    for i, line in enumerate(lines):
        # 如果已经标记跳过，跳过这一行
        if skip_next:
            # 检查是否应该继续跳过
            if i < len(lines) - 1:
                next_line = lines[i + 1]
                if '回合：冷却剩' in next_line or 'cooldownLeft := cooldown' in line:
                    skip_next = True
                    continue
            skip_next = False
        
        # 修复第381-384行（0-based: 380-383）
        if i == 380 and '假设' in line and '回合使用了技能' in line:
            # 检查接下来的几行
            if i + 3 < len(lines):
                next3 = ''.join(lines[i:i+4])
                if 'cooldownLeft := cooldown' in next3:
                    # 替换这4行
                    fixed_lines.append('					// 假设第1回合使用了技能，冷却时间3，那么：\n')
                    fixed_lines.append('					// 第2回合：冷却剩余2，不可用\n')
                    fixed_lines.append('					// 第3回合：冷却剩余1，不可用\n')
                    fixed_lines.append('					// 第4回合：冷却剩余0，可用\n')
                    # 跳过接下来的3行
                    skip_next = True
                    for j in range(1, 4):
                        if i + j < len(lines):
                            if 'cooldownLeft := cooldown' not in lines[i + j]:
                                skip_next = False
                                break
                    if skip_next:
                        # 找到 cooldownLeft 行并添加
                        for j in range(1, 10):
                            if i + j < len(lines) and 'cooldownLeft := cooldown' in lines[i + j]:
                                fixed_lines.append(lines[i + j])
                                skip_next = False
                                break
                        continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_lines()
