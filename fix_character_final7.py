#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final7():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第525行：注释和代码混在一起
    old1 = b'\xe4\xbe\xbf\xe4\xba\x8e\xe5\x90\x8e\xe7\xbb\xad\xe6\x81\xa2\xef\xbf\xbd\tif explicitHP > 0'
    new1 = b'\xe4\xbe\xbf\xe4\xba\x8e\xe5\x90\x8e\xe7\xbb\xad\xe6\x81\xa2\xe5\xa4\x8d\n\tif explicitHP > 0'
    content = content.replace(old1, new1)
    
    # 修复第1974行：注释和代码混在一起
    old2 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe4\xb8\x8d\xe5\x9c\xa8\xe6\x88\x98\xe6\x96\x97\xe4\xb8\xad\xef\xbc\x8c\xe5\xba\x94\xe8\xaf\xa5\xe4\xb8\xba0\xef\xbf\xbd\tif char.ResourceType'
    new2 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe4\xb8\x8d\xe5\x9c\xa8\xe6\x88\x98\xe6\x96\x97\xe4\xb8\xad\xef\xbc\x8c\xe5\xba\x94\xe8\xaf\xa5\xe4\xb8\xba0\n\tif char.ResourceType'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final7()
