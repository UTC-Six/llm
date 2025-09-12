package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/UTC-Six/llm/internal/eino/chain"
	"github.com/UTC-Six/llm/internal/eino/graph"
	"github.com/UTC-Six/llm/internal/eino/workflow"
	"github.com/UTC-Six/llm/internal/svc"
	"github.com/UTC-Six/llm/internal/types"
)

// EinoHandler Eino框架处理器
type EinoHandler struct {
	svcCtx               *svc.ServiceContext
	chainOrchestrator    *chain.ChainOrchestrator
	graphOrchestrator    *graph.GraphOrchestrator
	workflowOrchestrator *workflow.WorkflowOrchestrator
}

// NewEinoHandler 创建Eino处理器
func NewEinoHandler(svcCtx *svc.ServiceContext) *EinoHandler {
	return &EinoHandler{
		svcCtx:               svcCtx,
		chainOrchestrator:    chain.NewChainOrchestrator(svcCtx.Config),
		graphOrchestrator:    graph.NewGraphOrchestrator(svcCtx.Config),
		workflowOrchestrator: workflow.NewWorkflowOrchestrator(svcCtx.Config),
	}
}

// ChainInvoke Chain同步聊天接口
func (h *EinoHandler) ChainInvoke(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Chain请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Chain消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Chain编排聊天请求: %s", req.Message)

	resp, err := h.chainOrchestrator.ChatInvoke(r.Context(), &req)
	if err != nil {
		log.Printf("处理eino Chain编排聊天请求失败: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, types.ErrorResponse{
			Code:    500,
			Message: "处理请求失败: " + err.Error(),
		})
		return
	}

	log.Printf("eino Chain编排聊天请求处理成功，返回内容: %s", resp.Data)
	writeJSONResponse(w, http.StatusOK, resp)
}

// ChainStream Chain流式聊天接口
func (h *EinoHandler) ChainStream(w http.ResponseWriter, r *http.Request) {
	var req types.StreamChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Chain流式请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Chain流式消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Chain编排流式聊天请求: %s", req.Message)

	// 设置流式响应的HTTP头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	w.WriteHeader(http.StatusOK)

	err := h.chainOrchestrator.ChatStream(r.Context(), &req, w)
	if err != nil {
		log.Printf("处理eino Chain编排流式聊天请求失败: %v", err)
		io.WriteString(w, "data: {\"code\":500,\"message\":\"处理请求失败: "+err.Error()+"\",\"done\":true}\n\n")
		return
	}

	log.Printf("eino Chain编排流式聊天请求处理完成")
}

// GraphInvoke Graph同步聊天接口
func (h *EinoHandler) GraphInvoke(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Graph请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Graph消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Graph编排聊天请求: %s", req.Message)

	resp, err := h.graphOrchestrator.ChatInvoke(r.Context(), &req)
	if err != nil {
		log.Printf("处理eino Graph编排聊天请求失败: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, types.ErrorResponse{
			Code:    500,
			Message: "处理请求失败: " + err.Error(),
		})
		return
	}

	log.Printf("eino Graph编排聊天请求处理成功，返回内容: %s", resp.Data)
	writeJSONResponse(w, http.StatusOK, resp)
}

// GraphStream Graph流式聊天接口
func (h *EinoHandler) GraphStream(w http.ResponseWriter, r *http.Request) {
	var req types.StreamChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Graph流式请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Graph流式消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Graph编排流式聊天请求: %s", req.Message)

	// 设置流式响应的HTTP头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	w.WriteHeader(http.StatusOK)

	err := h.graphOrchestrator.ChatStream(r.Context(), &req, w)
	if err != nil {
		log.Printf("处理eino Graph编排流式聊天请求失败: %v", err)
		io.WriteString(w, "data: {\"code\":500,\"message\":\"处理请求失败: "+err.Error()+"\",\"done\":true}\n\n")
		return
	}

	log.Printf("eino Graph编排流式聊天请求处理完成")
}

// WorkflowInvoke Workflow同步聊天接口
func (h *EinoHandler) WorkflowInvoke(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Workflow请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Workflow消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Workflow编排聊天请求: %s", req.Message)

	resp, err := h.workflowOrchestrator.ChatInvoke(r.Context(), &req)
	if err != nil {
		log.Printf("处理eino Workflow编排聊天请求失败: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, types.ErrorResponse{
			Code:    500,
			Message: "处理请求失败: " + err.Error(),
		})
		return
	}

	log.Printf("eino Workflow编排聊天请求处理成功，返回内容: %s", resp.Data)
	writeJSONResponse(w, http.StatusOK, resp)
}

// WorkflowStream Workflow流式聊天接口
func (h *EinoHandler) WorkflowStream(w http.ResponseWriter, r *http.Request) {
	var req types.StreamChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("解析Workflow流式请求参数失败: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if req.Message == "" {
		log.Printf("Workflow流式消息内容不能为空")
		writeJSONResponse(w, http.StatusBadRequest, types.ErrorResponse{
			Code:    400,
			Message: "消息内容不能为空",
		})
		return
	}

	log.Printf("收到eino Workflow编排流式聊天请求: %s", req.Message)

	// 设置流式响应的HTTP头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	w.WriteHeader(http.StatusOK)

	err := h.workflowOrchestrator.ChatStream(r.Context(), &req, w)
	if err != nil {
		log.Printf("处理eino Workflow编排流式聊天请求失败: %v", err)
		io.WriteString(w, "data: {\"code\":500,\"message\":\"处理请求失败: "+err.Error()+"\",\"done\":true}\n\n")
		return
	}

	log.Printf("eino Workflow编排流式聊天请求处理完成")
}
