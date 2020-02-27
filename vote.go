package vote

import (
	"net/http"
	"google.golang.org/appengine"
    "goy/back/controller/vote/response"
	"google.golang.org/appengine/datastore"
    "golang.org/x/net/context"
)
var VoteKindName = "Vote"
var VoteId = "1"

type VoteDto struct {
	Id string
    Yes int `datastore:",noindex"`
    No int `datastore:",noindex"`
}

// API: GET /vote/info
func GetVoteInfo(writer http.ResponseWriter, request *http.Request) {
	// 定義api回傳結構
    type Data struct {
		Yes int
		No int
	}
	resp := Data{
		Yes: 0,
		No: 0,
	}
	ctx := appengine.NewContext(request)
	qData, _, err := QueryByID(ctx)
	if err != nil {
        response.Error(writer, ctx, http.StatusInternalServerError, 400, "GetVoteInfo", err)
        return
    }
	if len(qData) == 0 {
		response.Print(writer, http.StatusOK, 100, resp)
		return
	}
	resp.Yes = qData[0].Yes
	resp.No = qData[0].No
	response.Print(writer, http.StatusOK, 100, resp)
}

// API: GET /vote
func GetVote(writer http.ResponseWriter, request *http.Request) {
	// 定義api回傳結構
    type Data struct {
        Result string
	}
	ctx := appengine.NewContext(request)

	// 取得參數 op
	param := request.URL.Query()
	opArr, ok := param["op"]
	op := "no"
	if ok {
		op = opArr[0]
	}
	if op != "yes" {
		op = "no"
	}

	// 啟用 transaction
    err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		qData, keys, err := QueryByID(ctx)
        if err != nil {
            return err
		}
		if len(qData) == 0 {
			emptyDto := VoteDto{
				Id: VoteId,
				Yes: 0,
				No: 0,
			}
			if op == "yes" {
				emptyDto.Yes = 1
			} else {
				emptyDto.No = 1
			}
			Put(ctx, emptyDto)
			return nil
		}
		if op == "yes" {
			qData[0].Yes += 1
		} else {
			qData[0].No += 1
		}
		UpdateWithEntities(ctx, qData, keys) 
        return nil
	}, nil)
	
	// 當資料庫寫入失敗, 回傳error
    if err != nil {
        response.Error(writer, ctx, http.StatusInternalServerError, 400, "GetVote", err)
        return
	}
	
	// 回傳結果
    resp := Data {
        Result: "you vote " + op + " success",
    }
    response.Print(writer, http.StatusOK, 100, resp)
}


// 搜尋資料
func QueryByID(ctx context.Context) ([]VoteDto, []*datastore.Key, error){
    var dst = VoteDto{}
    key := datastore.NewKey(ctx, VoteKindName, VoteId, 0, nil)
	err := datastore.Get(ctx, key, &dst)
	
	if err == datastore.ErrNoSuchEntity {
        return []VoteDto{}, []*datastore.Key{}, nil
    }

    if err != nil {
        return nil, nil, err
    } else {
        return []VoteDto{dst}, []*datastore.Key{key}, nil
    }
}


// 新建一筆資料
func Put(ctx context.Context, voteDto VoteDto) (*datastore.Key, error) {
            
    key := datastore.NewKey(ctx, VoteKindName, VoteId, 0, nil)
    
    if key, err := datastore.Put(ctx, key, &voteDto); err != nil {
        return nil, err
    } else {
        return key, nil
    }
}

// 更新資料
func UpdateWithEntities(ctx context.Context, entities []VoteDto, keys []*datastore.Key) (int, error) {
    if len(entities) > 0 {
        _, err := datastore.PutMulti(ctx, keys, entities)
        if err != nil {
            return 0, err
        }
        return len(entities), nil
    }
    
    return 0, nil
}
