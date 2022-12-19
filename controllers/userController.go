package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"task-vix-btpns/models"
	"task-vix-btpns/app"
	"task-vix-btpns/app/auth"
	"task-vix-btpns/helpers/errorformat"
	"task-vix-btpns/helpers/hash"
)

//Function to be used for user login
func Login(c *gin.Context) {
	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Read body form
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Convert json to object
	user_model := models.User{}
	err = json.Unmarshal(body, &user_model)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Init user
	user_model.Init()
	err = user_model.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Check if user exist
	var user_login app.UserLogin

	err = db.Debug().Table("users").Select("*").Joins("LEFT JOIN photos ON photos.user_id = users.id").
		Where("users.email = ?", user_model.Email).Find(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + user_model.Email + " not found",
			"data":    nil,
		})
		return
	}

	//Verify password
	err = hash.CheckPasswordHash(user_login.Password, user_model.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	//Generate token when success login
	token, err := auth.GenerateJWT(user_login.Email, user_login.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserData{
		ID: user_login.ID, Username: user_login.Username, Email: user_login.Email, Token: token,
		Photos: app.Photo{Title: user_login.Title, Caption: user_login.Caption, PhotoUrl: user_login.PhotoUrl},
	}

	//Return response
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "Success",
		"message": "Login successfully",
		"data":    data,
	})
}

//Function to register user
func CreateUser(c *gin.Context) {
	//set database
	db := c.MustGet("db").(*gorm.DB)

	// read body form
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Convert json to object
	user_model := models.User{}
	err = json.Unmarshal(body, &user_model)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user_model.Init() //Inisialize user

	err = user_model.Validate("update") //Validate user
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = user_model.HashPassword() //Hash password
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Create(&user_model).Error //Create user to database
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserRegister{ //data to be used for response
		ID:        user_model.ID,
		Username:  user_model.Username,
		Email:     user_model.Email,
		CreatedAt: user_model.CreatedAt,
		UpdatedAt: user_model.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User registered succesfully",
		"data":    data,
	}) //Response success
}

//Function to update user
func UpdateUser(c *gin.Context) {

	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Check if user exist
	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + c.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	//Read body form
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//Convert json to object
	user_model := models.User{}

	user_model.ID = user.ID
	err = json.Unmarshal(body, &user_model)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Validate user
	err = user_model.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Hashing password
	err = user_model.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	//Update user
	err = db.Debug().Model(&user).Updates(&user_model).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserRegister{ //data to be used for response
		ID:        user_model.ID,
		Username:  user_model.Username,
		Email:     user_model.Email,
		CreatedAt: user_model.CreatedAt,
		UpdatedAt: user_model.UpdatedAt,
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Error",
		"message": "User updated succesfully",
		"data":    data,
	})
}

//Function to delete user
func DeleteUser(c *gin.Context) {

	//Set database
	db := c.MustGet("db").(*gorm.DB)

	//Check if user exist
	var user models.User

	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + c.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	//Delete user
	err = db.Debug().Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Response success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User deleted succesfully",
		"data":    nil,
	})
}
