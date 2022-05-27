# MySQL Sample
MySQLについての学習記録

## MySQL立ち上げ
以下コマンドを実行
```
docker compose build
docker compose up -d
docker compose exec mysql bash
mysql -u root -D app_develop -p
root
```

## サンプルデータファイル
secure_file_privの値として設定されているパスのみファイルの読み込みや書き込みが可能  
以下のパスがMySQLのDockerコンテナで設定されているsecure_file_privの値である為、このフォルダに読み込むCSVを配置する  
`/var/lib/mysql-files`  
設定値は下記コマンドで確認できる  
```
SHOW VARIABLES LIKE "secure_file_priv";
```
  
以下コマンドでデータを挿入する  
```
LOAD DATA INFILE "/var/lib/mysql-files/sample_data.csv" INTO TABLE users FIELDS TERMINATED BY ",";
```
  
LOAD DATA INFILE: サンプルデータのファイルパスを指定  
INTO TABLE: 挿入するテーブルを指定  
FIELDS TERMINATED BY: カラムの区切り文字を指定  
レコードの区切り文字はデフォルトで\nに設定されているので、環境によっては変更が必要  
NULLを挿入する場合は、\Nを使用する  

## コマンド
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
```
SELECT * FROM users WHERE name = "jjdd";
```
実行速度 0.05 → 0.01  
type: ALL → ref  
key: NULL → name_index  
rows: 99822 → 1  

postsテーブルのuser_idカラムにindexを貼る  
```
select id, name, post_count from users as u inner join (select user_id, count(user_id) as post_count from posts group by user_id) as p on u.id = p.user_id;
```
1.36 sec → 0.40 sec  

## 外部キー制約
```
FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
```
usersテーブルのid列を参照するuser_idという外部キーを定義する  
テーブルの参照整合性が保たれる  
`ON DELETE CASCADE`でuserを削除した時に子のカラムを一緒に削除する  
`ON DELETE SET NULL`等もある  
`ON UPDATE`で更新時の挙動を定義できる  

## LAST_INSERT_ID
最後に`INSERT`した`AUTO INCREMENT`の値が保存されている  
新しいユーザを登録した直後に  
`SELECT LAST_INSERT_ID()`  
を実行すると登録したユーザのIDが返ってくる  
したがって、IDを指定して`INSERT`した場合は`AUTO INCREMENT`が働かないので、`LAST_INSERT_ID`の値は変わらない 

## CHECK制約
```
`status` INT(10) NOT NULL CHECK (0 < `status` AND `status` <= 3
``` 
等でカラムの値を制限することができるが、MySQL8.0.16で追加された機能みたいなので、本環境では動作まで確認できない  

## パーティショニング
テーブルを分割する  
テーブルは物理的に分割されているが、あたかも1つのテーブルであるかのように扱うことができる
パーティションではなく、本当に別テーブルを定義すると、外部キー制約や、疑似キーがうまく機能しなくなる
データ量が多くなってきた時にパフォーマンスの向上が見込める
ビューと変わらないかと思ったが、ビューはテーブルの実体を持たないので、パフォーマンスは向上しない？

## COALESCE
引数に指定したもののうち、NULLでない最初の値を返す関数
```
SELECT id, name, COALESCE(landline_phone, mobile_phone, 'unregistered') FROM users;
```
mobile_phoneのみ登録している場合はmobile_phoneの値を返し、何も登録していない場合はunregisteredを返す
