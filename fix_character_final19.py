#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final19():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1303行：注释和代码混在一起
    old1 = b'\xe4\xbb\x8eVariables\xe6\x81\xa2\xe5\xa4\x8d\xef\xbf\xbd\tif strengthVal'
    new1 = b'\xe4\xbb\x8eVariables\xe6\x81\xa2\xe5\xa4\x8d\n\tif strengthVal'
    content = content.replace(old1, new1)
    
    # 修复第1284行：注释中的编码问题
    old2 = b'\xe5\xa6\x82\xe6\x9e\x9cHP\xef\xbf\xbd\xe6\x88\x96\xe8\xb6\x85\xe8\xbf\x87MaxHP'
    new2 = b'\xe5\xa6\x82\xe6\x9e\x9cHP\xe4\xb8\xba0\xe6\x88\x96\xe8\xb6\x85\xe8\xbf\x87MaxHP'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final19()
