package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/bson"
)

const ResourceFolder = "./res/"
const letterRunes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterRunes) {
			b[i] = letterRunes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func IfFileExists(fileName string) (bool, error) {
	s, err := os.Lstat(fileName)
	if err != nil {
		if err != os.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	switch mode := s.Mode(); {
	case mode.IsRegular():
		return true, nil
	default:
		return false, fmt.Errorf("%s has File Mode %s", fileName, mode.String())
	}
}

func CreateResourceFolderIfNotExists() error {
	s, err := os.Stat(ResourceFolder)
	if err != nil {
		if err != os.ErrNotExist {
			os.Mkdir(ResourceFolder, 0700)
			return nil
		}
		return err
	}
	if s.IsDir() {
		return nil
	}
	return fmt.Errorf("%s is not a Directory", ResourceFolder)
}

func ReadJsonObjectAsMap(obj *fastjson.Object) map[string]interface{} {
	objMap := make(map[string]interface{})
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
				objMap[keyString] = arr
			case fastjson.TypeNumber:
				arr := make([]int64, arrayLen)
				for i := 0; i < arrayLen; i++ {
					arr[i] = valueArray[i].GetInt64()
				}
				objMap[keyString] = arr
			}
			break
		case fastjson.TypeNumber:
			objMap[keyString] = v.GetInt64()
			break
		case fastjson.TypeString:
			byteValue := v.GetStringBytes()
			objMap[keyString] = string(byteValue[:])
			break
		case fastjson.TypeNull:
			objMap[keyString] = nil
			break
		case fastjson.TypeObject:
			o, err := v.Object()
			if err == nil {
				subObj := ReadJsonObjectAsMap(o)
				objMap[keyString] = subObj
			}
			break
		}
	})
	return objMap
}

func GetJsonFromReader(reader io.ReadCloser) (map[string]interface{}, error) {
	dataBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var p fastjson.Parser
	jValue, err := p.ParseBytes(dataBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println(jValue.String())
	jObj, err := jValue.Object()
	if err != nil {
		return nil, err
	}
	m := ReadJsonObjectAsMap(jObj)
	fmt.Printf("Parsed: %v", m)
	return m, nil
}

// ToBsonD transforms a given Interface to bson.D though binary marschalling
func ToBsonD(d interface{}) (doc *bson.D, err error) {
	bytes, err := bson.Marshal(d)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(bytes, &doc)
	if err != nil {
		return nil, err
	}
	return
}
