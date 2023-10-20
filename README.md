# notes---Kubernetes

1- Deployment vs ReplicaSet: Deployment provides higher-level abstractions and additional features such as rolling updates, rollbacks, and versioning of the application. ReplicaSet is a lower-level abstraction that provides basic scaling mechanisms, i.e., ensures that a specified number of Pod replicas are running at any given time. As a result, we are allowed to edit a limited number of fields in a Pod's definition with 'kubectl edit'. For other changes, we need to manually remove the existing Pod first before creating a new Pod (replace). However, in Deployments' definitions, we can modify those fields with `kubectl edit' since Deployment performs [rolling update](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/).

2- Specified "command" and "args" fields, for each container, override "ENTRYPOINT" and "CMD" of its image, respectively. When "command" is provided, Dockerfile "ENTRYPOINT" and "CMD" are both ignored. When "args" is only set, "CMD" in Dockerfile is replaced with "args". The first argument of "command" (or "ENTRYPOINT") must be an executable (We may need to update PATH env. variable to be able to run the executable). All arguments in "args" (or "CMD") will be appended to "command" (or "ENTRYPOINT"). Since providing long options is possible in both "--key=value" and "--key value" formats ([reference](https://unix.stackexchange.com/questions/573377/do-command-line-options-take-an-equals-sign-between-option-name-and-value), args/CMD can be written in both ways: '["--key=value"]' and '["--key". "value"]'

3- ConfigMap stores non-confidential data in key-value pairs. ConfigMap can be injected into Pods as **environment variables**, command-line arguments, or as configuration files in a **volume**.
ConfigMap accepts both single line property values and multi-line file-like values. When creating ConfigMap using "kubectl create ConfigMap", "--from-literal" option creates single line property values and "--from-file" creates multi-line file-like values. In the later case, a key will be created from the file's name with its content as the value. We can control multi-line values format through header as explained in [this link](https://yaml-multiline.info/). We can inject the entire ConfigMap data into a Pod as environment variable(s) (under Pod's "spec.containers[].envFrom" section) or volume(s) (under Pod's "spec.volumes" section). We can also inject only specific keys into a Pod as environment variables (under Pod's "spec.containers[].env" section) or volumes (using items under Pod's "spec.volumes" section).
  
4- Different types of Secret exist. "Opaque" is the default Secret type. "kubernetes.io/dockerconfigjson" type is used to store credentials for accessing a container image registry. To pull an image from registry using "kubernetes.io/dockerconfigjson" type secrets, in Pod definition we need to add "imagePullSecrets" field under "spec" section.


---
Helpful 'kubectl' commands:
kubectl explain <resource-type>: to find out about api version for a resource type. 

---
To do:
kubectl replace + kubectl replace -f -

After creating a container with kind, inside the container we have containerd client command line tool (ctl) and docker client is not provided. What is containerd?

uid, gid, group, and capabilities
