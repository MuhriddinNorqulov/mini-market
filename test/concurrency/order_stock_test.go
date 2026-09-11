package concurrency_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	container "mini-market/cmd/container"
)

func TestMain(m *testing.M) {
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
		_ = os.Chdir(repoRoot)
	}
	os.Exit(m.Run())
}

func baseURL() string {
	if v := os.Getenv("APP_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func testJWT(t *testing.T) (token string) {
	t.Helper()
	if v := os.Getenv("TEST_JWT"); v != "" {
		return v
	}

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("could not self-issue a test JWT — is Postgres reachable at localhost (e.g. via `docker compose up`)? panic: %v", r)
		}
	}()

	username := fmt.Sprintf("cc-test-%d", time.Now().UnixNano()%1_000_000_000)
	issued, err := container.InitSeedApp().IssueTestToken(username)
	if err != nil {
		t.Skipf("could not self-issue a test JWT — is Postgres reachable at localhost (e.g. via `docker compose up`)? %v", err)
	}
	return issued
}

func adminJWT(t *testing.T) string {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"username": "admin", "password": "admin123"})
	req, _ := http.NewRequest(http.MethodPost, baseURL()+"/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var parsed struct {
		Payload struct {
			AccessToken string `json:"access_token"`
		} `json:"payload"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	if parsed.Payload.AccessToken == "" {
		t.Skipf("could not log in as the default admin user — is the stack fully up (docker compose up, including seed-users)?")
	}
	return parsed.Payload.AccessToken
}

type orderResponse struct {
	Code    string          `json:"code"`
	Success bool            `json:"success"`
	Payload json.RawMessage `json:"payload"`
	Message string          `json:"message"`
}

type placeOrderResult struct {
	status int
	body   orderResponse
}

func TestFiftyConcurrentOrdersAgainstStockOfTen(t *testing.T) {
	token := testJWT(t)

	productID := createProduct(t, adminJWT(t), 10)

	var wg sync.WaitGroup
	results := make([]placeOrderResult, 50)

	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = placeOrder(t, token, productID, fmt.Sprintf("concurrency-test-%d", i))
		}(i)
	}
	wg.Wait()

	var successCount, conflictCount int64
	for i, r := range results {
		switch r.status {
		case http.StatusCreated:
			successCount++
			assertSuccessfulOrderBody(t, i, r.body)
		case http.StatusConflict, http.StatusUnprocessableEntity:
			conflictCount++
			assertInsufficientStockBody(t, i, productID, r.body)
		default:
			t.Fatalf("result %d: unexpected HTTP status %d, body %+v", i, r.status, r.body)
		}
	}

	if successCount != 10 {
		t.Fatalf("expected exactly 10 successful orders, got %d", successCount)
	}
	if conflictCount != 40 {
		t.Fatalf("expected exactly 40 rejected orders, got %d", conflictCount)
	}

	remaining := getProductStock(t, token, productID)
	if remaining != 0 {
		t.Fatalf("expected final stock_quantity == 0, got %d", remaining)
	}
}

func assertSuccessfulOrderBody(t *testing.T, i int, body orderResponse) {
	t.Helper()
	if !body.Success {
		t.Fatalf("result %d: expected success=true on a 201 response, got %+v", i, body)
	}
	if body.Code != "OK" {
		t.Fatalf("result %d: expected code=OK on a 201 response, got %q", i, body.Code)
	}
	var payload struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body.Payload, &payload); err != nil {
		t.Fatalf("result %d: could not parse order payload: %v (raw: %s)", i, err, body.Payload)
	}
	if payload.ID == 0 {
		t.Fatalf("result %d: expected a non-zero order id in the payload, got %+v", i, payload)
	}
	if payload.Status != "pending" {
		t.Fatalf("result %d: expected status=pending in the payload, got %q", i, payload.Status)
	}
}

func assertInsufficientStockBody(t *testing.T, i int, productID uint, body orderResponse) {
	t.Helper()
	if body.Success {
		t.Fatalf("result %d: expected success=false on a rejection, got %+v", i, body)
	}
	if body.Code != "INSUFFICIENT_STOCK" {
		t.Fatalf("result %d: expected code=INSUFFICIENT_STOCK, got %q (body: %+v)", i, body.Code, body)
	}
	var failedIDs []uint
	if err := json.Unmarshal(body.Payload, &failedIDs); err != nil {
		t.Fatalf("result %d: expected payload to be a list of product IDs, got %s: %v", i, body.Payload, err)
	}
	if len(failedIDs) != 1 || failedIDs[0] != productID {
		t.Fatalf("result %d: expected payload [%d], got %v", i, productID, failedIDs)
	}
	if !strings.Contains(strings.ToLower(body.Message), "insufficient stock") {
		t.Fatalf("result %d: expected message to mention insufficient stock, got %q", i, body.Message)
	}
	if !strings.Contains(body.Message, strconv.Itoa(int(productID))) {
		t.Fatalf("result %d: expected message to name product id %d, got %q", i, productID, body.Message)
	}
}

func createProduct(t *testing.T, token string, stock int) uint {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"name": "concurrency-test-product", "price": 100, "stock_quantity": stock})
	req, _ := http.NewRequest(http.MethodPost, baseURL()+"/api/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var parsed struct {
		Payload struct {
			ID uint `json:"id"`
		} `json:"payload"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	if parsed.Payload.ID == 0 {
		t.Fatalf("failed to create product, status %d", resp.StatusCode)
	}
	return parsed.Payload.ID
}

func placeOrder(t *testing.T, token string, productID uint, idempotencyKey string) placeOrderResult {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"items": []map[string]any{{"product_id": productID, "quantity": 1}},
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL()+"/api/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", idempotencyKey+"-"+strconv.Itoa(int(productID)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var parsed orderResponse
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	return placeOrderResult{status: resp.StatusCode, body: parsed}
}

func getProductStock(t *testing.T, token string, productID uint) int64 {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, baseURL()+"/api/products/"+strconv.Itoa(int(productID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var parsed struct {
		Payload struct {
			StockQuantity int64 `json:"stock_quantity"`
		} `json:"payload"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	return parsed.Payload.StockQuantity
}
