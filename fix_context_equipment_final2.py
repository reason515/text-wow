#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_context_equipment_final2():
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
    
    # 查找第270行附近的内容（从函数开始计算）
    lines = content[func_start:].split(b'\n')
    target_line_idx = 270 - (content[:func_start].count(b'\n') + 1)
    
    if target_line_idx < len(lines):
        print(f"  第270行内容: {repr(lines[target_line_idx][:200])}")
        print(f"  第269行内容: {repr(lines[target_line_idx-1][:200])}")
        print(f"  第271行内容: {repr(lines[target_line_idx+1][:200])}")
    
    # 检查是否有UTF-8替换字符
    if b'\xef\xbf\xbd' in content[func_start:func_start+20000]:
        print("  发现UTF-8替换字符，尝试修复...")
        # 查找并修复常见的编码问题
        # 修复注释中的编码问题
        content = content.replace(b'\xef\xbf\xbd', b'')
        if content != original_content:
            print("  已移除UTF-8替换字符")
    
    # 检查是否有未关闭的字符串或注释
    # 查找可能的问题模式 - 移除UTF-8替换字符
    if b'\xef\xbf\xbd' in content:
        # 简单地移除替换字符可能会破坏代码，所以只移除注释中的
        # 查找注释中的替换字符
        pattern = rb'//[^\n]*\xef\xbf\xbd[^\n]*'
        def replace_comment(m):
            # 保留注释前缀，移除替换字符
            comment = m.group(0).replace(b'\xef\xbf\xbd', b'')
            return comment
        
        new_content = re.sub(pattern, replace_comment, content)
        if new_content != content:
            content = new_content
            print(f"  修复了注释中的编码问题")
    
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
        # 查找第164行附近的内容
        lines2 = content2[func_start2:].split(b'\n')
        target_line_idx2 = 164 - (content2[:func_start2].count(b'\n') + 1)
        
        if target_line_idx2 < len(lines2):
            print(f"  第164行内容: {repr(lines2[target_line_idx2][:200])}")
            print(f"  第163行内容: {repr(lines2[target_line_idx2-1][:200])}")
            print(f"  第165行内容: {repr(lines2[target_line_idx2+1][:200])}")
    
    # 检查是否有UTF-8替换字符
    if b'\xef\xbf\xbd' in content2[func_start2:func_start2+10000]:
        print("  发现UTF-8替换字符，尝试修复...")
        content2 = content2.replace(b'\xef\xbf\xbd', b'')
        if content2 != original_content2:
            print("  已移除UTF-8替换字符")
    
    if content2 != original_content2:
        with open(filepath2, 'wb') as f:
            f.write(content2)
        print(f"  已保存更改 ({len(content2)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_context_equipment_final2()
