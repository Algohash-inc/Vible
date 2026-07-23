#!/usr/bin/env python3
"""
Realtime WebSocket（全雙工）—— 需 `pip install websockets`。
  python3 samples/realtime.py
"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

v = Vible()
v.realtime("講三樣關於海洋嘅有趣事。",
           on_delta=lambda d: (sys.stdout.write(d), sys.stdout.flush()))
print("\n[完成]")
