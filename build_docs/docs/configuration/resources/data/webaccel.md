# ウェブアクセラレータ サイト(sakuracloud_webaccel)

---

### 設定例

```hcl
# サイト情報の参照用
data sakuracloud_webaccel "site" {
  name = "example"
  # または
  # domain = "www.example.com"
}
```

### パラメーター

| パラメーター     | 必須    | 名称                   | 初期値        | 設定値    | 補足                                             |
|-- -------- | :---: | -------------------- | :--------: | ------ | -------------------------------------------- --|
| `name`     | △     | サイト名                 | -          | 文字列    | `name`または`domain`いずれか必須                        |
| `domain`   | △     | ドメイン                 | -          | 文字列    | `name`または`domain`いずれか必須                        |

### 属性

| 属性名                       | 名称                        | 補足                                           |
|-- ----------------------- | ------------------------- | ------------------------------------------ --|
| `id`                      | ID                        | -                                            |
| `site_id`                 | サイトID                     | -                                            |
| `origin`                  | オリジン                      | -                                            |
| `subdomain`               | サブドメイン                    | -                                            |
| `domain_type`             | ドメインタイプ                   | -                                            |
| `has_certificate`         | 証明書有無                     | -                                            |
| `host_header`             | ホストヘッダー                   | -                                            |
| `status`                  | ステータス                     | -                                            |
| `cname_record_value`      | CNAMEレコード値                | 独自ドメイン利用時にCNAMEレコードを利用する場合のレコード値 |
| `txt_record_value`        | TXTレコード値                  | 独自ドメイン利用時にTXTレコードを利用する場合のレコード値 |