#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 修复第151行：注释和type定义在同一行
content = content.replace('// 收集所有参与者（角色和怪物', '// 收集所有参与者（角色和怪物）')
content = content.replace('怪物	type participant struct {', '怪物）\n\ttype participant struct {')

# 修复第1597行和第1602行的字符串乱码
import re

# 修复 strings.Contains(instruction, "?) 为 strings.Contains(instruction, "在")
content = re.sub(r'strings\.Contains\(instruction, "([^"]*?)\?\)', r'strings.Contains(instruction, "\1在")', content)

# 修复 strings.Split(instruction, "?) 为 strings.Split(instruction, "在")
content = re.sub(r'strings\.Split\(instruction, "([^"]*?)\?\)', r'strings.Split(instruction, "\1在")', content)

# 修复其他可能的字符串乱码
content = re.sub(r'strings\.Contains\(instruction, "([^"]*?)\?\)', r'strings.Contains(instruction, "\1")', content)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
