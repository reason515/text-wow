#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
全面修复 character.go 文件中的所有编码问题
"""

def fix_comprehensive():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换所有 UTF-8 替换字符序列
    replacements = [
        (b'\xef\xbf\xbd)[0]', b'\xef\xbc\x89")[0]'),  # "）")[0]
        (b'\xef\xbf\xbd);', b'\xef\xbc\x89");'),  # "）");
        (b'\xef\xbf\xbd"', b'\xef\xbc\x89"'),  # "）"
        (b'\xef\xbf\xbd=', b'='),  # 移除替换字符
    ]
    
    for old, new in replacements:
        content_bytes = content_bytes.replace(old, new)
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 修复所有字符串分割中的编码问题
    import re
    content = re.sub(r'"\)\[0\]\)', '"）")[0])', content)
    
    # 修复重复引号
    content = re.sub(r'"）"）"', '"）"', content)
    content = re.sub(r'"%"）"', '"%"', content)
    content = re.sub(r'","）"', '","', content)
    
    # 修复注释和代码混在一起的情况
    fixes = [
        (r'// 解析主属性（[^）]*力量=20"[^）]*敏捷=10"[^）]*等）', '// 解析主属性（如"力量=20"或"敏捷=10"等）'),
        (r'// 去掉括号和注释（[^）]*1000（理论上暴击率会超过50%[^）]*）[^\n]*\t\tif idx := strings.Index\(value, "', 
         '// 去掉括号和注释（如"1000（理论上暴击率会超过50%）"）\n\t\tif idx := strings.Index(value, "（'),
        (r'// 更新已存在角色的ClassID（如果指令中指定了不同的职业[^\n]*\t\t\tif classIDVal', 
         '// 更新已存在角色的ClassID（如果指令中指定了不同的职业）\n\t\t\tif classIDVal'),
    ]
    
    for pattern, replacement in fixes:
        content = re.sub(pattern, replacement, content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_comprehensive()
