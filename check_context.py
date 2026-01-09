#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
# 统一换行符
content = content.replace('\r\n', '\n').replace('\r', '\n')
lines = content.split('\n')

print(f'文件总行数: {len(lines)}')

# 检查错误行 - 查找实际有内容的行
def find_non_empty_around(line_num, context=20):
    """查找指定行号附近有内容的行"""
    start = max(0, line_num - context - 1)
    end = min(len(lines), line_num + context)
    non_empty = []
    for i in range(start, end):
        if lines[i].strip():
            non_empty.append((i+1, lines[i].strip()[:150]))
    return non_empty

# 检查第8995行
print('\n=== 检查第8995行 (executeTeardown) ===')
for num, line in find_non_empty_around(8995, 30):
    marker = '>>>' if num == 8995 else '   '
    print(f'{marker} {num:5d}: {line}')

# 检查第10099行
print('\n=== 检查第10099行 (else) ===')
for num, line in find_non_empty_around(10099, 30):
    marker = '>>>' if num == 10099 else '   '
    print(f'{marker} {num:5d}: {line}')

# 检查第10243行
print('\n=== 检查第10243行 (safeSetContext) ===')
for num, line in find_non_empty_around(10243, 30):
    marker = '>>>' if num == 10243 else '   '
    print(f'{marker} {num:5d}: {line}')

# 检查第36277行
print('\n=== 检查第36277行 (else) ===')
for num, line in find_non_empty_around(36277, 30):
    marker = '>>>' if num == 36277 else '   '
    print(f'{marker} {num:5d}: {line}')

# 检查第41717行
print('\n=== 检查第41717行 (break) ===')
for num, line in find_non_empty_around(41717, 30):
    marker = '>>>' if num == 41717 else '   '
    print(f'{marker} {num:5d}: {line}')
