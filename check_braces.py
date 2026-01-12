#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_braces():
    for filepath in ['server/internal/test/runner/context.go', 'server/internal/test/runner/equipment.go']:
        print(f"\n检查文件: {filepath}")
        with open(filepath, 'rb') as f:
            lines = f.readlines()
        
        brace_count = 0
        func_start = None
        
        for i, line in enumerate(lines):
            # 查找函数定义
            if b'func ' in line and b'(' in line:
                if b'{' in line:
                    brace_count = line.count(b'{') - line.count(b'}')
                    func_start = i + 1
                else:
                    func_start = i + 1
                    brace_count = 0
            else:
                brace_count += line.count(b'{') - line.count(b'}')
                if func_start and brace_count == 0:
                    print(f"  函数从第 {func_start} 行开始，在第 {i+1} 行结束")
                    func_start = None
                    brace_count = 0
            
            # 检查目标行
            target_line = 269 if 'context.go' in filepath else 164
            if i + 1 == target_line:
                print(f"  第 {target_line} 行: brace_count={brace_count}, func_start={func_start}")
                print(f"    内容: {line[:100]}")

if __name__ == '__main__':
    check_braces()
