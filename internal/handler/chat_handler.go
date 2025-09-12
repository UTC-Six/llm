package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/UTC-Six/llm/internal/svc"
	"github.com/UTC-Six/llm/internal/types"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	svcCtx *svc.ServiceContext
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(svcCtx *svc.ServiceContext) *ChatHandler {
	return &ChatHandler{
		svcCtx: svcCtx,
	}
}

// ChatInvoke 同步聊天接口
// 处理传统的请求-响应模式的聊天请求
func (h *ChatHandler) ChatInvoke(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证请求参数
	if req.Message == "" {
		log.Printf("消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到聊天请求: %s", req.Message)

	// 调用业务逻辑
	resp, err := h.svcCtx.LLMService.ChatInvoke(r.Context(), &req)
	if err != nil {
		log.Printf("处理聊天请求失败: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, types.ErrorResponse{
			Code:    500,
			Message: "处理请求失败: " + err.Error(),
		})
		return
	}

	log.Printf("聊天请求处理成功，返回内容: %s", resp.Data)
	writeJSONResponse(w, http.StatusOK, resp)
}

// ChatStream 流式聊天接口
// 处理流式响应的聊天请求
func (h *ChatHandler) ChatStream(w http.ResponseWriter, r *http.Request) {
	var req types.StreamChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析流式请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证请求参数
	if req.Message == "" {
		log.Printf("消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到流式聊天请求: %s", req.Message)

	// 设置流式响应的HTTP头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 发送初始响应
	w.WriteHeader(http.StatusOK)

	// 调用流式业务逻辑
	err := h.svcCtx.LLMService.ChatStream(r.Context(), &req, w)
	if err != nil {
		log.Printf("处理流式聊天请求失败: %v", err)
		// 发送错误响应
		io.WriteString(w, "data: {\"code\":500,\"message\":\"处理请求失败: "+err.Error()+"\",\"done\":true}\n\n")
		return
	}

	log.Printf("流式聊天请求处理完成")
}

// writeJSONResponse 写入JSON响应
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
