#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中剩余的编码问题
"""

def fix_final2():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第4673行：注释和代码混在一起
    # 查找：检查该slot是否已存在角 + 替换字符 + tab + if
    old1 = b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xaf\xa5slot\xe6\x98\xaf\xe5\x90\xa6\xe5\xb7\xb2\xe5\xad\x98\xe5\x9c\xa8\xe8\xa7\x92\xef\xbf\xbd\tif'
    # 替换为：检查该slot是否已存在角色\n\t\tif
    new1 = b'\xe6\xa3\x80\xe6\x9f\xa5\xe8\xaf\xa5slot\xe6\x98\xaf\xe5\x90\xa6\xe5\xb7\xb2\xe5\xad\x98\xe5\x9c\xa8\xe8\xa7\x92\xe8\x89\xb2\n\t\tif'
    content = content.replace(old1, new1)
    
    # 修复第4701行：注释和代码混在一起
    # 查找：恢复保存的属性 + 替换字符 + tab + char.Agility
    old2 = b'\xe6\x81\xa2\xe5\xa4\x8d\xe4\xbf\x9d\xe5\xad\x98\xe7\x9a\x84\xe5\xb1\x9e\xe6\x80\xa7\xef\xbf\xbd\tchar.Agility'
    # 替换为：恢复保存的属性\n\t\t\tchar.Agility
    new2 = b'\xe6\x81\xa2\xe5\xa4\x8d\xe4\xbf\x9d\xe5\xad\x98\xe7\x9a\x84\xe5\xb1\x9e\xe6\x80\xa7\n\t\t\tchar.Agility'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final2()
