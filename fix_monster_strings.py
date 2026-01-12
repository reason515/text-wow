#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_monster_strings():
    filepath = 'server/internal/test/runner/monster.go'
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    # 修复第53行的注释问题
    # 查找 "解析数量（如"创建3个怪物"	count := 1" 的模式
    pattern1 = b'\xe8\xa7\xa3\xe6\x9e\x90\xe6\x95\xb0\xe9\x87\x8f\xef\xbc\x88\xe5\xa6\x82"\xe5\x88\x9b\xe5\xbb\xba3\xe4\xb8\xaa\xe6\x80\xaa\xe7\x89\xa9"\xef\xbf\xbd\tcount := 1'
    replacement1 = b'\xe8\xa7\xa3\xe6\x9e\x90\xe6\x95\xb0\xe9\x87\x8f\xef\xbc\x88\xe5\xa6\x82"\xe5\x88\x9b\xe5\xbb\xba3\xe4\xb8\xaa\xe6\x80\xaa\xe7\x89\xa9"\xef\xbc\x89\n\tcount := 1'
    
    if pattern1 in content:
        content = content.replace(pattern1, replacement1)
        print("修复了第53行的注释问题")
    
    # 修复第55行的字符串问题
    # 查找 "if strings.Contains(instruction, ") {" 的模式
    pattern2 = b'if strings.Contains(instruction, "\xef\xbf\xbd) {'
    replacement2 = b'if strings.Contains(instruction, "\xe4\xb8\xaa") {'
    
    if pattern2 in content:
        content = content.replace(pattern2, replacement2)
        print("修复了第55行的字符串问题")
    else:
        # 尝试其他可能的编码
        pattern2a = b'if strings.Contains(instruction, "'
        pos = content.find(pattern2a, 0)
        if pos != -1:
            # 查找这个模式后面跟着特殊字符的位置
            after_pattern = content[pos+len(pattern2a):pos+len(pattern2a)+20]
            if b'\xef\xbf\xbd)' in after_pattern:
                # 替换为正确的格式
                new_pattern = b'if strings.Contains(instruction, "\xe4\xb8\xaa") {'
                end_pos = content.find(b') {', pos)
                if end_pos != -1:
                    content = content[:pos] + new_pattern + content[end_pos+3:]
                    print("修复了第55行的字符串问题（方法2）")
    
    # 修复第57行的字符串问题
    # 查找 'parts := strings.Split(instruction, ")' 的模式
    pattern3 = b'parts := strings.Split(instruction, "\xef\xbf\xbd)'
    replacement3 = b'parts := strings.Split(instruction, "\xe4\xb8\xaa")'
    
    if pattern3 in content:
        content = content.replace(pattern3, replacement3)
        print("修复了第57行的字符串问题")
    
    with open(filepath, 'wb') as f:
        f.write(content)

if __name__ == '__main__':
    fix_monster_strings()
