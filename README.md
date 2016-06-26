# go-imgur-downloader
concurrently download images given an imgur link

Usage

```
go get github.com/1vn/go-imgur-downloader 
go-imgur-downloader -url=IMGUR_URL -d=PATH_TO_SAVE -o=MAINTAIN_ORDER(true/false)
```

TODO
- download in other formats, currently defaults to .jpg
- offer more naming otpions
- support imgur galleries
- buffer goroutines to avoid being blocked when used on large albums 
