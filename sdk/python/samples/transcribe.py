#!/usr/bin/env python3
"""
STT：音訊 → 文字。
  python3 samples/speech.py       # 先整 out.wav
  python3 samples/transcribe.py    # 或 transcribe.py 你自己.wav
"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible

path = sys.argv[1] if len(sys.argv) > 1 else "out.wav"
try:
    audio = open(path, "rb").read()
except OSError as e:
    raise SystemExit(f"讀唔到 {path}：{e}（先跑 samples/speech.py）")

print(Vible().transcribe(audio, language="yue"))
