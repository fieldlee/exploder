module explorer

go 1.12

require (
        github.com/go-sql-driver/mysql v1.4.1 // indirect
        github.com/hyperledger/fabric-sdk-go v1.0.0-beta1 // indirect
        google.golang.org/appengine v1.4.0 // indirect
)

replace (
        golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191029031824-8986dd9e96cf
        golang.org/x/net => github.com/golang/net v0.0.0-20191028085509-fe3aa8a45271
        golang.org/x/sys => github.com/golang/sys v0.0.0-20191029155521-f43be2a4598c
)
