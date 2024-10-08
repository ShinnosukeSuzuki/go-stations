package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP handles HTTP requests and routes them to the appropriate method.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodPost:
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := h.Create(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)

	case http.MethodGet:
		var req model.ReadTODORequest
		// クエリパラメータの取得
		query := r.URL.Query()
		// rからクエリパラメータのprevIDとsizeを取得してreqにセット
		// クエリパラメータがない場合は、reqのPrevIDとSizeに0をセット
		if query.Get("prev_id") == "" {
			req.PrevID = 0
		} else {
			prevID, err := strconv.Atoi(query.Get("prev_id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			req.PrevID = int64(prevID)
		}

		if query.Get("size") == "" {
			// クエリパラメータにsizeがない場合は、defaultで5をセット
			req.Size = 5
		} else {
			size, err := strconv.Atoi(query.Get("size"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			req.Size = int64(size)
		}

		resp, err := h.Read(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// リクエストのidが0または、subjectが空文字の場合はBadRequestを返す
		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		resp, err := h.Update(ctx, &req)
		// ErrNotFound errorの場合は404を返す
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(resp)

	case http.MethodDelete:
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// リクエストのids(配列)が空の場合はBadRequestを返す
		if len(req.IDs) == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		resp, err := h.Delete(ctx, &req)

		if err != nil {
			// ErrNotFound errorの場合は404を返す
			if errors.Is(err, &model.ErrNotFound{}) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	// reqからSubjectとDescriptionを取得してsvc.CreateTODOに渡す
	subject, description := req.Subject, req.Description
	todo, err := h.svc.CreateTODO(ctx, subject, description)
	if err != nil {
		return nil, err
	}
	// 取得したtodoからCreateTODOResponseを生成して返す
	var resp model.CreateTODOResponse
	resp.TODO = *todo
	return &resp, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}

	// []*model.TODO を []model.TODO に変換
	todoList := make([]model.TODO, len(todos))
	for i, todo := range todos {
		todoList[i] = *todo
	}

	var resp model.ReadTODOResponse
	resp.TODOs = todoList
	return &resp, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	id, subject, description := req.ID, req.Subject, req.Description
	todo, err := h.svc.UpdateTODO(ctx, id, subject, description)
	if err != nil {
		return nil, err
	}
	var resp model.UpdateTODOResponse
	resp.TODO = *todo

	return &resp, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	ids := req.IDs
	err := h.svc.DeleteTODO(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
