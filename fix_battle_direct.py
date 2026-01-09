#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')

# 直接替换所有乱码字符串
# 修复 strings.Contains(instruction, "?) 为 strings.Contains(instruction, "在")
content = content.replace('strings.Contains(instruction, "?)', 'strings.Contains(instruction, "在")')

# 修复 strings.Split(instruction, "?) 为 strings.Split(instruction, "在")
content = content.replace('strings.Split(instruction, "?)', 'strings.Split(instruction, "在")')

# 修复 strings.Split(parts[1], "?)[0] 为 strings.Split(parts[1], "个")[0]
content = content.replace('strings.Split(parts[1], "?)[0]', 'strings.Split(parts[1], "个")[0]')

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
