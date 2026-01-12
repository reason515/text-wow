#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整修复 calculation.go 文件：编码问题和空行清理
"""

def fix_complete():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    fixed_lines = []
    in_import = False
    prev_was_blank = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        is_blank = not stripped
        
        # 检测 import 块
        if stripped == 'import (':
            in_import = True
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        if in_import and stripped == ')':
            in_import = False
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # import 块内：移除所有空行
        if in_import:
            if is_blank:
                continue
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # 修复所有字符串分割中的编码问题
        if '"' in line and ')[0])' in line and 'strings.Split' in line:
            line = line.replace('")[0])', '"）")[0])')
        
        # 处理空行：函数体内最多保留一个空行
        if is_blank:
            if not prev_was_blank:
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or last_line.startswith('//'):
                        fixed_lines.append(line)
                        prev_was_blank = True
            continue
        
        # 非空行：如果前面有多个空行，减少到一个
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
    print(f"修复完成！原文件 {original_count} 行，现在 {new_count} 行，移除了 {removed} 行空行")

if __name__ == '__main__':
    fix_complete()
