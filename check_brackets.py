#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_brackets():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"检查文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 查找 executeInstruction 函数
    func_start = content.find(b'func (tr *TestRunner) executeInstruction')
    if func_start == -1:
        print("  未找到函数")
        return
    
    func_end = content.find(b'func (tr *TestRunner) executeTeardown', func_start)
    if func_end == -1:
        print("  未找到函数结束")
        return
    
    func_content = content[func_start:func_end]
    lines = func_content.split(b'\n')
    
    # 检查括号匹配
    paren_count = 0
    bracket_count = 0
    brace_count = 0
    
    for i, line in enumerate(lines):
        paren_count += line.count(b'(') - line.count(b')')
        bracket_count += line.count(b'[') - line.count(b']')
        brace_count += line.count(b'{') - line.count(b'}')
        
        if i > len(lines) - 10:
            print(f"  第 {i+96} 行: paren={paren_count}, bracket={bracket_count}, brace={brace_count}")
    
    if paren_count == 0 and bracket_count == 0 and brace_count == 0:
        print("  所有括号都匹配")
    else:
        print(f"  括号不匹配: paren={paren_count}, bracket={bracket_count}, brace={brace_count}")

if __name__ == '__main__':
    check_brackets()
