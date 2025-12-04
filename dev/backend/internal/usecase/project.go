package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type projectUsecase struct {
	projectRepo repository.ProjectRepository
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
	edgeRepo    repository.EdgeRepository
	txManager   repository.TransactionManager
}

func NewProjectUsecase(
	projectRepo repository.ProjectRepository,
	chatRepo repository.ChatRepository,
	messageRepo repository.MessageRepository,
	edgeRepo repository.EdgeRepository,
	txManager repository.TransactionManager,
) usecase.ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		edgeRepo:    edgeRepo,
		txManager:   txManager,
	}
}

// プロジェクト取得処理
func (u *projectUsecase) GetProjects(ctx context.Context, userUUID string) ([]*model.Project, error) {
	slog.InfoContext(ctx, "プロジェクト一覧取得処理を開始", "user_uuid", userUUID)
	projects, err := u.projectRepo.FindAllByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の取得に失敗: %w", err)
	}
	slog.InfoContext(ctx, "プロジェクト一覧取得処理を完了", "user_uuid", userUUID)
	return projects, nil
}

// プロジェクト作成処理
func (u *projectUsecase) CreateProject(ctx context.Context, userUUID, initialMessage string) (*model.Project, *model.Chat, *model.Message, error) {
	slog.InfoContext(ctx, "プロジェクト作成処理を開始", "user_uuid", userUUID)

	projectID := uuid.New().String()
	chatID := uuid.New().String()
	messageID := uuid.New().String()
	now := time.Now()

	project := &model.Project{
		UUID:      projectID,
		UserUUID:  userUUID,
		Title:     initialMessage, // タイトルは最初のメッセージとする（要件によるが一旦）
		CreatedAt: now,
		UpdatedAt: now,
	}

	chat := &model.Chat{
		UUID:        chatID,
		ProjectUUID: projectID,
		Title:       initialMessage, // チャットのタイトルも最初のメッセージとする
		Status:      "open",         // 初期状態はopen
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	message := &model.Message{
		UUID:      messageID,
		ChatUUID:  chatID,
		Role:      "user",
		Content:   initialMessage,
		CreatedAt: now,
	}

	err := u.txManager.Do(ctx, func(ctx context.Context) error {
		if err := u.projectRepo.Create(ctx, project); err != nil {
			return fmt.Errorf("プロジェクト作成失敗: %w", err)
		}
		if err := u.chatRepo.Create(ctx, chat); err != nil {
			return fmt.Errorf("チャット作成失敗: %w", err)
		}
		if err := u.messageRepo.Create(ctx, message); err != nil {
			return fmt.Errorf("メッセージ作成失敗: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	slog.InfoContext(ctx, "プロジェクト作成処理を完了", "project_id", projectID)
	return project, chat, message, nil
}

// プロジェクトの親チャット取得処理
func (u *projectUsecase) GetParentChat(ctx context.Context, projectUUID string) (*model.Chat, error) {
	slog.InfoContext(ctx, "プロジェクトの親チャット取得処理を開始", "project_uuid", projectUUID)
	chat, err := u.chatRepo.FindOldestByProjectUUID(ctx, projectUUID)
	if err != nil {
		return nil, fmt.Errorf("プロジェクトの親チャット取得に失敗: %w", err)
	}
	slog.InfoContext(ctx, "プロジェクトの親チャット取得処理を完了", "project_uuid", projectUUID)
	return chat, nil
}

// プロジェクトツリー取得処理
func (u *projectUsecase) GetProjectTree(ctx context.Context, projectUUID string) (*model.ProjectTree, error) {
	slog.InfoContext(ctx, "プロジェクトツリー取得処理を開始", "project_uuid", projectUUID)

	// 1. プロジェクトのルートチャットを取得
	rootChat, err := u.chatRepo.FindOldestByProjectUUID(ctx, projectUUID)
	if err != nil {
		return nil, fmt.Errorf("ルートチャットの取得に失敗: %w", err)
	}

	nodes := []model.ProjectNode{}
	edges := []model.ProjectEdge{}
	visitedChats := make(map[string]bool)
	queue := []string{rootChat.UUID}

	for len(queue) > 0 {
		currentChatUUID := queue[0]
		queue = queue[1:]

		if visitedChats[currentChatUUID] {
			continue
		}
		visitedChats[currentChatUUID] = true

		// 2. チャットのメッセージを取得
		messages, err := u.messageRepo.FindMessagesByChatID(ctx, currentChatUUID)
		if err != nil {
			return nil, fmt.Errorf("メッセージの取得に失敗 (chat_uuid: %s): %w", currentChatUUID, err)
		}

		// 3. メッセージからノードを作成
		for i := 0; i < len(messages); i++ {
			msg := messages[i]
			var node *model.ProjectNode

			// ユーザーメッセージの場合、次のメッセージがアシスタントならペアにする
			if msg.Role == "user" {
				if i+1 < len(messages) && messages[i+1].Role == "assistant" {
					assistantMsg := messages[i+1]
					userContent := msg.Content
					node = &model.ProjectNode{
						ID: assistantMsg.UUID,
						Data: model.ProjectNodeData{
							UserMessage: &userContent,
							Assistant:   assistantMsg.Content,
						},
						Position: model.ProjectNodePosition{
							X: assistantMsg.PositionX,
							Y: assistantMsg.PositionY,
						},
					}

					// アシスタントメッセージのフォークも確認
					for _, fork := range assistantMsg.Forks {
						if !visitedChats[fork.ChatUUID] {
							queue = append(queue, fork.ChatUUID)
						}
					}

					i++ // アシスタントメッセージをスキップ
				}
			} else if msg.Role == "assistant" {
				// アシスタントメッセージ単体の場合
				node = &model.ProjectNode{
					ID: msg.UUID,
					Data: model.ProjectNodeData{
						UserMessage: nil,
						Assistant:   msg.Content,
					},
					Position: model.ProjectNodePosition{
						X: msg.PositionX,
						Y: msg.PositionY,
					},
				}
			}

			if node != nil {
				nodes = append(nodes, *node)
			}

			// 子チャットがあればキューに追加
			for _, fork := range msg.Forks {
				if !visitedChats[fork.ChatUUID] {
					queue = append(queue, fork.ChatUUID)
				}
			}
		}

		// 4. チャットのエッジを取得
		chatEdges, err := u.edgeRepo.FindEdgesByChatID(ctx, currentChatUUID)
		if err != nil {
			return nil, fmt.Errorf("エッジの取得に失敗 (chat_uuid: %s): %w", currentChatUUID, err)
		}

		for _, edge := range chatEdges {
			edges = append(edges, model.ProjectEdge{
				ID:     edge.UUID,
				Source: edge.SourceMessageUUID,
				Target: edge.TargetMessageUUID,
			})
		}
	}

	slog.InfoContext(ctx, "プロジェクトツリー取得処理を完了", "project_uuid", projectUUID)
	return &model.ProjectTree{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
