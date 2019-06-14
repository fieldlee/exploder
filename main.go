package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    _ "github.com/lib/pq"
    "strconv"
    "time"
)

// refer to: https://www.cnblogs.com/ficow/p/6537238.html

var db *sql.DB
 
func sqlOpen() {
    var err error
    db, err = sql.Open("postgres", "port=5432 user=postgres password=postgres dbname=fabricexplorer sslmode=disable")
    //port是数据库的端口号，默认是5432，如果改了，这里一定要自定义；
    //user就是你数据库的登录帐号;
    //dbname就是你在数据库里面建立的数据库的名字;
    //sslmode就是安全验证模式;
 
    //还可以是这种方式打开
    //db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
    checkErr(err)
}
func sqlInsert(txid string, _type uint8, time, sender, receiver, token string, amount float64) error {
    //插入数据
    stmt, err := db.Prepare("INSERT INTO transactions(txid,type,time,sender,receiver,amount,token) VALUES($1,$2,$3,$4,$5,$6,$7)")
    if err != nil {
        return err
    }
 
    _, err = stmt.Exec(txid, _type, time, sender, receiver, amount, token)
    //这里的三个参数就是对应上面的$1,$2,$3了
 
    return err
}

func sqlSelect() {
    //查询数据
    rows, err := db.Query("SELECT * FROM transactions")
    checkErr(err)
 
    println("-----------")
    for rows.Next() {
	var txid string
        var _type int
        var time,sender,receiver,token string
        var amount float64
        err = rows.Scan(&txid,&_type,&time,&sender,&receiver,&amount,&token)
        checkErr(err)
        fmt.Println("txid = ", txid, "\namount = ", amount, "\ntime = ", time, "\nsender = ", sender, "\n-----------")
    }
}

func sqlClose() {
    db.Close()
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
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
 
 
func main() {
 
    sqlOpen()
    defer sqlClose()

    configFile := "./config.yaml"
    sdk, err := fabsdk.New(config.FromFile(configFile))
    if err != nil {
        fmt.Println("实例化Fabric SDK失败: %v\n", err)
        return
    }

    chan_pvd := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"), fabsdk.WithOrg("Org2"))

    {
        ec, err := event.New(chan_pvd, event.WithBlockEvents(), event.WithSeekType("from"), event.WithBlockNum(10))
        if err != nil {
            fmt.Println("failed to create client: %v\n", err)
            return
        }

        registration, notifier, err := ec.RegisterBlockEvent()
        if err != nil {
            fmt.Println("failed to register block event: %v\n", err)
            return
        }
        defer ec.Unregister(registration)

        fmt.Println("block event registered successfully\n")
        select {
        case ccEvent := <-notifier:
            fmt.Printf("received chaincode event %v\n", ccEvent)
        case <-time.After(time.Second * 50):
            fmt.Println("timeout while waiting for chaincode event\n")
        }
    }

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

        fmt.Println("chaincode event registered successfully\n")
        for ;; {
            select {
            case ccEvent := <-notifier:
                fmt.Printf("received chaincode event %v\n", ccEvent)
                var e ledgerEvent
                err := json.Unmarshal(ccEvent.Payload, &e)
                if err != nil {
                    fmt.Printf( err.Error() )
                    break
                }

                value, convErr := strconv.ParseFloat(e.Amount, 64)
                if convErr != nil {
                    fmt.Printf( convErr.Error() )
                    break
                }

                t := time.Unix(e.Time, 0)

                err = sqlInsert(e.Txid, e.Type, t.Format("2006-01-02 15:04:05"), e.From, e.To, e.Token, value)
                if err != nil {
                    fmt.Printf( err.Error() )
                }
            case <-time.After(time.Second * 5):
                fmt.Println("timeout while waiting for chaincode event\n")
            }
        }
    }
}
