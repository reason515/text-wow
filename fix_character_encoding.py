#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复 character.go 文件中的编码问题
"""

def fix_character():
    file_path = 'server/internal/test/runner/character.go'
    
    # 读取为字节
    with open(file_path, 'rb') as f:
        content_bytes = f.read()
    
    # 替换 UTF-8 替换字符序列 "\xef\xbf\xbd)[0]" 为 "）")[0]
    content_bytes = content_bytes.replace(b'\xef\xbf\xbd)[0]', b'\xef\xbc\x89")[0]')  # "）")[0]
    
    # 解码
    try:
        content = content_bytes.decode('utf-8')
    except:
        content = content_bytes.decode('utf-8', errors='replace')
    
    # 清理 import 块中的空行和函数体内的多余空行
    lines = content.split('\n')
    fixed_lines = []
    in_import = False
    prev_was_blank = False
    
    for line in lines:
        stripped = line.strip()
        is_blank = not stripped
        
        # 检测 import 块
        if stripped == 'import (':
            in_import = True
            fixed_lines.append('import (\n')
            continue
        
        if in_import and stripped == ')':
            in_import = False
            fixed_lines.append(')\n')
            continue
        
        # import 块内：移除所有空行
        if in_import:
            if is_blank:
                continue
            fixed_lines.append(line)
            prev_was_blank = False
            continue
        
        # 处理空行：函数体内最多保留一个空行
        if is_blank:
            if not prev_was_blank:
                if fixed_lines:
                    last_line = fixed_lines[-1].strip()
                    if last_line == '}' or last_line.startswith('//'):
                        fixed_lines.append(line)
                        prev_was_blank = True
            continue
        
        # 非空行
        if prev_was_blank and fixed_lines:
            last_line = fixed_lines[-1].strip()
            if last_line and not last_line.endswith('{') and not last_line.endswith('('):
                pass
            else:
                while fixed_lines and not fixed_lines[-1].strip():
                    fixed_lines.pop()
                if fixed_lines and fixed_lines[-1].strip():
                    fixed_lines.append('')
        
        fixed_lines.append(line)
        prev_was_blank = False
    
    # 移除文件末尾的空行
    while fixed_lines and not fixed_lines[-1].strip():
        fixed_lines.pop()
    
    # 写入文件
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write('\n'.join(fixed_lines))
        if fixed_lines:
            f.write('\n')
    
    print("修复完成！")

if __name__ == '__main__':
    fix_character()
