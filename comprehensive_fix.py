#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')
fixed_count = 0
original_lines = lines[:]

# 修复策略：检查每个错误行之前的代码，查找未闭合的结构

# 1. 修复第1579行 - 检查前面是否有未闭合的方括号
if len(lines) > 1578:
    # 从第1550行开始检查
    for i in range(1550, 1579):
        line = lines[i]
        # 检查是否有未闭合的方括号（不在字符串中）
        if '[' in line:
            # 简单检查：如果[的数量大于]的数量，可能是问题
            open_brackets = line.count('[')
            close_brackets = line.count(']')
            # 检查是否在字符串字面量中
            in_string = False
            escape = False
            for j, char in enumerate(line):
                if escape:
                    escape = False
                    continue
                if char == '\\':
                    escape = True
                    continue
                if char == '"' and not escape:
                    in_string = not in_string
                if not in_string and char == '[':
                    # 检查后面是否有对应的]
                    rest = line[j+1:]
                    # 简单检查：如果后面没有]，可能是问题
                    if ']' not in rest[:50]:  # 检查后面50个字符
                        # 在行末添加 ]
                        if line.rstrip()[-1] not in [']', ')', '}']:
                            lines[i] = line.rstrip() + ']' + '\n'
                            fixed_count += 1
                            print(f'修复了第 {i+1} 行的未闭合方括号')
                            break

# 2. 修复所有"non-declaration statement outside function body"错误
# 这些错误通常是因为前面的函数没有正确闭合
error_lines = [1790, 3220, 3660, 3876, 6398, 6677, 7361]
for err_line in sorted(error_lines, reverse=True):  # 从后往前修复
    if len(lines) > err_line - 1:
        # 检查之前是否有未闭合的函数
        brace_count = 0
        func_start = -1
        for i in range(max(0, err_line - 200), err_line):
            line = lines[i]
            if 'func ' in line and '{' in line:
                func_start = i
                brace_count = 0
            brace_count += line.count('{') - line.count('}')
        if brace_count > 0 and func_start >= 0:
            # 在错误行之前添加闭合大括号
            indent = '\t' * (brace_count - 1) if brace_count > 1 else '\t'
            lines.insert(err_line - 1, indent + '}\n')
            fixed_count += 1
            print(f'修复了第{err_line}行之前的未闭合函数（函数开始于第{func_start+1}行）')

# 3. 修复第6335行的else错误
if len(lines) > 6334:
    # 检查第6320-6335行之间是否有未闭合的if
    brace_count = 0
    found_if = False
    for i in range(6320, 6335):
        line = lines[i]
        if 'if ' in line and '{' in line:
            found_if = True
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if found_if and brace_count > 0:
        lines.insert(6334, '\t}\n')
        fixed_count += 1
        print(f'修复了第6335行之前的未闭合if语句')

# 4. 修复第7277行的break错误
if len(lines) > 7276:
    line_7277 = lines[7276]
    if 'break' in line_7277 and '//' not in line_7277.split('break')[0]:
        # 检查前面是否有for/switch/select
        found_loop = False
        brace_count = 0
        for i in range(max(0, 7250), 7277):
            line = lines[i]
            if 'for ' in line or 'switch ' in line or 'select ' in line:
                found_loop = True
                brace_count = 0
            brace_count += line.count('{') - line.count('}')
            if found_loop and brace_count == 0:
                # break不在循环中，注释掉
                lines[7276] = line_7277.replace('break', '// break  // 修复：break不在循环中')
                fixed_count += 1
                print(f'修复了第7277行的break语句（不在循环中）')
                break

# 写入修复后的文件
if fixed_count > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(lines)
    print(f'\n修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
