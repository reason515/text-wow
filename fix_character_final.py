#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中剩余的编码问题
"""

def fix_final():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 修复第562行：注释和代码混在一起
    # 查找 "基础HP=35" 后面跟着 if 的模式
    content_bytes = content_bytes.replace(
        b'\xe5\x9f\xba\xe7\xa1\x80HP=35"\xef\xbf\xbd\tif',
        b'\xe5\x9f\xba\xe7\xa1\x80HP=35"）\n\tif'
    )
    
    # 修复第4541行：字符串未闭合
    # 查找 "一个角" 后面跟着 , 1) 的模式
    content_bytes = content_bytes.replace(
        b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xef\xbf\xbd, 1)',
        b'\xe4\xb8\x80\xe4\xb8\xaa\xe8\xa7\x92\xe8\x89\xb2", 1)'
    )
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content_bytes)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final()
