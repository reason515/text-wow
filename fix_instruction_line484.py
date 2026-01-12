#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_line484():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复第484行（索引483）
    if len(lines) > 483:
        line = lines[483]
        # 检查注释中的字符串是否未关闭
        if b'//' in line and b'"' in line:
            # 检查引号数量
            quote_count = line.count(b'"')
            if quote_count % 2 != 0:
                # 字符串未关闭，在行尾添加引号
                fixed_line = line.rstrip(b'\r\n') + b'"' + b'\r\n'
                lines[483] = fixed_line
                changed = True
                print(f"  修复了第 484 行的字符串问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_line484()
