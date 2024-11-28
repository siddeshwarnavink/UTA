import {Handle, Position} from '@xyflow/react';

import clientIcon from '../../assets/client.png'
import Node from './Node';

const ClientNode = ({data, selected}) => {
  return (
    <Node>
      <Handle type="target" position={Position.Left} />
      <img src={clientIcon} alt='client' />
      <div>
        <strong>{data.to_ip}</strong>
      </div>
    </Node>
  );
}

export default ClientNode;
