# MySQL Sample
MySQLについての学習記録

### MySQL立ち上げ
以下コマンドを実行
```
docker compose build
docker compose up -d
docker compose exec mysql bash
mysql -u root -D app_develop -p
root
```

### サンプルデータファイル
secure_file_privの値として設定されているパスのみファイルの読み込みや書き込みが可能
以下のパスがMySQLのDockerコンテナで設定されているsecure_file_privの値である為、このフォルダに読み込むCSVを配置する
/var/lib/mysql-files
設定値は下記コマンドで確認できる
SHOW VARIABLES LIKE "secure_file_priv";

以下コマンドでデータを挿入する
LOAD DATA INFILE "/var/lib/mysql-files/sample_data.csv" INTO TABLE users FIELDS TERMINATED BY ",";

LOAD DATA INFILEでサンプルデータのファイルパスを指定
INTO TABLEで挿入するテーブルを指定
FIELDS TERMINATED BYでカラムの区切り文字を指定
レコードの区切り文字はデフォルトで\nに設定されているので、環境によっては変更が必要

NULLを挿入する場合は、\Nを使用する

### コマンド
指定したテーブルの全データを削除する
TRUNCATE TABLE users;

### EXPLAIN
詳細はこの辺が分かりやすい
http://nippondanji.blogspot.com/2009/03/mysqlexplain.html

- type: 対象のテーブルに対してどのような方法でアクセスするか
ALL: フルテーブルスキャンになっているので注意が必要
ref: ユニークでないインデックスを使って等価検索
- possible_keys: 利用可能なインデックスの候補
- key: 選択されたキー
- key_len: 選択されたキーの長さ。短い方が高速である
- そのテーブルからフェッチされる行数の見積もり

10万行のusersテーブルのnameカラムにindexを貼った場合
SELECT * FROM users WHERE name = "jjdd";
実行速度 0.05 → 0.01
type: ALL → ref
key: NULL → name_index
rows: 99822 → 1

postsテーブルのuser_idカラムにindexを貼る
select id, name, post_count from users as u inner join (select user_id, count(user_id) as post_count from posts group by user_id) as p on u.id = p.user_id;
1.36 sec → 0.40 sec
