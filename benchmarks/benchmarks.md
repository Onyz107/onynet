```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/auth
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkAuth
BenchmarkAuth
BenchmarkAuth-8
     348           3798166 ns/op           11019 B/op        186 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/auth 2.671s
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/crypto
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkEncryptAESGCM/size_64-8                 3209142               372.2 ns/op       171.97 MB/s         112 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_512-8                2086849               583.3 ns/op       877.80 MB/s         592 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_1024-8               1256118              1370 ns/op         747.50 MB/s        1168 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_4096-8                328135              3682 ns/op        1112.51 MB/s        4880 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_16384-8               110348             10320 ns/op        1587.64 MB/s       18448 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_65536-8                31131             39273 ns/op        1668.75 MB/s       73744 B/op          2 allocs/op
BenchmarkEncryptAESGCM/size_262144-8                7550            154362 ns/op        1698.24 MB/s      270353 B/op          2 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/crypto       11.717s
```

```bash
goos: windows       
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/crypto
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkDecryptAESGCM/size_64-8                 8972324               133.0 ns/op       481.25 MB/s           0 B/op          0 allocs/op
BenchmarkDecryptAESGCM/size_512-8                5897886               207.3 ns/op      2470.12 MB/s           0 B/op          0 allocs/op
BenchmarkDecryptAESGCM/size_1024-8               3670878               330.6 ns/op      3096.96 MB/s           0 B/op          0 allocs/op
BenchmarkDecryptAESGCM/size_4096-8               1000000              1004 ns/op        4080.27 MB/s           0 B/op          0 allocs/op
BenchmarkDecryptAESGCM/size_16384-8               150904              7245 ns/op        2261.49 MB/s       16390 B/op          1 allocs/op
BenchmarkDecryptAESGCM/size_65536-8                42973             26702 ns/op        2454.34 MB/s       65558 B/op          1 allocs/op
BenchmarkDecryptAESGCM/size_262144-8               12049             97121 ns/op        2699.15 MB/s      262229 B/op          1 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/crypto       10.438s
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkDial
BenchmarkDial
BenchmarkDial-8
    9447            220321 ns/op           30588 B/op        160 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkNewServer
BenchmarkNewServer
BenchmarkNewServer-8
   10000            190210 ns/op            6102 B/op         44 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkServer_Accept
BenchmarkServer_Accept-8   	    4141	    256245 ns/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkServer_Accept
BenchmarkServer_Accept
BenchmarkServer_Accept-8
    3380            424679 ns/op           58679 B/op        269 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkManager_Open
BenchmarkManager_Open-8   	      22	  50581868 ns/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkStream_Send
BenchmarkStream_Send-8   	   15684	     77384 ns/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkStream_NewStreamedSender
BenchmarkStream_NewStreamedSender-8   	   15816	     76355 ns/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkStream_SendSerialized
BenchmarkStream_SendSerialized-8   	   10000	    157570 ns/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkManager_Accept
BenchmarkManager_Accept
BenchmarkManager_Accept-8
      22          50999477 ns/op           11370 B/op        145 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkManager_Open
BenchmarkManager_Open
BenchmarkManager_Open-8
      22          50745432 ns/op           14745 B/op        146 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_Send
BenchmarkStream_Send
BenchmarkStream_Send-8
   15777             76600 ns/op            2111 B/op         33 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedSender
BenchmarkStream_NewStreamedSender
BenchmarkStream_NewStreamedSender-8
   15742             76836 ns/op            1848 B/op         21 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_SendSerialized
BenchmarkStream_SendSerialized
BenchmarkStream_SendSerialized-8
   10000            157351 ns/op            3122 B/op         65 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_SendEncrypted
BenchmarkStream_SendEncrypted
BenchmarkStream_SendEncrypted-8
   10000            157087 ns/op            9062 B/op         72 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedEncryptedSender
BenchmarkStream_NewStreamedEncryptedSender
BenchmarkStream_NewStreamedEncryptedSender-8
   15721             79072 ns/op            2901 B/op         22 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_Receive
BenchmarkStream_Receive
BenchmarkStream_Receive-8
   10000            100654 ns/op            2212 B/op         34 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedReceiver
BenchmarkStream_NewStreamedReceiver
BenchmarkStream_NewStreamedReceiver-8
   15690             77329 ns/op            1881 B/op         21 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_ReceiveSerialized
BenchmarkStream_ReceiveSerialized
BenchmarkStream_ReceiveSerialized-8
   10000            158961 ns/op            3114 B/op         65 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_ReceiveEncrypted
BenchmarkStream_ReceiveEncrypted
BenchmarkStream_ReceiveEncrypted-8
   10000            157493 ns/op            9257 B/op         72 allocs/op
PASS
```

```bash
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedEncryptedReceiver
BenchmarkStream_NewStreamedEncryptedReceiver
BenchmarkStream_NewStreamedEncryptedReceiver-8
   15700             76887 ns/op            2924 B/op         22 allocs/op
PASS
```