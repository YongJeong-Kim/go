package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) friendAccept(c *gin.Context) {
	var req struct {
		FromRequest string `json:"from_request"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("user")
	err := s.Service.Friend.Accept(c.Request.Context(), req.FromRequest, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) friendCount(c *gin.Context) {
	userID := c.GetHeader("user")
	cnt, err := s.Service.Friend.Count(c.Request.Context(), userID)
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
	userID := c.GetHeader("user")
	friends, err := s.Service.Friend.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (s *Server) mutualFriends(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("user")
	friends, err := s.Service.Friend.ListMutuals(c.Request.Context(), userID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, friends)
}

func (s *Server) listFriendRequests(c *gin.Context) {
	userID := c.GetHeader("user")
	reqs, err := s.Service.Friend.ListRequests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reqs)
}

func (s *Server) mutualFriendCount(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("user")
	count, err := s.Service.Friend.MutualCount(c.Request.Context(), userID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (s *Server) friendRequest(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("user")
	err := s.Service.Friend.Request(c.Request.Context(), userID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (s *Server) friendRequestCount(c *gin.Context) {
	userID := c.GetHeader("user")
	count, err := s.Service.Friend.RequestCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
