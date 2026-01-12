#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整修复 character.go 文件中的所有编码问题
"""

def fix_complete():
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
        
        # 修复所有字符串分割中的编码问题
        if 'strings.Split(instruction, "' in line and '"' not in line.split('strings.Split(instruction, "')[1].split('"')[0] if '"' in line.split('strings.Split(instruction, "')[1] else '':
            # 检查是否有未闭合的引号
            if line.count('"') % 2 != 0:
                if '创建多个角色' in line or ':' in line:
                    line = line.replace('strings.Split(instruction, "', 'strings.Split(instruction, "创建多个角色:"')
                elif '攻击' in line:
                    line = line.replace('strings.Split(instruction, "', 'strings.Split(instruction, "攻击=")
                elif '防御' in line:
                    line = line.replace('strings.Split(instruction, "', 'strings.Split(instruction, "防御=")
        
        # 修复 parts[1] 分割中的编码问题
        if 'strings.Split(parts[1], "' in line and line.count('"') % 2 != 0:
            if ',' in line or '，' in line:
                line = line.replace('strings.Split(parts[1], "', 'strings.Split(parts[1], ","')
            elif '）' in line:
                line = line.replace('strings.Split(parts[1], "', 'strings.Split(parts[1], "）")')
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_complete()
