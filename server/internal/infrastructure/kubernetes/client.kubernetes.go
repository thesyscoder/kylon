package kubernetes

import (
	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"k8s.io/client-go/kubernetes"
)

var clientSet *kubernetes.Clientset

func InitKubernetesClient(cfg config.Config) {
}
