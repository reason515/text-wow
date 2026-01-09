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

def find_unclosed_structure(start_line, end_line, lines):
    """查找未闭合的结构"""
    brace_count = 0
    paren_count = 0
    bracket_count = 0
    func_start = -1
    if_start = -1
    
    for i in range(start_line, min(end_line, len(lines))):
        line = lines[i]
        stripped = line.strip()
        
        # 检查函数定义
        if 'func ' in line and '{' in line:
            func_start = i
            brace_count = 0
        
        # 检查if语句
        if stripped.startswith('if ') and '{' in line:
            if_start = i
        
        # 计数括号
        brace_count += line.count('{') - line.count('}')
        paren_count += line.count('(') - line.count(')')
        bracket_count += line.count('[') - line.count(']')
    
    return {
        'brace_count': brace_count,
        'paren_count': paren_count,
        'bracket_count': bracket_count,
        'func_start': func_start,
        'if_start': if_start
    }

# 1. 修复第8995行的错误 - 检查之前是否有未闭合的数组或切片
if len(lines) > 8994:
    result = find_unclosed_structure(max(0, 8970), 8995, lines)
    if result['bracket_count'] > 0:
        # 查找未闭合的 [
        for i in range(max(0, 8970), min(8995, len(lines))):
            line = lines[i]
            if '[' in line and line.count('[') > line.count(']'):
                # 在行末添加 ]
                if ']' not in line.rstrip()[-10:]:  # 检查最后10个字符
                    lines[i] = line.rstrip() + ']'
                    fixed_count += 1
                    print(f'修复了第{i+1}行的未闭合 [')
                    break

# 2. 修复第10099行的错误 - 检查第10085行的else if是否有闭合大括号
if len(lines) > 10098:
    if len(lines) > 10084:
        line_10085 = lines[10084]
        if '} else if' in line_10085 and '{' in line_10085:
            # 检查之后是否有闭合大括号
            result = find_unclosed_structure(10085, min(10150, len(lines)), lines)
            if result['brace_count'] > 0:
                # 在第10099行之后添加闭合大括号
                insert_pos = 10099
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第10085行的未闭合else if语句，在第{insert_pos+1}行添加了闭合大括号')

# 3. 修复第10243行的错误 - 检查是否有未闭合的函数
if len(lines) > 10242:
    result = find_unclosed_structure(max(0, 10200), 10243, lines)
    if result['brace_count'] > 0 and result['func_start'] >= 0:
        # 在第10243行之前添加闭合大括号
        lines.insert(10243, '\t}')
        fixed_count += 1
        print(f'修复了第{result["func_start"]+1}行的未闭合函数，在第{10243+1}行之前添加了闭合大括号')

# 4. 修复第35620行的错误 - 检查第35600-35620行之间是否有未闭合的括号
if len(lines) > 35619:
    result = find_unclosed_structure(max(0, 35600), 35620, lines)
    if result['paren_count'] > 0:
        # 查找未闭合的 (
        for i in range(max(0, 35600), min(35620, len(lines))):
            line = lines[i]
            if '(' in line and line.count('(') > line.count(')'):
                # 检查是否是函数调用
                if 'TrimSuffix' in line or 'TrimPrefix' in line or 'Replace' in line or 'Contains' in line:
                    if ')' not in line.rstrip()[-10:]:  # 检查最后10个字符
                        lines[i] = line.rstrip() + ')'
                        fixed_count += 1
                        print(f'修复了第{i+1}行的未闭合括号')
                        break

# 5. 修复第36276行的错误 - 检查第36261行的if是否有闭合大括号
if len(lines) > 36275:
    if len(lines) > 36260:
        line_36261 = lines[36260]
        if 'if char, ok := tr.context.Characters["character"]; ok {' in line_36261:
            # 检查之后是否有闭合大括号
            result = find_unclosed_structure(36261, min(36350, len(lines)), lines)
            if result['brace_count'] > 0:
                # 在第36293行之后添加闭合大括号
                insert_pos = 36300
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第36261行的未闭合if语句，在第{insert_pos+1}行添加了闭合大括号')

# 6. 修复其他"non-declaration statement outside function body"错误
error_lines = [18483, 21044, 22356, 36628, 38292]
for error_line in sorted(error_lines, reverse=True):  # 从后往前修复，避免行号变化
    if len(lines) > error_line - 1:
        # 检查之前是否有未闭合的函数
        result = find_unclosed_structure(max(0, error_line - 200), error_line, lines)
        if result['brace_count'] > 0 and result['func_start'] >= 0:
            # 在错误行之前添加闭合大括号
            lines.insert(error_line - 1, '\t}')
            fixed_count += 1
            print(f'修复了第{result["func_start"]+1}行的未闭合函数，在第{error_line}行之前添加了闭合大括号')

content = '\n'.join(lines)

if content != '\n'.join(original_lines):
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
