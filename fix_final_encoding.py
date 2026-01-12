#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_final_encoding():
    # 修复 instruction.go 中的双引号问题
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复双引号问题: "\xe8\xa3\x85\xe5\xa4\x87"" 应该是 "\xe8\xa3\x85\xe5\xa4\x87")
    content = content.replace(b'\xe8\xa3\x85\xe5\xa4\x87""', b'\xe8\xa3\x85\xe5\xa4\x87")')
    if content != original_content:
        changes += 1
        print(f"  修复了双引号问题")
        original_content = content
    
    # 查找所有类似的模式: "text"" 应该是 "text")
    pattern = rb'("[\x20-\x7e\x80-\xff]+)""'
    def replace_double_quotes(m):
        text = m.group(1)
        return text + b'")'
    
    new_content = re.sub(pattern, replace_double_quotes, content)
    if new_content != content:
        changes += 1
        print(f"  修复了所有双引号问题")
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
        
        # 检查函数定义和闭合
        lines = content.split(b'\n')
        func_start = None
        brace_count = 0
        
        for i, line in enumerate(lines):
            # 查找函数定义
            if b'func ' in line:
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
            # 尝试找到问题所在
            if filepath.endswith('context.go'):
                # 检查第269行附近
                print(f"  检查第269行附近...")
                for i in range(265, min(275, len(lines))):
                    print(f"    行 {i+1}: {lines[i][:100]}")
        else:
            print(f"  函数结构正常")

if __name__ == '__main__':
    fix_final_encoding()
