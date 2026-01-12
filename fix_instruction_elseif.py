#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_instruction_elseif():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: else if 语句后换行，然后下一行是 { 的情况
    # 模式: ...))\n\t\t {\r\r\n 应该是 ...)) {\n
    pattern1 = rb'(else\s+if\s+[^)]+\))\)\r?\n\t\t\s*\{\r?\r?\n'
    def replace1(m):
        return m.group(1) + b') {\n'
    
    new_content = re.sub(pattern1, replace1, content)
    if new_content != content:
        print(f"  修复了else if语句后的换行和大括号")
        content = new_content
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_elseif()
