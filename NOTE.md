
nginx.yaml中replicas为2,通过label标记pod属于的podgroup， 当前该group有两个pod，不满足最小调度数量3个，
因此这两个pod因为没有满足最小调度pod数而pending. 上面的打印是在framework的Prefilter，没有满足MinMember参数而返回framework.Unschedulable，中断pod的调度流程，所以两个pod都报上面的event.

## PreFilter扩展点逻辑
1. 检查lastDeniedPG对象，是否已经拒绝过该pg, 若是直接返回决绝过(last failed in 3s)， 若设置了cache有效期（deniedPGExpirationTimeSeconds）该判断只在有效期内生效。 有效期的计算是上次Add的时间+这个超时时间，通过go-cache库实现
2. 检查MinMember条件是否符合
3. 检查permittedPG是否允许过该pg, 若是, 直接返回。 内存有效期和lastDeniedPG一样，也是和deniedPGExpirationTimeSeconds参数
4. 检查集群中node的资源是否满足MinResources的请求。只要有一个节点满足即ok，如果所有节点均不满足，则会将pg加入到lastDeniedPG中，后续会走1的判断。
   4.1 在计算pg资源时，是按照组资源去计算，当满足组资源才可以调度，否则不会调度任何一个pod, 这里有个特殊case, 如果nodes上已经有属于该pg的pod,则MinResources会减去这一部分资源。
5. 最后将pg加入到permittedPG对象中

小结：
1. 将pod调度提升到group级别。要么都调度，要么都不调度
2. lastDeniedPG，防御性，防止调度过热
3. permittedPG，提速，防止重复判断

## permit扩展点

将replicas增大为3个， 3个pod依然处于pending状态. 通过debug发现，通过scheduler snasphot拿到的node对象中下面下面待调度的pod，计算出来的assigned pod个数是0, 导致Permit无法他通过。
调度周期的WaitOnPermit失败，触发unserve流程回退资源。原因有待进一步分析

正常调度流程:
1. 在filter,score后会选出来一个最优节点，如果调度成功，会继续走后面的流程，即assumed流程。 
2. 如果调度失败，filter没有选出节点，会触发postFilter，在该流程进行抢占（DefaultPreemption插件）。根据QOS及优先级将节点上pod删除、驱逐并返回被提名的节点名字nominatedNode，更新到pod.Status.NominatedNodeName，
   此时assume不会执行，而是直接返回，停止本次的调度流程，等待下一次调度执行。
3. 关于assumePod 在选出节点后，reserve前执行。 其实是在scheduleCache.AssumePod中将pod和node信息绑定，让其他pod调度感知到比如Filter。
4. reserve 估计是将绑定结果通知到plugin, pod和node绑定了，pod上其他资源也需要和node进行绑定，比如pvc要和pv进行内存的绑定。 
   看名字寓意assumePodVolumes(), assume和podVolume两个词即scheduler cache中的pv cache。 真正的pv和pvc绑定是在BindPodVolume， 也就是在preBind扩展点，等pod依赖的资源都ready后，才会真正完成pod和node的绑定。
5. 下面就是Permit扩展点了，是否允许本地的调度，否则要wait.当等待其他pod满足条件后，一起触发调度。对于wait的pod此步骤不会触发unReserve.只有在WaitOnPermit时，还有pendingplugin则会reject. 
   此时会执行wait pod的unReserve. 因为等待超时了，在等待时间里还没有allow信号. 所以wait失败。释放资源。等待下一次被放到调度队列中
6. 在绑定协程中，阻塞在WaitOnPermit，// WaitOnPermit will block, if the pod is a waiting pod, until the waiting pod is rejected or allowed.

初步规避方式是，通过client-go获取node信息，并去掉pod.spec.nodename的检测，因为此时pod是未绑定节点的。

```
# kubectl get pod --show-labels
NAME                                READY   STATUS    RESTARTS   AGE    LABELS
nginx-56vsz                         0/1     Pending   0          25m    app=nginx,pod-group.scheduling.ming.io=nginx
nginx-cbmx5                         0/1     Pending   0          15h    app=nginx,pod-group.scheduling.ming.io=nginx
nginx-hrr6q                         0/1     Pending   0          15h    app=nginx,pod-group.scheduling.ming.io=nginx

对于新建的pod
# kubectl describe pod nginx-56vsz
...
Events:
  Type     Reason            Age   From              Message
  ----     ------            ----  ----              -------
  Warning  FailedScheduling  48s   sample-scheduler  pod "nginx-56vsz" rejected while waiting on permit: rejected due to timeout after waiting 10s at plugin Coscheduling
  Warning  FailedScheduling  47s   sample-scheduler  0/3 nodes are available: 3 pod with pgName: default/nginx last failed in 3s, deny.

对于老的pod
# kubectl describe pod nginx-cbmx5
Events:
  Type     Reason            Age   From              Message
  ----     ------            ----  ----              -------
  Warning  FailedScheduling  40m   sample-scheduler  0/3 nodes are available: 3 pre-filter pod nginx-cbmx5 cannot find enough sibling pods, current pods number: 2, minMember of group: 3.
  Warning  FailedScheduling  37m   sample-scheduler  pod "nginx-cbmx5" rejected while waiting on permit: rejected due to timeout after waiting 10s at plugin Coscheduling
```
在permit扩展点，通过nodeinfo检测集群中该podgroup中的pod已经assigned给node的pod个数。 这里的client是framework的，并不是apiServer的，是snapshot的。

