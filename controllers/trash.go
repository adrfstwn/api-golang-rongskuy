package controllers

import (
	"backend-rongskuy/models"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type TrashController struct {
	beego.Controller
}

// Create handles the creation of a new trash post
func (c *TrashController) Create() {
	log.Println("Starting Create method")

	// Log the raw JSON body (for debugging purposes)
	rawBody := c.Ctx.Input.RequestBody
	log.Printf("Raw JSON request body: %s", string(rawBody))

	// Proses Multipart Form untuk menangani unggahan file
	file, header, err := c.GetFile("gambar")
	if err != nil {
		log.Println("Failed to get uploaded file:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Failed to upload file"}
		c.ServeJSON()
		return
	}
	defer file.Close()
	log.Println("Uploaded file received:", header.Filename)

	// Simpan file di folder upload
	uploadPath := "upload/" + header.Filename
	err = c.SaveToFile("gambar", uploadPath)
	if err != nil {
		log.Println("Failed to save file:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to save file"}
		c.ServeJSON()
		return
	}
	log.Println("File saved to:", uploadPath)

	// Membaca data form dari input fields setelah file diunggah
	description := c.GetString("description")
	weight, _ := c.GetFloat("weight")
	latitude, _ := c.GetFloat("latitude")
	longitude, _ := c.GetFloat("longitude")
	whatsapp := c.GetString("whatsapp")

	// Ambil user ID dari context yang sudah disiapkan oleh JWTMiddleware
	userID, ok := c.Ctx.Input.GetData("id").(int64)
	if !ok || userID == 0 {
		log.Println("Failed to retrieve user ID from context")
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not authenticated"}
		c.CustomAbort(401, "Invalid User ID")
		c.ServeJSON()
		return
	}
	log.Println("User authenticated, ID:", userID)

	// Inisialisasi User di dalam model Trash
	trash := models.Trash{
		UserID:      userID,
		Gambar:      uploadPath,
		Description: description,
		Weight:      int64(weight),
		Latitude:    latitude,
		Longitude:   longitude,
		Whatsapp:    whatsapp,
		Status:      "available",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Simpan trash baru ke database
	if id, err := models.AddTrash(&trash); err == nil {
		log.Println("Trash entry created with ID:", id)
		c.Ctx.Output.SetStatus(http.StatusCreated)
		c.Data["json"] = map[string]interface{}{"id": id, "status": "created"}
	} else {
		log.Println("Failed to create trash entry:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
	}
	c.ServeJSON()
}

func (c *TrashController) Edit() {
	// Log the start of the Edit method
	log.Println("Starting Edit method")

	// Parse the multipart form with a max memory of 10MB
	if err := c.Ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error": "Failed to parse form data",
		}
		c.ServeJSON()
		return
	}

	// Retrieve the trash ID from the URL
	trashIDStr := c.Ctx.Input.Param(":id")
	trashID, err := strconv.Atoi(trashIDStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // 400 Bad Request
		c.Data["json"] = map[string]string{
			"error": "Invalid trash ID",
		}
		c.ServeJSON()
		return
	}

	// Retrieve the trash entry from the database
	trash, err := models.GetTrashById(int64(trashID))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound) // 404 Not Found
		c.Data["json"] = map[string]string{
			"error": "Trash entry not found",
		}
		c.ServeJSON()
		return
	}

	// Update fields from the form data
	description := c.GetString("description")
	if description != "" {
		trash.Description = description
	}

	weight, err := c.GetFloat("weight")
	if err == nil {
		trash.Weight = int64(weight)
	}

	latitude, err := c.GetFloat("latitude")
	if err == nil {
		trash.Latitude = latitude
	}

	longitude, err := c.GetFloat("longitude")
	if err == nil {
		trash.Longitude = longitude
	}

	whatsapp := c.GetString("whatsapp")
	if whatsapp != "" {
		trash.Whatsapp = whatsapp
	}

	// Handle file upload
	file, header, err := c.GetFile("gambar")
	if err == nil {
		defer file.Close()

		// Create a unique file name and save the file
		uploadPath := filepath.Join("upload", header.Filename)
		out, err := os.Create(uploadPath)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
			c.Data["json"] = map[string]string{
				"error": "Failed to save uploaded file",
			}
			c.ServeJSON()
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
			c.Data["json"] = map[string]string{
				"error": "Failed to save uploaded file",
			}
			c.ServeJSON()
			return
		}

		// Update the Gambar field
		trash.Gambar = uploadPath
	}

	// Update the updated_at timestamp
	trash.UpdatedAt = time.Now()

	// Save the updated trash entry to the database
	if err := models.UpdateTrashById(trash); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // 500 Internal Server Error
		c.Data["json"] = map[string]string{
			"error":   "Failed to update trash entry",
			"details": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Return success response
	c.Ctx.Output.SetStatus(http.StatusOK) // 200 OK
	c.Data["json"] = map[string]string{
		"success": "Trash entry updated successfully",
	}
	c.ServeJSON()
}

// ReadAll retrieves all trash posts for a specific user
func (c *TrashController) Read() {
	userID, ok := c.Ctx.Input.GetData("id").(int64)
	if !ok || userID == 0 {
		log.Println("Failed to retrieve user ID from context")
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]string{"error": "User not authenticated"}
		c.ServeJSON()
		return
	}

	trashes, err := models.GetTrashesByUserID(userID)
	if err != nil {
		log.Println("Failed to retrieve trashes:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve trashes"}
		c.ServeJSON()
		return
	}

	// Prepare response data with only user ID, excluding 'paid' status
	response := make([]map[string]interface{}, 0)
	for _, trash := range trashes {
		if trash.Status != "paid" { // Exclude trashes with status 'paid'
			response = append(response, map[string]interface{}{
				"Id":          trash.Id,
				"Description": trash.Description,
				"Weight":      trash.Weight,
				"Whatsapp":    trash.Whatsapp,
				"Latitude":    trash.Latitude,
				"Longitude":   trash.Longitude,
				"Gambar":      trash.Gambar,
				"UserId":      trash.UserID,
				"Status":      trash.Status,
				"Created_at":  trash.CreatedAt,
			})
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]interface{}{
		"trashes": response,
	}
	c.ServeJSON()
}

// ReadAll retrieves all trash posts for a specific user
func (c *TrashController) ReadAll() {
	// Fetch all trash posts
	trashes, err := models.GetAllTrash()
	if err != nil {
		log.Println("Failed to retrieve trashes:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to retrieve trashes"}
		c.ServeJSON()
		return
	}

	// Prepare response data, excluding 'paid' status
	response := make([]map[string]interface{}, 0)
	for _, trash := range trashes {
		if trash.Status != "paid" { // Exclude trashes with status 'paid'
			response = append(response, map[string]interface{}{
				"Id":          trash.Id,
				"Description": trash.Description,
				"Weight":      trash.Weight,
				"Whatsapp":    trash.Whatsapp,
				"Latitude":    trash.Latitude,
				"Longitude":   trash.Longitude,
				"Gambar":      trash.Gambar,
				"UserId":      trash.UserID,
				"Status":      trash.Status,
				"Created_at":  trash.CreatedAt,
			})
		}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]interface{}{
		"trashes": response,
	}
	c.ServeJSON()
}

// DeleteTrash deletes a Trash by its Id
func (c *TrashController) Delete() {
	// Get the Id from the URL
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid ID"}
		c.ServeJSON()
		return
	}

	// Call the model function to delete the Trash
	err = models.DeleteTrash(id)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = map[string]string{"error": "Trash not found"}
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = map[string]string{"message": "Trash deleted successfully"}
	}

	c.ServeJSON()
}
