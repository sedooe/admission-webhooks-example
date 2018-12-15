package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"k8s.io/api/admission/v1beta1"
	auth "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func patchLabels(userInfo auth.UserInfo, labels map[string]string) ([]byte, error) {
	org, project := getOrganizationAndProject(userInfo)

	if labels == nil {
		labels = map[string]string{"organization": org, "project": project}
	} else {
		labels["organization"] = org
		labels["project"] = project
	}

	patch := []patchOperation {
		{
			Op:    "replace",
			Path:  "/metadata/labels",
			Value: labels,
		},
	}

	return json.Marshal(patch)
}

func mutate(req *v1beta1.AdmissionRequest, namespace corev1.Namespace) *v1beta1.AdmissionResponse {
	patch, err := patchLabels(req.UserInfo, namespace.Labels)
	if err != nil {
		return &v1beta1.AdmissionResponse {
			Result: &metav1.Status {
				Message: err.Error(),
			},
		}
	}

	return &v1beta1.AdmissionResponse {
		Allowed: true,
		Patch:   patch,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

func mutationHandler(c *gin.Context) {
	handleReq(c, mutate)
}
