package journals

import (
	"log"
	"strconv"

	"github.com/fine-track/api-app/utils"
	"github.com/gin-gonic/gin"
)

func GetJournalsHandler(c *gin.Context) {
	res, user := getInitialData(c)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	rType := c.Query("type")
	if journals, err := GetJournals(user["_id"], rType, int32(page)); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
	} else {
		res.Message = "Journals found"
		res.Data = journals
		res.Ok(c)
	}
}

func CreateJournalHandler(c *gin.Context) {
	res, _ := getInitialData(c)
	payload := CreateJournalPayload{}
	if err := c.BindJSON(&payload); err != nil {
		log.Printf("unable to unmarshal JSON data from body\nerror: %s\nfile: journals/handlers.go:43", err)
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}
	if record, err := CreateJournal(&payload); err != nil {
		log.Printf("unable to create journal\nerror: %s\nfile: journals/handlers.go:60", err)
		res.Message = err.Error()
		res.BadRequest(c)
	} else {
		if !record.Success {
			res.Message = "Unable to create record"
			res.BadRequest(c)
		} else {
			res.Data = record.Record
			res.Message = "New record created"
			res.Created(c)
		}
	}

}

func UpdateJournalHandler(c *gin.Context) {
	res, _ := getInitialData(c)

	recordId := c.Param("id")
	if recordId == "" {
		res.Message = "Invalid record id"
		res.NotFound(c)
		return
	}

	p := UpdateJournalPayload{}
	if err := c.BindJSON(&p); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}

	if updatedData, err := UpdateJournal(recordId, &p); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
	} else {
		res.Message = "Record updated"
		res.Data = updatedData
		res.Ok(c)
	}
}

func RemoveJournalHandler(c *gin.Context) {
	res, _ := getInitialData(c)
	recordId := c.Param("id")
	if recordId == "" {
		res.Message = "No id specified"
		res.BadRequest(c)
		return
	}
	if err := DeleteJournal(recordId); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
	} else {
		res.Message = "Record removed"
		res.Ok(c)
	}
}

func getInitialData(c *gin.Context) (*utils.HTTPResponse, map[string]string) {
	res := utils.HTTPResponse{}
	u, _ := c.Get("user")
	user := u.(map[string]string)
	return &res, user
}
