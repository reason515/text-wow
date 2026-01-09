#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    content = f.read()

fixed_count = 0

# 常见的乱码模式修复
replacements = [
    ('创建一�?', '创建一个'),
    ('个角�?', '个角色'),
    ('创建一个角�?', '创建一个角色'),
    ('创建N个角�?', '创建N个角色'),
    ('�?', '在'),
    ('敏�?', '敏捷'),
    ('速度=60�?', '速度=60'),
    ('�?', '个'),
    ('技�?', '技能'),
    ('创建一�?人队伍', '创建一个多人队伍'),
    ('牧�?', '牧师'),
    ('法�?', '法师'),
    ('排�?', '排除'),
    ('包�?', '包含'),
    ('指�?', '指令'),
    ('处�?', '处理'),
    ('�?', '一'),
    ('�?', '个'),
]

for old, new in replacements:
    count = content.count(old)
    if count > 0:
        content = content.replace(old, new)
        fixed_count += count
        print(f'修复了 {count} 处: {repr(old)} -> {repr(new)}')

if fixed_count > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'\n修复完成！共修复了 {fixed_count} 处乱码')
else:
    print('没有发现需要修复的乱码')
