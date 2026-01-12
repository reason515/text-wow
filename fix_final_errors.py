#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_final_errors():
    # 修复 context.go 第270行的问题
    filepath = 'server/internal/test/runner/context.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 查找函数定义
    func_start = content.find(b'func (tr *TestRunner) updateAssertionContext()')
    if func_start == -1:
        print("  未找到函数定义")
        return
    
    # 查找第270行（if weaponID）
    target_pattern = b'if weaponID, ok := tr.context.Variables["weapon_id"]'
    target_pos = content.find(target_pattern, func_start)
    
    if target_pos != -1:
        print(f"  找到目标行，位置: {target_pos}")
        
        # 检查前后200字节的内容
        start_check = max(func_start, target_pos - 200)
        end_check = min(len(content), target_pos + 200)
        check_content = content[start_check:end_check]
        
        # 查找可能的编码问题
        # 1. UTF-8替换字符
        if b'\xef\xbf\xbd' in check_content:
            print("  发现UTF-8替换字符")
            # 移除替换字符（但要小心，不要破坏正常的UTF-8字符）
            # 只移除明显是替换字符的情况
            content = content.replace(b'\xef\xbf\xbd', b'')
        
        # 2. 检查是否有未关闭的字符串或注释
        # 查找可能的问题模式：注释中有特殊字符
        # 检查第270行之前的注释
        lines_before = content[func_start:target_pos].split(b'\n')
        for i, line in enumerate(lines_before[-10:]):  # 检查最后10行
            if b'//' in line:
                # 检查注释中是否有问题
                if b'\xef\xbf\xbd' in line:
                    print(f"  发现注释中的编码问题（第 {len(lines_before) - 10 + i} 行）")
                    # 修复注释
                    fixed_line = line.replace(b'\xef\xbf\xbd', b'')
                    # 替换原内容
                    old_line_pos = content.rfind(line, func_start, target_pos)
                    if old_line_pos != -1:
                        content = content[:old_line_pos] + fixed_line + content[old_line_pos + len(line):]
    
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
        # 查找第164行（level := 1）
        target_pattern2 = b'\tlevel := 1'
        target_pos2 = content2.find(target_pattern2, func_start2)
        
        if target_pos2 != -1:
            print(f"  找到目标行，位置: {target_pos2}")
            
            # 检查前后200字节的内容
            start_check2 = max(func_start2, target_pos2 - 200)
            end_check2 = min(len(content2), target_pos2 + 200)
            check_content2 = content2[start_check2:end_check2]
            
            # 查找可能的编码问题
            if b'\xef\xbf\xbd' in check_content2:
                print("  发现UTF-8替换字符")
                content2 = content2.replace(b'\xef\xbf\xbd', b'')
            
            # 检查第164行之前的注释
            lines_before2 = content2[func_start2:target_pos2].split(b'\n')
            for i, line in enumerate(lines_before2[-10:]):  # 检查最后10行
                if b'//' in line:
                    if b'\xef\xbf\xbd' in line:
                        print(f"  发现注释中的编码问题（第 {len(lines_before2) - 10 + i} 行）")
                        fixed_line = line.replace(b'\xef\xbf\xbd', b'')
                        old_line_pos = content2.rfind(line, func_start2, target_pos2)
                        if old_line_pos != -1:
                            content2 = content2[:old_line_pos] + fixed_line + content2[old_line_pos + len(line):]
    
    if content2 != original_content2:
        with open(filepath2, 'wb') as f:
            f.write(content2)
        print(f"  已保存更改 ({len(content2)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_final_errors()
