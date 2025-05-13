"use client";
import React, { useMemo } from "react";
import ReactFlow, {
    Background,
    Controls,
    MiniMap,
    Node,
    Edge,
    useNodesState,
    useEdgesState,
} from "react-flow-renderer";

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
    const { initialNodes, initialEdges } = useMemo(() => {
        const nodeMap = new Map<number, NodeData>();
        const childrenMap = new Map<number, NodeData[]>();

        // Flatten all nodes into a map
        for (const triplet of data.content) {
            for (const node of triplet) {
                nodeMap.set(node.id, node);
                if (!childrenMap.has(node.parent)) {
                    childrenMap.set(node.parent, []);
                }
                childrenMap.get(node.parent)!.push(node);
            }
        }

        const nodes: Node[] = [];
        const edges: Edge[] = [];

        const layoutTree = (
            id: number,
            depth: number,
            xPos: number
        ): { width: number; centerX: number } => {
            const children = childrenMap.get(id) || [];

            if (children.length === 0) {
                const node = nodeMap.get(id);
                const isMerger = node?.name.toLowerCase() === "merger";

                const position = { x: xPos, y: depth * 150 };
                nodes.push({
                    id: id.toString(),
                    data: { label: isMerger ? "" : node?.name },
                    position,
                    style: isMerger
                        ? {
                            width: 10,
                            height: 10,
                            backgroundColor: "#000",
                            borderRadius: "50%",
                        }
                        : {
                            padding: 10,
                            border: "1px solid #555",
                        },
                });
                return { width: 160, centerX: xPos + 80 };
            }

            let totalWidth = 0;
            const childCenters: number[] = [];
            let currentX = xPos;

            for (const child of children) {
                const childLayout = layoutTree(child.id, depth + 1, currentX);
                totalWidth += childLayout.width;
                childCenters.push(childLayout.centerX);
                currentX += childLayout.width + 20;
            }

            const centerX =
                childCenters.length > 0
                    ? childCenters.reduce((a, b) => a + b, 0) / childCenters.length
                    : xPos;

            const node = nodeMap.get(id) || { id: 0, name: data.name, parent: -1 };
            const isMerger = node.name.toLowerCase() === "merger";

            nodes.push({
                id: id.toString(),
                data: { label: isMerger ? "" : node.name },
                position: { x: centerX, y: depth * 150 },
                style: isMerger
                    ? {
                        width: 10,
                        height: 10,
                        backgroundColor: "#000",
                        borderRadius: "50%",
                    }
                    : {
                        padding: 10,
                        border: "1px solid #555",
                    },
            });

            for (const child of children) {
                edges.push({
                    id: `e${id}-${child.id}`,
                    source: id.toString(),
                    target: child.id.toString(),
                    animated: true,
                });
            }

            return { width: totalWidth, centerX };
        };

        layoutTree(0, 0, 0);

        return { initialNodes: nodes, initialEdges: edges };
    }, [data]);

    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

    return (
        <div style={{ width: "100%", height: "100%" }}>
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
                <MiniMap />
                <Controls />
                <Background />
            </ReactFlow>
        </div>
    );
};

export default FlowDiagram;
