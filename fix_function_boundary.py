#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_function_boundary():
    # 修复 context.go 第270行的问题
    # 问题可能是第306行之后，函数应该继续，但Go编译器认为函数已经结束
    filepath = 'server/internal/test/runner/context.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 查找第306行附近的代码
    # 第306行是 `}`，关闭了第53行开始的if语句
    # 但函数应该继续到第408行
    # 检查第306行之后是否有编码问题
    
    # 查找函数定义
    func_start = content.find(b'func (tr *TestRunner) updateAssertionContext()')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    # 查找第306行（关闭if语句的地方）
    lines = content[func_start:].split(b'\n')
    line306_idx = 306 - (content[:func_start].count(b'\n') + 1)
    
    if line306_idx < len(lines):
        print(f"  第306行内容: {repr(lines[line306_idx][:200])}")
        print(f"  第307行内容: {repr(lines[line306_idx+1][:200])}")
        print(f"  第308行内容: {repr(lines[line306_idx+2][:200])}")
        
        # 检查第306行之后是否有编码问题
        # 第306行之后应该是函数继续，不应该有函数外的代码
        # 但如果有编码问题，Go编译器可能误判
        
        # 检查是否有隐藏的字符或编码问题
        for i in range(line306_idx, min(len(lines), line306_idx + 10)):
            line = lines[i]
            # 检查是否有非ASCII字符（除了正常的UTF-8字符）
            if b'\xef\xbf\xbd' in line:
                print(f"  第 {306 + i - line306_idx} 行发现UTF-8替换字符")
                # 移除替换字符
                lines[i] = line.replace(b'\xef\xbf\xbd', b'')
        
        # 重新组合内容
        new_content = content[:func_start] + b'\n'.join(lines)
        if new_content != content:
            content = new_content
            print("  已修复编码问题")
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")
    
    # 修复 equipment.go 第164行的问题
    filepath2 = 'server/internal/test/runner/equipment.go'
    print(f"\n处理文件: {filepath2}")
    
    with open(filepath2, 'rb') as f:
        content2 = f.read()
    
    original_content2 = content2
    
    # 查找函数定义
    func_start2 = content2.find(b'func (tr *TestRunner) generateMultipleEquipments')
    if func_start2 == -1:
        print("  未找到函数定义")
    else:
        # 查找第164行
        lines2 = content2[func_start2:].split(b'\n')
        line164_idx = 164 - (content2[:func_start2].count(b'\n') + 1)
        
        if line164_idx < len(lines2):
            print(f"  第164行内容: {repr(lines2[line164_idx][:200])}")
            print(f"  第163行内容: {repr(lines2[line164_idx-1][:200])}")
            
            # 检查是否有编码问题
            for i in range(max(0, line164_idx-5), min(len(lines2), line164_idx+5)):
                line = lines2[i]
                if b'\xef\xbf\xbd' in line:
                    print(f"  第 {164 + i - line164_idx} 行发现UTF-8替换字符")
                    lines2[i] = line.replace(b'\xef\xbf\xbd', b'')
            
            # 重新组合内容
            new_content2 = content2[:func_start2] + b'\n'.join(lines2)
            if new_content2 != content2:
                content2 = new_content2
                print("  已修复编码问题")
    
    if content2 != original_content2:
        with open(filepath2, 'wb') as f:
            f.write(content2)
        print(f"  已保存更改 ({len(content2)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_function_boundary()
