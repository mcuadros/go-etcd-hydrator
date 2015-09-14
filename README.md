go-etcd-hydrator [![Build Status](https://travis-ci.org/mcuadros/go-etcd-hydrator.png?branch=master)](https://travis-ci.org/mcuadros/go-etcd-hydrator) [![GoDoc](http://godoc.org/github.com/mcuadros/go-etcd-hydrator?status.png)](http://godoc.org/github.com/mcuadros/go-etcd-hydrator)
==============================


Installation
------------
The recommended way to install go-etcd-hydrator

```
go get github.com/mcuadros/go-etcd-hydrator
```

Examples
--------
This is a very basic example reading a couple of keys from the etcd server:

```go
import (
    "fmt"

    "github.com/coreos/go-etcd/etcd"
    "github.com/mcuadros/go-etcd-hydrator"
)

type MyAppConfig struct {
    Rest string //<-- the hydrator will ask to etcd for `rest` key
    MongoDB string `etcd:"mongo.host"` //<--  ask to etcd for `mongo.host` key
}

func NewMyAppConfig() *MyAppConfig {
    example := new(MyAppConfig)

    h := NewHydrator(etcd.NewClient([]string{"http://127.0.0.1:2379"}))
    h.Hydrate(foo)

    return example
}

...
config := NewMyAppConfig()
fmt.Println(config.Rest) //Prints: http://rest.company.com
fmt.Println(config.MongoDB) //Prints: 127.0.0.1:27017
```

Required keys on the `etcd` server:

```sh
curl -L -X PUT http://127.0.0.1:4001/v2/keys/rest -d value="http://rest.company.com"
curl -L -X PUT http://127.0.0.1:4001/v2/keys/mongo.host -d value="127.0.0.1:27017"
```

Debug
-----
You can enable the debug mode of this library setting the environment variable
`ETCD_HYDRATOR_DEBUG` with the value `true`. (in fact any value is valid)

When this variable is set a small log will be printed to the stdout:
```
Hydrating var: "*hydrator.Example"
Recovered key "testing/string" with value "foo"
Recovered key "testing/integer" with value "42"
...
```

License
-------

MIT, see [LICENSE](LICENSE)
