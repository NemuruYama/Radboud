#!/usr/bin/env python3
import sys

if __name__ == '__main__':
    data = sys.stdin.readlines()
    elements = list(map(int, data[1].split(' ')))
    print(sum(elements), file=sys.stdout)
