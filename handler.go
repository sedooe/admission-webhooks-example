package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type namespaceRequestHandler func(*v1beta1.AdmissionRequest, corev1.Namespace) *v1beta1.AdmissionResponse

func handleReq(c *gin.Context, reqHandler namespaceRequestHandler) {
	var admissionResponse *v1beta1.AdmissionResponse
	admissionReview := v1beta1.AdmissionReview{}

	if err := c.ShouldBindJSON(&admissionReview); err != nil {
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	var namespace corev1.Namespace
	req := admissionReview.Request

	if err := json.Unmarshal(req.Object.Raw, &namespace); err != nil {
		glog.Errorf("Could not unmarshal raw object: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse {
			Result: &metav1.Status {
				Message: err.Error(),
			},
		}
	}

	if admissionResponse == nil {
		admissionResponse = reqHandler(req, namespace)
	}

	admissionReview.Response = admissionResponse
	admissionReview.Response.UID = admissionReview.Request.UID
	c.JSON(http.StatusOK, admissionReview)
}
