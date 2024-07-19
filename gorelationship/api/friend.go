package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) friendCount(c *gin.Context) {
	username := c.GetHeader("username")
	cnt, err := s.Service.Friend.Count(c.Request.Context(), map[string]any{
		"name": username,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": cnt,
	})
}

func (s *Server) listFriends(c *gin.Context) {
	username := c.GetHeader("username")
	friends, err := s.Service.Friend.List(c.Request.Context(), map[string]any{
		"name": username,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, friends)
}
