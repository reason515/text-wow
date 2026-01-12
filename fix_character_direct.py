#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
直接修复 character.go 文件中的编码问题
"""

def fix_direct():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复第622行：注释和代码混在一起
        if i == 621 and '// 解析攻击力（' in line and '	if strings.Contains' in line:
            fixed_lines.append('	// 解析攻击力（如"攻击=20"）\n')
            fixed_lines.append('	if strings.Contains(instruction, "攻击=") {\n')
            continue
        
        # 修复第4540行：编码问题
        if i == 4539 and '一个角' in line and ', 1)' in line:
            line = line.replace('一个角, 1)', '一个角色", 1)')
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_direct()
