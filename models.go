package ucodesdk

type (
	Request struct {
		Data     map[string]any `json:"data"`
		IsCached bool           `json:"is_cached"`
	}
)

// Response structures
type (
	// Create function response body >>>>> CREATE
	Datas struct {
		Data struct {
			Data map[string]any `json:"data"`
		} `json:"data"`
	}

	// ClientApiResponse This is get single api response >>>>> GET_SINGLE_BY_ID, GET_SLIM_BY_ID
	ClientApiResponse struct {
		Data ClientApiData `json:"data"`
	}

	ClientApiData struct {
		Data ClientApiResp `json:"data"`
	}

	ClientApiResp struct {
		Response map[string]any `json:"response"`
	}

	Response struct {
		Status string         `json:"status"`
		Error  string         `json:"error"`
		Data   map[string]any `json:"data"`
	}

	// GetListClientApiResponse This is get list api response >>>>> GET_LIST, GET_LIST_SLIM
	GetListClientApiResponse struct {
		Data GetListClientApiData `json:"data"`
	}

	GetListClientApiData struct {
		Data GetListClientApiResp `json:"data"`
	}

	GetListClientApiResp struct {
		Count    int32            `json:"count"`
		Response []map[string]any `json:"response"`
	}
	// GetListAggregationClientApiResponse  This is get list aggregation response
	GetListAggregationClientApiResponse struct {
		Data struct {
			Data struct {
				Data []map[string]any `json:"data"`
			} `json:"data"`
		} `json:"data"`
	}

	// ClientApiUpdateResponse This is single update api response >>>>> UPDATE
	ClientApiUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			TableSlug string         `json:"table_slug"`
			Data      map[string]any `json:"data"`
		} `json:"data"`
	}

	// ClientApiMultipleUpdateResponse This is multiple update api response >>>>> MULTIPLE_UPDATE
	ClientApiMultipleUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Data struct {
				Objects []map[string]any `json:"objects"`
			} `json:"data"`
		} `json:"data"`
	}

	ResponseError struct {
		StatusCode         int
		Description        any
		ErrorMessage       string
		ClientErrorMessage string
	}
)

type ActionBody struct {
	Body        map[string]any `json:"data"`
	DisableFaas bool           `json:"disable_faas"`
}

type AuthRequest struct {
	Body    map[string]any    `json:"data"`
	Headers map[string]string `json:"headers"`
}

type APIItem struct {
	collection string
	config     *Config
}

type CreateItem struct {
	collection string
	config     *Config
	data       ActionBody
}

type DeleteItem struct {
	collection  string
	config      *Config
	disableFaas bool
	id          string
}

type DeleteMultipleItem struct {
	collection  string
	config      *Config
	disableFaas bool
	ids         []string
}

type UpdateItem struct {
	collection string
	config     *Config
	data       ActionBody
}

type GetSingleItem struct {
	collection string
	config     *Config
	guid       string
}

type GetListItem struct {
	collection string
	config     *Config
	request    Request
	limit      int
	page       int
}

type GetListAggregation struct {
	collection string
	config     *Config
	request    Request
}

type Register struct {
	config *Config
	data   AuthRequest
}

type ResetPassword struct {
	config *Config
	data   AuthRequest
}

type Login struct {
	config *Config
	data   AuthRequest
}

type SendCode struct {
	config *Config
	data   AuthRequest
}

type APIAuth struct {
	config *Config
}

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

type APIFunction struct {
	config  *Config
	request Request
	path    string
}

type User struct {
	Id           string `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	ProjectId    string `json:"project_id"`
	RoleId       string `json:"role_id"`
	ClientTypeId string `json:"client_type_id"`
}

type Token struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	ExpiresAt        string `json:"expires_at"`
	RefreshInSeconds int32  `json:"refresh_in_seconds"`
}

type RegisterResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		UserFound      bool   `json:"user_found"`
		UserId         string `json:"user_id"`
		Token          *Token `json:"token"`
		LoginTableSlug string `json:"login_table_slug"`
		EnvironmentId  string `json:"environment_id"`
		User           *User  `json:"user"`
		UserIdAuth     string `json:"user_id_auth"`
	} `json:"data"`
}

type CreateFileResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		ID               string `json:"id"`
		Title            string `json:"title"`
		Storage          string `json:"storage"`
		FileNameDisk     string `json:"file_name_disk"`
		FileNameDownload string `json:"file_name_download"`
		Link             string `json:"link"`
		FileSize         int    `json:"file_size"`
	} `json:"data"`
	CustomMessage string `json:"custom_message"`
}

type FunctionResponse struct {
	Status        string `json:"status"`
	Description   string `json:"description"`
	Data          any    `json:"data"`
	CustomMessage any    `json:"custom_message"`
}

type LoginWithOptionResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		UserFound bool           `json:"user_found"`
		UserId    string         `json:"user_id"`
		Token     *Token         `json:"token"`
		Sessions  []*Session     `json:"sessions"`
		UserData  map[string]any `json:"user_data"`
	} `json:"data"`
}

type Session struct {
	Id           string `json:"id"`
	ProjectId    string `json:"project_id"`
	ClientTypeId string `json:"client_type_id"`
	UserId       string `json:"user_id"`
	RoleId       string `json:"role_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	UserIdAuth   string `json:"user_id_auth"`
}

type LoginResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		UserFound        bool             `json:"user_found"`
		ClientType       map[string]any   `json:"client_type"`
		UserId           string           `json:"user_id"`
		Role             map[string]any   `json:"role"`
		Token            *Token           `json:"token"`
		Permissions      []map[string]any `json:"permissions"`
		Sessions         []*Session       `json:"sessions"`
		LoginTableSlug   string           `json:"login_table_slug"`
		AppPermissions   []map[string]any `json:"app_permissions"`
		ResourceId       string           `json:"resource_id"`
		EnvironmentId    string           `json:"environment_id"`
		User             *User            `json:"user"`
		GlobalPermission map[string]any   `json:"global_permission"`
		UserData         map[string]any   `json:"user_data"`
		UserIdAuth       string           `json:"user_id_auth"`
	} `json:"data"`
}

type SendCodeResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		SmsId       string `json:"sms_id"`
		GoogleAcces bool   `json:"google_acces"`
		UserFound   bool   `json:"user_found"`
	} `json:"data"`
}
