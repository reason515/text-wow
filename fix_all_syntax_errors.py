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

# 修复策略：检查每个错误行之前的代码，查找未闭合的结构

# 1. 修复第8995行的错误 - 检查之前是否有未闭合的数组或切片
if len(lines) > 8994:
    # 检查第8970-8995行之间是否有未闭合的 [
    for i in range(max(0, 8970), min(8995, len(lines))):
        line = lines[i]
        if '[' in line and line.count('[') > line.count(']'):
            # 找到未闭合的 [
            print(f'发现第{i+1}行有未闭合的 [')
            # 尝试修复
            if ']' not in line:
                lines[i] = line.rstrip() + ']'
                fixed_count += 1
                print(f'修复了第{i+1}行的未闭合 [')

# 2. 修复第10099行的错误 - 检查第10085行的else if是否有闭合大括号
if len(lines) > 10098:
    # 检查第10085行的else if语句
    if len(lines) > 10084:
        line_10085 = lines[10084]
        if '} else if' in line_10085 and '{' in line_10085:
            # 检查之后是否有闭合大括号
            found_close = False
            for i in range(10085, min(10150, len(lines))):
                if lines[i].strip() == '}':
                    found_close = True
                    break
            if not found_close:
                # 在第10099行之后添加闭合大括号
                insert_pos = 10099
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第10085行的未闭合else if语句，在第{insert_pos+1}行添加了闭合大括号')

# 3. 修复第10243行的错误 - 检查是否有未闭合的函数
if len(lines) > 10242:
    # 检查第10200-10243行之间是否有未闭合的函数
    brace_count = 0
    func_start = -1
    for i in range(max(0, 10200), min(10243, len(lines))):
        line = lines[i]
        if 'func ' in line:
            func_start = i
            brace_count = 0
        brace_count += line.count('{') - line.count('}')
    if brace_count > 0 and func_start >= 0:
        # 在第10243行之前添加闭合大括号
        lines.insert(10243, '\t}')
        fixed_count += 1
        print(f'修复了第{func_start+1}行的未闭合函数，在第{10243+1}行之前添加了闭合大括号')

# 4. 修复第35619行的错误 - 检查第35603行是否有未闭合的括号
if len(lines) > 35618:
    # 检查第35600-35619行之间是否有未闭合的括号
    for i in range(max(0, 35600), min(35619, len(lines))):
        line = lines[i]
        if '(' in line and line.count('(') > line.count(')'):
            # 找到未闭合的 (
            if 'TrimSuffix' in line or 'TrimPrefix' in line or 'Replace' in line:
                lines[i] = line.rstrip() + ')'
                fixed_count += 1
                print(f'修复了第{i+1}行的未闭合括号')

# 5. 修复第36275行的错误 - 检查第36261行的if是否有闭合大括号
if len(lines) > 36274:
    # 检查第36261行的if语句
    if len(lines) > 36260:
        line_36261 = lines[36260]
        if 'if char, ok := tr.context.Characters["character"]; ok {' in line_36261:
            # 检查之后是否有闭合大括号
            found_close = False
            for i in range(36261, min(36350, len(lines))):
                if lines[i].strip() == '}':
                    found_close = True
                    break
            if not found_close:
                # 在第36293行之后添加闭合大括号
                insert_pos = 36300
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第36261行的未闭合if语句，在第{insert_pos+1}行添加了闭合大括号')

# 6. 修复其他"non-declaration statement outside function body"错误
# 这些错误通常是因为前面的函数没有正确闭合
error_lines = [18483, 21043, 22355, 36627, 38291]
for error_line in error_lines:
    if len(lines) > error_line - 1:
        # 检查之前是否有未闭合的函数
        brace_count = 0
        func_start = -1
        for i in range(max(0, error_line - 100), min(error_line, len(lines))):
            line = lines[i]
            if 'func ' in line:
                func_start = i
                brace_count = 0
            brace_count += line.count('{') - line.count('}')
        if brace_count > 0 and func_start >= 0:
            # 在错误行之前添加闭合大括号
            lines.insert(error_line - 1, '\t}')
            fixed_count += 1
            print(f'修复了第{func_start+1}行的未闭合函数，在第{error_line}行之前添加了闭合大括号')
            # 更新后续错误行号
            for j in range(len(error_lines)):
                if error_lines[j] > error_line:
                    error_lines[j] += 1

content = '\n'.join(lines)

if content != '\n'.join(original_lines):
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
