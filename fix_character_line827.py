#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件第827行的编码问题
"""

def fix_line827():
    file_path = 'server/internal/test/runner/character.go'
    
    with open(file_path, 'rb') as f:
        content = f.read()
    
    # 修复第827行：整个代码行被破坏
    # 查找：解析暴击率（ + 替换字符 + 物理暴击率= + strings.Split...
    old = b'\t// \xe8\xa7\xa3\xe6\x9e\x90\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87\xef\xbc\x88\xef\xbf\xbd\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87= strings.Split(instruction, "\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87= strings.TrimSpace(strings.Split(parts[1], "%")[0])'
    # 替换为完整的代码
    new = b'\t// \xe8\xa7\xa3\xe6\x9e\x90\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87\xef\xbc\x88\xe5\xa6\x82"\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87=30%"\xef\xbc\x89\n\tif strings.Contains(instruction, "\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87=") {\n\t\tparts := strings.Split(instruction, "\xe7\x89\xa9\xe7\x90\x86\xe6\x9a\xb4\xe5\x87\xbb\xe7\x8e\x87=")\n\t\tif len(parts) > 1 {\n\t\t\tcritStr := strings.TrimSpace(strings.Split(parts[1], "%")[0])'
    content = content.replace(old, new)
    
    # 写入文件
    with open(file_path, 'wb') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_line827()
