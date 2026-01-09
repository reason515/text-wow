#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 清理多余空行：最多保留2个连续空行
cleaned_lines = []
empty_count = 0
max_empty = 2  # 最多保留2个连续空行

for line in lines:
    if line.strip() == '':
        empty_count += 1
        if empty_count <= max_empty:
            cleaned_lines.append('')
    else:
        empty_count = 0
        cleaned_lines.append(line)

content = '\n'.join(cleaned_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(cleaned_lines)}')
print(f'移除了 {len(lines) - len(cleaned_lines)} 行空行')
