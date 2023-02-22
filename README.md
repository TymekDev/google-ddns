# Google Domains DDNS Updater

A simple CLI for updating dynamic DNS at Google Domains.
Before an update current IP is checked to prevent being blocked from DDNS API due to redundant requests.

[Reference](https://support.google.com/domains/answer/6147083)

## Usage
```
-d, --domain strings
-n, --interval duration    (default 5m0s)
-p, --password string
-u, --username string
```

## Tips
If you need several domains pointing to the same IP, then there is no need to make multiple DDNS records.
Use a domain with DDNS record and point the other domains to the DDNS one using CNAME records.

Note: a top level domain cannot have a CNAME record. You might need to make that one DDNS as well.
