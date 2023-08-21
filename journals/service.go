package journals

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fine-track/api-app/pb"
	"github.com/fine-track/api-app/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateJournalPayload struct {
	UserId      string `json:"user_id"`
	Type        string `json:"type"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int32  `json:"amount"`
}

type UpdateJournalPayload struct {
	Type        string `json:"type"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Amount      int32  `json:"amount"`
	Description string `json:"description"`
}

var (
	journalsService pb.RecordsServiceClient
)

func PrepareConn() *grpc.ClientConn {
	addr := os.Getenv("JOURNALS_SERVICE_ADDRESS")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.FailOnError(err, "unable to stablish connection with journals service", nil)

	journalsService = pb.NewRecordsServiceClient(conn)

	// test the connection by pinning the journals service
	if _, err := journalsService.Ping(context.TODO(), &pb.PingRequest{Message: "Ping"}); err != nil {
		utils.FailOnError(err, "unable to ping journals service", nil)
	}

	return conn
}

func GetJournals(userId string, recordType string, page int32) ([]*pb.Record, error) {
	rt := utils.InlineConditional(strings.ToLower(recordType) == "income", pb.RecordType_INCOME, pb.RecordType_EXPENSE)
	payload := &pb.GetRecordsRequest{
		Type:   rt,
		Page:   page,
		UserId: userId,
	}
	results, err := journalsService.GetRecords(context.TODO(), payload)
	if err != nil {
		return nil, err
	}
	fmt.Println(results)
	if !results.Success {
		return nil, errors.New(results.Message)
	}
	if results.Records == nil {
		d := []*pb.Record{}
		return d, nil
	} else {
		return results.Records, nil
	}
}

func CreateJournal(p *CreateJournalPayload) (*pb.UpdateRecordResponse, error) {
	if p.Type != "EXPENSE" && p.Type != "INCOME" {
		return nil, fmt.Errorf("'%s' isn't a valid type", p.Type)
	}
	pbRecordType := utils.InlineConditional(p.Type == "EXPENSE", pb.RecordType_EXPENSE, pb.RecordType_INCOME)
	record, err := journalsService.Create(context.TODO(), &pb.CreateRecordRequest{
		Date:        p.Date, // YYYY-MM-DD
		Amount:      p.Amount,
		Title:       p.Title,
		Description: p.Description,
		UserId:      p.UserId,
		Type:        pbRecordType,
	})
	return record, err
}

func UpdateJournal(id string, p *UpdateJournalPayload) (*pb.Record, error) {
	if p.Type != "EXPENSE" && p.Type != "INCOME" {
		return nil, fmt.Errorf("'%s' isn't a valid type", p.Type)
	}
	pbRecordType := utils.InlineConditional(p.Type == "EXPENSE", pb.RecordType_EXPENSE, pb.RecordType_INCOME)
	record, err := journalsService.Update(context.TODO(), &pb.Record{
		Id:          id,
		Type:        pbRecordType,
		Amount:      p.Amount,
		Title:       p.Title,
		Date:        p.Date,
		Description: p.Description,
	})
	return record.Record, err
}

func DeleteJournal(id string) error {
	res, err := journalsService.Delete(context.TODO(), &pb.DeleteRecordRequest{Id: id})
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Message)
	}
	return nil
}
