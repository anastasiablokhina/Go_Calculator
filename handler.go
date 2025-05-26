package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type JRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type JRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JRPCError  `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type JRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleJRPC godoc
// @Summary Обработка JSON-RPC запроса
// @Description Принимает JSON-RPC запрос, выполняет команды и возвращает результат вычислений
// @Accept json
// @Produce json
// @Param request body JRPCRequest true "JSON-RPC запрос"
// @Success 200 {object} JRPCResponse
// @Failure 400 {string} string "Ошибка в формате запроса"
// @Router /jrpc [post]
func handleJRPC(w http.ResponseWriter, r *http.Request) {
	var req JRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJRPCError(w, nil, -32700, "Parse error")
		return
	}

	if req.Method != "calculate" {
		sendJRPCError(w, req.ID, -32601, "Method not found")
		return
	}

	var commands []Command
	if err := json.Unmarshal(req.Params, &commands); err != nil {
		sendJRPCError(w, req.ID, -32602, "Invalid params")
		return
	}

	calculator := NewCalculator(commands)
	output := calculator.Process()

	sendJRPCResponse(w, req.ID, output)
}

func sendJRPCError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JRPCResponse{
		JSONRPC: "2.0",
		Error: &JRPCError{
			Code:    code,
			Message: message,
		},
		ID: id,
	})
}

func sendJRPCResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	})
}

// handleCompute godoc
// @Summary Вычисление выражений
// @Description Принимает список команд и возвращает результат вычислений
// @Accept json
// @Produce json
// @Param commands body []Command true "Список команд"
// @Success 200 {object} Output
// @Failure 400 {string} string "Ошибка в формате запроса"
// @Router /compute [post]
func handleCompute(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var commands []Command
	if err := json.Unmarshal(body, &commands); err != nil {
		http.Error(w, "JSON parsing error: "+err.Error(), http.StatusBadRequest)
		return
	}

	calculator := NewCalculator(commands)
	output := calculator.Process()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
