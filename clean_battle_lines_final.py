#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
彻底清理battle.go中的空行
- import块内完全移除空行
- 函数内部代码块之间最多保留1个空行
- 函数之间只保留1个空行
- 移除所有其他多余空行
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
prev_was_comment = False

for i, line in enumerate(lines):
    stripped = line.strip()
    is_empty = stripped == ''
    is_comment = stripped.startswith('//') and not stripped.startswith('// ')
    
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
    
    # import块内：完全移除空行
    if in_import_block:
        if not is_empty:
            cleaned_lines.append(line)
            prev_was_empty = False
        continue
    
    # 处理空行
    if is_empty:
        # 如果前一行是函数定义，保留1个空行
        if prev_was_func:
            if not prev_was_empty:
                cleaned_lines.append('')
                prev_was_empty = True
        # 如果前一行是右大括号且是函数结束，保留1个空行
        elif prev_was_brace and i + 1 < len(lines):
            next_line = lines[i + 1].strip() if i + 1 < len(lines) else ''
            # 如果下一行是函数定义，保留1个空行
            if next_line.startswith('func ') or next_line.startswith('type '):
                if not prev_was_empty:
                    cleaned_lines.append('')
                    prev_was_empty = True
            # 否则跳过多余空行
        # 如果前一行是注释，且下一行是代码，保留1个空行
        elif prev_was_comment and i + 1 < len(lines):
            next_line = lines[i + 1].strip() if i + 1 < len(lines) else ''
            if next_line and not next_line.startswith('//') and not next_line.startswith('}'):
                if not prev_was_empty:
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
        
        # 检查是否是注释
        prev_was_comment = is_comment

# 移除文件末尾的空行
while cleaned_lines and cleaned_lines[-1].strip() == '':
    cleaned_lines.pop()

# 再次清理：移除连续的空行（最多保留1个）
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