是的，官网注解有下面两句话，是根据snapshot来计算pods assigned to node. 而且是没有计算当前正在调度的pod, 所以是assigned+1
```
	// The number of pods that have been assigned nodes is calculated from the snapshot.
   	// The current pod in not included in the snapshot during the current scheduling cycle.
    ready := int32(assigned)+1 >= pg.Spec.MinMember
```
assign是在内置的assume及reserve扩展点完成。也就是pod assigned to node.
1. 第一个pod完成assigned后，在permit因为该条件，变为waiting pod.
2. 第二个pod完成assigned后，在permit因为该条件，变为waiting pod.
3. 第三个pod完成assigned后, 因为+1，不会变成waiting pod. 从而继续调度流程。
4. 第一、二 pod重试，在Permit的条件达到assigned的数量，由waiting pod变成可以继续调度的pod.

Premit扩展点逻辑
1. 查看pods assigned to node 的数量，
2. 对于刚开始调度的pod，会变成waitingPod，在framework侧实现， Coscheduling在Permit返回`framework.NewStatus(framework.Wait, ""), waitTime`，返回标记为framework.Wait及wait的时间。
   等待时间优先级由高到低：
   2.1 每个pg中设置的scheduleTimeoutSeconds (每个pg)
   2.2 调度器的args参数PermitWaitingTimeSeconds（全局）
   2.3 默认等待时间60s

在framework侧
1.运行所有插件的Permit,对于返回值是framework.Wait的plugin，通过一个map记录，插件名字及等待时间及标志statusCode
2.所有插件都运行完后，只要有一个插件返回的是Wait, 设置statusCode标志为，需要创建waitingPods 则遍历这个map创建waitingPod对象，并add到framework的waitingPodsMap结构体中（map+RWMutex）
```
// pkg/scheduler/framework/runtime/framework.go

    pluginsWaitTime[pl.Name()] = timeout

	if statusCode == framework.Wait {
		waitingPod := newWaitingPod(pod, pluginsWaitTime)
		f.waitingPods.add(waitingPod)
		msg := fmt.Sprintf("one or more plugins asked to wait and no plugin rejected pod %q", pod.Name)
		klog.V(4).Infof(msg)
		return framework.NewStatus(framework.Wait, msg)
	}
```
1. `newWaitingPod`返回wp对象，里面定义了当前pod和`pendingPlugins`, 令其处于wait状态的`plugins`和倒计时`waittime`后会执行的Timer。 为什么要有倒计时，因为在等待期间，可能Permit条件就达到了。
   1.1 该pod对应多个waiting的plugin, 每个plugin一个timer,当第一个wait plugin倒计时完后执行`wp.Reject`,其会将是所有的wait plugin的timer倒计时停止。
    ```
    // Reject declares the waiting pod unschedulable.
    func (w *waitingPod) Reject(msg string) {
        w.mu.RLock()
        defer w.mu.RUnlock()
        for _, timer := range w.pendingPlugins {
            timer.Stop()
        }

        // The select clause works as a non-blocking send.
        // If there is no receiver, it's a no-op (default case).
        select {
        case w.s <- framework.NewStatus(framework.Unschedulable, msg):
        // w.s没有接受者，就执行default, 非阻塞运行
        default:
        }
    }
    ```
2. `f.waitingPods.add(waitingPod)`, 将本次等待的pod加入到framework的waitingPods，这个`map`中会纳管所有wait的pod信息。
3. `return framework.NewStatus(framework.Wait, msg)`, 等待本次的调度。调度不成功时，对于返回wait的不会执行Unreserve去释放资源。延时到WaitOnPermit时（下一个调度周期触发allow or still reject）,若还是失败，去执行reserve回滚。
   其他的不成功的调度会运行Unreserve去释放/回滚资源

WaitOnPermit
1. 对于调度成功的Pod和处于Wait状态的pod可以继续后面的流程，即WaitOnPermit
```
func (f *frameworkImpl) WaitOnPermit(ctx context.Context, pod *v1.Pod) (status *framework.Status) {

    ...
    // 绑定协程会阻塞在这里
	s := <-waitingPod.s
	...
}
```
该channel有两个输入者，`waitingPod.Reject()`和`waitingPod.Allow()`。
1. Reject是waitingPod中许多的plugin中只要有一个plugin返回wait，在超过waittime后，就会往channel中发送`Unschedulable`. 试想该pod被两个plugin设置为wait, 有一个允许另一个还是wait，也就是pendingPlugins不是0,是不允许调度。
2. Allow()是waitingPod中许多的plugin中，当前plugin的Permit不返回wait而是success. 当该waitingPod的所有pendingPlugins=0时，才会向channel发送Success信号。

waitingPod.s这个channel 实现了调度周期和绑定周期的通信。绑定周期会阻塞在改channel,等待调度周期Permit的结果。

关于wp.reject。 运行到Permit扩展点时，触发wp对象创建，在waittime超时后，触发reject. 下一次调度周期仍是这个逻辑。

