# New Packages

Table of Contents:

- [New Packages](#new-packages)
  - [1. unique package](#1-unique-package)

This section is about the new packages be added.

## 1. unique package

Source: <https://go.dev/blog/unique>

The standard library of Go 1.23 now includes the [new unique package](https://pkg.go.dev/unique). The purposes behind this package is to enable the canonicalization of comparable values. In other words, this package lets you deduplicate values so that they point to a single, canonical, unique copy, while efficiently managing the canonical copies under the hood ([interning](<https://en.wikipedia.org/wiki/Interning_(computer_science)>)).

At high level, interning is very simple:

- Interning is re-using objects of equal value on-demand instead of creating new objects.
- For interning is a common application of interning, where many strings with identical values are needed in the same program. For example, if the name "Kien" appears 100 times, by interning you ensure only one "Kien" is actually allocated memory.

```go
var internPool map[string]string

// Intern returns a string that is equal to s but that may share storage with
// a string previously passed to Intern.
func Intern(s string) string {
    pooled, ok := internPool[s]
    if !ok {
        pooled = strings.Clone(s)
        internPool[pooled] = pooled
    }
    return pooled
}
```

This implementation is super simple and works well enough for some cases, but it has a few problems:

- It never removes strings from the pool.
- It cannot be safely used by multiple goroutines concurrently.
- It only works with strings, even though the idea is quite general.

The new `unique` package introduces a function similar to `Intern` called [Make](https://pkg.go.dev/unique#Make). But it also differs from `Intern` in two important ways:

- It accepts values of any comparable type.
- It returns a wrapper value, a [Handle[T]](https://pkg.go.dev/unique#Handle), from which the canonical value can be retrieved.

A real-world example: Look no further than the `net/netip` package in the standard library, which interns values of type `addrDetail`, part of the [netip.Addr](https://pkg.go.dev/net/netip#Addr) structure.

```go
// Addr represents an IPv4 or IPv6 address (with or without a scoped
// addressing zone), similar to net.IP or net.IPAddr.
type Addr struct {
    // Other irrelevant unexported fields...

    // Details about the address, wrapped up together and canonicalized.
    z unique.Handle[addrDetail]
}

// addrDetail indicates whether the address is IPv4 or IPv6, and if IPv6,
// specifies the zone name for the address.
type addrDetail struct {
    isV6   bool   // IPv4 is false, IPv6 is true.
    zoneV6 string // May be != "" if IsV6 is true.
}

var z6noz = unique.Make(addrDetail{isV6: true})

// WithZone returns an IP that's the same as ip but with the provided
// zone. If zone is empty, the zone is removed. If ip is an IPv4
// address, WithZone is a no-op and returns ip unchanged.
func (ip Addr) WithZone(zone string) Addr {
    if !ip.Is6() {
        return ip
    }
    if zone == "" {
        ip.z = z6noz
        return ip
    }
    ip.z = unique.Make(addrDetail{isV6: true, zoneV6: zone})
    return ip
}
```
