#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_strings_close():
    filepath = 'server/internal/test/runner/instruction.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复: "\xe8\xa3\x85\xe5\xa4\x87")) 应该是 "\xe8\xa3\x85\xe5\xa4\x87")
    # 查找所有 "text")) 模式，应该是 "text")
    import re
    pattern = rb'("[\x20-\x7e\x80-\xff]+"\))\)'
    def replace(m):
        text = m.group(1)
        return text
    
    new_content = re.sub(pattern, replace, content)
    if new_content != content:
        changes += 1
        print(f"  修复了多余的右括号")
        content = new_content
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_strings_close()
