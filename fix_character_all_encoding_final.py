#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
最终修复 character.go 文件中的所有编码问题
"""

def fix_all_encoding():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换所有 UTF-8 替换字符序列
    content_bytes = content_bytes.replace(b'\xef\xbf\xbd)[0]', b'\xef\xbc\x89")[0]')  # "）")[0]
    content_bytes = content_bytes.replace(b'\xef\xbf\xbd);', b'\xef\xbc\x89");')  # "）");
    content_bytes = content_bytes.replace(b'\xef\xbf\xbd"', b'\xef\xbc\x89"')  # "）"
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 修复所有 ")[0]) 为 "）")[0])
    import re
    content = re.sub(r'"\)\[0\]\)', '"）")[0])', content)
    
    # 修复重复引号
    content = re.sub(r'"）"）"', '"）"', content)
    content = re.sub(r'"%"）"', '"%"', content)
    content = re.sub(r'","）"', '","', content)
    
    # 修复注释中的编码问题
    content = re.sub(r'（[^）]*力量=20"[^）]*敏捷=10"[^）]*等）', '（如"力量=20"或"敏捷=10"等）', content)
    content = re.sub(r'（[^）]*1000（理论上暴击率会超过50%[^）]*）', '（如"1000（理论上暴击率会超过50%）"）', content)
    content = re.sub(r'物理暴击[^=]*=', '物理暴击率=', content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_all_encoding()
