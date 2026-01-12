#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_comment_encoding():
    # 修复 context.go 第337行的编码问题
    filepath = 'server/internal/test/runner/context.go'
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    
    # 修复: "也设" + 制表符 + "// 这需要" 应该是 "也设置\n\t\t\t// 这需要"
    # 查找包含 "也设" 和 "这需要" 的行
    import re
    pattern1 = rb'(\xe4\xb9\x9f\xe8\xae\xbe)[^\n]*(\xe8\xbf\x99\xe9\x9c\x80\xe8\xa6\x81)'
    def replace1(m):
        return b'\xe4\xb9\x9f\xe8\xae\xbe\xe7\xbd\xae\n\t\t\t// ' + m.group(2)
    
    new_content = re.sub(pattern1, replace1, content)
    if new_content != content:
        content = new_content
        print(f"  修复了注释中的编码问题")
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_comment_encoding()
