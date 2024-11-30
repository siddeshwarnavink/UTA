import { useCallback, useEffect, useMemo } from 'react';
import {
  ReactFlow,
  Background,
  useNodesState,
  useEdgesState,
  addEdge
} from '@xyflow/react';
import Dagre from '@dagrejs/dagre';
import '@xyflow/react/dist/base.css';

import ServerNode from '../graph/ServerNode';
import ClientNode from '../graph/ClientNode';
import AnimatedEdge from '../graph/AnimatedEdge';

const NetworkGraph = ({ routingTable, transmission }) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  const layoutHandler = (nodes, edges) => {
    const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
    g.setGraph({ rankdir: "LR", ranksep: 330 });

    edges.forEach((edge) => g.setEdge(edge.source, edge.target));
    nodes.forEach((node) =>
      g.setNode(node.id, {
        ...node,
        width: node.measured?.width ?? 120,
        height: node.measured?.height ?? 120,
      }),
    );

    Dagre.layout(g);

    return nodes.map((node) => {
      const position = g.node(node.id);

      const x = position.x - (node.measured?.width ?? 0) / 2;
      const y = position.y - (node.measured?.height ?? 0) / 2;

      return { ...node, position: { x, y } };
    });
  };

  const tableLength = useMemo(() => Object.keys(routingTable).length, [routingTable]);

  useEffect(() => {
    if (nodes.length !== tableLength && tableLength > 0) {
      let newNodes = [...nodes], newEdges = [...edges], update = false;

      Object.entries(routingTable)
        .filter(([_, item]) => item.role !== "wizard")
        .forEach(([_, peer]) => {
          // TCP IP of actual server/client whatever I'm processing.
          const id = peer.role === "adapter-server" ? peer.from_ip : peer.to_ip;

          if (!newNodes.find(n => n.id === id)) {
            console.log("pushing new node");
            update = true;
            newNodes.push({
              id,
              type: peer.role === "adapter-server" ? "server" : "client",
              position: {
                x: 0,
                y: 0
              },
              data: { ...peer }
            })
          }

          const edgeId = peer.from_ip + "-" + peer.to_ip;
          if (!newEdges.find(e => e.id === edgeId)) {
            console.log("pushing new edge");
            update = true;
            newEdges.push({
              id: edgeId,
              type: transmission && transmission.find(t => edgeId.indexOf(t.ip) !== 0) ? "dataFlow" : "smoothstep",
              source: peer.from_ip,
              target: peer.to_ip,
            })
          }
        });

      if (update) {
        newNodes = layoutHandler(newNodes, newEdges);
        setNodes(newNodes);
        setEdges(newEdges);
      }
    }
  }, [tableLength]);

  useEffect(() => {
    if (transmission) {
      Object.keys(transmission).forEach(ip => {
        setEdges(prevEdges => {
          return prevEdges.map(edge => {
            if (transmission && edge.id.indexOf(ip) !== 0 && transmission[ip]) {
              console.log("edge data flow");
              return { ...edge, type: "dataFlow" }
            }
            return edge;
          });
        });
      });
    }
  }, [transmission]);

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
        edgeTypes={{
          dataFlow: AnimatedEdge
        }}
      >
        <Background variant="dots" gap={12} size={1} />
      </ReactFlow>
    </div>
  );
}

export default NetworkGraph;
