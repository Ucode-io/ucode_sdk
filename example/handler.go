package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cast"
	sdk "github.com/ucode-io/ucode_sdk"
)

func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request            sdk.Request
			response           sdk.Response
			clientErrorMessage string
		)

		newsdk := sdk.New(&sdk.Config{
			BaseURL:   "",
			AppId:     "",
			ProjectId: "",
		})
		body := map[string]any{
			"title": fmt.Sprintf("%d", time.Now().Unix()),
		}

		createResp, _, err := newsdk.Items("order").Create(body).DisableFaas(true).Exec()
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}

		updateBody := map[string]any{
			"title": fmt.Sprintf("%d %s", time.Now().Unix(), "updated"),
			"guid":  createResp.Data.Data["guid"],
		}

		// Add is_new: true for multiple create
		updateResp, _, err := newsdk.Items("order").Update(updateBody).DisableFaas(true).ExecSingle()
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}

		fmt.Println(updateResp)

		_, err = newsdk.Items("order").Delete().Single(cast.ToString(createResp.Data.Data["guid"])).DisableFaas(true).Exec()
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}

		getListResp, _, err := newsdk.Items("order").
			GetList().
			Page(1).
			Limit(20).
			Sort(map[string]any{"created_at": -1}).
			Filter(map[string]any{"status": []string{"new"}}).
			Exec()
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}
		fmt.Println(getListResp)

		getListResp2, _, err := newsdk.Items("order_product").
			GetList().
			Page(1).
			Limit(20).
			Filter(map[string]any{
				"quantity": map[string]any{
					"$gte": 4,
				}},
			).
			Exec()
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}
		fmt.Println(getListResp2)

		requestByte, err := io.ReadAll(r.Body)
		if err != nil {
			clientErrorMessage = "Error on getting request body"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestByte, &request)
		if err != nil {
			clientErrorMessage = "Error on unmarshal request"
			handleResponse(w, returnError(clientErrorMessage, err.Error()), http.StatusInternalServerError)
			return
		}

		response.Status = "done"
		handleResponse(w, response, http.StatusOK)
	}
}

func handleResponse(w http.ResponseWriter, body any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	bodyByte, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`
			{
				"error": "Error marshalling response"
			}
		`))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bodyByte)
}

func returnError(clientError string, errorMessage string) interface{} {
	return sdk.Response{
		Status: "error",
		Data:   map[string]interface{}{"message": clientError, "error": errorMessage},
	}
}
