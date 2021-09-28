package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	db "github.com/amallick86/Go_bank/db/sqlc"
	"github.com/amallick86/Go_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type transferPspRequest struct {
	AccountID    int64  `json:"from_account_id" binding:"required,min=1"`
	PspAccountID int64  `json:"pspAccount_id" binding:"required,min=1"`
	Type         string `json:"type" binding:"required"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Currency     string `json:"currency" binding:"required,currency"`
}
type Transaction struct {
	Citizenship string `json:"citizenship"`
	From        string `json:"from"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
}

func (server *Server) createTransferPsp(ctx *gin.Context) {
	var tran Transaction
	var req transferPspRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccountUser(ctx, req.AccountID, req.Currency, req.Amount)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	valid = server.validBankAccount(ctx, req.PspAccountID, req.Currency, authorizationPayloadKey)
	if !valid {
		return
	}

	arg := db.TransferPspTxParams{
		AccountID:    req.AccountID,
		PspAccountID: req.PspAccountID,
		Amount:       req.Amount,
	}
	marshalRequest, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	var bearer = "Bearer " + authorizationPayloadKey
	//resp, err := http.Post("http://127.0.0.1:9000/transfer/psp/add", "application/json", bytes.NewBuffer(marshalRequest))

	resp, err := http.NewRequest("POST", "http://127.0.0.1:8080/transfer/psp/add", bytes.NewBuffer(marshalRequest))
	resp.Header.Add("Authorization", bearer)
	resp.Header.Add("Accept", "application/json")
	if err != nil {
		logrus.Error(err)
	}
	client := &http.Client{}
	response, err := client.Do(resp)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {

		result, err := server.store.TransferPspTx(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		tran.Citizenship = fromAccount.Citizenship
		tran.Type = req.Type
		tran.From = "Bank"
		tran.Amount = req.Amount
		marshalRequests, err := json.Marshal(tran)
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
		resp, err := http.NewRequest("POST", "http://localhost:8081/transaction", bytes.NewBuffer(marshalRequests))

		if err != nil {
			logrus.Error(err)
		}
		client := &http.Client{}
		response, err := client.Do(resp)
		if err != nil {
			fmt.Println("HTTP call failed:", err)
		}
		defer response.Body.Close()
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			ctx.JSON(http.StatusOK, result)
		} else {
			ctx.JSON(http.StatusOK, "file not saved")
		}
	} else {

		ctx.JSON(http.StatusInternalServerError, "")
	}
}

func (server *Server) validAccountUser(ctx *gin.Context, accountID int64, currency string, amount int64) (db.Account, bool) {

	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	if account.Balance < amount {
		return account, false
	}

	return account, true
}

func (server *Server) validBankAccount(ctx *gin.Context, accountID int64, currency string, token string) bool {
	var req struct {
		AccountID int64 `json:"id" binding:"required,min=1"`
	}
	req.AccountID = accountID
	marshalRequest, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var bearer = "Bearer " + token
	//resp, err := http.Post("http://127.0.0.1:9000/accounts/id", "application/json", bytes.NewBuffer(marshalRequest))
	resp, err := http.NewRequest("POST", "http://:8080/accounts/id", bytes.NewBuffer(marshalRequest))
	// ...

	resp.Header.Add("Authorization", bearer)
	resp.Header.Add("Accept", "application/json")
	if err != nil {
		logrus.Error(err)
	}
	client := &http.Client{}
	response, err := client.Do(resp)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusUnauthorized, response.Body)
		return false
	}
	return false
}

type addtransferPspRequest struct {
	AccountID    int64  `json:"bankAccountID" binding:"required,min=1"`
	PspAccountID int64  `json:"accountID" binding:"required,min=1"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Type         string `json:"type" binding:"required"`
	Currency     string `json:"currency" binding:"required,currency"`
}

func (server *Server) addTransferPsp(ctx *gin.Context) {
	var req addtransferPspRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ReceivePspTxParams{
		AccountID:    req.AccountID,
		PspAccountID: req.PspAccountID,
		Amount:       req.Amount,
	}

	result, err := server.store.ReceivePspTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
