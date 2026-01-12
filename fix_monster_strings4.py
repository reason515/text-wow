#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_strings4():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 修复字符串字面量问题：将 '")"[0]' 替换为 '")"[0]'
    # 实际上应该是 strings.Split(..., ")")[0]
    # 查找模式：strings.Split(..., ")"[0])
    # 替换为：strings.Split(..., ")")[0]
    
    # 使用正则表达式或简单的字符串替换
    import re
    
    # 模式1: strings.Split(parts[1], ")"[0])
    pattern1 = b'strings\\.Split\\(parts\\[1\\], "\\"\\)"\\[0\\]\\)'
    replacement1 = b'strings.Split(parts[1], ")")[0]'
    
    # 模式2: strings.Split(defenseStr, ")"[0])
    pattern2 = b'strings\\.Split\\(defenseStr, "\\"\\)"\\[0\\]\\)'
    replacement2 = b'strings.Split(defenseStr, ")")[0]'
    
    # 模式3: strings.Split(parts[1], ")"[0])
    pattern3 = b'strings\\.Split\\(parts\\[1\\], "\\"\\)"\\[0\\]\\)'
    replacement3 = b'strings.Split(parts[1], ")")[0]'
    
    # 简单的字节级替换
    content = content.replace(b'")"[0])', b'")"[0]')
    print("使用字节级替换修复了字符串字面量问题")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_strings4()
