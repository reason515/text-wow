#!/usr/bin/env python3
# -*- coding: utf-8 -*-

battle_file = 'server/internal/test/runner/battle.go'

with open(battle_file, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 修复所有乱码模式
replacements = [
    ('怪物?', '怪物）'),
    ('伤?', '伤害'),
    ('上下?', '上下文'),
    ('减成?', '减成'),
    ('加成?', '加成'),
    ('获?', '获得'),
]

for old, new in replacements:
    if old in content:
        content = content.replace(old, new)
        print(f'修复: {repr(old)} -> {repr(new)}')

# 修复第31行的特殊问题：注释后的type定义
content = content.replace('// 收集所有参与者（角色和怪物?', '// 收集所有参与者（角色和怪物）')

# 修复fmt.Sprintf中的乱码
import re
# 修复所有fmt.Sprintf中的"伤?"为"伤害"
content = re.sub(r'fmt\.Sprintf\("([^"]*?)伤\?([^"]*?)"', r'fmt.Sprintf("\1伤害\2"', content)

with open(battle_file, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
