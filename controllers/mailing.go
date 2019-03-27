package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"fmt"

	"pinehq/models"
	
	"github.com/jinzhu/gorm"
	"github.com/gin-contrib/sessions"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	gomail "gopkg.in/gomail.v2"
	"crypto/tls"
	"github.com/Sirupsen/logrus"
	
	"net/http"

)



// MAKE DB ORIGINAL
func MakeDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	// Create testdata if not existing
	db.AutoMigrate(&Notification{})
	var notification Notification
	result := db.First(&notification, "slug = ?", "abc")
	if result.Error != nil {
		db.Create(&Notification{Slug: "abc", EmailTo: "testemail3213d@mailinator.com"})
	}

	return db
}





func Mailer() gin.HandlerFunc {
	mailer := gomail.NewDialer(Config.MailServer, Config.MailPort, Config.MailUsername, Config.MailPassword)
	db := MakeDB()
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return func(c *gin.Context) {
		c.Set("mailer", mailer)
		c.Set("db", db)
		c.Next()
	}
}



//EmailEdit handles GET /notifications/:id/edit route
func EmailEdit(c *gin.Context) {
	db := models.GetDB()
	email := models.NotificationPG{}
	db.Preload("EmailTo").First(&email, c.Param("id"))
	if email.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit Email"
	h["Email"] = email
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "notifications/form", h)
}

//EmailUpdate handles POST /notifications/:id/edit route
func EmailUpdate(c *gin.Context) {
	h := DefaultH(c)
	psqldb := models.GetDB()
	email := models.NotificationPG{}
	if err := c.ShouldBind(&email); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/notifications/list/%s/edit", c.Param("id")))
		return
	}
 
	
	session := sessions.Default(c)
   	if err := psqldb.Save(&email).Error; err != nil {
		session.AddFlash("Error in adding email to database.")
		session.Save()
		logrus.Errorf("Error: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}




	var notification Notification

	sqlitedb, _ := c.Keys["db"].(*gorm.DB)
	result := sqlitedb.First(&notification, "slug = ?", "abc")
	if result.Error != nil {
		log.Printf("Some error happend %s", result.Error)
		c.HTML(404, "errors/404", h)
		return
	}





    description := c.PostForm("description")
	m := gomail.NewMessage()
	m.SetHeader("From", "notificationmaster@mg.baun.io")
	m.SetHeader("To", notification.EmailTo)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello, there was a notification: "+ description)

	go func() {
		log.Printf("Sending email to %s", notification.EmailTo)

		d, _ := c.Keys["mailer"].(*gomail.Dialer)

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}()






	c.Redirect(http.StatusFound, "/admin/notifications/list")
}
