#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def find_unclosed_blocks():
    filepath = 'server/internal/test/runner/instruction.go'
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    # 检查第96行到第508行之间的块匹配
    func_start = 95  # 第96行，索引从0开始
    func_end = 507   # 第508行
    
    brace_count = 0
    brace_stack = []  # 记录每个大括号的位置
    
    for i in range(func_start, func_end + 1):
        line = lines[i]
        line_num = i + 1
        
        # 忽略注释
        comment_pos = line.find(b'//')
        if comment_pos != -1:
            code_part = line[:comment_pos]
        else:
            code_part = line
        
        # 检查大括号
        for j, char in enumerate(code_part):
            if char == ord(b'{'):
                brace_count += 1
                brace_stack.append((line_num, j, 'open'))
            elif char == ord(b'}'):
                brace_count -= 1
                if brace_stack:
                    brace_stack.pop()
                else:
                    print(f"第 {line_num} 行：多余的关闭大括号")
        
        if brace_count < 0:
            print(f"第 {line_num} 行：大括号计数变为负: {brace_count}")
            print(f"  内容: {repr(line[:150])}")
    
    print(f"\n函数结束时：大括号计数: {brace_count}")
    if brace_count > 0:
        print(f"  有 {brace_count} 个未关闭的大括号")
        print(f"  未关闭的大括号位置: {brace_stack[-brace_count:] if brace_stack else '未知'}")

if __name__ == '__main__':
    find_unclosed_blocks()
