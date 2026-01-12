#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_final_issues():
    # 修复 instruction.go 第55行的问题
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复: "\xe9\xa5\xb0\xe5\x93\x81")" 应该是 "\xe9\xa5\xb0\xe5\x93\x81")
    pattern1 = b'\xe9\xa5\xb0\xe5\x93\x81"\)"'
    replacement1 = b'\xe9\xa5\xb0\xe5\x93\x81")'
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        changes += 1
        print(f"  修复了第55行的字符串问题")
    
    # 修复所有类似的模式: "text")" 应该是 "text")
    pattern2 = rb'("[\x20-\x7e\x80-\xff]+"\))"'
    def replace2(m):
        return m.group(1)
    
    new_content = re.sub(pattern2, replace2, content)
    if new_content != content:
        changes += 1
        print(f"  修复了所有多余的引号")
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
        
        # 查找函数定义
        lines = content.split(b'\n')
        func_stack = []
        brace_count = 0
        
        for i, line in enumerate(lines):
            # 查找函数定义
            if b'func ' in line:
                # 检查是否有开括号在同一行
                if b'{' in line:
                    brace_count = line.count(b'{') - line.count(b'}')
                    func_stack.append((i+1, brace_count))
                else:
                    # 函数定义在下一行
                    func_stack.append((i+1, 0))
                    brace_count = 0
            else:
                # 更新大括号计数
                brace_count += line.count(b'{') - line.count(b'}')
                if func_stack:
                    func_start, func_brace = func_stack[-1]
                    if func_brace + brace_count == 0:
                        # 函数结束
                        func_stack.pop()
                        brace_count = 0
        
        if func_stack:
            print(f"  警告: 有未关闭的函数")
            for func_start, _ in func_stack:
                print(f"    函数从第 {func_start} 行开始")
        else:
            print(f"  函数结构正常")

if __name__ == '__main__':
    fix_final_issues()
