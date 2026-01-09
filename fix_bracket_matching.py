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

def count_braces(lines, start, end):
    """计算大括号数量"""
    brace_count = 0
    for i in range(start, min(end, len(lines))):
        line = lines[i]
        brace_count += line.count('{') - line.count('}')
    return brace_count

def find_last_non_empty_line(lines, start, end):
    """查找最后一个非空行"""
    for i in range(end - 1, start - 1, -1):
        if lines[i].strip():
            return i
    return -1

# 1. 修复第8995行的错误 - 检查第8977行的结构体字段
if len(lines) > 8994:
    # 检查第8970-8995行之间是否有未闭合的结构体
    brace_count = count_braces(lines, 8970, 8995)
    if brace_count > 0:
        # 查找最后一个非空行
        last_line = find_last_non_empty_line(lines, 8970, 8995)
        if last_line >= 0 and 'TestName:' in lines[last_line]:
            # 在第8977行之后添加闭合大括号
            insert_pos = 8978
            if insert_pos < len(lines):
                lines.insert(insert_pos, '\t\t}')
                fixed_count += 1
                print(f'修复了第{last_line+1}行的未闭合结构体，在第{insert_pos+1}行添加了闭合大括号')

# 2. 修复第10099行的错误 - 检查第10081行之后是否有else if未闭合
if len(lines) > 10098:
    # 检查第10081行之后是否有未闭合的else if
    if len(lines) > 10080:
        if lines[10080].strip() == '}':
            # 检查之后是否有else if未闭合
            brace_count = count_braces(lines, 10081, 10099)
            if brace_count > 0:
                # 在第10099行之后添加闭合大括号
                insert_pos = 10099
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第10081行之后的未闭合else if，在第{insert_pos+1}行添加了闭合大括号')

# 3. 修复第10243行的错误 - 检查第10225行的函数调用
if len(lines) > 10242:
    # 检查第10200-10243行之间是否有未闭合的函数
    brace_count = count_braces(lines, 10200, 10243)
    if brace_count > 0:
        # 查找最后一个非空行
        last_line = find_last_non_empty_line(lines, 10200, 10243)
        if last_line >= 0:
            # 在第10225行之后添加闭合大括号
            insert_pos = 10226
            if insert_pos < len(lines):
                lines.insert(insert_pos, '\t}')
                fixed_count += 1
                print(f'修复了第{last_line+1}行的未闭合函数，在第{insert_pos+1}行添加了闭合大括号')

# 4. 修复第35620行的错误 - 检查第35575行的if语句
if len(lines) > 35619:
    # 检查第35575行的if语句是否有闭合大括号
    if len(lines) > 35574:
        line_35575 = lines[35574]
        if 'if strings.Contains' in line_35575 and '{' in line_35575:
            # 检查之后是否有闭合大括号
            brace_count = count_braces(lines, 35575, 35620)
            if brace_count > 0:
                # 在第35623行之后添加闭合大括号
                insert_pos = 35624
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t\t\t}')
                    fixed_count += 1
                    print(f'修复了第35575行的未闭合if语句，在第{insert_pos+1}行添加了闭合大括号')

# 5. 修复第36276行的错误 - 检查第36247行的for循环
if len(lines) > 36275:
    # 检查第36247行的for循环是否有闭合大括号
    if len(lines) > 36246:
        line_36247 = lines[36246]
        if 'for key, monster := range tr.context.Monsters {' in line_36247:
            # 检查之后是否有闭合大括号
            brace_count = count_braces(lines, 36247, 36276)
            if brace_count > 0:
                # 在第36276行之后添加闭合大括号
                insert_pos = 36276
                if insert_pos < len(lines):
                    lines.insert(insert_pos, '\t}')
                    fixed_count += 1
                    print(f'修复了第36247行的未闭合for循环，在第{insert_pos+1}行添加了闭合大括号')

# 6. 修复其他"non-declaration statement outside function body"错误
error_lines = [18483, 21044, 22356, 36628, 38292]
for error_line in sorted(error_lines, reverse=True):
    if len(lines) > error_line - 1:
        # 检查之前是否有未闭合的函数
        brace_count = count_braces(lines, max(0, error_line - 200), error_line)
        if brace_count > 0:
            # 在错误行之前添加闭合大括号
            lines.insert(error_line - 1, '\t}')
            fixed_count += 1
            print(f'修复了第{error_line}行之前的未闭合函数，在第{error_line}行之前添加了闭合大括号')

content = '\n'.join(lines)

if content != '\n'.join(original_lines):
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
