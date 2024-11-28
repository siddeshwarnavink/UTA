import { useCallback, useEffect } from 'react';
import { Container } from 'react-bootstrap';
import {
  ReactFlow,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
} from '@xyflow/react';

import '@xyflow/react/dist/base.css';

import ServerNode from '../graph/ServerNode';
import ClientNode from '../graph/ClientNode';

const NetworkGraph = ({ routingTable, transmission }) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  useEffect(() => {
    const newNodes = [...nodes], newEdges = [...edges];

    if (Object.keys(routingTable).length > 0) {
      Object.entries(routingTable)
        .filter(([_, item]) => item.role !== "wizard")
        .forEach(([_, peer]) => {
          // TCP IP of actual server/client whatever I'm processing.
          const id = peer.role === "adapter-server" ? peer.from_ip : peer.to_ip;

          if (!newNodes.find(n => n.id === id)) {
            newNodes.push({
              id,
              type: peer.role === "adapter-server" ? "server" : "client",
              position: {
                x: 0,
                y: 0
              },
              data: { ...peer }
            })
            newEdges.push({
              id: peer.from_ip + "-" + peer.to_ip,
              source: peer.from_ip,
              target: peer.to_ip
            })
          }
        });

        setNodes(newNodes);
        setEdges(newEdges);
    }
  }, [routingTable, transmission])

  const onConnect = useCallback(
    (params) => setEdges((eds) => addEdge(params, eds)),
    [setEdges],
  );

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        fitView
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={{
          server: ServerNode,
          client: ClientNode
        }}
      >
        <Background variant="dots" gap={12} size={1} />
      </ReactFlow>
    </div>
  );
}

export default NetworkGraph;
