package handler

import (
	"backend/internal/domain/usecase"
	"backend/internal/handler/model"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type projectHandler struct {
	projectUsecase usecase.ProjectUsecase
}

// projectHandlerの新しいインスタンスを作成する処理
func NewProjectHandler(projectUsecase usecase.ProjectUsecase) *projectHandler {
	return &projectHandler{
		projectUsecase: projectUsecase,
	}
}

// プロジェクト一覧を取得する処理
func (h *projectHandler) GetProjects(c echo.Context) error {
	ctx := c.Request().Context()
	userUUID, ok := c.Get("user_uuid").(string)
	if !ok {
		slog.WarnContext(ctx, "ユーザーUUIDの取得に失敗")
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  "error",
			Message: "ユーザーUUIDの取得に失敗しました",
		})
	}

	projects, err := h.projectUsecase.GetProjects(ctx, userUUID)
	if err != nil {
		slog.ErrorContext(ctx, "プロジェクト一覧の取得に失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := make([]*model.ProjectResponse, len(projects))
	for i, p := range projects {
		res[i] = &model.ProjectResponse{
			UUID:      p.ID,
			Title:     p.Title,
			UpdatedAt: p.UpdatedAt,
		}
	}

	slog.InfoContext(ctx, "プロジェクト一覧の取得に成功", "user_uuid", userUUID, "count", len(res))
	return c.JSON(http.StatusOK, res)
}
