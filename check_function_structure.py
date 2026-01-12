#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_function_structure():
    filepath = 'server/internal/test/runner/instruction.go'
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    # 检查第96行到第508行之间的括号匹配
    func_start = 95  # 第96行，索引从0开始
    func_end = 507   # 第508行
    
    brace_count = 0
    paren_count = 0
    
    for i in range(func_start, func_end + 1):
        line = lines[i]
        # 忽略注释
        comment_pos = line.find(b'//')
        if comment_pos != -1:
            code_part = line[:comment_pos]
        else:
            code_part = line
        
        brace_count += code_part.count(b'{') - code_part.count(b'}')
        paren_count += code_part.count(b'(') - code_part.count(b')')
        
        if brace_count < 0 or paren_count < 0:
            print(f"第 {i+1} 行：括号不匹配")
            print(f"  大括号计数: {brace_count}, 小括号计数: {paren_count}")
            print(f"  内容: {repr(line[:100])}")
    
    print(f"函数结束时：大括号计数: {brace_count}, 小括号计数: {paren_count}")
    
    # 检查第512行
    print(f"\n第512行内容: {repr(lines[511][:100])}")

if __name__ == '__main__':
    check_function_structure()
