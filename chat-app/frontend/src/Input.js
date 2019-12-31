import React, { Component } from 'react'
import PropTypes from 'prop-types'
import TextField from '@material-ui/core/TextField';
import SendIcon from '@material-ui/icons/Send';
import './Chat.css'
class Input extends Component {
  static propTypes = {
    onSubmitMessage: PropTypes.func.isRequired,
  }
  constructor(props){
    super(props);
    this.state = {
      message: '',
    }
}
  render() {
    return (
      <form
      //On form submit calls  SubmitMessage function in parent and send messages
          onSubmit={event => {
          event.preventDefault()
          this.props.onSubmitMessage(this.state.message)
          this.setState({ message: '' })
        }}
      >
      <TextField
        type="text"
        className="textboxmessage"
        placeholder={'Enter message...'}
        value={this.state.message}

        onChange={event => this.setState({ message: event.target.value })}
      />
        
        
        <SendIcon  onClick={event => {
          event.preventDefault()
          this.props.onSubmitMessage(this.state.message)
          this.setState({message:''})
        } }/> 
      </form>
      
    )
  }
}

export default Input
