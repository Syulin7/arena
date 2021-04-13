## arena serve delete

Delete a serving job and its associated instances

### Synopsis

Delete a serving job and its associated instances

```
arena serve delete JOB1 JOB2 ...JOBn [-T JOB_TYPE] [--version JOB_VERSION] [flags]
```

### Options

```
  -h, --help             help for delete
  -T, --type string      The serving type to delete, the possible option is [trt(Tensorrt),seldon(Seldon),custom(Custom),kf(KFServing),tf(Tensorflow)]. (optional)
  -v, --version string   The serving version to delete.
```

### Options inherited from parent commands

```
      --arena-namespace string   The namespace of arena system service, like tf-operator (default "arena-system")
      --config string            Path to a kube config. Only required if out-of-cluster
      --loglevel string          Set the logging level. One of: debug|info|warn|error (default "info")
  -n, --namespace string         the namespace of the job
      --pprof                    enable cpu profile
      --trace                    enable trace
```

### SEE ALSO

* [arena serve](arena_serve.md)	 - Serve a job.

###### Auto generated by spf13/cobra on 5-Mar-2021