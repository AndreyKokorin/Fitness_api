package users

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/models"
	"awesomeProject/pkg/validate"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

func UpdateUser(ctx *gin.Context) {
	var userid, _ = strconv.Atoi(ctx.Param("id"))

	var updateData models.UpdateUser
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("User request Failed")
		return
	}

	err := validate.ValidAndTrim(&updateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("User valid Failed")
		return
	}

	var userWithNewData models.User
	query := "UPDATE users SET name = COALESCE(NULLIF($1, ''), name), email = COALESCE(NULLIF($2, ''), email), age = COALESCE(NULLIF($3, 0), age) WHERE id = $4 RETURNING id, name, email, age;"
	err = database.DB.QueryRow(query, updateData.Name, updateData.Email, updateData.Age, userid).Scan(&userWithNewData.Id, &userWithNewData.Name, &userWithNewData.Email, &userWithNewData.Age)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("User sql Failed")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": userWithNewData})
}
