#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_strings6():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复第342行（索引341）
    if len(lines) > 341:
        line = lines[341]
        # 修复字符串未关闭的问题
        pattern = b'\xe4\xb8\xaa\xe6\x80\xaa\xe7\x89\xa9)'
        if pattern in line:
            fixed_line = line.replace(
                pattern,
                b'\xe4\xb8\xaa\xe6\x80\xaa\xe7\x89\xa9")'
            )
            if fixed_line != line:
                lines[341] = fixed_line
                changed = True
                print(f"  修复了第 342 行的字符串问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_strings6()
