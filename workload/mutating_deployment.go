package workload

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	appv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type DeployMutationWrapper struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *DeployMutationWrapper) Handle(ctx context.Context, req admission.Request) admission.Response {
	deploy := &appv1.Deployment{}
	if a.decoder != nil {
		if err := a.decoder.Decode(req, deploy); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("========a.decoder is nil")
	}

	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	deploy.Name = u2.String()
	fmt.Println(deploy.Name)
	fmt.Println(deploy.Namespace)
	return admission.Allowed("ok")
}
func (a *DeployMutationWrapper) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
