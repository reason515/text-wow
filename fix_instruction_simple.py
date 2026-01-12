#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_instruction_simple():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复: "\xe8\xa3\x85\xe5\xa4\x87"))\r\r\n\t\t { 应该是 "\xe8\xa3\x85\xe5\xa4\x87")) {
    pattern1 = b'\xe8\xa3\x85\xe5\xa4\x87")\\)\r\r\n\t\t {'
    replacement1 = b'\xe8\xa3\x85\xe5\xa4\x87")) {'
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        changes += 1
        print(f"  修复了模式1")
    
    # 修复: "\xe8\xa3\x85\xe5\xa4\x87"))\r\r\n 应该是 "\xe8\xa3\x85\xe5\xa4\x87"))
    pattern2 = b'\xe8\xa3\x85\xe5\xa4\x87")\\)\r\r\n'
    replacement2 = b'\xe8\xa3\x85\xe5\xa4\x87"))'
    if pattern2 in content:
        content = content.replace(pattern2, replacement2)
        changes += 1
        print(f"  修复了模式2")
    
    # 修复所有类似的模式: "text"))\r\r\n\t\t { 应该是 "text")) {
    import re
    pattern3 = rb'("[\x20-\x7e\x80-\xff]+"\))\)\r\r\n\t\t\s*\{'
    def replace3(m):
        return m.group(1) + b') {'
    
    new_content = re.sub(pattern3, replace3, content)
    if new_content != content:
        changes += 1
        print(f"  修复了通用模式")
        content = new_content
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_simple()
