# dbスキーマ

### administrators

アンケートの管理者 (編集等ができる人)  
`user_id`、`group_id`のいずれかは必ず非null

| Field            | Type        | Null | Key | Default | Extra | 説明など |
|------------------|-------------|------|-----|---------|-------|------|
| questionnaire_id | int(11)     | NO   | PRI | _NULL_  |       |      |
| user_id          | varchar(32) | YES  | PRI | _NULL_  |       |      |
| group_id         | varchar(32) | YES  | PRI | _NULL_  |       |      |

### targets

アンケートの対象者  
`user_id`、`group_id`のいずれかは必ず非null

| Field            | Type        | Null | Key | Default | Extra | 説明など |
|------------------|-------------|------|-----|---------|-------|------|
| questionnaire_id | int(11)     | NO   | PRI | _NULL_  |       |      |
| user_id          | varchar(32) | YES  | PRI | _NULL_  |       |      |
| group_id         | varchar(32) | YES  | PRI | _NULL_  |       |      |

### options

選択肢

| Field       | Type    | Null | Key | Default | Extra          | 説明など      |
|-------------|---------|------|-----|---------|----------------|-----------|
| id          | int(11) | NO   | PRI | _NULL_  | AUTO_INCREMENT |
| question_id | int(11) | NO   | MUL | _NULL_  |                | どの質問の選択肢か |
| order       | int(11) | NO   |     | _NULL_  |                | 何番目の選択肢か  |
| text        | text    | YES  |     | _NULL_  |                | 選択肢の内容    |

### questions

質問内容

| Field            | Type     | Null | Key | Default           | Extra          | 説明など                         |
|------------------|----------|------|-----|-------------------|----------------|------------------------------|
| id               | int(11)  | NO   | PRI | _NULL_            | AUTO_INCREMENT |                              |
| questionnaire_id | int(11)  | YES  |     | _NULL_            |                | どのアンケートの質問か                  |
| order            | int(11)  | NO   |     | _NULL_            |                | アンケートの質問のうち、何問目か             |
| type             | int(11)  | NO   | MUL | _NULL_            |                | どのタイプの質問か                    |
| text             | text     | YES  |     | _NULL_            |                | 質問の内容                        |
| is_required      | boolean  | NO   |     | 0                 |                | 回答が必須でか                      |
| created_at       | datetime | NO   |     | CURRENT_TIMESTAMP |                | 質問が作成された日時                   |
| updated_at       | datetime | NO   |     | CURRENT_TIMESTAMP |                | 質問が更新された日時                   |
| deleted_at       | datetime | YES  |     | _NULL_            |                | 質問が削除された日時 (削除されていない場合はNULL) |

### question_types

質問の種類  
`Text`、`TextArea`、`Number`、`MultipleChoice`、`Checkbox`、`Dropdown`、`LinearScale`、`Date`、`Time`

| Field | Type        | Null | Key | Default | Extra          | 説明など |
|-------|-------------|------|-----|---------|----------------|------|
| id    | int(11)     | NO   | PRI | _NULL_  | AUTO_INCREMENT |      |
| name  | varchar(30) | NO   |     | _NULL_  |                |      |

### questionnaires

アンケートの情報

| Field          | Type     | Null | Key | Default           | Extra          | 説明など                            |
|----------------|----------|------|-----|-------------------|----------------|---------------------------------|
| id             | int(11)  | NO   | PRI | _NULL_            | AUTO_INCREMENT |
| title          | char(50) | NO   | MUL | _NULL_            |                | アンケートのタイトル                      |
| description    | text     | NO   |     | _NULL_            |                | アンケートの説明                        |
| deadline       | datetime | YES  |     | _NULL_            |                | 回答の締切日時 (締切がない場合はNULL)          |
| res_visibility | int(11)  | NO   | MUL | _NULL_            |                |                                 |
| created_at     | datetime | NO   |     | CURRENT_TIMESTAMP |                | アンケートが作成された日時                   |
| is_multiple    | boolean  | NO   |     | 0                 |                | 複数回答を許すか                        |
| is_anonymous   | boolean  | NO   |     | 0                 |                | 匿名のアンケートか                       |
| is_editable    | boolean  | NO   |     | 1                 |                | 回答の編集を許すか                       |
| is_draft       | boolean  | NO   |     | 0                 |                | アンケートが下書き状態か                    |
| is_public      | boolean  | NO   |     | 0                 |                | 外部公開アンケートか                      |
| created_at     | datetime | NO   |     | CURRENT_TIMESTAMP |                | アンケートが作成された日時                   |
| updated_at     | datetime | NO   |     | CURRENT_TIMESTAMP |                | アンケートが更新された日時                   |
| deleted_at     | datetime | YES  |     | _NULL_            |                | アンケートが削除された日時 (削除されていない場合はNULL) |

### res_visibility_types

アンケート結果の公開範囲の種類  
アンケートの結果を、運営は見られる (`administrators`)、回答済みの人は見られる (`respondents`)、誰でも見られる (`public`)

