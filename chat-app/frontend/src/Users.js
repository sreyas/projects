import React,{ Component } from 'react'
import Table from './Components/Table';
import Button from '@material-ui/core/Button';
import AddUser from './AddUser';
class Users extends Component {
  constructor(props){
    super(props);
    this.state={
      userView:[]
    }
    
    this.handleAddButton=this.handleAddButton.bind(this)
  }
  componentWillMount(){
    var userview=[];
    userview.push(<UserList handleAddButton={this.handleAddButton}/>)
        this.setState({
          userView:userview
        });
  }
  handleAddButton(event){
    console.log("hehehhe");
    var userview=[];
    userview.push(<AddUser />);
    this.setState({
      userView:userview
    });
  }

  handleSubmitButton(){
    
    var userview=[];
    userview.push(<UserList />);
    this.setState({
      userView:userview
    });
  }
  



  render(){
    return (
        this.state.userView
    );
  }
}

export function UserList(props){
  return (
  <div>
          <div>
             <Button color="primary" onClick={(event) => props.handleAddButton(event)} >Add User</Button>
          </div>

             <Table
              tableHeaderColor="primary"
              tableHead={["Name", "Country", "City", "Salary"]}
              tableData={[
                ["Dakota Rice", "Niger", "Oud-Turnhout", "$36,738"],
                ["Minerva Hooper", "Curaçao", "Sinaai-Waas", "$23,789"],
                ["Sage Rodriguez", "Netherlands", "Baileux", "$56,142"],
                ["Philip Chaney", "Korea, South", "Overland Park", "$38,735"],
                ["Doris Greene", "Malawi", "Feldkirchen in Kärnten", "$63,542"],
                ["Mason Porter", "Chile", "Gloucester", "$78,615"]
              ]}
            />
        </div>
  )
}
export default Users
