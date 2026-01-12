#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 calculation.go 文件中的编码问题和空行，注意不要产生多余空行
"""

def fix_calculation():
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
        
        # 修复编码问题：第84行
        if i == 83 and '怒气获得不需要角' in line:
            fixed_lines.append('	// 怒气获得不需要角色存在\n')
            continue
        
        # 修复编码问题：第89行 - 检查是否有编码问题
        if i == 88 and 'if !ok || char == nil {' in line:
            # 检查下一行是否有编码问题
            if i + 1 < len(lines) and 'return fmt.Errorf("character not found")' in lines[i + 1]:
                fixed_lines.append(line)
                continue
        
        # 修复编码问题：检查字符串中的编码问题
        if 'strings.Contains(instruction, "' in line or 'strings.Split(instruction, "' in line:
            # 检查是否有未闭合的引号或编码问题
            if line.count('"') % 2 != 0:
                # 可能有编码问题，尝试修复
                line = line.replace('", "', '", "')
                line = line.replace('")', '")')
        
        # 处理空行：函数体内最多保留一个空行
        if is_blank:
            # 如果前面不是空行，且不是函数定义前，保留一个空行
            if not prev_was_blank:
                # 检查是否需要空行（函数定义前或注释后）
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or last_line.startswith('//') or last_line.startswith('func '):
                        fixed_lines.append(line)
                        prev_was_blank = True
            continue
        
        # 非空行：如果前面有多个空行，减少到一个
        if prev_was_blank and fixed_lines:
            last_line = fixed_lines[-1].strip()
            if last_line and not last_line.endswith('{') and not last_line.endswith('('):
                # 保留空行
                pass
            else:
                # 移除多余的空行
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
    fix_calculation()
