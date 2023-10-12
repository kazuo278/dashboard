# README

Self-Hosted Runnerで実行されるジョブと紐づくActionやReusableWorkflow情報を保管するアプリ

## 使い方

- ダッシュボード
  - http://${DASHBOARD_APP_HOST}/dashboard

- 実行履歴登録・更新REST API  
  Self-Hosted Runnerの`ACTIONS_RUNNER_HOOK_JOB_STARTED`および`ACTIONS_RUNNER_HOOK_JOB_COMPLETED`を利用してジョブ情報を登録・更新する想定。  
  詳細はAPI仕様は、[openapi.yaml](/docs/rest/openapi.yaml)を参照。
  - ジョブ実行開始履歴の登録

    ```sh
    $ WORKFLOW_REF=$(echo $GITHUB_WORKFLOW_REF | sed "s%$GITHUB_REPOSITORY/%%")
    $ curl -X POST ${DASHBOARD_APP_HOST}/actions/history -H 'Content-Type: application/json' -d @- <<EOM
    {
      "repository_id":"$GITHUB_REPOSITORY_ID",
      "repository_name":"$GITHUB_REPOSITORY",
      "run_id":"$GITHUB_RUN_ID",
      "workflow_ref":"$WORKFLOW_REF",
      "job_name":"$GITHUB_JOB",
      "run_attempt":"$GITHUB_RUN_ATTEMPT"
    }
    EOM
    ```

  - ジョブ終了履歴の登録

    ```sh
    $ curl -X PUT ${DASHBOARD_APP_HOST}/actions/history -H 'Content-Type: application/json' -d @- <<EOM
    {
      "repository_id":"$GITHUB_REPOSITORY_ID",
      "run_id":"$GITHUB_RUN_ID",
      "job_name":"$GITHUB_JOB",
      "run_attempt":"$GITHUB_RUN_ATTEMPT"
    }
    EOM
    ```

## 開発方法

- 事前準備  
  `.devcontainer/organization-token.txt`にAction実行履歴を取得したいリポジトリ閲覧権限をもったトークンを記載する。

- 開発モード起動

  ```sh
  air -c .air.toml
  ```

- ビルド

  ```sh
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo
  ```

- イメージ作成

  ```sh
  docker build -f docker/Dockerfile .
  ```
