### 编译

```shell
export GOPATH=/your/go/path/directory  #设置GOPATH路径
cd $GOPATH/src
git clone https://github.com/okblockchainlab/bytom.git ./github.com/bytom
cd ./github.com/bytom
./build.sh #run this script only if you first time build the project
./runbuild.sh
ls *.so
ls *.dylib
```

### 参数说明
**createrawtransaction**:  
输入参数为一个json字符串，格式如下：
{
	"base_transaction":null,
	"actions":[
		{
			"type":"spend_account_unspent_output",
			"utxo":{
        //一个utxo,可以从"list-unspent-outputs"命令中获取
			}
		},
		{
			"type":"spend_account_unspent_output",
			"utxo":{
        //另一个utxo,可以从"list-unspent-outputs"命令中获取
			}
		}
	],
	"ttl":0,
	"time_range":0,
	"xpub": "" //xpub为上面所有utxo所属的账户的xpub，可以从getaddressbyprivate中获取，或从"list-keys"等命令中获取
}`

**getaddressbyprivatekey**:
**signrawtransaction**:
比较简单，查看测试代码即可。
