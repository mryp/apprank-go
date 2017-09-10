FORMAT: 1A

# apprank-go API

# Group ランキングデータAPI

## 最新ランキング取得 [/v1/now{?country,kind}]
### GET

* 現在の最新ランキングを取得する
* データ種別番号は下記の一覧から指定する
    * 1 - iPhone無料
    * 2 - iPhone有料
    * 3 - iPhoneセールス

+ Parameters
    + country: jp (string, required) - 国
    + kind: 1 (number, required) - データ種別

+ Response 200 (application/json)
    + Attributes
        + updated: 2017-09-01 10:00:00 (string, required) - 順位データ取得日時
        + apps (array) - アプリ情報リスト
            + (object)
                + id: 1000000 (number)  - アプリID
                + name: ほげほげアプリ (string) - アプリ名
                + artwork_url: http://www.hoge/hoge.jpg (string)  - アイコン画像URL


## アプリ詳細取得 [/v1/appinfo{?id}]
### GET

* 指定したアプリの詳細情報を取得する

+ Parameters
    + id: 10000000 (number, required) - アプリID

+ Response 200 (application/json)
    + Attributes
        + name: ほげほげアプリ (string, required) - アプリ名
        + info_url: http://www/hoge/hoge (string, required) - アプリ詳細URL
        + artwork_url: http://www/hoge/hoge.jpg (string, required) - アイコン画像URL
        + artist_name: ほげほげ会社 (string, required) - 著作者名
        + artist_url: http://www/hoge (string, required) - 著作者URL
        + copyright: @ほげほげ (string, required) - 著作権表示


リリース日
## アプリ順位取得 [/v1/apprank{?id,country,kind,start,end}]
### GET

* 指定したアプリの順位一覧を取得する

+ Parameters
    + id: 10000000 (number, required) - アプリID
    + country: jp (string, required) - 国
    + kind: 1 (number, required) - データ種別
    + start: 2017-01-01 (string, required) - 取得開始日
    + end: 2017-09-06 (string, required) - 取得終了日

+ Response 200 (application/json)
    + Attributes
        + apps (array) - 順位情報リスト
            + (object)
                + rank: 9 (number)  - 順位
                + updated: 2017-09-01 10:00:00 (string) - 順位データ取得日時

