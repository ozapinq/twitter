Tweetserver
===

[API documentation](http://twitter.totallyfakedomain.xyz:32000/docs)


## Building

### Build image

```sh
$ make
```

To specify image version:

```sh
$ TAG=1.0 make
```

### Push image

```sh
$ make push
```

## Running tests

### Unit tests

```sh
$ cd $GOPATH/github.com/ozapinq/twitter/
$ go test ./...
```

### Functional API tests

```sh
$ kubectl delete -f k8s/tester.yaml & kubectl apply -f k8s/tester.yaml
```
