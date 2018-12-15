package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	auth "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func labelsValid(userInfo auth.UserInfo, labels map[string]string) bool {
	org, project := getOrganizationAndProject(userInfo)
	return (labels["organization"] == org) && (labels["project"] == project)
}

func validate(req *v1beta1.AdmissionRequest, namespace corev1.Namespace) *v1beta1.AdmissionResponse {
	if labelsValid(req.UserInfo, namespace.Labels) {
		return &v1beta1.AdmissionResponse {
			Allowed: true,
		}
	}

	glog.Infof("Denying the namespace update request. Request has forbidden label(s). " +
		"Username: " + req.UserInfo.Username)

	return &v1beta1.AdmissionResponse {
		Allowed: false,
		Result: &metav1.Status {
			Message: "Namespace update request rejected. You should not allowed to use \"organization\" " +
				"and \"project\" labels.",
		},
	}
}

func validationHandler(c *gin.Context) {
	handleReq(c, validate)
}
