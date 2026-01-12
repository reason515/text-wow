#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中最后的编码问题
"""

def fix_final10():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第584行：注释和代码混在一起
    old1 = b'\xe4\xb8\x8d\xe8\xa6\x81\xe8\xa6\x86\xe7\x9b\x96\xef\xbf\xbd\t\tif savedHP == 0'
    new1 = b'\xe4\xb8\x8d\xe8\xa6\x81\xe8\xa6\x86\xe7\x9b\x96\n\t\tif savedHP == 0'
    content = content.replace(old1, new1)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final10()
