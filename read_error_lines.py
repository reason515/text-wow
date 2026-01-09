#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')

error_lines = [1579, 1765, 1789, 3220, 3658, 3874, 6333, 6396, 6675, 7275]

for err_line in error_lines:
    if err_line <= len(lines):
        print(f'\n=== 第 {err_line} 行 ===')
        start = max(0, err_line - 5)
        end = min(len(lines), err_line + 4)
        for i in range(start, end):
            marker = '>>>' if i == err_line - 1 else '   '
            line_content = lines[i].rstrip()
            print(f'{marker} {i+1:5d}: {repr(line_content[:120])}')
