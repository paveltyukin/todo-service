package handler

import (
	"github.com/gin-gonic/gin"
	"todo/pkg/service"
)

type Handler struct {
	services *service.Service
}

func (h *Handler) InitRouts() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.UserIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.GET("/", h.getAllList)
			lists.GET("/:id", h.getListById)
			lists.POST("/", h.createList)
			lists.PATCH("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			innerItems := lists.Group("/:id/items")
			{
				innerItems.GET("/", h.getAllItem)
				innerItems.POST("/", h.createItem)
			}
		}

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PATCH("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}
