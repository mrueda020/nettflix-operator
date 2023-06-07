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

	appsScheme  = runtime.NewScheme()
	appsCodecs  = serializer.NewCodecFactory(appsScheme)
	corev1Codec = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDecoder(corev1.SchemeGroupVersion)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}

	if err := corev1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
}

func GetDeploymentFromFile(name string) (*appsv1.Deployment, error) {
	deploymentBytes, err := manifests.ReadFile(name)
	if err != nil {
		return nil, err
	}

	deploymentObject := &appsv1.Deployment{}
	if _, _, err := appsCodecs.UniversalDeserializer().Decode(deploymentBytes, nil, deploymentObject); err != nil {
		return nil, err
	}

	return deploymentObject, nil
}

func GetServiceFromFile(name string) (*corev1.Service, error) {
	serviceBytes, err := manifests.ReadFile(name)
	if err != nil {
		return nil, err
	}

	serviceObject := &corev1.Service{}
	if _, _, err := corev1Codec.Decode(serviceBytes, nil, serviceObject); err != nil {
		return nil, err
	}

	return serviceObject, nil
}
