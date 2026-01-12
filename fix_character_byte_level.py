#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
字节级别修复 character.go 文件
"""

def fix_byte_level():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换 UTF-8 替换字符序列
    # 修复 "一个角, 1) 为 "一个角色", 1)
    content_bytes = content_bytes.replace(b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xef\xbf\xbd, 1)', b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xe8\x89\xb2", 1)')
    
    # 修复注释和代码混在一起的问题
    # 修复 "基础HP=35"	if 为 "基础HP=35"）\n	if
    content_bytes = content_bytes.replace(b'\xe5\x9f\xba\xe7\xa1\x80HP=35"\xef\xbf\xbd\tif', b'\xe5\x9f\xba\xe7\xa1\x80HP=35"）\n\tif')
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_byte_level()
