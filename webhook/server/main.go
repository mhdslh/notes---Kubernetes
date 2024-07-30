package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type Update struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

var (
	deserializer runtime.Decoder
)

// mutate adds "createdBy" and "purpose" annotations for pods
func mutate(admissionreview *admissionv1.AdmissionReview) {
	log.Println("mutating resource")
	admissionresponse := &admissionv1.AdmissionResponse{}
	admissionreview.Response = admissionresponse
	admissionresponse.UID = admissionreview.Request.UID

	// before adding annotations, we must check if any already exist. If the metadata.annotations field is not set, we must initialize it as an empty object to enable the addition of new key-value pairs.
	obj, _, err := deserializer.Decode(admissionreview.Request.Object.Raw, nil, nil)
	if err != nil {
		log.Printf("unable to deserialize object from raw data: %s\n", err.Error())
		admissionresponse.Allowed = false
		return
	}
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		log.Println("object is not of type 'Pod': type assertion failed")
		admissionresponse.Allowed = false
		return
	}
	annotations := pod.ObjectMeta.Annotations

	admissionresponse.PatchType = new(admissionv1.PatchType)
	*admissionresponse.PatchType = admissionv1.PatchTypeJSONPatch
	updates := make([]Update, 0)
	if len(annotations) == 0 {
		updates = append(updates, Update{Op: "add", Path: "/metadata/annotations", Value: struct{}{}})
	}
	updates = append(updates, Update{Op: "add", Path: "/metadata/annotations/createdBy", Value: "Mohammad Salehi"})
	updates = append(updates, Update{Op: "add", Path: "/metadata/annotations/purpose", Value: "demo"})
	patch, err := json.Marshal(updates)
	if err != nil {
		log.Printf("failed to marshal JSON patches: %s\n", err.Error())
		admissionresponse.Allowed = false
		return
	}
	admissionresponse.Patch = patch

	admissionresponse.Allowed = true
}

// validate accepts a pod only if all container images include version tags
func validate(admissionreview *admissionv1.AdmissionReview) {
	log.Println("validating resource")
	admissionresponse := &admissionv1.AdmissionResponse{}
	admissionreview.Response = admissionresponse
	admissionresponse.UID = admissionreview.Request.UID

	obj, _, err := deserializer.Decode(admissionreview.Request.Object.Raw, nil, nil)
	if err != nil {
		log.Printf("unable to deserialize object from raw data: %s\n", err.Error())
		admissionresponse.Allowed = false
		return
	}
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		log.Println("object is not of type 'Pod': type assertion failed")
		admissionresponse.Allowed = false
		return
	}

	for _, container := range pod.Spec.InitContainers {
		if len(strings.Split(container.Image, ":")) == 1 {
			log.Printf("version tag for init container %s in %s pod is missing\n", container.Image, pod.ObjectMeta.Name)
			admissionresponse.Allowed = false
			return
		}
	}
	for _, container := range pod.Spec.Containers {
		if len(strings.Split(container.Image, ":")) == 1 {
			log.Printf("version tag for container %s in %s pod is missing\n", container.Image, pod.ObjectMeta.Name)
			admissionresponse.Allowed = false
			admissionresponse.Result = &metav1.Status{
				Message: fmt.Sprintf("version tag for container %s in %s pod is missing", container.Image, pod.ObjectMeta.Name),
			}
			return
		}
	}

	admissionresponse.Allowed = true
}

func init() {
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	deserializer = serializer.NewCodecFactory(scheme).UniversalDeserializer()
}

func main() {
	log.Println("registering handler for webhook requests")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received webhook request")
		var (
			reqBody         []byte
			respBody        []byte
			err             error
			admissionreview = &admissionv1.AdmissionReview{}
		)

		if reqBody, err = io.ReadAll(r.Body); err != nil {
			log.Printf("failed to read request body: %s\n", err.Error())
		}
		if err = json.Unmarshal(reqBody, admissionreview); err != nil {
			log.Printf("failed to unmarshal admission review: %s\n", err.Error())
		}

		switch r.URL.Path {
		case "/mutate":
			mutate(admissionreview)
		case "/validate":
			validate(admissionreview)
		default:
		}

		if respBody, err = json.Marshal(admissionreview); err != nil {
			log.Printf("failed to marshal admission review: %s\n", err.Error())
		}
		log.Printf("resp body: %s\n", respBody)
		w.Write(respBody)
	})

	log.Println("starting validating webhook server")
	log.Fatal(http.ListenAndServeTLS(":8080", "ca-cert.pem", "ca-key.pem", nil))
}
