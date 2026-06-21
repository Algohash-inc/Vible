#!/usr/bin/env python3
"""
雲端來電：外部 agent 主動打電話畀用戶，裝置即時響鈴。
call key 由 env VIBLE_CALL_KEY 或本機 operator 自動鑄。要裝置連住 hub 至會真響。
  python3 samples/call.py
"""
import os, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # 令 `vible` 可 import（免安裝）
from vible import Vible, VibleError

try:
    r = Vible().call("提你 3 點開個會。", title="會議提醒")
    print(r)   # {call_id, ringing, delivered}
except VibleError as e:
    print("call 失敗：", e)
