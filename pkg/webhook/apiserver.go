/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func RegistryHttpServer(port uint16) error {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.POST("/webhook", webhookHandler)

	return r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", port))
}

func webhookHandler(c *gin.Context) {
	signature := c.GetHeader("X-Hub-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}

	if !validateSignature(os.Getenv("SECRET_TOKEN"), signature, body) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
		return
	}

	logger.Info("Received valid webhook: %s", string(body))

	// Signature is valid, process the webhook payload
	c.String(http.StatusOK, "Received valid webhook")
}

func validateSignature(secret, signature string, payload []byte) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
