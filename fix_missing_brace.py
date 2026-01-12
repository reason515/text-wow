#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def fix_missing_brace():
    filepath = 'server/internal/test/runner/instruction.go'
    
    with open(filepath, 'rb') as f:
        lines = f.readlines()
    
    # 检查第212行到第264行之间的块
    # 第212行开始了一个 else if 块
    # 第263行有 return nil
    # 第265行有下一个 else if
    
    # 检查第264行是否有 }
    if len(lines) > 263:
        line_264 = lines[263]  # 第264行，索引从0开始
        print(f"第264行内容: {repr(line_264)}")
        
        # 如果第264行没有 }，需要添加
        if b'}' not in line_264.strip():
            print("第264行缺少 }，需要添加")
            # 在第264行添加 }
            lines.insert(263, b'\t}\n')
            print("已添加 }")
    
    with open(filepath, 'wb') as f:
        f.writelines(lines)

if __name__ == '__main__':
    fix_missing_brace()
