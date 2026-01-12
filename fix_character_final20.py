#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final20():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1262行：注释和代码混在一起
    old1 = b'\xe4\xbd\xbf\xe7\x94\xa8\xe8\xae\xa1\xe7\xae\x97\xef\xbf\xbd\t\t\t\tchar.MaxHP'
    new1 = b'\xe4\xbd\xbf\xe7\x94\xa8\xe8\xae\xa1\xe7\xae\x97\xe5\x80\xbc\n\t\t\t\tchar.MaxHP'
    content = content.replace(old1, new1)
    
    # 修复第1266行：注释和代码混在一起
    old2 = b'\xe7\xa1\xae\xe5\xae\x9a\xe6\x9c\x80\xe7\xbb\x88\xe7\x9a\x84HP\xef\xbf\xbd\t\t\t\tif restoreExplicitHP'
    new2 = b'\xe7\xa1\xae\xe5\xae\x9a\xe6\x9c\x80\xe7\xbb\x88\xe7\x9a\x84HP\n\t\t\t\tif restoreExplicitHP'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final20()
