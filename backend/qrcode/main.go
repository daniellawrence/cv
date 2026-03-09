package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	"google.golang.org/protobuf/encoding/protojson"

	qrcodev1 "github.com/daniellawrence/cv/gen/go/qrcode/v1"
)

type bufferWriteCloser struct {
	*bytes.Buffer
}

func (b bufferWriteCloser) Close() error {
	return nil
}

func getQRCode(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	qrc, err := qrcode.New(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	writer := standard.NewWithWriter(bufferWriteCloser{buf}, standard.WithQRWidth(5))

	if err := qrc.Save(writer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	resp := &qrcodev1.GenerateQRCodeResponse{
		Url:         url,
		ImageBase64: encoded,
	}

	data, err := m.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/qrcode", getQRCode)

	log.Println("qrcode service listening on :8084")

	err := http.ListenAndServe(":8084", common.CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
