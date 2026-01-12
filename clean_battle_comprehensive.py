#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
全面清理 battle.go 文件中的空行和修复语法错误
"""

import re

def clean_battle_file():
    file_path = 'server/internal/test/runner/battle.go'
    
    with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    
    lines = content.split('\n')
    cleaned_lines = []
    in_import_block = False
    import_block_start = -1
    import_block_end = -1
    prev_blank = False
    in_function = False
    brace_count = 0
    
    i = 0
    while i < len(lines):
        line = lines[i]
        stripped = line.strip()
        
        # 检测 import 块
        if stripped == 'import (':
            in_import_block = True
            import_block_start = len(cleaned_lines)
            cleaned_lines.append(line)
            prev_blank = False
            i += 1
            continue
        
        if in_import_block:
            if stripped == ')' or (stripped.startswith(')') and not stripped.startswith(')')):
                # 检查是否是 import 块的结束
                if stripped == ')':
                    in_import_block = False
                    import_block_end = len(cleaned_lines)
                    cleaned_lines.append(line)
                    prev_blank = False
                    i += 1
                    continue
                else:
                    # 可能是其他情况，继续处理
                    pass
        
        # 在 import 块内：移除所有空行
        if in_import_block:
            if not stripped:
                i += 1
                continue
            cleaned_lines.append(line)
            prev_blank = False
            i += 1
            continue
        
        # 检测函数开始
        if re.match(r'^\s*func\s+', stripped):
            # 如果前面有多个空行，只保留一个
            if prev_blank:
                # 移除最后一个空行
                cleaned_lines.pop()
            # 添加一个空行（如果前面不是空行）
            if cleaned_lines and cleaned_lines[-1].strip():
                cleaned_lines.append('')
            cleaned_lines.append(line)
            prev_blank = False
            in_function = True
            brace_count = stripped.count('{') - stripped.count('}')
            i += 1
            continue
        
        # 检测函数结束
        if in_function and stripped == '}':
            brace_count += stripped.count('}') - stripped.count('{')
            if brace_count <= 0:
                in_function = False
                brace_count = 0
        
        # 修复第1331行的编码问题
        if i == 1330:  # 0-based index
            if '// 更新上下' in line and 'tr.context.Characters["character"]' in line:
                # 分离注释和代码
                cleaned_lines.append('	// 更新上下文')
                cleaned_lines.append('	tr.context.Characters["character"] = char')
                prev_blank = False
                i += 1
                continue
        
        # 处理空行
        if not stripped:
            # 在函数体内：最多保留一个空行（如果前面不是空行）
            if in_function:
                if not prev_blank:
                    cleaned_lines.append(line)
                    prev_blank = True
            else:
                # 顶级声明之间：只保留一个空行
                if not prev_blank:
                    cleaned_lines.append(line)
                    prev_blank = True
            i += 1
            continue
        
        # 非空行
        # 如果前面有多个空行，减少到最多一个
        if prev_blank and cleaned_lines:
            # 检查是否真的需要空行
            last_line = cleaned_lines[-1].strip()
            if last_line and not last_line.endswith('{') and not last_line.endswith('('):
                # 保留空行
                pass
            else:
                # 移除多余的空行
                while cleaned_lines and not cleaned_lines[-1].strip():
                    cleaned_lines.pop()
                if cleaned_lines and cleaned_lines[-1].strip():
                    cleaned_lines.append('')
        
        cleaned_lines.append(line)
        prev_blank = False
        i += 1
    
    # 移除文件末尾的空行
    while cleaned_lines and not cleaned_lines[-1].strip():
        cleaned_lines.pop()
    
    # 写入文件
    output = '\n'.join(cleaned_lines)
    if not output.endswith('\n'):
        output += '\n'
    
    with open(file_path, 'w', encoding='utf-8', errors='replace') as f:
        f.write(output)
    
    print(f"清理完成！移除了大量空行，文件现在有 {len(cleaned_lines)} 行")

if __name__ == '__main__':
    clean_battle_file()
