package handler

import (
	"user_service/internal/entity"
	"user_service/internal/interface"
    "net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService _interface.UserService
	authService _interface.AuthSerice
}

func NewUserHandler(userService _interface.UserService, authService _interface.AuthSerice) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}


// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(c *gin.Context) {
    var req entity.RegsisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
        return
    }

    user := &entity.User{
        ID:    uuid.New(),
        Name:  req.Name,
        Email: req.Email,
        Role:  req.Role,
		Password: req.Password,
    }

    if err := h.userService.RegisterUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusCreated, entity.SuccessResponse{
        Message: "User registered successfully",
        Data:    user,
    })
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req entity.LoginRequest

	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		return
	}

	
	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	
	authenticatedUser, err := h.userService.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(authenticatedUser)
	if err != nil {
		// Nếu token không tạo được, trả về lỗi 500 (Internal Server Error)
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, entity.LoginResponse{
		Token: token,
		User:  authenticatedUser,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return
	}

	id, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.userService.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "User profile retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, ok := userIDRaw.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	var req entity.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Profile updated successfully",
		Data:    user,
	})
}

func (h *UserHandler) ToggleUserLockByAdmin(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	err = h.userService.ToggleUserLockByAdmin(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{Message: "User locked successfully"})
}


func (h *UserHandler) GetAllUsersByAdmin(c *gin.Context) {
	users, err := h.userService.GetAllUsersByAdmin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Error: "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    users,
	})
}