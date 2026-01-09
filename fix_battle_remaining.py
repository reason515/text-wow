#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 修复第301行：注释和type定义在同一行
content = content.replace('// 收集所有参与者（角色和怪物）', '// 收集所有参与者（角色和怪物）')
# 查找并修复注释和type定义在同一行的问题
import re

# 修复注释和type定义在同一行的问题
content = re.sub(
    r'// 收集所有参与者（角色和怪物）[^\n]*?type participant struct',
    '// 收集所有参与者（角色和怪物）\n\ttype participant struct',
    content
)

# 修复strings.Contains中的乱码
content = re.sub(
    r'strings\.Contains\(instruction, "([^"]*?)\?\)',
    r'strings.Contains(instruction, "\1在")',
    content
)

# 修复strings.Split中的乱码
content = re.sub(
    r'strings\.Split\(instruction, "([^"]*?)\?\)',
    r'strings.Split(instruction, "\1在")',
    content
)

# 修复strings.Split中的另一个乱码模式
content = re.sub(
    r'strings\.Split\(parts\[1\], "([^"]*?)\?\)\[0\]',
    r'strings.Split(parts[1], "\1个")[0]',
    content
)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
