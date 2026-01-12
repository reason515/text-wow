#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_strings5():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复第317行（索引316）
    if len(lines) > 316:
        line = lines[316]
        # 修复字符串未关闭的问题
        pattern = b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xa7\x92\xe8\x89\xb2\xe5\xb1\x9e)'
        if pattern in line:
            fixed_line = line.replace(
                pattern + b' || strings.Contains(instruction, "' + b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xa7\x92"',
                b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xa7\x92\xe8\x89\xb2\xe5\xb1\x9e\xe6\x80\xa7") || strings.Contains(instruction, "' + b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xa7\x92\xe8\x89\xb2"'
            )
            if fixed_line != line:
                lines[316] = fixed_line
                changed = True
                print(f"  修复了第 317 行的字符串问题")
    
    # 修复第318行（索引317）
    if len(lines) > 317:
        line = lines[317]
        # 修复注释和代码在同一行的问题
        if b'\t\t//' in line and b'\t\treturn' in line:
            # 分割为两行
            comment_pos = line.find(b'\t\t//')
            return_pos = line.find(b'\t\treturn', comment_pos)
            if return_pos != -1:
                comment_line = line[comment_pos:return_pos].rstrip()
                return_line = line[return_pos:]
                # 修复 instruction" 应该是 instruction)
                return_line = return_line.replace(b'instruction")', b'instruction)')
                lines[317] = comment_line + b'\n'
                lines.insert(318, return_line)
                changed = True
                print(f"  修复了第 318 行的格式问题")
    
    # 修复所有包含 (" 的行，应该是 (instruction)
    for i, line in enumerate(lines):
        if b'(")' in line and b'return tr.' in line:
            fixed_line = line.replace(b'(")', b'(instruction)')
            if fixed_line != line:
                lines[i] = fixed_line
                changed = True
                print(f"  修复了第 {i+1} 行的 (\" 问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_strings5()
