#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
file_path = os.path.join('server', 'internal', 'test', 'runner', 'battle.go')

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 清理策略：
# 1. import块内：完全移除空行
# 2. 函数内部：完全移除空行
# 3. 函数之间：只保留1个空行

final = []
in_import = False
in_func = False
brace_level = 0
prev_was_func = False

for i, line in enumerate(lines):
    stripped = line.strip()
    is_empty = stripped == ''
    
    # 检测import块
    if stripped == 'import (':
        in_import = True
        final.append(line)
        continue
    elif in_import and stripped == ')':
        in_import = False
        final.append(line)
        continue
    
    # import块内：移除所有空行
    if in_import:
        if not is_empty:
            final.append(line)
        continue
    
    # 检测函数开始
    if stripped.startswith('func '):
        in_func = True
        brace_level = 0
        # 如果前一行不是空行，添加一个空行分隔函数
        if final and final[-1].strip() != '':
            final.append('')
        final.append(line)
        prev_was_func = True
        continue
    
    # 计算大括号层级
    if '{' in stripped:
        brace_level += stripped.count('{')
    if '}' in stripped:
        brace_level -= stripped.count('}')
        if brace_level <= 0 and in_func:
            in_func = False
            final.append(line)
            continue
    
    # 处理空行
    if is_empty:
        # 只在函数之间保留空行（不在函数内部）
        if not in_func:
            # 检查下一行是否是函数定义
            if i + 1 < len(lines):
                next_stripped = lines[i + 1].strip()
                if next_stripped.startswith('func ') or next_stripped.startswith('type '):
                    if final and final[-1].strip() != '':
                        final.append('')
    else:
        final.append(line)

# 移除文件末尾的空行
while final and not final[-1].strip():
    final.pop()

# 最后清理：确保没有连续的空行
result = []
prev_empty = False
for line in final:
    is_empty = line.strip() == ''
    if is_empty:
        if not prev_empty:
            result.append('')
            prev_empty = True
    else:
        result.append(line)
        prev_empty = False

content = '\n'.join(result)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(result)}')
print(f'移除了 {len(lines) - len(result)} 行空行')
