package assets

import (
	"embed"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	//go:embed manifests/*
	manifests embed.FS

	appsScheme   = runtime.NewScheme()
	appsCodecs   = serializer.NewCodecFactory(appsScheme)
	coresScheme  = runtime.NewScheme()
	corev1Codecs = serializer.NewCodecFactory(coresScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}

	if err := corev1.AddToScheme(coresScheme); err != nil {
		panic(err)
	}
}

func GetDeploymentFromFile(name string) *appsv1.Deployment {
	deploymentBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}

	deploymentObject, err := runtime.Decode(appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion), deploymentBytes)
	if err != nil {
		panic(err)
	}

	return deploymentObject.(*appsv1.Deployment)
}

func GetServiceFromFile(name string) *corev1.Service {
	servicesBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}

	serviceObject, err := runtime.Decode(corev1Codecs.UniversalDecoder(corev1.SchemeGroupVersion), servicesBytes)
	if err != nil {
		panic(err)
	}

	return serviceObject.(*corev1.Service)
}
