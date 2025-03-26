package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "go_app/config"  // 添加配置包导入
)

// ImageResponse ImgBB响应结构
type ImageResponse struct {
    Data struct {
        URL string `json:"url"`
        DisplayURL string `json:"display_url"`
    } `json:"data"`
    Success bool   `json:"success"`
    Status  int    `json:"status"`
}

// UploadToImageHost 上传图片到图床
func UploadToImageHost(file *multipart.FileHeader) (string, error) {
    // 获取配置
    cfg, err := config.LoadConfig()
    if err != nil {
        return "", fmt.Errorf("加载配置失败: %v", err)
    }

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // 添加API key
    writer.WriteField("key", cfg.ImageHost.Token)

    // 创建文件表单字段
    part, err := writer.CreateFormFile("image", file.Filename)
    if err != nil {
        return "", err
    }

    // 读取文件内容
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // 写入文件内容
    if _, err = io.Copy(part, src); err != nil {
        return "", err
    }
    writer.Close()

    // 创建请求
    req, err := http.NewRequest("POST", "https://api.imgbb.com/1/upload", body)
    if err != nil {
        return "", err
    }

    // 设置请求头
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 解析响应
    var result ImageResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    if !result.Success {
        return "", fmt.Errorf("上传失败: 状态码 %d", result.Status)
    }

    return result.Data.URL, nil
}