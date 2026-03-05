# `gofix` vars

I like very explicit var declarations. It's perfectly normal to be this fixated
on minor things.

## Usage

(a.k.a., how to be as annoying about variable declarations as I am)

```bash
go install github.com/jechols/go-fixvars@v0.0.1
go fix -fixtool=$(go env GOPATH)/bin/go-fixvars /path/to/ugly/code
```
