#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_context_equipment_final():
    # 修复 context.go 第337行的编码问题
    filepath = 'server/internal/test/runner/context.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: "也设" 应该是 "也设置"
    # 查找: "也设" + 制表符 + "// 这需要"
    pattern1 = b'\xe4\xb9\x9f\xe8\xae\xbe\t\t\t// \xe8\xbf\x99\xe9\x9c\x80\xe8\xa6\x81'
    replacement1 = b'\xe4\xb9\x9f\xe8\xae\xbe\xe7\xbd\xae\n\t\t\t// \xe8\xbf\x99\xe9\x9c\x80\xe8\xa6\x81'
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        print(f"  修复了第337行的编码问题")
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")
    
    # 检查 equipment.go - 第164行看起来正常，可能是前面的函数没有正确关闭
    filepath2 = 'server/internal/test/runner/equipment.go'
    print(f"\n检查文件: {filepath2}")
    with open(filepath2, 'rb') as f:
        content2 = f.read()
    
    # 检查第164行之前是否有编码问题
    lines = content2.split(b'\n')
    if len(lines) >= 164:
        print(f"  第164行内容: {repr(lines[163][:150])}")
        print(f"  第163行内容: {repr(lines[162][:150])}")
        print(f"  第162行内容: {repr(lines[161][:150])}")

if __name__ == '__main__':
    fix_context_equipment_final()
