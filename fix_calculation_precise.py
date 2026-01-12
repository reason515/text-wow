#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
精确修复 calculation.go 文件中的字符串分割问题
"""

def fix_precise():
    file_path = 'server/internal/test/runner/calculation.go'
    
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 尝试解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    lines = content.split('\n')
    fixed_lines = []
    
    for i, line in enumerate(lines):
        # 修复所有包含 strings.Split 和 )[0]) 的行
        if 'strings.Split' in line and ')[0])' in line:
            # 使用正则表达式替换
            import re
            # 匹配 ")[0]) 并替换为 "）")[0])
            line = re.sub(r'"\)\[0\]\)', '"）")[0])', line)
        
        fixed_lines.append(line)
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write('\n'.join(fixed_lines))
        if fixed_lines and not fixed_lines[-1].endswith('\n'):
            f.write('\n')
    
    print("修复完成！")

if __name__ == '__main__':
    fix_precise()
