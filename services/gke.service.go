package services

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func DeletePod() {
	namespace := os.Getenv("POD_NAMESPACE") // Pod の Namespace
	if namespace == "" {
		namespace = "default" // デフォルト Namespace
	}

	podName, _ := os.Hostname() // ✅ 自身の Pod 名を取得
	log.Println("Deleting Pod:", podName)

	// Kubernetes API Server のエンドポイント
	kubeAPIURL := fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods/%s", namespace, podName)

	// ServiceAccount のトークンを取得
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Println("Error reading token:", err)
		return
	}

	// 証明書検証をスキップ
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Kubernetes API に DELETE リクエストを送信
	req, _ := http.NewRequest("DELETE", kubeAPIURL, nil)
	req.Header.Set("Authorization", "Bearer "+string(token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error deleting pod:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)
}
