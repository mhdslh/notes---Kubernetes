
-   On host machine, create certificate for webhook server with `openssl`  using a SAN configuration file as below:
```
cat > san.cnf <<EOF
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = webhook-svc.admission-webhook.svc

[v3_req]
subjectAltName = @alt_names

[alt_names]
DNS.1 = webhook-svc.admission-webhook.svc
EOF
``` 
```
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -config san.cnf
``` 
Then copy both cert and key files to webhook server image (needed for webhook server to reply to HTTPS requests sent by the API server). clientConfig.caBundle field in ValidatingWebhookConfiguration manifest also needs base64 encoded of ca-cert.pem 

-   Admission webhooks are deployed in a separate namespace. Therefore, we must create the namespace before installing the webhook server. We must deploy the webhook configuration in the end. Installing configuration before server may block installing the webhook server because validating admission webhooks try to validate the webhook but will fail.

- To deserialize raw data Kubernetes objects into Go structs:

Approach I (using json deserializer): works well when type of the data we're deserializing is known
```
    pod := &corev1.Pod{}
    err := json.Unmarshal(raw, pod)
```
Approach II (using runtime package's decoding functions): aware of Kubernetes' API conventions and can decode serialized Kubernetes objects into the appropriate Go structs, even if the exact type isn't known at compile time.
```
	// import
	// corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/serializer"

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	deserializer := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	obj, _, _ := deserializer.Decode(admissionreview.Request.Object.Raw, nil, nil)
	switch o := obj.(type) {
	case *corev1.Pod:
		log.Printf("object of Kind Pod with spec: %v\n", o.Spec)
	...
	}
```
- Webhooks are sent as POST requests, with Content-Type: application/json, with an AdmissionReview API object in the admission.k8s.io API group serialized to JSON as the body. Webhooks respond with a 200 HTTP status code, Content-Type: application/json, and a body containing an AdmissionReview object (in the same version they were sent), with the response stanza populated, serialized to JSON. (When rejecting a request, the webhook can customize the http code and message returned to the user.)

- In AdmissionReview response, patch field contains a base64-encoded array of JSON patch operations. In our implementation, at first glance, it may seem that this field is not encoded. In golang, JSON marshaling of []byte is base64 encoded by default (i.e., output is the base64 encoded version of the input string when it is passed as []byte). Similarly, when unmarshaling into a []byte it will base64 decode first. Patch field in ```type AdmissionResponse struct``` is of type []byte and, in our code, we have
```
respBody, err := json.Marshal(admissionreview)
w.Write(respBody)
```
Therefore, Patch field will be base64 encoded when sending the HTTP response to API server.

TODO:
-	check pod and deployment
-   Use kustomize/helm for image name... (Kustomize (for passing image name andnamespace + for monitoring)? in makefile)

Useful Links:
- https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
- https://jsonpatch.com/
- https://www.civo.com/learn/kubernetes-admission-controllers-for-beginners
