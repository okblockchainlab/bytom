package okwallet

import (
	"testing"
)

func TestGetAddressByPrivateKey(t *testing.T) {
	prvkey1 := `98b94afc04a4dcff39e63bf48d01da8c6057f2564b35def6210add648316584510fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a`
	const expect1 = `{"XPub":"f7bf3b320f828566dff61ec222d4cdbcea032892c8c542b3b6a1c3c725ec356210fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a","Address":"bm1qjvj4hy7tkycp99yuv6u9xjl5n0l0n3uwtzesf6"}`
	prvkey2 := `50024008ed99dfe4680e695fa0467c2960d1a8efbe206d986ba763611288674b9216b09d23f8c96d7ead0309fd53db524194d0a871a0ff718eac85308ffd7415`
	const expect2 = `{"XPub":"8b833f46f53b67e9ee87ca8be70e75da5edf91cfc6c427615ad9d6d96e3ec4e19216b09d23f8c96d7ead0309fd53db524194d0a871a0ff718eac85308ffd7415","Address":"bm1qecew0p8lrftfcc0l8sdppyk07fwk9gpcvnr6e8"}`

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

func TestCreateRawTransaction(t *testing.T) {
	const inputs = `
	{
		"base_transaction":null,
		"actions":[
			{
				"type":"spend_account_unspent_output",
				"utxo":{
					"account_alias": "default",
					"account_id": "0BKBR2D2G0A02",
					"address": "bm1qx7ylnhszg24995d5e0nftu9e87kt9vnxcn633r",
					"amount": 624000000000,
					"asset_alias": "BTM",
					"asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
					"change": false,
					"control_program_index": 12,
					"id": "5af9d3c9b69470983377c1fc0c9125c4ac3bfd32c8d505f2a6042aade8503bc9",
					"program": "00143789f9de0242aa52d1b4cbe695f0b93facb2b266",
					"source_id": "233d1dd49e591980f98e11f333c6c28a867e78448e272011f045131df5aa260b",
					"source_pos": 0,
					"valid_height": 12
				}
			}
		],
		"ttl":0,
		"time_range":0,
		"xpub": "f7bf3b320f828566dff61ec222d4cdbcea032892c8c542b3b6a1c3c725ec356210fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a"
	}`
	const expect = `{"raw_transaction":"07010001016c016a233d1dd49e591980f98e11f333c6c28a867e78448e272011f045131df5aa260bffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c0b1ca9412000121d34d78dfbf3d7fd75ed36e3669ae767756f871b7baf797f46fdddf69c6f66f6eba010000","signing_instructions":[{"position":0,"witness_components":[{"type":"raw_tx_signature","quorum":1,"keys":[{"xpub":"f7bf3b320f828566dff61ec222d4cdbcea032892c8c542b3b6a1c3c725ec356210fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a","derivation_path":["010100000000000000","0c00000000000000"]}],"signatures":null},{"type":"data","value":"a42b4998ff44281e6594c7e83b055bd7324a99f5fa27169f5cb5f92e6bb7d922"}]}],"allow_additional_actions":false}`

	rawTx, err := CreateRawTransaction(inputs)
	if err != nil {
		t.Fatal(err)
	}
	if expect != rawTx {
		t.Fatal("CreateRawTransaction failed. expect " + expect + " but return " + rawTx)
	}
}

func TestSignRawTransaction(t *testing.T) {
	const prvkey = `98b94afc04a4dcff39e63bf48d01da8c6057f2564b35def6210add648316584510fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a`
	const rawTx = `{"raw_transaction":"07010001016c016a233d1dd49e591980f98e11f333c6c28a867e78448e272011f045131df5aa260bffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c0b1ca9412000121d34d78dfbf3d7fd75ed36e3669ae767756f871b7baf797f46fdddf69c6f66f6eba010000","signing_instructions":[{"position":0,"witness_components":[{"type":"raw_tx_signature","quorum":1,"keys":[{"xpub":"f7bf3b320f828566dff61ec222d4cdbcea032892c8c542b3b6a1c3c725ec356210fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a","derivation_path":["010100000000000000","0c00000000000000"]}],"signatures":null},{"type":"data","value":"a42b4998ff44281e6594c7e83b055bd7324a99f5fa27169f5cb5f92e6bb7d922"}]}],"allow_additional_actions":false}`
	const expect = `{"transaction":{"raw_transaction":"07010001016c016a233d1dd49e591980f98e11f333c6c28a867e78448e272011f045131df5aa260bffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c0b1ca9412000121d34d78dfbf3d7fd75ed36e3669ae767756f871b7baf797f46fdddf69c6f66f6eba63024021b0d16ed75ac59b29eafed2a9881694bfdada041683cb21bb31e4942095e9fc646eefe27b455cb97ff312e8248ae318a1b344dab25d262437975536d60dac0d20a42b4998ff44281e6594c7e83b055bd7324a99f5fa27169f5cb5f92e6bb7d92200","signing_instructions":[{"position":0,"witness_components":[{"type":"raw_tx_signature","quorum":1,"keys":[{"xpub":"f7bf3b320f828566dff61ec222d4cdbcea032892c8c542b3b6a1c3c725ec356210fd3d91011aeba9d41f1654c0c8232d1f94eb807a89a35e1501c4effb5c4d5a","derivation_path":["010100000000000000","0c00000000000000"]}],"signatures":["21b0d16ed75ac59b29eafed2a9881694bfdada041683cb21bb31e4942095e9fc646eefe27b455cb97ff312e8248ae318a1b344dab25d262437975536d60dac0d"]},{"type":"data","value":"a42b4998ff44281e6594c7e83b055bd7324a99f5fa27169f5cb5f92e6bb7d922"}]}],"allow_additional_actions":false},"sign_complete":true}`

	resp, err := SignRawTransaction(prvkey, rawTx)
	if err != nil {
		t.Fatal(err)
	}
	if expect != resp {
		t.Fatal("SignRawTransaction failed. expect " + expect + " but return " + resp)
	}
}
