<div align="center">

[English](README.md) · [繁體中文](README.zh-Hant.md) · [简体中文](README.zh-Hans.md) · [日本語](README.ja.md) · **한국어** · [Español](README.es.md) · [Français](README.fr.md) · [Deutsch](README.de.md)

</div>

<div align="center">

<img src="banner.png" alt="Vible — Speak · Polish · Translate" width="100%" />

# Vible — 당신의 목소리를 완벽한 텍스트로

**말하고. 다듬고. 번역하세요.** 온디바이스로, 즉시, 그리고 프라이빗하게 — Mac과 Windows PC를 위한 완전한 음성 레이어.

<br/>

[![⬇ Download for macOS](https://img.shields.io/badge/⬇_Download-macOS_(.dmg)-C7CED9?style=for-the-badge&labelColor=0A0B0E)](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg)
[![⬇ Download for Windows](https://img.shields.io/badge/⬇_Download-Windows_(.exe)-C7CED9?style=for-the-badge&labelColor=0A0B0E)](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible-Setup.exe)

![macOS 12+](https://img.shields.io/badge/macOS-12%2B-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Windows 10/11](https://img.shields.io/badge/Windows-10%2F11-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Apple Silicon](https://img.shields.io/badge/Apple_Silicon-native-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Signed & Notarized](https://img.shields.io/badge/Signed_%26_Notarized-✓-34d399?style=for-the-badge&labelColor=0A0B0E)

[**웹사이트**](https://vible.space) · [**다운로드**](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) · [**설치 가이드**](#-설치) · [**자주 묻는 질문**](#-자주-묻는-질문)

</div>

Vible은 당신의 목소리를 깔끔하게 교정되고 번역된 텍스트로 바꿔줍니다 — **어떤 앱에서든**, 생각의 속도로. 로컬 모델은 Apple Silicon에서 속도와 프라이버시를 제공하며, 원할 때는 스위치 하나로 클라우드 GPU의 성능을 켤 수 있습니다. 단축키를 누르고, 말하고, 놓으세요 — Vible이 받아쓰고, 문법을 고치고, 번역하고, 심지어 소리 내어 읽어주기까지 합니다.

---

## ✨ 주요 기능

- ⚡ **생각의 속도로** — `⌃⌘` 를 누르고, 말하고, 놓으세요. 로컬 Whisper가 즉시 받아쓰고, 다듬어진 텍스트가 지금 사용 중인 앱에 바로 붙여넣어집니다.
- 🧹 **자동 문법 다듬기** — 모든 받아쓰기 결과는 화면에 나타나기 전에 정리됩니다: 문장 부호, 대소문자, 그리고 군더더기 말("음", "어", "그니까")이 사라집니다.
- 🌍 **실시간 번역, 20개 이상 언어** — 한 언어로 말하고, 다른 언어로 입력하세요. 실시간 양방향 대화를 위한 2-way 대화 모드까지.
- ⌘ **`⌘C` 두 번으로 어디서나 번역** — 어떤 앱에서든 텍스트를 선택하고 `⌘C` 를 두 번 누르면, 팝업이 그 자리에서 바로 번역해 줍니다 — 복사하거나, 소리 내어 읽거나, 선택 영역을 대체하세요.
- 🗣️ **자연스러운 텍스트 음성 변환 + 보이스 클론** — 무엇이든 스튜디오급 음성으로 읽어주거나, 30초를 녹음해 *당신 자신의* 목소리로 Vible이 읽어주게 하세요.
- 🪶 **기본값은 온디바이스** — 로컬 **MLX + ONNX** 모델이 Apple Silicon에서 실행됩니다. 빠르고, 오프라인에서도 작동하며, 프라이빗합니다 — 당신의 오디오는 Mac을 떠날 필요가 없습니다.
- ☁️ **원할 때는 클라우드** — 스위치 하나로 GPU 기반 모델을 켜세요: 다운로드도, 예열도 없이, 더 풍부한 뉴럴 음성을.
- 🔒 **프라이빗하고 암호화** — 클라우드를 켜지 않는 한 아무것도 업로드되지 않습니다. 온디바이스 및 전송 중 데이터는 암호화되며, 당신의 음성과 텍스트는 절대 판매·공유되거나 학습에 사용되지 않습니다.
- 🍎 **서명 및 공증 완료** — Apple Silicon용으로 제작되어 **Algohash Inc.** 가 서명하고 Apple이 공증했으므로 — 보안 경고 없이 열립니다.

---

## 🎙️ 음성 루프 — 단축키 하나, 네 가지 초능력

키를 누르고, 말하고, 놓으세요. 이것이 전체 동작입니다.

| | 단계 | 무슨 일이 일어나나요 |
|---|---|---|
| 1 | **캡처** | 어디서든 `⌃⌘` 를 누른 채 말하세요 — 실시간 파형이 듣고 있음을 보여줍니다. |
| 2 | **받아쓰기** | 로컬 Whisper가 음성을 즉시, 온디바이스에서 원본 텍스트로 바꿉니다. |
| 3 | **다듬기** | 문법, 대소문자, 군더더기 말이 한 번에 정리됩니다. |
| 4 | **번역 / 읽기** | 선택적으로 번역한 뒤 앱에 자동으로 붙여넣거나 — 소리 내어 읽어줍니다. |

> 당신은 타이핑보다 약 3배 빠르게 말합니다. Vible은 그 속도를 Mac의 모든 텍스트 필드에서 쓸 수 있게 만들어 줍니다.

---

## ⌘ 무엇이든, 지금 있는 그 자리에서 번역

열 창도, 붙여넣을 것도 없습니다:

1. 어떤 앱에서든 텍스트를 **선택**하세요 — Mail, Slack, PDF, 브라우저.
2. **`⌘C` 를 두 번 빠르게 누르세요** — Vible이 그 복사를 포착해, 커서 바로 옆에 팝업이 뜹니다.
3. **언어와 어조를 고르세요** — 원본 언어는 자동으로 감지됩니다. 대상 언어와 어조(*격식*, *자연스럽게*, 또는 *발음 포함*)를 선택한 뒤 **복사**하거나, 소리 내어 **읽거나**, 선택 영역을 그 자리에서 **대체**하세요.

`20개 이상 언어` · `격식 / 자연스럽게` · `그 자리에서 대체`

---

## 🗣️ 목소리를 내는 세 가지 방법

| 음성 | 등급 | 설명 |
|---|---|---|
| **Standard** | 온디바이스 · 무료 | Apple Silicon에서 오프라인으로 실행되는 선명한 로컬 TTS — 즉각적이고 프라이빗한 읽기. |
| **Premium** | 클라우드 GPU | 20개 이상 언어에 걸쳐 감정을 담은 스튜디오급 뉴럴 음성을 클라우드에서 스트리밍. |
| **Zero-shot 보이스 클론** | Studio | 30초를 녹음하면 Vible이 *당신 자신의* 목소리로 무엇이든 읽어줍니다. |

---

## 🔒 기본값이 프라이버시

Vible은 **완전히 로컬로** 실행됩니다 — 음성-텍스트 변환과 텍스트-음성 변환이 **MLX & ONNX** 를 통해 Apple Silicon에서 전적으로 처리됩니다. 인터넷도, 서버도, 새어 나갈 것도 없습니다 — Wi-Fi를 꺼도 프라이빗합니다.

- ✅ 완전히 오프라인에서 작동
- ✅ 아무것도 업로드하지 않음 — 절대로 (클라우드를 켜지 않는 한)
- ✅ 온디바이스 **및** 전송 중 암호화(TLS)

최대 처리량과 가장 풍부한 뉴럴 음성이 필요하신가요? 스위치 하나로 **GPU 클라우드** 를 켜세요 — 예열 없는 즉각적인 합성. 어느 쪽이든 데이터는 암호화되며 절대 학습에 사용되지 않습니다.

---

## 📊 타이핑은 느립니다. Vible은 아닙니다.

| 기능 | 타이핑 | 일반 받아쓰기 | **Vible** |
|---|:---:|:---:|:---:|
| 원본 입력 속도 | ~40 wpm | ~120 wpm | **~150 wpm** |
| 자동 문법 & 문장 부호 | 수동 | — | **✓** |
| 인라인 번역 | — | — | **✓ 20개 이상 언어** |
| 모든 앱에서 작동 | ✓ | 일부 | **✓ 시스템 전체** |
| 오프라인 / 온디바이스 실행 | ✓ | 클라우드 전용 | **✓ + 클라우드** |
| 결과를 소리 내어 읽기(TTS) | — | — | **✓ + 보이스 클론** |

---

## 🌍 언어

**인터페이스(8):** English · 繁體中文 · 简体中文 · 日本語 · 한국어 · Español · Français · Deutsch

**번역 & 음성:** 20개 이상 언어, 2-way 대화 모드 포함.

---

## 🚀 설치

> [**⬇ Vible.dmg 다운로드**](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) — 그런 다음 아래에서 당신의 언어로 된 단계를 따라 하세요.
>
> **Windows 사용자이신가요?** 대신 [`Vible-Setup.exe`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible-Setup.exe) 를 다운로드하고 설치 프로그램을 실행한 뒤 Vible 을 실행하세요. 마이크를 허용하고 **`Ctrl+Alt`** 를 누른 채 말하세요(macOS 의 `⌃⌘` 에 해당). 번역은 **`Ctrl+C`** 더블 탭.

<details open>
<summary><b>English</b></summary>

1. **Download** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg).
2. **Open** the downloaded `Vible.dmg` (double-click it).
3. **Drag** the **Vible** icon onto the **Applications** folder in the window that appears.
4. **Launch** Vible from Applications or Launchpad. It's signed and notarized by Apple, so it opens with **no security warning**.
5. **Grant the mic** — allow microphone + input monitoring once. Nothing leaves your Mac unless you turn cloud on.
6. **Hold `⌃⌘` and talk** — release, and the polished text pastes right where your cursor is. On first launch, Vible sets up its on-device models — give it a moment.

</details>

<details>
<summary><b>繁體中文</b></summary>

1. **下載** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg)。
2. **打開** 下載好的 `Vible.dmg`（雙擊它）。
3. 在彈出的視窗中，把 **Vible** 圖示**拖曳**到 **Applications（應用程式）**資料夾。
4. 從「應用程式」或 Launchpad **啟動** Vible。它已由 Apple 簽署並公證，因此開啟時**不會有任何安全警告**。
5. **授予麥克風權限** — 允許「麥克風」與「輸入監控」一次。除非你開啟雲端，否則資料不會離開你的 Mac。
6. **按住 `⌃⌘` 說話** — 放開後，整理好的文字會貼到游標所在處。首次啟動時，Vible 會設定本機模型，請稍候片刻。

</details>

<details>
<summary><b>简体中文</b></summary>

1. **下载** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg)。
2. **打开** 下载好的 `Vible.dmg`（双击它）。
3. 在弹出的窗口中，把 **Vible** 图标**拖到** **Applications（应用程序）**文件夹。
4. 从「应用程序」或 Launchpad **启动** Vible。它已由 Apple 签名并公证，因此打开时**不会出现任何安全警告**。
5. **授予麦克风权限** — 允许「麦克风」与「输入监控」一次。除非你开启云端，否则数据不会离开你的 Mac。
6. **按住 `⌃⌘` 说话** — 松开后，整理好的文字会粘贴到光标所在处。首次启动时，Vible 会设置本地模型，请稍候片刻。

</details>

<details>
<summary><b>日本語</b></summary>

1. [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) を **ダウンロード** します。
2. ダウンロードした `Vible.dmg` を **開きます**（ダブルクリック）。
3. 表示されたウインドウで、**Vible** アイコンを **「アプリケーション」フォルダ** に **ドラッグ** します。
4. 「アプリケーション」または Launchpad から Vible を **起動** します。Apple の署名・公証済みなので、**警告なしで** 開きます。
5. **マイクを許可** — 「マイク」と「入力監視」を一度だけ許可します。クラウドを有効にしない限り、データが Mac の外に出ることはありません。
6. **`⌃⌘` を押しながら話す** — 離すと、整えられたテキストがカーソル位置に貼り付けられます。初回起動時、Vible は端末内モデルを設定します。少しお待ちください。

</details>

<details>
<summary><b>한국어</b></summary>

1. [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) 를 **다운로드**합니다.
2. 다운로드한 `Vible.dmg` 를 **엽니다**(더블 클릭).
3. 나타나는 창에서 **Vible** 아이콘을 **응용 프로그램(Applications)** 폴더로 **드래그**합니다.
4. 응용 프로그램 또는 Launchpad에서 Vible 을 **실행**합니다. Apple의 서명·공증을 받았으므로 **보안 경고 없이** 열립니다.
5. **마이크 권한 허용** — '마이크'와 '입력 모니터링'을 한 번 허용합니다. 클라우드를 켜지 않는 한 데이터가 Mac을 떠나지 않습니다.
6. **`⌃⌘` 를 누른 채 말하기** — 손을 떼면 다듬어진 텍스트가 커서 위치에 붙여넣어집니다. 처음 실행 시 Vible 은 온디바이스 모델을 설정합니다. 잠시 기다려 주세요.

</details>

<details>
<summary><b>Español</b></summary>

1. **Descarga** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg).
2. **Abre** el `Vible.dmg` descargado (haz doble clic).
3. **Arrastra** el icono de **Vible** a la carpeta **Aplicaciones** en la ventana que aparece.
4. **Inicia** Vible desde Aplicaciones o Launchpad. Está firmado y notarizado por Apple, así que se abre **sin ninguna advertencia de seguridad**.
5. **Concede el micrófono** — permite el micrófono y la monitorización de entrada una vez. Nada sale de tu Mac salvo que actives la nube.
6. **Mantén `⌃⌘` y habla** — al soltar, el texto pulido se pega justo donde está tu cursor. En el primer inicio, Vible configura sus modelos en el dispositivo; espera un momento.

</details>

<details>
<summary><b>Français</b></summary>

1. **Téléchargez** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg).
2. **Ouvrez** le `Vible.dmg` téléchargé (double-cliquez).
3. **Faites glisser** l'icône **Vible** vers le dossier **Applications** dans la fenêtre qui s'affiche.
4. **Lancez** Vible depuis Applications ou le Launchpad. Il est signé et notarisé par Apple, il s'ouvre donc **sans aucun avertissement de sécurité**.
5. **Autorisez le micro** — autorisez le microphone et la surveillance des entrées une fois. Rien ne quitte votre Mac sauf si vous activez le cloud.
6. **Maintenez `⌃⌘` et parlez** — au relâchement, le texte nettoyé se colle là où se trouve votre curseur. Au premier lancement, Vible configure ses modèles sur l'appareil ; patientez un instant.

</details>

<details>
<summary><b>Deutsch</b></summary>

1. **Lade** [`Vible.dmg`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) herunter.
2. **Öffne** die heruntergeladene `Vible.dmg` (Doppelklick).
3. **Ziehe** das **Vible**-Symbol im erscheinenden Fenster auf den Ordner **Programme (Applications)**.
4. **Starte** Vible über Programme oder das Launchpad. Es ist von Apple signiert und notarisiert und öffnet sich **ohne Sicherheitswarnung**.
5. **Mikrofon erlauben** — erlaube Mikrofon + Eingabeüberwachung einmalig. Nichts verlässt deinen Mac, außer du aktivierst die Cloud.
6. **Halte `⌃⌘` und sprich** — beim Loslassen wird der bereinigte Text direkt an der Cursorposition eingefügt. Beim ersten Start richtet Vible seine geräteinternen Modelle ein – einen Moment Geduld.

</details>

> macOS가 앱을 "열 수 없습니다"라고 표시하면, 응용 프로그램에서 Vible을 마우스 오른쪽 버튼으로 클릭 → **열기** → **열기** 를 선택하세요. 이는 다운로드가 중단된 경우에만 발생하며 — 다시 다운로드하면 해결됩니다.

---

## 🖥️ 요구 사항

- **Apple Silicon**(M1 이상) 탑재 Mac
- **macOS 12 Monterey** 이상
- 약 600 MB 다운로드 · 몇 초 만에 설치 · **Windows 10/11 지원** · Linux 곧 지원 예정

---

## 💳 요금제

| | Free | Pro | Studio |
|---|---|---|---|
| 온디바이스 음성 입력 | ✓ | ✓ | ✓ |
| 문법 다듬기 | ✓ | ✓ | ✓ |
| 번역 | 기본 | 20개 이상 언어 + 대화 | 20개 이상 언어 + 대화 |
| 음성 | Standard | + Premium / 감정 | + 보이스 클론 |
| 클라우드 GPU 속도 | — | ✓ | ✓ |
| 에이전트 모드 & 도구 동작 | — | — | ✓ |

온디바이스 기능은 무료로, 계속 당신의 것입니다. 클라우드 GPU 속도와 무제한 사용을 원할 때만 결제하세요. 현재 가격은 [vible.space](https://vible.space) 에서 확인하세요.

---

## ❓ 자주 묻는 질문

<details>
<summary><b>정말 온디바이스에서 실행되나요?</b></summary>

네. 기본 음성-텍스트 변환 및 TTS 모델은 MLX와 ONNX를 통해 Apple Silicon에서 로컬로 실행됩니다 — 빠르고 오프라인에서도 작동합니다. 클라우드 GPU는 선택 사항으로, 더 풍부한 음성이나 최대 속도를 원할 때 스위치 하나로 켜면 됩니다.
</details>

<details>
<summary><b>어떤 Mac이 지원되나요?</b></summary>

macOS 12 이상을 실행하는 모든 Apple Silicon Mac(M1 이상)입니다. 앱은 서명 및 공증되어 있습니다. 서명 및 공증된 Windows 10/11 빌드를 지금 사용할 수 있습니다. Linux 빌드는 로드맵에 있습니다.
</details>

<details>
<summary><b>Vible은 모든 앱에서 작동하나요?</b></summary>

네 — 시스템 전체에서 작동합니다. 어떤 텍스트 필드(Mail, Slack, 에디터, 브라우저)에서든 단축키를 누르면 다듬어진 결과가 커서 위치에 바로 붙여넣어집니다.
</details>

<details>
<summary><b>제 목소리로 말할 수 있나요?</b></summary>

Studio 요금제에서는 짧은 샘플을 녹음하면 Vible이 표준 및 프리미엄 뉴럴 음성과 함께 당신의 목소리를 복제해 텍스트를 읽어줍니다.
</details>

---

<div align="center">

**[vible.space](https://vible.space)**  ·  Speak. Polish. Translate.

Copyright © 2026 **Algohash Inc.** · Vible is a signed, notarized macOS application. All rights reserved.

</div>
