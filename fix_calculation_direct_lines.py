#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
直接修复 calculation.go 文件中特定行的编码问题
"""

def fix_direct():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复第95-96行（0-based: 94-95）
        if i == 94 and 'gainStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('				gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        if i == 95 and 'gainStr = strings.TrimSpace(strings.Split(gainStr, "' in line:
            fixed_lines.append('				gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])\n')
            continue
        
        # 修复第150-151行（0-based: 149-150）
        if i == 149 and 'regenStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('			regenStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        if i == 150 and 'regenStr = strings.TrimSpace(strings.Split(regenStr, "' in line:
            fixed_lines.append('			regenStr = strings.TrimSpace(strings.Split(regenStr, "）")[0])\n')
            continue
        
        # 修复第166-167行（0-based: 165-166）
        if i == 165 and 'gainStr := strings.TrimSpace(strings.Split(parts[1], "' in line:
            fixed_lines.append('			gainStr := strings.TrimSpace(strings.Split(parts[1], "）")[0])\n')
            continue
        
        if i == 166 and 'gainStr = strings.TrimSpace(strings.Split(gainStr, "' in line:
            fixed_lines.append('			gainStr = strings.TrimSpace(strings.Split(gainStr, "）")[0])\n')
            continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_direct()
