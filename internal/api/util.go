package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/c4t-but-s4d/fastad/pkg/httpext"
)

func ProtoJSON(c echo.Context, code int, msg proto.Message) error {
	raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return httpext.NewErrorf(
			http.StatusInternalServerError,
			"marshaling teams: %v",
			err,
		)
	}

	return c.JSONBlob(code, raw)
}

func ProtoRaw(c echo.Context, code int, msg proto.Message) error {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return httpext.NewErrorf(
			http.StatusInternalServerError,
			"marshaling teams: %v",
			err,
		)
	}

	return c.Blob(code, "application/octet-stream", raw)
}
