#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复其他文件中的编码问题
"""

def fix_other_files5():
    # 修复 context.go
    context_path = 'server/internal/test/runner/context.go'
    with open(context_path, 'rb') as f:
        content = f.read()
    
    # 修复第316行：注释和代码混在一起
    old1 = b'\xe8\xae\xbe\xe7\xbd\xae\xe8\xa7\x92\xe8\x89\xb2\xe7\x9a\x84\xe5\x9f\xba\xe6\x9c\xac\xe5\xb1\x9e\xef\xbf\xbd\t\t\ttr.safeSetContext'
    new1 = b'\xe8\xae\xbe\xe7\xbd\xae\xe8\xa7\x92\xe8\x89\xb2\xe7\x9a\x84\xe5\x9f\xba\xe6\x9c\xac\xe5\xb1\x9e\xe6\x80\xa7\n\t\t\ttr.safeSetContext'
    content = content.replace(old1, new1)
    
    with open(context_path, 'wb') as f:
        f.write(content)
    
    # 修复 equipment.go
    equipment_path = 'server/internal/test/runner/equipment.go'
    with open(equipment_path, 'rb') as f:
        content = f.read()
    
    # 检查第164行的问题 - 从字节输出看，第164行是 `level := 1`，但它在函数外
    # 但实际上第164行应该是 `TeamSlot: 1,`，这应该是结构体的一部分
    # 让我检查前面的代码结构
    
    with open(equipment_path, 'wb') as f:
        f.write(content)
    
    # 修复 instruction.go
    instruction_path = 'server/internal/test/runner/instruction.go'
    with open(instruction_path, 'rb') as f:
        content = f.read()
    
    # 修复第128行：字符串缺少右引号
    old2 = b'\xe6\x9c\x89\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd\xef\xbf\xbd) {'
    new2 = b'\xe6\x9c\x89\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd") {'
    content = content.replace(old2, new2)
    
    # 修复第130行：字符串缺少右引号
    old3 = b'\xe9\x98\x9f\xe4\xbc\x8d\xe6\x94\xbb\xe5\x87\xbb\xef\xbf\xbd) &&'
    new3 = b'\xe9\x98\x9f\xe4\xbc\x8d\xe6\x94\xbb\xe5\x87\xbb") &&'
    content = content.replace(old3, new3)
    
    # 修复第143行：字符串缺少右引号
    old4 = b'\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd\xef\xbf\xbd) &&'
    new4 = b'\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd") &&'
    content = content.replace(old4, new4)
    
    # 修复第130行：注释和代码混在一起
    old5 = b'\xe8\xa7\x92\xe8\x89\xb2\xe5\x8d\xb8\xe4\xb8\x8b\xe6\xad\xa6\xe5\x99\xa8"\xef\xbf\xbd'
    new5 = b'\xe8\xa7\x92\xe8\x89\xb2\xe5\x8d\xb8\xe4\xb8\x8b\xe6\xad\xa6\xe5\x99\xa8"\xe6\x88\x96"'
    content = content.replace(old5, new5)
    
    # 修复第136行：注释和代码混在一起
    old6 = b'\xe8\xa7\x92\xe8\x89\xb2\xe4\xbe\x9d\xe6\xac\xa1\xe7\xa9\xbf\xe6\x88\xb4\xe6\x89\x80\xe6\x9c\x89\xe8\xa3\x85\xef\xbf\xbd'
    new6 = b'\xe8\xa7\x92\xe8\x89\xb2\xe4\xbe\x9d\xe6\xac\xa1\xe7\xa9\xbf\xe6\x88\xb4\xe6\x89\x80\xe6\x9c\x89\xe8\xa3\x85\xe5\xa4\x87"'
    content = content.replace(old6, new6)
    
    # 修复第142行：注释和代码混在一起
    old7 = b'\xe5\xb7\xb2\xe7\xbb\x8f\xe5\x9c\xa8updateAssertionContext\xe4\xb8\xad\xe5\xa4\x84\xef\xbf\xbd\treturn'
    new7 = b'\xe5\xb7\xb2\xe7\xbb\x8f\xe5\x9c\xa8updateAssertionContext\xe4\xb8\xad\xe5\xa4\x84\xe7\x90\x86\n\t\treturn'
    content = content.replace(old7, new7)
    
    with open(instruction_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_other_files5()
