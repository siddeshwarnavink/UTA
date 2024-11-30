import { getSmoothStepPath } from '@xyflow/react';

const AnimatedEdge = ({
  id,
  sourceX,
  sourceY,
  sourcePosition,
  targetX,
  targetY,
  targetPosition,
  style,
  markerEnd,
}) => {
  const [edgePath] = getSmoothStepPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  return (
    <g>
      <path
        id={id}
        d={edgePath}
        style={{
          ...style,
          stroke: 'black',
          strokeWidth: 2,
          strokeDasharray: '8, 6',
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          animation: 'dashmove 1s linear infinite', // in App.css
          fill: 'none',
        }}
        markerEnd={markerEnd}
      />
    </g>
  );
};
export default AnimatedEdge;
