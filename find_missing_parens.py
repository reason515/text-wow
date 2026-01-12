#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def find_missing_parens():
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
    
    # 检查括号匹配，找到缺少括号的位置
    paren_count = 0
    
    for i, line in enumerate(lines):
        old_count = paren_count
        paren_count += line.count(b'(') - line.count(b')')
        
        # 如果括号计数变为负数，说明之前有地方缺少右括号
        if paren_count < 0 and old_count >= 0:
            print(f"  第 {i+96} 行后括号计数变为负: {paren_count}")
            print(f"    内容: {repr(line[:200])}")
        
        if i > len(lines) - 20:
            if paren_count != 0:
                print(f"  第 {i+96} 行: paren_count={paren_count}")
                print(f"    内容: {repr(line[:200])}")

if __name__ == '__main__':
    find_missing_parens()
