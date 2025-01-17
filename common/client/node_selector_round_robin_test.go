package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goplugin/pluginv3.0/v2/common/types"
)

func TestRoundRobinNodeSelectorName(t *testing.T) {
	selector := newNodeSelector[types.ID, RPCClient[types.ID, Head]](NodeSelectionModeRoundRobin, nil)
	assert.Equal(t, selector.Name(), NodeSelectionModeRoundRobin)
}

func TestRoundRobinNodeSelector(t *testing.T) {
	t.Parallel()

	type nodeClient RPCClient[types.ID, Head]
	var nodes []Node[types.ID, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(nodeStateOutOfSync)
		} else {
			// second & third nodes are alive
			node.On("State").Return(nodeStateAlive)
		}
		nodes = append(nodes, node)
	}

	selector := newNodeSelector(NodeSelectionModeRoundRobin, nodes)
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
	assert.Same(t, nodes[1], selector.Select())
	assert.Same(t, nodes[2], selector.Select())
}

func TestRoundRobinNodeSelector_None(t *testing.T) {
	t.Parallel()

	type nodeClient RPCClient[types.ID, Head]
	var nodes []Node[types.ID, nodeClient]

	for i := 0; i < 3; i++ {
		node := newMockNode[types.ID, nodeClient](t)
		if i == 0 {
			// first node is out of sync
			node.On("State").Return(nodeStateOutOfSync)
		} else {
			// others are unreachable
			node.On("State").Return(nodeStateUnreachable)
		}
		nodes = append(nodes, node)
	}

	selector := newNodeSelector(NodeSelectionModeRoundRobin, nodes)
	assert.Nil(t, selector.Select())
}
