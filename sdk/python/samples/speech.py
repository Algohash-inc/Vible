#!/usr/bin/env python3
"""
TTS：text → out.wav。
voice："alloy"/留空 → 預設聲；"clone:名" → 本機 clone；speaker 名 → qwen_voice。
想直接喺裝置講出（唔回 bytes）：speak=True。
  python3 samples/speech.py
"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

v = Vible()
wav = v.speech("你好，我係你個 agent 嘅聲音。", voice="alloy")
with open("out.wav", "wb") as f:
    f.write(wav)
print(f"已儲存 out.wav（{len(wav)} bytes）。播放：afplay out.wav")
