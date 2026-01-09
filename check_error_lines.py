#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

print(f'文件总行数: {len(lines)}')

# 检查错误行
error_lines = [8995, 10099, 10243, 18483, 21044, 22356, 36277, 36629, 38293, 41717]

for err_line in error_lines:
    if err_line <= len(lines):
        line = lines[err_line - 1]
        print(f'\n第 {err_line} 行: {repr(line[:150])}')
        # 显示前后5行
        start = max(0, err_line - 6)
        end = min(len(lines), err_line + 4)
        print(f'上下文 ({start+1}-{end}):')
        for i in range(start, end):
            marker = '>>>' if i == err_line - 1 else '   '
            print(f'{marker} {i+1:5d}: {repr(lines[i][:100])}')
