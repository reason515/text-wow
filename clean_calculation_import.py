#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
清理 calculation.go 文件 import 块中的空行
"""

def clean_import():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    in_import = False
    
    for line in lines:
        stripped = line.strip()
        
        if stripped == 'import (':
            in_import = True
            fixed_lines.append('import (\n')
            continue
        
        if in_import and stripped == ')':
            in_import = False
            fixed_lines.append(')\n')
            continue
        
        if in_import:
            if not stripped:
                continue
            fixed_lines.append(line)
            continue
        
        fixed_lines.append(line)
    
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    print("清理完成！")

if __name__ == '__main__':
    clean_import()
