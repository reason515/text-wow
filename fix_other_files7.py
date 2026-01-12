#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_encoding_issues():
    files_to_fix = [
        'server/internal/test/runner/context.go',
        'server/internal/test/runner/equipment.go',
        'server/internal/test/runner/instruction.go',
    ]
    
    for filepath in files_to_fix:
        print(f"\n处理文件: {filepath}")
        
        with open(filepath, 'rb') as f:
            content = f.read()
        
        original_content = content
        changes = 0
        
        # 修复 strings.Split 中的 \xef\xbf\xbd) 模式
        # 查找 patterns like: strings.Split(..., "...\xef\xbf\xbd)")
        pattern1 = rb'strings\.Split\([^,]+,\s*"([^"]*)\xef\xbf\xbd\)'
        def replace_split1(m):
            text = m.group(1).decode('utf-8', errors='ignore')
            return f'strings.Split({m.group(0).split(b",")[0].decode("utf-8", errors="ignore")}, "{text}")'.encode('utf-8')
        
        # 更简单的替换：直接替换 \xef\xbf\xbd) 为 )
        content = content.replace(b'\xef\xbf\xbd)', b')')
        if content != original_content:
            changes += 1
            print(f"  修复了 strings.Split 中的替换字符问题")
        
        # 修复 strings.Contains 中的 \xef\xbf\xbd) 模式
        content = content.replace(b'\xef\xbf\xbd)', b')')
        
        # 修复注释中的编码问题 - 查找被截断的注释
        # 修复 "处" 为 "处理"
        content = content.replace(b'\xe5\xa4\x84\xef\xbf\xbd', b'\xe5\xa4\x84\xe7\x90\x86')  # 处 -> 处理
        if content != original_content:
            changes += 1
            print(f"  修复了注释中的编码问题")
        
        # 修复其他常见的编码问题
        # 修复 "一" 为 "一个"
        content = content.replace(b'\xe4\xb8\x80\xef\xbf\xbd', b'\xe4\xb8\x80\xe4\xb8\xaa')
        # 修复 "牧" 为 "牧师"
        content = content.replace(b'\xe7\x89\xa7\xef\xbf\xbd', b'\xe7\x89\xa7\xe5\xb8\x88')
        # 修复 "法" 为 "法师"
        content = content.replace(b'\xe6\xb3\x95\xef\xbf\xbd', b'\xe6\xb3\x95\xe5\xb8\x88')
        # 修复 "角" 为 "角色"
        content = content.replace(b'\xe8\xa7\x92\xef\xbf\xbd', b'\xe8\xa7\x92\xe8\x89\xb2')
        # 修复 "包" 为 "包含"
        content = content.replace(b'\xe5\x8c\x85\xef\xbf\xbd', b'\xe5\x8c\x85\xe5\x90\xab')
        # 修复 "个" 为 "一个"
        content = content.replace(b'\xef\xbf\xbd\xe4\xb8\xaa', b'\xe4\xb8\x80\xe4\xb8\xaa')
        # 修复 "敏" 为 "敏捷"
        content = content.replace(b'\xe6\x95\x8f\xef\xbf\xbd', b'\xe6\x95\x8f\xe6\x8d\xb7')
        # 修复 "50" 为 "50）"
        content = content.replace(b'50\xef\xbf\xbd', b'50\xef\xbc\x89')
        # 修复 "指" 为 "指令"
        content = content.replace(b'\xe6\x8c\x87\xef\xbf\xbd', b'\xe6\x8c\x87\xe4\xbb\xa4')
        # 修复 "排" 为 "排在"
        content = content.replace(b'\xe6\x8e\x92\xef\xbf\xbd', b'\xe6\x8e\x92\xe5\x9c\xa8')
        
        # 修复注释中缺少换行的问题
        # 查找 "处\t\treturn" 模式，应该是 "处理\n\t\treturn"
        content = re.sub(rb'(\xe5\xa4\x84\xe7\x90\x86)\t\treturn', rb'\1\n\t\treturn', content)
        # 查找 ")\t\treturn" 模式在注释中，应该是 ")\n\t\treturn"
        content = re.sub(rb'\)\t\treturn', rb')\n\t\treturn', content)
        
        if content != original_content:
            with open(filepath, 'wb') as f:
                f.write(content)
            print(f"  已保存更改 ({len(content)} 字节)")
        else:
            print(f"  无需更改")

if __name__ == '__main__':
    fix_encoding_issues()
