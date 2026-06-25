package one

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MasterRoleLabel       string = "node-role.kubernetes.io/master"
	ControlPlaneRoleLabel string = "node-role.kubernetes.io/control-plane"
)

// GetMasters returns all master nodes (either master or control-plane role label).
func (one *One) GetMasters() ([]corev1.Node, error) {
	ctx := context.Background()
	seen := make(map[string]struct{})
	var masters []corev1.Node

	// master 节点可能带有 node-role.kubernetes.io/master 或 node-role.kubernetes.io/control-plane
	// 标签（新旧版本 Kubernetes 命名不同），满足任一即可视为 master。
	//
	// Kubernetes API 的 label selector 仅支持 AND，无法在单次 List 中对两个不同 key 做 OR，
	// 因此分别按两个 label 查询，再按节点名去重合并结果。
	for _, label := range []string{MasterRoleLabel, ControlPlaneRoleLabel} {
		nodes, err := one.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: label})
		if err != nil {
			return nil, err
		}
		for _, node := range nodes.Items {
			if _, ok := seen[node.Name]; ok {
				continue
			}
			seen[node.Name] = struct{}{}
			masters = append(masters, node)
		}
	}
	return masters, nil
}

// GetWorkers return all worker nodes (not including master)
func (one *One) GetWorkers() ([]corev1.Node, error) {
	ctx := context.Background()
	nodes, err := one.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: "!" + MasterRoleLabel})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

// GetNodes return all nodes (masters + workers)
func (one *One) GetNodes() ([]corev1.Node, error) {
	ctx := context.Background()
	nodes, err := one.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}
