# dbスキーマ

### administrators

アンケートの管理者 (編集等ができる人)

| Field            | Type        | Null | Key | Default | Extra | 説明など |
| ---------------- | ----------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)     | NO   | PRI | _NULL_  |       |          |
| user_traqid      | varchar(32) | NO   | PRI | _NULL_  |       |          |

### options

選択肢

| Field       | Type    | Null | Key | Default | Extra          | 説明など           |
| ----------- | ------- | ---- | --- | ------- | -------------- | ------------------ |
| id          | int(11) | NO   | PRI | _NULL_  | AUTO_INCREMENT |
| question_id | int(11) | NO   | MUL | _NULL_  |                | どの質問の選択肢か |
| option_num  | int(11) | NO   |     | _NULL_  |                | 何番目の選択肢か   |
| body        | text    | YES  |     | _NULL_  |                | 選択肢の内容       |

### questions

質問内容

| Field            | Type       | Null | Key  | Default           | Extra          | 説明など                                                     |
| ---------------- | ---------- | ---- | ---- | ----------------- | -------------- | ------------------------------------------------------------ |
| id               | int(11)    | NO   | PRI  | _NULL_            | AUTO_INCREMENT |                                                              |
| questionnaire_id | int(11)    | YES  |      | _NULL_            |                | どのアンケートの質問か                                       |
| page_num         | int(11)    | NO   |      | _NULL_            |                | アンケートの何ページ目の質問か                               |
| question_num     | int(11)    | NO   |      | _NULL_            |                | アンケートの質問のうち、何問目か                             |
| type             | int(11)    | NO   | MUL  | _NULL_            |                | どのタイプの質問か |
| body             | text       | YES  |      | _NULL_            |                | 質問の内容                                                   |
| is_required      | boolean | NO   |      | 0                 |                | 回答が必須でか                               |
| created_at       | datetime  | NO   |      | CURRENT_TIMESTAMP |                | 質問が作成された日時                                         |
| deleted_at       | datetime  | YES  |      | _NULL_            |                | 質問が削除された日時 (削除されていない場合はNULL)           |

### question_types

質問の種類。
'Text'、'TextArea'、'Number'、'MultipleChoice'、'Checkbox', 'Dropdown', 'LinearScale', 'Date', 'Time'。

| Field  | Type       | Null | Key  | Default           | Extra          | 説明など |
| ------ | ---------- | ---- | ---- | ----------------- | -------------- | -------- |
| id     | int(11)    | NO   | PRI  | _NULL_            | AUTO_INCREMENT |          |
| name   | varchar(30)    | NO   |      | _NULL_            |                |          |
| active | boolean    | NO   |      | true              |                |          |

### questionnaires

アンケートの情報

| Field          | Type     | Null | Key | Default           | Extra          | 説明など                                                                                                                |
| -------------- | -------- | ---- | --- | ----------------- | -------------- | ----------------------------------------------------------------------------------------------------------------------- |
| id             | int(11)  | NO   | PRI | _NULL_            | AUTO_INCREMENT |
| title          | char(50) | NO   | MUL | _NULL_            |                | アンケートのタイトル |
| description    | text     | NO   |     | _NULL_            |                | アンケートの説明                                                                                                        |
| res_time_limit | datetime | YES  |     | _NULL_            |                | 回答の締切日時 (締切がない場合はNULL) |
| res_shared_to  | int(11)  | NO   | MUL | _NULL_            |                |                                       |
| created_at     | datetime | NO   |     | CURRENT_TIMESTAMP |                | アンケートが作成された日時                                                                                              |
| updated_at     | datetime | NO   |     | CURRENT_TIMESTAMP |                | アンケートが更新された日時                                                                                              |
| deleted_at     | datetime | YES  |     | _NULL_            |                | アンケートが削除された日時 (削除されていない場合はNULL)                                                                |

### res_share_types

