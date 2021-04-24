package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Komplementariteten/lutra/auth"
	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/util"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/bson"
)

const HeaderNameContentType = "Content-Type"
const contentItemTypeName = "Type"
const contentItemContentName = "Content"
const contentItemCollectionName = "Collection"
const contentItemOwnerName = "Owner"
const contentTypeMap = "map"
const contentTypeFile = "file"
const fileContentReferenceName = "R"
const fileContentPayloadName = "C"
const fileNameFormat = "%X_%X_%s"

var FileCachePath string

type apiHandle struct {
	db *db.Db
	Cleanabel
}

type ContentType interface {
	ToBsonD() *bson.D
	Read(*fastjson.Object) error
}

type Content struct {
	Content    ContentType
	Owner      string
	Type       string
	Collection string
}

type ContentMap struct {
	C map[string]interface{}
}

type ContentFile struct {
	C []byte
	R string
	T string
}

type ApiRequestResponse struct {
	Error          bool
	HttpStatusCode int
	ContentType    string
	Message        string
}

func (a *apiHandle) ShouldBeExecuted(n int) bool {
	targetTime := 4
	return n%targetTime == 0
}
func (a *apiHandle) Cleanup(ctx context.Context) {
	os.RemoveAll(FileCachePath)
	os.Mkdir(FileCachePath, 0777)
	fmt.Println("apiHandle Cleanup called")
}

func responseToAddition(response *ApiRequestResponse, item *db.TrackedItem) {
	response.HttpStatusCode = http.StatusOK
	response.ContentType = "application/json"
	response.Message = fmt.Sprintf(`{"id": "%s"}`, item.InsertedID.Hex())
	response.Error = false
}

func (api *apiHandle) AddToDb(ctx context.Context, c *Content) *ApiRequestResponse {
	response := &ApiRequestResponse{
		Error:          true,
		HttpStatusCode: http.StatusNotFound,
		ContentType:    "text/plain",
		Message:        "Not Found",
	}

	if c.Type == contentTypeMap {
		pool, err := api.db.Collection(ctx, c.Collection)
		if err != nil {
			response.Message = err.Error()
			return response
		}
		item, err := pool.Add(ctx, c.toBsonD())
		if err != nil {
			response.Error = true
			response.Message = err.Error()
		} else {
			responseToAddition(response, item)
		}
		return response
	}
	if c.Type == contentTypeFile {
		fileTemplate := c.Content.(*ContentFile)
		item, err := api.db.AddFile(ctx, fmt.Sprintf(fileNameFormat, c.Collection, c.Owner, fileTemplate.R), fileTemplate.C)
		if err != nil {
			response.Error = true
			response.Message = err.Error()
		} else {
			responseToAddition(response, item)
		}
	}
	return response
}

func (c *Content) toBsonD() *bson.D {
	doc := &bson.D{
		{contentItemOwnerName, c.Owner},
		{contentItemContentName, c.Content.ToBsonD()},
	}
	return doc
}

func (m *ContentMap) Read(obj *fastjson.Object) error {
	m.C = make(map[string]interface{})
	obj.Visit(func(k []byte, v *fastjson.Value) {
		keyString := string(k[:])
		switch v.Type() {
		case fastjson.TypeArray:
			valueArray, err := v.Array()
			if err != nil {
				break
			}
			valueType := valueArray[0].Type()
			arrayLen := len(valueArray)
			switch valueType {
			case fastjson.TypeString:
				arr := make([]string, arrayLen)
				for i := 0; i < arrayLen; i++ {
					arr[i] = valueArray[i].String()
				}
				m.C[keyString] = arr
			case fastjson.TypeNumber:
				arr := make([]int64, arrayLen)
				for i := 0; i < arrayLen; i++ {
					arr[i] = valueArray[i].GetInt64()
				}
				m.C[keyString] = arr
			}
			break
		case fastjson.TypeNumber:
			m.C[keyString] = v.GetInt64()
			break
		case fastjson.TypeString:
			byteValue := v.GetStringBytes()
			m.C[keyString] = string(byteValue[:])
			break
		case fastjson.TypeNull:
			m.C[keyString] = nil
			break
		case fastjson.TypeObject:
			o, err := v.Object()
			if err == nil {
				objMap := util.ReadJsonObjectAsMap(o)
				m.C[keyString] = objMap
			}
			break
		}
	})
	return nil
}

func (m *ContentMap) ToBsonD() *bson.D {
	var doc bson.D
	for k, v := range m.C {
		doc = append(doc, bson.E{k, v})
	}
	return &doc
}

