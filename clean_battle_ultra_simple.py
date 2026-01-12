#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
彻底清理 battle.go 文件中的空行
"""

def clean_battle_file():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        lines = f.readlines()
    
    cleaned_lines = []
    in_import = False
    prev_was_blank = False
    prev_was_func = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        is_blank = not stripped
        
        # 检测 import 块开始
        if stripped == 'import (':
            in_import = True
            cleaned_lines.append(line)
            prev_was_blank = False
            prev_was_func = False
            continue
        
        # 检测 import 块结束
        if in_import and stripped == ')':
            in_import = False
            cleaned_lines.append(line)
            prev_was_blank = False
            prev_was_func = False
            continue
        
        # import 块内：移除所有空行
        if in_import:
            if is_blank:
                continue
            cleaned_lines.append(line)
            prev_was_blank = False
            prev_was_func = False
            continue
        
        # 修复第1331行的编码问题（0-based: 1330）
        if i == 1330 and '// 更新上下' in line:
            # 分离注释和代码
            cleaned_lines.append('	// 更新上下文\n')
            cleaned_lines.append('	tr.context.Characters["character"] = char\n')
            prev_was_blank = False
            prev_was_func = False
            continue
        
        # 检测函数定义
        is_func = line.strip().startswith('func ') or line.strip().startswith('//')
        
        # 处理空行
        if is_blank:
            # 顶级声明之间：只保留一个空行
            if not prev_was_blank:
                # 检查是否需要空行（函数定义前或注释后）
                if cleaned_lines:
                    last_line = cleaned_lines[-1].strip()
                    # 如果上一行是函数结束或注释，可能需要空行
                    if last_line == '}' or last_line.startswith('//'):
                        cleaned_lines.append(line)
                        prev_was_blank = True
                        continue
            # 其他情况跳过空行
            continue
        
        # 非空行处理
        # 如果是函数定义，前面可能需要一个空行
        if is_func and cleaned_lines and not prev_was_blank:
            last_line = cleaned_lines[-1].strip()
            if last_line and last_line != '}' and not last_line.startswith('//'):
                # 在函数定义前添加一个空行
                cleaned_lines.append('\n')
        
        # 函数体内部：移除多余空行，但保留必要的空行用于逻辑分组
        # 简单的策略：连续空行只保留一个
        if prev_was_blank and not is_blank:
            # 检查是否真的需要空行
            if cleaned_lines:
                last_line = cleaned_lines[-1].strip()
                # 如果上一行是 } 或注释，可能需要空行
                if last_line == '}' or last_line.startswith('//'):
                    # 保留空行
                    pass
                else:
                    # 移除多余的空行
                    while cleaned_lines and not cleaned_lines[-1].strip():
                        cleaned_lines.pop()
        
        cleaned_lines.append(line)
        prev_was_blank = is_blank
        prev_was_func = is_func
    
    # 移除文件末尾的空行
    while cleaned_lines and not cleaned_lines[-1].strip():
        cleaned_lines.pop()
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.writelines(cleaned_lines)
    
    original_count = len(lines)
    new_count = len(cleaned_lines)
    removed = original_count - new_count
    print(f"清理完成！原文件 {original_count} 行，现在 {new_count} 行，移除了 {removed} 行空行")

if __name__ == '__main__':
    clean_battle_file()
