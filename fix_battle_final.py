#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

print(f'Total lines: {len(lines)}')

# 查找buildTurnOrder函数
for i, line in enumerate(lines):
    if 'func (tr *TestRunner) buildTurnOrder' in line:
        print(f'\nbuildTurnOrder found at line {i+1}')
        # 显示前后10行
        start = max(0, i-5)
        end = min(len(lines), i+15)
        for j in range(start, end):
            marker = '>>>' if j == i else '   '
            print(f'{marker} {j+1:4d}: {repr(lines[j][:80])}')

# 查找第202行附近的问题
print('\n=== Lines 200-210 ===')
for i in range(199, min(210, len(lines))):
    print(f'{i+1:4d}: {repr(lines[i][:100])}')

# 查找包含battle_logs的行
print('\n=== Searching for battle_logs ===')
for i, line in enumerate(lines):
    if 'battle_logs' in line:
        print(f'Line {i+1}: {repr(line[:100])}')
        # 显示前后3行
        start = max(0, i-3)
        end = min(len(lines), i+4)
        for j in range(start, end):
            marker = '>>>' if j == i else '   '
            print(f'{marker} {j+1:4d}: {repr(lines[j][:100])}')
