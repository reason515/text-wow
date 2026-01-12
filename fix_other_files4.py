#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复其他文件中的编码问题
"""

def fix_other_files4():
    # 修复 instruction.go
    instruction_path = 'server/internal/test/runner/instruction.go'
    with open(instruction_path, 'rb') as f:
        content = f.read()
    
    # 修复第119行：字符串缺少右引号
    old1 = b'\xe6\xac\xa1\xe6\x94\xbb\xef\xbf\xbd) {'
    new1 = b'\xe6\xac\xa1\xe6\x94\xbb\xe5\x87\xbb") {'
    content = content.replace(old1, new1)
    
    # 修复第125行：字符串缺少右引号
    old2 = b'\xe8\xae\xa1\xe7\xae\x97\xe9\x98\x9f\xe4\xbc\x8d\xe6\x80\xbb\xe7\x94\x9f\xe5\x91\xbd\xef\xbf\xbd) {'
    new2 = b'\xe8\xae\xa1\xe7\xae\x97\xe9\x98\x9f\xe4\xbc\x8d\xe6\x80\xbb\xe7\x94\x9f\xe5\x91\xbd") {'
    content = content.replace(old2, new2)
    
    # 修复第124行：注释和代码混在一起
    old3 = b'\xe8\xa7\x92\xe8\x89\xb2\xe7\xa9\xbf\xe6\x88\xb4\xe6\xad\xa6\xe5\x99\xa8"\xef\xbf\xbd'
    new3 = b'\xe8\xa7\x92\xe8\x89\xb2\xe7\xa9\xbf\xe6\x88\xb4\xe6\xad\xa6\xe5\x99\xa8"\xe6\x88\x96"'
    content = content.replace(old3, new3)
    
    with open(instruction_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_other_files4()
