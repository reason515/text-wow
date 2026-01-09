#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复所有拆分后文件中的编码问题
"""

import re
import os

# 需要修复的文件列表
files_to_fix = [
    'server/internal/test/runner/battle.go',
    'server/internal/test/runner/character.go',
    'server/internal/test/runner/monster.go',
    'server/internal/test/runner/team.go',
    'server/internal/test/runner/equipment.go',
    'server/internal/test/runner/calculation.go',
    'server/internal/test/runner/instruction.go',
    'server/internal/test/runner/context.go',
    'server/internal/test/runner/test_runner.go',
    'server/internal/test/runner/types.go',
]

# 常见的乱码替换模式
replacements = [
    # 常见的中文字符乱码
    ('伤?', '伤害'),
    ('怪物?', '怪物）'),
    ('上下?', '上下文'),
    ('减成?', '减成'),
    ('加成?', '加成'),
    ('获?', '获得'),
    ('个角?', '个角色'),
    ('创建一个角?', '创建一个角色'),
    ('创建N个角?', '创建N个角色'),
    ('创建一?', '创建一个'),
    ('创建一?人队伍', '创建一个多人队伍'),
    ('牧?', '牧师'),
    ('法?', '法师'),
    ('排?', '排除'),
    ('包?', '包含'),
    ('指?', '指令'),
    ('处?', '处理'),
    ('敏?', '敏捷'),
    ('技?', '技能'),
    ('?', ''),  # 删除单独的乱码字符
]

# 特殊模式：修复fmt.Sprintf中的乱码
def fix_sprintf_encoding(content):
    # 修复fmt.Sprintf("...伤?...") -> fmt.Sprintf("...伤害...")
    content = re.sub(r'fmt\.Sprintf\("([^"]*?)伤\?([^"]*?)"', r'fmt.Sprintf("\1伤害\2"', content)
    # 修复其他常见的字符串乱码
    content = re.sub(r'fmt\.Sprintf\("([^"]*?)怪物\?([^"]*?)"', r'fmt.Sprintf("\1怪物）\2"', content)
    return content

# 修复注释和字符串中的乱码
def fix_comment_encoding(content):
    # 修复注释中的乱码
    content = re.sub(r'// ([^\\n]*?)伤\?', r'// \1伤害', content)
    content = re.sub(r'// ([^\\n]*?)怪物\?', r'// \1怪物）', content)
    content = re.sub(r'// ([^\\n]*?)上下\?', r'// \1上下文', content)
    return content

total_fixed = 0

for file_path in files_to_fix:
    if not os.path.exists(file_path):
        print(f'文件不存在: {file_path}')
        continue
    
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    content = content_bytes.decode('utf-8', errors='replace')
    original_content = content
    
    # 应用所有替换
    for old, new in replacements:
        if old in content:
            count = content.count(old)
            content = content.replace(old, new)
            if count > 0:
                print(f'{file_path}: 修复了 {count} 处 {repr(old)} -> {repr(new)}')
                total_fixed += count
    
    # 修复fmt.Sprintf中的乱码
    content = fix_sprintf_encoding(content)
    
    # 修复注释中的乱码
    content = fix_comment_encoding(content)
    
    # 修复特殊模式：注释和代码在同一行的情况
    # 例如：// 收集所有参与者（角色和怪物?	type participant struct {
    content = re.sub(r'// ([^\\n]*?)怪物\?\s+type\s+', r'// \1怪物）\n\ttype ', content)
    content = re.sub(r'// ([^\\n]*?)伤\?\s+', r'// \1伤害\n\t', content)
    
    # 如果内容有变化，写入文件
    if content != original_content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f'{file_path}: 已更新\n')
    else:
        print(f'{file_path}: 无需修复\n')

print(f'\n总共修复了 {total_fixed} 处编码问题')
