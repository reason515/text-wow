#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
直接修复 calculation.go 文件中的编码问题
"""

def fix_direct():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复第89行
        if i == 88 and 'gainStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('				gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        # 修复第90行
        if i == 89 and 'gainStr = strings.TrimSpace(strings.Split(gainStr, "' in line:
            fixed_lines.append('				gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])\n')
            continue
        
        # 修复第126行
        if i == 125 and '允许char为nil（用于测试nil情况' in line:
            fixed_lines.append('	// 允许char为nil（用于测试nil情况）\n')
            continue
        
        # 修复第127行
        if i == 126 and '解析基础恢复值（' in line:
            fixed_lines.append('	// 解析基础恢复值（如"计算法力恢复（基础恢复=10）"）\n')
            fixed_lines.append('	baseRegen := 0\n')
            continue
        
        # 修复第131行
        if i == 130 and 'regenStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('			regenStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        # 修复第132行
        if i == 131 and 'regenStr = strings.TrimSpace(strings.Split(regenStr, "' in line:
            fixed_lines.append('			regenStr = strings.TrimSpace(strings.Split(regenStr, "）")[0])\n')
            continue
        
        # 修复第139行
        if i == 138 and '解析基础获得值（' in line:
            fixed_lines.append('	// 解析基础获得值（如"计算怒气获得（基础获得=10）"）\n')
            fixed_lines.append('	baseGain := 0\n')
            continue
        
        # 修复第143行
        if i == 142 and 'gainStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('			gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        # 修复第144行
        if i == 143 and 'gainStr = strings.TrimSpace(strings.Split(gainStr, "' in line:
            fixed_lines.append('			gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])\n')
            continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_direct()
