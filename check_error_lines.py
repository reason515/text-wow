#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def check_error_lines():
    for filepath, target_line in [
        ('server/internal/test/runner/context.go', 270),
        ('server/internal/test/runner/equipment.go', 164)
    ]:
        print(f"\n检查文件: {filepath}, 目标行: {target_line}")
        with open(filepath, 'rb') as f:
            lines = f.readlines()
        
        for i in range(max(0, target_line-5), min(len(lines), target_line+3)):
            print(f"Line {i+1}: {repr(lines[i][:150])}")

if __name__ == '__main__':
    check_error_lines()