| Field | Type        | Null | Key | Default | Extra          | 説明など |
|-------|-------------|------|-----|---------|----------------|------|
| id    | int(11)     | NO   | PRI | _NULL_  | AUTO_INCREMENT |      |
| name  | varchar(30) | NO   |     | _NULL_  |                |      |

### responses

アンケートごとの回答  
質問に対する回答 (`answer`) を束ねる

| Field            | Type        | Null | Key | Default           | Extra          | 説明など                         |
|------------------|-------------|------|-----|-------------------|----------------|------------------------------|
| id               | int(11)     | NO   | PRI | _NULL_            | AUTO_INCREMENT |                              |
| questionnaire_id | int(11)     | NO   | MUL | _NULL_            |                | どのアンケートへの回答か                 |
| user_id          | varchar(32) | YES  | MUL | _NULL_            |                | 回答者のtraQID                   |
| submitted_at     | datetime    | YES  |     | _NULL_            |                | 回答が送信された日時 (未送信の場合はNULL)     |
| updated_at       | datetime    | NO   |     | CURRENT_TIMESTAMP |                | 回答が変更された日時                   |
| deleted_at       | datetime    | YES  |     | _NULL_            |                | 回答が破棄された日時 (破棄されていない場合はNULL) |

### normal_answers

質問 (選択肢形式以外) に対する回答

| Field       | Type    | Null | Key | Default | Extra | 説明など                           |
|-------------|---------|------|-----|---------|-------|--------------------------------|
| response_id | int(11) | NO   | MUL | _NULL_  |       | どの回答(`response`)の回答(`answer`)か |
| question_id | int(11) | NO   | MUL | _NULL_  |       | どの質問への回答か                      |
| text        | text    | YES  |     | _NULL_  |       | 回答の内容                          |

### options_answers

質問 (選択肢形式) に対する回答

| Field       | Type    | Null | Key | Default | Extra | 説明など                           |
|-------------|---------|------|-----|---------|-------|--------------------------------|
| response_id | int(11) | NO   | MUL | _NULL_  |       | どの回答(`response`)の回答(`answer`)か |
| question_id | int(11) | NO   | MUL | _NULL_  |       | どの質問への回答か                      |
| option_id   | int(11) | YES  | MUL | _NULL_  |       | 回答選択肢のID                       |

### scales

目盛り (LinearScale) 形式の質問の目盛りの情報と左右のラベル

| Field       | Type        | Null | Key | Default | Extra | 説明など               |
|-------------|-------------|------|-----|---------|-------|--------------------|
| question_id | int(11)     | NO   | PRI | _NULL_  |       | どの質問のラベルか          |
| left_label  | varchar(50) | YES  |     | _NULL_  |       | 左側のラベル (ない場合はNULL) |
| right_label | varchar(50) | YES  |     | _NULL_  |       | 右側のラベル (ない場合はNULL) |
| min         | int(11)     | YES  |     | _NULL_  |       | スケールの最小値           |
| max         | int(11)     | YES  |     | _NULL_  |       | スケールの最大値           |
| step_width  | int(11)     | YES  |     | _NULL_  |       | スケールの最大値           |

### validations

`Number`の値制限、`Text`の正規表現によるパターンマッチング

| Field         | Type    | Null | Key | Default | Extra | 説明など      |
|---------------|---------|------|-----|---------|-------|-----------|
| question_id   | int(11) | YES  | PRI | _NULL_  |       | どの質問についてか |
| regex_pattern | text    | YES  |     | _NULL_  |       | 正規表現      |
| min_bound     | text    | YES  |     | _NULL_  |       | 数値の下界     |
| max_bound     | text    | YES  |     | _NULL_  |       | 数値の上界     |

### tags

アンケートに付けるタグ

| Field           | Type        | Null | Key | Default | Extra          | 説明など |
|-----------------|-------------|------|-----|---------|----------------|------|
| id              | int(11)     | NO   | PRI | _NULL_  | AUTO_INCREMENT |      |
| name            | varchar(30) | NO   |     | _NULL_  |                |      |
| name_lower_case | varchar(30) | NO   |     | _NULL_  |                |      |


### questionnaires_tags

アンケートとタグの関連

| Field            | Type        | Null | Key | Default | Extra | 説明など |
|------------------|-------------|------|-----|---------|-------|------|
| questionnaire_id | int(11)     | NO   | PRI | _NULL_  |       |      |
| tag_id           | int(11)     | NO   | PRI | _NULL_  |       |      |

### change_logs

アンケートの更新履歴  

| Field            | Type    | Null | Key | Default | Extra | 説明など          |
|------------------|---------|------|-----|---------|-------|---------------|
| questionnaire_id | int(11) | NO   | PRI | _NULL_  |       | どのアンケートの更新履歴か |
| body             | text    | YES  |     | _NULL_  |       | JSON形式の更新内容   |
