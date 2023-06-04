package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}
