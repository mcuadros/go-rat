rat - <i>random access tar</i>
=======================
[![Build Status](https://travis-ci.org/mcuadros/go-rat.png?branch=master)](https://travis-ci.org/mcuadros/go-rat) [![Coverage Status](https://coveralls.io/repos/mcuadros/go-rat/badge.svg?branch=master)](https://coveralls.io/r/mcuadros/go-rat?branch=master) [![GoDoc](http://godoc.org/github.com/mcuadros/go-rat?status.png)](http://godoc.org/github.com/mcuadros/go-rat) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-rat.svg)](https://github.com/mcuadros/go-rat/releases)


rat is an extension to the classical tar archive, focused on allowing constant-time random file access with linear memory consumption increase. <b>t</b>ape <b>ar</b>chive, was originally developed to write and read streamed sources, making random access to the content very inefficient.

Based on the [benchmarks](#benchmarks), we found that rat is **4x** to **60x** times faster over SSD and HDD than the classic tar file, when reading a single file from a tar archive.

> Any tar file produced by **rat** is compatible with standard tar implementation.

Installation
------------

The recommended way to install rat

```
go get -u github.com/mcuadros/go-rat/...
```

Example
-------

Import the package:

```go
import "github.com/mcuadros/go-rat"
```

Converting a standard `tar` file to a `rat` file:
```go
src, _ := os.Open("standard.file.tar")
dst, _ := os.Create("extended.rat.file.tar")
defer src.Close()
defer dst.Close()

if err = AddIndexToTar(src, dst); err != nil {
    panic(err)
}
```

Searching a specific file in a `rat` file:

```go
archive, _ := os.Open("extended.rat.file.tar")

content, _ := archive.ReadFile("foo.txt")
fmt.Println(string(content))
//Prints: foo
```


<a name="benchmarks"></a>Benchmarks Results
----------------------

These are some of the benchmarking results over differrent storage systems.
> Fixture name explanation: `5_1.0KB_102KB`  means a tar containing 5 files with a size between 1kb and 102kb.

| SSD              | TAR (ns)   | RAT (ns)  | times  |
|------------------|------------|-----------|--------|
| 5_1.0KB_102KB    | 367838     | 77236     | 4.76   |
| 100_1.0KB_102KB  | 5925036    | 350116    | 16.92  |
| 1000_1.0KB_102KB | 58735369   | 3503317   | 16.77  |
| 6000_1.0KB_102KB | 349484665  | 20064072  | 17.42  |
| 60_1.0MB_21MB    | 146302392  | 3402651   | 43.00  |


| HDD              | TAR (ns)   | RAT (ns)  | times  |
|------------------|------------|-----------|--------|
| 5_1.0KB_102KB    | 253406     | 54472     | 4.65   |
| 100_1.0KB_102KB  | 3682796    | 282085    | 13.06  |
| 1000_1.0KB_102KB | 37834628   | 2396239   | 15.79  |
| 6000_1.0KB_102KB | 210841382  | 13913158  | 15.15  |
| 60_1.0MB_21MB    | 166405959  | 2783659   | 59.78  |


| GlusterFS        | TAR (ns)   | RAT (ns)  | times  |
|------------------|------------|-----------|--------|
| 5_1.0KB_102KB    | 293252     | 130652    | 2.24   |
| 100_1.0KB_102KB  | 4292723    | 362399    | 11.85  |
| 1000_1.0KB_102KB | 39632581   | 4468976   | 8.87   |
| 6000_1.0KB_102KB | 2413057504 | 16586371  | 145.48 |
| 60_1.0MB_21MB    | 623461320  | 112529704 | 5.54   |

License
-------

MIT, see [LICENSE](LICENSE)
