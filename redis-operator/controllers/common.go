package controllers

import (
	"fmt"

	redisv1 "github.com/ybooks240/redisSentinel/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	redisType    = "redis"
	sentinelType = "sentinel"
)

var (
	serviceType = map[string]string{
		redisType:    "redis",
		sentinelType: "sentinel",
	}
	redisSentinelServiceType  = "type"
	RedisSentinelCommonLabel  = "app"
	redisSentinelClusterLabel = "redis.buleye.com/redis-operatiorredisSentinel"
)

func MutateHeadlessSvc(redisSentinel *redisv1.RedisSentinel, svc *corev1.Service, serviceType string) {
	var clusterLabel = redisSentinel.Namespace + "-" + redisSentinel.Name
	svc.Labels = map[string]string{
		RedisSentinelCommonLabel: "RedisSentinel",
		redisSentinelServiceType: serviceType,
	}
	var port corev1.ServicePort
	switch serviceType {
	case sentinelType:
		port = corev1.ServicePort{
			Name: serviceType,
			Port: redisSentinel.Spec.Sentinel.Port,
		}
	case redisType:
		port = corev1.ServicePort{
			Name: serviceType,
			Port: redisSentinel.Spec.Redis.Port,
		}
	}
	svc.Spec = corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Selector: map[string]string{
			redisSentinelClusterLabel: clusterLabel,
			redisSentinelServiceType:  serviceType,
		},
		Ports: []corev1.ServicePort{
			port,
		},
	}
}

func MutateStatefulSet(redisSentinel *redisv1.RedisSentinel, sts *appsv1.StatefulSet, serviceType string) {

	var clusterLabel = redisSentinel.Namespace + "-" + redisSentinel.Name
	sts.Labels = map[string]string{
		RedisSentinelCommonLabel: "RedisSentinel",
		redisSentinelServiceType: serviceType,
	}
	var replicas int32
	switch serviceType {
	case redisType:
		replicas = redisSentinel.Spec.Redis.Replicas
	case sentinelType:
		replicas = redisSentinel.Spec.Sentinel.Replicas
	}
	sts.Spec = appsv1.StatefulSetSpec{
		Replicas:    &replicas,
		ServiceName: sts.Name,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				redisSentinelClusterLabel: clusterLabel,
				redisSentinelServiceType:  serviceType,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					redisSentinelClusterLabel: clusterLabel,
					redisSentinelServiceType:  serviceType,
				},
			},
			Spec: corev1.PodSpec{
				Containers: newContainers(redisSentinel, serviceType),
			},
		},
	}

}

func newContainers(redisSentinel *redisv1.RedisSentinel, serviceType string) []corev1.Container {
	var image string
	var port int32
	var command string
	var initCommand string

	switch serviceType {
	case sentinelType:
		image = redisSentinel.Spec.Sentinel.Image
		port = redisSentinel.Spec.Sentinel.Port
		initCommand = fmt.Sprintf("cat >sentinel.conf<<EOF\nport %d\ndir /tmp\nsentinel monitor mymaster %s-0.%s %d 2\nsentinel auth-pass mymaster %s\nsentinel down-after-milliseconds mymaster 30000\nsentinel parallel-syncs mymaster 1\nsentinel failover-timeout mymaster 180000\nsentinel deny-scripts-reconfig yes\nsentinel resolve-hostnames yes\nEOF\n", port, redisSentinel.GetRedisSentinelName(redisType), redisSentinel.GetRedisSentinelName(redisType), redisSentinel.Spec.Redis.Port, redisSentinel.Spec.Redis.PassWd)

		command = initCommand + "sleep 2 ;redis-sentinel sentinel.conf"
		fmt.Println(command)

	case redisType:
		image = redisSentinel.Spec.Redis.Image
		port = redisSentinel.Spec.Redis.Port
		// command = "sleep 2;redis-server --requirepass 123.com  --masterauth 123.com"
		command = fmt.Sprintf("NodeNum=$(hostname |awk -F'-' '{print $NF}');if [ $NodeNum -eq 0 ];then redis-server --requirepass $PASSWD  --masterauth $PASSWD --port %d ;else sleep 2;redis-server --replicaof %s-0.%s %d --requirepass $PASSWD --masterauth $PASSWD --appendonly yes --port %d;fi", port, redisSentinel.GetRedisSentinelName(redisType), redisSentinel.GetRedisSentinelName(redisType), port, port)
		fmt.Println(command)
	}
	return []corev1.Container{
		corev1.Container{
			Name:            serviceType,
			Image:           image,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Ports: []corev1.ContainerPort{
				corev1.ContainerPort{
					Name:          serviceType,
					ContainerPort: port,
				},
			},
			Env: []corev1.EnvVar{
				corev1.EnvVar{
					Name:  "PASSWD",
					Value: redisSentinel.Spec.Redis.PassWd,
				},
			},
			Command: []string{
				"/bin/sh", "-ec",
				command,
			},
			// Lifecycle: &corev1.Lifecycle{
			// 	PreStop: &corev1.Handler{
			// 		Exec: &corev1.ExecAction{
			// 			Command: []string{},
			// 		},
			// 	},
			// },
		},
	}
}
