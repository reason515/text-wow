#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_if():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 检查所有 if 语句后换行然后 { 的情况
    i = 0
    while i < len(lines) - 1:
        line = lines[i]
        next_line = lines[i + 1] if i + 1 < len(lines) else b''
        
        # 检查是否是 if 语句，且下一行是缩进的 {
        if line.strip().startswith(b'if') and b'strings.Contains' in line:
            # 检查下一行是否是缩进的 {
            if next_line.strip() == b'{':
                # 合并到同一行
                fixed_line = line.rstrip(b'\r\n') + b' {'
                lines[i] = fixed_line
                # 移除下一行
                lines.pop(i + 1)
                changed = True
                print(f"  修复了第 {i+1} 行的 if 语句")
                # 继续检查，不增加i
                continue
        i += 1
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_if()
