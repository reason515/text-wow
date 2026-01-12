#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复其他文件中的编码问题
"""

def fix_other_files2():
    # 修复 context.go
    context_path = 'server/internal/test/runner/context.go'
    with open(context_path, 'rb') as f:
        content = f.read()
    
    # 修复第59行：注释和代码混在一起
    old1 = b'\xe5\x8f\xaf\xe8\x83\xbd\xe5\x9c\xa8setup\xe4\xb8\xad\xe8\xae\xbe\xe7\xbd\xae\xe4\xba\x86\xef\xbf\xbd\t\t\tif goldVal'
    new1 = b'\xe5\x8f\xaf\xe8\x83\xbd\xe5\x9c\xa8setup\xe4\xb8\xad\xe8\xae\xbe\xe7\xbd\xae\xe4\xba\x86\xef\xbc\x89\n\t\t\tif goldVal'
    content = content.replace(old1, new1)
    
    with open(context_path, 'wb') as f:
        f.write(content)
    
    # 修复 equipment.go
    equipment_path = 'server/internal/test/runner/equipment.go'
    with open(equipment_path, 'rb') as f:
        content = f.read()
    
    # 修复第163行：检查是否有未闭合的块
    # 从字节输出看，第163行是 `level := 1`，但它在函数外
    # 让我检查前面的代码结构
    
    with open(equipment_path, 'wb') as f:
        f.write(content)
    
    # 修复 instruction.go
    instruction_path = 'server/internal/test/runner/instruction.go'
    with open(instruction_path, 'rb') as f:
        content = f.read()
    
    # 修复第105行：字符串缺少右引号
    old2 = b'\xe8\xae\xa1\xe7\xae\x97\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xef\xbf\xbd) {'
    new2 = b'\xe8\xae\xa1\xe7\xae\x97\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb") {'
    content = content.replace(old2, new2)
    
    # 修复第107行：字符串缺少右引号
    old3 = b'\xe8\xae\xa1\xe7\xae\x97\xe6\xb3\x95\xe6\x9c\xaf\xe6\x9a\xb4\xe5\x87\xbb\xef\xbf\xbd) {'
    new3 = b'\xe8\xae\xa1\xe7\xae\x97\xe6\xb3\x95\xe6\x9c\xaf\xe6\x9a\xb4\xe5\x87\xbb") {'
    content = content.replace(old3, new3)
    
    with open(instruction_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_other_files2()
