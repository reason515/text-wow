#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_strings():
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
        if b'"检查战斗初始状)' in line:
            # 替换为正确的字符串
            fixed_line = line.replace(
                b'"检查战斗初始状) || strings.Contains(instruction, "检查战斗状")',
                b'"检查战斗初始状态") || strings.Contains(instruction, "检查战斗状态")'
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
    fix_instruction_strings()
