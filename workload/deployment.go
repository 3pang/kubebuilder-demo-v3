package workload

import (
	"context"
	"fmt"

	v1 "k8s.io/api/admission/v1"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var logger = ctrl.Log.WithName("workload/deployment")

type DeployWrapper struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *DeployWrapper) Handle(ctx context.Context, req admission.Request) admission.Response {
	fmt.Println(">>>>start new webhook handle>>>>>>")
	fmt.Printf(">>>>>work load -> deployment Handle: res name=%s,user=%s,Operation=%s \n", req.Name, req.UserInfo.Username, req.Operation)

	fmt.Println("--------------------1-req-------")
	fmt.Println(req)
	fmt.Println("--------------------1-req-end-------")

	fmt.Println("--------------------2-a.decoder-------")
	fmt.Println(a.decoder)
	fmt.Println("--------------------2-a.decoder-end-------")

	fmt.Println("--------------------3-a.decoder.Decode(req, deploy)-------")
	//err := a.decoder.Decode(req, deploy)
	fmt.Println("--------------------3-a.decoder.Decode(req, deploy)-end-------")
	deploy := &appv1.Deployment{}
	if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, deploy); err != nil {
		fmt.Println(err)
	}
	fmt.Println("--------------------4-deploy &appv1.Deployment{}-------")
	fmt.Println(deploy)
	fmt.Println("--------------------4-deploy &appv1.Deployment{}-end-------")

	//originObj, err := json.Marshal(deploy)
	//if err != nil {
	//	return admission.Errored(http.StatusBadRequest, err)
	//}
	// 将新的资源副本数量改为1
	//newobj := deploy.DeepCopy()
	//replicas := int32(1)
	//newobj.Spec.Replicas = &replicas
	//currentObj, err := json.Marshal(newobj)
	//if err != nil {
	//	return admission.Errored(http.StatusBadRequest, err)
	//}
	// 对比之前的资源类型和之后的资源类型的差异生成返回数据
	//resp := admission.PatchResponseFromRaw(originObj, currentObj)
	//if err != nil {
	//	return admission.Errored(http.StatusBadRequest, err)
	//}
	//respJson, err := json.Marshal(resp.AdmissionResponse)
	//fmt.Println(string(respJson))

	if req.Operation == v1.Create { // 如果是更新，判断是否image有变化。通过判断path的路径是否是/spec/containers
		//判断lable 是否具备  labels.owner/name: wy

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

func (a *DeployWrapper) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
