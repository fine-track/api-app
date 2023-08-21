package ledgers

import (
	"os"

	"github.com/fine-track/api-app/pb"
	"github.com/fine-track/api-app/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ledgersService pb.LedgerServiceClient
)

func PrepareConn() {
	addr := os.Getenv("LEDGER_SERVICE_ADDRESS")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.FailOnError(err, "unable to stablish connection with ledgers service", nil)
	defer conn.Close()
	ledgersService = pb.NewLedgerServiceClient(conn)
}

func GetLedgers(userId string) {

}
