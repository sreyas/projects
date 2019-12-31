import React from 'react'




const Message = props => {
    const { name, message} = props
    return (
        <p>
				<strong>{name}</strong> <em>{message}</em>
			</p>
    )
}
export default Message;