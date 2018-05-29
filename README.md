# IMEOS.one News Mixin bot
Go implementation of Mixin bot that posts IMEOS News stories to Mixin Messeger

## Build
1. [**Obtain your own Mixin bot app key**](https://developers.mixin.one/dashboard), fill `config/config.go` api key and secret...
2. Got dependencies and run `go build`
3. Launch two services `./imeos-news-mixin-bot -service blaze` and `./imeos-news-mixin-bot -service news`
