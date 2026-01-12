#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中重复引号的问题
"""

def fix_duplicate_quotes():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    # 修复所有重复引号
    import re
    # 修复 "）"）" 为 "）"
    content = re.sub(r'"）"）"', '"）"', content)
    # 修复 ","）" 为 ","
    content = re.sub(r'","）"', '","', content)
    # 修复其他可能的重复引号模式
    content = re.sub(r'"([^"]*)"）"', r'"\1"', content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_duplicate_quotes()
