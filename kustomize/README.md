![image](https://github.com/user-attachments/assets/c1647fb0-c4dd-4818-a6e6-956d7f103a8a)

-  [Built-in generators and transformers](https://kubectl.docs.kubernetes.io/references/kustomize/builtins/)

-  A fieldSpec list, in a transformer's configuration, determines which resource types and which fields within those types the transformer can modify.
```
group: some-group
version: some-version
kind: some-kind
path: path/to/the/field
create: false
```
If create is set to true, the transformer creates the path to the field in the resource if the path is not already found. This is most useful for label and annotation transformers, where the path for labels or annotations may not be set before the transformation. ([refrence](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/transformerconfigs/README.md))

- Strategic Merge Patch supports special operations through directives. To learn more refer to [link](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/strategic-merge-patch.md#basic-patch-format)