func (f *ContentFile) Read(obj *fastjson.Object) error {
	ref := obj.Get(fileContentReferenceName)
	b64 := obj.Get(fileContentPayloadName)
	if !ref.Exists() {
		return fmt.Errorf("Reference Identifier not set in json")
	}
	if !b64.Exists() {
		return fmt.Errorf("Base64 Content not set in json")
	}
	f.R = string(ref.GetStringBytes())
	b64StrBytes := b64.GetStringBytes()
	b64Str := string(b64StrBytes)
	b64Str2 := ""
	sepIndex := strings.Index(b64Str, ",")
	if sepIndex > 0 {
		typeSlice := b64Str[0:sepIndex]
		f.T = typeSlice
		b64Str2 = b64Str[sepIndex+1:]
	} else if sepIndex > -1 {
		b64Str2 = b64Str[sepIndex+1:]
	}
	b64bytes, err := base64.StdEncoding.DecodeString(b64Str2)

	if err != nil {
		return fmt.Errorf("Content is not Base64 Encoded: %s", err.Error())
	}
	f.C = b64bytes
	return nil
}

func (f *ContentFile) ToBsonD() *bson.D {
	panic("File can not be serialized")
}

func ResponseToHttp(w http.ResponseWriter, response *ApiRequestResponse) {
	w.Header().Set("Content-Type", response.ContentType)
	w.WriteHeader(response.HttpStatusCode)
	io.WriteString(w, response.Message)
}

func (c *Content) SetContentObject() error {
	typeString := strings.ToLower(string(c.Type[:]))
	switch typeString {
	case contentTypeMap:
		c.Content = &ContentMap{}
		break
	case contentTypeFile:
		c.Content = &ContentFile{}
		break
	default:
		return fmt.Errorf("Content Type not known")
	}
	return nil
}

func checkJson(json []byte, response *ApiRequestResponse) *Content {
	var p fastjson.Parser
	v, err := p.ParseBytes(json)
	if err != nil {
		response.Error = true
		response.Message = err.Error()
		return nil
	}
	content := &Content{}
	if !v.Exists(contentItemTypeName) {
		response.Error = true
		response.Message = "Type not found"
		return nil
	}

	if !v.Exists(contentItemContentName) {
		response.Error = true
		response.Message = "Content not found"
		return nil
	}

	if !v.Exists(contentItemCollectionName) {
		response.Error = true
		response.Message = "Collection is missing"
		return nil
	}

	content.Type = string(v.GetStringBytes(contentItemTypeName))
	content.Collection = string(v.GetStringBytes(contentItemCollectionName))
	content.Owner = auth.Environment.User
	err = content.SetContentObject()
	if err != nil {
		response.Error = true
		response.Message = "Invalid Content JSON"
		return nil
	}

	err = content.Content.Read(v.GetObject(contentItemContentName))
	if err != nil {
		response.Error = true
		response.Message = "Valid to read JSON as content type"
		return nil
	}
	response.Error = false

	return content
}

func checkJsonInHttp(r *http.Request) (*ApiRequestResponse, *Content) {
	response := &ApiRequestResponse{
		Error:          true,
		HttpStatusCode: http.StatusUnsupportedMediaType,
		ContentType:    "text/plain",
		Message:        "Request contains no valid json or has maleformed header",
	}
	contentType := r.Header.Get(HeaderNameContentType)
	if contentType == "" {
		response.Message = "You need to set the content-type header"
		return response, nil
	}
	if !strings.HasPrefix(contentType, "application/json") {
		return response, nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Message = err.Error()
		return response, nil
	}
	content := checkJson(body, response)
	if response.Error {
		return response, nil
	}
	return nil, content
}

func (api *apiHandle) writeToDb(w http.ResponseWriter, r *http.Request) {
	resp, content := checkJsonInHttp(r)
	if content == nil {
		ResponseToHttp(w, resp)
		return
	}

	switch r.Method {
	case http.MethodPost:
		resp = api.AddToDb(r.Context(), content)
		break
	case http.MethodDelete:
		break
	case http.MethodPut:
		break
	default:
		resp.Error = true
		resp.ContentType = "text/plain"
		resp.HttpStatusCode = http.StatusNotFound
		resp.Message = fmt.Sprintf("http method %s is not implemented for %s", r.Method, r.URL.Path)
		break
	}
	ResponseToHttp(w, resp)
}
func (api *apiHandle) getFileFromDb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	referenceIndex := strings.LastIndex(r.URL.Path, "/")
	reference := r.URL.Path[referenceIndex+1:]
	filename := path.Join(FileCachePath, reference)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		b, err := api.db.ReadFile(r.Context(), reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = ioutil.WriteFile(filename, b, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(b)
		return
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(b)
	//api.db.AddFile()

}

func (api *apiHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(FileCachePath) == 0 {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			cacheDir = os.TempDir()
		}
		FileCachePath = path.Join(cacheDir, "adv_cache")
	}
	_, err := os.Stat(FileCachePath)
	if os.IsNotExist(err) {
		os.Mkdir(FileCachePath, 0777)
	}

	_, err = auth.CheckRefresh(w, r)
	if err != nil {
		w.Write([]byte("HTTP Error"))
		return
	}
	if strings.HasPrefix(r.URL.Path, "write") {
		api.writeToDb(w, r)
	}
	if strings.HasPrefix(r.URL.Path, "files") {
		api.getFileFromDb(w, r)
	}
	if strings.HasPrefix(r.URL.Path, "read") {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}

}
