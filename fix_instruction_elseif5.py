#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_instruction_elseif5():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: else if 语句后换行，然后下一行是 { 的情况
    # 模式: ...))\r\n\t\t {\r\r\n 应该是 ...)) {\n
    # 更精确的模式：匹配 else if 语句，后面有换行和缩进的 {
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
    
    # 尝试直接替换已知的模式
    # 从之前的检查看，第55行是 else if strings.Contains(instruction, "获得")...
    # 查找所有 else if 后跟换行和 { 的情况
    lines = content.split(b'\n')
    fixed_lines = []
    changed = False
    
    for i, line in enumerate(lines):
        # 检查是否是 else if 语句，且下一行是缩进的 {
        if line.strip().startswith(b'else if') and i + 1 < len(lines):
            next_line = lines[i + 1]
            if next_line.strip() == b'{':
                # 合并到同一行
                fixed_line = line.rstrip(b'\r') + b' {'
                fixed_lines.append(fixed_line)
                # 跳过下一行
                if i + 2 < len(lines):
                    fixed_lines.extend(lines[i + 2:])
                changed = True
                print(f"  修复了第 {i+1} 行的 else if 语句")
                break
        fixed_lines.append(line)
    
    if changed:
        content = b'\n'.join(fixed_lines)
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_elseif5()
