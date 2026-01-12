#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import re

file_path = os.path.join('server', 'internal', 'test', 'runner', 'battle.go')

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 清理策略：完全移除import块和函数内部的所有空行，只保留函数之间的1个空行
result = []
in_import = False
in_func = False
brace_level = 0

for i, line in enumerate(lines):
    stripped = line.strip()
    is_empty = stripped == ''
    
    # 检测import块
    if stripped == 'import (':
        in_import = True
        result.append(line)
        continue
    elif in_import and stripped == ')':
        in_import = False
        result.append(line)
        continue
    
    # import块内：完全移除所有空行
    if in_import:
        if not is_empty:
            result.append(line)
        continue
    
    # 检测函数开始
    if stripped.startswith('func '):
        in_func = True
        brace_level = 0
        # 如果前一行不是空行，添加一个空行分隔函数
        if result and result[-1].strip() != '':
            result.append('')
        result.append(line)
        continue
    
    # 计算大括号层级
    if '{' in stripped:
        brace_level += stripped.count('{')
    if '}' in stripped:
        brace_level -= stripped.count('}')
        if brace_level <= 0 and in_func:
            in_func = False
            result.append(line)
            continue
    
    # 处理空行
    if is_empty:
        # 完全移除函数内部的空行
        if in_func:
            continue
        # 只在函数之间保留空行（不在函数内部）
        # 检查下一行是否是函数定义
        if i + 1 < len(lines):
            next_stripped = lines[i + 1].strip()
            if next_stripped.startswith('func '):
                if result and result[-1].strip() != '':
                    result.append('')
        # 其他空行都移除
    else:
        result.append(line)

# 移除文件末尾的空行
while result and not result[-1].strip():
    result.pop()

# 最后清理：确保没有连续的空行（最多保留1个）
final = []
prev_empty = False
for line in result:
    is_empty = line.strip() == ''
    if is_empty:
        if not prev_empty:
            final.append('')
            prev_empty = True
    else:
        final.append(line)
        prev_empty = False

content = '\n'.join(final)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(final)}')
print(f'移除了 {len(lines) - len(final)} 行空行')
