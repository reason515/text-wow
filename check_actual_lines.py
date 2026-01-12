#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_actual_lines():
    # 检查 context.go 第270行的实际内容
    filepath = 'server/internal/test/runner/context.go'
    print(f"检查文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 查找函数定义
    func_start = content.find(b'func (tr *TestRunner) updateAssertionContext()')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    # 计算行号
    func_line = content[:func_start].count(b'\n') + 1
    print(f"  函数从第 {func_line} 行开始")
    
    # 查找第270行（从函数开始计算）
    lines = content[func_start:].split(b'\n')
    target_line_idx = 270 - func_line
    
    if target_line_idx < len(lines):
        print(f"  第270行（函数内第 {target_line_idx} 行）:")
        print(f"    内容: {repr(lines[target_line_idx])}")
        print(f"    长度: {len(lines[target_line_idx])}")
        print(f"    字节: {lines[target_line_idx].hex()}")
        
        # 检查前后几行
        print(f"  前后上下文:")
        for i in range(max(0, target_line_idx-3), min(len(lines), target_line_idx+4)):
            marker = ">>>" if i == target_line_idx else "   "
            actual_line = func_line + i
            print(f"    {marker} Line {actual_line}: {repr(lines[i][:100])}")
    
    # 检查 equipment.go 第164行
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
    
    lines2 = content2[func_start2:].split(b'\n')
    target_line_idx2 = 164 - func_line2
    
    if target_line_idx2 < len(lines2):
        print(f"  第164行（函数内第 {target_line_idx2} 行）:")
        print(f"    内容: {repr(lines2[target_line_idx2])}")
        print(f"    长度: {len(lines2[target_line_idx2])}")
        print(f"    字节: {lines2[target_line_idx2].hex()}")
        
        # 检查前后几行
        print(f"  前后上下文:")
        for i in range(max(0, target_line_idx2-3), min(len(lines2), target_line_idx2+4)):
            marker = ">>>" if i == target_line_idx2 else "   "
            actual_line = func_line2 + i
            print(f"    {marker} Line {actual_line}: {repr(lines2[i][:100])}")

if __name__ == '__main__':
    check_actual_lines()
