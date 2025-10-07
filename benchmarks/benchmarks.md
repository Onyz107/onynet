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
BenchmarkEncryptAESGCM/size_64-8         	 1567856	       873.4 ns/op	  73.27 MB/s	    1392 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_512-8        	  863097	      1323 ns/op	 387.07 MB/s	    1872 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_1024-8       	  725403	      1685 ns/op	 607.81 MB/s	    2448 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_4096-8       	  341286	      3627 ns/op	1129.23 MB/s	    6160 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_16384-8      	  109008	     10781 ns/op	1519.71 MB/s	   19728 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_65536-8      	   30702	     39107 ns/op	1675.81 MB/s	   75024 B/op	       4 allocs/op
BenchmarkEncryptAESGCM/size_262144-8     	    8428	    143324 ns/op	1829.03 MB/s	  271632 B/op	       4 allocs/op
BenchmarkDecryptAESGCM/size_64-8         	 1436178	       836.5 ns/op	  76.51 MB/s	    1344 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_512-8        	 1000000	      1112 ns/op	 460.39 MB/s	    1792 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_1024-8       	  811314	      1471 ns/op	 696.33 MB/s	    2304 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_4096-8       	  304518	      3479 ns/op	1177.34 MB/s	    5376 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_16384-8      	  101857	     10607 ns/op	1544.67 MB/s	   17664 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_65536-8      	   31417	     38024 ns/op	1723.55 MB/s	   66816 B/op	       3 allocs/op
BenchmarkDecryptAESGCM/size_262144-8     	    8800	    145891 ns/op	1796.84 MB/s	  263424 B/op	       3 allocs/op
PASS
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