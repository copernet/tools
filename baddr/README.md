baddr
---
baddr is a convenient tool for generating bitcoin address according to the specified network. At the current version, the generated address is suitable for bitcoin core and bitcoin cash use.

Installation:
```Go
go get github.com/qshuai/Tools/baddr
```

Usage:
```
baddr testnet

// possible output:
privkey key:            cQKU3Xe2Z4KZKfMaaq4NYBsvmNDtkyAoHqCtcunm3ARG4iVkr39X
base58 encoded address: n2vUKbgiXoXLKRvC5Qs4J9M84cp38NZd8Q
bech32 encoded address: bchtest:qr4v63s836w29d0emvp983xff6gkmmvvevjrmqh78z
```
