package ucodesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type UcodeApis2 interface {
	Items(collection string) ItemsI
	Auth() AuthI
	Files() FilesI
	Function() FunctionI
	Config() *Config
	DoRequest(url string, method string, body interface{}, headers map[string]string) ([]byte, error)
}

func NewSDK(cfg *Config) *UcodeAPI {
	return &UcodeAPI{
		config: cfg,
	}
}

// UcodeAPI struct implements UcodeAPIInterface
type UcodeAPI struct {
	config *Config
}

func (u *UcodeAPI) Items(collection string) ItemsI {
	return &APIItem{
		collection: collection,
		config:     u.config,
	}
}

func (u *UcodeAPI) Config() *Config {
	return u.config
}

func (u *UcodeAPI) Auth() AuthI {
	return &APIAuth{
		config: u.config,
	}
}

func (u *UcodeAPI) Files() FilesI {
	return &APIFiles{
		config: u.config,
	}
}

func (u *UcodeAPI) Function() FunctionI {
	return &APIFunction{}
}

// CREATE ITEM

// Items interface defines methods related to item operations
type ItemsI interface {
	Create(data map[string]any) *CreateItem
	Update(data map[string]any) *UpdateItem
	Delete() *DeleteItem
	GetList() *GetListItem
	GetSingle(id string) *GetSingleItem
}

// CREATE ITEM EXEC
func (a *APIItem) Create(data map[string]any) *CreateItem {
	return &CreateItem{
		collection: a.collection,
		config:     a.config,
		data:       ActionBody{Body: data},
	}
}

func (c *CreateItem) DisableFaas(isDisable bool) *CreateItem {
	c.data.DisableFaas = isDisable
	return c
}

