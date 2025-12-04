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
			UUID:      p.UUID,
			Title:     p.Title,
			UpdatedAt: p.UpdatedAt,
		}
	}

	slog.InfoContext(ctx, "プロジェクト一覧の取得に成功", "user_uuid", userUUID, "count", len(res))
	return c.JSON(http.StatusOK, res)
}

// プロジェクト作成処理
func (h *projectHandler) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()
	userUUID, ok := c.Get("user_uuid").(string)
	if !ok {
		slog.WarnContext(ctx, "ユーザーUUIDの取得に失敗")
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  "error",
			Message: "ユーザーUUIDの取得に失敗しました",
		})
	}

	var req model.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		slog.WarnContext(ctx, "リクエストボディのパースに失敗", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "リクエストボディの形式が正しくありません",
		})
	}

	if req.InitialMessage == "" {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "initial_messageは必須です",
		})
	}

	project, chat, message, err := h.projectUsecase.CreateProject(ctx, userUUID, req.InitialMessage)
	if err != nil {
		slog.ErrorContext(ctx, "プロジェクト作成に失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := model.CreateProjectResponse{
		ProjectUUID: project.UUID,
		ChatUUID:    chat.UUID,
		MessageInfo: model.MessageInfo{
			MessageUUID: message.UUID,
			Message:     message.Content,
		},
		UpdatedAt: project.UpdatedAt,
	}

	slog.InfoContext(ctx, "プロジェクト作成に成功", "project_id", project.UUID)
	return c.JSON(http.StatusCreated, res)
}

// プロジェクトの親チャット取得処理
func (h *projectHandler) GetParentChat(c echo.Context) error {
	ctx := c.Request().Context()
	projectUUID := c.Param("project_uuid")
	if projectUUID == "" {
		slog.WarnContext(ctx, "project_uuidが指定されていません")
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "project_uuidは必須です",
		})
	}

	chat, err := h.projectUsecase.GetParentChat(ctx, projectUUID)
	if err != nil {
		slog.ErrorContext(ctx, "親チャットの取得に失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	res := model.GetParentChatResponse{
		ChatUUID: chat.UUID,
	}

	slog.InfoContext(ctx, "親チャットの取得に成功", "project_uuid", projectUUID, "chat_uuid", chat.UUID)
	return c.JSON(http.StatusOK, res)
}

// プロジェクトツリー取得処理
func (h *projectHandler) GetProjectTree(c echo.Context) error {
	ctx := c.Request().Context()
	projectUUID := c.Param("project_uuid")
	if projectUUID == "" {
		slog.WarnContext(ctx, "project_uuidが指定されていません")
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "project_uuidは必須です",
		})
	}

	tree, err := h.projectUsecase.GetProjectTree(ctx, projectUUID)
	if err != nil {
		slog.ErrorContext(ctx, "プロジェクトツリーの取得に失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	nodes := make([]model.ProjectNode, len(tree.Nodes))
	for i, n := range tree.Nodes {
		nodes[i] = model.ProjectNode{
			ID: n.ID,
			Data: model.ProjectNodeData{
				UserMessage: n.Data.UserMessage,
				Assistant:   n.Data.Assistant,
			},
			Position: model.ProjectNodePosition{
				X: n.Position.X,
				Y: n.Position.Y,
			},
		}
	}

	edges := make([]model.ProjectEdge, len(tree.Edges))
	for i, e := range tree.Edges {
		edges[i] = model.ProjectEdge{
			ID:     e.ID,
			Source: e.Source,
			Target: e.Target,
		}
	}

	res := model.GetProjectTreeResponse{
		Nodes: nodes,
		Edges: edges,
	}

	slog.InfoContext(ctx, "プロジェクトツリーの取得に成功", "project_uuid", projectUUID, "nodes", len(nodes), "edges", len(edges))
	return c.JSON(http.StatusOK, res)
}
