#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
进一步清理battle.go中的空行
- 函数之间只保留1个空行
- import块内不保留空行
- 代码块之间最多保留1个空行
"""

file_path = 'server/internal/test/runner/battle.go'

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

cleaned_lines = []
in_import_block = False
prev_was_func = False
prev_was_brace = False
prev_was_empty = False

for i, line in enumerate(lines):
    stripped = line.strip()
    is_empty = stripped == ''
    
    # 检测是否在import块中
    if stripped == 'import (':
        in_import_block = True
        cleaned_lines.append(line)
        prev_was_empty = False
        continue
    elif in_import_block and stripped == ')':
        in_import_block = False
        cleaned_lines.append(line)
        prev_was_empty = False
        continue
    
    # import块内：不保留空行
    if in_import_block:
        if not is_empty:
            cleaned_lines.append(line)
            prev_was_empty = False
        continue
    
    # 处理空行
    if is_empty:
        # 如果前一行是函数定义，保留1个空行
        if prev_was_func:
            if not prev_was_empty:  # 避免重复添加
                cleaned_lines.append('')
                prev_was_empty = True
        # 如果前一行是右大括号，保留1个空行
        elif prev_was_brace:
            if not prev_was_empty:
                cleaned_lines.append('')
                prev_was_empty = True
        # 如果前一行不是空行且不是函数/大括号，保留1个空行（用于分隔代码块）
        elif not prev_was_empty and len(cleaned_lines) > 0 and cleaned_lines[-1].strip() != '':
            cleaned_lines.append('')
            prev_was_empty = True
        # 其他情况跳过多余空行
    else:
        cleaned_lines.append(line)
        prev_was_empty = False
        
        # 检查是否是函数定义
        prev_was_func = (stripped.startswith('func ') or 
                        stripped.startswith('type ') or 
                        stripped.startswith('package '))
        
        # 检查是否是右大括号
        prev_was_brace = stripped == '}'

# 移除文件末尾的空行
while cleaned_lines and cleaned_lines[-1].strip() == '':
    cleaned_lines.pop()

content = '\n'.join(cleaned_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(cleaned_lines)}')
print(f'移除了 {len(lines) - len(cleaned_lines)} 行空行')
