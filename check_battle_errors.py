#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'Total lines: {len(lines)}')

# 检查第30-35行
print('\n=== Lines 30-35 ===')
for i in range(29, min(35, len(lines))):
    print(f'{i+1:4d}: {repr(lines[i][:100])}')

# 检查第200-210行
print('\n=== Lines 200-210 ===')
for i in range(199, min(210, len(lines))):
    print(f'{i+1:4d}: {repr(lines[i][:100])}')

# 检查第235-240行
print('\n=== Lines 235-240 ===')
for i in range(234, min(240, len(lines))):
    print(f'{i+1:4d}: {repr(lines[i][:100])}')

# 查找buildTurnOrder函数
print('\n=== Searching for buildTurnOrder ===')
for i, line in enumerate(lines):
    if 'buildTurnOrder' in line:
        print(f'Line {i+1}: {repr(line[:100])}')
        # 显示前后5行
        start = max(0, i-5)
        end = min(len(lines), i+6)
        for j in range(start, end):
            marker = '>>>' if j == i else '   '
            print(f'{marker} {j+1:4d}: {repr(lines[j][:100])}')
