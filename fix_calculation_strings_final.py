#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 calculation.go 文件中字符串分割的编码问题
"""

def fix_strings():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复所有包含 ")[0]) 的行，替换为 "）")[0])
        # 同时移除多余的空行
        if '"' in line and ')[0])' in line and 'strings.Split' in line:
            # 替换 ")[0]) 为 "）")[0])
            line = line.replace('")[0])', '"）")[0])')
            # 移除行尾的换行符，我们会在后面统一处理
            line = line.rstrip() + '\n'
        
        # 移除函数体内的多余空行（连续两个空行只保留一个）
        if not line.strip():
            if fixed_lines and not fixed_lines[-1].strip():
                # 如果上一行也是空行，跳过这一行
                continue
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_strings()
