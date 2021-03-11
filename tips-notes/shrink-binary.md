# Shrink Golang binary

- The simple Go code:

```go
package main

import "fmt"

func main() {
   fmt.Println("hello-world")
}
```

- Default build:

```bash
❯ go build -o bin/origin main.go                     testing/.../binary-tricks -> master |• $
❯ ls -la bin/origin                                testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 2063944 Thg 3  11 14:13 bin/origin
❯ ls -lh bin/origin                                testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 2,0M Thg 3  11 14:13 bin/origin
❯                                                  testing/.../binary-tricks -> master ! |• $
```

- Strip the debugging information with [-s and -w linker flags](https://golang.org/cmd/link/):

```bash
❯ go build -o bin/stripped -ldflags="-s -w" main.go
❯ ls -la bin/stripped                              testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 1458176 Thg 3  11 14:15 bin/stripped
❯ ls -lh bin/stripped                              testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 1,4M Thg 3  11 14:15 bin/stripped
```

- Use [UPX](http://upx.sourceforge.net/):

```bash
❯ upx --brute bin/stripped -o bin/upx-stripped   testing/.../binary-tricks -> master ! — |• $
                       Ultimate Packer for eXecutables
                          Copyright (C) 1996 - 2020
UPX 3.96        Markus Oberhumer, Laszlo Molnar & John Reiser   Jan 23rd 2020

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
   1458176 ->    454684   31.18%   linux/amd64   upx-stripped

Packed 1 file.
❯ ls -la bin/upx-stripped                          testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 454684 Thg 3  11 14:15 bin/upx-stripped
❯ ls -lh bin/upx-stripped                          testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 445K Thg 3  11 14:15 bin/upx-stripped
❯                                                  testing/.../binary-tricks -> master ! |• $
❯ upx --brute bin/origin -o bin/upx-origin                                                                                                               testing/.../binary-tricks -> master ! — |• $
                       Ultimate Packer for eXecutables
                          Copyright (C) 1996 - 2020
UPX 3.96        Markus Oberhumer, Laszlo Molnar & John Reiser   Jan 23rd 2020

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
   2063944 ->    959772   46.50%   linux/amd64   upx-origin

Packed 1 file.
❯ ls -la bin/upx-origin                                                                                                                                    testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 959772 Thg 3  11 14:13 bin/upx-origin
❯ ls -lh bin/upx-origin                                                                                                                                    testing/.../binary-tricks -> master ! |• $
-rwxrwxr-x 1 kiennt kiennt 938K Thg 3  11 14:13 bin/upx-origin
❯                                                                                                                                                          testing/.../binary-tricks -> master ! |• $
```

- Conclusion:

```bash
❯ ls -la bin                                                                                                                                               testing/.../binary-tricks -> master ! |• $
total 4836
drwxrwxr-x 2 kiennt kiennt    4096 Thg 3  11 14:16 .
drwxrwxr-x 3 kiennt kiennt    4096 Thg 3  11 14:03 ..
-rwxrwxr-x 1 kiennt kiennt 2063944 Thg 3  11 14:13 origin
-rwxrwxr-x 1 kiennt kiennt 1458176 Thg 3  11 14:15 stripped
-rwxrwxr-x 1 kiennt kiennt  959772 Thg 3  11 14:13 upx-origin
-rwxrwxr-x 1 kiennt kiennt  454684 Thg 3  11 14:15 upx-stripped
```
