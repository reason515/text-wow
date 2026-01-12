#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_all_parens():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 检查所有包含 else if 的行
    for i, line in enumerate(lines):
        if b'} else if' in line and b'{' in line:
            # 检查括号匹配
            open_count = line.count(b'(')
            close_count = line.count(b')')
            
            if open_count > close_count:
                # 缺少右括号，在 { 之前添加
                brace_pos = line.find(b'{')
                if brace_pos != -1:
                    # 在 { 之前添加缺失的右括号
                    fixed_line = line[:brace_pos] + b')' + line[brace_pos:]
                    lines[i] = fixed_line
                    changed = True
                    print(f"  修复了第 {i+1} 行的括号匹配（添加了 {open_count - close_count} 个右括号）")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_all_parens()
