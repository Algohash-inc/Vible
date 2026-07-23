#!/usr/bin/env python3
"""Chat（SSE 串流）—— token 一粒粒出，邊收邊印。  python3 samples/stream.py"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

v = Vible()
v.chat_stream("數 1 到 5，每個數字之間停一下。",
              on_delta=lambda d: (sys.stdout.write(d), sys.stdout.flush()))
print()
