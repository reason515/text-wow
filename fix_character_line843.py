#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件第843行的编码问题
"""

def fix_line843():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第843行：注释和代码混在一起
    # 查找：防止后续被覆 + 替换字符 + tab + tr.context.Variables
    old = b'\xe9\x98\xb2\xe6\xad\xa2\xe5\x90\x8e\xe7\xbb\xad\xe8\xa2\xab\xe8\xa6\x86\xef\xbf\xbd\t\t\t\ttr.context.Variables'
    # 替换为：防止后续被覆盖\n\t\t\t\ttr.context.Variables
    new = b'\xe9\x98\xb2\xe6\xad\xa2\xe5\x90\x8e\xe7\xbb\xad\xe8\xa2\xab\xe8\xa6\x86\xe7\x9b\x96\n\t\t\t\ttr.context.Variables'
    content = content.replace(old, new)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_line843()
