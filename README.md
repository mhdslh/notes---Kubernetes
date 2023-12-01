# notes---Kubernetes

1- Deployment vs ReplicaSet: Deployment provides higher-level abstractions and additional features such as rolling updates, rollbacks, and versioning of the application. ReplicaSet is a lower-level abstraction that provides basic scaling mechanisms, i.e., ensures that a specified number of Pod replicas are running at any given time. As a result, we are allowed to edit a limited number of fields in a Pod's definition with 'kubectl edit'. For other changes, we need to manually remove the existing Pod first before creating a new Pod (replace). However, in Deployments' definitions, we can modify those fields with `kubectl edit' since Deployment performs [rolling update](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/).

2- Specified "command" and "args" fields, for each container, override "ENTRYPOINT" and "CMD" of its image, respectively. When "command" is provided, Dockerfile "ENTRYPOINT" and "CMD" are both ignored. When "args" is only set, "CMD" in Dockerfile is replaced with "args". The first argument of "command" (or "ENTRYPOINT") must be an executable (We may need to update PATH env. variable to be able to run the executable). All arguments in "args" (or "CMD") will be appended to "command" (or "ENTRYPOINT"). Since providing long options is possible in both "--key=value" and "--key value" formats ([reference](https://unix.stackexchange.com/questions/573377/do-command-line-options-take-an-equals-sign-between-option-name-and-value), args/CMD can be written in both ways: '["--key=value"]' and '["--key". "value"]'

3- ConfigMap stores non-confidential data in key-value pairs. ConfigMap can be injected into Pods as **environment variables**, command-line arguments, or as configuration files in a **volume**.
ConfigMap accepts both single line property values and multi-line file-like values. When creating ConfigMap using "kubectl create ConfigMap", "--from-literal" option creates single line property values and "--from-file" creates multi-line file-like values. In the later case, a key will be created from the file's name with its content as the value. We can control multi-line values format through header as explained in [this link](https://yaml-multiline.info/). We can inject the entire ConfigMap data into a Pod as environment variable(s) (under Pod's "spec.containers[].envFrom" section) or volume(s) (under Pod's "spec.volumes" section). We can also inject only specific keys into a Pod as environment variables (under Pod's "spec.containers[].env" section) or volumes (using items under Pod's "spec.volumes" section).
  
4- Different types of Secret exist. "Opaque" is the default Secret type. "kubernetes.io/dockerconfigjson" type is used to store credentials for accessing a container image registry. To pull an image from registry using "kubernetes.io/dockerconfigjson" type secrets, in Pod definition we need to add "imagePullSecrets" field under "spec" section.

5- ServiceAccounts are used to authenticate to Kubernetes API server. Before v1.22, for every ServiceAccount a long-lived static token was created using Secrets. Then by setting "spec.serviceAccountName" inside Pod, Kubernetes mounted that specific ServiceAccount's token, instead of default ServiceAccount's token, as a volume inside the Pod. From v1.22, kubernetes gets short-lived automatically rotating (instead of long-lived static) tokens using the TokenRequest API (instead of Secret) and mounts it in Pod as a projected volume. These tokens are time and audiance bounded (their lifetime depends on the Pod rather than the ServiceAccount). From v1.24, Kubernetes no longer generates tokens automatically. Administrators are responsible for that, for instance by running "kubectl create token <service-account-name>". To prevent kubernetes from automatically injecting credentials (for a specified ServiceAccount or the default ServiceAccount) in the Pod, we must set "spec.automountServiceAccountToken" to false. Note that we can still get long-lived static tokens (similar to what we had before v1.22) using Secrets of type "kubernetes.io/service-account-token". Finally, by default, ServiceAccounts are granted default permissions. We must use RBAC to grant required permissions.

6, 7, 8, and 9 describe **Pod** scheduling onto nodes:

  6- ResourceQuotas limit aggregate resource consumption (limits.cpu, limits.memory, requests.cpu, requests.memory) per namespace. LimitRange is a policy to constrain the resource allocations (limits and requests) specified for each applicable object kind (such as Pod or PersistentVolumeClaim) in a namespace. For instance, if resource requests and limits are specified for a Pod they must be in the range [min, max] defined in LimitRange, and if not specified they will use default values defined in LimitRange. 

  7- Taints allow nodes to repel Pods. A Pod can be scheduled on a node only if it tolerates the taint, i.e., Tolerations applied to the Pod match the taints. Note that Tolerations allow scheduling but don't guarantee scheduling. In otherwise, tolerating a node's taints is the necessary condition to be able to schedule Pods on that node.

  8- With adding nodeSelector field in the Pod's definition and specifying node labels, a pod will **only** be scheduled onto the nodes that have all of the specified labels.

  9- Node affinity is conceptually similar to nodeSelector, allowing you to constrain which nodes your Pod can be scheduled on based on node labels. Compared to nodeSelector, affinity/anti-affinity is more expressive and provides more control over the selection logic. It also allows soft (preferred) rules. Similar to node taints, anti-affinity repel Pods from specific nodes.

10- In the absence of readinessProbe, a Pod is considered ready when all of its containers are created (**reference?**). With introducing a readinessProbe under .spec.containers[] for a container, application itself can decide on its readiness. An unready Pod does not receive traffic through Kubernetes Services. 
Liveness probes can be used to detect when to restart a container. For example, liveness probes can catch a deadlock, where an application is running and ready (receives traffic), but unable to make progress. Note that usually there is no need to consider application crash for liveness probe because upon main application (PID 1) crash the Pod is subjected to its restart policy (set by .spec.restartPolicy field).

11- Labels are key/value pairs that are attached to objects. Labels do not provide uniqueness. In general, we expect many objects to carry the same label(s). Via a label selector, the client/user can identify a set of objects. One usage scenario for label requirement is for Pods to specify node selection criteria (check Point 9). The API currently supports two types of selectors: equality-based and set-based ([more information](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)). In kubectl commands, --selector option (or -l in short) allows filtering by label keys and values. --field-selector option also allows selecting Kubernetes objects based on the value of one or more resource fields. For example, metadata.namespace!=default or status.phase=Pending. It is worth mentioning that, unlike --selector option, set-based selectors are not supported for --field-selector option.

12- A Deployment revision is created when a Deployment rollout is triggered, i.e., if and only if Deployment's Pod template (.spec.template) is changed. Other updates such as scaling the Deployment do not create a revision. Rolling update is the default rollout strategy where Pods are updated incrementally (there is a limit on maximum number of new Pods and unavailable Pods during update).  Updates can be reverted. Deployment's revision history is stored in the ReplicaSet it controls. Once an old ReplicaSet is deleted, we lose the ability to rollback to that revision of Deployment. By default, 10 old ReplicaSets will be kept (to change this value set .spec.revisionHistoryLimit).'

13- There are two Kubernetes-native strategies for updating Deployments, namely Recreate and RollingUpdate (default value). The updating strategy can be specified by setting .spec.strategy field. As explained earlier, in RollingUpdate, Pods are updated incrementally, while, for Recreate, all existing Pods are killed before new ones are created. With RollingUpdate, we can set maxUnavailable (absolute number/percentage of desired Pods unavailable at all times during update) and maxSurge (absolute number/percentage of Pods that can be created over the desired number of Pods) to control the process. As a Kubernetes developer, we can employ two other updating strategies, namely [Blue/Green](https://docs.aws.amazon.com/whitepapers/latest/overview-deployment-options/bluegreen-deployments.html) and [Canary](https://docs.aws.amazon.com/whitepapers/latest/introduction-devops-aws/canary-deployments.html). They both can be implemented through labels and services

14- A Job creates one or more Pods (.spec.parallelism) and will continue to retry execution of the Pods until a specified number of them (.spec.completions) successfully terminate. When a specified number of successful completions is reached, the task (ie, Job) is complete. To run a Job (either a single task or several in parallel) on a schedule use CronJob. For Jobs and CronJobs, only a RestartPolicy equal to Never or OnFailure is allowed in the Pod's spec (Unlike ReplicaSets and Deployments where RestartPolicy is Always). To specify the number of retries before considering a Job as failed set .spec.backoffLimit (default value is 6). Failed Pods associated with the Job are recreated by the Job controller with an exponential back-off delay.

15- A NodePort can span across multiple nodes, i.e., it can send traffic to backing Pods even when they are on different nodes.

16- 
![image](https://github.com/mhdslh/notes---Kubernetes/assets/61638154/c35b6761-61bd-4334-a324-42ddf33cbcd9)

17- Remember namespaces provide a mechanism for isolating groups of resources within a single cluster. For instance, a service directs traffic to the Pods that match its selector within the same namespace. Pods in namespace 'my-ns' can call the service 'my-svc' in that namespace by using its name,i.e., 'my-svc'. However, Pods in other namespaces must call this service 'my-svc.my-ns'. This is how DNS records can be used to contact services.

18- Ingress controller and Network Plugin must be configured for minikube and kind cluster to be able to use ingress resources and network policies.

19- In a deployment with multiple containers, if a container crashes the Pod will not restart, only the crashed container will.

multiple nodes in kind (containers in the background) and volumes on the same path in each node

---
Helpful 'kubectl' commands:
- kubectl explain <resource-type>: to find out about api version for a resource type.
- kubectl top <node or pod>: to see the resource consumption for nodes or pods.

---
Links:
- [Kubernetes-Sigs](https://github.com/kubernetes-sigs)

---
To do:
https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/
After creating a container with kind, inside the container we have containerd client command line tool (ctl) and docker client is not provided. What is containerd?
core dump
kubectl port-forward
