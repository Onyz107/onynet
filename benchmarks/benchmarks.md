## OnyNet Benchmarks

---

### Authentication
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/auth
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkAuth
BenchmarkAuth
BenchmarkAuth-8
     388           3002299 ns/op           10998 B/op        185 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/auth 2.449s
```

---

### Crypto

#### Encryption
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/crypto
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkEncryptAESGCM/size_64-8         	 3505194	       361.3 ns/op	 177.16 MB/s	     112 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_512-8        	 2007733	       569.5 ns/op	 898.99 MB/s	     592 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_1024-8       	 1408666	       856.6 ns/op	1195.40 MB/s	    1168 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_4096-8       	  520514	      2374 ns/op	1725.66 MB/s	    4880 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_16384-8      	  152104	      7843 ns/op	2088.98 MB/s	   18448 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_65536-8      	   39111	     30585 ns/op	2142.77 MB/s	   73744 B/op	       2 allocs/op
BenchmarkEncryptAESGCM/size_262144-8     	    9202	    137108 ns/op	1911.95 MB/s	  270352 B/op	       2 allocs/op
PASS
ok  	github.com/Onyz107/onynet/internal/crypto	21.224s
```

#### Decryption
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/crypto
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
BenchmarkDecryptAESGCM/size_64-8         	 6682563	       162.4 ns/op	 394.14 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecryptAESGCM/size_512-8        	 4858365	       252.6 ns/op	2027.14 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecryptAESGCM/size_1024-8       	 3135314	       389.7 ns/op	2627.85 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecryptAESGCM/size_4096-8       	  986744	      1211 ns/op	3383.37 MB/s	       0 B/op	       0 allocs/op
BenchmarkDecryptAESGCM/size_16384-8      	  136388	      9959 ns/op	1645.10 MB/s	   16390 B/op	       1 allocs/op
BenchmarkDecryptAESGCM/size_65536-8      	   29734	     39399 ns/op	1663.41 MB/s	   65560 B/op	       1 allocs/op
BenchmarkDecryptAESGCM/size_262144-8     	    9742	    146200 ns/op	1793.06 MB/s	  262238 B/op	       1 allocs/op
PASS
ok  	github.com/Onyz107/onynet/internal/crypto	21.224s
```

---

### KCP

#### Dial
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkDial
BenchmarkDial
BenchmarkDial-8
   10000            134718 ns/op           29867 B/op        134 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/kcp  2.213s
```

#### NewServer
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkNewServer
BenchmarkNewServer
BenchmarkNewServer-8
   15054             89510 ns/op            6112 B/op         44 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/kcp  2.394s
```

#### Server.Accept
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/kcp
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkServer_Accept
BenchmarkServer_Accept
BenchmarkServer_Accept-8
    4278            301513 ns/op           58295 B/op        252 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/kcp  1.883s
```

---

### SMUX

#### manager.Accept
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkManager_Accept
BenchmarkManager_Accept
BenchmarkManager_Accept-8
      22          50754545 ns/op           14883 B/op        144 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.639s
```

#### manager.Open
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkManager_Open
BenchmarkManager_Open
BenchmarkManager_Open-8
      22          50708905 ns/op           14634 B/op        145 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.675s
```

#### Send
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_Send
BenchmarkStream_Send
BenchmarkStream_Send-8
   15115             79286 ns/op            2148 B/op         34 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.804
```

#### NewStreamedSender
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedSender
BenchmarkStream_NewStreamedSender
BenchmarkStream_NewStreamedSender-8
   15789             76561 ns/op            1853 B/op         21 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.828s
```

#### SendSerialized
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_SendSerialized
BenchmarkStream_SendSerialized
BenchmarkStream_SendSerialized-8
   10000            169124 ns/op            3085 B/op         65 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 2.278s
```

#### SendEncrypted
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_SendEncrypted
BenchmarkStream_SendEncrypted
BenchmarkStream_SendEncrypted-8
   10000            157199 ns/op            5407 B/op         67 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 2.165s
```

#### NewStreamedEncryptedSender
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedEncryptedSender
BenchmarkStream_NewStreamedEncryptedSender
BenchmarkStream_NewStreamedEncryptedSender-8
   15718             76623 ns/op            2897 B/op         22 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.789s
```

#### Receive
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_Receive
BenchmarkStream_Receive
BenchmarkStream_Receive-8
   15729             76933 ns/op            2124 B/op         33 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.780s
```

#### NewStreamedReceiver
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedReceiver
BenchmarkStream_NewStreamedReceiver
BenchmarkStream_NewStreamedReceiver-8
   15760             76615 ns/op            1872 B/op         21 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.827s
```

#### ReceiveSerialized
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_ReceiveSerialized
BenchmarkStream_ReceiveSerialized
BenchmarkStream_ReceiveSerialized-8
   10000            168521 ns/op            3103 B/op         65 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 2.318s
```

#### ReceiveEncrypted
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_ReceiveEncrypted
BenchmarkStream_ReceiveEncrypted
BenchmarkStream_ReceiveEncrypted-8
   10000            157413 ns/op            5426 B/op         67 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 2.169s
```

#### NewStreamedEncryptedReceiver
```shell
goos: windows
goarch: amd64
pkg: github.com/Onyz107/onynet/internal/smux
cpu: Intel(R) Core(TM) i7-8665U CPU @ 1.90GHz
=== RUN   BenchmarkStream_NewStreamedEncryptedReceiver
BenchmarkStream_NewStreamedEncryptedReceiver
BenchmarkStream_NewStreamedEncryptedReceiver-8
   15698             77034 ns/op            2928 B/op         22 allocs/op
PASS
ok      github.com/Onyz107/onynet/internal/smux 1.807s
```