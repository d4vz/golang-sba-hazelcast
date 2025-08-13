package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "golang-sba-hazelcast/internal/auction"
)

type API struct {
    svc *auction.Service
}

func NewAPI(svc *auction.Service) *API {
    return &API{svc: svc}
}

func (a *API) Routes(mux *http.ServeMux) {
    mux.HandleFunc("/auctions", a.createAuction)
    mux.HandleFunc("/auctions/active", a.activeAuctions)
    mux.HandleFunc("/work", a.work)
}

// createAuction godoc
// @Summary      Create an auction
// @Description  Creates a new auction
// @Accept       json
// @Produce      json
// @Param        auction body auction.Auction true "Auction"
// @Success      201
// @Failure      400
// @Router       /auctions [post]
func (a *API) createAuction(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var in auction.Auction
    dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
    dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    if err := a.svc.CreateAuction(ctx, in); err != nil {
        http.Error(w, "cannot create", http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

// activeAuctions godoc
// @Summary      List active auctions
// @Produce      json
// @Success      200 {array} auction.Auction
// @Failure      500
// @Router       /auctions/active [get]
func (a *API) activeAuctions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    list, err := a.svc.ActiveAuctions(ctx, time.Now())
    if err != nil {
        http.Error(w, "cannot fetch", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(list)
}

// work burns CPU for a specified duration (ms) to help HPA-based autoscaling tests.
// @Summary      Consume CPU cycles for autoscaling tests
// @Param        ms query int false "Duration in milliseconds" default(50)
// @Success      200
// @Router       /work [get]
func (a *API) work(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    ms := 50
    if q := r.URL.Query().Get("ms"); q != "" {
        if v, err := time.ParseDuration(q + "ms"); err == nil {
            ms = int(v / time.Millisecond)
        }
    }
    deadline := time.Now().Add(time.Duration(ms) * time.Millisecond)
    for time.Now().Before(deadline) {
        select {
        case <-r.Context().Done():
            http.Error(w, "canceled", http.StatusRequestTimeout)
            return
        default:
        }
        // Busy loop to consume CPU cycles intentionally.
    }
    w.WriteHeader(http.StatusOK)
}


