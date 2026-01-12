#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_unclosed_strings():
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
    
    # 检查是否有未关闭的字符串
    in_string = False
    string_char = None
    
    for i, line in enumerate(lines):
        j = 0
        while j < len(line):
            char = line[j:j+1]
            
            # 检查是否是字符串开始或结束
            if char == b'"' and (j == 0 or line[j-1:j] != b'\\'):
                if not in_string:
                    in_string = True
                    string_char = b'"'
                elif string_char == b'"':
                    in_string = False
                    string_char = None
            
            j += 1
        
        # 如果行结束但字符串未关闭，报告
        if in_string and i > len(lines) - 20:
            print(f"  第 {i+96} 行后字符串未关闭")
    
    if not in_string:
        print("  所有字符串都已正确关闭")

if __name__ == '__main__':
    check_unclosed_strings()
