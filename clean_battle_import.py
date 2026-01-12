#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import re

file_path = os.path.join('server', 'internal', 'test', 'runner', 'battle.go')

with open(file_path, 'rb') as f:
    content = f.read().decode('utf-8', errors='replace')

lines = content.split('\n')
print(f'原始行数: {len(lines)}')

# 清理import块内的空行
result = []
in_import = False

for line in lines:
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
    
    # 其他行保持原样
    result.append(line)

content = '\n'.join(result)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'清理后行数: {len(result)}')
print(f'移除了 {len(lines) - len(result)} 行空行')
