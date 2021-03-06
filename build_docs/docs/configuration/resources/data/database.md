# データベース(sakuracloud_database)

---

### 設定例

```hcl
# データベースの参照
data "sakuracloud_database" "foobar" {
  name_selectors = ["foobar"]
}
```

## `sakuracloud_database`

データベースアプライアンスを表します。

### パラメーター

|パラメーター       |必須  |名称           |初期値     |設定値                         |補足                                          |
|-----------------|:---:|----------------|:--------:|-------------------------------|----------------------------------------------|
| `name_selectors`  | -   | 検索条件(名称)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `tag_selectors`   | -   | 検索条件(タグ)      | -        | リスト(文字列)           | 複数指定した場合はAND条件  |
| `filter`          | -   | 検索条件(その他)    | -        | オブジェクト             | APIにそのまま渡されます。検索条件を指定してもAPI側が対応していない場合があります。 |
| `zone`          | -   | ゾーン          | -        | `tk1a`<br />`is1b`<br />`is1a` | - |


### 属性

|属性名          | 名称             | 補足                  |
|---------------|------------------|----------------------|
| `id`            | データベースID | -                    |
| `name`          | データベース名   |  - |
| `database_type` | データベースタイプ|  - |
| `plan`          | プラン           | - |
| `user_name`     | ユーザー名       |  - |
| `user_password` | パスワード       |  - |
| `replica_user`     | レプリケーションユーザー名       |  - |
| `replica_password` | レプリケーションパスワード       |  - |
| `allow_networks`| 送信元ネットワーク | 接続を許可するネットワークアドレス(範囲)のリスト |
| `port`          | ポート番号       |  - |
| `backup_time`   | バックアップ開始時刻   | - |
| `backup_weekdays`   | バックアップ取得曜日   | - |
| `switch_id`     | スイッチID      | - |
| `ipaddress1`    | IPアドレス1     | - |
| `nw_mask_len`   | ネットマスク     | - |
| `default_route` | ゲートウェイ     | - |
| `icon_id`       | アイコンID         | - |
| `description`   | 説明           | - |
| `tags`          | タグ           | - |
| `zone`          | ゾーン          | - |

