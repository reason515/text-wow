#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_all_elseif():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 查找所有 else if 语句后直接跟了 { 和注释的情况
    # 模式: } else if (...) {\t\t// ...
    # 需要将 { 后面的内容移到下一行
    lines = content.split(b'\n')
    changed = False
    
    i = 0
    while i < len(lines):
        line = lines[i]
        # 检查是否是 else if 语句，且包含 { 和注释在同一行
        if b'} else if' in line and b'{\t\t//' in line:
            # 找到 { 的位置
            brace_pos = line.find(b'{\t\t//')
            if brace_pos != -1:
                # 分割行
                before_brace = line[:brace_pos+1].rstrip()  # 包含 {
                after_brace = line[brace_pos+1:]  # 包含 \t\t// ...
                
                # 重新组织：第一行是 else if ... {，第二行是注释
                lines[i] = before_brace
                lines.insert(i + 1, after_brace)
                print(f"  修复了第 {i+1} 行的 else if 语句")
                changed = True
                # 继续检查，不break
        i += 1
    
    if changed:
        new_content = b'\n'.join(lines)
        with open(filepath, 'wb') as f:
            f.write(new_content)
        print(f"  已保存更改 ({len(new_content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_all_elseif()
