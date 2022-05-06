package responseutil

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nmluci/sumber-sari-garden/internal/global/util/jsonutil"
	"github.com/nmluci/sumber-sari-garden/pkg/dto"
	"github.com/nmluci/sumber-sari-garden/pkg/errors"
)

func WriteSuccessResponse(rw http.ResponseWriter, status int, data interface{}) {
	BaseResponseWriter(rw, status, data, nil)
}

func WriteErrorResponse(rw http.ResponseWriter, err error) {
	errMetadata := errors.GetErrorResponseMetadata(err)
	BaseResponseWriter(rw, errMetadata.Status, nil, &dto.ErrorData{Code: errMetadata.Code, Message: errMetadata.Message})
}

func WriteFileResponse(rw http.ResponseWriter, status int, filename string, data bytes.Buffer) {
	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	rw.Header().Set("Content-Type", "text/csv")
	rw.Header().Set("Content-Length", strconv.Itoa(data.Len()))
	rw.WriteHeader(status)

	rw.Write(data.Bytes())
}

func BaseResponseWriter(rw http.ResponseWriter, status int, data interface{}, er *dto.ErrorData) {
	res := dto.BaseResponse{Data: data, Error: er}
	jsonData, err := jsonutil.ConvertToJSON(res)
	if err != nil {
		log.Printf("[WriteSuccessResponse] json conversion error: %v", err)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(jsonData)
}