func (c *CreateItem) Exec() (Datas, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		createdObject Datas
		url           = fmt.Sprintf("%s/v2/items/%s?from-ofs=%t", c.config.BaseURL, c.collection, c.data.DisableFaas)
	)

	var appId = c.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	createObjectResponseInByte, err := DoRequest(url, http.MethodPost, c.data, header)
	if err != nil {
		response.Data = map[string]any{"description": string(createObjectResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	err = json.Unmarshal(createObjectResponseInByte, &createdObject)
	if err != nil {
		response.Data = map[string]any{"description": string(createObjectResponseInByte), "message": "Error while unmarshalling create object", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	return createdObject, response, nil
}

// UPDATE ITEM EXEC
func (a *APIItem) Update(data map[string]any) *UpdateItem {
	return &UpdateItem{
		collection: a.collection,
		config:     a.config,
		data:       ActionBody{Body: data},
	}
}

func (a *UpdateItem) DisableFaas(isDisable bool) *UpdateItem {
	a.data.DisableFaas = isDisable
	return a
}

func (u *UpdateItem) ExecSingle() (ClientApiUpdateResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		updateObject ClientApiUpdateResponse
		url          = fmt.Sprintf("%s/v2/items/%s?from-ofs=%t", u.config.BaseURL, u.collection, u.data.DisableFaas)
	)

	var appId = u.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	updateObjectResponseInByte, err := DoRequest(url, http.MethodPut, u.data, header)
	if err != nil {
		response.Data = map[string]any{"description": string(updateObjectResponseInByte), "message": "Error while updating object", "error": err.Error()}
		response.Status = "error"
		return ClientApiUpdateResponse{}, response, err
	}

	err = json.Unmarshal(updateObjectResponseInByte, &updateObject)
	if err != nil {
		response.Data = map[string]any{"description": string(updateObjectResponseInByte), "message": "Error while unmarshalling update object", "error": err.Error()}
		response.Status = "error"
		return ClientApiUpdateResponse{}, response, err
	}

	return updateObject, response, nil
}

func (a *UpdateItem) ExecMultiple() (ClientApiMultipleUpdateResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		multipleUpdateObject ClientApiMultipleUpdateResponse
		url                  = fmt.Sprintf("%s/v1/object/multiple-update/%s?from-ofs=%t", a.config.BaseURL, a.collection, a.data.DisableFaas)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	multipleUpdateObjectsResponseInByte, err := DoRequest(url, http.MethodPut, a.data, header)
	if err != nil {
		response.Data = map[string]any{"description": string(multipleUpdateObjectsResponseInByte), "message": "Error while multiple updating objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, err
	}

	err = json.Unmarshal(multipleUpdateObjectsResponseInByte, &multipleUpdateObject)
	if err != nil {
		response.Data = map[string]any{"description": string(multipleUpdateObjectsResponseInByte), "message": "Error while unmarshalling multiple update objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, err
	}

	return multipleUpdateObject, response, nil
}

// DELETE ITEM EXEC
func (a *APIItem) Delete() *DeleteItem {
	return &DeleteItem{
		collection: a.collection,
		config:     a.config,
	}
}

func (a *DeleteItem) DisableFaas(disable bool) *DeleteItem {
	a.disableFaas = disable
	return a
}

func (a *DeleteItem) Single(id string) *DeleteItem {
	a.id = id
	return a
}

func (a *DeleteMultipleItem) Multiple(ids []string) *DeleteMultipleItem {
	a.ids = ids
	return a
}

func (a *DeleteItem) Exec() (Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		url = fmt.Sprintf("%s/v2/items/%s/%v?from-ofs=%t", a.config.BaseURL, a.collection, a.id, a.disableFaas)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodDelete, Request{Data: map[string]any{}}, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while deleting object", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (a *DeleteMultipleItem) Exec() (Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		url = fmt.Sprintf("%s/v1/object/%s/?from-ofs=%t", a.config.BaseURL, a.collection, a.disableFaas)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodDelete, map[string]any{"ids": a.ids}, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while deleting objects", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

// GET SINGLE ITEM EXEC
func (a *APIItem) GetSingle(id string) *GetSingleItem {
	return &GetSingleItem{
		collection: a.collection,
		config:     a.config,
		guid:       id,
	}
}

func (a *GetSingleItem) Exec() (ClientApiResponse, Response, error) {
	var (
		response  = Response{Status: "done"}
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v2/items/%s/%v?from-ofs=%t", a.config.BaseURL, a.collection, a.guid, true)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	resByte, err := DoRequest(url, http.MethodGet, nil, header)
	if err != nil {
		response.Data = map[string]any{"description": string(resByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]any{"description": string(resByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

// GET SINGLE SLIM ITEM EXEC

func (a *GetSingleItem) ExecSlim() (ClientApiResponse, Response, error) {
	var (
		response  = Response{Status: "done"}
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v1/object-slim/%s/%v?from-ofs=%t", a.config.BaseURL, a.collection, a.guid, true)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	resByte, err := DoRequest(url, http.MethodGet, nil, header)
	if err != nil {
		response.Data = map[string]any{"description": string(resByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]any{"description": string(resByte), "message": "Error while unmarshalling to object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

// GET LIST ITEM EXEC
func (a *APIItem) GetList() *GetListItem {
	return &GetListItem{
		collection: a.collection,
		config:     a.config,
		request:    Request{Data: map[string]any{}},
	}
}

func (a *GetListItem) Limit(limit int) *GetListItem {
	if limit <= 0 {
		limit = 10
	}
	a.limit = limit
	a.request.Data["offset"] = (a.page - 1) * limit
	a.request.Data["limit"] = limit
	return a
}

func (a *GetListItem) Page(page int) *GetListItem {
	if page <= 0 {
		page = 1
	}
	a.page = page
	a.request.Data["offset"] = (page - 1) * a.limit
	return a
}

func (a *GetListItem) Filter(filter map[string]any) *GetListItem {
	for key, value := range filter {
		a.request.Data[key] = value
	}
	return a
}

func (a *GetListItem) Search(search string) *GetListItem {
	a.request.Data["search"] = search
	return a
}

func (a *GetListItem) Sort(sort map[string]any) *GetListItem {
	a.request.Data["order"] = sort
	return a
}

func (a *GetListItem) ViewFields(fields []string) *GetListItem {
	a.request.Data["view_fields"] = fields
	return a
}

func (a *GetListItem) Pipelines(query map[string]any) *GetListAggregation {
	return &GetListAggregation{
		collection: a.collection,
		config:     a.config,
		request:    Request{Data: query},
	}
}

type GetListArggItem struct {
}

func (a *GetListItem) WithRelations(with bool) *GetListItem {
	a.request.Data["with_relations"] = with
	return a
}

func (a *GetListItem) Exec() (GetListClientApiResponse, Response, error) {
	var (
		response      = Response{Status: "done"}
		getListObject GetListClientApiResponse
		url           = fmt.Sprintf("%s/v2/object/get-list/%s?from-ofs=%t", a.config.BaseURL, a.collection, true)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListResponseInByte, err := DoRequest(url, http.MethodPost, a.request, header)
	if err != nil {
		response.Data = map[string]any{"description": string(getListResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &getListObject)
	if err != nil {
		response.Data = map[string]any{"description": string(getListResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	return getListObject, response, nil
}

func (a *GetListItem) ExecSlim() (GetListClientApiResponse, Response, error) {
	var (
		response = Response{Status: "done"}
		listSlim GetListClientApiResponse
		url      = fmt.Sprintf("%s/v2/object-slim/get-list/%s?from-ofs=%t", a.config.BaseURL, a.collection, true)
	)

	reqObject, err := json.Marshal(a.request.Data)
	if err != nil {
		response.Data = map[string]any{"message": "Error while marshalling request getting list slim object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	url = fmt.Sprintf("%s&data=%s&offset=%d&limit=%d", url, string(reqObject), (a.page-1)*a.limit, a.limit)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListResponseInByte, err := DoRequest(url, http.MethodGet, nil, header)
	if err != nil {
		response.Data = map[string]any{"description": string(getListResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &listSlim)
	if err != nil {
		response.Data = map[string]any{"description": string(getListResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	return listSlim, response, nil
}

func (a *GetListAggregation) ExecAggregation() (GetListAggregationClientApiResponse, Response, error) {
	var (
		response           = Response{Status: "done"}
		getListAggregation GetListAggregationClientApiResponse
		url                = fmt.Sprintf("%s/v2/items/%s/aggregation", a.config.BaseURL, a.collection)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListAggregationResponseInByte, err := DoRequest(url, http.MethodPost, a.request, header)
	if err != nil {
		response.Data = map[string]any{"description": string(getListAggregationResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListAggregationClientApiResponse{}, response, err
	}
	fmt.Println("asdfadfas", string(getListAggregationResponseInByte))

	err = json.Unmarshal(getListAggregationResponseInByte, &getListAggregation)
	if err != nil {
		response.Data = map[string]any{"description": string(getListAggregationResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListAggregationClientApiResponse{}, response, err
	}

	return getListAggregation, response, nil
}

type AuthI interface {
	Login()
	Register(data AuthRequest) *Register
	SendCode(data AuthRequest)
	ResetPassword(data AuthRequest) *ResetPassword
}

type Register struct {
	config *Config
	data   AuthRequest
}

type ResetPassword struct {
	config *Config
	data   AuthRequest
}

type APIAuth struct {
	config *Config
}

func (a *APIAuth) Login() {
}

func (a *APIAuth) Register(data AuthRequest) *Register {
	return &Register{
		config: a.config,
		data:   data,
	}
}

func (a *Register) Exec() (RegisterResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		registerObject RegisterResponse
		url            = fmt.Sprintf("%s/v2/register?project-id=%s", a.config.AuthBaseURL, a.config.ProjectId)
	)

	registerResponseInByte, err := DoRequest(url, http.MethodPost, a.data.Body, a.data.Headers)
	if err != nil {
		response.Data = map[string]any{"description": string(registerResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return RegisterResponse{}, response, err
	}

	err = json.Unmarshal(registerResponseInByte, &registerObject)
	if err != nil {
		response.Data = map[string]any{"description": string(registerResponseInByte), "message": "Error while unmarshalling register object", "error": err.Error()}
		response.Status = "error"
		return RegisterResponse{}, response, err
	}

	return registerObject, response, nil
}

func (a *APIAuth) ResetPassword(data AuthRequest) *ResetPassword {
	return &ResetPassword{
		config: a.config,
		data:   data,
	}
}

func (a *ResetPassword) Exec() (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v2/reset-password", a.config.AuthBaseURL)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodPut, a.data.Body, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while reset password", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (a *APIAuth) SendCode(data AuthRequest) {
}

// Files interface defines methods for file operations
type FilesI interface {
	Upload(filePath string) *UploadFile
	Delete(fileID string) *DeleteFile
}

// APIFiles struct implements FilesInterface
type APIFiles struct {
	config *Config
}

type UploadFile struct {
	config *Config
	path   string
}

type DeleteFile struct {
	config *Config
	id     string
}

func (f *APIFiles) Upload(filePath string) *UploadFile {
	return &UploadFile{
		config: f.config,
		path:   filePath,
	}
}

func (c *UploadFile) Exec() (CreateFileResponse, Response, error) {
	var (
		file          *os.File
		fileBuffer    bytes.Buffer
		writer        *multipart.Writer
		response      = Response{Status: "done"}
		createdObject CreateFileResponse
		url           = fmt.Sprintf("%s/v1/files/folder_upload?folder_name=Media", c.config.BaseURL)
	)

	file, err := os.Open(c.path)
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't open file by path", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}
	defer file.Close()

	writer = multipart.NewWriter(&fileBuffer)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't create from file", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't copy file", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	err = writer.Close()
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't close writer", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	var appId = c.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	createFileInByte, err := DoFileRequest(url, http.MethodPost, header, fileBuffer, writer)
	if err != nil {
		response.Data = map[string]any{"description": string(createFileInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	err = json.Unmarshal(createFileInByte, &createdObject)
	if err != nil {
		response.Data = map[string]any{"description": string(createFileInByte), "message": "Error while unmarshalling create file object", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	return createdObject, response, nil
}

func (f *APIFiles) Delete(fileID string) *DeleteFile {
	return &DeleteFile{
		config: f.config,
		id:     fileID,
	}
}

func (a *DeleteFile) Exec() (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v1/files/%s", a.config.BaseURL, a.id)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodDelete, Request{Data: map[string]any{}}, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while deleting file", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

// Function interface defines methods for invoking functions
type FunctionI interface {
	InvokeByPath(path string)
}

// APIFunction struct implements FunctionInterface
type APIFunction struct{}

func (f *APIFunction) InvokeByPath(path string) {
}

func DoRequest(url string, method string, body any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Add headers from the map
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}

func (a *UcodeAPI) DoRequest(url string, method string, body any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Add headers from the map
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}

func DoFileRequest(url, method string, headers map[string]string, body bytes.Buffer, writer *multipart.Writer) ([]byte, error) {
	request, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)

	return respByte, err
}
