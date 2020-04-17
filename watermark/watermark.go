package main

import (
    "net/http"
    "os"
    "io"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/pdfcpu/pdfcpu/pkg/api"
    pdf "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

const WMConfig = "points:48, color:0.203 0.286 0.333, op:.60"

type WatermarkRequest struct {
    Url         string `json:"url"`
    Text        string `json:"text"`
    Path        string `json:"path"`
    Filename    string `json:"filename"`
}

type WatermarkResponse struct {
    Url         string `json:"url"`
}

func (s *Server) handleWatermark() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var wmreq WatermarkRequest
        decodeBody(r, &wmreq)
        resp, err := http.Get(wmreq.Url)
        if err != nil {
            respondErr(w, r, http.StatusBadRequest, err)
            return
        }
        in, err := os.Create("./in.pdf")
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        defer in.Close()
        _, err = io.Copy(in, resp.Body)
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        out, err := os.Create("./out.pdf")
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        defer out.Close()
        wm, err := pdf.ParseTextWatermarkDetails(wmreq.Text, WMConfig, true)
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        err = api.AddWatermarks(in, out, nil, wm, nil)
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        result, err := os.Open("./out.pdf")
        if err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        defer result.Close()
        object := s3.PutObjectInput{
            Bucket: aws.String(s.bucket),
            Key:    aws.String(wmreq.Path + wmreq.Filename),
            Body:   result,
            ACL:    aws.String("public-read"),
        }
        if _, err := s.s3.PutObject(&object); err != nil {
            respondErr(w, r, http.StatusInternalServerError, err)
            return
        }
        wmres := WatermarkResponse{s.endpoint + s.bucket + wmreq.Path + wmreq.Filename}
        respond(w, r, http.StatusOK, wmres)
    }
}
