#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_bytes():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 查找模式：")"[0] (字节序列：22 29 22 5b 30 5d)
    # 应该改为：")"[0] (字节序列：22 29 22 29 5b 30 5d)
    # 即：在 ")" 后面添加一个 " 和 )
    
    # 模式：")"[0]
    pattern = b'")"[0]'
    replacement = b'")"[0]'
    
    if pattern in content:
        content = content.replace(pattern, replacement)
        print(f"修复了 {content.count(replacement)} 处字符串字面量问题")
    else:
        print("未找到模式")
        # 尝试查找类似的模式
        # 查找包含 ")"[0] 的行
        lines = content.split(b'\n')
        for i, line in enumerate(lines):
            if b'strings.Split' in line and b'")' in line and b'[0]' in line:
                print(f"第 {i+1} 行可能有问题: {repr(line[:100])}")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_bytes()
