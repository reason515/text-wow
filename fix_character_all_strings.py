#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中所有字符串分割的编码问题
"""

def fix_all_strings():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    # 使用正则表达式替换所有 ")[0]) 为 "）")[0])
    import re
    content = re.sub(r'"\)\[0\]\)', '"）")[0])', content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_all_strings()
