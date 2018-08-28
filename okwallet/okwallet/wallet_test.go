package okwallet

import (
  "testing"
)

func TestGetAddressByPrivateKey(t *testing.T) {
  prvkey1 := `98b94afc04a4dcff39e63bf48d01da8c6057f2564b35def6210add648316584510fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a`
  const expect1 = `bm1qjvj4hy7tkycp99yuv6u9xjl5n0l0n3uwtzesf6`
  prvkey2 := `50024008ed99dfe4680e695fa0467c2960d1a8efbe206d986ba763611288674b9216b09d23f8c96d7ead0309fd53db524194d0a871a0ff718eac85308ffd7415`
  const expect2 = `bm1qecew0p8lrftfcc0l8sdppyk07fwk9gpcvnr6e8`

  addr, err := GetAddressByPrivateKey(prvkey1)
  if err != nil {
    t.Fatal(err)
  }
  if expect1 != addr {
    t.Fatal("GetAddressByPrivateKey failed. expect " + expect1 + " but return " + addr)
  }

  addr, err = GetAddressByPrivateKey(prvkey2)
  if err != nil {
    t.Fatal(err)
  }
  if expect2 != addr {
    t.Fatal("GetAddressByPrivateKey failed. expect " + expect2 + " but return " + addr)
  }
}
