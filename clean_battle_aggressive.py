#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
激进清理battle.go中的空行
- import块内完全移除空行
- 函数内部完全移除空行（除了必要的分隔）
- 函数之间只保留1个空行
"""

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

cleaned_lines = []
in_import_block = False
in_function = False
brace_level = 0
prev_was_func_def = False

for i, line in enumerate(lines):
    stripped = line.strip()
    is_empty = stripped == ''
    
    # 检测是否在import块中
    if stripped == 'import (':
        in_import_block = True
        cleaned_lines.append(line)
        continue
    elif in_import_block and stripped == ')':
        in_import_block = False
        cleaned_lines.append(line)
        continue
    
    # import块内：完全移除空行
    if in_import_block:
        if not is_empty:
            cleaned_lines.append(line)
        continue
    
    # 检测函数开始和结束
    if stripped.startswith('func '):
        in_function = True
        brace_level = 0
        # 如果前一行不是空行，添加一个空行分隔函数
        if cleaned_lines and cleaned_lines[-1].strip() != '':
            cleaned_lines.append('')
        cleaned_lines.append(line)
        prev_was_func_def = True
        continue
    
    # 计算大括号层级
    if '{' in stripped:
        brace_level += stripped.count('{')
    if '}' in stripped:
        brace_level -= stripped.count('}')
        if brace_level == 0 and in_function:
            # 函数结束
            cleaned_lines.append(line)
            in_function = False
            continue
    
    # 处理空行
    if is_empty:
        # 只在函数之间保留空行
        if not in_function and prev_was_func_def:
            if cleaned_lines and cleaned_lines[-1].strip() != '':
                cleaned_lines.append('')
            prev_was_func_def = False
        # 函数内部完全移除空行
        # 其他情况也移除空行
    else:
        cleaned_lines.append(line)
        prev_was_func_def = False

# 移除文件末尾的空行
while cleaned_lines and cleaned_lines[-1].strip() == '':
    cleaned_lines.pop()

# 最后清理：移除连续的空行（最多保留1个）
final_lines = []
prev_empty = False
for line in cleaned_lines:
    is_empty = line.strip() == ''
    if is_empty:
        if not prev_empty:
            final_lines.append('')
            prev_empty = True
    else:
        final_lines.append(line)
        prev_empty = False

content = '\n'.join(final_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(final_lines)}')
print(f'移除了 {len(lines) - len(final_lines)} 行空行')
