#!/usr/bin/env python3
# -*- coding: utf-8 -*-

file_path = 'server/internal/test/runner/test_runner.go'

with open(file_path, 'r', encoding='utf-8', errors='replace') as f:
    lines = f.readlines()

print(f'File total lines: {len(lines)}')

# Check brace matching between lines 1550-1580
brace_count = 0
for i in range(1550, 1580):
    line = lines[i]
    brace_count += line.count('{') - line.count('}')
    if line.strip():
        print(f'Line {i+1:5d}: brace={brace_count:2d}')

print(f'\nBrace count at line 1580: {brace_count}')

# If brace_count < 0, need to add closing brace before line 1579
if brace_count < 0:
    # Find suitable position to insert }
    # Insert after line 1573
    if len(lines) > 1573:
        lines.insert(1573, '\t}\n')
        print('Added closing brace after line 1573')
        
        # Write back to file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.writelines(lines)
        print('Fix completed!')
    else:
        print('Cannot fix: line number out of range')
else:
    print('No unclosed braces found')
