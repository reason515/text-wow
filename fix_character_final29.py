#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final29():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1558行：注释和代码混在一起
    old1 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe8\xae\xbe\xe7\xbd\xae\xe4\xba\x86\xef\xbf\xbd\tif goldVal'
    new1 = b'\xe5\xa6\x82\xe6\x9e\x9c\xe8\xae\xbe\xe7\xbd\xae\xe4\xba\x86\n\tif goldVal'
    content = content.replace(old1, new1)
    
    # 修复第1562行：注释和代码混在一起，SQL语句也有问题
    old2 = b'\xe6\x95\xb0\xe6\x8d\xae\xe5\xba\x93\xe4\xb8\xad\xe7\x9a\x84\xe7\x94\xa8\xe6\x88\xb7\xe9\x87\x91\xef\xbf\xbd\t\t\t_, err := database.DB.Exec(`UPDATE users SET gold =  WHERE id = `, gold, char.UserID)'
    new2 = b'\xe6\x95\xb0\xe6\x8d\xae\xe5\xba\x93\xe4\xb8\xad\xe7\x9a\x84\xe7\x94\xa8\xe6\x88\xb7\xe9\x87\x91\xe5\xb8\x81\n\t\t\t_, err := database.DB.Exec(`UPDATE users SET gold = ? WHERE id = ?`, gold, char.UserID)'
    content = content.replace(old2, new2)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final29()
