#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_all_comment_strings():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 检查所有注释中的字符串
    for i, line in enumerate(lines):
        if b'//' in line and b'"' in line:
            # 检查注释中的字符串是否未关闭
            comment_start = line.find(b'//')
            comment_part = line[comment_start:]
            
            # 检查注释部分中的引号数量
            quote_count = comment_part.count(b'"')
            if quote_count % 2 != 0:
                # 字符串未关闭，在行尾添加引号
                fixed_line = line.rstrip(b'\r\n') + b'"' + b'\r\n'
                lines[i] = fixed_line
                changed = True
                print(f"  修复了第 {i+1} 行的注释字符串问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_all_comment_strings()
