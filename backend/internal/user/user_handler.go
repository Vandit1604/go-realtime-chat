package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(ginContext *gin.Context) {
	var userReq CreateUserReq
	if err := ginContext.ShouldBindJSON(&userReq); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// creating user
	res, err := h.Service.CreateUser(ginContext, &userReq)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, res)
}

func (h *Handler) Login(ginContext *gin.Context) {
	var loginUserReq LoginUserReq

	// fills the json data to the structs in golang
	if err := ginContext.ShouldBindJSON(&loginUserReq); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	res, err := h.Service.Login(ginContext, &loginUserReq)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3600 is maxAge for the cookie which is 1 hour
	ginContext.SetCookie("jwt", res.accessToken, 3600, "/", "localhost", false, true)

	resp := &LoginUserRes{
		ID:       res.ID,
		Username: res.Username,
	}

	ginContext.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(ginContext *gin.Context) {
	ginContext.SetCookie("jwt", "", -1, "", "", false, true)
	ginContext.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
