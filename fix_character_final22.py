#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final22():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1158行：注释缺少右括号
    old1 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe5\xae\x83\xe4\xbb\xac\xe9\x83\xbd\xe4\xb8\x8d\xe4\xb8\xba0\r\n'
    new1 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe5\xae\x83\xe4\xbb\xac\xe9\x83\xbd\xe4\xb8\x8d\xe4\xb8\xba0\xef\xbc\x89\r\n'
    content = content.replace(old1, new1)
    
    # 修复第1208行：注释和代码混在一起
    old2 = b'\xe7\x9a\x84MaxHP\xef\xbf\xbd\t\t\trestoreExplicitMaxHP'
    new2 = b'\xe7\x9a\x84MaxHP\n\t\t\trestoreExplicitMaxHP'
    content = content.replace(old2, new2)
    
    # 修复第1220行：注释和代码混在一起
    old3 = b'\xe7\x9a\x84HP\xef\xbf\xbd\t\t\trestoreExplicitHP'
    new3 = b'\xe7\x9a\x84HP\n\t\t\trestoreExplicitHP'
    content = content.replace(old3, new3)
    
    # 修复第1267行：注释和代码混在一起
    old4 = b'\xe7\x9a\x84HP\xef\xbf\xbd\t\t\tif restoreExplicitHP'
    new4 = b'\xe7\x9a\x84HP\n\t\t\tif restoreExplicitHP'
    content = content.replace(old4, new4)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final22()
