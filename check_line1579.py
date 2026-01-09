#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')

# 检查第1579行
print('\n=== 第1579行及其上下文 ===')
for i in range(max(0, 1570), min(len(lines), 1590)):
    line = lines[i].rstrip()
    marker = '>>>' if i == 1578 else '   '
    print(f'{marker} {i+1:5d}: {repr(line)}')

# 检查第1579行是否有方括号问题
line_1579 = lines[1578]
print(f'\n第1579行内容: {repr(line_1579)}')
print(f'包含 [: {"[" in line_1579}')
print(f'包含 ]: {"]" in line_1579}')
print(f'[ 数量: {line_1579.count("[")}')
print(f'] 数量: {line_1579.count("]")}')

# 检查前面的代码是否有未闭合的方括号
print('\n检查前面50行是否有未闭合的方括号:')
for i in range(max(0, 1530), 1579):
    line = lines[i]
    if '[' in line:
        open_count = line.count('[')
        close_count = line.count(']')
        if open_count > close_count:
            print(f'第 {i+1} 行可能有未闭合的 [: {repr(line[:100])}')
