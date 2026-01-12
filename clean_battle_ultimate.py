#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 第一步：移除所有连续的空行，只保留单个空行
cleaned = []
prev_empty = False
for line in lines:
    is_empty = line.strip() == ''
    if is_empty:
        if not prev_empty:
            cleaned.append('')
            prev_empty = True
    else:
        cleaned.append(line)
        prev_empty = False

# 第二步：移除import块和函数内部的空行
final = []
in_import = False
in_func = False
brace_count = 0

for i, line in enumerate(cleaned):
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
        brace_count = 0
        # 如果前一行不是空行，添加一个空行分隔函数
        if final and final[-1].strip() != '':
            final.append('')
        final.append(line)
        continue
    
    # 计算大括号层级
    if '{' in stripped:
        brace_count += stripped.count('{')
    if '}' in stripped:
        brace_count -= stripped.count('}')
        if brace_count <= 0 and in_func:
            in_func = False
            final.append(line)
            continue
    
    # 处理空行
    if is_empty:
        # 只在函数之间保留空行（不在函数内部）
        if not in_func:
            # 检查下一行是否是函数定义
            if i + 1 < len(cleaned):
                next_line = cleaned[i + 1].strip()
                if next_line.startswith('func ') or next_line.startswith('type '):
                    if final and final[-1].strip() != '':
                        final.append('')
    else:
        final.append(line)

# 移除文件末尾的空行
while final and not final[-1].strip():
    final.pop()

content = '\n'.join(final)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(final)}')
print(f'移除了 {len(lines) - len(final)} 行空行')
