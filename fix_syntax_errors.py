#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')
original_lines = lines[:]

fixed_count = 0

# 统计实际行数（包括空行）
print(f'文件总行数: {len(lines)}')

# 检查错误行号是否超出范围
error_lines = [8995, 10099, 10243, 18483, 21044, 22356, 36277, 36629, 38293, 41717]
for err_line in error_lines:
    if err_line > len(lines):
        print(f'警告: 错误行号 {err_line} 超出文件范围 (文件只有 {len(lines)} 行)')
    else:
        print(f'检查第 {err_line} 行: {repr(lines[err_line-1][:100])}')

# 由于文件只有6631行，错误行号可能不准确
# 让我们搜索可能的问题模式

# 1. 查找未闭合的括号和方括号
print('\n搜索可能的语法错误...')

# 查找可能的未闭合结构
for i, line in enumerate(lines):
    # 检查是否有未闭合的方括号
    if '[' in line and line.count('[') > line.count(']'):
        # 检查是否在字符串中
        if not (line.count('"') % 2 == 1 or line.count("'") % 2 == 1):
            print(f'第 {i+1} 行可能有未闭合的 [: {repr(line[:100])}')
    
    # 检查是否有未闭合的括号
    if '(' in line and line.count('(') > line.count(')'):
        # 检查是否在字符串中
        if not (line.count('"') % 2 == 1 or line.count("'") % 2 == 1):
            # 跳过函数定义
            if 'func ' not in line and 'if ' not in line and 'for ' not in line:
                print(f'第 {i+1} 行可能有未闭合的 (: {repr(line[:100])}')

# 2. 查找可能的else语句问题
for i, line in enumerate(lines):
    stripped = line.strip()
    if stripped == 'else {' or stripped.startswith('} else'):
        # 检查前面是否有if语句
        found_if = False
        for j in range(max(0, i-50), i):
            if 'if ' in lines[j] and '{' in lines[j]:
                found_if = True
                break
        if not found_if:
            print(f'第 {i+1} 行可能有孤立的else: {repr(line[:100])}')

# 3. 查找可能的break语句问题
for i, line in enumerate(lines):
    if 'break' in line and '//' not in line.split('break')[0]:
        # 检查是否在循环中
        found_loop = False
        brace_count = 0
        for j in range(max(0, i-200), i):
            if 'for ' in lines[j] or 'switch ' in lines[j] or 'select ' in lines[j]:
                brace_count = 0
            brace_count += lines[j].count('{') - lines[j].count('}')
            if brace_count > 0 and ('for ' in lines[j] or 'switch ' in lines[j] or 'select ' in lines[j]):
                found_loop = True
                break
        if not found_loop and 'break' in line:
            print(f'第 {i+1} 行可能有break不在循环中: {repr(line[:100])}')

print('\n检查完成')
