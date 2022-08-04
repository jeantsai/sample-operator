package assets

import (
	"embed"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	//go:embed manifests/*
	menifests  embed.FS
	appsScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
}

func GetLegacyDeploymentFromEmbeddedFile(name string) *appsv1.Deployment {
	depBytes, err := menifests.ReadFile(name)
	if err != nil {
		panic(err)
	}
	depObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion),
		depBytes,
	)
	if err != nil {
		panic(err)
	}
	return depObject.(*appsv1.Deployment)
}
