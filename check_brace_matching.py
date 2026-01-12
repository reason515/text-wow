#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_brace_matching():
    # 检查 context.go 的大括号匹配
    filepath = 'server/internal/test/runner/context.go'
    print(f"检查文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 查找函数定义
    func_start = content.find(b'func (tr *TestRunner) updateAssertionContext()')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    func_line = content[:func_start].count(b'\n') + 1
    print(f"  函数从第 {func_line} 行开始")
    
    # 从函数开始计算大括号
    brace_count = 0
    lines = content[func_start:].split(b'\n')
    
    for i, line in enumerate(lines):
        brace_count += line.count(b'{') - line.count(b'}')
        current_line = func_line + i
        
        # 检查第269行和第270行
        if current_line == 269:
            print(f"  第269行: brace_count={brace_count}")
            print(f"    内容: {repr(line)}")
        elif current_line == 270:
            print(f"  第270行: brace_count={brace_count}")
            print(f"    内容: {repr(line)}")
            if brace_count == 0:
                print(f"    警告: 大括号计数为0，函数可能已结束")
            elif brace_count < 0:
                print(f"    错误: 大括号计数为负")
            else:
                print(f"    正常: 大括号计数为正")
        
        # 如果大括号计数为0，说明函数结束
        if brace_count == 0 and i > 0:
            print(f"  函数在第 {current_line} 行结束")
            if current_line < 270:
                print(f"    ⚠️ 警告: 函数在目标行之前结束！")
            break
    
    # 检查 equipment.go 的大括号匹配
    filepath2 = 'server/internal/test/runner/equipment.go'
    print(f"\n检查文件: {filepath2}")
    
    with open(filepath2, 'rb') as f:
        content2 = f.read()
    
    func_start2 = content2.find(b'func (tr *TestRunner) generateMultipleEquipments')
    if func_start2 == -1:
        print("  未找到函数定义")
        return
    
    func_line2 = content2[:func_start2].count(b'\n') + 1
    print(f"  函数从第 {func_line2} 行开始")
    
    brace_count2 = 0
    lines2 = content2[func_start2:].split(b'\n')
    
    for i, line in enumerate(lines2):
        brace_count2 += line.count(b'{') - line.count(b'}')
        current_line = func_line2 + i
        
        # 检查第161行和第164行
        if current_line == 161:
            print(f"  第161行: brace_count={brace_count2}")
            print(f"    内容: {repr(line)}")
        elif current_line == 164:
            print(f"  第164行: brace_count={brace_count2}")
            print(f"    内容: {repr(line)}")
            if brace_count2 == 0:
                print(f"    警告: 大括号计数为0，函数可能已结束")
            elif brace_count2 < 0:
                print(f"    错误: 大括号计数为负")
            else:
                print(f"    正常: 大括号计数为正")
        
        # 如果大括号计数为0，说明函数结束
        if brace_count2 == 0 and i > 0:
            print(f"  函数在第 {current_line} 行结束")
            if current_line < 164:
                print(f"    ⚠️ 警告: 函数在目标行之前结束！")
            break

if __name__ == '__main__':
    check_brace_matching()
