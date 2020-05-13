# frangipanni
Program to convert lines of text into a tree structure.

<img src="frangipanni.jpg" alt="A Tree" width="200">

# Usage

```
cat file.txt | go run frangipanni.go
```

or from binary 

```
cat file.txt | frangipanni
```

# Example

Given input:

```
/var
/boot/vmlinuz-4.9.0-11-amd64
/boot/vmlinuz-4.9.0-8-amd64
/boot/vmlinuz-4.9.0-12-amd64
/boot/initrd.img-4.9.0-12-amd64
/boot/System.map-4.9.0-8-amd64
/boot/config-4.9.0-12-amd64
/boot/System.map-4.9.0-11-amd64
/boot/initrd.img-4.9.0-8-amd64
```

Output:

```
$ <example-input.txt ./frangipanni 

/boot
  /System.map-4.9.0
    -11-amd64
    -8-amd64
  /config-4.9.0-12-amd64
  /initrd.img-4.9.0
    -12-amd64
    -8-amd64
  /vmlinuz-4.9.0
    -11-amd64
    -12-amd64
    -8-amd64
/var
```