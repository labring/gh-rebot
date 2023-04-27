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
	"bufio"
	"context"
	"github.com/cuisongliu/logger"
	"github.com/gin-gonic/gin"
	"github.com/labring-actions/gh-rebot/pkg/client-go/kubernetes"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"io"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type GetPodLogReq struct {
	Container    string `form:"container"`
	Follow       bool   `form:"follow"`
	TailLines    int64  `form:"tailLines"`
	SinceSeconds int64  `form:"sinceSeconds"`
	Previous     bool   `form:"previous"`
}

var cli kubernetes.Client

func init() {
	client, err := kubernetes.NewKubernetesClient("", "")
	if err != nil {
		panic(err)
	}
	cli = client
}

func podLogs(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("podName")
	getLogReq := GetPodLogReq{}

	if err := c.ShouldBindQuery(&getLogReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request params error", "msg": err.Error()})
		return
	}

	podLogOption := &v1.PodLogOptions{
		TailLines: &getLogReq.TailLines,
		Container: getLogReq.Container,
		Follow:    getLogReq.Follow,
		Previous:  getLogReq.Previous,
	}

	if getLogReq.SinceSeconds != 0 {
		podLogOption.SinceSeconds = &getLogReq.SinceSeconds
	}

	kubeClient := cli.Kubernetes()

	_, err := kubeClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, v12.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get pod failed", "msg": err.Error(), "pod": podName, "namespace": namespace})
		return
	}

	req := kubeClient.CoreV1().Pods(namespace).GetLogs(podName, podLogOption)

	if !podLogOption.Follow {
		logs, err := req.DoRaw(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "get pod logs failed", "msg": err.Error()})
			return
		}
		c.String(http.StatusOK, string(logs))
		return
	}

	stream, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(err, "init stream failed")
		return
	}
	defer stream.Close()

	buf := bufio.NewReader(stream)

	ws, err := utils.NewWSLogger(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err, "ws初始化失败")
		return
	}

	for {
		bytes, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				logger.Error(err, "read steam failed")
				return
			}
			return
		}
		//写入ws流
		if err := ws.Write(bytes); err != nil {
			logger.Error(err, "ws write error")
			return
		}
	}
}
