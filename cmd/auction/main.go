package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    "golang-sba-hazelcast/cmd/auction/handler"
    "golang-sba-hazelcast/internal/auction"
    docs "golang-sba-hazelcast/internal/docs"
    ihz "golang-sba-hazelcast/internal/platform/hazelcast"
    "golang-sba-hazelcast/internal/platform/testdouble"
)

// @title           Auction SBA API
// @version         0.1.0
// @description     Online auction API built with Space Based Architecture using Hazelcast.
// @BasePath        /
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    clusterName := getenv("HZ_CLUSTERNAME", "auction-cluster")
    members := getenv("HZ_MEMBERS", "127.0.0.1:5701")
    httpAddr := getenv("APP_HTTP_ADDR", ":8080")

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    srv := &http.Server{
        Addr:              httpAddr,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
    }

    go func() {
        log.Printf("auction app listening on %s", httpAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("http server error: %v", err)
        }
    }()

    var space ihz.Space
    if os.Getenv("HZ_SKIP") == "1" {
        space = testdouble.NewFakeSpace()
        log.Printf("running with fake space (HZ_SKIP=1)")
    } else {
        client, err := ihz.New(ctx, clusterName, splitAndTrim(members))
        if err != nil {
            log.Fatalf("failed to start hazelcast client: %v", err)
        }
        defer func() { _ = client.Shutdown(ctx) }()
        space = client
        log.Printf("hazelcast connected: cluster=%s members=%s", clusterName, members)
    }

    svc := auction.NewService(space)
    api := handler.NewAPI(svc)
    api.Routes(mux)
    // Serve swagger JSON and UI from the same server
    docs.Handler(mux)
    handler.SwaggerUI(mux)

    // Block main until context timeout or server error (simulate run-forever)
    <-ctx.Done()
}

func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func splitAndTrim(s string) []string {
    out := []string{}
    cur := ""
    for i := 0; i < len(s); i++ {
        if s[i] == ',' {
            if cur != "" {
                out = append(out, cur)
                cur = ""
            }
            continue
        }
        cur += string(s[i])
    }
    if cur != "" {
        out = append(out, cur)
    }
    return out
}


