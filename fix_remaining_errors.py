#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'文件总行数: {len(lines)}')
fixed_count = 0
original_lines = lines[:]

# 检查第1579行之前的代码，查找未闭合的结构
print('\n检查第1579行之前的代码...')
for i in range(max(0, 1550), 1579):
    line = lines[i]
    # 检查是否有未闭合的方括号（简单检查）
    if '[' in line and ']' not in line:
        # 检查是否在字符串中
        if line.count('"') % 2 == 0:  # 不在字符串中
            # 检查是否是函数参数类型（如 []string）
            if '[]' not in line and 'func' not in line:
                # 可能是问题，但先不修复，因为可能是误报
                pass

# 检查括号匹配
print('检查括号匹配...')
brace_count = 0
paren_count = 0
bracket_count = 0

for i in range(1500, 1580):
    line = lines[i]
    brace_count += line.count('{') - line.count('}')
    paren_count += line.count('(') - line.count(')')
    bracket_count += line.count('[') - line.count(']')
    
    if bracket_count < 0:
        print(f'第 {i+1} 行: 方括号不匹配 - bracket={bracket_count}')
        print(f'  内容: {repr(line[:100])}')

print(f'\n第1580行时的括号状态: brace={brace_count}, paren={paren_count}, bracket={bracket_count}')

# 如果bracket_count < 0，说明有未闭合的]
if bracket_count < 0:
    # 查找第一个未闭合的]
    for i in range(1500, 1580):
        line = lines[i]
        if ']' in line and '[' not in line:
            # 检查是否在字符串中
            if line.count('"') % 2 == 0:  # 不在字符串中
                # 删除这个]
                lines[i] = line.replace(']', '', 1)
                fixed_count += 1
                print(f'修复了第 {i+1} 行的多余方括号')
                bracket_count += 1
                if bracket_count >= 0:
                    break

# 修复所有"non-declaration statement outside function body"错误
error_lines = [1790, 3221, 3661, 3878, 6400, 6679, 7363]
for err_line in sorted(error_lines, reverse=True):
    if len(lines) > err_line - 1:
        brace_count = 0
        func_start = -1
        for i in range(max(0, err_line - 300), err_line):
            line = lines[i]
            if 'func ' in line and '{' in line:
                func_start = i
                brace_count = 0
            brace_count += line.count('{') - line.count('}')
        if brace_count > 0 and func_start >= 0:
            # 计算缩进
            indent_level = 0
            for j in range(func_start, err_line):
                if lines[j].strip().startswith('}'):
                    indent_level -= 1
                elif '{' in lines[j]:
                    indent_level += 1
            indent = '\t' * max(0, indent_level)
            lines.insert(err_line - 1, indent + '}\n')
            fixed_count += 1
            print(f'修复了第{err_line}行之前的未闭合函数（函数开始于第{func_start+1}行）')

# 修复第6337行的else错误
if len(lines) > 6336:
    brace_count = 0
    found_if = False
    for i in range(6320, 6337):
        line = lines[i]
        if 'if ' in line and '{' in line:
            found_if = True
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if found_if and brace_count > 0:
        lines.insert(6336, '\t}\n')
        fixed_count += 1
        print(f'修复了第6337行之前的未闭合if语句')

# 修复第7279行的break错误
if len(lines) > 7278:
    line_7279 = lines[7278]
    if 'break' in line_7279 and '//' not in line_7279.split('break')[0]:
        lines[7278] = line_7279.replace('break', '// break  // 修复：break不在循环中')
        fixed_count += 1
        print(f'修复了第7279行的break语句（不在循环中）')

# 写入修复后的文件
if fixed_count > 0:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(lines)
    print(f'\n修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
