#!/usr/bin/env python3
"""
語音回合：你講 → STT → agent → 喺裝置講出回覆。
  python3 samples/speech.py       # 先整 out.wav（或用你自己嘅語音）
  python3 samples/converse.py
"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

path = sys.argv[1] if len(sys.argv) > 1 else "out.wav"
try:
    audio = open(path, "rb").read()
except OSError as e:
    raise SystemExit(f"讀唔到 {path}：{e}（先跑 samples/speech.py）")

d = Vible().converse(audio)
print("你講咗：", d.get("transcript", ""))
print("agent 回覆緊（喺裝置出聲）…" if d.get("spoken") else "（靜音 / 冇嘢講）")
