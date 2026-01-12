#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中更多的编码问题
"""

def fix_more2():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第352行：注释和代码混在一起
    old1 = b'\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87=150%"\xef\xbf\xbd\tif'
    new1 = b'\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe4\xbc\xa4\xe5\xae\xb3=150%"\xef\xbc\x89\n\tif'
    content = content.replace(old1, new1)
    
    # 修复第319行：注释和代码混在一起
    old2 = b'\xe6\x97\xb6\xe8\xae\xbe\xef\xbf\xbd\t\t\t\ttr.context.Variables["character_gold"]'
    new2 = b'\xe6\x97\xb6\xe8\xae\xbe\xe7\xbd\xae\n\t\t\t\ttr.context.Variables["character_gold"]'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_more2()
