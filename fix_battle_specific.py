#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 修复第91行的问题：注释和type定义在同一行
# 查找模式：// 收集所有参与者（角色和怪物?	type participant struct {
import re

# 修复注释和type定义在同一行的问题
content = re.sub(
    r'// 收集所有参与者（角色和怪物[^\n]*?\)\s*type participant struct \{',
    '// 收集所有参与者（角色和怪物）\n\ttype participant struct {',
    content
)

# 修复fmt.Sprintf中的乱码
content = re.sub(
    r'fmt\.Sprintf\("角色攻击怪物，造成%d点伤([^"]*?)"',
    r'fmt.Sprintf("角色攻击怪物，造成%d点伤害"',
    content
)

content = re.sub(
    r'fmt\.Sprintf\("怪物攻击角色，造成%d点伤([^"]*?)"',
    r'fmt.Sprintf("怪物攻击角色，造成%d点伤害"',
    content
)

# 修复其他可能的fmt.Sprintf乱码
content = re.sub(
    r'fmt\.Sprintf\("([^"]*?)伤([^"]*?)"',
    lambda m: f'fmt.Sprintf("{m.group(1)}伤害{m.group(2)}"' if '伤' in m.group(1) or '伤' in m.group(2) else m.group(0),
    content
)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
