# First Go Web Application Tutorial

This is from the [tutorial here](https://golang.org/doc/articles/wiki/).

To run:

```
go run wiki.go
```

Then visit [localhost:8080](http://localhost:8080)

## Learnings

* Function wrappers seem like a decent way to create abstractions/middleware in Go.  This is encouraging since roughly the same approach is used in JS.
* Might not need much in the way of frameworks, standard lib gets you mostly there.
* Keep services small, orchestrate with Kubernetes.
* Gorilla has some useful tools that can be used piecemeal.
* Could potentially [build Go docker containers that are under 10MB](https://rollout.io/blog/building-minimal-docker-containers-for-go-applications/).

