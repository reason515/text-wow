#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')
original = content

fixed_count = 0

# 修复第36261行的未闭合if语句 - 在第36293行之后添加闭合大括号
# 检查第36261行到第36300行之间是否有闭合大括号
if len(lines) > 36260:
    if 'if char, ok := tr.context.Characters["character"]; ok {' in lines[36260]:
        # 检查第36293行之后是否有闭合大括号
        found_close = False
        for i in range(36293, min(36350, len(lines))):
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

# 修复第10085行的else if语句 - 检查是否有闭合大括号
if len(lines) > 10084:
    if '} else if' in lines[10084] and '{' in lines[10084]:
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

# 修复第10245行的return nil - 检查是否在函数体内
if len(lines) > 10244:
    if 'return nil' in lines[10244]:
        # 检查之前是否有函数定义
        has_func = False
        for i in range(max(0, 10220), 10245):
            if 'func ' in lines[i]:
                has_func = True
                break
        if not has_func:
            # 检查是否有未闭合的大括号
            brace_count = 0
            for i in range(max(0, 10200), 10245):
                brace_count += lines[i].count('{') - lines[i].count('}')
            if brace_count > 0:
                # 在第10245行之前添加闭合大括号
                lines.insert(10245, '\t}')
                fixed_count += 1
                print(f'修复了第10245行的return nil，在第{10245+1}行之前添加了闭合大括号')

# 修复第35619行的语法错误 - 检查第35603行是否有未闭合的括号
if len(lines) > 35618:
    # 检查第35603行是否有未闭合的括号
    if len(lines) > 35602:
        line_35603 = lines[35602]
        # 检查是否有未闭合的TrimSuffix调用
        if 'TrimSuffix' in line_35603 and line_35603.count('(') > line_35603.count(')'):
            # 修复未闭合的括号
            lines[35602] = line_35603.rstrip() + ')'
            fixed_count += 1
            print(f'修复了第35603行的未闭合括号')

content = '\n'.join(lines)

if content != original:
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'修复完成！共修复了 {fixed_count} 处')
else:
    print('没有需要修复的内容')
