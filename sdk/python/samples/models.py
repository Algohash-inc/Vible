#!/usr/bin/env python3
"""列出工作區可用 model。  python3 samples/models.py"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

for m in Vible().models():
    print(m)
