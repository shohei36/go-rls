# go-rls
GoとPostgreSQLのRLSを使って、マルチテナント環境で安全にアクセス制御する実装サンプル

# プロジェクトの構成
Clean Architecture をベースにしている
```
.
├── Makefile
├── README.md
├── db.go
├── go.mod
├── go.sum
├── main.go  // Main関数（controllerを含む）
├── model.go // ドメインモデル
├── pgsql
│   ├── docker-compose.yaml
│   └── init
│       ├── 001_ddl.sql   // 初期化DDL
│       └── 002_data.sql  // サンプルデータ
├── repository.go   // DBなど外部サービスとのアダプター
├── transaction.go  // トランザクション
└── usecase.go      // アプリケーションのビジネスロジック
```

# 試してみる
PostgreSQLコンテナを立てる
```
$ cd pgsql
$ docker compose up -d
[+] Running 2/2
 ⠿ Volume "pgsql_pgsqldata"  Created                                                                                                              0.0s
 ⠿ Container my_pgsql        Started
```
Wabサーバーを起動する
```
$ go run .
```

ユーザーを取得するAPIを実行
```
$ curl -X GET http://localhost:8080/users/1000 -H "tenant-id:tenant01" | python -m json.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    60  100    60    0     0   8571      0 --:--:-- --:--:-- --:--:--  8571
{
    "user": {
        "id": "1000",
        "name": "Bob",
        "gender": "male",
        "age": 22
    }
}
```
ヘッダにセットするテナントIDをtenant02にすると、同じユーザーIDでもtenant02が持つユーザーのみにアクセス可能
```
$ curl -X GET http://localhost:8080/users/1000 -H "tenant-id:tenant02" | python -m json.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    62  100    62    0     0   8857      0 --:--:-- --:--:-- --:--:--  8857
{
    "user": {
        "id": "1000",
        "name": "James",
        "gender": "male",
        "age": 44
    }
}
```