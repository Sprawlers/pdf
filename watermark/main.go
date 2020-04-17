package main

import (
    "os"
    "net/http"
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

type Server struct {
    s3       *s3.S3
    bucket   string
    endpoint string
}

func main() {
    var (
        key = os.Getenv("SPACES_KEY")
        secret = os.Getenv("SPACES_SECRET")
        bucket = os.Getenv("SPACES_BUCKET")
        endpoint = os.Getenv("SPACES_ENDPOINT")
        region = os.Getenv("SPACES_REGION")
    )
    s3config := &aws.Config{
        Credentials:    credentials.NewStaticCredentials(key, secret, ""),
        Endpoint:       aws.String(endpoint),
        Region:         aws.String(region),
    }
    newSession := session.New(s3config)
    s3 := s3.New(newSession)
    s := &Server{s3, bucket, endpoint}
    mux := http.NewServeMux()
    mux.HandleFunc("/", s.handleHealth())
    mux.HandleFunc("/watermark", s.handleWatermark())
    log.Fatal(http.ListenAndServe(":80", mux))
}
