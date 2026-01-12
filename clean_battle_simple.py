#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
简单直接清理battle.go中的空行
- 移除所有连续的空行，只保留单个空行
- import块内也清理
"""

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 第一步：移除所有连续的空行，只保留单个空行
cleaned_lines = []
prev_empty = False

for line in lines:
    is_empty = line.strip() == ''
    if is_empty:
        if not prev_empty:
            cleaned_lines.append('')
            prev_empty = True
    else:
        cleaned_lines.append(line)
        prev_empty = False

# 第二步：进一步清理 - 移除函数内部和import块内的空行
final_lines = []
in_import = False
in_function = False
brace_count = 0
prev_was_func = False

for i, line in enumerate(cleaned_lines):
    stripped = line.strip()
    is_empty = stripped == ''
    
    # 检测import块
    if stripped == 'import (':
        in_import = True
        final_lines.append(line)
        continue
    elif in_import and stripped == ')':
        in_import = False
        final_lines.append(line)
        continue
    
    # import块内：移除所有空行
    if in_import:
        if not is_empty:
            final_lines.append(line)
        continue
    
    # 检测函数
    if stripped.startswith('func '):
        in_function = True
        brace_count = 0
        # 如果前一行不是空行，添加一个空行分隔函数
        if final_lines and final_lines[-1].strip() != '':
            final_lines.append('')
        final_lines.append(line)
        prev_was_func = True
        continue
    
    # 计算大括号
    if '{' in stripped:
        brace_count += stripped.count('{')
    if '}' in stripped:
        brace_count -= stripped.count('}')
        if brace_count <= 0 and in_function:
            in_function = False
            final_lines.append(line)
            continue
    
    # 处理空行
    if is_empty:
        # 只在函数之间保留空行
        if not in_function:
            if prev_was_func and final_lines and final_lines[-1].strip() != '':
                final_lines.append('')
            prev_was_func = False
        # 函数内部移除所有空行
    else:
        final_lines.append(line)
        prev_was_func = False

# 移除文件末尾的空行
while final_lines and final_lines[-1].strip() == '':
    final_lines.pop()

content = '\n'.join(final_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(final_lines)}')
print(f'移除了 {len(lines) - len(final_lines)} 行空行')
