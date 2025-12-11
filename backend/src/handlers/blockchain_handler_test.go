package handlers_test

import (
	"app/blockchain"
	"app/handlers"
	"app/requestModels"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockchainStatus_ErrorConnecting(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
	

    // Override NewClient to simulate error
    handlers.NewClient = func() (*blockchain.Client, error) {
        return nil, errors.New("connection failed")
    }

    h := handlers.BlockchainHandler{}
    h.GetBlockchainStatus(c)

    assert.Equal(t, http.StatusOK, w.Code)

    var resp requestModels.BlockchainStatusResponse
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.False(t, resp.Connected)
    assert.Equal(t, "connection failed", resp.Error)
}

func TestDeployContract_Error(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    handlers.DeployContract = func() (common.Address, error) {
        return common.Address{}, errors.New("deploy failed")
    }

    h := handlers.BlockchainHandler{}
    h.DeployContract(c)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.Contains(t, w.Body.String(), "deploy failed")
}

func TestDeployContract_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    handlers.DeployContract = func() (common.Address, error) {
        return common.HexToAddress("0x123"), nil
    }

    h := handlers.BlockchainHandler{}
    h.DeployContract(c)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "Contract deployed")
    assert.Contains(t, w.Body.String(), "0x0000000000000000000000000000000000000123")
}
