#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
最终修复 calculation.go 文件中的编码问题
"""

def fix_final():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 将内容解码为字符串，处理编码错误
    try:
        content_str = content.decode('utf-8')
    except:
        content_str = content.decode('utf-8', errors='replace')
    
    lines = content_str.split('\n')
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复所有包含 ")[0]) 的行，替换为 "）")[0])
        if '"' in line and ')[0])' in line and 'strings.Split' in line:
            # 替换 ")[0]) 为 "）")[0])
            line = line.replace('")[0])', '"）")[0])')
        
        # 修复注释和代码混在一起的情况
        if '// 允许char为nil（用于测试nil情况' in line and 'baseRegen := 0' not in line:
            fixed_lines.append('	// 允许char为nil（用于测试nil情况）')
            continue
        
        if '// 解析基础恢复值（' in line and 'baseRegen := 0' in line:
            fixed_lines.append('	// 解析基础恢复值（如"计算法力恢复（基础恢复=10）"）')
            fixed_lines.append('	baseRegen := 0')
            continue
        
        if '// 解析基础获得值（' in line and 'baseGain := 0' in line:
            fixed_lines.append('	// 解析基础获得值（如"计算怒气获得（基础获得=10）"）')
            fixed_lines.append('	baseGain := 0')
            continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write('\n'.join(fixed_lines))
        if not fixed_lines[-1].endswith('\n'):
            f.write('\n')
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final()
