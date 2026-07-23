<div align="center">

[English](README.md) · [繁體中文](README.zh-Hant.md) · [简体中文](README.zh-Hans.md) · **日本語** · [한국어](README.ko.md) · [Español](README.es.md) · [Français](README.fr.md) · [Deutsch](README.de.md)

</div>

<div align="center">

<img src="banner.png" alt="Vible — Speak · Polish · Translate" width="100%" />

# Vible — あなたの声を、完璧なテキストへ磨き上げる

**話す。磨く。翻訳する。** オンデバイスで、瞬時に、そしてプライベートに — Mac と Windows PC のための、まるごとの音声レイヤー。

<br/>

[![⬇ Download for macOS](https://img.shields.io/badge/⬇_Download-macOS_(.dmg)-C7CED9?style=for-the-badge&labelColor=0A0B0E)](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg)
[![⬇ Download for Windows](https://img.shields.io/badge/⬇_Download-Windows_(.exe)-C7CED9?style=for-the-badge&labelColor=0A0B0E)](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible-Setup.exe)

![macOS 12+](https://img.shields.io/badge/macOS-12%2B-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Windows 10/11](https://img.shields.io/badge/Windows-10%2F11-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Apple Silicon](https://img.shields.io/badge/Apple_Silicon-native-0A0B0E?style=for-the-badge&labelColor=0A0B0E)
![Signed & Notarized](https://img.shields.io/badge/Signed_%26_Notarized-✓-34d399?style=for-the-badge&labelColor=0A0B0E)

[**ウェブサイト**](https://vible.space) · [**ダウンロード**](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) · [**セットアップガイド**](#-インストール) · [**FAQ**](#-faq)

</div>

Vible は、あなたの声を、クリーンで、修正され、翻訳されたテキストへと変えます — **どんなアプリでも**、思考のスピードで。ローカルモデルが Apple Silicon 上でスピードとプライバシーをもたらし、必要なときはスイッチひとつでクラウド GPU のパワーへ。ホットキーを押しながら話して、離すだけ — Vible が文字起こしし、文法を整え、翻訳し、さらに読み上げることもできます。

---

## ✨ ハイライト

- ⚡ **思考のスピードで** — `⌃⌘` を押しながら話して、離す。ローカルの Whisper が瞬時に文字起こしし、磨き上げられたテキストが、いま使っているアプリへそのまま貼り付けられます。
- 🧹 **自動の文法磨き** — すべての文字起こしは、着地する前にクリーンに。句読点、大文字・小文字、そして口ぐせ（「えっと」「あの」「みたいな」）はすべて取り除かれます。
- 🌍 **リアルタイム翻訳、20 以上の言語** — ある言語で話し、別の言語で入力。双方向の会話モードで、リアルタイムのやり取りも。
- ⌘ **`⌘C` のダブルタップで、どこでも翻訳** — どんなアプリでもテキストを選択し、`⌘C` を 2 回タップすると、ポップアップがその場で翻訳します — コピー、読み上げ、または選択範囲の置き換えも。
- 🗣️ **自然なテキスト読み上げ + ボイスクローン** — スタジオ品質の声で何でも読み上げ。あるいは 30 秒録音すれば、Vible が *あなた自身の* 声で読み返します。
- 🪶 **標準でオンデバイス** — ローカルの **MLX + ONNX** モデルが Apple Silicon 上で動作。高速で、オフライン対応、そしてプライベート — あなたの音声が Mac の外に出る必要はありません。
- ☁️ **必要なときはクラウドで** — スイッチひとつで GPU バックのモデルが有効に。ダウンロードも、ウォームアップも不要で、より豊かなニューラルボイスを。
- 🔒 **プライベートで暗号化済み** — クラウドをオンにしない限り、何もアップロードされません。オンデバイスと通信中のデータは暗号化され、あなたの声とテキストが販売・共有・学習に使われることは決してありません。
- 🍎 **署名・公証済み** — Apple Silicon 向けに構築され、**Algohash Inc.** によって署名、Apple によって公証済み — セキュリティ警告なしで開きます。

---

## 🎙️ 音声のループ — ひとつのホットキー、4 つのスーパーパワー

キーを押しながら話して、離す。やり取りはこれだけ。

| | ステップ | 何が起きるか |
|---|---|---|
| 1 | **キャプチャ** | どこでも `⌃⌘` を押しながら話すと — ライブの波形が、聴いていることを示します。 |
| 2 | **文字起こし** | ローカルの Whisper が、音声を瞬時に生のテキストへ、オンデバイスで変換します。 |
| 3 | **磨き上げ** | 文法、大文字・小文字、口ぐせが、ひと処理でクリーンに。 |
| 4 | **翻訳 / 読み上げ** | 任意で翻訳し、アプリへ自動貼り付け — または声に出して読み上げます。 |

> 人はタイプするより約 3 倍速く話せます。Vible は、そのスピードを Mac のあらゆるテキストフィールドで使えるものにします。

---

## ⌘ どこにいても、その場で何でも翻訳

開くウインドウも、貼り付けも不要:

1. どんなアプリでも任意のテキストを **選択** — メール、Slack、PDF、ブラウザ。
2. **`⌘C` をダブルタップ** — すばやく。Vible がコピーを捉え、カーソルのすぐそばにポップアップが現れます。
3. **言語とトーンを選ぶ** — 元の言語は自動検出され、翻訳先の言語とトーン（*フォーマル*、*ナチュラル*、*発音つき*）を選んだら、**コピー**、**読み上げ**、または選択範囲へその場で **置き換え** ます。

`20 以上の言語` · `フォーマル / ナチュラル` · `その場で置き換え`

---

## 🗣️ 声で伝える 3 つの方法

| ボイス | ティア | 説明 |
|---|---|---|
| **スタンダード** | オンデバイス · 無料 | Apple Silicon 上でオフライン動作する、クリアなローカル TTS — 瞬時で、プライベートな読み上げ。 |
| **プレミアム** | クラウド GPU | 20 以上の言語にわたる、感情豊かなスタジオ品質のニューラルボイスを、クラウドからストリーミング。 |
| **ゼロショット・ボイスクローン** | スタジオ | 30 秒録音すれば、Vible が *あなた自身の* 声で何でも読み返します。 |

---

## 🔒 標準でプライベート

Vible は **すべてローカルで** 動作します — 音声からテキスト、テキストから音声への処理は、**MLX & ONNX** を通じて Apple Silicon 上で完結します。インターネットも、サーバーも、漏れるものも何もなく — Wi-Fi をオフにしていてもプライベートです。

- ✅ 完全にオフラインで動作
- ✅ 何もアップロードされません — 一切（クラウドを有効にしない限り）
- ✅ オンデバイス **かつ** 通信中も暗号化（TLS）

最大のスループットと、いちばん豊かなニューラルボイスが必要ですか？ スイッチひとつで **GPU クラウド** へ — ウォームアップなしの瞬時の合成。どちらの場合も、あなたのデータは暗号化され、学習に使われることは決してありません。

---

## 📊 タイピングは遅い。Vible は違う。

| 機能 | タイピング | 一般的な音声入力 | **Vible** |
|---|:---:|:---:|:---:|
| 生の入力スピード | 約 40 wpm | 約 120 wpm | **約 150 wpm** |
| 自動の文法・句読点 | 手動 | — | **✓** |
| インライン翻訳 | — | — | **✓ 20 以上の言語** |
| すべてのアプリで動作 | ✓ | 一部 | **✓ システム全体** |
| オフライン / オンデバイスで動作 | ✓ | クラウドのみ | **✓ + クラウド** |
| 結果を読み上げる（TTS） | — | — | **✓ + ボイスクローン** |

---

## 🌍 言語

**インターフェース（8）:** English · 繁體中文 · 简体中文 · 日本語 · 한국어 · Español · Français · Deutsch

**翻訳 & ボイス:** 20 以上の言語、双方向の会話モードつき。

---

## 🚀 インストール

> [**⬇ Vible.dmg をダウンロード**](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible.dmg) — その後、下記のあなたの言語の手順に従ってください。
>
> **Windows をお使いですか？** 代わりに [`Vible-Setup.exe`](https://github.com/Algohash-inc/Vible/releases/latest/download/Vible-Setup.exe) をダウンロードし、インストーラーを実行して Vible を起動します。マイクを許可し、**`Ctrl+Alt`** を押しながら話します（macOS の `⌃⌘` に相当）。翻訳は **`Ctrl+C`** をダブルタップ。

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

> macOS がアプリを「開けません」と表示した場合は、「アプリケーション」内の Vible を右クリック → **開く** → **開く**。これはダウンロードが中断された場合にのみ起こります — 再ダウンロードで解決します。

---

## 🖥️ 動作環境

- **Apple Silicon**（M1 以降）搭載の Mac
- **macOS 12 Monterey** 以降
- 約 600 MB のダウンロード · 数秒でインストール · **Windows 10/11 対応済み** · Linux は近日対応予定

---

## 💳 プラン

| | Free | Pro | Studio |
|---|---|---|---|
| オンデバイス音声入力 | ✓ | ✓ | ✓ |
| 文法磨き | ✓ | ✓ | ✓ |
| 翻訳 | 基本 | 20 以上の言語 + 会話 | 20 以上の言語 + 会話 |
| ボイス | スタンダード | + プレミアム / 感情表現 | + ボイスクローン |
| クラウド GPU スピード | — | ✓ | ✓ |
| エージェントモード & ツール操作 | — | — | ✓ |

オンデバイス機能は、無料でずっとあなたのもの。クラウド GPU のスピードと無制限の利用が欲しいときだけ、お支払いください。最新の価格は [vible.space](https://vible.space) をご覧ください。

---

## ❓ FAQ

<details>
<summary><b>本当にオンデバイスで動作するのですか？</b></summary>

はい。標準の音声からテキスト、および TTS モデルは、MLX と ONNX を通じて Apple Silicon 上でローカルに動作します — 高速で、オフライン対応です。クラウド GPU はオプトインで、より豊かなボイスや最高のスピードが欲しいときに、スイッチひとつで切り替えます。
</details>

<details>
<summary><b>どの Mac に対応していますか？</b></summary>

macOS 12 以降を搭載した、すべての Apple Silicon Mac（M1 以降）に対応します。アプリは署名・公証済みです。署名・公証済みの Windows 10/11 版を現在提供中です。Linux 版はロードマップにあります。
</details>

<details>
<summary><b>Vible はすべてのアプリで動作しますか？</b></summary>

はい — システム全体で動作します。どんなテキストフィールド（メール、Slack、エディタ、ブラウザ）でもホットキーを押せば、磨き上げられた結果がカーソルのある場所にそのまま貼り付けられます。
</details>

<details>
<summary><b>自分の声で読み上げられますか？</b></summary>

Studio プランでは、短いサンプルを録音すると、スタンダードやプレミアムのニューラルボイスに加えて、Vible があなたの声のクローンでテキストを読み返します。
</details>

---

<div align="center">

**[vible.space](https://vible.space)**  ·  Speak. Polish. Translate.

Copyright © 2026 **Algohash Inc.** · Vible は署名・公証済みの macOS アプリケーションです。All rights reserved.

</div>
