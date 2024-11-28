import { Handle, Position } from '@xyflow/react';

import serverIcon from '../../assets/server.png'
import Node from './Node';

const ServerNode = ({ data, selected }) => {
  return (
    <Node>
      <img src={serverIcon} alt='Server' />
      <div>
        <strong>{data.from_ip}</strong>
      </div>
      <Handle type='source' position={Position.Right} />
    </Node>
  );
}

export default ServerNode;
