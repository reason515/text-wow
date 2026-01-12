#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
最终修复 character.go 文件中的所有编码问题
"""

def fix_final():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换所有 UTF-8 替换字符序列
    replacements = [
        (b'\xef\xbf\xbd)', b':'),  # 替换字符后跟 ) 改为 :
        (b'\xef\xbf\xbd=', b'='),  # 替换字符后跟 = 改为 =
        (b'\xef\xbf\xbd[0]', b'[0]'),  # 修复 [0]
    ]
    
    for old, new in replacements:
        content_bytes = content_bytes.replace(old, new)
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 修复特定的编码问题
    import re
    # 修复 "30" 相关的问题
    content = re.sub(r'if strings\.Contains\(instruction, "30[^\"]*"\)', 
                     'if strings.Contains(instruction, "30级")', content)
    
    # 修复 "创建多个角色" 相关的问题
    content = re.sub(r'if strings\.Contains\(instruction, "[^\"]*\) \{', 
                     lambda m: m.group(0).replace(')', '":') if '' in m.group(0) else m.group(0), content)
    
    # 修复 strings.Split 中的编码问题
    content = re.sub(r'strings\.Split\(instruction, "[^\"]*\)', 
                     lambda m: m.group(0).replace(')', '":') if '' in m.group(0) else m.group(0), content)
    content = re.sub(r'strings\.Split\(parts\[1\], "[^\"]*\)', 
                     lambda m: m.group(0).replace(')', '",') if '' in m.group(0) else m.group(0), content)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(content)
    
    print("修复完成！")

if __name__ == '__main__':
    fix_final()
