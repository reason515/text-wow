#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_newline_brace():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: "\xe9\xa5\xb0\xe5\x93\x81"))\r\n\t\t { 应该是 "\xe9\xa5\xb0\xe5\x93\x81")) {
    pattern1 = b'\xe9\xa5\xb0\xe5\x93\x81"))\r\n\t\t {'
    replacement1 = b'\xe9\xa5\xb0\xe5\x93\x81")) {'
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        print(f"  修复了模式1")
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_newline_brace()
