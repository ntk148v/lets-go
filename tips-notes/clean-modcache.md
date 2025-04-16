# go clean -modcache

- The -modcache flag causes clean to remove the entire module download cache, including unpacked source code of versioned dependencies. These packages can be found in the `$GOPATH/pkg/mod` directory.
