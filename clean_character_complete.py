#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整清理 character.go 文件中的空行并修复编码问题
"""

def clean_complete():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    in_import = False
    prev_was_blank = False
    
    for line in lines:
        stripped = line.strip()
        is_blank = not stripped
        
        # 检测 import 块
        if stripped == 'import (':
            in_import = True
            fixed_lines.append('import (\n')
            prev_was_blank = False
            continue
        
        if in_import and stripped == ')':
            in_import = False
            fixed_lines.append(')\n')
            prev_was_blank = False
            continue
        
        # import 块内：移除所有空行
        if in_import:
            if is_blank:
                continue
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # 处理空行：函数体内最多保留一个空行
        if is_blank:
            if not prev_was_blank:
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or (last_line.startswith('//') and not last_line.endswith(')')):
                        fixed_lines.append(line)
                        prev_was_blank = True
            continue
        
        # 非空行
        if prev_was_blank and fixed_lines:
            last_line = fixed_lines[-1].strip()
            if last_line and not last_line.endswith('{') and not last_line.endswith('('):
                pass
            else:
                while fixed_lines and not fixed_lines[-1].strip():
                    fixed_lines.pop()
                if fixed_lines and fixed_lines[-1].strip():
                    fixed_lines.append('\n')
        
        fixed_lines.append(line)
        prev_was_blank = False
    
    # 移除文件末尾的空行
    while fixed_lines and not fixed_lines[-1].strip():
        fixed_lines.pop()
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(fixed_lines)
    
    original_count = len(lines)
    new_count = len(fixed_lines)
    removed = original_count - new_count
    print(f"清理完成！原文件 {original_count} 行，现在 {new_count} 行，移除了 {removed} 行空行")

if __name__ == '__main__':
    clean_complete()
