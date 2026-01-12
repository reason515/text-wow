#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_correct():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 修复字符串字面量问题：将 '")"[0]' 替换为 '")"[0]'
    # 模式：strings.Split(..., ")"[0]
    # 应该改为：strings.Split(..., ")")[0]
    
    # 使用字节级替换
    content = content.replace(b'")"[0]', b'")"[0]')
    print("修复了字符串字面量问题")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_correct()
