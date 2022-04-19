package grpcgateway

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const cookieName = "Authorization"

func metadataAnnotator(_ context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return metadata.New(md)
	}

	md[cookieName] = cookie.Value

	return metadata.New(md)
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()

	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		delete(w.Header(), "Grpc-Metadata-Location")
		w.Header().Set("location", location[0])
	}

	if codeRaw, ok := headers["Grpc-Metadata-X-Http-Code"]; ok {
		code, err := strconv.Atoi(codeRaw[0])
		if err != nil {
			return fmt.Errorf("converting status code to type int: %v", err)
		}

		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		defer w.WriteHeader(code)
	}

	cookie := http.Cookie{
		Name:  cookieName,
		Value: headers.Get("grpc-metadata-authorization"),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	delete(w.Header(), "Grpc-Metadata-Authorization")
	delete(w.Header(), "Grpc-Metadata-Content-Type")

	return nil
}
