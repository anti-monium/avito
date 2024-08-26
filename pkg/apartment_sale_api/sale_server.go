package apartment_sale_api

//lint:file-ignore ST1006 'this' is package-specific

import (
	"avito_bootcamp/pkg/database"
	"avito_bootcamp/pkg/sender"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SaleServer struct {
	db database.IApartmentStorage
}

func NewSaleServer(db database.IApartmentStorage) *SaleServer {
	return &SaleServer{db}
}

func getUserType(c *gin.Context) (userType string) {
	value, exists := c.Get("user_type")
	if exists {
		userType = value.(string)
	} else {
		userType = "client"
	}
	return
}

func (this *SaleServer) GetHouseById(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	userType := getUserType(c)

	flats, err := this.db.GetFlatList(id, userType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, flats)
}

func (this *SaleServer) PostHouseSubscribe(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	type requestInput struct {
		Email string `json:"email" binding:"required"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	this.db.AddSubscriber(id, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (this *SaleServer) PostFlatCreate(c *gin.Context) {
	type requestInput struct {
		HouseId int `json:"house_id" binding:"required"`
		Price   int `json:"price" binding:"required"`
		Rooms   int `json:"rooms" binding:"-"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	flat, err := this.db.CreateFlat(input.HouseId, input.Price, input.Rooms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, flat)
}

func (this *SaleServer) PostHouseCreate(c *gin.Context) {
	type requestInput struct {
		Address   string `json:"address" binding:"required"`
		Year      int    `json:"year" binding:"required"`
		Developer string `json:"developer" binding:"-"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userType := getUserType(c)
	if userType != "moderator" {
		c.Status(http.StatusUnauthorized)
		return
	}

	house, err := this.db.CreateHouse(input.Address, input.Developer, input.Year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, house)
}

func (this *SaleServer) PostFlatUpdate(c *gin.Context) {
	type requestInput struct {
		HouseId int    `json:"house_id" binding:"required"`
		Id      int    `json:"id" binding:"required"`
		Status  string `json:"status" binding:"-"`
	}
	var input requestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userType := getUserType(c)
	if userType != "moderator" {
		c.Status(http.StatusUnauthorized)
		return
	}

	flat, err := this.db.ModerateFlat(input.HouseId, input.Id, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	if flat.Status == database.APPROVED {
		emails, err := this.db.GetSubscribers(input.HouseId)
		if err != nil {
			log.WithError(err).Error("can't send mails to subscribers")
		}
		for _, email := range emails {
			go func() {
				s := sender.New()
				msg := "A new apartment has appeared in house number" + string(flat.HouseId) + "don't miss it!"
				s.SendEmail(context.Background(), email, msg)
			}()
		}
	}
	c.JSON(http.StatusOK, flat)
}
