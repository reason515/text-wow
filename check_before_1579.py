#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')

# 检查第1573行之前的代码，查找未闭合的结构
print('\n=== 检查第1573行之前的代码 ===')
for i in range(max(0, 1550), 1580):
    line = lines[i].rstrip()
    if line.strip():
        print(f'{i+1:5d}: {line[:120]}')

# 检查括号匹配
print('\n=== 检查括号匹配 ===')
brace_count = 0
paren_count = 0
bracket_count = 0

for i in range(1500, 1580):
    line = lines[i]
    brace_count += line.count('{') - line.count('}')
    paren_count += line.count('(') - line.count(')')
    bracket_count += line.count('[') - line.count(']')
    
    if brace_count < 0 or paren_count < 0 or bracket_count < 0:
        print(f'第 {i+1} 行: 括号不匹配 - brace={brace_count}, paren={paren_count}, bracket={bracket_count}')
        print(f'  内容: {repr(line[:100])}')

print(f'\n第1580行时的括号状态: brace={brace_count}, paren={paren_count}, bracket={bracket_count}')

# 查找可能的未闭合方括号
print('\n=== 查找未闭合的方括号 ===')
for i in range(1500, 1580):
    line = lines[i]
    if '[' in line:
        # 检查是否在字符串中
        in_string = False
        escape = False
        for j, char in enumerate(line):
            if escape:
                escape = False
                continue
            if char == '\\':
                escape = True
                continue
            if char == '"':
                in_string = not in_string
            if not in_string and char == '[':
                # 检查是否有对应的 ]
                rest = line[j+1:]
                if ']' not in rest.split('"')[0] if '"' in rest else rest:
                    print(f'第 {i+1} 行可能有未闭合的 [: {repr(line[:100])}')
