# Tips and tricks for writing unit test in Golang

[Source](https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742)

## Put your tests in a different package

Moving your test code out of the package allows you to write tests as though youwere a real user of the package.

## Internal tests go in a different files

## Run all tests on save

Go builds and runs very quickly, sp there's little to no reason why you shouldn't run your entire test suite every time you hit save. While you're at it, why not run go vet, golint and goimports at the same time.

## Write table driven tests

Anonymous structs and composite literals allow us to write every clear and simple table tests without relying on any external package.

## Mock things using Go code

If you need to mock something that your code relies on in order to properly test it, chances are it is a good candidate for an interface. Even if you're relying on an external package that you cannot change, your code can still take interface that the external types will satisfy.
