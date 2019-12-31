import React, { Component } from 'react'
import Input from './Input'
import Message from './Message'
import TextField from '@material-ui/core/TextField';
import InputLabel from '@material-ui/core/InputLabel';
import './Chat.css'

const URL = 'ws://localhost:4000'

class Chat extends Component {
  constructor(props){
  super(props);
  this.state = {
    name: '',
    messages: [],
    
  }
}
  ws = new WebSocket(URL)

  componentDidMount() {
    //Receieves Message and updates message list view
    this.ws.onmessage = evt => {
      const message = JSON.parse(evt.data)
      this.sendMessage(message)
    }
    //When connection refused reconnects to socket
    this.ws.onclose = () => {
      this.setState({
        ws: new WebSocket(URL),
      })
    }
  }
//Update messages on the list view
  sendMessage = message =>
    this.setState(state => ({ messages: [message, ...state.messages]}))

//Send message and updates in the list view
  submitMessage = messageString => {
    const message = { name: this.state.name, message: messageString }
    this.ws.send(JSON.stringify(message))
    this.sendMessage(message)
  }

  render() {
    return (
      <div>
        <InputLabel htmlFor="name">
          Name:&nbsp;
          <TextField
          className="textBoxName"
            type="text"
            id={'name'}
            placeholder={'Enter your name...'}
            onChange={e => this.setState({ name: e.target.value })}
          />
        </InputLabel>
        <div className="message-view">
        {this.state.messages.map((message, index) =>
          <Message
            key={index}
            message={message.message}
            name={message.name}
          />,
        )}
        </div>
        <Input
          ws={this.ws}
          onSubmitMessage={messageString => this.submitMessage(messageString)}
        />
      </div>
    )
  }
}

export default Chat
