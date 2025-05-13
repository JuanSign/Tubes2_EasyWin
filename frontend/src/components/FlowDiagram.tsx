'use client';
import React, { useMemo, useState, useEffect } from 'react';
import ReactFlow, {
    Background,
    Controls,
    MiniMap,
    useNodesState,
    useEdgesState,
    Node,
    Edge,
} from 'react-flow-renderer';

interface NodeData {
    name: string;
    id: number;
    parent: number;
}

interface DiagramData {
    name: string;
    content: NodeData[][];
}

interface Props {
    data: DiagramData;
}

const FlowDiagram: React.FC<Props> = ({ data }) => {
    const [step, setStep] = useState(0);
    const maxStep = data.content.length;

    const { nodes: initialNodes, edges: initialEdges } = useMemo(() => {
        const nodeMap = new Map<number, NodeData>();
        const childrenMap = new Map<number, NodeData[]>();
        const slicedContent = data.content.slice(0, step + 1).flat();

        const root: NodeData = { name: data.name, id: 0, parent: -1 };
        slicedContent.push(root);

        for (const node of slicedContent) {
            nodeMap.set(node.id, node);
            if (!childrenMap.has(node.parent)) {
                childrenMap.set(node.parent, []);
            }
            childrenMap.get(node.parent)!.push(node);
        }

        const nodes: Node[] = [];
        const edges: Edge[] = [];

        const layoutTree = (
            id: number,
            depth: number,
            xPos: number
        ): { width: number; centerX: number } => {
            const children = childrenMap.get(id) || [];
            const minNodeWidth = 160;
            const spacing = 60;

            if (children.length === 0) {
                const node = nodeMap.get(id);
                const isMerger = node?.name.toLowerCase() === 'merger';
                const position = { x: xPos, y: depth * 150 };

                nodes.push({
                    id: id.toString(),
                    data: { label: isMerger ? '' : node?.name },
                    position,
                    style: isMerger
                        ? {
                            width: 12,
                            height: 12,
                            backgroundColor: '#03A9F4',
                            borderRadius: '50%',
                        }
                        : {
                            padding: 10,
                            border: '1px solid #03A9F4',
                            background: '#1e1e1e',
                            color: '#fff',
                            borderRadius: 6,
                            boxShadow: '0 2px 10px rgba(0,0,0,0.3)',
                        },
                });

                return { width: minNodeWidth, centerX: xPos + minNodeWidth / 2 };
            }

            let totalWidth = 0;
            const childCenters: number[] = [];
            let currentX = xPos;

            for (let i = 0; i < children.length; i++) {
                const child = children[i];
                const childLayout = layoutTree(child.id, depth + 1, currentX);
                totalWidth += childLayout.width;
                childCenters.push(childLayout.centerX);
                currentX += childLayout.width + spacing;
            }

            const centerX =
                childCenters.reduce((a, b) => a + b, 0) / childCenters.length;

            const node = nodeMap.get(id) || root;
            const isMerger = node.name.toLowerCase() === 'merger';

            nodes.push({
                id: id.toString(),
                data: { label: isMerger ? '' : node.name },
                position: { x: centerX - minNodeWidth / 2, y: depth * 150 },
                style: isMerger
                    ? {
                        width: 12,
                        height: 12,
                        backgroundColor: '#03A9F4',
                        borderRadius: '50%',
                    }
                    : {
                        padding: 10,
                        border: '1px solid #03A9F4',
                        background: '#1e1e1e',
                        color: '#fff',
                        borderRadius: 6,
                        boxShadow: '0 2px 10px rgba(0,0,0,0.3)',
                    },
            });

            for (const child of children) {
                edges.push({
                    id: `e${id}-${child.id}`,
                    source: id.toString(),
                    target: child.id.toString(),
                    animated: false,
                    style: { stroke: '#03A9F4', strokeWidth: 2 },
                });
            }

            return {
                width: totalWidth + spacing * (children.length - 1),
                centerX,
            };
        };

        layoutTree(0, 0, 0);

        return { nodes, edges };
    }, [data, step]);

    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

    useEffect(() => {
        setNodes(initialNodes);
        setEdges(initialEdges);
    }, [initialNodes, initialEdges, setNodes, setEdges]);

    return (
        <div
            style={{
                width: '100%',
                height: '100%',
                position: 'relative',
                backgroundColor: '#121212',
            }}
        >
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                fitView
                panOnDrag
                zoomOnScroll
                zoomOnPinch
            >
                <MiniMap
                    nodeColor={() => '#03A9F4'}
                    style={{ backgroundColor: '#1a1a1a' }}
                />
                <Controls />
                <Background color="#333" gap={16} />
            </ReactFlow>

            {/* Step Controls */}
            <div
                style={{
                    position: 'absolute',
                    bottom: 20,
                    left: 20,
                    display: 'flex',
                    gap: 10,
                    zIndex: 10,
                    pointerEvents: 'auto',
                    backgroundColor: '#1e1e1e',
                    padding: '8px 16px',
                    borderRadius: '8px',
                    border: '1px solid #444',
                    boxShadow: '0 2px 6px rgba(0,0,0,0.4)',
                }}
            >
                <button
                    onClick={() => setStep((s) => Math.max(0, s - 1))}
                    disabled={step === 0}
                    style={{
                        padding: '6px 12px',
                        borderRadius: 4,
                        border: 'none',
                        backgroundColor: '#2c2c2c',
                        color: '#fff',
                        cursor: 'pointer',
                        opacity: step === 0 ? 0.5 : 1,
                    }}
                >
                    ◀ Prev
                </button>
                <span style={{ color: '#fff', fontSize: 14 }}>Step {step + 1}</span>
                <button
                    onClick={() => setStep((s) => Math.min(maxStep - 1, s + 1))}
                    disabled={step >= maxStep - 1}
                    style={{
                        padding: '6px 12px',
                        borderRadius: 4,
                        border: 'none',
                        backgroundColor: '#2c2c2c',
                        color: '#fff',
                        cursor: 'pointer',
                        opacity: step >= maxStep - 1 ? 0.5 : 1,
                    }}
                >
                    Next ▶
                </button>
            </div>
        </div>
    );
};

export default FlowDiagram;
