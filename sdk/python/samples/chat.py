#!/usr/bin/env python3
"""Chat（非串流）—— 一次過攞齊回覆。  python3 samples/chat.py"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

v = Vible()
print(v.chat("用一句話自我介紹。"))
