#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_remaining_encoding():
    files_to_fix = [
        ('server/internal/test/runner/instruction.go', [
            # 修复未关闭的字符串 - 使用十六进制编码
            (rb'\xe8\xae\xa1\xe7\xae\x97\xe5\x87\x8f\xe4\xbc\xa4\xe5\x90\x8e\xe4\xbc\xa4\)\s*{', b'\xe8\xae\xa1\xe7\xae\x97\xe5\x87\x8f\xe4\xbc\xa4\xe5\x90\x8e\xe4\xbc\xa4\xe5\xae\xb3")\n\t\t{'),
            (rb'\xe5\xad\xa6\xe4\xb9\xa0\xe6\x8a\x80\)\s*\|\|', b'\xe5\xad\xa6\xe4\xb9\xa0\xe6\x8a\x80\xe8\x83\xbd")\n\t\t||'),
            (rb'\xe8\xa7\x92\xe8\x89\xb2\xe5\xad\xa6\xe4\xb9\xa0\xe6\x8a\x80\)\s*{', b'\xe8\xa7\x92\xe8\x89\xb2\xe5\xad\xa6\xe4\xb9\xa0\xe6\x8a\x80\xe8\x83\xbd")\n\t\t{'),
            (rb'\xe6\x8a\x80\)\s*{', b'\xe6\x8a\x80\xe8\x83\xbd")\n\t\t{'),
            (rb'\xe4\xbd\xbf\xe7\x94\xa8\xe6\x8a\x80\)\s*\|\|', b'\xe4\xbd\xbf\xe7\x94\xa8\xe6\x8a\x80\xe8\x83\xbd")\n\t\t||'),
            (rb'\xe8\xa7\x92\xe8\x89\xb2\xe4\xbd\xbf\xe7\x94\xa8\xe6\x8a\x80\)\s*\|\|', b'\xe8\xa7\x92\xe8\x89\xb2\xe4\xbd\xbf\xe7\x94\xa8\xe6\x8a\x80\xe8\x83\xbd")\n\t\t||'),
            (rb'\xe6\x8a\x80\)\)\s*{', b'\xe6\x8a\x80\xe8\x83\xbd"))\n\t\t{'),
            (rb'\xe5\x88\x9b\xe5\xbb\xba\xe4\xb8\x80\)\s*&&', b'\xe5\x88\x9b\xe5\xbb\xba\xe4\xb8\x80\xe4\xb8\xaa")\n\t\t&&'),
            (rb'\xe6\x89\xa7\xe8\xa1\x8c\)\s*&&', b'\xe6\x89\xa7\xe8\xa1\x8c")\n\t\t&&'),
        ]),
    ]
    
    for filepath, patterns in files_to_fix:
        print(f"\n处理文件: {filepath}")
        
        with open(filepath, 'rb') as f:
            content = f.read()
        
        original_content = content
        changes = 0
        
        for pattern, replacement in patterns:
            new_content = re.sub(pattern, replacement, content)
            if new_content != content:
                content = new_content
                changes += 1
                print(f"  修复了模式: {pattern[:50]}...")
        
        # 修复 instruction.go 中所有未关闭的字符串
        # 查找 patterns like: "text) { 或 "text) ||
        # 应该是: "text") { 或 "text") ||
        pattern_general = rb'("[\x20-\x7e\x80-\xff]+)\)(\s*\{|\s*\|\|)'
        def replace_general(m):
            text = m.group(1)
            suffix = m.group(2)
            if suffix.strip() == b'{':
                return text + b'")' + b'\n\t\t' + suffix
            else:
                return text + b'")' + suffix
        
        new_content = re.sub(pattern_general, replace_general, content)
        if new_content != content:
            changes += 1
            print(f"  修复了通用未关闭字符串模式")
            content = new_content
        
        if content != original_content:
            with open(filepath, 'wb') as f:
                f.write(content)
            print(f"  已保存更改 ({len(content)} 字节)")
        else:
            print(f"  无需更改")
    
    # 检查 context.go 和 equipment.go 的函数结构
    for filepath in ['server/internal/test/runner/context.go', 'server/internal/test/runner/equipment.go']:
        print(f"\n检查文件: {filepath}")
        with open(filepath, 'rb') as f:
            content = f.read()
        
        # 检查是否有未关闭的函数
        # 查找函数定义，确保它们都有正确的闭合
        lines = content.split(b'\n')
        brace_count = 0
        in_function = False
        for i, line in enumerate(lines):
            if b'func ' in line and b'{' in line:
                in_function = True
                brace_count = line.count(b'{') - line.count(b'}')
            elif in_function:
                brace_count += line.count(b'{') - line.count(b'}')
                if brace_count == 0:
                    in_function = False
        
        if brace_count != 0:
            print(f"  警告: 可能有不匹配的大括号 (计数: {brace_count})")
        else:
            print(f"  大括号匹配正常")

if __name__ == '__main__':
    fix_remaining_encoding()