アンケート結果の公開範囲の種類。
アンケートの結果を、運営は見られる ("administrators")、回答済みの人は見られる ("respondents")、誰でも見られる ("public")。

| Field  | Type       | Null | Key  | Default           | Extra          | 説明など |
| ------ | ---------- | ---- | ---- | ----------------- | -------------- | -------- |
| id     | int(11)    | NO   | PRI  | _NULL_            | AUTO_INCREMENT |          |
| name   | varchar(30)    | NO   |      | _NULL_            |                |          |
| active | boolean    | NO   |      | true              |                |          |

### respondents

アンケートごとの回答者

| Field            | Type      | Null | Key | Default           | Extra          | 説明など                                            |
| ---------------- | --------- | ---- | --- | ----------------- | -------------- | --------------------------------------------------- |
| response_id      | int(11)   | NO   | PRI | _NULL_            | AUTO_INCREMENT | 1つのアンケートに対する1つの回答ごとに振られるID |
| questionnaire_id | int(11)   | NO   | MUL | _NULL_            |                | どのアンケートへの回答か                            |
| user_traqid      | varchar(32)  | YES  | MUL | _NULL_            |                | 回答者のtraQID                                     |
| submitted_at     | datetime | YES  |     | _NULL_            |                | 回答が送信された日時 (未送信の場合はNULL)          |
| updated_at       | datetime | NO   |     | CURRENT_TIMESTAMP |                | 回答が変更された日時                                |
| deleted_at       | datetime | YES  |     | _NULL_            |                | 回答が破棄された日時 (破棄されていない場合はNULL)  |

### responses

回答

| Field       | Type      | Null | Key | Default           | Extra | 説明など                                            |
| ----------- | --------- | ---- | --- | ----------------- | ----- | --------------------------------------------------- |
| response_id | int(11)   | NO   | MUL | _NULL_            |       | 1つのアンケートに対する1つの回答ごとに振られるID |
| question_id | int(11)   | NO   | MUL | _NULL_            |       | どの質問への回答か                                  |
| body        | text      | YES  |     | _NULL_            |       | 回答の内容                                          |
| updated_at  | datetime  | NO   |     | CURRENT_TIMESTAMP |       | 回答が変更された日時                                |
| deleted_at  | datetime  | YES  |     | _NULL_            |       | 回答が破棄された日時 (破棄されていない場合はNULL)  |

### scale_labels

目盛り (LinearScale) 形式の質問の左右のラベル

| Field             | Type        | Null | Key | Default | Extra | 説明など                       |
| ----------------- | ----------- | ---- | --- | ------- | ----- | ------------------------------ |
| question_id       | int(11)     | NO   | PRI | _NULL_  |       | どの質問のラベルか             |
| scale_label_left  | varchar(50) | YES  |     | _NULL_  |       | 左側のラベル (ない場合はNULL) |
| scale_label_right | varchar(50) | YES  |     | _NULL_  |       | 右側のラベル (ない場合はNULL) |
| scale_min         | int(11)     | YES  |     | _NULL_  |       | スケールの最小値               |
| scale_max         | int(11)     | YES  |     | _NULL_  |       | スケールの最大値               |

### validations

`Number`の値制限、`Text`の正規表現によるパターンマッチング。

| Field         | Type    | Null | Key  | Default | Extra | 説明など           |
| ------------- | ------- | ---- | ---- | ------- | ----- | ------------------ |
| question_id   | int(11) | YES  | PRI  | _NULL_  |       | どの質問についてか |
| regex_pattern | text    | YES  |      | _NULL_  |       | 正規表現           |
| min_bound     | text    | YES  |      | _NULL_  |       | 数値の下界         |
| max_bound     | text    | YES  |      | _NULL_  |       | 数値の上界         |

### targets

アンケートの対象者

| Field            | Type        | Null | Key | Default | Extra | 説明など |
| ---------------- | ----------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)     | NO   | PRI | _NULL_  |       |          |
| user_traqid      | varchar(32) | NO   | PRI | _NULL_  |       |          |
