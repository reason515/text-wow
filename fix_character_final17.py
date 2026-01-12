#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final17():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1247行：注释和代码混在一起
    old1 = b'\xe7\xa1\xae\xe5\xae\x9a\xe6\x9c\x80\xe7\xbb\x88\xe7\x9a\x84MaxHP\xef\xbf\xbd\t\t\tif restoreExplicitMaxHP'
    new1 = b'\xe7\xa1\xae\xe5\xae\x9a\xe6\x9c\x80\xe7\xbb\x88\xe7\x9a\x84MaxHP\n\t\t\tif restoreExplicitMaxHP'
    content = content.replace(old1, new1)
    
    # 修复第1253行：注释和代码混在一起
    old2 = b'\xe4\xbd\xbf\xe7\x94\xa8\xe4\xbf\x9d\xe5\xad\x98\xe7\x9a\x84\xef\xbf\xbd\t\t\t\tchar.MaxHP'
    new2 = b'\xe4\xbd\xbf\xe7\x94\xa8\xe4\xbf\x9d\xe5\xad\x98\xe7\x9a\x84\xe5\x80\xbc\n\t\t\t\tchar.MaxHP'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final17()
