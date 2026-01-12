#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_instruction_elseif4():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: else if 语句后换行，然后下一行是 { 的情况
    # 模式: ...))\r\n\t\t {\r\r\n 应该是 ...)) {\n
    # 注意：字符串可能很长，需要匹配到换行符
    pattern1 = rb'(else\s+if\s+strings\.Contains\([^)]+\)[^)]*\))\)\r\n\t\t\s*\{\r\r\n'
    def replace1(m):
        return m.group(1) + b') {\n'
    
    new_content = re.sub(pattern1, replace1, content)
    if new_content != content:
        print(f"  修复了else if语句后的换行和大括号（模式1）")
        content = new_content
    
    # 更通用的模式：匹配任何 else if 后换行然后 {
    pattern2 = rb'(else\s+if[^)]+\))\)\r\n\t\t\s*\{\r\r\n'
    def replace2(m):
        return m.group(1) + b') {\n'
    
    new_content = re.sub(pattern2, replace2, content)
    if new_content != content:
        print(f"  修复了else if语句后的换行和大括号（模式2）")
        content = new_content
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_elseif4()
