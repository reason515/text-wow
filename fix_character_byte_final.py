#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
字节级别修复 character.go 文件
"""

def fix_byte_final():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换 UTF-8 替换字符序列
    content_bytes = content_bytes.replace(b'\xef\xbf\xbd)[0]', b'\xef\xbc\x89")[0]')  # "）")[0]
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 修复重复引号
    import re
    content = re.sub(r'"）"）"', '"）"', content)
    content = re.sub(r'","）"', '","', content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_byte_final()
