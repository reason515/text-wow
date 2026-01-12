#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
使用正则表达式修复 calculation.go 文件中的编码问题
"""

import re

def fix_regex():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    # 修复字符串分割中的编码问题
    content = re.sub(r'strings\.Split\(parts\[1\], "[^"]*\)\[0\]\)', 
                     lambda m: m.group(0).replace('")[0])', '"）")[0])') if '")[0])' in m.group(0) else m.group(0),
                     content)
    
    # 更精确的修复：替换所有 ")[0]) 为 "）")[0])
    content = re.sub(r'"\)\[0\]\)', '"）")[0])', content)
    
    # 修复注释和代码混在一起的情况
    content = re.sub(r'// 允许char为nil（用于测试nil情况[^\n]*\n// 解析基础恢复值（[^\n]*\tbaseRegen := 0',
                     '// 允许char为nil（用于测试nil情况）\n\t// 解析基础恢复值（如"计算法力恢复（基础恢复=10）"）\n\tbaseRegen := 0',
                     content)
    
    content = re.sub(r'// 解析基础获得值（[^\n]*\tbaseGain := 0',
                     '// 解析基础获得值（如"计算怒气获得（基础获得=10）"）\n\tbaseGain := 0',
                     content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_regex()
