#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final15():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1108行：注释和代码混在一起
    old1 = b'\xe4\xbb\xa5\xe9\x98\xb2Create\xe5\x90\x8e\xe4\xb8\xa2\xef\xbf\xbd\t\t\t\tsavedPhysicalAttack'
    new1 = b'\xe4\xbb\xa5\xe9\x98\xb2Create\xe5\x90\x8e\xe4\xb8\xa2\xe5\xa4\xb1\n\t\t\t\tsavedPhysicalAttack'
    content = content.replace(old1, new1)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final15()
