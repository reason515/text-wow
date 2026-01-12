#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final25():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第1402行：注释和代码混在一起，缺少右括号
    old1 = b'\xe5\x8d\xb3\xe4\xbd\xbf\xe5\xbd\x93\xe5\x89\x8d\xe5\x80\xbc\xe4\xb8\x8d\xe4\xb8\xba0\xef\xbf\xbd\tif !explicitPhysicalAttack'
    new1 = b'\xe5\x8d\xb3\xe4\xbd\xbf\xe5\xbd\x93\xe5\x89\x8d\xe5\x80\xbc\xe4\xb8\x8d\xe4\xb8\xba0\xef\xbc\x89\n\tif !explicitPhysicalAttack'
    content = content.replace(old1, new1)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final25()
