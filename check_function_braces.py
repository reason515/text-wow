#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_function_braces():
    for filepath, func_name, target_line in [
        ('server/internal/test/runner/context.go', 'updateAssertionContext', 270),
        ('server/internal/test/runner/equipment.go', 'generateMultipleEquipments', 164)
    ]:
        print(f"\n检查文件: {filepath}, 函数: {func_name}, 目标行: {target_line}")
        with open(filepath, 'rb') as f:
            content = f.read()
        
        # 查找函数定义
        func_pattern = f'func (tr *TestRunner) {func_name}'.encode('utf-8')
        func_start = content.find(func_pattern)
        
        if func_start == -1:
            print("  未找到函数定义")
            continue
        
        # 从函数开始计算大括号
        brace_count = 0
        lines = content[func_start:].split(b'\n')
        func_line_num = content[:func_start].count(b'\n') + 1
        
        print(f"  函数从第 {func_line_num} 行开始")
        
        for i, line in enumerate(lines):
            brace_count += line.count(b'{') - line.count(b'}')
            
            # 检查目标行
            current_line = func_line_num + i
            if current_line == target_line:
                print(f"  第 {target_line} 行: brace_count={brace_count}")
                print(f"    内容: {line[:150]}")
                
                if brace_count == 0:
                    print(f"    警告: 大括号计数为0，函数可能已结束")
                elif brace_count < 0:
                    print(f"    错误: 大括号计数为负，可能有未匹配的右括号")
                else:
                    print(f"    正常: 大括号计数为正，函数仍在进行中")
                
                # 显示前后几行
                print(f"    前后上下文:")
                for j in range(max(0, i-3), min(len(lines), i+4)):
                    marker = ">>>" if j == i else "   "
                    line_num = func_line_num + j
                    print(f"    {marker} Line {line_num}: {lines[j][:100]}")
                
                break
            
            # 如果大括号计数为0且已经过了函数开始，说明函数结束
            if brace_count == 0 and i > 0:
                line_num = func_line_num + i
                if line_num < target_line:
                    print(f"  警告: 函数在第 {line_num} 行结束（在目标行之前）")
                    print(f"    内容: {line[:150]}")

if __name__ == '__main__':
    check_function_braces()
