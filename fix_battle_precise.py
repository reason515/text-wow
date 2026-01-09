#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

# 修复第3194行
for i, line in enumerate(lines):
    if i == 3193:  # 第3194行（0-based索引）
        # 查找并修复 strings.Contains(instruction, "?) 为 strings.Contains(instruction, "在")
        if 'strings.Contains(instruction, "' in line:
            # 替换乱码字符
            line = line.replace('"?)', '"在")')
            lines[i] = line
            print(f'Line {i+1}: Fixed strings.Contains')
    
    if i == 3203:  # 第3204行
        # 修复 strings.Split(instruction, "?) 为 strings.Split(instruction, "在")
        if 'strings.Split(instruction, "' in line:
            line = line.replace('"?)', '"在")')
            lines[i] = line
            print(f'Line {i+1}: Fixed strings.Split')
    
    if i == 4463:  # 第4464行
        # 修复 strings.Split(parts[1], "?)[0] 为 strings.Split(parts[1], "个")[0]
        if 'strings.Split(parts[1], "' in line:
            line = line.replace('"?)', '"个")')
            lines[i] = line
            print(f'Line {i+1}: Fixed strings.Split parts')

content = '\n'.join(lines)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print('修复完成')
