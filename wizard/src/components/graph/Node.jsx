import classes from './Node.module.css';

const Node = ({children}) => {
  return (
    <div className={classes.node}>
      {children}
    </div>
  )
}

export default Node;
