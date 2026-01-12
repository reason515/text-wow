#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_line80():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 查找第80行的问题：else if 语句后直接跟了 { 和注释
    # 模式: } else if (...) {\t\t// ...
    # 需要将 { 后面的内容移到下一行
    lines = content.split(b'\n')
    
    for i, line in enumerate(lines):
        # 检查是否是 else if 语句，且包含 { 和注释
        if b'} else if' in line and b'{\t\t//' in line:
            # 找到 { 的位置
            brace_pos = line.find(b'{\t\t//')
            if brace_pos != -1:
                # 分割行
                before_brace = line[:brace_pos+1]  # 包含 {
                after_brace = line[brace_pos+1:]  # 包含 \t\t// ...
                
                # 重新组织：第一行是 else if ... {，第二行是注释
                new_line1 = before_brace.rstrip()
                new_line2 = after_brace
                
                # 替换
                lines[i] = new_line1
                lines.insert(i + 1, new_line2)
                print(f"  修复了第 {i+1} 行的 else if 语句")
                break
    
    if lines != content.split(b'\n'):
        new_content = b'\n'.join(lines)
        with open(filepath, 'wb') as f:
            f.write(new_content)
        print(f"  已保存更改 ({len(new_content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_line80()
