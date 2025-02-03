package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project1/internal/service"
	"strconv"
)

type OperationController struct {
	financialOperator *service.FinancialOperator
}

func NewOperationController(fo *service.FinancialOperator) *OperationController {
	return &OperationController{
		financialOperator: fo,
	}
}

func (oc *OperationController) HandleDeposit(c *gin.Context) {
	var req struct {
		Id     int     `json:"id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := oc.financialOperator.Deposit(c, req.Id, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Deposit successful"})
}

func (oc *OperationController) HandleTransfer(c *gin.Context) {
	var req struct {
		FromId int     `json:"from_id"`
		ToId   int     `json:"to_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.ToId == req.FromId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The same sender and reciever"})
		return
	}

	err := oc.financialOperator.Transfer(c, req.FromId, req.ToId, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Transfer successful"})
}

func (oc *OperationController) HandleGetTransactions(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	transactions, err := oc.financialOperator.GetLastTransactions(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
