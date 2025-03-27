package apitest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/k1LoW/runn"
	apiserver "github.com/mochibuta/apitest-example/cmd/api-server/server"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:17",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
		},
		WaitingFor: wait.NewExecStrategy([]string{"pg_isready"}).WithStartupTimeout(60 * time.Second),
	}

	psqlContainer, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{ContainerRequest: req})
	if err != nil {
		log.Fatal(err)
	}

	if err := psqlContainer.Start(t.Context()); err != nil {
		log.Fatal(err)
	}

	psqlContainerPort, err := psqlContainer.MappedPort(t.Context(), "5432")
	if err != nil {
		log.Fatal(err)
	}

	t.Setenv("DB_PORT", psqlContainerPort.Port())

}

func setupMockServer(t *testing.T) *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/posts" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id":    1,
					"title": "test",
				},
			})
		}
	}))

	t.Setenv("EXTERNAL_API_URL", mockServer.URL)

	return mockServer
}

func TestAPISenario(t *testing.T) {
	setupTestContainer(t)
	mockServer := setupMockServer(t)
	defer mockServer.Close()

	srv, err := apiserver.InitServer(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	apiSrv := httptest.NewServer(srv)

	defer apiSrv.Close()
	defer apiserver.CloseDB()

	opts := []runn.Option{
		runn.T(t),
		runn.Runner("req", apiSrv.URL),
	}

	op, err := runn.Load("scenario/example.yaml", opts...)
	if err != nil {
		t.Fatal(err)
	}

	if err := op.RunN(t.Context()); err != nil {
		t.Fatal(err)
	}
}
