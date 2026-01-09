#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
清理battle.go中的多余空行
保留函数之间的分隔（最多2个空行），移除其他连续空行
"""

import os

# 获取脚本所在目录
script_dir = os.path.dirname(os.path.abspath(__file__))
file_path = os.path.join(script_dir, 'server', 'internal', 'test', 'runner', 'battle.go')

with open(file_path, 'rb') as f:
    content_bytes = f.read()

content = content_bytes.decode('utf-8', errors='replace')
lines = content.split('\n')

print(f'原始行数: {len(lines)}')

# 清理多余空行
cleaned_lines = []
empty_count = 0
prev_was_func = False
prev_was_brace = False

for i, line in enumerate(lines):
    is_empty = line.strip() == ''
    stripped = line.strip()
    
    if is_empty:
        empty_count += 1
        # 如果前一行是函数定义或类型定义，允许最多2个空行
        if prev_was_func and empty_count <= 2:
            cleaned_lines.append('')
        # 如果前一行是右大括号，允许1个空行
        elif prev_was_brace and empty_count == 1:
            cleaned_lines.append('')
        # 如果前一行不是空行，允许1个空行（用于分隔代码块）
        elif empty_count == 1 and len(cleaned_lines) > 0 and cleaned_lines[-1].strip() != '':
            cleaned_lines.append('')
        # 其他情况跳过多余空行
    else:
        empty_count = 0
        cleaned_lines.append(line)
        # 检查是否是函数定义或类型定义
        prev_was_func = (stripped.startswith('func ') or 
                        stripped.startswith('type ') or 
                        stripped.startswith('package ') or
                        stripped.startswith('import ('))
        # 检查是否是右大括号
        prev_was_brace = stripped == '}'

content = '\n'.join(cleaned_lines)

# 写入文件
with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(cleaned_lines)}')
print(f'移除了 {len(lines) - len(cleaned_lines)} 行空行')
