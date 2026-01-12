#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_strings2():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    original_lines = lines[:]
    changed = False
    
    # 修复第300行（索引299）
    if len(lines) > 299:
        line = lines[299]
        # 修复字符串未关闭的问题
        # 查找 "检查战斗初始状) 模式
        pattern = b'\xe6\xa3\x80\xe6\x9f\xa5\xe6\x88\x98\xe6\x96\x97\xe5\x88\x9d\xe5\xa7\x8b\xe7\x8a\xb6)'
        if pattern in line:
            # 替换为正确的字符串
            fixed_line = line.replace(
                pattern + b' || strings.Contains(instruction, "' + b'\xe6\xa3\x80\xe6\x9f\xa5\xe6\x88\x98\xe6\x96\x97\xe7\x8a\xb6"',
                b'\xe6\xa3\x80\xe6\x9f\xa5\xe6\x88\x98\xe6\x96\x97\xe5\x88\x9d\xe5\xa7\x8b\xe7\x8a\xb6\xe6\x80\x81") || strings.Contains(instruction, "' + b'\xe6\xa3\x80\xe6\x9f\xa5\xe6\x88\x98\xe6\x96\x97\xe7\x8a\xb6\xe6\x80\x81"'
            )
            if fixed_line != line:
                lines[299] = fixed_line
                changed = True
                print(f"  修复了第 300 行的字符串问题")
    
    # 修复第303行（索引302）
    if len(lines) > 302:
        line = lines[302]
        # 修复 instruction" 应该是 instruction)
        if b'instruction")' in line:
            fixed_line = line.replace(b'instruction")', b'instruction)')
            if fixed_line != line:
                lines[302] = fixed_line
                changed = True
                print(f"  修复了第 303 行的字符串问题")
    
    if changed:
        with open(filepath, 'wb') as f:
            f.writelines(lines)
        print(f"  已保存更改")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_strings2()
