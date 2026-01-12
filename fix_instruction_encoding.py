#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

def fix_instruction_encoding():
    filepath = 'server/internal/test/runner/instruction.go'
    
    print(f"处理文件: {filepath}")
    
    with open(filepath, 'rb') as f:
        content = f.read()
    
    original_content = content
    changes = 0
    
    # 修复 strings.Split 中字符串未正确关闭的问题
    # 查找 patterns like: strings.Split(instruction, "队伍攻击)\r\r\n
    # 应该是: strings.Split(instruction, "队伍攻击")
    
    # 修复 "队伍攻击)" 模式 - 查找并修复未关闭的字符串
    pattern1 = rb'strings\.Split\(instruction,\s*"(\xe9\x98\x9f\xe4\xbc\x8d\xe6\x94\xbb\xe5\x87\xbb)\)\r\r\n'
    def replace1(m):
        return b'strings.Split(instruction, "\xe9\x98\x9f\xe4\xbc\x8d\xe6\x94\xbb\xe5\x87\xbb")\n'
    
    content = re.sub(pattern1, replace1, content)
    if content != original_content:
        changes += 1
        print(f"  修复了 '队伍攻击' 字符串问题")
        original_content = content
    
    # 修复 "队伍生命)" 模式
    pattern2 = rb'strings\.Split\(instruction,\s*"(\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd)\)\r\r\n'
    def replace2(m):
        return b'strings.Split(instruction, "\xe9\x98\x9f\xe4\xbc\x8d\xe7\x94\x9f\xe5\x91\xbd")\n'
    
    content = re.sub(pattern2, replace2, content)
    if content != original_content:
        changes += 1
        print(f"  修复了 '队伍生命' 字符串问题")
        original_content = content
    
    # 修复 strings.Contains 中字符串未正确关闭的问题
    # 查找 patterns like: strings.Contains(instruction, "计算减伤后伤) {
    # 应该是: strings.Contains(instruction, "计算减伤后伤害") {
    
    # 修复 "计算减伤后伤)" 模式
    pattern3 = rb'strings\.Contains\(instruction,\s*"(\xe8\xae\xa1\xe7\xae\x97\xe5\x87\x8f\xe4\xbc\xa4\xe5\x90\x8e\xe4\xbc\xa4)\)\s*{'
    def replace3(m):
        return b'strings.Contains(instruction, "\xe8\xae\xa1\xe7\xae\x97\xe5\x87\x8f\xe4\xbc\xa4\xe5\x90\x8e\xe4\xbc\xa4\xe5\xae\xb3")\n\t\t{'
    
    content = re.sub(pattern3, replace3, content)
    if content != original_content:
        changes += 1
        print(f"  修复了 '计算减伤后伤害' 字符串问题")
        original_content = content
    
    # 更通用的修复：查找所有 strings.Split 或 strings.Contains 中未正确关闭的字符串
    # 查找 patterns like: "text)\r\r\n 或 "text) {
    # 应该是: "text")\n 或 "text") {
    
    # 修复所有以 )\r\r\n 结尾的字符串（在 strings.Split 或 strings.Contains 中）
    pattern4 = rb'("[\x20-\x7e\x80-\xff]+)\)\r\r\n'
    def replace4(m):
        text = m.group(1)
        # 检查是否在 strings.Split 或 strings.Contains 调用中
        return text + b'")'
    
    # 更精确的修复：只修复在函数调用中的情况
    # 查找: strings.Split(..., "text)\r\r\n 或 strings.Contains(..., "text) {
    pattern5 = rb'(strings\.(Split|Contains)\([^,]+,\s*")([^"]+)\)(\r\r\n|\s*{)'
    def replace5(m):
        prefix = m.group(1)
        text = m.group(3)
        suffix = m.group(4)
        if suffix == b'\r\r\n':
            return prefix + text + b'")\n'
        else:
            return prefix + text + b'")' + suffix
    
    # 先尝试更精确的修复
    content_new = re.sub(pattern5, replace5, content)
    if content_new != content:
        changes += 1
        print(f"  修复了函数调用中未关闭的字符串")
        content = content_new
    
    if content != original_content:
        with open(filepath, 'wb') as f:
            f.write(content)
        print(f"  已保存更改 ({len(content)} 字节)")
    else:
        print(f"  无需更改")

if __name__ == '__main__':
    fix_instruction_encoding()
