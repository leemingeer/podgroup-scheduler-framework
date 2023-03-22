
# podgroup-scheduler-framework

Based on 1.20-release branch of scheduler-plugin

## start

### 创建crd
```shell
# kubectl apply -f config/crd/scheduling.ming.io_podgroups.yaml
```
### 部署sample-scheduler
```
# helm install podgroup-scheduler ./install/charts/as-a-second-scheduler/
# kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
sample-scheduler-598cf4cf56-ntptg   1/1     Running   0          15h
```

## 测试
```
# kubectl apply -f config/example/

# kubectl get podgroup
NAME    AGE
nginx   18h

# kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
nginx-cbmx5                         0/1     Pending   0          15h
nginx-hrr6q                         0/1     Pending   0          15h
sample-scheduler-598cf4cf56-ntptg   1/1     Running   0          15h

# kubectl get pod
NAME                                READY   STATUS              RESTARTS   AGE
nginx-8sgkm                         0/1     Pending   0          20s
nginx-fw4lf                         0/1     Pending   0          20s

# kubectl describe pod nginx-8sgkm 
...
Events:
  Type     Reason            Age   From              Message
  ----     ------            ----  ----              -------
  Warning  FailedScheduling  59m   sample-scheduler  0/3 nodes are available: 3 pre-filter pod nginx-8sgkm cannot find enough sibling pods, current pods number: 2, minMember of group: 3.
```
在PreFilter扩展点不满足调度条件

将replicas从2修改为3后，达到群组调度的限制
```
# kubectl get pod
NAME                                READY   STATUS              RESTARTS   AGE
nginx-8sgkm                         0/1     ContainerCreating   0          86s
nginx-fw4lf                         0/1     ContainerCreating   0          86s
nginx-qhrg2                         1/1     Running             0          50s
```

## 编译镜像
```
make release-image.arm64
docker push leemingeer/sample-scheduler:v20230317-v0.0.1-arm64
```