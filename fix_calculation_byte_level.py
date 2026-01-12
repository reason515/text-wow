#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
字节级别修复 calculation.go 文件
"""

def fix_byte_level():
    file_path = 'server/internal/test/runner/calculation.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 替换所有 ")[0]) 为 "）")[0])
    # 使用多种方式尝试
    replacements = [
        ('")[0])', '"）")[0])'),
        ('" )[0])', '"）")[0])'),
        ('" )[0] )', '"）")[0])'),
    ]
    
    for old, new in replacements:
        content = content.replace(old, new)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_byte_level()
