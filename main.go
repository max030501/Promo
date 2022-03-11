package main

import (
	"awesomeProject1/participant"
	"awesomeProject1/prize"
	"strings"

	_ "awesomeProject1/docs"
	"awesomeProject1/promo"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files" // swagger embed files
	"github.com/swaggo/gin-swagger"        // gin-swagger middleware
	"net/http"
)

type storeServer struct {
	promo  *promo.PromoStore
}

// createPromoHandler godoc
// @Summary Create new promo
// @Accept json
// @Param json_body body promo.Promo true "json body"
// @Success 200 {string} string "id"
// @Failure 400 {string} string "Field validation failed on the 'required' tag"
// @Router /promo [post]
func (s storeServer) createPromoHandler(c *gin.Context) {
	var pr promo.Promo
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id := s.promo.CreatePromo(pr.Name,pr.Description)
	c.String(http.StatusOK,strconv.Itoa(id))
}

// getAllPromoHandler godoc
// @Summary Retrieves all promos
// @Produce json
// @Success 200 {object} []promo.Promo
// @Router /promo [get]
func (s storeServer) getAllPromoHandler(c *gin.Context) {
	promos := s.promo.GetPromos()
	c.JSON(http.StatusOK, promos)
}

// getPromoHandler godoc
// @Summary Retrieves promo by id
// @Produce json
// @Param id path integer true "Promo ID"
// @Success 200 {object} promo.ResponsePromo
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id} [get]
func (s storeServer) getPromoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	promo, err := s.promo.GetPromo(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, promo)
}

// putPromoHandler godoc
// @Summary Edit promo
// @Accept json
// @Param id path integer true "Promo ID"
// @Param json_body body promo.RequestPromo true "json body"
// @Success 200
// @Failure 400 {string} string "Field validation failed on the 'required' tag"
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id} [put]
func (s storeServer) putPromoHandler(c *gin.Context) {
	var rp promo.RequestPromo
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ShouldBindJSON(&rp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

		err = s.promo.PutPromo(id, rp.Name, rp.Description)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}
// deletePromoHandler godoc
// @Summary Delete promo
// @Param id path integer true "Promo ID"
// @Success 200
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id} [delete]
func (s storeServer) deletePromoHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = s.promo.DeletePromo(id); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}
// createParticipantHandler godoc
// @Summary Create new participant
// @Accept json
// @Param id path integer true "Promo ID"
// @Param json_body body participant.Participant true "json body"
// @Success 200 {string} string "id"
// @Failure 400 {string} string "Field validation failed on the 'required' tag"
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id}/participant [post]
func (s storeServer) createParticipantHandler(c *gin.Context) {
	var pr participant.Participant
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	idPart, err := s.promo.CreateParticipant(id, pr.Name)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK,strconv.Itoa(idPart))
}

// deleteParticipantHandler godoc
// @Summary Delete participant
// @Param id path integer true "Promo ID"
// @Param partId path integer true "Participant ID"
// @Success 200
// @Failure 400 {string} string "There is no promo with id"
// @Failure 400 {string} string "There is no participant with id"
// @Router /promo/{id}/participant/{partId} [delete]
func (s storeServer) deleteParticipantHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	partId, err := strconv.Atoi(c.Param("partId"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = s.promo.DeleteParticipant(id,partId); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}

// createPrizeHandler godoc
// @Summary Create new prize
// @Accept json
// @Param id path integer true "Promo ID"
// @Param json_body body prize.Prize true "json body"
// @Success 200 {string} string "id"
// @Failure 400 {string} string "Field validation failed on the 'required' tag"
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id}/prize [post]
func (s storeServer) createPrizeHandler(c *gin.Context) {
	var pr prize.Prize
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	idPrize, err := s.promo.CreatePrize(id, pr.Description)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK,strconv.Itoa(idPrize))
}

// deletePrizeHandler godoc
// @Summary Delete prize
// @Param id path integer true "Promo ID"
// @Param prizeId path integer true "Prize ID"
// @Success 200
// @Failure 400 {string} string "There is no promo with id"
// @Failure 400 {string} string "There is no prize with id"
// @Router /promo/{id}/prize/{prizeId} [delete]
func (s storeServer) deletePrizeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	prizeId, err := strconv.Atoi(c.Param("prizeId"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = s.promo.DeletePrize(id,prizeId); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
}

// raffleHander godoc
// @Summary Do raffle
// @Param id path integer true "Promo ID"
// @Produce json
// @Success 200 {object} []promo.RaffleResponse
// @Failure 409 {string} string "Count of participants and prizes is not equal or empty"
// @Failure 400 {string} string "There is no promo with id"
// @Router /promo/{id} [post]
func (s storeServer) raffleHander(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	raf,err := s.promo.Raffle(id)
	if err != nil {
		if strings.Contains(err.Error(),"equal"){
			c.String(http.StatusConflict, err.Error())
		}else {
			c.String(http.StatusBadRequest, err.Error())
		}
	}
	c.JSON(http.StatusOK, raf)

}


func NewStoreServer() *storeServer {
	promoStore := promo.New()
	return &storeServer{promo: promoStore}
}


// @title Products_Store
// @version 1.0

// @BasePath /

func main() {
	router := gin.Default()
	server := NewStoreServer()

	router.POST("/promo", server.createPromoHandler)
	router.GET("/promo", server.getAllPromoHandler)
	router.GET("/promo/:id", server.getPromoHandler)
	router.PUT("/promo/:id", server.putPromoHandler)
	router.DELETE("/promo/:id", server.deletePromoHandler)
	router.POST("/promo/:id/participant", server.createParticipantHandler)
	router.DELETE("/promo/:id/participant/:partId", server.deleteParticipantHandler)
	router.POST("/promo/:id/prize", server.createPrizeHandler)
	router.DELETE("/promo/:id/prize/:prizeId", server.deletePrizeHandler)
	router.POST("/promo/:id/raffle", server.raffleHander)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("localhost:8080")
}
