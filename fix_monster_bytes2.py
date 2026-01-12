#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_bytes2():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 从字节序列看：22 29 22 5b 30 5d 是 ")"[0]
    # 应该改为：22 29 22 29 5b 30 5d 是 ")"[0]
    # 即：在 ")" 后面添加 " 和 )
    
    # 模式：")"[0] (字节：22 29 22 5b 30 5d)
    pattern = bytes([0x22, 0x29, 0x22, 0x5b, 0x30, 0x5d])
    # 替换为：")"[0] (字节：22 29 22 29 5b 30 5d)
    replacement = bytes([0x22, 0x29, 0x22, 0x29, 0x5b, 0x30, 0x5d])
    
    if pattern in content:
        count = content.count(pattern)
        content = content.replace(pattern, replacement)
        print(f"修复了 {count} 处字符串字面量问题")
    else:
        print("未找到模式，尝试其他方法")
        # 直接查找并替换
        content = content.replace(b'")"[0]', b'")"[0]')
        print("使用直接替换方法")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_bytes2()
