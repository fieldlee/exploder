package main

// sudo docker pull mysql:5.6.35
// docker run --name mysql -v /home/czp/SVN/ST/explorer/MySQL/db:/initdb -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:5.6.35
// docker exec -it mysql bash

// code sample: https://blog.csdn.net/m1179457922/article/details/80797480
import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"strconv"
	"time"
)
//下面填写自己的数据库信息，看不懂英文？这都看不懂还学啥编程。不知道数据库信息？那还用啥数据库。
var (
	dbhostip="192.168.1.100"
	dbusername="root"
	dbpassword="123456"
	dbname="mmchannel"
)

func checkErr(err error){
	if err!=nil{
		panic(err)
	}
}

var db *sql.DB

func sqlOpen() {
	var err error
	db, err = sql.Open("mysql",dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8")
	checkErr(err)
}

func sqlClose() {
	db.Close()
}

func sqlInsert(txid string, _type uint8, date int32, time, sender, receiver, token string, amount float64) error {
	stmt, err := db.Prepare("INSERT INTO transactions(txid,type,date,time,sender,receiver,amount,token) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(txid, _type, date, time, sender, receiver, amount, token)

	return err
}



type ledgerEvent struct {
	Type   uint8  `json:"type"`
	Txid   string `json:"txid"`
	Time   int64  `json:"time"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Token  string `json:"token"`
}

func onEvent(ccEvent *fab.CCEvent) {
	var e ledgerEvent
	err := json.Unmarshal(ccEvent.Payload, &e)
	if err != nil {
		fmt.Printf( err.Error() )
		return
	}

	value, convErr := strconv.ParseFloat(e.Amount, 64)
	if convErr != nil {
		fmt.Printf( convErr.Error() )
		return
	}

	t := time.Unix(e.Time, 0)

	y,m,d:=t.Date()

	err = sqlInsert(e.Txid, e.Type, int32(y*10000+int(m)*100+d), t.Format("2006-01-02 15:04:05"), e.From, e.To, e.Token, value)
	if err != nil {
		fmt.Printf( err.Error() )
	}
}


func main(){
	sqlOpen()
	defer sqlClose()

	configFile := "/var/yaml/config_event.yaml"
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		fmt.Println("实例化Fabric SDK失败: %v\n", err)
		return
	}

	chan_pvd := sdk.ChannelContext("mmchannel", fabsdk.WithUser("Admin"), fabsdk.WithOrg("mmOrg"))

	fmt.Println("=====ChaincodeEvent=======\n")
	{
		ec, err := event.New(chan_pvd, event.WithBlockEvents(), event.WithSeekType("from"), event.WithBlockNum(10))
		if err != nil {
			fmt.Println("failed to create client: %v\n", err)
			return
		}

		registration, notifier, err := ec.RegisterChaincodeEvent("ledger", `LEDGER_TX_[^ ]+`)
		if err != nil {
			fmt.Println("failed to register chaincode event: %v\n", err)
			return
		}
		defer ec.Unregister(registration)

		registration, notifier2, err := ec.RegisterChaincodeEvent("payment", `[0-9a-f]{64}`)
		if err != nil {
			fmt.Println("failed to register chaincode event: %v\n", err)
			return
		}
		defer ec.Unregister(registration)

		fmt.Println("chaincode event registered successfully\n")
		for ;; {
			select {
			case ccEvent := <-notifier:
				fmt.Printf("received ledger event %v\n", ccEvent)
				onEvent(ccEvent)
			case ccEvent := <-notifier2:
				fmt.Printf("received payment event %v\n", ccEvent)
				onEvent(ccEvent)
			case <-time.After(time.Second * 600):
				fmt.Println("timeout while waiting for chaincode event\n")
			}
		}
	}
}
