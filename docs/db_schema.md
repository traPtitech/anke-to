# db スキーマ

### administrators

アンケートの運営 (編集等ができる人)

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| user_traqid      | char(32) | NO   | PRI | _NULL_  |

### administrator_groups

アンケートの運営 (編集等ができるグループ)（実際の管理はadministratorsで行い、これは前回選択した内容を提示するためのみに使用される）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| group_id         | char(36) | NO   | PRI | _NULL_  |

### administrator_users

アンケートの運営 (編集等ができるユーザー)（実際の管理はadministratorsで行い、これは前回選択した内容を提示するためのみに使用される）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| user_traqid      | char(32) | NO   | PRI | _NULL_  |

### options

選択肢

| Field       | Type    | Null | Key | Default | Extra          | 説明など           |
| ----------- | ------- | ---- | --- | ------- | -------------- | ------------------ |
| id          | int(11) | NO   | PRI | _NULL_  | auto_increment |
| question_id | int(11) | NO   | MUL | _NULL_  |                | どの質問の選択肢か |
| option_num  | int(11) | NO   |     | _NULL_  |                | 何番目の選択肢か   |
| body        | text    | YES  |     | _NULL_  |                | 選択肢の内容       |

### question

質問内容

| Field            | Type       | Null | Key  | Default           | Extra          | 説明など                                                     |
| ---------------- | ---------- | ---- | ---- | ----------------- | -------------- | ------------------------------------------------------------ |
| id               | int(11)    | NO   | PRI  | _NULL_            | auto_increment |                                                              |
| questionnaire_id | int(11)    | YES  |      | _NULL_            |                | どのアンケートの質問か                                       |
| page_num         | int(11)    | NO   |      | _NULL_            |                | アンケートの何ページ目の質問か                               |
| question_num     | int(11)    | NO   |      | _NULL_            |                | アンケートの質問のうち、何問目か                             |
| type             | char(20)   | NO   |      | _NULL_            |                | どのタイプの質問か ("Text","TextArea",  "Number", "MultipleChoice", "Checkbox", "Dropdown", "LinearScale", "Date", "Time") |
| body             | text       | YES  |      | _NULL_            |                | 質問の内容(title)(v1との互換性のためfield nameはbodyのまま)                                               |
| description      | text       | YES  |      | _NULL_            |                | 質問の内容(description)                                        |
| is_required      | tinyint(4) | NO   |      | 0                 |                | 回答が必須である (1) , ない(0)                               |
| deleted_at       | timestamp  | YES  |      | _NULL_            |                | 質問が削除された日時 (削除されていない場合は NULL)           |
| created_at       | timestamp  | NO   |      | CURRENT_TIMESTAMP |                | 質問が作成された日時                                         |

### questionnaires

アンケートの情報

| Field          | Type      | Null | Key | Default           | Extra          | 説明など                                                                                                                |
| -------------- | --------- | ---- | --- | ----------------- | -------------- | ----------------------------------------------------------------------------------------------------------------------- |
| id             | int(11)   | NO   | PRI | _NULL_            | auto_increment |
| title          | char(50)  | NO   | MUL | _NULL_            |                | アンケートのタイトル                                                                                                    |
| description    | text      | NO   |     | _NULL_            |                | アンケートの説明                                                                                                        |
| res_time_limit | timestamp | YES  |     | _NULL_            |                | 回答の締切日時 (締切がない場合は NULL)                                                                                  |
| deleted_at     | timestamp | YES  |     | _NULL_            |                | アンケートが削除された日時 (削除されていない場合は NULL)                                                                |
| res_shared_to  | char(30)  | NO   |     | administrators    |                | アンケートの結果を, 運営は見られる ("administrators"), 回答済みの人は見られる ("respondents") 誰でも見られる ("public") |
| is_anonymous   | boolean   | NO   |     | false             |                | アンケートが匿名解答かどうか                                                                                            |
| created_at     | timestamp | NO   |     | CURRENT_TIMESTAMP |                | アンケートが作成された日時                                                                                              |
| modified_at    | timestamp | NO   |     | CURRENT_TIMESTAMP |                | アンケートが更新された日時                                                                                              |
| is_published   | boolean   | NO   |     | false             |                | アンケートが公開かどうか                                                                                                |

