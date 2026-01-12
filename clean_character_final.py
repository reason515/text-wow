#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
清理 character.go 文件中的空行并修复剩余编码问题
"""

def clean_and_fix():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    cleaned_lines = []
    prev_empty = False
    
    for i, line in enumerate(lines):
        # 检查是否是空行（只包含空白字符）
        is_empty = line.strip() == ''
        
        # 如果连续多个空行，只保留一个
        if is_empty:
            if not prev_empty:
                cleaned_lines.append('\n')
            prev_empty = True
        else:
            cleaned_lines.append(line)
            prev_empty = False
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(cleaned_lines)
    
    print(f"清理完成！原始行数: {len(lines)}, 清理后行数: {len(cleaned_lines)}")

if __name__ == '__main__':
    clean_and_fix()
