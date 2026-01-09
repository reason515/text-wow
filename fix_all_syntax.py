#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')
fixed_count = 0

# 1. 修复第1579行 - 检查前面是否有未闭合的方括号
if len(lines) > 1578:
    # 检查第1570-1579行之间
    for i in range(1570, 1579):
        line = lines[i]
        # 查找未闭合的方括号
        if '[' in line and line.count('[') > line.count(']'):
            # 检查是否在字符串中
            in_string = False
            for char in line:
                if char == '"' and (i == 0 or line[i-1] != '\\'):
                    in_string = not in_string
                if not in_string and char == '[':
                    # 在行末添加 ]
                    lines[i] = line.rstrip() + ']' + '\n'
                    fixed_count += 1
                    print(f'修复了第 {i+1} 行的未闭合方括号')
                    break

# 2. 修复第1765行 - 检查前面是否有未闭合的if
if len(lines) > 1764:
    # 检查第1750-1765行之间
    brace_count = 0
    found_if = False
    for i in range(1750, 1765):
        line = lines[i]
        if 'if ' in line and '{' in line:
            found_if = True
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if found_if and brace_count > 0:
        # 在第1765行之前添加闭合大括号
        lines.insert(1764, '\t\t}\n')
        fixed_count += 1
        print(f'修复了第1765行之前的未闭合if语句')

# 3. 修复第1789行 - 检查是否在函数体内
if len(lines) > 1788:
    # 检查第1770-1789行之间是否有函数定义
    brace_count = 0
    func_start = -1
    for i in range(1770, 1789):
        line = lines[i]
        if 'func ' in line and '{' in line:
            func_start = i
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if brace_count > 0 and func_start >= 0:
        # 在第1789行之前添加闭合大括号
        lines.insert(1788, '\t}\n')
        fixed_count += 1
        print(f'修复了第1789行之前的未闭合函数')

# 4. 修复第3220行 - 检查是否在函数体内
if len(lines) > 3219:
    brace_count = 0
    func_start = -1
    for i in range(3200, 3220):
        line = lines[i]
        if 'func ' in line and '{' in line:
            func_start = i
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if brace_count > 0 and func_start >= 0:
        lines.insert(3219, '\t}\n')
        fixed_count += 1
        print(f'修复了第3220行之前的未闭合函数')

# 5. 修复第6333行 - 检查前面是否有未闭合的if
if len(lines) > 6332:
    brace_count = 0
    found_if = False
    for i in range(6320, 6333):
        line = lines[i]
        if 'if ' in line and '{' in line:
            found_if = True
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if found_if and brace_count > 0:
        lines.insert(6332, '\t}\n')
        fixed_count += 1
        print(f'修复了第6333行之前的未闭合if语句')

# 6. 修复第7275行 - 检查break是否在循环中
if len(lines) > 7274:
    line_7275 = lines[7274]
    if 'break' in line_7275:
        # 检查前面是否有for/switch/select
        found_loop = False
        brace_count = 0
        for i in range(max(0, 7250), 7275):
            line = lines[i]
            if 'for ' in line or 'switch ' in line or 'select ' in line:
                found_loop = True
                brace_count = 0
            brace_count += line.count('{') - line.count('}')
            if found_loop and brace_count == 0:
                # break不在循环中，删除或注释掉
                lines[7274] = '\t\t\t\t// break  // 修复：break不在循环中\n'
                fixed_count += 1
                print(f'修复了第7275行的break语句（不在循环中）')
                break

# 写入修复后的文件
if fixed_count > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(lines)
    print(f'\n修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
