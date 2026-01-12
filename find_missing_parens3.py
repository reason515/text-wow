#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def find_missing_parens3():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"检查文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
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
    
    # 检查括号匹配，忽略注释中的括号
    paren_count = 0
    
    for i, line in enumerate(lines[:146]):
        # 检查是否是注释行
        comment_pos = line.find(b'//')
        if comment_pos != -1:
            # 这是注释行，只计算注释之前的括号
            code_part = line[:comment_pos]
            paren_count += code_part.count(b'(') - code_part.count(b')')
        else:
            # 不是注释行，计算所有括号
            paren_count += line.count(b'(') - line.count(b')')
        
        if paren_count < 0:
            print(f"  第 {i+96} 行后括号计数变为负: {paren_count}")
            print(f"    内容: {repr(line[:150])}")
            break

if __name__ == '__main__':
    find_missing_parens3()
