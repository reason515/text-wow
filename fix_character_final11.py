#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final11():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第631行：注释和代码混在一起
    old1 = b'\xe7\xa1\xae\xe4\xbf\x9d\xe8\xa7\x92\xe8\x89\xb2\xe6\x9c\x89\xe5\xbf\x85\xe9\x9c\x80\xe7\x9a\x84\xe5\xad\x97\xef\xbf\xbd\tif char.RaceID'
    new1 = b'\xe7\xa1\xae\xe4\xbf\x9d\xe8\xa7\x92\xe8\x89\xb2\xe6\x9c\x89\xe5\xbf\x85\xe9\x9c\x80\xe7\x9a\x84\xe5\xad\x97\xe6\xae\xb5\n\tif char.RaceID'
    content = content.replace(old1, new1)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final11()
