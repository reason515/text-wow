#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
使用正确的字节序列修复 character.go 文件
"""

def fix_correct():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第562行：注释和代码混在一起
    # 查找：基础HP=35" + 替换字符(0xef 0xbf 0xbd) + tab + if
    old1 = b'\xe5\x9f\xba\xe7\xa1\x80HP=35"\xef\xbf\xbd\tif'
    # 替换为：基础HP=35"）\n\tif
    new1 = b'\xe5\x9f\xba\xe7\xa1\x80HP=35"\xef\xbc\x89\n\tif'
    content = content.replace(old1, new1)
    
    # 修复第4541行：字符串未闭合
    # 查找：一个角 + 替换字符(0xef 0xbf 0xbd) + , 1)
    old2 = b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xef\xbf\xbd, 1)'
    # 替换为：一个角色", 1)
    new2 = b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xe8\x89\xb2", 1)'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_correct()
