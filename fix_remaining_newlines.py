#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_remaining_newlines():
    # 修复 instruction.go 中的换行问题
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复: else if 语句后换行，然后下一行是 { 的情况
    # 模式: ...))\n\t\t {\r\r\n 应该是 ...)) {\n
    pattern1 = rb'(else\s+if\s+[^)]+\))\)\r?\n\t\t\s*\{\r?\r?\n'
    def replace1(m):
        return m.group(1) + b') {\n'
    
    new_content = re.sub(pattern1, replace1, content)
    if new_content != content:
        changes += 1
        print(f"  修复了else if语句后的换行和大括号")
        content = new_content
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")
    
    # 检查 context.go 和 equipment.go 的函数结构
    for filepath in ['server/internal/test/runner/context.go', 'server/internal/test/runner/equipment.go']:
        print(f"\n检查文件: {filepath}")
        with open(filepath, 'rb') as f:
            content = f.read()
        
        # 查找函数定义和闭合
        lines = content.split(b'\n')
        func_start = None
        brace_count = 0
        
        for i, line in enumerate(lines):
            # 查找函数定义
            if b'func ' in line and b'(' in line:
                # 检查是否有开括号在同一行
                if b'{' in line:
                    brace_count = line.count(b'{') - line.count(b'}')
                    func_start = i + 1
                else:
                    # 函数定义在下一行
                    func_start = i + 1
                    brace_count = 0
            elif func_start is not None:
                # 在函数内部
                brace_count += line.count(b'{') - line.count(b'}')
                if brace_count == 0:
                    # 函数结束
                    func_start = None
        
        if brace_count != 0 and func_start is not None:
            print(f"  警告: 函数从第 {func_start} 行开始，大括号不匹配 (计数: {brace_count})")
            # 检查第269行（context.go）或第164行（equipment.go）
            target_line = 269 if 'context.go' in filepath else 164
            if target_line <= len(lines):
                print(f"  检查第 {target_line} 行:")
                print(f"    {lines[target_line-1][:100]}")
        else:
            print(f"  函数结构正常")

if __name__ == '__main__':
    fix_remaining_newlines()
