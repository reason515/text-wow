#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_strings3():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复所有包含 instruction" 的行
    for i, line in enumerate(lines):
        if b'instruction")' in line:
            fixed_line = line.replace(b'instruction")', b'instruction)')
            if fixed_line != line:
                lines[i] = fixed_line
                changed = True
                print(f"  修复了第 {i+1} 行的 instruction\" 问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_strings3()
