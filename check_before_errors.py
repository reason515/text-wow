#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')

# 检查第1579行之前
print('\n=== 检查第1579行之前 (executeTeardown) ===')
for i in range(max(0, 1570), 1580):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')

# 检查第1765行之前
print('\n=== 检查第1765行之前 (else) ===')
for i in range(max(0, 1755), 1766):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')

# 检查第1789行之前
print('\n=== 检查第1789行之前 (safeSetContext) ===')
for i in range(max(0, 1780), 1790):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')

# 检查第3220行之前
print('\n=== 检查第3220行之前 ===')
for i in range(max(0, 3210), 3221):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')

# 检查第6333行之前
print('\n=== 检查第6333行之前 (else) ===')
for i in range(max(0, 6325), 6334):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')

# 检查第7275行之前
print('\n=== 检查第7275行之前 (break) ===')
for i in range(max(0, 7265), 7276):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:100]}')
