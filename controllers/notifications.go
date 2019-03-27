package controllers

import (
	"github.com/gin-gonic/gin"
	"log"


	"pinehq/models"
	
	"github.com/jinzhu/gorm"
	"github.com/gin-contrib/sessions"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	gomail "gopkg.in/gomail.v2"
	"github.com/Sirupsen/logrus"
	"strings"
	"net/http"

)




var Config = models.Configuration{
	MailServer: "localhost",
	MailPort:  5587,
	MailUsername: "test@user.name",
	MailPassword: "passie",
}


//EmailGet handles GET /notifications/list route
func EmailGet(c *gin.Context) {
	db := models.GetDB()
	var notifications []models.NotificationPG
	db.Preload("Description").Find(&notifications)



// SEARCH ONE handles /notifications/:id

	// db.Preload("Description", func(db *gorm.DB) *gorm.DB {
	// 	return db.Order("comments.created_at DESC")
	// }).First(&notification, c.Param("id"))


	h := DefaultH(c)
	h["Title"] = "Notifications you sent!"
	h["Notifications"] = notifications
	c.HTML(http.StatusOK, "notifications/list", h)
}

func ShowMail(c *gin.Context) {
    h := DefaultH(c)
	// slug := c.Param("slug")
	// db := MakeDB()
	// var notification Notification
	// result := db.First(&notification, "slug = ?", slug)

	// if result.Error != nil {
	// 	log.Printf("failed getting %s with error %s", slug, result.Error)
	// 	c.HTML(404, "errors/404", h)
	// 	return
	// }
	
	h["Title"] = "Send an Email!"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	c.HTML(200, "notifications/notification", h)
}



type Notification struct {
	gorm.Model
	Slug    string
	EmailTo string
}


func SendMail(c *gin.Context) {
	// slug := c.Param("slug")
	session := sessions.Default(c)
	h := DefaultH(c)
	
	var notification Notification

	db, _ := c.Keys["db"].(*gorm.DB)
	result := db.First(&notification, "slug = ?", "abc")
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
	m.SetBody("text/html", "Hello, there was a notification: "+description)

	go func() {
		log.Printf("Sending email to %s", notification.EmailTo)

		d, _ := c.Keys["mailer"].(*gomail.Dialer)

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}()

	notificationB := models.NotificationPG{}
    psqldb := models.GetDB()
	if err := c.ShouldBind(&notificationB); err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}


    notificationB.EmailTo = notification.EmailTo
    notificationB.Description = strings.ToLower(description)

   	if err := psqldb.Create(&notificationB).Error; err != nil {
		session.AddFlash("Error whilst adding email to database.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}
// MAKING BACKUP ENDS
	
	h["Title"] = "Notification sent: "+description
	h["Flash"] = session.Flashes()
	c.HTML(200, "notifications/show", h)
}


