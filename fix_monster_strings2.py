#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_strings2():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 修复字符串字面量问题：将 "\xef\xbf\xbd)[0]" 替换为 ")[0]"
    pattern1 = b'strings.Split(parts[1], "\\xef\\xbf\\xbd)[0])'
    replacement1 = b'strings.Split(parts[1], ")")[0]'
    
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        print("修复了 strings.Split(parts[1], ...) 的问题")
    
    # 修复字符串字面量问题：将 "\xef\xbf\xbd)[0]" 替换为 ")[0]"
    pattern2 = b'strings.Split(defenseStr, "\\xef\\xbf\\xbd)[0])'
    replacement2 = b'strings.Split(defenseStr, ")")[0]'
    
    if pattern2 in content:
        content = content.replace(pattern2, replacement2)
        print("修复了 strings.Split(defenseStr, ...) 的问题")
    
    # 修复字符串字面量问题：将 "\xef\xbf\xbd)[0]" 替换为 ")[0]"
    pattern3 = b'strings.Split(parts[1], "\\xef\\xbf\\xbd)[0])'
    replacement3 = b'strings.Split(parts[1], ")")[0]'
    
    if pattern3 in content:
        content = content.replace(pattern3, replacement3)
        print("修复了 strings.Split(parts[1], ...) 的问题（模式3）")
    
    # 使用字节级替换
    content = content.replace(b'\xef\xbf\xbd)[0]', b')"[0]')
    print("使用字节级替换修复了所有字符串字面量问题")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_strings2()
