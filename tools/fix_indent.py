#!/usr/bin/env python3
import os
import sys

SPACES = 4


def is_binary(file_path):
    try:
        with open(file_path, 'rb') as _f:
            _chunk = _f.read(8192)
        return b'\x00' in _chunk
    except OSError:
        return True


def convert_tabs(file_path):
    _indent = ' ' * SPACES
    try:
        with open(file_path, 'r', encoding='utf-8') as _f:
            _lines = _f.readlines()
    except UnicodeDecodeError:
        print(f"  [skip] cannot decode: {file_path}")
        return False

    _converted = []
    for _line in _lines:
        _leading = len(_line) - len(_line.lstrip('\t'))
        _converted.append(_indent * _leading + _line[_leading:])

    with open(file_path, 'w', encoding='utf-8') as _f:
        _f.writelines(_converted)
    return True


def ask_extension(ext):
    while True:
        _answer = input(f"  Convert *.{ext} files? [y/n]: ").strip().lower()
        if _answer in ('y', 'n'):
            return _answer == 'y'


def process_file(file_path, ext_cache):
    _ext = os.path.splitext(file_path)[1]

    if not _ext:
        print(f"  [skip/no-ext]  {file_path}")
        return False

    _ext_clean = _ext.lstrip('.')

    if is_binary(file_path):
        print(f"  [skip/binary]  {file_path}")
        return False

    if _ext_clean not in ext_cache:
        print(f"\nNew extension found: .{_ext_clean}")
        ext_cache[_ext_clean] = ask_extension(_ext_clean)

    if not ext_cache[_ext_clean]:
        return False

    print(f"  [convert]      {file_path}")
    return convert_tabs(file_path)


def walk_and_convert(root, ext_cache):
    _converted = 0
    _skipped = 0

    for _dirpath, _dirnames, _filenames in os.walk(root):
        for _fname in _filenames:
            _fpath = os.path.join(_dirpath, _fname)
            if process_file(_fpath, ext_cache):
                _converted += 1
            else:
                _skipped += 1

    return _converted, _skipped


def main():
    if len(sys.argv) < 2:
        print("Usage: python convert_indent.py <file_or_dir> [file_or_dir ...]")
        sys.exit(1)

    _ext_cache = {}
    _total_converted = 0
    _total_skipped = 0

    for _target in sys.argv[1:]:
        if os.path.isfile(_target):
            if process_file(_target, _ext_cache):
                _total_converted += 1
            else:
                _total_skipped += 1
        elif os.path.isdir(_target):
            _c, _s = walk_and_convert(_target, _ext_cache)
            _total_converted += _c
            _total_skipped += _s
        else:
            print(f"  [error] not found: '{_target}'")

    print(f"\nDone: {_total_converted} file(s) converted, {_total_skipped} skipped.")


if __name__ == '__main__':
    main()