package controllers

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"task-vix-btpns/models"
	"task-vix-btpns/app"
	"task-vix-btpns/app/auth"
	"task-vix-btpns/helpers/errorformat"
)

//GFunction to get photo profile
func GetPhoto(c *gin.Context) {
	//Create list photo
	photos := []models.Photo{}

	//Set database
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&models.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	//Init list photo
	if len(photos) > 0 {
		for i := range photos {
			user := models.User{}
			err := db.Model(&models.User{}).Where("id = ?", photos[i].UserID).Take(&user).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": err.Error(),
					"data":    nil,
				})
				return
			}

			photos[i].Owner = app.Owner{
				ID: user.ID, Username: user.Username, Email: user.Email,
			}
		}
	}

	//Return response
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Data retrieved successfully",
		"data":    photos,
	})
}

//Function to create photo profile
func CreatePhoto(c *gin.Context) {
	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Get bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "T"})
		return
	}

	//Get user mail from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Get user data from database
	var user_has_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_has_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil,
		})
		return
	}

	// Read body request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	// log.Println(string(body))
	//Convert json to object
	input_photo := models.Photo{}
	err = json.Unmarshal(body, &input_photo)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Init photo
	input_photo.Init()
	input_photo.UserID = user_has_login.ID
	input_photo.Owner = app.Owner{
		ID:       user_has_login.ID,
		Username: user_has_login.Username,
		Email:    user_has_login.Email,
	}
	err = input_photo.Validate("upload") //Validate photo
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Check if photo already exist
	var old_photo models.Photo
	err = db.Debug().Model(&models.Photo{}).Where("user_id = ?", user_has_login.ID).Find(&old_photo).Error
	if err != nil {
		if err.Error() == "Data not found" {
			err = db.Debug().Create(&input_photo).Error //Create photo to database
			if err != nil {
				formattedError := errorformat.ErrorMessage(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": formattedError.Error(),
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  "Success",
				"message": "Photo uploaded successfully",
				"data":    input_photo,
			})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Update photo with new data
	input_photo.ID = old_photo.ID
	err = db.Debug().Model(&old_photo).Updates(&input_photo).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo changed successfully",
		"data":    input_photo,
	}) //Return response

}

//Function to update photo profile
func UpdatePhoto(c *gin.Context) {
	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Get bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Token not found"})
		return
	}

	//Get user mail from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Get user data from JWT
	var user_has_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_has_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil,
		})
		return
	}

	// Read body request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Convert json to object
	photo_input := models.Photo{}
	err = json.Unmarshal(body, &photo_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Validate photo
	err = photo_input.Validate("change")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Check if photo already exist
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo with id " + c.Param("photoId") + " not found",
			"data":    nil,
		})
		return
	}

	//Validate user id
	if user_has_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "You can't change photo of another user",
			"data":    nil,
		})
		return
	}

	//Updating photo to database
	err = db.Model(&photo).Updates(&photo_input).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	photo.Owner = app.Owner{
		ID:       user_has_login.ID,
		Username: user_has_login.Username,
		Email:    user_has_login.Email,
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo updated successfully",
		"data":    photo,
	})
}

//Function to delete photo
func DeletePhoto(c *gin.Context) {

	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Get bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Token not found"})
		return
	}

	//Get user mail from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Get user data from JWT
	var user_has_login models.User
	if err := db.Debug().Where("email = ?", email).First(&user_has_login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil})
		return
	}

	//Check if photo already exist
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	//Validate user id
	if user_has_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "You can't delete photo of another user",
			"data":    nil,
		})
		return
	}

	//Delete photo from database
	err = db.Debug().Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo deleted successfully",
		"data":    nil}) //Return response
}
