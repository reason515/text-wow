#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import re

file_path = 'server/internal/test/runner/calculation.go'
with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

fixed_lines = []
for i, line in enumerate(lines):
    # 修复所有注释和代码混在一起的情况
    if '\t// ' in line and '\t' in line[line.find('\t// ') + 4:] and not line.strip().endswith(')'):
        # 检查是否有代码跟在注释后面
        parts = line.split('\t// ')
        if len(parts) > 1:
            comment_part = parts[1].split('\t')
            if len(comment_part) > 1:
                # 分离注释和代码
                comment = comment_part[0].rstrip()
                code = '\t' + '\t'.join(comment_part[1:])
                fixed_lines.append(f'\t// {comment}\n')
                fixed_lines.append(code)
                continue
    
    # 修复所有 ")[0]) 为 "）")[0])
    line = line.replace('")[0])', '"）")[0])')
    
    fixed_lines.append(line)

with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
    f.writelines(fixed_lines)

print('修复完成')
