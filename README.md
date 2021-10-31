# send-to-kindle

Given a page, visit all the links and create one ebook from it.
Then send this to you kindle.

Have fun

```bash
Usage:
  send-to-kindle send [flags]

Flags:
  -c, --cookies string                 Cookie file
  -h, --help                           help for send
  -k, --kindle-email string            Kindle email address
  -l, --language string                Book main language (default "eng")
  -u, --url string                     URL to parse
  -e, --origin-email string            Origin email address
  -p, --origin-email-password string   Origin email password
  -t, --title string                   Book Title
```

## Cookies
If the site is using cookies, use [Get cookies.txt](https://chrome.google.com/webstore/detail/get-cookiestxt/bgaddhkoddajcdgocldbbfleckgcbcid?hl=en) to download the cookies as file

## Email

### Kindle
Get you kindle email from [amazon](https://www.amazon.com/hz/mycd/digital-console/alldevices)

### Origin Email
* Supporting Gmail ATM
* If your gmail account is protected with 2FA, use [app password](https://myaccount.google.com/apppasswords). Generate for mail app and other device.