package controllers

import (
	"backend-rongskuy/models"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

const (
	charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength = 12 // Adjust length as needed
)

type TransactionController struct {
	beego.Controller
}

// GenerateRandomCode generates a random code with numbers and uppercase letters
func GenerateRandomCode(length int) (string, error) {
	var result []byte
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result = append(result, charset[num.Int64()])
	}
	return string(result), nil
}

func (c *TransactionController) BuyTrash() {
	log.Println("Starting BuyTrash method")

	// Read and parse the JSON request body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Failed to read request body"}
		c.ServeJSON()
		return
	}

	var transaction models.Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	// Retrieve user ID from the context
	userID, ok := c.Ctx.Input.GetData("id").(int64)
	if !ok || userID == 0 {
		log.Println("Failed to retrieve user ID from context")
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not authenticated"}
		c.ServeJSON()
		return
	}
	transaction.UserID = userID
	transaction.CreatedAt = time.Now() // Set creation time
	transaction.UpdatedAt = time.Now() // Set update time

	code, err := GenerateRandomCode(codeLength)
	if err != nil {
		log.Println("Failed to generate transaction code:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to generate transaction code"}
		c.ServeJSON()
		return
	}
	transaction.KodeTransaksi = code

	// Set the description and status for the transaction
	transaction.Keterangan = "Pembelian sampah"
	transaction.Status = "paid"

	adminID := int64(1)
	if _, err := models.GetUserModelById(adminID); err != nil {
		log.Println("Failed to retrieve admin user:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve admin user"}
		c.ServeJSON()
		return
	}

	// Retrieve trash details to get the weight
	trash, err := models.GetTrashById(transaction.TrashID)
	if err != nil {
		log.Println("Failed to retrieve trash:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve trash"}
		c.ServeJSON()
		return
	}

	// Check if the trash is already paid
	if trash.Status == "paid" {
		log.Println("Trash already purchased")
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = map[string]string{"error": "Trash already purchased"}
		c.ServeJSON()
		return
	}

	// Assume 10 coins per unit weight
	coinsPerWeight := int64(10)
	transaction.CoinsEarned = trash.Weight * coinsPerWeight

	// Update the user's coins
	if err := models.UpdateUserCoins(trash.UserID, transaction.CoinsEarned); err != nil {
		log.Println("Failed to update admin coins:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to update user coins"}
		c.ServeJSON()
		return
	}

	// Update the admin's coins
	if err := models.UpdateUserCoins(userID, -transaction.CoinsEarned); err != nil {
		log.Println("Failed to update user coins:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to update admin coins"}
		c.ServeJSON()
		return
	}

	// Save the transaction to the database
	if id, err := models.AddTransaction(&transaction); err == nil {
		// Update trash status to 'paid'
		trashID := transaction.TrashID // Assuming TrashID is part of the transaction
		if err := models.UpdateTrashStatus(trashID, "paid"); err != nil {
			log.Println("Failed to update trash status:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Failed to update trash status"}
			c.ServeJSON()
			return
		}

		log.Println("Transaction created with ID:", id)
		c.Ctx.Output.SetStatus(http.StatusCreated)
		c.Data["json"] = map[string]interface{}{"id": id, "status": "created"}
	} else {
		log.Println("Failed to create transaction:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
	}
	c.ServeJSON()
}

// ReadUserTransactions handles the retrieval of transactions based on the user_id in trash
func (c *TransactionController) Read() {
	log.Println("Starting ReadUserTransactions method")

	// Retrieve user ID from the context or from some authentication token
	userID, ok := c.Ctx.Input.GetData("id").(int64)
	if !ok || userID == 0 {
		log.Println("Failed to retrieve user ID from context")
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not authenticated"}
		c.ServeJSON()
		return
	}

	// Retrieve trash entries for the user
	trashes, err := models.GetTrashesByUserID(userID)
	if err != nil {
		log.Println("Failed to retrieve trashes:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve trashes"}
		c.ServeJSON()
		return
	}

	// Retrieve transaction details based on trash IDs
	var transactions []models.Transaction
	for _, trash := range trashes {
		trx, err := models.GetTransactionsByTrashID(trash.Id)
		if err != nil {
			log.Println("Failed to retrieve transaction for trash ID:", trash.Id, err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Failed to retrieve transactions"}
			c.ServeJSON()
			return
		}
		transactions = append(transactions, trx...)
	}

	// Return the transactions as JSON
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = transactions
	c.ServeJSON()
}

// ReadAdmin handles the retrieval of transactions for a specific user based on user_id in transactions
func (c *TransactionController) ReadAdmin() {
	log.Println("Starting ReadAdmin method")

	// Retrieve user ID from the context or from some authentication token
	userID, ok := c.Ctx.Input.GetData("id").(int64)
	if !ok || userID == 0 {
		log.Println("Failed to retrieve user ID from context")
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not authenticated"}
		c.ServeJSON()
		return
	}

	// Retrieve transactions for the user directly from the transactions table
	transactions, err := models.GetTransactionsByUserID(userID)
	if err != nil {
		log.Println("Failed to retrieve transactions:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve transactions"}
		c.ServeJSON()
		return
	}

	// Return the transactions as JSON
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = transactions
	c.ServeJSON()
}
