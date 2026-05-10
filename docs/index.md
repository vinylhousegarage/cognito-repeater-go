# Cognito Repeater Go (Go 製 Cognito 認証中継 API)

### 1. 概要
  - **目的**
    - 本 API は、FastAPI製認証中継API のコールドスタート起動時間を短縮することを目的としています。

  - **技術選定**
    - インタプリタ言語（Python/FastAPI）からコンパイル言語（Go）に刷新することで、AWS Lambda の起動速度を 1,329ms から 127ms に短縮しました。

      #### FastAPI (改善前)
      ![FastAPI Log](./docs/images/cut-cognito-repeater_InitDuration_20260426.png)

      #### Go (改善後)
      ![Go Log](./docs/images/cut-cognito-repeater-go_InitDuration_20260426.png)

  - **提供機能**
    - 本 API は、AWS Cognito を利用した認証処理を中継し、複数アプリ間で共通利用できるログイン機能を提供します。
      - Cognito へのログインおよびログアウト
      - Cognito が発行する id_token の署名および標準クレーム（iss・aud・exp）の検証
      - Cognito が発行する access_token によるユーザーアカウントの有効確認
      - Cognito が発行する refresh_token の強制失効
      - 本 API による OpenAPI 3.0.3 仕様ドキュメント（JSON 形式）

### 2. ルートURL
  - ### [https://cognito-repeater-go.com](https://cognito-repeater-go.com)
  -  ルートURL にアクセスすると自動的に /login に転送され、Cognito ログイン画面が表示されます。

### 3. エンドポイント
  - すべてのエンドポイントは、ルートURL に対する相対パスです。

|メソッド|パス|用途|正常時の戻り値|Content-Type|Status Code|
|-------|----|----|-------------|-----------|----|
|GET|/health|稼働確認|{"status": "healthy"}|application/json|200|
|GET|/login|Cognito へのログイン処理|Cognito ログイン画面へリダイレクト|(リダイレクト)|302|
|GET|/logout|Cognito からのログアウト処理|"Logout successful"|text/plain |200|
|POST|/revoke|refresh_token の強制失効|(なし)|(なし)|204|
|GET|/me|署名および標準クレーム (iss・aud・exp) の検証|{"sub": "<ユーザーID>"}|application/json|200|
|GET|/whoami|Cognito ユーザーアカウントの有効確認|{"sub": "<ユーザーID>"}|application/json|200|
|GET|/openapi.json|OpenAPI仕様の取得|OpenAPI 3.0.3 ドキュメント|application/json|200|

### 4. トークン指定方法
  - 以下のエンドポイントでは、認証用トークンの指定が必要です。
  - 各トークンの種類と指定方法は、エンドポイントによって異なります。

    - **/logout**
      - クエリパラメーターに id_token を指定
    - **/revoke**
      - リクエストボディに refresh_token を JSON 形式で指定
    - **/me**
      - Authorization ヘッダーに Bearer id_token を指定
    - **/whoami**
      - Authorization ヘッダーに Bearer access_token を指定

### 5. システム構成
  - **技術スタック**
    - プログラミング言語：Go 1.24.3
    - 認証機能：Cognito
    - 仮想環境構築：Docker
      - 開発環境：Docker でコンテナを起動
      - 本番環境：AWS Lambda で Dockerイメージを使用
    - テスト環境：GitHub Actions
    - ローカル環境のバージョン管理：Git
    - リモートリポジトリのホスティング：GitHub
    - CI/CD：GitHub Actions

  - **インフラ構成**
    - API 実行クラウド：AWS
    - API 実行環境：Lambda
    - イメージ管理：ECR
    - ハンドラー：aws-lambda-go + aws-lambda-go-api-proxy/httpadapter
    - API Gateway：HTTP API
    - ドメイン・DNS管理：Route 53

### 6. アクセス情報
  - **APIエンドポイント**
    - [https://cognito-repeater-go.com](https://cognito-repeater-go.com)
  - **API仕様（OpenAPI）**
    - [https://cognito-repeater-go.com/openapi.json](https://cognito-repeater-go.com/openapi.json)

### 7. ライセンス
  - 本 API は [MIT License](https://opensource.org/licenses/MIT) のもとで公開されています。
