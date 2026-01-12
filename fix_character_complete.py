#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整修复 character.go 文件中的编码问题
"""

def fix_complete():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复第63行：注释和代码混在一起
        if '// 保存当前指令到上下文，以便后续判断是否明确设置了某些属' in line and 'tr.context.Variables["last_instruction"]' in line:
            fixed_lines.append('	// 保存当前指令到上下文，以便后续判断是否明确设置了某些属性\n')
            fixed_lines.append('	tr.context.Variables["last_instruction"] = instruction\n')
            continue
        
        # 修复所有字符串分割中的编码问题
        if '"' in line and ')[0])' in line and 'strings.Split' in line:
            line = line.replace('")[0])', '"）")[0])')
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_complete()
