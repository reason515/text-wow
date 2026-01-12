#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
最终完整修复 calculation.go 文件
"""

def fix_final():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    # 使用正则表达式替换所有 ")[0]) 为 "）")[0])
    import re
    content = re.sub(r'"\)\[0\]\)', '"）")[0])', content)
    
    # 移除函数体内的多余空行（连续两个空行只保留一个）
    lines = content.split('\n')
    fixed_lines = []
    in_import = False
    prev_was_blank = False
    
    for line in lines:
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
        
        # 处理空行
        if is_blank:
            if not prev_was_blank:
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or last_line.startswith('//'):
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
                    fixed_lines.append('')
        
        fixed_lines.append(line)
        prev_was_blank = False
    
    # 移除文件末尾的空行
    while fixed_lines and not fixed_lines[-1].strip():
        fixed_lines.pop()
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write('\n'.join(fixed_lines))
        if fixed_lines:
            f.write('\n')
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final()