### respondents

アンケートごとの回答者

| Field            | Type      | Null | Key | Default           | Extra          | 説明など                                            |
| ---------------- | --------- | ---- | --- | ----------------- | -------------- | --------------------------------------------------- |
| response_id      | int(11)   | NO   | PRI | _NULL_            | auto_increment | 一つのアンケートに対する一つの回答ごとに振られる ID |
| questionnaire_id | int(11)   | NO   | MUL | _NULL_            |                | どのアンケートへの回答か                            |
| user_traqid      | char(32)  | YES  | MUL | _NULL_            |                | 回答者の traQID                                     |
| modified_at      | timestamp | NO   |     | CURRENT_TIMESTAMP |                | 回答が変更された日時                                |
| submitted_at     | timestamp | YES  |     | _NULL_            |                | 回答が送信された日時 (未送信の場合は NULL)          |
| deleted_at       | timestamp | YES  |     | _NULL_            |                | 回答が破棄された日時 (破棄されていない場合は NULL)  |

### response

回答

| Field       | Type      | Null | Key | Default           | Extra | 説明など                                            |
| ----------- | --------- | ---- | --- | ----------------- | ----- | --------------------------------------------------- |
| response_id | int(11)   | NO   | MUL | _NULL_            |       | 一つのアンケートに対する一つの回答ごとに振られる ID |
| question_id | int(11)   | NO   | MUL | _NULL_            |       | どの質問への回答か                                  |
| body        | text      | YES  |     | _NULL_            |       | 回答の内容                                          |
| modified_at | timestamp | NO   |     | CURRENT_TIMESTAMP |       | 回答が変更された日時                                |
| deleted_at  | timestamp | YES  |     | _NULL_            |       | 回答が破棄された日時 (破棄されていない場合は NULL)  |

### scale_labels

目盛り (LinearScale) 形式の質問の左右のラベル

| Field             | Type    | Null | Key | Default | Extra | 説明など                       |
| ----------------- | ------- | ---- | --- | ------- | ----- | ------------------------------ |
| question_id       | int(11) | NO   | PRI | _NULL_  |       | どの質問のラベルか             |
| scale_label_left  | text    | YES  |     | _NULL_  |       | 左側のラベル (ない場合は NULL) |
| scale_label_right | text    | YES  |     | _NULL_  |       | 右側のラベル (ない場合は NULL) |
| scale_min         | int(11) | YES  |     | _NULL_  |       | スケールの最小値               |
| scale_max         | int(11) | YES  |     | _NULL_  |       | スケールの最大値               |

### validations

`Number`の値制限，`Text`の正規表現によるパターンマッチング．

| Field         | Type    | Null | Key  | Default | Extra | 説明など           |
| ------------- | ------- | ---- | ---- | ------- | ----- | ------------------ |
| question_id   | int(11) | YES  | PRI  | _NULL_  |       | どの質問についてか |
| regex_pattern | text    | YES  |      | _NULL_  |       | 正規表現           |
| min_bound     | text    | YES  |      | _NULL_  |       | 数値の下界         |
| max_bound     | text    | YES  |      | _NULL_  |       | 数値の上界         |

### targets

アンケートの対象者

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| user_traqid      | char(32) | NO   | PRI | _NULL_  |
| is_canceled      | boolean  | NO   |     | false   |       | アンケートの対象者がキャンセルしたかどうか |

### target_groups

選択したアンケートの対象者（グループ）（実際の管理はtargetsで行い、これは前回選択した内容を提示するためのみに使用される）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| group_id         | char(36) | NO   | PRI | _NULL_  |

### target_groups

選択したアンケートの対象者（ユーザー）（実際の管理はtargetsで行い、これは前回選択した内容を提示するためのみに使用される）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| questionnaire_id | int(11)  | NO   | PRI | _NULL_  |
| user_traqid      | char(32) | NO   | PRI | _NULL_  |
