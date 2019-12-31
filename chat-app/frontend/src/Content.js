import React, { Component } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Chat from './Chat';
import  './Content.css';
import Logo from './logo.png'
const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  
}));



class Content extends Component {
  constructor(props) {
    super(props)
  
    this.state = {
       defaultView:[],
       mainContent:[]
    }
  }
  componentWillMount(){
    const classes = this.props.myHookValue;
    var defaultView=[];
    defaultView.push(
      <div className={classes.root}>
      <CssBaseline />
      <AppBar position="fixed" className={classes.appBar}>
        
        <img src={Logo} className="App-logo" alt="logo" />
      </AppBar>
     
      
    </div>

    );
    this.setState({
      defaultView:defaultView,
      mainContent:<Chat />
    });
  }
    
  render() {
    return (
      <div>
       {this.state.defaultView}
       <div className="Main-div">
       {this.state.mainContent}
       </div>
      </div>
    )
  }


}
function withMyHook(Component) {
  return function WrappedComponent(props) {
    const myHookValue = useStyles();
    return <Component {...props} myHookValue={myHookValue} />;
  }
}
  Content = withMyHook(Content);


  
export default Content;