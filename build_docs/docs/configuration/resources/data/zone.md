# ゾーン(sakuracloud_zone)

---

### 設定例

```hcl
data sakuracloud_zone "current" {}

data sakuracloud_zone "is1a" {
  name = "is1a"
}
```

### パラメーター

|パラメーター|必須  |名称                |初期値     |設定値 |補足                                          |
|----------|:---:|--------------------|:--------:|------|----------------------------------------------|
| `name`  | -   | ゾーン名      | -        | 文字列           | 省略した場合はプロバイダー設定が利用される|

### 属性

|属性名                    | 名称                     | 補足                                        |
|-------------------------|-------------------------|--------------------------------------------|
| `id`          | ID         | -                                          |
| `name`        | ゾーン名    | - |
| `zone_id`     | ゾーンID    | - |
| `description` | 説明        | - |
| `region_id`   | リージョンID | - |
| `region_name` | リージョン名 | - |
| `dns_servers` | リージョンのDNSサーバIPアドレスのリスト | - |
