package workload

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/admission/v1"
	appv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type DeployWrapper struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *DeployWrapper) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx)
	logger.Info(">>>>start new webhook handle")
	logger.Info(fmt.Sprintf(">>>>>work load -> deployment Handle: res name=%s,user=%s,Operation=%s \n", req.Name, req.UserInfo.Username, req.Operation))
	logger.Info("--------------------1")
	obj := &appv1.Deployment{}
	logger.Info("--------------------2")
	err1 := a.decoder.Decode(req, obj)
	logger.Info("--------------------3")
	if err1 != nil {
		return admission.Errored(http.StatusBadRequest, err1)
	}
	logger.Info(">>>>obj.Namee", obj.Name)
	originObj, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	// 将新的资源副本数量改为1
	newobj := obj.DeepCopy()
	replicas := int32(1)
	newobj.Spec.Replicas = &replicas
	currentObj, err := json.Marshal(newobj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	// 对比之前的资源类型和之后的资源类型的差异生成返回数据
	resp := admission.PatchResponseFromRaw(originObj, currentObj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	respJson, err := json.Marshal(resp.AdmissionResponse)
	logger.Info(string(respJson))

	if req.Operation == v1.Update { // 如果是更新，判断是否image有变化。通过判断path的路径是否是/spec/containers
		respPatch := admission.PatchResponseFromRaw(req.OldObject.Raw, req.Object.Raw)
		if !respPatch.Allowed {
			return respPatch
		}
		patches := respPatch.Patches
		if len(patches) > 0 {
			for i := 0; i < len(patches); i++ {
				p := patches[i]
				if p.Operation != "add" && p.Operation != "replace" {
					logger.Info("work load -> deployment- >>!add and ! replace")
					return admission.Allowed("ok") //
				} else if p.Operation == "add" {
					logger.Info("work load -> deployment- >>add ")
					return admission.Allowed("ok") //
				} else if p.Operation == "replace" {
					logger.Info("work load -> deployment- >>replace")
					return admission.Allowed("ok") //
				}
			}
		}
		return admission.Allowed("ok")
	}
	return admission.Allowed("ok")
}
